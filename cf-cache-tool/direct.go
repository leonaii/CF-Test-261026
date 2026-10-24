package main

// 首屏直接输入 CF 卡密时的上号流程：不建主界面，弹一个小窗口显示进度，
// 直接登录卡密平台 → 查账号 → 起本地服务器 → 启动上号器 → 等它退出后恢复 host.dat。

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// uiDo 把函数投递到 mw 的消息循环里执行（walk 的 Synchronize 不会自动唤醒循环）。
func uiDo(mw *walk.MainWindow, f func()) {
	if mw == nil {
		return
	}
	mw.Synchronize(f)
	win.PostMessage(mw.Handle(), win.WM_NULL, 0, 0)
}

// runDirectLaunch 显示独立的上号窗口并执行流程，返回后进程即可退出。
func runDirectLaunch(cfg *appConfig, card string) error {
	var mw *walk.MainWindow
	var status *walk.Label
	var detail *walk.Label

	decl := MainWindow{
		AssignTo: &mw,
		Title:    "CG租号 · 上号中",
		Size:     Size{Width: 460, Height: 170},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "CG租号", TextColor: walk.RGB(24, 96, 186), Font: Font{Family: "Microsoft YaHei", PointSize: 14, Bold: true}},
			Label{AssignTo: &status, Text: "准备中…", Font: Font{PointSize: 10}},
			Label{AssignTo: &detail, Text: card, TextColor: walk.RGB(120, 120, 120), Font: Font{PointSize: 9}},
			VSpacer{},
			Label{Text: "上号器退出后本窗口自动关闭", TextColor: walk.RGB(150, 150, 150), Font: Font{PointSize: 8}},
		},
	}
	if err := decl.Create(); err != nil {
		return err
	}

	setStatus := func(s string) {
		uiDo(mw, func() {
			if status != nil {
				status.SetText(s)
			}
		})
	}

	go func() {
		err := doDirectLaunch(cfg, card, setStatus)
		uiDo(mw, func() {
			if err != nil {
				walk.MsgBox(mw, "上号失败", err.Error(), walk.MsgBoxIconError)
			}
			_ = mw.Close()
		})
	}()

	mw.Run()
	return nil
}

// doDirectLaunch 是真正的上号流程，status 用于更新窗口里的进度文字。
func doDirectLaunch(cfg *appConfig, card string, status func(string)) error {
	dir := cfg.WorkDir
	prof := findProfile(cfg.Profile)
	if prof == nil {
		return fmt.Errorf("上号器类型无效")
	}

	launcher := filepath.Join(dir, prof.ExeName)
	if _, err := os.Stat(launcher); err != nil {
		status("正在获取上号器…")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		p, e := syncLoader(ctx, dir, prof, func(string, ...any) {})
		if e != nil {
			return fmt.Errorf("获取上号器失败：%v", e)
		}
		launcher = p
	}

	status("正在登录卡密平台…")
	api := newCardClient()
	if e := api.login(cardPassDefault); e != nil {
		return fmt.Errorf("登录卡密平台失败：%v", e)
	}

	status("正在查询账号…")
	var hit cardItem
	var found bool
	for page := 1; page <= 5 && !found; page++ {
		items, _, e := api.listCards(cardQuery{Q: card, Page: page, PageSize: 100})
		if e != nil {
			return fmt.Errorf("查询卡密失败：%v", e)
		}
		for _, it := range items {
			if strings.EqualFold(it.CardKey, card) {
				hit, found = it, true
				break
			}
		}
		if len(items) < 100 {
			break
		}
	}
	if !found {
		return fmt.Errorf("卡密平台里没有这个卡密：%s", card)
	}
	if strings.TrimSpace(hit.QQ) == "" {
		return fmt.Errorf("卡密 %s 还没有绑定 QQ", card)
	}

	status("正在准备上号数据…")
	row := ymmRow{QQ: hit.QQ, CardKey: hit.CardKey, Note: hit.Note}
	row.ymmParseNote()

	d, e := ymmBuildOffline(api, []ymmRow{row}, card, "") // 卡密直登固定走在线接口
	if e != nil {
		return e
	}

	if _, _, e := writeOfflineConfig(filepath.Dir(launcher)); e != nil {
		return fmt.Errorf("写 config.ini 失败：%v", e)
	}

	status("正在启动上号器…")
	if e := ymmLaunchViaLauncher(launcher, d); e != nil {
		return e
	}
	status("上号器已退出，host.dat 已还原")
	return nil
}
