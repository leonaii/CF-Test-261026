//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// killProcessesUsing 结束可执行文件位于 dir 目录下的进程（不动自己），返回被结束的进程。
// 用于切换上号器时删除旧目录：先把占着文件的程序关掉。
func killProcessesUsing(dir string) []string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil
	}
	prefix := strings.ToLower(strings.TrimRight(abs, `\`) + `\`)

	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil
	}
	defer windows.CloseHandle(snap)

	self := uint32(os.Getpid())
	var killed []string

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	for err = windows.Process32First(snap, &entry); err == nil; err = windows.Process32Next(snap, &entry) {
		pid := entry.ProcessID
		if pid == 0 || pid == self {
			continue
		}
		name := windows.UTF16ToString(entry.ExeFile[:])

		h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION|windows.PROCESS_TERMINATE, false, pid)
		if err != nil {
			continue
		}
		buf := make([]uint16, windows.MAX_PATH)
		size := uint32(len(buf))
		if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err == nil {
			path := strings.ToLower(windows.UTF16ToString(buf[:size]))
			if strings.HasPrefix(path, prefix) {
				if windows.TerminateProcess(h, 1) == nil {
					killed = append(killed, fmt.Sprintf("%s(%d)", name, pid))
				}
			}
		}
		windows.CloseHandle(h)
	}
	return killed
}
