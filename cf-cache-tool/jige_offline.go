package main

// 牛逼上号器（鸡哥直登）离线桥接。
//
// 已还原的协议（IDA sub_14001AA40 + 动态抓包）：
//
//	POST http://111.170.7.81:8787/api/card/boot    → {"ok":true,"banned":false}
//	POST http://111.170.7.81:8787/api/card/redeem  {"card_key":"..."}
//	     → {"ok":true,"qq":"...","zone":"...","key":RSA-OAEP-SHA256(AESKey),
//	        "iv":12B,"tag":16B,"data":AES-256-GCM(ini 明文)}
//	POST http://111.170.7.81:8787/api/card/log     → {"ok":true}
//
// RSA 是 2048 位，密钥对与圆梦上号器不同；私钥 blob 在
// 鸡哥直登.exe 的 .rdata 0x1400d8480（magic "RSA3"，1179 字节），
// 这里只用它对应的公钥封信封。
//
// 该程序走 WinHTTP（不走 IE 代理），所以用 netsh winhttp set proxy
// 把请求劫持到本地；退出后按备份的注册表值原样还原。

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const jigeRSA2048N = "c400c74a53126b57cc3fba38f4530ab3cbe656f32fb91229079214b883e467990dae5485c1495067a1f83609f8c72e2499eec824de465fcacd04842cdad559624510e8688695eb6c8afa087c598f3a2c053d5ffed8a77db380748afeb213c00d678c49bf9dabdbca3b79785ee71332c06019a8b6c8d3427a97b1f1fe1543d62c4482f292355d4d883175375b4e7bbaaaca7ac58aab403a0af8775a1bdf24c624cb50e782bf4a3ec40fe2c6baa87ca696657375b55c489cd7001fda06c6b9d06cb21e6b2fbf414238d2562e0e354f9eb5e83188400b6d9a2d4d84294a79360ab5556fd69fead886282e0499a1cf96f710629e62c0d15df87985ba83fe530e6ebb"

func jigePub() *rsa.PublicKey {
	n, _ := new(big.Int).SetString(jigeRSA2048N, 16)
	return &rsa.PublicKey{N: n, E: 65537}
}

// jigeSeal 把 ini 明文封成鸡哥能解的信封。
func jigeSeal(plain []byte) (map[string]string, error) {
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, err
	}
	keyCt, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, jigePub(), aesKey, nil)
	if err != nil {
		return nil, err
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
	b64 := base64.StdEncoding.EncodeToString
	return map[string]string{"key": b64(keyCt), "iv": b64(iv), "tag": b64(tag), "data": b64(data)}, nil
}

// jigeStartServer 起本地假卡密服务，返回监听地址。
func jigeStartServer(ini, qq, zone string) (string, func(), error) {
	env, err := jigeSeal([]byte(ini))
	if err != nil {
		return "", nil, err
	}
	mux := http.NewServeMux()
	reply := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	mux.HandleFunc("/api/card/boot", func(w http.ResponseWriter, r *http.Request) {
		reply(w, map[string]any{"ok": true, "banned": false})
	})
	mux.HandleFunc("/api/card/redeem", func(w http.ResponseWriter, r *http.Request) {
		out := map[string]any{
			"ok": true, "qq": qq, "zone": zone, "note": zone, "extra": zone,
			"activated": true, "valid_hours": 87600,
			"expires_at": time.Now().Add(10 * 365 * 24 * time.Hour).UnixMilli(),
		}
		for k, v := range env {
			out[k] = v
		}
		reply(w, out)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reply(w, map[string]any{"ok": true})
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}
	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()
	return ln.Addr().String(), func() { _ = srv.Close() }, nil
}

// ---- WinHTTP 代理备份/设置/还原（注册表二进制值最稳，不受系统语言影响） ----

const (
	hkeyLocalMachine = 0x80000002
	regSZ            = 1
	regBinary        = 3
)

var (
	advapi32          = syscall.NewLazyDLL("advapi32.dll")
	pRegOpenKeyEx     = advapi32.NewProc("RegOpenKeyExW")
	pRegQueryValueEx  = advapi32.NewProc("RegQueryValueExW")
	pRegSetValueEx    = advapi32.NewProc("RegSetValueExW")
	pRegDeleteValue   = advapi32.NewProc("RegDeleteValueW")
	pRegCloseKey      = advapi32.NewProc("RegCloseKey")
)

const winhttpConnKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\Internet Settings\Connections`

func winhttpOpen() syscall.Handle {
	sub, _ := syscall.UTF16PtrFromString(winhttpConnKey)
	var h syscall.Handle
	r, _, _ := pRegOpenKeyEx.Call(hkeyLocalMachine, uintptr(unsafe.Pointer(sub)), 0, 0x0002|0x0001, uintptr(unsafe.Pointer(&h)))
	if r != 0 {
		return 0
	}
	return h
}

// winhttpBackup 读原始 WinHttpSettings；existed=false 表示本来就没有（直连）。
func winhttpBackup() (data []byte, existed bool) {
	h := winhttpOpen()
	if h == 0 {
		return nil, false
	}
	defer pRegCloseKey.Call(uintptr(h))
	name, _ := syscall.UTF16PtrFromString("WinHttpSettings")
	var typ, size uint32
	if r, _, _ := pRegQueryValueEx.Call(uintptr(h), uintptr(unsafe.Pointer(name)), 0,
		uintptr(unsafe.Pointer(&typ)), 0, uintptr(unsafe.Pointer(&size))); r != 0 || size == 0 {
		return nil, false
	}
	buf := make([]byte, size)
	if r, _, _ := pRegQueryValueEx.Call(uintptr(h), uintptr(unsafe.Pointer(name)), 0,
		uintptr(unsafe.Pointer(&typ)), uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size))); r != 0 {
		return nil, false
	}
	return buf[:size], true
}

func winhttpRestore(data []byte, existed bool) {
	h := winhttpOpen()
	if h == 0 {
		return
	}
	defer pRegCloseKey.Call(uintptr(h))
	name, _ := syscall.UTF16PtrFromString("WinHttpSettings")
	if !existed || len(data) == 0 {
		pRegDeleteValue.Call(uintptr(h), uintptr(unsafe.Pointer(name)))
		return
	}
	pRegSetValueEx.Call(uintptr(h), uintptr(unsafe.Pointer(name)), 0, uintptr(regBinary),
		uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
}

func winhttpSetProxy(addr string) error {
	cmd := exec.Command("netsh", "winhttp", "set", "proxy", addr)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("netsh winhttp set proxy 失败: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func winhttpResetProxy() error {
	cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("netsh winhttp reset proxy 失败: %v (%s)", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// 当前是否占用着 WinHTTP 代理：关窗口时后台 goroutine 会被直接杀掉，
// defer 不一定来得及跑，所以退出路径要能自己还原。
var (
	jigeProxyMu    sync.Mutex
	jigeProxyOn    bool
	jigeProxyData  []byte
	jigeProxyHad   bool
)

func jigeMarkProxyActive(data []byte, existed bool) {
	jigeProxyMu.Lock()
	jigeProxyOn, jigeProxyData, jigeProxyHad = true, data, existed
	jigeProxyMu.Unlock()
}

func jigeMarkProxyIdle() {
	jigeProxyMu.Lock()
	jigeProxyOn = false
	jigeProxyMu.Unlock()
}

// jigeRestoreProxyIfActive 退出兜底：还占着 WinHTTP 代理就按备份还原。
func jigeRestoreProxyIfActive() bool {
	jigeProxyMu.Lock()
	active, data, had := jigeProxyOn, jigeProxyData, jigeProxyHad
	jigeProxyOn = false
	jigeProxyMu.Unlock()
	if !active {
		return false
	}
	_ = winhttpResetProxy()
	winhttpRestore(data, had)
	return true
}

// jigeLaunchOffline 起服务 → 接管 WinHTTP 代理 → 启动鸡哥 → 退出后还原。
func jigeLaunchOffline(launcher, ini, qq, zone string, logf func(string, ...any)) error {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	if strings.TrimSpace(ini) == "" {
		return fmt.Errorf("离线票据为空")
	}
	addr, stop, err := jigeStartServer(ini, qq, zone)
	if err != nil {
		return fmt.Errorf("起本地卡密服务失败: %v", err)
	}
	defer stop()
	logf("本地卡密服务 %s（QQ %s / %s）", addr, qq, zone)

	backup, existed := winhttpBackup()
	jigeMarkProxyActive(backup, existed)
	defer func() {
		jigeMarkProxyIdle()
		winhttpRestore(backup, existed)
	}()

	if err := winhttpSetProxy(addr); err != nil {
		return err
	}
	logf("WinHTTP 代理已指向 %s", addr)

	cmd := exec.Command(launcher)
	cmd.Dir = filepath.Dir(launcher)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动上号器失败: %v", err)
	}
	logf("已启动上号器: %s", launcher)
	go jigeAutoDrive(uint32(cmd.Process.Pid), jigeRandomCard(), logf)
	err = cmd.Wait()
	logf("上号器已退出，正在还原 WinHTTP 代理")
	if e := winhttpResetProxy(); e != nil {
		// 已经按备份还原过注册表，这里 reset 只是让 netsh 的缓存同步，失败不影响
		_ = e
	}
	winhttpRestore(backup, existed)
	return err
}

// ymmLaunchJige 是「新上号」：离线票据 + 本地假卡密服务 + WinHTTP 代理。
func (a *guiApp) ymmLaunchJige(row ymmRow) {
	if row.QQ == "" || row.CardKey == "" {
		a.ymm.logf("该行缺少 QQ 或卡密")
		return
	}
	offline, err := ymmBuildOffline(a.ymm.api, a.ymm.snapshot(), row.CardKey, a.ymm.isDump())
	if err != nil {
		a.ymm.logf("准备失败: %v", err)
		return
	}
	zone := row.Area
	if zone == "" {
		zone = ymmZoneFromNote(row.Note)
	}
	if zone == "" {
		zone = "河南一区"
	}
	a.ymm.logf("新上号（牛逼）：QQ %s（卡密 %s / %s）", row.QQ, row.CardKey, zone)

	go func() {
		a.ui(a.spinBegin)
		defer a.ui(a.spinEnd)

		launcher, err := a.ensureLoaderFor(findProfile("2"))
		if err != nil {
			// 云端拿不到时退回本地原始文件名
			launcher = jigeFindLauncher(a.dir, findProfile("2"))
			if launcher == "" {
				a.ui(func() { a.ymm.logf("获取上号器失败: %v", err) })
				return
			}
		}

		text, qq, err := offline.ticketFor(row.CardKey)
		if err != nil {
			a.ui(func() { a.ymm.logf("取票据失败: %v", err) })
			return
		}
		appendHistory(a.dir, histEntry{
			TS:     time.Now().UnixMilli(),
			Action: "offline",
			QQ:     row.QQ,
			Card:   row.CardKey,
			Note:   row.Note,
			Zone:   zone,
			Result: "已启动",
		})
		err = jigeLaunchOffline(launcher, text, qq, zone, func(f string, args ...any) {
			a.ui(func() { a.ymm.logf(f, args...) })
		})
		if err != nil {
			a.ui(func() { a.ymm.logf("失败: %v", err) })
			return
		}
		a.ui(func() { a.ymm.logf("牛逼上号器已退出，WinHTTP 代理已还原") })
	}()
}

// jigeRandomCard 生成 16 位纯英文卡密：鸡哥只做格式校验，真数据由本地服务器给。
func jigeRandomCard() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 16)
	raw := make([]byte, 16)
	_, _ = rand.Read(raw)
	for i := range b {
		b[i] = letters[int(raw[i])%len(letters)]
	}
	return string(b)
}

// jigeAutoDrive 等鸡哥主窗口出来后自动填卡密并点「验证上号」。
func jigeAutoDrive(pid uint32, card string, logf func(string, ...any)) {
	u := syscall.NewLazyDLL("user32.dll")
	enumWindows := u.NewProc("EnumWindows")
	enumChild := u.NewProc("EnumChildWindows")
	getText := u.NewProc("GetWindowTextW")
	getClass := u.NewProc("GetClassNameW")
	getPid := u.NewProc("GetWindowThreadProcessId")
	isVisible := u.NewProc("IsWindowVisible")
	sendMsg := u.NewProc("SendMessageW")

	text := func(h uintptr, proc *syscall.LazyProc) string {
		var buf [256]uint16
		n, _, _ := proc.Call(h, uintptr(unsafe.Pointer(&buf[0])), 256)
		return syscall.UTF16ToString(buf[:n])
	}
	owner := func(h uintptr) uint32 {
		var p uint32
		getPid.Call(h, uintptr(unsafe.Pointer(&p)))
		return p
	}

	deadline := time.Now().Add(25 * time.Second)
	for time.Now().Before(deadline) {
		var mainWnd uintptr
		enumWindows.Call(syscall.NewCallback(func(h uintptr, _ uintptr) uintptr {
			if owner(h) != pid {
				return 1
			}
			vis, _, _ := isVisible.Call(h)
			if vis == 1 && text(h, getText) == "CF上号器" {
				mainWnd = h
				return 0
			}
			return 1
		}), 0)
		if mainWnd != 0 {
			var edit, btn uintptr
			enumChild.Call(mainWnd, syscall.NewCallback(func(ch uintptr, _ uintptr) uintptr {
				cls := text(ch, getClass)
				if cls == "Edit" && edit == 0 {
					edit = ch
				}
				if cls == "Button" && text(ch, getText) == "验证上号" {
					btn = ch
				}
				return 1
			}), 0)
			if edit != 0 && btn != 0 {
				cv, _ := syscall.UTF16PtrFromString(card)
				sendMsg.Call(edit, 0x000C, 0, uintptr(unsafe.Pointer(cv)))
				sendMsg.Call(btn, 0x00F5, 0, 0)
				logf("已自动填入卡密并点上号（%s）", card)
				return
			}
		}
		time.Sleep(400 * time.Millisecond)
	}
	logf("自动填卡密超时，请手动在窗口里点「验证上号」")
}

// jigeFindLauncher 找牛逼上号器 exe：优先当前 loader，其次工作目录/资源目录。
func jigeFindLauncher(workDir string, p *loaderProfile) string {
	if p != nil {
		if exe, err := ensureLocalLoader(workDir, p); err == nil {
			return exe
		}
	}
	if p != nil {
		cand := filepath.Join(workDir, p.ExeName)
		if _, err := os.Stat(cand); err == nil {
			return cand
		}
	}
	return ""
}
