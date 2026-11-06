# 枪神纪 / TCLS 完整逆向分析与目录 Dump 报告

- 分析日期：2026-11-04
- 分析对象：`C:\WeGameApps\枪神纪`（game_id=33）、`C:\WeGame\tcls`、官方 WeGame 启动链
- 分析方式：全目录只读 dump + 字符串挖掘（ASCII / UTF-16LE）+ IDA Pro 9.0 自动分析与 Hex-Rays 反编译 + 官方日志取证 + 运行时共享内存实测抓取
- 样本性质：本机已安装文件、本机运行进程、本机日志；**未修改任何游戏/WeGame 文件，未启动任何游戏进程，未注入任何 DLL**
- 配套原始证据：`C:\Users\Administrator\Desktop\cF\cf-pentest-ym\evidence\tcls_dump\`
- 配套脚本：`C:\Users\Administrator\Desktop\cF\cf-pentest-ym\scripts\`

---

## 0. 结论先行

| # | 结论 | 置信度 |
|---|---|---|
| C1 | 上一轮逆向**找错了组件**。`TCLS_CORE_SHARE_MEMORY_*`（0x49B / 1179B）是 `tcls_core.exe` 内 `Tenio::DNFLauncherByQQLogin`（DNF / QQ 登录分支）的通道，**不是枪神纪 `-src=tgp` 的票据通道**。日志里「票据走落地文件，不依赖此通道」是错误推断。 | 高（IDA 反编译 + 实测） |
| C2 | **账号替换根本没被执行**。生效函数明确禁止修改游戏目录，整套写票据文件 / 改 `mmog_data.xml` / 自启 `TPS.exe` 的代码是**全库无调用点的死代码**；且它假设的 `CF_Gamest` / `TPS_Gamest` 文件在枪神纪目录里**根本不存在**（从 CF 上号器照抄）。 | 高 |
| C3 | **真正的票据通道是 `Local\TCLS_SHAREDMEMEMORY<游戏进程PID>`**，0x80000（524288）字节、字段 ID 索引表，在游戏进程被挂起创建后、`ResumeThread` 之前写入。本机实测抓到 `TCLS_SHAREDMEMEMORY4620`，`4620` 与官方日志 `gameProcessID: 4620` 完全一致。 | 高（实测 + IDA + 日志） |
| C4 | 该共享内存**同时被游戏侧（`TPSBase.dll`）与多个反作弊组件（`ACE-Base32/64.dll`、`TerSafe.dll`、`tersafe32/64.dll`、`SGuard64Base.dll`、`TGuardBase.dll`）、代理组件（`TenProxy.dll` / `TenTcls.dll`）读取**。外部改写会直接进入 ACE 的校验视野。 | 高（字符串证据） |
| C5 | `LoginByQq` 换出的 `openid / open_key / ctt` 与 WeGame 当前登录会话**不同源**。官方链把 WeGame 会话派生的 `OpenID / JumpSignature / launch_key / TGPSignature` 灌进字段表，游戏侧 `CLaunchUI::CheckShareMemoryValue` 逐字段与进程参数/命令行/区服配置交叉校验。因此「写文件 / 改几个共享内存字段」在架构上不成立。 | 高（反编译 + 日志） |
| C6 | 记录的三次测试里**游戏从未成功启动**（`exit code 2` / `exit code 4 / result:false`），没有一次成功登录可供「替换」。 | 高（日志） |
| C7 | 代码里存在多处与「逆向深浅」无关的工程 bug：`OpenProcess` 参数颠倒、`-q` 传错值、不校验启动结果、清理文案与真实状态不符。 | 高（源码） |

> 一句话：**上一轮观测的是 DNF 分支的 1179 字节通道，而枪神纪走的是 512KB 的 `TCLS_SHAREDMEMEMORY<PID>`；同时替换代码被禁用、票据落地文件模型是错的、官方启动本身也失败了——三者叠加，所以「逆向没到位、账号没替换成功」。**

---

# 第一部分 · 完整 Dump

## 1.1 Dump 范围与规模

| 范围 | 文件数 | 说明 |
|---|---:|---|
| `C:\WeGameApps\枪神纪`（递归） | **14,145** | 含 Binaries / Engine / UDKGame / TCLS / WeGameLauncher / components / x64 等 |
| 其中 PE（exe/dll/ocx/sys） | **624** | 全部记录节表、架构、SHA256 |
| ASCII + UTF-16LE 字符串提取 | **400 份** | 每个 PE 一份完整 `.txt` |
| 关键词命中（票据 / API / 文件名） | **93,502** 条 | 游戏目录；另 WeGame\tcls 6,425 条 |
| `C:\WeGame\tcls`（递归） | **70** | 38 个 PE |
| 运行时用户数据（`%APPDATA%\Tencent\*`） | **538** | 含 TCLSCore / WeGame / TPS |
| 官方日志扫描 | 1,048 行命中 | `wegame_log_hits.txt` |
| IDA 反汇编 / 反编译报告 | 5 份，约 470 KB | 含 `tcls_core.exe.i64` 数据库 |

**没有破坏现场**：所有二进制先复制到 `evidence\tcls_dump\bin\` 再做 IDA 分析；IDA 数据库另存，原目录只读。

## 1.2 枪神纪根目录

```text
C:\WeGameApps\枪神纪\
  Binaries\           游戏主程序（Win32 / Win64 / Cross）
  components\         tclsqmfix 等组件
  DirectX\
  Engine\
  rail_files\         Rail 资源
  TCLS\               ★ TCLS 启动器（本报告重点）
  tiny_cache\         (hidden)
  UDKGame\            UE3 脚本 / 配置 / BinData
  WeGameLauncher\     WeGame 侧启动器资源
  x64\
  amdcpusetup.exe                 2,966,160
  枪神纪卸载.exe                     225,376
  枪神纪预发布服卸载.exe             225,376
  说明.rtf                         170,662
```

关键可执行文件（全盘定位结果 `key_exes.txt`）：

| 路径 | 大小 | 角色 |
|---|---:|---|
| `C:\WeGame\client_launcher.exe` | 479,784 | WeGame 启动转发（`SleepOnly`） |
| `C:\WeGame\wegame.exe` | 1,223,216 | WeGame 主进程 |
| `C:\WeGame\tcls\tcls_core.exe` | 1,971,752 | **TCLS 宿主，票据共享内存的创建者**（PE32 / x86，4 节） |
| `C:\WeGameApps\枪神纪\TCLS\Client.exe` | 917,704 | 游戏侧 TCLS 客户端 |
| `C:\WeGameApps\枪神纪\Binaries\Win64\TPS.exe` | 20,232,824 | 游戏主程序 x64 |
| `C:\WeGameApps\枪神纪\Binaries\Win32\TPS.exe` | 11,831,928 | 游戏主程序 x86 |
| `C:\WeGameApps\rail_apps\CF：穿越火线(2)\TCLS\Client.exe` | 2,266,064 | CF 的 TCLS 客户端（**不同**） |
| `C:\WeGameApps\rail_apps\CF：穿越火线(2)\NTCLS\tcls\tcls_core.exe` | 859,704 | CF 的 tcls_core（**不同版本**） |

`tcls_core.exe` PE 指纹：`Machine=0x14C (i386)`、`OptMagic=0x10B (PE32)`、`ImageBase=0x400000`，PDB 路径 `E:\dailybuild_fix\wegame_client\build\bin\Release\tcls_core.pdb`（符号服务器未命中，全部靠 FLIRT + 字符串 + 反编译还原）。

IDA 段表：

```text
0x401000 - 0x556000  .text
0x556000 - 0x5564A8  .idata
0x5564A8 - 0x5A6000  .rdata
0x5A6000 - 0x5BB000  .data
```

文件偏移 → VA 换算：`VA = 文件偏移 + 0x401000`（本样本）。例：`0x156270 → 0x557270 = L"TCLS_CORE_SHARE_MEMORY_"`。

## 1.3 `C:\WeGameApps\枪神纪\TCLS` 顶层文件全清单

| 文件 | 大小 | 文件 | 大小 |
|---|---:|---|---:|
| `Assistant.exe` | 359,368 | `Client.exe` | 917,704 |
| `BackgroundDownloader.exe` | 769,696 | `Repair.exe` | 11,775,624 |
| `bugreport.exe` | 279,968 | `dbghelp.dll` | 1,037,880 |
| `DirInfo3_0.dll` | 1,674,896 | `Lua51.dll` | 212,536 |
| `SSOCommon.dll` | 1,531,448 | `SSOPlatform.dll` | 1,678,904 |
| `tacc_2_1.dll` | 2,309,776 | `TCLSClient.dll` | 56,888 |
| `tinyxml.dll` | 101,944 | `ver_30.dll` | 1,932,944 |
| `flash.ocx` | 15 | `mmog_data.tdr` | 9,878 |
| `tacc_datadef.tdr` | 4,566 | `Tenio.ini` | 480 |
| `version.ini` | 136 | `mmog_data.xml` | 433 |
| `__mmog_data.xml.bak` | 433 | `__mmog_data.xml.tmp` | 350 |
| `b-45f7ee19...-context.log` | 0 | `b-45f7ee19...-flag.log` | 16 |
| `b-45f7ee19...-result.log` | 0 | `b-45f7ee19-result.log` | 0 |

> `__mmog_data.xml.bak` / `.tmp` 是 TCLS 版本更新器自己留下的备份，不是票据文件。

子目录：

```text
TCLS\config\          commonconfig.cfg(4232) dirserver.xml(962) dirserverdynamic.ini(56)
                      dirsystemconfig.cfg(1024) gic_cs.tdr(7242) name2id.xml(2683)
                      prompt.ini(9811) tacc_client_protocol.tdr(9482) uiconfig.xml(1739)
                      updateconfig.cfg(26880) UserImages.bmp(0) VersionServerConfig.xml(662)
                      __commonconfig.cfg.bak/.tmp
                      Dynamic\gamesetting.ini(119)
TCLS\plugins\         ATL80.dll dr.dll frame.ini GameCooperation.zip(4,692,407)
                      Microsoft.VC80.*.manifest msvcm80/msvcp80/msvcr80.dll
                      QMCommon.dll QMIpc.dll QMNetworkMgr.dll QmTCLSPatch.dll
                      QmTCLSPlugin.dll QQPCDetector.dll QQPCfix.dll
TCLS\Tenio\           TenioDL\
TCLS\tenprotect\      反作弊组件
TCLS\ui\              Layout\ res\game\dialog\mibao2\ res\tck\ res\TPS\smzc\
TCLS\music\
```

## 1.4 `C:\WeGame\tcls` 全清单（70 文件，摘要）

```text
tcls_core.exe(1,971,752)   common.dll(4,122,160)   ProtocolCenter.dll(2,089,520)
TSG.dll(67,680)            client_proxy.dll(119,856) adapt_for_imports.dll(413,744)
GameDiffer.dll(387,744)    QQSpeedUpdate.dll(170,008) UpdateInstaller.exe(488,608)
bugreport.exe(4,080,744)   7z.dll(3,293,344)        dbghelp.dll(1,256,080)
ucrtbase.dll(900,288)      normaliz.dll(23,552)    crash_report_helper.dll(3,176,552)
TCLSCoreConfig.ini(992)    Tenio.ini(450)          perfix.ini(17)
tinydl.ini(18)             tinydl_tcls.ini(43)     VFSConfig.ini(184)
bugreport_config.ini(63)   config\*.cfg / *.xml / *.ini
Tenio\                    tcls.dll tenbase.dll TenioTPF.dll June.dll ScriptManager.dll
                          MemoryAlloctor.dll VFS.dll TinyXml.dll TPFCustom.dll
                          TPFAdapter.dll TenFact.dll TenBaseForCross.dll
                          xRender.dll log.dll GdiPlus.dll
Tenio\TenioDL\            TenioDL.dll TenioDL.exe TenioDL_core.dll p2papp.dll
                          p2plog.dll msvcp100.dll msvcr100.dll normaliz.dll
Tenio\Config\             Cursor.xml CustomEvent.xml name2id.xml uiconfig.xml
TLog\                     6 份运行日志（含 20261104_164133.log，即测试窗口）
UIShellRes\               DMLQ1 / QQSPEED / QQXJ 的 TCLS\config
```

完整逐文件清单（含时间戳）：`evidence\tcls_dump\wegame_tcls_files.txt` / `.json`。

## 1.5 关键配置原文

**`TCLS\mmog_data.xml`（433 B，2026-11-04 16:30:44）**

```xml
<?xml version="1.0" encoding="gb2312" ?>
<VersionUpdateData>
    <VersionConf>16:16:16:16</VersionConf>
    <ClientConf>
        <VersionData>
            <Version>1.2.29.652</Version>
            <Uin>0</Uin>
            <WorldID>0.1.3.8</WorldID>
            <BranchId>1</BranchId>
        </VersionData>
        <VersionDemoStr>ver-1.2.29.652</VersionDemoStr>
        <resv1>0.1.2.2</resv1>
    </ClientConf>
</VersionUpdateData>
```

`<Uin>` 恒为 `0` —— 版本/更新配置，**不是登录票据**。这与官方日志 `current version: 281483568546444(1.2.29.652)` 对应。

**`TCLS\version.ini`（UTF-16LE）**

```ini
[LiveUpdateVersion]
version=2016032301
[TCLS]
version=2016032301
```

**`TCLS\Tenio.ini`（UTF-16LE）**

```ini
[Log]
OutputDebugString=E
WriteToLogFile=E
logFileCount=20
[DefaultValue]
DllPath="Tenio"
Dlls=TCLS.dll,TenBase.dll,Dir.dll,TenTPF.dll,MemoryAlloctor.dll,June.dll,log.dll
[361]  Dll=TenFact.dll
[1121] DLL=TinyXml.dll
```

**`TCLS\config\Dynamic\gamesetting.ini`**

```ini
[TPSSetting]
FullScreen=0
ScreenPix=1024*640
ScreenIndex=3
FirstLogin=0
IsDefault=1
MovieIndex=0
WindowIndex=0
```

**`%APPDATA%\Tencent\TCLSCore\UserInfo.ini`（52 B）**

```ini
[USER_INFO]
ACCOUNT_ID=2485853957
UID=3839768809
```

对照日志：11-03 19:08 启动 `uin: 2637938199`，11-04 13:31 / 16:41 启动 `uin: 3839768809`。该文件随登录账号变化，是**链条的输出/缓存**，不是身份输入。

**`C:\WeGame\tcls\TCLSCoreConfig.ini`（992 B，GBK）** —— 按 game_id 分节，决定 `TCLSPath` 与 `TCLSClientFileName`：

```ini
[1]   ; dnf     TCLSPath=TCLS  TCLSClientFileName=Client.exe  UsesTCLSLoginOnly=0
[2]   ; cf
[10]  ; x5      TCLSPath=TCLS  TCLSClientFileName=Client.exe  UsesTCLSLoginOnly=0
[11]  ; x5_2    TCLSPath=.     TCLSClientFileName=X52.exe      UsesTCLSLoginOnly=1
[21]  ; yxd     TCLSPath=TCLS  TCLSClientFileName=Client.exe  UsesTCLSLoginOnly=1
[32]  ; cof     TCLSPath=TCLS  TCLSClientFileName=Client.exe  UsesTCLSLoginOnly=1
[48]  ; bns     TCLSPath=TCLS  TCLSClientFileName=Launcher.exe
[57]  ; hon     TCLSPath=TCLS  TCLSClientFileName=honlogin.exe
[512] ; lcu     TCLSPath=LCU\TCLS  TCLSClientFileName=Client.exe
[99]  ; BN      TCLSPath=.    TCLSClientFileName=dz.exe
```

**`[33]`（枪神纪）不在表里** → 走默认 `TCLSPath=TCLS`、`Client.exe`、`UsesTCLSLoginOnly=0`，与实测 `TCLS\Client.exe` 一致。

## 1.6 证伪：上一轮假想的「票据落地文件」全部不存在

对 `C:\WeGameApps\枪神纪` 递归 glob `*Gamest*` / `*gamest*` / `*Ticket*` / `*ticket*`：

```text
（无 CF_Gamest / TPS_Gamest / op.json / nt.ini / Token010A.txt 之类文件）
只命中：UDKGame\Config\DefaultGameStat.ini、UDKGame\Config\UDKGameStat.ini、
        UDKGame\BinData\RuleTicket.bin、Binaries\Win32|Win64\Cross\mmog_data.xml、
        TCLS\mmog_data.xml、TCLS\version.ini
```

`%APPDATA%\Tencent\TCLSCore\` 下唯一文件是 `UserInfo.ini`（52 B）。

**结论：`tclsTicketPaths()` / `tclsMmogPaths()` / `tclsBuildGamest()` 那套 `k=v` 明文票据文件模型是 CF 上号器的，套到枪神纪完全没有消费方。**

## 1.7 Dump 产物清单

证据根目录 `C:\Users\Administrator\Desktop\cF\cf-pentest-ym\evidence\tcls_dump\`：

| 文件 | 内容 |
|---|---|
| `game_qsj_files.txt` / `.json` | 14,145 文件清单（1.9 MB / 2.6 MB） |
| `game_qsj_pe_manifest.json` | 624 个 PE 的 machine / 节表 / SHA256 |
| `game_qsj_keyword_hits.txt` | 93,502 条命中（11.5 MB） |
| `game_qsj_configs.txt` | 199 份文本配置全文（5.2 MB） |
| `strings_*.txt`（400 份） | 每 PE 的 ASCII + UTF16 字符串及文件偏移 |
| `focus_*.all.txt` / `.hits.txt` | 启动链 15 个关键 PE 的定向扫描 |
| `wegame_tcls_*` | `C:\WeGame\tcls` 清单 / PE / 配置 / 命中 |
| `key_exes.txt` | 全盘关键 exe 定位 |
| `userdata_files.txt` / `userdata_*.txt` | 运行时用户数据清单与内容 |
| `wegame_log_hits.txt` | 官方日志 1,048 行命中 |
| `bin\tcls_core.exe` + `.i64` | IDA 分析副本与数据库（7.8 MB） |
| `bin\TCLS_Client.exe` / `client_launcher.exe` | 只读副本 |
| `ida_tcls2/3/4/5_report.txt` | IDA 反汇编 + Hex-Rays 全集 |
| `ida_keyfuncs.txt` / `ida_struct.txt` | 关键函数摘录 |
| `TCLS_SHAREDMEMEMORY.bin` | **实测 524,288 字节票据结构原始 dump** |
| `TCLS_SHAREDMEMEMORY_fields.json` / `_head.hex` | 字段表解析 / 头部十六进制 |
| `shm_field_map.txt` | `CheckShareMemoryValue` 日志标签 → 文件偏移 / VA |
| `live_shm_probe.txt` / `live_handles.txt` | 运行时探测记录 |

脚本根目录 `C:\Users\Administrator\Desktop\cF\cf-pentest-ym\scripts\`：

```text
tcls_dump.py            全目录 dump + 字符串挖掘
tcls_focus.py           启动链关键 PE 定向扫描
tcls_handles.py         句柄枚举（NtQuerySystemInformation）
tcls_bruteforce_shm.py  共享内存名暴力探测
log_hits.py             官方日志命中提取
find_shm_name.py        容错内容搜索
ida_tcls.py ~ ida_tcls5.py  IDA 批处理脚本（v5 为最终版）
tcls_shm_fields.py      字段标签提取 + 新通道重探
tcls_dump_shm.py        ★ 共享内存 dump + 字段表解析（敏感值脱敏）
```

---

# 第二部分 · 官方到底怎么「用现有数据上号」

## 2.1 完整进程链与数据流

```text
┌──────────────────────────────────────────────────────────────────────┐
│ 用户已登录官方 WeGame（账号 X）                                       │
│  wegame.exe (pid 15240)                                              │
│   └─ wgcore.dll: LoginMgr::RequestGameST / OnQueryGameStRsp          │
│      向 WeGame 后端换取本次游戏的 openid / open_key / pf_key / ST     │
│      （字符串证据见 2.2）                                            │
└───────────────────────────┬──────────────────────────────────────────┘
                            │ 用户点「开始游戏」或
                            │ C:\WeGame\client_launcher.exe -game_id=33
                            ▼
┌──────────────────────────────────────────────────────────────────────┐
│ C:\WeGame\tcls\tcls_core.exe  (pid 16244, ppid=<wegame pid>)          │
│  1) TCLSHostApp::Init / RecvLoginInfo()                              │
│  2) AddGame(gameID=33, path=C:\WeGameApps\枪神纪)                     │
│  3) ServerDir::Fetch → ReadServiceID() = 2103297                      │
│  4) TCLSHostApp::SaveLoginInfoToTCLSApp(sub_460472)                   │
│       SetOpenID / SetNetbarInfo / SetLoaderArguments                   │
│       JumpSignatureLen=72, tgpSingLen=96                              │
│  5) BeginToLaunchGame → GameLauncher::LaunchGameImpl                   │
│  6) sub_4764BF: CreateProcessW(..., CREATE_SUSPENDED)                  │
│       ★ 游戏进程 PID = 4620                                           │
│  7) sub_472FF1 TCLS_Launcher::SaveGameInfo（Rail 分支 sub_474377）    │
│       CreateFileMappingW(L"TCLS_SHAREDMEMEMORY4620", 0x80000, NULL DACL)
│       MapViewOfFile(FILE_MAP_ALL_ACCESS)                              │
│       sub_483F02(view, 64) 初始化 64 项字段索引表                     │
│       sub_483FAC / sub_483FF0 按字段 ID 写入票据                      │
│       UnmapViewOfFile / CloseHandle（句柄由宿主持有）                  │
│  8) ResumeThread(ProcessInformation.hThread)                          │
│  9) 命令行: "…\TPS.exe" -q <UIN> -src=tgp -game_id 33 -area 1        │
│             -zone_id 18220545 <附加参数×3>                            │
└───────────────────────────┬──────────────────────────────────────────┘
                            ▼
┌──────────────────────────────────────────────────────────────────────┐
│ 游戏进程（挂起恢复后）                                                │
│  TPSBase.dll / TenProxy.dll(CTclsProxy::Initialize)                    │
│      → OpenFileMappingW("TCLS_SHAREDMEMEMORY<own PID>")               │
│  ACE-Base32/64.dll、TerSafe.dll、tersafe32/64.dll、                    │
│  SGuard64Base.dll、TGuardBase.dll、TenTcls.dll                        │
│      → 同样打开该共享内存（反作弊消费方）                             │
│  CLaunchUI::CheckShareMemoryValue(sub_478666)                          │
│      → 逐字段与进程参数/命令行/区服配置交叉校验                       │
│  TCLS\Client.exe（TCLSCoreConfig.ini 默认值）                         │
│  TPS.exe 用票据登录游戏服务器                                          │
└──────────────────────────────────────────────────────────────────────┘
```

**核心时序（E2）**：票据在游戏进程**第一条指令执行之前**就写完；游戏与反作弊组件启动后按自己的 PID 拼名字打开映射。这解释了为什么名字要带 PID —— 它是父子进程的 rendezvous key，不是可预测的常量。

## 2.2 WeGame 侧生成票据的字符串证据（`wgcore.dll`，18.7 MB）

```text
0x00c34568  A  openid
0x00c34e04  A  open_key
0x00c3b270  A  third_openid
0x00c34c40  A  [login_mgr]OnQueryGameStRsp, game_id: {}, result:{}, rsp_code:{}, uin: {}, ...
0x00c3c5d8  A  [LoginMgr]RequestGameST, serviceID: %u, openSvrID: %u, UIN: %u, request_access_token:%u
0x00c3c630  A  [LoginMgr]RequestGameST, NT qq user backend request serviceID: {}, ...
0x00c5cc80  A  [gl][game][%llu]LaunchWebAsMMOGGame, open_service_id:%d, open_id:%s, open_key:%s, open_pfkey:%s
```

`wegame.exe` 另有 `TPF_Share_Mem`、`WEGMAE_QBLINK_SHARE_MEMORY_INFO_%d` 等**无关**通道（QoS / 链路），不要与票据通道混淆。

`client_launcher.exe` 含 `SleepOnly` 字符串（对应 `--SleepOnly=1` 转发模式）。

## 2.3 官方日志原文（测试窗口 16:41，来自 `TCLSHost_2026-11-4-1793780698.log`）

```text
[16:24:58][TRACE] Succeed to create mutex which named TCLSHOST_MUTEX_3326F17117BDA9AA14086B4D69E27089
[16:24:59][TRACE] TCLSHostApp::RecvLoginInfo()
[16:24:59][TRACE] TCLSHostApp::RecvLoginInfo(), login ok
[16:41:32][TRACE] [TCLSHostApp::SaveLoginInfoToTCLSApp]
[16:41:32][TRACE] [TCLSHostApp::SaveLoginInfoToTCLSApp] SetOpenID, JumpSignatureLen:72
[16:41:32][TRACE] [TCLSApp::SetOpenID]
[16:41:32][TRACE] [TCLSApp::SetNetbarInfo]
[16:41:32][TRACE] [TCLSApp::SetLoaderArguments]
[16:41:32][TRACE] ServerDir::Fetch() tclsPath: C:\WeGameApps\枪神纪\TCLS, gameID: 33
[16:41:32][TRACE] ServerDir::ReadServiceID() ServiceID = 2103297
[16:41:33][TRACE] TCLSHostApp::OnFetchServerDirEnd(), gameID: 33, result: true, forTgp: false
[16:41:33][TRACE] TCLSHostApp::BeginToLaunchGame()
[16:41:33][TRACE] ServerDir::GetSelectedAPPID() Selected serverid = 18220545, AppID = 0
[16:41:33][TRACE] GameLauncher::LaunchGame(), gameID: 33, server_id: 1537, uin: 3839768809,
                  tclsPath: C:\WeGameApps\枪神纪\TCLS, appID: 1, tgpSingLen: 96
[16:41:33][TRACE] GameLauncher::LaunchGameImpl(), gameID: 33, server_id: 1537, uin: 3839768809, ...
[16:41:33][TRACE] Use TCLS 10 launch game:C:\WeGameApps\枪神纪\TCLS
[16:41:33][TRACE] GameLauncher::Start(), branch:1
[16:41:33][ERROR][TCLS]Begin launch, game id: 33 AppID: 1. ServiceID: 1537, UIN: 3839768809 TCLS: ...
[16:42:22][ERROR] GameLauncher::OnNotifyExit(), exit code: 2
[16:42:22][ERROR] TCLSHostApp::OnLaunchEnd(), gameID: 33, gameProcessID: 4620,
                  result: true, errorCode: 2, tp_download_result: 0, tp_download_time: 0
[16:42:27][TRACE] TCLSHostApp::OnGameExit(), game_id: 33
```

其它时间点的失败记录：

```text
[19:08:43] (11-03) [TCLS]Begin launch ... UIN: 2637938199  → [19:09:43] exit code: 2
[13:31:41] (11-04) [TCLS]Begin launch ... UIN: 3839768809  → [13:32:17] OnNotifyExit exit code: 4,
                   OnLaunchEnd result: false, errorCode: 4
```

> **注意 `uin` 是 TGP UID（`UserInfo.ini` 的 `UID=3839768809`），不是 QQ 号（`ACCOUNT_ID=2485853957`）。**

## 2.4 IDA 关键函数表（`tcls_core.exe`，RVA 即 VA，ImageBase 0x400000）

| 地址 | 名称 / 作用 |
|---|---|
| `0x40BB57` | `Tenio::DNFLauncherByQQLogin` 析构（vftable 定位） |
| `0x40BC15` | DNF 分支入口：`qmemcpy(this+24, Src, 0x49B)` → 建 mutex → `WriteLoginInfo` |
| `0x40BF99` | `CreateProcessW("start\DNFchina.exe")` ← **证明 1179B 通道属于 DNF** |
| `0x40C370` | `DNFLauncher::WriteLoginInfo`：建 0x49B 映射 + memmove |
| `0x40C5B7` | 名字构造：`L"%s%s"` = `TCLS_CORE_SHARE_MEMORY_` + 后缀 |
| `0x40C6BC` | mutex 名字构造：`TCLS_CORE_MUTEX_` + 后缀 |
| `0x40C92D` | 从结构 `+666` 取窄串构造 `std::string`（`strlen`） |
| `0x44B8A5` | `MultiByteToWideChar(CP_UTF8=65001, ...)` UTF-8 → wstring |
| `0x44BBEE` | `sub_44B8A5` 包装 |
| `0x44DFD2` | `TCLSApp::SetOpenID`（`memmove(a1+4836, Src, 0xA2)`） |
| `0x4521D3` / `0x45DB68` | `USER_INFO` / `ACCOUNT_ID` 相关读写 |
| `0x45C863` | 创建 `TCLSHOST_MUTEX_<32hex>` |
| `0x460472` | `TCLSHostApp::SaveLoginInfoToTCLSApp`（把会话信息推给 TCLS 组件） |
| `0x472FF1` | **`TCLS_Launcher::SaveGameInfo`：创建 0x80000 票据共享内存** |
| `0x474377` | `TCLS_Launcher::RailSaveGameInfo`（Rail 分支，同一通道另一套写入原语） |
| `0x4764BF` | **挂起创建游戏进程 → SaveGameInfo → ResumeThread** |
| `0x478666` | `CLaunchUI::CheckShareMemoryValue`：逐字段交叉校验 |
| `0x480084` | `L"%d"` 数字 → 宽串（**共享内存名后缀**） |
| `0x48210D` | `wcscat_s` 拼接名字后缀 |
| `0x483F02` | 字段表头初始化（64 项） |
| `0x483FAC` / `0x483FF0` | 按字段 ID 写 DWORD / 字节块 |
| `0x483D4A` / `0x483DB8` / `0x483E70` | 按字段 ID 读 DWORD / 字节块 / 宽串（校验侧） |
| `0x48443A` | 按 `序号 + 字段 ID` 读（区服数组校验） |
| `0x44E799` | `_WinMain@16`（`ppid` 处理） |

## 2.5 关键反编译片段（Hex-Rays 原文摘录）

**(a) `sub_40C370` — 1179B 通道（DNF 分支）**

```c
char __thiscall sub_40C370(HANDLE *this, char *Src)
{
  ...
  if ( !WaitForSingleObject(this[5], 0x2710u) ) {
    InitializeSecurityDescriptor(pSecurityDescriptor, 1u);
    SetSecurityDescriptorDacl(pSecurityDescriptor, 1, 0, 1);   // NULL DACL
    FileMappingAttributes.nLength = 12;
    FileMappingAttributes.lpSecurityDescriptor = pSecurityDescriptor;
    sub_40C92D(Src + 666);                 // ★ 结构 +666 处的 UTF-8 串
    v12 = sub_44BBEE(v17, v18);
    sub_40C5B7(v19, v12);                  // 拼名字
    sub_46F23B(&off_5571D8, v2);           // 日志 shareMemoryName: %s
    lpName = (LPCWSTR)sub_40A631(v19);
    v3 = CreateFileMappingW((HANDLE)0xFFFFFFFF, &FileMappingAttributes,
                            4u, 0, 0x49Bu, lpName);
    this[301] = v3;
    if ( this[301] ) {
      lpBaseAddress = MapViewOfFile(this[301], 0xF001Fu, 0, 0, 0);
      if ( lpBaseAddress ) {
        memmove((void *)lpBaseAddress, Src, 0x49Bu);   // 整块 1179B 拷入
        UnmapViewOfFile(lpBaseAddress);
      }
    }
    ReleaseMutex(this[5]);
  }
  ...
}
```

**(b) `sub_40C5B7` — 名字 = 前缀 + 结构内令牌**

```c
int __stdcall sub_40C5B7(int a1, int a2)
{
  sub_4857B8(v3, a2);
  memset(Buffer, 0, sizeof(Buffer));       // wchar_t[512]
  sub_40A631(v3);                          // c_str()
  sub_4077F0(Buffer, 0x200u, L"%s%s", L"TCLS_CORE_SHARE_MEMORY_", suffix);
  sub_409093(Buffer);                      // 赋给输出 wstring
  return a1;
}
```
汇编实证：
```asm
0040C623  call    sub_40A631          ; 后缀 c_str
0040C628  push    eax
0040C629  push    offset off_557270  ; L"TCLS_CORE_SHARE_MEMORY_"
0040C62E  push    offset aSS_0       ; L"%s%s"
0040C633  push    200h
0040C638  lea     eax, [ebp+Buffer]
0040C63F  call    sub_4077F0          ; swprintf 包装
```

**(c) `sub_40BC15` — 该分支的下游是 DNF**

```c
qmemcpy((void *)(this + 24), Src, 0x49Bu);
...
if ( sub_40C370((HANDLE *)this, Src) ) {
   sub_40C92D(v11, Src + 666);
   v4 = sub_44BBEE(v10, v11);
   v5 = (unsigned __int8)sub_40BF99(v4) == 0;   // ← 起 start\DNFchina.exe
   ...
}
```

**(d) `sub_472FF1` — 真正的 512KB 通道**

```c
memset(Name, 0, sizeof(Name));            // WCHAR Name[260]
if ( *((_DWORD *)v51 + 140105) ) {        // 网吧分支
    wcscpy_s(Name, 0x104u, L"NETBAR_LAUNCHGAMESHAREDMEM");
    sub_480084(*((_DWORD *)v51 + 125849), Source, 0x104u);
    sub_48210D(Name, 260, Source);
} else {                                  // ★ TGP 分支
    wcscpy_s(Name, 0x104u, L"TCLS_SHAREDMEMEMORY");   // 官方拼写就是 MEMEMORY
    memset(Destination, 0, sizeof(Destination));
    sub_480084(ArgList, Destination, 0x104u);         // L"%d" 数字后缀
    sub_48210D(Name, 260, Destination);
    sub_480A0E(v51 + 105, &off_5647FC, Name);         // 打日志
}
InitializeSecurityDescriptor(...);
SetSecurityDescriptorDacl(pSecurityDescriptor, 1, 0, 1);   // NULL DACL
hFileMappingObject = CreateFileMappingW((HANDLE)0xFFFFFFFF,
                        &FileMappingAttributes, 4u, 0, 0x80000u, Name);
...
lpBaseAddress = MapViewOfFile(hFileMappingObject, 0xF001Fu, 0, 0, 0);
if ( lpBaseAddress ) {
    sub_483F02(lpBaseAddress, 64);                    // 64 项字段表
    sub_483FAC(1, *((_DWORD *)v51 + 125849), 4);
    sub_483FAC(2, *((_DWORD *)v51 + 39),    4);
    sub_483FAC(6, *((_DWORD *)v51 + 148),   4);
    sub_483FF0(3, v51 + 204, 1u);
    sub_483FF0(4, v51 + 205, (unsigned __int8)v51[204]);
    sub_483FF0(5, v51 + 160, 0x10u);
    ...
}
```

**(e) `sub_480084` / `sub_48210D` — 后缀是十进制数字**

```c
int __stdcall sub_480084(int ArgList, wchar_t *Destination, int SizeInWords)
{
  sub_4077F0(Source, 0x104u, L"%d", ArgList);
  wcscpy_s(Destination, SizeInWords, Source);
}
int __stdcall sub_48210D(wchar_t *Destination, int a2, wchar_t *Source)
{
  wcscat_s(Destination, a2, Source);
}
```

**(f) `sub_4764BF` — 挂起 → 写入 → 恢复**

```c
CreateProcessW(..., CREATE_SUSPENDED, ..., &StartupInfo, &ProcessInformation);
...
if ( *((_DWORD *)v65 + 140120) )
    v64 = sub_474377(v65, v65 + 503400, ..., v56, v55);   // RailSaveGameInfo
else
    v64 = sub_472FF1(v65, v65 + 503400, ..., v56, v55);   // SaveGameInfo
...
if ( v64 )
    v21 = ResumeThread(ProcessInformation.hThread);
```

**(g) 命令行模板（`tcls_core.exe` 字符串 @ 文件偏移 `0x161940`）**

```text
"%s" -q %u -src=tgp -game_id %llu -area %d -zone_id %u %s %s %s
```

对应实测：`gameID: 33`、`server_id: 1537`、`uin: 3839768809`、`Selected serverid = 18220545`、`AppID = 1`、`Use TCLS 10 launch game`。

## 2.6 谁在读这个共享内存（关键）

全目录字符串扫描，`TCLS_SHAREDMEMEMORY` 出现在下列二进制中：

| 二进制 | 位置 | 证据字符串 |
|---|---|---|
| `C:\WeGame\tcls\tcls_core.exe` | 0x15824C 等 6 处 | `TCLS_SHAREDMEMEMORY` / `NETBAR_LAUNCHGAMESHAREDMEM` |
| `TenProxy.dll` | 0x075B8 / 0x075E0 / 0x07690 | `[CTclsProxy::Initialize] OpenFileMapping TCLS_SHAREDMEMEMORY+%d fail!Error code = %d` / `... Success!` |
| `TenTcls.dll` | 0xD0A50 | `TCLS_SHAREDMEMEMORY` / `[CTclsProxy] OpenFileMapping` |
| `TPSBase.dll` | 0x325AE8 / 0x325B60 | `TCLS_SHAREDMEMEMORY`（游戏侧） |
| `ACE-Base32.dll` / `ACE-Base64.dll` | 0x2AE8B0 / 0x31D308 | `TCLS_SHAREDMEMEMORY` + `OpenFileMappingW` |
| `TerSafe.dll` / `tersafe32.dll` / `tersafe64.dll` | 各自 | `TCLS_SHAREDMEMEMORY` + `OpenFileMappingA` |
| `SGuard64Base.dll` | 0x245A78 / 0x24B4B6 | `TCLS_SHAREDMEMEMORY`、`TenProtect::ShareMemory` |
| `TGuardBase.dll` | 0x7C08CC / 0x7C09D8 | `TCLS_SHAREDMEMEMORY` + 多个 `OpenFileMappingW/A` |
| `TPSMusicBase.dll` | 0x611560 | `TCLS_SHAREDMEMEMORY` |

> `TenProxy.dll` 的日志字符串 `TCLS_SHAREDMEMEMORY+%d` 里的 `+` 只是消息里的分隔符；实际对象名不含 `+`（实测 `TCLS_SHAREDMEMEMORY4620` 打开成功）。

**这是「写共享内存必然被 ACE 看到」的直接证据**：反作弊模块本身就是这个映射的消费方。

## 2.7 `UserInfo.ini` / `mmog_data.xml` 的真实角色

| 文件 | 读/写方 | 作用 | 能否决定登录身份 |
|---|---|---|---|
| `%APPDATA%\Tencent\TCLSCore\UserInfo.ini` | `tcls_core.exe`（`USER_INFO` / `ACCOUNT_ID` 字符串在 `sub_45DB68` / `sub_4521D3` 附近） | 记录本次登录的 `ACCOUNT_ID` / `UID`，供 UI 与下次预填 | **否**，是输出 |
| `TCLS\mmog_data.xml` | TCLS 版本更新器（`TCLS_Launcher::*` 多处引用） | `<Version>` / `<Uin>0</Uin>` / `<WorldID>` 版本与区服配置 | **否** |
| `TCLS\version.ini` | 更新器 | 2016032301 版本号 | 否 |
| `TCLSCoreConfig.ini` | `tcls_core.exe` | 按 game_id 决定 `TCLSPath` / `TCLSClientFileName` | 否 |
| `CF_Gamest` / `TPS_Gamest` | **不存在** | 无 | — |

---

# 第三部分 · 票据共享内存完整结构（实测）

## 3.1 两条通道对比

| 项 | 通道 A（枪神纪 TGP，实际使用） | 通道 B（DNF / QQ 登录分支） |
|---|---|---|
| 名字 | `TCLS_SHAREDMEMEMORY<游戏进程PID>`<br>网吧：`NETBAR_LAUNCHGAMESHAREDMEM<数字>` | `TCLS_CORE_SHARE_MEMORY_<结构+666 的 UTF-8 令牌>` |
| 大小 | `0x80000` = 524,288 B | `0x49B` = 1,179 B |
| 创建者 | `sub_472FF1` / `sub_474377`（`tcls_core.exe`） | `sub_40C370`（同 exe 的 DNFLauncher 类） |
| 写入方式 | 64 项字段 ID 索引表，按 ID 逐字段写 | `memmove(view, Src, 0x49B)` 整块拷 |
| 时机 | 游戏进程 `CREATE_SUSPENDED` 之后、`ResumeThread` 之前 | DNF 分支建 mutex 后，`start\DNFchina.exe` 之前 |
| 校验 | `CLaunchUI::CheckShareMemoryValue` 逐字段比对 | `CheckShareMemoryValue` 同源函数族 |
| 配套 mutex | `TCLSHOST_MUTEX_<32hex>`（实测 `3326F17117BDA9AA14086B4D69E27089`） | `TCLS_CORE_MUTEX_<结构+690 串>`、`_TGP_EXISTS_MUTEX_NAME_` |
| 安全属性 | `InitializeSecurityDescriptor` + `SetSecurityDescriptorDacl(...,1,0,1)` = **NULL DACL** | 同 |
| 消费方 | 游戏 `TPSBase.dll` + `TenProxy.dll` + ACE / TerSafe / SGuard / TGuard | DNF 启动器 |
| 本机实测 | **抓到 `Local\TCLS_SHAREDMEMEMORY4620`，dump 成功** | 960 次候选名探测全部 0 命中 |

## 3.2 字段表二进制格式

由 `sub_483F02` / `sub_483FAC` / `sub_483FF0` 与读侧 `sub_483D4A` / `sub_483DB8` / `sub_483E70` 互证：

```text
偏移        内容
+0          DWORD  fieldCount = 64
+4          DWORD  offsets[64]     // offsets[id-1] = 字段数据区相对 view 起点的偏移
+0x104      BYTE   data[]          // 每字段: DWORD size + size 字节 payload
...         （总映射 0x80000，尾部留空）
```

反编译原文：

```c
// sub_483F02  初始化
this[5] = a3;                    // 64
*this = a2;                      // view 基址
*(_DWORD *)*this = 64;           // +0 写 fieldCount
this[1] = a2 + 4;                // 偏移表起点
this[2] = this[1] + 4 * a3;      // 数据区起点 (+260)
this[3] = this[2];               // 当前写指针

// sub_483FAC  写 DWORD
sub_4844D4(a2);                  // 记录 offsets[id-1]
v5 = (_DWORD *)this[3];
*v5 = 4; v5[1] = a3;            // [size=4][value]
this[3] += 8;

// sub_483FF0  写字节块
sub_4844D4(a2);
v5 = (_DWORD *)this[3];
*v5 = Size; memmove(v5 + 1, Src, Size);
this[3] += Size + 4;

// sub_483D4A  读 DWORD（校验侧）
v5 = (_DWORD *)(*(_DWORD *)(this[1] + 4*a2 - 4) + *this);
if ( *v5 == 4 ) { *a3 = v5[1]; return 1; }

// sub_483DB8  读字节块
v6 = (size_t *)(*(_DWORD *)(this[1] + 4*a2 - 4) + *this);
v7 = *v6;  memmove(a3, v6 + 1, min(*a4, v7));

// sub_483E70  读宽串
v6 = *(_DWORD *)(this[1] + 4*a2 - 4) + *this;
v7 = *(_DWORD *)v6;  wcsncpy(Destination, (wchar_t*)(v6+4), v7 >> 1);
```

## 3.3 实测字段表（`TCLS_SHAREDMEMEMORY4620`，敏感值脱敏）

用 `OpenFileMappingW(FILE_MAP_READ)` + `MapViewOfFile` 读取 0x80000 字节，按上述格式解析：

| 字段 ID | 数据区偏移 | size | 非零字节 | 含义 / 对应证据 |
|---:|---:|---:|---:|---|
| 1 | `0x104` | 4 | 3 | ServiceID = `2103297`（日志 `ReadServiceID() ServiceID = 2103297`） |
| 2 | `0x10C` | 4 | 4 | UIN = `3839768809`（日志 `uin: 3839768809`） |
| 3 | `0x11C` | 1 | 1 | 长度标志 byte |
| 4 | `0x121` | 120 | 0 | 票据块（本次为空） |
| 5 | `0x19D` | 16 | 0 | 16B 密钥块（本次为空） |
| 6 | `0x114` | 4 | 4 | `_time64` 时间戳 = `1793781692`（2026-11-04） |
| 7 | `0x1612` | 96 | 96 | `launch_key`（对应日志 `tgpSingLen: 96`） |
| 14 | `0x168A` | 68 | 36 | 宽串路径 |
| 15 | `0x16DA` | 20 | 16 | 宽串路径/参数 |
| 16 | `0x16F2` | 2 | 0 | 空 |
| 17 | `0x13C7` | 4 | 1 | `2` |
| 18 | `0x13CF` | 4 | 1 | `1` |
| 19 | `0x13D7` | 36 | 23 | 宽串 |
| 20 | `0x1450` | 5208 | 340 | 扩展 blob |
| 21 | `0x1306` | 1 | 0 | 长度前缀 |
| 22 | `0x130B` | 0 | 0 | 空 |
| 23 | `0x130F` | 1 | 0 | 长度前缀 |
| 24 | `0x1314` | 0 | 0 | 空 |
| 25 | `0x1318` | 16 | 0 | LoginKey 槽（本次为空） |
| 26 | `0x132C` | 1 | 1 | 标志 |
| 27 | `0x1331` | 72 | 70 | **JumpSignature**（对应 `JumpSignatureLen: 72`） |
| 32 | `0x137D` | 4 | 1 | `4` |
| 33 | `0x1385` | 62 | 41 | 宽串 |
| 34 | `0x13FF` | 4 | 1 | `37` |
| 35 | `0x1407` | 37 | 37 | GUID 形态启动令牌（`D696A827-…-…`，已脱敏） |
| 36 | `0x1430` | 4 | 1 | `2` |
| 37 / 40 / 41 / 42 | — | 4 | 0 | 0 |
| 43 | `0x1B1` | 1 | 1 | 标志 |
| 44 | `0x1B6` | 120 | 118 | 票据主体 |
| 45 | `0x232` | 4 | 1 | **`33` = game_id** |
| 46 | `0x124E` | 32 | 32 | 32 位十六进制摘要（已脱敏） |
| 47 | `0x12F6` | 4 | 1 | `32` |
| 48 | `0x1272` | 128 | 32 | 块 |
| 49 | `0x12FE` | 4 | 1 | `128` |
| 50 | `0x23A` | 4 | 1 | `1` |
| 56 | `0x242` | 4096 | 585 | 扩展数据 |
| 57 | `0x1246` | 4 | 2 | **`55555` = gid** |
| 60 | `0x1676` | 16 | 16 | 尾部 |

首 512 字节十六进制见 `TCLS_SHAREDMEMEMORY_head.hex`，原始 512KB 见 `TCLS_SHAREDMEMEMORY.bin`，结构化结果见 `TCLS_SHAREDMEMEMORY_fields.json`。

> 映射能被 `OpenFileMappingW(FILE_MAP_READ)` 从外部进程打开，是因为创建时用了 NULL DACL；**但「能读」不等于「能改后不被发现」** —— ACE/TerSafe/SGuard 是同一映射的消费方（见 2.6）。

## 3.4 `CheckShareMemoryValue` 的字段 ID → 校验项映射

`sub_478666` 反编译里逐项读字段并与 `this` 对象内对应成员比较，日志标签位于 `.rdata`：

| 字段 ID | 读取原语 | 比较对象（对象内偏移） | 日志标签（文件偏移 / VA） |
|---:|---|---|---|
| 2 | `sub_483D4A(2,&v)` | `*(DWORD*)(v21+39*4)` | `[CLaunchUI::CheckShareMemoryValue]QQ` (0x164730 / 0x565730) |
| 1 | `sub_483D4A(1,&v)` | `*(DWORD*)(v21+125849*4)` | `...ServiceID` (0x164798 / 0x565798) |
| 18 | `sub_483D4A(18,&v)` | `*(DWORD*)(v21+1300*4)+1` | (0x164820 / 0x565820) |
| 17 | `sub_483D4A(17,&v)` | `*(DWORD*)(v21+1299*4)` | (0x1648B0 / 0x5658B0) |
| 6 | `sub_483D4A(6,&v)` | `*(DWORD*)(v21+148*4)` | (0x164930 / 0x565930) |
| 3 | `sub_483DB8(3,&v16,&v20)` | `v21[204]` | (0x1649A8 / 0x5659A8) |
| 4 | `sub_483DB8(4,v25[128],&v20)` | `v21[205..]` | (0x164A18 / 0x565A18) |
| 5 | `sub_483DB8(5,v30[16],&v20)` | `v21[160..175]` | (0x164A70 / 0x565A70) |
| 21 | `sub_483DB8(21,&v15,&v20)` | `v21[4984]` | `...OpenID` (0x164B10 / 0x565B10) |
| 22 | `sub_483DB8(22,v29[16],&v20)` | `v21[4985..]` | `...OpenID` (0x164B88 / 0x565B88) |
| 23 | `sub_483DB8(23,&v14,&v20)` | `v21[5001]` | `...OpenID` (0x164BE8 / 0x565BE8) |
| 24 | `sub_483DB8(24,v24[128],&v20)` | `v21[5002..]` | `...OpenID` (0x164C80 / 0x565C80) |
| 25 | `sub_483DB8(25,v28[16],&v20)` | `v21[5130..]` | `...OpenID LoginKey` (0x164D08 / 0x565D08) |
| 7 | `sub_483DB8(7,v23[512],&v20)` | 与 `a2/a3`（命令行参数）比较 | `...NetBar Signature` (0x164D78 / 0x565D78) |
| 60 | `sub_483DB8(60,v27[16],&v20)` | `*(DWORD*)(v21+140639*4)` | `...launch_key size invalid!` (0x164E08 / 0x565E08) |
| 19 / 15 / 14 | `sub_483E70` 宽串 | `v21+2602` / `+244298` / `+244558` (wchar*) | 0x164E98 / 0x164EF8 / 0x164F58 |
| 32 | `sub_483D4A(32,&v7)` | `*(DWORD*)(v21+139780*4)` | 0x165FB8 / 0x566070 |
| 33 | `sub_483E70(33, Destination[26000], 52000)` | `v21+253560` (wchar*) | 0x1660C8 / 0x566130 |
| 3 / 5 / 6（按区服序号循环） | `sub_48443A(id, &v, jj+1)` | `v21[3744*jj + 7916/7920/7924]` | 0x1661A8 / 0x166248 / 0x1662E0 |

**任一不符即 `v17 = 0`（校验失败）并打日志。** 这是「只改某几个票据字段」过不去的根本原因。

## 3.5 互斥体与窗口名（完整命名族）

来自 `tcls_core.exe` 的 UTF-16 字符串与实测日志：

```text
TCLS_SHAREDMEMEMORY          (+ 十进制数字 = 游戏进程 PID)
NETBAR_LAUNCHGAMESHAREDMEM   (+ 数字，网吧分支)
TCLS_CORE_SHARE_MEMORY_      (+ 结构 +666 的 UTF-8 令牌，DNF 分支)
TCLS_CORE_MUTEX_             (+ 结构 +690 的串，DNF 分支)
TCLSHOST_MUTEX_              (+ 32 位十六进制，实测 3326F17117BDA9AA14086B4D69E27089)
TCLSHOST_UPDATE_MUTEX_
_TGP_EXISTS_MUTEX_NAME_
TCLS_CORE_WND_%u  /  gTCLS_CORE_WND_%u  /  ITCLS_CORE_WND_%u
_TCLS_CORE_APP_WND_NAME_
Rail.SharedMemory.%u / rail_pipe_15 ...（CF 族，与枪神纪无关）
WeGameIPCShareMem10 / TPF_Share_Mem / WEGMAE_QBLINK_SHARE_MEMORY_INFO_%d（WeGame 其他模块）
```

## 3.6 实测抓取过程（可复现）

```text
1) 用 tcls_bruteforce_shm.py 对旧假设 TCLS_CORE_SHARE_MEMORY_* 做 960 次组合探测 → 0 命中
2) 从 IDA sub_472FF1 得到新名 TCLS_SHAREDMEMEMORY + %d
3) 从官方日志 OnLaunchEnd 得到 gameProcessID: 4620
4) 用 tcls_shm_fields.py 重探：
     FOUND TCLS_SHAREDMEMEMORY4620
     FOUND Local\TCLS_SHAREDMEMEMORY4620
     FOUND Session\1\TCLS_SHAREDMEMEMORY4620
5) 用 tcls_dump_shm.py 以 FILE_MAP_READ 打开并 dump 0x80000 字节 → 解析出 64 项字段表
```

---

# 第四部分 · 为什么上一轮「逆向没到位 / 账号没替换成功」

| # | 问题 | 证据 | 性质 |
|---|---|---|---|
| 1 | 观测对象整个错了：探 `TCLS_CORE_SHARE_MEMORY_*`（1179B、DNF 分支），实际通道是 `TCLS_SHAREDMEMEMORY<PID>`（512KB） | C1 / 第三章；960 次旧名 0 命中、新名一次命中 | 逆向结论错误 |
| 2 | 候选后缀也错了：代码猜 `"" / "33" / "2" / PID / 888PID`；真后缀是**游戏进程 PID**，DNF 分支后缀则是结构内 UTF-8 令牌 | `tps_tcls_windows.go:33-52`；`sub_480084` | 名字空间错误 |
| 3 | **替换代码是死代码**：`tclsWriteTicketFiles` / `tclsStartGame` / `tpsWGLaunchOffline` 全库无调用点；生效函数写明「不修改游戏目录、不自行启动进程」 | `tps_tcls.go:185-191`；`tps_wg.go:261` | 设计上就没做 |
| 4 | 落地文件模型从 CF 照抄：`CF_Gamest` / `TPS_Gamest` 在枪神纪不存在；改 `<Uin>` 无效（该字段恒 0，是版本配置） | 1.5 / 1.6 | 模型错误 |
| 5 | `OpenProcess` 参数写反，恒「已强制结束 0 个进程」 | `tps_wg_windows.go:210` | 代码 bug |
| 6 | `tpsWGStartGame` 只 `cmd.Start()`，不 Wait、不记 pid、不验证游戏进程是否起来 | `tps_wg_windows.go:280-285` | 无启动校验 |
| 7 | `-q` 传 QQ 号，官方传 TGP UID（两者不同） | `tps_tcls_windows.go:152` vs 日志 `uin: 3839768809` / `ACCOUNT_ID=2485853957` | 参数错误 |
| 8 | 清理文案「票据文件已还原」是硬编码，实际既没写也没还原 | `tps_wg.go:250-255` | 日志撒谎 |
| 9 | `tpsWGDisableSwitches()` 清 `Local\ValWgInjectAcc_v1`，但 TCLS 流程从未创建该映射 → 空操作 | `tps_wg_windows.go:148-172` | 遗留死路径 |
| 10 | WEG 注入方案自诊断「60 秒未拦到 cmd100 帧，这套 hook 对枪神纪无效」 | `tps_wg.go:334-338` | 方案不适用 |
| 11 | 观测 60 秒超时后把「没抓到」写成「票据走落地文件」结论 | `tps_tcls_windows.go:119` | 由失败反推的错误结论 |
| 12 | 记录中游戏从未成功启动（`exit code 2` / `4`），没有一次成功登录可供替换 | 2.3 | 前置条件不成立 |

---

# 第五部分 · 为什么「账号替换」在这条链上架构级不成立

1. **票据不是从卡密平台那侧直接喂给游戏的。**
   `tpsWGGetOP()` 用 `LoginByQq` 换出的 `openid / open_key / ctt` 是**协议层票据**；而官方链里进入 512KB 字段表的是 **WeGame 当前登录会话**派生出来的 `OpenID / JumpSignature / launch_key / TGPSignature`（`TCLSHostApp::SaveLoginInfoToTCLSApp` → `SetOpenID` / `SetLoaderArguments`，日志 `JumpSignatureLen:72`、`tgpSingLen:96`）。两者不同源，也不在同一层。

2. **父进程把会话信息写死在挂起期。**
   `sub_4764BF` 在 `CREATE_SUSPENDED` 与 `ResumeThread` 之间调用 `SaveGameInfo` 写映射。外部进程即使拿到句柄，也只能在游戏恢复之后改；此时 `CheckShareMemoryValue` 与 ACE 已经开始按 PID 打开并校验同一块内存。

3. **字段级一致性校验。**
   `CheckShareMemoryValue` 把字段 1/2/3/4/5/6/7/14/15/17/18/19/21/22/23/24/25/32/33/60 以及按区服序号的 3/5/6，逐项与对象内成员、命令行参数、区服配置比对。单独替换 `openid` / `open_key` / `ctt` 会立刻造成多处不一致。

4. **反作弊是同一映射的消费方。**
   `ACE-Base32/64.dll`、`TerSafe.dll`、`tersafe32/64.dll`、`SGuard64Base.dll`、`TGuardBase.dll` 都含 `TCLS_SHAREDMEMEMORY` + `OpenFileMapping*`。因此「改共享内存」不是绕过校验，而是直接改在校验者眼前。

5. **写游戏目录 / 自启 TPS.exe 已被实测拒绝。**
   代码注释记录 v1.3.21 实测：写 `CF_Gamest` / 改 `mmog_data.xml` / `UserInfo.ini`，或用降权令牌自拉 `TPS.exe`，都被 ACE 判为「安全组件运行异常 (13,131083,140)」。

6. **因此可实现面只剩「改变 WeGame 自身的登录会话」。**
   那需要在 WeGame 进程内完成登录态替换 / 注入 hook，属于绕过第三方游戏的身份认证与反作弊控制。

> 本报告到此为止：完整的目录 dump、官方链还原、票据结构、失败根因已全部给出；**不提供账号替换 / ACE 绕过 / 校验逻辑 patch / DLL 注入的实现**。

---

# 第六部分 · 与逆向无关的源码 bug 清单（可直接修）

| 位置 | 问题 | 建议 |
|---|---|---|
| `cf-cache-tool/tps_wg_windows.go:210` | `openProc.Call(uintptr(p.pid), 0, 0x0001)` 参数顺序反了。`OpenProcess(dwDesiredAccess, bInheritHandle, dwProcessId)` 应为 `openProc.Call(0x0001, 0, uintptr(p.pid))` | 修正顺序 |
| `cf-cache-tool/selfop_launch.go:265, 297` | 同类 `OpenProcess` 参数颠倒（复制传播） | 一并修正 |
| `cf-cache-tool/tps_tcls_windows.go:152` | `-q` 用 QQ 号；官方是 TGP UID | 换 `op.Tgpid` / `op.Uin` 并核对语义 |
| `cf-cache-tool/tps_wg_windows.go:280-285` | 启动后不 Wait、不记 pid、不检查 `tps.exe` / `Client.exe` 是否出现 | 增加进程与退出码校验，日志输出真实结果 |
| `cf-cache-tool/tps_wg.go:250-255` | 「票据文件已还原」是硬编码文案 | 改为按实际写入/还原计数输出 |
| `cf-cache-tool/tps_wg_windows.go:148-172` | `tpsWGDisableSwitches` 操作未创建的 `ValWgInjectAcc_v1` | 删除或仅在注入路径调用 |
| `cf-cache-tool/tps_tcls_windows.go:119` | 超时文案把「未捕获」写成机制结论 | 改为「本次未命中候选名，不能推断通道」 |
| `cf-cache-tool/tps_tcls.go:39-166` | 整套票据落地文件备份/还原为死代码 | 删除或标注废弃，避免误导 |
| `cf-cache-tool/tps_wg.go:261-345` | `tpsWGLaunchOffline`（WEG 注入）无调用点且自诊断无效 | 删除或标注废弃 |
| `cf-cache-tool/tps_tcls_windows.go:30-52` | 共享内存候选名基于错误通道 | 若保留只读观测，改为从游戏进程 PID 拼 `TCLS_SHAREDMEMEMORY<PID>` |

> 本轮**未修改任何源码，未重新打包/发布**（没有代码改动可发版）。

---

# 第七部分 · 证据索引（E#）

| 编号 | 发现 | 证据 | 复现 |
|---|---|---|---|
| E1 | 官方启动链与失败退出码 | `TCLSHost_2026-11-4-1793780698.log` 16:41:32–16:42:27 | `python scripts/log_hits.py` |
| E2 | 挂起创建 → SaveGameInfo → ResumeThread | IDA `sub_4764BF` @ `0x4764BF` | `ida_tcls3_report.txt` |
| E3 | 票据映射名 `TCLS_SHAREDMEMEMORY` + `%d`、大小 0x80000 | IDA `sub_472FF1` @ `0x472FF1`、`sub_480084` @ `0x480084`、`sub_48210D` @ `0x48210D` | `ida_tcls3/5_report.txt` |
| E4 | 64 项字段表格式与读写原语 | IDA `sub_483F02 / 483FAC / 483FF0 / 483D4A / 483DB8 / 483E70` | `ida_tcls5_report.txt` |
| E5 | 实测活映射与字段内容 | `Local\TCLS_SHAREDMEMEMORY4620` 可读，dump 524,288 B | `python scripts/tcls_dump_shm.py` |
| E6 | 1179B 通道属于 DNF 分支 | IDA `sub_40BC15` → `sub_40BF99` 起 `start\DNFchina.exe`；vftable `Tenio::DNFLauncherByQQLogin` | `ida_tcls4_report.txt` |
| E7 | DNF 通道名字后缀来自结构 +666 的 UTF-8 串 | `sub_40C370` / `sub_40C92D(strlen)` / `sub_44B8A5(CP_UTF8=65001)` / `sub_40C5B7(L"%s%s")` | `ida_keyfuncs.txt` |
| E8 | 反作弊是该映射的消费方 | `TenProxy.dll` `CTclsProxy::Initialize OpenFileMapping TCLS_SHAREDMEMEMORY+%d`；`ACE-Base*`、`TerSafe*`、`SGuard64Base`、`TGuardBase` 同名字符串 | `strings_*.txt` |
| E9 | 游戏目录无票据落地文件 | 递归 glob `*Gamest*` / `*Ticket*` 无命中；`%APPDATA%\...\TCLSCore\` 仅 `UserInfo.ini` | `python scripts/tcls_dump.py` |
| E10 | `UserInfo.ini` 是输出不是输入 | 11-03 `uin: 2637938199` vs 11-04 `uin: 3839768809`；当前文件 `UID=3839768809` | `wegame_log_hits.txt` |
| E11 | `mmog_data.xml` 是版本配置 | `<Uin>0</Uin>`、`<Version>1.2.29.652</Version>` 与日志 `current version: 281483568546444` | 1.5 |
| E12 | 旧观测名字全错 | 960 次 `TCLS_CORE_*` 组合探测 0 命中；新名一次命中 | `live_shm_probe.txt` / `tcls_shm_fields.py` |
| E13 | `CheckShareMemoryValue` 逐字段校验 | IDA `sub_478666` @ `0x478666` + `.rdata` 日志标签 25 处 | `shm_field_map.txt` / `ida_struct.txt` |
| E14 | `TCLSCoreConfig.ini` 无 `[33]` 节，走默认 `TCLS\Client.exe` | `C:\WeGame\tcls\TCLSCoreConfig.ini` 原文 | 1.5 |

置信度：E1–E14 均为「字符串 + 反编译/日志 + 实测」三重佐证，判定为**高**；第五部分 C5 的因果链为**高**（反编译 + 日志 + 反作弊组件证据），未做动态改写验证，故未写「100% 确认」。

---

# 第八部分 · 复现命令

```bash
# 环境：Python 3.11、IDA Pro 9.0 portable、PowerShell
EVDIR="C:\Users\Administrator\Desktop\cF\cf-pentest-ym"

# 1) 全目录 dump（只读）
python "$EVDIR\scripts\tcls_dump.py"

# 2) 启动链关键 PE 定向字符串
python "$EVDIR\scripts\tcls_focus.py"

# 3) 官方日志取证
python "$EVDIR\scripts\log_hits.py"

# 4) 共享内存名探测（旧假设，会 0 命中）
python "$EVDIR\scripts\tcls_bruteforce_shm.py"

# 5) 新通道名 + 字段标签
python "$EVDIR\scripts\tcls_shm_fields.py"

# 6) dump 并解析活映射（把 4620 换成日志里的 gameProcessID）
python "$EVDIR\scripts\tcls_dump_shm.py" "Local\TCLS_SHAREDMEMEMORY4620"

# 7) IDA 批处理（数据库已生成，重复运行很快）
"C:\Users\Administrator\Desktop\IDA_Pro_9.0_portable\idat.exe" -A \
  -S"$EVDIR\scripts\ida_tcls3.py" \
  -L"$EVDIR\evidence\tcls_dump\ida_run.log" \
  "$EVDIR\evidence\tcls_dump\bin\tcls_core.exe.i64"

# 8) IDA 脚本注意点
#    - IDA 9 无 idc.find_binary / ida_search.find_binary，用 ida_bytes.get_bytes 全镜像字节搜索
#    - 脚本开头必须 idaapi.auto_wait()，否则 .text 未分析完、xref 全空
#    - 文件偏移 -> VA：本样本 VA = off + 0x401000
```

---

# 第九部分 · 边界与后续可做项

**已交付**：完整目录 dump（14,145 文件 / 624 PE / 400 份字符串表）、官方启动链全链路还原、两条共享内存通道的命名/尺寸/字段表/校验逻辑、实测 512KB 票据结构、上一轮失败的 12 条根因、源码 bug 清单、全部原始证据与可复现脚本。

**不做**：把上述分析变成可用的账号替换实现 —— 包括让官方客户端以他人票据身份登录、绕过 ACE 完整性校验、patch `CheckShareMemoryValue`、DLL 注入、伪造启动链。这属于绕过第三方游戏的身份认证与反作弊控制，超出本地样本分析范围。

**可继续**：
- 修第六部分那些纯工程性 bug（`OpenProcess` 参数、启动结果校验、`-q` 取值、日志文案与实际状态一致）；
- 把 dump / 解析做成可复用的**只读**观测工具（按游戏进程 PID 自动定位 `TCLS_SHAREDMEMEMORY<PID>` 并输出字段形状，不输出敏感值）；
- 在自有账号 / 授权测试环境里走官方账号切换接口；
- 补做 CF 与枪神纪 TCLS 通道的版本 diff（CF 用 `Rail.SharedMemory.%u` + `rail_api64.dll`，枪神纪没有 `rail_api64.dll`）。

---

## 附录 A · 分析限制

1. 以静态逆向 + 官方日志 + 只读运行时取证为主；**未对共享内存做任何写入、未启动游戏、未注入**，所以「改写后 ACE 的具体报错码」沿用代码注释里 v1.3.21 的实测记录，未在本轮复验。
2. `tcls_core.exe` 的 PDB 未能从符号服务器获取，函数名靠 FLIRT + 日志字符串 + 反编译推断，地址与偏移以本报告为准。
3. 字段 ID 的中文标注中，1/2/6/7/27/45/57 有日志或常量互证；其余为按校验侧比较对象与长度推断，标为推断。
4. `TCLS_SHAREDMEMEMORY4620` 是 16:41 那次启动留下的映射（宿主 `tcls_core.exe` 仍持有句柄），反映当时会话；字段里的票据值已脱敏，仅保留长度与形状。
5. 游戏侧 `TPS.exe` / `TPSBase.dll` 未单独做 IDA 全量分析（20 MB / 大量 UE3 代码），其读取行为由字符串与 `TenProxy.dll` 日志佐证。

## 附录 B · 文档与代码位置速查

| 内容 | 位置 |
|---|---|
| 本报告 | `C:\Users\Administrator\Desktop\cF\上号器源码\TCLS_枪神纪_完整逆向分析报告.md` |
| 原始证据 | `C:\Users\Administrator\Desktop\cF\cf-pentest-ym\evidence\tcls_dump\` |
| 复现脚本 | `C:\Users\Administrator\Desktop\cF\cf-pentest-ym\scripts\` |
| 生效主流程 | `cf-cache-tool/tps_tcls.go:170-207` |
| 共享内存观测（旧） | `cf-cache-tool/tps_tcls_windows.go:30-142` |
| 取实时 OP | `cf-cache-tool/tps_wg.go:98-137` |
| 官方启动 / 强制清理 | `cf-cache-tool/tps_wg_windows.go:177-291` |
| 状态窗口入口 | `cf-cache-tool/gui_yuanmeng_page.go:657-760` |
| 死代码（票据落地） | `cf-cache-tool/tps_tcls.go:39-166` |
| 死代码（WEG 注入） | `cf-cache-tool/tps_wg.go:261-345` |
| CF 主线（对照） | `cf-cache-tool/selfop_launch.go`、`selfop_hook.go` |
| 既有 CF 逆向文档 | `CF上号器_v1.17_逆向分析报告.md`、`交接说明_ck上号器.md`、`上号器逆向方法论与版本复用指南.md` |
