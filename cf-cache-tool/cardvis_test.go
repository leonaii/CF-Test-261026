package main

import (
	"strings"
	"testing"
)

func TestCardVisRoundTrip(t *testing.T) {
	cards := []string{
		"CFTVNKFQZBYTGXYO",
		"CFYAWOJVXAROMRZA",
		"AAAAAAAAAAAAAAAA",
		"ZZZZZZZZZZZZZZZZ",
		"CGABCDEFGHIJKLMN",
	}
	for _, c := range cards {
		v := renderCardKey(c)
		if len(v) != cardVisLen {
			t.Fatalf("%s 渲染长度 %d，期望 %d（%s）", c, len(v), cardVisLen, v)
		}
		if !strings.HasPrefix(v, cardVisPrefix) {
			t.Fatalf("%s 渲染结果没有 CG 前缀: %s", c, v)
		}
		if got, want := v[len(cardVisPrefix):len(cardVisPrefix)+cardVisPadWidth], strings.Repeat("0", cardVisPadWidth); got != want {
			t.Fatalf("%s 补齐位不对: %q", c, got)
		}
		back, ok := decodeCardKey(v)
		if !ok || back != c {
			t.Fatalf("%s -> %s 还原失败: %q ok=%v", c, v, back, ok)
		}
		// 大小写、首尾空格不影响还原
		if back, ok := decodeCardKey("  " + strings.ToLower(v) + "  "); !ok || back != c {
			t.Fatalf("%s 小写还原失败: %q ok=%v", c, back, ok)
		}
	}
}

func TestCardVisDeterministic(t *testing.T) {
	const card = "CFTVNKFQZBYTGXYO"
	a, b := renderCardKey(card), renderCardKey(card)
	if a != b {
		t.Fatalf("同一张卡两次渲染不一致: %s / %s", a, b)
	}
	if a == renderCardKey("CFTVNKFQZBYTGXYP") {
		t.Fatal("不同卡密不应渲染成同一个串")
	}
}

func TestCardVisRejectsTamper(t *testing.T) {
	v := []byte(renderCardKey("CFTVNKFQZBYTGXYO"))
	for i := 2; i < len(v); i++ {
		bad := append([]byte(nil), v...)
		if bad[i] == 'A' {
			bad[i] = 'B'
		} else {
			bad[i] = 'A'
		}
		if got, ok := decodeCardKey(string(bad)); ok {
			t.Fatalf("改坏第 %d 位仍被接受: %q -> %s", i, string(bad), got)
		}
	}
	for _, s := range []string{"", "CG", strings.Repeat("0", 32), "CG000000" + strings.Repeat("A", 24)} {
		if _, ok := decodeCardKey(s); ok {
			t.Fatalf("非法串被接受: %q", s)
		}
	}
}

func TestRenderCardKeyLeavesNames(t *testing.T) {
	for _, s := range []string{"leonai", "LEONAI", "1", "CFTVNKFQZBYTGXY", ""} {
		if got := renderCardKey(s); got != s {
			t.Fatalf("%q 不该被渲染: %q", s, got)
		}
	}
}

func TestResolveCardKey(t *testing.T) {
	const card = "CFTVNKFQZBYTGXYO"
	if got, ok := resolveCardKey(renderCardKey(card)); !ok || got != card {
		t.Fatalf("可视卡密没还原: %q ok=%v", got, ok)
	}
	if got, ok := resolveCardKey(" leonai "); !ok || got != "leonai" {
		t.Fatalf("普通名称被改动: %q ok=%v", got, ok)
	}
	// 像可视卡密但解不开：必须报错，不能当名称放过去
	if _, ok := resolveCardKey("CG000000" + strings.Repeat("A", 24)); ok {
		t.Fatal("解不开的可视卡密应按非法处理")
	}
}
