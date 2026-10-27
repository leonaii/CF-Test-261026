package main

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCloudClientBypassesProxy(t *testing.T) {
	t.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")
	c := newCloudHTTPClient(time.Second)
	tr, ok := c.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport 类型错误: %T", c.Transport)
	}
	if tr.Proxy != nil {
		t.Fatal("51i.cc 客户端不应读取或使用代理")
	}
}

func TestCloudNetworkErrorHidesEndpoint(t *testing.T) {
	err := cloudNetworkError("连接云端", &url.Error{
		Op:  "Get",
		URL: "https://cloud.51i.cc/remote.php/dav/files/leonai/CG/mm.md",
		Err: context.DeadlineExceeded,
	})
	msg := err.Error()
	for _, secret := range []string{"51i.cc", "remote.php", "https://", "/dav/"} {
		if strings.Contains(msg, secret) {
			t.Fatalf("错误提示泄露接口信息 %q: %s", secret, msg)
		}
	}
	if !strings.Contains(msg, "超时") {
		t.Fatalf("超时提示不明确: %s", msg)
	}
}
