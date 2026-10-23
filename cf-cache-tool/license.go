package main

// 卡密白名单：dav 上的 CG/mm.md，一行一个卡密。
// 空行、# 开头的注释行忽略；比较时忽略大小写和首尾空格。
// 不在名单里的卡密不允许进入主界面。

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const licenseRemote = "mm.md"

// fetchLicenseList 拉取白名单，返回统一大写后的卡密集合。
func fetchLicenseList(timeout time.Duration) (map[string]struct{}, error) {
	cli := &http.Client{Timeout: timeout}
	req, err := http.NewRequest(http.MethodGet, cloudBaseURL+cloudDir+"/"+licenseRemote, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cloudUser, cloudPassword)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("连接云端失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("云端返回 HTTP %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	set := map[string]struct{}{}
	for _, line := range strings.Split(string(raw), "\n") {
		s := strings.TrimSpace(strings.TrimPrefix(line, "\ufeff"))
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		set[strings.ToUpper(s)] = struct{}{}
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("授权名单为空")
	}
	return set, nil
}

// looksLikeCard 判断输入是不是 16 位纯大写字母的卡密（如 CFTVNKFQZBYTGXYO）。
// 命中则首屏直接走离线上号，不再当成普通名称。
func looksLikeCard(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) != 16 {
		return false
	}
	for _, c := range s {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

// checkCardLicense 校验卡密是否在授权名单里。
// 网络失败会重试一次；拉不到名单一律按拒绝处理。
func checkCardLicense(card string) error {
	card = strings.ToUpper(strings.TrimSpace(card))
	if card == "" {
		return fmt.Errorf("卡密为空")
	}

	// CF 卡密直登：授权由卡密平台自己把关（查不到账号就上不了号），不受 mm.md 限制
	if looksLikeCard(card) {
		return nil
	}

	var lastErr error
	for try := 1; try <= 2; try++ {
		set, err := fetchLicenseList(8 * time.Second)
		if err != nil {
			lastErr = err
			if try == 1 {
				time.Sleep(600 * time.Millisecond)
			}
			continue
		}
		if _, ok := set[card]; ok {
			return nil
		}
		return fmt.Errorf("卡密 %s 不在授权名单里", card)
	}
	return fmt.Errorf("无法校验卡密（联网失败）：%v", lastErr)
}
