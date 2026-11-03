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
