package main

// 动态代理池：按需从 tianqiip 取 1 条，缓存 3 分钟。
// 只有「勾选代理」并且真的要发请求时才去取，离线操作不触发。
//
// 取 IP 的节奏：有效期 180s，正常到点即换；两次取 IP 之间至少间隔 120s，
// 避免代理接口异常时每个请求都去打一次。

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	proxyFetchURL    = "http://api.tianqiip.com/getip?secret=yek5w5or&num=1&type=txt&port=1&mr=1&sign=441bfbbfa9dea28120618b3c3477e721"
	proxyTTL         = 3 * time.Minute  // 一个 IP 用 180s
	proxyMinInterval = 2 * time.Minute  // 两次取 IP 的最小间隔
	proxyFetchTimout = 12 * time.Second
)

type proxyPool struct {
	mu        sync.Mutex
	enabled   bool
	addr      string // ip:port
	expire    time.Time
	lastFetch time.Time
	cachePath string // 本地缓存文件，重启后 180s 内可直接复用
	loaded    bool
}

var ymmProxy = &proxyPool{}

// proxyCacheFile 是落盘的代理缓存。
type proxyCacheFile struct {
	Addr      string `json:"addr"`
	Expire    int64  `json:"expire"`     // 毫秒
	LastFetch int64  `json:"last_fetch"` // 毫秒
}

// setCachePath 指定缓存文件位置（启动时调用一次）。
func (p *proxyPool) setCachePath(path string) {
	p.mu.Lock()
	p.cachePath = path
	p.mu.Unlock()
}

// loadCacheLocked 首次取用时从磁盘恢复上次的 IP（调用方持锁）。
func (p *proxyPool) loadCacheLocked() {
	if p.loaded {
		return
	}
	p.loaded = true
	if p.cachePath == "" {
		return
	}
	raw, err := os.ReadFile(p.cachePath)
	if err != nil {
		return
	}
	var c proxyCacheFile
	if json.Unmarshal(raw, &c) != nil || strings.TrimSpace(c.Addr) == "" {
		return
	}
	if _, _, err := splitHostPort(c.Addr); err != nil {
		return
	}
	p.addr = c.Addr
	p.expire = time.UnixMilli(c.Expire)
	if c.LastFetch > 0 {
		p.lastFetch = time.UnixMilli(c.LastFetch)
	}
}

func (p *proxyPool) saveCacheLocked() {
	if p.cachePath == "" {
		return
	}
	c := proxyCacheFile{
		Addr:      p.addr,
		Expire:    p.expire.UnixMilli(),
		LastFetch: p.lastFetch.UnixMilli(),
	}
	b, err := json.Marshal(c)
	if err != nil {
		return
	}
	_ = os.WriteFile(p.cachePath, b, 0o600)
}

func (p *proxyPool) dropCacheLocked() {
	if p.cachePath != "" {
		_ = os.Remove(p.cachePath)
	}
}

// SetEnabled 由界面勾选框调用。
func (p *proxyPool) SetEnabled(v bool) {
	p.mu.Lock()
	p.enabled = v
	if !v {
		p.addr = ""
		p.expire = time.Time{}
		p.loaded = false // 重新开启时允许读磁盘缓存
	}
	p.mu.Unlock()
}

func (p *proxyPool) Enabled() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.enabled
}

// get 返回可用代理地址，必要时向服务端取一条。
func (p *proxyPool) get() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.enabled {
		return "", nil
	}
	// 首次取用时恢复上次落盘的 IP（180s 内直接复用，省一次代理接口往返）
	p.loadCacheLocked()
	// 未过期：直接用
	if p.addr != "" && time.Now().Before(p.expire.Add(-10*time.Second)) {
		return p.addr, nil
	}
	// 已过期但距上次取 IP 不足 120s：先拿旧 IP 顶着，不重复打代理接口
	if !p.lastFetch.IsZero() && time.Since(p.lastFetch) < proxyMinInterval {
		if p.addr != "" {
			return p.addr, nil
		}
		wait := int((proxyMinInterval - time.Since(p.lastFetch)).Seconds())
		return "", fmt.Errorf("代理获取过于频繁，%ds 后可重试", wait)
	}

	p.lastFetch = time.Now()
	client := &http.Client{Timeout: proxyFetchTimout}
	resp, err := client.Get(proxyFetchURL)
	if err != nil {
		return "", fmt.Errorf("取代理失败: %w", err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	txt := strings.TrimSpace(string(b))

	// 正常返回 "ip:port"；异常时是 JSON 错误
	if !strings.Contains(txt, ":") || len(txt) > 64 {
		return "", fmt.Errorf("代理接口异常: %s", selfClip(txt, 120))
	}
	// 只保留一行
	if i := strings.IndexAny(txt, "\r\n"); i >= 0 {
		txt = strings.TrimSpace(txt[:i])
	}
	if _, _, err := splitHostPort(txt); err != nil {
		return "", fmt.Errorf("代理格式异常: %s", txt)
	}

	p.addr = txt
	p.expire = time.Now().Add(proxyTTL)
	p.saveCacheLocked()
	return p.addr, nil
}

func splitHostPort(s string) (string, string, error) {
	i := strings.LastIndex(s, ":")
	if i <= 0 || i == len(s)-1 {
		return "", "", errors.New("缺少端口")
	}
	return s[:i], s[i+1:], nil
}

// client 返回本次请求应使用的 HTTP 客户端。
// 未勾选代理时返回 nil，调用方用默认直连。
func (p *proxyPool) client() *http.Client {
	addr, err := p.get()
	if err != nil || addr == "" {
		return nil
	}
	u, err := url.Parse("http://" + addr)
	if err != nil {
		return nil
	}
	return &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			Proxy:              http.ProxyURL(u),
			DisableKeepAlives:  true,
			IdleConnTimeout:    10 * time.Second,
		},
	}
}

// proxyHTTPClient 给外部直接调用，未启用时返回 nil。
func proxyHTTPClient() *http.Client { return ymmProxy.client() }

// warmup 预先取一条 IP：自动切代理时先拉好，避免首个请求等取 IP。
func (p *proxyPool) warmup() error {
	_, err := p.get()
	return err
}

// drop 丢弃当前 IP（判定不可用时调用），下次 get 会重新取。
// 只在"探测到 IP 不通"时用，调用方控制次数，避免打爆代理接口。
func (p *proxyPool) drop() {
	p.mu.Lock()
	p.addr = ""
	p.expire = time.Time{}
	p.lastFetch = time.Time{}
	p.loaded = true // 别再从磁盘恢复这条坏 IP
	p.dropCacheLocked()
	p.mu.Unlock()
}

// addrNow 返回当前缓存的代理地址（不触发取 IP）。
func (p *proxyPool) addrNow() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.addr
}

// proxyDesc 用于日志展示。
func proxyDesc() string {
	ymmProxy.mu.Lock()
	defer ymmProxy.mu.Unlock()
	if !ymmProxy.enabled {
		return "直连"
	}
	if ymmProxy.addr == "" {
		return "代理（待取）"
	}
	left := int(time.Until(ymmProxy.expire).Seconds())
	if left < 0 {
		left = 0
	}
	return fmt.Sprintf("代理 %s（剩 %ds）", ymmProxy.addr, left)
}
