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
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

//go:embed qwh_cf_hook.dll
var selfHookDLL []byte

// selfRailHello 是样本 sub_140018D10 在管道连上后先写出的 8817 字节
// Rail 初始包（样本日志 send100，last_cf_rail.bin 原样）。
//
//go:embed rail_hello.bin
var selfRailHello []byte

// 离线器 last_cf_pipe.txt 只有一次 send100；8 个管道都发首包会让
// 游戏端协议状态错乱，连接后一个 recv 都不回（仍 2310）。
var selfRailHelloOnce sync.Once

// 离线器 last_cf_pipe.txt 的顺序：
//
//	listen ... -> inject_hook cmd=1 -> listen ... -> send100 cmd=100
//
// 即注入由游戏命令触发，不是启动时主动注入。
var (
	selfHookInjectOnce sync.Once
	selfHookPath       string
	selfHookProcHandle syscall.Handle
	selfHookGamePID    int
)

func selfResetHookInjection() {
	selfHookInjectOnce = sync.Once{}
	selfHookPath = ""
	selfHookProcHandle = 0
	selfHookGamePID = 0
}

// selfInjectHookOnDemand 在游戏发来 cmd=1（inject_hook）时调用。
func selfInjectHookOnDemand() {
	selfHookInjectOnce.Do(func() {
		if selfHookProcHandle == 0 || selfHookPath == "" {
			selfLog("inject_hook cmd=1 但缺少进程句柄/dll")
			return
		}
		if err := selfInjectHookHandle(selfHookProcHandle, selfHookPath); err != nil {
			selfLog("inject_hook 失败 pid=%d: %v", selfHookGamePID, err)
			return
		}
		selfLog("inject_hook 成功 pid=%d dll=%s", selfHookGamePID, selfHookPath)
		syscall.CloseHandle(selfHookProcHandle)
		selfHookProcHandle = 0
	})
}

// selfLoadRailHello 使用离线器 sub_14001AEE0 每次启动重新生成的 last_cf_rail.bin。
// 固定旧模板会带旧账号/旧 Rail 会话，导致 Rail 初始化失败。
func selfLoadRailHello(qq string, args ...string) {
	selfRailHelloOnce = sync.Once{}
	tgpid := ""
	var ssn string
	if len(args) > 0 {
		tgpid = args[0]
	}
	if len(args) > 1 {
		ssn = args[1]
	}
	p := filepath.Join(os.TempDir(), "last_cf_rail.bin")
	// 优先用离线器已经为该账号生成好的 Rail 文件，绝不覆盖。
	// 旧版本会用别的区模板补出 8817 这类假文件；长度不等于本区模板就丢弃。
	if b, err := os.ReadFile(p); err == nil && len(b) >= 8 &&
		int(binary.LittleEndian.Uint32(b[:4]))+8 == len(b) && selfRailTemplateUin(b) == qq &&
		selfRailFileMatchesZone(ssn, len(b)) {
		b2, _ := selfPatchRailStr1(b, tgpid)
		selfRailHello = b2
		selfLog("复用离线 Rail 文件 qq=%s tgpid=%s ssn=%s len=%d", qq, tgpid, ssn, len(b2))
		return
	}
	if ssn != "" {
		if z := selfZoneRailTemplate(ssn); z != nil {
			if patched, err := selfPatchRailAccount(z, qq, tgpid); err == nil {
				selfRailHello = patched
				selfLog("使用区服 Rail 模板 SSN=%s qq=%s tgpid=%s len=%d", ssn, qq, tgpid, len(patched))
				_ = os.WriteFile(p, patched, 0o644)
				return
			} else {
				selfLog("区服模板 %s 打账号补丁失败: %v", ssn, err)
			}
		}
	}
	b, err := os.ReadFile(p)
	if err == nil && len(b) >= 8 && int(binary.LittleEndian.Uint32(b[:4]))+8 == len(b) {
		selfRailHello = b
	}
	if qq != "" && !bytes.Contains(selfRailHello, []byte(qq)) {
		oldUin := selfRailTemplateUin(selfRailHello)
		if oldUin != "" {
			if patched, err := selfPatchRailUin(selfRailHello, oldUin, qq); err == nil {
				selfRailHello = patched
				selfLog("Rail Uin 动态替换 %s -> %s len=%d", oldUin, qq, len(patched))
			} else {
				selfLog("Rail Uin 替换失败 %s -> %s: %v", oldUin, qq, err)
			}
		}
	}
	// 样本在 0x14001C60A 生成/写 last_cf_rail.bin，游戏侧也会读它。
	if len(selfRailHello) > 0 {
		if err := os.WriteFile(p, selfRailHello, 0o644); err == nil {
			selfLog("写入 last_cf_rail.bin len=%d qq=%s", len(selfRailHello), qq)
		}
	}
}

func selfRailTemplateUin(b []byte) string {
	marker := []byte(`"str2":"`)
	i := bytes.Index(b, marker)
	if i < 0 {
		return ""
	}
	start := i + len(marker)
	end := bytes.IndexByte(b[start:], '"')
	if end <= 0 {
		return ""
	}
	return string(b[start : start+end])
}

// selfPatchRailUin 按 sub_140014CB0 的 protobuf 结构替换 Uin，并修正每层长度。
func selfPatchRailUin(b []byte, oldUin, newUin string) ([]byte, error) {
	if len(b) < 8 || oldUin == "" || newUin == "" {
		return nil, fmt.Errorf("invalid rail template")
	}
	oldB, newB := []byte(oldUin), []byte(newUin)
	body, err := selfPatchProtoMessage(b[8:], oldB, newB)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, 8+len(body))
	out = append(out, b[:8]...)
	out = append(out, body...)
	binary.LittleEndian.PutUint32(out[:4], uint32(len(out)-8))
	return out, nil
}

func selfReadUvarint(b []byte) (v uint64, n int, ok bool) {
	for i := 0; i < len(b) && i < 10; i++ {
		c := b[i]
		if i == 9 && c > 1 {
			return 0, 0, false
		}
		v |= uint64(c&0x7f) << (7 * i)
		if c < 0x80 {
			return v, i + 1, true
		}
	}
	return 0, 0, false
}

func selfAppendUvarint(dst []byte, v uint64) []byte {
	for v >= 0x80 {
		dst = append(dst, byte(v)|0x80)
		v >>= 7
	}
	return append(dst, byte(v))
}

// selfPatchProtoMessage 只处理 protobuf wire type 0/1/2/5。
// 对长度字段，若内容本身是合法嵌套 message 就递归；否则按原始字符串替换。
func selfPatchProtoMessage(data, oldB, newB []byte) ([]byte, error) {
	var out []byte
	pos := 0
	for pos < len(data) {
		start := pos
		key, n, ok := selfReadUvarint(data[pos:])
		if !ok || key>>3 == 0 {
			return nil, fmt.Errorf("bad protobuf key at %d", pos)
		}
		pos += n
		keyEnd := pos
		wire := key & 7
		switch wire {
		case 0:
			_, n, ok = selfReadUvarint(data[pos:])
			if !ok {
				return nil, fmt.Errorf("bad varint at %d", pos)
			}
			pos += n
			out = append(out, data[start:pos]...)
		case 1:
			if pos+8 > len(data) {
				return nil, fmt.Errorf("bad fixed64 at %d", pos)
			}
			pos += 8
			out = append(out, data[start:pos]...)
		case 5:
			if pos+4 > len(data) {
				return nil, fmt.Errorf("bad fixed32 at %d", pos)
			}
			pos += 4
			out = append(out, data[start:pos]...)
		case 2:
			ln, n, ok := selfReadUvarint(data[pos:])
			if !ok {
				return nil, fmt.Errorf("bad length at %d", pos)
			}
			pos += n
			if ln > uint64(len(data)-pos) {
				return nil, fmt.Errorf("length overflow at %d", pos)
			}
			end := pos + int(ln)
			payload := data[pos:end]
			newPayload := payload
			if bytes.Contains(payload, oldB) {
				if nested, err := selfPatchProtoMessage(payload, oldB, newB); err == nil && !bytes.Equal(nested, payload) {
					newPayload = nested
				} else {
					newPayload = bytes.ReplaceAll(payload, oldB, newB)
				}
			}
			out = append(out, data[start:keyEnd]...)
			out = selfAppendUvarint(out, uint64(len(newPayload)))
			out = append(out, newPayload...)
			pos = end
		default:
			return nil, fmt.Errorf("unsupported wire type %d at %d", wire, pos)
		}
	}
	return out, nil
}

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

// selfRailPipeBases 是样本 sub_14001A7B0 的 4 个管道基名。
func selfRailPipeBases() []string {
	return []string{
		"rail_pipe_",
		"rail_named_pipe_",
		"rail_named_pipe_host_",
		"rail_named_pipe_guest_",
	}
}

// selfServeStaticRailPipes 对应样本启动前创建的 4 个 *_15 管道。
func selfServeStaticRailPipes() func() {
	var names []string
	for _, b := range selfRailPipeBases() {
		names = append(names, `\\.\pipe\`+b+"15")
	}
	return selfServeRailPipeSet(names)
}

// selfServeGameRailPipes 对应样本 sub_14001AB70(gamePid)：
// 4 个基名 × {gamePid, gamePid+1500000000}，全部由 host 建服务端。
func selfServeGameRailPipes(gamePID int) func() {
	var names []string
	for _, b := range selfRailPipeBases() {
		names = append(names,
			fmt.Sprintf(`\\.\pipe\%s%d`, b, gamePID),
			fmt.Sprintf(`\\.\pipe\%s%d`, b, uint32(gamePID)+1500000000),
		)
	}
	return selfServeRailPipeSet(names)
}

// selfWatchGameRailPipes 对应样本 sub_14001AEE0 里每 50ms 的进程循环：
// crossfire.exe 可能重启成子进程，子进程的 rail_pipe_<childPid> 也必须建。
func selfWatchGameRailPipes(initialPID int, timeout time.Duration) {
	_ = initialPID
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for _, p := range selfProcessList() {
			if p.pid == 0 || !strings.EqualFold(p.name, "crossfire.exe") {
				continue
			}
			// 样本每 50ms 对每个 crossfire PID 调一次 sub_14001AB70；
			// 已建好的管道由 selfRailRoles 去重，失败的下一轮继续重试。
			selfServeGameRailPipes(int(p.pid))
		}
		time.Sleep(50 * time.Millisecond)
	}
}

type selfProcessInfo struct {
	pid  uint32
	name string
}

func selfProcessList() []selfProcessInfo {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	createSnap := k32.NewProc("CreateToolhelp32Snapshot")
	first := k32.NewProc("Process32FirstW")
	next := k32.NewProc("Process32NextW")
	closeHandle := k32.NewProc("CloseHandle")

	snap, _, _ := createSnap.Call(0x2, 0) // TH32CS_SNAPPROCESS
	if snap == 0 || syscall.Handle(snap) == syscall.InvalidHandle {
		return nil
	}
	defer closeHandle.Call(snap)

	type processEntry32 struct {
		Size            uint32
		CntUsage        uint32
		ProcessID       uint32
		DefaultHeapID   uintptr
		ModuleID        uint32
		CntThreads      uint32
		ParentProcessID uint32
		PriClassBase    int32
		Flags           uint32
		ExeFile         [260]uint16
	}
	var e processEntry32
	e.Size = uint32(unsafe.Sizeof(e))
	var out []selfProcessInfo
	ok, _, _ := first.Call(snap, uintptr(unsafe.Pointer(&e)))
	for ok != 0 {
		out = append(out, selfProcessInfo{pid: e.ProcessID, name: syscall.UTF16ToString(e.ExeFile[:])})
		e.Size = uint32(unsafe.Sizeof(e))
		ok, _, _ = next.Call(snap, uintptr(unsafe.Pointer(&e)))
	}
	return out
}

var (
	selfRailRoleMu sync.Mutex
	selfRailRoles  = map[string]string{} // name -> "server" / "client"
	selfRailLogMu  sync.Mutex
	selfRailLogged = map[string]bool{}
)

func selfRailLogOnce(key, format string, args ...any) {
	selfRailLogMu.Lock()
	if selfRailLogged[key] {
		selfRailLogMu.Unlock()
		return
	}
	selfRailLogged[key] = true
	selfRailLogMu.Unlock()
	selfLog(format, args...)
}

// selfServeRailPipeSet 按样本的角色处理每个管道名：
// 先尝试 CreateFileW 作为客户端连游戏已建的管道；连不上才 CreateNamedPipeW 当服务端。
// 样本 sub_140018D10 只在连上后发 8817 字节首包，所以角色反了就一直等不到 hello。
func selfServeRailPipeSet(names []string) func() {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	createPipe := k32.NewProc("CreateNamedPipeW")
	createFile := k32.NewProc("CreateFileW")
	connectPipe := k32.NewProc("ConnectNamedPipe")
	setPipeState := k32.NewProc("SetNamedPipeHandleState")
	readFile := k32.NewProc("ReadFile")
	writeFile := k32.NewProc("WriteFile")
	closeHandle := k32.NewProc("CloseHandle")

	var (
		mu      sync.Mutex
		handles []syscall.Handle
	)
	run := func(p syscall.Handle, server bool, name string) {
		if server {
			connectPipe.Call(uintptr(p), 0)
			mode := uint32(1)
			setPipeState.Call(uintptr(p), uintptr(unsafe.Pointer(&mode)), 0, 0)
		}
		selfLog("rail pipe %s role=%s", name, selfRailRoles[name])
		buf := make([]byte, 0x1000)
		for {
			var n uint32
			ok, _, _ := readFile.Call(uintptr(p), uintptr(unsafe.Pointer(&buf[0])),
				uintptr(len(buf)), uintptr(unsafe.Pointer(&n)), 0)
			if ok == 0 || n < 8 {
				return
			}
			ln := binary.LittleEndian.Uint32(buf[0:4])
			cmd := binary.LittleEndian.Uint32(buf[4:8])
			if ln > 0x400000 || n < 8+ln {
				return
			}
			payload := buf[8 : 8+ln]
			selfDumpFrame("host:"+fmt.Sprint(p), cmd, payload)

			// 离线器顺序：inject_hook cmd=1 -> 注入 -> ack -> send100。
			needHello := cmd == 1
			if needHello {
				selfInjectHookOnDemand()
			}

			resp := make([]byte, 12+len(payload))
			binary.LittleEndian.PutUint32(resp[0:4], uint32(4+len(payload)))
			binary.LittleEndian.PutUint32(resp[4:8], cmd)
			binary.LittleEndian.PutUint32(resp[8:12], 0)
			copy(resp[12:], payload)
			var w uint32
			writeFile.Call(uintptr(p), uintptr(unsafe.Pointer(&resp[0])),
				uintptr(len(resp)), uintptr(unsafe.Pointer(&w)), 0)
			if needHello {
				selfRailHelloOnce.Do(func() {
					if len(selfRailHello) == 0 {
						return
					}
					var hw uint32
					ok, _, err := writeFile.Call(uintptr(p), uintptr(unsafe.Pointer(&selfRailHello[0])),
						uintptr(len(selfRailHello)), uintptr(unsafe.Pointer(&hw)), 0)
					selfLog("rail hello sent %d %s ok=%v err=%v", len(selfRailHello), name, ok != 0, err)
				})
			}
		}
	}

	for _, name := range names {
		namePtr, _ := syscall.UTF16PtrFromString(name)

		selfRailRoleMu.Lock()
		role := selfRailRoles[name]
		selfRailRoleMu.Unlock()
		if role != "" {
			continue
		}

		// 先当客户端连游戏已创建的管道。
		fh, _, ferr := createFile.Call(
			uintptr(unsafe.Pointer(namePtr)),
			uintptr(0xC0000000), 0, 0, 3, 0, 0, // GENERIC_READ|WRITE, OPEN_EXISTING
		)
		if fh != 0 && syscall.Handle(fh) != syscall.InvalidHandle {
			pipe := syscall.Handle(fh)
			selfRailRoleMu.Lock()
			selfRailRoles[name] = "client"
			selfRailRoleMu.Unlock()
			mu.Lock()
			handles = append(handles, pipe)
			mu.Unlock()
			go run(pipe, false, name)
			continue
		}
		selfRailLogOnce(name+"|client", "rail pipe client miss %s err=%v", name, ferr)

		h, _, herr := createPipe.Call(
			uintptr(unsafe.Pointer(namePtr)),
			3, 0, 0xFF, uintptr(0x10000), uintptr(0x10000), 0, 0,
		)
		if h == 0 || syscall.Handle(h) == syscall.InvalidHandle {
			selfRailLogOnce(name+"|server", "rail pipe server fail %s err=%v", name, herr)
			continue
		}
		pipe := syscall.Handle(h)
		selfRailRoleMu.Lock()
		selfRailRoles[name] = "server"
		selfRailRoleMu.Unlock()
		mu.Lock()
		handles = append(handles, pipe)
		mu.Unlock()
		go run(pipe, true, name)
	}
	return func() {
		mu.Lock()
		for _, h := range handles {
			closeHandle.Call(uintptr(h))
		}
		handles = nil
		mu.Unlock()
	}
}

// selfEnsureHookDLL 把内嵌的 qwh_cf_hook.dll 释放到 %TEMP% 再注入。
// IDA: sub_14001AEE0 0x14001CE73 GetTempPathW + "qwh_cf_hook.dll"。
// 样本不往游戏目录写 hook；游戏目录受 TP/ACE 保护时 LoadLibrary 会被拦。
func selfEnsureHookDLL(gameDir string) (string, error) {
	if len(selfHookDLL) == 0 {
		return "", fmt.Errorf("未内嵌 hook DLL")
	}
	_ = gameDir
	p := filepath.Join(os.TempDir(), selfHookDLLName)
	if st, err := os.Stat(p); err == nil && st.Size() == int64(len(selfHookDLL)) {
		selfLog("hook 已存在 %s (%d)", p, st.Size())
		return p, nil
	}
	if err := os.WriteFile(p, selfHookDLL, 0o644); err != nil {
		selfLog("写 hook 失败 %s: %v", p, err)
		return "", err
	}
	selfLog("写 hook 成功 %s (%d)", p, len(selfHookDLL))
	return p, nil
}

// selfRemoveFakeRail 删除以前版本在 OP 路径留下的 1 字节假 rail。
// 假 rail 是 sub_140019270 的 NT 路径产物；OP 路径 sub_14001AEE0 不写它。
// 一旦 CF 搜索到坏 PE，就会直接 WeGame Init Failed (2310)。
func selfRemoveFakeRail(gameDir string) {
	if gameDir == "" {
		return
	}
	for _, p := range []string{
		gameDir + `\ntop_v17_dll\x64\rail_api64.dll`,
		filepath.Dir(gameDir) + `\ntop_v17_dll\x64\rail_api64.dll`,
	} {
		st, err := os.Stat(p)
		if err == nil && st.Size() == 1 {
			if rmErr := os.Remove(p); rmErr == nil {
				selfLog("已删除假 rail %s", p)
			}
		}
	}
	// 旧版本把 hook 放游戏目录；样本只从 %TEMP% 注入，这里清掉避免双注入。
	for _, p := range []string{
		gameDir + `\ntop_v17_dll\x64\qwh_cf_hook.dll`,
		filepath.Dir(gameDir) + `\ntop_v17_dll\x64\qwh_cf_hook.dll`,
	} {
		st, err := os.Stat(p)
		if err == nil && st.Size() == int64(len(selfHookDLL)) {
			if rmErr := os.Remove(p); rmErr == nil {
				selfLog("已删除旧 hook %s", p)
			}
		}
	}
}

// selfEnsureFakeRail 在 <游戏根>\ntop_v17_dll\x64 下放一个 1 字节的假 rail_api64.dll。
//
// 新版上号器（IDA sub_140011E50 @0x140011E50）的顺序：
//
//	sub_1400ABB10(L"ntop_v17_dll\\x64")  → 建子目录
//	sub_1400B10F0(..., L"\\rail_api64.dll") → fopen("wb") + fwrite("x",1,1)
//	sub_1400B10F0(..., L"\\crossfire.exe")  → 再从该目录启动
//
// CF 在这个路径拿到无效镜像后 LoadLibrary 失败，就不再走 WeGame 初始化，
// 因此不会弹 "WeGame Init Failed (2310)"。注意不要写游戏自己的 x64\rail_api64.dll，
// 那是共享内存/管道要用的真文件。
func selfEnsureFakeRail(gameDir string) {
	if strings.TrimSpace(gameDir) == "" {
		return
	}
	p := filepath.Join(gameDir, "ntop_v17_dll", "x64", "rail_api64.dll")
	if st, err := os.Stat(p); err == nil && st.Size() == 1 {
		selfLog("假 rail 已存在 %s", p)
		return
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		selfLog("建 ntop 目录失败 %s: %v", filepath.Dir(p), err)
		return
	}
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		selfLog("写假 rail 失败 %s: %v", p, err)
		return
	}
	selfLog("写假 rail 成功 %s", p)
}

// selfValidateVersionIni 复刻新版 sub_140011E50 的启动前校验：
// 读 x64\version.ini，必须同时含 -SSN:<当前区号>、-rotate:0、-cf:1，
// 且不能带 -HostIP=；不满足时新版直接放弃启动，这里提前报错，避免等到 2310。
func selfValidateVersionIni(dir, ssn string) error {
	cands := []string{
		filepath.Join(dir, "x64", "version.ini"),
		filepath.Join(dir, "version.ini"),
	}
	var body string
	var used string
	for _, p := range cands {
		if raw, err := os.ReadFile(p); err == nil {
			body, used = string(raw), p
			break
		}
	}
	if used == "" {
		return fmt.Errorf("version.ini 缺失（期望 %s）", cands[0])
	}
	if ssn != "" && !strings.Contains(body, "-SSN:"+ssn) {
		return fmt.Errorf("version.ini 的 SSN 与当前大区不一致（%s 里没有 -SSN:%s）", used, ssn)
	}
	for _, need := range []string{"-SSN:", "-rotate:0", "-cf:1"} {
		if !strings.Contains(body, need) {
			return fmt.Errorf("version.ini 缺少 %s（%s）", need, used)
		}
	}
	if strings.Contains(body, "-HostIP=") {
		return fmt.Errorf("version.ini 不应包含 -HostIP=（%s）", used)
	}
	return nil
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
