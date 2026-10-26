package main

// 云缓存：用 Nextcloud WebDAV 备份 / 恢复 cache 目录（打包成 zip）。
// Nextcloud 的文件入口固定为 /remote.php/dav/files/<用户名>/。

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	cloudBaseURL = "https://cloud.51i.cc/remote.php/dav/files/leonai/"
	cloudDir     = "CG"
	cloudFile    = "cache.zip"
	cloudRemote  = cloudDir + "/" + cloudFile
)

// dav 凭据不打进源码：正式包由 build.sh 生成 zz_build_secret.go（-tags cg_secret）写死，
// 这样 garble -literals 能把它加密；本地调试用环境变量。
var (
	cloudUser     = buildCloudUser()
	cloudPassword = buildCloudPass()
)

func init() {
	if cloudUser == "" {
		cloudUser = os.Getenv("CG_DAV_USER")
	}
	if cloudPassword == "" {
		cloudPassword = os.Getenv("CG_DAV_PASS")
	}
	warmRuntimeConf()
}

var errCloudNotFound = errors.New("云端还没有备份")

type cloudClient struct {
	base string
	user string
	pass string
	http *http.Client
}

func newCloudClient() *cloudClient {
	return &cloudClient{
		base: cloudBaseURL,
		user: cloudUser,
		pass: cloudPassword,
		http: &http.Client{Timeout: 90 * time.Second},
	}
}

func (c *cloudClient) urlFor(path string) string {
	p := strings.Trim(path, "/")
	if p == "" {
		return c.base
	}
	segs := strings.Split(p, "/")
	for i := range segs {
		segs[i] = url.PathEscape(segs[i])
	}
	u := c.base + strings.Join(segs, "/")
	if strings.HasSuffix(path, "/") {
		u += "/"
	}
	return u
}

func (c *cloudClient) do(ctx context.Context, method, path string, body io.Reader, size int64) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.urlFor(path), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("User-Agent", "CG/"+toolVersion)
	if size >= 0 {
		req.ContentLength = size
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接 Nextcloud 失败: %w", err)
	}
	return resp, nil
}

func drain(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<16))
	_ = resp.Body.Close()
}

// ensureDir 创建远端目录，已存在时 WebDAV 返回 405。
func (c *cloudClient) ensureDir(ctx context.Context, dir string) error {
	resp, err := c.do(ctx, "MKCOL", dir+"/", nil, 0)
	if err != nil {
		return err
	}
	defer drain(resp)
	switch resp.StatusCode {
	case 200, 201, 204, 301, 302, 405:
		return nil
	case 401:
		return errors.New("Nextcloud 账号或应用密码错误（401）")
	}
	return fmt.Errorf("创建云端目录失败: HTTP %d", resp.StatusCode)
}

func (c *cloudClient) upload(ctx context.Context, path string, data []byte) error {
	resp, err := c.do(ctx, "PUT", path, bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	defer drain(resp)
	switch {
	case resp.StatusCode/100 == 2:
		return nil
	case resp.StatusCode == 401:
		return errors.New("Nextcloud 账号或应用密码错误（401）")
	case resp.StatusCode == 507:
		return errors.New("云端空间不足（507）")
	}
	return fmt.Errorf("上传失败: HTTP %d", resp.StatusCode)
}

func (c *cloudClient) download(ctx context.Context, path string) ([]byte, error) {
	resp, err := c.do(ctx, "GET", path, nil, -1)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
		if err != nil {
			return nil, fmt.Errorf("读取云端文件失败: %w", err)
		}
		return data, nil
	case http.StatusNotFound:
		return nil, errCloudNotFound
	case http.StatusUnauthorized:
		return nil, errors.New("Nextcloud 账号或应用密码错误（401）")
	}
	return nil, fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
}

// stat 返回云端备份大小与时间。
// HEAD 在部分 WebDAV 服务端不返回真实大小，统一改用 PROPFIND Depth:0 取。
func (c *cloudClient) stat(ctx context.Context, path string) (int64, time.Time, error) {
	req, err := http.NewRequestWithContext(ctx, "PROPFIND", c.urlFor(path), nil)
	if err != nil {
		return 0, time.Time{}, err
	}
	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("Depth", "0")
	req.Header.Set("User-Agent", "CG/"+toolVersion)

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("连接 Nextcloud 失败: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return 0, time.Time{}, errCloudNotFound
	case http.StatusUnauthorized:
		return 0, time.Time{}, errors.New("Nextcloud 账号或应用密码错误（401）")
	case http.StatusMultiStatus, http.StatusOK:
	default:
		return 0, time.Time{}, fmt.Errorf("查询云端文件失败: HTTP %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, time.Time{}, err
	}
	var ms struct {
		Responses []struct {
			Propstat []struct {
				Prop struct {
					Len string `xml:"getcontentlength"`
					Mod string `xml:"getlastmodified"`
				} `xml:"prop"`
			} `xml:"propstat"`
		} `xml:"response"`
	}
	if err := xml.Unmarshal(raw, &ms); err != nil {
		return 0, time.Time{}, fmt.Errorf("解析云端信息失败: %w", err)
	}
	for _, r := range ms.Responses {
		for _, ps := range r.Propstat {
			if strings.TrimSpace(ps.Prop.Len) == "" {
				continue
			}
			size, _ := strconv.ParseInt(strings.TrimSpace(ps.Prop.Len), 10, 64)
			mod, _ := http.ParseTime(strings.TrimSpace(ps.Prop.Mod))
			return size, mod, nil
		}
	}
	return 0, time.Time{}, errCloudNotFound
}

// normalizeBackupName 校验用户输入的备份名，返回去掉 .zip 的纯名字。
func normalizeBackupName(s string) (string, error) {
	name := strings.TrimSpace(s)
	name = strings.TrimSuffix(name, ".zip")
	name = strings.TrimSuffix(name, ".ZIP")
	if name == "" {
		return "", errors.New("名称不能为空")
	}
	if len([]rune(name)) > 40 {
		return "", errors.New("名称太长（最多 40 个字）")
	}
	if name == "." || name == ".." {
		return "", errors.New("名称不合法")
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return "", errors.New(`名称不能包含 / \ : * ? " < > | 这些字符`)
	}
	return name, nil
}

// packCache 把 cache 目录下的 .body / .meta.json 打成一个 zip。
func packCache(cacheDir string) ([]byte, int, error) {
	info, err := os.Stat(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, errors.New("本地还没有缓存（cache 目录不存在），请先在“录制模式”上号一次")
		}
		return nil, 0, fmt.Errorf("读取缓存目录失败: %w", err)
	}
	if !info.IsDir() {
		return nil, 0, fmt.Errorf("缓存路径不是目录: %s", cacheDir)
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	count := 0

	walkErr := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil // 扫描期间被删掉的文件直接跳过
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if !strings.HasSuffix(name, ".body") && !strings.HasSuffix(name, ".meta.json") {
			return nil
		}
		rel, err := filepath.Rel(cacheDir, path)
		if err != nil {
			return err
		}
		w, err := zw.Create(filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(w, f); err != nil {
			return err
		}
		count++
		return nil
	})
	if walkErr != nil {
		return nil, 0, fmt.Errorf("打包缓存失败: %w", walkErr)
	}
	if err := zw.Close(); err != nil {
		return nil, 0, err
	}
	return buf.Bytes(), count, nil
}

// unpackCache 把 zip 解到 cache 目录，同名文件覆盖。
func unpackCache(data []byte, cacheDir string) (int, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return 0, fmt.Errorf("云端备份不是有效的 zip: %w", err)
	}

	// 备份里出现过的卡密目录先整体清掉，再解压：同卡密覆盖，云端没有的本地卡密保留（增量）
	cards := map[string]bool{}
	for _, f := range zr.File {
		name := filepath.ToSlash(f.Name)
		if i := strings.IndexByte(name, '/'); i > 0 {
			card := name[:i]
			if card != "" && !strings.Contains(card, "..") && !strings.Contains(card, `\`) {
				cards[card] = true
			}
		}
	}
	for card := range cards {
		_ = os.RemoveAll(filepath.Join(cacheDir, card))
	}

	count := 0
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := filepath.ToSlash(f.Name)
		if strings.Contains(name, "..") || strings.HasPrefix(name, "/") {
			continue // 防目录穿越
		}
		if !strings.HasSuffix(name, ".body") && !strings.HasSuffix(name, ".meta.json") {
			continue
		}
		dst := filepath.Join(cacheDir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
			return count, err
		}
		rc, err := f.Open()
		if err != nil {
			return count, err
		}
		out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
		if err != nil {
			rc.Close()
			return count, err
		}
		_, copyErr := io.Copy(out, io.LimitReader(rc, 8<<20))
		out.Close()
		rc.Close()
		if copyErr != nil {
			return count, copyErr
		}
		count++
	}
	return count, nil
}

// cloudBackup 上传本地缓存，返回文件数与字节数。
func cloudBackup(ctx context.Context, cacheDir, remotePath string) (int, int64, error) {
	if remotePath == "" {
		remotePath = cloudRemote
	}
	data, files, err := packCache(cacheDir)
	if err != nil {
		return 0, 0, err
	}
	if files == 0 {
		return 0, 0, errors.New("本地没有可备份的缓存")
	}
	c := newCloudClient()
	if err := c.ensureDir(ctx, filepath.ToSlash(filepath.Dir(remotePath))); err != nil {
		return 0, 0, err
	}
	if err := c.upload(ctx, remotePath, data); err != nil {
		return 0, 0, err
	}
	return files, int64(len(data)), nil
}

// cloudRestore 下载云端备份并合并到本地缓存目录。
func cloudRestore(ctx context.Context, cacheDir, remotePath string) (int, int64, error) {
	if remotePath == "" {
		remotePath = cloudRemote
	}
	c := newCloudClient()
	data, err := c.download(ctx, remotePath)
	if err != nil {
		return 0, 0, err
	}
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		return 0, 0, err
	}
	n, err := unpackCache(data, cacheDir)
	return n, int64(len(data)), err
}

// cloudInfo 查询云端备份信息。
func cloudInfo(ctx context.Context, remotePath string) (int64, time.Time, error) {
	if remotePath == "" {
		remotePath = cloudRemote
	}
	return newCloudClient().stat(ctx, remotePath)
}

// listEntries 列出目录下的 zip 文件名与子目录名。
func (c *cloudClient) listEntries(ctx context.Context, dir string) (files, dirs []string, err error) {
	req, err := http.NewRequestWithContext(ctx, "PROPFIND", c.urlFor(dir)+"/", nil)
	if err != nil {
		return nil, nil, err
	}
	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("Depth", "1")
	req.Header.Set("User-Agent", "CG/"+toolVersion)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("连接 Nextcloud 失败: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusMultiStatus, http.StatusOK:
	case http.StatusNotFound:
		return nil, nil, nil
	case http.StatusUnauthorized:
		return nil, nil, errors.New("Nextcloud 账号或应用密码错误（401）")
	default:
		return nil, nil, fmt.Errorf("列出云端备份失败: HTTP %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, nil, err
	}
	var ms struct {
		Responses []struct {
			Href string `xml:"href"`
		} `xml:"response"`
	}
	if err := xml.Unmarshal(raw, &ms); err != nil {
		return nil, nil, fmt.Errorf("解析云端列表失败: %w", err)
	}

	self := strings.Trim(strings.TrimRight(dir, "/"), "/")
	seenFile := map[string]bool{}
	seenDir := map[string]bool{}
	for _, r := range ms.Responses {
		href := r.Href
		isDir := strings.HasSuffix(href, "/")
		base := strings.TrimSuffix(href, "/")
		if i := strings.LastIndex(base, "/"); i >= 0 {
			base = base[i+1:]
		}
		if unescaped, e := url.PathUnescape(base); e == nil {
			base = unescaped
		}
		if base == "" {
			continue
		}
		if isDir {
			full := strings.Trim(strings.TrimPrefix(strings.TrimSuffix(href, "/"), "/remote.php/dav/files/"), "/")
			if strings.HasSuffix(full, self) || full == self {
				continue
			}
			if !seenDir[base] {
				seenDir[base] = true
				dirs = append(dirs, base)
			}
			continue
		}
		if !strings.HasSuffix(strings.ToLower(base), ".zip") {
			continue
		}
		name := base[:len(base)-4]
		if !seenFile[name] {
			seenFile[name] = true
			files = append(files, name)
		}
	}
	sort.Strings(files)
	sort.Strings(dirs)
	return files, dirs, nil
}

// cloudList 列出某个名称目录下的备份（不含扩展名）。
func cloudList(ctx context.Context, folder string) ([]string, error) {
	folder = strings.Trim(strings.TrimSpace(folder), "/")
	if folder == "" {
		folder = cloudDir
	} else if !strings.HasPrefix(folder, cloudDir+"/") {
		folder = cloudDir + "/" + folder
	}
	files, _, err := newCloudClient().listEntries(ctx, folder)
	return files, err
}

// cloudFolders 列出根目录（CG/）下的账号名称目录。
func cloudFolders(ctx context.Context) ([]string, error) {
	_, dirs, err := newCloudClient().listEntries(ctx, cloudDir)
	return dirs, err
}

// cloudPathFor 云端备份路径：CG/<云端名称>/<上号器类型>.zip，例如 CG/leonai/1.zip。
func cloudPathFor(folder, profileID string) string {
	folder = strings.Trim(strings.TrimSpace(folder), "/")
	if folder == "" {
		folder = "default"
	}
	if profileID == "" {
		profileID = "1"
	}
	return cloudDir + "/" + folder + "/" + profileID + ".zip"
}
