package main

// 内嵌的上号器分发资源。新增上号器类型时：
//  1. 把 loader.exe / host.dat 放到 assets/<id>/ 下
//  2. 加一行 //go:embed
//  3. 在 loaderProfiles 里注册

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
)

//go:embed assets/yuanmeng/loader.exe
var yuanmengLoader []byte

//go:embed assets/yuanmeng/host.dat
var yuanmengHostDat []byte

type loaderProfile struct {
	ID       string
	Name     string
	ExeName  string
	SiteName string // 卡网名称
	Site     string // 卡网链接
	Loader   []byte
	HostDat  []byte
}

var loaderProfiles = []loaderProfile{
	{
		ID:       "1",
		Name:     "圆梦科级上号器",
		ExeName:  "CF上号器_v1.17.exe",
		SiteName: "圆梦卡网",
		Site:     "https://www.517fk.cn/links/8B1B512A",
		Loader:   yuanmengLoader,
		HostDat:  yuanmengHostDat,
	},
}

func findProfile(id string) *loaderProfile {
	for i := range loaderProfiles {
		if loaderProfiles[i].ID == id {
			return &loaderProfiles[i]
		}
	}
	if len(loaderProfiles) > 0 {
		return &loaderProfiles[0]
	}
	return nil
}

func profileNames() []string {
	names := make([]string, len(loaderProfiles))
	for i, p := range loaderProfiles {
		names[i] = p.Name
	}
	return names
}

func profileByIndex(idx int) *loaderProfile {
	if idx < 0 || idx >= len(loaderProfiles) {
		return findProfile("")
	}
	return &loaderProfiles[idx]
}

// ensure 把内嵌的上号器与 host.dat 解压到工作目录。
// force 为真时覆盖上号器；host.dat 只在缺失或已被本工具改成 127.0.0.1 时才重写。
func (p *loaderProfile) ensure(dir string, force bool) (string, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	exePath := filepath.Join(dir, p.ExeName)
	if _, err := os.Stat(exePath); force || os.IsNotExist(err) {
		if err := os.WriteFile(exePath, p.Loader, 0o700); err != nil {
			return "", err
		}
	}

	hostPath := filepath.Join(dir, "host.dat")
	needHost := false
	if _, err := os.Stat(hostPath); os.IsNotExist(err) {
		needHost = true
	} else if force {
		if cur, err := DecryptHostDat(hostPath); err != nil || strings.Contains(cur, "127.0.0.1") || strings.Contains(cur, "localhost") {
			needHost = true
		}
	}
	if needHost {
		if err := os.WriteFile(hostPath, p.HostDat, 0o600); err != nil {
			return "", err
		}
	}
	return exePath, nil
}
