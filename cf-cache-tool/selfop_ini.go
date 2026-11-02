package main

// 用接口取回的现成票据（ct/st/biz_ticket）走「明文通道」换 openid。
//
// 明文通道：不带 _ws，body 直接是 JSON，Cookie 带 uin/skey/p_skey/tgp_biz_ticket/tgp_id。
// 实测服务端会返回明文 JSON（含 openid / open_key / ctt），不再包 _ep。

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const selfGameTicketURL = "https://www.wegame.com.cn/api/v1/wegame.platform.game.TicketProxy/GetGameTicket"

var selfPlainHTTPClient = &http.Client{Timeout: 20 * time.Second}

// selfTickFromIni 把 ini 票据转成 LoginByQq 需要的 selfTicket。
func selfTickFromIni(tk *ymmTicket) *selfTicket {
	return &selfTicket{
		QQ:        tk.QQ,
		UID:       tk.UID,
		Sig:       tk.Sig,
		MachineID: tk.MachineID,
		CT:        tk.CT,
		Ctid:      tk.CTID,
		Raw: map[string]string{
			"ct":         tk.CT,
			"st":         tk.ST,
			"sig":        tk.Sig,
			"dt":         tk.DT,
			"sk":         tk.SessionKey,
			"machine_id": tk.MachineID,
			"client_key": tk.ClientKey,
			"st_key":     tk.STKey,
			"biz_ticket": tk.BizTicket,
		},
	}
}

// selfParseOPPlain 从 LoginByQq 解出来的 _ep 明文里取 OP 字段。
// 真上号器的回包层级不稳定，所以递归找字段；error_code 用 anyStr 避免 8.000102e+06。
func selfParseOPPlain(plain []byte, tk *ymmTicket) (*selfGameTicketData, error) {
	var v any
	if err := json.Unmarshal(plain, &v); err != nil {
		return nil, fmt.Errorf("_ep 不是 JSON：%s", selfClip(selfSanitize(string(plain)), 200))
	}
	str := func(names ...string) string {
		f := findField(v, names...)
		if f == nil {
			return ""
		}
		s := strings.TrimSpace(anyStr(f))
		if s == "-" {
			return ""
		}
		return s
	}

	uin := firstNonEmpty(str("uin"), tk.QQ)
	if uin == "0" {
		uin = tk.QQ
	}
	data := &selfGameTicketData{
		Uin:    uin,
		Tgpid:  firstNonEmpty(str("tgpid", "tgp_id"), tk.Tgpid),
		OpenID: str("openid"),
		// EP 里没有 open_key 时，qq_info.st_key 与离线器 CF_Gamest 的
		// open_key 同为 32 位大写 hex；样本写的就是这个值。
		OpenKey:   firstNonEmpty(str("open_key", "openkey"), str("st_key")),
		CTT:       str("game_client_key", "ctt"),
		GameST:    str("game_st"),
		GameSTKey: str("game_st_key"),
		PfKey:     str("pf_key"),
		ClientKey: str("client_key"),
		STKey:     str("st_key"),
		QQInfoSig: str("sig"),
		SK:        str("sk"),
		CT:        str("ct"),
		DT:        str("dt"),
		BizTicket: str("biz_ticket"),
	}
	if data.OpenID == "" {
		code := firstNonEmpty(str("error_code"), "-")
		msg := firstNonEmpty(str("error_message", "message"), "-")
		return nil, fmt.Errorf("LoginByQq 未返回 openid（code=%s msg=%s）", code, msg)
	}
	return data, nil
}

// selfLoginResultData 从一次 LoginByQq 的结果里取 OP 字段。
// 加密请求回 _ep（能解开就用）；明文请求服务端直接回 JSON（可能就是 result.error_code）。
func selfLoginResultData(res *selfLoginResult, tk *ymmTicket) (*selfGameTicketData, error) {
	if res == nil {
		return nil, errors.New("没有响应")
	}
	if len(res.Plain) > 0 {
		return selfParseOPPlain(res.Plain, tk)
	}
	if res.EPHex != "" {
		return nil, fmt.Errorf("回包 _ep 解不开（%d 字节，HTTP %d）", len(res.EPHex)/2, res.HTTPStatus)
	}
	body := bytes.TrimSpace(res.RawBody)
	if len(body) > 0 && body[0] == '{' {
		return selfParseOPPlain(body, tk)
	}
	return nil, fmt.Errorf("响应不是 JSON（HTTP %d）：%s", res.HTTPStatus, selfClip(selfSanitize(string(res.RawBody)), 160))
}

// selfLoginByIni 复刻真上号器的转 OP 第一步：LoginByQq 取 openid/open_key/ctt。
// _ep 解不开时改用明文 body 再发一次（服务端会把 error_code 直接给出来）；
// mcode 校验失败（8010170）时再用抓包模板里的本机 machine_id 试一次。
func selfLoginByIni(tk *ymmTicket, saveDir string) (*selfGameTicketData, error) {
	stk := selfTickFromIni(tk)
	// 新版流程：只要卡里有 Token0106，就先走 NT 免密登录（KEX → SsoNTLoginEasyLogin
	// → QQConnectLogin.auth）刷新 sig，而不是直接吃 ini 里可能已过期的 sig。
	// 刷新失败才沿用 ini 的 sig，避免 Token0106 失效时整条链路直接断。
	if tok := strings.TrimSpace(tk.Token0106); tok != "" {
		uin, err := strconv.ParseUint(strings.TrimSpace(tk.QQ), 10, 64)
		if err != nil || uin == 0 {
			return nil, errors.New("NT 登录需要有效的 QQ 号")
		}
		var machineID []byte
		if s := strings.TrimSpace(tk.MachineID); s != "" {
			machineID, _ = hex.DecodeString(s)
			if len(machineID) != 16 {
				machineID = nil // 回退到 MD5(uin+"ComputerID")
			}
		}
		token0106, derr := hex.DecodeString(tok)
		if derr == nil && len(token0106) > 0 {
			if nt, e := selfNTLogin(uin, machineID, token0106); e == nil {
				stk.Sig = nt.Sig
				stk.Raw["sig"] = nt.Sig
				tk.Sig = nt.Sig
				if len(nt.Token0106) > 0 {
					tk.Token0106 = strings.ToUpper(hex.EncodeToString(nt.Token0106))
				}
				selfLog("NT 免密登录成功 QQ=%s sig=%s", tk.QQ, shortHex(nt.Sig, 8))
			} else {
				selfLog("NT 免密登录失败，沿用 ini 的 sig（%s）: %v", shortHex(stk.Sig, 8), e)
			}
		} else {
			selfLog("Token0106 不是合法 hex，沿用 ini 的 sig")
		}
	}
	if strings.TrimSpace(stk.Sig) == "" {
		return nil, errors.New("这张卡既没有可用 sig，也无法完成 NT 免密登录")
	}

	hasNT := strings.TrimSpace(stk.Raw["client_key"]) != "" && strings.TrimSpace(stk.Raw["st_key"]) != ""
	login := func(t *selfTicket, encrypt bool) (*selfLoginResult, error) {
		switch {
		case hasNT && encrypt:
			return selfLoginByQqV6(t, saveDir)
		case hasNT:
			return selfLoginByQqV6Plain(t, saveDir)
		}
		return selfLoginByQq(t, saveDir)
	}

	attempt := func(t *selfTicket) (*selfGameTicketData, error) {
		res, err := login(t, true)
		if err != nil {
			return nil, err
		}
		if data, derr := selfLoginResultData(res, tk); derr == nil {
			return data, nil
		} else if len(res.Plain) > 0 {
			return nil, derr // 解开了是业务错误，没必要再试
		}

		// _ep 解不开：换明文 body 再发一次，服务端会直接返回 result
		res2, err := login(t, false)
		if err != nil {
			return nil, fmt.Errorf("明文重试失败：%w", err)
		}
		return selfLoginResultData(res2, tk)
	}

	data, err := attempt(stk)
	if err == nil {
		// type4 的 _ep 只给 openid/client_key/st_key，没有真实 open_key。
		// 样本拿到这些字段后再走一次 qq_info_type=6，才有 open_key。
		if data.ClientKey != "" && data.STKey != "" && (data.OpenKey == "" || data.OpenKey == data.STKey) {
			if v6, e := selfLoginV6FromTicket(stk, tk, data, saveDir); e == nil {
				return v6, nil
			} else {
				selfLog("V6 二次登录失败 QQ=%s: %v", tk.QQ, e)
			}
		}
		return data, nil
	}
	if !strings.Contains(err.Error(), "8010170") {
		return nil, err
	}

	// mcode 与设备绑定：mid_list 里的 blob 是从 selfTemplateMachineID 那台机抓的
	stk.MachineID = selfTemplateMachineID
	return attempt(stk)
}

// selfLoginV6FromTicket 用 type4 返回的 openid/client_key/st_key/ctt/dt/sk
// 再执行样本的 qq_info_type=6 LoginByQq，换取真实 open_key。
func selfLoginV6FromTicket(stk *selfTicket, tk *ymmTicket, d *selfGameTicketData, saveDir string) (*selfGameTicketData, error) {
	v6 := *stk
	v6.Raw = make(map[string]string, len(stk.Raw)+8)
	for k, v := range stk.Raw {
		v6.Raw[k] = v
	}
	v6.UID = firstNonEmpty(d.Tgpid, tk.Tgpid, stk.UID)
	// head.credential.ct 用原始卡密 ct；d.CT 是 type4 回包的新 ct，
	// 填进去会 8010014 认证失败。
	v6.CT = firstNonEmpty(stk.CT, d.CT)
	if d.QQInfoSig != "" {
		v6.Sig = d.QQInfoSig
	}
	v6.Raw["openid"] = d.OpenID
	v6.Raw["client_key"] = d.ClientKey
	v6.Raw["st_key"] = d.STKey
	v6.Raw["ctt"] = d.CTT
	v6.Raw["dt"] = d.DT
	v6.Raw["sk"] = d.SK
	v6.Raw["biz_ticket"] = d.BizTicket
	res, err := selfLoginByQqV6(&v6, saveDir)
	if err != nil {
		return nil, err
	}
	selfLog("V6 回包 %s", selfClip(selfSanitize(string(res.Plain)), 600))
	out, err := selfLoginResultData(res, tk)
	if err != nil {
		return nil, err
	}
	if out.OpenKey == "" {
		return nil, errors.New("V6 回包仍没有 open_key")
	}
	if out.OpenID == "" {
		out.OpenID = d.OpenID
	}
	if out.Tgpid == "" {
		out.Tgpid = d.Tgpid
	}
	return out, nil
}

// parseYmmTicket 把卡密平台下发的 NT 缓存（首行 [QQ]，之后 Key=Value）解析成 ymmTicket。
// 键名不区分大小写、忽略下划线：st_key → stkey，machine_id → machineid。
func parseYmmTicket(iniText, qq string) (*ymmTicket, error) {
	kv, section := parsePayload([]byte(iniText))
	if len(kv) == 0 {
		return nil, errors.New("ini 内容为空或不是 Key=Value 文本")
	}

	qq = firstNonEmpty(qq, section, kv["qq"])
	if qq == "" {
		return nil, errors.New("ini 里没有 QQ 号")
	}

	t := &ymmTicket{
		QQ:         qq,
		UID:        kv["uid"],
		Tgpid:      firstNonEmpty(kv["tgpid"], kv["uid"]),
		CTID:       kv["ctid"],
		MachineID:  strings.ToUpper(kv["machineid"]),
		CT:         kv["ct"],
		ST:         kv["st"],
		DT:         kv["dt"],
		BizTicket:  kv["bizticket"],
		ClientKey:  kv["clientkey"],
		STKey:      kv["stkey"],
		Sig:        kv["sig"],
		SessionKey: kv["sessionkey"],
		Token0106:  kv["token0106"],
		Token010A:  kv["token010a"],
		Token0143:  kv["token0143"],
	}
	if t.CT == "" || t.ST == "" {
		return nil, errors.New("ini 缺少 ct/st，无法换 openid")
	}
	return t, nil
}

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
		OpenKey: firstNonEmpty(pick("open_key", "openkey"), pick("st_key")),
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
				code = anyStr(v) // 不能用 fmt.Sprint：float64 会变成 8.000102e+06
			}
			if s, ok := r["error_message"].(string); ok && s != "" {
				msg = s
			}
		}
		return nil, fmt.Errorf("未拿到 openid（code=%s msg=%s）%s", code, msg, selfClip(string(raw), 160))
	}
	return out, nil
}
