//go:build windows

package main

import (
	"os/exec"
	"strings"
	"syscall"
)

// scheduleDelayedDelete 用一个分离的 cmd 等本进程退出后再删文件（用于 CG.exe.old 这类自身镜像）。
func scheduleDelayedDelete(paths []string) {
	if len(paths) == 0 {
		return
	}
	var sb strings.Builder
	sb.WriteString("ping -n 4 127.0.0.1 >nul")
	for _, p := range paths {
		sb.WriteString(` & del /f /q "`)
		sb.WriteString(p)
		sb.WriteString(`"`)
	}
	c := exec.Command("cmd")
	c.SysProcAttr = &syscall.SysProcAttr{
		CmdLine:       "cmd /c " + sb.String(),
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	_ = c.Start()
}
