package main

// yuanmeng 页的控件与按钮回调：数据来自卡密平台接口，不再读本地 dump。

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// ymmPage 构建 yuanmeng 页。
func (a *guiApp) ymmPage() Widget {
	a.ymm.app = a
	a.ymm.model = &ymmModel{}
	if a.ymm.api == nil {
		a.ymm.api = newCardClient()
	}

	dim := walk.RGB(110, 110, 110)
	sp, btnH, rowH := 8, 32, 26
	if a.slim {
		sp, btnH, rowH = 6, 28, 24
	}
	// 默认按创建时间倒序（列索引 5）
	a.ymm.sortCol = 5
	a.ymm.sortAsc = false

	ops := []Widget{
		PushButton{AssignTo: &a.ymm.btnVerify, Text: "①查封", MinSize: Size{Width: 90, Height: btnH}, OnClicked: a.ymmQuerySelected},
		PushButton{Text: "②上号", MinSize: Size{Width: 90, Height: btnH}, OnClicked: a.ymmLaunchViaLauncherClick},
		PushButton{AssignTo: &a.ymm.btnSelfOp, Text: "③自研上号", MinSize: Size{Width: 110, Height: btnH}, OnClicked: a.ymmSelfOpClick},
		PushButton{Text: "复制卡密", MinSize: Size{Width: 90, Height: btnH}, OnClicked: a.ymmCopyCard},
		PushButton{AssignTo: &a.ymm.btnSource, Text: "数据源：在线", MinSize: Size{Width: 110, Height: btnH}, OnClicked: func() { a.ymm.toggleSource() }},
		CheckBox{
			AssignTo: &a.ymm.chkProxy,
			Text:     "使用代理",
			OnCheckedChanged: func() {
				on := a.ymm.chkProxy.Checked()
				ymmProxy.SetEnabled(on)
				// 记住开关状态，下次启动不用再手动勾
				if a.cfg != nil && a.cfg.proxyDefault() != on {
					a.cfg.setProxy(on)
					_ = a.cfg.save()
				}
				a.ymm.updateProxyStatus()
				go func(on bool) {
					desc := "直连"
					if on {
						if err := ymmProxy.warmup(); err != nil {
							desc = fmt.Sprintf("代理获取失败：%v", err)
						} else {
							desc = proxyDesc()
						}
					}
					a.ui(func() {
						a.ymm.updateProxyStatus()
						a.ymm.logf("网络：%s", desc)
					})
				}(on)
			},
		},
		Label{
			AssignTo:  &a.ymm.proxyStatus,
			TextColor: walk.RGB(30, 110, 200),
			Font:      Font{PointSize: 9},
			MinSize:   Size{Width: 220},
			Visible:   false,
		},
		HSpacer{},
	}
	statusKids := []Widget{
		Label{AssignTo: &a.ymm.label, Text: "", TextColor: walk.RGB(120, 120, 120), Font: Font{PointSize: 9}},
		HSpacer{},
	}
	// 精简模式没有顶部栏：把「切换上号器」挂到这一页，否则退不回第一界面
	if a.slim {
		ops = append(ops, PushButton{Text: "切换上号器", MinSize: Size{Width: 110, Height: btnH}, OnClicked: a.reconfigure})
	}

	children := []Widget{
		// 第一行：筛选 / 搜索 / 刷新
		Composite{
			Layout: HBox{Spacing: sp},
			Children: []Widget{
				Label{Text: "状态", TextColor: dim, Font: Font{PointSize: 9}},
				ComboBox{
					AssignTo:              &a.ymm.filter,
					Model:                 ymmFilterLabels(),
					CurrentIndex:          0,
					MinSize:               Size{Width: 130},
					OnCurrentIndexChanged: func() { a.ymm.refresh() },
				},
				LineEdit{
					AssignTo:  &a.ymm.search,
					MinSize:   Size{Width: 200},
					CueBanner: "QQ / 卡密 / 编号 / 备注",
					OnKeyDown: func(key walk.Key) {
						if key == walk.KeyReturn {
							a.ymm.refresh()
						}
					},
				},
				PushButton{AssignTo: &a.ymm.btnSearch, Text: "搜索", MinSize: Size{Width: 70, Height: btnH}, OnClicked: func() { a.ymm.refresh() }},
				HSpacer{},
				// 右上角转圈：有接口请求时转，空闲时清空
				Label{
					AssignTo:  &a.spinLabel,
					Text:      "",
					TextColor: walk.RGB(30, 110, 200),
					Font:      Font{Family: "Segoe UI Symbol", PointSize: 12},
					MinSize:   Size{Width: 22},
				},
			},
		},
		// 第二行：操作
		Composite{
			Layout:   HBox{Spacing: sp},
			Children: ops,
		},
		TableView{
			AssignTo:            &a.ymm.tv,
			AlternatingRowBG:    true,
			CustomRowHeight:     rowH,
			LastColumnStretched: false,
			Columns: []TableViewColumn{
				{Title: "备注", Width: 365},
				{Title: "等级", Width: 55},
				{Title: "VIP", Width: 50},
				{Title: "英雄级", Width: 65},
				{Title: "战区", Width: 95},
				{Title: "创建时间 ▼", Width: 150},
				{Title: "激活时间", Width: 135},
			},
			Model: a.ymm.model,
		},
		// 底部状态行：查询结果 / 错误信息都放这里
		Composite{
			Layout:   HBox{Spacing: 6},
			Children: statusKids,
		},
	}

	// 精简模式：左上角品牌名 + 版本号
	if a.slim {
		children = append([]Widget{Composite{
			Layout: HBox{Spacing: 8},
			Children: []Widget{
				Label{
					Text:      "CG租号",
					TextColor: walk.RGB(24, 96, 186),
					Font:      Font{Family: "Microsoft YaHei", PointSize: 14, Bold: true},
				},
				HSpacer{},
				Label{Text: "v" + toolVersion, TextColor: walk.RGB(150, 150, 150), Font: Font{PointSize: 8}},
			},
		}}, children...)
	}

	if !a.slim {
		children = append(children, GroupBox{
			Title:  "yuanmeng 日志",
			Layout: VBox{},
			Children: []Widget{
				TextEdit{AssignTo: &a.ymm.logBox, ReadOnly: true, VScroll: true, MinSize: Size{Height: 140}},
			},
		})
	}

	return Composite{Layout: VBox{}, Children: children}
}

// updateProxyStatus 在勾选框后显示当前代理及剩余时效。
func (s *ymmPageState) updateProxyStatus() {
	if s.proxyStatus == nil {
		return
	}
	enabled, addr, left := ymmProxy.status()
	if !enabled {
		s.proxyStatus.SetText("")
		s.proxyStatus.SetVisible(false)
		return
	}
	s.proxyStatus.SetVisible(true)
	if addr == "" {
		s.proxyStatus.SetText("代理 IP：待获取")
		return
	}
	seconds := int(left / time.Second)
	if seconds < 0 {
		seconds = 0
	}
	s.proxyStatus.SetText(fmt.Sprintf("IP：%s｜剩余：%d 秒", addr, seconds))
}

// proxyStatusLoop 每秒刷新倒计时，窗口关闭后自动退出。
func (a *guiApp) proxyStatusLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if a.closing.Load() {
			return
		}
		a.ui(a.ymm.updateProxyStatus)
	}
}

// ymmQuerySelected 查封按钮：不分数据源，都用行里的卡密调验号接口拿实时封号状态。
// 离线（dump）列表的行同样带真卡密（来自 cards.jsonl），不需要先切到在线数据源。
func (a *guiApp) ymmQuerySelected() {
	// 离线数据源没有 ck，查封做不了（按钮已禁用，这里兜底防其它入口触发）
	if a.ymm.isDump() {
		a.ymm.logf("离线数据源暂不支持查封，请点「数据源：离线」切回在线再查")
		return
	}
	row, ok := a.ymm.selected()
	if !ok {
		a.ymm.logf("请先选中一行")
		return
	}
	if row.CardKey == "" {
		a.ymm.logf("该行没有卡密，无法查封")
		return
	}
	a.ymmVerifyRow(row, a.ymm.isDump())
}

// ymmVerifyRow 用卡密查封号。
// 在线/离线走同一条实时查询路径，区别只在失败兜底：
// 离线还会补一条免登录的实名/人脸检查，并退到上次成功的本地快照。
func (a *guiApp) ymmVerifyRow(row ymmRow, offline bool) {
	qq := strings.TrimSpace(row.QQ)
	a.ymm.logf("查封 %s（%s）", dash(qq), row.CardKey)

	go func() {
		a.ui(a.spinBegin)
		defer a.ui(a.spinEnd)

		// 验号接口按 IP 限频/封禁：直连不通就先切代理，别等接口回 rate_limited
		if !ymmProxy.Enabled() && !a.ymm.api.probeDirect(3*time.Second) {
			ymmProxy.SetEnabled(true)
			_ = ymmProxy.warmup()
			a.ui(func() {
				if a.ymm.chkProxy != nil {
					a.ymm.chkProxy.SetChecked(true)
				}
				a.ymm.logf("直连不通，改用代理查询")
			})
		}
		a.ui(func() { a.ymm.logf("网络：%s", proxyDesc()) })

		st := queryVerifyStatus(a.cfmHost(), row.CardKey)

		// 离线兜底：实时查不通（或卡密失效）时用上次快照 + 免登录的实名/人脸
		if offline && (st.Err != "" || (!st.OK && !st.RateLimited)) {
			rail := railAntiAddictText(qq)
			snap, hasSnap := getVerifySnap(qq)
			reason := firstNonEmpty(st.Err, st.Message, "接口没返回封号信息")
			a.ui(func() {
				if hasSnap {
					at := time.UnixMilli(snap.At)
					a.ymm.logf("实时查封失败（%s），改用快照 %s", gsClip(reason, 40), at.Format("01-02 15:04"))
					a.ymmVerifyDialog(verifyDialog{card: row.CardKey, st: snap.snapStatus(), offline: true, snapAt: at, rail: rail})
					return
				}
				a.ymm.logf("查封失败：%s", reason)
				a.ymmVerifyDialog(verifyDialog{card: row.CardKey, st: st, offline: true, rail: rail})
			})
			return
		}

		if st.Err != "" {
			a.ui(func() { a.ymm.logf("失败: %s", st.Err) })
			return
		}
		if st.RateLimited {
			a.ui(func() {
				a.ymm.logf("查询过于频繁，%s 后再试", st.RetryAfter.Round(time.Second))
				a.ymmVerifyDialog(verifyDialog{card: row.CardKey, st: st})
			})
			return
		}
		if !st.OK {
			msg := firstNonEmpty(st.Message, "接口返回 ok=false")
			expired := strings.Contains(msg, "失效") || strings.Contains(msg, "过期")
			a.ui(func() {
				if expired {
					a.ymm.logf("查询失败：%s —— 卡密已失效，接口不再返回封号信息，请用「用上号器(离线)」进游戏验号", msg)
				} else {
					a.ymm.logf("查询失败：%s", msg)
				}
				a.ymmVerifyDialog(verifyDialog{card: row.CardKey, st: st})
			})
			return
		}

		saveVerifySnap(row.CardKey, st)
		a.ui(func() {
			a.ymm.logf("QQ=%s 状态=%s 封号=%s 冻结=%s 在线=%s 人脸=%s 令牌=%s 区服=%s",
				st.QQ, st.Status, st.Fenghao, st.Dongjie, st.Zaixian, st.Renlian, st.Lingpai, st.Zone)
			a.ymmVerifyDialog(verifyDialog{card: row.CardKey, st: st})
		})
	}()
}

// verifyDialog 是一次查封弹窗的数据：在线实时查询、离线快照两种来源。
type verifyDialog struct {
	card    string
	st      verifyStatus
	offline bool
	snapAt  time.Time // 非零 = 这是离线快照
	rail    string    // 免登录的实名/人脸实时检查结果
}

// ymmVerifyDialog 弹框显示验号结果（状态 / 封号 / 冻结 / 在线 / 人脸 / 令牌 / 区服）。
func (a *guiApp) ymmVerifyDialog(v verifyDialog) {
	card, st := v.card, v.st
	var dlg *walk.Dialog
	var edit *walk.TextEdit

	var b strings.Builder
	if v.offline {
		b.WriteString("【离线数据源】dump 里只有游戏登录票据，没有查封要的 ck，查不到实时封号。\r\n\r\n")
		if v.snapAt.IsZero() {
			b.WriteString("本地还没有这个 QQ 的查封快照。\r\n")
			b.WriteString("想看封号：切到「数据源：在线」查封一次，结果会自动存成本地快照。\r\n\r\n")
		} else {
			fmt.Fprintf(&b, "以下为上次在线查封的快照（%s）：\r\n\r\n", v.snapAt.Format("2006-01-02 15:04"))
		}
	} else {
		if !st.OK {
			fmt.Fprintf(&b, "查询失败：%s\r\n", dash(st.Message))
			b.WriteString("（接口未返回封号详情，卡密可能已失效）\r\n\r\n")
		} else if st.RateLimited {
			b.WriteString("查询过于频繁，请稍后再试\r\n\r\n")
		}
	}
	fmt.Fprintf(&b, "卡密：%s\r\n", dash(card))
	fmt.Fprintf(&b, "QQ：%s\r\n", dash(st.QQ))
	fmt.Fprintf(&b, "区服：%s\r\n", dash(st.Zone))
	fmt.Fprintf(&b, "状态：%s\r\n", dash(st.Status))
	fmt.Fprintf(&b, "封号：%s\r\n", dash(st.Fenghao))
	fmt.Fprintf(&b, "冻结：%s\r\n", dash(st.Dongjie))
	fmt.Fprintf(&b, "在线：%s\r\n", dash(st.Zaixian))
	fmt.Fprintf(&b, "人脸：%s\r\n", dash(st.Renlian))
	fmt.Fprintf(&b, "令牌：%s\r\n", dash(st.Lingpai))
	if st.Nickname != "" {
		fmt.Fprintf(&b, "昵称：%s\r\n", st.Nickname)
	}
	if st.Level != "" {
		fmt.Fprintf(&b, "等级：%s\r\n", st.Level)
	}
	if !v.offline && !st.At.IsZero() {
		fmt.Fprintf(&b, "查询时间：%s\r\n", st.At.Format("2006-01-02 15:04:05"))
	}
	if v.rail != "" {
		fmt.Fprintf(&b, "\r\n实名/人脸（实时，免登录）：%s\r\n", v.rail)
	}

	d := Dialog{
		AssignTo: &dlg,
		Title:    "①查封",
		Size:     Size{Width: 400, Height: 430},
		Layout:   VBox{},
		Children: []Widget{
			TextEdit{AssignTo: &edit, Text: b.String(), ReadOnly: true, VScroll: true, MinSize: Size{Width: 360, Height: 340}},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{Text: "复制", MinSize: Size{Width: 96, Height: 30}, OnClicked: func() {
						if err := copyToClipboard(edit.Text()); err != nil {
							walk.MsgBox(dlg, "提示", "复制失败: "+err.Error(), walk.MsgBoxIconWarning)
						}
					}},
					PushButton{Text: "关闭", MinSize: Size{Width: 96, Height: 30}, OnClicked: func() { dlg.Accept() }},
				},
			},
		},
	}
	if err := d.Create(a.mw); err != nil {
		a.ymm.logf("弹窗失败: %v", err)
		return
	}
	dlg.Run()
}

// ymmDeleteSelectedHistory 删除历史列表里选中的那条上号记录。
func (a *guiApp) ymmDeleteSelectedHistory() {
	row, ok := a.ymm.selected()
	if !ok || row.CardKey == "" {
		return
	}

	msg := fmt.Sprintf("删除这条上号记录？\n\nQQ：%s\n时间：%s", dash(row.QQ), row.createdText())
	if walk.MsgBox(a.mw, "确认删除", msg, walk.MsgBoxYesNo|walk.MsgBoxIconQuestion) != walk.DlgCmdYes {
		return
	}
	if e := deleteHistory(a.dir, row.CardKey, row.CreatedAt); e != nil {
		walk.MsgBox(a.mw, "删除失败", e.Error(), walk.MsgBoxIconWarning)
		return
	}
	a.ymm.logf("已删除历史：%s", dash(row.QQ))
	a.ymm.refresh()
}

// ymmCopyCard 复制选中行的卡密。
func (a *guiApp) ymmCopyCard() {
	row, ok := a.ymm.selected()
	if !ok {
		a.ymm.logf("请先选中一行")
		return
	}
	if row.CardKey == "" {
		a.ymm.logf("该账号没有卡密")
		return
	}
	if err := copyToClipboard(row.CardKey); err != nil {
		a.ymm.logf("复制失败: %v", err)
		return
	}
	a.ymm.logf("已复制卡密 %s", row.CardKey)
}

// ymmLaunchViaLauncherClick 离线上号：CG 冒充卡密服务器，让上号器自己上号。
func (a *guiApp) ymmLaunchViaLauncherClick() {
	row, ok := a.ymm.selected()
	if !ok {
		a.ymm.logf("请先选中一行")
		return
	}
	a.ymmLaunchOffline(row)
}

// ymmLaunchOffline 用指定行启动离线上号（ini 由本地服务器按需从接口拉）。
func (a *guiApp) ymmLaunchOffline(row ymmRow) {
	if row.QQ == "" || row.CardKey == "" {
		a.ymm.logf("该行缺少 QQ 或卡密")
		return
	}

	launcher := a.loader
	if launcher == "" {
		cands := []string{
			filepath.Join(exeDir(), "..", "资源", "上号器", "yuanmeng", "CF上号器_v1.17.exe"),
			filepath.Join(exeDir(), "CF上号器_v1.17.exe"),
		}
		for _, c := range cands {
			if _, err := os.Stat(c); err == nil {
				launcher = filepath.Clean(c)
				break
			}
		}
	}
	if launcher == "" {
		a.ymm.logf("找不到上号器 exe，请把 CF上号器_v1.17.exe 放到 CG 同级或 资源\\上号器\\yuanmeng\\")
		return
	}
	if _, err := os.Stat(launcher); err != nil {
		a.ymm.logf("上号器不存在: %s", launcher)
		return
	}

	offline, err := ymmBuildOffline(a.ymm.api, a.ymm.snapshot(), row.CardKey, a.ymm.isDump())
	if err != nil {
		a.ymm.logf("准备失败: %v", err)
		return
	}

	// 与回放流程一致：config.ini 里写占位卡密，游戏目录由上号器界面选
	// 离线上号额外把 auto_zone 设成 1（自动选区）
	dir := filepath.Dir(launcher)
	if _, changed, e := writeOfflineConfig(dir); e != nil {
		a.ymm.logf("写 config.ini 失败: %v", e)
	} else if changed {
		a.ymm.logf("已写入 config.ini（card=%s, auto_zone=1，首次需在上号器里选游戏目录）", placeholderCard)
	} else {
		a.ymm.logf("config.ini 就绪（card=%s, auto_zone=1，保留已选游戏目录）", placeholderCard)
	}

	a.ymm.logf("离线上号：%s（卡密 %s）", launcher, row.CardKey)
	a.ymm.logf("本地服务器 %s（host.dat 会临时切到本地，上号器退出后自动还原）", ymmLocalPrefix)

	go func() {
		a.ui(a.spinBegin)
		defer a.ui(a.spinEnd)

		// 取 ini 要联网，先把公用代理备好，避免直连 IP 被限频
		if ymmProxy.Enabled() {
			if e := ymmProxy.warmup(); e != nil {
				a.ui(func() { a.ymm.logf("代理取 IP 失败: %v", e) })
			}
		}

		appendHistory(a.dir, histEntry{
			TS:     time.Now().UnixMilli(),
			Action: "offline",
			QQ:     row.QQ,
			Card:   row.CardKey,
			Note:   row.Note,
			Zone:   ymmZoneFromNote(row.Note),
			Result: "已启动",
		})
		if err := ymmLaunchViaLauncher(launcher, offline); err != nil {
			a.ui(func() { a.ymm.logf("失败: %v", err) })
			return
		}
		a.ui(func() { a.ymm.logf("上号器已退出，host.dat 已还原") })
	}()
}

// ymmSelfOpClick ③自研上号：只用这一行的票据自己换 openid 再启动游戏，
// 不经过真上号器。按钮只挂在第二界面，第一界面的「自研上号」不受影响。
func (a *guiApp) ymmSelfOpClick() {
	row, ok := a.ymm.selected()
	if !ok {
		a.ymm.setErr("③自研上号：请先在列表里选中一行")
		walk.MsgBox(a.mw, "③自研上号", "请先在列表里选中一行", walk.MsgBoxIconWarning)
		return
	}
	if row.QQ == "" || row.CardKey == "" {
		a.ymm.setErr("③自研上号：该行缺少 QQ 或卡密")
		walk.MsgBox(a.mw, "③自研上号", "该行缺少 QQ 或卡密", walk.MsgBoxIconWarning)
		return
	}
	// 游戏目录：配置里没有或已失效就弹框选，选完记住（跟离线上号器一样）
	dir, err := a.selfOpGameDir()
	if err != nil {
		a.ymm.setErr("③自研上号：" + err.Error())
		return
	}
	if !a.selfBusy.CompareAndSwap(false, true) {
		a.ymm.setErr("③自研上号：上一个任务还在执行，稍后再试")
		return
	}
	if a.ymm.btnSelfOp != nil {
		a.ymm.btnSelfOp.SetEnabled(false)
	}
	a.ymm.setLabel("③自研上号：游戏目录 " + dir)
	go a.ymmSelfOpRun(row, dir)
}

// ymmSelfOpRun 是 ③ 的实际流程：取 ini → 解析票据 → 换 openid → 启动并注入。
// 精简模式没有日志框、也不写 run.log，所以进度和失败原因必须写到页面底部状态行，
// 失败再弹一次框，否则用户看到的就是"点了没反应"。
func (a *guiApp) ymmSelfOpRun(row ymmRow, dir string) {
	logf := func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)
		a.ui(func() { a.ymm.logf("%s", msg) })
	}
	// step 更新底部状态行（同时留一行日志，非精简模式可见）
	step := func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)
		logf("%s", msg)
		a.ui(func() { a.ymm.setLabel("③自研上号：" + msg) })
	}
	// fail 是终止性失败：状态行变红 + 弹框，让精简模式也能看见
	fail := func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)
		logf("失败：%s", msg)
		a.ui(func() {
			a.ymm.setErr("③自研上号失败：" + msg)
			walk.MsgBox(a.mw, "③自研上号失败", msg, walk.MsgBoxIconError)
		})
	}

	defer func() {
		a.selfBusy.Store(false)
		a.ui(func() {
			a.spinEnd()
			if a.ymm.btnSelfOp != nil {
				a.ymm.btnSelfOp.SetEnabled(true)
			}
		})
		if r := recover(); r != nil {
			fail("内部异常：%v", r)
		}
	}()
	a.ui(a.spinBegin)
	step("开始 %s（QQ %s）", row.CardKey, row.QQ)
	step("游戏目录：%s", dir)

	// 用客户自己游戏目录里的 rail_api64.dll（不内置、不从 WeGame 取）
	step("校验游戏目录的 rail_api64.dll…")
	if desc, err := selfRailDLLInfo(dir); err != nil {
		fail("%v", err)
		return
	} else {
		step("rail64 就绪：%s", desc)
	}

	src := "卡密平台"
	if a.ymm.isDump() {
		src = "离线数据源"
	}
	step("取票据（%s）…", src)
	iniText, err := a.ymmSelfIni(row.QQ)
	if err != nil {
		fail("取 ini 失败：%v", err)
		return
	}
	tk, err := parseYmmTicket(iniText, row.QQ)
	if err != nil {
		fail("解析 ini 失败：%v", err)
		return
	}
	step("票据就绪 tgpid=%s machine=%s ct=%s… st=%s…",
		dash(tk.Tgpid), dash(tk.MachineID), shortHex(tk.CT, 10), shortHex(tk.ST, 10))

	// 真上号器的顺序：LoginByQq 换 OP，再拿返回的 game_st/ctt 去启动。
	// 直接拿 ini 里的 NT st 打 GetGameTicket 必然 702（票据不是游戏票）。
	// 真上号器的容错顺序：本地 OP 缓存 → LoginByQq（退避重试）→ CF_Gamest（同 QQ 才认）
	step("取 OP（缓存 → LoginByQq → CF_Gamest）…")
	data, err := a.selfOpLoginWithRetry(tk, row.QQ, step)
	if err != nil {
		fail("换 OP 失败：%v", err)
		return
	}
	step("OP 票据就绪 openid=%s open_key=%s… ctt=%s…",
		data.OpenID, shortHex(data.OpenKey, 8), shortHex(data.CTT, 10))
	if data.CTT == "" {
		step("注意：回包没有 ctt，游戏可能拒绝登录")
	}

	step("启动游戏并注入（用该目录的 rail_api64.dll）…")
	pid, err := selfLaunchGame(dir, data)
	if err != nil {
		fail("启动失败：%v", err)
		return
	}
	appendHistory(a.dir, histEntry{
		TS:     time.Now().UnixMilli(),
		Action: "selfop",
		QQ:     row.QQ,
		Card:   row.CardKey,
		Note:   row.Note,
		Zone:   ymmZoneFromNote(row.Note),
		Result: fmt.Sprintf("已启动 pid=%d", pid),
	})
	step("已启动 crossfire.exe pid=%d", pid)
}

// ymmSelfIni 取某个 QQ 的票据明文：离线读 dump，在线走卡密平台接口。
func (a *guiApp) ymmSelfIni(qq string) (string, error) {
	if a.ymm.isDump() {
		return ensureDumpIni(qq)
	}
	if a.ymm.api == nil {
		return "", fmt.Errorf("卡密平台未初始化")
	}
	if !a.ymm.api.loggedIn() {
		if err := a.ymm.api.login(adminPass()); err != nil {
			return "", fmt.Errorf("登录卡密平台: %w", err)
		}
	}
	// 启用了代理就先备一条可用 IP，避免直连被限频
	if ymmProxy.Enabled() {
		_ = ymmProxy.warmup()
	}
	return a.ymm.api.getData(qq)
}

// cfmHost 返回卡密服务器地址。
func (a *guiApp) cfmHost() string {
	if a.ymm.api != nil {
		return a.ymm.api.getBase()
	}
	return cardSrvDefault
}

// ymmStatusLabel 把状态码转中文（日志用）。
func ymmStatusLabel(s string) string {
	switch s {
	case "unused":
		return "未用"
	case "active":
		return "使用中"
	case "expired":
		return "过期"
	case "revoked":
		return "作废"
	}
	return s
}
