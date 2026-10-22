package main

// 启动过程留痕：界面启动失败或“点完没反应”时，可直接看这个文件定位。
// 位置：%AppData%\CG\startup.log，每次启动覆盖重写。

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type startupTrace struct {
	f *os.File
}

func openStartupTrace() *startupTrace {
	p := filepath.Join(configDirPath(), "startup.log")
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return &startupTrace{}
	}
	return &startupTrace{f: f}
}

func (t *startupTrace) printf(format string, args ...any) {
	if t == nil || t.f == nil {
		return
	}
	fmt.Fprintf(t.f, "%s  %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
	_ = t.f.Sync()
}

func (t *startupTrace) close() {
	if t != nil && t.f != nil {
		_ = t.f.Close()
	}
}

func (t *startupTrace) path() string {
	if t == nil || t.f == nil {
		return ""
	}
	return t.f.Name()
}
