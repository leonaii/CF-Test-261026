package main

// 卡密白名单：dav 上的 CG/mm.md，一行一个卡密。
// 空行、# 开头的注释行忽略；比较时忽略大小写和首尾空格。
// 不在名单里的卡密不允许进入主界面。

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const licenseRemote = "mm.md"

var (
	errLicenseUnavailable = errors.New("卡密授权服务暂不可用")
	errLicenseDenied      = errors.New("卡密不在授权名单")
)

func isLicenseUnavailable(err error) bool { return errors.Is(err, errLicenseUnavailable) }
func isLicenseDenied(err error) bool      { return errors.Is(err, errLicenseDenied) }

// fetchLicenseList 拉取白名单，返回统一大写后的卡密集合。
func fetchLicenseList(timeout time.Duration) (map[string]struct{}, error) {
	cli := newCloudHTTPClient(timeout)
	req, err := http.NewRequest(http.MethodGet, cloudBaseURL+cloudDir+"/"+licenseRemote, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cloudUser, cloudPassword)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, cloudNetworkError("连接云端", err)
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
// 只请求一次并短超时；网络失败返回 errLicenseUnavailable，不能被当成“未授权”。
func checkCardLicense(card string) error {
	card = strings.ToUpper(strings.TrimSpace(card))
	if card == "" {
		return fmt.Errorf("卡密为空")
	}

	// CF 卡密直登：授权由卡密平台自己把关（查不到账号就上不了号），不受 mm.md 限制
	if looksLikeCard(card) {
		return nil
	}

	set, err := fetchLicenseList(4 * time.Second)
	if err != nil {
		return fmt.Errorf("%w：%v", errLicenseUnavailable, err)
	}
	if _, ok := set[card]; !ok {
		return fmt.Errorf("%w：%s", errLicenseDenied, card)
	}
	return nil
}
