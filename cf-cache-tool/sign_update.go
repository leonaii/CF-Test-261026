//go:build ignore

// 给 version.json 生成签名（build.sh 调用）。
//
//	go run sign_update.go <version> <sha256> <size>
//
// 私钥：%USERPROFILE%\.cg-keys\update.key（hex 的 ed25519 seed）
package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "用法: go run sign_update.go <version> <sha256> <size>")
		os.Exit(2)
	}
	ver := os.Args[1]
	sha := strings.ToLower(strings.TrimSpace(os.Args[2]))
	size, err := strconv.ParseInt(strings.TrimSpace(os.Args[3]), 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "size 不是数字:", os.Args[3])
		os.Exit(2)
	}

	keyPath := filepath.Join(os.Getenv("USERPROFILE"), ".cg-keys", "update.key")
	raw, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "读私钥失败（先跑一次 keygen）:", err)
		os.Exit(1)
	}
	seed, err := hex.DecodeString(strings.TrimSpace(string(raw)))
	if err != nil || len(seed) != ed25519.SeedSize {
		fmt.Fprintln(os.Stderr, "私钥格式错误")
		os.Exit(1)
	}

	priv := ed25519.NewKeyFromSeed(seed)
	msg := []byte(ver + "\n" + sha + "\n" + strconv.FormatInt(size, 10))
	fmt.Print(base64.StdEncoding.EncodeToString(ed25519.Sign(priv, msg)))
}
