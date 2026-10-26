package main

// 查封快照：在线查封成功后把结果按 QQ 落到本地，
// 离线数据源（dump 里没有 ck）也能看到"上次查封"的结果。
//
// 文件：%AppData%\CG\data\verify_snapshot.json
// 只存展示字段，不存 ck / cookie。

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const verifySnapMax = 8000

type verifySnap struct {
	QQ       string `json:"qq"`
	Card     string `json:"card"`
	Status   string `json:"status"`
	Fenghao  string `json:"fenghao"`
	Dongjie  string `json:"dongjie"`
	Zaixian  string `json:"zaixian"`
	Renlian  string `json:"renlian"`
	Shiming  string `json:"shiming"`
	Lingpai  string `json:"lingpai"`
	Nickname string `json:"nickname"`
	Level    string `json:"level"`
	Daqu     string `json:"daqu"`
	Zhanqu   string `json:"zhanqu"`
	Zone     string `json:"zone"`
	Message  string `json:"message"`
	At       int64  `json:"at"` // 毫秒
}

var verifySnapMu sync.Mutex

func verifySnapPath() string {
	return filepath.Join(configDirPath(), "data", "verify_snapshot.json")
}

// loadVerifySnaps 读全部快照，键是 QQ。
func loadVerifySnaps() map[string]verifySnap {
	out := map[string]verifySnap{}
	raw, err := os.ReadFile(verifySnapPath())
	if err != nil || len(raw) == 0 {
		return out
	}
	var list []verifySnap
	if json.Unmarshal(raw, &list) != nil {
		return out
	}
	for _, s := range list {
		if s.QQ != "" {
			out[s.QQ] = s
		}
	}
	return out
}

// saveVerifySnap 记录一条快照；同一个 QQ 只保留最新。
func saveVerifySnap(card string, st verifyStatus) {
	qq := strings.TrimSpace(st.QQ)
	if !st.OK || qq == "" {
		return
	}

	verifySnapMu.Lock()
	defer verifySnapMu.Unlock()

	all := loadVerifySnaps()
	entry := verifySnap{
		QQ:       qq,
		Card:     card,
		Status:   st.Status,
		Fenghao:  st.Fenghao,
		Dongjie:  st.Dongjie,
		Zaixian:  st.Zaixian,
		Renlian:  st.Renlian,
		Shiming:  st.Shiming,
		Lingpai:  st.Lingpai,
		Nickname: st.Nickname,
		Level:    st.Level,
		Daqu:     st.Daqu,
		Zhanqu:   st.Zhanqu,
		Zone:     st.Zone,
		Message:  st.Message,
		At:       st.At.UnixMilli(),
	}
	if entry.At == 0 {
		entry.At = time.Now().UnixMilli()
	}
	all[qq] = entry

	list := make([]verifySnap, 0, len(all))
	for _, s := range all {
		list = append(list, s)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].At > list[j].At })
	if len(list) > verifySnapMax {
		list = list[:verifySnapMax]
	}

	data, err := json.Marshal(list)
	if err != nil {
		return
	}
	p := verifySnapPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return
	}
	_ = os.Rename(tmp, p)
}

// getVerifySnap 按 QQ 取一条快照。
func getVerifySnap(qq string) (verifySnap, bool) {
	qq = strings.TrimSpace(qq)
	if qq == "" {
		return verifySnap{}, false
	}
	verifySnapMu.Lock()
	defer verifySnapMu.Unlock()
	s, ok := loadVerifySnaps()[qq]
	return s, ok
}

// snapStatus 把快照还原成 verifyStatus，方便复用同一个弹窗。
func (s verifySnap) snapStatus() verifyStatus {
	st := verifyStatus{
		OK:       true,
		HasData:  true,
		QueryCode: s.Card,
		QQ:       s.QQ,
		Zone:     s.Zone,
		Status:   s.Status,
		Fenghao:  s.Fenghao,
		Dongjie:  s.Dongjie,
		Zaixian:  s.Zaixian,
		Renlian:  s.Renlian,
		Shiming:  s.Shiming,
		Lingpai:  s.Lingpai,
		Nickname: s.Nickname,
		Level:    s.Level,
		Daqu:     s.Daqu,
		Zhanqu:   s.Zhanqu,
		Message:  s.Message,
	}
	if s.At > 0 {
		st.At = time.UnixMilli(s.At)
	}
	return st
}
