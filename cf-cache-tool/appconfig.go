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
