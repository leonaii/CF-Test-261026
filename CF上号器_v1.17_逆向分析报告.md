# CF上号器 v1.17 逆向分析报告

- 分析对象：`CF上号器_v1.17.exe`
- 配套文件：`host.dat`、`config.ini`、`一定要看使用说明.txt`
- 配套方法论：`上号器逆向方法论与版本复用指南.md`
- 分析方式：静态逆向（PE 结构、导入表、字符串、函数反汇编、内嵌资源提取）
- 样本来源：用户本地提供
- 分析日期：2026-10-22
- 分析结论：该程序不是普通的游戏启动器，而是一个“卡密平台 + QQ/WeGame 自动登录 + 伪造 Rail 启动链 + DLL 注入”的完整上号客户端

---

## 1. 样本指纹

### 1.1 文件信息

| 文件 | 大小 | 类型 | SHA-256 |
|---|---:|---|---|
| `CF上号器_v1.17.exe` | 1,597,952 字节 | PE32+ executable (GUI) x86-64 | `2d855012ee4ed2a8df6ecb034af2c51bbd395207fa338b94eb94512c4daae878` |
| `host.dat` | 103 字节 | data（自定义加密容器） | `f72f8d0317756365e0b6e740595091915a9c509b6b570df4f5c151dd97777894` |
| `config.ini` | 80 字节 | INI 文本 | `9c767f953892d1c597d7e2738192ade069041f8fd7bcc36b35a2ae5a93b321f1` |
| `一定要看使用说明.txt` | 356 字节 | UTF-8 文本 | `e542b4f594f41c873080ff8ee79f41812c55d34d3d70526c77376d86ec23ed00` |

### 1.2 可执行文件哈希

```text
MD5    83773f545f3e5c63227c145bdca1096e
SHA256 2d855012ee4ed2a8df6ecb034af2c51bbd395207fa338b94eb94512c4daae878
```

### 1.3 基础属性

| 属性 | 值 |
|---|---|
| 架构 | x86-64 |
| 子系统 | Windows GUI |
| 编译语言 | C++（MinGW/GCC，存在 libstdc++ 错误字符串与 `.eh_frame`） |
| 镜像基址 | `0x140000000` |
| 入口点 | 由 PE 可选头指定 |
| 加壳 | 未发现常见壳（UPX 等），字符串可直接读取 |
| 数字签名 | 无（Security Directory 为 0） |
| 版本资源 | 无 `RT_VERSION` |
| 资源 | 图标 + 两份 Manifest |
| 权限清单 | 同时存在 `requireAdministrator` 与 `asInvoker` 两份，配置不规范；从实际行为看，注入和写游戏目录通常需要管理员权限 |

---

## 2. 结论摘要

该程序的核心链路如下：

1. 读取 `config.ini` 和注册表，定位 CF 安装目录。
2. 校验卡密格式（16 位纯英文）。
3. 通过 WinHTTP 调用卡密平台的 `/api/card/boot`、`/api/card/peek`、`/api/card/redeem`、`/api/card/log` 接口。
4. 卡密平台下发两类账号数据：
   - OP 数据：`openid`、`open_key`、`client_key`、`ctt` 等；
   - NT 缓存：`Token010A`、`Token0143`、`sessionKey` 等。
5. 如果是 NT 缓存，程序使用内置的 QQ NT 登录协议（ECDH P-256、TEA、AES-GCM、SSO over TCP）完成登录，换取 `clientkey`、`skey`、`openkey`。
6. 将 QQ 登录态转换为 WeGame OP 登录态（`openid`、`open_key`、`client_key`、`game_ticket`、`tgp_id`），也就是说明书中说的“转 OP”。
7. 伪造 WeGame Rail 运行环境：创建共享内存、命名管道，写入/使用 `mmog_data.xml`、`rail_user_data`、`version.ini` 等。
8. 直接启动 `crossfire.exe`，并通过 `VirtualAllocEx + WriteProcessMemory + CreateRemoteThread(LoadLibraryW)` 注入内嵌的 `qwh_cf_hook.dll`。
9. 等待 Rail 管道和 CShell，向卡密平台 `/api/card/log` 上报启动结果。

一句话概括：

> 用卡密从平台换取账号票据，用内置 QQ NT 登录把票据换成 WeGame OP 登录态，伪造 Rail 启动链直接拉起 `crossfire.exe`，再注入 Hook DLL 让游戏接受这套假启动环境。

---

## 3. 样本构成

| 文件 | 作用 |
|---|---|
| `CF上号器_v1.17.exe` | 主程序，包含 UI、卡密协议、QQ NT 登录、WeGame 授权、Rail 启动、DLL 注入 |
| `host.dat` | 自定义加密配置，解密后作为卡密服务器地址使用 |
| `config.ini` | 本地配置：卡密、游戏目录、区服号、自动选区 |
| `一定要看使用说明.txt` | 使用说明，包含“转OP失败就是死号”等关键提示 |

### 3.1 `config.ini` 字段

```ini
[ui]
card=
dir=D:\WeGameApps\rail_apps\CF：穿越火线(2)
ssn=333
auto_zone=0
```

| 字段 | 含义 |
|---|---|
| `card` | 卡密，程序要求 16 位纯英文 |
| `dir` | CF 游戏目录 |
| `ssn` | 区服号，会拼进启动参数 `-SSN:<ssn>` |
| `auto_zone` | 是否自动选区 |

### 3.2 `host.dat` 结构

文件头：

```text
d8 3a 91 c7 02 ...
```

已知特征：

| 偏移 | 内容 | 说明 |
|---:|---|---|
| `0x00` | `d8 3a 91 c7` | 魔数 |
| `0x04` | `02` | 版本 |
| `0x05` | 16 字节 | 初始化向量/盐 |
| `0x15` | `uint16` | 长度字段（样本中为 32） |
| `0x17` 起 | 密文 + HMAC | 自定义加密数据 |

解密相关代码：

- 读取与解析：RVA `0x90130`
- 解密入口：RVA `0x8F550`
- 密钥派生：RVA `0x8E300`，使用 `SHA256` 和两个标签：
  - `xor-host-t-v2`
  - `mac-host-t-v2`
- 完整性校验：RVA `0x8E140`，HMAC-SHA256

如果 `host.dat` 缺失、损坏或校验失败，程序会回退到：

```text
http://203.0.113.10:8788
```

`203.0.113.0/24` 属于 TEST-NET-3 文档保留地址，不是公网真实地址。这说明：

- 二进制中的该地址是开发占位/兜底；
- 真实卡密服务器地址被加密在 `host.dat` 中；
- 替换或删除 `host.dat` 会导致程序连不上卡密服务。

---

## 4. 导入表与能力画像

### 4.1 网络与 HTTP

- `WINHTTP.dll`
  - `WinHttpOpen`
  - `WinHttpConnect`
  - `WinHttpOpenRequest`
  - `WinHttpSendRequest`
  - `WinHttpReceiveResponse`
  - `WinHttpReadData`
  - `WinHttpQueryHeaders`
  - `WinHttpCrackUrl`
- `WS2_32.dll`
  - `socket`、`connect`、`send`、`recv`
  - 用于 QQ NT SSO 的原始 TCP 连接
- 支持 `http://`、`https://`、`socks5://`

### 4.2 加密

- `bcrypt.dll`
  - ECDH P-256：`BCryptSecretAgreement`、`BCryptDeriveKey`
  - AES-GCM / AES-CBC：`BCryptEncrypt`、`BCryptDecrypt`
  - SHA256：`BCryptCreateHash`、`BCryptHashData`
  - 随机数：`BCryptGenRandom`

### 4.3 进程与注入

- `CreateProcessW`
- `OpenProcess`
- `ReadProcessMemory`
- `VirtualAllocEx`
- `WriteProcessMemory`
- `CreateRemoteThread`
- `VirtualProtectEx`
- `CreateToolhelp32Snapshot`
- `Process32FirstW` / `Process32NextW`

### 4.4 游戏/WeGame 交互

- 命名管道：
  - `CreateNamedPipeW`
  - `ConnectNamedPipe`
  - `\\.\pipe\rail_pipe_%u`
  - `\\.\pipe\rail_named_pipe_%u`
  - `\\.\pipe\rail_named_pipe_host_%u`
  - `\\.\pipe\rail_named_pipe_guest_%u`
- 共享内存：
  - `CreateFileMappingW`
  - `MapViewOfFile`
  - `Rail.SharedMemory.%u`
- 窗口：
  - `FindWindowW`
  - `EnumWindows`
  - `PostMessageW`
  - `SendMessageW`
  - `WE_GAME_SDK_WINDOW_NAME_CLASS`
  - `rail_dev_wnd`
  - `_TGP_MAIN_WINDOW_WND_NAME_`

### 4.5 本地文件与注册表

- `GetPrivateProfileStringW` / `WritePrivateProfileStringW`
- `CreateFileW` / `ReadFile` / `WriteFile`
- `RegOpenKeyExW` / `RegQueryValueExW`
- 注册表路径：
  - `SOFTWARE\Tencent\CrossFire`
  - `SOFTWARE\WOW6432Node\Tencent\CrossFire`
  - 值名：`InstallPath`、`default_path`

---

## 5. 总体架构

```text
+---------------------------------------------------------------+
|                         CF上号器 v1.17                        |
+---------------------------------------------------------------+
|  UI 层（Win32 / Microsoft YaHei UI / STATIC / EDIT / BUTTON）  |
|  - 卡密输入、游戏目录、区服选择、自动选区、状态显示             |
+-------------------------------+-------------------------------+
                                |
+-------------------------------v-------------------------------+
|  卡密协议层（DfCard/1.0）                                      |
|  - /api/card/boot   /api/card/peek                            |
|  - /api/card/redeem /api/card/log                             |
|  - RSA + AES-GCM + Base64                                     |
+-------------------------------+-------------------------------+
                                |
+-------------------------------v-------------------------------+
|  账号数据层                                                    |
|  - OP：openid/open_key/client_key/ctt                         |
|  - NT：Token010A/Token0143/sessionKey                         |
|  - 本地缓存：op.json / last_op.json / Token010A.txt / nt.ini   |
+-------------------------------+-------------------------------+
                                |
+-------------------------------v-------------------------------+
|  QQ NT 登录层                                                  |
|  - ECDH P-256（BCrypt）                                       |
|  - TEA / AES-GCM / TLV                                        |
|  - SsoKeyExchange                                             |
|  - SsoNTLoginEasyLogin / SsoNTLoginPasswordLogin              |
|  - SsoGetClientKey                                            |
+-------------------------------+-------------------------------+
                                |
+-------------------------------v-------------------------------+
|  WeGame OP 转换层（“转 OP”）                                   |
|  - UnionAuth/LoginByQq                                        |
|  - Auth/GetGameToken                                          |
|  - TicketProxy/GetGameTicket                                  |
|  - GameAdapter/GetGameLaunchInfo                              |
+-------------------------------+-------------------------------+
                                |
+-------------------------------v-------------------------------+
|  Rail 启动层                                                   |
|  - Rail.SharedMemory.%u                                       |
|  - rail_pipe_15 / rail_named_pipe_15 / host / guest           |
|  - mmog_data.xml / rail_user_data / version.ini               |
|  - crossfire.exe                                              |
+-------------------------------+-------------------------------+
                                |
+-------------------------------v-------------------------------+
|  DLL 注入层                                                    |
|  - qwh_cf_hook.dll（内嵌 101,888 字节）                        |
|  - VirtualAllocEx + WriteProcessMemory                        |
|  - CreateRemoteThread(LoadLibraryW)                           |
+---------------------------------------------------------------+
```

---

## 6. 启动与上号主流程

### 6.1 主流程

```text
1. 双击程序
2. 读取 config.ini
3. 读取注册表 SOFTWARE\Tencent\CrossFire\InstallPath
4. 校验卡密格式（16 位纯英文）
5. POST /api/card/boot      激活卡密、开始计时、更新缓存
6. POST /api/card/peek      查询卡密状态、到期时间、区服、备注
7. POST /api/card/redeem    领取账号数据
8. 判断返回类型：
   A. OP 数据
      -> 直接写 op.json / last_op.json
   B. NT 缓存
      -> 走 QQ NT 登录
      -> 换取 clientkey/skey/openkey
      -> 转 WeGame OP
      -> 写 op.json / last_cf_op.json
9. 读取 ssn / auto_zone，确定区服
10. 构造 Rail 运行环境
11. 启动 crossfire.exe
12. 注入 qwh_cf_hook.dll
13. 等待 Rail 管道和 CShell
14. POST /api/card/log 上报结果
```

### 6.2 流水线函数

主流程位于 RVA `0x7220` 的大函数中，调用顺序可以概括为：

| 顺序 | 调用/RVA | 作用 |
|---:|---|---|
| 1 | `0x8B1E0` | 卡密 peek |
| 2 | `0x8A0D0` | 卡密 log |
| 3 | `0x8C140` | 卡密 redeem |
| 4 | `0x2110` | 读取本地缓存（op.json/Token010A.txt/nt.ini 等） |
| 5 | `0x22DD0` | NT 缓存转 CF OP |
| 6 | `0x1AEE0` | 伪造 Rail + 启动 + 注入 |
| 7 | `0x1D540` | 等待管道/收尾 |
| 8 | `0x8A0D0` | 再次上报 log |

---

## 7. 模块详解

### 7.1 配置与目录探测

程序会：

- 读取 `config.ini` 的 `card`、`dir`、`ssn`、`auto_zone`；
- 读取注册表 `SOFTWARE\Tencent\CrossFire` 和 `SOFTWARE\WOW6432Node\Tencent\CrossFire`；
- 读取 `InstallPath`、`default_path`；
- 检查 `crossfire.exe`、`crossfire_x64.exe`、`x64\crossfire.exe`；
- 检查 `rail_api64.dll`、`rail_sdk_tgp_platform.dll`、`rail_sdk_platform.dll`、`rail_api.dll`；
- 搜索常见目录：
  - `D:\WeGameApps\rail_apps\CF`
  - `C:\WeGameApps\rail_apps\CF`
  - `D:\Program Files (x86)\WeGame`
  - `G:\Program Files (x86)\WeGame`
  - `C:\Users\Administrator\Desktop\WeGame`

### 7.2 卡密协议

#### 7.2.1 基本参数

| 项目 | 值 |
|---|---|
| User-Agent | `DfCard/1.0` |
| 方法 | `POST` |
| Content-Type | `application/json` |
| 传输 | WinHTTP |
| 加密 | RSA + AES-GCM + Base64 |
| 卡密长度 | 16 位纯英文 |

#### 7.2.2 接口

| 接口 | 作用 |
|---|---|
| `/api/card/boot` | 激活卡密、开始计时、更新最新缓存 |
| `/api/card/peek` | 查询卡密状态、到期时间、有效小时、区服、备注 |
| `/api/card/redeem` | 领取账号数据（OP 或 NT 缓存） |
| `/api/card/log` | 上报事件、QQ、卡密、启动结果 |
| `text-upload` | 文本上传 |
| `agent-upload` | 代理/客户端信息上传 |

#### 7.2.3 返回字段

- `expires_at`
- `valid_hours`
- `activated`
- `zone`
- `note`
- `balance_cents`
- `available_cents`
- `reserved_cents`

#### 7.2.4 关键提示

| 字符串 | 含义 |
|---|---|
| `卡密需为16位纯英文` | 卡密格式校验 |
| `卡密无效` | 卡密不存在/已过期 |
| `兑换失败` | redeem 失败 |
| `服务端仍在下发OP` | 服务端仍返回 OP 数据 |
| `解密内容不是 NT 缓存` | 返回数据不是预期的 NT 缓存 |
| `兑换成功` | redeem 成功 |
| `缓存转 CF OP 失败` | NT 缓存转 OP 失败，即“转OP失败” |

### 7.3 QQ NT 登录

程序内置了完整的 QQ NT 登录协议实现，相关字符串与 API：

#### 7.3.1 协议命令

- `trpc.login.ecdh.EcdhService.SsoKeyExchange`
- `trpc.login.ecdh.EcdhService.SsoNTLoginEasyLogin`
- `trpc.login.ecdh.EcdhService.SsoNTLoginPasswordLogin`
- `trpc.login.ecdh.EcdhService.SsoGetClientKey`
- `ConfigPushSvc.GehSvc.GetIpDirect`

#### 7.3.2 加密与编码

- ECDH P-256
- SHA256
- AES-GCM
- AES-CBC
- TEA
- TLV
- Base64

#### 7.3.3 票据与缓存

- `Token010A`
- `Token0143`
- `sessionKey`
- `tgt`
- `d2`
- `d2_key`
- `client_key`
- `st_key`
- `skey`
- `p_skey`
- `pt4_token`

本地缓存文件：

- `Token010A.txt`
- `nt.ini`
- `last_nt.ini`
- `data.txt`
- `cache.ini`

#### 7.3.4 登录分支

- `SsoNTLoginEasyLogin`：使用缓存票据快速登录
- `SsoNTLoginPasswordLogin`：密码登录分支
- `QQConnectLogin.getAccessToken`
- `QQConnectLogin.getClientKey`
- `QQConnectLogin.getSkey`
- `QQConnectLogin.getOpenKey`
- `QQConnectLogin.getST`
- `QQConnectLogin.getTicket`

#### 7.3.5 Web 票据

程序还保留了 QQ 网页登录/快速登录链路：

- `https://ptlogin2.qq.com/jump?clientuin=...&clientkey=...&keyindex=18&daid=8&u1=http://x5.qq.com/main.shtml`
- `https://accounts.qq.com/safe/connect/captcha?aid=2009945610`
- `https://configsvr.msf.3g.qq.com/configsvr/serverlist.jsp?mType=getssolist`
- `msfwifi.3g.qq.com`
- Cookie 字段：
  - `skey`
  - `p_skey`
  - `pt4_token`
  - `uin=o...; skey=...; p_skey=...; p_uin=o...; pt4_token=...`

### 7.4 WeGame OP 转换（“转 OP”）

#### 7.4.1 接口

- `https://www.wegame.com.cn/api/v2/wegame.auth.game.UnionAuth/LoginByQq`
- `.../wegame.auth.Auth/GetGameToken`
- `.../wegame.platform.game.TicketProxy/GetGameTicket`
- `.../wegame.gameadapter.GameAdapter/GetGameLaunchInfo`
- `.../UnionAuth/RefreshClientTicket`

#### 7.4.2 请求参数

- `appid = 2112513`
- `mappid = 10001`
- `clienttype = 1000000`
- `qq_info_type = 4/6`
- `openid`
- `client_key`
- `sig`
- `st_key`
- `ct`
- `ctt`
- `machine_id`
- `version_id`

#### 7.4.3 返回字段

- `openid`
- `open_key`
- `client_key`
- `game_ticket`
- `tgp_id`
- `biz_ticket`
- `game_st`
- `game_st_key`
- `game_client_key`
- `trace_id`
- `sdk_traceid`
- `expire_interval`
- `is_real_name`
- `adult_type`

#### 7.4.4 失败含义

二进制中的提示：

```text
缓存转 CF OP 失败
OP字段无效(ctt=...,openid=...,key=...)
data 缺少 Token010A/0143/sessionKey
```

对应说明书：

```text
如果提示转OP失败那就是百分百彻底死号了
```

也就是说，NT 票据无法换成 WeGame OP 登录态时，账号基本不可用。

### 7.5 Rail 启动链

程序围绕 WeGame Rail 运行环境做了大量伪造：

#### 7.5.1 共享内存与管道

- `Rail.SharedMemory.%u`
- `\\.\pipe\rail_pipe_%u`
- `\\.\pipe\rail_named_pipe_%u`
- `\\.\pipe\rail_named_pipe_host_%u`
- `\\.\pipe\rail_named_pipe_guest_%u`
- 样本中实际使用：
  - `\\.\pipe\rail_pipe_15`
  - `\\.\pipe\rail_named_pipe_15`
  - `\\.\pipe\rail_named_pipe_host_15`
  - `\\.\pipe\rail_named_pipe_guest_15`

#### 7.5.2 文件

- `\TCLS\mmog_data.xml`
- `\NTCLS\mmog_data.xml`
- `\NTCLS\ntcls_config\mmog_data.xml`
- `\NTCLS\tgp_config\mmog_data.xml`
- `rail_user_data\2`
- `version.ini`
- `last_cf_rail.bin`
- `last_cf_pipe.txt`

#### 7.5.3 启动参数

- `-SSN:318`
- `-SSN:333`
- `-cf:1`
- `-HostIP=`
- `-q`
- `-q --disable-gpu --disable-gpu-compositing`

#### 7.5.4 错误提示

- `假Rail`
- `WeGame Init Failed. (2310-x)`
- `WeGame Init Failed. (2311)`
- `WeGame Init Failed 2310/2311:管道未连上`
- `闪崩:管道未连上`
- `闪崩:CShell退出`
- `没有该区服Rail模板`
- `Rail模板坏了`
- `Rail模板缺启动包`
- `WeGame共享内存被占用,请退出 WeGame`
- `请先退出 WeGame 再上穿越火线`
- `环境异常`
- `没有QQ`
- `OP缺openid`
- `找不到rail_api64.dll`

### 7.6 DLL 注入

#### 7.6.1 内嵌 DLL

- 名称：`qwh_cf_hook.dll`
- 大小：101,888 字节（`0x18E00`）
- 位置：主程序 `.rdata` 文件偏移 `0xDB8E0`
- 导出：无
- 导入：仅 `KERNEL32.dll`
- 关键字符串：
  - `LdrRegisterDllNotification`
  - `RailNeedRestartAppForCheckingEnvironment`

从字符串和导入看，它是一个运行时解析 API 的 Hook DLL，主要目标是 `rail_api64.dll` 的加载与环境检查。

#### 7.6.2 落盘与加载

主程序在 RVA `0x1AEE0` 的启动函数中：

1. 获取自身路径；
2. 拼接 `qwh_cf_hook.dll`；
3. 使用 `_wfopen(..., L"wb")` 写文件；
4. `fwrite` 写入 `0x18E00` 字节；
5. `LoadLibrary` 加载。

#### 7.6.3 远程注入

注入调用点：

| 调用 | RVA |
|---|---|
| `VirtualAllocEx` | `0x1CFD8` |
| `WriteProcessMemory` | `0x1D004` |
| `CreateRemoteThread` | `0x1D058` |
| `GetProcAddress("LoadLibraryW")` | `0x1D029` |

注入流程：

```text
VirtualAllocEx(目标进程, NULL, 路径长度*2, MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)
WriteProcessMemory(目标进程, 远程地址, DLL路径, 路径长度*2, NULL)
CreateRemoteThread(目标进程, NULL, 0, LoadLibraryW, 远程地址, 0, NULL)
WaitForSingleObject(线程, 2000)
CloseHandle(线程)
```

目标进程是同一函数中通过 `CreateProcessW` 创建的进程，启动路径指向 `crossfire.exe`。

### 7.7 区服与换区

程序内置 42 个 CF 区服：

```text
广东一区、广东二区、广东三区、广东四区
河南一区、河南二区
山东一区、山东二区
湖北一区、湖北二区
四川一区、四川二区
江苏一区、江苏二区
辽宁一区、辽宁二区、辽宁三区
北京一区、北京二区、北京三区、北京四区
上海一区、上海二区
浙江一区、浙江二区
湖南一区、湖南二区
陕西一区
广西一区
南方大区、北方大区
福建一区
江西一区
重庆一区
安徽一区
云南一区
黑龙江区
吉林一区
山西一区
河北一区
移动专区
教育网区
```

- `ssn` 会拼进 `-SSN:<ssn>`；
- `auto_zone` 控制是否自动选区；
- 说明书中的“最右边的自动左边不要勾选就可以换区创角色”，对应自动选区逻辑。

### 7.8 风控、打码、代理与信用查询

程序集成了多类第三方服务：

#### 7.8.1 QQ 信用分/信誉查询

- `http://www.budqq.com/`
- `http://38.47.219.48/`
- `http://38.47.219.53/`
- `api/1/v1/see?token=`
- `api/1/v1/query?token=`
- `score=`

对应 UI 字段：

- `信用分`
- `信誉`

#### 7.8.2 打码/验证码平台

- `http://www.xmdm.pro`
- `http://103.236.98.220:64398`
- `/api/v1/user/balance`
- `/api/v1/user/captcha`
- `/api/v1/user/tasks/`
- `/api/v1/user/tasks/upload`
- `/api/order/add`
- `/api/order/get?apikey=`
- `/api/getBalance?apikey=`

字段/标识：

- `Qx Token`
- `XM ApiKey`
- `XM aid`
- `task_id`
- `status`

#### 7.8.3 代理

- `socks5://`
- `SOCKS5://`
- `CONNECT`
- `Proxy-Authorization: Basic`
- `api/1/v1/add_proxy`
- `&proxy=`
- `&proxy_username=`
- `&proxy_password=`

#### 7.8.4 登录验证提示

- `验证`
- `短信`
- `手机`
- `设备环境`
- `is_use_mobile`
- `mobile_number_`
- `down_sms_flag_`
- `wg_safe_chk_`
- `mibao_cmd`
- `mibao_rsp`
- `wx_offical_bind`

### 7.9 UI 层

- 字体：`Microsoft YaHei UI`
- 控件类型：
  - `STATIC`
  - `EDIT`
  - `BUTTON`
  - `COMBOBOX`
- 状态字段：
  - `信用分`
  - `信誉`
  - `验证`
  - `短信`
  - `手机`
  - `设备环境`
  - `满单`
  - `余额`
  - `积分`
  - `有误`
- 关键提示：
  - `读卡密`
  - `正在读取卡密`
  - `请输入卡密`
  - `请选择游戏根目录`
  - `先关掉其他上号器`
  - `客户端完整性校验失败`
- 版本号：
  - `v1.17`

### 7.10 调试与自测痕迹

样本中保留了明显的开发/测试痕迹：

- 构建路径：
  - `D:\re_embed\nt_login_build\last_ep.json`
  - `D:\re_embed\nt_login_build\last_mid.txt`
  - `D:\re_embed\nt_login_build\launch_op_win.json`
  - `D:\re_embed\nt_login_build\web_ticket_win.json`
  - `D:\re_embed\nt_login_build\webticket_probe.json`
- 自测模式：
  - `--self-test`
  - `http://127.0.0.1:8790`
  - `NtSelfTest`
- 硬编码测试数据：
  - `{"qq":"123","openid":"0123456789ABCDEF0123456789ABCDEF","open_key":"FEDCBA9876543210FEDCBA9876546543210","ctt":"...`
  - `{"qq_info":{"user_id":"85488237"},"user_id":"u_wtjwweIVUptJkvjDMrpR8A","tgpid":85488237}`
  - `{"mappid":10001,"clienttype":1000000,...}`

这些痕迹说明该样本可能是开发版/测试版，正式版行为可能不同。

---

## 8. 四种启动策略

程序硬编码了四个启动策略名：

| 策略 | 推断作用 | 证据 |
|---|---|---|
| `CfLaunch` | 走官方 `rail_api64.dll` / Rail SDK 启动 | 字符串 `CfLaunch`、`rail_api64.dll` |
| `QwhCfOpLaunch` | 伪造 Rail + OP 数据启动 | 字符串 `QwhCfOpLaunch`、`假Rail`、共享内存/管道 |
| `QwhNtOpLaunch` | 通过 `ntop_v17` 模块把 NT 票据转 OP 后启动 | `ntop_v17_ini`、`ntop_v17_dll\x64`、`-SSN:318`、`-cf:1` |
| `DpCfWin32Launch` | 直接 Win32 拉起 `crossfire.exe` | 字符串 `DpCfWin32Launch`、`x64\crossfire.exe` |

### 8.1 `QwhCfOpLaunch` / 伪造 Rail 路径

对应 RVA `0x1AEE0` 的大函数，主要行为：

- 检查 QQ、OP `openid`、`x64\crossfire.exe`、`rail_api64.dll`；
- 检查 WeGame 是否在运行；
- 读取/构造 `mmog_data.xml`；
- 创建 `Rail.SharedMemory.15`；
- 创建 `rail_pipe_15` 等命名管道；
- 写 `rail_user_data\2`；
- `CreateProcessW` 启动 `crossfire.exe`；
- 注入 `qwh_cf_hook.dll`；
- 等待管道和 CShell。

### 8.2 `QwhNtOpLaunch` / ntop 路径

对应 RVA `0x19270` 的函数，主要行为：

- 读取 `ntop_v17_ini`；
- 使用 `ntop_v17_dll\x64`；
- 检查 `rail_api64.dll` 和 `crossfire.exe`；
- 启动参数：
  - `-SSN:318`
  - `-cf:1`
  - `-HostIP=`
  - `-q`
- 错误提示：
  - `WeGame Init Failed. (2310-x)`
  - `WeGame Init Failed. (2311)`
  - `假Rail`

---

## 9. 关键函数与地址

> 地址均为 RVA，镜像基址 `0x140000000`。文件偏移 = RVA - 节 VA + 节 RAW。

| RVA | 作用 |
|---:|---|
| `0x30D0` | 主 UI/消息处理，读取配置、注册表、启动流程入口 |
| `0x6030` | “上号”按钮处理/启动方式选择，检查卡密、目录、实例、版本 |
| `0x7220` | 上号主流水线：peek -> redeem -> 缓存 -> 转 OP -> 启动 -> 注入 |
| `0x896F0` | WinHTTP 通用请求封装（UA `DfCard/1.0`） |
| `0x89E30` | `/api/card/boot` |
| `0x8A0D0` | `/api/card/log` |
| `0x8B1E0` | `/api/card/peek` |
| `0x8C140` | `/api/card/redeem` |
| `0x90130` | 读取/解析 `host.dat` |
| `0x8F550` | `host.dat` 解密入口 |
| `0x8E300` | `host.dat` 密钥派生（HMAC-SHA256，标签 `xor-host-t-v2` / `mac-host-t-v2`） |
| `0x8E140` | HMAC-SHA256 完整性校验 |
| `0x2110` | 读取本地缓存：`op.json`、`last_op.json`、`Token010A.txt`、`last_cf_op.json`、`nt.ini`、`data.txt` |
| `0x22DD0` | NT 缓存转 CF OP（“转 OP”核心） |
| `0x1AEE0` | 伪造 Rail + `CreateProcessW` + 注入 `qwh_cf_hook.dll` |
| `0x19270` | `ntop_v17` 启动路径 |
| `0x1CFD8` | `VirtualAllocEx` 调用点 |
| `0x1D004` | `WriteProcessMemory` 调用点 |
| `0x1D058` | `CreateRemoteThread` 调用点 |
| `0x1D029` | `GetProcAddress("LoadLibraryW")` 调用点 |
| `0x1CAEB` | `CreateProcessW` 调用点 |

---

## 10. 关键字符串证据

> RVA 对应主程序虚拟地址 `0x140000000 + RVA`。

| RVA | 字符串 | 说明 |
|---:|---|---|
| `0xFB7A0` | `host.dat` | 加密配置文件 |
| `0xFAC22` | `DfCard/1.0` | 卡密客户端 UA |
| `0xFACC6` | `/api/card/boot` | 卡密激活/启动 |
| `0xFAD18` | `/api/card/log` | 卡密日志上报 |
| `0xFAD7C` | `/api/card/peek` | 卡密状态查询 |
| `0xFAE28` | `/api/card/redeem` | 领取账号数据 |
| `0xFAD66` | `卡密需为16位纯英文` | 卡密格式校验 |
| `0xFADE8` | `服务端仍在下发OP` | redeem 返回 OP |
| `0xFAE0E` | `解密内容不是 NT 缓存` | redeem 返回类型异常 |
| `0xFAE90` | `兑换成功` | redeem 成功 |
| `0xF6CBC` | `缓存转 CF OP 失败` | 转 OP 失败 |
| `0xDC7AA` | `ntop_v17_ini` | ntop 配置 |
| `0xDCF28` | `qwh_cf_hook.dll` | 内嵌 Hook DLL |
| `0xDC99A` | `假Rail` | 伪造 Rail |
| `0xDCCA8` | `Rail.SharedMemory.%u` | Rail 共享内存 |
| `0xDCDC0` | `\\.\pipe\rail_pipe_15` | Rail 管道 |
| `0xDCAEC` | `没有QQ` | 启动前检查 |
| `0xDCB0C` | `OP缺openid` | OP 数据不完整 |
| `0xDCB74` | `环境异常` | 环境检查 |
| `0xDCB90` | `请先退出 WeGame 再上穿越火线` | 冲突检查 |
| `0xDCCE0` | `WeGame共享内存被占用,请退出 WeGame` | 共享内存冲突 |
| `0xDCF7C` | `闪崩:管道未连上` | 管道未连上 |
| `0xDCF8E` | `闪崩:CShell退出` | CShell 退出 |
| `0xDCFAE` | `已拉起游戏,管道等待中` | 启动成功 |
| `0xDCFC8` | `WeGame Init Failed 2310/2311:管道未连上` | WeGame 初始化失败 |
| `0xF6F90` | `trpc.login.ecdh.EcdhService.SsoKeyExchange` | NT 登录密钥交换 |
| `0xF6FF0` | `trpc.login.ecdh.EcdhService.SsoNTLoginEasyLogin` | NT 快速登录 |
| `0xF71B0` | `trpc.login.ecdh.EcdhService.SsoNTLoginPasswordLogin` | NT 密码登录 |
| `0xF7728` | `trpc.login.ecdh.EcdhService.SsoGetClientKey` | 获取 clientkey |
| `0xF5BB0` | `/api/v2/wegame.auth.Auth/GetGameToken` | WeGame 取 token |
| `0xF5C00` | `/api/v1/wegame.gameadapter.GameAdapter/GetGameLaunchInfo` | WeGame 启动信息 |
| `0xF61E0` | `/api/v2/wegame.auth.game.UnionAuth/RefreshClientTicket` | 刷新票据 |
| `0xF6E50` | `/api/v1/wegame.platform.game.TicketProxy/GetGameTicket` | 获取游戏票据 |
| `0xF8428` | `http://www.budqq.com/` | 信用分查询 |
| `0xF8608` | `/api/order/add` | 订单接口 |
| `0xF8650` | `http://www.xmdm.pro` | 打码/订单平台 |
| `0xF877E` | `http://103.236.98.220:64398/api/v1/user/balance?user_token=usr_BAD` | 打码平台余额 |
| `0xF8C50` | `/api/v1/user/captcha` | 验证码接口 |
| `0xF8C90` | `/api/v1/user/tasks/` | 任务接口 |
| `0xF8CD8` | `/api/v1/user/tasks/upload` | 任务上传 |
| `0xF8D38` | `api/1/v1/add_proxy` | 添加代理 |
| `0xF7978` | `https://accounts.qq.com/safe/connect/captcha?aid=2009945610` | QQ 验证码 |
| `0xF6BE4` | `msfwifi.3g.qq.com` | QQ 登录服务器 |
| `0xF6B50` | `https://configsvr.msf.3g.qq.com/configsvr/serverlist.jsp?mType=getssolist` | QQ 服务器列表 |

---

## 11. 行为特征与风险

### 11.1 高风险行为

1. **管理员权限**
   - Manifest 包含 `requireAdministrator`；
   - 同时存在 `asInvoker`，配置不规范；
   - 注入、写游戏目录、创建全局命名管道通常需要管理员权限。

2. **无数字签名**
   - Security Directory 为 0；
   - 无法验证发布者；
   - 容易被篡改/二次打包。

3. **注入游戏进程**
   - `VirtualAllocEx`
   - `WriteProcessMemory`
   - `CreateRemoteThread(LoadLibraryW)`
   - 注入 `qwh_cf_hook.dll`

4. **修改/伪造游戏与 WeGame 数据**
   - `mmog_data.xml`
   - `rail_user_data`
   - `version.ini`
   - `Rail.SharedMemory.15`
   - `rail_pipe_15` 等命名管道

5. **绕过官方启动链**
   - 不走 WeGame 正常登录；
   - 直接拉起 `crossfire.exe`；
   - 伪造 Rail 环境；
   - 这类方式容易被 CF 的 CShell/反作弊判定为异常环境，存在封号风险。

6. **账号票据外发/落盘**
   - 卡密、QQ、启动结果上报 `/api/card/log`；
   - `op.json`、`Token010A.txt`、`nt.ini` 等票据缓存到本地；
   - 真实服务器地址藏在加密的 `host.dat` 中；
   - 从静态看，程序具备收集并外发账号票据的能力。

7. **第三方打码/代理**
   - `xmdm.pro`
   - `103.236.98.220:64398`
   - `socks5://`
   - 说明它考虑了对批量登录、验证码、风控的处理。

8. **调试/测试痕迹**
   - `D:\re_embed\nt_login_build\...`
   - `--self-test`
   - `127.0.0.1:8790`
   - 硬编码测试 `openid`/`qq`
   - 说明该样本可能是开发版/测试版，正式版可能不同。

### 11.2 安全建议

- 不要在存有高价值账号的主力机上运行。
- 如需分析，建议放在虚拟机/沙箱中，断网或限制出网。
- 不要使用真实卡密、真实 QQ 账号测试。
- 不要在生产环境替换 `host.dat` 或修改程序文件。
- 如果必须运行，先备份游戏目录和 WeGame 配置。
- 关注游戏账号是否出现异常登录、封禁、资产变动。

---

## 12. 分析限制

- 本次分析以静态逆向为主，未在 Windows 环境中动态运行样本。
- `host.dat` 的自定义加密已完整还原（见第 15 节），真实卡密服务器地址为 `http://111.170.141.190:8790`。
- 部分启动策略（`CfLaunch`、`DpCfWin32Launch`）的具体分支未逐条动态验证，结论基于字符串、导入表和函数调用关系推断。
- 样本包含开发/测试痕迹，正式版行为可能与当前样本不同。
- 文中“推断”部分均基于静态证据，未做动态抓包验证。

---

## 13. 附录：复现命令

### 13.1 文件类型与哈希

```bash
file CF上号器_v1.17.exe host.dat
shasum -a 256 CF上号器_v1.17.exe host.dat
md5 CF上号器_v1.17.exe host.dat
```

### 13.2 PE 结构

```bash
objdump -f CF上号器_v1.17.exe
objdump -d --start-address=0x14001AEE0 --stop-address=0x14001D540 CF上号器_v1.17.exe
```

### 13.3 字符串

```bash
strings -a -n 6 CF上号器_v1.17.exe | grep -iE "card|token|wegame|rail|qwh|ntop"
```

### 13.4 内嵌 DLL 提取

```python
data = open("CF上号器_v1.17.exe", "rb").read()
open("qwh_cf_hook.dll", "wb").write(data[0xDB8E0:0xDB8E0 + 0x18E00])
```

### 13.5 导入表查看

```bash
python3 -c "import pefile; pe=pefile.PE('CF上号器_v1.17.exe'); pe.parse_data_directories(); [print(hex(i.address), e.dll.decode(), i.name.decode()) for e in pe.DIRECTORY_ENTRY_IMPORT for i in e.imports if i.name]"
```

---

## 14. 最终判断

该样本的本质是：

> **卡密平台配套的 CF 自动上号客户端。**

它把“卡密授权、账号票据、QQ NT 登录、WeGame OP 转换、Rail 启动、DLL 注入”打包成了一个完整工具。其技术含量不低，但行为特征也非常明确：

- 非官方启动链；
- 进程注入；
- 伪造 WeGame Rail 环境；
- 账号票据本地缓存与远端上报；
- 第三方打码/代理；
- 无签名、要求管理员权限。

如果只是研究其工作原理，以上内容已经覆盖主要链路；如果要实际运行，建议按高风险样本处理，不要在有价值账号的环境中直接使用。

---

## 15. 补充：host.dat 完整解密与卡密协议细节（二次分析）

### 15.1 host.dat 算法还原

文件结构（已通过代码与 HMAC 双重验证）：

```text
0x00  magic    d8 3a 91 c7
0x04  version  02
0x05  iv       16 字节
0x15  len      uint16，密文长度
0x17  ct       len 字节（AES-256-CBC / PKCS7）
0x17+len       hmac 32 字节（HMAC-SHA256(macKey, 文件头 + 密文)）
```

密钥派生（RVA `0x8E300`）：

- 程序内固定 32 字节种子由 RVA `0x8E300` 的混淆循环生成，常量表位于 RVA `0x1400FB770` / `0x1400FB780` / `0x1400FB790`；
- `aesKey = SHA256(seed + "aes-host-v2")`
- `xorKey = SHA256(seed + "xor-host-v2")`
- `macKey = SHA256(seed + "mac-host-v2")`

解密（RVA `0x8F550`）：

```text
1. 校验 magic / version
2. HMAC-SHA256(macKey, raw[0:0x17+len]) == raw[0x17+len : 0x17+len+32]
3. p = AES-256-CBC-Decrypt(aesKey, iv, ct)，去 PKCS7
4. p[i] ^= xorKey[(3+7*i) & 31] ^ ((xorKey[i & 31] + 17*i) & 0xFF)
5. 每字节高低 4 位互换（rol 4）
```

解密结果：

```text
http://111.170.141.190:8790
```

该算法完全可逆：按相反顺序重新加密并重算 HMAC，程序会正常接受。这也是本工具改写
`host.dat` 指向本地服务、且不修改 exe 的依据。

### 15.2 卡密接口协议

请求：

```http
POST /api/card/redeem HTTP/1.1
User-Agent: DfCard/1.0
Content-Type: application/json

{"card_key":"<16位卡密>"}
```

请求体不含随机数、时间戳或签名字段，仅卡密本身。

响应（JSON 信封）：

| 字段 | 含义 |
|---|---|
| `success` | 是否成功 |
| `key` | Base64，RSA-2048(PKCS1 v1.5) 加密的 AES 会话密钥 |
| `iv` | Base64，12 字节 GCM nonce |
| `tag` | Base64，16 字节 GCM tag |
| `data` | Base64，AES-GCM 密文 |
| `verify` | 校验字段 |
| `ok` | 是否成功（线上用 `ok`，不是 `success`） |
| `qq` | 账号 QQ（明文放在信封外层） |
| `zone` / `note` / `extra` | 区服与备注（`note` 形如 `编号:xxx----卡密----战区----区服----段位:宗师----...`） |
| `expires_at` | 毫秒时间戳 |

解密（RVA `0x8D700`）：

```text
aesKey = RSA-DecryptOAEP-SHA256(内嵌私钥, base64(key))
plain  = AES-GCM-Open(aesKey, iv, data || tag)
```

内嵌私钥位于 RVA `0x1400FB140`，格式为 `BCRYPT_RSAFULLPRIVATE_BLOB`（magic `RSA3`，
2048 位，e = 65537），已用 `p*q == n` 校验通过。

### 15.3 对“本地缓存”的意义

- 响应信封自带 RSA 加密的会话密钥，客户端用内嵌私钥即可解开，缓存原始字节后重放仍能正常解密；
- 请求只依赖卡密，不含一次性随机数，因此同一卡密的响应可以整包复用；
- 结论：无需重写协议，只要在 `host.dat` 指向的本地 HTTP 服务里按卡密缓存 `redeem` 的
  原始响应体即可。

配套实现见 `cf-cache-tool/`（Go）与 `dist/`（Windows 可执行文件及说明）。
