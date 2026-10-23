package main

// 自研上号（单页一键）：选中账号 → 按钮点亮 → 点击后自动解票并登录。

import (
	"encoding/json"
	"fmt"
	"strings"
)

// selfSanitize 去掉 NUL 与控制字符：walk 的 SetText/AppendText 最终走
// syscall.StringToUTF16，字符串里带 \x00 会直接 panic。
func selfSanitize(s string) string {
	if !strings.ContainsAny(s, "\x00") && !hasControl(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == 0:
			// 丢弃
		case r == '\n' || r == '\t':
			b.WriteRune(r)
		case r < 0x20 || r == 0x7f:
			b.WriteByte('.')
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func hasControl(s string) bool {
	for _, r := range s {
		if (r < 0x20 && r != '\n' && r != '\t') || r == 0x7f {
			return true
		}
	}
	return false
}

func selfClip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + " …"
}

func prettyOrRaw(b []byte) string {
	var v any
	if json.Unmarshal(b, &v) == nil {
		if out, err := json.MarshalIndent(v, "", "  "); err == nil {
			return selfClip(selfSanitize(string(out)), 1200)
		}
	}
	return selfClip(selfSanitize(string(b)), 1200)
}

// selfSelected 返回当前选中行（无选中返回 nil）。
func (a *guiApp) selfSelected() *accountRow {
	if a.model == nil || a.tv == nil {
		return nil
	}
	i := a.tv.CurrentIndex()
	if i < 0 || i >= len(a.model.rows) {
		return nil
	}
	r := a.model.rows[i]
	return &r
}

// updateSelfOpButton 让「自研上号」按钮跟随选中行与忙碌状态。
func (a *guiApp) updateSelfOpButton() {
	if a.btnSelfOp == nil || a.tv == nil {
		return
	}
	a.btnSelfOp.SetEnabled(a.selfSelected() != nil && !a.selfBusy.Load())
}

func (a *guiApp) onSelfOpClick() {
	row := a.selfSelected()
	if row == nil {
		a.logf("自研上号：请先在列表里选中一个账号")
		return
	}
	if a.selfBusy.Load() {
		a.logf("自研上号：上一个任务还在执行")
		return
	}
	a.selfBusy.Store(true)
	a.updateSelfOpButton()
	card := row.Card
	a.markSelfState(card, "进行中", "")
	a.logf("自研上号 %s：开始", card)
	go a.runSelfOp(card)
}

// markSelfState 更新某行的自研状态并刷新表格。
func (a *guiApp) markSelfState(card, state, errText string) {
	a.ui(func() {
		if a.model == nil {
			return
		}
		for i := range a.model.rows {
			if a.model.rows[i].Card == card {
				a.model.rows[i].SelfState = state
				a.model.rows[i].SelfErr = errText
				break
			}
		}
		a.model.PublishRowsReset()
	})
}

// runSelfOp 是一条完整链路：自检 → 解票 → 登录。
func (a *guiApp) runSelfOp(card string) {
	defer func() {
		a.selfBusy.Store(false)
		a.updateSelfOpButton()
		if r := recover(); r != nil {
			a.logf("自研上号 %s：内部异常 %v", card, r)
			a.markSelfState(card, "失败", fmt.Sprint(r))
		}
	}()

	if !a.selfChecked {
		a.selfChecked = true
		for _, line := range selfSelfCheck() {
			a.logf("自研上号 自检 %s", line)
		}
	}

	t, err := loadSelfTicket(NewCache(a.cacheDir), card)
	if err != nil {
		a.logf("自研上号 %s：解票失败 %v", card, err)
		a.markSelfState(card, "失败", err.Error())
		return
	}
	a.selfTk = t
	a.logf("自研上号 %s：解票成功 QQ=%s uid=%s machine=%s", card, t.QQ, t.UID, t.MachineID)
	a.logf("自研上号 %s：client_dev_name=%s mid=%s", card, t.ClientDevName(), t.Mid())

	res, err := selfLoginByQq(t, a.cacheDir)
	if err != nil {
		a.logf("自研上号 %s：登录请求失败 %v", card, err)
		a.markSelfState(card, "失败", err.Error())
		return
	}
	a.logf("自研上号 %s：HTTP %d 响应 %d 字节 trace=%s",
		card, res.HTTPStatus, len(res.RawBody), res.TraceID)

	switch {
	case res.Plain != nil:
		a.logf("自研上号 %s：登录成功 %s", card, selfClip(selfSanitize(string(res.Plain)), 500))
		a.markSelfState(card, "成功", "")
	case res.EPHex != "":
		a.logf("自研上号 %s：LoginByQq 收到 _ep（%d 字节），继续尝试转 OP", card, len(res.EPHex)/2)
		a.markSelfState(card, "失败", "_ep 未解密")
		a.selfOpTransferOP(card, t)
	default:
		a.logf("自研上号 %s：响应里没有 _ep：%s",
			card, selfClip(selfSanitize(string(res.RawBody)), 300))
		a.markSelfState(card, "失败", "响应格式未知")
	}
	if res.DiagPath != "" {
		a.logf("自研上号 %s：诊断 %s", card, res.DiagPath)
	}
}

// selfOpTransferOP 用 redeem 里的 st/ct 直接走 GetGameTicket（转 OP），
// 不依赖 LoginByQq 的响应。
func (a *guiApp) selfOpTransferOP(card string, t *selfTicket) {
	st := firstNonEmpty(t.Raw["st"], t.Raw["sgt"], t.Raw["st_key"])
	if st == "" {
		a.logf("自研上号 %s：redeem 里没有 st，跳过转 OP", card)
		return
	}
	a.logf("自研上号 %s：转 OP 开始 st=%s… ct=%s…", card, shortHex(st, 12), shortHex(t.CT, 12))

	res, err := selfGetGameTicket(t, t.CT, st, t.Ctid, a.cacheDir)
	if err != nil {
		a.logf("自研上号 %s：转 OP 请求失败 %v", card, err)
		return
	}
	a.logf("自研上号 %s：转 OP HTTP %d 响应 %d 字节 trace=%s",
		card, res.HTTPStatus, len(res.RawBody), res.TraceID)

	switch {
	case res.OpenID != "":
		a.logf("自研上号 %s：转 OP 成功 openid=%s open_key=%s expire=%ds",
			card, res.OpenID, res.OpenKey, res.Expire)
		a.markSelfState(card, "成功", "")
	case res.EPHex != "":
		a.logf("自研上号 %s：转 OP 收到 _ep（%d 字节），改走 CF_Gamest 通道", card, len(res.EPHex)/2)
		a.selfOpLaunchFromGamest(card)
	default:
		a.logf("自研上号 %s：转 OP 响应异常 %s",
			card, selfClip(selfSanitize(string(res.RawBody)), 200))
	}
	if res.DiagPath != "" {
		a.logf("自研上号 %s：转 OP 诊断 %s", card, res.DiagPath)
	}
}

// selfOpLaunchFromGamest 读上号器写下的 TCLS\CF_Gamest（明文票据），
// 再用 CG 自己的实现重写 mmog_data.xml + 启动 crossfire.exe。
func (a *guiApp) selfOpLaunchFromGamest(card string) {
	t, p, err := selfReadCFGamest("")
	if err != nil {
		a.logf("自研上号 %s：%v", card, err)
		return
	}
	a.logf("自研上号 %s：读到 %s", card, p)
	a.logf("自研上号 %s：uin=%s tgpid=%s openid=%s open_key=%s ctt=%s…",
		card, t.Uin, t.Tgpid, t.OpenID, t.OpenKey, shortHex(t.CTT, 12))

	pid, err := selfLaunchGame("", t)
	if err != nil {
		a.logf("自研上号 %s：启动失败 %v", card, err)
		a.markSelfState(card, "失败", err.Error())
		return
	}
	a.logf("自研上号 %s：已启动 crossfire.exe pid=%d", card, pid)
	a.markSelfState(card, "成功", "")
}
