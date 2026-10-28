#![cfg_attr(not(feature = "console"), windows_subsystem = "windows")]

mod cng;
mod crypto;
mod gzip;
mod net;

use std::ffi::OsStr;
use std::fs;
use std::os::windows::ffi::OsStrExt;
use std::path::{Path, PathBuf};
use std::process::Command;
use std::time::Duration;

const MANIFEST_PATH: &str = "/public.php/webdav/app/version.json";
const APP_GZ_PATH: &str = "/public.php/webdav/app/CG-app.exe.gz";
const ECDSA_X: &str = "e6dc0ae2dcdc968b29e429c3741dcbb263a784b6dc6e593cfe3c3b72899b2eaa";
const ECDSA_Y: &str = "81b77ca002cb17327b831011d500d3949a9802ac2e868ba33fbc3a03107a67e5";

include!(concat!(env!("OUT_DIR"), "/secret.rs"));

#[link(name = "user32")]
extern "system" {
    fn MessageBoxW(hwnd: isize, text: *const u16, caption: *const u16, flags: u32) -> i32;
}

struct Manifest {
    version: String,
    sha256: String,
    md5: String,
    size: u64,
    raw: Vec<u8>,
}

impl Manifest {
    fn parse(raw: &[u8]) -> Result<Self, String> {
        let text = String::from_utf8_lossy(raw);
        let version = json_string(&text, "version").ok_or("清单缺少版本号")?;
        let sha256 = json_string(&text, "sha256").ok_or("清单缺少 SHA256")?;
        let md5 = json_string(&text, "md5").ok_or("清单缺少 MD5")?;
        let sig = json_string(&text, "sig_ecdsa").ok_or("清单尚未启用 ECDSA 签名")?;
        let size = json_u64(&text, "size").ok_or("清单缺少文件大小")?;
        if !valid_version(&version) {
            return Err("清单版本号不合法".into());
        }
        if sha256.len() != 64 || md5.len() != 32 {
            return Err("清单校验值长度错误".into());
        }
        if size == 0 || size > 128 * 1024 * 1024 {
            return Err("清单文件大小不合法".into());
        }
        cng::verify_ecdsa(&version, &sha256, size, ECDSA_X, ECDSA_Y, &sig)?;
        Ok(Self {
            version,
            sha256: sha256.to_ascii_lowercase(),
            md5: md5.to_ascii_lowercase(),
            size,
            raw: raw.to_vec(),
        })
    }
}

fn main() {
    if let Err(e) = run() {
        fatal(&e);
    }
}

fn run() -> Result<(), String> {
    let token = crypto::obf_decode(SHARE_TOKEN_ENC).unwrap_or_default();
    let pass = crypto::obf_decode(SHARE_PASS_ENC).unwrap_or_default();
    if token.is_empty() {
        return Err("只读分享 token 未配置".into());
    }

    cleanup_old_files();
    let (app, state) = install_paths()?;

    let fallback = || -> Option<PathBuf> {
        let raw = fs::read(&state).ok()?;
        let m = Manifest::parse(&raw).ok()?;
        if file_md5(&app).eq_ignore_ascii_case(&m.md5) {
            Some(app.clone())
        } else {
            None
        }
    };

    let selected = match net::http_get(MANIFEST_PATH, &token, &pass, 1 << 20)
        .and_then(|raw| Manifest::parse(&raw))
    {
        Ok(manifest) => match ensure_app(&manifest, &app, &state, &token, &pass) {
            Ok(()) => app.clone(),
            Err(e) => fallback().ok_or_else(|| format!("下载主程序失败：{e}"))?,
        },
        Err(e) => fallback().ok_or_else(|| format!("获取云端主程序清单失败：{e}"))?,
    };

    #[cfg(feature = "console")]
    {
        println!("main={}", selected.display());
        if std::env::var_os("CG_LAUNCHER_CHECK_ONLY").is_some() {
            return Ok(());
        }
    }

    #[cfg(feature = "console")]
    let selected = std::env::var_os("CG_LAUNCHER_TEST_APP")
        .map(PathBuf::from)
        .unwrap_or(selected);
    let args: Vec<_> = std::env::args_os().skip(1).collect();
    Command::new(&selected)
        .args(args)
        .spawn()
        .map_err(|e| format!("启动主程序失败：{e}"))?;
    Ok(())
}

fn ensure_app(
    m: &Manifest,
    app: &Path,
    state: &Path,
    token: &str,
    pass: &str,
) -> Result<(), String> {
    if let Ok(local_raw) = fs::read(state) {
        if let Ok(local) = Manifest::parse(&local_raw) {
            if local.version == m.version
                && local.md5.eq_ignore_ascii_case(&m.md5)
                && file_md5(app).eq_ignore_ascii_case(&m.md5)
            {
                cleanup_old_app_files();
                return Ok(());
            }
        }
    }
    if file_md5(app).eq_ignore_ascii_case(&m.md5) {
        fs::write(state, &m.raw).map_err(|e| format!("写入版本状态失败：{e}"))?;
        cleanup_old_app_files();
        return Ok(());
    }

    let gz = net::http_get(APP_GZ_PATH, token, pass, 64 << 20)?;
    let data = gzip::gunzip(&gz, m.size as usize)?;
    let sha = crypto::hex_lower(&crypto::sha256(&data));
    if !sha.eq_ignore_ascii_case(&m.sha256) {
        return Err("主程序 SHA256 校验失败".into());
    }
    let md5 = crypto::hex_lower(&crypto::md5(&data));
    if !md5.eq_ignore_ascii_case(&m.md5) {
        return Err("主程序 MD5 校验失败".into());
    }
    write_atomic(app, &data)?;
    fs::write(state, &m.raw).map_err(|e| format!("写入版本状态失败：{e}"))?;
    cleanup_old_app_files();
    Ok(())
}

fn install_paths() -> Result<(PathBuf, PathBuf), String> {
    let base = std::env::var_os("APPDATA").ok_or("无法定位 AppData\\Roaming 目录")?;
    let dir = PathBuf::from(base).join("CG");
    fs::create_dir_all(&dir).map_err(|e| format!("创建 CG 目录失败：{e}"))?;
    Ok((dir.join("CG.exe"), dir.join("CG.version.json")))
}

fn cleanup_old_files() {
    if let Ok(exe) = std::env::current_exe() {
        if let Some(dir) = exe.parent() {
            remove_old_in_dir(dir, true);
        }
    }
    if let Some(base) = std::env::var_os("APPDATA") {
        remove_old_in_dir(&PathBuf::from(base).join("CG"), false);
    }
    if let Some(base) = std::env::var_os("LOCALAPPDATA") {
        let _ = fs::remove_dir_all(PathBuf::from(base).join("CG").join("app"));
    }
}

fn cleanup_old_app_files() {
    if let Some(base) = std::env::var_os("APPDATA") {
        let dir = PathBuf::from(base).join("CG");
        for name in ["CG-app.exe", "CG-app.version.json"] {
            let _ = fs::remove_file(dir.join(name));
        }
        remove_old_in_dir(&dir, false);
    }
}

fn remove_old_in_dir(dir: &Path, include_update_tmp: bool) {
    for attempt in 0..10 {
        let mut left = false;
        if let Ok(rd) = fs::read_dir(dir) {
            for ent in rd.flatten() {
                let name = ent.file_name();
                let name = name.to_string_lossy();
                let hit = name == "CG.exe.old"
                    || name.starts_with("CG.exe.old.")
                    || (include_update_tmp && name == "CG.update.tmp");
                if hit {
                    match fs::remove_file(ent.path()) {
                        Ok(_) => {}
                        Err(_) => left = true,
                    }
                }
            }
        }
        if !left {
            break;
        }
        if attempt < 9 {
            std::thread::sleep(Duration::from_millis(100));
        }
    }
}

fn write_atomic(path: &Path, data: &[u8]) -> Result<(), String> {
    let dir = path.parent().ok_or("主程序目录无效")?;
    let tmp = dir.join(format!(".CG-{}.tmp", std::process::id()));
    fs::write(&tmp, data).map_err(|e| format!("写入主程序失败：{e}"))?;
    let _ = fs::remove_file(path);
    fs::rename(&tmp, path).map_err(|e| {
        let _ = fs::remove_file(&tmp);
        format!("替换主程序失败：{e}")
    })
}

fn file_md5(path: &Path) -> String {
    match fs::read(path) {
        Ok(data) => crypto::hex_lower(&crypto::md5(&data)),
        Err(_) => String::new(),
    }
}

fn json_string(src: &str, key: &str) -> Option<String> {
    let needle = format!("\"{}\"", key);
    let at = src.find(&needle)? + needle.len();
    let rest = &src[at..];
    let colon = rest.find(':')?;
    let rest = rest[colon + 1..].trim_start();
    if !rest.starts_with('"') {
        return None;
    }
    let mut out = String::new();
    let mut chars = rest[1..].chars();
    while let Some(c) = chars.next() {
        match c {
            '"' => return Some(out),
            '\\' => match chars.next()? {
                '"' => out.push('"'),
                '\\' => out.push('\\'),
                '/' => out.push('/'),
                'n' => out.push('\n'),
                'r' => out.push('\r'),
                't' => out.push('\t'),
                'u' => {
                    let h: String = chars.by_ref().take(4).collect();
                    if let Ok(v) = u16::from_str_radix(&h, 16) {
                        if let Some(ch) = char::from_u32(v as u32) {
                            out.push(ch);
                        }
                    }
                }
                _ => {}
            },
            _ => out.push(c),
        }
    }
    None
}

fn json_u64(src: &str, key: &str) -> Option<u64> {
    let needle = format!("\"{}\"", key);
    let at = src.find(&needle)? + needle.len();
    let rest = &src[at..];
    let colon = rest.find(':')?;
    let rest = rest[colon + 1..].trim_start();
    let digits: String = rest.chars().take_while(|c| c.is_ascii_digit()).collect();
    if digits.is_empty() {
        return None;
    }
    digits.parse().ok()
}

fn valid_version(v: &str) -> bool {
    !v.is_empty()
        && v.len() <= 64
        && v
            .chars()
            .all(|c| c.is_ascii_alphanumeric() || c == '.' || c == '-' || c == '_')
}

fn wide(s: &str) -> Vec<u16> {
    OsStr::new(s).encode_wide().chain(std::iter::once(0)).collect()
}

fn fatal(msg: &str) {
    #[cfg(feature = "console")]
    eprintln!("CG 启动失败: {msg}");
    #[cfg(not(feature = "console"))]
    {
        let text = wide(msg);
        let caption = wide("CG 启动失败");
        unsafe {
            MessageBoxW(0, text.as_ptr(), caption.as_ptr(), 0x10);
        }
    }
    std::process::exit(1);
}
