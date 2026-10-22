package main

// 卡密接口响应信封解析。
//
// 响应形如：
//   {"success":true,"key":"<base64 RSA(PKCS1v15) 加密的 AES 密钥>",
//    "iv":"<base64 12 字节>","tag":"<base64 16 字节>",
//    "data":"<base64 密文>","verify":"..."}
//
// 程序用内嵌 RSA 私钥解出 AES 密钥，再做 AES-GCM 解密。
// 该私钥可从原程序 RVA 0x1400FB140 的 BCRYPT_RSAFULLPRIVATE_BLOB 提取。

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type envelope struct {
	Success *bool  `json:"success"`
	Key     string `json:"key"`
	IV      string `json:"iv"`
	Tag     string `json:"tag"`
	Data    string `json:"data"`
	Verify  string `json:"verify"`
}

func decodeB64(s string) ([]byte, error) {
	encs := []*base64.Encoding{
		base64.StdEncoding, base64.RawStdEncoding,
		base64.URLEncoding, base64.RawURLEncoding,
	}
	var lastErr error = errors.New("空字符串")
	for _, enc := range encs {
		if b, err := enc.DecodeString(s); err == nil {
			return b, nil
		} else {
			lastErr = err
		}
	}
	return nil, lastErr
}

var (
	privOnce   sync.Once
	cachedPriv *rsa.PrivateKey
)

func privateKey() *rsa.PrivateKey {
	privOnce.Do(func() { cachedPriv = loadPrivateKey() })
	return cachedPriv
}

// decryptSessionKey 依次尝试 OAEP-SHA256（线上协议实际使用）、OAEP-SHA1、PKCS#1 v1.5（兼容旧缓存/测试样本）。
func decryptSessionKey(ct []byte) ([]byte, error) {
	priv := privateKey()
	if k, err := rsa.DecryptOAEP(sha256.New(), nil, priv, ct, nil); err == nil && len(k) > 0 {
		return k, nil
	}
	if k, err := rsa.DecryptOAEP(sha1.New(), nil, priv, ct, nil); err == nil && len(k) > 0 {
		return k, nil
	}
	if k, err := rsa.DecryptPKCS1v15(nil, priv, ct); err == nil && len(k) > 0 {
		return k, nil
	}
	return nil, errors.New("RSA 解密会话密钥失败（OAEP-SHA256 / OAEP-SHA1 / PKCS1 都不匹配）")
}

// OpenEnvelope 解出响应信封中的明文 JSON（账号数据）。
func OpenEnvelope(body []byte) ([]byte, error) {
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("响应不是合法 JSON: %w", err)
	}
	if env.Key == "" || env.Data == "" || env.IV == "" || env.Tag == "" {
		return nil, errors.New("响应缺少 key/iv/tag/data 字段")
	}
	keyCt, err := decodeB64(env.Key)
	if err != nil {
		return nil, fmt.Errorf("key base64 解码失败: %w", err)
	}
	aesKey, err := decryptSessionKey(keyCt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("AES 密钥长度非法(%d): %w", len(aesKey), err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	iv, err := decodeB64(env.IV)
	if err != nil {
		return nil, fmt.Errorf("iv base64 解码失败: %w", err)
	}
	tag, err := decodeB64(env.Tag)
	if err != nil {
		return nil, fmt.Errorf("tag base64 解码失败: %w", err)
	}
	data, err := decodeB64(env.Data)
	if err != nil {
		return nil, fmt.Errorf("data base64 解码失败: %w", err)
	}
	if len(iv) != gcm.NonceSize() {
		return nil, fmt.Errorf("iv 长度 %d，期望 %d", len(iv), gcm.NonceSize())
	}
	sealed := append(append([]byte{}, data...), tag...)
	plain, err := gcm.Open(nil, iv, sealed, nil)
	if err != nil {
		return nil, fmt.Errorf("AES-GCM 解密失败: %w", err)
	}
	return plain, nil
}

// PrettyJSON 格式化 JSON，失败则原样返回。
func PrettyJSON(b []byte) string {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return string(b)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(out)
}
