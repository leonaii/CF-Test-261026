package main

// 自研上号的启动自检：用抓包向量校验 KDF / _ws / mid_list，避免把错误带到线上。

import (
	"encoding/hex"
	"strings"
)

func selfSelfCheck() []string {
	var out []string

	key, iv := selfKDF("3702734984951639852597")
	if hexUpper(key) == "EC0190C8800F66D330D648BC95522838" && hexUpper(iv) == "C355247721D24F1DE8736ED91841C5AD" {
		out = append(out, "KDF(MD5/IV) ✓")
	} else {
		out = append(out, "KDF(MD5/IV) ✗")
	}

	if plain, err := selfEpDecode(
		"fec60ec4af786980b3aea293085abf5214510e9940f8668cd999aa7562bc53af29be2ba3250f438df412953d9b7910ed063d3424a6ed06b693503a6c285a5506",
		"370273498495512170951",
	); err == nil && strings.Contains(string(plain), "8000150") {
		out = append(out, "_ep 解密 ✓")
	} else {
		out = append(out, "_ep 解密 ✗")
	}

	if mid, err := buildMidList("8AC936BFCC316D3D660CB564B40F0188"); err == nil && mid == selfMidListTemplate {
		out = append(out, "mid_list 结构 ✓")
	} else {
		out = append(out, "mid_list 结构 ✗")
	}

	const probe = "ext=probe&gid=55555"
	if ws, err := selfWsEncode(probe); err == nil {
		if ct, e := hex.DecodeString(ws); e == nil {
			if p, e2 := aesCBCDecrypt(selfWsKey, selfWsIV, ct); e2 == nil && string(p) == probe {
				out = append(out, "_ws 封装 ✓")
			} else {
				out = append(out, "_ws 封装 ✗")
			}
		}
	} else {
		out = append(out, "_ws 封装 ✗")
	}

	return out
}
