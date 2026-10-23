package main

// 回放/离线上号会重写 config.ini：卡密固定成占位值，游戏目录留提示语，由使用者在界面上选。
// 原文件首次覆盖前备份为 config.ini.orig。
// 如果 card 已经是占位值，说明模板写过、游戏目录也选好了，此时不覆盖，
// 只按需调整 auto_zone 这一行（离线上号要 1，回放保持 0）。

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const placeholderCard = "CFAAAAAAAAAAAAAA"

// writeReplayConfig 回放用：auto_zone=0。
func writeReplayConfig(dir string) (backup string, changed bool, err error) {
	return writeConfig(dir, false)
}

// writeOfflineConfig 离线上号用：auto_zone=1（自动选区）。
func writeOfflineConfig(dir string) (backup string, changed bool, err error) {
	return writeConfig(dir, true)
}

// writeConfig 用固定模板覆盖 config.ini（GBK 编码，与原文件一致）。
// 返回备份路径以及是否真的写了文件。
func writeConfig(dir string, autoZone bool) (backup string, changed bool, err error) {
	p := filepath.Join(dir, "config.ini")
	backup = p + ".orig"
	zone := "0"
	if autoZone {
		zone = "1"
	}

	// 已经是占位卡密：模板早就写过，游戏目录也选好了，只调整 auto_zone
	if strings.EqualFold(strings.TrimSpace(configCard(dir)), placeholderCard) {
		ch, e := setConfigAutoZone(p, zone)
		return backup, ch, e
	}

	if _, statErr := os.Stat(backup); os.IsNotExist(statErr) {
		if old, readErr := os.ReadFile(p); readErr == nil {
			_ = os.WriteFile(backup, old, 0o600)
		}
	}
	tpl := "[ui]\r\n" +
		"card=" + placeholderCard + "\r\n" +
		"dir=请选择游戏目录\r\n" +
		"ssn=346\r\n" +
		"auto_zone=" + zone + "\r\n"
	data, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(tpl))
	if err != nil {
		return backup, false, err
	}
	return backup, true, os.WriteFile(p, data, 0o600)
}

// setConfigAutoZone 只改 auto_zone 行，其他内容（游戏目录等）原样保留。
// 没有这一行就追加到末尾。
func setConfigAutoZone(p, zone string) (bool, error) {
	raw, err := os.ReadFile(p)
	if err != nil {
		return false, err
	}
	txt, err := simplifiedchinese.GBK.NewDecoder().Bytes(raw)
	if err != nil {
		return false, err
	}

	want := "auto_zone=" + zone
	lines := strings.Split(string(txt), "\n")
	found := false
	for i, line := range lines {
		trim := strings.TrimSpace(strings.TrimRight(line, "\r"))
		if len(trim) >= 9 && strings.EqualFold(trim[:9], "auto_zone") {
			if strings.EqualFold(trim, want) {
				return false, nil // 已经是目标值
			}
			cr := ""
			if strings.HasSuffix(line, "\r") {
				cr = "\r"
			}
			lines[i] = want + cr
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, want)
	}

	out, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(strings.Join(lines, "\n")))
	if err != nil {
		return false, err
	}
	return true, os.WriteFile(p, out, 0o600)
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
