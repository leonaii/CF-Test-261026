package main

// 可视卡密：界面上不直接显示 16 位真卡密，而是按固定规则渲染成
//
//	CG + 补齐的 0 + 加密串 = 32 位
//
// 真卡密（16 位 A-Z）按 5 bit/字符打包成 10 字节，加上版本字节、分隔字节和
// CRC16 组成 15 字节明文，AES-256-CTR 加密后按自定义 base32 字母表编成 24 位，
// 前面补 6 个 0 凑满 32 位。
//
// 渲染是确定性的：同一张卡每次渲染结果相同，复制出来贴回首屏登录框能还原成
// 真卡密再去上号；上号器、接口、云端目录仍然只用真卡密，不受影响。

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"strings"
)

const (
	cardVisPrefix   = "CG"
	cardVisPadWidth = 6  // 补齐的 0 的个数
	cardVisBodyLen  = 24 // 15 字节 × 8 / 5
	cardVisLen      = len(cardVisPrefix) + cardVisPadWidth + cardVisBodyLen // 32
	cardVisPlainLen = 15
	cardVisPadChar  = "0"
)

// cardVisAlphabet 去掉 I/L/O/U 的 Crockford 风格字母表，避免 0/O、1/I 误读。
const cardVisAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// cardVisSecret 参与密钥派生。**改了历史渲染串就全部失效**，不要动。
const cardVisSecret = "cg.card.vis.v1|5f2b8d41e7a9c3"

var (
	cardVisKey = sha256.Sum256([]byte(cardVisSecret))
	cardVisIV  = [16]byte{'C', 'G', 'c', 'a', 'r', 'd', 'V', 'i', 's', 0x1F, 0x27, 0x5A, 0x93, 0xC4, 0x6E, 0xB8}
)

// renderCardKey 返回卡密在界面上显示的形式：CG + 补齐的 0 + 加密串，共 32 位。
// 不是 16 位 A-Z 的卡密（例如管理员名 leonai）原样返回。
func renderCardKey(card string) string {
	c := strings.ToUpper(strings.TrimSpace(card))
	if !looksLikeCard(c) {
		return card
	}
	sealed := cardVisSeal(c)
	if sealed == nil {
		return card
	}
	body := cardVisBase32Encode(cardVisCrypt(sealed))
	if len(body) != cardVisBodyLen {
		return card
	}
	return cardVisPrefix + strings.Repeat(cardVisPadChar, cardVisPadWidth) + body
}

// decodeCardKey 把 32 位可视卡密还原成 16 位真卡密。
func decodeCardKey(s string) (string, bool) {
	in := strings.ToUpper(strings.TrimSpace(s))
	if len(in) != cardVisLen || !strings.HasPrefix(in, cardVisPrefix) {
		return "", false
	}
	pad := in[len(cardVisPrefix) : len(cardVisPrefix)+cardVisPadWidth]
	if strings.Trim(pad, cardVisPadChar) != "" {
		return "", false
	}
	body := in[len(cardVisPrefix)+cardVisPadWidth:]
	raw, ok := cardVisBase32Decode(body)
	if !ok || len(raw) != cardVisPlainLen {
		return "", false
	}
	return cardVisOpen(cardVisCrypt(raw))
}

// storeCardKey 是卡密在本地文件里的形态：与界面渲染串一致（CG 开头 32 位）。
// cache 目录名、meta.json、history.jsonl、config.json、verify_snapshot.json
// 以及离线列表缓存都只存这个形态，磁盘上不出现明文真卡密。
func storeCardKey(card string) string { return renderCardKey(card) }

// loadCardKey 还原本地文件里的卡密；老版本直接存真卡密的文件原样返回，
// 迁移时由 migrateLocalCardStore 统一改写。
func loadCardKey(s string) string {
	if real, ok := decodeCardKey(s); ok {
		return real
	}
	return s
}

// isCardVisForm 判断输入长得像不像可视卡密（CG 开头、32 位）。
func isCardVisForm(s string) bool {
	t := strings.ToUpper(strings.TrimSpace(s))
	return len(t) == cardVisLen && strings.HasPrefix(t, cardVisPrefix)
}

// resolveCardKey 统一首屏输入：可视卡密还原成真卡密，其它输入原样返回。
// 长得像可视卡密却解不开时返回 ok=false，避免把输错的串当成普通名称。
func resolveCardKey(s string) (string, bool) {
	t := strings.TrimSpace(s)
	if !isCardVisForm(t) {
		return t, true
	}
	return decodeCardKey(t)
}

// cardVisSeal 把 16 位卡密封成 15 字节明文。
func cardVisSeal(card string) []byte {
	if len(card) != 16 {
		return nil
	}
	p := make([]byte, cardVisPlainLen)
	p[0] = 0xC5 // 版本/魔数
	for i := 0; i < 16; i++ {
		cardVisPutBits(p[1:11], i*5, card[i]-'A', 5)
	}
	p[11] = 0x3C
	sum := cardVisCRC16(p[:11])
	p[12] = byte(sum >> 8)
	p[13] = byte(sum)
	p[14] = byte(sum) ^ 0x5A // 尾字节跟着内容变，同时当结束标记
	return p
}

// cardVisOpen 校验并拆出 16 位卡密。
func cardVisOpen(p []byte) (string, bool) {
	if len(p) != cardVisPlainLen || p[0] != 0xC5 || p[11] != 0x3C || p[14] != p[13]^0x5A {
		return "", false
	}
	if want := cardVisCRC16(p[:11]); want != uint16(p[12])<<8|uint16(p[13]) {
		return "", false
	}
	var b strings.Builder
	b.Grow(16)
	for i := 0; i < 16; i++ {
		v := cardVisGetBits(p[1:11], i*5, 5)
		if v > 25 {
			return "", false
		}
		b.WriteByte('A' + v)
	}
	out := b.String()
	if !looksLikeCard(out) {
		return "", false
	}
	return out, true
}

// renderCardInText 把一段文本里出现的 16 位 A-Z 卡密换成可视卡密，只用于显示
// （列表备注、原始备注这类夹着卡密的文本）。整段必须是 16 位，避免切坏其它字段。
func renderCardInText(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		c := s[i]
		if c < 'A' || c > 'Z' {
			b.WriteByte(c)
			i++
			continue
		}
		j := i
		for j < len(s) && s[j] >= 'A' && s[j] <= 'Z' {
			j++
		}
		token := s[i:j]
		if len(token) == 16 {
			b.WriteString(renderCardKey(token))
		} else {
			b.WriteString(token)
		}
		i = j
	}
	return b.String()
}

// cardVisCrypt 用 AES-256-CTR 变换；CTR 是对称的，加密解密同一个函数。
func cardVisCrypt(in []byte) []byte {
	blk, err := aes.NewCipher(cardVisKey[:])
	if err != nil {
		return nil
	}
	out := make([]byte, len(in))
	cipher.NewCTR(blk, cardVisIV[:]).XORKeyStream(out, in)
	return out
}

func cardVisPutBits(dst []byte, bitPos int, v uint8, n int) {
	for i := n - 1; i >= 0; i-- {
		if v&(1<<uint(i)) != 0 {
			dst[bitPos/8] |= 1 << uint(7-bitPos%8)
		}
		bitPos++
	}
}

func cardVisGetBits(src []byte, bitPos, n int) uint8 {
	var v uint8
	for i := 0; i < n; i++ {
		v <<= 1
		if src[bitPos/8]&(1<<uint(7-bitPos%8)) != 0 {
			v |= 1
		}
		bitPos++
	}
	return v
}

func cardVisCRC16(b []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, x := range b {
		crc ^= uint16(x) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = crc<<1 ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

func cardVisBase32Encode(b []byte) string {
	var sb strings.Builder
	var acc uint32
	bits := 0
	for _, x := range b {
		acc = acc<<8 | uint32(x)
		bits += 8
		for bits >= 5 {
			bits -= 5
			sb.WriteByte(cardVisAlphabet[(acc>>uint(bits))&31])
		}
	}
	if bits > 0 {
		sb.WriteByte(cardVisAlphabet[(acc<<uint(5-bits))&31])
	}
	return sb.String()
}

func cardVisBase32Decode(s string) ([]byte, bool) {
	out := make([]byte, 0, len(s)*5/8)
	var acc uint32
	bits := 0
	for i := 0; i < len(s); i++ {
		v := strings.IndexByte(cardVisAlphabet, s[i])
		if v < 0 {
			return nil, false
		}
		acc = acc<<5 | uint32(v)
		bits += 5
		if bits >= 8 {
			bits -= 8
			out = append(out, byte(acc>>uint(bits)))
		}
	}
	return out, true
}
