package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestLicenseErrorClassification(t *testing.T) {
	netErr := fmt.Errorf("%w：%v", errLicenseUnavailable, errors.New("连接云端超时"))
	if !isLicenseUnavailable(netErr) || isLicenseDenied(netErr) {
		t.Fatal("网络失败必须归为暂不可用，不能归为未授权")
	}
	denyErr := fmt.Errorf("%w：%s", errLicenseDenied, "LEONAI")
	if isLicenseUnavailable(denyErr) || !isLicenseDenied(denyErr) {
		t.Fatal("名单拒绝必须归为未授权")
	}
}

func TestCheckCardLicenseBypassesNetwork(t *testing.T) {
	// CF 卡密由卡密平台校验，不应触发白名单联网请求。
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:1")
	if err := checkCardLicense("ABCDEFGHIJKLMNOP"); err != nil {
		t.Fatalf("CF 卡密不应被白名单拦截: %v", err)
	}
}
