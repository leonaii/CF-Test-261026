package main

// 离线上号：CG 起一个本地 HTTP 服务器，冒充卡密服务器，
// 按卡密对应的 QQ 从卡密平台接口取 ini 明文，回给 CF上号器_v1.17.exe，由它自己启动游戏。
//
// 协议（实测 2026-10-23，代理抓包）：
//
//	POST /api/card/redeem  {"card_key":"CF..."}
//	→ 200 {
//	     "ok": true, "qq": "1978628973", "zone": "江苏电信一区",
//	     "note": "...", "extra": "...", "expires_at": 1792795881839,
//	     "key": "<base64 RSA-OAEP-SHA256(AESKey)>",
//	     "iv":  "<base64 12B>", "tag": "<base64 16B>",
//	     "data":"<base64 AES-256-GCM(明文)>"
//	  }
//
// 明文（2260 B）= "[QQ]\r\nToken0106=…\r\n…sig=…\r\n"，与 ini 文件内容同构。
//
// 上号器读 host.dat 拿 base_url，所以还要把 host.dat 临时指到本地。

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	ymmLocalPrefix = "http://127.0.0.1:18790"
	// 离线上号不看卡密真实状态：有效期统一按 10 年给，过期卡也能上号验号
	ymmValidSpan  = 10 * 365 * 24 * time.Hour
	ymmValidHours = 10 * 365 * 24
)

// ymmZoneFromNote 从备注里取区服（格式：编号----QQ----卡密----等级----区服----战区----...）。
func ymmZoneFromNote(note string) string {
	parts := strings.Split(note, "----")
	if len(parts) > 4 {
		return strings.TrimSpace(parts[4])
	}
	return ""
}

// ymmSealEnvelope 把票据明文封成上号器能解的信封。
func ymmSealEnvelope(plain []byte) (map[string]any, error) {
	priv := privateKey()
	if priv == nil {
		return nil, fmt.Errorf("RSA 私钥未就绪")
	}

	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, err
	}

	keyCt, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &priv.PublicKey, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA 加密: %w", err)
	}

	blk, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(blk)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, g.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	sealed := g.Seal(nil, iv, plain, nil)
	data, tag := sealed[:len(sealed)-g.Overhead()], sealed[len(sealed)-g.Overhead():]

	b64 := func(b []byte) string { return base64.StdEncoding.EncodeToString(b) }
	return map[string]any{
		"key":  b64(keyCt),
		"iv":   b64(iv),
		"tag":  b64(tag),
		"data": b64(data),
	}, nil
}

// ymmOfflineData 保存 card_key → (qq, note) 映射；ini 按需从接口拉。
type ymmOfflineData struct {
	api    *cardClient
	pin    string            // 选中的卡密：config.ini 里是占位卡密，所有请求都按它取数据
	cards  map[string]ymmRow // card_key -> row
	byQQ   map[string]string // qq -> card_key
	mu     sync.Mutex
	ticket map[string]string // qq -> ini 内容缓存
}

// resolve 把请求里的卡密（可能是 config.ini 的占位卡密）统一换成选中的卡密。
func (d *ymmOfflineData) resolve(code string) string {
	if d.pin != "" {
		return d.pin
	}
	return code
}

// ymmBuildOffline 用当前列表快照构建离线上号所需的数据。
// pin 是界面上选中的卡密：上号器读 config.ini 的占位卡密发请求，统一用 pin 取数据。
func ymmBuildOffline(api *cardClient, rows []ymmRow, pin string) (*ymmOfflineData, error) {
	if api == nil || !api.loggedIn() {
		return nil, fmt.Errorf("未登录卡密服务器")
	}
	d := &ymmOfflineData{
		api:    api,
		pin:    strings.TrimSpace(pin),
		cards:  map[string]ymmRow{},
		byQQ:   map[string]string{},
		ticket: map[string]string{},
	}
	for _, r := range rows {
		if r.CardKey == "" || r.QQ == "" {
			continue
		}
		d.cards[r.CardKey] = r
		d.byQQ[r.QQ] = r.CardKey
	}
	if len(d.cards) == 0 {
		return nil, fmt.Errorf("当前列表里没有可用卡密")
	}
	return d, nil
}

// ticketFor 按卡密（或 QQ）从接口取 ini 明文，取到后缓存。
func (d *ymmOfflineData) ticketFor(code string) (string, string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", "", fmt.Errorf("卡密为空")
	}

	qq := code
	if row, ok := d.cards[code]; ok {
		qq = row.QQ
	} else if k, ok := d.byQQ[code]; ok {
		qq = k
	}
	if qq == "" {
		return "", "", fmt.Errorf("未知卡密: %s", code)
	}

	d.mu.Lock()
	if t, ok := d.ticket[qq]; ok {
		d.mu.Unlock()
		return t, qq, nil
	}
	d.mu.Unlock()

	text, err := d.api.getData(qq)
	if err != nil {
		return "", "", fmt.Errorf("接口取 %s 的 ini: %w", qq, err)
	}
	d.mu.Lock()
	d.ticket[qq] = text
	d.mu.Unlock()
	return text, qq, nil
}

// ymmStartOfflineServer 启动本地服务器，返回关闭函数。
func ymmStartOfflineServer(d *ymmOfflineData) (func(), string, error) {
	if d == nil {
		return nil, "", fmt.Errorf("离线上号数据未就绪")
	}

	mux := http.NewServeMux()

	// boot：直接返回固定内容，不查数据（与回放模式一致）
	mux.HandleFunc("/api/card/boot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"banned":false}`))
	})

	handleRedeem := func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		code := ""
		for _, k := range []string{"card_key", "code", "card"} {
			if v, ok := req[k].(string); ok && v != "" {
				code = v
				break
			}
		}
		if code == "" {
			http.Error(w, `{"ok":false,"message":"empty code"}`, 200)
			return
		}
		// config.ini 里写的是占位卡密，这里统一换成界面上选中的卡密
		code = d.resolve(code)

		plain, qq, err := d.ticketFor(code)
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "message": err.Error()})
			return
		}

		env, err := ymmSealEnvelope([]byte(plain))
		if err != nil {
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "message": err.Error()})
			return
		}

		row := d.cards[code]
		out := map[string]any{
			"ok":          true,
			"qq":          qq,
			"zone":        ymmZoneFromNote(row.Note),
			"note":        row.Note,
			"extra":       row.Note,
			"activated":   true,
			"valid_hours": ymmValidHours,
			// 过期卡也照给：有效期一律 10 年后
			"expires_at": time.Now().Add(ymmValidSpan).UnixMilli(),
		}
		for k, v := range env {
			out[k] = v
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	}

	// redeem 返信封；peek 返明文（实测 peek 是 566B 明文，redeem 是 3979B 信封）
	mux.HandleFunc("/api/card/redeem", handleRedeem)
	mux.HandleFunc("/api/card/peek", func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		code := ""
		for _, k := range []string{"card_key", "code", "card"} {
			if v, ok := req[k].(string); ok && v != "" {
				code = v
				break
			}
		}
		if code == "" {
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": false, "message": "empty code"})
			return
		}
		code = d.resolve(code)

		row, ok := d.cards[code]
		qq := row.QQ
		if !ok {
			qq = code // 允许直接传 QQ
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":          true,
			"qq":          qq,
			"zone":        ymmZoneFromNote(row.Note),
			"note":        row.Note,
			"extra":       row.Note,
			"activated":   true, // 一律按已激活处理，避免上号器再走激活校验
			"expires_at":  time.Now().Add(ymmValidSpan).UnixMilli(),
			"valid_hours": ymmValidHours,
		})
	})
	mux.HandleFunc("/api/card/query", handleRedeem)
	// 日志接口：直接 200
	mux.HandleFunc("/api/card/log", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	// 兜底
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	ln, err := net.Listen("tcp", "127.0.0.1:18790")
	if err != nil {
		return nil, "", fmt.Errorf("本地端口 18790 被占用: %w", err)
	}
	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()

	return func() { _ = srv.Close() }, ymmLocalPrefix, nil
}

// ymmSwitchHostDat 把上号器目录下的 host.dat 指向本地，返回恢复函数。
func ymmSwitchHostDat(launcherDir string) (func(), error) {
	p := filepath.Join(launcherDir, "host.dat")
	orig, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("读 host.dat: %w", err)
	}

	blob, err := EncryptHostDat(ymmLocalPrefix)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(p, blob, 0o644); err != nil {
		return nil, err
	}

	return func() { _ = os.WriteFile(p, orig, 0o644) }, nil
}

// ymmLaunchViaLauncher 离线上号主流程：
// 起本地服务器 → 改 host.dat → 启动上号器 → 等它退出 → 恢复 host.dat。
func ymmLaunchViaLauncher(launcherPath string, d *ymmOfflineData) error {
	stopSrv, addr, err := ymmStartOfflineServer(d)
	if err != nil {
		return err
	}
	defer stopSrv()

	dir := filepath.Dir(launcherPath)
	restore, err := ymmSwitchHostDat(dir)
	if err != nil {
		return err
	}
	defer restore()

	fmt.Printf("本地服务器 %s，host.dat 已切换\n", addr)
	cmd := exec.Command(launcherPath)
	cmd.Dir = dir
	if err := cmd.Start(); err != nil {
		return err
	}
	// 上号器运行期间 host.dat 必须保持指向本地，退出后才还原
	_ = cmd.Wait()
	return nil
}
