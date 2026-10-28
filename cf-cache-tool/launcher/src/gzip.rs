pub fn gunzip(data: &[u8], expected: usize) -> Result<Vec<u8>, String> {
    if data.len() < 18 || data[0] != 0x1f || data[1] != 0x8b || data[2] != 8 {
        return Err("主程序压缩包损坏".into());
    }
    let flags = data[3];
    let mut pos = 10usize;
    if flags & 0x04 != 0 {
        if pos + 2 > data.len() {
            return Err("主程序压缩包损坏".into());
        }
        let xlen = u16::from_le_bytes([data[pos], data[pos + 1]]) as usize;
        pos += 2 + xlen;
    }
    if flags & 0x08 != 0 {
        while pos < data.len() && data[pos] != 0 {
            pos += 1;
        }
        pos += 1;
    }
    if flags & 0x10 != 0 {
        while pos < data.len() && data[pos] != 0 {
            pos += 1;
        }
        pos += 1;
    }
    if flags & 0x02 != 0 {
        pos += 2;
    }
    if pos + 8 > data.len() {
        return Err("主程序压缩包损坏".into());
    }
    let deflate = &data[pos..data.len() - 8];
    let out = miniz_oxide::inflate::decompress_to_vec_with_limit(deflate, expected.saturating_add(1))
        .map_err(|_| "主程序解压失败".to_string())?;
    if out.len() != expected {
        return Err(format!("主程序解压大小不一致（{}/{}）", out.len(), expected));
    }
    Ok(out)
}
