//go:build ignore

// build.sh 用来把凭据混淆成 base64+XOR 串（配合 secret_str.go 的 deob）。
//
//	go run obfstr.go "明文"
package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

const mask = "cgx0r"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "用法: go run obfstr.go <明文>")
		os.Exit(2)
	}
	b := []byte(os.Args[1])
	for i := range b {
		b[i] ^= mask[i%len(mask)]
	}
	fmt.Print(base64.StdEncoding.EncodeToString(b))
}
