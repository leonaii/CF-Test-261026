//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// tpsWGWeGameRunning 检查用户是否已经打开并登录官方 WeGame。
func tpsWGWeGameRunning() bool {
	want := map[string]bool{
		"wegame.exe":      true,
		"rail.exe":        true,
		"railservice.exe": true,
	}
	for _, p := range selfProcessList() {
		if want[strings.ToLower(p.name)] {
			return true
		}
	}
	return false
}

// tpsWGMapAndWrite 创建并写入 WEG 同款共享内存，返回释放函数。
// 布局（IDA sub_1403B1253 主 exe + wg_hook.dll sub_10004070 / sub_10002870）：
//
//	+0     DWORD 魔数 "VALW" (0x574C4156)
//	+8     DWORD game_id（WEG 写死 2001715；枪神纪 33）
//	+16    票据 payload（VAL={...} / CF={...} JSON 文本）
//	+1712  64B  hook DLL 回写状态：hook_install / hook_ok / patch_fail / dump_ok …
//	+1776  DWORD 启用开关，非 0 时 hook 才处理帧（sub_10004070 开头就检查它）
//	+1780  DWORD dump 序号计数器（hook 会 InterlockedIncrement）
//	+1784  DWORD gid 过滤器（0=不过滤；非 0 只处理该 game_id 的帧）
//	+1788  200B  hook DLL 回写 dump 文件 ANSI 路径（不是输入！）
//
// hook DLL 用 OpenFileMappingW(Local\ValWgInjectAcc_v1) + MapViewOfFile(0x29C4)
// 打开，映射必须在注入器与游戏启动期间保持打开。
func tpsWGMapAndWrite(payload, hookPath string) (func(), error) {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	createMap := k32.NewProc("CreateFileMappingW")
	mapView := k32.NewProc("MapViewOfFile")
	unmap := k32.NewProc("UnmapViewOfFile")

	name, err := syscall.UTF16PtrFromString(tpsWGMemName)
	if err != nil {
		return nil, err
	}
	const invalidHandle = ^uintptr(0)
	h, _, e := createMap.Call(invalidHandle, 0, 4 /*PAGE_READWRITE*/, 0,
		uintptr(tpsWGMemSize), uintptr(unsafe.Pointer(name)))
	if h == 0 {
		return nil, fmt.Errorf("CreateFileMappingW 失败: %v", e)
	}
	view, _, e := mapView.Call(h, 0x000F001F /*FILE_MAP_ALL_ACCESS*/, 0, 0, uintptr(tpsWGMemSize))
	if view == 0 {
		syscall.CloseHandle(syscall.Handle(h))
		return nil, fmt.Errorf("MapViewOfFile 失败: %v", e)
	}

	mem := unsafe.Slice((*byte)(unsafe.Pointer(view)), tpsWGMemSize)
	for i := range mem {
		mem[i] = 0
	}
	// 头部字段
	putU32(mem[0:4], tpsWGMagic)             // "VALW"
	putU32(mem[4:8], 1)                      // 替换开关：非 0 才调用 sub_100046B0 换票
	putU32(mem[8:12], uint32(tpsWGGameID))   // game_id = 33
	putU32(mem[16:20], uint32(len(payload))) // +16 QWORD 低位，替换函数要求非 0
	putU32(mem[1776:1780], 1)                // dump 开关：非 0 才调用 sub_10004070 录帧
	putU32(mem[1780:1784], 0)                // dump 序号，交给 hook 自增
	putU32(mem[1784:1788], 0)                // gid 过滤器：0=不过滤。填错会导致 hook 全部 return 0

	// 替换函数 sub_100046B0 / sub_100028F0 读的是两个 C 字符串槽：
	//   +24  主票据串（strlen 读取，放进 protobuf 帧）
	//   +152 次票据串（可选，非空才读）
	// 之前把整段 JSON 从 +16 开始铺过去，会覆盖这两个槽的边界，且 +16 是 DWORD 长度域。
	const str1Off, str2Off, stateOff = 24, 152, 1712
	p1 := payload
	if len(p1) > str2Off-str1Off-1 {
		p1 = p1[:str2Off-str1Off-1]
	}
	copy(mem[str1Off:], p1)
	// +152 放一份同样的票据串，确保该槽非空（替换函数要求 *(_BYTE*)(+152) != 0）
	p2 := payload
	if len(p2) > stateOff-str2Off-1 {
		p2 = p2[:stateOff-str2Off-1]
	}
	copy(mem[str2Off:], p2)

	// +1712 状态槽 / +1788 路径槽都由 hook DLL 回写，此处保持清零。
	_ = hookPath

	return func() {
		unmap.Call(view)
		syscall.CloseHandle(syscall.Handle(h))
	}, nil
}

func putU32(b []byte, v uint32) {
	b[0], b[1], b[2], b[3] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24)
}

// tpsWGReadStatus 重新打开共享内存，读 hook DLL 回写的状态与 dump 路径。
// +1712 是 ASCII 状态（hook_ok / patch_fail / dump_ok / dump_open_fail …），
// +1788 是 ANSI 的 cmd100 dump 文件路径。
func tpsWGReadStatus() (string, string) {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	openMap := k32.NewProc("OpenFileMappingW")
	mapView := k32.NewProc("MapViewOfFile")
	unmap := k32.NewProc("UnmapViewOfFile")

	name, err := syscall.UTF16PtrFromString(tpsWGMemName)
	if err != nil {
		return "", ""
	}
	h, _, _ := openMap.Call(0x000F001F, 0, uintptr(unsafe.Pointer(name)))
	if h == 0 {
		return "", ""
	}
	view, _, _ := mapView.Call(h, 0x000F001F, 0, 0, uintptr(tpsWGMemSize))
	if view == 0 {
		syscall.CloseHandle(syscall.Handle(h))
		return "", ""
	}
	mem := unsafe.Slice((*byte)(unsafe.Pointer(view)), tpsWGMemSize)
	readC := func(off, n int) string {
		b := mem[off : off+n]
		end := 0
		for end < len(b) && b[end] != 0 {
			end++
		}
		return string(b[:end])
	}
	status, path := readC(1712, 32), readC(1788, 200)
	unmap.Call(view)
	syscall.CloseHandle(syscall.Handle(h))
	return status, path
}

// tpsWGDisableSwitches 把共享内存里的「替换开关 +4」和「录帧开关 +1776」清 0，
// 让还挂在进程里的 hook 立刻停手（即使进程暂时没被杀掉）。
func tpsWGDisableSwitches() {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	openMap := k32.NewProc("OpenFileMappingW")
	mapView := k32.NewProc("MapViewOfFile")
	unmap := k32.NewProc("UnmapViewOfFile")

	name, err := syscall.UTF16PtrFromString(tpsWGMemName)
	if err != nil {
		return
	}
	h, _, _ := openMap.Call(0x000F001F, 0, uintptr(unsafe.Pointer(name)))
	if h == 0 {
		return
	}
	view, _, _ := mapView.Call(h, 0x000F001F, 0, 0, uintptr(tpsWGMemSize))
	if view == 0 {
		syscall.CloseHandle(syscall.Handle(h))
		return
	}
	mem := unsafe.Slice((*byte)(unsafe.Pointer(view)), tpsWGMemSize)
	putU32(mem[4:8], 0)
	putU32(mem[1776:1780], 0)
	unmap.Call(view)
	syscall.CloseHandle(syscall.Handle(h))
}

// tpsWGForceCleanup 关闭状态窗口时执行：清开关 -> 强制结束 WeGame / 枪神纪进程。
// wg_hook.dll 是通过 LoadLibraryW 注入到 wegame.exe / rail.exe 的，没有单独的
// FreeLibrary 通道；进程退出即等于卸载 hook，这也是最彻底、不会留下 inline patch 的做法。
func tpsWGForceCleanup(logf func(string, ...any)) {
	tpsWGDisableSwitches()

	// 官方 WeGame 全家 + 枪神纪 TCLS / 反作弊子进程。
	// 不含 zuhaoclient.exe / CG租号.exe，避免把自己杀掉。
	targets := map[string]bool{
		"wegame.exe":          true,
		"wegame_env.exe":      true,
		"rail.exe":            true,
		"railservice.exe":     true,
		"client_launcher.exe": true,
		"tps.exe":             true,
		"tpswebbrowser.exe":   true,
		"tpsmusic.exe":        true,
		"tpsreport.exe":       true,
		"client.exe":          true, // TCLS\Client.exe
		"crossproxy.exe":      true,
		"tp3helper.exe":       true,
		"sguard64.exe":        true,
		"sguardsvc64.exe":     true,
		"valwginject.exe":     true,
	}

	k32 := syscall.NewLazyDLL("kernel32.dll")
	openProc := k32.NewProc("OpenProcess")
	terminate := k32.NewProc("TerminateProcess")
	closeHandle := k32.NewProc("CloseHandle")

	n := 0
	for _, p := range selfProcessList() {
		if !targets[strings.ToLower(p.name)] {
			continue
		}
		h, _, err := openProc.Call(uintptr(p.pid), 0, 0x0001 /*PROCESS_TERMINATE*/)
		if h == 0 {
			continue
		}
		ok, _, _ := terminate.Call(h, 1)
		closeHandle.Call(h)
		if ok != 0 {
			n++
		}
		_ = err
	}
	logf("已强制结束 %d 个 WeGame / 枪神纪进程。", n)
}

// tpsWGRunInjector 运行 ValWgInject.exe，隐藏控制台窗口，并以 inject.log 判定成败。
// 故意不传命令行参数：IDA sub_401FC0 显示无参数时注入器会自动拼 <自身目录>\wg_hook.dll，
// 传目录会被当成 DLL 路径交给 LoadLibraryW，必然返回 NULL（any=0）。
func tpsWGRunInjector(injector, dir string) (string, error) {
	logPath := filepath.Join(dir, "inject.log")
	_ = os.Remove(logPath) // 清掉上一轮日志，避免把旧的 any=1 误判成本次成功

	cmd := exec.Command(injector)
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW，去掉黑色控制台框
	}
	out, err := cmd.CombinedOutput()

	// 注入器把详细日志写进 inject.log（UTF-8 追加），stdout 可能为空
	logText := ""
	if b, e := os.ReadFile(logPath); e == nil {
		logText = string(b)
	}
	if strings.TrimSpace(logText) == "" {
		logText = string(out)
	}

	if err != nil {
		return logText, fmt.Errorf("注入器执行失败: %w", err)
	}
	if !strings.Contains(logText, "inject done any=1") {
		if strings.Contains(logText, "err=5") {
			return logText, fmt.Errorf("注入失败：OpenProcess 被拒绝（权限不足，请以管理员身份运行本程序）")
		}
		return logText, fmt.Errorf("注入失败：inject done any=0，wg_hook.dll 未加载进 wegame/rail")
	}
	return logText, nil
}

// tpsWGStartGame 让官方 WeGame 拉起枪神纪（不能自己 CreateProcess TPS.exe，
// 否则 ACE/TP 安全组件会拦，已实测）。
func tpsWGStartGame(gameDir string, logf func(string, ...any)) error {
	wegame := `C:\WeGame\wegame.exe`
	client := `C:\WeGame\client_launcher.exe`

	var lastErr error
	for _, c := range []struct {
		exe  string
		args []string
	}{
		{client, []string{"-game_id=33"}},
		{wegame, []string{"-game_id=33"}},
	} {
		if _, err := os.Stat(c.exe); err != nil {
			lastErr = err
			continue
		}
		cmd := exec.Command(c.exe, c.args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
		if err := cmd.Start(); err != nil {
			lastErr = err
			continue
		}
		logf("已请求官方启动枪神纪: %s %v", c.exe, c.args)
		return nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("未找到 WeGame 启动器")
	}
	return lastErr
}
