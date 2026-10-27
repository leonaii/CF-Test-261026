package main

// NT 登录第一层：ECDH P-256 密钥协商。
//
// 逐字节复刻样本 sub_1400B3960（IDA 反编译 + 汇编确认）：
//   1. 服务端公钥必须是 65 字节未压缩点，首字节 0x04；
//   2. 导入 CNG 的 ECCPUBLICBLOB（magic "ECK1" + cbKey=32 + X + Y）；
//   3. BCryptSecretAgreement(本机私钥, 对端公钥)；
//   4. BCryptDeriveKey(secret, "TRUNCATE" = BCRYPT_KDF_RAW_SECRET, 32 字节)；
//   5. 把 32 字节整体反转（样本里那段 for 循环就是 reverse）。
//
// 第 4、5 步已用 Windows CNG 做过 5 组随机对比，结果与 crypto/ecdh 完全一致，
// 见 selfnt_ecdh_test.go。

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"fmt"
)

// selfNTKeyPair 是一次 NT 登录握手的本机 P-256 密钥对。
type selfNTKeyPair struct {
	priv *ecdh.PrivateKey
}

// selfNTNewKeyPair 生成握手用的 P-256 密钥对。
func selfNTNewKeyPair() (*selfNTKeyPair, error) {
	p, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成 P-256 密钥对: %w", err)
	}
	return &selfNTKeyPair{priv: p}, nil
}

// Uncompressed 返回 65 字节 0x04||X||Y，对应样本里 a5==65 && *a4==4 的校验。
func (k *selfNTKeyPair) Uncompressed() []byte {
	return k.priv.PublicKey().Bytes()
}

// PublicBLOB 返回 CNG 的 ECCPUBLICBLOB：magic "ECK1" + cbKey=32 + X + Y。
// 样本里是把这段直接当 BCryptImportKeyPair 的输入（pbInput = 0x20314B4345）。
func (k *selfNTKeyPair) PublicBLOB() []byte {
	pub := k.priv.PublicKey().Bytes()
	out := make([]byte, 0, 8+64)
	out = append(out, 'E', 'C', 'K', '1', 32, 0, 0, 0)
	out = append(out, pub[1:]...)
	return out
}

// selfNTParsePeerPoint 把服务端公钥归一成 65 字节未压缩点。
// 样本只接受 65 字节 0x04 开头；这里额外容忍 64 字节裸 X||Y，方便调试。
func selfNTParsePeerPoint(b []byte) ([]byte, error) {
	switch {
	case len(b) == 65 && b[0] == 4:
		return b, nil
	case len(b) == 64:
		out := make([]byte, 65)
		out[0] = 4
		copy(out[1:], b)
		return out, nil
	}
	return nil, fmt.Errorf("对端公钥长度/格式不对（len=%d，首字节=%#x）", len(b), firstByte(b))
}

func firstByte(b []byte) byte {
	if len(b) == 0 {
		return 0
	}
	return b[0]
}

// selfNTSharedSecret 计算 reverse(ECDH(A, B))，与样本的
// BCryptDeriveKey(RAW_SECRET) + 反转 完全一致（32 字节）。
func selfNTSharedSecret(priv *ecdh.PrivateKey, peerPoint []byte) ([]byte, error) {
	if priv == nil {
		return nil, errors.New("本机私钥为空")
	}
	pt, err := selfNTParsePeerPoint(peerPoint)
	if err != nil {
		return nil, err
	}
	pub, err := ecdh.P256().NewPublicKey(pt)
	if err != nil {
		return nil, fmt.Errorf("对端公钥不是合法 P-256 点: %w", err)
	}
	secret, err := priv.ECDH(pub)
	if err != nil {
		return nil, fmt.Errorf("ECDH 协商失败: %w", err)
	}
	// 样本在 BCryptDeriveKey 之后逐字节反转
	out := make([]byte, len(secret))
	for i := range secret {
		out[i] = secret[len(secret)-1-i]
	}
	return out, nil
}

// SharedSecret 是对外入口：用本机密钥对 + 服务端公钥拿到 32 字节会话密钥材料。
func (k *selfNTKeyPair) SharedSecret(peerPoint []byte) ([]byte, error) {
	if k == nil || k.priv == nil {
		return nil, errors.New("密钥对未初始化")
	}
	return selfNTSharedSecret(k.priv, peerPoint)
}
