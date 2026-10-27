package main

import (
	"strconv"
	"time"
)

// NT SSO 的 trpc 帧（逐条复刻 sub_140059330）。
//
// 所有 LP(x) 都表示 BE32(len(x)+4) ‖ x，长度包含自身 4 字节：
//   sub_14002F600  = LP(vector)
//   sub_14002FFC0  = LP(raw pointer)
//   sub_14002C420  = BE32
//   sub_14002E150  = 裸追加
//
// 组帧顺序：
//   body  = BE32(cmd) ‖ BE32(0x200471C7) ‖ BE32(2052)
//         ‖ 02 00 00 00 00 00 00 00 00 00 00 00
//         ‖ LP("") ‖ LP(service) ‖ LP("") ‖ LP(a5) ‖ LP("")
//         ‖ 00 07 "8<:.8" ‖ LP(header)
//   plain = LP(body) ‖ LP(payload)
//   outer = BE32(12) ‖ 02 ‖ LP("") ‖ 00 ‖ LP("0") ‖ TEA(TEAWrap(plain))
//   frame = LP(outer)
//
// 这里最后的 LP(outer) 就是 sub_14003E5A0 接收时读取的 BE32 总长度。
const (
	selfNTFrameMagic5 = "\x38\x3c\x3a\x2e\x38" // "8<:.8"
	selfNTCmdMagic    = uint32(0x200471C7)
	selfNTLocale      = uint32(2052)
)

// selfNTLenPrefixed 是 sub_14002F600/sub_14002FFC0 的共同格式。
func selfNTLenPrefixed(d []byte) []byte {
	out := make([]byte, 4, 4+len(d))
	selfNTPutBE32(out, uint32(len(d)+4))
	return append(out, d...)
}

func selfNTBE32Bytes(v uint32) []byte {
	out := make([]byte, 4)
	selfNTPutBE32(out, v)
	return out
}

// selfNTCallFrame 保留旧入口；没有 ComputerID 时 a5 为空。
func selfNTCallFrame(cmd uint32, svc string, payload []byte, connSeq string) []byte {
	return selfNTCallFrameWithA5(cmd, svc, payload, nil, connSeq)
}

// selfNTCallFrameWithA5 按 0x14006A495 的实参组帧。
// a5 是 ComputerID；connSeq 为空时取当前 Unix 秒（样本 0x140069CC9 起）。
func selfNTCallFrameWithA5(cmd uint32, svc string, payload, a5 []byte, connSeq string) []byte {
	if connSeq == "" {
		connSeq = strconv.FormatInt(time.Now().Unix(), 10)
	}
	hdr := selfNTHeader(payload, []byte(selfNTRandHeader15()), []byte(connSeq))
	return selfNTFrame(cmd, []byte(svc), payload, a5, []byte(connSeq), hdr)
}

// selfNTRandHeader15 复刻 sub_140030060：
// "700-" + hex(16 随机字节) + "-" + hex(8 随机字节) + "-01"。
func selfNTRandHeader15() string {
	hexOf := func(b []byte) string {
		const hexdig = "0123456789abcdef"
		out := make([]byte, len(b)*2)
		for i, v := range b {
			out[i*2] = hexdig[v>>4]
			out[i*2+1] = hexdig[v&0x0f]
		}
		return string(out)
	}
	a := make([]byte, 16)
	b := make([]byte, 8)
	_, _ = randRead(a)
	_, _ = randRead(b)
	return "700-" + hexOf(a) + "-" + hexOf(b) + "-01"
}

// selfNTHeader 生成 header protobuf。
func selfNTHeader(field13, field15, connSeq []byte) []byte {
	out := make([]byte, 0, 256)
	if len(field13) > 0 {
		out = selfNTAppendVarint(out, 13<<3|2)
		out = selfNTAppendBytes(out, field13)
	}
	if len(field15) > 0 {
		out = selfNTAppendVarint(out, 15<<3|2)
		out = selfNTAppendBytes(out, field15)
	}

	kv := make([]byte, 0, 64)
	kv = selfNTAppendVarint(kv, 1<<3|2)
	kv = selfNTAppendBytes(kv, []byte("client_conn_seq"))
	kv = selfNTAppendVarint(kv, 2<<3|2)
	kv = selfNTAppendBytes(kv, connSeq)
	out = selfNTAppendVarint(out, 23<<3|2)
	out = selfNTAppendBytes(out, kv)

	for _, pair := range [][2]string{{"nt_core", "nt_pc_sdk"}, {"os", "win"}} {
		m := make([]byte, 0, 48)
		m = selfNTAppendVarint(m, 1<<3|2)
		m = selfNTAppendBytes(m, []byte(pair[0]))
		m = selfNTAppendVarint(m, 2<<3|2)
		m = selfNTAppendBytes(m, []byte(pair[1]))
		out = selfNTAppendVarint(out, 23<<3|2)
		out = selfNTAppendBytes(out, m)
	}
	return out
}

// selfNTFrame 组装一帧。a3=服务名，a4=载荷，a5=ComputerID，a6=connSeq。
func selfNTFrame(cmd uint32, a3, a4, a5, connSeq []byte, header []byte) []byte {
	_ = connSeq // connSeq 已在 header 中，保留形参以对应样本 a6

	body := make([]byte, 0, 256+len(a3)+len(a5)+len(header))
	body = append(body, selfNTBE32Bytes(cmd)...)
	body = append(body, selfNTBE32Bytes(selfNTCmdMagic)...)
	body = append(body, selfNTBE32Bytes(selfNTLocale)...)
	body = append(body, 0x02)
	body = append(body, make([]byte, 11)...)
	body = append(body, selfNTLenPrefixed(nil)...)
	body = append(body, selfNTLenPrefixed(a3)...)
	body = append(body, selfNTLenPrefixed(nil)...)
	body = append(body, selfNTLenPrefixed(a5)...)
	body = append(body, selfNTLenPrefixed(nil)...)
	body = append(body, 0x00, 0x07)
	body = append(body, selfNTFrameMagic5...)
	body = append(body, selfNTLenPrefixed(header)...)

	plain := selfNTLenPrefixed(body)
	plain = append(plain, selfNTLenPrefixed(a4)...)
	sealed := selfNTTEAEncrypt(selfNTTEAWrap(plain))

	outer := make([]byte, 0, 15+len(sealed))
	outer = append(outer, selfNTBE32Bytes(12)...)
	outer = append(outer, 0x02)
	outer = append(outer, selfNTLenPrefixed(nil)...)
	outer = append(outer, 0x00)
	outer = append(outer, selfNTLenPrefixed([]byte("0"))...)
	outer = append(outer, sealed...)
	return selfNTLenPrefixed(outer)
}

// selfNTTEAWrap 复刻 sub_14003DD30 的 TEA 前置包装（不是 PKCS7）。
//
//	v6 = (-2 - len) & 7
//	byte(v6 | 0xF8) ‖ (v6+2) 个 0xDC ‖ 原文 ‖ 7 个 0x00
func selfNTTEAWrap(b []byte) []byte {
	v6 := (-2 - len(b)) & 7
	n := (v6 + 2) & 0xF
	out := make([]byte, 0, 1+n+len(b)+7)
	out = append(out, byte(v6)|0xF8)
	for i := 0; i < n; i++ {
		out = append(out, 0xDC)
	}
	out = append(out, b...)
	for i := 0; i < 7; i++ {
		out = append(out, 0x00)
	}
	return out
}

// selfNTTEAUnwrap 是 sub_140030E30 的解封逻辑。
// 服务端填充可以是随机字节，因此不能像客户端那样要求 0xDC；
// 样本只使用 start=(首字节&7)+3，并去掉末尾 7 个零。
func selfNTTEAUnwrap(b []byte) ([]byte, bool) {
	if len(b) < 10 || len(b)%8 != 0 {
		return nil, false
	}
	start := int(b[0]&7) + 3
	end := len(b) - 7
	if start > end {
		return nil, false
	}
	for _, c := range b[end:] {
		if c != 0 {
			return nil, false
		}
	}
	out := make([]byte, end-start)
	copy(out, b[start:end])
	return out, true
}
