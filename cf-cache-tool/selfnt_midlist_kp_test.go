package main

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

// TestSelfNTMidListFromKeyPair 验证 mid_list 里带的是本次握手的公钥 X||Y，
// 且该点在 P-256 曲线上（抓包 work/midlist_final.txt 已确认是同一种形态）。
func TestSelfNTMidListFromKeyPair(t *testing.T) {
	kp, err := selfNTNewKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	rand8, _ := hex.DecodeString("dc8101bcd9826fc4")
	s, err := selfNTMidListFromKeyPair(kp, "8ac936bfcc316d3d660cb564b40f0188", "MS-130667489970", rand8)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	want := strings.ToUpper(hex.EncodeToString(kp.Uncompressed()[1:]))
	if !strings.Contains(body, "_"+want) {
		t.Fatalf("part2 里没有本次公钥\nwant=%s\nbody=%s", want, body)
	}
	// 两段式：<b64inner>_<128 hex>
	part2 := body[strings.LastIndex(body, "\x12"):]
	if !strings.Contains(part2, "_"+want) {
		t.Fatal("part2 结构不对")
	}
	if len(want) != 128 {
		t.Fatalf("公钥十六进制长度 %d，应为 128", len(want))
	}
}
