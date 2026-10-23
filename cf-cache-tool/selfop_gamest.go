package main

// 读取 / 复用 TCLS\CF_Gamest（上号器写下的明文票据），
// 让 CG 可以在不依赖 _ep 的前提下拿到 openid/open_key/ctt 并自行启动游戏。

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// selfReadCFGamest 解析游戏目录下的 CF_Gamest。
func selfReadCFGamest(dir string) (*selfGameTicketData, string, error) {
	if dir == "" {
		var err error
		dir, err = selfGameDir()
		if err != nil {
			return nil, "", err
		}
	}
	p := filepath.Join(dir, selfGameRelDir, "CF_Gamest")
	raw, err := os.ReadFile(p)
	if err != nil {
		return nil, p, fmt.Errorf("读不到 CF_Gamest（需先用上号器登录一次）: %w", err)
	}
	t := &selfGameTicketData{}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimRight(line, "\r")
		i := strings.Index(line, "=")
		if i <= 0 {
			continue
		}
		k, v := strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:])
		switch k {
		case "uin":
			t.Uin = v
		case "tgpid":
			t.Tgpid = v
		case "openid":
			t.OpenID = v
		case "open_key":
			t.OpenKey = v
		case "pf_key":
			t.PfKey = v
		case "game_client_key":
			t.CTT = v
		case "game_st":
			t.GameST = v
		case "game_st_key":
			t.GameSTKey = v
		case "client_key":
			t.ClientKey = v
		case "ctt":
			if t.CTT == "" {
				t.CTT = v
			}
		}
	}
	if t.OpenID == "" || t.OpenKey == "" {
		return nil, p, errors.New("CF_Gamest 里缺少 openid/open_key（上号器未登录成功）")
	}
	return t, p, nil
}
