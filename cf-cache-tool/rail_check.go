package main

// 实名 / 人脸 / 防沉迷实时检查。
//
// rail.tgp.qq.com 的 JudgeLoginQQUin 不需要登录态：skey 传空串也能返回，
// 只要 body 里的 uin 是数字（实测 2026-10-25）。离线数据源没有 ck 时，
// 查封弹窗用这一条补"当前是否需要人脸/防沉迷拦截"。

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const railAntiAddictURL = "https://rail.tgp.qq.com/rpc/rail.AntiAddictServer.JudgeLoginQQUin"

// railAntiAddictText 返回一句可直接显示的中文结论；失败也给一句，不返回空。
func railAntiAddictText(qq string) string {
	qq = strings.TrimSpace(qq)
	if qq == "" {
		return "未查询（没有 QQ 号）"
	}
	for _, c := range qq {
		if c < '0' || c > '9' {
			return "未查询（QQ 号格式不对）"
		}
	}

	body := fmt.Sprintf(`{"gameid":3,"skey":"","uin":%s}`, qq)
	req, err := http.NewRequest(http.MethodPost, railAntiAddictURL, strings.NewReader(body))
	if err != nil {
		return "查询失败（构造请求失败）"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	cli := &http.Client{Timeout: 8 * time.Second}
	if p := proxyHTTPClient(); p != nil {
		cli = p
	}
	resp, err := cli.Do(req)
	if err != nil {
		return "查询失败（网络不通）"
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return "查询失败（读取响应失败）"
	}

	var m struct {
		Result       int      `json:"result"`
		Msg          string   `json:"msg"`
		Instructions []any    `json:"instructions"`
	}
	if json.Unmarshal(raw, &m) != nil {
		return fmt.Sprintf("查询失败（HTTP %d）", resp.StatusCode)
	}
	if m.Result != 0 {
		msg := m.Msg
		if msg == "" {
			msg = fmt.Sprintf("result=%d", m.Result)
		}
		return "查询失败：" + msg
	}
	if len(m.Instructions) == 0 {
		return "正常（无需人脸 / 防沉迷拦截）"
	}
	b, err := json.Marshal(m.Instructions)
	if err != nil {
		return "需要处理（返回内容无法解析）"
	}
	s := string(b)
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	return "需要处理：" + s
}
