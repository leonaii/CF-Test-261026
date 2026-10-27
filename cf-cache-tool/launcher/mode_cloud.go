//go:build !cg_local

package main

// 正式启动器只从云端只读分享获取主程序。
func localMode() bool { return false }
