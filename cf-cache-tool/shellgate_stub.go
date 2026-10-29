//go:build !windows

package main

// 非 Windows 只用于跨平台编译，不做启动来源校验。
func requireShell() bool { return true }
