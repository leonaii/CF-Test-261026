# 交接说明：独立 ck 上号器

- 整理时间：2026-10-23
- 用途：换电脑后直接在新机器上继续，不需要重新摸一遍
- 当前阶段：协议验证阶段，代码只落在源码目录，未打包发布

---

## 0. 一句话目标

做一个独立的上号工具，脱离 `CF上号器_v1.17.exe`，自己完成「账号数据 -> WeGame OP 票据 -> 启动 CF」的闭环。

已经验证到：换 OP 这一步被腾讯的 `mcode`（本机设备指纹）挡住，必须依赖本机 WeGame 的 rail 组件。这是当前唯一的技术卡点。

---

## 1. 目录与关键文件

```text
新1.17上号器/
  CF上号器_v1.17_逆向分析报告.md        上号器完整逆向结论（1105 行，先看这个）
  上号器逆向方法论与版本复用指南.md      方法论（1467 行）
  REMARK.md                            CG 工具自身架构说明
  交接说明_ck上号器.md                  本文件

  资源/上号器/yuanmeng/
    CF上号器_v1.17.exe                 主样本（1,597,952 字节）
    host.dat                           卡密服务器地址（加密容器）

  资源/验证/yuanmeng/
    CF验号器_v1.3.exe                  验号器样本（367,616 字节）
    host.dat / config.ini

  cf-cache-tool/                       Go 源码（CG.exe 的来源，本次改动都在这里）
  dist/                                产物目录（当前 CG.exe 是旧版本，未重新打包）
  work/                                逆向中间产物（asm、脚本）
```

`work/` 里可复用的脚本：

| 文件 | 作用 | 状态 |
|---|---|---|
| `decrypt_host.py` | host.dat v2 解密 | 已修正 XOR 公式，可跑 |
| `pe.py` | PE 读取工具（RVA -> 文件偏移） | 已修正样本路径 |
| `cf_check.py` | 验号接口直连 CLI | 可用 |
| `full.asm` | 上号器完整反汇编（12MB） | 可用 |
| `card_redeem.asm` / `card_http.asm` / `codec.asm` / `hostdat.asm` | 分模块反汇编 | 可用 |
| `rsa_priv.txt` | 内嵌 RSA 私钥（n/e/d/p/q） | 可用 |

换机后重新提取的临时文件（本机 `/tmp` 下的会丢）：

```bash
# qwh_cf_hook.dll：文件偏移 0xDB8E0，长度 0x18E00，101,888 字节
python3 -c "d=open('资源/上号器/yuanmeng/CF上号器_v1.17.exe','rb').read();open('/tmp/qwh_cf_hook.dll','wb').write(d[0xDB8E0:0xDB8E0+0x18E00])"
```

---

## 2. 已验证的事实（别重复验证）

### 2.1 服务器与接口

`host.dat` 三份（根目录、上号器、验号器）解密后明文一致：

```text
http://111.170.141.190:8790
```

算法：`seed` 由 exe 内混淆循环生成，密钥 `SHA256(seed + "aes-host-v2" / "xor-host-v2" / "mac-host-v2")`，
`明文 = nibbleSwap(AES-256-CBC(ct) ^ xorMask)`，XOR 索引是 `(3 + 7*i) & 31`（不是 `(i+3)&31`）。

验号器只调两个接口，均为 POST JSON、无签名、无时间戳、UA 不校验：

| 接口 | 请求体 | 说明 |
|---|---|---|
| `/api/card/boot` | `{}` | 探测服务器与封禁，返回 `{"ok":true,"banned":false}` |
| `/api/query/status` | `{"code":"<编号或卡密>"}` | 验号主接口 |

参数规则（实测）：卡密 16 位大写英文；编号 8 位小写英文（可含下划线，如 `HY_KYEJK`、`HY_XLRVS`）。

服务端按 IP 限频约 8 秒，触发返回：

```json
{"ok":false,"rate_limited":true,"retry_after_ms":7477,"message":"查询过于频繁，请 8 秒后再试"}
```

成功响应（完整字段，含敏感 ck）：

```json
{"ok":true,"query_code":"HY_KYEJK","qq":"1978628973","zone":"江苏电信一区","has_data":true,
 "query_at":1792728018726,"ck":"1978628973----Tencent----uin=o1978628973;skey=...;p_skey=...",
 "cookie":"uin=o1978628973;skey=...;p_skey=...","ck_ok":true,
 "status":"正常","fenghao":"未封禁","dongjie":"未冻结","zaixian":"离线","renlian":"不要人脸",
 "shiming":"已实名","lingpai":"无令牌","nickname":"","level":"","daqu":"江苏电信一区","zhanqu":"",
 "elapsed_ms":1133,"message":"查询完成"}
```

注意：每次查询 `p_skey` 都不同，平台每次都重新登录一次腾讯。

### 2.2 redeem 信封

请求：`POST /api/card/redeem`，body `{"card_key":"<16位卡密>"}`，UA `DfCard/1.0`。

响应字段：`ok / qq / zone / note / extra / expires_at / key / iv / tag / data`。

```text
key : 256 字节 = RSA-2048，填充 RSA-OAEP-SHA256，解出 32 字节随机 AES-256 密钥
iv  : 12 字节 GCM nonce
tag : 16 字节
data: AES-256-GCM 密文
```

用内嵌私钥（RVA `0x1400FB140`，`BCRYPT_RSAFULLPRIVATE_BLOB`）解密。**key 与 ck 无关**，请求里没有 ck，解出的密钥也不是卡密哈希。

明文是 QQ NT 缓存（KV 文本，首行 `[QQ号]`）：

```text
Token0106 / Token010A / Token0143 / sessionKey
uid / tgpid / ctid / ct / dt / st / st_key / sig / client_key / biz_ticket / machine_id
```

### 2.3 ck 与 OP 的关系（已证伪的错误假设）

上号器不是「ck 转 OP」，ck 是 NT 登录的**产物**。真实链路：

```text
redeem -> NT 缓存 -> 本机 QQ NT 协议登录（ECDH P-256 + TEA + AES-GCM + SSO over TCP）
       -> clientkey / skey / openkey（即 ck）
       -> LoginByQq(qq_info_type=6, openid/sig/client_key/st_key) -> OP 票据
       -> RefreshClientTicket / GetGameTicket / GetGameLaunchInfo -> 启动
```

用 ck 直接打 LoginByQq 的结果：

```text
{"result":{"error_code":8000001,"error_message":"sig empty"}}
```

### 2.4 换票卡在 mcode（当前卡点）

用真实 redeem 票据（即刻领取，非缓存）打 LoginByQq：

```text
POST https://www.wegame.com.cn/api/v2/wegame.auth.game.UnionAuth/LoginByQq
{"mappid":10001,"clienttype":1000000,
 "login_info":{"qq_info_type":6,"appid":2112513,"openid":"<uid>",
               "sig":"<sig>","client_key":"<client_key>","st_key":"<st_key>"},
 "config_params":{"lang_type":0,"login_type":1,"channel_id":0,
                  "promote_id":"0","client_version":"6.2.2.12044"},
 "head":{"credential":{"game_ver":"7.9.8.0","ct":"<ct>","mid_list":"<machine_id>"}}}

HTTP 200
{"result":{"error_code":8010170,"error_message":"ret: 0, func_ret: 8010170, err: mcode error"}}
```

四组对照（同一票据）：

| head.credential | 返回 |
|---|---|
| `mid_list=<machine_id>` | `mcode error` |
| `mid_list=""` | 认证失败 |
| `mid_list=<machine_id>` + `mcode=<machine_id>` | `mcode error` |
| head 全空 | 认证失败 |

结论：有值才进入 mcode 校验，值是关键。`mcode` 是 WeGame 本机组件用硬件信息 + 内置密钥生成的设备指纹，不是 `machine_id` 明文，协议层伪造不出来。

推论：换 OP 必须在装了 WeGame 的目标机上、由 `rail_api64.dll` 参与完成。这也是上号器必须伪造 Rail 环境并注入 `qwh_cf_hook.dll` 的原因。

### 2.5 上号器 WeGame 协议模板（RVA 与字符串）

```text
0xF7741  login_info":{"qq_info_type":6,"appid":2112513,"openid":"
0xF777A  ","sig":"
0xF7784  ","client_key":"
0xF7795  ","st_key":"
0xF7E49  "mappid":10001,"clienttype":1000000,"login_info":{"qq_info_type":4,"appid":2112513,"openid":"
0xF7EA9  },"config_params":{"lang_type":0,"login_type":1,"channel_id":0,"promote_id":"0","client_version":"
0xF7F10  "},"head":{"credential":{"game_ver":"","ct":"","mid_list":"
0xF7A70  /api/v2/wegame.auth.game.UnionAuth/
0xF77B0  /api/v2/wegame.auth.Auth/GetGameToken
0xF7801  api/v1/wegame.gameadapter.GameAdapter/GetGameLaunchInfo
0xF7DE0  /api/v2/wegame.auth.game.UnionAuth/RefreshClientTicket
0xF7BC9  "head":{"credential":{"ct":"...","game_ver":"7.9.8.0","machine_id":"...","wg_client_ver":"6.2.2.12044"}
0xF7CD8  "userid_type":0,"version_id":1970363492204544}
0xF7903  GetWebTicket
0xF79CE  LoginWG / 0xF79E2 RequestCTT / 0xF79F8 GetWGCTT / 0xF7A01 GetClientTmpTicket
0xF79B1  ChangeTempToken
0xF8689  in=o%llu; skey=%s; p_skey=%s; p_uin=o%llu; pt4_token=%s
0xF8732  ----uin=o
0xDBEEA  ...FEDCBA9876543210...,"ctt":"
```

常量：`appid=2112513`、`mappid=10001`、`clienttype=1000000`、`gid=55555`、`version_id=1970363492204544`、`game_ver=7.9.8.0`、`wg_client_ver=6.2.2.12044`。

### 2.6 Rail 启动链（上号器实现）

```text
共享内存  Rail.SharedMemory.%u
管道      \\.\pipe\rail_pipe_%u
          \\.\pipe\rail_named_pipe_%u
          \\.\pipe\rail_named_pipe_host_%u
          \\.\pipe\rail_named_pipe_guest_%u
          样本实际用 15
文件      \TCLS\mmog_data.xml、\NTCLS\mmog_data.xml
          \NTCLS\ntcls_config\mmog_data.xml、\NTCLS\tgp_config\mmog_data.xml
          \rail_user_data\2、version.ini、last_cf_rail.bin
参数      -SSN:<区服>、-cf:1、-HostIP=、-q
注入      VirtualAllocEx + WriteProcessMemory + CreateRemoteThread(LoadLibraryW)
          qwh_cf_hook.dll，101,888 字节，无导出，仅导入 KERNEL32
```

关键错误串：`没有该区服Rail模板`、`Rail模板坏了`、`Rail模板缺启动包`、`没有完整CF启动签名`、`请选择 x64\crossfire.exe`、`找不到rail_api64.dll`、`WeGame Init Failed. (2310/2311)`、`闪崩:管道未连上`。

说明：Rail 模板来自本机 `rail_apps` 目录，不是内嵌的。所以独立实现的前提是本机已装 WeGame + CF。

---

## 3. 已完成的工作（代码）

改动全部在 `cf-cache-tool/`，已通过 `GOOS=windows go build`、`go vet`、darwin `go build`。

| 文件 | 改动 |
|---|---|
| `verify.go`（新增） | `/api/query/status` 客户端；`verifyStatus` 结构；`verifySummary` 列摘要；`verifyDetailText` 详情文本；`prettyJSON` 公共工具 |
| `pick.go` | `accountRow` 增加 `Verify / VerifyRaw / VerifyErr / VerifyAt`（仅内存，不落盘） |
| `gui_windows.go` | 第 9 列「验证数据（双击查看）」宽 260；工具栏「验证」按钮；`verifySelected` 手动验证；`storeVerifyResult`；`showVerifyError` 限频弹框；`runVerifyDialog` 完整 JSON 弹窗；`verifyColumn` 常量；`clickedCol` 双击分流 |
| `setup_windows.go` | `prettyJSON` 迁出到 `verify.go`，移除 `encoding/json` 导入 |

行为约定：

- 验证只在点击按钮时发起，单次串行，不自动刷新
- 未选中行时按钮禁用
- 失败（含限频）弹提示框，成功弹多行文本框显示完整 JSON（含 ck）
- 双击「验证数据」列回看上次结果；双击其它列仍是原「原始备注」窗口（旧行为保留）

**未做**：没有重新打包 `dist/CG.exe`，`dist/` 里还是旧版本。

打包命令：

```bash
cd cf-cache-tool && ./build.sh --no-loader
```

`build.sh` 会自动把 `version.txt` +1，并覆盖 `dist/CG.exe`、`dist/CG-cli.exe`。

---

## 4. 未完成的任务（按优先级）

### P0：拿到 mcode 的真实来源

需要在新机器（Windows）上做：

1. 取本机 `rail_api64.dll`
   - 位置：WeGame 安装目录，或 `%WeGameApps%\rail_apps\...`
   - 拿到后分析导出表，找 machine code 相关函数
2. 抓一次真实 WeGame 登录 CF 的包
   - 工具：Fiddler / mitmproxy / Charles，抓 `www.wegame.com.cn`
   - 目标接口：`/api/v2/wegame.auth.game.UnionAuth/LoginByQq`
   - 只要字段结构：`head.credential` 里的 `mid_list` 长什么样、有没有 `mcode`、长度与编码
   - 敏感值可以打码

判定标准：

- 若 `mid_list` 是几十到几百字节的 base64/hex，则确认是加密指纹，必须调用 rail_api64.dll 生成
- 若是 32 位 hex，则可能还有别的校验字段，需要对比完整 body

### P1：NT 登录协议（M2）

只有 P0 通了才有意义。三档深度：

| 深度 | 内容 | 用途 |
|---|---|---|
| 浅 | 复刻到拿 `client_key / st_key / ct / sig` | 能换 OP，票过期要重登 |
| 中 | 加 RefreshClientTicket / GetWGCTT | 长会话 |
| 深 | 完整 SSO（心跳、重连、多票据） | 最稳，工作量最大 |

上号器内实现位置参考逆向报告第 5 节（QQ NT 登录协议：ECDH P-256、TEA、AES-GCM、SSO over TCP）。

### P2：Rail 环境与启动（M4）

- 复用本机 `rail_apps` 里的 CF 模板
- 造 `mmog_data.xml`（用 OP 票据填字段）
- 创建共享内存与命名管道（编号 15）
- `CreateProcess crossfire.exe -SSN:<区服> -cf:1`
- 注入 `qwh_cf_hook.dll`（已可提取）

区服映射：上号器内置 42 个区服，`ssn` 拼进 `-SSN:`。示例：山东网通一区、江苏电信一区。

### P3：CG 工具的收尾

- 验证结果目前只在内存，重启即丢；需要的话可落盘到 `cache/<卡密>/verify.json`
- 未打包 dist
- `work/cf_check.py` 可加 `--json`、批量模式

---

## 5. 需要的输入清单（换机前准备）

| 项 | 说明 | 现状 |
|---|---|---|
| Windows 机器 | 必须装 WeGame + CF | 待准备 |
| `rail_api64.dll` | 从 WeGame 目录取 | 未拿到 |
| LoginByQq 真实抓包 | 结构即可，不必真值 | 未拿到 |
| Go 1.26+ | 编译 CG 工具 | 新机装 |
| 卡密 | 测试用，会消耗领取次数 | 有限，省着用 |
| 云端凭据 | 在 `cloud.go` 常量里（Nextcloud WebDAV） | 已有 |

测试卡密（已用过，可再次尝试，注意额度）：

```text
CFHEEUGNCSHXODBN   编号 HY_KYEJK   江苏电信一区
CFYAWOJVXAROMRZA   编号 HY_XLRVS   山东网通一区
```

---

## 6. 新机器上的操作顺序

```bash
# 1. 拷贝整个项目目录
#    新1.17上号器/

# 2. 装 Go 1.26+，然后在源码目录编译验证
cd cf-cache-tool
GOOS=windows GOARCH=amd64 go build -o /tmp/cg.exe .
go vet ./...

# 3. 重新提取 hook DLL（临时文件不会跟着过去）
python3 -c "d=open('../资源/上号器/yuanmeng/CF上号器_v1.17.exe','rb').read();open('/tmp/qwh_cf_hook.dll','wb').write(d[0xDB8E0:0xDB8E0+0x18E00])"

# 4. 确认 host.dat 解密仍然可用
cd ../work && python3 decrypt_host.py

# 5. 在 Windows 上抓包 + 取 rail_api64.dll（P0）

# 6. 分析结果出来后，再决定 M2 的实现深度
```

验号接口快速验证（不需要编译）：

```bash
cd work
python3 cf_check.py --boot
python3 cf_check.py --code HY_KYEJK
```

---

## 7. 风险与注意事项

1. **封号风险**：CF 有 TP 保护，注入 DLL 与伪造 Rail 环境都可能触发检测。上号器本身就在做这个，测试用小号。
2. **顶号**：同一账号多地登录会互踢，测试时别和平台侧并发。
3. **限频**：验号接口 8 秒一次，redeem 可能消耗卡密领取额度，测试前确认。
4. **ck 敏感性**：`skey` / `p_skey` 是可登录凭证，弹窗会完整显示，注意截图外传。
5. **mcode 是本机绑定**：不要试图跨机器复用票据，实测必然 `mcode error`。
6. **不要删 `.codebuddy` 目录**（项目记忆）。
7. 云端索引里目前的账号：`CFHEEUGNCSHXODBN`（QQ 1978628973，江苏电信一区，录制于 2026-10-23 06:49）。

---

## 8. 一句话交接

协议层已经摸到最后一环：NT 票据能解、LoginByQq 能打通到 mcode 校验，只差 `mcode` 的生成方式。新机器上第一件事是拿 `rail_api64.dll` 和一次真实抓包，其余都建立在它之上。
