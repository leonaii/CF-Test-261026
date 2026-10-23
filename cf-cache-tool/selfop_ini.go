package main

// 用接口取回的现成票据（ct/st/biz_ticket）走「明文通道」换 openid。
//
// 明文通道：不带 _ws，body 直接是 JSON，Cookie 带 uin/skey/p_skey/tgp_biz_ticket/tgp_id。
// 实测服务端会返回明文 JSON（含 openid / open_key / ctt），不再包 _ep。

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const selfGameTicketURL = "https://www.wegame.com.cn/api/v1/wegame.platform.game.TicketProxy/GetGameTicket"

var selfPlainHTTPClient = &http.Client{Timeout: 20 * time.Second}

// selfGetGameTicketFromIni 用 ini 票据 + ck 换 openid。
func selfGetGameTicketFromIni(t *ymmTicket, ck string, qq string) (*selfGameTicketData, error) {
	if t == nil {
		return nil, errors.New("ini 为空")
	}
	if t.CT == "" || t.ST == "" {
		return nil, errors.New("ini 缺少 ct/st")
	}

	// 组一个 selfTicket 供 selfBuildGameBody 复用
	stk := &selfTicket{
		QQ:        qq,
		UID:       t.Tgpid,
		Sig:       t.Sig,
		MachineID: t.MachineID,
		CT:        t.CT,
		Ctid:      t.CTID,
		Raw: map[string]string{
			"st":         t.ST,
			"ct":         t.CT,
			"sig":        t.Sig,
			"machine_id": t.MachineID,
			"biz_ticket": t.BizTicket,
			"ctid":       t.CTID,
		},
	}

	body, err := selfBuildGameBody(stk, t.CT, t.ST)
	if err != nil {
		return nil, fmt.Errorf("构造 body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, selfGameTicketURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 WeGame")

	// Cookie：ck（uin/skey/p_skey）+ tgp_biz_ticket + tgp_id
	cookie := strings.TrimSpace(ck)
	if cookie != "" && !strings.HasSuffix(cookie, ";") {
		cookie += ";"
	}
	if t.BizTicket != "" {
		cookie += "tgp_biz_ticket=" + t.BizTicket + ";"
	}
	if t.Tgpid != "" {
		cookie += "tgp_id=" + t.Tgpid + ";"
	}
	req.Header.Set("Cookie", cookie)

	client := selfPlainHTTPClient
	if p := proxyHTTPClient(); p != nil {
		client = p
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, selfClip(string(raw), 200))
	}

	// 明文里直接找 openid / open_key / ctt
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("响应不是 JSON: %s", selfClip(string(raw), 200))
	}

	pick := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := m[k]; ok {
				if s, ok := v.(string); ok && s != "" {
					return s
				}
			}
		}
		return ""
	}

	out := &selfGameTicketData{
		Uin:     qq,
		Tgpid:   t.Tgpid,
		OpenID:  pick("openid", "open_id"),
		OpenKey: pick("open_key", "openkey"),
		CTT:     pick("ctt", "game_client_key"),
		PfKey:   pick("pf_key"),
	}

	// 少数情况下字段在 result 子对象里
	if out.OpenID == "" {
		if r, ok := m["result"].(map[string]any); ok {
			for _, k := range []string{"openid", "open_id"} {
				if s, ok := r[k].(string); ok && s != "" {
					out.OpenID = s
					break
				}
			}
			for _, k := range []string{"open_key", "openkey"} {
				if s, ok := r[k].(string); ok && s != "" {
					out.OpenKey = s
					break
				}
			}
			if s, ok := r["ctt"].(string); ok && s != "" {
				out.CTT = s
			}
		}
	}

	if out.OpenID == "" {
		// 把错误码带出来，便于定位
		msg := pick("error_message", "message")
		code := ""
		if r, ok := m["result"].(map[string]any); ok {
			if v, ok := r["error_code"]; ok {
				code = fmt.Sprint(v)
			}
			if s, ok := r["error_message"].(string); ok && s != "" {
				msg = s
			}
		}
		return nil, fmt.Errorf("未拿到 openid（code=%s msg=%s）%s", code, msg, selfClip(string(raw), 160))
	}
	return out, nil
}
