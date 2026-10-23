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
	var statusLabel *walk.Label
	var startBtn *walk.PushButton

	// applyUpdateLock 按更新状态锁/解锁输入：「正在更新」时不能输入也不能启动，
	// 新版本替换完成后只能点重启，重启进新进程才恢复。
	applyUpdateLock := func() {
		if nameEdit == nil || startBtn == nil {
			return
		}
		if appliedVersion() != "" {
			nameEdit.SetEnabled(false)
			startBtn.SetEnabled(true)
			startBtn.SetText("重启到新版本")
			if statusLabel != nil {
				statusLabel.SetText("新版本已就绪，点「重启到新版本」生效")
			}
			return
		}
		if updateBusy.Load() {
			nameEdit.SetEnabled(false)
			startBtn.SetEnabled(false)
			startBtn.SetText("正在更新…")
			return
		}
		// 没在更新（含下载失败）：恢复正常输入
		nameEdit.SetEnabled(true)
		startBtn.SetEnabled(true)
		startBtn.SetText("开始使用")
	}

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

	gray := walk.RGB(140, 140, 140)
	dim := walk.RGB(110, 110, 110)

	d := Dialog{
		AssignTo: &dlg,
		Title:    "CG",
		Size:     Size{Width: 560, Height: 330},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "CG租号", TextColor: walk.RGB(24, 96, 186), Font: Font{Family: "Microsoft YaHei", PointSize: 24, Bold: true}},
			VSpacer{Size: 10},

			Label{Text: "上号器类型", TextColor: dim, Font: Font{PointSize: 9}},
			ComboBox{AssignTo: &combo, Model: names, MinSize: Size{Width: 320}},

			VSpacer{Size: 6},
			Label{Text: "卡密 / 名称", TextColor: dim, Font: Font{PointSize: 9}},
			LineEdit{AssignTo: &nameEdit, Text: defName},

			VSpacer{},
			Label{AssignTo: &statusLabel, Text: "", TextColor: gray, Font: Font{PointSize: 8}},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "v" + toolVersion, TextColor: gray, Font: Font{PointSize: 8}},
					HSpacer{},
					PushButton{AssignTo: &startBtn, Text: "开始使用", MinSize: Size{Width: 132, Height: 38}, OnClicked: func() {
						// 已经替换好新版本：直接重启进新版本
						if appliedVersion() != "" {
							if exe, e := os.Executable(); e == nil {
								pendingRestart.Store(exe)
								dlg.Cancel()
								return
							}
						}
						if updateBusy.Load() {
							walk.MsgBox(dlg, "提示", "正在更新到新版本，请稍候…", walk.MsgBoxIconInformation)
							return
						}

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
	applyUpdateLock() // 窗口打开时可能已经在下载更新了

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
						if dialogClosed.Load() {
							return
						}
						if statusLabel != nil {
							statusLabel.SetText(msg)
						}
						applyUpdateLock()
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

// runRawNoteDialog 查看/复制完整的原始备注。
func runRawNoteDialog(owner walk.Form, row accountRow) {
	var dlg *walk.Dialog
	var edit *walk.TextEdit

	text := row.RawNote
	if text == "" || text == "-" {
		text = "（这个账号没有原始备注）"
	} else {
		text = prettyJSON(text)
	}

	d := Dialog{
		AssignTo: &dlg,
		Title:    "原始备注",
		Size:     Size{Width: 720, Height: 520},
		Layout:   VBox{},
		Children: []Widget{
			Label{Text: "卡密 " + row.Card + "    QQ " + row.QQ, TextColor: walk.RGB(120, 120, 120), Font: Font{PointSize: 9}},
			TextEdit{AssignTo: &edit, Text: text, ReadOnly: true, VScroll: true, MinSize: Size{Width: 660, Height: 300}},
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
