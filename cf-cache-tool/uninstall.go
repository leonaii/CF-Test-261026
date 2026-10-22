package main

// 切换上号器时清理旧目录里的本地数据（默认按用户选择执行）。

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// uninstallLocal 删除工作目录里由本工具生成的文件。
// withCache 为真时连账号缓存一起删；为假时只删上号器与配置。
func uninstallLocal(dir string, withCache bool) ([]string, error) {
	if dir == "" {
		return nil, errors.New("目录为空")
	}
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return nil, nil
	}

	var removed []string
	removeFile := func(name string) {
		p := filepath.Join(dir, name)
		if st, err := os.Stat(p); err != nil || st.IsDir() {
			return
		}
		if os.Remove(p) == nil {
			removed = append(removed, name)
		}
	}
	removeDir := func(name string) {
		p := filepath.Join(dir, name)
		if st, err := os.Stat(p); err != nil || !st.IsDir() {
			return
		}
		if os.RemoveAll(p) == nil {
			removed = append(removed, name+"/")
		}
	}

	// 所有类型的上号器 exe 都清掉，避免旧文件残留占空间
	for _, prof := range loaderProfiles {
		if prof.ExeName != "" {
			removeFile(prof.ExeName)
		}
	}
	for _, name := range []string{"host.dat", "host.dat.orig", "config.ini", "config.ini.orig"} {
		removeFile(name)
	}
	if withCache {
		removeDir("cache")
		removeDir(stagingDirName)
	}

	// 目录已经空了就一并删掉（只删空目录，不动用户其它文件）
	if entries, err := os.ReadDir(dir); err == nil && len(entries) == 0 {
		if os.Remove(dir) == nil {
			removed = append(removed, "(空目录)")
		}
	}
	return removed, nil
}

// removeDirForce 删除目录；被占用时先结束目录里的进程再重试。
func removeDirForce(dir string) ([]string, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, nil
	}
	if err := os.RemoveAll(dir); err == nil {
		return nil, nil
	}

	killed := killProcessesUsing(dir)
	var err error
	for i := 0; i < 12; i++ {
		time.Sleep(150 * time.Millisecond)
		if err = os.RemoveAll(dir); err == nil {
			return killed, nil
		}
	}
	return killed, fmt.Errorf("目录仍被占用，删除失败：%s\n%v", dir, err)
}
