package main

// CG gs：用网页 ck 查腾讯侧账号状态。
//
// 用法：
//   CG gs -ck "uin=o0123456; skey=@xxxx"
//   CG gs -qq 123456 -skey @xxxx -pskey yyyy
//   CG gs -ck "..." -json          机器可读输出
//   CG gs -ck "..." -raw           附带每条接口的原始响应（排查用）
//
// 注意：处罚/冻结/在线/登录保护 都需要真正的网页 skey；
// 离线 ini 里的 Token010A/0143/0106 不能当 skey 用（实测全部 auth failed）。

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
)

func cmdGs(args []string) error {
	fs := flag.NewFlagSet("gs", flag.ExitOnError)
	ck := fs.String("ck", "", "完整 cookie（uin=o0xxx; skey=@xxx; p_skey=xxx）")
	qq := fs.String("qq", "", "QQ 号（-ck 里没带 uin 时必填）")
	skey := fs.String("skey", "", "skey（不用 -ck 时单独传）")
	pskey := fs.String("pskey", "", "p_skey（可选）")
	asJSON := fs.Bool("json", false, "输出 JSON")
	withRaw := fs.Bool("raw", false, "JSON 里附带原始响应")
	_ = fs.Parse(args)

	cookieStr := strings.TrimSpace(*ck)
	if cookieStr == "" {
		if strings.TrimSpace(*skey) == "" {
			return fmt.Errorf("需要 -ck，或同时给 -qq 和 -skey")
		}
		c := ntCookie(*qq, *skey, *pskey)
		if !c.valid() {
			return fmt.Errorf("-qq 必须是纯数字 QQ 号")
		}
		cookieStr = c.header()
	}

	c := parseGSCookie(cookieStr, *qq)
	if c.QQ == "" && strings.TrimSpace(*qq) != "" {
		c.QQ = gsDigits(*qq)
	}
	if !c.valid() {
		return fmt.Errorf("cookie 里没解析出 QQ 或 skey，原串：%s", gsClip(cookieStr, 80))
	}

	rep := gsQueryAll(c)

	if *asJSON {
		res := map[string]any{
			"qq":            rep.QQ,
			"punish_ok":     rep.PunishDone,
			"punish_list":   rep.Punish,
			"frozen_ok":     rep.FrozenDone,
			"is_frozen":     rep.Frozen,
			"online_ok":     rep.OnlineDone,
			"is_online":     rep.Online,
			"protect_ok":    rep.ProtectDone,
			"protected":     rep.Protect,
			"antiaddict_ok": rep.AntiAddictDone,
			"antiaddict":    rep.AntiAddict,
			"errors":        rep.Errs,
		}
		if *withRaw {
			res["raw"] = rep.Raw
		}
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	}

	fmt.Print(rep.Text())
	return nil
}
