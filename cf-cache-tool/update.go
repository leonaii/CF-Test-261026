package main

// 云端自动更新：
//
//	CG/update/version.json   清单（版本、文件名、sha256、大小）
//	CG/update/CG.exe         新版本二进制
//
// 启动后后台检查，发现新版本就下载 → 校验 → 替换当前 exe，下次启动生效。
// 全程不阻塞界面，任何失败都只写日志。

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	updateDir      = cloudDir + "/update"
	updateManifest = updateDir + "/version.json"
	updateBinary   = updateDir + "/CG.exe"
)

type updateManifestFile struct {
	Version string `json:"version"`
	File    string `json:"file"`
	SHA256  string `json:"sha256"`
	Size    int64  `json:"size"`
	Notes   string `json:"notes"`
}

// fetchUpdateManifest 读取云端清单；不存在返回 errCloudNotFound。
func fetchUpdateManifest(ctx context.Context) (*updateManifestFile, error) {
	data, err := newCloudClient().download(ctx, updateManifest)
	if err != nil {
		return nil, err
	}
	var m updateManifestFile
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("更新清单格式错误: %w", err)
	}
	if m.Version == "" {
		return nil, errors.New("更新清单缺少 version")
	}
	return &m, nil
}

// downloadUpdate 下载新版本到 exe 同目录的临时文件，并校验大小与 sha256。
func downloadUpdate(ctx context.Context, curExe string, m *updateManifestFile) (string, error) {
	path := updateBinary
	if strings.TrimSpace(m.File) != "" {
		path = updateDir + "/" + strings.TrimPrefix(m.File, "/")
	}

	client := &http.Client{Timeout: 10 * time.Minute}
	c := newCloudClient()
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(path), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("User-Agent", "CG/"+toolVersion)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("下载更新失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", errCloudNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载更新失败: HTTP %d", resp.StatusCode)
	}

	tmp := filepath.Join(filepath.Dir(curExe), "CG.update.tmp")
	_ = os.Remove(tmp)
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o700)
	if err != nil {
		return "", fmt.Errorf("写入临时文件失败（可能需要管理员权限）: %w", err)
	}

	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(resp.Body, 200<<20))
	closeErr := f.Close()
	if err != nil {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("下载中断: %w", err)
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return "", closeErr
	}

	if m.Size > 0 && n != m.Size {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("大小不一致（期望 %d，实际 %d）", m.Size, n)
	}
	if want := strings.ToLower(strings.TrimSpace(m.SHA256)); want != "" {
		if got := hex.EncodeToString(h.Sum(nil)); got != want {
			_ = os.Remove(tmp)
			return "", errors.New("sha256 校验失败，已丢弃下载文件")
		}
	}
	return tmp, nil
}

// applyUpdate 用新文件替换当前 exe：先把运行中的 exe 改名，再把新文件就位。
// 失败会自动回滚，保证原来还能启动。
func applyUpdate(curExe, newPath string) error {
	old := curExe + ".old"
	_ = os.Remove(old) // 上一次的残留；删不掉就换个名字
	if _, err := os.Stat(old); err == nil {
		old = fmt.Sprintf("%s.%d", old, time.Now().Unix())
	}

	// 杀软/索引服务可能短暂占用，重试几次
	var err error
	for i := 0; i < 5; i++ {
		if err = os.Rename(curExe, old); err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		return fmt.Errorf("无法替换当前程序（可能没有写权限或被占用）: %w", err)
	}

	for i := 0; i < 5; i++ {
		if err = os.Rename(newPath, curExe); err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	_ = os.Rename(old, curExe) // 回滚
	return fmt.Errorf("替换失败，已回滚: %w", err)
}

// cleanupOldExe 启动时删掉上一次替换留下的 .old（失败忽略）。
func cleanupOldExe() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	_ = os.Remove(exe + ".old")
	if matches, err := filepath.Glob(exe + ".old.*"); err == nil {
		for _, p := range matches {
			_ = os.Remove(p)
		}
	}
	_ = os.Remove(filepath.Join(filepath.Dir(exe), "CG.update.tmp"))
}
