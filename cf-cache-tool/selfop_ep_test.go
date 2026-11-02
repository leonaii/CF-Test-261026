package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDecodeLastLoginEP(t *testing.T) {
	raw, err := os.ReadFile(os.Getenv("APPDATA") + `\CG\data\cache\selfop\last_login.json`)
	if err != nil {
		t.Skip(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	ep, _ := m["ep"].(string)
	r, _ := m["used_r"].(string)
	if ep == "" || r == "" {
		t.Fatalf("ep/r missing len ep=%d r=%d", len(ep), len(r))
	}
	plain, err := selfEpDecode(ep, r)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	t.Logf("plain(%d): %s", len(plain), string(plain))
	var v any
	if err := json.Unmarshal(plain, &v); err != nil {
		t.Logf("not json: %v", err)
	}
}
