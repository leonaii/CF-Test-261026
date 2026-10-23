package main

// 内嵌 Hook DLL 与 Rail 管道。
//
// 管道名（.rdata 原文，UTF-16）：
//
//	rail_named_pipe_host_%u      0x1400DCA72 / 0x1400DCE3A
//	rail_named_pipe_guest_%u     0x1400DCABA / 0x1400DCE82
//	rail_pipe_15                 0x1400DCDD2
//	rail_named_pipe_15           0x1400DCE02
//	rail_named_pipe_host_15      0x1400DCE40
//	rail_named_pipe_guest_15     0x1400DCE82
//
// 带 %u 的那一组才是运行时使用的：%u 是各自进程的 PID。
// guest（游戏）是服务端 → 创建 rail_named_pipe_guest_<gamePid>；
// host（WeGame / 上号器 / 我们）是客户端 → CreateFileW 去连它。
//
// 帧格式（0x140018D20 处解析）：
//
//	[uint32 len][uint32 cmd][payload:len]
//	响应回写同样 cmd，payload 首个 uint32 为状态码。

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

//go:embed qwh_cf_hook.dll
var selfHookDLL []byte

// selfConnectGuestPipe 作为 host 连上游戏的 guest 管道，并在有数据时回帧。
// 游戏启动后才会创建管道，所以要重试。
func selfConnectGuestPipe(gamePID int) {
	candidates := []string{
		fmt.Sprintf(`\\.\pipe\rail_named_pipe_guest_%d`, gamePID),
		fmt.Sprintf(`\\.\pipe\rail_named_pipe_guest_15%08d`, gamePID),
		fmt.Sprintf(`\\.\pipe\rail_named_pipe_%d`, gamePID),
		fmt.Sprintf(`\\.\pipe\rail_named_pipe_15%08d`, gamePID),
		fmt.Sprintf(`\\.\pipe\rail_pipe_%d`, gamePID),
		fmt.Sprintf(`\\.\pipe\rail_pipe_15%08d`, gamePID),
	}
	name := candidates[0]
	go func() {
		k32 := syscall.NewLazyDLL("kernel32.dll")
		createFile := k32.NewProc("CreateFileW")
		readFile := k32.NewProc("ReadFile")
		writeFile := k32.NewProc("WriteFile")
		closeHandle := k32.NewProc("CloseHandle")

		// 等到连上为止（最多 60 秒），逐个候选名尝试
		var h syscall.Handle
		deadline := time.Now().Add(60 * time.Second)
		for time.Now().Before(deadline) && h == 0 {
			for _, cand := range candidates {
				p, _ := syscall.UTF16PtrFromString(cand)
				v, _, _ := createFile.Call(
					uintptr(unsafe.Pointer(p)),
					uintptr(0xC0000000), // GENERIC_READ|GENERIC_WRITE
					0, 0,
					uintptr(3), // OPEN_EXISTING
					0, 0,
				)
				if v != 0 && syscall.Handle(v) != syscall.InvalidHandle {
					h = syscall.Handle(v)
					name = cand
					break
				}
			}
			if h == 0 {
				time.Sleep(200 * time.Millisecond)
			}
		}
		if h == 0 {
			fmt.Printf("rail: 未能连上任一 guest 管道\n")
			return
		}
		defer closeHandle.Call(uintptr(h))
		fmt.Printf("rail: 已连上 %s\n", name)

		buf := make([]byte, 0x1000)
		for {
			var n uint32
			ok, _, _ := readFile.Call(
				uintptr(h),
				uintptr(unsafe.Pointer(&buf[0])),
				uintptr(len(buf)),
				uintptr(unsafe.Pointer(&n)),
				0,
			)
			if ok == 0 || n < 8 {
				return
			}
			ln := binary.LittleEndian.Uint32(buf[0:4])
			cmd := binary.LittleEndian.Uint32(buf[4:8])
			if ln > 0x400000 || n < 8+ln {
				return
			}
			payload := buf[8 : 8+ln]

			// 把游戏发来的原始帧落盘，用于比对协议
			selfDumpFrame("guest", cmd, payload)

			resp := make([]byte, 12+len(payload))
			binary.LittleEndian.PutUint32(resp[0:4], uint32(4+len(payload)))
			binary.LittleEndian.PutUint32(resp[4:8], cmd)
			binary.LittleEndian.PutUint32(resp[8:12], 0)
			copy(resp[12:], payload)

			var w uint32
			writeFile.Call(
				uintptr(h),
				uintptr(unsafe.Pointer(&resp[0])),
				uintptr(len(resp)),
				uintptr(unsafe.Pointer(&w)),
				0,
			)
		}
	}()
}

// selfDumpFrame 把管道帧追加写到 <exe目录>\rail_frames.log，用于还原协议。
func selfDumpFrame(dir string, cmd uint32, payload []byte) {
	f, err := os.OpenFile(exeDir()+`\rail_frames.log`, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	head := fmt.Sprintf("[%s] cmd=%d len=%d payload=", dir, cmd, len(payload))
	fmt.Fprintf(f, "%s%x\n", head, payload)
}

// selfServeHostPipe 同时也建一条 host 管道（部分 Rail 版本由 host 当服务端）。
// 名字必须带自己的 PID。
func selfServeHostPipe() func() {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	createPipe := k32.NewProc("CreateNamedPipeW")
	connectPipe := k32.NewProc("ConnectNamedPipe")
	disconnPipe := k32.NewProc("DisconnectNamedPipe")
	readFile := k32.NewProc("ReadFile")
	writeFile := k32.NewProc("WriteFile")
	closeHandle := k32.NewProc("CloseHandle")

	var (
		mu      sync.Mutex
		handles []syscall.Handle
	)
	stop := make(chan struct{}, 1)

	serve := func(name string) {
		namePtr, _ := syscall.UTF16PtrFromString(name)
		h, _, _ := createPipe.Call(
			uintptr(unsafe.Pointer(namePtr)),
			3, 0, uintptr(0x10000), uintptr(0x10000), 0, 0,
		)
		if h == 0 || syscall.Handle(h) == syscall.InvalidHandle {
			return
		}
		pipe := syscall.Handle(h)
		mu.Lock()
		handles = append(handles, pipe)
		mu.Unlock()

		go func() {
			for {
				select {
				case <-stop:
					return
				default:
				}
				connectPipe.Call(uintptr(pipe), 0)
				go func(p syscall.Handle) {
					buf := make([]byte, 0x1000)
					for {
						var n uint32
						ok, _, _ := readFile.Call(uintptr(p), uintptr(unsafe.Pointer(&buf[0])),
							uintptr(len(buf)), uintptr(unsafe.Pointer(&n)), 0)
						if ok == 0 || n < 8 {
							disconnPipe.Call(uintptr(p))
							return
						}
						ln := binary.LittleEndian.Uint32(buf[0:4])
						cmd := binary.LittleEndian.Uint32(buf[4:8])
						if ln > 0x400000 || n < 8+ln {
							disconnPipe.Call(uintptr(p))
							return
						}
						payload := buf[8 : 8+ln]
						resp := make([]byte, 12+len(payload))
						binary.LittleEndian.PutUint32(resp[0:4], uint32(4+len(payload)))
						binary.LittleEndian.PutUint32(resp[4:8], cmd)
						binary.LittleEndian.PutUint32(resp[8:12], 0)
						copy(resp[12:], payload)
						var w uint32
						writeFile.Call(uintptr(p), uintptr(unsafe.Pointer(&resp[0])),
							uintptr(len(resp)), uintptr(unsafe.Pointer(&w)), 0)
					}
				}(pipe)
			}
		}()
	}

	// 实测上号器会同时建出 12 个管道（4 个基名 × 3 种后缀）：
	//   rail_pipe_15 / rail_pipe_15%08d(pid) / rail_pipe_%d(pid)
	//   rail_named_pipe_* / rail_named_pipe_host_* / rail_named_pipe_guest_*
	// 少一个 Rail 就可能连不上 → CShell 报 2310。
	pid := os.Getpid()
	bases := []string{
		"rail_pipe_",
		"rail_named_pipe_",
		"rail_named_pipe_host_",
		"rail_named_pipe_guest_",
	}
	for _, b := range bases {
		serve(`\\.\pipe\` + b + "15")
		serve(fmt.Sprintf(`\\.\pipe\%s15%08d`, b, pid))
		serve(fmt.Sprintf(`\\.\pipe\%s%d`, b, pid))
	}

	return func() {
		select {
		case stop <- struct{}{}:
		default:
		}
		mu.Lock()
		for _, h := range handles {
			disconnPipe.Call(uintptr(h))
			closeHandle.Call(uintptr(h))
		}
		handles = nil
		mu.Unlock()
	}
}

// selfEnsureHookDLL 把内嵌的 qwh_cf_hook.dll 释放到 exe 同级目录。
//
// 样本里的相对路径是 "ntop_v17_dll\x64"（.rdata 0x1400DC8F8），
// 基准是上号器自身目录 —— 不是 %TEMP%。ACE 对 TEMP 下的 DLL 更敏感，
// 放到 exe 同级才与上号器一致。
func selfEnsureHookDLL() (string, error) {
	if len(selfHookDLL) == 0 {
		return "", fmt.Errorf("未内嵌 hook DLL")
	}
	base := exeDir()
	dir := base + `\ntop_v17_dll\x64`
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	p := dir + `\qwh_cf_hook.dll`
	if st, err := os.Stat(p); err == nil && st.Size() == int64(len(selfHookDLL)) {
		return p, nil
	}
	if err := os.WriteFile(p, selfHookDLL, 0o644); err != nil {
		return "", err
	}
	return p, nil
}

// exeDir 返回当前可执行文件所在目录。
func exeDir() string {
	p, err := os.Executable()
	if err != nil {
		return "."
	}
	i := len(p)
	for i > 0 {
		c := p[i-1]
		if c == '\\' || c == '/' {
			break
		}
		i--
	}
	if i > 0 {
		return p[:i-1]
	}
	return "."
}
