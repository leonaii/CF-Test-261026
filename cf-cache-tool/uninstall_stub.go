//go:build !windows

package main

// killProcessesUsing 只在 Windows 上实现（结束占用目录的进程）。
func killProcessesUsing(dir string) []string { return nil }
