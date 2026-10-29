//go:build ignore

// 把参考 exe 的图标 / 版本信息 / 数字签名 / 文件时间复制到目标 exe。
//
//	go run stamp_exe.go -from exe/zuhaoclient.exe -to ../dist/CG租号.exe
//	go run stamp_exe.go -from 源.exe -to 目标.exe -icon -version   // 只写资源（UPX 前）
//	go run stamp_exe.go -from 源.exe -to 目标.exe -sig -time       // 只补签名/时间（UPX 后）
//	go run stamp_exe.go -dump 源.exe                              // 看参考文件里有什么
//
// 图标与版本信息用 Win32 UpdateResource 整块替换：先删掉目标里同类型资源，再把参考 exe 的
// 原始资源字节写进去，不重新编码，字段完全一致。
// 数字签名把参考 exe 的 PE 证书表整块搬到目标文件尾部并改写 DataDirectory[4]；签名里没有
// 目标文件的哈希，验证结果为“哈希不匹配”，但证书主体、序列号、签名时间与参考文件一致。
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	rtIcon      = 3
	rtMenu      = 4
	rtDialog    = 5
	rtString    = 6
	rtRCData    = 10
	rtGroupIcon = 14
	rtVersion   = 16
	rtManifest  = 24

	// 窗口类图标的资源 ID（walk 注册窗口类时硬编码 7）
	windowIconID = 7
)

var typeNames = map[uint32]string{
	rtIcon: "RT_ICON", rtMenu: "RT_MENU", rtDialog: "RT_DIALOG", rtString: "RT_STRING",
	rtRCData: "RT_RCDATA", rtGroupIcon: "RT_GROUP_ICON", rtVersion: "RT_VERSION", rtManifest: "RT_MANIFEST",
}

// ---------------------------------------------------------------------------
// 资源标识
// ---------------------------------------------------------------------------

type resName struct {
	isStr bool
	id    uint32
	str   string
}

func (n resName) String() string {
	if n.isStr {
		return n.str
	}
	return fmt.Sprintf("%d", n.id)
}

type resource struct {
	Type, Name     resName
	Lang, CodePage uint32
	Data           []byte
}

// ---------------------------------------------------------------------------
// PE 解析
// ---------------------------------------------------------------------------

type section struct {
	vaddr, span, rawptr, rawsize uint32
}

type peFile struct {
	path      string
	data      []byte
	ntOff     int
	ddOff     int
	csumOff   int
	headersSz uint32
	machine   uint16
	sections  []section
}

func loadPE(path string) (*peFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) < 0x40 || data[0] != 'M' || data[1] != 'Z' {
		return nil, fmt.Errorf("%s: 不是 PE 文件", path)
	}
	p := &peFile{path: path, data: data}
	p.ntOff = int(u32(data, 0x3C))
	if p.ntOff+24 > len(data) || string(data[p.ntOff:p.ntOff+4]) != "PE\x00\x00" {
		return nil, fmt.Errorf("%s: PE 头损坏", path)
	}
	coff := p.ntOff + 4
	p.machine = u16(data, coff)
	nsec := int(u16(data, coff+2))
	optSz := int(u16(data, coff+16))
	opt := coff + 20
	magic := u16(data, opt)
	switch magic {
	case 0x10b:
		p.ddOff = opt + 96
	case 0x20b:
		p.ddOff = opt + 112
	default:
		return nil, fmt.Errorf("%s: 未知可选头格式 0x%X", path, magic)
	}
	p.csumOff = opt + 64
	p.headersSz = u32(data, opt+60)
	sect := opt + optSz
	for i := 0; i < nsec; i++ {
		o := sect + i*40
		if o+40 > len(data) {
			break
		}
		vsize := u32(data, o+8)
		vaddr := u32(data, o+12)
		rawsize := u32(data, o+16)
		rawptr := u32(data, o+20)
		span := vsize
		if rawsize > span {
			span = rawsize
		}
		p.sections = append(p.sections, section{vaddr, span, rawptr, rawsize})
	}
	return p, nil
}

func u16(b []byte, off int) uint16 { return binary.LittleEndian.Uint16(b[off:]) }
func u32(b []byte, off int) uint32 { return binary.LittleEndian.Uint32(b[off:]) }

func (p *peFile) rva2off(rva uint32) (uint32, error) {
	for _, s := range p.sections {
		if rva >= s.vaddr && rva < s.vaddr+s.span {
			d := rva - s.vaddr
			if d >= s.rawsize {
				return 0, fmt.Errorf("%s: RVA 0x%X 落在未初始化区", p.path, rva)
			}
			return s.rawptr + d, nil
		}
	}
	if rva < p.headersSz {
		return rva, nil
	}
	return 0, fmt.Errorf("%s: RVA 0x%X 不在任何节内", p.path, rva)
}

func (p *peFile) nameAt(base, field uint32) (resName, error) {
	if field&0x80000000 == 0 {
		return resName{id: field}, nil
	}
	off := base + (field & 0x7fffffff)
	if int(off)+2 > len(p.data) {
		return resName{}, errors.New("资源名字越界")
	}
	n := int(u16(p.data, int(off)))
	start := int(off) + 2
	if start+n*2 > len(p.data) {
		return resName{}, errors.New("资源名字越界")
	}
	raw := p.data[start : start+n*2]
	runes := make([]uint16, n)
	for i := 0; i < n; i++ {
		runes[i] = binary.LittleEndian.Uint16(raw[i*2:])
	}
	return resName{isStr: true, str: syscall.UTF16ToString(runes)}, nil
}

func (p *peFile) resources() ([]resource, error) {
	rsrcRVA := u32(p.data, p.ddOff+2*8)
	if rsrcRVA == 0 {
		return nil, nil
	}
	base, err := p.rva2off(rsrcRVA)
	if err != nil {
		return nil, err
	}
	var out []resource
	var walk func(off uint32, path []resName) error
	walk = func(off uint32, path []resName) error {
		if int(base+off+16) > len(p.data) {
			return errors.New("资源目录越界")
		}
		named := u16(p.data, int(base+off+12))
		ids := u16(p.data, int(base+off+14))
		for i := 0; i < int(named+ids); i++ {
			ent := base + off + 16 + uint32(i)*8
			if int(ent)+8 > len(p.data) {
				return errors.New("资源目录项越界")
			}
			nf := u32(p.data, int(ent))
			tf := u32(p.data, int(ent+4))
			key, err := p.nameAt(base, nf)
			if err != nil {
				return err
			}
			if tf&0x80000000 != 0 {
				if err := walk(tf&0x7fffffff, append(path, key)); err != nil {
					return err
				}
				continue
			}
			if len(path) != 2 {
				continue
			}
			de := base + (tf & 0x7fffffff)
			if int(de)+16 > len(p.data) {
				return errors.New("资源数据项越界")
			}
			dataRVA := u32(p.data, int(de))
			size := u32(p.data, int(de+4))
			cp := u32(p.data, int(de+8))
			doff, err := p.rva2off(dataRVA)
			if err != nil {
				return err
			}
			if int(doff+size) > len(p.data) {
				return errors.New("资源数据越界")
			}
			out = append(out, resource{
				Type: path[0], Name: path[1], Lang: key.id, CodePage: cp,
				Data: p.data[doff : doff+size],
			})
		}
		return nil
	}
	if err := walk(0, nil); err != nil {
		return nil, err
	}
	return out, nil
}

func (p *peFile) cert() []byte {
	off := u32(p.data, p.ddOff+4*8)
	size := u32(p.data, p.ddOff+4*8+4)
	if off == 0 || size == 0 || int(off+size) > len(p.data) {
		return nil
	}
	return p.data[off : off+size]
}

// ---------------------------------------------------------------------------
// 校验和
// ---------------------------------------------------------------------------

func peChecksum(data []byte, skip int) uint32 {
	var total uint64
	n := len(data)
	i := 0
	for ; i+1 < n; i += 2 {
		if i >= skip && i < skip+4 {
			continue
		}
		total += uint64(data[i]) | uint64(data[i+1])<<8
		total = (total & 0xFFFF) + (total >> 16)
	}
	if i < n {
		total += uint64(data[i])
		total = (total & 0xFFFF) + (total >> 16)
	}
	total = (total & 0xFFFF) + (total >> 16)
	return uint32(total) + uint32(n)
}

func writeChecksum(path string) (uint32, error) {
	p, err := loadPE(path)
	if err != nil {
		return 0, err
	}
	d := append([]byte(nil), p.data...)
	binary.LittleEndian.PutUint32(d[p.csumOff:], 0)
	sum := peChecksum(d, p.csumOff)
	binary.LittleEndian.PutUint32(d[p.csumOff:], sum)
	if err := os.WriteFile(path, d, 0644); err != nil {
		return 0, err
	}
	return sum, nil
}

// ---------------------------------------------------------------------------
// 资源写入（kernel32!UpdateResource）
// ---------------------------------------------------------------------------

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procBeginUpdateResourceW = kernel32.NewProc("BeginUpdateResourceW")
	procUpdateResourceW      = kernel32.NewProc("UpdateResourceW")
	procEndUpdateResourceW   = kernel32.NewProc("EndUpdateResourceW")
)

func resPtr(n resName) (uintptr, *uint16, error) {
	if !n.isStr {
		return uintptr(n.id), nil, nil
	}
	p, err := syscall.UTF16PtrFromString(n.str)
	if err != nil {
		return 0, nil, err
	}
	return uintptr(unsafe.Pointer(p)), p, nil
}

func updateResources(path string, toDelete []resource, entries []resource) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	h, _, e := procBeginUpdateResourceW.Call(uintptr(unsafe.Pointer(pathPtr)), 0)
	if h == 0 {
		return fmt.Errorf("BeginUpdateResource 失败: %v", e)
	}
	commit := false
	var keep []*uint16
	defer func() {
		discard := uintptr(0)
		if !commit {
			discard = 1
		}
		procEndUpdateResourceW.Call(h, discard)
	}()

	// 逐个删掉目标里同类型的资源（lpData=NULL + cbData=0）
	for _, r := range toDelete {
		tp, tpPtr, err := resPtr(r.Type)
		if err != nil {
			return err
		}
		np, npPtr, err := resPtr(r.Name)
		if err != nil {
			return err
		}
		if tpPtr != nil {
			keep = append(keep, tpPtr)
		}
		if npPtr != nil {
			keep = append(keep, npPtr)
		}
		if ok, _, e := procUpdateResourceW.Call(h, tp, np, uintptr(r.Lang), 0, 0); ok == 0 {
			return fmt.Errorf("删除资源 %v/%v[%d] 失败: %v", r.Type, r.Name, r.Lang, e)
		}
	}

	for _, r := range entries {
		tp, tpPtr, err := resPtr(r.Type)
		if err != nil {
			return err
		}
		np, npPtr, err := resPtr(r.Name)
		if err != nil {
			return err
		}
		if tpPtr != nil {
			keep = append(keep, tpPtr)
		}
		if npPtr != nil {
			keep = append(keep, npPtr)
		}
		var dataPtr uintptr
		if len(r.Data) > 0 {
			dataPtr = uintptr(unsafe.Pointer(&r.Data[0]))
		}
		ok, _, e := procUpdateResourceW.Call(h, tp, np, uintptr(r.Lang), dataPtr, uintptr(len(r.Data)))
		if ok == 0 {
			return fmt.Errorf("写入资源 %v/%v 失败: %v", r.Type, r.Name, e)
		}
	}
	if ok, _, e := procEndUpdateResourceW.Call(h, 0); ok == 0 {
		return fmt.Errorf("EndUpdateResource 失败: %v", e)
	}
	commit = true
	return nil
}

// ---------------------------------------------------------------------------
// 版本信息（读取用 version.dll，和资源管理器显示一致）
// ---------------------------------------------------------------------------

type verInfo struct {
	fileVersion    string
	productVersion string
	strings        map[string]string
	tables         []string
}

var (
	versionDLL                  = syscall.NewLazyDLL("version.dll")
	procGetFileVersionInfoSizeW = versionDLL.NewProc("GetFileVersionInfoSizeW")
	procGetFileVersionInfoW     = versionDLL.NewProc("GetFileVersionInfoW")
	procVerQueryValueW          = versionDLL.NewProc("VerQueryValueW")
)

var verKeys = []string{
	"CompanyName", "FileDescription", "FileVersion", "InternalName", "LegalCopyright",
	"LegalTrademarks", "OriginalFilename", "ProductName", "ProductVersion", "Comments",
}

func versionStrings(path string) verInfo {
	vi := verInfo{strings: map[string]string{}}
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return vi
	}
	var dummy uint32
	size, _, _ := procGetFileVersionInfoSizeW.Call(uintptr(unsafe.Pointer(pathPtr)), uintptr(unsafe.Pointer(&dummy)))
	if size == 0 {
		return vi
	}
	buf := make([]byte, size)
	ok, _, _ := procGetFileVersionInfoW.Call(uintptr(unsafe.Pointer(pathPtr)), 0, uintptr(size), uintptr(unsafe.Pointer(&buf[0])))
	if ok == 0 {
		return vi
	}

	query := func(sub string) (unsafe.Pointer, uint32, bool) {
		subPtr, e := syscall.UTF16PtrFromString(sub)
		if e != nil {
			return nil, 0, false
		}
		var ptr unsafe.Pointer
		var n uint32
		r, _, _ := procVerQueryValueW.Call(uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(subPtr)), uintptr(unsafe.Pointer(&ptr)), uintptr(unsafe.Pointer(&n)))
		return ptr, n, r != 0
	}
	queryString := func(sub string) string {
		ptr, n, good := query(sub)
		if !good || ptr == nil || n == 0 {
			return ""
		}
		return syscall.UTF16ToString(unsafe.Slice((*uint16)(ptr), n))
	}

	if ptr, n, good := query("\\"); good && n >= 52 {
		ffi := unsafe.Slice((*uint32)(ptr), 13)
		vi.fileVersion = fmt.Sprintf("%d.%d.%d.%d", ffi[2]>>16, ffi[2]&0xFFFF, ffi[3]>>16, ffi[3]&0xFFFF)
		vi.productVersion = fmt.Sprintf("%d.%d.%d.%d", ffi[4]>>16, ffi[4]&0xFFFF, ffi[5]>>16, ffi[5]&0xFFFF)
	}

	type trans struct{ lang, cp uint16 }
	var list []trans
	if ptr, n, good := query("\\VarFileInfo\\Translation"); good && ptr != nil {
		// 该块返回的字节数，两个 WORD 为一组
		words := unsafe.Slice((*uint16)(ptr), n/2)
		for i := 0; i+1 < len(words); i += 2 {
			list = append(list, trans{words[i], words[i+1]})
		}
	}
	if len(list) == 0 {
		list = []trans{{0x0409, 0x04B0}}
	}
	for _, t := range list {
		vi.tables = append(vi.tables, fmt.Sprintf("%04x%04x", t.lang, t.cp))
		for _, k := range verKeys {
			if _, done := vi.strings[k]; done {
				continue
			}
			if v := queryString(fmt.Sprintf("\\StringFileInfo\\%04x%04x\\%s", t.lang, t.cp, k)); v != "" {
				vi.strings[k] = v
			}
		}
	}
	return vi
}

// ---------------------------------------------------------------------------
// 数字签名
// ---------------------------------------------------------------------------

type certEntry struct {
	length   uint32
	revision uint16
	kind     uint16
}

func certEntries(b []byte) []certEntry {
	var out []certEntry
	for off := 0; off+8 <= len(b); {
		length := u32(b, off)
		if length < 8 || int(off)+int(length) > len(b) {
			break
		}
		out = append(out, certEntry{length, u16(b, off+4), u16(b, off+6)})
		off += int((length + 7) &^ 7)
	}
	return out
}

func stripCert(path string) error {
	p, err := loadPE(path)
	if err != nil {
		return err
	}
	d := append([]byte(nil), p.data...)
	off := u32(d, p.ddOff+4*8)
	size := u32(d, p.ddOff+4*8+4)
	if off != 0 && size != 0 && int(off+size) <= len(d) {
		d = d[:off]
	} else if off != 0 && int(off) <= len(d) {
		d = d[:off]
	}
	binary.LittleEndian.PutUint32(d[p.ddOff+4*8:], 0)
	binary.LittleEndian.PutUint32(d[p.ddOff+4*8+4:], 0)
	binary.LittleEndian.PutUint32(d[p.csumOff:], 0)
	return os.WriteFile(path, d, 0644)
}

func appendCert(path string, cert []byte) (uint32, error) {
	p, err := loadPE(path)
	if err != nil {
		return 0, err
	}
	d := append([]byte(nil), p.data...)
	off := u32(d, p.ddOff+4*8)
	if off != 0 && int(off) <= len(d) {
		d = d[:off]
	}
	binary.LittleEndian.PutUint32(d[p.ddOff+4*8:], 0)
	binary.LittleEndian.PutUint32(d[p.ddOff+4*8+4:], 0)
	for len(d)%8 != 0 {
		d = append(d, 0)
	}
	certOff := uint32(len(d))
	d = append(d, cert...)
	for len(d)%8 != 0 {
		d = append(d, 0)
	}
	binary.LittleEndian.PutUint32(d[p.ddOff+4*8:], certOff)
	binary.LittleEndian.PutUint32(d[p.ddOff+4*8+4:], uint32(len(cert)))
	return certOff, os.WriteFile(path, d, 0644)
}

// ---------------------------------------------------------------------------
// 文件时间
// ---------------------------------------------------------------------------

func copyTimes(src, dst string) (time.Time, error) {
	st, err := os.Stat(src)
	if err != nil {
		return time.Time{}, err
	}
	if err := os.Chtimes(dst, st.ModTime(), st.ModTime()); err != nil {
		return time.Time{}, err
	}

	pathPtr, err := syscall.UTF16PtrFromString(dst)
	if err != nil {
		return time.Time{}, err
	}
	const fileWriteAttributes = 0x0100
	h, err := syscall.CreateFile(pathPtr, fileWriteAttributes, 0, nil,
		syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return st.ModTime(), nil
	}
	defer syscall.CloseHandle(h)

	ft := syscall.NsecToFiletime(st.ModTime().UnixNano())
	creation := ft
	access := ft
	write := ft
	if info, e := st.Sys().(*syscall.Win32FileAttributeData); e {
		creation = info.CreationTime
		access = info.LastAccessTime
	}
	_ = syscall.SetFileTime(h, &creation, &access, &write)
	return st.ModTime(), nil
}

// ---------------------------------------------------------------------------
// 展示
// ---------------------------------------------------------------------------

func listResources(path string) error {
	p, err := loadPE(path)
	if err != nil {
		return err
	}
	res, err := p.resources()
	if err != nil {
		return err
	}
	fmt.Printf("文件: %s\n", path)
	for _, r := range res {
		name := typeNames[r.Type.id]
		if r.Type.isStr {
			name = r.Type.str
		}
		fmt.Printf("  %-16s %-12s lang=0x%04X %8d 字节\n", name, r.Name.String(), r.Lang, len(r.Data))
	}
	return nil
}

func human(n int) string {
	switch {
	case n < 1024:
		return fmt.Sprintf("%d B", n)
	case n < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(n)/1024)
	default:
		return fmt.Sprintf("%.1f MB", float64(n)/1024/1024)
	}
}

func describe(path string) error {
	p, err := loadPE(path)
	if err != nil {
		return err
	}
	res, err := p.resources()
	if err != nil {
		return err
	}
	counts := map[string]int{}
	order := []string{}
	for _, r := range res {
		name := typeNames[r.Type.id]
		if r.Type.isStr {
			name = r.Type.str
		}
		name = strings.Map(func(ch rune) rune {
			if ch < 0x20 || ch == 0x7F {
				return -1
			}
			return ch
		}, name)
		if name == "" {
			name = "(未命名)"
		}
		if _, ok := counts[name]; !ok {
			order = append(order, name)
		}
		counts[name]++
	}
	sort.Strings(order)

	arch := fmt.Sprintf("0x%X", p.machine)
	switch p.machine {
	case 0x8664:
		arch = "x64"
	case 0x14c:
		arch = "x86"
	}
	fmt.Printf("文件: %s\n", path)
	fmt.Printf("  大小: %s    架构: %s\n", human(len(p.data)), arch)
	parts := make([]string, 0, len(order))
	for _, k := range order {
		parts = append(parts, fmt.Sprintf("%s x%d", k, counts[k]))
	}
	fmt.Printf("  资源: %s\n", strings.Join(parts, "，"))

	vi := versionStrings(path)
	if vi.fileVersion != "" {
		fmt.Printf("  版本信息: 文件 %s / 产品 %s  (语言表 %s)\n",
			vi.fileVersion, vi.productVersion, strings.Join(vi.tables, ","))
		keys := make([]string, 0, len(vi.strings))
		for k := range vi.strings {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("    %-16s: %s\n", k, vi.strings[k])
		}
	} else {
		fmt.Println("  版本信息: 无")
	}
	for _, r := range res {
		if r.Type.id == rtGroupIcon && !r.Type.isStr && len(r.Data) >= 6 {
			fmt.Printf("  图标组: %d 个图像\n", u16(r.Data, 4))
			break
		}
	}

	cert := p.cert()
	if len(cert) == 0 {
		fmt.Println("  数字签名: 无")
	} else {
		ents := certEntries(cert)
		fmt.Printf("  数字签名: %s，证书项 %d\n", human(len(cert)), len(ents))
		for _, e := range ents {
			fmt.Printf("    证书项: %d 字节 rev=0x%04X type=%d\n", e.length, e.revision, e.kind)
		}
	}

	st, _ := os.Stat(path)
	creation := st.ModTime()
	if info, ok := st.Sys().(*syscall.Win32FileAttributeData); ok {
		creation = time.Unix(0, info.CreationTime.Nanoseconds())
	}
	fmt.Printf("  时间: 修改 %s / 创建 %s\n",
		st.ModTime().Format("2006-01-02 15:04:05"), creation.Format("2006-01-02 15:04:05"))
	return nil
}

// ---------------------------------------------------------------------------
// 校验
// ---------------------------------------------------------------------------

func resKey(r resource) string {
	return fmt.Sprintf("%v|%v|%d", r.Type, r.Name, r.Lang)
}

func verify(src, dst string, doIcon, doVersion, doSig bool) bool {
	sp, err1 := loadPE(src)
	dp, err2 := loadPE(dst)
	if err1 != nil || err2 != nil {
		fmt.Printf("    校验失败: %v %v\n", err1, err2)
		return false
	}
	sres, _ := sp.resources()
	dres, _ := dp.resources()
	dmap := map[string][]byte{}
	for _, r := range dres {
		dmap[resKey(r)] = r.Data
	}
	ok := true
	check := func(label string, pred func(resource) bool) {
		total, same := 0, 0
		for _, r := range sres {
			if !pred(r) {
				continue
			}
			total++
			if bytes.Equal(dmap[resKey(r)], r.Data) {
				same++
			}
		}
		status := "OK"
		if total == 0 || same != total {
			status = "不一致"
			ok = false
		}
		fmt.Printf("    %s: %s（%d/%d 一致）\n", label, status, same, total)
	}
	if doIcon {
		check("图标", func(r resource) bool {
			return !r.Type.isStr && (r.Type.id == rtIcon || r.Type.id == rtGroupIcon)
		})
	}
	if doVersion {
		check("版本信息", func(r resource) bool {
			return !r.Type.isStr && r.Type.id == rtVersion
		})
	}
	if doSig {
		sc, dc := sp.cert(), dp.cert()
		good := len(sc) > 0 && bytes.Equal(sc, dc)
		status := "OK"
		if !good {
			status = "不一致"
			ok = false
		}
		fmt.Printf("    数字签名: %s（%d 字节）\n", status, len(dc))
	}
	return ok
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	var (
		src    = flag.String("from", "", "参考 exe")
		dst    = flag.String("to", "", "目标 exe")
		dump   = flag.String("dump", "", "只打印参考文件的信息")
		list   = flag.String("list", "", "列出资源类型/名字/语言")
		doIcon = flag.Bool("icon", false, "复制图标")
		doVer  = flag.Bool("version", false, "复制版本信息")
		doSig  = flag.Bool("sig", false, "复制数字签名")
		doTime = flag.Bool("time", false, "复制文件时间")
	)
	flag.Parse()

	if *dump != "" {
		if err := describe(*dump); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if *list != "" {
		if err := listResources(*list); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if *src == "" || *dst == "" {
		flag.Usage()
		os.Exit(2)
	}
	for _, p := range []string{*src, *dst} {
		if _, err := os.Stat(p); err != nil {
			fmt.Fprintf(os.Stderr, "错误: 找不到 %s\n", p)
			os.Exit(2)
		}
	}
	icon, ver, sig, tm := *doIcon, *doVer, *doSig, *doTime
	if !icon && !ver && !sig && !tm {
		icon, ver, sig, tm = true, true, true, true
	}

	sp, err := loadPE(*src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
	sres, err := sp.resources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("==> 复制到 %s\n", *dst)
	fmt.Printf("    参考: %s\n", *src)

	if icon || ver {
		var entries []resource
		want := func(r resource) bool {
			if r.Type.isStr {
				return false
			}
			if icon && (r.Type.id == rtIcon || r.Type.id == rtGroupIcon) {
				return true
			}
			return ver && r.Type.id == rtVersion
		}
		var groups []resource
		for _, r := range sres {
			if want(r) {
				entries = append(entries, r)
			}
			if icon && !r.Type.isStr && r.Type.id == rtGroupIcon {
				groups = append(groups, r)
			}
		}
		// walk 注册窗口类时只 LoadIcon(hInst, 7)，参考程序的图标组是 133。
		// 补一份同内容的 ID 7 图标组，标题栏 / 任务栏 / Alt-Tab 才会用同一个图标。
		if icon {
			for _, g := range groups {
				if !g.Name.isStr && g.Name.id == windowIconID {
					continue
				}
				dup := g
				dup.Name = resName{id: windowIconID}
				entries = append(entries, dup)
			}
		}
		if len(entries) == 0 {
			fmt.Println("    警告: 参考文件里没有对应资源，跳过")
		} else {
			dp, err := loadPE(*dst)
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				os.Exit(1)
			}
			dres, err := dp.resources()
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				os.Exit(1)
			}
			var deletes []resource
			for _, r := range dres {
				if want(r) {
					deletes = append(deletes, r)
				}
			}
			if err := updateResources(*dst, deletes, entries); err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				os.Exit(1)
			}
			if icon {
				icons, groups := 0, 0
				for _, r := range entries {
					switch r.Type.id {
					case rtIcon:
						icons++
					case rtGroupIcon:
						groups++
					}
				}
				fmt.Printf("    图标: 写入 %d 个图像 / %d 个图标组\n", icons, groups)
			}
			if ver {
				if vi := versionStrings(*dst); vi.fileVersion != "" {
					fmt.Printf("    版本信息: %s  公司=%s  说明=%s\n",
						vi.fileVersion, vi.strings["CompanyName"], vi.strings["FileDescription"])
				}
			}
		}
	}

	if sig {
		cert := sp.cert()
		if len(cert) == 0 {
			fmt.Println("    警告: 参考文件没有数字签名，跳过")
		} else {
			if err := stripCert(*dst); err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				os.Exit(1)
			}
			off, err := appendCert(*dst, cert)
			if err != nil {
				fmt.Fprintf(os.Stderr, "错误: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("    数字签名: 写入 %s（证书表偏移 0x%X）\n", human(len(cert)), off)
		}
	}

	// 校验和要写在设置时间之前，否则最后写文件会把修改时间刷成现在
	sum, err := writeChecksum(*dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("    校验和: 0x%08X\n", sum)

	if tm {
		when, err := copyTimes(*src, *dst)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("    文件时间: %s\n", when.Format("2006-01-02 15:04:05"))
	}

	fmt.Println("==> 校验")
	if !verify(*src, *dst, icon, ver, sig) {
		os.Exit(1)
	}
}
