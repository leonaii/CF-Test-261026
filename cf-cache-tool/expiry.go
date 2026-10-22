package main

// 把响应里的 expires_at 统一延长到 10 年后。
// 缓存回放的意义就是长期可用，而服务端给的到期时间可能只有几小时/几天。

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const extendYears = 10

// 匹配 "expires_at": 1234567890123 / "expires_at": "2026-10-23 22:29:20"
// 兼容 expiresAt / expire_at / expireat 等写法。
var expiresRe = regexp.MustCompile(`(?i)("expires?_?at"\s*:\s*)("([^"]*)"|(-?\d+))`)

var timeLayouts = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"2006/01/02 15:04:05",
}

// extendExpiry 把 body 里所有 expires_at 改成 now+10 年，保持原来的类型与格式。
// 返回新的 body 以及是否发生改动。
func extendExpiry(body []byte, now time.Time) ([]byte, bool) {
	if len(body) == 0 || !strings.Contains(strings.ToLower(string(body)), "expire") {
		return body, false
	}
	target := now.AddDate(extendYears, 0, 0)
	changed := false

	out := expiresRe.ReplaceAllFunc(body, func(m []byte) []byte {
		sub := expiresRe.FindSubmatch(m)
		if len(sub) < 5 {
			return m
		}
		prefix := sub[1]

		// 字符串形式："2026-10-23 22:29:20"
		if raw := string(sub[3]); raw != "" {
			for _, layout := range timeLayouts {
				if _, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
					changed = true
					return append(append([]byte{}, prefix...), []byte(strconv.Quote(target.Format(layout)))...)
				}
			}
			if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
				changed = true
				return append(append([]byte{}, prefix...), []byte(strconv.Quote(formatByUnit(n, target)))...)
			}
			return m
		}

		// 数字形式：1792765760283（毫秒）或 1792765760（秒）
		n, err := strconv.ParseInt(string(sub[4]), 10, 64)
		if err != nil {
			return m
		}
		changed = true
		return append(append([]byte{}, prefix...), []byte(formatByUnit(n, target))...)
	})
	if !changed {
		return body, false
	}
	return out, true
}

func formatByUnit(n int64, target time.Time) string {
	switch {
	case n >= 1e12: // 毫秒时间戳
		return strconv.FormatInt(target.UnixMilli(), 10)
	case n >= 1e9: // 秒时间戳
		return strconv.FormatInt(target.Unix(), 10)
	}
	return strconv.FormatInt(n, 10)
}
