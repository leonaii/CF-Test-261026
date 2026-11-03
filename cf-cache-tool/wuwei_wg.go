package main

// 无畏-WEG.exe（WeGame 注入版）离线桥接。
//
// 与 无畏上号器.exe（ACLOS 代理版）同源，静态还原结论：
//   - 卡密协议完全一样：/api/card/boot、redeem、log，服务器地址同为
//     http://111.170.7.81:8750，24 字节 XOR 槽位 + FNV-1a64 校验常量一致；
//   - exe 内嵌的 RSA 私钥与 jigeRSA2048N 对应，所以直接复用 jigeSeal 信封；
//   - 解出的明文同样是 [QQ] ini（Token010A/Token0143/sessionKey/ct/st/tgpid/
//     machine_id …），wuweiConvertTicket 生成的 6 键真票照样能解。
//
// 区别只在落地方式：WEG 不启动 AclosGameProxy，而是释放内嵌的
// ValWgInject.exe + wg_hook.dll，通过共享内存把票据注入官方 WeGame 进程
// （wegame.exe / rail.exe / railservice.exe）。所以硬前提是**官方 WeGame 已经
// 打开并登录**，否则注入必然失败；CG 侧只做提示，不弹框拦人。
//
// 另外 WEG 不读 config.ini 里的 exe= 键（那是老版选 ACLOS 代理用的），
// 也不读 host.dat（地址硬编码在 exe 里），所以走 patch 地址槽这一条路。

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	wuweiWGExeName    = "无畏-WEG.exe"
	wuweiWGPatchedExe = "无畏-WEG.local.exe"
	wuweiWGListenAddr = "127.0.0.1:18794"
	// 必须正好 24 字节，和 wuweiStubURL 等长（exe 里是定长槽位）
	wuweiWGLocalURL = "http://127.0.0.1:18794/x"
)

// wuweiWGWorkDir 单独一个工作目录：它和老版无畏 / 圆梦共用 config.ini 会互相覆盖。
func wuweiWGWorkDir(base string) string { return filepath.Join(base, "wuwei-weg") }

// wuweiWGLocalCandidates 本地已有的 WEG 上号器位置（开发者目录 / 绿色包）。
func wuweiWGLocalCandidates() []string {
	base, _ := os.Executable()
	exeDir := filepath.Dir(base)
	return []string{
		filepath.Join(exeDir, wuweiWGExeName),
		filepath.Join(exeDir, "other-tools", "无畏", wuweiWGExeName),
		filepath.Join(exeDir, "..", "other-tools", "无畏", wuweiWGExeName),
		filepath.Join(exeDir, "..", "..", "other-tools", "无畏", wuweiWGExeName),
		filepath.Join(exeDir, "资源", "上号器", "wuwei01", wuweiWGExeName),
		filepath.Join(exeDir, "..", "资源", "上号器", "wuwei01", wuweiWGExeName),
		filepath.Join(exeDir, "..", "..", "资源", "上号器", "wuwei01", wuweiWGExeName),
	}
}

// wuweiWGEnsureLocal 云端拿不到时，从本地候选目录拷一份到工作目录。
func wuweiWGEnsureLocal(workDir string) (string, error) {
	dst := filepath.Join(workDir, wuweiWGExeName)
	if st, err := os.Stat(dst); err == nil && st.Size() > 0 {
		return dst, nil
	}
	for _, c := range wuweiWGLocalCandidates() {
		st, err := os.Stat(c)
		if err != nil || st.IsDir() || st.Size() == 0 {
			continue
		}
		if err := os.MkdirAll(workDir, 0o700); err != nil {
			return "", err
		}
		if err := copyFileLocal(c, dst); err != nil {
			return "", fmt.Errorf("复制上号器失败: %w", err)
		}
		return dst, nil
	}
	return "", fmt.Errorf("本地没有 %s", wuweiWGExeName)
}

// wuweiWGPrepareExe 生成/复用打过本地地址的副本，原始 exe 不动。
func wuweiWGPrepareExe(src, dir string) (string, error) {
	raw, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("读上号器失败: %w", err)
	}
	patched, err := wuweiPatchExeURL(raw, wuweiWGLocalURL)
	if err != nil {
		return "", err
	}
	dst := filepath.Join(dir, wuweiWGPatchedExe)
	if cur, err := os.ReadFile(dst); err == nil && bytes.Equal(cur, patched) {
		return dst, nil
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(dst, patched, 0o700); err != nil {
		return "", fmt.Errorf("写入上号器副本失败: %w", err)
	}
	return dst, nil
}

// wuweiWGWriteConfig 写 [ui] 段：占位卡密 + remember，不写 exe=。
// WEG 没有「选 ACLOS 代理」这一步，多写 exe= 只会让配置变脏。
func wuweiWGWriteConfig(dir string) (card string, err error) {
	path := filepath.Join(dir, "config.ini")
	vals := map[string]string{}
	if raw, rerr := os.ReadFile(path); rerr == nil {
		section := ""
		for _, line := range strings.Split(string(raw), "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, ";") || strings.HasPrefix(t, "#") {
				continue
			}
			if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
				section = strings.ToLower(strings.Trim(t, "[]"))
				continue
			}
			if section != "ui" {
				continue
			}
			if k, v, ok := strings.Cut(t, "="); ok {
				vals[strings.ToLower(strings.TrimSpace(k))] = strings.TrimSpace(v)
			}
		}
	}
	card = jigeRandomCard()
	vals["card"] = card
	vals["remember"] = "1"
	delete(vals, "exe")

	var b strings.Builder
	b.WriteString("[ui]\r\n")
	for _, k := range []string{"card", "remember"} {
		b.WriteString(k + "=" + vals[k] + "\r\n")
		delete(vals, k)
	}
	for k, v := range vals {
		b.WriteString(k + "=" + v + "\r\n")
	}
	if err = os.WriteFile(path, []byte(b.String()), 0o600); err != nil {
		return card, err
	}
	return card, nil
}

// wuweiWGLaunchOffline 离线上号主流程：补丁副本 → 本地假卡密服务 → config.ini → 启动。
func wuweiWGLaunchOffline(launcher, workDir, ticket, qq string, logf func(string, ...any)) error {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	if strings.TrimSpace(ticket) == "" {
		return fmt.Errorf("离线票据为空")
	}

	exe, err := wuweiWGPrepareExe(launcher, workDir)
	if err != nil {
		return err
	}
	stop, err := wuweiStartServerAt(ticket, qq, wuweiWGListenAddr)
	if err != nil {
		return err
	}
	defer stop()

	card, err := wuweiWGWriteConfig(workDir)
	if err != nil {
		logf("写 config.ini 失败: %v", err)
	} else {
		logf("占位卡密 %s", card)
	}

	cmd := exec.Command(exe)
	cmd.Dir = workDir
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动上号器失败: %w", err)
	}
	logf("无畏WG 已启动（本地卡密服务 %s）", wuweiWGListenAddr)
	return cmd.Wait()
}
