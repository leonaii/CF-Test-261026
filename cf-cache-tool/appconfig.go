package main

// 本地配置：上号器类型、缓存（工作）目录、账号状态标记。
// 存在 %AppData%\CG\config.json，换目录/换机器互不影响。

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	markNormal  = "normal"
	markBan     = "ban"
	markBlack   = "black"
	markExpired = "expired"
	markOther   = "other"
)

type markOption struct {
	ID    string
	Label string
}

var markOptions = []markOption{
	{markNormal, "正常"},
	{markBan, "封禁"},
	{markBlack, "黑屋"},
	{markExpired, "过期"},
	{markOther, "其他"},
}

func markLabel(id string) string {
	for _, o := range markOptions {
		if o.ID == id {
			return o.Label
		}
	}
	return "正常"
}

type appConfig struct {
	Profile    string            `json:"profile"`
	ProfileVer int               `json:"profile_ver,omitempty"` // 老配置迁移标记
	WorkDir    string            `json:"work_dir"`
	GameDir    string            `json:"game_dir,omitempty"` // 自研上号：用户选过的游戏目录
	CloudName  string            `json:"cloud_name"`
	AutoCloud  bool              `json:"auto_cloud"`
	ProxyOn    *bool             `json:"proxy_on,omitempty"`    // nil = 默认开启
	DataSource string            `json:"data_source,omitempty"` // "" / online / dav
	Marks      map[string]string `json:"marks,omitempty"`
}

// defaultProfileID 是默认上号器：牛逼上号器（本地直登）。
const defaultProfileID = "2"

func configDirPath() string {
	d, err := os.UserConfigDir() // Windows: %AppData%
	if err != nil || d == "" {
		return "."
	}
	d = filepath.Join(d, "CG")
	_ = os.MkdirAll(d, 0o700)
	return d
}

func configFilePath() string { return filepath.Join(configDirPath(), "config.json") }

func loadAppConfig() *appConfig {
	raw, err := os.ReadFile(configFilePath())
	if err != nil {
		return nil
	}
	var c appConfig
	if json.Unmarshal(raw, &c) != nil {
		return nil
	}
	// 默认上号器升级为「牛逼」：只迁移一次，用户之后手动改回圆梦不会再被覆盖
	if c.ProfileVer < 2 {
		c.Profile = defaultProfileID
		c.ProfileVer = 2
		_ = c.save()
	}
	// 配置里只存可视卡密：读进来还原成真卡密，老文件顺手改写成可视形态
	c.normalizeCards()
	disk := c.diskView()
	if b, err := json.MarshalIndent(&disk, "", "  "); err == nil && !bytes.Equal(b, raw) {
		_ = c.save()
	}
	return &c
}

// normalizeCards 把配置里的可视卡密还原成真卡密（内存里一律用真卡密）。
func (c *appConfig) normalizeCards() {
	c.CloudName = loadCardKey(c.CloudName)
	if len(c.Marks) == 0 {
		return
	}
	m := make(map[string]string, len(c.Marks))
	for k, v := range c.Marks {
		m[loadCardKey(k)] = v
	}
	c.Marks = m
}

// diskView 返回落盘用的副本：cloud_name 与 marks 的键都换成可视卡密。
func (c *appConfig) diskView() appConfig {
	cp := *c
	cp.CloudName = storeCardKey(c.CloudName)
	if len(c.Marks) > 0 {
		m := make(map[string]string, len(c.Marks))
		for k, v := range c.Marks {
			m[storeCardKey(k)] = v
		}
		cp.Marks = m
	}
	return cp
}

func (c *appConfig) save() error {
	disk := c.diskView()
	raw, err := json.MarshalIndent(&disk, "", "  ")
	if err != nil {
		return err
	}
	p := configFilePath()
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

func (c *appConfig) mark(card string) string {
	if c == nil || c.Marks == nil {
		return markNormal
	}
	if s := c.Marks[card]; s != "" {
		return s
	}
	return markNormal
}

func (c *appConfig) setMark(card, status string) {
	if c.Marks == nil {
		c.Marks = map[string]string{}
	}
	if status == "" || status == markNormal {
		delete(c.Marks, card)
		return
	}
	c.Marks[card] = status
}

// dataSource 返回数据源（默认在线）。
func (c *appConfig) dataSource() string {
	if c == nil || strings.TrimSpace(c.DataSource) == "" {
		return srcOnline
	}
	if c.DataSource == srcDav {
		return srcDav
	}
	return srcOnline
}

// proxyDefault 代理开关默认值：没配置过就默认不开启（直连优先，不通再自动切代理）。
func (c *appConfig) proxyDefault() bool {
	if c == nil || c.ProxyOn == nil {
		return false
	}
	return *c.ProxyOn
}

func (c *appConfig) setProxy(on bool) {
	if c == nil {
		return
	}
	v := on
	c.ProxyOn = &v
}

// cloudNameOrDefault 自动备份/恢复使用的名称：优先用户填的，否则用上号器类型 ID。
func (c *appConfig) cloudNameOrDefault() string {
	if c != nil && c.CloudName != "" {
		return c.CloudName
	}
	if c != nil && c.Profile != "" {
		return c.Profile
	}
	return "cache"
}

// defaultWorkDir 默认工作目录：%AppData%\CG\data（上号器、host.dat、cache 都放这里）
func defaultWorkDir() string {
	return filepath.Join(configDirPath(), "data")
}

// migrateWorkDir 把老版本自定义的安装目录迁到默认目录（只复制 cache，不动原目录）。
func migrateWorkDir(cfg *appConfig) (string, bool) {
	if cfg == nil || cfg.WorkDir == "" {
		return "", false
	}
	def := defaultWorkDir()
	old := filepath.Clean(cfg.WorkDir)
	if old == filepath.Clean(def) {
		return "", false
	}

	if entries, err := os.ReadDir(filepath.Join(old, "cache")); err == nil {
		_ = os.MkdirAll(filepath.Join(def, "cache"), 0o700)
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			dst := filepath.Join(def, "cache", e.Name())
			if _, err := os.Stat(dst); err == nil {
				continue
			}
			_ = copyDir(filepath.Join(old, "cache", e.Name()), dst)
		}
	}

	cfg.WorkDir = def
	_ = cfg.save()
	return old, true
}
