//go:build ignore

// 生成 Go / Rust 用的 Windows 资源对象（rsrc.syso）。
//
//	go run build_resources.go <version>
//
// 这里只负责“编译期资源”：管理员清单 + 兜底图标/版本，保证链接出来的 exe 是合法 PE。
// 真正的“和参考程序一致”由 stamp_exe.go 在链接后完成：图标与版本信息在 UPX 前原样
// 取参考 exe 的资源块，数字签名与修改时间在 UPX 后追加。
//
// 依赖 go-winres：go install github.com/tc-hib/go-winres@latest
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	company     = "湖南南澳网络科技有限公司"
	copyright   = "Copyright (C) 2017-2027 GGzuhao.All Rights Reserved"
	description = "游戏租号客户端"
	product     = "游戏租号客户端"
)

func version4(ver string) string {
	nums := []string{}
	for _, part := range strings.FieldsFunc(ver, func(r rune) bool { return r == '.' || r == '-' }) {
		if _, err := strconv.Atoi(part); err == nil {
			nums = append(nums, part)
		}
	}
	for len(nums) < 4 {
		nums = append(nums, "0")
	}
	return strings.Join(nums[:4], ".")
}

func config(ver, original, internal string) map[string]any {
	v4 := version4(ver)
	manifest := map[string]any{
		"identity":                              map[string]any{"name": "cfcache.app", "version": v4},
		"description":                           description,
		"minimum-os":                            "win7",
		"execution-level":                       "administrator",
		"ui-access":                             false,
		"auto-elevate":                          false,
		"dpi-awareness":                         "per monitor v2",
		"disable-theming":                       false,
		"disable-window-filtering":              false,
		"high-resolution-scrolling-aware":       false,
		"ultra-high-resolution-scrolling-aware": false,
		"long-path-aware":                       false,
		"printer-driver-isolation":              false,
		"gdi-scaling":                           false,
		"segment-heap":                          false,
		"use-common-controls-v6":                true,
	}
	info := map[string]any{
		"CompanyName":      company,
		"FileDescription":  description,
		"FileVersion":      v4,
		"InternalName":     internal,
		"LegalCopyright":   copyright,
		"OriginalFilename": original,
		"ProductName":      product,
		"ProductVersion":   v4,
	}
	return map[string]any{
		"RT_GROUP_ICON": map[string]any{
			"APP": map[string]any{"0000": []string{"icon_256.png", "icon_64.png", "icon_48.png"}},
		},
		"RT_MANIFEST": map[string]any{
			"#1": map[string]any{"0804": manifest},
		},
		"RT_VERSION": map[string]any{
			"#1": map[string]any{
				"0000": map[string]any{
					"fixed": map[string]any{"file_version": v4, "product_version": v4},
					"info":  map[string]any{"0804": info},
				},
			},
		},
	}
}

func findWinres() (string, error) {
	if p := os.Getenv("GO_WINRES"); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	if p, err := exec.LookPath("go-winres"); err == nil {
		return p, nil
	}
	home, _ := os.UserHomeDir()
	candidate := filepath.Join(home, "go", "bin", "go-winres.exe")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}
	return "", fmt.Errorf("找不到 go-winres，请先执行: go install github.com/tc-hib/go-winres@latest")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run build_resources.go <version>")
		os.Exit(2)
	}
	ver := strings.TrimSpace(os.Args[1])
	iconDir := "assets"
	for _, name := range []string{"zuhaoclient.ico", "icon_256.png", "icon_64.png", "icon_48.png"} {
		if _, err := os.Stat(filepath.Join(iconDir, name)); err != nil {
			fmt.Fprintf(os.Stderr, "缺少图标: %s\n", filepath.Join(iconDir, name))
			os.Exit(1)
		}
	}
	winres, err := findWinres()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	jobs := []struct {
		original string
		internal string
		out      string
	}{
		{"zuhaoclient.exe", "zuhaoclient", "rsrc.syso"},
		{"CG租号.exe", "CG租号", filepath.Join("launcher", "rsrc.syso")},
	}

	td, err := os.MkdirTemp("", "cg-res-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建临时目录失败: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(td)

	for _, name := range []string{"icon_256.png", "icon_64.png", "icon_48.png"} {
		if err := copyFile(filepath.Join(iconDir, name), filepath.Join(td, name)); err != nil {
			fmt.Fprintf(os.Stderr, "复制图标失败: %v\n", err)
			os.Exit(1)
		}
	}

	for i, job := range jobs {
		cfg := filepath.Join(td, fmt.Sprintf("winres%d.json", i))
		raw, err := json.MarshalIndent(config(ver, job.original, job.internal), "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "生成配置失败: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(cfg, raw, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "写配置失败: %v\n", err)
			os.Exit(1)
		}
		produced := filepath.Join(td, fmt.Sprintf("rsrc%d.syso", i))
		cmd := exec.Command(winres, "make", "--in", cfg, "--arch", "amd64", "--out", produced, "--no-suffix")
		cmd.Dir = td
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "go-winres 执行失败: %v\n", err)
			os.Exit(1)
		}
		if _, err := os.Stat(produced); err != nil {
			names := []string{}
			entries, _ := os.ReadDir(td)
			for _, e := range entries {
				if strings.HasSuffix(e.Name(), ".syso") {
					names = append(names, e.Name())
				}
			}
			fmt.Fprintf(os.Stderr, "资源生成失败: %s; generated=%v\n", produced, names)
			os.Exit(1)
		}
		if err := copyFile(produced, job.out); err != nil {
			fmt.Fprintf(os.Stderr, "写出 %s 失败: %v\n", job.out, err)
			os.Exit(1)
		}
		st, _ := os.Stat(job.out)
		fmt.Printf("resource: %s -> %s (%d bytes)\n", job.original, job.out, st.Size())
	}
}
