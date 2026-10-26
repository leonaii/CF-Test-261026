//go:build !cg_secret

package main

// 本地调试/测试用的默认值。
// 正式包由 build.sh 生成 zz_build_secret.go 并加 -tags cg_secret 覆盖，
// 这样 dav 凭据只存在于生成文件里，不会进仓库，也能被 garble -literals 加密。

func buildToolVersion() string { return "0.0.0-dev" }
func buildCloudUser() string   { return "" }
func buildCloudPass() string   { return "" }
