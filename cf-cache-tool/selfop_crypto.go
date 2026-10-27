package main

// 自研上号：WeGame UnionAuth 的传输层加密
//
// 请求：
//   R      = 22 位十进制随机串
//   key    = MD5(R)
//   iv     = MD5(key ‖ R)
//   body   = {"_eq": hex(AES-128-CBC(明文JSON, key, iv))}
//   sign   = HMAC-SHA256(key=R, data=body ‖ extBase)
//   _ws    = hex(AES-128-CBC(extBase+"&sign="+sign, 固定key, 固定IV))
//
// 响应：
//   _ep    = AES-128-CBC 密文，key=MD5(R'), iv=MD5(key' ‖ R')

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var (
	selfWsKey = selfMustHex("5FCE2F0E2EE844DA9A7147AFEBC940EF")
	selfWsIV  = selfMustHex("AB3DE93F7E1C7653FE5F3F0269B33634")
)

func selfMustHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("selfop: 非法 hex 常量: " + s)
	}
	return b
}

func selfPKCS7Pad(b []byte, size int) []byte {
	n := size - len(b)%size
	out := make([]byte, len(b)+n)
	copy(out, b)
	for i := len(b); i < len(out); i++ {
		out[i] = byte(n)
	}
	return out
}

func selfPKCS7Unpad(b []byte, size int) ([]byte, error) {
	if len(b) == 0 || len(b)%size != 0 {
		return nil, errors.New("密文长度不是块大小的整数倍")
	}
	n := int(b[len(b)-1])
	if n == 0 || n > size || n > len(b) {
		return nil, errors.New("padding 非法")
	}
	for _, c := range b[len(b)-n:] {
		if int(c) != n {
			return nil, errors.New("padding 非法")
		}
	}
	return b[:len(b)-n], nil
}

func aesCBCEncrypt(key, iv, plain []byte) ([]byte, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src := selfPKCS7Pad(plain, blk.BlockSize())
	dst := make([]byte, len(src))
	cipher.NewCBCEncrypter(blk, iv).CryptBlocks(dst, src)
	return dst, nil
}

func aesCBCDecrypt(key, iv, ct []byte) ([]byte, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ct) == 0 || len(ct)%blk.BlockSize() != 0 {
		return nil, errors.New("密文长度非法")
	}
	dst := make([]byte, len(ct))
	cipher.NewCBCDecrypter(blk, iv).CryptBlocks(dst, ct)
	return selfPKCS7Unpad(dst, blk.BlockSize())
}

// selfKDF 由 22 位随机串派生 AES key/iv。
func selfKDF(r string) (key, iv []byte) {
	k := md5.Sum([]byte(r))
	key = k[:]
	h := md5.New()
	h.Write(key)
	h.Write([]byte(r))
	iv = h.Sum(nil)
	return
}

// selfEpDerivedR 复刻样本的确定性 R：decimal(m) + decimal(hash)。
// m 由 timestamp 和 0x1400F97E8 表派生，滚动哈希使用 64 位自然溢出。
func selfEpDerivedR(timestamp int64, input string) string {
	table := [...]uint32{0x58880A08, 0x6C4B6A28, 0x0DB2D93B}
	x := uint32(timestamp) ^ 0xD2AD4E72
	r := (x - 0x5E295AB4) ^ table[(x+0x56)%3]
	m := uint64(0x5600000000) + uint64(r)
	var h uint64
	for i := 0; i < len(input); i++ {
		h = h*m + uint64(input[i])
	}
	h &= 0x7FFFFFFF
	return strconv.FormatUint(m, 10) + strconv.FormatUint(h, 10)
}

// selfEpRequestR 对应 sub_140038080/sub_140050CE0：
// 请求签名和 _ep 响应都使用 hash(C10+nonce+machineID+"0")。
func selfEpRequestR(timestamp int64, nonce, machineID string) string {
	return selfEpDerivedR(timestamp, "C10"+nonce+strings.ToUpper(machineID)+"0")
}

func selfEpResponseR(timestamp int64, nonce, machineID string) string {
	return selfEpRequestR(timestamp, nonce, machineID)
}

// selfRandR 生成 22 位十进制随机串。
func selfRandR() (string, error) {
	limit := new(big.Int).Exp(big.NewInt(10), big.NewInt(22), nil)
	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return "", err
	}
	s := n.String()
	for len(s) < 22 {
		s = "0" + s
	}
	return s, nil
}

// selfEqEncode 把明文 JSON 编码成 {"_eq":"..."} 并给出 R。
func selfEqEncode(plain []byte) (body string, r string, err error) {
	r, err = selfRandR()
	if err != nil {
		return "", "", err
	}
	body, err = selfEqEncodeWithR(plain, r)
	return body, r, err
}

func selfEqEncodeWithR(plain []byte, r string) (string, error) {
	key, iv := selfKDF(r)
	ct, err := aesCBCEncrypt(key, iv, plain)
	if err != nil {
		return "", err
	}
	return `{"_eq":"` + hexUpper(ct) + `"}`, nil
}

// selfSign 计算 HMAC-SHA256(R, body‖extBase)。
func selfSign(r, body, extBase string) string {
	m := hmac.New(sha256.New, []byte(r))
	m.Write([]byte(body))
	m.Write([]byte(extBase))
	return hexUpper(m.Sum(nil))
}

// selfWsEncode 把 ext 串 + sign 封装成 URL 里的 _ws 参数。
func selfWsEncode(wsPlain string) (string, error) {
	ct, err := aesCBCEncrypt(selfWsKey, selfWsIV, []byte(wsPlain))
	if err != nil {
		return "", err
	}
	return hexUpper(ct), nil
}

// selfEpDecode 用给定 R 解密响应体 _ep。
func selfEpDecode(epHex, r string) ([]byte, error) {
	ct, err := hex.DecodeString(epHex)
	if err != nil {
		return nil, fmt.Errorf("_ep 不是十六进制: %w", err)
	}
	key, iv := selfKDF(r)
	return aesCBCDecrypt(key, iv, ct)
}

func randRead(b []byte) (int, error) { return rand.Read(b) }

func hexUpper(b []byte) string {
	const tab = "0123456789ABCDEF"
	out := make([]byte, len(b)*2)
	for i, c := range b {
		out[i*2] = tab[c>>4]
		out[i*2+1] = tab[c&0x0f]
	}
	return string(out)
}
