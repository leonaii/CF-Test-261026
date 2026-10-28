//go:build ignore

package main

// 构建机工具：维护 CG 根目录的只读公开分享。
// 只在 build.sh 中通过 go run share_provision.go 执行，不进入客户端二进制。
//
// 输出 .cg-share.env，供构建脚本生成只读 token/口令常量。

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ocsOrigin   = "https://cloud.51i.cc"
	ocsShares   = ocsOrigin + "/ocs/v2.php/apps/files_sharing/api/v1/shares"
	sharePath   = "/CG"
	shareNote   = "CG readonly client"
	statePath   = ".cg-share.env"
	publicDav   = ocsOrigin + "/public.php/webdav/"
	shareFolder = "CG"
)

type shareState struct {
	ID    string
	Token string
	Pass  string
}

type shareData struct {
	ID          json.RawMessage `json:"id"`
	Token       string          `json:"token"`
	ShareType   int             `json:"share_type"`
	Permissions int             `json:"permissions"`
	Expiration  string          `json:"expiration"`
	Note        string          `json:"note"`
}

var httpClient = &http.Client{Timeout: 30 * time.Second}

func main() {
	user := strings.TrimSpace(os.Getenv("DAV_USER"))
	pass := strings.TrimSpace(os.Getenv("DAV_PASS"))
	if user == "" || pass == "" {
		fail(errors.New("缺少 DAV_USER / DAV_PASS"))
	}

	st, _ := loadState()
	var data shareData

	if st.ID != "" {
		d, err := getShare(user, pass, st.ID)
		if err != nil {
			fail(fmt.Errorf("读取现有只读分享失败，拒绝自动更换口令: %w", err))
		}
		data = d
		if d.Token != "" {
			st.Token = d.Token
		}
		if st.Pass == "" {
			st.Pass = newPassword()
			if err := setShareField(user, pass, st.ID, "password", st.Pass); err != nil {
				fail(fmt.Errorf("设置只读分享口令失败: %w", err))
			}
		} else if !verifyPublic(st.Token, st.Pass) {
			fail(errors.New("现有只读分享口令自检失败；为避免让已发布客户端失效，本次不自动轮换口令"))
		}
	}

	if data.ID == nil {
		if d, ok := findExistingShare(user, pass); ok {
			data = d
			st.ID = rawString(d.ID)
			st.Token = d.Token
			st.Pass = newPassword()
			if err := setShareField(user, pass, st.ID, "password", st.Pass); err != nil {
				fail(fmt.Errorf("设置只读分享口令失败: %w", err))
			}
		}
	}

	if data.ID == nil {
		st.Pass = newPassword()
		d, err := createShare(user, pass, st.Pass)
		if err != nil {
			fail(fmt.Errorf("创建只读分享失败: %w", err))
		}
		data = d
		st.ID = rawString(d.ID)
		st.Token = d.Token
	}

	if st.ID == "" || st.Token == "" || st.Pass == "" {
		fail(errors.New("只读分享返回数据不完整"))
	}

	if err := setShareField(user, pass, st.ID, "permissions", "1"); err != nil {
		fail(fmt.Errorf("设置只读权限失败: %w", err))
	}
	if err := setShareField(user, pass, st.ID, "expireDate", ""); err != nil {
		fail(fmt.Errorf("取消分享有效期失败: %w", err))
	}
	if err := setShareField(user, pass, st.ID, "note", shareNote); err != nil {
		fail(fmt.Errorf("写入分享备注失败: %w", err))
	}
	if !verifyPublic(st.Token, st.Pass) {
		fail(errors.New("只读分享公开 WebDAV 自检失败"))
	}
	if err := saveState(st); err != nil {
		fail(fmt.Errorf("保存分享配置失败: %w", err))
	}

	fmt.Printf("==> [share] 只读分享已就绪（id=%s，无有效期，已设访问密码）\n", st.ID)
}

func loadState() (shareState, error) {
	var st shareState
	raw, err := os.ReadFile(statePath)
	if err != nil {
		return st, err
	}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(k) {
		case "CG_SHARE_ID":
			st.ID = strings.TrimSpace(v)
		case "CG_SHARE_TOKEN":
			st.Token = strings.TrimSpace(v)
		case "CG_SHARE_PASS":
			st.Pass = strings.TrimSpace(v)
		}
	}
	return st, nil
}

func saveState(st shareState) error {
	body := strings.Join([]string{
		"# build.sh 自动生成，只读分享凭据",
		"CG_SHARE_ID=" + st.ID,
		"CG_SHARE_TOKEN=" + st.Token,
		"CG_SHARE_PASS=" + st.Pass,
		"",
	}, "\n")
	return os.WriteFile(statePath, []byte(body), 0o600)
}

func newPassword() string {
	b := make([]byte, 18)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func findExistingShare(user, pass string) (shareData, bool) {
	u := ocsShares + "?path=" + url.QueryEscape(sharePath) + "&format=json"
	env, err := ocsRequest(http.MethodGet, u, user, pass, nil)
	if err != nil {
		return shareData{}, false
	}
	var list []shareData
	if err := json.Unmarshal(env.OCS.Data, &list); err != nil {
		return shareData{}, false
	}
	for _, d := range list {
		if d.ShareType == 3 && d.Note == shareNote && rawString(d.ID) != "" && d.Token != "" {
			return d, true
		}
	}
	return shareData{}, false
}

func getShare(user, pass, id string) (shareData, error) {
	u := ocsShares + "/" + url.PathEscape(id) + "?format=json"
	env, err := ocsRequest(http.MethodGet, u, user, pass, nil)
	if err != nil {
		return shareData{}, err
	}
	var d shareData
	if err := json.Unmarshal(env.OCS.Data, &d); err != nil || rawString(d.ID) == "" {
		var list []shareData
		if e := json.Unmarshal(env.OCS.Data, &list); e != nil || len(list) == 0 {
			return shareData{}, errors.New("分享详情为空或格式错误")
		}
		d = list[0]
	}
	if rawString(d.ID) == "" {
		return shareData{}, errors.New("分享详情为空")
	}
	return d, nil
}

func createShare(user, pass, sharePass string) (shareData, error) {
	form := url.Values{
		"path":        {sharePath},
		"shareType":   {"3"},
		"permissions": {"1"},
		"password":    {sharePass},
		"note":        {shareNote},
	}
	env, err := ocsRequest(http.MethodPost, ocsShares+"?format=json", user, pass, form)
	if err != nil {
		return shareData{}, err
	}
	var d shareData
	if err := json.Unmarshal(env.OCS.Data, &d); err != nil {
		return shareData{}, err
	}
	if rawString(d.ID) == "" || d.Token == "" {
		return shareData{}, errors.New("创建结果缺少 id/token")
	}
	return d, nil
}

func setShareField(user, pass, id, key, value string) error {
	form := url.Values{key: {value}}
	u := ocsShares + "/" + url.PathEscape(id) + "?format=json"
	_, err := ocsRequest(http.MethodPut, u, user, pass, form)
	return err
}

type ocsEnvelope struct {
	OCS struct {
		Meta struct {
			StatusCode any    `json:"statuscode"`
			Message    string `json:"message"`
		} `json:"meta"`
		Data json.RawMessage `json:"data"`
	} `json:"ocs"`
}

func ocsRequest(method, endpoint, user, pass string, form url.Values) (*ocsEnvelope, error) {
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, pass)
	req.Header.Set("OCS-APIRequest", "true")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "CG-build/"+shareFolder)
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("OCS HTTP %d: %s", resp.StatusCode, clip(string(raw), 160))
	}
	var env ocsEnvelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("OCS 返回格式错误: %s", clip(string(raw), 160))
	}
	code := anyString(env.OCS.Meta.StatusCode)
	if code != "100" && code != "200" {
		return nil, fmt.Errorf("OCS 状态 %s: %s", code, env.OCS.Meta.Message)
	}
	return &env, nil
}

func verifyPublic(token, pass string) bool {
	req, err := http.NewRequest("PROPFIND", publicDav, nil)
	if err != nil {
		return false
	}
	req.SetBasicAuth(token, pass)
	req.Header.Set("Depth", "0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "CG-build/"+shareFolder)
	resp, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	return resp.StatusCode == http.StatusMultiStatus || resp.StatusCode == http.StatusOK
}

func rawString(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" {
		return ""
	}
	if strings.HasPrefix(s, `"`) {
		var v string
		if json.Unmarshal(raw, &v) == nil {
			return v
		}
	}
	return s
}

func anyString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case json.Number:
		return x.String()
	default:
		return fmt.Sprint(x)
	}
}

func clip(s string, n int) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "错误:", err)
	os.Exit(1)
}
