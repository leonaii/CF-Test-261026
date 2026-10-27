package main

// 退出清理：把为绕过 WeGame 2310 而写进游戏目录的伪造文件删掉。
//
// 写这些文件的依据是样本 sub_140019270：
//   <base>\ntop_v17_dll\x64\rail_api64.dll  = 1 字节 "x"（0x14001A51C，fwrite("x",1,1)）
//   <base>\ntop_v17_dll\x64\qwh_cf_hook.dll = 注入用 hook
// 其中 <base> 是 crossfire.exe 所在目录（<游戏>\x64），早期版本错放在 <游戏>。
//
// 清理是"校验后删除"：只有长度对得上才删，避免误删游戏自带文件。

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	selfFakeRailName = "rail_api64.dll"
	selfHookDLLName  = "qwh_cf_hook.dll"
)

// selfRailDirs 返回可能被我们写过的 ntop_v17_dll\x64 目录（去重）。
func selfRailDirs() []string {
	seen := map[string]bool{}
	var out []string
	add := func(base string) {
		if base == "" {
			return
		}
		p := filepath.Join(base, "ntop_v17_dll", "x64")
		k := strings.ToLower(p)
		if seen[k] {
			return
		}
		seen[k] = true
		out = append(out, p)
	}
	if d, err := selfGameDir(); err == nil && d != "" {
		add(filepath.Join(d, "x64")) // 样本用的位置（crossfire.exe 同级）
		add(d)                       // 旧版本误放的位置
	}
	add(exeDir())
	return out
}

// selfCleanupRailArtifacts 删除伪造文件，返回删除个数。
func selfCleanupRailArtifacts() int {
	n := 0
	for _, dir := range selfRailDirs() {
		n += selfRemoveIfSize(dir, selfFakeRailName, func(sz int64) bool { return sz == 1 })
		n += selfRemoveIfSize(dir, selfHookDLLName, func(sz int64) bool {
			return len(selfHookDLL) > 0 && sz == int64(len(selfHookDLL))
		})
	}
	return n
}

// selfRemoveIfSize 只在文件长度满足条件时删除。
func selfRemoveIfSize(dir, name string, ok func(int64) bool) int {
	p := filepath.Join(dir, name)
	st, err := os.Stat(p)
	if err != nil || st.IsDir() {
		return 0
	}
	if !ok(st.Size()) {
		return 0
	}
	if err := os.Remove(p); err != nil {
		return 0
	}
	return 1
}
