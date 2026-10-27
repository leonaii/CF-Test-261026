package main

// 自研上号：GetGameTicket（转 OP）
//
// 与 LoginByQq 不同，这一步的 body 是明文 JSON，只靠 URL 里的 _ws 做签名：
//   login_check = {"channel_id":0,...,"ct":<RefreshClientTicket 返回>,"st":<票据>,"uin":<QQ>}
//   body        = {"from_src":"client","game_id":2,"login_check":"<转义后的 login_check>",...}
//   ext         = ext=<machine_id>&gid=2&nonce=<6字母>&session_id=<ctid>&timestamp=<秒>&uid=<uid>&version=C10
//   sign        = HMAC-SHA256(R, body ‖ ext)
//   _ws         = hex(AES-128-CBC(ext+"&sign="+sign, 5FCE…, AB3D…))
//
// 成功响应里的关键字段：
//   openid / open_key / ctt / trace_id(HP…) / expire_interval

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	selfGamePath    = "/api/v1/wegame.platform.game.TicketProxy/GetGameTicket"
	selfGameID      = 2
	selfGameWGVer   = "5.7.5.8271"
	selfGameOpenSvc = 8702647
)

// selfGameResult 是一次 GetGameTicket 的结果。
type selfGameResult struct {
	HTTPStatus int
	TraceID    string
	EPHex      string
	Plain      []byte
	OpenID     string
	OpenKey    string
	CTT        string
	Expire     int
	RawBody    []byte
	DiagPath   string
}

// selfGamePlain 是 _ep 解出来后的结构（只取关心的字段）。
type selfGamePlain struct {
	Result struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	} `json:"result"`
	GameST     string `json:"game_st"`
	GameSTKey  string `json:"game_st_key"`
	GameClient string `json:"game_client_key"`
	TraceID    string `json:"trace_id"`
	OpenID     string `json:"openid"`
	OpenKey    string `json:"open_key"`
	CTT        string `json:"ctt"`
	SDKTrace   string `json:"sdk_traceid"`
	Expire     int    `json:"expire_interval"`
	InnerRes   int    `json:"inner_result"`
}

// selfBuildLoginCheck 生成内层 login_check 字符串（字段顺序与上号器一致）。
func selfBuildLoginCheck(t *selfTicket, ct, st string) string {
	return fmt.Sprintf(
		`{"channel_id":0,"client_dev_name":"%s","ct":"%s","mac_id":"%s","machine_id":"%s",`+
			`"promote_id":"0","qq_sec_chk_time":0,"st":"%s","tgt":"","uin":%s,"wg_client_ver":"%s"}`,
		t.ClientDevName(), ct, strings.ToUpper(t.MachineID), strings.ToUpper(t.MachineID),
		st, t.QQ, selfGameWGVer,
	)
}

// selfBuildGameBody 生成外层 body（login_check 需要 JSON 转义）。
func selfBuildGameBody(t *selfTicket, ct, st string) (string, error) {
	lc := selfBuildLoginCheck(t, ct, st)
	esc, err := json.Marshal(lc)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		`{"from_src":"client","game_id":%d,"login_check":%s,"mibao_cmd":1,"mibao_params":"",`+
			`"need_game_ticket_type":1,"open_service_id":%d,"pf":"","st":"%s","tgpid":%s,`+
			`"tgt":"","user_id":"%s","userid_type":0}`,
		selfGameID, string(esc), selfGameOpenSvc, st, t.UID, t.QQ,
	), nil
}

// selfGetGameTicket 走一次转 OP 请求。
func selfGetGameTicket(t *selfTicket, ct, st, sessionID string, saveDir string) (*selfGameResult, error) {
	body, err := selfBuildGameBody(t, ct, st)
	if err != nil {
		return nil, err
	}
	nonce, err := selfNewNonce()
	if err != nil {
		return nil, err
	}
	timestamp := time.Now().Unix()
	r := selfEpRequestR(timestamp, nonce, t.MachineID)

	uid := t.UID
	if uid == "" {
		uid = "0"
	}
	if sessionID == "" {
		sessionID = "0"
	}
	extBase := fmt.Sprintf("ext=%s&gid=%d&nonce=%s&session_id=%s&timestamp=%d&uid=%s&version=%s",
		strings.ToUpper(t.MachineID), selfGameID, nonce, sessionID, timestamp, uid, selfVersionTag)
	sign := selfSign(r, body, extBase)
	ws, err := selfWsEncode(extBase + "&sign=" + sign)
	if err != nil {
		return nil, err
	}
	url := selfBaseURL + selfGamePath + "?_ws=" + ws

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Referer", url)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	res := &selfGameResult{HTTPStatus: resp.StatusCode, RawBody: raw}
	var parsed struct {
		EP      string `json:"_ep"`
		TraceID string `json:"trace_id"`
	}
	if json.Unmarshal(raw, &parsed) == nil {
		res.EPHex = parsed.EP
		res.TraceID = parsed.TraceID
	}

	if res.EPHex != "" {
		if p, e := selfEpDecode(res.EPHex, r); e == nil {
			res.Plain = p
			var gp selfGamePlain
			if json.Unmarshal(p, &gp) == nil {
				res.OpenID = gp.OpenID
				res.OpenKey = gp.OpenKey
				res.CTT = gp.CTT
				res.Expire = gp.Expire
				if gp.TraceID != "" {
					res.TraceID = gp.TraceID
				}
			}
		}
	}

	res.DiagPath = selfSaveGameDiag(saveDir, t, body, ws, url, r, res, raw)
	return res, nil
}

func selfSaveGameDiag(dir string, t *selfTicket, body, ws, url, r string, res *selfGameResult, raw []byte) string {
	if dir == "" {
		return ""
	}
	rec := map[string]any{
		"card":     t.Card,
		"qq":       t.QQ,
		"uid":      t.UID,
		"url":      url,
		"r":        r,
		"body":     body,
		"ws":       ws,
		"http":     res.HTTPStatus,
		"trace_id": res.TraceID,
		"ep":       res.EPHex,
		"plain":    string(res.Plain),
		"openid":   res.OpenID,
		"open_key": res.OpenKey,
		"ctt":      res.CTT,
		"expire":   res.Expire,
		"raw":      string(raw),
		"at":       time.Now().Format(time.RFC3339),
	}
	b, _ := json.MarshalIndent(rec, "", "  ")
	d := filepath.Join(dir, "selfop")
	if err := os.MkdirAll(d, 0o755); err != nil {
		return ""
	}
	p := filepath.Join(d, "last_game.json")
	_ = os.WriteFile(p, b, 0o600)
	return p
}
