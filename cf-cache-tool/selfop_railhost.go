package main

// Rail Host：新版上号器在启动 crossfire.exe 之前，会先在自己进程里把
// 游戏目录的 rail_api64.dll 加载进来（IDA sub_14009F540 检查这个模块，
// 返回 -1 就直接放弃启动）。rail 的 DllMain 会以当前进程为 Host 建立
//
//	Rail.SharedMemory.<HWND>、rail_named_pipe_host_<pid>、rail_pipe_15
//
// 之后再拉游戏，guest 侧握手才成立。自研之前只靠自己手写共享内存和管道，
// 少了这一步，rail 找不到 Host 就弹 "WeGame Init Failed. (2310)"。
//
// 顺序必须是：建 WE_GAME_SDK 窗口 → LoadLibrary(rail) → 再 CreateProcess。

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	selfRailHostOnce sync.Once
	selfRailHostDLL  syscall.Handle
)

// selfLoadRailHost 在自研进程里加载 rail_api64.dll，让它自己当 Rail Host。
// 必须在 CreateProcess 之前调用；失败时只记日志，不阻断（部分版本 rail 会自检失败）。
func selfLoadRailHost(gameDir string) {
	selfRailHostOnce.Do(func() {
		// rail 需要 WE_GAME_SDK_WINDOW_NAME_CLASS 窗口才能初始化 Host
		hwnd := selfEnsureRailWindow()
		selfLog("rail host 窗口 hwnd=%#x", hwnd)

		cands := []string{
			filepath.Join(gameDir, "x64", "rail_api64.dll"),
			filepath.Join(gameDir, "rail_api64.dll"),
			filepath.Join(gameDir, "ntop_v17_dll", "x64", "rail_api64.dll"),
		}
		for _, p := range cands {
			st, err := os.Stat(p)
			if err != nil || st.Size() < 1024 {
				continue
			}
			h, err := syscall.LoadLibrary(p)
			if err != nil {
				selfLog("LoadLibrary rail 失败 %s: %v", p, err)
				continue
			}
			selfRailHostDLL = h
			selfLog("已加载 rail host: %s", p)
			return
		}
		selfLog("rail host 未加载：游戏目录里没有可用的 rail_api64.dll")
	})
}
