package main

// 本地落盘的卡密一律用可视形态（CG 开头 32 位），磁盘上不出现明文真卡密：
//
//	cache/<可视卡密>/*.meta.json   card 字段写可视串，目录名也是可视串
//	cache.recording/<可视卡密>/    录制暂存同理
//	history.jsonl                  card 字段写可视串
//	config.json                    cloud_name / marks 的键写可视串
//	verify_snapshot.json           card 字段写可视串
//	data/dump/cards.jsonl           card_key 写可视串，note/zone 里的卡密一并渲染
//
// 老版本直接存真卡密，启动时 migrateLocalCardStore 一次性改名/改写，
// 之后再读写都走 storeCardKey / loadCardKey，内存里始终是真卡密。

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// migrateStartupWorkDir 在启动时迁移当前工作目录（配置还没写时用默认目录）。
func migrateStartupWorkDir(cfg *appConfig) {
	dir := defaultWorkDir()
	if cfg != nil && strings.TrimSpace(cfg.WorkDir) != "" {
		dir = cfg.WorkDir
	}
	migrateLocalCardStore(dir)
}

// migratedDirs 记住已经迁移过的工作目录，避免每次启动重复扫一遍缓存。
var migratedDirs sync.Map

// migrateLocalCardStore 迁移一个工作目录下所有本地落盘（同一目录只跑一次）。
func migrateLocalCardStore(dir string) {
	if dir == "" {
		return
	}
	dir = filepath.Clean(dir)
	if _, loaded := migratedDirs.LoadOrStore(dir, true); loaded {
		return
	}
	migrateCardDirs(filepath.Join(dir, "cache"))
	migrateCardDirs(filepath.Join(dir, stagingDirName))
	migrateHistoryStore(historyPath(dir))
	migrateVerifySnapStore(verifySnapPath())
	migrateDumpCardsStore(dumpCardsPath())
}

// migrateCardDirs 把 cache 下明文卡密目录改名成可视形态，并改写 meta 里的 card。
func migrateCardDirs(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		card, ok := decodeCardKey(name)
		if !ok {
			if !looksLikeCard(name) {
				continue // 不是卡密目录（例如 unknown），不动
			}
			card = name
		}
		want := cardDirName(card)
		if want == "" {
			continue
		}
		src := filepath.Join(root, name)
		dst := filepath.Join(root, want)
		if want != name {
			if _, err := os.Stat(dst); err == nil {
				mergeLegacyDir(src, dst)
			} else if err := os.Rename(src, dst); err != nil {
				mergeLegacyDir(src, dst)
			}
			src = dst
		}
		rewriteMetaCards(src, card)
	}
}

// mergeLegacyDir 把老目录里目标没有的文件补过去，然后删掉老目录（名字带明文卡密）。
func mergeLegacyDir(src, dst string) {
	_ = os.MkdirAll(dst, 0o700)
	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}
	ok := true
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		to := filepath.Join(dst, e.Name())
		if _, err := os.Stat(to); err == nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			ok = false
			continue
		}
		if err := os.WriteFile(to, data, 0o600); err != nil {
			ok = false
		}
	}
	// 有文件没搬过去就保留老目录，宁可留下明文也不能丢数据
	if ok {
		_ = os.RemoveAll(src)
	}
}

// rewriteMetaCards 把目录里所有 meta 的 card 字段改成可视形态。
func rewriteMetaCards(dir, card string) {
	want := storeCardKey(card)
	metas, err := filepath.Glob(filepath.Join(dir, "*.meta.json"))
	if err != nil {
		return
	}
	for _, mp := range metas {
		raw, err := os.ReadFile(mp)
		if err != nil {
			continue
		}
		var m Meta
		if json.Unmarshal(raw, &m) != nil {
			continue
		}
		changed := false
		if m.Card != want {
			m.Card = want
			changed = true
		}
		// 请求体（card_key）里也可能带明文卡密
		if req := renderCardInText(m.Request); req != m.Request {
			m.Request = req
			changed = true
		}
		if !changed {
			continue
		}
		if b, err := json.MarshalIndent(&m, "", "  "); err == nil {
			_ = writeFileAtomic(mp, b)
		}
	}
}

// migrateHistoryStore 把 history.jsonl 每行的 card 改写成可视形态。
func migrateHistoryStore(p string) {
	lines := readHistoryLines(p)
	if len(lines) == 0 {
		return
	}
	changed := false
	for i, line := range lines {
		var e histEntry
		if json.Unmarshal([]byte(line), &e) != nil {
			continue
		}
		wantCard := storeCardKey(loadCardKey(e.Card))
		wantNote := renderCardInText(e.Note)
		if e.Card == wantCard && e.Note == wantNote {
			continue
		}
		e.Card, e.Note = wantCard, wantNote
		b, err := json.Marshal(e)
		if err != nil {
			continue
		}
		lines[i] = string(b)
		changed = true
	}
	if !changed {
		return
	}
	writeLinesAtomic(p, lines)
}

// migrateVerifySnapStore 把查封快照里的 card 改写成可视形态。
func migrateVerifySnapStore(p string) {
	raw, err := os.ReadFile(p)
	if err != nil || len(raw) == 0 {
		return
	}
	var list []verifySnap
	if json.Unmarshal(raw, &list) != nil {
		return
	}
	changed := false
	for i := range list {
		want := storeCardKey(loadCardKey(list[i].Card))
		msg := renderCardInText(list[i].Message)
		if list[i].Card == want && list[i].Message == msg {
			continue
		}
		list[i].Card = want
		list[i].Message = msg
		changed = true
	}
	if !changed {
		return
	}
	if b, err := json.Marshal(list); err == nil {
		writeBytesAtomic(p, b)
	}
}

// migrateDumpCardsStore 把离线列表缓存里的明文卡密改成可视形态。
func migrateDumpCardsStore(p string) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 256<<10), 8<<20)
	var out []string
	changed := false
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) > 0 {
			fixed := transformDumpCardLine(line)
			if string(fixed) != string(line) {
				changed = true
			}
			line = fixed
		}
		out = append(out, string(line))
	}
	_ = f.Close()
	if !changed || sc.Err() != nil {
		return
	}
	writeLinesAtomic(p, out)
}

// transformDumpCardLine 处理 cards.jsonl 的一行：card_key 换成可视串，
// note/zone 里的卡密一并渲染。解析失败原样返回，不破坏离线数据。
func transformDumpCardLine(line []byte) []byte {
	if len(line) == 0 {
		return line
	}
	var obj map[string]json.RawMessage
	if json.Unmarshal(line, &obj) != nil {
		return line
	}
	changed := false
	put := func(k, v string) {
		b, err := json.Marshal(v)
		if err != nil {
			return
		}
		obj[k] = b
		changed = true
	}
	if raw, ok := obj["card_key"]; ok {
		var s string
		if json.Unmarshal(raw, &s) == nil {
			if v := storeCardKey(loadCardKey(s)); v != s {
				put("card_key", v)
			}
		}
	}
	for _, k := range []string{"note", "zone"} {
		raw, ok := obj[k]
		if !ok {
			continue
		}
		var s string
		if json.Unmarshal(raw, &s) != nil {
			continue
		}
		if v := renderCardInText(s); v != s {
			put(k, v)
		}
	}
	if !changed {
		return line
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return out
}

// writeDumpCards 把远端 cards.jsonl 逐行转成可视形态后落盘。
func writeDumpCards(dst io.Writer, src io.Reader) error {
	sc := bufio.NewScanner(src)
	sc.Buffer(make([]byte, 0, 256<<10), 8<<20)
	w := bufio.NewWriter(dst)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) > 0 {
			line = transformDumpCardLine(line)
		}
		if _, err := w.Write(line); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return w.Flush()
}

func writeLinesAtomic(p string, lines []string) {
	writeBytesAtomic(p, []byte(strings.Join(lines, "\n")+"\n"))
}

func writeBytesAtomic(p string, data []byte) {
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return
	}
	_ = os.Rename(tmp, p)
}
