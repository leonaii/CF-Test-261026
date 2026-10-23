package main

// Rail host 窗口。
//
// 证据（CF上号器_v1.17.exe @ 0x140009C00）：
//
//	7242  GetModuleHandleW(NULL)
//	7246  lea rax, "rail_dev_wnd"          ← 窗口类名
//	7256  WndProc = 0x140009BC0
//	7265  RegisterClassW(&wc)
//	7270  r8 = "WeGame"                    ← 窗口标题
//	7271  rdx = "rail_dev_wnd"
//	7286  CreateWindowExW(0,"rail_dev_wnd","WeGame",...)
//	7295  [0x1401741F8] = hwnd             ← 句柄存全局
//	7298  GetMessageW/DispatchMessageW     ← 消息循环
//
// 随后 0x14001BAxx 创建共享内存时：
//
//	22519 mov rdi,[0x1401741F8]            ← 用 HWND，不是 PID
//	22526 swprintf(buf,80,"Rail.SharedMemory.%u", rdi)
//
// rail_api64.dll 会用 FindWindowW 找这两个窗口，并打开
// "Rail.SharedMemory.<HWND>"，缺一样就弹 2310。

import (
	"sync"
	"syscall"
	"unsafe"
)

const (
	selfRailWndClass = "rail_dev_wnd"
	selfRailWndTitle = "WeGame"
	selfTgpWndClass  = "_TGP_MAIN_WINDOW_WND_NAME_"
)

var (
	selfRailWndOnce sync.Once
	selfRailWndHWND uintptr
	selfRailWndProc uintptr
)

// selfEnsureRailWindow 创建 Rail host 窗口，返回 HWND。只创建一次。
func selfEnsureRailWindow() uintptr {
	selfRailWndOnce.Do(func() {
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
		*(*uintptr)(unsafe.Pointer(&wc[0])) = 0
		*(*uintptr)(unsafe.Pointer(&wc[8])) = selfRailWndProc
		*(*uintptr)(unsafe.Pointer(&wc[24])) = hInst

		cls, _ := syscall.UTF16PtrFromString(selfRailWndClass)
		*(*uintptr)(unsafe.Pointer(&wc[64])) = uintptr(unsafe.Pointer(cls))

		registerClass.Call(uintptr(unsafe.Pointer(&wc[0])))

		title, _ := syscall.UTF16PtrFromString(selfRailWndTitle)
		// CreateWindowExW(0, cls, title, WS_OVERLAPPED=0, 0,0,0,0, 0,0,0,0)
		hwnd, _, _ := createWindowEx.Call(
			0,
			uintptr(unsafe.Pointer(cls)),
			uintptr(unsafe.Pointer(title)),
			0,
			0, 0, 0, 0,
			0, 0, hInst, 0,
		)
		selfRailWndHWND = hwnd

		// 第二个窗口：_TGP_MAIN_WINDOW_WND_NAME_（同样注册+创建，句柄不保存）
		if hwnd != 0 {
			wc2 := make([]byte, 72)
			*(*uintptr)(unsafe.Pointer(&wc2[8])) = selfRailWndProc
			*(*uintptr)(unsafe.Pointer(&wc2[24])) = hInst
			cls2, _ := syscall.UTF16PtrFromString(selfTgpWndClass)
			*(*uintptr)(unsafe.Pointer(&wc2[64])) = uintptr(unsafe.Pointer(cls2))
			registerClass.Call(uintptr(unsafe.Pointer(&wc2[0])))
			createWindowEx.Call(0, uintptr(unsafe.Pointer(cls2)), uintptr(unsafe.Pointer(cls2)),
				0, 0, 0, 0, 0, 0, 0, hInst, 0)
		}

		// 消息循环：DefWindowProc 需要有人泵消息，否则窗口不响应 SendMessage
		go func() {
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
