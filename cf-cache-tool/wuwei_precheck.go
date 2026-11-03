package main

// 无畏上号预检：拉起上号器之前只看本机环境，不猜账号类型、不弹框拦人。
//
//   - 本机没有 %LOCALAPPDATA%\Riot Games，说明这台机器还没成功登录过一次无畏。
//   - 首次登录要拉 riot-client.dyn.riotcdn.net 的 ToS，连不上同样会 VAL 29。

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type wuweiPrecheck struct {
	Proxy   string // ACLOS 代理路径，空=没找到
	HasRiot bool   // 本机有 Riot 登录缓存
}

// wuweiRunPrecheck 只做本地快速检查，不联网、不看账号备注。
func wuweiRunPrecheck() *wuweiPrecheck {
	return &wuweiPrecheck{
		Proxy:   wuweiDetectProxy(),
		HasRiot: wuweiHasRiotCache(),
	}
}

// summary 是写进日志的一行预检结果。
func (p *wuweiPrecheck) summary() string {
	proxy := "未找到 ACLOS 代理"
	if p.Proxy != "" {
		proxy = "ACLOS 代理已找到"
	}
	riot := "本机没有 Riot 登录缓存（还没成功登录过无畏）"
	if p.HasRiot {
		riot = "本机有 Riot 登录缓存"
	}
	return proxy + " · " + riot
}

func wuweiHasRiotCache() bool {
	for _, env := range []string{"LOCALAPPDATA", "APPDATA"} {
		base := strings.TrimSpace(os.Getenv(env))
		if base == "" {
			continue
		}
		if st, err := os.Stat(filepath.Join(base, "Riot Games")); err == nil && st.IsDir() {
			return true
		}
	}
	return false
}

// wuweiRiotCDNReachable 查 Riot ToS 域名，首次登录/取名流程要用到。
func wuweiRiotCDNReachable() error {
	c, err := net.DialTimeout("tcp", "riot-client.dyn.riotcdn.net:443", 2*time.Second)
	if err != nil {
		return err
	}
	return c.Close()
}
