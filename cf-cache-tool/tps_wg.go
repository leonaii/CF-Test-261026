package main

// 枪神纪 WEG 式上号：用户先自己打开并登录官方 WeGame，
// CG 把实时 OP 票据写进 Local\ValWgInjectAcc_v1 共享内存，
// 用 ValWgInject.exe 把 wg_hook.dll 注入 wegame.exe / rail.exe / RailService.exe，
// 最后由官方 WeGame 拉起 TPS.exe。
//
// 逆向依据（IDA sub_1403B1253 + wg_hook.dll 字符串）：
//   - 共享内存大小 0x29C4 = 10692；+0 DWORD 魔数 "VALW"(0x574C4156)
//   - +8 DWORD game_id：WEG 写死 2001715（无畏），枪神纪为 33
//   - 票据 payload 形如 VAL={...} / DNF={...} / CF={...} 的 JSON
//   - +1788 起 200 字节 UTF16 hook dll 路径
//   - hook DLL 按 cmd100_gid<gid>_*.bin 拦截 rail 的票据读取
//   - ValWgInject.exe 枚举 wegame.exe / rail.exe / RailService.exe 做远程线程注入

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	tpsWGMemName       = `Local\ValWgInjectAcc_v1`
	tpsWGMemSize       = 0x29C4     // 10692
	tpsWGMagic         = 0x574C4156 // "VALW"
	tpsWGGameID        = 33
	tpsWGPayloadPrefix = "CF=" // WEG 支持 VAL=/DNF=/CF=；枪神纪与 CF 同构，走 CF 分支

	tpsWGInjectorExe = "ValWgInject.exe"
	tpsWGHookDLL     = "wg_hook.dll"
	tpsWGRemoteDir   = cloudDir + "/tpswg"
)

// tpsWGWorkDir 独立工作目录，避免和 wuwei / wuwei-weg 的 config.ini 互覆盖。
func tpsWGWorkDir(base string) string { return filepath.Join(base, "tpswg") }

type tpsWGMeta struct {
	InjectorSHA string `json:"injector_sha256"`
	HookSHA     string `json:"hook_sha256"`
	Updated     string `json:"updated"`
}

func tpsWGSHA(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

// tpsWGEnsureResources 从云端按需下载注入器与 hook DLL，sha256 校验。
func tpsWGEnsureResources(ctx context.Context, dir string, logf func(string, ...any)) (string, string, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", "", err
	}
	c := newCloudClient()
	raw, err := c.download(ctx, tpsWGRemoteDir+"/meta.json")
	if err != nil {
		return "", "", fmt.Errorf("获取枪神纪注入器清单失败: %w", err)
	}
	var m tpsWGMeta
	if err := json.Unmarshal(raw, &m); err != nil {
		return "", "", fmt.Errorf("注入器清单格式错误: %w", err)
	}

	want := map[string]string{
		tpsWGInjectorExe: strings.ToLower(m.InjectorSHA),
		tpsWGHookDLL:     strings.ToLower(m.HookSHA),
	}
	paths := map[string]string{}
	for name, sum := range want {
		dst := filepath.Join(dir, name)
		if b, e := os.ReadFile(dst); e == nil && (sum == "" || strings.EqualFold(tpsWGSHA(b), sum)) {
			paths[name] = dst
			continue
		}
		logf("按需下载枪神纪注入组件 %s…", name)
		data, e := c.download(ctx, tpsWGRemoteDir+"/"+name)
		if e != nil {
			return "", "", fmt.Errorf("下载 %s 失败: %w", name, e)
		}
		if sum != "" && !strings.EqualFold(tpsWGSHA(data), sum) {
			return "", "", fmt.Errorf("%s 校验失败，已丢弃", name)
		}
		if e := os.WriteFile(dst, data, 0o700); e != nil {
			return "", "", fmt.Errorf("写入 %s 失败: %w", name, e)
		}
		paths[name] = dst
	}
	return paths[tpsWGInjectorExe], paths[tpsWGHookDLL], nil
}

// tpsWGGetOP 复用已验证的链路：实时 ini -> NT 免密刷新 sig -> LoginByQq -> OP。
func tpsWGGetOP(qq string, logf func(string, ...any)) (*selfGameTicketData, error) {
	api := newCardClient()
	if err := api.login(adminPass()); err != nil {
		return nil, fmt.Errorf("登录卡密平台失败: %w", err)
	}
	ini, err := api.getData(qq)
	if err != nil {
		return nil, fmt.Errorf("取实时 ini 失败: %w", err)
	}
	tk, err := parseYmmTicket(ini, qq)
	if err != nil {
		return nil, err
	}
	if tok := strings.TrimSpace(tk.Token0106); tok != "" {
		uin, e := strconv.ParseUint(strings.TrimSpace(tk.QQ), 10, 64)
		if e == nil && uin != 0 {
			var machineID []byte
			if s := strings.TrimSpace(tk.MachineID); s != "" {
				if b, de := hex.DecodeString(s); de == nil && len(b) == 16 {
					machineID = b
				}
			}
			if t0106, de := hex.DecodeString(tok); de == nil && len(t0106) > 0 {
				if nt, ne := selfNTLogin(uin, machineID, t0106); ne == nil {
					tk.Sig = nt.Sig
					logf("NT 免密刷新 sig 成功")
				} else {
					logf("NT 免密刷新失败，沿用 ini 的 sig: %v", ne)
				}
			}
		}
	}
	op, err := selfLoginByIni(tk, "")
	if err != nil {
		return nil, fmt.Errorf("LoginByQq 换 OP 失败: %w", err)
	}
	logf("OP 就绪 openid=%dB open_key=%dB ctt=%dB", len(op.OpenID), len(op.OpenKey), len(op.CTT))
	return op, nil
}

// tpsWGBuildPayload 组织共享内存里的票据 JSON。
func tpsWGBuildPayload(op *selfGameTicketData, qq string) string {
	m := map[string]any{
		"qq":          qq,
		"uin":         op.Uin,
		"game_id":     tpsWGGameID,
		"openid":      op.OpenID,
		"open_key":    op.OpenKey,
		"ctt":         op.CTT,
		"game_ticket": op.GameST,
		"tgp_id":      op.Tgpid,
		"client_key":  op.ClientKey,
		"st_key":      op.STKey,
		"pf_key":      op.PfKey,
		"sdk_traceid": "",
		"biz_ticket":  op.BizTicket,
	}
	b, _ := json.Marshal(m)
	return tpsWGPayloadPrefix + string(b)
}

// tpsWGSession 把「弹框生命周期」和「注入流程」绑在一起：
// 共享内存必须一直保持打开（hook DLL 随时会回来读），所以不能由流程自己定时释放，
// 而是等用户关闭状态窗口时统一清理：清开关 -> 关映射 -> 强制结束 WeGame/枪神纪 -> hook 随进程卸载。
type tpsWGSession struct {
	mu       sync.Mutex
	logf     func(string, ...any)
	uiStatus func(string)
	closeMem func()
	stopped  chan struct{}
	once     sync.Once
	status   string
	armed    bool // 已写入共享内存并注入，关闭时才需要强杀进程
}

func newTpsWGSession(logf func(string, ...any)) *tpsWGSession {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &tpsWGSession{logf: logf, stopped: make(chan struct{}), status: "准备中…"}
}

func (s *tpsWGSession) setStatus(v string) {
	s.mu.Lock()
	s.status = v
	f := s.logf
	ui := s.uiStatus
	s.mu.Unlock()
	f("%s", v)
	if ui != nil {
		ui(v)
	}
}

func (s *tpsWGSession) statusText() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

// setCloseMem 交给会话托管共享内存映射的释放函数。
func (s *tpsWGSession) setCloseMem(f func()) {
	s.mu.Lock()
	s.closeMem = f
	s.armed = true
	s.mu.Unlock()
}

// arm 标记会话已经改动了系统状态，关闭时需要执行强制清理。
func (s *tpsWGSession) arm() {
	s.mu.Lock()
	s.armed = true
	s.mu.Unlock()
}

// wait 可被关闭事件打断的等待；返回 false 表示用户已关闭窗口。
func (s *tpsWGSession) wait(d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-s.stopped:
		return false
	}
}

func (s *tpsWGSession) isStopped() bool {
	select {
	case <-s.stopped:
		return true
	default:
		return false
	}
}

// stop 幂等清理：把共享内存里的替换/录帧开关清 0（让残留 hook 立刻停手），
// 关掉映射，再强制结束 WeGame 与枪神纪进程。hook DLL 是注入到这些进程里的，
// 进程退出即等于卸载，没有单独的 FreeLibrary 通道。
func (s *tpsWGSession) stop() {
	s.once.Do(func() {
		close(s.stopped)
		s.mu.Lock()
		armed := s.armed
		cm := s.closeMem
		s.closeMem = nil
		s.mu.Unlock()
		if !armed {
			s.setStatus("已关闭：本次未写入票据，无需结束进程。")
			return
		}
		s.setStatus("正在清理：结束 WeGame / 枪神纪并还原票据文件…")
		tpsWGForceCleanup(s.logf)
		if cm != nil {
			cm()
		}
		s.setStatus("已清理：进程已结束、票据文件已还原。")
	})
}

// tpsWGLaunchOffline 主流程：前置检查 -> 取 OP -> 写共享内存 -> 注入 -> 官方启动游戏。
// 共享内存保持到 sess.stop()（用户关闭状态窗口）才释放。
func tpsWGLaunchOffline(dir, qq, gameDir string, sess *tpsWGSession) error {
	logf := sess.logf
	sess.setStatus("检查官方 WeGame…")
	if !tpsWGWeGameRunning() {
		return fmt.Errorf("未检测到官方 WeGame 进程，请先打开并登录 WeGame 再上号")
	}
	logf("官方 WeGame 已就绪")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	sess.setStatus("获取注入组件…")
	injector, hook, err := tpsWGEnsureResources(ctx, dir, logf)
	if err != nil {
		return err
	}

	sess.setStatus("换取登录票据…")
	op, err := tpsWGGetOP(qq, logf)
	if err != nil {
		return err
	}
	payload := tpsWGBuildPayload(op, qq)
	logf("票据 payload %d 字节", len(payload))

	closeMem, err := tpsWGMapAndWrite(payload, hook)
	if err != nil {
		return fmt.Errorf("写共享内存失败: %w", err)
	}
	sess.setCloseMem(closeMem)
	logf("共享内存 %s 已写入（game_id=%d）", tpsWGMemName, tpsWGGameID)

	sess.setStatus("注入 wg_hook.dll…")
	injectLog, err := tpsWGRunInjector(injector, dir)
	if err != nil {
		return fmt.Errorf("%w（日志：%s）", err, tpsWGClip(injectLog, 400))
	}
	lines := strings.Split(strings.TrimSpace(injectLog), "\n")
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			logf("注入器: %s", strings.TrimSpace(l))
		}
	}

	sess.setStatus("由官方 WeGame 启动枪神纪…")
	if err := tpsWGStartGame(gameDir, logf); err != nil {
		logf("官方启动枪神纪失败（可在 WeGame 里手动点「开始游戏」）: %v", err)
	}

	// hook_ok 只代表 NtWriteFile 的 detour 装上了，不代表拦到了帧；
	// 必须继续等到 dump_ok（拦到 cmd100 帧）或超时，才能判断协议是否适用。
	sess.setStatus("等待 hook 回报状态…")
	logf("等待 hook DLL 回报状态（最多 60 秒）…")
	last := ""
	sawFrame := false
	for i := 0; i < 12; i++ {
		if !sess.wait(5 * time.Second) {
			return nil // 用户已关闭窗口，stop() 已负责清理
		}
		status, path := tpsWGReadStatus()
		if status == "" || status == last {
			continue
		}
		last = status
		extra := ""
		if path != "" {
			extra = "，dump=" + tpsWGClip(path, 160)
		}
		logf("hook 状态: %s%s", status, extra)
		if strings.Contains(status, "dump_ok") {
			sawFrame = true
			break
		}
	}
	if !sawFrame {
		logf("诊断：hook 已安装，但 60 秒内未拦到 cmd100 帧（没有 dump_ok）。")
		logf("说明枪神纪的登录链路不走 WEG 所 hook 的 rail NtWriteFile/cmd100 协议，票据替换不会触发。")
		logf("可查 %%TEMP%%\\rail_cmd100\\ 是否有 .bin；若为空则这套 hook 对枪神纪无效。")
	}

	// 共享内存一直保持，直到用户关闭状态窗口
	sess.setStatus("注入已就绪。关闭本窗口即强制结束 WeGame / 枪神纪并卸载 hook。")
	logf("共享内存保持打开，等待你关闭状态窗口…")
	<-sess.stopped
	return nil
}

func tpsWGClip(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
