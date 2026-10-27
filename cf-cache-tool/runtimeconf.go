package main

// 运行时配置：dav 上的 CG/conf.json（只读）。
// 启动时后台拉一次，全程内存使用、不落盘；拉不到就用内置兜底值，不影响任何功能。
// 这样改口令/换服务地址只要改 dav 上的一个文件，不用重发客户端。

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// runtimeConf 是 conf.json 的结构，缺字段就用默认值。
type runtimeConf struct {
	AdminPass string `json:"admin_pass"` // 卡密平台管理口令
	APIBase   string `json:"api_base"`   // 卡密平台地址
}

var (
	rtCfg  atomic.Value // runtimeConf
	rtOnce sync.Once
)

func runtimeConfURL() string { return cloudBaseURL + cloudDir + "/conf.json" }

// loadRuntimeConf 拉一次配置，失败静默（后续按缺省处理）。
func loadRuntimeConf() {
	resp, err := davGet(runtimeConfURL(), 3*time.Second)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil || len(raw) == 0 {
		return
	}
	var c runtimeConf
	if json.Unmarshal(raw, &c) != nil {
		return
	}
	rtCfg.Store(c)
}

// ensureRuntimeConf 保证配置拉过一次（并发只拉一次）。
func ensureRuntimeConf() {
	rtOnce.Do(loadRuntimeConf)
}

// warmRuntimeConf 启动时异步预热，不阻塞界面。
func warmRuntimeConf() {
	go ensureRuntimeConf()
}

func rtCurrent() runtimeConf {
	if v, ok := rtCfg.Load().(runtimeConf); ok {
		return v
	}
	return runtimeConf{}
}

// adminPass 卡密平台管理口令：只从 dav 的 conf.json 取。
// 预热还没完成时同步等一次（最多 3 秒）；拉不到返回空串，登录会提示口令未加载。
func adminPass() string {
	if p := rtCurrent().AdminPass; p != "" {
		return p
	}
	ensureRuntimeConf()
	return rtCurrent().AdminPass
}

// cardServer 卡密平台地址：conf.json 优先，否则用内置兜底。
func cardServer() string {
	if b := rtCurrent().APIBase; b != "" {
		return b
	}
	return cardSrvDefault
}
