package main

// 用抓包取得的真实向量校验自研上号的密码学实现。

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestSelfKDFVector(t *testing.T) {
	key, iv := selfKDF("3702734984951639852597")
	if got, want := hexUpper(key), "EC0190C8800F66D330D648BC95522838"; got != want {
		t.Fatalf("key= %s want %s", got, want)
	}
	if got, want := hexUpper(iv), "C355247721D24F1DE8736ED91841C5AD"; got != want {
		t.Fatalf("iv= %s want %s", got, want)
	}
}

func TestSelfEpDecodeVector(t *testing.T) {
	ep := "fec60ec4af786980b3aea293085abf5214510e9940f8668cd999aa7562bc53af29be2ba3250f438df412953d9b7910ed063d3424a6ed06b693503a6c285a5506"
	plain, err := selfEpDecode(ep, "370273498495512170951")
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	if !strings.Contains(string(plain), "8000150") {
		t.Fatalf("明文不含预期错误码: %s", plain)
	}
	t.Logf("_ep 明文: %s", plain)
}

func TestSelfEpDerivedRVector(t *testing.T) {
	const machineID = "CC815BEE32490F9CA4F6CC7F116E4469"
	const want = "3702800218351216143717"
	if got := selfEpRequestR(1793097317, "BXYBPD", machineID); got != want {
		t.Fatalf("request R=%s want %s", got, want)
	}
	if got := selfEpResponseR(1793097317, "BXYBPD", machineID); got != want {
		t.Fatalf("response R=%s want %s", got, want)
	}
}

func TestSelfEpDecodeDerivedVector(t *testing.T) {
	const machineID = "CC815BEE32490F9CA4F6CC7F116E4469"
	ep := "7160cf42b8a6e4fd97397388ea5ab70822f3695446f5d18bbe841a0f03ee62c3caad4079a8711db67ec469f8e39ac8033430ec366e68f5d082a3d5072d8964fc"
	r := selfEpResponseR(1793095598, "PZBZVN", machineID)
	plain, err := selfEpDecode(ep, r)
	if err != nil {
		t.Fatalf("派生 R 解密失败: %v", err)
	}
	if !strings.Contains(string(plain), "8000150") || !strings.Contains(string(plain), "not match") {
		t.Fatalf("明文不符合预期: %s", plain)
	}
}

func TestBuildMidListMatchesCapture(t *testing.T) {
	got, err := buildMidList("8AC936BFCC316D3D660CB564B40F0188")
	if err != nil {
		t.Fatal(err)
	}
	if got != selfMidListTemplate {
		t.Fatalf("mid_list 与抓包不一致\n got=%s\nwant=%s", got, selfMidListTemplate)
	}
}

func TestWsRoundTrip(t *testing.T) {
	const want = "ext=8AC936BFCC316D3D660CB564B40F0188&gid=55555&nonce=UVLTDM&session_id=0&timestamp=1792732663&uid=0&version=C10&sign=31F077DA"
	ws, err := selfWsEncode(want)
	if err != nil {
		t.Fatal(err)
	}
	ct, err := hex.DecodeString(ws)
	if err != nil {
		t.Fatal(err)
	}
	plain, err := aesCBCDecrypt(selfWsKey, selfWsIV, ct)
	if err != nil {
		t.Fatal(err)
	}
	if string(plain) != want {
		t.Fatalf("_ws 往返不一致: %s", plain)
	}
}

func TestEqEncodeRoundTrip(t *testing.T) {
	const sample = `{"mappid":10001,"clienttype":1000000}`
	body, r, err := selfEqEncode([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(body, `{"_eq":"`) || !strings.HasSuffix(body, `"}`) {
		t.Fatalf("body 外壳错误: %s", body)
	}
	epHex := strings.TrimSuffix(strings.TrimPrefix(body, `{"_eq":"`), `"}`)
	key, iv := selfKDF(r)
	ct, err := hex.DecodeString(epHex)
	if err != nil {
		t.Fatal(err)
	}
	plain, err := aesCBCDecrypt(key, iv, ct)
	if err != nil {
		t.Fatal(err)
	}
	if string(plain) != sample {
		t.Fatalf("_eq 往返不一致: %s", plain)
	}
}
