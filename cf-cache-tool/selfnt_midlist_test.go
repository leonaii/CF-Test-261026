package main

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

// 抓包里那条 mid_list 的内层 70 字节（work/midlist_final.txt 解出来的原样）。
const selfNTCapturedInnerHex = "01000000000000000200000000000000dc8101bcd9826fc4" +
	"0600000001000000b11d0000" +
	"4d0053002d00310033003000360036003700340038003900390037003000" +
	"00000000"

// TestSelfNTMidListInnerMatchesCapture 逐字节对齐抓包。
func TestSelfNTMidListInnerMatchesCapture(t *testing.T) {
	rand8, _ := hex.DecodeString("dc8101bcd9826fc4")
	got := selfNTMidListInner("MS-130667489970", rand8)
	want, _ := hex.DecodeString(selfNTCapturedInnerHex)
	if hex.EncodeToString(got) != hex.EncodeToString(want) {
		t.Fatalf("内层不一致\n got=%x\nwant=%x", got, want)
	}
	if len(got) != 70 {
		t.Fatalf("内层长度 %d，抓包是 70", len(got))
	}
}

// TestSelfNTMidListWrapper 校验外层是
// base64( field1{0,mid} + field1{1,"<b64inner>_<hex>_<hex>"} )。
func TestSelfNTMidListWrapper(t *testing.T) {
	mid := "8ac936bfcc316d3d660cb564b40f0188"
	rand8, _ := hex.DecodeString("dc8101bcd9826fc4")
	key := make([]byte, 32)
	ct := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
		ct[i] = byte(0xff - i)
	}
	s, err := selfNTMidList(mid, "MS-130667489970", rand8, key, ct)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	// 头两个字节必须是 0A 24（field1, len=36 = 4 + 32）
	if len(raw) < 2 || raw[0] != 0x0A || raw[1] != 0x24 {
		t.Fatalf("外层头不对: %x", raw[:4])
	}
	// 第二段里必须出现 "..._<64HEX>_<64HEX>"
	body := string(raw)
	if !strings.Contains(body, "_"+strings.ToUpper(hex.EncodeToString(key))+"_"+strings.ToUpper(hex.EncodeToString(ct))) {
		t.Fatal("两段十六进制没拼进 part2")
	}
	if !strings.Contains(body, base64.StdEncoding.EncodeToString(selfNTMidListInner("MS-130667489970", rand8))[:20]) {
		t.Fatal("内层 base64 没拼进 part2")
	}
}
