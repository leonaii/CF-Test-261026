package main

// 本机使用历史：每次离线上号 / 查询封号往 history.jsonl 追加一行。
// 只记本机行为，不上传，不做去重。文件超过 2MB 时自动裁剪到最近 5000 条。

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type histEntry struct {
	TS     int64  `json:"ts"`               // 毫秒
	Action string `json:"action"`           // offline / verify
	QQ     string `json:"qq"`
	Card   string `json:"card"`
	Note   string `json:"note,omitempty"`
	Zone   string `json:"zone,omitempty"`
	Result string `json:"result,omitempty"`
}

const (
	histFile    = "history.jsonl"
	histMaxLoad = 2000
	histMaxKeep = 5000
	histMaxSize = 2 << 20
)

var histMu sync.Mutex

func historyPath(dir string) string { return filepath.Join(dir, histFile) }

// appendHistory 追加一条记录，失败静默（历史记录不该影响主流程）。
func appendHistory(dir string, e histEntry) {
	if dir == "" {
		return
	}
	b, err := json.Marshal(e)
	if err != nil {
		return
	}

	histMu.Lock()
	defer histMu.Unlock()

	p := historyPath(dir)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return
	}
	_, _ = f.Write(append(b, '\n'))
	_ = f.Close()

	if st, err := os.Stat(p); err == nil && st.Size() > histMaxSize {
		trimHistoryLocked(p)
	}
}

// trimHistoryLocked 只保留最近 histMaxKeep 条（调用方已持锁）。
func trimHistoryLocked(p string) {
	lines := readHistoryLines(p)
	if len(lines) <= histMaxKeep {
		return
	}
	keep := lines[len(lines)-histMaxKeep:]
	_ = os.WriteFile(p, []byte(strings.Join(keep, "\n")+"\n"), 0o600)
}

func readHistoryLines(p string) []string {
	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()

	var out []string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

// loadHistory 读取历史，按时间倒序返回最近 limit 条。
func loadHistory(dir string, limit int) []histEntry {
	if dir == "" {
		return nil
	}
	histMu.Lock()
	lines := readHistoryLines(historyPath(dir))
	histMu.Unlock()

	out := make([]histEntry, 0, len(lines))
	for _, line := range lines {
		var e histEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		out = append(out, e)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].TS > out[j].TS })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// clearHistory 清空历史文件。
func clearHistory(dir string) error {
	if dir == "" {
		return nil
	}
	histMu.Lock()
	defer histMu.Unlock()
	return os.Remove(historyPath(dir))
}
