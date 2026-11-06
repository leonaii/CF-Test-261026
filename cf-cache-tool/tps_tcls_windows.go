//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// tclsUserInfoPath 是 tcls_core.exe 写 ACCOUNT_ID/UID 的文件。
func tclsUserInfoPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Tencent", "TCLSCore", "UserInfo.ini"), nil
}

var (
	tclsHex32 = regexp.MustCompile(`[0-9A-Fa-f]{32}`)
	tclsLong  = regexp.MustCompile(`[A-Za-z0-9+/=]{96,}`)
)

// tclsShmCandidates 生成 TCLS_CORE_SHARE_MEMORY_* 的候选名字。
// IDA sub_40C5B7 显示名字 = 前缀 + 结构 +666 处的字符串（通常是 ppid / game_id）。
// 这里用运行时拿到的进程 PID 兜底，并覆盖 Global/Local 命名空间。
func tclsShmCandidates() []string {
	suffixes := []string{"", "33", "2"}
	for _, p := range selfProcessList() {
		n := strings.ToLower(p.name)
		if n == "wegame.exe" || n == "tcls_core.exe" || n == "rail.exe" {
			suffixes = append(suffixes, fmt.Sprintf("%d", p.pid), fmt.Sprintf("888%d", p.pid))
		}
	}
	seen := map[string]bool{}
	var out []string
	for _, s := range suffixes {
		for _, prefix := range []string{"", "Local\\", "Global\\", `Session\1\`} {
			name := prefix + tclsMemMagic + s
			if !seen[name] {
				seen[name] = true
				out = append(out, name)
			}
		}
	}
	return out
}

// tclsOpenShm 以只读方式打开已存在的共享内存，返回 view 和释放函数。
func tclsOpenShm(name string) ([]byte, func(), bool) {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	openMap := k32.NewProc("OpenFileMappingW")
	mapView := k32.NewProc("MapViewOfFile")
	unmap := k32.NewProc("UnmapViewOfFile")

	np, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, nil, false
	}
	// 只读打开：共享内存由 wegame 创建，全访问权限可能被 DACL 拒绝。
	h, _, _ := openMap.Call(0x0004 /*FILE_MAP_READ*/, 0, uintptr(unsafe.Pointer(np)))
	if h == 0 {
		return nil, nil, false
	}
	view, _, _ := mapView.Call(h, 0x0004, 0, 0, uintptr(tclsMemSize))
	if view == 0 {
		syscall.CloseHandle(syscall.Handle(h))
		return nil, nil, false
	}
	mem := unsafe.Slice((*byte)(unsafe.Pointer(view)), tclsMemSize)
	return mem, func() {
		unmap.Call(view)
		syscall.CloseHandle(syscall.Handle(h))
	}, true
}

// tclsDescribeTicket 只读分析 1179B 结构，不改动它。
// 字段偏移尚未精确校准（IDA 里是 CheckShareMemoryValue 逐字段校验），
// 贸然按正则改结构可能让 tcls_core 崩溃，所以这里只记录特征位置，
// 供下一次用真实样本确定偏移后再做替换。
func tclsDescribeTicket(mem []byte, logf func(string, ...any)) {
	for i, hit := range tclsHex32.FindAllIndex(mem, -1) {
		if i >= 4 {
			break
		}
		logf("结构 +%d 处 32 位票据串: %s…", hit[0], string(mem[hit[0]:hit[0]+8]))
	}
	for i, hit := range tclsLong.FindAllIndex(mem, -1) {
		if i >= 2 {
			break
		}
		logf("结构 +%d 处长票据串: %d 字节", hit[0], hit[1]-hit[0])
	}
}

// tclsStartSharedMemProbe 周期性观测 TCLS_CORE_SHARE_MEMORY_*，命中就把
// 原始 1179B 结构 dump 到工作目录并打印 ASCII 摘要。只读，不破坏结构。
func tclsStartSharedMemProbe(dir string, logf func(string, ...any)) func() {
	done := make(chan struct{})
	go func() {
		reported := map[string]bool{}
		deadline := time.Now().Add(60 * time.Second)
		tick := time.NewTicker(500 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-done:
				return
			case <-tick.C:
			}
			if time.Now().After(deadline) {
				if len(reported) == 0 {
					logf("观测结束：未捕获 TCLS_CORE_SHARE_MEMORY_*（票据走落地文件，不依赖此通道）")
				}
				return
			}
			for _, name := range tclsShmCandidates() {
				mem, release, ok := tclsOpenShm(name)
				if !ok {
					continue
				}
				if !reported[name] {
					reported[name] = true
					snap := append([]byte(nil), mem...)
					p := filepath.Join(dir, "tcls_shm_orig.bin")
					_ = os.WriteFile(p, snap, 0o644)
					logf("捕获 TCLS 共享内存 %s（%d 字节），原始结构已存 %s", name, len(snap), p)
					logf("结构 ASCII: %s", tclsAsciiPreview(snap))
					tclsDescribeTicket(snap, logf)
				}
				release()
			}
		}
	}()
	return func() { close(done) }
}

// tclsStartGame 用降权令牌自己启动 TPS.exe，绕开官方 WeGame 已登录的原账号。
// 参数按 2026-11-04 ProcMon 抓到的官方命令行重建；-q 是 uin(QQ)。
// 必须降权：管理员令牌下 ACE/TP 会判启动环境异常（项目里已实测，见 selfop_runas.go）。
func tclsStartGame(gameDir, qq string, logf func(string, ...any)) (*selfGameProc, error) {
	exe := filepath.Join(gameDir, "Binaries", "Win64", "TPS.exe")
	if _, err := os.Stat(exe); err != nil {
		return nil, fmt.Errorf("找不到游戏主程序 %s: %w", exe, err)
	}
	args := fmt.Sprintf("-q %s -src=tgp -game_id %d -area 1 -zone_id 18220545", qq, tclsGameID)
	logf("降权启动 TPS.exe，参数: %s", args)
	proc, err := selfStartGameAsUser(exe, filepath.Dir(exe), args)
	if err != nil {
		return nil, err
	}
	logf("TPS.exe 已启动 pid=%d", proc.ProcessID)
	return proc, nil
}

func tclsAsciiPreview(b []byte) string {
	var s strings.Builder
	for _, c := range b {
		if c >= 32 && c < 127 {
			s.WriteByte(c)
		} else {
			s.WriteByte('.')
		}
	}
	out := s.String()
	if len(out) > 400 {
		out = out[:400] + "…"
	}
	return out
}
