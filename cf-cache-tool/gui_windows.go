//go:build windows

package main

// 图形界面：录制 / 回放两个 Tab，所有操作在一个 exe 内完成。

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// 表格列索引：验证数据固定放在最后一列，前面各列沿用原顺序。
const verifyColumn = 8

type accountModel struct {
	walk.TableModelBase
	rows []accountRow
	cfg  *appConfig
}

// StyleCell 给“状态”列上色：正常绿、封禁红、黑屋橙、过期灰、其他蓝。
func (m *accountModel) StyleCell(style *walk.CellStyle) {
	if m.cfg == nil {
		return
	}
	row := style.Row()
	if row < 0 || row >= len(m.rows) {
		return
	}
	r := m.rows[row]

	// 验号列：请求失败的显示红字，其余保持默认色，不受云端灰显影响
	if style.Col() == verifyColumn {
		if r.VerifyErr != "" {
			style.TextColor = walk.RGB(183, 28, 28)
		}
		return
	}

	// 只在云端、还没下载到本地的行：其它列灰显，提示“点启动会自动下载”
	if style.Col() != 0 {
		if r.Remote {
			style.TextColor = walk.RGB(140, 140, 140)
		}
		return
	}

	switch m.cfg.mark(r.Card) {
	case markBan:
		style.BackgroundColor = walk.RGB(198, 40, 40)
		style.TextColor = walk.RGB(255, 255, 255)
	case markBlack:
		style.BackgroundColor = walk.RGB(239, 108, 0)
		style.TextColor = walk.RGB(255, 255, 255)
	case markExpired:
		style.BackgroundColor = walk.RGB(97, 97, 97)
		style.TextColor = walk.RGB(255, 255, 255)
	case markOther:
		style.BackgroundColor = walk.RGB(21, 101, 192)
		style.TextColor = walk.RGB(255, 255, 255)
	default:
		style.BackgroundColor = walk.RGB(46, 125, 50)
		style.TextColor = walk.RGB(255, 255, 255)
	}
}

func (m *accountModel) RowCount() int { return len(m.rows) }

func (m *accountModel) Value(row, col int) interface{} {
	if row < 0 || row >= len(m.rows) {
		return ""
	}
	r := m.rows[row]
	switch col {
	case 0:
		if m.cfg != nil {
			return markLabel(m.cfg.mark(r.Card))
		}
		return "正常"
	case 1:
		return r.Card
	case 2:
		return r.QQ
	case 3:
		return r.OpenID
	case 4:
		return r.Note
	case 5:
		return r.Zone
	case 6:
		return r.Recorded.Format("01-02 15:04")
	case 7:
		return r.RawNote
	case verifyColumn:
		if r.Verify != "" {
			return r.Verify
		}
		return "-"
	default:
		return ""
	}
}

func (m *accountModel) SetValue(row, col int, value interface{}) error { return nil }

type guiApp struct {
	mw              *walk.MainWindow
	tv              *walk.TableView
	model           *accountModel
	logBox          *walk.TextEdit
	btnRecord       *walk.PushButton
	btnReplay       *walk.PushButton
	btnChoose       *walk.PushButton
	headLabel       *walk.Label
	headSub         *walk.Label
	siteEdit        *walk.LineEdit
	siteNote        *walk.Label
	btnCopySite     *walk.PushButton
	btnOpenSite     *walk.PushButton
	verLabel        *walk.Label
	btnMark         *walk.PushButton
	btnVerify       *walk.PushButton
	btnCloudRefresh *walk.PushButton

	dir      string
	cacheDir string
	logPath  string
	listen   string
	loader   string

	cfg      *appConfig
	profile  *loaderProfile
	firstRun bool

	host        *HostDatManager
	srvMu       sync.Mutex
	srv         *Server
	busy        bool
	runLog      *runLogger
	logErr      error
	closing     atomic.Bool
	loadGen     int64
	uiLogLines  int
	selectTop   atomic.Bool
	cloudIdx    *cloudIndex // 只在界面线程读写
	loaderReady atomic.Bool

	verifying  atomic.Bool // 验证请求进行中
	clickedCol int         // 最近一次点击的列，用于区分双击打开哪个详情
}

func runGUI() (err error) {
	runtime.LockOSThread()
	trace := openStartupTrace()
	defer trace.close()
	defer func() {
		if r := recover(); r != nil {
			trace.printf("panic: %v", r)
			walk.MsgBox(nil, "CG 启动失败", fmt.Sprintf("%v", r), walk.MsgBoxIconError)
			err = fmt.Errorf("%v", r)
		}
	}()
	trace.printf("启动 CG v%s", toolVersion)

	// 首次配置窗口也检查更新：后台执行，进度显示在第一界面上
	go func() {
		defer func() {
			if r := recover(); r != nil {
				trace.printf("检查更新异常: %v", r)
			}
		}()
		emitUpdate("正在检查更新…")
		v, err := checkUpdateCore(func(f string, args ...any) {
			msg := fmt.Sprintf(f, args...)
			trace.printf("%s", msg)
			emitUpdate(msg)
		})
		if err != nil {
			trace.printf("检查更新失败: %v", err)
			emitUpdate("检查更新失败：" + err.Error())
			return
		}
		if v != "" {
			trace.printf("已更新到 v%s（首次界面）", v)
			emitUpdate("已更新到 v" + v + "，完成配置后自动重启")
			return
		}
		emitUpdate("已是最新版本 v" + toolVersion)
	}()

	cleanupOldExe(trace.printf)

	cfg := loadAppConfig()
	if old, moved := migrateWorkDir(cfg); moved {
		trace.printf("安装目录已迁移: %s -> %s", old, cfg.WorkDir)
	}

	firstRun := false
	if cfg == nil || cfg.WorkDir == "" || findProfile(cfg.Profile) == nil {
		trace.printf("需要首次配置")
		sr, act := runSetupDialog(nil, cfg)
		trace.printf("配置窗口返回: act=%d result=%+v", act, sr)
		if act != setupOK {
			trace.printf("用户没有完成配置，退出")
			return nil
		}
		if cfg == nil {
			cfg = &appConfig{}
		}
		cfg.Profile = sr.ProfileID
		cfg.WorkDir = sr.WorkDir
		cfg.CloudName = sr.CloudName
		cfg.AutoCloud = true
		if e := cfg.save(); e != nil {
			walk.MsgBox(nil, "提示", "配置无法保存，下次启动会重新询问：\n"+e.Error(), walk.MsgBoxIconWarning)
		}
		firstRun = true

		// 首配期间已经替换好新版本：直接重启进新版本，不再用旧进程开主界面
		if v := appliedVersion(); v != "" {
			if exe, e := os.Executable(); e == nil {
				trace.printf("配置完成，重启到 v%s", v)
				pendingRestart.Store(exe)
				return nil
			}
		}
	}

	trace.printf("工作目录=%s 类型=%s 卡密=%s", cfg.WorkDir, cfg.Profile, cfg.CloudName)
	a := &guiApp{cfg: cfg, firstRun: firstRun}
	if e := a.init(); e != nil {
		trace.printf("init 失败: %v", e)
		walk.MsgBox(nil, "CG 启动失败", e.Error(), walk.MsgBoxIconError)
		return e
	}
	trace.printf("主窗口创建完成，进入消息循环")
	a.mw.Run()
	trace.printf("主窗口已退出")
	return nil
}

func (a *guiApp) init() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	a.logPath = filepath.Join(exeDir, "run.log")

	if a.cfg == nil {
		a.cfg = &appConfig{Profile: "yuanmeng", WorkDir: exeDir, AutoCloud: true}
	}
	a.cfg.AutoCloud = true // 录制完成即自动上云，界面暂不提供开关
	a.profile = findProfile(a.cfg.Profile)
	if a.profile != nil {
		a.cfg.Profile = a.profile.ID
	}
	a.dir = a.cfg.WorkDir
	if a.dir == "" {
		a.dir = exeDir
	}
	a.cacheDir = filepath.Join(a.dir, "cache")
	a.listen = "127.0.0.1:18800"
	a.host = NewHostDatManager(filepath.Join(a.dir, "host.dat"))
	if a.profile != nil {
		a.loader = filepath.Join(a.dir, a.profile.ExeName)
		if ok, reason := localLoaderReady(a.dir, a.profile); ok {
			a.loaderReady.Store(true)
		} else {
			a.logErr = fmt.Errorf("%s", reason)
			go a.syncLoaderAsync()
		}
	}
	a.model = &accountModel{rows: nil, cfg: a.cfg}

	if rl, e := openRunLogger(a.logPath); e == nil {
		a.runLog = rl
	} else {
		a.logErr = e
	}

	if err := a.createWindow(); err != nil {
		return err
	}
	a.tv.ColumnClicked().Attach(func(col int) { a.clickedCol = col })
	a.tv.ItemActivated().Attach(func() {
		// 双击“验证数据”列看验号详情，其它列沿用原来的原始备注窗口
		if a.clickedCol == verifyColumn {
			a.showVerifyDetail()
			return
		}
		a.showRawNote()
	})
	a.tv.CurrentIndexChanged().Attach(func() { a.updateSelectionState() })
	a.updateSelectionState()
	a.refreshHeader()
	a.mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		// 先标记关闭，后续后台协程的界面调用直接丢弃，避免关闭过程中再弹窗/改控件
		a.closing.Store(true)
		a.stopServerFast()
		if cur, err := DecryptHostDat(a.host.Path); err == nil && cur == "http://"+a.listen {
			_ = a.host.Restore()
		}
		a.runLog.line("程序退出")
		a.runLog.close()
	})
	a.reloadAccounts()
	a.refreshHostState()
	a.logf("工作目录: %s", a.dir)
	a.logf("缓存目录: %s", a.cacheDir)
	if a.loaderReady.Load() {
		a.logf("上号器: %s", a.loader)
	} else {
		a.logf("上号器还没准备好（%v），正在从云端获取…", a.logErr)
	}
	if a.runLog != nil {
		a.runLog.block("启动 CG v"+toolVersion,
			fmt.Sprintf("工作目录: %s\r\n缓存目录: %s\r\n上号器: %s\r\n", a.dir, a.cacheDir, a.loader))
		a.logf("运行日志: %s", a.logPath)
	} else {
		a.logf("警告: 无法写入 run.log: %v", a.logErr)
	}
	a.logf("自动云备份: CG/%s/%s.zip", a.cfg.cloudNameOrDefault(), a.profileID())
	if a.firstRun {
		a.logf("首次配置完成，拉取云端账号列表…")
	}
	a.refreshCloudList(false)
	if a.logBox != nil {
		a.logBox.SetTextColor(walk.RGB(0, 110, 60))
	}
	a.logf("使用步骤: 新增 -> 选中账号 -> 选中上号")

	// 后台检查云端更新（不阻塞界面，失败只写日志）
	go a.checkUpdate()

	return nil
}

func (a *guiApp) createWindow() error {
	return MainWindow{
		AssignTo: &a.mw,
		Title:    "CG",
		Size:     Size{Width: 1150, Height: 640},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{Spacing: 10},
				Children: []Widget{
					Label{AssignTo: &a.headLabel, Text: "CG", Font: Font{PointSize: 18, Bold: true}},
					Label{AssignTo: &a.headSub, Text: "", TextColor: walk.RGB(120, 120, 120), Font: Font{PointSize: 9}},
					HSpacer{},
					Label{AssignTo: &a.verLabel, Text: "v" + toolVersion, TextColor: walk.RGB(150, 150, 150), Font: Font{PointSize: 8}},
				},
			},
			Composite{
				Layout: HBox{Spacing: 8},
				Children: []Widget{
					Label{Text: "购买地址", TextColor: walk.RGB(110, 110, 110), Font: Font{PointSize: 9}},
					LineEdit{AssignTo: &a.siteEdit, ReadOnly: true, MinSize: Size{Width: 380}},
					PushButton{AssignTo: &a.btnCopySite, Text: "复制", MinSize: Size{Width: 64}, OnClicked: func() {
						if err := copyToClipboard(a.siteEdit.Text()); err != nil {
							a.logf("复制购买地址失败: %v", err)
						} else {
							a.logf("购买地址已复制")
						}
					}},
					PushButton{AssignTo: &a.btnOpenSite, Text: "打开", MinSize: Size{Width: 64}, OnClicked: func() {
						openBrowser(a.siteEdit.Text())
					}},
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{Spacing: 8},
				Children: []Widget{
					Label{Text: "购买说明", TextColor: walk.RGB(110, 110, 110), Font: Font{PointSize: 9}},
					Label{AssignTo: &a.siteNote, Text: "", TextColor: walk.RGB(120, 120, 120), Font: Font{PointSize: 9}},
					HSpacer{},
				},
			},
			VSeparator{},
			Composite{
				Layout: HBox{Spacing: 6},
				Children: []Widget{
					PushButton{AssignTo: &a.btnRecord, Text: "新增", MinSize: Size{Width: 96, Height: 34}, OnClicked: a.startRecord},
					PushButton{AssignTo: &a.btnReplay, Text: "选中上号", MinSize: Size{Width: 96, Height: 34}, OnClicked: a.startReplay},
					PushButton{AssignTo: &a.btnMark, Text: "标记", MinSize: Size{Width: 80, Height: 34}, OnClicked: a.markSelected},
					PushButton{AssignTo: &a.btnVerify, Text: "验证", MinSize: Size{Width: 80, Height: 34}, OnClicked: a.verifySelected},
					PushButton{AssignTo: &a.btnCloudRefresh, Text: "刷新", MinSize: Size{Width: 80, Height: 34}, OnClicked: func() { a.refreshCloudList(true) }},
					HSpacer{},
					PushButton{AssignTo: &a.btnChoose, Text: "切换上号器", OnClicked: a.reconfigure},
				},
			},
			TableView{
				AssignTo:            &a.tv,
				AlternatingRowBG:    true,
				CustomRowHeight:     26,
				LastColumnStretched: true,
				Columns: []TableViewColumn{
					{Title: "状态", Width: 60},
					{Title: "卡密", Width: 150},
					{Title: "QQ", Width: 100},
					{Title: "uid", Width: 110},
					{Title: "备注", Width: 110},
					{Title: "区服", Width: 110},
					{Title: "录制时间", Width: 105},
					{Title: "原始备注（双击查看）", Width: 320},
					{Title: "验证数据（双击查看）", Width: 260},
				},
				Model: a.model,
			},
			GroupBox{
				Title:  "运行日志",
				Layout: VBox{},
				Children: []Widget{
					TextEdit{AssignTo: &a.logBox, ReadOnly: true, VScroll: true, MinSize: Size{Height: 160}},
				},
			},
		},
	}.Create()
}

// ui 把函数排到界面线程执行，并投递一条空消息唤醒消息循环
// （walk 的 Synchronize 只入队，空闲时不会主动执行）。
func (a *guiApp) ui(f func()) {
	if a.mw == nil || a.closing.Load() {
		return
	}
	a.mw.Synchronize(f)
	win.PostMessage(a.mw.Handle(), win.WM_NULL, 0, 0)
}

// 界面日志区最多保留的行数，超出后清屏，避免长时间运行越用越卡（完整内容在 run.log）。
const maxLogViewLines = 600

func (a *guiApp) logf(format string, args ...any) {
	a.runLog.line(format, args...)
	line := time.Now().Format("15:04:05") + "  " + fmt.Sprintf(format, args...) + "\r\n"
	a.ui(func() {
		if a.uiLogLines >= maxLogViewLines {
			_ = a.logBox.SetText("(显示区已清空，完整内容见 run.log)\r\n")
			a.uiLogLines = 0
		}
		a.logBox.AppendText(line)
		a.uiLogLines++
	})
}

// onPayload 把每次 redeem 的原始响应解密后写进 run.log（等价 CLI 的 dump）。
func (a *guiApp) onPayload(action, card, endpoint string, body []byte) {
	dumpPayload(a.runLog, action, card, endpoint, body)
}

func (a *guiApp) msg(title, text string) {
	a.ui(func() {
		walk.MsgBox(a.mw, title, text, walk.MsgBoxIconWarning)
	})
}

func (a *guiApp) upstream() string {
	if u, err := a.host.Upstream(a.listen); err == nil {
		return u
	}
	return ""
}

func (a *guiApp) refreshHostState() {
	cur, err := DecryptHostDat(a.host.Path)
	if err != nil {
		a.logf("host.dat 读取失败: %v", err)
		return
	}
	if cur == "http://"+a.listen {
		a.logf("host.dat 当前指向本地服务（上次异常退出残留），下次启动会自动还原")
		return
	}
	if up := a.upstream(); up != "" {
		a.logf("host.dat 正常，上游 %s", up)
	}
}

// reloadAccounts 在后台线程读盘并解密缓存，再回界面线程刷新表格，避免界面卡住。
// 用代数号丢弃过期的刷新结果，防止连点“刷新”时旧结果覆盖新结果。
func (a *guiApp) reloadAccounts() {
	gen := atomic.AddInt64(&a.loadGen, 1)
	go func() {
		local := loadAccountRows(NewCache(a.cacheDir))
		a.ui(func() {
			if atomic.LoadInt64(&a.loadGen) != gen {
				return
			}
			// 刷新列表时把已有的验号结果带过来，避免每刷一次就清空
			rows := carryVerifyResults(a.model.rows, mergeRows(local, a.cloudIdx))
			a.model.rows = rows
			a.model.PublishRowsReset()
			// 录完自动选中最新一行（列表本来就按录制时间倒序，最新在最上面）
			if a.selectTop.CompareAndSwap(true, false) && len(rows) > 0 {
				a.tv.SetCurrentIndex(0)
			}
			a.updateSelectionState()
		})
	}()
}

// mergeRows 本地优先，云端独有的补在后面并标成“云端”。
func mergeRows(local []accountRow, idx *cloudIndex) []accountRow {
	if idx == nil {
		return local
	}
	seen := make(map[string]bool, len(local))
	rows := make([]accountRow, 0, len(local)+len(idx.Cards))
	for _, r := range local {
		seen[r.Card] = true
		if e := idx.find(r.Card); e != nil {
			if r.QQ == "-" && e.QQ != "" {
				r.QQ = e.QQ
			}
			if r.OpenID == "-" && e.OpenID != "" {
				r.OpenID = e.OpenID
			}
			if r.Zone == "-" && e.Zone != "" {
				r.Zone = e.Zone
			}
		}
		rows = append(rows, r)
	}
	for _, e := range idx.Cards {
		if seen[e.Card] {
			continue
		}
		rows = append(rows, entryToRow(e))
	}
	sortAccountRows(rows)
	return rows
}

func (a *guiApp) reloadAccountsAsync() { a.reloadAccounts() }

func (a *guiApp) setBusy(b bool) {
	a.busy = b
	a.btnRecord.SetEnabled(!b)
	a.btnChoose.SetEnabled(!b)
	a.btnCloudRefresh.SetEnabled(!b)
	a.updateSelectionState()
}

func (a *guiApp) setBusyAsync(b bool) {
	a.ui(func() { a.setBusy(b) })
}

// stopServer 优雅停止（在后台协程里调用，不阻塞界面）。
func (a *guiApp) stopServer() {
	srv := a.takeServer()
	if srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

// stopServerFast 关闭窗口时用，立即断开，避免等待连接导致界面卡住。
func (a *guiApp) stopServerFast() {
	srv := a.takeServer()
	if srv == nil {
		return
	}
	_ = srv.Close()
}

func (a *guiApp) takeServer() *Server {
	a.srvMu.Lock()
	srv := a.srv
	a.srv = nil
	a.srvMu.Unlock()
	return srv
}

// reconfigure 重新选择上号器类型与缓存目录。
func (a *guiApp) reconfigure() {
	if a.busy {
		return
	}

	// 收起第二界面，先把旧目录整个清掉，再回到第一界面重新配置
	a.mw.SetVisible(false)
	a.logf("切换上号器：卸载旧目录 %s", a.dir)
	killed, err := removeDirForce(a.dir)
	if len(killed) > 0 {
		a.logf("已结束占用进程: %s", strings.Join(killed, "、"))
	}
	if err != nil {
		a.logf("卸载旧目录失败: %v", err)
		a.msg("卸载失败", err.Error()+"\n\n请手动关闭占用该目录的程序后重试。")
	} else {
		a.logf("已删除旧目录: %s", a.dir)
	}
	// 配置也一并清掉：下次启动需要重新填写
	if err := os.Remove(configFilePath()); err != nil && !os.IsNotExist(err) {
		a.logf("清除本地配置失败: %v", err)
	} else {
		a.logf("已清除本地配置（%s），下次启动需重新填写", configFilePath())
	}

	// 旧数据已经不存在了，列表和云端索引一并清空
	a.cloudIdx = nil
	a.model.rows = nil
	a.model.PublishRowsReset()

	sr, act := runSetupDialog(nil, a.cfg)
	if act != setupOK {
		// “退出”= 退出软件（此时已经没有第二界面可回）
		a.logf("用户选择退出，关闭程序")
		_ = a.mw.Close()
		return
	}

	a.cfg.Profile = sr.ProfileID
	a.cfg.WorkDir = sr.WorkDir
	a.cfg.CloudName = sr.CloudName
	if e := a.cfg.save(); e != nil {
		a.logf("保存配置失败: %v", e)
	}
	a.applyConfig()
	a.mw.SetVisible(true)
}

// applyConfig 按当前配置切换工作目录、解压上号器、刷新界面。
func (a *guiApp) applyConfig() {
	a.cfg.AutoCloud = true // 录制完成即自动上云，界面暂不提供开关
	a.profile = findProfile(a.cfg.Profile)
	a.dir = a.cfg.WorkDir
	a.cacheDir = filepath.Join(a.dir, "cache")
	a.host = NewHostDatManager(filepath.Join(a.dir, "host.dat"))
	a.loaderReady.Store(false)
	if a.profile != nil {
		a.loader = filepath.Join(a.dir, a.profile.ExeName)
		if ok, reason := localLoaderReady(a.dir, a.profile); ok {
			a.loaderReady.Store(true)
		} else {
			a.logf("本地上号器需要重新获取：%s", reason)
			go a.syncLoaderAsync()
		}
	}
	a.updateLoaderButtons()
	a.logf("工作目录切换为: %s（%s）", a.dir, a.profileName())
	a.refreshHeader()
	a.refreshHostState()
	a.reloadAccounts()
	a.refreshCloudList(false)
}

func (a *guiApp) profileName() string {
	if a.profile == nil {
		return "未知"
	}
	return a.profile.Name
}

func (a *guiApp) startRecord() {
	if a.busy {
		return
	}
	a.launch(true, "")
}

func (a *guiApp) startReplay() {
	if a.busy {
		return
	}
	idx := a.tv.CurrentIndex()
	if idx < 0 || idx >= len(a.model.rows) {
		walk.MsgBox(a.mw, "提示", "请先在列表中选择一个账号。", walk.MsgBoxIconInformation)
		return
	}
	row := a.model.rows[idx]
	if !hasLocalCard(a.cacheDir, row.Card) {
		a.downloadThenReplay(row.Card)
		return
	}
	a.launch(false, row.Card)
}

// downloadThenReplay 本地没有该账号时，先从云端拉这一份再上号。
func (a *guiApp) downloadThenReplay(card string) {
	if a.busy {
		return
	}
	name, pid := a.cfg.cloudNameOrDefault(), a.profileID()
	cacheDir := a.cacheDir
	a.setBusy(true)
	a.logf("本地没有 %s 的缓存，开始按需下载…", card)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.logf("按需下载异常: %v", r)
			}
		}()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		files, size, err := downloadCardCloud(ctx, name, pid, cacheDir, card)
		if err != nil {
			a.setBusyAsync(false)
			a.logf("按需下载失败: %v", err)
			a.msg("下载失败", "无法从云端获取该账号缓存:\n"+err.Error())
			return
		}
		a.logf("已下载 %s: %d 个文件 / %.1f KB", card, files, float64(size)/1024)
		a.setBusyAsync(false)
		a.reloadAccounts()
		a.launch(false, card)
	}()
}

func (a *guiApp) launch(record bool, card string) {
	if !a.loaderReady.Load() {
		walk.MsgBox(a.mw, "提示", "上号器还没准备好（正在从云端下载或下载失败）。\n稍等片刻再试，或点“切换上号器”重新获取。", walk.MsgBoxIconInformation)
		return
	}
	up := a.upstream()
	if up == "" {
		walk.MsgBox(a.mw, "错误", "无法从 host.dat 读取服务器地址，请检查文件是否完整。", walk.MsgBoxIconError)
		return
	}
	if !record && card == "" {
		return
	}

	all := true        // 永远全量录制，保证离线可回放
	offline := !record // 回放永远走本地缓存
	modeName := "录制"
	if !record {
		modeName = "回放"
	}
	a.setBusy(true)
	a.logf("---- 开始%s：%s ----", modeName, card)

	loader := a.loader
	workDir, cacheDir, hostMgr := a.dir, a.cacheDir, a.host
	go func() {
		// 兜底：任何意外都在这里收住，避免整个界面闪退
		defer func() {
			if r := recover(); r != nil {
				a.logf("内部异常: %v", r)
				a.msg("内部异常", fmt.Sprintf("%v", r))
			}
		}()
		defer a.setBusyAsync(false)
		defer a.reloadAccountsAsync()

		a.stopServer()
		if !record {
			if _, changed, err := writeReplayConfig(workDir); err != nil {
				a.logf("写入 config.ini 失败: %v", err)
			} else if changed {
				a.logf("已重写 config.ini（card=%s, dir=请选择游戏目录, ssn=346, auto_zone=0）", placeholderCard)
			} else {
				a.logf("config.ini 已是占位卡密，保留现有内容（含已选好的游戏目录）")
			}
		}

		// 录制先写独立目录，结束后再按卡密合并，避免边录边清导致多个号互相覆盖
		serverCache := cacheDir
		staging := ""
		if record {
			staging = filepath.Join(workDir, stagingDirName)
			if hasContent(staging) {
				if m, _, err := mergeRecording(staging, cacheDir); err == nil && m > 0 {
					a.logf("发现上次未合并的录制，已并入缓存: %d 个账号", m)
				}
			}
			_ = os.RemoveAll(staging)
			if err := os.MkdirAll(staging, 0o700); err != nil {
				a.logf("创建录制目录失败: %v", err)
				a.msg("录制失败", "无法创建录制目录:\n"+err.Error())
				return
			}
			serverCache = staging
			a.logf("本次录制写入临时目录: %s", staging)
		}
		opts := Options{
			Listen:    a.listen,
			Upstream:  up,
			CacheDir:  serverCache,
			All:       all,
			Offline:   offline,
			Refresh:   record,
			Pin:       card,
			Log:       a.logf,
			OnPayload: a.onPayload,
		}
		srv := NewServer(opts)
		ln, err := net.Listen("tcp", a.listen)
		if err != nil {
			a.logf("端口 %s 被占用: %v", a.listen, err)
			a.msg("启动失败", "本地端口 "+a.listen+" 被占用。\n请先关闭其他上号器或缓存服务。")
			return
		}
		go func() {
			if err := srv.Serve(ln); err != nil {
				a.logf("服务异常: %v", err)
			}
		}()
		a.srvMu.Lock()
		a.srv = srv
		a.srvMu.Unlock()

		orig, changed, err := hostMgr.Patch(a.listen)
		if err != nil {
			a.logf("改写 host.dat 失败: %v", err)
			a.msg("启动失败", "改写 host.dat 失败:\n"+err.Error())
			return
		}
		if changed {
			a.logf("host.dat: %s -> http://%s", orig, a.listen)
			defer func() {
				if err := hostMgr.Restore(); err != nil {
					a.logf("还原 host.dat 失败: %v", err)
				} else {
					a.logf("host.dat 已还原")
				}
			}()
		}

		a.logf("启动上号器: %s", loader)
		cmd := exec.Command(loader)
		cmd.Dir = workDir
		if err := cmd.Start(); err != nil {
			if isElevationRequired(err) {
				a.logf("上号器需要管理员权限")
				a.msg("需要管理员权限", "上号器要求管理员权限。\n请退出后右键本程序，选择“以管理员身份运行”。")
			} else {
				a.logf("启动失败: %v", err)
				a.msg("启动失败", err.Error())
			}
			return
		}
		_ = cmd.Wait()
		a.logf("上号器已退出")

		if record && staging != "" {
			cards := stagedCards(staging)
			merged, skipped, err := mergeRecording(staging, cacheDir)
			switch {
			case err != nil:
				a.logf("合并录制失败: %v（数据保留在 %s）", err, staging)
				a.msg("合并失败", "录制数据已保留在:\n"+staging+"\n\n错误: "+err.Error())
			default:
				if merged > 0 {
					a.selectTop.Store(true)
					a.logf("已合并 %d 个账号到缓存，最新录制排在最前", merged)
				} else {
					a.logf("本次没有录制到完整账号（缺少 redeem 数据）")
				}
				if len(skipped) > 0 {
					a.logf("跳过未完成的录制: %s", strings.Join(skipped, ", "))
				}
				_ = os.RemoveAll(staging)

				if len(cards) > 0 && a.cfg != nil {
					a.uploadRecordedCards(cards)
				}
			}
			a.reloadAccounts()
		}
	}()
}

// ---------- 云缓存 ----------

func (a *guiApp) profileID() string {
	if a.profile != nil && a.profile.ID != "" {
		return a.profile.ID
	}
	return "1"
}

// markSelected 给选中账号打标记（封禁 / 黑屋 / 过期 / 其他 / 正常）。
func (a *guiApp) markSelected() {
	if a.busy {
		return
	}
	idx := a.tv.CurrentIndex()
	if idx < 0 || idx >= len(a.model.rows) {
		walk.MsgBox(a.mw, "提示", "请先在列表中选择一个账号。", walk.MsgBoxIconInformation)
		return
	}
	row := a.model.rows[idx]
	status, ok := runMarkDialog(a.mw, row, a.cfg.mark(row.Card))
	if !ok {
		return
	}
	a.cfg.setMark(row.Card, status)
	if err := a.cfg.save(); err != nil {
		a.logf("保存标记失败: %v", err)
	}
	a.logf("标记 %s (%s) -> %s", row.Card, row.QQ, markLabel(status))
	a.model.PublishRowsReset()

	name, pid := a.cfg.cloudNameOrDefault(), a.profileID()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := updateCloudMark(ctx, name, pid, row.Card, status); err != nil {
			a.logf("云端标记同步失败: %v", err)
			return
		}
		a.logf("云端标记已同步: %s -> %s", row.Card, markLabel(status))
	}()
}

// refreshCloudList 拉取云端 index.json 并合并进列表（本地优先）。
func (a *guiApp) refreshCloudList(manual bool) {
	if a.cfg == nil {
		return
	}
	name, pid := a.cfg.cloudNameOrDefault(), a.profileID()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.logf("刷新云端异常: %v", r)
			}
		}()
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		idx, err := fetchCloudIndex(ctx, name, pid)
		if err != nil {
			a.logf("拉取云端列表失败: %v", err)
			if manual {
				a.msg("刷新云端失败", err.Error())
			}
			return
		}
		a.ui(func() { a.applyCloudIndex(idx) })
	}()
}

// applyCloudIndex 记住云端索引，同步标记，然后重新合并表格。
func (a *guiApp) applyCloudIndex(idx *cloudIndex) {
	marked := 0
	for _, e := range idx.Cards {
		if e.Mark != "" {
			a.cfg.setMark(e.Card, e.Mark)
			marked++
		}
	}
	a.cloudIdx = idx
	if marked > 0 {
		if err := a.cfg.save(); err != nil {
			a.logf("保存标记失败: %v", err)
		}
	}
	a.reloadAccounts()
	a.logf("云端列表: %d 个账号", len(idx.Cards))
}

// uploadRecordedCards 录制完成后按账号上传，并合并更新云端索引。
func (a *guiApp) uploadRecordedCards(cards []string) {
	name, pid := a.cfg.cloudNameOrDefault(), a.profileID()
	cacheDir := a.cacheDir
	rows := loadAccountRows(NewCache(cacheDir))
	byCard := make(map[string]accountRow, len(rows))
	for _, r := range rows {
		byCard[r.Card] = r
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	a.logf("自动云备份 -> %s/（%d 个账号）", cloudBaseDir(name, pid), len(cards))
	for _, card := range cards {
		entry := cloudCardEntry{Card: card, Mark: a.cfg.mark(card)}
		if r, ok := byCard[card]; ok {
			entry.QQ = cleanField(r.QQ)
			entry.OpenID = cleanField(r.OpenID)
			entry.Zone = cleanField(r.Zone)
			entry.Note = cleanField(r.Note)
			entry.RawNote = cleanField(r.RawNote)
			if !r.Recorded.IsZero() {
				entry.Recorded = r.Recorded.Format("2006-01-02 15:04:05")
			}
		}
		if err := uploadCardCloud(ctx, name, pid, cacheDir, card, entry); err != nil {
			a.logf("云备份失败 %s: %v", card, err)
			continue
		}
		a.logf("云备份完成 %s", card)
	}
}

// checkUpdate 后台检查云端更新，发现新版本就下载、校验并替换当前 exe。
// 已替换过（比如首配窗口里换的）就只负责重启。
func (a *guiApp) checkUpdate() {
	defer func() {
		if r := recover(); r != nil {
			a.logf("检查更新异常: %v", r)
		}
	}()

	v, err := checkUpdateCore(a.logf)
	if err != nil {
		a.logf("检查更新失败: %v", err)
		return
	}
	if v == "" {
		a.logf("版本 v%s 已是最新", toolVersion)
		return
	}

	if a.busy || a.closing.Load() {
		a.logf("已更新到 v%s，下次启动生效", v)
		a.ui(func() {
			if a.headSub != nil {
				a.headSub.SetText("  " + a.profileName() + "  ·  已更新 v" + v)
			}
		})
		return
	}

	exe, err := os.Executable()
	if err != nil {
		a.logf("已更新到 v%s，下次启动生效", v)
		return
	}
	a.logf("已更新到 v%s，正在重启…", v)
	// 先正常退出（还原 host.dat、关日志），由 main 在 runGUI 返回后拉起新版本
	pendingRestart.Store(exe)
	a.ui(func() { _ = a.mw.Close() })
}

// restartIfNeeded 在界面退出后拉起刚更新好的新版本。
func restartIfNeeded() {
	v := pendingRestart.Load()
	if v == nil {
		return
	}
	exe, _ := v.(string)
	if exe == "" {
		return
	}
	pendingRestart.Store("")

	cmd := exec.Command(exe)
	cmd.Dir = filepath.Dir(exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000008} // DETACHED_PROCESS
	if err := cmd.Start(); err == nil {
		_ = cmd.Process.Release()
		logRestart("已启动新版本: " + exe)
		return
	} else {
		logRestart("直接启动失败: " + err.Error())
	}

	// 兜底：交给 shell 启动（CmdLine 原样传入，不经 exec 的引号转义）
	c := exec.Command("cmd")
	c.SysProcAttr = &syscall.SysProcAttr{
		CmdLine:       `cmd /c start "" "` + exe + `"`,
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	if err := c.Start(); err != nil {
		logRestart("shell 启动也失败: " + err.Error())
		return
	}
	_ = c.Process.Release()
	logRestart("已通过 shell 启动新版本: " + exe)
}

// logRestart 重启结果写进 startup.log（此时 run.log 已关闭）。
func logRestart(msg string) {
	p := filepath.Join(configDirPath(), "startup.log")
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "%s  [restart] %s\n", time.Now().Format("15:04:05"), msg)
}

// updateLoaderButtons 上号器没就绪时禁用“录制/回放”，避免点了报错。
func (a *guiApp) updateLoaderButtons() {
	ready := a.loaderReady.Load()
	if a.btnRecord != nil {
		a.btnRecord.SetEnabled(ready && !a.busy)
	}
	a.updateSelectionState()
}

// syncLoaderAsync 后台从云端获取上号器与 host.dat。
func (a *guiApp) syncLoaderAsync() {
	defer func() {
		if r := recover(); r != nil {
			a.logf("获取上号器异常: %v", r)
		}
	}()
	if a.profile == nil {
		return
	}
	a.ui(func() { a.updateLoaderButtons() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	path, err := syncLoader(ctx, a.dir, a.profile, a.logf)
	if err != nil {
		a.logf("获取上号器失败: %v", err)
		a.ui(func() {
			walk.MsgBox(a.mw, "上号器获取失败",
				"无法从云端获取上号器：\n"+err.Error()+"\n\n请检查网络，或点“切换上号器”重试。",
				walk.MsgBoxIconError)
		})
		return
	}
	a.loader = path
	a.loaderReady.Store(true)
	a.ui(func() { a.updateLoaderButtons() })
}

// refreshHeader 刷新顶部信息：类型名 + 购买地址（不显示本地目录）。
func (a *guiApp) refreshHeader() {
	if a.headSub != nil {
		a.headSub.SetText("  " + a.profileName())
	}
	if a.siteEdit != nil && a.profile != nil {
		a.siteEdit.SetText(a.profile.Site)
	}
	if a.siteNote != nil && a.profile != nil {
		a.siteNote.SetText(a.profile.SiteNote)
	}
}

// showRawNote 弹出原始备注窗口，可自由选中复制。
func (a *guiApp) showRawNote() {
	idx := a.tv.CurrentIndex()
	if idx < 0 || idx >= len(a.model.rows) {
		walk.MsgBox(a.mw, "提示", "请先在列表中选择一个账号。", walk.MsgBoxIconInformation)
		return
	}
	runRawNoteDialog(a.mw, a.model.rows[idx])
}

// updateSelectionState 没选中行时把“选中上号 / 标记”置灰，选中后恢复。
func (a *guiApp) updateSelectionState() {
	if a.tv == nil || a.model == nil {
		return
	}
	idx := a.tv.CurrentIndex()
	hasRow := idx >= 0 && idx < len(a.model.rows)

	if a.btnReplay != nil {
		a.btnReplay.SetEnabled(hasRow && !a.busy && a.loaderReady.Load())
	}
	if a.btnMark != nil {
		a.btnMark.SetEnabled(hasRow && !a.busy)
	}
	if a.btnVerify != nil {
		a.btnVerify.SetEnabled(hasRow && !a.busy && !a.verifying.Load())
	}
}

// verifySelected 手动验证：只查当前选中的账号，不自动刷新、不并发。
func (a *guiApp) verifySelected() {
	idx := a.tv.CurrentIndex()
	if idx < 0 || idx >= len(a.model.rows) {
		walk.MsgBox(a.mw, "提示", "请先在列表中选择一个账号。", walk.MsgBoxIconInformation)
		return
	}
	if a.busy || a.verifying.Load() {
		return
	}

	row := a.model.rows[idx]
	card := strings.TrimSpace(row.Card)
	if card == "" || card == "-" {
		walk.MsgBox(a.mw, "提示", "这一行没有卡密，无法验证。", walk.MsgBoxIconWarning)
		return
	}
	host := strings.TrimSpace(a.upstream())
	if host == "" {
		walk.MsgBox(a.mw, "提示", "拿不到验证服务器地址，请检查 host.dat。", walk.MsgBoxIconWarning)
		return
	}

	a.verifying.Store(true)
	a.setVerifyBusy(true)
	a.logf("验证 %s …", card)

	go func() {
		st := queryVerifyStatus(host, card)
		a.ui(func() {
			a.verifying.Store(false)
			a.setVerifyBusy(false)
			updated := a.storeVerifyResult(card, st)

			// 失败：只弹提示，不展示 JSON
			if st.Err != "" || st.RateLimited || !st.OK {
				a.showVerifyError(st)
				return
			}
			a.logf("验证 %s 完成：%s", card, verifySummary(st))
			runVerifyDialog(a.mw, updated)
		})
	}()
}

// setVerifyBusy 请求进行中禁用按钮，防止连点造成并发。
func (a *guiApp) setVerifyBusy(b bool) {
	if a.btnVerify == nil {
		return
	}
	if b {
		a.btnVerify.SetText("验证中…")
		a.btnVerify.SetEnabled(false)
		return
	}
	a.btnVerify.SetText("验证")
	a.updateSelectionState()
}

// storeVerifyResult 把结果写进对应行并刷新该行。
func (a *guiApp) storeVerifyResult(card string, st verifyStatus) accountRow {
	var updated accountRow
	for i := range a.model.rows {
		row := &a.model.rows[i]
		if row.Card != card {
			continue
		}
		row.Verify = verifySummary(st)
		row.VerifyRaw = st.Raw
		row.VerifyErr = st.Err
		row.VerifyAt = st.At
		a.model.PublishRowChanged(i)
		updated = *row
		break
	}
	return updated
}

// showVerifyError 失败弹框：限频时把需要等待的秒数说清楚。
func (a *guiApp) showVerifyError(st verifyStatus) {
	text := st.Err
	if st.RateLimited {
		wait := st.RetryAfter
		if wait <= 0 {
			wait = 8 * time.Second
		}
		text = fmt.Sprintf("查询过于频繁，请 %.0f 秒后再试。", wait.Seconds()+0.5)
		if st.Message != "" {
			text = st.Message
		}
	} else if text == "" {
		text = st.Message
	}
	if text == "" {
		text = "验证失败"
	}
	a.logf("验证失败：%s", text)
	walk.MsgBox(a.mw, "验证失败", text, walk.MsgBoxIconWarning)
}

// carryVerifyResults 列表重建时按卡密把旧验号结果带过来。
func carryVerifyResults(old, fresh []accountRow) []accountRow {
	if len(old) == 0 || len(fresh) == 0 {
		return fresh
	}
	saved := make(map[string]accountRow, len(old))
	for _, r := range old {
		if r.Verify != "" || r.VerifyRaw != "" || r.VerifyErr != "" {
			saved[r.Card] = r
		}
	}
	if len(saved) == 0 {
		return fresh
	}
	for i := range fresh {
		if o, ok := saved[fresh[i].Card]; ok {
			fresh[i].Verify = o.Verify
			fresh[i].VerifyRaw = o.VerifyRaw
			fresh[i].VerifyErr = o.VerifyErr
			fresh[i].VerifyAt = o.VerifyAt
		}
	}
	return fresh
}

// showVerifyDetail 双击“验证数据”列时打开详情。
func (a *guiApp) showVerifyDetail() {
	idx := a.tv.CurrentIndex()
	if idx < 0 || idx >= len(a.model.rows) {
		walk.MsgBox(a.mw, "提示", "请先在列表中选择一个账号。", walk.MsgBoxIconInformation)
		return
	}
	runVerifyDialog(a.mw, a.model.rows[idx])
}

// runVerifyDialog 查看验号结果与原始响应，交互与原始备注窗口一致。
func runVerifyDialog(owner walk.Form, row accountRow) {
	var dlg *walk.Dialog
	var edit *walk.TextEdit

	d := Dialog{
		AssignTo: &dlg,
		Title:    "验证数据",
		Size:     Size{Width: 720, Height: 520},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "卡密 " + row.Card + "    QQ " + row.QQ, TextColor: walk.RGB(120, 120, 120), Font: Font{PointSize: 9}},
			TextEdit{AssignTo: &edit, Text: verifyDetailText(row), ReadOnly: true, VScroll: true, MinSize: Size{Width: 660, Height: 300}},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "可以选中后 Ctrl+C，或点“复制全部”", TextColor: walk.RGB(150, 150, 150), Font: Font{PointSize: 8}},
					HSpacer{},
					PushButton{Text: "复制全部", MinSize: Size{Width: 100, Height: 34}, OnClicked: func() {
						if err := copyToClipboard(edit.Text()); err != nil {
							walk.MsgBox(dlg, "提示", "复制失败: "+err.Error(), walk.MsgBoxIconWarning)
						}
					}},
					PushButton{Text: "关闭", MinSize: Size{Width: 96, Height: 34}, OnClicked: func() { dlg.Accept() }},
				},
			},
		},
	}
	if err := d.Create(owner); err != nil {
		return
	}
	dlg.Run()
}
