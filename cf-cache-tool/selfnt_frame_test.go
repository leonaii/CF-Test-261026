package main

import (
	"bytes"
	"testing"
)

// 验证 sub_14003DD30 的 TEA 前置包装：总长恒为 8 的倍数，头字节 = (v6|0xF8)。
func TestSelfNTTEAWrapAlign(t *testing.T) {
	for n := 0; n <= 40; n++ {
		b := make([]byte, n)
		w := selfNTTEAWrap(b)
		if len(w)%8 != 0 {
			t.Fatalf("n=%d wrap len=%d 不是 8 的倍数", n, len(w))
		}
		v6 := (-2 - n) & 7
		if w[0] != byte(v6)|0xF8 {
			t.Fatalf("n=%d 头字节=%#x 期望 %#x", n, w[0], byte(v6)|0xF8)
		}
		cnt := 0
		for i := 1; i < len(w); i++ {
			if w[i] == 0xDC {
				cnt++
			} else {
				break
			}
		}
		if cnt != (v6+2)&0xF {
			t.Fatalf("n=%d 填充=%d 期望 %d", n, cnt, (v6+2)&0xF)
		}
		back, ok := selfNTTEAUnwrap(w)
		if !ok || !bytes.Equal(back, b) {
			t.Fatalf("n=%d unwrap 不一致", n)
		}
	}
}

// 验证最外层总长度、固定前缀，以及 TEA 解密后的 LP(body)‖LP(payload)。
func TestSelfNTFrameHead(t *testing.T) {
	payload := []byte("payload")
	cid := []byte("1234567890abcdef")
	hdr := selfNTHeader(payload, []byte("0123456789abcdef"), []byte("123"))
	f := selfNTFrame(0x2711, []byte("svc"), payload, cid, []byte("123"), hdr)
	if got := selfNTBE32(f); got != uint32(len(f)) {
		t.Fatalf("总长度前缀=%d，实际=%d", got, len(f))
	}
	want := []byte{0, 0, 0, 12, 2, 0, 0, 0, 4, 0, 0, 0, 0, 5, '0'}
	if !bytes.Equal(f[4:19], want) {
		t.Fatalf("外层固定前缀不对\n got=% x\nwant=% x", f[4:19], want)
	}
	if (len(f)-19)%8 != 0 {
		t.Fatalf("TEA 段长度 %d 不是 8 的倍数", len(f)-19)
	}
	wrapped := selfNTTEADecrypt(f[19:])
	plain, ok := selfNTTEAUnwrap(wrapped)
	if !ok {
		t.Fatal("TEA unwrap 失败")
	}
	if !bytes.HasSuffix(plain, selfNTLenPrefixed(payload)) {
		t.Fatal("明文末尾不是 LP(payload)")
	}
	bodyLen := selfNTBE32(plain)
	if int(bodyLen) > len(plain) || bodyLen < 4 {
		t.Fatalf("body 长度非法: %d", bodyLen)
	}
	body := plain[4:bodyLen]
	prefix := append(selfNTBE32Bytes(0x2711), selfNTBE32Bytes(selfNTCmdMagic)...)
	prefix = append(prefix, selfNTBE32Bytes(selfNTLocale)...)
	prefix = append(prefix, 0x02)
	prefix = append(prefix, make([]byte, 11)...)
	if !bytes.HasPrefix(body, prefix) {
		t.Fatalf("body 固定头不对\n got=% x\nwant=% x", body[:len(prefix)], prefix)
	}
}
