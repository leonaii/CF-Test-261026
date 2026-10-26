package main

// 腾讯侧账号状态查询（查封号）。
//
// 五个接口全部用网页 cookie（uin/skey/p_skey）鉴权：
//
//	处罚记录     GET  https://gamesafe.qq.com/api/proxy/punish_query?query_type=3&limit=10
//	冻结状态     GET  https://gamesafe.qq.com/api/frozen?_=<ts>
//	在线状态     GET  https://gamesafe.qq.com/api/online?_=<ts>
//	实名/防沉迷  POST https://rail.tgp.qq.com/rpc/rail.AntiAddictServer.JudgeLoginQQUin
//	登录保护     GET  https://wx.gamesafe.qq.com/cgi-bin/login_protect_status
//
// ck（Cookie）来源三选一：
//  1. 卡密平台验号接口返回的 ck（verifyStatus.CK，最省事，见 verify.go）
//  2. NT 登录产物：clientkey/skey 直接组 cookie（见 ntCookie）
//  3. 手工粘贴浏览器里的 cookie
//
// 解析字段：history_list[].game_id/start_stmp/duration/reason/strategy_desc、
// is_frozen、is_online、protected、result/instructions。
//
// 注意：接口返回层级不稳定（有时套 data、有时套 result），所以统一用
// gsFind 深度找键，不写死路径。

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	gsPunishURL  = "https://gamesafe.qq.com/api/proxy/punish_query?query_type=3&limit=10"
	gsFrozenURL  = "https://gamesafe.qq.com/api/frozen"
	gsOnlineURL  = "https://gamesafe.qq.com/api/online"
	gsProtectURL = "https://wx.gamesafe.qq.com/cgi-bin/login_protect_status"

	gsUA      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	gsTimeout = 15 * time.Second
)

// gsCookie 是调腾讯网页接口用的身份凭证。
type gsCookie struct {
	QQ    string // 纯数字 QQ
	Skey  string
	PSkey string
}

// parseGSCookie 从 "uin=o012345; skey=@abc; p_skey=xyz" 这类串里提取。
// qqHint 用于 cookie 里没有 uin 时兜底。
func parseGSCookie(raw, qqHint string) gsCookie {
	c := gsCookie{QQ: strings.TrimSpace(qqHint)}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return c
	}
	for _, seg := range strings.Split(raw, ";") {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		k, v, ok := strings.Cut(seg, "=")
		if !ok {
			continue
		}
		k = strings.ToLower(strings.TrimSpace(k))
		v = strings.TrimSpace(v)
		switch k {
		case "uin", "p_uin", "pt2gguin", "o":
			if n := gsDigits(strings.TrimPrefix(strings.TrimPrefix(v, "o"), "0")); n != "" && c.QQ == "" {
				c.QQ = n
			}
		case "skey", "p_skey", "r_skey":
			if k == "p_skey" {
				if c.PSkey == "" {
					c.PSkey = v
				}
			} else if c.Skey == "" {
				c.Skey = v
			}
		}
	}
	return c
}

// ntCookie 用 NT 登录产物组装网页 cookie。
// NT 侧拿到的是 uin + skey（+ 可选的域 p_skey）；p_skey 需要单独换，
// 没有就先不传——gamesafe 的处罚/冻结/在线接口只靠 skey 也能过。
func ntCookie(qq, skey, pskey string) gsCookie {
	return gsCookie{QQ: gsDigits(qq), Skey: strings.TrimSpace(skey), PSkey: strings.TrimSpace(pskey)}
}

func (c gsCookie) valid() bool { return c.QQ != "" && c.Skey != "" }

func (c gsCookie) header() string {
	qq := c.QQ
	var b strings.Builder
	b.WriteString("uin=o0")
	b.WriteString(qq)
	b.WriteString("; p_uin=o0")
	b.WriteString(qq)
	if c.Skey != "" {
		b.WriteString("; skey=")
		b.WriteString(c.Skey)
	}
	if c.PSkey != "" {
		b.WriteString("; p_skey=")
		b.WriteString(c.PSkey)
	}
	return b.String()
}

// gsPunishItem 是一条处罚记录。
type gsPunishItem struct {
	GameID       int
	StartStmp    int64
	Duration     int64
	Reason       string
	StrategyDesc string
}

func (p gsPunishItem) StartText() string {
	if p.StartStmp <= 0 {
		return "-"
	}
	ms := p.StartStmp
	if ms < 1e12 { // 秒
		ms *= 1000
	}
	return time.UnixMilli(ms).Format("2006-01-02 15:04:05")
}

func (p gsPunishItem) DurationText() string {
	switch {
	case p.Duration <= 0:
		return "永久/未知"
	case p.Duration >= 86400:
		return fmt.Sprintf("%d 天", p.Duration/86400)
	case p.Duration >= 3600:
		return fmt.Sprintf("%d 小时", p.Duration/3600)
	case p.Duration >= 60:
		return fmt.Sprintf("%d 分钟", p.Duration/60)
	default:
		return fmt.Sprintf("%d 秒", p.Duration)
	}
}

// gsAntiAddictResult 是防沉迷/人脸结果。
type gsAntiAddictResult struct {
	Result       int
	Msg          string
	Instructions []any
}

// gsReport 是一次完整查询的结果。
type gsReport struct {
	QQ string

	Punish     []gsPunishItem
	PunishDone bool

	Frozen     bool
	FrozenDone bool

	Online     bool
	OnlineDone bool

	AntiAddict     gsAntiAddictResult
	AntiAddictDone bool

	Protect     bool
	ProtectDone bool

	Errs []string
	Raw  map[string]string
}

// gsQueryAll 依次查五个接口。任一接口失败只记录到 Errs，不影响其它接口。
func gsQueryAll(c gsCookie) *gsReport {
	r := &gsReport{QQ: c.QQ, Raw: map[string]string{}}
	if !c.valid() {
		r.Errs = append(r.Errs, "缺少 QQ 或 skey，无法查询")
		return r
	}

	{
		items, raw, err := gsPunish(c)
		r.Raw["punish_query"] = raw
		if err != nil {
			r.Errs = append(r.Errs, "处罚记录: "+err.Error())
		} else {
			r.Punish, r.PunishDone = items, true
		}
	}

	{
		v, raw, err := gsBoolAPI(c, gsFrozenURL, "is_frozen")
		r.Raw["frozen"] = raw
		if err != nil {
			r.Errs = append(r.Errs, "冻结状态: "+err.Error())
		} else {
			r.Frozen, r.FrozenDone = v, true
		}
	}

	{
		v, raw, err := gsBoolAPI(c, gsOnlineURL, "is_online")
		r.Raw["online"] = raw
		if err != nil {
			r.Errs = append(r.Errs, "在线状态: "+err.Error())
		} else {
			r.Online, r.OnlineDone = v, true
		}
	}

	{
		aa, raw, err := gsAntiAddict(c.QQ)
		r.Raw["anti_addict"] = raw
		if err != nil {
			r.Errs = append(r.Errs, "防沉迷: "+err.Error())
		} else {
			r.AntiAddict, r.AntiAddictDone = aa, true
		}
	}

	{
		v, raw, err := gsProtect(c)
		r.Raw["login_protect"] = raw
		if err != nil {
			r.Errs = append(r.Errs, "登录保护: "+err.Error())
		} else {
			r.Protect, r.ProtectDone = v, true
		}
	}

	return r
}

func gsPunish(c gsCookie) ([]gsPunishItem, string, error) {
	raw, err := gsGet(gsPunishURL, c)
	if err != nil {
		return nil, raw, err
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return nil, raw, fmt.Errorf("响应不是 JSON（%s）", gsClip(raw, 120))
	}
	list := gsFindFirst(v, "history_list")
	rows, _ := list.([]any)
	items := make([]gsPunishItem, 0, len(rows))
	for _, row := range rows {
		m, _ := row.(map[string]any)
		if m == nil {
			continue
		}
		it := gsPunishItem{
			GameID:       int(gsNum(gsField(m, "game_id"))),
			StartStmp:    int64(gsNum(gsField(m, "start_stmp"))),
			Duration:     int64(gsNum(gsField(m, "duration"))),
			Reason:       gsText(gsField(m, "reason")),
			StrategyDesc: gsText(gsField(m, "strategy_desc")),
		}
		items = append(items, it)
	}
	return items, raw, nil
}

// gsBoolAPI 处理 frozen / online 这两个返回单个布尔字段的接口。
func gsBoolAPI(c gsCookie, endpoint, key string) (bool, string, error) {
	u := endpoint + "?_=" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	raw, err := gsGet(u, c)
	if err != nil {
		return false, raw, err
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return false, raw, fmt.Errorf("响应不是 JSON（%s）", gsClip(raw, 120))
	}
	val, ok := gsBoolValue(gsFindFirst(v, key))
	if !ok {
		return false, raw, fmt.Errorf("响应里没有 %s（%s）", key, gsClip(raw, 120))
	}
	return val, raw, nil
}

func gsProtect(c gsCookie) (bool, string, error) {
	raw, err := gsGet(gsProtectURL, c)
	if err != nil {
		return false, raw, err
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return false, raw, fmt.Errorf("响应不是 JSON（%s）", gsClip(raw, 120))
	}
	val, ok := gsBoolValue(gsFindFirst(v, "protected"))
	if !ok {
		return false, raw, fmt.Errorf("响应里没有 protected（%s）", gsClip(raw, 120))
	}
	return val, raw, nil
}

// gsAntiAddict 不需要登录态：skey 传空串也能返回，只要 uin 是纯数字。
func gsAntiAddict(qq string) (gsAntiAddictResult, string, error) {
	var out gsAntiAddictResult
	qq = gsDigits(qq)
	if qq == "" {
		return out, "", fmt.Errorf("QQ 号为空")
	}
	body := fmt.Sprintf(`{"gameid":3,"skey":"","uin":%s}`, qq)
	req, err := http.NewRequest(http.MethodPost, railAntiAddictURL, strings.NewReader(body))
	if err != nil {
		return out, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", gsUA)

	raw, err := gsDo(req)
	if err != nil {
		return out, raw, err
	}
	var m struct {
		Result       int    `json:"result"`
		Msg          string `json:"msg"`
		Instructions []any  `json:"instructions"`
	}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return out, raw, fmt.Errorf("响应不是 JSON（%s）", gsClip(raw, 120))
	}
	out.Result, out.Msg, out.Instructions = m.Result, m.Msg, m.Instructions
	return out, raw, nil
}

func gsGet(u string, c gsCookie) (string, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", gsUA)
	req.Header.Set("Referer", "https://gamesafe.qq.com/")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if h := c.header(); h != "" {
		req.Header.Set("Cookie", h)
	}
	return gsDo(req)
}

func gsDo(req *http.Request) (string, error) {
	client := &http.Client{Timeout: gsTimeout}
	if p := proxyHTTPClient(); p != nil {
		client = p
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	raw := string(b)
	if resp.StatusCode != http.StatusOK {
		return raw, fmt.Errorf("HTTP %d（%s）", resp.StatusCode, gsClip(raw, 100))
	}
	// 未登录时腾讯会 302 到登录页或直接吐 HTML
	t := strings.TrimSpace(raw)
	if strings.HasPrefix(t, "<") {
		return raw, fmt.Errorf("返回的是 HTML，ck 可能已失效")
	}
	return raw, nil
}

// ---------- JSON 取值工具 ----------

// gsFind 深度找键，返回所有同名值（大小写不敏感）。
func gsFind(v any, key string) []any {
	var out []any
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			if strings.EqualFold(k, key) {
				out = append(out, vv)
			}
			out = append(out, gsFind(vv, key)...)
		}
	case []any:
		for _, vv := range t {
			out = append(out, gsFind(vv, key)...)
		}
	}
	return out
}

func gsFindFirst(v any, key string) any {
	all := gsFind(v, key)
	if len(all) == 0 {
		return nil
	}
	return all[0]
}

func gsField(m map[string]any, key string) any {
	for k, v := range m {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return nil
}

func gsNum(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case json.Number:
		f, _ := t.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(t), 64)
		return f
	case bool:
		if t {
			return 1
		}
	}
	return 0
}

func gsText(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case bool:
		if t {
			return "true"
		}
		return "false"
	case nil:
		return ""
	}
	return fmt.Sprint(v)
}

// gsBoolValue 兼容 true/false、1/0、"true"/"false"、"1"/"0"。
func gsBoolValue(v any) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case float64:
		return t != 0, true
	case string:
		s := strings.ToLower(strings.TrimSpace(t))
		switch s {
		case "true", "1", "yes", "y", "是", "已冻结", "已保护":
			return true, true
		case "false", "0", "no", "n", "否", "未冻结", "未保护", "":
			return false, true
		}
	}
	return false, false
}

func gsDigits(s string) string {
	var b strings.Builder
	for _, c := range strings.TrimSpace(s) {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func gsClip(s string, n int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ---------- 展示 ----------

// Text 生成可直接显示的多行中文结果。
func (r *gsReport) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "QQ：%s\r\n", dash(r.QQ))

	if r.PunishDone {
		if len(r.Punish) == 0 {
			b.WriteString("处罚记录：无\r\n")
		} else {
			fmt.Fprintf(&b, "处罚记录：%d 条\r\n", len(r.Punish))
			for i, it := range r.Punish {
				fmt.Fprintf(&b, "  %d) 游戏=%d 时间=%s 时长=%s 原因=%s 策略=%s\r\n",
					i+1, it.GameID, it.StartText(), it.DurationText(), dash(it.Reason), dash(it.StrategyDesc))
			}
		}
	}

	if r.FrozenDone {
		fmt.Fprintf(&b, "冻结状态：%s\r\n", gsYN(r.Frozen, "已冻结", "未冻结"))
	}
	if r.OnlineDone {
		fmt.Fprintf(&b, "在线状态：%s\r\n", gsYN(r.Online, "在线", "离线"))
	}
	if r.AntiAddictDone {
		aa := r.AntiAddict
		if aa.Result != 0 {
			b.WriteString("实名/防沉迷：查询失败（" + dash(aa.Msg) + "）\r\n")
		} else if len(aa.Instructions) == 0 {
			b.WriteString("实名/防沉迷：正常（无需人脸）\r\n")
		} else {
			j, _ := json.Marshal(aa.Instructions)
			b.WriteString("实名/防沉迷：需要处理 " + gsClip(string(j), 200) + "\r\n")
		}
	}
	if r.ProtectDone {
		fmt.Fprintf(&b, "登录保护：%s\r\n", gsYN(r.Protect, "已开启", "未开启"))
	}
	for _, e := range r.Errs {
		b.WriteString("失败：" + e + "\r\n")
	}
	return b.String()
}

func gsYN(v bool, yes, no string) string {
	if v {
		return yes
	}
	return no
}

// gsQuery 是给界面/命令行用的便捷入口：传 ck 串即可。
func gsQuery(ck, qqHint string) *gsReport {
	c := parseGSCookie(ck, qqHint)
	if c.QQ == "" && qqHint != "" {
		c.QQ = gsDigits(qqHint)
	}
	return gsQueryAll(c)
}

// gsEncodedQQ 用于需要 URL 参数的场景（保底，当前接口都用 cookie）。
func gsEncodedQQ(c gsCookie) string { return url.QueryEscape("o0" + c.QQ) }
