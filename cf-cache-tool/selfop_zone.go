package main

// 大区 → SSN 与 Rail 首包。
//
// 证据：CF上号器_v1.17.exe 的 .rdata 0x1400DD240 是 42 项
// {u32 SSN, u32 pad, char* 大区名} 表；0x1400FB7C0 之后的打包表
// 按大区名存放 Rail 配置。不同 SSN 的首包不同，不能拿一个模板只换 Uin。
// 这里内置已提取的 318（广东一区）/333（四川一区）两套；其余 SSN 仍按
// 原逻辑回退。

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

//go:embed rail_zone_318.bin
var selfRailZone318 []byte

//go:embed rail_zone_333.bin
var selfRailZone333 []byte

var selfZoneSSNByName = map[string]string{
	"广东一区": "318", "北京一区": "319", "上海一区": "320", "北京二区": "321",
	"辽宁一区": "322", "辽宁二区": "323", "福建一区": "324", "浙江一区": "325",
	"上海二区": "326", "广东二区": "327", "湖北一区": "328", "湖北二区": "329",
	"陕西一区": "330", "重庆一区": "332", "四川一区": "333", "北京三区": "334",
	"北京四区": "335", "辽宁三区": "336", "广东三区": "338", "广东四区": "339",
	"湖南二区": "340", "湖南一区": "341", "南方大区": "342", "北方大区": "343",
	"江苏一区": "344", "河南一区": "345", "山东一区": "346", "安徽一区": "347",
	"云南一区": "348", "浙江二区": "349", "黑龙江区": "350", "吉林一区": "351",
	"江西一区": "352", "广西一区": "353", "山西一区": "354", "河北一区": "355",
	"四川二区": "356", "江苏二区": "357", "山东二区": "358", "河南二区": "359",
	"移动专区": "360", "教育网区": "361",
}

// selfNormalizeZone 把卡密备注里的“广东电信一区/四川电信一区”归一到
// 样本表里的“广东一区/四川一区”。
func selfNormalizeZone(s string) string {
	s = strings.TrimSpace(s)
	for _, w := range []string{"电信", "网通", "联通", "铁通", "移动", "区服"} {
		s = strings.ReplaceAll(s, w, "")
	}
	return strings.TrimSpace(s)
}

// selfResolveSSN 返回当前卡的大区 SSN；未知时保留 version.ini 现有值。
func selfResolveSSN(dir, zone string) string {
	if z := selfNormalizeZone(zone); z != "" {
		if ssn := selfZoneSSNByName[z]; ssn != "" {
			return ssn
		}
	}
	if ssn := selfReadSSN(dir); ssn != "" {
		return ssn
	}
	return ""
}

func selfReadSSN(dir string) string {
	for _, p := range []string{filepath.Join(dir, "version.ini"), filepath.Join(dir, "x64", "version.ini")} {
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(raw), "\n") {
			if !strings.Contains(line, "-SSN:") {
				continue
			}
			part := line[strings.Index(line, "-SSN:")+len("-SSN:"):]
			part = strings.TrimSpace(strings.Fields(part + " ")[0])
			part = strings.Trim(part, "\r\n ;")
			if part != "" {
				return part
			}
		}
	}
	return ""
}

// selfWriteVersionIni 按当前大区写 version.ini；未知 SSN 时不覆盖。
func selfWriteVersionIni(dir, ssn string) {
	if strings.TrimSpace(ssn) == "" {
		return
	}
	body := fmt.Sprintf("[AREA]\r\nUSN=1\r\nSSN=-rotate:0 -SSN:%s -cf:1\r\n", ssn)
	for _, p := range []string{
		filepath.Join(dir, "version.ini"),
		filepath.Join(dir, "x64", "version.ini"),
	} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			continue
		}
		_ = os.WriteFile(p, []byte(body), 0o644)
	}
	selfLog("写 version.ini SSN=%s", ssn)
}

// selfRailFileMatchesZone 判断现有 Rail 文件是否像本区离线器产物。
// 8817 是旧版本用别的区模板补出来的假文件，必须排除。
func selfRailFileMatchesZone(ssn string, n int) bool {
	if z := selfZoneRailTemplate(ssn); len(z) > 0 {
		if n == len(z) {
			return true
		}
		// 333 区离线器文件会带账号尾部，长度可能多 132 字节。
		if ssn == "333" && (n == 8669 || n == 8801 || n == 8671) {
			return true
		}
		return false
	}
	return true
}

// selfZoneRailTemplate 返回该 SSN 的真实首包。
// 优先从离线器 exe 的 .rdata 按大区名提取，这样 42 个大区都覆盖；
// 提取失败才退回内置的 318/333。
func selfZoneRailTemplate(ssn string) []byte {
	if b := selfZoneRailFromLoader(ssn); len(b) > 0 {
		return b
	}
	var src []byte
	switch ssn {
	case "318":
		src = selfRailZone318
	case "333":
		src = selfRailZone333
	default:
		return nil
	}
	out := make([]byte, len(src))
	copy(out, src)
	return out
}

// selfZoneRailFromLoader 解析 CF上号器_v1.17.exe：
//
//	0x1400DD240  42 项 {u32 ssn; u32 pad; char* 大区名}
//	0x1400FB7C0  打包表：1 字节 keyLen + key + u32 valLen + val
//
// val 就是该大区的 Rail 首包。
func selfZoneRailFromLoader(ssn string) []byte {
	p := selfFinderLoaderPath()
	if p == "" {
		return nil
	}
	raw, err := os.ReadFile(p)
	if err != nil || len(raw) < 0x1000 || raw[0] != 'M' || raw[1] != 'Z' {
		return nil
	}
	pe := uint64(binary.LittleEndian.Uint32(raw[0x3c:]))
	if pe+0x100 > uint64(len(raw)) {
		return nil
	}
	// PE32+ ImageBase 在 OptionalHeader+0x18，即 pe+0x30。
	base := uint64(binary.LittleEndian.Uint64(raw[pe+0x30:]))
	if base == 0 {
		base = 0x140000000
	}
	secOff := pe + 0x18 + uint64(binary.LittleEndian.Uint16(raw[pe+0x14:]))
	nsec := int(binary.LittleEndian.Uint16(raw[pe+0x06:]))
	rvaOff := func(rva uint64) uint64 {
		for i := 0; i < nsec; i++ {
			o := secOff + uint64(i)*40
			if o+40 > uint64(len(raw)) {
				return 0
			}
			va := uint64(binary.LittleEndian.Uint32(raw[o+12:]))
			vs := uint64(binary.LittleEndian.Uint32(raw[o+8:]))
			rp := uint64(binary.LittleEndian.Uint32(raw[o+20:]))
			if rva >= va && rva < va+vs {
				return rp + (rva - va)
			}
		}
		return 0
	}
	table := rvaOff(0xDD240)
	if table == 0 || table+672 > uint64(len(raw)) {
		return nil
	}
	name := ""
	for i := 0; i < 42; i++ {
		o := table + uint64(i)*16
		key := binary.LittleEndian.Uint32(raw[o:])
		if fmt.Sprint(key) != ssn {
			continue
		}
		va := binary.LittleEndian.Uint64(raw[o+8:])
		off := va - base
		// 反向找所属节
		for j := 0; j < nsec; j++ {
			so := secOff + uint64(j)*40
			sv := uint64(binary.LittleEndian.Uint32(raw[so+12:]))
			ss := uint64(binary.LittleEndian.Uint32(raw[so+8:]))
			sr := uint64(binary.LittleEndian.Uint32(raw[so+20:]))
			if off >= sv && off < sv+ss {
				fo := sr + (off - sv)
				if fo < uint64(len(raw)) {
					e := fo
					for e < uint64(len(raw)) && raw[e] != 0 {
						e++
					}
					name = string(raw[fo:e])
				}
				break
			}
		}
		break
	}
	if name == "" {
		return nil
	}
	off := rvaOff(0xFB7C0)
	end := rvaOff(0x15926E)
	if off == 0 || end == 0 || end <= off {
		return nil
	}
	for off < end && off+5 < uint64(len(raw)) {
		kl := int(raw[off])
		if kl <= 0 || kl > 96 || off+1+uint64(kl)+4 > uint64(len(raw)) {
			off++
			continue
		}
		key := string(raw[off+1 : off+1+uint64(kl)])
		vl := uint64(binary.LittleEndian.Uint32(raw[off+1+uint64(kl):]))
		vo := off + 1 + uint64(kl) + 4
		if vl == 0 || vl > 1<<20 || vo+vl > uint64(len(raw)) {
			off++
			continue
		}
		if key == name {
			out := make([]byte, vl)
			copy(out, raw[vo:vo+vl])
			return out
		}
		off = vo + vl
	}
	return nil
}

// selfFinderLoaderPath 找离线器 exe。
func selfFinderLoaderPath() string {
	cands := []string{
		filepath.Join(exeDir(), "资源", "上号器", "yuanmeng", "CF上号器_v1.17.exe"),
		filepath.Join(exeDir(), "CF上号器_v1.17.exe"),
		filepath.Join(os.Getenv("APPDATA"), "CG", "data", "CF上号器_v1.17.exe"),
	}
	for _, p := range cands {
		if st, err := os.Stat(p); err == nil && st.Size() > 1<<20 {
			return p
		}
	}
	return ""
}

// selfPatchRailAccount 把模板里的 str2（QQ）改成当前账号，并修正 protobuf 长度。
func selfPatchRailAccount(b []byte, qq, tgpid string) ([]byte, error) {
	if len(b) < 8 || qq == "" {
		return nil, errors.New("invalid rail template")
	}
	out := b
	if old := selfRailTemplateUin(b); old != "" && old != qq {
		var err error
		out, err = selfPatchRailUin(b, old, qq)
		if err != nil {
			return nil, err
		}
	}
	out, err := selfPatchRailStr1(out, tgpid)
	if err != nil {
		return nil, err
	}
	if binary.LittleEndian.Uint32(out[:4])+8 != uint32(len(out)) {
		return nil, errors.New("rail 长度修正失败")
	}
	return out, nil
}

// selfPatchRailStr1 原地替换 "str1":"GUID#TGPID#0#CN#1" 的 GUID/TGPID。
// str1 长度不变，所以不用改 protobuf 长度；长度变了就保留原值。
func selfPatchRailStr1(b []byte, tgpid string) ([]byte, error) {
	if strings.TrimSpace(tgpid) == "" {
		return b, nil
	}
	marker := []byte(`"str1":"`)
	i := bytes.Index(b, marker)
	if i < 0 {
		return b, nil
	}
	start := i + len(marker)
	end := bytes.IndexByte(b[start:], '"')
	if end <= 0 {
		return b, nil
	}
	old := string(b[start : start+end])
	parts := strings.Split(old, "#")
	if len(parts) < 2 {
		return b, nil
	}
	guid := parts[0]
	if g := selfMachineGUID(); g != "" {
		guid = g
	}
	parts[0], parts[1] = guid, strings.TrimSpace(tgpid)
	newVal := strings.Join(parts, "#")
	if newVal == old {
		return b, nil
	}
	body, changed, err := selfProtoReplace(b[8:], []byte(old), []byte(newVal))
	if err != nil || !changed {
		selfLog("rail str1 替换失败 %q -> %q: %v", old, newVal, err)
		return b, err
	}
	out := make([]byte, 0, 8+len(body))
	out = append(out, b[:8]...)
	out = append(out, body...)
	binary.LittleEndian.PutUint32(out[:4], uint32(len(out)-8))
	if len(out) != len(b) {
		selfLog("rail str1 长度 %d->%d", len(b), len(out))
	}
	return out, nil
}

// selfProtoReplace 递归替换 protobuf 长度字段里的子串，并重算每层长度。
// 字符串不是合法 protobuf 时退化成 bytes.ReplaceAll。
func selfProtoReplace(data, old, new []byte) ([]byte, bool, error) {
	var out []byte
	pos := 0
	changed := false
	for pos < len(data) {
		start := pos
		key, n, ok := selfReadUvarint(data[pos:])
		if !ok || key>>3 == 0 {
			return nil, false, errors.New("bad protobuf key")
		}
		pos += n
		keyEnd := pos
		wire := key & 7
		switch wire {
		case 0:
			_, n, ok = selfReadUvarint(data[pos:])
			if !ok {
				return nil, false, errors.New("bad varint")
			}
			pos += n
			out = append(out, data[start:pos]...)
		case 1:
			if pos+8 > len(data) {
				return nil, false, errors.New("bad fixed64")
			}
			pos += 8
			out = append(out, data[start:pos]...)
		case 5:
			if pos+4 > len(data) {
				return nil, false, errors.New("bad fixed32")
			}
			pos += 4
			out = append(out, data[start:pos]...)
		case 2:
			l, m, ok := selfReadUvarint(data[pos:])
			if !ok || l > uint64(len(data)-pos-m) {
				return nil, false, errors.New("bad length")
			}
			pos += m
			content := data[pos : pos+int(l)]
			pos += int(l)
			newContent := content
			if sub, ch, err := selfProtoReplace(content, old, new); err == nil && ch {
				newContent = sub
				changed = true
			} else if bytes.Contains(content, old) {
				newContent = bytes.ReplaceAll(content, old, new)
				changed = true
			}
			out = append(out, data[start:keyEnd]...)
			out = selfAppendUvarint(out, uint64(len(newContent)))
			out = append(out, newContent...)
		default:
			return nil, false, errors.New("unsupported wire type")
		}
	}
	if !changed {
		return data, false, nil
	}
	return out, true, nil
}

func selfMachineGUID() string {
	const (
		hkeyLocalMachine = 0x80000002
		keyRead          = 0x20019
		keyWow64_64      = 0x0100
	)
	advapi := syscall.NewLazyDLL("advapi32.dll")
	open := advapi.NewProc("RegOpenKeyExW")
	query := advapi.NewProc("RegQueryValueExW")
	closeKey := advapi.NewProc("RegCloseKey")
	sub, _ := syscall.UTF16PtrFromString(`SOFTWARE\Microsoft\Cryptography`)
	var h syscall.Handle
	if r, _, _ := open.Call(hkeyLocalMachine, uintptr(unsafe.Pointer(sub)), 0, keyRead|keyWow64_64, uintptr(unsafe.Pointer(&h))); r != 0 {
		return ""
	}
	defer closeKey.Call(uintptr(h))
	name, _ := syscall.UTF16PtrFromString("MachineGuid")
	buf := make([]uint16, 128)
	size := uint32(len(buf) * 2)
	if r, _, _ := query.Call(uintptr(h), uintptr(unsafe.Pointer(name)), 0, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size))); r != 0 {
		return ""
	}
	return strings.TrimSpace(syscall.UTF16ToString(buf))
}
