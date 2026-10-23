# CG 维护与扩展说明（REMARK）

面向“以后还要改这个工具”的人。看完应当能独立新增一个上号器类型、发一版，并在出问题时知道去哪看。

---

## 0. 一句话架构

`CG.exe` 是一个 Win32 界面程序（Go + lxn/walk），把**上号器 exe** 跑起来，
用本地 HTTP 服务（`127.0.0.1:18800`）截获上号器与服务器的接口应答并落盘；
回放时不再连真实服务器，直接把这些应答“喂”回去。
上号器本体、host.dat、账号缓存、程序更新全部走 Nextcloud WebDAV。

---

## 1. 目录结构

### 1.1 仓库

```text
新1.17上号器/
  REMARK.md              本文件
  dist/                  产物（CG.exe / CG-cli.exe / 使用说明.md）
  资源/
    上号器/
      yuanmeng/          上号器资源（不打进 exe，构建时上传云端）
        CF上号器_v1.17.exe
        host.dat
  cf-cache-tool/         源码
```

资源目录可以改，`build.sh` 认两个位置：

1. 环境变量 `LOADER_ASSET_DIR`
2. 默认 `../资源/上号器/yuanmeng`

```bash
LOADER_ASSET_DIR=../资源/上号器/xxx ./build.sh
```

### 1.2 源码（`cf-cache-tool/`）

```text
main.go              入口 + 命令行子命令（push/pull/cloud/run/list/clean…）
gui_windows.go       主窗口、录制/回放流程、本地代理服务、列表模型、更新接线
setup_windows.go     首次配置窗口、标记窗口、原始备注视窗
setup_stub.go        非 Windows 占位（只为跨平台编译）
launch_stub.go       非 Windows 的 restartIfNeeded 占位
loader.go            上号器类型注册表（新增类型改这里）
loaderfetch.go       云端下载上号器/host.dat、sha256 校验、本地标记、本机缓存
update.go            自动更新：清单、下载、替换、清理、单飞锁
version.go           版本号变量 + 语义版本比较
cloud.go             WebDAV 客户端（PROPFIND/PUT/GET/MKCOL/DELETE）
index.go             云端索引 index.json 与单卡 zip 的读写、合并
pick.go              扫描本地缓存生成列表行、排序
cache.go             缓存落盘格式、录制暂存目录、合并录制
hostdat.go           host.dat v2 编解码 + 改写/还原（原地址备份为 .orig）
configini.go         回放前重写 config.ini（占位卡密 + 提示语）
appconfig.go         本地配置、标记、目录迁移
uninstall.go         卸载本地数据（按名单删 / 整目录强删）
uninstall_windows.go 枚举并结束占用目录的进程
cleanup_windows.go   退出后延迟删除 CG.exe.old
build.sh             打包脚本（版本自增、编译、资源同步、发布）
version.txt          当前版本号
```

### 1.3 用户机器上的数据

```text
%AppData%\CG\
  config.json        类型 / 工作目录 / 卡密 / 标记
  startup.log        启动与重启留痕（每次启动覆盖写）
  loaders\<类型ID>\  下载缓存 exe-<sha256>、dat-<sha256>（换目录不用重下）
  data\              默认工作目录
    CF上号器_v1.17.exe
    host.dat / host.dat.orig
    config.ini
    .loader.json     本地已安装的上号器版本与 sha256
    cache\<卡密>\*.body + *.meta.json
    cache.recording\ 录制暂存
```

`run.log` 在 **exe 同级**（跟着绿色包走），不在 `%AppData%`。

### 1.4 云端（Nextcloud WebDAV）

```text
CG/
  update/
    version.json        程序更新清单
    CG.exe              新版本程序
  loaders/<类型ID>/
    meta.json           上号器清单（文件名、大小、两个 sha256）
    <上号器文件名>       上号器本体（约 1.5MB）
    host.dat
  <卡密>/<类型ID>/
    index.json          账号列表（约 2KB，列表页只下这个）
    cards/<卡密>.zip    单个账号的完整缓存（约 5KB，回放时才下）
```

根地址：`https://cloud.51i.cc/remote.php/dav/files/leonai/`
账号密码写死在 `cloud.go`，**拿到 exe 的人都能读写整个网盘，别外发**。

---

## 2. 三条数据流

### 2.1 录制

```text
点“新增”
  → 读 host.dat 得到真实上游，备份 .orig，改写为 http://127.0.0.1:18800
  → 启动上号器（工作目录 = data\）
  → 请求进本地服务：转发到真实上游，同时把 /api/card/* 应答写进 cache.recording\<卡密>\
  → 上号器退出 → 还原 host.dat
  → cache.recording 按卡密合并进 cache\（有 redeem 才算完整）
  → 逐账号上传 cards/<卡密>.zip，合并更新 index.json
```

### 2.2 回放

```text
选中一行 → 点“选中上号”
  → 本地没有这个卡密 → 先从云端下 cards/<卡密>.zip
  → host.dat 改写为本地地址（同上）
  → 写 config.ini（占位卡密）
  → 启动上号器；本地服务直接用 cache 里的 .body 应答，不访问外网
  → 退出后还原 host.dat
```

### 2.3 列表与同步

```text
启动 → 扫本地 cache → 渲染列表 → 后台拉 index.json 合并（云端独有 = 灰行）
标记 → 本地立即生效 → 后台 GET 最新 index → 只改这一条 → PUT
刷新 → 重拉 index.json + 重扫本地
```

流量：看列表 ≈2KB；回放一个号 ≈5KB；录一个号 ≈7KB；标记 ≈2KB。

---

## 3. 新增一个上号器（完整步骤）

### 3.1 准备两个文件

| 文件 | 说明 |
|---|---|
| 上号器 exe | 例如 `CF上号器_v1.18.exe`，原样，不要改成 loader.exe 这种通用名 |
| host.dat | 该版本自带的那份，且**未被本工具改写**（不能是已指向 127.0.0.1 的） |

放到 `资源/上号器/<名字>/`，例如 `资源/上号器/yuanmeng/`。
目录名只是本地组织用，云端用的是**类型 ID**。

### 3.2 注册类型（`loader.go`）

```go
var loaderProfiles = []loaderProfile{
    {
        ID:       "1",                       // 云端目录名，改了就换一份云端数据
        Name:     "圆梦科级上号器",             // 界面显示的类型名
        ExeName:  "CF上号器_v1.17.exe",        // 落盘文件名，必须与云端 meta 一致
        SiteName: "圆梦卡网",
        Site:     "https://www.517fk.cn/links/8B1B512A",       // 购买地址
        SiteNote: "0V-10V / 10V-19V / 20V-49V / 50-99V 这家的", // 购买说明
    },
}
```

规则：

- **ID 唯一且稳定**，决定云端路径 `CG/<卡密>/<ID>/`，改了等于换一套账号数据
- `ExeName` 必须和云端 `meta.json` 的 `exe` 完全一致，否则每次启动都判定“本地没有”
- 只往数组里加一项，界面下拉框自动出现新类型

### 3.3 上传资源

```bash
cd cf-cache-tool
./build.sh                 # 编译 + 比对 sha256，变了才上传
./build.sh --no-loader     # 离线打包，跳过资源同步
```

脚本流程：

1. 从 `资源/上号器/yuanmeng/` 取 `*.exe` 和 `host.dat`，算 sha256 与大小
2. 读云端 `CG/loaders/<ID>/meta.json`，两个 hash 都对得上 → “云端已是最新”并跳过
3. 否则 MKCOL 三级目录 → 上传 exe、host.dat → 写 meta.json

meta.json 格式（手传也要照这个写）：

```json
{
  "id": "1",
  "name": "圆梦科级上号器",
  "exe": "CF上号器_v1.17.exe",
  "exe_sha256": "2d855012...",
  "exe_size": 1597952,
  "dat_sha256": "f72f8d03...",
  "updated": "2026-10-23 06:26:59"
}
```

> 换版本（v1.18）：新文件放进同一目录、删旧的、改 `loader.go` 的 `ExeName`（ID 建议保持），
> `./build.sh --publish`。客户端下次启动按新 hash 下载，并自动删掉旧文件名。

### 3.4 验证清单

```text
[ ] 清空 %AppData%\CG\data 启动 → 日志“下载上号器 / 上号器已就绪”
[ ] 再启动一次 → 不再下载（或“上号器取自本地缓存”）
[ ] 录制一次 → cache 下出现该卡密目录，含 api_card_redeem.body
[ ] 云端出现 CG/<卡密>/<ID>/index.json 与 cards/<卡密>.zip
[ ] 换台机器启动 → 列表能看到账号（灰行），点“选中上号”能自动下载并成功上号
```

---

## 4. 新上号器的“规则”要求

本工具依赖的是**网络行为**，不是上号器内部实现。新增类型必须满足：

### 4.1 host.dat 必须是同一种格式

`hostdat.go` 实现 v2 格式：magic `d8 3a 91 c7` + ver 2 + 16 字节 IV + uint16 长度 +
AES-256-CBC 密文 + 32 字节 HMAC，明文再经 `xorMask` 与 `nibbleSwap`。
密钥由二进制内固定种子派生。

**新版如果换了格式/密钥，必须重写 codec**，否则 `DecryptHostDat` 失败，录制和回放都起不来。

自检：

```bash
CG-cli.exe hostinfo -dir <目录>    # 能打印上游地址就是兼容的
```

### 4.2 接口路径

已知接口（`cache.go: isKnownEndpoint`）：

```text
api_card_boot      固定应答，一般不用录
api_card_peek
api_card_log
api_card_redeem    ← 唯一判定“录制完整”的依据
```

其它 `/api/...` 也会被记录（文件名由路径 + 查询串 hash 生成），但**合并时只认 redeem**。
若新上号器的关键接口不叫 `api_card_redeem`，要改三处：
`isKnownEndpoint`、`mergeRecording` 的完整性判断、`pick.go` 读取用的 endpoint 名。

### 4.3 端口与协议

固定 `127.0.0.1:18800`（`gui_windows.go: a.listen`），只提供 **http**。
上号器必须按 host.dat 里的地址发请求；如果强制 https，需要改成“本地 HTTPS + 自签证书”，
并在 host.dat 里写 https 地址。

### 4.4 exe 的启动约定

- 以 `data\` 为工作目录启动，`config.ini` 与它同级
- 不解析退出码，进程结束即认为录制完成
- 若 exe 需要管理员权限，本地服务也要以管理员运行（否则 18800 可能被拦）

### 4.5 config.ini 模板

回放前用 `configini.go` 的模板覆盖：

```ini
card=CFAAAAAAAAAAAAAA
dir=请选择游戏目录
ssn=346
auto_zone=0
```

首次覆盖前备份为 `config.ini.orig`；若当前已是占位卡密则**不重复覆盖**（保留用户选好的目录）。
新上号器字段不同就改这里。

---

## 5. 打包与版本

```bash
cd cf-cache-tool

./build.sh                 # 版本 +1（1.0.24 -> 1.0.25），编译，同步上号器资源
./build.sh 1.2.0           # 指定版本
./build.sh --publish       # 再上传 CG/update/{CG.exe,version.json}
./build.sh 1.2.0 --publish
./build.sh --no-loader     # 跳过资源同步
```

- 版本存 `version.txt`，通过 `-ldflags "-X main.toolVersion=x.y.z"` 注入
- `toolVersion` 必须是 `var`，写成 `const` 会让 ldflags 失效
- 产物：`dist/CG.exe`（`-H=windowsgui -s -w -trimpath`）、`dist/CG-cli.exe`

**链接器会删死代码**：没被引用的函数连同字符串常量一起消失，功能会“静默失效”。
曾因漏写 `go a.checkUpdate()` 导致整个更新模块被优化掉。所以脚本里有一道字符串校验，
新增功能后建议补检查项：

```bash
python3 -c "d=open('../dist/CG.exe','rb').read(); print('关键字'.encode() in d)"
```

---

## 6. 自动更新

清单 `CG/update/version.json`：

```json
{ "version": "1.0.25", "file": "CG.exe", "sha256": "...", "size": 12011008, "notes": "build ..." }
```

```text
启动（第一界面之前就开始，后台）→ GET 清单
  → 版本不新：结束
  → 有新版本：下载到 CG.update.tmp → 校验大小 + sha256
       → CG.exe 改名 CG.exe.old → 新文件就位 → 记录“本进程已更新到 vX”
  → 不在录制/回放时：关闭窗口 → main() 里拉起新版本
```

必须知道的约束：

- **运行中的 exe 不能删、不能覆盖，只能改名**，所以更新当次一定存在 `CG.exe.old`
- `.old` 在下次启动清理：先删一次 → 失败每 1.5s 重试（15s）→ 仍失败则起分离 cmd，
  `ping -n 4` 后 `del /f /q`（本进程退出后由系统删除）
- `applyUpdate` 只用固定名 `.old`，**删不掉就跳过本次更新**，绝不生成第二个文件
- 首配窗口与主界面都会检查，用互斥锁 + `updateApplied` 保证**只下载一次**
- 重启走 `main()`：`runGUI()` 返回后 `restartIfNeeded()` 用 `CreateProcess`（DETACHED_PROCESS）拉起，
  不经过 `cmd`，避免引号转义

---

## 7. 云端约定与坑

| 项 | 约定 |
|---|---|
| 目录创建 | **先建父目录再建子目录**，否则 Nextcloud 返回 409 |
| 路径转义 | 中文文件名必须逐段 `url.PathEscape`，整段拼接会 404 |
| 覆盖上传 | PUT 直接覆盖，不存在也能创建 |
| 文件大小 | 用 `PROPFIND Depth:0` 的 `getcontentlength`（部分服务端 HEAD 返回 0） |
| 索引合并 | 先 GET 最新 `index.json` → 只改本机这条 → PUT，避免覆盖别人新增的账号 |
| 失败策略 | 上传失败只写日志，本地照常可用，下次录制再补 |

---

## 8. 本地存储

### 8.1 config.json

```json
{
  "profile": "1",
  "work_dir": "C:\\Users\\x\\AppData\\Roaming\\CG\\data",
  "cloud_name": "leonai",
  "auto_cloud": true,
  "marks": { "CFYAWOJVXAROMRZA": "ban" }
}
```

标记取值：`normal` / `ban` / `black` / `expired` / `other`。

### 8.2 缓存文件

```text
cache/<卡密>/
  api_card_peek.body / .meta.json
  api_card_log.body / .meta.json
  api_card_redeem.body / .meta.json
```

`.meta.json` 含 `card / endpoint / status / header / recorded_at / bytes / request`。
缺 `recorded_at` 的旧数据按文件 mtime 排序，保证“最新在前”。

### 8.3 卸载

“切换上号器”会**整个删掉**当前工作目录（`%AppData%\CG\data`），删不掉时枚举进程，
把 exe 路径在该目录内的进程结束再重试。只改卡密或只改类型（目录不变）时不删目录。

---

## 9. 踩过的坑（改代码前先看）

1. **walk 对话框：控件值必须在关闭前取**
   `dlg.Run()` 返回后控件已销毁，`combo.CurrentIndex()` 得到 `-1`、`Text()` 得到空串。
   曾导致“点开始使用没反应”“标记永远显示正常”。取值/校验都写在按钮的 `OnClicked` 里。

2. **不要用 `exec.Command("cmd", "/c", 拼接字符串)`**
   Go 会转义内层引号为 `\"`，cmd 不认，命令被拆坏（自动重启因此失效）。
   需要 shell 时用 `SysProcAttr.CmdLine` 原样传入。

3. **运行中的 exe 只能改名，不能删**——更新与清理逻辑都基于这一点。

4. **链接器死代码删除**——见第 5 节，发版前用字符串校验兜底。

5. **MKCOL 顺序 / 中文路径转义 / HEAD 返回 0**——见第 7 节。

6. **界面框架限制**：walk 用 Win32 原生控件，按钮/输入框背景不接受自定义颜色。
   视觉分层只能靠字号、灰阶、间距、分隔线；要深色皮肤必须自绘控件。
   多行文本框必须给 `MinSize.Height`，否则 VBox 只给一行高度（看起来像单行框）。

---

## 10. 排障速查

| 现象 | 先看哪里 |
|---|---|
| 点了没反应、闪一下就没了 | `%AppData%\CG\startup.log`（启动流程留痕） |
| 录制/回放失败 | `run.log`（与 exe 同级，含上号器路径、host.dat 上游、账号明文） |
| 列表没有云端账号 | `run.log` 里 `拉取云端列表失败` / 检查 `CG/<卡密>/<ID>/index.json` |
| 每次启动都下载上号器 | `run.log` 里 `本地上号器需要重新获取：<原因>`（大小/sha 不一致） |
| 更新后没自动重开 | `startup.log` 里 `[restart] ...` 行 |
| `.old` 文件残留 | `startup.log` 里 `旧版本文件暂未释放` |
| 上号器起不来 | `host.dat` 是否被占用/只读；`CG-cli.exe hostinfo -dir <目录>` |

命令行兜底（`CG-cli.exe`）：

```bat
CG-cli.exe push -name leonai      把本地全部账号传到云端
CG-cli.exe pull -name leonai      只补本地缺失的账号
CG-cli.exe cloud -name leonai     看云端账号与标记
CG-cli.exe cloud -list            看所有卡密
CG-cli.exe clean                  清理更新残留
```
