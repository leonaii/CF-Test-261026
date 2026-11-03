//go:build windows

package main

import (
	"os/exec"
	"strings"
	"syscall"
)

// wuweiWGWeGameState 只看本机官方 WeGame 有没有跑起来，不联网、不弹框。
// 隐藏子进程窗口，避免 GUI 程序里闪一下控制台。
func wuweiWGWeGameState() string {
	c := exec.Command("tasklist", "/FO", "CSV", "/NH", "/FI", "IMAGENAME eq wegame.exe")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	out, err := c.Output()
	if err != nil {
		return "无法检测官方 WeGame 进程"
	}
	if strings.Contains(strings.ToLower(string(out)), "wegame.exe") {
		return "官方 WeGame 已运行"
	}
	return "官方 WeGame 未运行（请先打开并登录，否则注入会失败）"
}
