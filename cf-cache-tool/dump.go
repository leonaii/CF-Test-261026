package main

// 离线数据源：dav 上的 CG/dump/dump.zip（cards.jsonl + data/<QQ>.ini）。
// 不整包下载：只读 zip 的中央目录 + cards.jsonl；ini 在「②上号」时按需 Range 读单个条目，
// 读到后缓存在本地，后续不再请求。

import (
	"archive/zip"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	dumpRemoteZip   = cloudDir + "/dump/dump.zip"
	dumpRemoteCards = cloudDir + "/dump/cards.jsonl.gz" // 214KB，一次拿完
)

// 数据源取值
const (
	srcOnline = "online" // 卡密平台接口（默认）
	srcDav    = "dav"    // dav 上的离线打包数据
)

// davClient 复用直连连接（TLS 握手很贵，每次新建连接会让按需读取慢到几十秒）。
var davClient = newCloudHTTPClient(90 * time.Second)

// rangeWindow 是预读窗口：zip 的 local header 和文件数据挨着，一次拉 64KB 通常够一个 ini。
const rangeWindow = 64 << 10

// httpRangeReader 让 archive/zip 按需读远端 zip（HTTP Range + 窗口缓存）。
type httpRangeReader struct {
	url  string
	size int64

	mu     sync.Mutex
	winOff int64
	winBuf []byte
}

func (r *httpRangeReader) fetch(off, end int64) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, r.url, nil)
	if err != nil {
		return nil, err
	}
	newCloudClient().authorize(req)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", off, end-1))

	resp, err := davClient.Do(req)
	if err != nil {
		return nil, cloudNetworkError("读取云端文件", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, end-off))
}

// readChunk 从窗口里拷一段，未命中就按窗口对齐拉新的一块。
func (r *httpRangeReader) readChunk(p []byte, off int64) (int, error) {
	r.mu.Lock()
	if r.winBuf != nil && off >= r.winOff && off < r.winOff+int64(len(r.winBuf)) {
		n := copy(p, r.winBuf[off-r.winOff:])
		r.mu.Unlock()
		return n, nil
	}
	r.mu.Unlock()

	wOff := off / rangeWindow * rangeWindow
	wEnd := wOff + rangeWindow
	if wEnd > r.size {
		wEnd = r.size
	}
	buf, err := r.fetch(wOff, wEnd)
	if err != nil {
		return 0, err
	}
	r.mu.Lock()
	r.winOff, r.winBuf = wOff, buf
	r.mu.Unlock()

	if off < wOff || off >= wOff+int64(len(buf)) {
		return 0, fmt.Errorf("range 越界")
	}
	return copy(p, buf[off-wOff:]), nil
}

func (r *httpRangeReader) ReadAt(p []byte, off int64) (int, error) {
	if off >= r.size {
		return 0, io.EOF
	}
	want := len(p)
	if remain := r.size - off; int64(want) > remain {
		want = int(remain)
	}

	got := 0
	for got < want {
		n, err := r.readChunk(p[got:want], off+int64(got))
		if n > 0 {
			got += n
		}
		if err != nil {
			if got == want {
				break
			}
			return got, err
		}
		if n == 0 {
			break
		}
	}
	if got < len(p) {
		return got, io.EOF
	}
	return got, nil
}

var (
	dumpMu   sync.Mutex
	dumpZip  *zip.Reader
	dumpSize int64
	dumpErr  error
	dumpInit bool
)

func dumpURL() string { return newCloudClient().urlFor(dumpRemoteZip) }

func davHeadSize(url string) (int64, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0, err
	}
	newCloudClient().authorize(req)
	resp, err := davClient.Do(req)
	if err != nil {
		return 0, cloudNetworkError("连接云端", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dump.zip HTTP %d", resp.StatusCode)
	}
	if resp.ContentLength <= 0 {
		return 0, fmt.Errorf("dump.zip 大小未知")
	}
	return resp.ContentLength, nil
}

// dumpSource 进程内只打开一次远端 zip（只读中央目录，几百 KB）。
func dumpSource() (*zip.Reader, int64, error) {
	dumpMu.Lock()
	defer dumpMu.Unlock()
	if dumpInit {
		return dumpZip, dumpSize, dumpErr
	}
	dumpInit = true

	size, err := davHeadSize(dumpURL())
	if err != nil {
		dumpErr = err
		return nil, 0, err
	}
	zr, err := zip.NewReader(&httpRangeReader{url: dumpURL(), size: size}, size)
	if err != nil {
		dumpErr = fmt.Errorf("读 zip 目录失败: %v", err)
		return nil, 0, dumpErr
	}
	dumpZip, dumpSize = zr, size
	return zr, size, nil
}

func dumpCacheRoot() string { return filepath.Join(configDirPath(), "data", "dump") }

func dumpCardsPath() string { return filepath.Join(dumpCacheRoot(), "cards.jsonl") }

func dumpIniPath(qq string) string { return filepath.Join(dumpCacheRoot(), "data", qq+".ini") }

func dumpStampPath() string { return filepath.Join(dumpCacheRoot(), "stamp") }

// findZipEntry 在 zip 里找条目（兼容 / 与 \ 分隔符）。
func findZipEntry(zr *zip.Reader, name string) *zip.File {
	name = strings.ReplaceAll(name, "\\", "/")
	for _, f := range zr.File {
		if strings.ReplaceAll(f.Name, "\\", "/") == name {
			return f
		}
	}
	return nil
}

// davGet 带 Basic Auth 的 GET（复用连接，并按调用方要求限制本次超时）。
func davGet(url string, timeout time.Duration) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	newCloudClient().authorize(req)
	cli := davClient
	if timeout > 0 {
		c := *davClient
		c.Timeout = timeout
		cli = &c
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, cloudNetworkError("连接云端", err)
	}
	return resp, nil
}

// ensureDumpCards 保证本地有 cards.jsonl：只下 214KB 的 gz，不碰整包。
func ensureDumpCards(logf func(string, ...any)) (string, error) {
	cached := func() bool {
		st, e := os.Stat(dumpCardsPath())
		return e == nil && st.Size() > 0
	}
	fallback := func(err error) (string, error) {
		if cached() {
			if logf != nil {
				logf("云端不可用，用本地缓存的离线列表: %v", err)
			}
			return dumpCardsPath(), nil
		}
		return "", err
	}

	url := newCloudClient().urlFor(dumpRemoteCards)
	size, err := davHeadSize(url)
	if err != nil {
		return fallback(err)
	}

	key := fmt.Sprintf("%d", size)
	if cached() {
		if raw, e := os.ReadFile(dumpStampPath()); e == nil && strings.TrimSpace(string(raw)) == key {
			return dumpCardsPath(), nil
		}
	}

	if logf != nil {
		logf("下载离线列表（%.0f KB）…", float64(size)/1024)
	}
	resp, err := davGet(url, 3*time.Minute)
	if err != nil {
		return fallback(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fallback(fmt.Errorf("cards.jsonl.gz HTTP %d", resp.StatusCode))
	}
	gz, err := gzip.NewReader(bufio.NewReader(resp.Body))
	if err != nil {
		return fallback(fmt.Errorf("解压离线列表失败: %v", err))
	}
	defer gz.Close()

	if err := os.MkdirAll(dumpCacheRoot(), 0o700); err != nil {
		return "", err
	}
	tmp := dumpCardsPath() + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(out, gz); err != nil {
		out.Close()
		_ = os.Remove(tmp)
		return "", err
	}
	out.Close()
	if err := os.Rename(tmp, dumpCardsPath()); err != nil {
		return "", err
	}
	_ = os.WriteFile(dumpStampPath(), []byte(key), 0o600)
	return dumpCardsPath(), nil
}

// ensureDumpIni 取某个 QQ 的 ini：先看本地缓存，没有再按需从 zip 读这一个条目。
func ensureDumpIni(qq string) (string, error) {
	if qq == "" {
		return "", fmt.Errorf("QQ 为空")
	}
	if raw, err := os.ReadFile(dumpIniPath(qq)); err == nil && len(raw) > 0 {
		return string(raw), nil
	}

	zr, _, err := dumpSource()
	if err != nil {
		return "", err
	}
	entry := findZipEntry(zr, "data/"+qq+".ini")
	if entry == nil {
		return "", fmt.Errorf("离线库里没有 %s 的 ini", qq)
	}
	rc, err := entry.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	raw, err := io.ReadAll(io.LimitReader(rc, 1<<20))
	if err != nil {
		return "", err
	}
	if len(raw) == 0 {
		return "", fmt.Errorf("离线 ini 为空: %s", qq)
	}
	if err := os.MkdirAll(filepath.Dir(dumpIniPath(qq)), 0o700); err == nil {
		_ = os.WriteFile(dumpIniPath(qq), raw, 0o600)
	}
	return string(raw), nil
}

// dumpRows 从本地 cards.jsonl 生成列表，状态/关键字/排序与在线模式一致。
func dumpRows(dir string, f ymmFilter, keyword string) ([]ymmRow, int, error) {
	file, err := os.Open(dir)
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
