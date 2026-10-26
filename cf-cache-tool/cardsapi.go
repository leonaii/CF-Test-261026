package main

// 卡密平台后台接口客户端（yuanmeng 页数据来源）。
//
// 接口（与卡密管理后台前端一致）：
//
//	POST /api/auth/login            {"password":"<管理口令或代理码>"} → {ok, role, token}
//	GET  /api/admin/cards?q&status&page&page_size   admin 卡密列表
//	GET  /api/agent/cards?q&status&page&page_size   agent 卡密列表
//	GET  /api/admin/data/<qq>                       某 QQ 的 ini 明文
//
// 鉴权头：admin 用 X-Admin-Token，agent 用 X-Agent-Token。
// 只读，不提供生成/删除/作废等写操作。

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 管理口令不再内置：从 dav 的 CG/conf.json（admin_pass）取，改口令不用重发客户端。
// 卡密平台地址同理，conf.json 的 api_base 可覆盖。
const cardSrvDefault = "http://111.170.141.190:8790"

type cardClient struct {
	mu    sync.Mutex
	base  string
	token string
	role  string
	cli   *http.Client
}

func newCardClient() *cardClient {
	return &cardClient{
		base: cardServer(),
		cli:  &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *cardClient) setBase(base string) {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		base = cardServer()
	}
	c.mu.Lock()
	c.base = base
	c.mu.Unlock()
}

func (c *cardClient) getBase() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.base
}

func (c *cardClient) getRole() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.role
}

func (c *cardClient) loggedIn() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.token != ""
}

// call 发一次请求并解析 JSON；返回状态码，非 2xx 也返回 body 便于取 message。
func (c *cardClient) call(method, path string, body []byte) (map[string]any, int, error) {
	c.mu.Lock()
	base, token, role := c.base, c.token, c.role
	c.mu.Unlock()

	var rd io.Reader
	if body != nil {
		rd = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, base+path, rd)
	if err != nil {
		return nil, 0, fmt.Errorf("构造请求失败")
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		if role == "agent" {
			req.Header.Set("X-Agent-Token", token)
		} else {
			req.Header.Set("X-Admin-Token", token)
		}
	}

	cli := c.cli
	if ymmProxy.Enabled() {
		pc := proxyHTTPClient()
		if pc == nil {
			// 不能悄悄退回直连：直连本来就不通，只会白等 20 秒
			return nil, 0, fmt.Errorf("代理未就绪（未取到 IP）")
		}
		cli = pc
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("读取响应失败")
	}

	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("响应不是 JSON（HTTP %d）", resp.StatusCode)
	}
	if resp.StatusCode == http.StatusUnauthorized {
		// token 失效（服务端重启等），清掉让下次自动重新登录
		c.mu.Lock()
		c.token = ""
		c.mu.Unlock()
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := verifyText(m, "message")
		if msg == "" {
			msg = "HTTP " + strconv.Itoa(resp.StatusCode)
		}
		return m, resp.StatusCode, fmt.Errorf("%s", msg)
	}
	return m, resp.StatusCode, nil
}

// login 用管理口令或代理码登录，成功后缓存 token / role。
func (c *cardClient) login(pass string) error {
	pass = strings.TrimSpace(pass)
	if pass == "" {
		return fmt.Errorf("管理口令未加载（dav 配置拉取失败，检查网络后重试）")
	}
	body, _ := json.Marshal(map[string]string{"password": pass})
	m, _, err := c.call(http.MethodPost, "/api/auth/login", body)
	if err != nil {
		return err
	}
	token := verifyText(m, "token")
	if token == "" {
		return fmt.Errorf("登录响应没有 token")
	}
	c.mu.Lock()
	c.token = token
	c.role = verifyText(m, "role")
	c.mu.Unlock()
	return nil
}

// cardQuery 是列表查询条件。
type cardQuery struct {
	Status   string // unused / active / expired / revoked / all
	Q        string
	Page     int
	PageSize int
}

// cardItem 是列表里的一张卡密。
type cardItem struct {
	CardKey     string
	QueryCode   string
	QQ          string
	Note        string
	Status      string
	ValidHours  int64
	CreatedAt   int64
	ActivatedAt int64
	ExpiresAt   int64
	Owner       string
	AgentName   string
	RedeemCount int64
}

// listCards 拉一页卡密。
func (c *cardClient) listCards(q cardQuery) ([]cardItem, int, error) {
	vals := url.Values{}
	if s := strings.TrimSpace(q.Q); s != "" {
		vals.Set("q", s)
	}
	if q.Status != "" && q.Status != "all" {
		vals.Set("status", q.Status)
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 200
	}
	vals.Set("page", strconv.Itoa(q.Page))
	vals.Set("page_size", strconv.Itoa(q.PageSize))

	path := "/api/admin/cards"
	if c.getRole() == "agent" {
		path = "/api/agent/cards"
	}
	m, _, err := c.call(http.MethodGet, path+"?"+vals.Encode(), nil)
	if err != nil {
		return nil, 0, err
	}

	list, _ := m["list"].([]any)
	out := make([]cardItem, 0, len(list))
	for _, v := range list {
		obj, ok := v.(map[string]any)
		if !ok {
			continue
		}
		out = append(out, cardItemFromMap(obj))
	}
	total := int(verifyInt(m, "total"))
	return out, total, nil
}

func cardItemFromMap(m map[string]any) cardItem {
	note := verifyText(m, "note")
	if note == "" {
		note = verifyText(m, "zone") // 后台列表把备注放在 zone
	}
	return cardItem{
		CardKey:     verifyText(m, "card_key"),
		QueryCode:   verifyText(m, "query_code"),
		QQ:          verifyText(m, "qq"),
		Note:        note,
		Status:      verifyText(m, "status"),
		ValidHours:  verifyInt(m, "valid_hours"),
		CreatedAt:   normalizeMS(verifyInt(m, "created_at")),
		ActivatedAt: normalizeMS(verifyInt(m, "activated_at")),
		ExpiresAt:   normalizeMS(verifyInt(m, "expires_at")),
		Owner:       verifyText(m, "owner"),
		AgentName:   verifyText(m, "agent_name"),
		RedeemCount: verifyInt(m, "redeem_count"),
	}
}

// normalizeMS 把秒级时间戳统一成毫秒。
func normalizeMS(v int64) int64 {
	if v > 0 && v < 100000000000 {
		return v * 1000
	}
	return v
}

// proxyProbe 用当前代理请求一次目标，确认这条 IP 真能到达服务器。
func (c *cardClient) proxyProbe(timeout time.Duration) bool {
	if ymmProxy.addrNow() == "" {
		return false
	}
	pc := proxyHTTPClient()
	if pc == nil {
		return false
	}
	pc.Timeout = timeout

	req, err := http.NewRequest(http.MethodGet, c.getBase()+"/api/auth/me", nil)
	if err != nil {
		return false
	}
	resp, err := pc.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	return resp.StatusCode < 500
}

// probeDirect 直连探测服务器是否可用（不走代理）。
// 2xx/4xx 都说明服务活着；502/504 或连接失败说明后端不可用。
func (c *cardClient) probeDirect(timeout time.Duration) bool {
	cli := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, c.getBase()+"/api/auth/me", nil)
	if err != nil {
		return false
	}
	resp, err := cli.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

	switch resp.StatusCode {
	case http.StatusOK, http.StatusBadRequest, http.StatusUnauthorized,
		http.StatusForbidden, http.StatusNotFound:
		return true
	}
	return false
}

// getData 取某个 QQ 的 ini 明文（需要 admin 权限）。
func (c *cardClient) getData(qq string) (string, error) {
	qq = strings.TrimSpace(qq)
	if qq == "" {
		return "", fmt.Errorf("QQ 为空")
	}
	m, _, err := c.call(http.MethodGet, "/api/admin/data/"+url.PathEscape(qq), nil)
	if err != nil {
		return "", err
	}
	text := verifyText(m, "text")
	if text == "" {
		return "", fmt.Errorf("接口返回空 ini")
	}
	return text, nil
}
