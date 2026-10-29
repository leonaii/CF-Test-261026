package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func usage() {
	fmt.Print(`CG - CF 上号器卡密数据本地缓存工具

用法:
  CG <命令> [参数]

命令:
  patch     改写 host.dat，把上号器指向本地缓存服务（自动备份为 host.dat.orig）
  restore   还原 host.dat
  status    查看当前 host.dat 指向
  serve     启动本地缓存服务（录制 / 回放 / 透传）
  launch    一键：改写 host.dat + 启动服务 + 拉起上号器 + 退出后自动还原
  list      列出已缓存的卡密数据
  dump      解密并打印缓存的账号数据
  push      把本地缓存打包上传到 Nextcloud（云备份）
  pull      从 Nextcloud 下载缓存并合并到本地（云下载）
  cloud     查看云端备份信息
  gs        查封号：处罚记录 / 冻结 / 在线 / 实名防沉迷 / 登录保护

示例:
  CG serve
  CG serve -pick                 交互选号，运行中输入序号随时切换
  CG launch -pick                交互选号 + 拉起上号器
  CG serve -pin ABCDEFGHIJKLMNOP 固定使用某个卡密的缓存
  CG list -detail                按卡密查看账号详情
  CG dump -card ABCDEFGHIJKLMNOP

每个命令都支持 -h 查看参数。
`)
}

func logf(format string, args ...any) {
	all := append([]any{time.Now().Format("15:04:05")}, args...)
	fmt.Printf("%s  "+format+"\n", all...)
}

func main() {
	// 旧版本留下的启动日志和上号器下载缓存都已停用，启动时顺手清掉。
	_ = os.Remove(filepath.Join(configDirPath(), "startup.log"))
	cleanupLoaderCache()
	if len(os.Args) < 2 {
		// 单实例：已经开了一个就不再开第二个
		if !ensureSingleInstance() {
			warnAlreadyRunning()
			return
		}
		// 无参数：Windows 下打开图形界面，其他平台给出提示
		if err := runGUI(); err != nil {
			fmt.Fprintln(os.Stderr, "错误:", err)
			usage()
			os.Exit(2)
		}
		// 退出前清掉为绕过 2310 写进游戏目录的伪造文件
		if n := selfCleanupRailArtifacts(); n > 0 {
			fmt.Printf("已清理 rail 残留文件 %d 个\n", n)
		}
		// 更新完成：界面已退出，这里拉起新版本
		restartIfNeeded()
		return
	}
	var err error
	switch os.Args[1] {
	case "patch":
		err = cmdPatch(os.Args[2:])
	case "restore":
		err = cmdRestore(os.Args[2:])
	case "status":
		err = cmdStatus(os.Args[2:])
	case "serve":
		err = cmdServe(os.Args[2:], false)
	case "launch":
		err = cmdServe(os.Args[2:], true)
	case "list":
		err = cmdList(os.Args[2:])
	case "dump":
		err = cmdDump(os.Args[2:])
	case "push":
		err = cmdPush(os.Args[2:])
	case "pull":
		err = cmdPull(os.Args[2:])
	case "cloud":
		err = cmdCloud(os.Args[2:])
	case "gs":
		err = cmdGs(os.Args[2:])
	case "clean":
		err = cmdClean()
	case "version", "-v", "--version":
		fmt.Println("游戏租号客户端", toolVersion)
	case "help", "-h", "--help":
		usage()
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "错误:", err)
		os.Exit(1)
	}
}

func cmdPatch(args []string) error {
	fs := flag.NewFlagSet("patch", flag.ExitOnError)
	dir := fs.String("dir", ".", "host.dat 所在目录")
	listen := fs.String("listen", "127.0.0.1:18800", "本地缓存服务监听地址")
	_ = fs.Parse(args)

	h := NewHostDatManager(filepath.Join(*dir, "host.dat"))
	orig, changed, err := h.Patch(*listen)
	if err != nil {
		return err
	}
	if !changed {
		logf("host.dat 已指向 http://%s，无需修改", *listen)
		return nil
	}
	back, err := DecryptHostDat(h.Path)
	if err != nil {
		return fmt.Errorf("写入后校验失败: %w", err)
	}
	logf("host.dat 已改写: %s -> %s", orig, back)
	logf("原文件已备份: %s", h.Bak)
	return nil
}

func cmdRestore(args []string) error {
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	dir := fs.String("dir", ".", "host.dat 所在目录")
	_ = fs.Parse(args)

	h := NewHostDatManager(filepath.Join(*dir, "host.dat"))
	if err := h.Restore(); err != nil {
		return fmt.Errorf("还原失败（备份不存在？）: %w", err)
	}
	u, _ := DecryptHostDat(h.Path)
	logf("host.dat 已还原，当前指向: %s", u)
	return nil
}

func cmdStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	dir := fs.String("dir", ".", "host.dat 所在目录")
	_ = fs.Parse(args)

	h := NewHostDatManager(filepath.Join(*dir, "host.dat"))
	cur, err := DecryptHostDat(h.Path)
	if err != nil {
		return err
	}
	fmt.Println("当前 host.dat :", cur)
	if _, err := os.Stat(h.Bak); err == nil {
		if u, e := DecryptHostDat(h.Bak); e == nil {
			fmt.Println("备份 host.dat.orig:", u)
		}
	}
	return nil
}

func buildOptions(fs *flag.FlagSet, args []string, dir, listen, upstream, cache, ttl *string,
	all, offline, refresh, verbose *bool) (Options, *HostDatManager, error) {

	_ = fs.Parse(args)

	h := NewHostDatManager(filepath.Join(*dir, "host.dat"))
	up := strings.TrimSpace(*upstream)
	if up == "" || strings.EqualFold(up, "auto") {
		resolved, err := h.Upstream(*listen)
		if err != nil {
			return Options{}, nil, fmt.Errorf("无法确定真实上游服务器: %w", err)
		}
		up = resolved
	}
	d, err := time.ParseDuration(*ttl)
	if err != nil {
		return Options{}, nil, fmt.Errorf("-ttl 解析失败: %w", err)
	}
	opts := Options{
		Listen:   *listen,
		Upstream: up,
		CacheDir: *cache,
		TTL:      d,
		All:      *all,
		Offline:  *offline,
		Refresh:  *refresh,
		Verbose:  *verbose,
		Log:      logf,
	}
	return opts, h, nil
}

func banner(opts Options) {
	fmt.Println("------------------------------------------------------------")
	fmt.Println("本地缓存服务 : http://" + opts.Listen)
	fmt.Println("上游服务器   : " + opts.Upstream)
	fmt.Println("缓存目录     : " + absPath(opts.CacheDir))
	if opts.Offline {
		fmt.Println("运行模式     : 纯离线（只回放，不访问上游）")
	} else if opts.All {
		fmt.Println("运行模式     : redeem + 其余接口均缓存（TTL 见 -ttl）")
	} else {
		fmt.Println("运行模式     : redeem 永久缓存回放，其余接口透传")
	}
	fmt.Println("------------------------------------------------------------")
}

func absPath(p string) string {
	if a, err := filepath.Abs(p); err == nil {
		return a
	}
	return p
}

func cmdServe(args []string, launch bool) error {
	name := "serve"
	if launch {
		name = "launch"
	}
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	dir := fs.String("dir", ".", "host.dat 与上号器所在目录")
	listen := fs.String("listen", "127.0.0.1:18800", "本地缓存服务监听地址")
	upstream := fs.String("upstream", "auto", "真实卡密服务器，auto 表示从 host.dat.orig 读取")
	cache := fs.String("cache", "cache", "缓存目录")
	ttl := fs.String("ttl", "0", "非 redeem 接口的缓存有效期，如 10m；0 表示永久")
	all := fs.Bool("all", true, "录制所有接口（默认开启，-all=false 只录 redeem）")
	offline := fs.Bool("offline", false, "纯离线：只回放缓存，绝不访问上游")
	refresh := fs.Bool("refresh", false, "本次运行忽略已有缓存，重新录制一次")
	verbose := fs.Bool("v", false, "输出更多日志")
	exe := fs.String("exe", "", "上号器路径（launch 模式），默认自动查找")
	pick := fs.Bool("pick", false, "启动时按卡密交互选择要使用的账号，运行中可切换")
	pin := fs.String("pin", "", "固定使用某个卡密的缓存，不交互")
	logPath := fs.String("log", "", "同时写入运行日志文件（含账号明文），例如 run.log")

	opts, h, err := buildOptions(fs, args, dir, listen, upstream, cache, ttl, all, offline, refresh, verbose)
	if err != nil {
		return err
	}
	if *logPath != "" {
		rl, err := openRunLogger(*logPath)
		if err != nil {
			return fmt.Errorf("无法写入日志 %s: %w", *logPath, err)
		}
		defer rl.close()
		base := opts.Log
		opts.Log = func(format string, args ...any) {
			base(format, args...)
			rl.line(format, args...)
		}
		opts.OnPayload = func(action, card, ep string, body []byte) {
			dumpPayload(rl, action, card, ep, body)
		}
	}

	var exePath string
	if launch {
		exePath = *exe
		if exePath == "" {
			matches, _ := filepath.Glob(filepath.Join(*dir, "CF上号器*.exe"))
			if len(matches) == 1 {
				exePath = matches[0]
			}
		}
		if exePath == "" {
			return fmt.Errorf("未找到上号器 exe，请用 -exe 指定")
		}
		if _, err := os.Stat(exePath); err != nil {
			return fmt.Errorf("上号器不存在: %s", exePath)
		}
		orig, changed, err := h.Patch(opts.Listen)
		if err != nil {
			return err
		}
		if changed {
			logf("host.dat 已改写: %s -> http://%s", orig, opts.Listen)
			// 只还原本次真正改过的文件，避免覆盖用户已有的本地指向
			defer func() {
				if err := h.Restore(); err != nil {
					logf("还原 host.dat 失败: %v", err)
				} else {
					logf("host.dat 已还原")
				}
			}()
		} else {
			logf("host.dat 已指向 http://%s，本次不做改动，退出时也不会覆盖", opts.Listen)
		}
		// 端口占用检查
		ln, err := net.Listen("tcp", opts.Listen)
		if err != nil {
			return fmt.Errorf("端口 %s 不可用: %w", opts.Listen, err)
		}
		_ = ln.Close()
	}

	banner(opts)
	if *pick || *pin != "" {
		if *pin != "" {
			logf("选号方式: 固定卡密 %s", *pin)
		} else {
			logf("选号方式: 交互选择（运行中输入序号可随时切换）")
		}
		return runPinned(opts, launch, exePath, *dir, *pick, *pin)
	}

	srv := NewServer(opts)
	if launch {
		go func() {
			if err := srv.Run(); err != nil {
				logf("服务退出: %v", err)
			}
		}()
	} else {
		return srv.Run()
	}

	// launch：等服务可用后拉起上号器
	if !waitPort(opts.Listen, 3*time.Second) {
		return fmt.Errorf("本地服务未能在 %s 启动", opts.Listen)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	done := make(chan error, 1)
	cmd := exec.Command(exePath)
	cmd.Dir = absPath(*dir)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		if isElevationRequired(err) {
			return fmt.Errorf("上号器要求管理员权限：请右键以管理员身份运行本工具，" +
				"或改用 1-启动缓存服务.bat 后手动打开上号器")
		}
		return err
	}
	logf("已启动上号器: %s", exePath)
	go func() { done <- cmd.Wait() }()

	select {
	case <-interrupt:
		_ = cmd.Process.Kill()
		<-done
	case <-done:
	}
	logf("上号器已退出，正在还原 host.dat")
	return nil
}

// isElevationRequired 判断 CreateProcess 是否因为目标要求管理员权限而失败（Win32 错误 740）。
func isElevationRequired(err error) bool {
	var pe *os.PathError
	if errors.As(err, &pe) {
		if errno, ok := pe.Err.(syscall.Errno); ok && uintptr(errno) == 740 {
			return true
		}
	}
	return false
}

func waitPort(addr string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func cmdList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	cache := fs.String("cache", "cache", "缓存目录")
	detail := fs.Bool("detail", false, "按卡密展示账号详情（QQ/openid/备注）")
	_ = fs.Parse(args)

	if *detail {
		rows := loadAccountRows(NewCache(*cache))
		if len(rows) == 0 {
			fmt.Println("（暂无 redeem 缓存）")
			return nil
		}
		printAccountTable(rows)
		return nil
	}

	metas, err := NewCache(*cache).List()
	if err != nil {
		return err
	}
	if len(metas) == 0 {
		fmt.Println("（暂无缓存）")
		return nil
	}
	fmt.Printf("%-18s %-26s %-6s %-10s %s\n", "卡密", "接口", "状态", "大小", "录制时间")
	for _, m := range metas {
		fmt.Printf("%-18s %-26s %-6d %-10d %s\n",
			m.Card, m.Endpoint, m.Status, m.Bytes, m.RecordedAt.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func cmdDump(args []string) error {
	fs := flag.NewFlagSet("dump", flag.ExitOnError)
	cache := fs.String("cache", "cache", "缓存目录")
	card := fs.String("card", "", "卡密；留空则取最近一条 redeem 缓存")
	_ = fs.Parse(args)

	c := NewCache(*cache)
	metas, err := c.List()
	if err != nil {
		return err
	}
	var target *Meta
	for i := range metas {
		m := &metas[i]
		if m.Endpoint != "api_card_redeem" {
			continue
		}
		if *card != "" && m.Card != *card {
			continue
		}
		if target == nil || m.RecordedAt.After(target.RecordedAt) {
			target = m
		}
	}
	if target == nil {
		return fmt.Errorf("没有找到 redeem 缓存")
	}
	body, _, err := c.Load(target.Card, target.Endpoint)
	if err != nil {
		return err
	}
	plain, err := OpenEnvelope(body)
	if err != nil {
		return fmt.Errorf("解密失败: %w", err)
	}
	fmt.Println("卡密   :", target.Card)
	fmt.Println("录制于 :", target.RecordedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("----------- 解密内容 -----------")
	fmt.Println(PrettyJSON(plain))
	return nil
}

func cmdPush(args []string) error {
	fs := flag.NewFlagSet("push", flag.ExitOnError)
	cache := fs.String("cache", "cache", "缓存目录")
	name := fs.String("name", "", "云端卡密（目录名），如 leonai")
	_ = fs.Parse(args)

	n, err := normalizeBackupName(*name)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	logf("云备份 -> %s", newCloudClient().urlFor(cloudBaseDir(n, "1")+"/"))
	count, err := cloudSyncAll(ctx, n, "1", *cache)
	if err != nil {
		return err
	}
	logf("上传完成: %d 个账号（每个账号一个 zip + index.json）", count)
	return nil
}

func cmdPull(args []string) error {
	fs := flag.NewFlagSet("pull", flag.ExitOnError)
	cache := fs.String("cache", "cache", "缓存目录")
	name := fs.String("name", "", "云端卡密（目录名），如 leonai")
	_ = fs.Parse(args)

	n, err := normalizeBackupName(*name)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	logf("云下载 <- %s/index.json", newCloudClient().urlFor(cloudBaseDir(n, "1")))
	total, err := cloudPullAll(ctx, n, "1", *cache, func(f string, a ...any) { logf(f, a...) })
	if err != nil {
		return err
	}
	logf("云端共 %d 个账号，已按需补齐到 %s", total, *cache)
	return nil
}

func cmdCloud(args []string) error {
	fs := flag.NewFlagSet("cloud", flag.ExitOnError)
	name := fs.String("name", "", "云端卡密（目录名）")
	listAll := fs.Bool("list", false, "列出云端所有卡密")
	_ = fs.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if *listAll {
		dirs, err := cloudFolders(ctx)
		if err != nil {
			return err
		}
		if len(dirs) == 0 {
			fmt.Println("（云端暂无备份）")
			return nil
		}
		for _, d := range dirs {
			idx, err := fetchCloudIndex(ctx, d, "1")
			if err != nil {
				fmt.Printf("%-16s 读取失败: %v\n", d, err)
				continue
			}
			fmt.Printf("%-16s %d 个账号  更新于 %s\n", d, len(idx.Cards), idx.Updated)
		}
		return nil
	}

	n, err := normalizeBackupName(*name)
	if err != nil {
		return err
	}
	idx, err := fetchCloudIndex(ctx, n, "1")
	if err != nil {
		return err
	}
	if len(idx.Cards) == 0 {
		return errCloudNotFound
	}
	fmt.Println("云端路径:", newCloudClient().urlFor(cloudBaseDir(n, "1")+"/"))
	fmt.Printf("账号数: %d   索引更新: %s\n", len(idx.Cards), idx.Updated)
	fmt.Printf("%-4s %-18s %-12s %-6s %s\n", "#", "CARD", "QQ", "标记", "录制时间")
	for i, e := range idx.Cards {
		fmt.Printf("%-4d %-18s %-12s %-6s %s\n", i+1, e.Card, e.QQ, markLabel(e.Mark), e.Recorded)
	}
	return nil
}

// cmdClean 清理更新残留文件（程序被占用时提示关闭后再试）。
func cmdClean() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Dir(exe)
	total := 0
	for _, base := range []string{"zuhaoclient.exe", "CG.exe", "CG-cli.exe"} {
		p := filepath.Join(dir, base)
		before := 0
		if matches, _ := filepath.Glob(p + ".old*"); len(matches) > 0 {
			before = len(matches)
		}
		left := cleanupArtifacts(p)
		done := before - len(left)
		total += done
		if done > 0 {
			logf("已清理 %s 的 %d 个残留文件", base, done)
		}
		for _, f := range left {
			logf("删除失败（可能仍在运行）: %s", f)
		}
	}
	if n := selfCleanupRailArtifacts(); n > 0 {
		logf("已清理 rail 残留文件 %d 个", n)
		total += n
	}
	if total == 0 {
		logf("没有需要清理的文件")
	}
	return nil
}
