//go:build !windows

package main

import "fmt"

func tclsUserInfoPath() (string, error) {
	return "", fmt.Errorf("仅支持 Windows")
}

func tclsStartSharedMemProbe(dir string, logf func(string, ...any)) func() {
	logf("非 Windows 平台跳过 TCLS 共享内存观测")
	return func() {}
}

func tclsStartGame(gameDir, qq string, logf func(string, ...any)) (*selfGameProc, error) {
	return nil, fmt.Errorf("仅支持 Windows")
}
