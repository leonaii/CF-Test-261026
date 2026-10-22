//go:build !windows

package main

// 非 Windows 平台没有界面向导，仅供命令行版本编译通过。

type setupResult struct {
	ProfileID string
	WorkDir   string
}

func runSetupDialog(owner any, cfg *appConfig) (*setupResult, setupAction) {
	return nil, setupQuit
}
