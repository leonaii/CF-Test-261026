//go:build cg_local

package main

// 测试启动器读取 exe 同目录 app/，用于未发布版本的本地验证。
func localMode() bool { return true }
