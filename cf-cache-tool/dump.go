package main

// 离线数据源：把 dump（cards.jsonl + data/<QQ>.ini）打包放在 dav 的 CG/dump/，
// 切到「离线」时下载到本地缓存，列表和上号用的 ini 都从这里读，不再打接口。

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	dumpRemoteZip  = cloudDir + "/dump/dump.zip"
	dumpRemoteMeta = cloudDir + "/dump/meta.json"
)

// 数据源取值
const (
	srcOnline = "online" // 卡密平台接口（默认）
	srcDav    = "dav"    // dav 上的离线打包数据
)

type dumpMeta struct {
	Size  int64  `json:"size"`
	Files int    `json:"files"`
	Built string `json:"built"`
}

func (m dumpMeta) key() string { return fmt.Sprintf("%d-%d", m.Size, m.Files) }

// dumpCacheDir 本地缓存目录。
func dumpCacheDir() string { return filepath.Join(configDirPath(), "data", "dump") }

func dumpStampPath() string { return dumpCacheDir() + ".stamp" }

func dumpZipPath() string { return filepath.Join(configDirPath(), "data", "dump.zip") }

func hasDumpData(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, "cards.jsonl"))
	return err == nil && st.Size() > 0
}

// davGet 带 Basic Auth 的 GET。
func davGet(url string, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cloudUser, cloudPassword)
	cli := &http.Client{Timeout: timeout}
	return cli.Do(req)
}

func fetchDumpMeta() (dumpMeta, error) {
	var m dumpMeta
	resp, err := davGet(cloudBaseURL+dumpRemoteMeta, 10*time.Second)
	if err != nil {
		return m, fmt.Errorf("连接云端失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return m, fmt.Errorf("meta.json HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		return m, fmt.Errorf("解析 meta.json: %v", err)
	}
	return m, nil
}

// ensureDump 保证本地缓存与云端一致，返回目录。云端不可用时退回本地缓存。
func ensureDump(logf func(string, ...any)) (string, error) {
	dir := dumpCacheDir()

	meta, err := fetchDumpMeta()
	if err != nil {
		if hasDumpData(dir) {
			if logf != nil {
				logf("云端不可用，用本地离线数据: %v", err)
			}
			return dir, nil
		}
		return "", err
	}

	if hasDumpData(dir) {
		if raw, e := os.ReadFile(dumpStampPath()); e == nil && strings.TrimSpace(string(raw)) == meta.key() {
			return dir, nil
		}
	}

	if logf != nil {
		logf("下载离线数据（%d 个 ini / %.1f MB）…", meta.Files, float64(meta.Size)/1024/1024)
	}
	if err := downloadDumpZip(); err != nil {
		if hasDumpData(dir) {
			if logf != nil {
				logf("下载失败，继续用本地缓存: %v", err)
			}
			return dir, nil
		}
		return "", err
	}

	if logf != nil {
		logf("解压离线数据…")
	}
	if err := unzipDump(dir); err != nil {
		return "", err
	}
	_ = os.WriteFile(dumpStampPath(), []byte(meta.key()), 0o600)
	return dir, nil
}

func downloadDumpZip() error {
	if err := os.MkdirAll(filepath.Dir(dumpZipPath()), 0o700); err != nil {
		return err
	}
	resp, err := davGet(cloudBaseURL+dumpRemoteZip, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dump.zip HTTP %d", resp.StatusCode)
	}
	f, err := os.OpenFile(dumpZipPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// unzipDump 解压到 dir，忽略越界路径。
func unzipDump(dir string) error {
	zr, err := zip.OpenReader(dumpZipPath())
	if err != nil {
		return fmt.Errorf("打开 zip: %v", err)
	}
	defer zr.Close()

	_ = os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	root := filepath.Clean(dir) + string(os.PathSeparator)

	for _, zf := range zr.File {
		name := filepath.Clean(zf.Name)
		if strings.HasPrefix(name, "..") || filepath.IsAbs(name) {
			continue
		}
		dst := filepath.Join(dir, name)
		if !strings.HasPrefix(filepath.Clean(dst)+string(os.PathSeparator), root) {
			continue
		}
		if zf.FileInfo().IsDir() {
			_ = os.MkdirAll(dst, 0o700)
			continue
		}
		_ = os.MkdirAll(filepath.Dir(dst), 0o700)
		rc, err := zf.Open()
		if err != nil {
			continue
		}
		out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err != nil {
			rc.Close()
			continue
		}
		_, _ = io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return nil
}

// dumpRows 从本地 cards.jsonl 生成列表，状态/关键字/排序与在线模式一致。
func dumpRows(dir string, f ymmFilter, keyword string) ([]ymmRow, int, error) {
	file, err := os.Open(filepath.Join(dir, "cards.jsonl"))
	if err != nil {
		return nil, 0, fmt.Errorf("读 cards.jsonl: %v", err)
	}
	defer file.Close()

	total := 0
	var out []ymmRow
	sc := bufio.NewScanner(file)
	sc.Buffer(make([]byte, 0, 256<<10), 8<<20)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var obj map[string]any
		if json.Unmarshal(line, &obj) != nil {
			continue
		}
		it := cardItemFromMap(obj)
		total++
		if f.Status != "" && f.Status != "all" && it.Status != f.Status {
			continue
		}
		if f.Since > 0 && it.ActivatedAt < time.Now().Add(-f.Since).UnixMilli() {
			continue
		}
		row := ymmRowFromItem(it)
		if !ymmRowMatch(row, keyword) {
			continue
		}
		out = append(out, row)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
	return out, total, nil
}

// dumpIni 读本地 data/<qq>.ini。
func dumpIni(dir, qq string) (string, bool) {
	if dir == "" || qq == "" {
		return "", false
	}
	raw, err := os.ReadFile(filepath.Join(dir, "data", qq+".ini"))
	if err != nil || len(raw) == 0 {
		return "", false
	}
	return string(raw), true
}
