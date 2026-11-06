//go:build !windows

package main

import "fmt"

func tpsWGWeGameRunning() bool { return false }

func tpsWGMapAndWrite(payload, hookPath string) (func(), error) {
	return nil, fmt.Errorf("枪神纪 WEG 上号仅支持 Windows")
}

func tpsWGRunInjector(injector, dir string) (string, error) {
	return "", fmt.Errorf("枪神纪 WEG 上号仅支持 Windows")
}

func tpsWGReadStatus() (string, string) { return "", "" }

func tpsWGForceCleanup(logf func(string, ...any)) {
	logf("非 Windows 平台无需清理")
}

func tpsWGStartGame(gameDir string, logf func(string, ...any)) error {
	return fmt.Errorf("枪神纪 WEG 上号仅支持 Windows")
}
