use std::ffi::{c_void, OsStr};
use std::os::windows::ffi::OsStrExt;
use std::ptr;

type Handle = *mut c_void;

const ACCESS_NO_PROXY: u32 = 1;
const FLAG_SECURE: u32 = 0x0080_0000;
const QUERY_STATUS: u32 = 19;
const QUERY_NUMBER: u32 = 0x2000_0000;
const OPTION_SECURE_PROTOCOLS: u32 = 84;
const TLS1_2: u32 = 0x0000_0800;
const AUTH_TARGET_SERVER: u32 = 0;
const AUTH_SCHEME_BASIC: u32 = 1;

#[link(name = "winhttp")]
extern "system" {
    fn WinHttpOpen(
        agent: *const u16,
        access_type: u32,
        proxy: *const u16,
        bypass: *const u16,
        flags: u32,
    ) -> Handle;
    fn WinHttpConnect(session: Handle, host: *const u16, port: u16, reserved: u32) -> Handle;
    fn WinHttpOpenRequest(
        connect: Handle,
        verb: *const u16,
        object: *const u16,
        version: *const u16,
        referrer: *const u16,
        accept_types: *const *const u16,
        flags: u32,
    ) -> Handle;
    fn WinHttpSetTimeouts(
        handle: Handle,
        resolve: i32,
        connect: i32,
        send: i32,
        receive: i32,
    ) -> i32;
    fn WinHttpSetOption(handle: Handle, option: u32, buffer: *mut c_void, len: u32) -> i32;
    fn WinHttpSetCredentials(
        request: Handle,
        auth_targets: u32,
        auth_scheme: u32,
        user: *const u16,
        pass: *const u16,
        auth_params: *mut c_void,
    ) -> i32;
    fn WinHttpSendRequest(
        request: Handle,
        headers: *const u16,
        headers_len: u32,
        optional: *mut c_void,
        optional_len: u32,
        total_len: u32,
        context: usize,
    ) -> i32;
    fn WinHttpReceiveResponse(request: Handle, reserved: *mut c_void) -> i32;
    fn WinHttpQueryHeaders(
        request: Handle,
        info_level: u32,
        name: *const u16,
        buffer: *mut c_void,
        buffer_len: *mut u32,
        index: *mut u32,
    ) -> i32;
    fn WinHttpReadData(request: Handle, buffer: *mut c_void, to_read: u32, read: *mut u32) -> i32;
    fn WinHttpCloseHandle(handle: Handle) -> i32;
}

struct Owned(Handle);

impl Drop for Owned {
    fn drop(&mut self) {
        if !self.0.is_null() {
            unsafe {
                WinHttpCloseHandle(self.0);
            }
        }
    }
}

fn wide(s: &str) -> Vec<u16> {
    OsStr::new(s).encode_wide().chain(std::iter::once(0)).collect()
}

fn last_error(what: &str) -> String {
    let code = std::io::Error::last_os_error()
        .raw_os_error()
        .unwrap_or(0) as u32;
    format!("{}失败（WinHTTP 错误 {}）", what, code)
}

pub fn http_get(path: &str, token: &str, pass: &str, limit: usize) -> Result<Vec<u8>, String> {
    let agent = wide("CG-launcher");
    let session = unsafe {
        WinHttpOpen(
            agent.as_ptr(),
            ACCESS_NO_PROXY,
            ptr::null(),
            ptr::null(),
            0,
        )
    };
    if session.is_null() {
        return Err(last_error("初始化网络"));
    }
    let session = Owned(session);
    unsafe {
        WinHttpSetTimeouts(session.0, 5000, 15000, 30000, 120000);
        let mut protocols = TLS1_2;
        let _ = WinHttpSetOption(
            session.0,
            OPTION_SECURE_PROTOCOLS,
            &mut protocols as *mut u32 as *mut c_void,
            4,
        );
    }

    let host = wide("cloud.51i.cc");
    let connect = unsafe { WinHttpConnect(session.0, host.as_ptr(), 443, 0) };
    if connect.is_null() {
        return Err(last_error("连接云端"));
    }
    let connect = Owned(connect);

    let verb = wide("GET");
    let object = wide(path);
    let request = unsafe {
        WinHttpOpenRequest(
            connect.0,
            verb.as_ptr(),
            object.as_ptr(),
            ptr::null(),
            ptr::null(),
            ptr::null(),
            FLAG_SECURE,
        )
    };
    if request.is_null() {
        return Err(last_error("创建请求"));
    }
    let request = Owned(request);

    let user = wide(token);
    let password = wide(pass);
    unsafe {
        let _ = WinHttpSetCredentials(
            request.0,
            AUTH_TARGET_SERVER,
            AUTH_SCHEME_BASIC,
            user.as_ptr(),
            password.as_ptr(),
            ptr::null_mut(),
        );
    }

    let headers = wide(
        "X-Requested-With: XMLHttpRequest\r\nUser-Agent: CG-launcher\r\nAccept: */*\r\n",
    );
    let ok = unsafe {
        WinHttpSendRequest(
            request.0,
            headers.as_ptr(),
            u32::MAX,
            ptr::null_mut(),
            0,
            0,
            0,
        )
    };
    if ok == 0 {
        return Err(last_error("发送请求"));
    }
    let ok = unsafe { WinHttpReceiveResponse(request.0, ptr::null_mut()) };
    if ok == 0 {
        return Err(last_error("接收响应"));
    }

    let mut status: u32 = 0;
    let mut status_len: u32 = 4;
    let ok = unsafe {
        WinHttpQueryHeaders(
            request.0,
            QUERY_STATUS | QUERY_NUMBER,
            ptr::null(),
            &mut status as *mut u32 as *mut c_void,
            &mut status_len,
            ptr::null_mut(),
        )
    };
    if ok == 0 {
        return Err(last_error("读取 HTTP 状态"));
    }
    match status {
        200 => {}
        401 => return Err("只读分享令牌或口令失效（401）".into()),
        404 => return Err("云端文件不存在（404）".into()),
        _ => return Err(format!("云端返回 HTTP {}", status)),
    }

    let mut out: Vec<u8> = Vec::new();
    let mut buf = vec![0u8; 65536];
    loop {
        let mut read: u32 = 0;
        let ok = unsafe {
            WinHttpReadData(
                request.0,
                buf.as_mut_ptr() as *mut c_void,
                buf.len() as u32,
                &mut read,
            )
        };
        if ok == 0 {
            return Err(last_error("读取云端数据"));
        }
        if read == 0 {
            break;
        }
        let read = read as usize;
        if out.len().saturating_add(read) > limit {
            return Err("云端文件超过大小限制".into());
        }
        out.extend_from_slice(&buf[..read]);
    }
    Ok(out)
}
