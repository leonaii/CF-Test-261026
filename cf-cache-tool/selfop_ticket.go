package main

// 自研上号：从本地录制缓存里解出可直接用于 LoginByQq 的票据。

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// selfMidListTemplate 是上号器为 machine_id=8AC936BFCC316D3D660CB564B40F0188
// 生成的真实 mid_list（抓包取得）。第一段是 machine_id，第二段是设备侧 blob。
const selfMidListTemplate = "CiQIABIgOGFjOTM2YmZjYzMxNmQzZDY2MGNiNTY0YjQwZjAxODgK5wIIARLiAkFRQUFBQUFBQUFBQ0FBQUFBQUFBQU55QkFielpnbS9FQmdBQUFBRUFBQUN4SFFBQVRRQlRBQzBBTVFBekFEQUFOZ0EyQURjQU5BQTRBRGtBT1FBM0FEQUFBQUFBQUE9PV83MzRFRDA0QTMzOTIxQjA1NTREQ0E2NzI1NkVDQkY3NDk5QjcxMUQxOTZGQzc5MTUxMEEzNUNGNUFFNDI0MDZENDQ5ODRFNzNGNEEzM0IzOTk3RUY5MjU1M0VBMkIxNzhBMkNDRTdEOTM1NEExQzBEMEFFQzg4MjhBNkM2ODVEN18yRjdGMTc2NTUwNDIzRDhEMTAxMEZDRTcyNzNERkVCQTAzRjY2NzQ1OEJCRDk2NkU3REEzNTI0NTBGN0E1ODcyNENDOEVBMDk0QjhDMjU0RjQzNjcyMjI0N0NFNDU2Mzg1N0UwMjc2MjJFMThCMzVGMTczRjQyRjU5NDM2OEZFRQ=="

// selfTicket 是 LoginByQq 需要的最小字段集合。
type selfTicket struct {
	Card      string
	QQ        string
	UID       string
	Sig       string
	MachineID string
	CT        string
	Ctid      string
	Raw       map[string]string
}

func (t *selfTicket) String() string {
	return fmt.Sprintf("QQ=%s uid=%s machine=%s sig=%s", t.QQ, t.UID, t.MachineID, shortHex(t.Sig, 8))
}

func shortHex(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// loadSelfTicket 读取并解密某个卡密的 redeem 缓存。
func loadSelfTicket(c *Cache, card string) (*selfTicket, error) {
	if c == nil {
		return nil, errors.New("缓存未初始化")
	}
	body, _, err := c.Load(card, "api_card_redeem")
	if err != nil {
		return nil, fmt.Errorf("没有该卡密的 redeem 缓存: %w", err)
	}
	plain, err := OpenEnvelope(body)
	if err != nil {
		return nil, fmt.Errorf("解信封失败: %w", err)
	}
	kv, section := parsePayload(plain)
	t := &selfTicket{Card: card, Raw: kv}
	t.QQ = firstNonEmpty(kv["qq"], section, kv["qqnumber"])
	t.UID = firstNonEmpty(kv["uid"], kv["openid"])
	t.Sig = kv["sig"]
	t.MachineID = strings.ToUpper(firstNonEmpty(kv["machine_id"], kv["machineid"]))
	t.CT = kv["ct"]
	t.Ctid = kv["ctid"]

	if t.QQ == "" {
		return nil, errors.New("redeem 明文里没有 QQ 号")
	}
	if t.Sig == "" {
		return nil, errors.New("redeem 明文里没有 sig")
	}
	if t.MachineID == "" {
		return nil, errors.New("redeem 明文里没有 machine_id")
	}
	return t, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// mid 是上号器用的设备编号：QQ 号 + 固定常量 7049668087。
const selfMidSuffix = "7049668087"

func (t *selfTicket) Mid() string {
	return t.QQ + selfMidSuffix
}

// ClientDevName = "MS-" + mid 前 12 位。
func (t *selfTicket) ClientDevName() string {
	mid := t.Mid()
	if len(mid) > 12 {
		mid = mid[:12]
	}
	return "MS-" + mid
}

func selfVarint(n int) []byte {
	var out []byte
	for n >= 0x80 {
		out = append(out, byte(n)|0x80)
		n >>= 7
	}
	return append(out, byte(n))
}

// buildMidList 按 protobuf 结构生成 mid_list：
//
//	field1(message) = { field1: 0, field2: "<machine_id 小写>" }
//	field1(bytes)   = 设备 blob（沿用模板，本机固定）
func buildMidList(machineID string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(selfMidListTemplate)
	if err != nil {
		return "", fmt.Errorf("mid_list 模板损坏: %w", err)
	}
	// 解析第一段：0A <len> <payload>
	if len(raw) < 2 || raw[0] != 0x0A {
		return "", errors.New("mid_list 模板格式异常")
	}
	n1 := int(raw[1])
	if 2+n1 > len(raw) {
		return "", errors.New("mid_list 模板长度越界")
	}
	rest := raw[2+n1:] // 第二段起（含 0A len blob）

	mid := strings.ToLower(strings.TrimSpace(machineID))
	inner := []byte{0x08, 0x00, 0x12, byte(len(mid))}
	inner = append(inner, []byte(mid)...)
	item1 := append([]byte{0x0A, byte(len(inner))}, inner...)

	out := append(item1, rest...)
	return base64.StdEncoding.EncodeToString(out), nil
}
