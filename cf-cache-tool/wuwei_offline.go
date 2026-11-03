package main

// 无畏上号器（无畏契约 / VALORANT）离线桥接。
//
// 协议与鸡哥同源（实测 2026-11-03，静态还原）：
//
//	POST {base}/api/card/boot    → {"ok":true,"banned":false}
//	POST {base}/api/card/redeem  {"card_key":"16 位纯英文"}
//	     → {"ok":true,"qq":"...","expires_at":毫秒,"zone":"...",
//	        "key":RSA-OAEP-SHA256(AESKey),"iv":12B,"tag":16B,"data":AES-256-GCM(ini 明文)}
//	POST {base}/api/card/log     → {"ok":true}
//
// 明文就是 ini（[QQ] + Token010A/0143 + machine_id + ct/st…），与 dump 里的票据同构；
// 程序只校验明文里含 "Token010A=" 或 "openid="。
//
// RSA 公钥模数与鸡哥直登完全相同，所以直接复用 jigeSeal 信封。
//
// 唯一区别：服务器地址硬编码在 exe 的 .rdata 里，24 字节 XOR 混淆 + FNV-1a64 自校验
// （http://111.170.7.81:8750）。启动前把这段改成本地地址并重算 FNV，退出不用还原，
// 因为改的是工作目录里的副本，原始 exe 不动。

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	wuweiExeName    = "无畏上号器.exe"
	wuweiPatchedExe = "无畏上号器.local.exe"
	wuweiListenAddr = "127.0.0.1:18793"
	// 必须正好 24 字节：exe 里这个字符串是定长槽位
	wuweiLocalURL = "http://127.0.0.1:18793/x"
)

var (
	// 原始地址的 XOR key 与明文，用于定位补丁点并做版本校验。
	wuweiXorKey  = []byte{0x3c, 0x91, 0x57, 0xe2, 0x0b, 0x6d, 0xa4}
	wuweiStubURL = "http://111.170.7.81:8750"
)

func wuweiFNV64(s string) uint64 {
	const (
		offset64 = 0xcbf29ce484222325
		prime64  = 0x100000001b3
	)
	h := uint64(offset64)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}

// wuweiEncodeURL 按 exe 里的方式做 XOR 混淆。
func wuweiEncodeURL(s string) []byte {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		out[i] = s[i] ^ wuweiXorKey[i%len(wuweiXorKey)]
	}
	return out
}

// wuweiPatchExe 返回改好地址的 exe 内容（不修改入参）。
func wuweiPatchExe(raw []byte) ([]byte, error) {
	if len(wuweiLocalURL) != len(wuweiStubURL) {
		return nil, fmt.Errorf("本地地址必须和原地址等长（%d 字节）", len(wuweiStubURL))
	}

	slot := bytes.Index(raw, wuweiEncodeURL(wuweiStubURL))
	if slot < 0 {
		return nil, fmt.Errorf("不认识的无畏上号器版本：找不到服务器地址")
	}

	oldSum := make([]byte, 8)
	binary.LittleEndian.PutUint64(oldSum, wuweiFNV64(wuweiStubURL))
	// movabs rcx, imm64 → 48 B9 <8 字节小端>
	check := bytes.Index(raw, append([]byte{0x48, 0xb9}, oldSum...))
	if check < 0 {
		return nil, fmt.Errorf("不认识的无畏上号器版本：找不到地址校验值")
	}

	out := append([]byte(nil), raw...)
	copy(out[slot:slot+len(wuweiStubURL)], wuweiEncodeURL(wuweiLocalURL))
	binary.LittleEndian.PutUint64(out[check+2:check+10], wuweiFNV64(wuweiLocalURL))
	return out, nil
}

// wuweiPrepareExe 读原始 exe，生成/复用打过补丁的副本。
func wuweiPrepareExe(src, dir string) (string, error) {
	raw, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("读上号器失败: %w", err)
	}
	patched, err := wuweiPatchExe(raw)
	if err != nil {
		return "", err
	}
	dst := filepath.Join(dir, wuweiPatchedExe)
	if cur, err := os.ReadFile(dst); err == nil && bytes.Equal(cur, patched) {
		return dst, nil
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(dst, patched, 0o700); err != nil {
		return "", fmt.Errorf("写入上号器副本失败: %w", err)
	}
	return dst, nil
}

// wuweiStartServer 起本地假卡密服务。所有路径用后缀匹配，避免路径前缀被客户端改写。
func wuweiStartServer(ini, qq, zone string) (func(), error) {
	env, err := jigeSeal([]byte(ini))
	if err != nil {
		return nil, fmt.Errorf("封信封失败: %w", err)
	}
	if strings.TrimSpace(zone) == "" {
		zone = "无畏契约"
	}

	reply := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	redeem := func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		out := map[string]any{
			"ok":          true,
			"qq":          qq,
			"zone":        zone,
			"note":        zone,
			"extra":       zone,
			"activated":   true,
			"valid_hours": ymmValidHours,
			// 离线不看卡密状态：有效期一律给 10 年
			"expires_at": time.Now().Add(ymmValidSpan).UnixMilli(),
		}
		for k, v := range env {
			out[k] = v
		}
		reply(w, out)
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		p := strings.ToLower(r.URL.Path)
		switch {
		case strings.HasSuffix(p, "/api/card/boot"):
			reply(w, map[string]any{"ok": true, "banned": false})
		case strings.HasSuffix(p, "/api/card/redeem"),
			strings.HasSuffix(p, "/api/card/query"),
			strings.HasSuffix(p, "/api/card/peek"):
			redeem(w, r)
		case strings.HasSuffix(p, "/api/card/log"):
			reply(w, map[string]any{"ok": true})
		default:
			reply(w, map[string]any{"ok": true})
		}
	}

	ln, err := net.Listen("tcp", wuweiListenAddr)
	if err != nil {
		return nil, fmt.Errorf("本地端口 %s 被占用: %w", wuweiListenAddr, err)
	}
	srv := &http.Server{Handler: http.HandlerFunc(h)}
	go func() { _ = srv.Serve(ln) }()
	return func() { _ = srv.Close() }, nil
}

// wuweiWriteConfig 写 [ui] 段：占位卡密 + 已选的 ACLOS 代理路径，其余键原样保留。
func wuweiWriteConfig(dir string) (card, proxy string, err error) {
	path := filepath.Join(dir, "config.ini")
	vals := map[string]string{}
	if raw, rerr := os.ReadFile(path); rerr == nil {
		section := ""
		for _, line := range strings.Split(string(raw), "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, ";") || strings.HasPrefix(t, "#") {
				continue
			}
			if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
				section = strings.ToLower(strings.Trim(t, "[]"))
				continue
			}
			if section != "ui" {
				continue
			}
			if k, v, ok := strings.Cut(t, "="); ok {
				vals[strings.ToLower(strings.TrimSpace(k))] = strings.TrimSpace(v)
			}
		}
	}

	proxy = vals["exe"]
	if proxy == "" {
		proxy = wuweiDetectProxy()
	}
	card = jigeRandomCard()
	vals["card"] = card
	vals["exe"] = proxy
	vals["remember"] = "1"

	var b strings.Builder
	b.WriteString("[ui]\r\n")
	for _, k := range []string{"card", "exe", "remember"} {
		b.WriteString(k + "=" + vals[k] + "\r\n")
		delete(vals, k)
	}
	for k, v := range vals {
		b.WriteString(k + "=" + v + "\r\n")
	}
	if err = os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		return card, proxy, err
	}
	return card, proxy, nil
}

// wuweiDetectProxy 在常见 WeGame 目录里找 AclosGameProxy.exe（找不到就留空，让上号器自己弹框选）。
func wuweiDetectProxy() string {
	var roots []string
	for _, env := range []string{"ProgramFiles(x86)", "ProgramFiles"} {
		if v := strings.TrimSpace(os.Getenv(env)); v != "" {
			roots = append(roots, filepath.Join(v, "WeGame"))
		}
	}
	for _, d := range []string{"C:", "D:", "E:", "F:", "G:", "H:"} {
		roots = append(roots,
			d+`\WeGameApps`,
			d+`\网络游戏`,
			d+`\游戏`,
			d+`\Games`,
			d+`\WeGame`,
		)
	}
	rel := []string{
		`rail_apps\无畏契约(2001715)\ACLOS\Proxy\AclosGameProxy.exe`,
		`rail_apps\无畏契约（异变轮战上线） (2001715)\ACLOS\Proxy\AclosGameProxy.exe`,
		`无畏契约(2001715)\ACLOS\Proxy\AclosGameProxy.exe`,
	}
	for _, root := range roots {
		for _, r := range rel {
			p := filepath.Join(root, r)
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				return p
			}
		}
		for _, pat := range []string{
			filepath.Join(root, `rail_apps`, `无畏契约*`, `ACLOS`, `Proxy`, `AclosGameProxy.exe`),
			filepath.Join(root, `无畏契约*`, `ACLOS`, `Proxy`, `AclosGameProxy.exe`),
		} {
			if hits, err := filepath.Glob(pat); err == nil && len(hits) > 0 {
				return hits[0]
			}
		}
	}
	return ""
}

// wuweiWorkDir 无畏上号器独立工作目录：它和圆梦共用 config.ini 会互相覆盖，必须分开放。
func wuweiWorkDir(base string) string { return filepath.Join(base, "wuwei") }

// wuweiLocalCandidates 本地已有的无畏上号器位置（开发者目录 / 绿色包）。
func wuweiLocalCandidates() []string {
	base, _ := os.Executable()
	exeDir := filepath.Dir(base)
	return []string{
		filepath.Join(exeDir, wuweiExeName),
		filepath.Join(exeDir, "other-tools", "无畏", wuweiExeName),
		filepath.Join(exeDir, "..", "other-tools", "无畏", wuweiExeName),
		filepath.Join(exeDir, "..", "..", "other-tools", "无畏", wuweiExeName),
		filepath.Join(exeDir, "资源", "上号器", "wuwei", wuweiExeName),
		filepath.Join(exeDir, "..", "资源", "上号器", "wuwei", wuweiExeName),
	}
}

// wuweiEnsureLocal 保证工作目录里有原始 无畏上号器.exe，返回路径。
func wuweiEnsureLocal(workDir string) (string, error) {
	dst := filepath.Join(workDir, wuweiExeName)
	if st, err := os.Stat(dst); err == nil && st.Size() > 0 {
		return dst, nil
	}
	for _, c := range wuweiLocalCandidates() {
		st, err := os.Stat(c)
		if err != nil || st.IsDir() || st.Size() == 0 {
			continue
		}
		if err := os.MkdirAll(workDir, 0o700); err != nil {
			return "", err
		}
		if err := copyFileLocal(c, dst); err != nil {
			return "", fmt.Errorf("复制上号器失败: %w", err)
		}
		return dst, nil
	}
	return "", fmt.Errorf("本地没有 %s", wuweiExeName)
}

// wuweiLaunchOffline 离线上号主流程：生成补丁副本 → 起本地卡密服务 → 写 config.ini → 启动上号器。
func wuweiLaunchOffline(launcher, workDir, ini, qq, zone string, logf func(string, ...any)) error {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	if strings.TrimSpace(ini) == "" {
		return fmt.Errorf("离线票据为空")
	}

	exe, err := wuweiPrepareExe(launcher, workDir)
	if err != nil {
		return err
	}
	stop, err := wuweiStartServer(ini, qq, zone)
	if err != nil {
		return err
	}
	defer stop()

	card, proxy, err := wuweiWriteConfig(workDir)
	if err != nil {
		logf("写 config.ini 失败: %v", err)
	} else {
		logf("占位卡密 %s", card)
		if proxy == "" {
			logf("首次运行：请在无畏上号器里选择 ACLOS\\Proxy\\AclosGameProxy.exe（选一次会记住）")
		}
	}

	cmd := exec.Command(exe)
	cmd.Dir = workDir
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动上号器失败: %w", err)
	}
	logf("无畏上号器已启动（本地卡密服务 %s）", wuweiListenAddr)
	return cmd.Wait()
}
