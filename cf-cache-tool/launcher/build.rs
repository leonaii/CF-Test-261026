use std::env;
use std::fs;
use std::path::PathBuf;

fn obfuscate(s: &str) -> String {
    let key = b"cgx0r";
    let raw: Vec<u8> = s
        .as_bytes()
        .iter()
        .enumerate()
        .map(|(i, b)| b ^ key[i % key.len()])
        .collect();
    base64_encode(&raw)
}

fn base64_encode(src: &[u8]) -> String {
    const T: &[u8; 64] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    let mut out = String::with_capacity((src.len() + 2) / 3 * 4);
    for chunk in src.chunks(3) {
        let b0 = chunk[0] as u32;
        let b1 = *chunk.get(1).unwrap_or(&0) as u32;
        let b2 = *chunk.get(2).unwrap_or(&0) as u32;
        let n = (b0 << 16) | (b1 << 8) | b2;
        out.push(T[((n >> 18) & 63) as usize] as char);
        out.push(T[((n >> 12) & 63) as usize] as char);
        out.push(if chunk.len() > 1 { T[((n >> 6) & 63) as usize] as char } else { '=' });
        out.push(if chunk.len() > 2 { T[(n & 63) as usize] as char } else { '=' });
    }
    out
}

fn main() {
    for key in ["CG_SHARE_TOKEN", "CG_SHARE_PASS", "CG_LAUNCHER_VERSION"] {
        println!("cargo:rerun-if-env-changed={key}");
    }
    let token = env::var("CG_SHARE_TOKEN").unwrap_or_default();
    let pass = env::var("CG_SHARE_PASS").unwrap_or_default();
    let ver = env::var("CG_LAUNCHER_VERSION").unwrap_or_else(|_| "0.0.0".into());
    let out = PathBuf::from(env::var("OUT_DIR").expect("OUT_DIR"));
    let secret = format!(
        "pub const SHARE_TOKEN_ENC: &str = \"{}\";\n\
         pub const SHARE_PASS_ENC: &str = \"{}\";\n\
         pub const LAUNCHER_VERSION: &str = \"{}\";\n",
        obfuscate(&token),
        obfuscate(&pass),
        ver
    );
    fs::write(out.join("secret.rs"), secret).expect("write secret.rs");

    let rsrc = PathBuf::from(env::var("CARGO_MANIFEST_DIR").expect("manifest dir")).join("rsrc.syso");
    if env::var_os("CG_NO_RSRC").is_none() && rsrc.exists() {
        println!("cargo:rustc-link-arg={}", rsrc.display());
    }
}
