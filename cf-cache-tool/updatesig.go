package main

// 更新包签名：version.json 带 ed25519 签名，客户端用内置公钥验签。
// 即使 dav 被写，没有私钥也推不出能被接受的假更新。
// 私钥只在构建机（%USERPROFILE%\.cg-keys\update.key），不进仓库、不进 exe。

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const updatePubKeyHex = "b1a9bd1f59c40a4671fa79f0c1dc16833403482045b580f3ebb53fe7249e40b5"

// updateSigPayload 是签名覆盖的内容：版本 + 文件 sha256 + 大小。
func updateSigPayload(version, sha string, size int64) []byte {
	return []byte(version + "\n" + strings.ToLower(strings.TrimSpace(sha)) + "\n" + strconv.FormatInt(size, 10))
}

// verifyUpdateSig 校验清单签名，失败说明清单被改过或不是我们发的。
func verifyUpdateSig(version, sha string, size int64, sigB64 string) error {
	pub, err := hex.DecodeString(updatePubKeyHex)
	if err != nil || len(pub) != ed25519.PublicKeySize {
		return fmt.Errorf("内置公钥无效")
	}
	if strings.TrimSpace(sigB64) == "" {
		return fmt.Errorf("更新清单缺少签名")
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(sigB64))
	if err != nil || len(sig) != ed25519.SignatureSize {
		return fmt.Errorf("签名格式错误")
	}
	if !ed25519.Verify(ed25519.PublicKey(pub), updateSigPayload(version, sha, size), sig) {
		return fmt.Errorf("更新包签名校验失败（可能被篡改）")
	}
	return nil
}
