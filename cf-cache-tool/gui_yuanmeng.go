package main

// yuanmeng 页：数据全部走卡密平台接口，不再读本地 dump。
//
//	GET /api/admin/cards   列表（状态/关键字筛选、分页）
//	GET /api/admin/data/<qq>  ini 明文（上号时按需拉取）
//
// 列表只展示 note（备注），卡密列用于复制/查号/上号。

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lxn/walk"
)

// ymmTicket 是 ini 解析结果。
type ymmTicket struct {
	QQ         string
	UID        string
	Tgpid      string
	CTID       string
	MachineID  string
	CT         string
	ST         string
	DT         string
	BizTicket  string
	ClientKey  string
	STKey      string
	Sig        string
	SessionKey string
	Token0106  string
	Token010A  string
	Token0143  string
}

// ymmRow 是列表里的一行（只保留页面要用的字段）。
type ymmRow struct {
	QQ          string
	CardKey     string
	QueryCode   string
	Note        string
	Status      string
	CreatedAt   int64 // 毫秒
	ActivatedAt int64 // 毫秒

	// 以下从 note 里解析，-1 表示没解析到（例如未激活的卡 note 为空）
	Level   int
	VIP     int
	Hero    int
	Area    string
	WarZone string

	Hist bool // 来自本机历史：时间显示到秒、排序固定按操作时间倒序
}

// ymmParseNote 从 note 解析等级 / VIP / 英雄级武器数量 / 区服 / 战区。
// note 形如：
//
//	编号:HY_X----847955166----CFXXX----等级:101----山东网通一区----北部战区----点券:6003----安全分:100----vip:0----心悦:游戏家Pro
//
// 带冒号的段按键名取，不带冒号的段按出现顺序当成区服、战区。
func (r *ymmRow) ymmParseNote() {
	r.Level, r.VIP, r.Hero = -1, -1, -1
	note := r.Note
	if note == "" {
		return
	}

	// 区服/战区是不带键名的段，位置在第 3/4 段（前两段是 QQ、卡密）。
	// 两种 note 格式顺序不同，按内容判断：
	//   老：...----等级:101----山东网通一区----北部战区----点券:...
	//   新：...----南部战区----湖北电信一区----段位:专家----...
	plain := 0
	for _, seg := range strings.Split(note, "----") {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		if i := strings.Index(seg, ":"); i > 0 {
			k := strings.TrimSpace(seg[:i])
			v := strings.TrimSpace(seg[i+1:])
			switch k {
			case "等级", "-等级":
				if n, err := strconv.Atoi(v); err == nil {
					r.Level = n
				}
			case "vip", "VIP":
				if n, err := strconv.Atoi(v); err == nil {
					r.VIP = n
				}
			case "英雄级武器数量":
				if n, err := strconv.Atoi(v); err == nil {
					r.Hero = n
				}
			}
			continue
		}

		plain++
		if plain <= 2 { // 前两段是 QQ 和卡密
			continue
		}
		if strings.Contains(seg, "战区") {
			if r.WarZone == "" {
				r.WarZone = seg
			}
		} else if strings.Contains(seg, "区") {
			if r.Area == "" {
				r.Area = seg
			}
		}
	}
}

// createdText 创建时间按本地时区显示；历史行精确到秒，方便区分同一分钟内的多次操作。
func (r ymmRow) createdText() string {
	if r.CreatedAt <= 0 {
		return "-"
	}
	if r.Hist {
		return time.UnixMilli(r.CreatedAt).Format("2006-01-02 15:04:05")
	}
	return time.UnixMilli(r.CreatedAt).Format("2006-01-02 15:04")
}

// activatedText 激活时间，未激活显示 -。
func (r ymmRow) activatedText() string {
	if r.ActivatedAt <= 0 {
		return "-"
	}
	return time.UnixMilli(r.ActivatedAt).Format("2006-01-02 15:04")
}

// ymmModel 是 TableView 的模型：只有「备注」「卡密」两列。
type ymmModel struct {
	walk.TableModelBase
	rows []ymmRow
}

func (m *ymmModel) RowCount() int { return len(m.rows) }

func (m *ymmModel) Value(row, col int) interface{} {
	if row < 0 || row >= len(m.rows) {
		return ""
	}
	r := m.rows[row]
	switch col {
	case 0:
		return r.Note
	case 1:
		if r.Level < 0 {
			return "-"
		}
		return r.Level
	case 2:
		if r.VIP < 0 {
			return "-"
		}
		return r.VIP
	case 3:
		if r.Hero < 0 {
			return "-"
		}
		return r.Hero
	case 4:
		return dash(r.WarZone)
	case 5:
		return r.createdText()
	case 6:
		return r.activatedText()
	}
	return ""
}

// ymmFilter 是状态下拉项。
type ymmFilter struct {
	Label  string
	Status string        // 传给接口的 status
	Since  time.Duration // >0 时按 activated_at 本地过滤
}

var ymmFilters = []ymmFilter{
	{"未使用", "unused", 0},
	{"24小时内激活", "active", 24 * time.Hour},
	{"使用中", "active", 0},
	{"已过期", "expired", 0},
	{"已作废", "revoked", 0},
	{"全部", "all", 0},
	{"本机历史", "history", 0},
}

// histKey 按卡密聚合（没有卡密时退回 QQ）。
func histKey(e histEntry) string {
	if s := strings.TrimSpace(e.Card); s != "" {
		return s
	}
	return strings.TrimSpace(e.QQ)
}

func histActionLabel(a string) string {
	switch a {
	case "offline":
		return "②上号"
	case "selfop":
		return "③自研上号"
	case "verify":
		return "①查封"
	}
	return a
}

func histEntryMatch(e histEntry, kw string) bool {
	for _, f := range []string{e.QQ, e.Card, e.Note, e.Zone, e.Result, histActionLabel(e.Action)} {
		if strings.Contains(strings.ToLower(f), kw) {
			return true
		}
	}
	return false
}

// ymmHistoryRows 把本机历史转成列表行：同一卡密只显示**最后一次操作**，
// 时间就是最后一次操作时间，后面附上累计次数。
func ymmHistoryRows(dir, kw string) []ymmRow {
	entries := loadHistory(dir, histMaxLoad) // 已按时间倒序
	kw = strings.ToLower(strings.TrimSpace(kw))

	// 先按关键字过滤（搜索时次数只统计命中的那些）
	matched := make([]histEntry, 0, len(entries))
	for _, e := range entries {
		if e.Action != "offline" && e.Action != "selfop" { // 历史只看上号记录，查封不进列表
			continue
		}
		if kw != "" && !histEntryMatch(e, kw) {
			continue
		}
		matched = append(matched, e)
	}

	count := make(map[string]int, len(matched))
	for _, e := range matched {
		count[histKey(e)]++
	}

	seen := make(map[string]bool, len(count))
	rows := make([]ymmRow, 0, len(count))
	for _, e := range matched { // matched 保持时间倒序，每个卡密取第一条即最后一次
		k := histKey(e)
		if k == "" || seen[k] {
			continue
		}
		seen[k] = true

		note := histActionLabel(e.Action) + " · " + dash(e.QQ) + " · " + dash(e.Result) + " · " + dash(e.Zone)
		if n := count[k]; n > 1 {
			note += fmt.Sprintf(" · 共 %d 次", n)
		}

		rows = append(rows, ymmRow{
			QQ:        e.QQ,
			CardKey:   e.Card,
			Note:      note,
			Status:    e.Action,
			CreatedAt: e.TS,
			Level:     -1,
			VIP:       -1,
			Hero:      -1,
			Hist:      true,
		})
	}
	return rows
}

func ymmFilterLabels() []string {
	out := make([]string, 0, len(ymmFilters))
	for _, f := range ymmFilters {
		out = append(out, f.Label)
	}
	return out
}

// ymmPageState 是页面上跨回调共享的状态。
type ymmPageState struct {
	tv          *walk.TableView
	model       *ymmModel
	label       *walk.Label
	logBox      *walk.TextEdit
	chkProxy    *walk.CheckBox
	proxyStatus *walk.Label
	filter      *walk.ComboBox
	search      *walk.LineEdit

	api *cardClient
	app *guiApp

	btnSearch *walk.PushButton

	mu   sync.Mutex
	rows []ymmRow
	busy bool

	sortCol  int  // 1=创建时间 2=激活时间
	sortAsc  bool // true 升序
	histMode bool // 当前展示的是本机历史（表头改叫「操作时间」）

	btnSource *walk.PushButton
	btnVerify *walk.PushButton
	btnSelfOp *walk.PushButton // ③自研上号（仅本页）
	source    string           // online / dav
	histMenu  *walk.Menu
}

// syncContextMenu 只在历史模式下挂右键删除菜单，其它模式清掉。
func (s *ymmPageState) syncContextMenu() {
	if s.tv == nil {
		return
	}
	if !s.histMode {
		s.tv.SetContextMenu(nil)
		return
	}
	if s.histMenu == nil {
		menu, err := walk.NewMenu()
		if err != nil {
			return
		}
		act := walk.NewAction()
		_ = act.SetText("删除这条历史")
		act.Triggered().Attach(func() {
			if s.app != nil {
				s.app.ymmDeleteSelectedHistory()
			}
		})
		_ = menu.Actions().Add(act)
		s.histMenu = menu
	}
	s.tv.SetContextMenu(s.histMenu)
}

// isDump 当前是否用离线（dav）数据源。
func (s *ymmPageState) isDump() bool { return s.source == srcDav }

func (s *ymmPageState) updateSourceBtn() {
	if s.btnSource != nil {
		if s.isDump() {
			s.btnSource.SetText("数据源：离线")
		} else {
			s.btnSource.SetText("数据源：在线")
		}
	}
	// 查封依赖卡密平台的 ck（CF 登录态），离线数据里没有，所以离线时直接禁用
	if s.btnVerify != nil {
		if s.isDump() {
			s.btnVerify.SetText("①查封")
			s.btnVerify.SetEnabled(false)
		} else {
			s.btnVerify.SetText("①查封")
			s.btnVerify.SetEnabled(true)
		}
	}
}

// toggleSource 在线 / 离线 切换，记忆到配置并立即刷新列表。
func (s *ymmPageState) toggleSource() {
	if s.isDump() {
		s.source = srcOnline
	} else {
		s.source = srcDav
	}
	if s.app != nil && s.app.cfg != nil {
		s.app.cfg.DataSource = s.source
		_ = s.app.cfg.save()
	}
	s.updateSourceBtn()
	if s.isDump() {
		s.logf("数据源切换为离线（dav 打包数据）")
	} else {
		s.logf("数据源切换为在线（卡密平台接口）")
	}
	s.refresh()
}

// applySort 按当前排序设置重排列表并刷新表头箭头。
func (s *ymmPageState) applySort() {
	if s.model == nil {
		return
	}
	col := s.sortCol
	if col < 0 || col > 6 {
		col = 5 // 默认创建时间
	}
	asc := s.sortAsc

	s.mu.Lock()
	rows := make([]ymmRow, len(s.rows))
	copy(rows, s.rows)
	s.mu.Unlock()

	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		var less, greater bool
		switch col {
		case 0:
			less, greater = a.Note < b.Note, a.Note > b.Note
		case 1:
			less, greater = a.Level < b.Level, a.Level > b.Level
		case 2:
			less, greater = a.VIP < b.VIP, a.VIP > b.VIP
		case 3:
			less, greater = a.Hero < b.Hero, a.Hero > b.Hero
		case 4:
			less, greater = a.WarZone < b.WarZone, a.WarZone > b.WarZone
		case 6:
			less, greater = a.ActivatedAt < b.ActivatedAt, a.ActivatedAt > b.ActivatedAt
		default:
			less, greater = a.CreatedAt < b.CreatedAt, a.CreatedAt > b.CreatedAt
		}
		if asc {
			return less
		}
		return greater
	})

	s.mu.Lock()
	s.rows = rows
	s.mu.Unlock()
	s.model.rows = rows
	s.model.PublishRowsReset()
	s.updateHeaders()
}

// ymmColumnTitles 是表头基名，顺序必须和 TableView.Columns 一致。
var ymmColumnTitles = []string{"备注", "等级", "VIP", "英雄级", "战区", "创建时间", "激活时间"}

// updateHeaders 在当前排序列的表头标 ▲/▼。
func (s *ymmPageState) updateHeaders() {
	if s.tv == nil {
		return
	}
	cols := s.tv.Columns()
	if cols == nil || cols.Len() < len(ymmColumnTitles) {
		return
	}
	for i, base := range ymmColumnTitles {
		if i == 5 && s.histMode {
			base = "最后操作时间"
		}
		title := base
		if s.sortCol == i {
			if s.sortAsc {
				title += " ▲"
			} else {
				title += " ▼"
			}
		}
		_ = cols.At(i).SetTitle(title)
	}
}

// onColumnClick 点列头切换升/降序。
func (s *ymmPageState) onColumnClick(col int) {
	if col < 0 || col >= len(ymmColumnTitles) {
		return
	}
	if s.sortCol == col {
		s.sortAsc = !s.sortAsc
	} else {
		s.sortCol = col
		s.sortAsc = false
	}
	s.applySort()
}

func (s *ymmPageState) logf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	// 精简模式没有日志框，信息仍然要落到 run.log
	if s.app != nil && s.app.runLog != nil {
		s.app.runLog.line("%s", line)
	}
	if s.logBox == nil {
		return
	}
	s.logBox.AppendText(line + "\r\n")
}

func (s *ymmPageState) setLabel(text string) {
	if s.label == nil {
		return
	}
	s.label.SetTextColor(walk.RGB(120, 120, 120))
	s.label.SetText(text)
}

// setErr 底部状态栏显示错误（红色）。
func (s *ymmPageState) setErr(text string) {
	if s.label == nil {
		return
	}
	s.label.SetTextColor(walk.RGB(190, 40, 40))
	s.label.SetText(text)
}

// ymmErrText 把接口错误压成一句，避免底部被长串 URL 撑爆。
func ymmErrText(err error) string {
	s := err.Error()
	switch {
	case strings.Contains(s, "context deadline exceeded"):
		return "连接超时，服务器或代理无响应"
	case strings.Contains(s, "代理未就绪"):
		return "代理未就绪，点搜索重试"
	case strings.Contains(s, "connection refused"):
		return "连接被拒绝"
	case strings.Contains(s, "no such host"):
		return "域名解析失败"
	}
	if r := []rune(s); len(r) > 60 {
		s = string(r[:60]) + "…"
	}
	return s
}

// snapshot 返回当前列表快照（离线上号要用）。
func (s *ymmPageState) snapshot() []ymmRow {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ymmRow, len(s.rows))
	copy(out, s.rows)
	return out
}

// refresh 按当前筛选条件重新拉列表。
func (s *ymmPageState) refresh() {
	if s.api == nil || s.model == nil {
		return
	}
	// 口令写死，未登录时由下面的 goroutine 先登录再拉列表

	idx := 0
	if s.filter != nil && s.filter.CurrentIndex() >= 0 {
		idx = s.filter.CurrentIndex()
	}
	if idx >= len(ymmFilters) {
		idx = 0
	}
	f := ymmFilters[idx]

	// 历史列表固定按操作时间倒序：最新的一次操作永远在最上面
	s.histMode = f.Status == "history"
	if s.histMode {
		s.sortCol = 5
		s.sortAsc = false
	}
	s.syncContextMenu()

	keyword := ""
	if s.search != nil {
		keyword = strings.TrimSpace(s.search.Text())
	}

	s.mu.Lock()
	if s.busy {
		s.mu.Unlock()
		return
	}
	s.busy = true
	s.mu.Unlock()

	s.setLabel("查询中…")
	if s.btnSearch != nil {
		s.btnSearch.SetEnabled(false)
	}
	app := s.app
	go func() {
		// 无论正常结束还是中途 panic，都把搜索按钮恢复
		defer func() {
			if app != nil {
				app.ui(func() {
					if s.btnSearch != nil {
						s.btnSearch.SetEnabled(true)
					}
				})
			}
		}()
		if app != nil {
			app.ui(app.spinBegin)
			defer app.ui(app.spinEnd)
		}
		// 本机历史：不联网，直接读本地 history.jsonl
		if f.Status == "history" {
			dir := ""
			if app != nil {
				dir = app.dir
			}
			rows := ymmHistoryRows(dir, keyword)
			s.mu.Lock()
			s.busy = false
			s.mu.Unlock()

			apply := func() {
				s.mu.Lock()
				s.rows = rows
				s.mu.Unlock()
				s.applySort()
				s.setLabel(fmt.Sprintf("本机历史：%d 条", len(rows)))
			}
			if app != nil {
				app.ui(apply)
			} else {
				apply()
			}
			return
		}

		notify := func(format string, args ...any) {
			if app == nil {
				return
			}
			msg := fmt.Sprintf(format, args...)
			app.ui(func() { s.logf("%s", msg) })
		}

		// 取一条真能打到服务器的代理 IP，死 IP 直接换（最多 3 条）
		ensureProxy := func() bool {
			for try := 1; try <= 3; try++ {
				if e := ymmProxy.warmup(); e != nil {
					notify("取代理失败(%d/3): %v", try, e)
					time.Sleep(time.Second)
					continue
				}
				if s.api.proxyProbe(4 * time.Second) {
					return true
				}
				notify("代理 %s 不通，换下一条(%d/3)", ymmProxy.addrNow(), try)
				ymmProxy.drop()
				time.Sleep(300 * time.Millisecond)
			}
			return false
		}

		online := s.source != srcDav

		// 离线数据源只读 dav 打包数据，不需要代理和登录
		if online {
			// 未开代理且直连 3s 不通 → 自动切代理
			if !ymmProxy.Enabled() && !s.api.probeDirect(3*time.Second) {
				ymmProxy.SetEnabled(true)
				if app != nil {
					app.ui(func() {
						if s.chkProxy != nil {
							s.chkProxy.SetChecked(true)
						}
						s.logf("直连 3s 不通，已自动开启代理")
					})
				}
			}
			// 代理开着但还没有 IP（含默认开启）：先备好一条可用的
			if ymmProxy.Enabled() && ymmProxy.addrNow() == "" {
				if !ensureProxy() {
					notify("代理均不可用，稍后点搜索重试")
				}
			}
			notify("网络：%s", proxyDesc())
		}

		var err error
		if online && !s.api.loggedIn() {
			err = s.api.login(adminPass())
			// 代理 IP 中途失效：换一条再登一次
			if err != nil && ymmProxy.Enabled() {
				ymmProxy.drop()
				if e := ymmProxy.warmup(); e == nil && s.api.proxyProbe(4*time.Second) {
					err = s.api.login(adminPass())
				}
			}
		}
		var rows []ymmRow
		var total int
		if err == nil {
			rows, total, err = s.fetch(f, keyword)
			// 列表请求失败也换一条 IP 再试
			if err != nil && online && ymmProxy.Enabled() {
				ymmProxy.drop()
				if e := ymmProxy.warmup(); e == nil && s.api.proxyProbe(4*time.Second) {
					rows, total, err = s.fetch(f, keyword)
				}
			}
		}
		s.mu.Lock()
		s.busy = false
		s.mu.Unlock()

		apply := func() {
			if err != nil {
				s.setErr("查询失败：" + ymmErrText(err))
				s.logf("查询失败: %v", err)
				return
			}
			s.mu.Lock()
			s.rows = rows
			s.mu.Unlock()
			s.applySort()

			if s.isDump() {
				s.setLabel(fmt.Sprintf("%s：%d 条（离线库内 %d）", f.Label, len(rows), total))
				return
			}
			role := s.api.getRole()
			if role == "" {
				role = "unknown"
			}
			if total > len(rows) {
				s.setLabel(fmt.Sprintf("%s：%d 条（库内 %d）· %s", f.Label, len(rows), total, role))
			} else {
				s.setLabel(fmt.Sprintf("%s：%d 条 · %s", f.Label, len(rows), role))
			}
		}
		if app != nil {
			app.ui(apply)
		} else {
			apply()
		}
	}()
}

// ymmRowMatch 本地模糊匹配：备注（含区服/战区）、QQ、卡密、编号、状态中文名都能搜。
func ymmRowMatch(r ymmRow, kw string) bool {
	kw = strings.ToLower(strings.TrimSpace(kw))
	if kw == "" {
		return true
	}
	for _, f := range []string{
		r.Note, r.QQ, r.CardKey, r.QueryCode, r.Area, r.WarZone,
		r.Status, ymmStatusLabel(r.Status),
	} {
		if strings.Contains(strings.ToLower(f), kw) {
			return true
		}
	}
	return false
}

// ymmRowFromItem 把接口/dump 的卡密条目转成列表行。
func ymmRowFromItem(it cardItem) ymmRow {
	row := ymmRow{
		QQ:          it.QQ,
		CardKey:     it.CardKey,
		QueryCode:   it.QueryCode,
		Note:        it.Note,
		Status:      it.Status,
		CreatedAt:   it.CreatedAt,
		ActivatedAt: it.ActivatedAt,
	}
	row.ymmParseNote()
	// 未使用的卡密 note 是空的（要等上号器 redeem 才回写），用编号占位便于区分
	if row.Note == "" {
		row.Note = it.QueryCode
	}
	return row
}

// fetch 分页拉取并按 Since 过滤，最多 10 页（2000 条）。
// 关键字不走服务端 q（它只认卡密/QQ，搜不到备注里的区服战区），改成本地模糊匹配。
func (s *ymmPageState) fetch(f ymmFilter, keyword string) ([]ymmRow, int, error) {
	// 离线数据源：只读 dav 上的 cards.jsonl（ini 等上号时再按需取）
	if s.source == srcDav {
		path, err := ensureDumpCards(nil)
		if err != nil {
			return nil, 0, err
		}
		// 顺手把 zip 目录读到内存，上号时少一趟往返
		go func() { _, _, _ = dumpSource() }()
		return dumpRows(path, f, keyword)
	}

	const pageSize = 200
	var out []ymmRow
	total := 0
	kw := strings.ToLower(strings.TrimSpace(keyword))

	for page := 1; page <= 10; page++ {
		items, t, err := s.api.listCards(cardQuery{
			Status:   f.Status,
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			return nil, 0, err
		}
		total = t
		for _, it := range items {
			row := ymmRowFromItem(it)
			if !ymmRowMatch(row, kw) {
				continue
			}
			out = append(out, row)
		}
		if len(items) < pageSize {
			break
		}
	}

	if f.Since > 0 {
		cut := time.Now().Add(-f.Since).UnixMilli()
		kept := out[:0]
		for _, r := range out {
			if r.ActivatedAt >= cut {
				kept = append(kept, r)
			}
		}
		out = kept
	}

	// 默认按创建时间倒序（分页是逐页拉的，跨页排序只能在这里做）
	sort.SliceStable(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	return out, total, nil
}

// selected 返回当前选中行。
func (s *ymmPageState) selected() (ymmRow, bool) {
	i := -1
	if s.tv != nil {
		i = s.tv.CurrentIndex()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if i < 0 || i >= len(s.rows) {
		return ymmRow{}, false
	}
	return s.rows[i], true
}
