package main

// 复刻样本 sub_14000FC70 / sub_14000F440：
// 启动后持续扫描所有进程，找到加载了 rail_api64.dll 的进程，把
// RailNeedRestartAppForCheckingEnvironment 的入口改成 trampoline：
//
//	sub rsp, 28h
//	mov rax, stub
//	call rax                 ; 执行原函数前 15 字节 + 跳回原函数 +15
//	xor eax, eax             ; 丢弃原返回值，强制返回 0
//	add rsp, 28h
//	ret
//
// 样本每 15ms 扫一次，持续 3000 次（约 45 秒）；这里保持一致。

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const selfRailPatchExport = "RailNeedRestartAppForCheckingEnvironment"

// selfPatchRailLoop 对应样本 sub_14000FC70。
func selfPatchRailLoop(railDLL string, timeout time.Duration) {
	rva, err := selfPEExportRVA(railDLL, selfRailPatchExport)
	if err != nil {
		selfLog("rail patch: 解析 %s 失败: %v", railDLL, err)
		return
	}
	selfLog("rail patch: %s %s RVA=0x%X", filepath.Base(railDLL), selfRailPatchExport, rva)

	deadline := time.Now().Add(timeout)
	patched := 0
	for time.Now().Before(deadline) {
		for _, pid := range selfEnumProcessIDs() {
			if selfPatchRailInProcess(pid, rva) {
				patched++
				selfLog("rail patch: pid=%d patched", pid)
			}
		}
		time.Sleep(15 * time.Millisecond)
	}
	selfLog("rail patch: 结束，patched=%d", patched)
}

// selfPatchRailInProcess 对应样本 sub_14000F440。
func selfPatchRailInProcess(pid uint32, exportRVA uint32) bool {
	k32 := syscall.NewLazyDLL("kernel32.dll")
	psapi := syscall.NewLazyDLL("psapi.dll")
	openProcess := k32.NewProc("OpenProcess")
	closeHandle := k32.NewProc("CloseHandle")
	readProcMem := k32.NewProc("ReadProcessMemory")
	writeProcMem := k32.NewProc("WriteProcessMemory")
	virtualAllocEx := k32.NewProc("VirtualAllocEx")
	virtualProtectEx := k32.NewProc("VirtualProtectEx")
	enumModules := psapi.NewProc("EnumProcessModules")
	getModuleName := psapi.NewProc("GetModuleFileNameExW")

	const (
		processQueryInformation = 0x0400
		processVMOperation      = 0x0008
		processVMWrite          = 0x0020
		processVMRead           = 0x0010
		processAccess           = processQueryInformation | processVMOperation | processVMWrite | processVMRead
		memCommitReserve        = 0x3000
		pageExecuteReadWrite    = 0x40
	)
	hProc, _, _ := openProcess.Call(uintptr(pid), 0, processAccess)
	if hProc == 0 {
		return false
	}
	defer closeHandle.Call(hProc)

	mods := make([]uintptr, 2048)
	var needed uint32
	ok, _, _ := enumModules.Call(hProc, uintptr(unsafe.Pointer(&mods[0])),
		uintptr(len(mods)*int(unsafe.Sizeof(mods[0]))), uintptr(unsafe.Pointer(&needed)))
	if ok == 0 || needed == 0 {
		return false
	}
	count := int(needed) / int(unsafe.Sizeof(mods[0]))
	if count > len(mods) {
		count = len(mods)
	}

	buf := make([]uint16, 512)
	for i := 0; i < count; i++ {
		n, _, _ := getModuleName.Call(hProc, mods[i], uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if n == 0 {
			continue
		}
		name := syscall.UTF16ToString(buf)
		if !strings.EqualFold(filepath.Base(name), "rail_api64.dll") {
			continue
		}
		target := mods[i] + uintptr(exportRVA)
		return selfPatchRailTarget(hProc, target,
			readProcMem, writeProcMem, virtualAllocEx, virtualProtectEx,
			memCommitReserve, pageExecuteReadWrite)
	}
	return false
}

func selfPatchRailTarget(hProc, target uintptr,
	readProcMem, writeProcMem, virtualAllocEx, virtualProtectEx *syscall.LazyProc,
	memCommitReserve, pageExecuteReadWrite uintptr) bool {

	var original [16]byte
	var n uintptr
	ok, _, _ := readProcMem.Call(hProc, target, uintptr(unsafe.Pointer(&original[0])), 16, uintptr(unsafe.Pointer(&n)))
	if ok == 0 || n < 15 {
		return false
	}
	// 已经是我们或样本写过的跳转 / 直接返回，不重复 patch。
	if original[0] == 0x48 && original[1] == 0xB8 {
		return true
	}
	if original[0] == 0x31 && original[1] == 0xC0 && original[2] == 0xC3 {
		return true
	}

	stubAddr, _, _ := virtualAllocEx.Call(hProc, 0, 0x40, memCommitReserve, pageExecuteReadWrite)
	if stubAddr == 0 {
		return false
	}

	stub := make([]byte, 0x40)
	copy(stub[0:15], original[:15])
	binary.LittleEndian.PutUint16(stub[15:], 0xB848) // mov rax, target+15
	binary.LittleEndian.PutUint64(stub[17:], uint64(target+15))
	stub[25], stub[26] = 0xFF, 0xE0                                 // jmp rax
	stub[32], stub[33], stub[34], stub[35] = 0x48, 0x83, 0xEC, 0x28 // sub rsp,28h
	binary.LittleEndian.PutUint16(stub[36:], 0xB848)                // mov rax, stub
	binary.LittleEndian.PutUint64(stub[38:], uint64(stubAddr))
	stub[46], stub[47] = 0xFF, 0xD0                                 // call rax
	stub[48], stub[49] = 0x31, 0xC0                                 // xor eax,eax
	stub[50], stub[51], stub[52], stub[53] = 0x48, 0x83, 0xC4, 0x28 // add rsp,28h
	stub[54] = 0xC3

	var written uintptr
	if ok, _, _ := writeProcMem.Call(hProc, stubAddr,
		uintptr(unsafe.Pointer(&stub[0])), uintptr(len(stub)),
		uintptr(unsafe.Pointer(&written))); ok == 0 || written != uintptr(len(stub)) {
		return false
	}

	patch := make([]byte, 15)
	patch[0], patch[1] = 0x48, 0xB8 // mov rax, stub+32
	binary.LittleEndian.PutUint64(patch[2:], uint64(stubAddr+32))
	patch[10], patch[11] = 0xFF, 0xE0 // jmp rax
	for i := 12; i < 15; i++ {
		patch[i] = 0x90
	}

	var oldProtect uint32
	if ok, _, _ := virtualProtectEx.Call(hProc, target, 16, pageExecuteReadWrite,
		uintptr(unsafe.Pointer(&oldProtect))); ok == 0 {
		return false
	}
	ok, _, _ = writeProcMem.Call(hProc, target, uintptr(unsafe.Pointer(&patch[0])),
		uintptr(len(patch)), uintptr(unsafe.Pointer(&written)))
	var restored uint32
	virtualProtectEx.Call(hProc, target, 16, uintptr(oldProtect), uintptr(unsafe.Pointer(&restored)))
	return ok != 0 && written == uintptr(len(patch))
}

// selfEnumProcessIDs 对应样本 sub_14000A130，使用 psapi!EnumProcesses。
func selfEnumProcessIDs() []uint32 {
	psapi := syscall.NewLazyDLL("psapi.dll")
	enumProcesses := psapi.NewProc("EnumProcesses")
	buf := make([]uint32, 4096)
	var needed uint32
	ok, _, _ := enumProcesses.Call(uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)*4), uintptr(unsafe.Pointer(&needed)))
	if ok == 0 {
		return nil
	}
	n := int(needed) / 4
	if n > len(buf) {
		n = len(buf)
	}
	return buf[:n]
}

// selfPEExportRVA 从磁盘 PE 的导出表里取指定导出函数的 RVA。
func selfPEExportRVA(path, want string) (uint32, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(data) < 0x40 || data[0] != 'M' || data[1] != 'Z' {
		return 0, errors.New("不是 PE 文件")
	}
	nt := int(selfU32(data, 0x3C))
	if nt+24 > len(data) || string(data[nt:nt+4]) != "PE\x00\x00" {
		return 0, errors.New("PE 头损坏")
	}
	coff := nt + 4
	nsec := int(selfU16(data, coff+2))
	optSz := int(selfU16(data, coff+16))
	opt := coff + 20
	magic := selfU16(data, opt)
	dd := 0
	switch magic {
	case 0x10b:
		dd = opt + 96
	case 0x20b:
		dd = opt + 112
	default:
		return 0, fmt.Errorf("未知可选头 0x%X", magic)
	}
	type sec struct {
		va, raw, size uint32
	}
	var secs []sec
	so := opt + optSz
	for i := 0; i < nsec; i++ {
		o := so + i*40
		if o+40 > len(data) {
			break
		}
		secs = append(secs, sec{
			va:   selfU32(data, o+12),
			raw:  selfU32(data, o+20),
			size: selfU32(data, o+16),
		})
	}
	rva2off := func(rva uint32) (int, bool) {
		for _, s := range secs {
			if rva >= s.va && rva < s.va+s.size {
				return int(s.raw + (rva - s.va)), true
			}
		}
		return 0, false
	}
	expRVA := selfU32(data, dd)
	if expRVA == 0 {
		return 0, errors.New("没有导出表")
	}
	exp, ok := rva2off(expRVA)
	if !ok || exp+40 > len(data) {
		return 0, errors.New("导出表越界")
	}
	nNames := selfU32(data, exp+24)
	funcs := selfU32(data, exp+28)
	names := selfU32(data, exp+32)
	ords := selfU32(data, exp+36)
	fo, ok1 := rva2off(funcs)
	no, ok2 := rva2off(names)
	oo, ok3 := rva2off(ords)
	if !ok1 || !ok2 || !ok3 {
		return 0, errors.New("导出目录越界")
	}
	for i := uint32(0); i < nNames; i++ {
		nameRVA := selfU32(data, no+int(i*4))
		off, ok := rva2off(nameRVA)
		if !ok {
			continue
		}
		end := off
		for end < len(data) && data[end] != 0 {
			end++
		}
		if string(data[off:end]) != want {
			continue
		}
		ord := selfU16(data, oo+int(i*2))
		fnOff := fo + int(ord)*4
		if fnOff+4 > len(data) {
			return 0, errors.New("导出函数越界")
		}
		return selfU32(data, fnOff), nil
	}
	return 0, fmt.Errorf("找不到导出 %s", want)
}

func selfU16(b []byte, off int) uint16 { return binary.LittleEndian.Uint16(b[off:]) }
func selfU32(b []byte, off int) uint32 { return binary.LittleEndian.Uint32(b[off:]) }
