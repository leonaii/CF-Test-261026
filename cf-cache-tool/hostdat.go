package main

// host.dat v2 编解码。
//
// 文件结构：
//   0x00  magic  d8 3a 91 c7
//   0x04  ver    02
//   0x05  iv     16 字节
//   0x15  len    uint16，密文长度
//   0x17  ct     len 字节（AES-256-CBC/PKCS7）
//   0x17+len     hmac 32 字节（HMAC-SHA256(macKey, 文件头+密文)）
//
// 明文 = nibbleSwap(aesDecrypt(ct) ^ xorMask)
// 密钥由二进制内固定种子派生：SHA256(seed + "aes-host-v2" / "xor-host-v2" / "mac-host-v2")
// 三者均从原程序 RVA 0x8E300 / 0x8F550 还原，HMAC 校验通过验证。

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

const (
	hostDatMagic = 0xC7913AD8
	hostDatVer   = 2
	hostDatMin   = 0x17 + 16 + 32
)

// 以下三张表来自原程序 .rdata（RVA 0x1400FB770 / 0x1400FB780 / 0x1400FB790），
// 用于生成固定种子，属于程序常量而非密钥泄露。
var (
	tblT3 = [4]uint32{1579979188, 150197882, 2979938501, 3597831953}
	tblT1 = [4]uint32{835758967, 3838351874, 2065037016, 2427077838}
	tblT2 = [4]uint32{2709278240, 1798560031, 3236829348, 1509063443}
)

// deriveSeed 复刻 RVA 0x8E300 中的 32 字节种子生成循环。
func deriveSeed() []byte {
	seed := make([]byte, 32)
	seed[0] = 0x92
	b, a := uint32(1), uint32(3)
	r8 := uint32(2)
	t := uint32(0xE4C89A02)
	esi := uint32(0x0D)
	for i := 1; i < 32; i++ {
		x := r8 & 3
		tOld := t
		r8 = a & 3
		t = tblT1[x]
		v := (tblT2[b&3] >> ((b & 3) * 8)) ^ (t >> ((a & 3) * 8))
		v ^= tblT3[a&3] >> ((^b & 3) * 8)
		v = ((b * b) ^ 0x5A + esi + v)
		v ^= tblT2[(b-1)&3] >> 11
		v ^= tOld >> 3
		seed[i] = byte(v)
		a++
		b++
		esi += 0x0D
	}
	return seed
}

func hostKeys() (aesKey, xorKey, macKey []byte) {
	seed := deriveSeed()
	derive := func(label string) []byte {
		h := sha256.New()
		h.Write(seed)
		h.Write([]byte(label))
		return h.Sum(nil)
	}
	return derive("aes-host-v2"), derive("xor-host-v2"), derive("mac-host-v2")
}

func xorMaskAt(xorKey []byte, i int) byte {
	return xorKey[(3+7*i)&31] ^ byte((int(xorKey[i&31])+17*i)&0xFF)
}

func nibbleSwap(b []byte) {
	for i := range b {
		b[i] = b[i]<<4 | b[i]>>4
	}
}

func pkcs7Unpad(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, errors.New("空明文")
	}
	n := int(b[len(b)-1])
	if n == 0 || n > aes.BlockSize || n > len(b) {
		return nil, errors.New("PKCS7 填充非法")
	}
	for _, v := range b[len(b)-n:] {
		if int(v) != n {
			return nil, errors.New("PKCS7 填充非法")
		}
	}
	return b[:len(b)-n], nil
}

func pkcs7Pad(b []byte) []byte {
	n := aes.BlockSize - len(b)%aes.BlockSize
	return append(b, bytes.Repeat([]byte{byte(n)}, n)...)
}

// DecryptHostDat 读取并还原 host.dat 中的服务器地址。
func DecryptHostDat(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(raw) < hostDatMin {
		return "", fmt.Errorf("host.dat 长度不足: %d", len(raw))
	}
	if binary.LittleEndian.Uint32(raw[0:4]) != hostDatMagic {
		return "", errors.New("host.dat 魔数不匹配")
	}
	if raw[4] != hostDatVer {
		return "", fmt.Errorf("host.dat 版本不支持: %d", raw[4])
	}
	iv := raw[5:0x15]
	n := int(binary.LittleEndian.Uint16(raw[0x15:0x17]))
	if len(raw) < 0x17+n+32 {
		return "", errors.New("host.dat 密文/校验长度越界")
	}
	ct := raw[0x17 : 0x17+n]
	tag := raw[0x17+n : 0x17+n+32]

	aesKey, xorKey, macKey := hostKeys()
	mac := hmac.New(sha256.New, macKey)
	mac.Write(raw[:0x17+n])
	if !hmac.Equal(mac.Sum(nil), tag) {
		return "", errors.New("host.dat HMAC 校验失败（文件被改动或密钥不匹配）")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	plain := make([]byte, len(ct))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, ct)
	plain, err = pkcs7Unpad(plain)
	if err != nil {
		return "", err
	}
	for i := range plain {
		plain[i] ^= xorMaskAt(xorKey, i)
	}
	nibbleSwap(plain)
	return string(plain), nil
}

// EncryptHostDat 生成一个指向 host 的合法 host.dat。
func EncryptHostDat(host string) ([]byte, error) {
	if host == "" {
		return nil, errors.New("host 不能为空")
	}
	aesKey, xorKey, macKey := hostKeys()

	plain := []byte(host)
	nibbleSwap(plain)
	for i := range plain {
		plain[i] ^= xorMaskAt(xorKey, i)
	}
	plain = pkcs7Pad(plain)

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	ct := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ct, plain)

	raw := make([]byte, 0, 0x17+len(ct)+32)
	var head [0x17]byte
	binary.LittleEndian.PutUint32(head[0:4], hostDatMagic)
	head[4] = hostDatVer
	copy(head[5:0x15], iv)
	binary.LittleEndian.PutUint16(head[0x15:0x17], uint16(len(ct)))
	raw = append(raw, head[:]...)
	raw = append(raw, ct...)

	mac := hmac.New(sha256.New, macKey)
	mac.Write(raw)
	raw = append(raw, mac.Sum(nil)...)
	return raw, nil
}

// HostDatManager 负责改写/还原 host.dat，并保留原始服务器地址。
type HostDatManager struct {
	Path string
	Bak  string
}

func NewHostDatManager(path string) *HostDatManager {
	return &HostDatManager{Path: path, Bak: path + ".orig"}
}

// Patch 将 host.dat 指向本地监听地址。重复执行安全。
func (m *HostDatManager) Patch(listen string) (original string, changed bool, err error) {
	original, err = DecryptHostDat(m.Path)
	if err != nil {
		return "", false, err
	}
	local := "http://" + listen
	if original == local {
		return original, false, nil
	}
	// 只有当前不是本地地址时才备份，避免把本地地址误存成上游。
	if _, statErr := os.Stat(m.Bak); os.IsNotExist(statErr) {
		src, readErr := os.ReadFile(m.Path)
		if readErr != nil {
			return "", false, readErr
		}
		if writeErr := os.WriteFile(m.Bak, src, 0o600); writeErr != nil {
			return "", false, writeErr
		}
	}
	data, err := EncryptHostDat(local)
	if err != nil {
		return "", false, err
	}
	if err := os.WriteFile(m.Path, data, 0o600); err != nil {
		return "", false, err
	}
	return original, true, nil
}

// Restore 用备份覆盖 host.dat。
func (m *HostDatManager) Restore() error {
	data, err := os.ReadFile(m.Bak)
	if err != nil {
		return err
	}
	return os.WriteFile(m.Path, data, 0o600)
}

// Upstream 返回真实上游地址：当前文件不是本地地址时以它为准，否则回退到备份。
func (m *HostDatManager) Upstream(listen string) (string, error) {
	local := "http://" + listen
	u, err := DecryptHostDat(m.Path)
	if err == nil && u != "" && u != local {
		return u, nil
	}
	if _, err := os.Stat(m.Bak); err == nil {
		if backup, err := DecryptHostDat(m.Bak); err == nil && backup != "" {
			return backup, nil
		}
	}
	if err != nil {
		return "", err
	}
	return "", errors.New("host.dat 已指向本地，但缺少 host.dat.orig 备份，无法确定真实服务器")
}
