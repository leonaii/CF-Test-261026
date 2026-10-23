package main

// 验号接口客户端：直接调用卡密平台的 /api/query/status。
// 与 CF验号器 v1.3 的请求等价：POST JSON，body 只有 code 字段，无签名、无时间戳。
// 该文件不依赖界面，命令行与图形界面都可以复用。
//
// 实测要点（2026-10-23）：
//   - 卡密（16 位大写英文）与编号（8 位小写英文，可含下划线）都能查
//   - 服务端按 IP 限频，约 8 秒一次，触发时返回 rate_limited + retry_after_ms
//   - 成功响应含 ck / cookie（uin、skey、p_skey），是可直接登录的凭证，
//     弹窗里完整展示，注意不要截图外传

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	verifyEndpoint  = "/api/query/status"
	verifyUserAgent = "CfCheck/1.3"
	verifyTimeout   = 20 * time.Second
)

// verifyStatus 是一次验号请求的结果。
type verifyStatus struct {
	OK          bool
	HasData     bool
	CKOK        bool
	RateLimited bool
	RetryAfter  time.Duration

	QueryCode string
	QQ        string
	Zone      string
	Status    string
	Fenghao   string
	Dongjie   string
	Zaixian   string
	Renlian   string
	Shiming   string
	Lingpai   string
	Nickname  string
	Level     string
	Daqu      string
	Zhanqu    string
	Message   string
	ElapsedMS int64

	Raw string
	Err string
	At  time.Time
}

// 复用同一个客户端，串行调用时连接可以复用；调用方负责控制请求间隔。
var verifyHTTPClient = &http.Client{Timeout: verifyTimeout}

// queryVerifyStatus 调用验号接口。host 形如 http://1.2.3.4:8790。
func queryVerifyStatus(host, code string) verifyStatus {
	st := verifyStatus{At: time.Now()}

	code = strings.TrimSpace(code)
	if code == "" {
		st.Err = "卡密为空"
		return st
	}

	base := strings.TrimRight(strings.TrimSpace(host), "/")
	if base == "" {
		st.Err = "服务器地址为空"
		return st
	}

	body, err := json.Marshal(map[string]string{"code": code})
	if err != nil {
		st.Err = "构造请求失败"
		return st
	}

	req, err := http.NewRequest(http.MethodPost, base+verifyEndpoint, bytes.NewReader(body))
	if err != nil {
		st.Err = "构造请求失败"
		return st
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", verifyUserAgent)

	resp, err := verifyHTTPClient.Do(req)
	if err != nil {
		st.Err = "连接失败"
		return st
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		st.Err = "读取响应失败"
		return st
	}
	st.Raw = string(data)

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		st.Err = "响应不是 JSON（HTTP " + resp.Status + "）"
		return st
	}

	st.OK = verifyBool(parsed, "ok")
	st.HasData = verifyBool(parsed, "has_data")
	st.CKOK = verifyBool(parsed, "ck_ok")
	st.RateLimited = verifyBool(parsed, "rate_limited")
	st.RetryAfter = time.Duration(verifyInt(parsed, "retry_after_ms")) * time.Millisecond

	st.QueryCode = verifyText(parsed, "query_code")
	st.QQ = verifyText(parsed, "qq")
	st.Zone = verifyText(parsed, "zone")
	st.Status = verifyText(parsed, "status")
	st.Fenghao = verifyText(parsed, "fenghao")
	st.Dongjie = verifyText(parsed, "dongjie")
	st.Zaixian = verifyText(parsed, "zaixian")
	st.Renlian = verifyText(parsed, "renlian")
	st.Shiming = verifyText(parsed, "shiming")
	st.Lingpai = verifyText(parsed, "lingpai")
	st.Nickname = verifyText(parsed, "nickname")
	st.Level = verifyText(parsed, "level")
	st.Daqu = verifyText(parsed, "daqu")
	st.Zhanqu = verifyText(parsed, "zhanqu")
	st.Message = verifyText(parsed, "message")
	st.ElapsedMS = verifyInt(parsed, "elapsed_ms")
	return st
}

// verifySummary 把结果压成一行，用于列表列显示。
func verifySummary(st verifyStatus) string {
	if st.RateLimited {
		return "查询过于频繁"
	}
	if st.Err != "" {
		return st.Err
	}
	if !st.OK {
		if st.Message != "" {
			return st.Message
		}
		return "验证失败"
	}

	parts := make([]string, 0, 8)
	if st.Status != "" {
		parts = append(parts, st.Status)
	}
	if st.Zaixian != "" {
		parts = append(parts, st.Zaixian)
	}
	if st.Dongjie != "" {
		parts = append(parts, st.Dongjie)
	}
	// 封禁只在异常时提，未封禁不占列宽
	if st.Fenghao != "" && st.Fenghao != "未封禁" {
		parts = append(parts, st.Fenghao)
	}
	if st.Renlian != "" {
		parts = append(parts, st.Renlian)
	}
	if st.Lingpai != "" {
		parts = append(parts, st.Lingpai)
	}

	if len(parts) == 0 {
		if st.Message != "" {
			return st.Message
		}
		return "正常"
	}
	return strings.Join(parts, " ")
}

// verifyDetailText 生成详情窗口文本：先给摘要，再给完整原始响应。
func verifyDetailText(row accountRow) string {
	var b strings.Builder
	fmt.Fprintf(&b, "卡密：%s\r\n", row.Card)
	if row.QQ != "" && row.QQ != "-" {
		fmt.Fprintf(&b, "QQ：%s\r\n", row.QQ)
	}
	if !row.VerifyAt.IsZero() {
		fmt.Fprintf(&b, "查询时间：%s\r\n", row.VerifyAt.Format("2006-01-02 15:04:05"))
	}

	if row.Verify == "" && row.VerifyRaw == "" && row.VerifyErr == "" {
		b.WriteString("\r\n这个账号还没有验证数据，后台会自动刷新。\r\n")
		return b.String()
	}

	fmt.Fprintf(&b, "\r\n验证结果：%s\r\n", dash(row.Verify))
	if row.VerifyErr != "" {
		fmt.Fprintf(&b, "错误原因：%s\r\n", row.VerifyErr)
	}
	if row.VerifyRaw != "" {
		b.WriteString("\r\n完整响应（含 ck / skey，可直接用于登录，请勿外传）：\r\n")
		b.WriteString(prettyJSON(row.VerifyRaw))
		b.WriteString("\r\n")
	}
	return b.String()
}

// prettyJSON 把单行 JSON 展开成多行，方便看清全部内容；不是 JSON 就原样返回。
func prettyJSON(s string) string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" || (trimmed[0] != '{' && trimmed[0] != '[') {
		return s
	}
	var v any
	if err := json.Unmarshal([]byte(trimmed), &v); err != nil {
		return s
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return s
	}
	return string(out)
}

// verifyBool 读取布尔字段，兼容 true/false 与 "true"/"false"。
func verifyBool(m map[string]any, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return strings.EqualFold(strings.TrimSpace(t), "true")
	case float64:
		return t != 0
	default:
		return false
	}
}

// verifyInt 读取整数字段，缺省返回 0。
func verifyInt(m map[string]any, key string) int64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return int64(t)
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

// verifyText 读取文本字段，数字按原样输出，不做占位符替换。
func verifyText(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case bool:
		if t {
			return "是"
		}
		return "否"
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}
