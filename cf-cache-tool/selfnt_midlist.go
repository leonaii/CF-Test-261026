package main

// NT 第三层：mid_list 生成（复刻 sub_14004FB70，IDA 反编译 + 抓包逐字节对齐）。
//
// 内层 70 字节设备块（抓包 midlist 解码，偏移见注释）：
//   0x00 int32 1
//   0x04 int32 0
//   0x08 int32 2
//   0x0C int32 0
//   0x10 8 字节随机（样本 BCryptGenRandom(8, flags=2)）
//   0x18 int32 6
//   0x1C int32 1
//   0x20 int32 7601        （Windows 6.1.7601 的 build 号，样本硬编码）
//   0x24 UTF-16LE 设备名    （"MS-" + mid 前 12 位）
//   末尾 int32 0
//
// 外层 protobuf（抓包 "CiQIABIg..." 解码）：
//   field1 { field1: 0, field2: <machine_id 小写 32 字节> }
//   field1 { field1: 1, field2: "<base64(内层)>_<hex>_<hex>" }
// 整包再做一次 base64，就是 credential.mid_list。

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

// selfNTAppendVarint 追加 protobuf varint。
func selfNTAppendVarint(b []byte, v uint64) []byte {
	for v >= 0x80 {
		b = append(b, byte(v)|0x80)
		v >>= 7
	}
	return append(b, byte(v))
}

// selfNTAppendBytes 追加 protobuf 的 length-delimited 字段（不含 tag）。
func selfNTAppendBytes(b, data []byte) []byte {
	b = selfNTAppendVarint(b, uint64(len(data)))
	return append(b, data...)
}

// selfNTMidListInner 生成 70 字节设备块。rand8 为空时用固定值（仅自检用）。
func selfNTMidListInner(devName string, rand8 []byte) []byte {
	out := make([]byte, 0, 64+len(devName)*2)
	put32 := func(v uint32) {
		var t [4]byte
		binary.LittleEndian.PutUint32(t[:], v)
		out = append(out, t[:]...)
	}
	put32(1)
	put32(0)
	put32(2)
	put32(0)
	if len(rand8) == 8 {
		out = append(out, rand8...)
	} else {
		out = append(out, 0, 0, 0, 0, 0, 0, 0, 0)
	}
	put32(6)
	put32(1)
	put32(7601) // 样本 .text 0x140050367 直接写 7601
	for _, r := range devName {
		var t [2]byte
		binary.LittleEndian.PutUint16(t[:], uint16(r))
		out = append(out, t[:]...)
	}
	put32(0)
	return out
}

// selfNTMidList 生成完整 mid_list。
// keyHex/ctHex 是握手阶段的两段 32 字节材料（十六进制大写），
// 与样本 "…_<hex>_<hex>" 的三段式保持一致。
func selfNTMidList(machineID, devName string, rand8, keyHex, ctHex []byte) (string, error) {
	return selfNTMidListWithPart2(machineID, devName, rand8, func(pubXY []byte) string {
		return strings.ToUpper(hex.EncodeToString(keyHex)) + "_" +
			strings.ToUpper(hex.EncodeToString(ctHex))
	})
}

// selfNTMidListFromKeyPair 用本次握手的密钥对生成 mid_list：
// part2 = base64(内层) + "_" + hex(公钥 X||Y)。
// 抓包验证：work/midlist_final.txt 里那 64 字节在 P-256 曲线上（true），
// 说明服务端就是拿它和硬编码私钥做 ECDH 解 mcode。
func selfNTMidListFromKeyPair(kp *selfNTKeyPair, machineID, devName string, rand8 []byte) (string, error) {
	if kp == nil {
		return "", errors.New("密钥对为空")
	}
	pub := kp.Uncompressed() // 65 字节 0x04||X||Y
	return selfNTMidListWithPart2(machineID, devName, rand8, func([]byte) string {
		return strings.ToUpper(hex.EncodeToString(pub[1:]))
	})
}

// selfNTMidListWithPart2 负责外层 protobuf 组包，part2 的后半段由回调决定。
func selfNTMidListWithPart2(machineID, devName string, rand8 []byte, tail func(pubXY []byte) string) (string, error) {
	mid := strings.ToLower(strings.TrimSpace(machineID))
	if len(mid) != 32 {
		return "", errors.New("machine_id 必须是 32 位十六进制")
	}
	if _, err := hex.DecodeString(mid); err != nil {
		return "", errors.New("machine_id 不是合法十六进制")
	}
	inner := selfNTMidListInner(devName, rand8)
	part2 := base64.StdEncoding.EncodeToString(inner) + "_" + tail(nil)

	// field1: {1: 0, 2: mid}
	f1 := selfNTAppendVarint(nil, 0x08)
	f1 = selfNTAppendVarint(f1, 0)
	f1 = selfNTAppendVarint(f1, 0x12)
	f1 = selfNTAppendBytes(f1, []byte(mid))

	// field1: {1: 1, 2: part2}
	f2 := selfNTAppendVarint(nil, 0x08)
	f2 = selfNTAppendVarint(f2, 1)
	f2 = selfNTAppendVarint(f2, 0x12)
	f2 = selfNTAppendBytes(f2, []byte(part2))

	out := selfNTAppendVarint(nil, 0x0A)
	out = selfNTAppendBytes(out, f1)
	out = selfNTAppendVarint(out, 0x0A)
	out = selfNTAppendBytes(out, f2)
	return base64.StdEncoding.EncodeToString(out), nil
}
