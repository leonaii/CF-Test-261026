package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// 零密钥链式 TEA（16 轮、BE 字）的已知向量：
// 第一块是 TEA("12345678")；第二块按样本继续与第一块明文/密文异或。
// 明文 "12345678abcdefgh" → 密文 16a1975f4bce3948b2f558bda28033c1
func TestSelfNTTEAKnownVector(t *testing.T) {
	plain := []byte("12345678abcdefgh")
	want, _ := hex.DecodeString("16a1975f4bce3948b2f558bda28033c1")
	got := selfNTTEAEncrypt(plain)
	if hex.EncodeToString(got) != hex.EncodeToString(want) {
		t.Fatalf("encrypt mismatch\n got=%x\nwant=%x", got, want)
	}
	back := selfNTTEADecrypt(got)
	if string(back) != string(plain) {
		t.Fatalf("roundtrip mismatch: %q", back)
	}
}

func TestSelfNTTEAKeyRoundTrip(t *testing.T) {
	key := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	plain := []byte("12345678abcdefgh")
	enc := selfNTTEAEncryptKey(plain, key)
	if bytes.Equal(enc, selfNTTEAEncrypt(plain)) {
		t.Fatal("非零密钥不应与零密钥输出相同")
	}
	back := selfNTTEADecryptKey(enc, key)
	if !bytes.Equal(back, plain) {
		t.Fatalf("keyed roundtrip mismatch: %x != %x", back, plain)
	}
}
