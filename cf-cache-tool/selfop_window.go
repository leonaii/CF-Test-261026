package main

// Rail host 窗口。
//
// IDA 证据（CF上号器_v1.17.exe sub_140009C00）：
//
//	类数组 = { "WE_GAME_SDK_WINDOW_NAME_CLASS", "rail_dev_wnd" }
//	CreateWindowExW(0, 类名, "W", 0x80000000, 0, 0, 1, 1, ...)
//	qword_1401741F8 = 第一个窗口（WE_GAME_SDK）的 HWND
//
// 随后 sub_14001AEE0 创建共享内存时：
//
//	mov rdi,[0x1401741F8]                  ← 用 HWND，不是 PID
//	swprintf(buf,80,"Rail.SharedMemory.%u", rdi)
//
// rail_api64.dll 会找 WE_GAME_SDK_WINDOW_NAME_CLASS，并打开
// "Rail.SharedMemory.<该 HWND>"；句柄或类名错一个就是 2310。

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

const (
	selfRailWndClass = "rail_dev_wnd"
	selfRailWndTitle = "W"
	selfTgpWndClass  = "WE_GAME_SDK_WINDOW_NAME_CLASS"
)

var (
	selfRailWndOnce sync.Once
	selfRailWndHWND uintptr
	selfRailWndProc uintptr
)

// selfEnsureRailWindow 创建 Rail host 窗口，返回 HWND。只创建一次。
func selfEnsureRailWindow() uintptr {
	selfRailWndOnce.Do(func() {
		ready := make(chan uintptr, 1)
		// Windows 窗口消息队列绑定创建线程，必须在同一个 LockOSThread 线程里
		// CreateWindowExW + GetMessageW；否则游戏 SendMessage 会超时，Rail 2310。
		go func() {
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()

			user32 := syscall.NewLazyDLL("user32.dll")
			k32 := syscall.NewLazyDLL("kernel32.dll")

			registerClass := user32.NewProc("RegisterClassW")
			createWindowEx := user32.NewProc("CreateWindowExW")
			defWndProc := user32.NewProc("DefWindowProcW")
			getMessage := user32.NewProc("GetMessageW")
			translate := user32.NewProc("TranslateMessage")
			dispatch := user32.NewProc("DispatchMessageW")
			getModule := k32.NewProc("GetModuleHandleW")

			hInst, _, _ := getModule.Call(0)
			selfRailWndProc = defWndProc.Addr()

			// WNDCLASSW: 0 style, 8 lpfnWndProc, 16 cbClsExtra, 20 cbWndExtra,
			// 24 hInstance, 32 hIcon, 40 hCursor, 48 hbrBackground,
			// 56 lpszMenuName, 64 lpszClassName
			wc := make([]byte, 72)
			*(*uintptr)(unsafe.Pointer(&wc[8])) = selfRailWndProc
			*(*uintptr)(unsafe.Pointer(&wc[24])) = hInst

			// 样本顺序：先 WE_GAME_SDK_WINDOW_NAME_CLASS，再 rail_dev_wnd。
			// qword_1401741F8 保存第一个窗口句柄，共享内存名用它。
			title, _ := syscall.UTF16PtrFromString(selfRailWndTitle) // "W"
			tgpCls, _ := syscall.UTF16PtrFromString(selfTgpWndClass)
			railCls, _ := syscall.UTF16PtrFromString(selfRailWndClass)

			*(*uintptr)(unsafe.Pointer(&wc[64])) = uintptr(unsafe.Pointer(tgpCls))
			registerClass.Call(uintptr(unsafe.Pointer(&wc[0])))
			tgp, _, _ := createWindowEx.Call(
				0,
				uintptr(unsafe.Pointer(tgpCls)),
				uintptr(unsafe.Pointer(title)),
				0x80000000, // 样本 CreateWindowExW 的 style
				0, 0, 1, 1,
				0, 0, hInst, 0,
			)
			if tgp == 0 {
				ready <- 0
				return
			}

			wc2 := make([]byte, 72)
			*(*uintptr)(unsafe.Pointer(&wc2[8])) = selfRailWndProc
			*(*uintptr)(unsafe.Pointer(&wc2[24])) = hInst
			*(*uintptr)(unsafe.Pointer(&wc2[64])) = uintptr(unsafe.Pointer(railCls))
			registerClass.Call(uintptr(unsafe.Pointer(&wc2[0])))
			createWindowEx.Call(
				0,
				uintptr(unsafe.Pointer(railCls)),
				uintptr(unsafe.Pointer(title)),
				0x80000000,
				0, 0, 1, 1,
				0, 0, hInst, 0,
			)

			selfRailWndHWND = tgp
			ready <- tgp

			// 消息循环必须留在创建窗口的线程上。
			msg := make([]byte, 64)
			for {
				r, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg[0])), 0, 0, 0)
				if int32(r) <= 0 {
					return
				}
				translate.Call(uintptr(unsafe.Pointer(&msg[0])))
				dispatch.Call(uintptr(unsafe.Pointer(&msg[0])))
			}
		}()
		selfRailWndHWND = <-ready
		if selfRailWndHWND == 0 {
			selfLog("创建 WE_GAME_SDK 窗口失败")
		}
	})
	return selfRailWndHWND
}

// selfSharedMemName 返回 Rail 共享内存名（用 HWND，与样本一致）。
func selfSharedMemName() uint32 {
	if h := selfEnsureRailWindow(); h != 0 {
		return uint32(h)
	}
	return uint32(syscall.Getpid())
}
