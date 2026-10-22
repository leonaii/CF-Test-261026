package main

// 从云端获取上号器本体与 host.dat，并做 sha256 校验、本地标记。
// 只要本地文件与标记一致，就不再联网。

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type loaderMeta struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Exe     string `json:"exe"`
	ExeSHA  string `json:"exe_sha256"`
	ExeSize int64  `json:"exe_size"`
	DatSHA  string `json:"dat_sha256"`
	Updated string `json:"updated"`
}

// installedLoader 是本地标记：记录当前安装的是哪一版。
type installedLoader struct {
	Meta      loaderMeta `json:"meta"`
	Installed string     `json:"installed"`
}

func loaderMarkerPath(dir string) string { return filepath.Join(dir, ".loader.json") }

// 跨工作目录的下载缓存：%AppData%\CG\loaders\<id>\<sha>，换目录/清目录都不用重新下载
func loaderCacheDir(id string) string {
	if id == "" {
		id = "1"
	}
	return filepath.Join(configDirPath(), "loaders", id)
}

func loaderCacheExe(id, sha string) string {
	return filepath.Join(loaderCacheDir(id), "exe-"+sha)
}

func loaderCacheDat(id, sha string) string {
	return filepath.Join(loaderCacheDir(id), "dat-"+sha)
}

// copyFile 复制文件并返回目标大小。
func copyFile(src, dst string) (int64, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return 0, err
	}
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o700)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(out, in)
	if cerr := out.Close(); err == nil {
		err = cerr
	}
	return n, err
}

func readLoaderMarker(dir string) *installedLoader {
	data, err := os.ReadFile(loaderMarkerPath(dir))
	if err != nil {
		return nil
	}
	var m installedLoader
	if json.Unmarshal(data, &m) != nil {
		return nil
	}
	return &m
}

func writeLoaderMarker(dir string, m *installedLoader) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(loaderMarkerPath(dir), data, 0o600)
}

func fileSHA256(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// localLoaderReady 本地是否已有可用且校验通过的上号器；false 时返回原因便于排查。
func localLoaderReady(dir string, p *loaderProfile) (bool, string) {
	if p == nil {
		return false, "未选择上号器类型"
	}
	m := readLoaderMarker(dir)
	if m == nil || m.Meta.Exe == "" {
		return false, "缺少本地标记（.loader.json）"
	}
	exe := filepath.Join(dir, m.Meta.Exe)
	st, err := os.Stat(exe)
	if err != nil || st.Size() == 0 {
		return false, "本地没有上号器文件"
	}
	if m.Meta.ExeSize > 0 && st.Size() != m.Meta.ExeSize {
		return false, fmt.Sprintf("大小不一致（本地 %d / 云端 %d）", st.Size(), m.Meta.ExeSize)
	}
	sha, err := fileSHA256(exe)
	if err != nil {
		return false, fmt.Sprintf("读取失败: %v", err)
	}
	if !strings.EqualFold(sha, m.Meta.ExeSHA) {
		return false, "sha256 不一致"
	}
	return true, ""
}

// syncLoader 拉取 meta，按需下载 exe / host.dat，校验后写标记。
// 返回最终的上号器路径。
func syncLoader(ctx context.Context, dir string, p *loaderProfile, logf func(string, ...any)) (string, error) {
	if p == nil {
		return "", errors.New("未指定上号器类型")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if logf == nil {
		logf = func(string, ...any) {}
	}

	c := newCloudClient()
	raw, err := c.download(ctx, loaderMetaPath(p.ID))
	if err != nil {
		if errors.Is(err, errCloudNotFound) {
			return "", errors.New("云端还没有这个上号器的文件")
		}
		return "", fmt.Errorf("获取上号器信息失败: %w", err)
	}
	var meta loaderMeta
	if err := json.Unmarshal(raw, &meta); err != nil {
		return "", fmt.Errorf("上号器信息格式错误: %w", err)
	}
	if meta.Exe == "" {
		meta.Exe = p.ExeName
	}

	exePath := filepath.Join(dir, meta.Exe)
	datPath := filepath.Join(dir, "host.dat")

	// 1) 上号器本体
	needExe := true
	if st, err := os.Stat(exePath); err == nil && st.Size() > 0 {
		if meta.ExeSize > 0 && st.Size() == meta.ExeSize {
			if sha, err := fileSHA256(exePath); err == nil && strings.EqualFold(sha, meta.ExeSHA) {
				needExe = false
			}
		}
	}
	if needExe {
		// 先看本机缓存（换工作目录/清目录时不用重新下载）
		cached := loaderCacheExe(p.ID, meta.ExeSHA)
		if _, err := os.Stat(cached); err == nil {
			if sha, e := fileSHA256(cached); e == nil && strings.EqualFold(sha, meta.ExeSHA) {
				if _, e := copyFile(cached, exePath); e == nil {
					logf("上号器取自本地缓存: %s", meta.Exe)
					needExe = false
				}
			}
		}
	}
	if needExe {
		logf("下载上号器: %s (%.1f MB)", meta.Exe, float64(meta.ExeSize)/1024/1024)
		data, err := c.download(ctx, loaderExePath(p.ID, meta.Exe))
		if err != nil {
			return "", fmt.Errorf("下载上号器失败: %w", err)
		}
		if meta.ExeSHA != "" {
			sum := sha256.Sum256(data)
			if !strings.EqualFold(hex.EncodeToString(sum[:]), meta.ExeSHA) {
				return "", errors.New("上号器校验失败，已丢弃")
			}
		}
		if err := os.WriteFile(exePath, data, 0o700); err != nil {
			return "", fmt.Errorf("写入上号器失败（可能没有权限）: %w", err)
		}
		// 存一份到本机缓存，下次换目录不用再下
		if meta.ExeSHA != "" {
			_ = os.MkdirAll(loaderCacheDir(p.ID), 0o700)
			_ = os.WriteFile(loaderCacheExe(p.ID, meta.ExeSHA), data, 0o700)
		}
	}

	// 2) host.dat：缺失、或当前被本工具改成 127.0.0.1、或云端版本变了才写
	needDat := false
	if _, err := os.Stat(datPath); os.IsNotExist(err) {
		needDat = true
	} else if cur, err := DecryptHostDat(datPath); err != nil {
		needDat = true
	} else if strings.Contains(cur, "127.0.0.1") || strings.Contains(cur, "localhost") {
		needDat = true
	} else if m := readLoaderMarker(dir); m != nil && meta.DatSHA != "" && !strings.EqualFold(m.Meta.DatSHA, meta.DatSHA) {
		needDat = true
	}
	if needDat {
		data, err := os.ReadFile(loaderCacheDat(p.ID, meta.DatSHA))
		if err != nil || meta.DatSHA == "" {
			data, err = c.download(ctx, loaderDatPath(p.ID))
			if err != nil {
				return "", fmt.Errorf("下载 host.dat 失败: %w", err)
			}
			if meta.DatSHA != "" {
				sum := sha256.Sum256(data)
				if !strings.EqualFold(hex.EncodeToString(sum[:]), meta.DatSHA) {
					return "", errors.New("host.dat 校验失败，已丢弃")
				}
				_ = os.MkdirAll(loaderCacheDir(p.ID), 0o700)
				_ = os.WriteFile(loaderCacheDat(p.ID, meta.DatSHA), data, 0o600)
			}
			logf("已下载 host.dat")
		} else {
			logf("host.dat 取自本地缓存")
		}
		if err := os.WriteFile(datPath, data, 0o600); err != nil {
			return "", fmt.Errorf("写入 host.dat 失败: %w", err)
		}
		// 旧备份作废，避免还原时把老版本写回去
		_ = os.Remove(datPath + ".orig")
	}

	// 3) 换了文件名的旧版本 exe 清掉
	if m := readLoaderMarker(dir); m != nil && m.Meta.Exe != "" && !strings.EqualFold(m.Meta.Exe, meta.Exe) {
		_ = os.Remove(filepath.Join(dir, m.Meta.Exe))
	}

	_ = writeLoaderMarker(dir, &installedLoader{Meta: meta, Installed: time.Now().Format("2006-01-02 15:04:05")})
	logf("上号器已就绪: %s", meta.Exe)
	return exePath, nil
}
