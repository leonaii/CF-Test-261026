package main

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

// 用合成数据验证地址槽位与 FNV 校验值的补丁逻辑（不依赖真实 exe）。
func TestWuweiPatchExe(t *testing.T) {
	oldSum := make([]byte, 8)
	binary.LittleEndian.PutUint64(oldSum, wuweiFNV64(wuweiStubURL))

	var buf bytes.Buffer
	buf.WriteString("head")
	buf.Write(wuweiEncodeURL(wuweiStubURL))
	buf.WriteString("mid")
	buf.Write([]byte{0x48, 0xb9})
	buf.Write(oldSum)
	buf.WriteString("tail")

	out, err := wuweiPatchExe(buf.Bytes())
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if !bytes.Contains(out, wuweiEncodeURL(wuweiLocalURL)) {
		t.Fatal("补丁后找不到本地地址槽位")
	}
	if len(wuweiLocalURL) != len(wuweiStubURL) {
		t.Fatalf("本地地址长度 %d != 原地址长度 %d", len(wuweiLocalURL), len(wuweiStubURL))
	}

	idx := bytes.Index(out, []byte{0x48, 0xb9})
	if idx < 0 {
		t.Fatal("补丁后找不到校验指令")
	}
	got := binary.LittleEndian.Uint64(out[idx+2 : idx+10])
	if got != wuweiFNV64(wuweiLocalURL) {
		t.Fatalf("校验值没更新：%x", got)
	}
	if strings.Contains(string(out), string(wuweiEncodeURL(wuweiStubURL))) {
		t.Fatal("原地址槽位没被替换")
	}
}

// 版本不认识时要报错，而不是写坏文件。
func TestWuweiPatchExeUnknown(t *testing.T) {
	if _, err := wuweiPatchExe([]byte("not a launcher")); err == nil {
		t.Fatal("期望报错")
	}
}

// WEG 版换端口但共用同一段 24 字节槽位与 FNV 校验常量。
func TestWuweiWGPatchExe(t *testing.T) {
	if len(wuweiWGLocalURL) != len(wuweiStubURL) {
		t.Fatalf("WEG 本地地址长度 %d != 原地址长度 %d", len(wuweiWGLocalURL), len(wuweiStubURL))
	}
	if wuweiWGListenAddr == wuweiListenAddr {
		t.Fatal("WEG 本地端口与老版相同，会互相顶掉")
	}

	oldSum := make([]byte, 8)
	binary.LittleEndian.PutUint64(oldSum, wuweiFNV64(wuweiStubURL))
	var buf bytes.Buffer
	buf.WriteString("head")
	buf.Write(wuweiEncodeURL(wuweiStubURL))
	buf.WriteString("mid")
	buf.Write([]byte{0x48, 0xb9})
	buf.Write(oldSum)
	buf.WriteString("tail")

	out, err := wuweiPatchExeURL(buf.Bytes(), wuweiWGLocalURL)
	if err != nil {
		t.Fatalf("patch: %v", err)
	}
	if !bytes.Contains(out, wuweiEncodeURL(wuweiWGLocalURL)) {
		t.Fatal("补丁后找不到 WEG 本地地址槽位")
	}
	idx := bytes.Index(out, []byte{0x48, 0xb9})
	if idx < 0 {
		t.Fatal("补丁后找不到校验指令")
	}
	got := binary.LittleEndian.Uint64(out[idx+2 : idx+10])
	if got != wuweiFNV64(wuweiWGLocalURL) {
		t.Fatalf("校验值没更新：%x", got)
	}
}

// dump ini（16 键小写）要转成真实卡密服务器那种 6 键大写票据。
func TestWuweiConvertTicket(t *testing.T) {
	dump := "[123456]\r\n" +
		"Token0106=AAAA\r\n" +
		"Token010A=010A\r\n" +
		"Token0143=0143\r\n" +
		"sessionKey=SK\r\n" +
		"uid=999\r\n" +
		"tgpid=888\r\n" +
		"ctid=CTID\r\n" +
		"ct=CTVALUE\r\n" +
		"machine_id=MID\r\n" +
		"sig=SIG\r\n"

	got, qq, err := wuweiConvertTicket(dump, "123456")
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	want := "[123456]\r\n" +
		"Token010A=010A\r\n" +
		"Token0143=0143\r\n" +
		"Sessionkey=SK\r\n" +
		"CT=CTVALUE\r\n" +
		"TGPID=888\r\n" +
		"machine_id=MID"
	if got != want {
		t.Fatalf("转换结果不对\ngot:  %q\nwant: %q", got, want)
	}
	if qq != "123456" {
		t.Fatalf("qq=%s", qq)
	}

	// 已经是真票格式的输入再转一次应该不变（排障用的 real_ticket.ini）
	got2, qq2, err := wuweiConvertTicket(got, "")
	if err != nil || got2 != want || qq2 != "123456" {
		t.Fatalf("不是幂等: err=%v qq=%s got=%q", err, qq2, got2)
	}

	// 缺字段要报错，不能把半张票喂给上号器
	if _, _, err := wuweiConvertTicket("[1]\r\nToken010A=x\r\n", "1"); err == nil {
		t.Fatal("缺字段时应该报错")
	}
}
