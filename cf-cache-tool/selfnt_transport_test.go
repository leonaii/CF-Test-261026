package main

import (
	"reflect"
	"testing"
)

func TestSelfNTParseSSOEndpoints(t *testing.T) {
	text := `xx 1.2.3.4:10012; 5.6.7.8, bad 1.2.3.4.5, 9.10.11.12:0, 13.14.15.16:65535`
	// 与样本一致：扫描到 "1.2.3.4.5" 时会先拒绝 4 段形式，
	// 随后从第二个字符继续，得到 "2.3.4.5"。
	want := []string{"1.2.3.4:10012", "5.6.7.8:8080", "2.3.4.5:8080", "9.10.11.12:8080", "13.14.15.16:65535"}
	if got := selfNTParseSSOEndpoints(text); !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
