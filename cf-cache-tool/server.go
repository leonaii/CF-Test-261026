package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Options 控制本地缓存服务的行为。
type Options struct {
	Listen   string        // 本地监听地址，例如 127.0.0.1:18800
	Upstream string        // 真实卡密服务器，例如 http://1.2.3.4:8790
	CacheDir string        // 缓存目录
	TTL      time.Duration // 非 redeem 接口的缓存有效期，0 表示永久
	All      bool          // 是否缓存 redeem 之外的接口
	Offline  bool          // 完全离线：不访问上游，只回放缓存
	Refresh  bool          // 本次运行忽略已有缓存，重新录制一次
	Pin      string        // 固定使用某个卡密的缓存（覆盖请求中的卡密）
	Verbose  bool
	Log      func(format string, args ...any)
	// OnPayload 在录制或回放成功后被调用，用于把原始响应落盘（run.log）
	OnPayload func(action, card, endpoint string, body []byte)
}

func (s *Server) payload(action, card, endpoint string, body []byte) {
	if s.opts.OnPayload != nil {
		s.opts.OnPayload(action, card, endpoint, body)
	}
}

// Server 是本工具的本地缓存 HTTP 服务。
type Server struct {
	opts      Options
	cache     *Cache
	client    *http.Client
	refreshed map[string]bool
	mu        sync.Mutex
	pinMu     sync.Mutex
	pin       string
	httpSrv   *http.Server
}

// SetPin 运行期切换要使用的缓存卡密，空字符串表示按请求卡密自动匹配。
func (s *Server) SetPin(card string) {
	s.pinMu.Lock()
	s.pin = card
	s.pinMu.Unlock()
}

// PinCard 返回当前固定的卡密。
func (s *Server) PinCard() string {
	s.pinMu.Lock()
	defer s.pinMu.Unlock()
	return s.pin
}

var cardRe = regexp.MustCompile(`"card_key"\s*:\s*"([^"]{1,64})"`)

func NewServer(opts Options) *Server {
	if opts.CacheDir == "" {
		opts.CacheDir = "cache"
	}
	if opts.Log == nil {
		opts.Log = func(string, ...any) {}
	}
	transport := &http.Transport{
		Proxy:               nil,
		DisableCompression:  true, // 保持上游原始字节，避免 gzip 头与已解压内容不一致
		MaxIdleConnsPerHost: 4,
		IdleConnTimeout:     60 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	return &Server{
		opts:      opts,
		cache:     NewCache(opts.CacheDir),
		client:    &http.Client{Transport: transport, Timeout: 30 * time.Second},
		refreshed: map[string]bool{},
		pin:       opts.Pin,
	}
}

func extractCard(body []byte) string {
	if m := cardRe.FindSubmatch(body); len(m) == 2 {
		return strings.TrimSpace(string(m[1]))
	}
	// 兼容 form 编码
	if vals, err := url.ParseQuery(string(body)); err == nil {
		if v := vals.Get("card_key"); v != "" {
			return v
		}
	}
	return ""
}

func hopHeader(k string) bool {
	switch strings.ToLower(k) {
	case "connection", "proxy-connection", "keep-alive", "proxy-authenticate",
		"proxy-authorization", "te", "trailer", "transfer-encoding", "upgrade", "host":
		return true
	}
	return false
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<20))
	if err != nil {
		s.fail(w, http.StatusBadRequest, "读取请求体失败: "+err.Error())
		return
	}
	card := extractCard(body)
	if card == "" {
		// 解析不到卡密时按请求体哈希分桶，避免不同卡密互相串数据
		sum := sha256.Sum256(body)
		card = "h" + hex.EncodeToString(sum[:8])
	}
	ep := endpointKey(r.URL.Path, r.URL.RawQuery)

	// boot 不录制、不回源，直接返回固定内容
	if ep == "api_card_boot" {
		s.opts.Log("[BOOT] %s 直接返回 {\"ok\":true,\"banned\":false}", r.URL.Path)
		s.write(w, &Meta{
			Status: http.StatusOK,
			Header: map[string][]string{"Content-Type": {"application/json"}},
		}, []byte(`{"ok":true,"banned":false}`))
		return
	}
	cacheOK := ep == "api_card_redeem" || s.opts.All      // 是否允许写缓存
	replayOK := ep == "api_card_redeem" || s.opts.Offline // 是否允许读缓存（boot/peek 只有纯离线才回放）

	// 回放时 config.ini 里是占位卡密，所有接口都按选中的缓存取
	override := false
	if pin := s.PinCard(); pin != "" && pin != card {
		s.opts.Log("[PIN] 请求卡密 %s -> 使用缓存 %s", card, pin)
		card = pin
		override = true
	}
	key := card + "|" + ep

	s.mu.Lock()
	forceRefresh := s.opts.Refresh && !s.refreshed[key]
	if forceRefresh {
		s.refreshed[key] = true
	}
	s.mu.Unlock()

	// 1. 命中缓存
	if replayOK && !forceRefresh {
		if data, meta, err := s.cache.Load(card, ep); err == nil {
			if meta.Status == 0 {
				meta.Status = http.StatusOK
			}
			if ep != "api_card_redeem" && s.opts.TTL > 0 && time.Since(meta.RecordedAt) > s.opts.TTL {
				s.opts.Log("[EXPIRE] %s %s card=%s 缓存超过 %s，走上游", r.Method, r.URL.Path, card, s.opts.TTL)
			} else {
				s.opts.Log("[REPLAY] %s %s card=%s size=%d", r.Method, r.URL.Path, card, len(data))
				s.payload("REPLAY", card, ep, data)
				s.write(w, meta, data)
				return
			}
		}
	}

	// 指定了别的卡密但本地没有它的缓存：直接报错，避免把别的卡密数据写进这份缓存
	if override && ep == "api_card_redeem" {
		s.opts.Log("[MISS] 固定卡密 %s 没有 redeem 缓存", card)
		s.fail(w, http.StatusBadGateway, "固定卡密 "+card+" 没有本地缓存，请先录制")
		return
	}

	// 2. 离线模式：优先回放；boot/peek/log 缺失时自动补录一次，redeem 缺失直接报错
	if s.opts.Offline {
		if cacheOK {
			if data, meta, err := s.cache.Load(card, ep); err == nil {
				s.opts.Log("[REPLAY-OFFLINE] %s %s card=%s size=%d", r.Method, r.URL.Path, card, len(data))
				s.payload("REPLAY", card, ep, data)
				s.write(w, meta, data)
				return
			}
		}
		if ep == "api_card_redeem" {
			s.opts.Log("[MISS] %s %s card=%s 无缓存且处于离线模式", r.Method, r.URL.Path, card)
			s.fail(w, http.StatusBadGateway, "本地没有该卡密的缓存，请先录制")
			return
		}
		s.opts.Log("[AUTO-FETCH] %s %s card=%s 本地无缓存，自动补录一次", r.Method, r.URL.Path, card)
	}

	// 3. 回源
	resp, data, err := s.forward(r, body)
	if err != nil {
		// 上游失败时尽量回放旧缓存，保证“后续步骤”可执行
		if cacheOK {
			if cached, meta, e2 := s.cache.Load(card, ep); e2 == nil {
				s.opts.Log("[FALLBACK] %s %s card=%s 上游失败(%v)，使用旧缓存", r.Method, r.URL.Path, card, err)
				s.payload("REPLAY", card, ep, cached)
				s.write(w, meta, cached)
				return
			}
		}
		s.opts.Log("[ERROR] %s %s card=%s 上游请求失败: %v", r.Method, r.URL.Path, card, err)
		s.fail(w, http.StatusBadGateway, "上游请求失败: "+err.Error())
		return
	}

	if out, ok := extendExpiry(data, time.Now()); ok {
		data = out
		s.opts.Log("[EXPIRE+] %s %s card=%s expires_at 已延长到 %d 年后",
			r.Method, r.URL.Path, card, extendYears)
	}
	meta := &Meta{Status: resp.StatusCode, Header: map[string][]string{}, Request: string(body)}
	for k, v := range resp.Header {
		if hopHeader(k) {
			continue
		}
		meta.Header[k] = v
	}
	const maxCacheBody = 4 << 20 // 单条响应超过 4MB 不缓存，避免磁盘被撑爆
	switch {
	case cacheOK && len(data) > maxCacheBody:
		s.opts.Log("[SKIP] %s %s card=%s 响应 %.1f MB 超过上限，未缓存",
			r.Method, r.URL.Path, card, float64(len(data))/(1<<20))
	case cacheOK:
		if err := s.cache.Save(card, ep, meta, data); err != nil {
			s.opts.Log("[WARN] 写入缓存失败: %v", err)
		} else {
			tag := ""
			if !isKnownEndpoint(ep) {
				tag = " (新接口)"
			}
			s.opts.Log("[RECORD] %s %s card=%s status=%d size=%d%s",
				r.Method, r.URL.Path, card, resp.StatusCode, len(data), tag)
			s.payload("RECORD", card, ep, data)
		}
	default:
		s.opts.Log("[PROXY] %s %s card=%s status=%d size=%d", r.Method, r.URL.Path, card, resp.StatusCode, len(data))
	}
	s.write(w, meta, data)
}

func (s *Server) forward(r *http.Request, body []byte) (*http.Response, []byte, error) {
	target := strings.TrimRight(s.opts.Upstream, "/") + r.URL.RequestURI()
	if s.opts.Verbose {
		s.opts.Log("[UPSTREAM] %s %s", r.Method, target)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, r.Method, target, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	for k, vv := range r.Header {
		if hopHeader(k) {
			continue
		}
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	if u, err := url.Parse(s.opts.Upstream); err == nil {
		req.Host = u.Host
	}
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, nil, err
	}
	return resp, data, nil
}

func (s *Server) write(w http.ResponseWriter, meta *Meta, body []byte) {
	// 回放时再保一次：老缓存里的 expires_at 也一并延长
	if out, ok := extendExpiry(body, time.Now()); ok {
		body = out
	}
	for k, vv := range meta.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	status := meta.Status
	if status == 0 {
		status = http.StatusOK
	}
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func (s *Server) fail(w http.ResponseWriter, code int, msg string) {
	payload, _ := json.Marshal(map[string]any{"success": false, "error": msg})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(payload)))
	w.WriteHeader(code)
	_, _ = w.Write(payload)
}

// Serve 在给定监听器上提供服务，便于界面程序精确控制端口占用与停止。
func (s *Server) Serve(ln net.Listener) error {
	srv := &http.Server{
		Handler:           s,
		ReadHeaderTimeout: 10 * time.Second,
	}
	s.mu.Lock()
	s.httpSrv = srv
	s.mu.Unlock()

	err := srv.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Run 启动服务，阻塞直到进程退出。
func (s *Server) Run() error {
	addr := s.opts.Listen
	if addr == "" {
		addr = "127.0.0.1:18800"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

// Close 立即关闭监听与所有连接。
func (s *Server) Close() error {
	s.mu.Lock()
	srv := s.httpSrv
	s.mu.Unlock()
	if srv == nil {
		return nil
	}
	return srv.Close()
}

// Shutdown 停止服务。
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	srv := s.httpSrv
	s.mu.Unlock()
	if srv == nil {
		return nil
	}
	return srv.Shutdown(ctx)
}
