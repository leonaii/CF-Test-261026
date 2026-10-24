//go:build windows

package main

// 单实例：用一个命名互斥量判断是否已经有 CG 在跑。
// 更新/切换卡密会重启进程，所以检测到已存在时先短重试，给旧进程退出留时间。

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/lxn/walk"
)

const instanceMutexName = "CG_Rent_SingleInstance_v1"

var (
	modKernel32      = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutexW = modKernel32.NewProc("CreateMutexW")
	procCloseHandle  = modKernel32.NewProc("CloseHandle")
)

// createMutex 返回句柄和 GetLastError；ERROR_ALREADY_EXISTS 表示已有实例。
func createMutex() (uintptr, uint32) {
	p, err := syscall.UTF16PtrFromString(instanceMutexName)
	if err != nil {
		return 0, 0
	}
	h, _, e := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(p)))
	if h == 0 {
		return 0, 0
	}
	code := uint32(0)
	if errno, ok := e.(syscall.Errno); ok {
		code = uint32(errno)
	}
	return h, code
}

// ensureSingleInstance 已有实例在跑时返回 false（并重试几次避开重启竞态）。
func ensureSingleInstance() bool {
	h, code := createMutex()
	if code != uint32(syscall.ERROR_ALREADY_EXISTS) {
		return true // 拿到了（句柄保持不关，进程退出自动释放）
	}
	procCloseHandle.Call(h)

	for i := 0; i < 4; i++ {
		time.Sleep(400 * time.Millisecond)
		h2, c2 := createMutex()
		if c2 != uint32(syscall.ERROR_ALREADY_EXISTS) {
			return true
		}
		procCloseHandle.Call(h2)
	}
	return false
}

// warnAlreadyRunning 第二个实例的提示。
func warnAlreadyRunning() {
	walk.MsgBox(nil, "CG租号",
		"软件已经在运行了，不要重复打开。\n\n请切到已打开的窗口。",
		walk.MsgBoxIconInformation)
}
