package main

// 交互式选号：按卡密列出所有已缓存的 redeem 数据，运行时也可随时切换。

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	errQuit    = errors.New("quit")
	errRefresh = errors.New("refresh")
)

type accountRow struct {
	Card     string
	QQ       string
	OpenID   string
	Note     string
	RawNote  string
	Zone     string
	Expires  string
	Recorded time.Time
	Decoded  bool
	Remote   bool // true = 只在云端，本地还没下载

	// 验号接口（/api/query/status）的异步结果，仅内存展示，不落盘
	Verify    string
	VerifyRaw string
	VerifyErr string
	VerifyAt  time.Time

	// 自研上号状态，仅内存
	SelfState string // ""=未开始 / 进行中 / 成功 / 失败
	SelfErr   string
}

func anyStr(v any) string {
	switch t := v.(type) {
	case nil:
		return "-"
	case string:
		if t == "" {
			return "-"
		}
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(t)
	default:
		return fmt.Sprint(t)
	}
}

func truncRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + ".."
}

// sortAccountRows 最新录制排最前；时间相同时按卡密倒序，保证顺序稳定。
func sortAccountRows(rows []accountRow) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Recorded.Equal(rows[j].Recorded) {
			return rows[i].Card > rows[j].Card
		}
		return rows[i].Recorded.After(rows[j].Recorded)
	})
}

// loadAccountRows 读取全部 redeem 缓存，尽力解出账号信息用于展示。
func loadAccountRows(c *Cache) []accountRow {
	metas, err := c.List()
	if err != nil {
		return nil
	}
	seen := map[string]bool{}
	rows := make([]accountRow, 0, len(metas))
	for _, m := range metas {
		if m.Endpoint != "api_card_redeem" || seen[m.Card] {
			continue
		}
		seen[m.Card] = true
		recorded := m.RecordedAt
		if recorded.IsZero() {
			// 老缓存/手工拷贝可能没有 recorded_at，退回文件修改时间，保证“最新在前”仍然成立
			metaPath := filepath.Join(c.Dir, sanitizeName(m.Card), m.Endpoint+".meta.json")
			if st, err := os.Stat(metaPath); err == nil {
				recorded = st.ModTime()
			}
		}
		row := accountRow{
			Card:     m.Card,
			QQ:       "-",
			OpenID:   "-",
			Note:     "-",
			Zone:     "-",
			Expires:  "-",
			Recorded: recorded,
		}
		body, _, err := c.Load(m.Card, m.Endpoint)
		if err == nil {
			// 外层信封：线上协议把 qq / zone / note / expires_at 直接放在这一层
			var env any
			if json.Unmarshal(body, &env) == nil {
				row.QQ = anyStr(findField(env, "qq", "qqnumber", "qqnum", "qqid"))
				row.RawNote = anyStr(findField(env, "note", "extra", "remark", "comment", "memo", "desc"))
				row.Note = noteSummary(row.RawNote)
				row.Zone = anyStr(findField(env, "zone", "zoneid", "areaid", "server"))
				row.Expires = formatExpires(findField(env, "expires_at", "expire_at", "expiretime", "expires"))
			}
			// 内层载荷：可能是 JSON，也可能是 [QQ] 开头的 Key=Value 文本
			if plain, err := OpenEnvelope(body); err == nil {
				row.Decoded = true
				kv, section := parsePayload(plain)
				if row.QQ == "-" {
					if section != "" {
						row.QQ = section
					} else if v := kv["qq"]; v != "" {
						row.QQ = v
					}
				}
				if v := kv["openid"]; v != "" {
					row.OpenID = v
				} else if v := kv["uid"]; v != "" {
					row.OpenID = v
				}
			}
			// 兜底：peek 响应是明文 JSON，通常也带 qq
			if row.QQ == "-" || row.OpenID == "-" {
				if pb, _, e := c.Load(m.Card, "api_card_peek"); e == nil {
					var pk any
					if json.Unmarshal(pb, &pk) == nil {
						if row.QQ == "-" {
							row.QQ = anyStr(findField(pk, "qq", "qqnumber"))
						}
						if row.OpenID == "-" {
							row.OpenID = anyStr(findField(pk, "openid", "open_id"))
						}
						if row.RawNote == "-" {
							row.RawNote = anyStr(findField(pk, "note", "extra"))
							row.Note = noteSummary(row.RawNote)
						}
					}
				}
			}
		}
		rows = append(rows, row)
	}
	sortAccountRows(rows)
	return rows
}

func printAccountTable(rows []accountRow) {
	fmt.Println()
	fmt.Printf("%-4s %-18s %-12s %-14s %-12s %-14s %-12s %s\n",
		"#", "CARD", "QQ", "OPENID/UID", "NOTE", "ZONE", "RECORDED", "RAW NOTE")
	fmt.Println(strings.Repeat("-", 150))
	for i, r := range rows {
		fmt.Printf("%-4d %-18s %-12s %-14s %-12s %-14s %-12s %s\n",
			i+1, r.Card, truncRunes(r.QQ, 12), truncRunes(r.OpenID, 14),
			truncRunes(r.Note, 12), truncRunes(r.Zone, 14),
			r.Recorded.Format("01-02 15:04"), truncRunes(r.RawNote, 48))
	}
	fmt.Println(strings.Repeat("-", 150))
}

func describeRow(r accountRow) string {
	return fmt.Sprintf("[%s] qq=%s openid=%s note=%s zone=%s 录制于 %s\n  原始备注: %s",
		r.Card, r.QQ, r.OpenID, r.Note, r.Zone, r.Recorded.Format("2006-01-02 15:04:05"), r.RawNote)
}

// chooseAccount 打印列表并读取用户输入，返回 0-based 下标。
func chooseAccount(rows []accountRow) (int, error) {
	printAccountTable(rows)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\n请选择要使用的账号 [1-%d，回车=1，r=刷新列表，q=退出]: ", len(rows))
		line, err := reader.ReadString('\n')
		s := strings.TrimSpace(line)
		if s == "" && err != nil {
			return 0, errors.New("当前环境无法读取键盘输入，请改用 -pin <卡密> 指定")
		}
		switch strings.ToLower(s) {
		case "":
			if len(rows) == 0 {
				return 0, errRefresh
			}
			return 0, nil
		case "q":
			return 0, errQuit
		case "r":
			return 0, errRefresh
		}
		n, convErr := strconv.Atoi(s)
		if convErr != nil || n < 1 || n > len(rows) {
			fmt.Println("输入无效，请输入列表中的序号。")
			continue
		}
		return n - 1, nil
	}
}

// switchLoop 运行期间持续读取序号，实现“热切换”账号。
// quitOnQ 为 true 时输入 q 返回 errQuit（serve 模式退出服务）；否则只停止切换监听。
// 输入被重定向/EOF 时返回 nil，调用方应保持服务继续运行。
func switchLoop(srv *Server, quitOnQ bool) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		if quitOnQ {
			fmt.Print("\n[切换] 序号 / r 刷新 / q 退出: ")
		} else {
			fmt.Print("\n[切换] 序号 / r 刷新 / 直接关闭窗口结束: ")
		}
		line, err := reader.ReadString('\n')
		s := strings.TrimSpace(line)
		if s == "" && err != nil {
			// 没有可读输入（例如被批处理重定向），保持服务运行
			return nil
		}
		switch strings.ToLower(s) {
		case "q":
			if quitOnQ {
				return errQuit
			}
			continue
		case "":
			continue
		case "r":
			rows := loadAccountRows(srv.cache)
			if len(rows) == 0 {
				fmt.Println("暂无缓存。")
				continue
			}
			printAccountTable(rows)
			continue
		}
		rows := loadAccountRows(srv.cache)
		n, convErr := strconv.Atoi(s)
		if convErr != nil || n < 1 || n > len(rows) {
			fmt.Println("输入无效。")
			continue
		}
		srv.SetPin(rows[n-1].Card)
		fmt.Println("已切换 -> " + describeRow(rows[n-1]))
	}
}

// runPinned 交互选号 / 固定卡密 模式的主流程。
func runPinned(opts Options, launch bool, exePath, dir string, pick bool, pin string) error {
	srv := NewServer(opts)

	if pin != "" {
		if _, _, err := srv.cache.Load(pin, "api_card_redeem"); err != nil && !opts.Refresh {
			return fmt.Errorf("卡密 %s 没有 redeem 缓存，无法固定；请先录制或去掉 -pin", pin)
		}
		srv.SetPin(pin)
		logf("已固定卡密: %s", pin)
	}

	ln, err := net.Listen("tcp", opts.Listen)
	if err != nil {
		return fmt.Errorf("端口 %s 不可用: %w", opts.Listen, err)
	}
	_ = ln.Close()

	go func() {
		if err := srv.Run(); err != nil {
			logf("服务退出: %v", err)
		}
	}()
	if !waitPort(opts.Listen, 3*time.Second) {
		return fmt.Errorf("本地服务未能在 %s 启动", opts.Listen)
	}

	if pick && pin == "" {
		for {
			rows := loadAccountRows(srv.cache)
			if len(rows) == 0 {
				return errors.New("还没有任何 redeem 缓存，请先正常上号录制一次")
			}
			idx, err := chooseAccount(rows)
			if err == errQuit {
				return nil
			}
			if err == errRefresh {
				continue
			}
			if err != nil {
				return err
			}
			srv.SetPin(rows[idx].Card)
			logf("已选择 -> %s", describeRow(rows[idx]))
			break
		}
	}

	if launch {
		cmd := exec.Command(exePath)
		cmd.Dir = absPath(dir)
		// 交互模式不把控制台输入交给上号器，保证还能接收切换指令
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Start(); err != nil {
			if isElevationRequired(err) {
				return fmt.Errorf("上号器要求管理员权限：请右键以管理员身份运行本工具")
			}
			return err
		}
		logf("已启动上号器: %s", exePath)
		go func() { _ = switchLoop(srv, false) }()
		if err := cmd.Wait(); err != nil {
			logf("上号器退出: %v", err)
		}
		return nil
	}

	if pin != "" {
		// 固定卡密模式不需要交互，直接常驻
		select {}
	}
	if err := switchLoop(srv, true); err == errQuit {
		return nil
	}
	// 输入不可用（重定向/EOF）时保持服务常驻，用 Ctrl+C 或关闭窗口结束
	select {}
}

// normKey 归一化字段名：转小写并去掉下划线、短横线等符号，
// 这样 QQ / qq / "QQ号" / OpenId / open_id 都能命中同一个名字。
func normKey(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// findField 按优先级在 JSON 结构里递归查找字段，忽略大小写与分隔符。
func findField(v any, names ...string) any {
	for _, name := range names {
		if got := searchKey(v, normKey(name), 0); got != nil {
			return got
		}
	}
	return nil
}

func searchKey(v any, want string, depth int) any {
	if depth > 5 || v == nil {
		return nil
	}
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			if normKey(k) == want && val != nil {
				return val
			}
		}
		for _, val := range t {
			switch val.(type) {
			case map[string]any, []any:
				if got := searchKey(val, want, depth+1); got != nil {
					return got
				}
			}
		}
	case []any:
		for _, item := range t {
			if got := searchKey(item, want, depth+1); got != nil {
				return got
			}
		}
	}
	return nil
}

// parsePayload 解析解密后的账号数据：
// 既支持 JSON，也支持线上实际使用的 [QQ]\nKey=Value 文本格式。
func parsePayload(plain []byte) (map[string]string, string) {
	kv := map[string]string{}
	var js any
	if json.Unmarshal(plain, &js) == nil {
		if m, ok := js.(map[string]any); ok {
			for k, v := range m {
				switch v.(type) {
				case string, float64, bool:
					kv[normKey(k)] = anyStr(v)
				}
			}
		}
		return kv, ""
	}

	section := ""
	for _, line := range strings.Split(string(plain), "\n") {
		line = strings.TrimSpace(strings.TrimRight(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if section == "" {
				section = strings.TrimSpace(line[1 : len(line)-1])
			}
			continue
		}
		if i := strings.IndexByte(line, '='); i > 0 {
			kv[normKey(line[:i])] = strings.TrimSpace(line[i+1:])
		}
	}
	return kv, section
}

// formatExpires 兼容毫秒/秒时间戳与字符串。
func formatExpires(v any) string {
	switch t := v.(type) {
	case float64:
		n := int64(t)
		switch {
		case n > 1e12:
			return time.UnixMilli(n).Local().Format("2006-01-02 15:04")
		case n > 1e9:
			return time.Unix(n, 0).Local().Format("2006-01-02 15:04")
		}
		return strconv.FormatInt(n, 10)
	case string:
		if t == "" {
			return "-"
		}
		return t
	}
	return anyStr(v)
}

// noteSummary 线上 note 是 "编号:xxx----卡密----战区----区服----段位:宗师----..."，
// 表格里优先显示段位，找不到就原样返回。
func noteSummary(note string) string {
	if note == "" || note == "-" {
		return "-"
	}
	for _, seg := range strings.Split(note, "----") {
		seg = strings.TrimSpace(seg)
		if strings.HasPrefix(seg, "段位") {
			return seg
		}
	}
	return note
}
