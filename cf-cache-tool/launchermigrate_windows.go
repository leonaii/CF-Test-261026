//go:build windows

package main

// 把 1.2.56 那代没有自更新的 CG.exe 外壳换成对外分发的 CG租号.exe。
//
// 旧外壳只做两件事：下载主程序、等主程序退出，没有任何自更新代码，所以它永远停在
// CG.exe 这个名字上。主程序是被它拉起来的子进程，能拿到它的路径，于是启动时顺手把
// 新外壳写到同一目录，并把旧文件改名（运行中的 exe 只能改名，删不掉），下次启动就
// 只剩 CG租号.exe。

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	legacyLauncherName = "CG.exe"
	launcherDistName   = "CG租号.exe"
	legacyMinSize      = 1 << 20 // 新外壳约 200KB，旧 Go 外壳十几 MB
)

var legacyLauncherOnce sync.Once

func migrateLegacyLauncher(logf func(string, ...any)) {
	legacyLauncherOnce.Do(func() {
		target, err := replaceLegacyLauncher()
		switch {
		case err != nil:
			if logf != nil {
				logf("旧外壳迁移跳过: %v", err)
			}
		case target != "":
			if logf != nil {
				logf("已把旧外壳换成 %s", filepath.Base(target))
			}
		}
	})
}

// 返回新外壳路径；不需要迁移时返回空串
func replaceLegacyLauncher() (string, error) {
	parent, err := parentProcessImage()
	if err != nil || parent == "" {
		return "", nil
	}
	if !strings.EqualFold(filepath.Base(parent), legacyLauncherName) {
		return "", nil
	}
	st, err := os.Stat(parent)
	if err != nil || st.Size() < legacyMinSize {
		return "", nil
	}

	dir := filepath.Dir(parent)
	target := filepath.Join(dir, launcherDistName)
	if fi, err := os.Stat(target); err == nil && fi.Size() > 0 {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	m, err := fetchUpdateManifest(ctx)
	if err != nil {
		return "", err
	}
	name := strings.TrimSpace(m.File)
	if name == "" {
		name = legacyLauncherName
	}
	data, err := newCloudClient().download(ctx, updateDir+"/"+strings.TrimPrefix(name, "/"))
	if err != nil {
		return "", err
	}
	if int64(len(data)) != m.Size {
		return "", fmt.Errorf("外壳大小不符: %d != %d", len(data), m.Size)
	}
	sum := sha256.Sum256(data)
	if hex.EncodeToString(sum[:]) != strings.ToLower(strings.TrimSpace(m.SHA256)) {
		return "", errors.New("外壳 SHA256 校验失败")
	}

	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, data, 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	// 旧外壳正在运行，只能改名；残留的 .old 由新外壳下次启动清理
	if err := os.Rename(parent, parent+".old"); err == nil {
		scheduleLegacyDelete(parent + ".old")
	}
	return target, nil
}

// parentProcessImage 取父进程的可执行文件路径（旧外壳 → 主程序 是直接父子关系）
func parentProcessImage() (string, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(snap)

	self := uint32(os.Getpid())
	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	var parent uint32
	for err = windows.Process32First(snap, &entry); err == nil; err = windows.Process32Next(snap, &entry) {
		if entry.ProcessID == self {
			parent = entry.ParentProcessID
			break
		}
	}
	if parent == 0 {
		return "", errors.New("找不到父进程")
	}

	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, parent)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(h)

	buf := make([]uint16, 32768)
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return "", err
	}
	return windows.UTF16ToString(buf[:size]), nil
}

// scheduleLegacyDelete 循环尝试删除（父进程要等主程序退出才结束，短时间内删不掉）
func scheduleLegacyDelete(path string) {
	c := exec.Command("cmd")
	c.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: fmt.Sprintf(
			`cmd /c for /l %%i in (1,1,30) do @(del /f /q "%s" >nul 2>&1 & ping -n 3 127.0.0.1 >nul)`,
			path,
		),
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	_ = c.Start()
}
