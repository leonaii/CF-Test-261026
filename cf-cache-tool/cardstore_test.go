package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testCard = "CFTVNKFQZBYTGXYO"
const testCard2 = "CFYAWOJVXAROMRZA"

func TestStoreCardKey(t *testing.T) {
	stored := storeCardKey(testCard)
	if stored != renderCardKey(testCard) || len(stored) != 32 {
		t.Fatalf("落盘形态不对: %s", stored)
	}
	if loadCardKey(stored) != testCard {
		t.Fatalf("还原失败: %s", loadCardKey(stored))
	}
	// 老数据/非卡密原样返回
	if loadCardKey(testCard) != testCard || loadCardKey("leonai") != "leonai" {
		t.Fatal("非可视形态不应被改动")
	}
}

func TestCacheStoresRenderedCard(t *testing.T) {
	dir := t.TempDir()
	c := NewCache(dir)

	req := `{"card_key":"` + testCard + `"}`
	if err := c.Save(testCard, "api_card_redeem", &Meta{Status: 200, Request: req}, []byte(`{"ok":true}`)); err != nil {
		t.Fatalf("写缓存失败: %v", err)
	}

	// 目录名必须是可视卡密，磁盘上不能出现明文目录
	if _, err := os.Stat(filepath.Join(dir, cardDirName(testCard))); err != nil {
		t.Fatalf("缓存目录不是可视卡密: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, testCard)); err == nil {
		t.Fatal("磁盘上出现了明文卡密目录")
	}

	metaPath := filepath.Join(dir, cardDirName(testCard), "api_card_redeem.meta.json")
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("读 meta 失败: %v", err)
	}
	if strings.Contains(string(raw), testCard) {
		t.Fatalf("meta 里还有明文卡密: %s", raw)
	}
	if !strings.Contains(string(raw), storeCardKey(testCard)) {
		t.Fatalf("meta 里没有可视卡密: %s", raw)
	}

	body, m, err := c.Load(testCard, "api_card_redeem")
	if err != nil || m.Card != testCard || string(body) != `{"ok":true}` {
		t.Fatalf("Load 还原失败: card=%q body=%q err=%v", m.Card, body, err)
	}

	list, err := c.List()
	if err != nil || len(list) != 1 || list[0].Card != testCard {
		t.Fatalf("List 还原失败: %+v err=%v", list, err)
	}
	if !hasLocalCard(dir, testCard) {
		t.Fatal("hasLocalCard 没命中可视卡密目录")
	}
}

func TestHistoryStoresRenderedCard(t *testing.T) {
	dir := t.TempDir()
	note := "编号:X----123----" + testCard + "----等级:101"
	appendHistory(dir, histEntry{TS: 1, Action: "offline", QQ: "123", Card: testCard, Note: note})

	raw, err := os.ReadFile(historyPath(dir))
	if err != nil {
		t.Fatalf("读 history 失败: %v", err)
	}
	if strings.Contains(string(raw), testCard) {
		t.Fatalf("history 里还有明文卡密: %s", raw)
	}

	list := loadHistory(dir, 10)
	if len(list) != 1 || list[0].Card != testCard {
		t.Fatalf("history 还原失败: %+v", list)
	}
	if err := deleteHistory(dir, testCard, 1); err != nil {
		t.Fatalf("按真卡密删除历史失败: %v", err)
	}
}

func TestConfigStoresRenderedCard(t *testing.T) {
	t.Setenv("AppData", t.TempDir())

	cfg := &appConfig{
		Profile:   "1",
		WorkDir:   filepath.Join(t.TempDir(), "data"),
		CloudName: testCard,
		Marks:     map[string]string{testCard2: markBan},
	}
	if err := cfg.save(); err != nil {
		t.Fatalf("保存配置失败: %v", err)
	}
	raw, err := os.ReadFile(configFilePath())
	if err != nil {
		t.Fatalf("读配置失败: %v", err)
	}
	if strings.Contains(string(raw), testCard) || strings.Contains(string(raw), testCard2) {
		t.Fatalf("config 里还有明文卡密: %s", raw)
	}

	back := loadAppConfig()
	if back == nil || back.CloudName != testCard {
		t.Fatalf("配置还原失败: %+v", back)
	}
	if back.mark(testCard2) != markBan {
		t.Fatalf("标记没还原: %q", back.mark(testCard2))
	}
	if back.mark(storeCardKey(testCard2)) != markNormal {
		t.Fatal("内存里不应该用可视卡密当键")
	}
}

func TestMigrateLocalCardStore(t *testing.T) {
	t.Setenv("AppData", t.TempDir())
	dir := t.TempDir()

	// 老版本：明文卡密目录 + 明文 meta
	old := filepath.Join(dir, "cache", testCard)
	if err := os.MkdirAll(old, 0o700); err != nil {
		t.Fatal(err)
	}
	meta, _ := json.Marshal(Meta{
		Card:     testCard,
		Endpoint: "api_card_redeem",
		Status:   200,
		Request:  `{"card_key":"` + testCard + `"}`,
	})
	if err := os.WriteFile(filepath.Join(old, "api_card_redeem.meta.json"), meta, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(old, "api_card_redeem.body"), []byte("body"), 0o600); err != nil {
		t.Fatal(err)
	}
	// 老版本：明文 history
	line, _ := json.Marshal(histEntry{
		TS: 1, Action: "offline", QQ: "123", Card: testCard,
		Note: "编号:X----123----" + testCard + "----等级:101",
	})
	if err := os.WriteFile(historyPath(dir), append(line, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}

	migrateLocalCardStore(dir)
	migrateLocalCardStore(dir) // 可重复执行

	if _, err := os.Stat(filepath.Join(dir, "cache", testCard)); err == nil {
		t.Fatal("迁移后明文目录还在")
	}
	if !hasLocalCard(filepath.Join(dir, "cache"), testCard) {
		t.Fatal("迁移后读不到缓存")
	}
	metaPath := filepath.Join(dir, "cache", cardDirName(testCard), "api_card_redeem.meta.json")
	rawMeta, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("迁移后读 meta 失败: %v", err)
	}
	if strings.Contains(string(rawMeta), testCard) {
		t.Fatalf("迁移后 meta 仍是明文: %s", rawMeta)
	}
	rawHist, err := os.ReadFile(historyPath(dir))
	if err != nil {
		t.Fatalf("迁移后读 history 失败: %v", err)
	}
	if strings.Contains(string(rawHist), testCard) {
		t.Fatalf("迁移后 history 仍是明文: %s", rawHist)
	}
	if list := loadHistory(dir, 5); len(list) != 1 || list[0].Card != testCard {
		t.Fatalf("迁移后 history 读不回真卡密: %+v", list)
	}
}

func TestDumpCardLineStoresRenderedCard(t *testing.T) {
	line := []byte(`{"card_key":"` + testCard + `","qq":"123","note":"编号:X----123----` + testCard + `----等级:101","created_at":1700000000000}`)
	out := transformDumpCardLine(line)
	if strings.Contains(string(out), testCard) {
		t.Fatalf("离线列表里还有明文卡密: %s", out)
	}
	if !strings.Contains(string(out), storeCardKey(testCard)) {
		t.Fatalf("离线列表没有可视卡密: %s", out)
	}

	var obj map[string]any
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("离线列表不是合法 JSON: %v", err)
	}
	it := cardItemFromMap(obj)
	if it.CardKey != storeCardKey(testCard) {
		t.Fatalf("磁盘上 card_key 应为可视卡密: %q", it.CardKey)
	}
	if got := loadCardKey(it.CardKey); got != testCard {
		t.Fatalf("离线列表还原失败: %q", got)
	}
	// 再跑一次不应变化
	if again := transformDumpCardLine(out); string(again) != string(out) {
		t.Fatal("重复处理应保持稳定")
	}
}
