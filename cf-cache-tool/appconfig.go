package main

// 本地配置：上号器类型、缓存（工作）目录、账号状态标记。
// 存在 %AppData%\CG\config.json，换目录/换机器互不影响。

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	Profile   string            `json:"profile"`
	WorkDir   string            `json:"work_dir"`
	CloudName string            `json:"cloud_name"`
	AutoCloud bool              `json:"auto_cloud"`
	ProxyOn   *bool             `json:"proxy_on,omitempty"` // nil = 默认开启
	Marks     map[string]string `json:"marks,omitempty"`
}

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
	return &c
}

func (c *appConfig) save() error {
	raw, err := json.MarshalIndent(c, "", "  ")
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

// proxyDefault 代理开关默认值：没配置过就默认开启（直连 IP 容易被验号接口限频/封）。
func (c *appConfig) proxyDefault() bool {
	if c == nil || c.ProxyOn == nil {
		return true
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
