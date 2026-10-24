//go:build !windows

package main

// 非 Windows 平台不做单实例限制。
func ensureSingleInstance() bool { return true }

func warnAlreadyRunning() {}
