package main

// 云端自动更新：
//
//	CG/update/version.json   清单（版本、文件名、sha256、大小）
//	CG/update/CG.exe         新版本二进制
//
// 启动后后台检查，发现新版本就下载 → 校验 → 替换当前 exe，下次启动生效。
// 全程不阻塞界面，任何失败都只写日志。

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// pendingRestart 记录“更新完成后要拉起的新版本”，由 main 在界面退出后启动。
var pendingRestart atomic.Value // string

// 本进程已经替换到的版本 + 单飞锁（首次界面和主界面可能同时触发检查）
var (
	updateMu      sync.Mutex
	updateApplied atomic.Value // string
	// updateBusy 为 true 表示正在下载/替换，首屏此时要锁住输入，避免用旧版本启动
	updateBusy atomic.Bool
)

const (
	updateDir      = cloudDir + "/update"
	updateManifest = updateDir + "/version.json"
	updateBinary   = updateDir + "/CG.exe"
)

type updateManifestFile struct {
	Version string `json:"version"`
	File    string `json:"file"`
	SHA256  string `json:"sha256"`
	Size    int64  `json:"size"`
	Notes   string `json:"notes"`
	Sig     string `json:"sig"` // ed25519 签名（构建机私钥），客户端用内置公钥验
}

// fetchUpdateManifest 读取云端清单；不存在返回 errCloudNotFound。
func fetchUpdateManifest(ctx context.Context) (*updateManifestFile, error) {
	data, err := newCloudClient().download(ctx, updateManifest)
	if err != nil {
		return nil, err
	}
	var m updateManifestFile
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("更新清单格式错误: %w", err)
	}
	if m.Version == "" {
		return nil, errors.New("更新清单缺少 version")
	}
	// 没有签名或验签失败一律不更新：dav 被写也推不了假包
	if err := verifyUpdateSig(m.Version, m.SHA256, m.Size, m.Sig); err != nil {
		return nil, err
	}
	return &m, nil
}

// downloadUpdate 下载新版本到 exe 同目录的临时文件，并校验大小与 sha256。
func downloadUpdate(ctx context.Context, curExe string, m *updateManifestFile) (string, error) {
	path := updateBinary
	if strings.TrimSpace(m.File) != "" {
		path = updateDir + "/" + strings.TrimPrefix(m.File, "/")
	}

	client := newCloudHTTPClient(10 * time.Minute)
	c := newCloudClient()
	req, err := http.NewRequestWithContext(ctx, "GET", c.urlFor(path), nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("User-Agent", "CG/"+toolVersion)

	resp, err := client.Do(req)
	if err != nil {
		return "", cloudNetworkError("下载更新", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", errCloudNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载更新失败: HTTP %d", resp.StatusCode)
	}

	tmp := filepath.Join(filepath.Dir(curExe), "CG.update.tmp")
	_ = os.Remove(tmp)
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o700)
	if err != nil {
		return "", fmt.Errorf("写入临时文件失败（可能需要管理员权限）: %w", err)
	}

	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(resp.Body, 200<<20))
	closeErr := f.Close()
	if err != nil {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("下载中断: %w", err)
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return "", closeErr
	}

	if m.Size > 0 && n != m.Size {
		_ = os.Remove(tmp)
		return "", fmt.Errorf("大小不一致（期望 %d，实际 %d）", m.Size, n)
	}
	if want := strings.ToLower(strings.TrimSpace(m.SHA256)); want != "" {
		if got := hex.EncodeToString(h.Sum(nil)); got != want {
			_ = os.Remove(tmp)
			return "", errors.New("sha256 校验失败，已丢弃下载文件")
		}
	}
	return tmp, nil
}

// applyUpdate 用新文件替换当前 exe：先把运行中的 exe 改名，再把新文件就位。
// 失败会自动回滚，保证原来还能启动。
func applyUpdate(curExe, newPath string) error {
	old := curExe + ".old"

	// 只使用固定的 .old 名字：删不掉就跳过本次更新，绝不生成第二个文件
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(old); os.IsNotExist(err) {
			break
		}
		if os.Remove(old) == nil {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	if _, err := os.Stat(old); err == nil {
		return errors.New("旧版本文件仍被占用，本次跳过更新")
	}

	var err error
	for i := 0; i < 5; i++ {
		if err = os.Rename(curExe, old); err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		return fmt.Errorf("无法替换当前程序（可能没有写权限或被占用）: %w", err)
	}

	for i := 0; i < 5; i++ {
		if err = os.Rename(newPath, curExe); err == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	_ = os.Rename(old, curExe) // 回滚
	return fmt.Errorf("替换失败，已回滚: %w", err)
}

// cleanupOldExe 启动时删掉上一次替换留下的 .old（失败忽略）。
// cleanupOldExe 清理上次更新留下的 CG.exe.old / CG.exe.old.* / CG.update.tmp。
// 新进程刚起来时旧进程可能还没完全退出，删不掉就后台持续重试，最多约 30 秒。
// cleanupOldExe 启动时清理上次更新留下的 CG.exe.old / .old.* / CG.update.tmp。
// 旧进程可能还没完全退出，删不掉就后台重试；仍删不掉就交给系统在本进程退出后删除。
func cleanupOldExe(logf func(string, ...any)) {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	if logf == nil {
		logf = func(string, ...any) {}
	}

	left := cleanupArtifacts(exe)
	if len(left) == 0 {
		return
	}
	logf("旧版本文件暂未释放，稍后重试删除: %s", strings.Join(left, ", "))

	go func() {
		for i := 0; i < 10 && len(left) > 0; i++ {
			time.Sleep(1500 * time.Millisecond)
			left = cleanupArtifacts(exe)
		}
		if len(left) == 0 {
			return
		}
		// 还在被占用（比如本进程就是从这个文件启动的）：交给系统在退出后删除
		scheduleDelayedDelete(left)
	}()
}

// cleanupArtifacts 删除指定 exe 对应的一批遗留文件，返回没删掉的。
func cleanupArtifacts(exe string) []string {
	targets := []string{exe + ".old", filepath.Join(filepath.Dir(exe), "CG.update.tmp")}
	if matches, e := filepath.Glob(exe + ".old.*"); e == nil {
		targets = append(targets, matches...)
	}
	var left []string
	for _, p := range targets {
		if _, e := os.Stat(p); e != nil {
			continue
		}
		_ = os.Chmod(p, 0o600) // 去掉只读属性，否则删不掉
		if os.Remove(p) != nil {
			left = append(left, p)
		}
	}
	return left
}

// updateEvents 把更新进度广播给第一界面（无人监听时丢弃）。
var updateEvents = make(chan string, 32)

func emitUpdate(msg string) {
	select {
	case updateEvents <- msg:
	default:
	}
}

// appliedVersion 返回本进程已替换到的版本（"" 表示没更新）。
func appliedVersion() string {
	v, _ := updateApplied.Load().(string)
	return v
}

// checkUpdateCore 检查并替换当前 exe，返回本次生效的版本（"" 表示没有更新）。
// 全程加锁：首次配置窗口和主界面同时触发时只会下载一次。
func checkUpdateCore(logf func(string, ...any)) (string, error) {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	updateMu.Lock()
	defer updateMu.Unlock()

	// 本进程已经换过一次了（比如首配窗口里换的），直接返回
	if v, _ := updateApplied.Load().(string); v != "" {
		return v, nil
	}

	// 清单请求必须短超时：云端抖动时不能把首屏一直卡在“检查更新”。
	start := time.Now()
	manifestCtx, cancelManifest := context.WithTimeout(context.Background(), 5*time.Second)
	m, err := fetchUpdateManifest(manifestCtx)
	cancelManifest()
	selfLog("更新清单耗时 %s", time.Since(start).Round(10*time.Millisecond))
	if err != nil {
		if errors.Is(err, errCloudNotFound) {
			return "", nil
		}
		return "", err
	}
	if !versionNewer(m.Version, toolVersion) {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	logf("发现新版本 v%s，开始下载…", m.Version)
	updateBusy.Store(true)

	tmp, err := downloadUpdate(ctx, exe, m)
	if err != nil {
		updateBusy.Store(false)
		return "", err
	}
	if err := applyUpdate(exe, tmp); err != nil {
		_ = os.Remove(tmp)
		updateBusy.Store(false)
		return "", err
	}
	updateApplied.Store(m.Version)
	return m.Version, nil
}
