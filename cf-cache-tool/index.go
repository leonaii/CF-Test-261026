package main

// 云端按需结构：
//
//	CG/<卡密>/<类型ID>/
//	    index.json                 账号列表（列表页只下这个，约 2KB）
//	    cards/<卡密>.zip           单个账号的完整缓存（回放时才下）
//
// 录制完成后只上传这一个账号的 zip，并合并更新 index.json。

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type cloudCardEntry struct {
	Card     string `json:"card"`
	QQ       string `json:"qq,omitempty"`
	OpenID   string `json:"openid,omitempty"`
	Zone     string `json:"zone,omitempty"`
	Note     string `json:"note,omitempty"`
	RawNote  string `json:"raw_note,omitempty"`
	Recorded string `json:"recorded,omitempty"`
	Mark     string `json:"mark,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Updated  string `json:"updated,omitempty"`
}

type cloudIndex struct {
	Version int              `json:"version"`
	Updated string           `json:"updated,omitempty"`
	Cards   []cloudCardEntry `json:"cards"`
}

func cloudBaseDir(name, pid string) string {
	name = strings.Trim(strings.TrimSpace(name), "/")
	if name == "" {
		name = "default"
	}
	if pid == "" {
		pid = "1"
	}
	return cloudDir + "/" + name + "/" + pid
}

func cloudIndexPath(name, pid string) string {
	return cloudBaseDir(name, pid) + "/index.json"
}

func cloudCardPath(name, pid, card string) string {
	return cloudBaseDir(name, pid) + "/cards/" + sanitizeName(card) + ".zip"
}

func (idx *cloudIndex) find(card string) *cloudCardEntry {
	for i := range idx.Cards {
		if idx.Cards[i].Card == card {
			return &idx.Cards[i]
		}
	}
	return nil
}

func (idx *cloudIndex) upsert(e cloudCardEntry) {
	e.Updated = time.Now().Format("2006-01-02 15:04:05")
	if old := idx.find(e.Card); old != nil {
		// 保留旧的标记（除非本次明确给了新的）
		if e.Mark == "" {
			e.Mark = old.Mark
		}
		*old = e
		return
	}
	idx.Cards = append(idx.Cards, e)
}

func (idx *cloudIndex) setMark(card, mark string) {
	if e := idx.find(card); e != nil {
		e.Mark = mark
		e.Updated = time.Now().Format("2006-01-02 15:04:05")
	}
}

// fetchCloudIndex 拉取云端索引；不存在时返回空索引（不是错误）。
func fetchCloudIndex(ctx context.Context, name, pid string) (*cloudIndex, error) {
	data, err := newCloudClient().download(ctx, cloudIndexPath(name, pid))
	if err != nil {
		if errors.Is(err, errCloudNotFound) {
			return &cloudIndex{Version: 1}, nil
		}
		return nil, err
	}
	var idx cloudIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("云端 index.json 格式错误: %w", err)
	}
	if idx.Version == 0 {
		idx.Version = 1
	}
	return &idx, nil
}

func putCloudIndex(ctx context.Context, name, pid string, idx *cloudIndex) error {
	idx.Version = 1
	idx.Updated = time.Now().Format("2006-01-02 15:04:05")
	sort.Slice(idx.Cards, func(i, j int) bool { return idx.Cards[i].Card < idx.Cards[j].Card })

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	c := newCloudClient()
	if err := c.ensureDir(ctx, cloudBaseDir(name, pid)+"/"); err != nil {
		return err
	}
	if err := c.ensureDir(ctx, cloudBaseDir(name, pid)+"/cards/"); err != nil {
		return err
	}
	return c.upload(ctx, cloudIndexPath(name, pid), data)
}

// packCardDir 把 cache/<card>/ 下的缓存文件打成 zip（内部仍用 card/ 前缀）。
func packCardDir(cacheDir, card string) ([]byte, int, error) {
	dir := filepath.Join(cacheDir, sanitizeName(card))
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		return nil, 0, errors.New("本地没有该卡密的缓存")
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	count := 0
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
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
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		w, err := zw.Create(card + "/" + filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
		count++
		return nil
	})
	if walkErr != nil {
		_ = zw.Close()
		return nil, 0, walkErr
	}
	if err := zw.Close(); err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, errors.New("本地没有该卡密的缓存")
	}
	return buf.Bytes(), count, nil
}

// unpackCardZip 解压单个账号的 zip 到 cache/<card>/，覆盖旧内容并校验完整性。
func unpackCardZip(cacheDir, card string, data []byte) (int, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return 0, fmt.Errorf("云端数据不是有效 zip: %w", err)
	}
	dst := filepath.Join(cacheDir, sanitizeName(card))
	_ = os.RemoveAll(dst)
	if err := os.MkdirAll(dst, 0o700); err != nil {
		return 0, err
	}

	count := 0
	for _, f := range zr.File {
		name := filepath.ToSlash(f.Name)
		base := name
		if i := strings.IndexByte(name, '/'); i >= 0 {
			base = name[i+1:]
		}
		if base == "" || strings.Contains(base, "..") || strings.HasSuffix(base, "/") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return count, err
		}
		content, err := io.ReadAll(io.LimitReader(rc, 8<<20))
		_ = rc.Close()
		if err != nil {
			return count, err
		}
		if err := os.WriteFile(filepath.Join(dst, base), content, 0o600); err != nil {
			return count, err
		}
		count++
	}

	if _, err := os.Stat(filepath.Join(dst, "api_card_redeem.body")); err != nil {
		_ = os.RemoveAll(dst)
		return 0, errors.New("云端缓存缺少 redeem 数据，可能不完整")
	}
	return count, nil
}

// uploadCardCloud 上传单个账号缓存，并把它写进索引。
func uploadCardCloud(ctx context.Context, name, pid, cacheDir, card string, entry cloudCardEntry) error {
	data, _, err := packCardDir(cacheDir, card)
	if err != nil {
		return err
	}
	c := newCloudClient()
	dir := cloudBaseDir(name, pid)
	if err := c.ensureDir(ctx, dir+"/"); err != nil {
		return err
	}
	if err := c.ensureDir(ctx, dir+"/cards/"); err != nil {
		return err
	}
	if err := c.upload(ctx, cloudCardPath(name, pid, card), data); err != nil {
		return err
	}

	idx, err := fetchCloudIndex(ctx, name, pid)
	if err != nil {
		return err
	}
	entry.Card = card
	entry.Size = int64(len(data))
	idx.upsert(entry)
	return putCloudIndex(ctx, name, pid, idx)
}

// downloadCardCloud 按需下载单个账号缓存。
func downloadCardCloud(ctx context.Context, name, pid, cacheDir, card string) (int, int64, error) {
	data, err := newCloudClient().download(ctx, cloudCardPath(name, pid, card))
	if err != nil {
		return 0, 0, err
	}
	files, err := unpackCardZip(cacheDir, card, data)
	return files, int64(len(data)), err
}

// updateCloudMark 只改索引里这条的标记。
func updateCloudMark(ctx context.Context, name, pid, card, mark string) error {
	idx, err := fetchCloudIndex(ctx, name, pid)
	if err != nil {
		return err
	}
	if idx.find(card) == nil {
		return nil // 云端没这条（还没上传过），下次录制时会带上
	}
	idx.setMark(card, mark)
	return putCloudIndex(ctx, name, pid, idx)
}

func entryRecorded(e cloudCardEntry) time.Time {
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02 15:04"} {
		if t, err := time.ParseInLocation(layout, e.Recorded, time.Local); err == nil {
			return t
		}
	}
	return time.Time{}
}

func entryToRow(e cloudCardEntry) accountRow {
	row := accountRow{
		Card:     e.Card,
		QQ:       dash(e.QQ),
		OpenID:   dash(e.OpenID),
		Note:     dash(e.Note),
		RawNote:  dash(e.RawNote),
		Zone:     dash(e.Zone),
		Expires:  "-",
		Recorded: entryRecorded(e),
		Remote:   true,
	}
	return row
}

// cleanField 把展示用的占位符 “-” 还原成空值。
func cleanField(s string) string {
	if s == "-" {
		return ""
	}
	return s
}

func dash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

// cloudSyncAll 把本地全部账号按新结构上传，并重建云端索引（也可用于旧 zip 迁移）。
func cloudSyncAll(ctx context.Context, name, pid, cacheDir string) (int, error) {
	rows := loadAccountRows(NewCache(cacheDir))
	if len(rows) == 0 {
		return 0, errors.New("本地没有可备份的缓存")
	}
	idx, err := fetchCloudIndex(ctx, name, pid)
	if err != nil {
		return 0, err
	}
	done := 0
	for _, r := range rows {
		e := cloudCardEntry{
			Card:    r.Card,
			QQ:      cleanField(r.QQ),
			OpenID:  cleanField(r.OpenID),
			Zone:    cleanField(r.Zone),
			Note:    cleanField(r.Note),
			RawNote: cleanField(r.RawNote),
		}
		if !r.Recorded.IsZero() {
			e.Recorded = r.Recorded.Format("2006-01-02 15:04:05")
		}
		if old := idx.find(r.Card); old != nil {
			e.Mark = old.Mark
		}
		data, _, err := packCardDir(cacheDir, r.Card)
		if err != nil {
			return done, err
		}
		c := newCloudClient()
		// 先建 <名称>/<类型>/，再建 cards/，否则父目录不存在会 409
		if err := c.ensureDir(ctx, cloudBaseDir(name, pid)+"/"); err != nil {
			return done, err
		}
		if err := c.ensureDir(ctx, cloudBaseDir(name, pid)+"/cards/"); err != nil {
			return done, err
		}
		if err := c.upload(ctx, cloudCardPath(name, pid, r.Card), data); err != nil {
			return done, err
		}
		e.Size = int64(len(data))
		idx.upsert(e)
		done++
	}
	if err := putCloudIndex(ctx, name, pid, idx); err != nil {
		return done, err
	}
	return done, nil
}

// cloudPullAll 只下载本地缺失的账号（按需），并返回云端条目数。
func cloudPullAll(ctx context.Context, name, pid, cacheDir string, log func(format string, args ...any)) (int, error) {
	idx, err := fetchCloudIndex(ctx, name, pid)
	if err != nil {
		return 0, err
	}
	if len(idx.Cards) == 0 {
		return 0, errCloudNotFound
	}
	for _, e := range idx.Cards {
		if hasLocalCard(cacheDir, e.Card) {
			continue
		}
		if _, _, err := downloadCardCloud(ctx, name, pid, cacheDir, e.Card); err != nil {
			return 0, err
		}
		if log != nil {
			log("已下载 %s", e.Card)
		}
	}
	return len(idx.Cards), nil
}
