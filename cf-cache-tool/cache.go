package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Meta 描述一次被缓存的 HTTP 响应。
type Meta struct {
	Card       string              `json:"card"`
	Endpoint   string              `json:"endpoint"`
	Status     int                 `json:"status"`
	Header     map[string][]string `json:"header"`
	RecordedAt time.Time           `json:"recorded_at"`
	Bytes      int                 `json:"bytes"`
	Request    string              `json:"request,omitempty"`
}

// Cache 以 卡密 + 接口 为键保存原始响应字节。
type Cache struct {
	Dir string
}

func NewCache(dir string) *Cache {
	return &Cache{Dir: dir}
}

func sanitizeName(s string) string {
	if s == "" {
		return "unknown"
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			sum := sha256.Sum256([]byte(s))
			return "h" + hex.EncodeToString(sum[:8])
		}
	}
	return b.String()
}

// endpointKey 把请求路径转成缓存文件名。
// 带查询参数时追加短哈希，避免 /a?x=1 与 /a?x=2 互相覆盖；超长路径同样折叠成哈希。
func endpointKey(path, query string) string {
	p := strings.Trim(path, "/")
	if p == "" {
		p = "root"
	}
	key := strings.ReplaceAll(p, "/", "_")
	if query != "" {
		sum := sha256.Sum256([]byte(query))
		key += "_q" + hex.EncodeToString(sum[:4])
	}
	if len(key) > 100 {
		sum := sha256.Sum256([]byte(path + "?" + query))
		key = "h" + hex.EncodeToString(sum[:8])
	}
	return sanitizeName(key)
}

// knownEndpoints 是逆向时确认过的接口，其余一律按“新接口”记录。
func isKnownEndpoint(ep string) bool {
	switch ep {
	case "api_card_boot", "api_card_peek", "api_card_log", "api_card_redeem":
		return true
	}
	return false
}

func (c *Cache) files(card, endpoint string) (bodyPath, metaPath string) {
	dir := filepath.Join(c.Dir, sanitizeName(card))
	return filepath.Join(dir, endpoint+".body"), filepath.Join(dir, endpoint+".meta.json")
}

// Load 读取缓存的响应体与元信息。
func (c *Cache) Load(card, endpoint string) ([]byte, *Meta, error) {
	bodyPath, metaPath := c.files(card, endpoint)
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, nil, err
	}
	var m Meta
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, nil, err
	}
	body, err := os.ReadFile(bodyPath)
	if err != nil {
		return nil, nil, err
	}
	return body, &m, nil
}

// Save 原子写入响应体与元信息。
func (c *Cache) Save(card, endpoint string, m *Meta, body []byte) error {
	bodyPath, metaPath := c.files(card, endpoint)
	if err := os.MkdirAll(filepath.Dir(bodyPath), 0o700); err != nil {
		return err
	}
	m.Card = card
	m.Endpoint = endpoint
	m.Bytes = len(body)
	m.RecordedAt = time.Now()
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err := writeFileAtomic(bodyPath, body); err != nil {
		return err
	}
	return writeFileAtomic(metaPath, raw)
}

// writeFileAtomic 先写临时文件再改名，避免并发读到半截内容。
func writeFileAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Delete 删除某个卡密的全部缓存。
func (c *Cache) Delete(card string) error {
	return os.RemoveAll(filepath.Join(c.Dir, sanitizeName(card)))
}

// List 返回全部缓存元信息。
func (c *Cache) List() ([]Meta, error) {
	var out []Meta
	entries, err := os.ReadDir(c.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		metas, err := filepath.Glob(filepath.Join(c.Dir, e.Name(), "*.meta.json"))
		if err != nil {
			continue
		}
		for _, mp := range metas {
			raw, err := os.ReadFile(mp)
			if err != nil {
				continue
			}
			var m Meta
			if json.Unmarshal(raw, &m) == nil {
				out = append(out, m)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Card == out[j].Card {
			return out[i].Endpoint < out[j].Endpoint
		}
		return out[i].Card < out[j].Card
	})
	return out, nil
}

func (m *Meta) String() string {
	return fmt.Sprintf("card=%s endpoint=%s status=%d size=%d at=%s",
		m.Card, m.Endpoint, m.Status, m.Bytes, m.RecordedAt.Format("2006-01-02 15:04:05"))
}

// Has 判断是否存在某个接口的缓存（不区分卡密）。
// 用于 boot 这类请求体里不带卡密的接口。
func (c *Cache) Has(endpoint string) bool {
	entries, err := os.ReadDir(c.Dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		body := filepath.Join(c.Dir, e.Name(), endpoint+".body")
		if st, err := os.Stat(body); err == nil && !st.IsDir() {
			return true
		}
	}
	return false
}

// ClearCard 删除某个卡密的整个缓存目录（重新录制时用，保证不残留旧数据）。
func (c *Cache) ClearCard(card string) error {
	if card == "" {
		return nil
	}
	return os.RemoveAll(filepath.Join(c.Dir, sanitizeName(card)))
}

// ClearEndpoint 删除所有卡密目录下某个接口的缓存（boot 改成固定返回后用来清掉旧数据）。
func (c *Cache) ClearEndpoint(endpoint string) error {
	entries, err := os.ReadDir(c.Dir)
	if err != nil {
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		base := filepath.Join(c.Dir, e.Name(), endpoint)
		_ = os.Remove(base + ".body")
		_ = os.Remove(base + ".meta.json")
	}
	return nil
}

// stagingDirName 是录制时的临时目录名，录制成功后再按卡密合并进正式缓存。
const stagingDirName = "cache.recording"

// mergeRecording 把录制目录里已完成的卡密合并进正式缓存：
// 有 api_card_redeem.body 才算录制成功，整目录覆盖；没有 redeem 的（只录到一半）跳过。
// 返回合并数量与跳过的目录名。
func mergeRecording(staging, target string) (int, []string, error) {
	entries, err := os.ReadDir(staging)
	if err != nil {
		return 0, nil, err
	}
	merged := 0
	var skipped []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		card := e.Name()
		src := filepath.Join(staging, card)
		if _, err := os.Stat(filepath.Join(src, "api_card_redeem.body")); err != nil {
			skipped = append(skipped, card)
			continue
		}
		dst := filepath.Join(target, card)
		if err := os.RemoveAll(dst); err != nil {
			return merged, skipped, err
		}
		if err := os.MkdirAll(target, 0o700); err != nil {
			return merged, skipped, err
		}
		if err := os.Rename(src, dst); err == nil {
			merged++
			continue
		}
		// 跨盘等情况下退化为拷贝
		if err := copyDir(src, dst); err != nil {
			return merged, skipped, err
		}
		merged++
	}
	return merged, skipped, nil
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o700); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(src, e.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dst, e.Name()), data, 0o600); err != nil {
			return err
		}
	}
	return nil
}

// hasContent 判断目录里是否有子目录（用于发现上次没合并的录制）。
func hasContent(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() {
			return true
		}
	}
	return false
}

// hasLocalCard 判断本地是否已有该卡密的 redeem 缓存。
func hasLocalCard(cacheDir, card string) bool {
	_, err := os.Stat(filepath.Join(cacheDir, sanitizeName(card), "api_card_redeem.body"))
	return err == nil
}

// stagedCards 列出录制临时目录里已完成的卡密（有 redeem 的）。
func stagedCards(staging string) []string {
	entries, err := os.ReadDir(staging)
	if err != nil {
		return nil
	}
	var cards []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(staging, e.Name(), "api_card_redeem.body")); err == nil {
			cards = append(cards, e.Name())
		}
	}
	sort.Strings(cards)
	return cards
}
