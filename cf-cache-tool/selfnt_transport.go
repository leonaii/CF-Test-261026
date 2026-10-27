package main

// NT SSO 传输层：复刻 sub_140065630（取地址/连接）、
// sub_14002B4C0（发送）、sub_14003E5A0（接收）和 sub_140031DA0 的解封入口。

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	selfNTSSOListURL     = "https://configsvr.msf.3g.qq.com/configsvr/serverlist.jsp?mType=getssolist"
	selfNTSSODefaultAddr = "msfwifi.3g.qq.com:8080"
	selfNTSSOMaxFrame    = 0x3FFFC
)

var selfNTSSOCache struct {
	sync.Mutex
	at   time.Time
	list []string
}

// selfNTParseSSOEndpoints 复刻 sub_140065010：在任意文本里扫描 IPv4，
// 可带 :port；无端口或端口非法时使用 8080。样本不校验每个 octet 的范围。
func selfNTParseSSOEndpoints(text string) []string {
	var out []string
	seen := map[string]bool{}
	add := func(host string, port int) {
		if port < 1 || port > 65535 {
			port = 8080
		}
		ep := net.JoinHostPort(host, strconv.Itoa(port))
		if !seen[ep] {
			seen[ep] = true
			out = append(out, ep)
		}
	}

	for i := 0; i < len(text); i++ {
		if text[i] < '0' || text[i] > '9' {
			continue
		}
		j, dots := i, 0
		for j < len(text) {
			c := text[j]
			if c == '.' {
				dots++
			} else if c < '0' || c > '9' {
				break
			}
			j++
		}
		if dots != 3 || j-i-7 > 8 {
			continue
		}
		host := text[i:j]
		port, next := 8080, j
		if j < len(text) && text[j] == ':' {
			k, p := j+1, 0
			for k < len(text) && text[k] >= '0' && text[k] <= '9' {
				p = p*10 + int(text[k]-'0')
				k++
			}
			if k > j+1 && p >= 1 && p <= 65535 {
				port, next = p, k
			}
		}
		add(host, port)
		i = next - 1
	}
	return out
}

// selfNTSSOEndpoints 拉取并缓存 10 分钟；失败时仍返回默认地址。
func selfNTSSOEndpoints() []string {
	selfNTSSOCache.Lock()
	defer selfNTSSOCache.Unlock()
	if time.Since(selfNTSSOCache.at) < 10*time.Minute && len(selfNTSSOCache.list) > 0 {
		return append([]string(nil), selfNTSSOCache.list...)
	}

	var list []string
	cli := &http.Client{Timeout: 3500 * time.Millisecond}
	resp, err := cli.Get(selfNTSSOListURL)
	if err == nil {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		list = selfNTParseSSOEndpoints(string(body))
	}
	seen := map[string]bool{}
	for _, ep := range list {
		seen[ep] = true
	}
	if !seen[selfNTSSODefaultAddr] {
		list = append(list, selfNTSSODefaultAddr)
	}
	selfNTSSOCache.at = time.Now()
	selfNTSSOCache.list = append([]string(nil), list...)
	return list
}

// selfNTConnectSSO 按列表顺序直连，返回第一个可用连接。
func selfNTConnectSSO(timeout time.Duration) (net.Conn, string, error) {
	if timeout <= 0 {
		timeout = 4 * time.Second
	}
	var last error
	for _, ep := range selfNTSSOEndpoints() {
		d := net.Dialer{Timeout: timeout}
		conn, err := d.Dial("tcp", ep)
		if err == nil {
			if tc, ok := conn.(*net.TCPConn); ok {
				_ = tc.SetNoDelay(true)
			}
			selfLog("NT SSO 连接成功 %s", ep)
			return conn, ep, nil
		}
		last = err
		selfLog("NT SSO 连接失败 %s: %v", ep, err)
	}
	if last == nil {
		last = fmt.Errorf("没有可用地址")
	}
	return nil, "", last
}

// selfNTSendFrame 对应 sub_14002B4C0：循环写直到全部发完。
func selfNTSendFrame(conn net.Conn, frame []byte) error {
	for len(frame) > 0 {
		n, err := conn.Write(frame)
		if err != nil {
			return err
		}
		if n <= 0 {
			return io.ErrUnexpectedEOF
		}
		frame = frame[n:]
	}
	return nil
}

// selfNTRecvFrame 对应 sub_14003E5A0：BE32 总长度（含自身），随后是正文。
func selfNTRecvFrame(conn net.Conn) ([]byte, error) {
	var head [4]byte
	if _, err := io.ReadFull(conn, head[:]); err != nil {
		return nil, err
	}
	total := selfNTBE32(head[:])
	if total < 4 || total-4 > selfNTSSOMaxFrame {
		return nil, fmt.Errorf("SSO 回包长度非法: %d", total)
	}
	out := make([]byte, total)
	copy(out, head[:])
	if _, err := io.ReadFull(conn, out[4:]); err != nil {
		return nil, err
	}
	return out, nil
}

// selfNTDecodeFrame 使用组帧零密钥解包；FirstAuth 使用 selfNTDecodeFrameKey。
func selfNTDecodeFrame(frame []byte) ([]byte, error) {
	return selfNTDecodeFrameKey(frame, nil)
}

// selfNTSealedOffset 跳过 BE32(12)、flag、两个 LP 和中间的零字节。
// 普通帧得到 15；FirstAuth 因 d2/uin 长度不同而得到动态偏移。
func selfNTSealedOffset(body []byte) (int, bool) {
	if len(body) < 15 || selfNTBE32(body) != 12 {
		return 0, false
	}
	off := 5
	skipLP := func() bool {
		if off+4 > len(body) {
			return false
		}
		n := int(selfNTBE32(body[off:]))
		if n < 4 || off+n > len(body) {
			return false
		}
		off += n
		return true
	}
	if !skipLP() || off >= len(body) {
		return 0, false
	}
	off++ // LP 之间的单字节标志
	if !skipLP() || off >= len(body) || (len(body)-off)%8 != 0 {
		return 0, false
	}
	return off, true
}

// selfNTDecodeFrameKey 解开外层 TEA；优先按 LP 结构定位，保留固定偏移扫描兜底。
func selfNTDecodeFrameKey(frame, key []byte) ([]byte, error) {
	if len(frame) < 4 {
		return nil, fmt.Errorf("帧太短")
	}
	body := frame[4:]
	if off, ok := selfNTSealedOffset(body); ok {
		dec := selfNTTEADecryptKey(body[off:], key)
		if plain, ok := selfNTTEAUnwrap(dec); ok {
			return plain, nil
		}
	}
	for off := 15; off <= 20 && off < len(body); off++ {
		if (len(body)-off)%8 != 0 {
			continue
		}
		dec := selfNTTEADecryptKey(body[off:], key)
		if plain, ok := selfNTTEAUnwrap(dec); ok {
			return plain, nil
		}
	}
	for off := 8; off < 15 && off < len(body); off++ {
		if (len(body)-off)%8 != 0 {
			continue
		}
		dec := selfNTTEADecryptKey(body[off:], key)
		if plain, ok := selfNTTEAUnwrap(dec); ok {
			return plain, nil
		}
	}
	return nil, fmt.Errorf("回包 TEA 解封失败")
}

// selfNTExtractPayload 取 LP(body) 后的 LP(payload)；空 payload 返回空切片。
func selfNTExtractPayload(plain []byte) ([]byte, error) {
	if len(plain) < 8 {
		return nil, errors.New("SSO 明文太短")
	}
	bodyLen := int(selfNTBE32(plain))
	if bodyLen < 4 || bodyLen+4 > len(plain) {
		return nil, fmt.Errorf("SSO body 长度非法: %d/%d", bodyLen, len(plain))
	}
	lp := plain[bodyLen:]
	n := int(selfNTBE32(lp))
	if n < 4 || n > len(lp) {
		return nil, fmt.Errorf("SSO payload 长度非法: %d/%d", n, len(lp))
	}
	return append([]byte(nil), lp[4:n]...), nil
}

// selfNTCall 发一帧并接收一帧，同时返回原始回包与解封后的明文。
func selfNTCall(conn net.Conn, frame []byte, timeout time.Duration) (raw, plain []byte, err error) {
	return selfNTCallKey(conn, frame, nil, timeout)
}

func selfNTCallKey(conn net.Conn, frame, key []byte, timeout time.Duration) (raw, plain []byte, err error) {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	_ = conn.SetDeadline(time.Now().Add(timeout))
	if err = selfNTSendFrame(conn, frame); err != nil {
		return nil, nil, err
	}
	raw, err = selfNTRecvFrame(conn)
	if err != nil {
		return nil, nil, err
	}
	plain, err = selfNTDecodeFrameKey(raw, key)
	return raw, plain, err
}
