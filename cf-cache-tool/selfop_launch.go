package main

// 自研上号：启动游戏（1:1 复刻上号器行为）
//
// 已还原的完整流程（证据来自样本反汇编 + 落盘文件）：
//
//  1) 读 HKCU\SOFTWARE\Tencent\CrossFire InstallPath 得到游戏目录
//  2) 写 <dir>\TCLS\mmog_data.xml（Uin=QQ / WorldID=0.0.0.0 / Version=7.3.3.0）
//  3) 写 <dir>\TCLS\CF_Gamest（uin/tgpid/openid/open_key/ctt 纯文本）
//  4) 创建共享内存 Rail.SharedMemory.<PID>（544 字节）
//  5) 创建命名管道 \\.\pipe\rail_pipe_15
//  6) CreateProcessW("<dir>\x64\crossfire.exe", "...\" -q --disable-gpu --disable-gpu-compositing")
//  7) 注入 qwh_cf_hook.dll（VirtualAllocEx + WriteProcessMemory + CreateRemoteThread(LoadLibraryW)）
//
// 其中 1/2/3/6 是必需项；4/5/7 是上号器附加的"假 Rail"通道。

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	selfGameRelDir   = "TCLS"
	selfGameExeRel   = "x64\\crossfire.exe"
	selfRailVersion  = "7.3.3.0"
	selfRailVerConf  = "16:16:16:16"
	selfGameArgs     = "-q --disable-gpu --disable-gpu-compositing"
	selfSharedMemFmt = "Rail.SharedMemory.%d"
	selfPipeGuest    = `\\.\pipe\rail_pipe_15`

	// selfEnableRailHost 控制是否复刻新版那套"自己当 Rail Host"的环境：
	// WE_GAME_SDK 窗口 + Rail.SharedMemory.<hwnd> + rail_pipe_* 管道 + LoadLibrary rail。
	//
	// 实测（2026-11-02 18:15）：只写 CF_Gamest/mmog/version.ini，然后直接
	//   CreateProcess("x64\crossfire.exe", "-q --disable-gpu --disable-gpu-compositing")
	// 游戏正常起窗口 "Crossfire"，没有 WeGame Init Failed (2310)。
	//
	// 反过来，自研把窗口/共享内存/管道全建好后再启动，游戏 rail 会找到这套
	// "半成品 Host" 去握手，协议对不上就弹 2310 —— 也就是说 2310 不是
	// "找不到 Host" 引起的，而是"找到了但握手失败"引起的。所以这里默认关闭。
	selfEnableRailHost = false

	// selfEnableRailPatch 控制是否给 rail_api64.dll 打
	// RailNeedRestartAppForCheckingEnvironment 补丁。裸启动验证时没有这一步也能进游戏，
	// 而补丁 RVA 是按旧版 rail 推出来的，打在别的版本上可能把 rail 内部状态改坏。
	selfEnableRailPatch = false
)

// selfGameTicketData 是写进 CF_Gamest 的字段。
type selfGameTicketData struct {
	Uin       string // QQ
	Tgpid     string // uid
	Zone      string // 卡密备注里的大区名
	OpenID    string
	OpenKey   string
	PfKey     string
	CTT       string // game_client_key 与 ctt 相同
	GameST    string
	GameSTKey string
	ClientKey string
	STKey     string // qq_info.st_key，用于 qq_info_type=6 二次登录
	QQInfoSig string
	SK        string
	CT        string
	DT        string
	BizTicket string
}

// selfGameDir 从注册表读游戏安装目录。
func selfGameDir() (string, error) {
	const (
		hkeyCurrentUser = 0x80000001
		keyQueryValue   = 0x0001
	)
	advapi := syscall.NewLazyDLL("advapi32.dll")
	regOpen := advapi.NewProc("RegOpenKeyExW")
	regQuery := advapi.NewProc("RegQueryValueExW")
	regClose := advapi.NewProc("RegCloseKey")

	sub, _ := syscall.UTF16PtrFromString(`SOFTWARE\Tencent\CrossFire`)
	var hKey syscall.Handle
	r, _, _ := regOpen.Call(hkeyCurrentUser, uintptr(unsafe.Pointer(sub)), 0, keyQueryValue, uintptr(unsafe.Pointer(&hKey)))
	if r != 0 {
		return "", errors.New("找不到注册表 HKCU\\SOFTWARE\\Tencent\\CrossFire")
	}
	defer regClose.Call(uintptr(hKey))

	name, _ := syscall.UTF16PtrFromString("InstallPath")
	buf := make([]uint16, 512)
	size := uint32(len(buf) * 2)
	var typ uint32
	r, _, _ = regQuery.Call(uintptr(hKey), uintptr(unsafe.Pointer(name)), 0,
		uintptr(unsafe.Pointer(&typ)), uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r != 0 {
		return "", errors.New("注册表缺少 InstallPath")
	}
	dir := syscall.UTF16ToString(buf)
	if dir == "" {
		return "", errors.New("InstallPath 为空")
	}
	return dir, nil
}

// selfWriteMmogXML 写 TCLS\mmog_data.xml（gb2312 兼容 ASCII 部分）。
func selfWriteMmogXML(dir, qq string) (string, error) {
	p := filepath.Join(dir, selfGameRelDir, "mmog_data.xml")
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="gb2312" ?>
<VersionUpdateData>
    <VersionConf>%s</VersionConf>
    <ClientConf>
        <VersionData>
            <Version>%s</Version>
            <Uin>%s</Uin>
            <WorldID>0.0.0.0</WorldID>
            <BranchId>1</BranchId>
        </VersionData>
        <VersionDemoStr>%s</VersionDemoStr>
        <resv1>%s</resv1>
    </ClientConf>
</VersionUpdateData>`, selfRailVerConf, selfRailVersion, qq, selfRailVersion, selfRailVersion)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(p, []byte(xml), 0o644); err != nil {
		return "", err
	}
	return p, nil
}

// selfPatchMmogXMLs 对应样本 sub_14001AEE0 的四个 mmog_data.xml 改写：
// 保留原文件格式，只把 <Uin> 换成当前 QQ；WorldID 不动。
func selfPatchMmogXMLs(dir, qq string) {
	for _, rel := range []string{
		`NTCLS\mmog_data.xml`,
		`NTCLS\ntcls_config\mmog_data.xml`,
		`NTCLS\tgp_config\mmog_data.xml`,
	} {
		p := filepath.Join(dir, rel)
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		s := string(raw)
		if next, ok := selfReplaceXMLTag(s, "Uin", qq); ok {
			if err := os.WriteFile(p, []byte(next), 0o644); err == nil {
				selfLog("mmog Uin 已写 %s", p)
			}
		}
	}
	// 样本 SHCreateDirectoryExW(<base>\rail_user_data\2)。
	// <base> 是 WeGameApps 根，不是游戏目录；两处都建，兼容不同版本。
	_ = os.MkdirAll(filepath.Join(dir, "rail_user_data", "2"), 0o755)
	weGamRoot := filepath.Dir(filepath.Dir(dir))
	_ = os.MkdirAll(filepath.Join(weGamRoot, "rail_user_data", "2"), 0o755)
}

func selfReplaceXMLTag(s, tag, val string) (string, bool) {
	open, close := "<"+tag+">", "</"+tag+">"
	i := strings.Index(s, open)
	if i < 0 {
		return s, false
	}
	start := i + len(open)
	j := strings.Index(s[start:], close)
	if j < 0 {
		return s, false
	}
	return s[:start] + val + s[start+j:], true
}

// selfCFTicketBase64 把 OP 的 hex ctt 转成离线器 CF_Gamest 的 Base64 格式。
// 非 hex 或长度异常时原样返回，避免破坏已经是 Base64 的票据。
func selfCFTicketBase64(s string) string {
	s = strings.TrimSpace(s)
	if len(s) < 32 || len(s)%2 != 0 {
		return s
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return s
		}
	}
	raw, err := hex.DecodeString(s)
	if err != nil {
		return s
	}
	return base64.StdEncoding.EncodeToString(raw)
}

// selfWriteCFGamest 按离线器 14:35 实际落盘位置写票据：
// <game>\CF_Gamest 和 <game>\x64\CF_Gamest；TCLS 仅作兼容。
// 只写 TCLS 会让 x64 里的旧账号票据继续生效，导致 2310。
func selfWriteCFGamest(dir string, t *selfGameTicketData) (string, error) {
	// 离线器 CF_Gamest 里的 game_client_key/ctt 是 Base64；OP 回包给的是 hex。
	// 直接写 hex，游戏按 Base64 解会失败并 2310。
	gck := selfCFTicketBase64(t.CTT)
	var b strings.Builder
	fmt.Fprintf(&b, "uin=%s\n", t.Uin)
	fmt.Fprintf(&b, "tgpid=%s\n", t.Tgpid)
	fmt.Fprintf(&b, "openid=%s\n", t.OpenID)
	fmt.Fprintf(&b, "open_key=%s\n", t.OpenKey)
	fmt.Fprintf(&b, "pf_key=%s\n", t.PfKey)
	fmt.Fprintf(&b, "game_client_key=%s\n", gck)
	fmt.Fprintf(&b, "game_st=%s\n", t.GameST)
	fmt.Fprintf(&b, "game_st_key=%s\n", t.GameSTKey)
	fmt.Fprintf(&b, "client_key=%s\n", t.ClientKey)
	fmt.Fprintf(&b, "ctt=%s\n", gck)
	body := []byte(b.String())
	paths := []string{
		filepath.Join(dir, "CF_Gamest"),
		filepath.Join(dir, "x64", "CF_Gamest"),
		filepath.Join(dir, selfGameRelDir, "CF_Gamest"),
	}
	for _, p := range paths {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(p, body, 0o644); err != nil {
			return "", err
		}
		selfLog("写 CF_Gamest %s", p)
	}
	return paths[1], nil
}

// selfEnsureVersionIni 只在 version.ini 缺失或没有 -SSN 时才写。
//
// 新版上号器 sub_140011E50 只读不写；里面的 -SSN 是这台机器上次成功上号留下的
// 真实区服。自研如果每次按卡密备注重写，备注一旦不准就会把区服改错，游戏侧
// rail 首包跟着错，直接 2310。所以这里只补缺，不覆盖已有值。
func selfEnsureVersionIni(dir, ssn string) {
	if strings.TrimSpace(selfReadSSN(dir)) != "" {
		selfLog("version.ini 已存在 SSN=%s，保留不覆盖", selfReadSSN(dir))
		return
	}
	if strings.TrimSpace(ssn) == "" {
		selfLog("version.ini 缺失且无法确定区服，跳过写入")
		return
	}
	selfWriteVersionIni(dir, ssn)
}

// selfKillWeGame 对应样本 sub_1400099D0：清理 WeGame 客户端/启动器，避免它抢 Rail。
func selfKillWeGame() {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	openProc := k32.NewProc("OpenProcess")
	terminate := k32.NewProc("TerminateProcess")
	waitFor := k32.NewProc("WaitForSingleObject")
	closeHandle := k32.NewProc("CloseHandle")
	for _, p := range selfProcessList() {
		name := strings.ToLower(p.name)
		if name != "wegame.exe" && name != "wegamelauncher.exe" {
			continue
		}
		h, _, err := openProc.Call(uintptr(p.pid), 0, 0x100001) // TERMINATE|SYNCHRONIZE
		if h == 0 {
			selfLog("清理 WeGame pid=%d OpenProcess 失败: %v", p.pid, err)
			continue
		}
		ok, _, terr := terminate.Call(h, 0)
		waitFor.Call(h, 150)
		closeHandle.Call(h)
		if ok != 0 {
			selfLog("已清理 WeGame pid=%d", p.pid)
		} else {
			selfLog("清理 WeGame pid=%d 失败: %v", p.pid, terr)
		}
	}
}

// selfKillExistingCrossfire 对应样本 sub_140009860：
// 启动前枚举 crossfire.exe / crossfire_x64.exe，清理残留进程，最多等 2.5 秒。
// 不清理旧进程时，新进程会被单实例/互斥量顶成 2310。
func selfKillExistingCrossfire() {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	openProc := k32.NewProc("OpenProcess")
	terminate := k32.NewProc("TerminateProcess")
	waitFor := k32.NewProc("WaitForSingleObject")
	closeHandle := k32.NewProc("CloseHandle")

	kill := func() int {
		n := 0
		for _, p := range selfProcessList() {
			if !strings.EqualFold(p.name, "crossfire.exe") && !strings.EqualFold(p.name, "crossfire_x64.exe") {
				continue
			}
			h, _, err := openProc.Call(uintptr(p.pid), 0, 0x100001) // TERMINATE|SYNCHRONIZE
			if h == 0 {
				selfLog("清理旧 crossfire pid=%d OpenProcess 失败: %v", p.pid, err)
				continue
			}
			ok, _, terr := terminate.Call(h, 0)
			waitFor.Call(h, 150)
			closeHandle.Call(h)
			if ok != 0 {
				n++
				selfLog("已清理旧 crossfire pid=%d", p.pid)
			} else {
				selfLog("清理旧 crossfire pid=%d 失败: %v", p.pid, terr)
			}
		}
		return n
	}
	kill()
	deadline := time.Now().Add(2500 * time.Millisecond)
	for time.Now().Before(deadline) {
		found := false
		for _, p := range selfProcessList() {
			if strings.EqualFold(p.name, "crossfire.exe") || strings.EqualFold(p.name, "crossfire_x64.exe") {
				found = true
				break
			}
		}
		if !found {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// selfCreateSharedMem 创建 544 字节共享内存：
// +0 uint32=2，+4 WCHAR[]=rail_api64.dll 完整路径，+0x20C uint64=WE_GAME_SDK 窗口句柄。
// IDA: sub_14001AEE0 0x14001BBA6 写 v164[0]；v164 随后赋给 FileName，
// sub_14000F440 用 _wfopen(FileName) 按 PE 解析，所以它只能是 rail 路径，不是 QQ。
func selfCreateSharedMem(railDLL string) (syscall.Handle, uintptr, error) {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	createMap := k32.NewProc("CreateFileMappingW")
	mapView := k32.NewProc("MapViewOfFile")

	// 名字里的 %u 是 Rail host 窗口句柄，不是 PID（样本 0x14001BAxx 处用
	// [0x1401741F8] 存的 HWND 去拼）。用 PID 拼会 2310。
	name, _ := syscall.UTF16PtrFromString(fmt.Sprintf(selfSharedMemFmt, selfSharedMemName()))
	h, _, err := createMap.Call(
		uintptr(syscall.InvalidHandle), // INVALID_HANDLE_VALUE = 页文件
		0, 4, 0, uintptr(0x220),        // flProtect=PAGE_READWRITE，样本 0x14001BB11 传 4
		uintptr(unsafe.Pointer(name)),
	)
	if h == 0 {
		return 0, 0, fmt.Errorf("CreateFileMappingW 失败: %v", err)
	}
	view, _, err := mapView.Call(h, 0x000F001F /* FILE_MAP_ALL_ACCESS */, 0, 0, uintptr(0x220))
	if view == 0 {
		syscall.CloseHandle(syscall.Handle(h))
		return 0, 0, fmt.Errorf("MapViewOfFile 失败: %v", err)
	}
	mem := unsafe.Slice((*byte)(unsafe.Pointer(view)), 0x220)
	for i := range mem {
		mem[i] = 0
	}
	binary.LittleEndian.PutUint32(mem[0:4], 2)
	// +4 起写 UTF-16 rail_api64.dll 完整路径，最多 129 个 WCHAR。
	rail := syscall.StringToUTF16(railDLL)
	for i := 0; i < len(rail) && 4+i*2+1 < 0x20C; i++ {
		binary.LittleEndian.PutUint16(mem[4+i*2:], rail[i])
	}
	// +0x20C 写 Rail host 窗口句柄（样本 0x14001BBAB: mov [rsi+20Ch], rdi）。
	// rail_api64 从这里取 HWND 再去找窗口；不写就是 0 → 初始化失败 → 2310。
	hwnd := selfEnsureRailWindow()
	if hwnd == 0 {
		hwnd = uintptr(selfSharedMemName())
	}
	binary.LittleEndian.PutUint64(mem[0x20C:], uint64(hwnd))
	return syscall.Handle(h), view, nil
}

// selfLaunchGame 写票据文件并启动 crossfire.exe。
func selfLaunchGame(dir string, t *selfGameTicketData) (int, error) {
	if dir == "" {
		var err error
		dir, err = selfGameDir()
		if err != nil {
			return 0, err
		}
	}
	if norm, err := selfNormalizeGameDir(dir); err == nil {
		dir = norm
	}
	exe := filepath.Join(dir, selfGameExeRel)
	if _, err := os.Stat(exe); err != nil {
		return 0, fmt.Errorf("找不到游戏主程序: %w", err)
	}
	// 新版 sub_140007640 / sub_140008050 启动前先枚举 crossfire.exe：
	// 已经有游戏在跑就不再拉起第二份，否则 rail 管道/共享内存会互相抢。
	for _, p := range selfProcessList() {
		if strings.EqualFold(p.name, "crossfire.exe") && p.pid != uint32(os.Getpid()) {
			return 0, fmt.Errorf("游戏已在运行（crossfire.exe pid=%d），请先退出再上号", p.pid)
		}
	}
	// 样本 sub_14001AEE0 读现有 CF_Gamest 的 uin= 直接启动；
	// 离线器 16:02 已经为当前账号写入有效 open_key/ctt，不能再被我们的
	// 二次登录票据覆盖，否则覆盖后立刻 2310。
	reuseGamest := false
	if old, p, err := selfReadCFGamest(dir); err == nil &&
		strings.TrimSpace(old.Uin) == strings.TrimSpace(t.Uin) &&
		strings.TrimSpace(old.OpenKey) != "" && strings.TrimSpace(old.CTT) != "" {
		selfLog("复用现有 CF_Gamest %s uin=%s open_key=%s", p, old.Uin, old.OpenKey)
		t = old
		reuseGamest = true
	}

	selfLog("==== 启动 start dir=%s exe=%s", dir, exe)
	ssn := selfResolveSSN(dir, t.Zone)
	selfEnsureVersionIni(dir, ssn)

	// 新版上号器的启动前检查：version.ini 必须带 -SSN:<当前区>/-rotate:0/-cf:1，
	// 且不含 -HostIP=；不满足它直接不启动，我们提前报错避免等到 2310。
	if err := selfValidateVersionIni(dir, ssn); err != nil {
		selfLog("version.ini 校验失败: %v", err)
		return 0, err
	}
	// 新版还会在 <游戏根>\ntop_v17_dll\x64 写一个 1 字节假 rail_api64.dll，
	// 让 CF 的 LoadLibrary 失败、跳过 WeGame 初始化（否则 2310）。
	selfEnsureFakeRail(dir)
	if selfEnableRailHost {
		// 新版在启动 crossfire.exe 之前先在自己进程里加载真 rail
		// （IDA sub_14009F540 会检查该模块，缺失直接放弃启动）。
		selfLoadRailHost(dir)
		selfLoadRailHello(t.Uin, t.Tgpid, ssn)
	}

	if _, err := selfWriteMmogXML(dir, t.Uin); err != nil {
		selfLog("写 mmog_data.xml 失败: %v", err)
		return 0, fmt.Errorf("写 mmog_data.xml 失败: %w", err)
	}
	selfPatchMmogXMLs(dir, t.Uin)
	if !reuseGamest {
		if _, err := selfWriteCFGamest(dir, t); err != nil {
			selfLog("写 CF_Gamest 失败: %v", err)
			return 0, fmt.Errorf("写 CF_Gamest 失败: %w", err)
		}
	}
	selfLog("票据文件已写 qq=%s hasSig=%v reuse=%v", t.Uin, t.GameSTKey != "", reuseGamest)

	if selfEnableRailHost {
		fakeRail := filepath.Join(dir, "ntop_v17_dll", "x64", "rail_api64.dll")
		if h, v, err := selfCreateSharedMem(fakeRail); err == nil {
			selfHoldSharedMem(h, v)
			selfLog("共享内存 ok name=Rail.SharedMemory.%d hwnd=%#x view=%#x", selfSharedMemName(), selfEnsureRailWindow(), v)
		} else {
			selfLog("共享内存失败: %v", err)
		}
	}

	// 样本顺序：先把 qwh_cf_hook.dll 落盘，再 CreateProcess。
	// 当前代码以前在启动后才写，游戏已经开始初始化 Rail，来不及。
	base := filepath.Dir(exe)
	dll, err := selfEnsureHookDLL(base)
	if err != nil {
		selfLog("释放 hook dll 失败: %v", err)
		return 0, fmt.Errorf("释放 hook dll 失败: %w", err)
	}
	selfLog("hook dll=%s", dll)

	// 不主动注入：离线器是等游戏发 inject_hook cmd=1 后才注入。
	// 先把 dll 路径挂到全局，游戏一连上静态管道就可能立刻发 cmd=1。
	selfResetHookInjection()
	selfHookPath = dll

	// 样本先建 4 个 *_15 静态管道，CreateProcess 后再建游戏 PID 的 8 个动态管道。
	staticStop := func() {}
	if selfEnableRailHost {
		staticStop = selfServeStaticRailPipes()
	}
	// 工作目录必须是游戏根目录；并且必须降权启动 —— 管理员令牌下 ACE 会判定
	// 启动环境异常（elevated=true → 0xDEADBEEF / 2310），普通令牌才能进。
	proc, err := selfStartGameAsUser(exe, dir, selfGameArgs)
	if err != nil {
		staticStop()
		selfLog("CreateProcessW 失败: %v", err)
		return 0, fmt.Errorf("启动游戏失败: %w", err)
	}
	pid := int(proc.ProcessID)
	selfLog("CreateProcessW ok pid=%d", pid)

	// 句柄必须在启动动态管道前挂好，否则 cmd=1 可能先到。
	selfHookProcHandle = proc.Process
	selfHookGamePID = pid
	if selfEnableRailHost {
		selfLog("等待游戏 inject_hook cmd=1 pid=%d", pid)
		_ = selfServeGameRailPipes(pid)
		go selfWatchGameRailPipes(pid, 45*time.Second)
	}

	// 样本在 CreateProcess 后启动 sub_14000FC70，持续给所有加载 rail_api64.dll
	// 的进程打 RailNeedRestartAppForCheckingEnvironment 补丁。少了这段，子进程
	// 重新加载 rail 后仍会弹 2310。
	if selfEnableRailPatch {
		go selfPatchRailLoop(filepath.Join(base, "rail_api64.dll"), 10*time.Minute)
	}

	syscall.CloseHandle(proc.Thread)
	return pid, nil
}

// selfGameProc 保存 CreateProcessW 的原始句柄；注入必须复用这些句柄。
type selfGameProc struct {
	Process   syscall.Handle
	Thread    syscall.Handle
	ProcessID uint32
}

// selfStartGame 对应样本 sub_14001AEE0 的 CreateProcessW：
// appName=<exe>，commandLine="<exe>" <args>，creationFlags=0x400。
func selfStartGame(exe, workDir, args string) (*selfGameProc, error) {
	app, err := syscall.UTF16PtrFromString(exe)
	if err != nil {
		return nil, err
	}
	cmdLine, err := syscall.UTF16PtrFromString(fmt.Sprintf(`"%s" %s`, exe, args))
	if err != nil {
		return nil, err
	}
	dir, err := syscall.UTF16PtrFromString(workDir)
	if err != nil {
		return nil, err
	}
	var si syscall.StartupInfo
	si.Cb = uint32(unsafe.Sizeof(si))
	var pi syscall.ProcessInformation
	if err := syscall.CreateProcess(app, cmdLine, nil, nil, false, 0x400, nil, dir, &si, &pi); err != nil {
		return nil, err
	}
	return &selfGameProc{Process: pi.Process, Thread: pi.Thread, ProcessID: pi.ProcessId}, nil
}

var (
	selfSharedMemMu sync.Mutex
	selfSharedMemH  syscall.Handle
	selfSharedMemV  uintptr
)

// selfHoldSharedMem 把共享内存句柄挂在进程级，避免被 GC/函数返回回收。
func selfHoldSharedMem(h syscall.Handle, v uintptr) {
	selfSharedMemMu.Lock()
	old, oldV := selfSharedMemH, selfSharedMemV
	selfSharedMemH, selfSharedMemV = h, v
	selfSharedMemMu.Unlock()
	if old != 0 {
		k32 := syscall.NewLazyDLL("kernel32.dll")
		k32.NewProc("UnmapViewOfFile").Call(oldV)
		syscall.CloseHandle(old)
	}
}

// selfInjectHookHandle 对应样本 sub_14001AEE0 的注入段：
// 直接复用 CreateProcessW 返回的 hProcess，绝不重新 OpenProcess。
func selfInjectHookHandle(hProcess syscall.Handle, dllPath string) error {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	virtualAllocEx := k32.NewProc("VirtualAllocEx")
	writeProcMem := k32.NewProc("WriteProcessMemory")
	createRemote := k32.NewProc("CreateRemoteThread")
	waitFor := k32.NewProc("WaitForSingleObject")
	getExitCode := k32.NewProc("GetExitCodeThread")
	getModule := k32.NewProc("GetModuleHandleW")
	getProcAddr := k32.NewProc("GetProcAddress")

	const (
		memCommitReserve  = 0x3000
		pageExecReadWrite = 0x40
	)
	hProc := uintptr(hProcess)
	if hProc == 0 {
		return errors.New("进程句柄为空")
	}

	dllUTF16, _ := syscall.UTF16FromString(dllPath)
	size := uintptr(len(dllUTF16) * 2)

	remote, _, _ := virtualAllocEx.Call(hProc, 0, size, memCommitReserve, pageExecReadWrite)
	if remote == 0 {
		return errors.New("VirtualAllocEx 失败")
	}
	var written uintptr
	ok, _, _ := writeProcMem.Call(hProc, remote, uintptr(unsafe.Pointer(&dllUTF16[0])), size, uintptr(unsafe.Pointer(&written)))
	if ok == 0 {
		return errors.New("WriteProcessMemory 失败")
	}

	kern, _ := syscall.UTF16PtrFromString("kernel32.dll")
	hKern, _, _ := getModule.Call(uintptr(unsafe.Pointer(kern)))
	if hKern == 0 {
		return errors.New("GetModuleHandleW(kernel32) 失败")
	}
	procName, _ := syscall.BytePtrFromString("LoadLibraryW")
	loadLib, _, _ := getProcAddr.Call(hKern, uintptr(unsafe.Pointer(procName)))
	if loadLib == 0 {
		return errors.New("GetProcAddress(LoadLibraryW) 失败")
	}

	hThread, _, _ := createRemote.Call(hProc, 0, 0, loadLib, remote, 0, 0)
	if hThread == 0 {
		return errors.New("CreateRemoteThread 失败")
	}
	defer syscall.CloseHandle(syscall.Handle(hThread))
	waitFor.Call(hThread, 10000)
	var exitCode uint32
	if ok, _, _ := getExitCode.Call(hThread, uintptr(unsafe.Pointer(&exitCode))); ok != 0 {
		selfLog("LoadLibraryW 远程返回 %#x", exitCode)
		if exitCode == 0 {
			return errors.New("LoadLibraryW 返回 0（注入未生效）")
		}
	}
	return nil
}
