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
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
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
)

// selfGameTicketData 是写进 CF_Gamest 的字段。
type selfGameTicketData struct {
	Uin      string // QQ
	Tgpid    string // uid
	OpenID   string
	OpenKey  string
	PfKey    string
	CTT      string // game_client_key 与 ctt 相同
	GameST   string
	GameSTKey string
	ClientKey string
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

// selfWriteCFGamest 写 TCLS\CF_Gamest（纯文本票据，行序与上号器一致）。
func selfWriteCFGamest(dir string, t *selfGameTicketData) (string, error) {
	p := filepath.Join(dir, selfGameRelDir, "CF_Gamest")
	var b strings.Builder
	fmt.Fprintf(&b, "uin=%s\n", t.Uin)
	fmt.Fprintf(&b, "tgpid=%s\n", t.Tgpid)
	fmt.Fprintf(&b, "openid=%s\n", t.OpenID)
	fmt.Fprintf(&b, "open_key=%s\n", t.OpenKey)
	fmt.Fprintf(&b, "pf_key=%s\n", t.PfKey)
	fmt.Fprintf(&b, "game_client_key=%s\n", t.CTT)
	fmt.Fprintf(&b, "game_st=%s\n", t.GameST)
	fmt.Fprintf(&b, "game_st_key=%s\n", t.GameSTKey)
	fmt.Fprintf(&b, "client_key=%s\n", t.ClientKey)
	fmt.Fprintf(&b, "ctt=%s\n", t.CTT)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(p, []byte(b.String()), 0o644); err != nil {
		return "", err
	}
	return p, nil
}

// selfCreateSharedMem 创建 544 字节共享内存：+0 uint32=2，+4 WCHAR[]=QQ，+0x20C uint64。
// 返回句柄与视图指针，调用方负责释放。
func selfCreateSharedMem(qq string) (syscall.Handle, uintptr, error) {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	createMap := k32.NewProc("CreateFileMappingW")
	mapView := k32.NewProc("MapViewOfFile")

	// 名字里的 %u 是 Rail host 窗口句柄，不是 PID（样本 0x14001BAxx 处用
	// [0x1401741F8] 存的 HWND 去拼）。用 PID 拼会 2310。
	name, _ := syscall.UTF16PtrFromString(fmt.Sprintf(selfSharedMemFmt, selfSharedMemName()))
	h, _, err := createMap.Call(
		uintptr(syscall.InvalidHandle), // INVALID_HANDLE_VALUE = 页文件
		0, 0, 0, uintptr(0x220),
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
	// +4 起写 UTF-16 QQ，最多 129 个 WCHAR
	q := syscall.StringToUTF16(qq)
	for i := 0; i < len(q) && 4+i*2+1 < 0x20C; i++ {
		binary.LittleEndian.PutUint16(mem[4+i*2:], q[i])
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
	selfLog("==== 启动 start dir=%s exe=%s", dir, exe)

	if _, err := selfWriteMmogXML(dir, t.Uin); err != nil {
		selfLog("写 mmog_data.xml 失败: %v", err)
		return 0, fmt.Errorf("写 mmog_data.xml 失败: %w", err)
	}
	if _, err := selfWriteCFGamest(dir, t); err != nil {
		selfLog("写 CF_Gamest 失败: %v", err)
		return 0, fmt.Errorf("写 CF_Gamest 失败: %w", err)
	}
	selfLog("票据文件已写 qq=%s hasSig=%v", t.Uin, t.GameSTKey != "")

	// 共享内存必须活到游戏退出，不能像之前那样 defer 掉。
	// Rail 在初始化后会反复读这块内存，提前 Unmap 会让它 2310。
	if h, v, err := selfCreateSharedMem(t.Uin); err == nil {
		selfHoldSharedMem(h, v)
		selfLog("共享内存 ok name=Rail.SharedMemory.%d hwnd=%#x view=%#x", selfSharedMemName(), selfEnsureRailWindow(), v)
	} else {
		selfLog("共享内存失败: %v", err)
	}

	// 管道：host 侧名字带自己的 PID（rail_named_pipe_host_<pid>），
	// 游戏起来后会创建 rail_named_pipe_guest_<gamePid>，我们作为 host 去连。
	pipeStop := selfServeHostPipe()

	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine:       fmt.Sprintf(`"%s" %s`, exe, selfGameArgs),
		CreationFlags: 0x00000400, // CREATE_UNICODE_ENVIRONMENT
	}
	if err := cmd.Start(); err != nil {
		selfLog("CreateProcess 失败: %v", err)
		pipeStop()
		return 0, fmt.Errorf("启动游戏失败: %w", err)
	}
	pid := cmd.Process.Pid
	selfLog("CreateProcess ok pid=%d cmdline=\"%s\"", pid, cmd.SysProcAttr.CmdLine)

	// 作为 host 连游戏的 guest 管道（游戏自己创建，名字带它的 PID）
	selfConnectGuestPipe(pid)

	// 注入 Hook DLL（上号器在第 7 步做，缺这一步 Rail 仍会 2310）
	//
	// 注意 base：样本 sub_140019270 是先拼 "<base>\crossfire.exe" 再落到
	// "<base>\ntop_v17_dll\x64\rail_api64.dll"（0x14001A578 建 exe 路径，
	// 0x14001A51C 写假 dll）。所以 base 是 exe 所在目录（<游戏>\x64），
	// 不是游戏根目录。放错一级 CF 就看不见假 dll，仍然 2310。
	base := filepath.Dir(exe)
	if dll, err := selfEnsureHookDLL(base); err != nil {
		selfLog("释放 hook dll 失败: %v", err)
		fmt.Printf("注入跳过: %v\n", err)
	} else {
		fake := filepath.Join(base, "ntop_v17_dll", "x64", "rail_api64.dll")
		st, serr := os.Stat(fake)
		selfLog("hook dll=%s fakeRail=%s statErr=%v", dll, fake, serr)
		if serr == nil {
			selfLog("fakeRail size=%d", st.Size())
		}
		// 兼容旧位置（游戏根目录），老版本已经在那儿放过文件
		_, _ = selfEnsureHookDLL(dir)
		if err := selfInjectHook(pid, dll); err != nil {
			selfLog("注入失败: %v", err)
			fmt.Printf("注入失败: %v\n", err)
		} else {
			selfLog("注入成功 pid=%d dll=%s", pid, dll)
		}
	}

	return pid, nil
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



// selfInjectHook 把 qwh_cf_hook.dll 注入目标进程（VirtualAllocEx + WriteProcessMemory + CreateRemoteThread）。
func selfInjectHook(pid int, dllPath string) error {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	openProc := k32.NewProc("OpenProcess")
	virtualAllocEx := k32.NewProc("VirtualAllocEx")
	writeProcMem := k32.NewProc("WriteProcessMemory")
	createRemote := k32.NewProc("CreateRemoteThread")
	waitFor := k32.NewProc("WaitForSingleObject")
	getModule := k32.NewProc("GetModuleHandleW")
	getProcAddr := k32.NewProc("GetProcAddress")

	const (
		processAllAccess = 0x001FFFFF
		memCommitReserve = 0x3000
		pageExecReadWrite = 0x40
	)

	hProc, _, _ := openProc.Call(uintptr(pid), 0, processAllAccess)
	if hProc == 0 {
		return errors.New("OpenProcess 失败")
	}
	defer syscall.CloseHandle(syscall.Handle(hProc))

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
	return nil
}
