//go:build windows

package main

// 主程序只允许由外壳拉起，双击 zuhaoclient.exe 时给出提示并退出。
// 新外壳通过环境变量标记；旧 1.2.56 外壳没有标记，但它是父进程且会等待主程序退出，
// 所以再按父进程文件名兜底（CG.exe / CG租号.exe）。

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

const viaLauncherEnv = "CG_VIA_LAUNCHER"

func launchedByShell() bool {
	if os.Getenv(viaLauncherEnv) == "1" {
		return true
	}
	parent, err := parentProcessImage()
	if err != nil || parent == "" {
		return false
	}
	switch strings.ToLower(filepath.Base(parent)) {
	case "cg租号.exe", "cg.exe":
		return true
	}
	return false
}

// requireShell 返回 false 表示不是外壳拉起的，调用方应直接退出
func requireShell() bool {
	if launchedByShell() {
		return true
	}
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	title, _ := syscall.UTF16PtrFromString("游戏租号客户端")
	body, _ := syscall.UTF16PtrFromString("请通过 CG租号.exe 启动游戏租号客户端。")
	const mbOK, mbIconInformation = 0x0, 0x40
	_, _, _ = proc.Call(0,
		uintptr(unsafe.Pointer(body)),
		uintptr(unsafe.Pointer(title)),
		uintptr(mbOK|mbIconInformation))
	return false
}
