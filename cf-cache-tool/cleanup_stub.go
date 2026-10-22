//go:build !windows

package main

// scheduleDelayedDelete 只有 Windows 版需要（删除自身镜像的旧版本文件）。
func scheduleDelayedDelete(paths []string) {}
