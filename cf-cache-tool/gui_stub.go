//go:build !windows

package main

import "errors"

// runGUI 在非 Windows 平台不可用，主程序仍可使用命令行子命令。
func runGUI() error {
	return errors.New("图形界面仅支持 Windows；当前平台请使用命令行子命令")
}
