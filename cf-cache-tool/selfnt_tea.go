package main

// NT 组帧的 TEA 层（复刻 sub_14003DD30）。
//
// 反编译 + 汇编确认：
//
//	v6 = (-2 - len) & 7                 补位到 8 字节
//	v10 = qword_1400F7409               固定 8 字节 XOR 掩码 —— 实测全 0
//	key = xmm6                          sub_140059330 里 `pxor xmm6,xmm6` 后未再装载 —— 实测全 0
//	do { sum -= 0x61C88647; … } while (sum != 0xE3779B90);   16 次迭代 = 标准 TEA
//
// 即：零密钥 TEA，块内两个 32 位字按大端读写，无掩码。

const selfNTTEADelta uint32 = 0x9E3779B9

// selfNTTEAEncrypt 加密（len 必须是 8 的倍数）。
//
// 这不是 ECB。按 sub_14003DD30 的 0x14003DE8E～0x14003DFE1：
//
//	X[i] = P[i] XOR O[i-1]
//	O[i] = TEA(X[i]) XOR X[i-1]
//
// 初始 O[0] = qword_1400F7409 = 0，X[0] = 0；密钥也是 16 字节全零。
func selfNTTEAEncrypt(b []byte) []byte {
	return selfNTTEAEncryptKey(b, nil)
}

// selfNTTEAEncryptKey 是带 128 位密钥的版本；nil 表示组帧使用的零密钥。
func selfNTTEAEncryptKey(b, key []byte) []byte {
	var k [4]uint32
	if len(key) >= 16 {
		for i := 0; i < 4; i++ {
			k[i] = selfNTBE32(key[i*4:])
		}
	}
	out := make([]byte, 0, len(b))
	var prevOut, prevX [8]byte
	for i := 0; i+8 <= len(b); i += 8 {
		var x [8]byte
		for j := 0; j < 8; j++ {
			x[j] = b[i+j] ^ prevOut[j]
		}
		v0 := selfNTBE32(x[0:])
		v1 := selfNTBE32(x[4:])
		var sum uint32
		for r := 0; r < 16; r++ {
			sum += selfNTTEADelta
			v0 += (v1 + sum) ^ (k[1] + (v1 >> 5)) ^ (k[0] + (v1 << 4))
			v1 += (v0 + sum) ^ (k[2] + (v0 << 4)) ^ (k[3] + (v0 >> 5))
		}
		var enc, o [8]byte
		selfNTPutBE32(enc[0:], v0)
		selfNTPutBE32(enc[4:], v1)
		for j := 0; j < 8; j++ {
			o[j] = enc[j] ^ prevX[j]
		}
		out = append(out, o[:]...)
		prevX = x
		prevOut = o
	}
	return out
}

// selfNTTEADecrypt 是 selfNTTEAEncrypt 的逆运算。
func selfNTTEADecrypt(b []byte) []byte {
	return selfNTTEADecryptKey(b, nil)
}

func selfNTTEADecryptKey(b, key []byte) []byte {
	var k [4]uint32
	if len(key) >= 16 {
		for i := 0; i < 4; i++ {
			k[i] = selfNTBE32(key[i*4:])
		}
	}
	out := make([]byte, 0, len(b))
	var prevOut, prevX [8]byte
	for i := 0; i+8 <= len(b); i += 8 {
		var o [8]byte
		copy(o[:], b[i:i+8])
		var enc [8]byte
		for j := 0; j < 8; j++ {
			enc[j] = o[j] ^ prevX[j]
		}
		v0 := selfNTBE32(enc[0:])
		v1 := selfNTBE32(enc[4:])
		var sum uint32
		for r := 0; r < 16; r++ {
			sum += selfNTTEADelta
		}
		for r := 0; r < 16; r++ {
			v1 -= (v0 + sum) ^ (k[2] + (v0 << 4)) ^ (k[3] + (v0 >> 5))
			v0 -= (v1 + sum) ^ (k[0] + (v1 << 4)) ^ (k[1] + (v1 >> 5))
			sum -= selfNTTEADelta
		}
		var x, p [8]byte
		selfNTPutBE32(x[0:], v0)
		selfNTPutBE32(x[4:], v1)
		for j := 0; j < 8; j++ {
			p[j] = x[j] ^ prevOut[j]
		}
		out = append(out, p[:]...)
		prevX = x
		prevOut = o
	}
	return out
}

func selfNTBE32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func selfNTPutBE32(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}
