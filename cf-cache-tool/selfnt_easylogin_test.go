package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestSelfNTSealTickets 校验 0x01‖GCM‖0x00 封装可逆，且帧头尾正确。
func TestSelfNTSealTickets(t *testing.T) {
	shared, _ := hex.DecodeString("983930e02b8bce2b61d775a83d82e39929b25c0cf60bacd12aff8a3238a882e5")
	ticketA := "aq_base64_ticket_A"
	ticketB, _ := hex.DecodeString("d0e82e34023ae5e0")

	payload, err := selfNTSealTickets(shared, ticketA, ticketB)
	if err != nil {
		t.Fatal(err)
	}
	if payload[0] != 0x01 || payload[len(payload)-1] != 0x00 {
		t.Fatalf("载荷首尾不对: %x ... %x", payload[0], payload[len(payload)-1])
	}
	// 1 + 12(nonce) + 密文 + 16(tag) + 1
	if len(payload) < 1+12+16+1 {
		t.Fatalf("载荷太短: %d", len(payload))
	}
	gotA, gotB, err := selfNTOpenTickets(shared, payload)
	if err != nil {
		t.Fatal(err)
	}
	if gotA != ticketA || !bytes.Equal(gotB, ticketB) {
		t.Fatalf("解回不一致: %q / %x", gotA, gotB)
	}
}

// TestSelfNTEasyDevice 校验设备块的字段顺序与硬编码值。
func TestSelfNTEasyDevice(t *testing.T) {
	dev := selfNTEasyDevice([]byte{0x01, 0xAA, 0x00})
	s := string(dev)
	for _, want := range []string{"Windows", "Administrator", "6.1.7601", "MS-20220623RAMD"} {
		if !bytes.Contains(dev, []byte(want)) {
			t.Fatalf("设备块缺少 %q: %x", want, dev)
		}
	}
	// 字段 3 是 varint 4：tag 0x18 + 0x04
	if !bytes.Contains(dev, []byte{0x18, 0x04}) {
		t.Fatal("字段3 应为 4")
	}
	_ = s
}

func selfNTTestBytesField(out []byte, field int, data []byte) []byte {
	out = selfNTAppendVarint(out, uint64(field)<<3|2)
	return selfNTAppendBytes(out, data)
}

func selfNTTestPlain(payload []byte) []byte {
	body := selfNTLenPrefixed([]byte{0})
	out := selfNTLenPrefixed(body)
	return append(out, selfNTLenPrefixed(payload)...)
}

func TestSelfNTParseEasyLoginResponse(t *testing.T) {
	ticketA := make([]byte, 32)
	for i := range ticketA {
		ticketA[i] = byte(i + 1)
	}
	token := []byte{0xaa, 0xbb}
	tgt := []byte{1, 2, 3, 4, 5}
	d2 := make([]byte, 16)
	d2k := make([]byte, 16)
	sm := selfNTTestBytesField(nil, 3, token)
	sm = selfNTTestBytesField(sm, 4, tgt)
	sm = selfNTTestBytesField(sm, 5, d2)
	sm = selfNTTestBytesField(sm, 6, d2k)
	body := selfNTTestBytesField(nil, 1, sm)
	dec := selfNTTestBytesField(nil, 1, nil)
	dec = selfNTTestBytesField(dec, 2, body)
	enc, err := selfNTGCMSeal(ticketA, dec)
	if err != nil {
		t.Fatal(err)
	}
	plain := selfNTTestPlain(selfNTTestBytesField(nil, 3, enc))
	sess, err := selfNTParseEasyLoginResponse(plain, ticketA)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(sess.Token0106, token) || !bytes.Equal(sess.TGT, tgt) || !bytes.Equal(sess.D2, d2) || !bytes.Equal(sess.D2K, d2k) {
		t.Fatalf("会话字段解析不一致: %+v", sess)
	}
}

func TestSelfNTFirstAuthResponse(t *testing.T) {
	const sig = "56135b1187d955f6035e16e48c102209"
	plain := selfNTTestPlain(selfNTTestBytesField(nil, 4, []byte(sig)))
	got, err := selfNTParseFirstAuthResponse(plain)
	if err != nil {
		t.Fatal(err)
	}
	if got != "56135B1187D955F6035E16E48C102209" {
		t.Fatalf("sig mismatch: %s", got)
	}
}

func TestSelfNTFirstAuthFrameRoundTrip(t *testing.T) {
	guid := make([]byte, 16)
	tgt := []byte{1, 2, 3, 4}
	d2 := make([]byte, 88)
	d2k := make([]byte, 16)
	frame, err := selfNTFirstAuthFrame(1306674899, []byte(hex.EncodeToString(guid)), "1793000000", "1306674899", tgt, d2, d2k, 1600001614)
	if err != nil {
		t.Fatal(err)
	}
	plain, err := selfNTDecodeFrameKey(frame, d2k)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(plain, []byte("QQConnectLogin.auth")) {
		t.Fatal("缺少 FirstAuth 服务名")
	}
	payload, err := selfNTExtractPayload(plain)
	if err != nil {
		t.Fatal(err)
	}
	login, err := selfNTKexBytesField(payload, 65)
	if err != nil {
		t.Fatal(err)
	}
	got, err := selfNTKexBytesField(login, 2)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != hex.EncodeToString(tgt) {
		t.Fatalf("tgt hex mismatch: %s", got)
	}
}
