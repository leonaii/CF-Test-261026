package main

// 枪神纪 TCLS 上号器。
//
// 依据 2026-11-04 的 ProcMon 全量抓取（395 万行）+ IDA 静态逆向，已确认的官方链路：
//
//	wegame.exe
//	  -> client_launcher.exe --SleepOnly=1
//	     -> tcls_core.exe ppid=<wegame pid>
//	        -> tps.exe -q <UID> -src=tgp -game_id 33 -area 1 -zone_id 18220545
//
// 票据不在文件里，而是由 wegame 通过共享内存 TCLS_CORE_SHARE_MEMORY_*
// （0x49B = 1179 字节）传给 tcls_core：IDA sub_40C370 创建并 memmove 整块，
// sub_40C1C8 读回，CLaunchUI::CheckShareMemoryValue 校验 OpenID / LoginKey，
// TCLSHostApp::SaveLoginInfoToTCLSApp 写 SetOpenID + JumpSignature。
// tcls_core 另写 AppData\Roaming\Tencent\TCLSCore\UserInfo.ini（ACCOUNT_ID / UID）。
//
// 与 CF 的区别：枪神纪没有 rail_api64.dll，CF 上号器的 qwh_cf_hook.dll 那套
// （hook RailNeedRestartAppForCheckingEnvironment）在这里没有目标，照搬无效。
// 所以本实现走「票据落地 + 官方启动链 + TCLS 共享内存接管」，由官方进程
// 拉起 tps.exe，避免自己 CreateProcess 被 ACE 判为安全组件异常。

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tclsGameDir  = `C:\WeGameApps\枪神纪`
	tclsGameID   = 33
	tclsMemSize  = 0x49B // 1179 字节，IDA sub_40C370 的 CreateFileMappingW 大小
	tclsMemMagic = "TCLS_CORE_SHARE_MEMORY_"
)

// tclsTicketPath 是票据落地文件清单。CF 上号器把明文票据写 3 份，
// 枪神纪沿用同一套 TCLS 启动器，这里覆盖 tcls_core / tps.exe 会读的位置。
func tclsTicketPaths(gameDir string) []string {
	return []string{
		filepath.Join(gameDir, "TCLS", "CF_Gamest"),
		filepath.Join(gameDir, "Binaries", "Win64", "CF_Gamest"),
		filepath.Join(gameDir, "CF_Gamest"),
		filepath.Join(gameDir, "TCLS", "TPS_Gamest"),
	}
}

func tclsMmogPaths(gameDir string) []string {
	return []string{
		filepath.Join(gameDir, "TCLS", "mmog_data.xml"),
		filepath.Join(gameDir, "Binaries", "Win64", "mmog_data.xml"),
	}
}

// tclsBuildGamest 按 CF 上号器成功写下的明文票据格式构造。
// 20:53 那次 CF 成功上号时 CF_Gamest 的字段就是这些。
func tclsBuildGamest(op *selfGameTicketData, qq string) string {
	var b strings.Builder
	put := func(k, v string) {
		if v != "" {
			fmt.Fprintf(&b, "%s=%s\n", k, v)
		}
	}
	put("uin", qq)
	put("tgpid", firstNonEmpty(op.Tgpid, op.Uin))
	put("openid", op.OpenID)
	put("open_key", op.OpenKey)
	put("pf_key", op.PfKey)
	put("game_client_key", op.CTT)
	put("game_st", op.GameST)
	put("game_st_key", op.GameSTKey)
	put("client_key", op.ClientKey)
	put("st_key", op.STKey)
	put("ctt", op.CTT)
	put("biz_ticket", op.BizTicket)
	return b.String()
}

// tclsReplaceUin 保留 mmog_data.xml 原有结构（枪神纪的 Version/WorldID
// 与 CF 不同，不能用 CF 的模板覆盖），只替换 <Uin> 标签内容。
func tclsReplaceUin(s, qq string) (string, bool) {
	const open, close = "<Uin>", "</Uin>"
	i := strings.Index(s, open)
	if i < 0 {
		return s, false
	}
	j := strings.Index(s[i+len(open):], close)
	if j < 0 {
		return s, false
	}
	j += i + len(open)
	return s[:i+len(open)] + qq + s[j:], true
}

// tclsWriteTicketFiles 把目标账号票据写进游戏目录，返回实际写入的路径。
// 写之前备份原文件，关闭窗口时由 tclsRestoreTicketFiles 还原。
func tclsWriteTicketFiles(gameDir string, op *selfGameTicketData, qq string, logf func(string, ...any)) ([]string, error) {
	var written []string
	st := tclsBuildGamest(op, qq)
	for _, p := range tclsTicketPaths(gameDir) {
		if err := tclsBackupAndWrite(p, []byte(st), logf); err == nil {
			written = append(written, p)
		}
	}
	for _, p := range tclsMmogPaths(gameDir) {
		raw, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if next, ok := tclsReplaceUin(string(raw), qq); ok {
			if err := tclsBackupAndWrite(p, []byte(next), logf); err == nil {
				written = append(written, p)
			}
		}
	}
	// tcls_core 会重写 UserInfo.ini，这里先预置目标 QQ，避免它读到旧账号。
	if p, err := tclsUserInfoPath(); err == nil {
		body := fmt.Sprintf("[USER_INFO]\r\nACCOUNT_ID=%s\r\nUID=%s\r\n", qq, firstNonEmpty(op.Tgpid, qq))
		if err := tclsBackupAndWrite(p, []byte(body), logf); err == nil {
			written = append(written, p)
		}
	}
	if len(written) == 0 {
		return nil, fmt.Errorf("没有任何票据落地文件写入成功（检查游戏目录 %s）", gameDir)
	}
	logf("票据落地 %d 个文件（QQ %s）", len(written), qq)
	return written, nil
}

func tclsBackupAndWrite(p string, data []byte, logf func(string, ...any)) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	bak := p + ".cgbak"
	if _, err := os.Stat(bak); os.IsNotExist(err) {
		if old, e := os.ReadFile(p); e == nil {
			_ = os.WriteFile(bak, old, 0o644)
		} else {
			// 原本不存在：写个空标记，清理时删掉新文件而不是留残骸
			_ = os.WriteFile(bak, nil, 0o644)
		}
	}
	if err := os.WriteFile(p, data, 0o644); err != nil {
		return err
	}
	logf("已写 %s（%d 字节）", p, len(data))
	return nil
}

// tclsRestoreTicketFiles 关闭窗口时还原现场，避免污染用户下次正常上号。
func tclsRestoreTicketFiles(paths []string, logf func(string, ...any)) {
	for _, p := range paths {
		bak := p + ".cgbak"
		old, err := os.ReadFile(bak)
		if err != nil {
			continue
		}
		if len(old) == 0 {
			_ = os.Remove(p) // 原本不存在
		} else {
			_ = os.WriteFile(p, old, 0o644)
		}
		_ = os.Remove(bak)
	}
	logf("票据落地文件已还原（%d 个）", len(paths))
}

// tclsLaunchOffline 主流程：取实时 OP -> 写票据 -> 接管 TCLS 共享内存
// -> 由官方 WeGame 启动 -> 保持到用户关闭状态窗口。
func tclsLaunchOffline(dir, qq string, sess *tpsWGSession) error {
	logf := sess.logf
	sess.setStatus("检查官方 WeGame…")
	if !tpsWGWeGameRunning() {
		return fmt.Errorf("未检测到官方 WeGame 进程，请先打开并登录 WeGame 再上号")
	}
	logf("官方 WeGame 已就绪")

	sess.setStatus("换取登录票据…")
	op, err := tpsWGGetOP(qq, logf)
	if err != nil {
		return err
	}
	logf("OP 就绪 openid=%dB open_key=%dB ctt=%dB", len(op.OpenID), len(op.OpenKey), len(op.CTT))

	// 不要碰游戏目录、不要自己启动 TPS.exe。
	// v1.3.21 实测：往枪神纪目录写 CF_Gamest / 改 mmog_data.xml / UserInfo.ini，
	// 或用降权令牌自己拉起 TPS.exe 绕开 client_launcher -> tcls_core 链，
	// 都会被 ACE 完整性校验判为「安全组件运行异常 (13,131083,140)」直接拒绝启动。
	// 换号点必须在官方认可的路径上，靠静态逆向确认后再动，不做对抗。
	logf("实时 OP 已就绪（openid/open_key/ctt 各 %d/%d/%d 字节）", len(op.OpenID), len(op.OpenKey), len(op.CTT))
	logf("当前实现不修改游戏目录、不自行启动进程，避免触发 ACE 完整性校验。")

	sess.setStatus("只读观测 TCLS_CORE_SHARE_MEMORY_* 票据通道…")
	stopProbe := tclsStartSharedMemProbe(dir, logf)
	defer stopProbe()

	sess.arm()
	sess.setStatus("由官方 WeGame 启动枪神纪…")
	if err := tpsWGStartGame(tclsGameDir, logf); err != nil {
		logf("官方启动失败（可在 WeGame 里手动点「开始游戏」）: %v", err)
	}

	sess.setStatus("观测中（未做账号替换）。关闭本窗口即强制结束 WeGame / 枪神纪。")
	logf("等待你关闭状态窗口…")
	<-sess.stopped
	return nil
}
