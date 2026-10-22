//go:build !windows

package main

// restartIfNeeded 只有 Windows 版需要（更新后重启）。
func restartIfNeeded() {}
