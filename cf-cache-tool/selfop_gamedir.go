package main

// 自研上号的游戏目录：选一次、校验、记住（config.json 的 game_dir），
// 之后启动与 DLL 处理都以客户自己的游戏目录为准，不再依赖注册表/WeGame。

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/walk"
)

// selfNormalizeGameDir 把用户选中的目录（可能选到 x64、TCLS 或子目录）归一到游戏根目录。
func selfNormalizeGameDir(raw string) (string, error) {
	dir := strings.TrimSpace(strings.Trim(strings.TrimSpace(raw), `"`))
	if dir == "" {
		return "", errors.New("没有选择游戏目录")
	}
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	dir = strings.TrimRight(dir, `\/`)
	if dir == "" { // 例如选了盘符根目录
		return "", fmt.Errorf("目录无效：%s", raw)
	}

	for i := 0; i < 4; i++ {
		if selfHasGameExe(dir) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("目录里没有 x64\\crossfire.exe：%s", raw)
}

func selfHasGameExe(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, selfGameExeRel))
	return err == nil && !st.IsDir()
}

// selfGameRailDLL 返回客户自己游戏目录里的 rail_api64.dll。
func selfGameRailDLL(dir string) string {
	return filepath.Join(dir, "x64", "rail_api64.dll")
}

// selfInspectGameDir 校验并归一目录。
// 返回值：归一后的目录、提示（缺少 rail_api64.dll 之类，可继续）、错误（不可用）。
func selfInspectGameDir(raw string) (string, string, error) {
	dir, err := selfNormalizeGameDir(raw)
	if err != nil {
		return "", "", err
	}
	warn := ""
	if st, err := os.Stat(selfGameRailDLL(dir)); err != nil || st.IsDir() {
		warn = "缺少 x64\\rail_api64.dll（游戏可能不完整，Rail 会初始化失败）"
	}
	if st, err := os.Stat(filepath.Join(dir, selfGameRelDir)); err != nil || !st.IsDir() {
		if warn == "" {
			warn = "没有 TCLS 目录（游戏可能没装完整）"
		}
	}
	return dir, warn, nil
}

// selfRailDLLInfo 校验客户自己游戏目录里的 rail_api64.dll。
// 只读文件确认是 Rail SDK（导出表里有 RailInitialize），不加载执行。
func selfRailDLLInfo(dir string) (string, error) {
	p := selfGameRailDLL(dir)
	st, err := os.Stat(p)
	if err != nil || st.IsDir() {
		return "", fmt.Errorf("找不到 %s", p)
	}
	if st.Size() < 1<<20 {
		return "", fmt.Errorf("rail_api64.dll 大小异常（%d 字节），可能不是真的", st.Size())
	}
	raw, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("读 rail_api64.dll 失败: %w", err)
	}
	if !bytes.Contains(raw, []byte("RailInitialize")) {
		return "", errors.New("rail_api64.dll 里没有 RailInitialize，不是 Rail SDK")
	}
	return fmt.Sprintf("%s（%.1f MB）", p, float64(st.Size())/(1<<20)), nil
}

// selfPickGameDir 弹「选择文件夹」对话框。
func selfPickGameDir(owner walk.Form, initial string) (string, bool) {
	dlg := &walk.FileDialog{Title: "选择穿越火线游戏目录（里面要有 x64\\crossfire.exe）"}
	if initial != "" {
		dlg.InitialDirPath = initial
	}
	ok, err := dlg.ShowBrowseFolder(owner)
	if err != nil || !ok {
		return "", false
	}
	return dlg.FilePath, true
}

// selfPickFolderAtDrives 是 shell 的「浏览文件夹」对话框，根节点固定为「此电脑」，
// 第一次选目录时从盘符开始展开，不会一打开就扎进某个游戏目录。
// 返回：路径、是否确认、对话框是否成功弹出。
func selfPickFolderAtDrives(owner uintptr, title string) (string, bool, bool) {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	ole32 := syscall.NewLazyDLL("ole32.dll")
	oleInit := ole32.NewProc("OleInitialize")
	oleUninit := ole32.NewProc("OleUninitialize")
	coFree := ole32.NewProc("CoTaskMemFree")
	getFolder := shell32.NewProc("SHGetSpecialFolderLocation")
	browse := shell32.NewProc("SHBrowseForFolderW")
	getPath := shell32.NewProc("SHGetPathFromIDListW")

	// BIF_NEWDIALOGSTYLE 需要 OLE 初始化
	res, _, _ := oleInit.Call(0)
	if int32(res) >= 0 {
		defer oleUninit.Call()
	}

	const csidlDrives = 0x0011 // 「我的电脑 / 此电脑」
	var root uintptr
	if hr, _, _ := getFolder.Call(0, csidlDrives, uintptr(unsafe.Pointer(&root))); hr != 0 || root == 0 {
		return "", false, false
	}
	defer coFree.Call(root)

	t, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return "", false, false
	}
	type browseInfoW struct {
		hwndOwner      uintptr
		pidlRoot       uintptr
		pszDisplayName *uint16
		lpszTitle      *uint16
		ulFlags        uint32
		lpfn           uintptr
		lParam         uintptr
		iImage         int32
	}
	const (
		bifReturnOnlyFSDirs = 0x0001
		bifNewDialogStyle   = 0x0040
	)
	disp := make([]uint16, 260)
	bi := browseInfoW{
		hwndOwner:      owner,
		pidlRoot:       root,
		pszDisplayName: &disp[0],
		lpszTitle:      t,
		ulFlags:        bifReturnOnlyFSDirs | bifNewDialogStyle,
	}
	pidl, _, _ := browse.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return "", false, true // 弹出来了，用户取消
	}
	defer coFree.Call(pidl)

	path := make([]uint16, 1024)
	if ok, _, _ := getPath.Call(pidl, uintptr(unsafe.Pointer(&path[0]))); ok == 0 {
		return "", false, true
	}
	p := syscall.UTF16ToString(path)
	return p, p != "", true
}

// selfOpGameDir 按「配置 → 弹框让用户选」的顺序拿一个可用目录，选完记住。
// 只在界面线程调用（要弹模态框）。
func (a *guiApp) selfOpGameDir() (string, error) {
	if a.cfg != nil && strings.TrimSpace(a.cfg.GameDir) != "" {
		if dir, _, err := selfInspectGameDir(a.cfg.GameDir); err == nil {
			return dir, nil
		}
	}

	// 第一次点 ③ 必须让用户自己选：对话框从「此电脑」开始，不预设任何游戏路径。
	var ownerForm walk.Form
	owner := uintptr(0)
	if a.mw != nil {
		ownerForm = a.mw
		owner = uintptr(a.mw.Handle())
	}
	for i := 0; i < 3; i++ {
		raw, accepted, shown := selfPickFolderAtDrives(owner, "选择穿越火线游戏目录（里面要有 x64\\crossfire.exe）")
		if !shown { // shell 对话框弹不出来时退回 walk 的实现
			raw, accepted = selfPickGameDir(ownerForm, "")
		}
		if !accepted {
			return "", errors.New("已取消选择游戏目录")
		}
		dir, warn, err := selfInspectGameDir(raw)
		if err == nil {
			a.rememberGameDir(dir)
			if warn != "" {
				a.ymm.setLabel("游戏目录提示：" + warn)
			}
			return dir, nil
		}
		walk.MsgBox(a.mw, "游戏目录不合适", err.Error()+"\n\n请选择包含 x64\\crossfire.exe 的穿越火线安装目录。", walk.MsgBoxIconWarning)
	}
	return "", errors.New("连续 3 次选择的目录都不对")
}

// rememberGameDir 把目录写进 config.json，下次点 ③ 不再问。
func (a *guiApp) rememberGameDir(dir string) {
	if a.cfg == nil || strings.TrimSpace(dir) == "" {
		return
	}
	if a.cfg.GameDir == dir {
		return
	}
	a.cfg.GameDir = dir
	_ = a.cfg.save()
}
