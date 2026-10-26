package main

// 构建期字符串混淆：build.sh 用 obfstr.go 把凭据 encode 成 base64+XOR，
// 生成的 zz_build_secret.go 里只有密文，strings 搜不到明文。

import "encoding/base64"

const secretMask = "cgx0r"

// deob 还原 build.sh 混淆过的字符串；解码失败返回空串（构建异常时不会崩）。
func deob(s string) string {
	if s == "" {
		return ""
	}
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	k := []byte(secretMask)
	for i := range raw {
		raw[i] ^= k[i%len(k)]
	}
	return string(raw)
}
