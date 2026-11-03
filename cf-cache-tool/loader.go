package main

// 上号器类型注册表。
//
// 上号器本体（exe）和 host.dat 不再打进 CG.exe，改由云端分发：
//
//	CG/loaders/<类型ID>/meta.json     版本、文件名、两个文件的 sha256
//	CG/loaders/<类型ID>/<exe>         上号器本体
//	CG/loaders/<类型ID>/host.dat       服务器列表
//
// 运行时按需下载并校验；本地已就绪则直接用，不再联网。

type loaderProfile struct {
	ID       string
	Name     string
	ExeName  string
	SiteName string // 卡网名称
	Site     string // 卡网链接
	SiteNote string // 购买说明

	// LocalOnly 表示本机已有 exe 时优先直用，不走云端校验（开发者本机）。
	LocalOnly bool
	// NoHostDat 表示这套上号器不读 host.dat（CG 用 WinHTTP 代理喂票据）。
	NoHostDat bool
}

var loaderProfiles = []loaderProfile{
	{
		ID:        "2",
		Name:      "牛逼上号器",
		ExeName:   "CF上号器.exe",
		SiteName:  "牛逼直登",
		SiteNote:  "离线票据桥接，不需要在线卡密",
		NoHostDat: true,
	},
	{
		ID:       "1",
		Name:     "圆梦上号器",
		ExeName:  "CF上号器_v1.17.exe",
		SiteName: "圆梦卡网",
		Site:     "https://www.517fk.cn/links/8B1B512A",
		SiteNote: "0V-10V / 10V-19V / 20V-49V / 50-99V 这家的",
	},
	{
		ID:        "3",
		Name:      "无畏上号器",
		ExeName:   wuweiExeName,
		SiteNote:  "无畏契约上号器，离线票据桥接",
		NoHostDat: true, // 不读 host.dat，靠改 exe 内置地址 + 本地假卡密服务
	},
	{
		ID:        "4",
		Name:      "无畏WG上号器",
		ExeName:   wuweiWGExeName,
		SiteNote:  "无畏契约 WeGame 注入版，离线票据桥接（需先登录官方 WeGame）",
		NoHostDat: true, // 与老版无畏同源，改 exe 内置地址 + 本地假卡密服务
	},
}

func findProfile(id string) *loaderProfile {
	for i := range loaderProfiles {
		if loaderProfiles[i].ID == id {
			return &loaderProfiles[i]
		}
	}
	if len(loaderProfiles) > 0 {
		return &loaderProfiles[0]
	}
	return nil
}

func profileNames() []string {
	names := make([]string, len(loaderProfiles))
	for i, p := range loaderProfiles {
		names[i] = p.Name
	}
	return names
}

func profileByIndex(idx int) *loaderProfile {
	if idx < 0 || idx >= len(loaderProfiles) {
		return findProfile("")
	}
	return &loaderProfiles[idx]
}

// 云端路径
func loaderRemoteDir(id string) string {
	if id == "" {
		id = "1"
	}
	return cloudDir + "/loaders/" + id
}

func loaderMetaPath(id string) string { return loaderRemoteDir(id) + "/meta.json" }

func loaderExePath(id, exe string) string { return loaderRemoteDir(id) + "/" + exe }

func loaderDatPath(id string) string { return loaderRemoteDir(id) + "/host.dat" }
