package main

// 降权启动 crossfire.exe。
//
// 实测（2026-11-02 18:2x）：
//
//	普通权限 CreateProcess → 游戏正常起窗口 "Crossfire"，无 2310
//	管理员权限 CreateProcess → 进程 elevated=true，12 秒后退出码 0xDEADBEEF / 或弹 2310
//
// 游戏自己的 manifest 是 asInvoker，但 ACE/TP 在**提升令牌**下会认为启动环境异常
// （正常 CF 永远由 WeGame 以普通用户身份拉起），于是 rail 初始化直接失败。
//
// CG 需要管理员（netsh WinHTTP 代理、杀进程），所以从 CG 里直接 CreateProcess
// 出来的游戏一定是 elevated。新版 CF上号器 的 manifest 没有 requireAdministrator，
// 它自己就是普通权限，所以它拉起的游戏天然是普通权限 —— 这就是"新上号行、自研不行"
// 的真正原因。
//
// 这里用 explorer.exe 的令牌（桌面用户的普通令牌）CreateProcessWithTokenW 启动游戏。

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

const (
	tokenAssignPrimary = 0x0001
	tokenDuplicate     = 0x0002
	tokenQuery         = 0x0008
	tokenAdjustDefault = 0x0080
	tokenAdjustSession = 0x0100
	tokenAllAccess     = 0xF01FF

	securityImpersonation = 2
	tokenPrimary          = 1
	logonWithProfile      = 0x00000001
)

// selfExplorerPID 找当前会话的 explorer.exe。
func selfExplorerPID() uint32 {
	for _, p := range selfProcessList() {
		if strings.EqualFold(p.name, "explorer.exe") {
			return p.pid
		}
	}
	return 0
}

// selfUserToken 复制 explorer.exe 的主令牌，返回可用于 CreateProcessWithTokenW 的句柄。
func selfUserToken() (syscall.Handle, error) {
	pid := selfExplorerPID()
	if pid == 0 {
		return 0, fmt.Errorf("找不到 explorer.exe")
	}
	k32 := syscall.NewLazyDLL("kernel32.dll")
	adv := syscall.NewLazyDLL("advapi32.dll")
	openProcess := k32.NewProc("OpenProcess")
	openToken := adv.NewProc("OpenProcessToken")
	dupToken := adv.NewProc("DuplicateTokenEx")
	closeHandle := k32.NewProc("CloseHandle")

	ph, _, err := openProcess.Call(0x1000, 0, uintptr(pid)) // PROCESS_QUERY_LIMITED_INFORMATION
	if ph == 0 {
		return 0, fmt.Errorf("OpenProcess(explorer) 失败: %v", err)
	}
	defer closeHandle.Call(ph)

	var tok syscall.Handle
	if r, _, e := openToken.Call(ph, tokenDuplicate|tokenQuery|tokenAssignPrimary|tokenAdjustDefault|tokenAdjustSession,
		uintptr(unsafe.Pointer(&tok))); r == 0 {
		return 0, fmt.Errorf("OpenProcessToken 失败: %v", e)
	}
	defer closeHandle.Call(uintptr(tok))

	var dup syscall.Handle
	if r, _, e := dupToken.Call(uintptr(tok), tokenAllAccess, 0,
		securityImpersonation, tokenPrimary, uintptr(unsafe.Pointer(&dup))); r == 0 {
		return 0, fmt.Errorf("DuplicateTokenEx 失败: %v", e)
	}
	return dup, nil
}

// selfStartGameAsUser 用 explorer 的普通令牌启动游戏；失败时回退普通 CreateProcess。
func selfStartGameAsUser(exe, workDir, args string) (*selfGameProc, error) {
	tok, err := selfUserToken()
	if err != nil {
		selfLog("降权令牌不可用（%v），回退普通 CreateProcess", err)
		return selfStartGame(exe, workDir, args)
	}
	defer func() {
		k32 := syscall.NewLazyDLL("kernel32.dll")
		k32.NewProc("CloseHandle").Call(uintptr(tok))
	}()

	app, _ := syscall.UTF16PtrFromString(exe)
	cmdLine, _ := syscall.UTF16PtrFromString(fmt.Sprintf(`"%s" %s`, exe, args))
	dir, _ := syscall.UTF16PtrFromString(workDir)
	var si syscall.StartupInfo
	si.Cb = uint32(unsafe.Sizeof(si))
	var pi syscall.ProcessInformation

	adv := syscall.NewLazyDLL("advapi32.dll")
	createWithToken := adv.NewProc("CreateProcessWithTokenW")
	r, _, e := createWithToken.Call(
		uintptr(tok),
		logonWithProfile,
		uintptr(unsafe.Pointer(app)),
		uintptr(unsafe.Pointer(cmdLine)),
		0x400, // CREATE_UNICODE_ENVIRONMENT
		0,
		uintptr(unsafe.Pointer(dir)),
		uintptr(unsafe.Pointer(&si)),
		uintptr(unsafe.Pointer(&pi)),
	)
	if r == 0 {
		selfLog("CreateProcessWithTokenW 失败（%v），回退普通 CreateProcess", e)
		return selfStartGame(exe, workDir, args)
	}
	selfLog("已降权启动游戏 pid=%d（explorer 令牌）", pi.ProcessId)
	return &selfGameProc{Process: pi.Process, Thread: pi.Thread, ProcessID: pi.ProcessId}, nil
}
