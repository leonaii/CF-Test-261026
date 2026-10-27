// CG 轻量启动器：启动前检查云端版本和 MD5，必要时更新 %AppData%\CG\CG-app.exe。
package main

import (
	"bytes"
	"compress/gzip"
	"crypto/ed25519"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	cloudShareBase = "https://cloud.51i.cc/public.php/webdav/"
	appManifest    = "app/version.json"
	appGzip        = "app/CG-app.exe.gz"
	installApp     = "CG-app.exe"
	installState   = "CG-app.version.json"
	updatePubKey   = "b1a9bd1f59c40a4671fa79f0c1dc16833403482045b580f3ebb53fe7249e40b5"
	mask           = "cgx0r"
)

var (
	shareToken = buildShareToken()
	sharePass  = buildSharePass()
)

type appManifestFile struct {
	Version string `json:"version"`
	File    string `json:"file"`
	SHA256  string `json:"sha256"`
	MD5     string `json:"md5"`
	Size    int64  `json:"size"`
	Notes   string `json:"notes"`
	Sig     string `json:"sig"`
}

func init() {
	if shareToken == "" {
		shareToken = strings.TrimSpace(os.Getenv("CG_SHARE_TOKEN"))
	}
	if sharePass == "" {
		sharePass = strings.TrimSpace(os.Getenv("CG_SHARE_PASS"))
	}
}

func main() {
	if r := recover(); r != nil {
		fatal(fmt.Errorf("内部异常: %v", r))
	}
	app, err := resolveApp()
	if err != nil {
		fatal(err)
	}
	if err := launch(app); err != nil {
		fatal(err)
	}
}

func resolveApp() (string, error) {
	var (
		m   *appManifestFile
		raw []byte
		err error
	)
	if localMode() {
		m, raw, err = loadLocalManifest()
	} else {
		m, raw, err = loadCloudManifest()
	}
	if err != nil {
		if p, ok := installedApp(); ok {
			return p, nil
		}
		return "", fmt.Errorf("获取云端主程序清单失败: %w", err)
	}
	download := downloadCloudGzip
	if localMode() {
		download = downloadLocalGzip
	}
	app, err := ensureApp(m, raw, download)
	if err == nil {
		return app, nil
	}
	if p, ok := installedApp(); ok {
		return p, nil
	}
	return "", fmt.Errorf("下载主程序失败: %w", err)
}

func loadLocalManifest() (*appManifestFile, []byte, error) {
	base, err := exeDir()
	if err != nil {
		return nil, nil, err
	}
	raw, err := os.ReadFile(filepath.Join(base, "app", "version.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("读取本地清单失败: %w", err)
	}
	m, err := parseManifest(raw)
	return m, raw, err
}

func loadCloudManifest() (*appManifestFile, []byte, error) {
	raw, err := httpGet(appManifest, 1<<20)
	if err != nil {
		return nil, nil, err
	}
	m, err := parseManifest(raw)
	return m, raw, err
}

func downloadLocalGzip() ([]byte, error) {
	base, err := exeDir()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(base, "app", "CG-app.exe.gz"))
}

func parseManifest(raw []byte) (*appManifestFile, error) {
	var m appManifestFile
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("清单格式错误: %w", err)
	}
	if !validVersion(m.Version) {
		return nil, errors.New("清单版本号不合法")
	}
	if m.Size <= 0 || m.Size > 128<<20 {
		return nil, errors.New("清单文件大小不合法")
	}
	if len(strings.TrimSpace(m.MD5)) != 32 {
		return nil, errors.New("清单缺少有效 MD5")
	}
	if err := verifySig(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// ensureApp 固定安装到 %AppData%\CG\CG-app.exe。
// 只有版本号或 MD5 不一致时才下载覆盖；启动器每次都会进入这里检查。
func ensureApp(m *appManifestFile, rawManifest []byte, download func() ([]byte, error)) (string, error) {
	dir, err := installDir()
	if err != nil {
		return "", err
	}
	app := filepath.Join(dir, installApp)
	state := filepath.Join(dir, installState)

	if localRaw, e := os.ReadFile(state); e == nil {
		if local, e := parseManifest(localRaw); e == nil &&
			local.Version == m.Version &&
			strings.EqualFold(local.MD5, m.MD5) &&
			strings.EqualFold(fileMD5(app), m.MD5) {
			return app, nil
		}
	}
	// 状态文件丢失但文件与云端 MD5 完全一致：直接补状态，不重复下载。
	if strings.EqualFold(fileMD5(app), m.MD5) {
		_ = os.WriteFile(state, rawManifest, 0o600)
		cleanupLegacyCache()
		return app, nil
	}
	if download == nil {
		return "", errors.New("本地主程序与云端版本或 MD5 不一致")
	}

	gz, err := download()
	if err != nil {
		return "", err
	}
	data, err := gunzipApp(gz, m.Size)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	if !strings.EqualFold(hex.EncodeToString(sum[:]), m.SHA256) {
		return "", errors.New("主程序 sha256 校验失败")
	}
	if !strings.EqualFold(bytesMD5(data), m.MD5) {
		return "", errors.New("主程序 MD5 校验失败")
	}
	if err := writeAtomic(app, data); err != nil {
		return "", err
	}
	if err := os.WriteFile(state, rawManifest, 0o600); err != nil {
		return "", err
	}
	cleanupLegacyCache()
	return app, nil
}

// installedApp 在云端不可用时使用上次校验通过的主程序。
func installedApp() (string, bool) {
	dir, err := installDir()
	if err != nil {
		return "", false
	}
	raw, err := os.ReadFile(filepath.Join(dir, installState))
	if err != nil {
		return "", false
	}
	m, err := parseManifest(raw)
	if err != nil {
		return "", false
	}
	app := filepath.Join(dir, installApp)
	if !strings.EqualFold(fileMD5(app), m.MD5) {
		return "", false
	}
	return app, true
}

func downloadCloudGzip() ([]byte, error) {
	return httpGet(appGzip, 64<<20)
}

func httpGet(path string, limit int64) ([]byte, error) {
	if shareToken == "" {
		return nil, errors.New("只读分享 token 未配置")
	}
	req, err := http.NewRequest(http.MethodGet, cloudShareBase+path, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(shareToken, sharePass)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "CG-launcher")
	cli := &http.Client{
		Timeout: 10 * time.Minute,
		Transport: &http.Transport{
			Proxy:               nil,
			ForceAttemptHTTP2:   true,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
		if err != nil {
			return nil, err
		}
		if int64(len(data)) > limit {
			return nil, errors.New("云端文件超过大小限制")
		}
		return data, nil
	case http.StatusNotFound:
		return nil, errors.New("云端文件不存在（404）")
	case http.StatusUnauthorized:
		return nil, errors.New("只读分享令牌或口令失效（401）")
	default:
		return nil, fmt.Errorf("云端返回 HTTP %d", resp.StatusCode)
	}
}

func gunzipApp(data []byte, size int64) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("主程序压缩包损坏: %w", err)
	}
	defer zr.Close()
	out, err := io.ReadAll(io.LimitReader(zr, size+1))
	if err != nil {
		return nil, err
	}
	if int64(len(out)) != size {
		return nil, fmt.Errorf("主程序解压大小不一致（%d/%d）", len(out), size)
	}
	return out, nil
}

func verifySig(m *appManifestFile) error {
	pub, err := hex.DecodeString(updatePubKey)
	if err != nil || len(pub) != ed25519.PublicKeySize {
		return errors.New("内置公钥无效")
	}
	sig, err := base64.StdEncoding.DecodeString(strings.TrimSpace(m.Sig))
	if err != nil || len(sig) != ed25519.SignatureSize {
		return errors.New("清单签名格式错误")
	}
	payload := m.Version + "\n" + strings.ToLower(strings.TrimSpace(m.SHA256)) + "\n" + strconv.FormatInt(m.Size, 10)
	if !ed25519.Verify(ed25519.PublicKey(pub), []byte(payload), sig) {
		return errors.New("清单签名校验失败（可能被篡改）")
	}
	return nil
}

// installDir 固定使用 Roaming，即 C:\Users\<用户>\AppData\Roaming\CG。
func installDir() (string, error) {
	base := strings.TrimSpace(os.Getenv("APPDATA"))
	if base == "" {
		d, err := os.UserConfigDir()
		if err != nil || d == "" {
			return "", errors.New("无法定位 AppData\\Roaming 目录")
		}
		base = d
	}
	dir := filepath.Join(base, "CG")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("创建 CG 目录失败: %w", err)
	}
	return dir, nil
}

// cleanupLegacyCache 删除上一版启动器放在 LocalAppData 里的主程序缓存。
func cleanupLegacyCache() {
	if d, err := os.UserCacheDir(); err == nil && d != "" {
		_ = os.RemoveAll(filepath.Join(d, "CG", "app"))
	}
}

func bytesMD5(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

func fileMD5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".CG-app-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	_ = os.Remove(path)
	return os.Rename(name, path)
}

func launch(app string) error {
	cmd := exec.Command(app, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if wd, err := os.Getwd(); err == nil {
		cmd.Dir = wd
	}
	err := cmd.Run()
	if err == nil {
		return nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		os.Exit(ee.ExitCode())
	}
	return err
}

func exeDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func validVersion(v string) bool {
	if v == "" || len(v) > 64 || v == "." || v == ".." {
		return false
	}
	for _, c := range v {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '.' || c == '-' || c == '_' {
			continue
		}
		return false
	}
	return true
}

func deob(s string) string {
	if s == "" {
		return ""
	}
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	k := []byte(mask)
	for i := range raw {
		raw[i] ^= k[i%len(k)]
	}
	return string(raw)
}

func fatal(err error) {
	msg := err.Error()
	fmt.Fprintln(os.Stderr, "CG 启动失败:", msg)
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("MessageBoxW")
	title, _ := syscall.UTF16PtrFromString("CG 启动失败")
	text, _ := syscall.UTF16PtrFromString(msg)
	_, _, _ = proc.Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), 0x10)
	os.Exit(1)
}
