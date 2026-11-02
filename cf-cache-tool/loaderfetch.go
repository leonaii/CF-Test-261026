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

// cleanupLoaderCache 清理旧版本遗留的本地上号器下载缓存。
func cleanupLoaderCache() {
	_ = os.RemoveAll(filepath.Join(configDirPath(), "loaders"))
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

func copyFileLocal(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
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

// ensureLocalLoader 本地直供上号器：目标目录没有就从 CG 同级 / other-tools 拷一份。
func ensureLocalLoader(dir string, p *loaderProfile) (string, error) {
	if p == nil {
		return "", errors.New("未指定上号器类型")
	}
	dst := filepath.Join(dir, p.ExeName)
	if st, err := os.Stat(dst); err == nil && st.Size() > 0 {
		return dst, nil
	}
	base, _ := os.Executable()
	exeDir := filepath.Dir(base)
	names := []string{p.ExeName}
	if p.NoHostDat {
		names = append(names, "鸡哥直登.exe")
	}
	for _, name := range names {
		cands := []string{
			filepath.Join(exeDir, name),
			filepath.Join(exeDir, "other-tools", "kunge", name),
			filepath.Join(exeDir, "..", "other-tools", "kunge", name),
			filepath.Join(exeDir, "..", "..", "other-tools", "kunge", name),
			filepath.Join(exeDir, "资源", "上号器", name),
			filepath.Join(exeDir, "..", "资源", "上号器", name),
		}
		cands = append(cands,
			filepath.Join(exeDir, "资源", "上号器", "niubi", name),
			filepath.Join(exeDir, "..", "资源", "上号器", "niubi", name),
			filepath.Join(exeDir, "..", "..", "资源", "上号器", "niubi", name),
		)
		for _, c := range cands {
			st, err := os.Stat(c)
			if err != nil || st.Size() == 0 {
				continue
			}
			if err := os.MkdirAll(dir, 0o700); err != nil {
				return "", err
			}
			if err := copyFileLocal(c, dst); err != nil {
				return "", err
			}
			return dst, nil
		}
	}
	return "", fmt.Errorf("本地没有 %s，请把它放到 %s", p.ExeName, dir)
}

// localLoaderReady 本地是否已有可用且校验通过的上号器；false 时返回原因便于排查。
func localLoaderReady(dir string, p *loaderProfile) (bool, string) {
	if p == nil {
		return false, "未选择上号器类型"
	}
	if p.LocalOnly || p.NoHostDat {
		st, err := os.Stat(filepath.Join(dir, p.ExeName))
		if err != nil || st.Size() == 0 {
			return false, "本地没有上号器文件（" + p.ExeName + "）"
		}
		return true, ""
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
	if p.LocalOnly {
		exe, err := ensureLocalLoader(dir, p)
		if err != nil {
			return "", err
		}
		logf("上号器已就绪: %s（本地直供）", filepath.Base(exe))
		return exe, nil
	}

	c := newCloudClient()
	raw, err := c.download(ctx, loaderMetaPath(p.ID))
	if err != nil {
		// 云端还没发布时，本地有就直接用（开发者本机/离线包）
		if p.NoHostDat {
			if exe, e2 := ensureLocalLoader(dir, p); e2 == nil {
				logf("云端暂无 %s，使用本地文件: %s", p.ExeName, filepath.Base(exe))
				return exe, nil
			}
		}
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
	}

	// 2) host.dat：缺失、或当前被本工具改成 127.0.0.1、或云端版本变了才写
	if !p.NoHostDat {
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
			data, err := c.download(ctx, loaderDatPath(p.ID))
			if err != nil {
				return "", fmt.Errorf("下载 host.dat 失败: %w", err)
			}
			if meta.DatSHA != "" {
				sum := sha256.Sum256(data)
				if !strings.EqualFold(hex.EncodeToString(sum[:]), meta.DatSHA) {
					return "", errors.New("host.dat 校验失败，已丢弃")
				}
			}
			if err := os.WriteFile(datPath, data, 0o600); err != nil {
				return "", fmt.Errorf("写入 host.dat 失败: %w", err)
			}
			logf("已下载 host.dat")
			// 旧备份作废，避免还原时把老版本写回去
			_ = os.Remove(datPath + ".orig")
		}
	}

	// 3) 换了文件名的旧版本 exe 清掉
	if m := readLoaderMarker(dir); m != nil && m.Meta.Exe != "" && !strings.EqualFold(m.Meta.Exe, meta.Exe) {
		_ = os.Remove(filepath.Join(dir, m.Meta.Exe))
	}

	_ = writeLoaderMarker(dir, &installedLoader{Meta: meta, Installed: time.Now().Format("2006-01-02 15:04:05")})
	logf("上号器已就绪: %s", meta.Exe)
	return exePath, nil
}
