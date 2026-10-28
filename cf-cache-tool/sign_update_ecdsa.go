//go:build ignore

package main

// 启动器清单的 ECDSA P-256 签名工具。
// 私钥只保存在构建机 %USERPROFILE%\.cg-keys\launcher_ecdsa.key。
//
// 用法：
//   go run sign_update_ecdsa.go keygen
//   go run sign_update_ecdsa.go pub
//   go run sign_update_ecdsa.go sign <version> <sha256> <size>

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func keyPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cg-keys", "launcher_ecdsa.key")
}

func loadKey() (*ecdsa.PrivateKey, error) {
	raw, err := os.ReadFile(keyPath())
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("私钥 PEM 无效")
	}
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := k.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("不是 ECDSA 私钥")
	}
	return priv, nil
}

func keygen() error {
	p := keyPath()
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("私钥已存在，拒绝覆盖: %s", p)
	}
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(p, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0o600); err != nil {
		return err
	}
	fmt.Printf("X=%s\nY=%s\n", hex32(priv.PublicKey.X), hex32(priv.PublicKey.Y))
	return nil
}

func printPub() error {
	priv, err := loadKey()
	if err != nil {
		return err
	}
	fmt.Printf("X=%s\nY=%s\n", hex32(priv.PublicKey.X), hex32(priv.PublicKey.Y))
	return nil
}

func makeSig(version, sha string, size int64) (string, error) {
	priv, err := loadKey()
	if err != nil {
		return "", err
	}
	payload := version + "\n" + strings.ToLower(strings.TrimSpace(sha)) + "\n" + strconv.FormatInt(size, 10)
	digest := sha256.Sum256([]byte(payload))
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
	if err != nil {
		return "", err
	}
	out := make([]byte, 64)
	r.FillBytes(out[:32])
	s.FillBytes(out[32:])
	return base64.StdEncoding.EncodeToString(out), nil
}

func sign(version, sha string, size int64) error {
	sig, err := makeSig(version, sha, size)
	if err != nil {
		return err
	}
	fmt.Print(sig)
	return nil
}

func addManifest(input, output string) error {
	raw, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}
	var version, sha string
	var size int64
	if err := json.Unmarshal(m["version"], &version); err != nil {
		return fmt.Errorf("清单缺少 version")
	}
	if err := json.Unmarshal(m["sha256"], &sha); err != nil {
		return fmt.Errorf("清单缺少 sha256")
	}
	if err := json.Unmarshal(m["size"], &size); err != nil {
		return fmt.Errorf("清单缺少 size")
	}
	sig, err := makeSig(version, sha, size)
	if err != nil {
		return err
	}
	encSig, _ := json.Marshal(sig)
	m["sig_ecdsa"] = encSig
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(output, data, 0o644)
}

func hex32(v *big.Int) string {
	return fmt.Sprintf("%064x", v)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: sign_update_ecdsa keygen|pub|sign")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "keygen":
		err = keygen()
	case "pub":
		err = printPub()
	case "sign":
		if len(os.Args) != 5 {
			err = fmt.Errorf("usage: sign <version> <sha256> <size>")
			break
		}
		size, e := strconv.ParseInt(os.Args[4], 10, 64)
		if e != nil {
			err = e
			break
		}
		err = sign(os.Args[2], os.Args[3], size)
	case "add":
		if len(os.Args) != 4 {
			err = fmt.Errorf("usage: add <input.json> <output.json>")
			break
		}
		err = addManifest(os.Args[2], os.Args[3])
	default:
		err = fmt.Errorf("unknown command")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}
