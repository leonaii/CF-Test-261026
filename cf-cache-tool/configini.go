package main

// 回放模式会重写 config.ini：卡密固定成占位值，游戏目录留提示语，由使用者在界面上选。
// 原文件首次覆盖前备份为 config.ini.orig。
// 如果 card 已经是占位值，说明模板写过、游戏目录也选好了，此时不再覆盖，保留现有内容。

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	placeholderCard = "CFAAAAAAAAAAAAAA"
	replayConfigTpl = "[ui]\r\n" +
		"card=" + placeholderCard + "\r\n" +
		"dir=请选择游戏目录\r\n" +
		"ssn=346\r\n" +
		"auto_zone=0\r\n"
)

// writeReplayConfig 用固定模板覆盖 config.ini（GBK 编码，与原文件一致）。
// 返回备份路径以及是否真的写了文件。
func writeReplayConfig(dir string) (backup string, changed bool, err error) {
	p := filepath.Join(dir, "config.ini")
	backup = p + ".orig"

	// 已经是占位卡密：模板早就写过，游戏目录也选好了，不动它
	if strings.EqualFold(strings.TrimSpace(configCard(dir)), placeholderCard) {
		return backup, false, nil
	}

	if _, statErr := os.Stat(backup); os.IsNotExist(statErr) {
		if old, readErr := os.ReadFile(p); readErr == nil {
			_ = os.WriteFile(backup, old, 0o600)
		}
	}
	data, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(replayConfigTpl))
	if err != nil {
		return backup, false, err
	}
	return backup, true, os.WriteFile(p, data, 0o600)
}

// configCard 读取 config.ini 里的 card= 值（纯 ASCII，不需要 GBK 解码）。
func configCard(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "config.ini"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.TrimRight(line, "\r"))
		if len(line) > 5 && strings.EqualFold(line[:5], "card=") {
			return strings.TrimSpace(line[5:])
		}
	}
	return ""
}
