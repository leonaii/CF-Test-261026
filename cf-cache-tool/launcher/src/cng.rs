use crate::crypto;
use std::ffi::{c_void, OsStr};
use std::os::windows::ffi::OsStrExt;

const ECDSA_P256: &str = "ECDSA_P256";
const ECCPUBLICBLOB: &str = "ECCPUBLICBLOB";
const MAGIC_P256_PUBLIC: u32 = 0x31534345;

#[link(name = "bcrypt")]
extern "system" {
    fn BCryptOpenAlgorithmProvider(
        ph_algorithm: *mut *mut c_void,
        alg_id: *const u16,
        implementation: *const u16,
        flags: u32,
    ) -> i32;
    fn BCryptCloseAlgorithmProvider(algorithm: *mut c_void, flags: u32) -> i32;
    fn BCryptImportKeyPair(
        algorithm: *mut c_void,
        import_key: *mut c_void,
        blob_type: *const u16,
        ph_key: *mut *mut c_void,
        input: *mut u8,
        input_len: u32,
        flags: u32,
    ) -> i32;
    fn BCryptVerifySignature(
        key: *mut c_void,
        padding: *const c_void,
        hash: *const u8,
        hash_len: u32,
        signature: *const u8,
        signature_len: u32,
        flags: u32,
    ) -> i32;
    fn BCryptDestroyKey(key: *mut c_void) -> i32;
}

fn wide(s: &str) -> Vec<u16> {
    OsStr::new(s).encode_wide().chain(std::iter::once(0)).collect()
}

pub fn verify_ecdsa(
    version: &str,
    sha256_hex: &str,
    size: u64,
    public_x_hex: &str,
    public_y_hex: &str,
    signature_b64: &str,
) -> Result<(), String> {
    let sig = crypto::b64_decode(signature_b64).ok_or("清单 ECDSA 签名格式错误")?;
    if sig.len() != 64 {
        return Err("清单 ECDSA 签名长度错误".into());
    }
    let x = crypto::hex_decode(public_x_hex).ok_or("内置公钥 X 无效")?;
    let y = crypto::hex_decode(public_y_hex).ok_or("内置公钥 Y 无效")?;
    if x.len() != 32 || y.len() != 32 {
        return Err("内置公钥长度错误".into());
    }
    let mut blob = Vec::with_capacity(72);
    blob.extend_from_slice(&MAGIC_P256_PUBLIC.to_le_bytes());
    blob.extend_from_slice(&32u32.to_le_bytes());
    blob.extend_from_slice(&x);
    blob.extend_from_slice(&y);

    let payload = format!(
        "{}\n{}\n{}",
        version,
        sha256_hex.trim().to_ascii_lowercase(),
        size
    );
    let digest = crypto::sha256(payload.as_bytes());

    let alg_name = wide(ECDSA_P256);
    let blob_name = wide(ECCPUBLICBLOB);
    let mut alg: *mut c_void = std::ptr::null_mut();
    let status = unsafe { BCryptOpenAlgorithmProvider(&mut alg, alg_name.as_ptr(), std::ptr::null(), 0) };
    if status != 0 {
        return Err(format!("打开 ECDSA_P256 失败（0x{:08x}）", status as u32));
    }
    let mut key: *mut c_void = std::ptr::null_mut();
    let status = unsafe {
        BCryptImportKeyPair(
            alg,
            std::ptr::null_mut(),
            blob_name.as_ptr(),
            &mut key,
            blob.as_mut_ptr(),
            blob.len() as u32,
            0,
        )
    };
    if status != 0 {
        unsafe {
            BCryptCloseAlgorithmProvider(alg, 0);
        }
        return Err(format!("导入 ECDSA 公钥失败（0x{:08x}）", status as u32));
    }
    let status = unsafe {
        BCryptVerifySignature(
            key,
            std::ptr::null(),
            digest.as_ptr(),
            digest.len() as u32,
            sig.as_ptr(),
            sig.len() as u32,
            0,
        )
    };
    unsafe {
        BCryptDestroyKey(key);
        BCryptCloseAlgorithmProvider(alg, 0);
    }
    if status != 0 {
        return Err("清单签名校验失败（可能被篡改）".into());
    }
    Ok(())
}
