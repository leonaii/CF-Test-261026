package main

// 自研上号：WeGame UnionAuth 接口调用。

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
	selfBaseURL     = "https://www.wegame.com.cn"
	selfLoginPath   = "/api/v2/wegame.auth.game.UnionAuth/LoginByQq"
	selfMappid      = 10001
	selfClientType  = 1000000
	selfAppID       = 2112513
	selfWGClientVer = "1688862745363712"
	selfVersionTag  = "C10"
	selfUDIDGid     = 55555
)

// selfLoginResult 是一次 LoginByQq 的完整结果。
type selfLoginResult struct {
	HTTPStatus int
	TraceID    string
	EPHex      string
	Plain      []byte
	UsedR      string
	RawBody    []byte
	DiagPath   string
}

func selfNewNonce() (string, error) {
	b := make([]byte, 6)
	if _, err := randRead(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = byte('A' + int(b[i])%26)
	}
	return string(b), nil
}

// selfBuildLoginPlain 生成 LoginByQq 的明文 JSON（字段顺序与上号器一致）。
func selfBuildLoginPlain(t *selfTicket, midList string) string {
	return fmt.Sprintf(
		`{"mappid":%d,"clienttype":%d,"login_info":{"qq_info_type":4,"appid":%d,"openid":"%s","sig":"%s"},`+
			`"config_params":{"lang_type":0,"login_type":1,"channel_id":0,"promote_id":"0",`+
			`"client_version":"%s","client_dev_name":"%s"},`+
			`"head":{"credential":{"game_ver":"","ct":"","mid_list":"%s","wg_client_ver":"%s"}}}`,
		selfMappid, selfClientType, selfAppID, t.QQ, t.Sig,
		selfWGClientVer, t.ClientDevName(), midList, selfWGClientVer,
	)
}

// selfBuildLoginPlainV6 是 qq_info_type=6 的模板（样本 .rdata 0xF7741 / 0xF7E49）：
// openid/sig/client_key/st_key 都是 NT 登录产物；head.credential 还要带 ct/ctt/dt/sk。
func selfBuildLoginPlainV6(t *selfTicket, midList string) string {
	return fmt.Sprintf(
		`{"mappid":%d,"clienttype":%d,"login_info":{"qq_info_type":6,"appid":%d,"openid":"%s","sig":"%s","client_key":"%s","st_key":"%s"},`+
			`"config_params":{"lang_type":0,"login_type":1,"channel_id":0,"promote_id":"0",`+
			`"client_version":"%s","client_dev_name":"%s"},`+
			`"head":{"credential":{"game_ver":"","ct":"%s","ctt":"%s","dt":"%s","sk":"%s","mid_list":"%s","wg_client_ver":"%s"}}}`,
		selfMappid, selfClientType, selfAppID, firstNonEmpty(t.UID, t.QQ), t.Sig,
		strings.TrimSpace(t.Raw["client_key"]), strings.TrimSpace(t.Raw["st_key"]),
		selfWGClientVer, t.ClientDevName(),
		t.CT, strings.TrimSpace(t.Raw["ctt"]), strings.TrimSpace(t.Raw["dt"]), strings.TrimSpace(t.Raw["sk"]),
		midList, selfWGClientVer,
	)
}

// selfLoginByQq 走一次完整的三层封装并发送请求（qq_info_type=4，纯 QQ + sig）。
func selfLoginByQq(t *selfTicket, saveDir string) (*selfLoginResult, error) {
	return selfLoginByQqPlain(t, false, saveDir, true)
}

// selfLoginByQqV6 用 qq_info_type=6（client_key/st_key）登录，真上号器 NT 登录后走这条。
func selfLoginByQqV6(t *selfTicket, saveDir string) (*selfLoginResult, error) {
	return selfLoginByQqPlain(t, true, saveDir, true)
}

// selfLoginByQqV6Plain 同 V6，但 body 不套 _eq 直接发明文：
// 服务端对明文请求回的是明文 JSON，报错码能直接看（_ep 解不开时的兜底）。
func selfLoginByQqV6Plain(t *selfTicket, saveDir string) (*selfLoginResult, error) {
	return selfLoginByQqPlain(t, true, saveDir, false)
}

func selfLoginByQqPlain(t *selfTicket, useOpenID bool, saveDir string, encrypt bool) (*selfLoginResult, error) {
	midList, err := buildMidList(t.MachineID)
	if err != nil {
		return nil, err
	}
	plain := selfBuildLoginPlain(t, midList)
	if useOpenID {
		plain = selfBuildLoginPlainV6(t, midList)
	}
	nonce, err := selfNewNonce()
	if err != nil {
		return nil, err
	}
	timestamp := time.Now().Unix()
	r := selfEpRequestR(timestamp, nonce, t.MachineID)
	var body string
	if encrypt {
		body, err = selfEqEncodeWithR([]byte(plain), r)
	} else {
		body = plain
	}
	if err != nil {
		return nil, err
	}
	extBase := fmt.Sprintf("ext=%s&gid=%d&nonce=%s&session_id=0&timestamp=%d&uid=0&version=%s",
		strings.ToUpper(t.MachineID), selfUDIDGid, nonce, timestamp, selfVersionTag)
	sign := selfSign(r, body, extBase)
	ws, err := selfWsEncode(extBase + "&sign=" + sign)
	if err != nil {
		return nil, err
	}
	url := selfBaseURL + selfLoginPath + "?_ws=" + ws

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", "")

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

	res := &selfLoginResult{HTTPStatus: resp.StatusCode, RawBody: raw}
	var parsed struct {
		EP      string `json:"_ep"`
		TraceID string `json:"trace_id"`
	}
	if json.Unmarshal(raw, &parsed) == nil {
		res.EPHex = parsed.EP
		res.TraceID = parsed.TraceID
	}

	// 与样本一致：先试请求 R，再试由 timestamp/nonce/machineID 派生的响应 R。
	if res.EPHex != "" {
		if p, e := selfEpDecode(res.EPHex, r); e == nil {
			res.Plain = p
			res.UsedR = r
		}
		if len(res.Plain) == 0 {
			responseR := selfEpResponseR(timestamp, nonce, t.MachineID)
			if p, e := selfEpDecode(res.EPHex, responseR); e == nil {
				res.Plain = p
				res.UsedR = responseR
			}
		}
	}

	res.DiagPath = selfSaveDiag(saveDir, t, body, ws, url, r, res)
	return res, nil
}

// selfSaveDiag 把请求与响应落盘，便于离线分析 _ep 的密钥规则。
func selfSaveDiag(dir string, t *selfTicket, body, ws, url, r string, res *selfLoginResult) string {
	if dir == "" {
		return ""
	}
	d := filepath.Join(dir, "selfop")
	if err := os.MkdirAll(d, 0o755); err != nil {
		return ""
	}
	rec := map[string]any{
		"card":       t.Card,
		"qq":         t.QQ,
		"machine_id": t.MachineID,
		"mid_list":   "",
		"url":        url,
		"r":          r,
		"body":       body,
		"ws":         ws,
		"http":       res.HTTPStatus,
		"trace_id":   res.TraceID,
		"ep":         res.EPHex,
		"used_r":     res.UsedR,
		"raw":        string(res.RawBody),
		"at":         time.Now().Format(time.RFC3339),
	}
	b, _ := json.MarshalIndent(rec, "", "  ")
	p := filepath.Join(d, "last_login.json")
	_ = os.WriteFile(p, b, 0o600)
	return p
}

// selfDecodeWithCandidates 用一批候选 R 尝试解 _ep。
func selfDecodeWithCandidates(epHex string, rs []string) ([]byte, string) {
	for _, r := range rs {
		if r == "" {
			continue
		}
		if p, err := selfEpDecode(epHex, r); err == nil && len(p) > 0 && p[0] == '{' {
			return p, r
		}
	}
	return nil, ""
}
