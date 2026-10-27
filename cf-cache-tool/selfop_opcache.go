package main

// 真上号器不是"一次请求定生死"：它把 OP 票据落盘（op.json / last_op.json /
// qwh_cf_data.ini），网络抖动时重试，实在不行才回退到上一份票据。
// 这里按同样的模式做：内存/磁盘缓存 + 退避重试 + CF_Gamest 兜底（仅限同一个 QQ）。

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// selfOPCacheTTL 是 OP 缓存有效期。真上号器缓存到下次登录前都复用，
// 这里保守取 2 小时，过期就重新走 LoginByQq。
const selfOPCacheTTL = 2 * time.Hour

type selfOPCache struct {
	QQ        string `json:"qq"`
	OpenID    string `json:"openid"`
	OpenKey   string `json:"open_key"`
	CTT       string `json:"ctt"`
	GameST    string `json:"game_st"`
	GameSTKey string `json:"game_st_key"`
	PfKey     string `json:"pf_key"`
	At        int64  `json:"at"`
}

func selfOPCachePath(dir, qq string) string {
	return filepath.Join(dir, "selfop", "op_"+qq+".json")
}

func selfLoadOPCache(dir, qq string) (*selfOPCache, bool) {
	if strings.TrimSpace(qq) == "" || strings.TrimSpace(dir) == "" {
		return nil, false
	}
	raw, err := os.ReadFile(selfOPCachePath(dir, qq))
	if err != nil {
		return nil, false
	}
	var c selfOPCache
	if json.Unmarshal(raw, &c) != nil {
		return nil, false
	}
	if c.QQ != qq || c.OpenID == "" {
		return nil, false
	}
	if time.Since(time.UnixMilli(c.At)) > selfOPCacheTTL {
		return nil, false
	}
	return &c, true
}

func selfSaveOPCache(dir string, qq string, d *selfGameTicketData) {
	if d == nil || d.OpenID == "" || strings.TrimSpace(qq) == "" || strings.TrimSpace(dir) == "" {
		return
	}
	p := selfOPCachePath(dir, qq)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return
	}
	b, err := json.Marshal(selfOPCache{
		QQ:        qq,
		OpenID:    d.OpenID,
		OpenKey:   d.OpenKey,
		CTT:       d.CTT,
		GameST:    d.GameST,
		GameSTKey: d.GameSTKey,
		PfKey:     d.PfKey,
		At:        time.Now().UnixMilli(),
	})
	if err != nil {
		return
	}
	_ = os.WriteFile(p, b, 0o600)
}

func (c *selfOPCache) data(tk *ymmTicket) *selfGameTicketData {
	d := &selfGameTicketData{
		OpenID:    c.OpenID,
		OpenKey:   c.OpenKey,
		CTT:       c.CTT,
		GameST:    c.GameST,
		GameSTKey: c.GameSTKey,
		PfKey:     c.PfKey,
		Uin:       c.QQ,
	}
	if tk != nil {
		d.Uin = firstNonEmpty(d.Uin, tk.QQ)
		d.Tgpid = firstNonEmpty(tk.Tgpid, tk.UID)
		d.ClientKey = tk.ClientKey
	}
	return d
}

// selfRetryableErr 判断错误值不值得重试：网络类重试，业务类（有错误码/没 sig）直接放弃。
func selfRetryableErr(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	if strings.Contains(s, "code=") || strings.Contains(s, "sig") || strings.Contains(s, "mcode") {
		return false
	}
	for _, k := range []string{
		"请求失败", "timeout", "timed out", "eof", "connection", "reset",
		"no such host", "http 5", "http 429", "temporary", "网络",
	} {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}

// selfOpLoginWithRetry 复刻上号器的容错顺序：
// 缓存 → LoginByQq（最多 3 次退避重试）→ 游戏目录里的 CF_Gamest（必须同一个 QQ）。
func (a *guiApp) selfOpLoginWithRetry(tk *ymmTicket, qq string, step func(string, ...any)) (*selfGameTicketData, error) {
	if c, ok := selfLoadOPCache(a.cacheDir, qq); ok {
		step("命中 OP 缓存（%s 前），跳过登录", time.Since(time.UnixMilli(c.At)).Truncate(time.Minute))
		return c.data(tk), nil
	}

	var lastErr error
	for i := 1; i <= 3; i++ {
		data, err := selfLoginByIni(tk, a.cacheDir)
		if err == nil {
			selfSaveOPCache(a.cacheDir, qq, data)
			return data, nil
		}
		lastErr = err
		if !selfRetryableErr(err) {
			break
		}
		if i < 3 {
			step("登录失败（%v），%d 秒后重试", err, i)
			time.Sleep(time.Duration(i) * time.Second)
		}
	}

	// 兜底：真上号器写过的 CF_Gamest，但必须确认是同一个 QQ，不能拿别人的票启动
	if g, _, err2 := selfReadCFGamest(""); err2 == nil && g != nil {
		if uin := strings.TrimSpace(g.Uin); uin != "" && uin == qq && g.OpenID != "" {
			step("改用游戏目录里的 CF_Gamest（同一个 QQ）")
			return g, nil
		}
	}
	if lastErr == nil {
		lastErr = errors.New("登录失败")
	}
	return nil, lastErr
}
