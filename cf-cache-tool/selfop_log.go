package main

// 启动流程日志：GUI 版没有控制台，出问题只能靠文件。
// 位置：%TEMP%\CG_launch.log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var selfLogMu sync.Mutex

// selfLogPath 返回日志文件路径。
func selfLogPath() string {
	return filepath.Join(os.TempDir(), "CG_launch.log")
}

// selfLog 追加一行日志（带毫秒时间戳）。
func selfLog(format string, args ...any) {
	selfLogMu.Lock()
	defer selfLogMu.Unlock()
	f, err := os.OpenFile(selfLogPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = fmt.Fprintf(f, "%s  %s\n", time.Now().Format("15:04:05.000"), fmt.Sprintf(format, args...))
}
