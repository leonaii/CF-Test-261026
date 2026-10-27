package main

import (
	"bytes"
	"crypto/ecdh"
	"encoding/hex"
	"testing"
)

// TestSelfNTServerPubOnCurve 确认样本里那把硬编码公钥确实是合法 P-256 点。
func TestSelfNTServerPubOnCurve(t *testing.T) {
	pub, err := selfNTServerPub()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ecdh.P256().NewPublicKey(pub); err != nil {
		t.Fatalf("内置服务端公钥不在 P-256 曲线上: %v", err)
	}
	t.Log("server pub ok:", hex.EncodeToString(pub[:8]))
}

// TestSelfNTKexRoundTrip 验证 KeyExchange 的 CBC+PKCS7 封装可逆。
func TestSelfNTKexRoundTrip(t *testing.T) {
	kp, err := selfNTNewKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	plain := selfNTKexPlain("130667489970499", 1793000000)
	key, ct, err := selfNTSealKex(kp, plain)
	if err != nil {
		t.Fatal(err)
	}
	if len(key) != 16 {
		t.Fatalf("key 长度 %d", len(key))
	}
	if len(ct)%16 != 0 {
		t.Fatalf("密文长度 %d 不是 16 的倍数", len(ct))
	}
	got, err := selfNTOpenKex(key, ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("解回不一致\n got=%s\nwant=%s", got, plain)
	}
}

func TestSelfNTComputerIDMD5(t *testing.T) {
	got := selfNTComputerID(2510904083, nil)
	want, _ := hex.DecodeString("1663760f4b70e6c0c237ced2129c05d8")
	if !bytes.Equal(got, want) {
		t.Fatalf("ComputerID mismatch\n got=%x\nwant=%x", got, want)
	}
	if len(got) != 16 {
		t.Fatalf("ComputerID 长度应为 16，实际 %d", len(got))
	}
}
