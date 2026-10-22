//go:build windows

package main

// 首次配置向导：选择上号器类型 + 缓存目录，之后本地记住。

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type setupResult struct {
	ProfileID string
	WorkDir   string
	CloudName string
}

func runSetupDialog(owner walk.Form, cfg *appConfig) (*setupResult, setupAction) {
	var dlg *walk.Dialog
	var combo *walk.ComboBox
	var nameEdit *walk.LineEdit
	var siteEdit *walk.LineEdit
	var statusLabel *walk.Label

	var (
		accepted bool
		result   setupResult
	)

	names := profileNames()
	if len(names) == 0 {
		walk.MsgBox(owner, "错误", "程序内没有可用的上号器类型。", walk.MsgBoxIconError)
		return nil, setupQuit
	}

	defIdx := 0
	defName := ""
	if cfg != nil {
		for i := range loaderProfiles {
			if loaderProfiles[i].ID == cfg.Profile {
				defIdx = i
			}
		}
		defName = cfg.CloudName
	}

	site := profileByIndex(defIdx).Site
	siteName := profileByIndex(defIdx).SiteName

	gray := walk.RGB(140, 140, 140)
	dim := walk.RGB(110, 110, 110)

	d := Dialog{
		AssignTo: &dlg,
		Title:    "CG",
		Size:     Size{Width: 620, Height: 430},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "CG", Font: Font{PointSize: 24, Bold: true}},
			Label{Text: "上号器集成工具", TextColor: gray, Font: Font{PointSize: 9}},
			VSpacer{Size: 10},

			Label{Text: "上号器类型", TextColor: dim, Font: Font{PointSize: 9}},
			ComboBox{AssignTo: &combo, Model: names, MinSize: Size{Width: 320}},

			VSpacer{Size: 6},
			Label{Text: "卡密", TextColor: dim, Font: Font{PointSize: 9}},
			LineEdit{AssignTo: &nameEdit, Text: defName},

			VSpacer{Size: 12},
			GroupBox{
				Title:  "购买地址",
				Layout: VBox{},
				Children: []Widget{
					LineEdit{AssignTo: &siteEdit, ReadOnly: true, Text: site},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{Text: "复制链接", MinSize: Size{Width: 96}, OnClicked: func() {
								if err := copyToClipboard(site); err != nil {
									walk.MsgBox(dlg, "提示", "复制失败，请手动选中链接复制。", walk.MsgBoxIconWarning)
								}
							}},
							PushButton{Text: "打开卡网", MinSize: Size{Width: 96}, OnClicked: func() {
								openBrowser(site)
							}},
							HSpacer{},
						},
					},
				},
			},

			VSpacer{},
			Label{AssignTo: &statusLabel, Text: "", TextColor: gray, Font: Font{PointSize: 8}},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "v" + toolVersion, TextColor: gray, Font: Font{PointSize: 8}},
					HSpacer{},
					PushButton{Text: "开始使用", MinSize: Size{Width: 132, Height: 38}, OnClicked: func() {
						name, err := normalizeBackupName(nameEdit.Text())
						if err != nil {
							walk.MsgBox(dlg, "卡密不合法", err.Error(), walk.MsgBoxIconWarning)
							return
						}
						p := profileByIndex(combo.CurrentIndex())
						if p == nil {
							walk.MsgBox(dlg, "提示", "请选择上号器类型。", walk.MsgBoxIconWarning)
							return
						}
						dir := defaultWorkDir()
						if err := os.MkdirAll(dir, 0o700); err != nil {
							walk.MsgBox(dlg, "目录不可用", fmt.Sprintf("无法创建目录:\n%s\n\n%v", dir, err), walk.MsgBoxIconError)
							return
						}
						result = setupResult{ProfileID: p.ID, WorkDir: dir, CloudName: name}
						accepted = true
						dlg.Accept()
					}},
					PushButton{Text: "退出", MinSize: Size{Width: 96, Height: 38}, OnClicked: func() { dlg.Cancel() }},
				},
			},
		},
	}
	if err := d.Create(owner); err != nil {
		return nil, setupQuit
	}
	_ = combo.SetCurrentIndex(defIdx)

	// 切换类型时同步购买地址
	combo.CurrentIndexChanged().Attach(func() {
		p := profileByIndex(combo.CurrentIndex())
		if p == nil {
			return
		}
		site = p.Site
		siteName = p.SiteName
		siteEdit.SetText(p.Site)
	})
	_ = siteName

	// 监听更新进度，实时显示在窗口里
	var dialogClosed atomic.Bool
	go func() {
		for {
			select {
			case msg := <-updateEvents:
				if dialogClosed.Load() {
					return
				}
				if msg == "" {
					continue
				}
				func() {
					defer func() { _ = recover() }()
					dlg.Synchronize(func() {
						if !dialogClosed.Load() && statusLabel != nil {
							statusLabel.SetText(msg)
						}
					})
				}()
			case <-time.After(300 * time.Millisecond):
				if dialogClosed.Load() {
					return
				}
			}
		}
	}()

	dlg.Run()
	dialogClosed.Store(true)

	if !accepted {
		return nil, setupQuit
	}
	return &result, setupOK
}

// copyToClipboard 把文本放进 Windows 剪贴板。
func copyToClipboard(text string) error {
	cmd := exec.Command("cmd", "/c", "clip")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// openBrowser 用系统默认浏览器打开链接。
func openBrowser(url string) {
	if url == "" {
		return
	}
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}

// runMarkDialog 选择账号状态标记。
func runMarkDialog(owner walk.Form, row accountRow, current string) (string, bool) {
	var dlg *walk.Dialog
	var combo *walk.ComboBox

	// 同首配：所有取值在点“确定”时完成，对话框关闭后不再访问控件
	var chosen string
	var accepted bool

	labels := make([]string, len(markOptions))
	idx := 0
	for i, o := range markOptions {
		labels[i] = o.Label
		if o.ID == current {
			idx = i
		}
	}

	d := Dialog{
		AssignTo: &dlg,
		Title:    "CG",
		Size:     Size{Width: 430, Height: 290},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "卡密: " + row.Card},
			Label{Text: "账号: " + row.QQ + "    备注: " + row.Note},
			Label{Text: "状态:"},
			ComboBox{AssignTo: &combo, Model: labels, MinSize: Size{Width: 220}},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{Text: "确定", MinSize: Size{Width: 96, Height: 34}, OnClicked: func() {
						i := combo.CurrentIndex()
						if i < 0 || i >= len(markOptions) {
							i = 0
						}
						chosen = markOptions[i].ID
						accepted = true
						dlg.Accept()
					}},
					PushButton{Text: "取消", MinSize: Size{Width: 96, Height: 34}, OnClicked: func() { dlg.Cancel() }},
				},
			},
		},
	}
	if err := d.Create(owner); err != nil {
		return "", false
	}
	_ = combo.SetCurrentIndex(idx)
	dlg.Run()
	if !accepted {
		return "", false
	}
	return chosen, true
}
