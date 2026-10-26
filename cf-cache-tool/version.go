package main

// 版本号：正式包由 build.sh 生成 zz_build_secret.go（-tags cg_secret）写死，
// 源码里不留构建产物信息；本地调试用 buildvars_default.go 的默认值。

import (
	"strconv"
	"strings"
)

var toolVersion = buildToolVersion()

// parseVersion 解析 1.2.3，解析失败返回 ok=false。
func parseVersion(s string) (major, minor, patch int, ok bool) {
	s = strings.TrimSpace(strings.TrimPrefix(strings.ToLower(s), "v"))
	if s == "" {
		return 0, 0, 0, false
	}
	parts := strings.Split(s, ".")
	if len(parts) > 3 {
		parts = parts[:3]
	}
	nums := make([]int, 3)
	for i := 0; i < 3; i++ {
		if i >= len(parts) {
			nums[i] = 0
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(parts[i]))
		if err != nil {
			return 0, 0, 0, false
		}
		nums[i] = n
	}
	return nums[0], nums[1], nums[2], true
}

// versionNewer 判断 remote 是否比 local 新。
func versionNewer(remote, local string) bool {
	ra, rb, rc, ok1 := parseVersion(remote)
	la, lb, lc, ok2 := parseVersion(local)
	if !ok1 || !ok2 {
		return false
	}
	if ra != la {
		return ra > la
	}
	if rb != lb {
		return rb > lb
	}
	return rc > lc
}
