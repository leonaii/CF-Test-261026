//go:build !windows

package main

// 非 Windows 平台没有官方 WeGame，也不做进程检测。
func wuweiWGWeGameState() string { return "非 Windows，跳过官方 WeGame 检测" }
