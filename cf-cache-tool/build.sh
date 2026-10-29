#!/usr/bin/env bash
# CG 打包脚本
#
#   ./build.sh              自动把版本号 +1（1.0.1 -> 1.0.2）并编译
#   ./build.sh 1.2.0        指定版本号编译
#   ./build.sh --publish    自动 +1 并上传到云端更新目录
#   ./build.sh 1.2.0 --publish
#   ./build.sh --no-loader  跳过上号器资源同步（离线打包时用）
#   ./build_launcher.sh     只重编译/发布 Rust 启动器，主程序不变
#
# 凭据：.build.env 里的 DAV_USER / DAV_PASS 只用于构建机上传和创建只读分享，
#       不进入任何客户端二进制。客户端只拿只读分享 token / 口令。
#
# ============================ 外观复制顺序（不可调换） ============================
# 参考文件：exe/zuhaoclient.exe（GG租号，EV 签名有效），工具：stamp_exe.go
#
#   启动器：链接 → -icon -version → UPX → -sig -time
#   主程序：go build → 四项一次写完 → 算 sha256/md5/大小 → 签名清单 → gzip
#
# 1) 图标/版本信息走 UpdateResource，会重写 .rsrc，必须在 UPX 之前；先压缩再改会
#    破坏 UPX 的节布局。
# 2) 签名是往文件尾追加 PE 证书表并改写 DataDirectory[4]，之后任何写文件都会让
#    内容与签名彻底不符，只能放最后。
# 3) 主程序的签名会改变文件内容，必须在算 sha256 之前写，否则清单校验必失败。
# 4) 文件时间要排在 writeChecksum 之后，写文件会把 mtime 刷成当前时间。
# =================================================================================
#
# 产物：../dist/CG.exe（唯一对外分发的启动器，Rust + UPX）
#       ../dist/cloud/* 仅为本地上传暂存，不对外分发
# 云端：CG/app/* 为新版主程序清单；CG/update/* 保留用于旧客户端迁移。

set -euo pipefail
cd "$(dirname "$0")"

# ---------- 凭据 ----------
if [ -f ".build.env" ]; then
    set -a
    # shellcheck disable=SC1091
    . ./.build.env
    set +a
fi
if [ -z "${DAV_USER:-}" ] || [ -z "${DAV_PASS:-}" ]; then
    echo "错误: 缺少 DAV_USER / DAV_PASS"
    echo "      写进 .build.env（DAV_USER=... / DAV_PASS=...）或导出为环境变量"
    exit 1
fi

VER_FILE="version.txt"
PUBLISH=0
SYNC_LOADER=1
NEW=""

for arg in "$@"; do
    case "$arg" in
        --publish) PUBLISH=1 ;;
        --no-loader) SYNC_LOADER=0 ;;
        -h|--help) sed -n '2,30p' "$0"; exit 0 ;;
        *) NEW="$arg" ;;
    esac
done

if [ ! -f "$VER_FILE" ]; then
    echo "1.0.0" > "$VER_FILE"
fi

CUR="$(tr -d ' \t\r\n' < "$VER_FILE")"

if [ -z "$NEW" ]; then
    MAJOR="${CUR%%.*}"
    REST="${CUR#*.}"
    if [ "$REST" = "$CUR" ]; then
        MINOR=0
        PATCH=0
    else
        MINOR="${REST%%.*}"
        if [ "$MINOR" = "$REST" ]; then
            PATCH=0
        else
            PATCH="${REST#*.}"
        fi
    fi
    MAJOR=${MAJOR:-1}
    MINOR=${MINOR:-0}
    PATCH=${PATCH:-0}
    NEW="${MAJOR}.${MINOR}.$((PATCH + 1))"
fi

echo "$NEW" > "$VER_FILE"
echo "==> 版本: $CUR -> $NEW"

# ---------- 只读分享 ----------
# 构建机用 DAV 账号维护 CG 根目录的只读公开分享；客户端只拿到分享 token / 口令。
SHARE_STATE=".cg-share.env"
echo "==> 同步 Nextcloud 只读分享"
CGO_ENABLED=0 go run share_provision.go
if [ -f "$SHARE_STATE" ]; then
    set -a
    # shellcheck disable=SC1091
    . "./$SHARE_STATE"
    set +a
fi
if [ -z "${CG_SHARE_TOKEN:-}" ] || [ -z "${CG_SHARE_PASS:-}" ]; then
    echo "错误: 只读分享凭据未生成（检查 $SHARE_STATE）"
    exit 1
fi

# ---------- 生成构建期常量 ----------
# 只把只读分享 token / 口令和版本号写进生成文件；DAV 主账号密码不进入二进制。
TOKEN_ENC="$(CGO_ENABLED=0 go run obfstr.go "$CG_SHARE_TOKEN")"
SHARE_ENC="$(CGO_ENABLED=0 go run obfstr.go "$CG_SHARE_PASS")"

SECRET_FILE="zz_build_secret.go"
LAUNCHER_RSRC="launcher/rsrc.syso"
cat > "$SECRET_FILE" <<EOF
//go:build cg_secret

package main

func buildToolVersion() string    { return "${NEW}" }
func buildCloudShareToken() string { return deob("${TOKEN_ENC}") }
func buildCloudSharePass() string  { return deob("${SHARE_ENC}") }
EOF
cleanup() { rm -f "$SECRET_FILE" "$LAUNCHER_RSRC"; }
trap cleanup EXIT

LDFLAGS="-s -w"
DIST="../dist"
CLOUD="${DIST}/cloud"
mkdir -p "$DIST" "$CLOUD/app"
# 旧版本遗留的独立主程序/CLI/测试启动器不再分发。
rm -f "${DIST}/CG-app.exe" "${DIST}/CG-cli.exe" "${DIST}/CG-test.exe" "${DIST}/CG.exe"
rm -rf "${DIST}/app"

# ---------- 编译器 ----------
# 不用 garble：它会把依赖包（lxn/walk declarative）的未导出字段名一起混淆，
# 而 declarative 靠反射按字段名赋值，混淆后启动即 panic（实测 v0.18.0）。
# 凭据改用 obfstr.go 的 base64+XOR 混淆，strings 同样搜不到明文。
BUILD_TOOL="go"

find_python() {
    for c in python3 python; do
        if command -v "$c" >/dev/null 2>&1 && "$c" -V >/dev/null 2>&1; then
            command -v "$c"
            return 0
        fi
    done
    return 1
}
# Python 只用于可选的字符串自检；资源生成已改成 Go，不再强制依赖
PYTHON_BIN="${PYTHON_BIN:-$(find_python || true)}"

# macOS 只有 shasum，Git Bash/Linux 只有 sha256sum
sha256() {
    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$1" | cut -d' ' -f1
    else
        sha256sum "$1" | cut -d' ' -f1
    fi
}

md5_of() {
    if command -v md5sum >/dev/null 2>&1; then
        md5sum "$1" | cut -d' ' -f1
    else
        md5 -q "$1"
    fi
}

sha256_stdin() {
    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 | cut -d' ' -f1
    else
        sha256sum | cut -d' ' -f1
    fi
}

find_upx() {
    if [ -n "${UPX_BIN:-}" ] && [ -x "$UPX_BIN" ]; then
        printf '%s' "$UPX_BIN"
        return 0
    fi
    if command -v upx >/dev/null 2>&1; then
        command -v upx
        return 0
    fi
    if command -v cygpath >/dev/null 2>&1 && [ -n "${LOCALAPPDATA:-}" ]; then
        find "$(cygpath "$LOCALAPPDATA")/Microsoft/WinGet/Packages" -name upx.exe -type f 2>/dev/null | head -1
        return 0
    fi
    return 0
}

echo "==> 生成图标/版本资源"
CGO_ENABLED=0 go run build_resources.go "$NEW"

echo "==> 编译主程序（云端下发）"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $BUILD_TOOL build -trimpath -tags cg_secret \
    -ldflags "-H=windowsgui ${LDFLAGS}" -o "${CLOUD}/zuhaoclient.exe" .

# 参考 exe：图标 / 版本信息 / 数字签名 / 文件时间都从这里复制
DONOR="exe/zuhaoclient.exe"

# 主程序不经过 UPX，四项一次写完；必须在算 sha256 之前，签名会影响文件内容
if [ -f "$DONOR" ]; then
    echo "==> 复制图标/版本信息/签名/时间（主程序）"
    CGO_ENABLED=0 go run stamp_exe.go -from "$DONOR" -to "${CLOUD}/zuhaoclient.exe"
else
    echo "==> 跳过外观复制：缺少 $DONOR"
fi

echo "==> 编译 CG租号.exe (Rust + UPX)"
(
    cd launcher
    export CG_SHARE_TOKEN CG_SHARE_PASS CG_LAUNCHER_VERSION="$NEW"
    export VCVARS64="${VCVARS64:-C:\\Program Files (x86)\\Microsoft Visual Studio\\2022\\BuildTools\\VC\\Auxiliary\\Build\\vcvars64.bat}"
    cmd.exe //c build.cmd
)
cp launcher/target/release/cg-launcher.exe "${DIST}/CG租号.exe"

# 图标 / 版本信息必须在 UPX 之前写入：UpdateResource 会重写 .rsrc，
# 先压缩再改会破坏 UPX 的节布局。签名只能放在最后（追加到文件尾）。
if [ -f "$DONOR" ]; then
    echo "==> 复制图标/版本信息（参考 $(basename "$DONOR")）"
    CGO_ENABLED=0 go run stamp_exe.go -from "$DONOR" -to "${DIST}/CG租号.exe" -icon -version
else
    echo "==> 跳过图标/版本复制：缺少 $DONOR"
fi

UPX="$(find_upx)"
if [ -n "$UPX" ]; then
    "$UPX" --best --lzma -q "${DIST}/CG租号.exe"
    "$UPX" -t -q "${DIST}/CG租号.exe"
fi

if [ -f "$DONOR" ]; then
    echo "==> 复制数字签名/修改时间"
    CGO_ENABLED=0 go run stamp_exe.go -from "$DONOR" -to "${DIST}/CG租号.exe" -sig -time
fi
if ! grep -aq "requireAdministrator" "${DIST}/CG租号.exe"; then
    echo "错误: 启动器没有嵌入管理员清单"
    exit 1
fi

# 分发包指纹只看真正影响启动器行为的输入，不受 PE 时间戳影响。
TOOLCHAIN="$(rustc -Vv 2>/dev/null | tr '\n' ' ' || true)"
SECRET_FP="$(printf '%s\n%s' "$CG_SHARE_TOKEN" "$CG_SHARE_PASS" | sha256_stdin)"
LAUNCHER_FP="$(
    {
        printf '%s\n%s\n%s\n' "$TOOLCHAIN" "$SECRET_FP" "$NEW"
        for f in launcher/Cargo.toml launcher/Cargo.lock launcher/build.rs launcher/build.cmd \
                 launcher/src/main.rs launcher/src/crypto.rs launcher/src/cng.rs launcher/src/net.rs \
                 launcher/src/gzip.rs build_resources.go stamp_exe.go exe/zuhaoclient.exe assets/zuhaoclient.ico \
                 assets/icon_256.png assets/icon_64.png assets/icon_48.png; do
            sha256 "$f"
        done
    } | sha256_stdin
)"

echo "==> 生成云下发主程序包"
APP_SHA="$(sha256 "${CLOUD}/zuhaoclient.exe")"
APP_MD5="$(md5_of "${CLOUD}/zuhaoclient.exe")"
APP_SIZE="$(wc -c < "${CLOUD}/zuhaoclient.exe" | tr -d ' ')"
APP_SIG="$(CGO_ENABLED=0 go run sign_update.go "$NEW" "$APP_SHA" "$APP_SIZE" 2>/dev/null || true)"
if [ -z "$APP_SIG" ]; then
    echo "错误: 主程序清单签名失败（检查 ${HOME}/.cg-keys/update.key）"
    exit 1
fi
APP_ECDSA="$(CGO_ENABLED=0 go run sign_update_ecdsa.go sign "$NEW" "$APP_SHA" "$APP_SIZE" 2>/dev/null || true)"
if [ -z "$APP_ECDSA" ]; then
    echo "错误: 主程序 ECDSA 签名失败（检查 ${HOME}/.cg-keys/launcher_ecdsa.key）"
    exit 1
fi
gzip -9 -n -c "${CLOUD}/zuhaoclient.exe" > "${CLOUD}/app/CG-app.exe.gz"
cat > "${CLOUD}/app/version.json" <<EOF
{
  "version": "${NEW}",
  "file": "CG-app.exe.gz",
  "sha256": "${APP_SHA}",
  "md5": "${APP_MD5}",
  "size": ${APP_SIZE},
  "notes": "build $(date '+%Y-%m-%d %H:%M')",
  "sig": "${APP_SIG}",
  "sig_ecdsa": "${APP_ECDSA}"
}
EOF

echo "==> 校验内嵌资源"
check_str() {
    if grep -aq -- "$2" "$1"; then
        echo "    $3: OK"
    else
        echo "    $3: 缺失"
    fi
}
check_str "${CLOUD}/zuhaoclient.exe" "$NEW" "版本号"
check_str "${CLOUD}/zuhaoclient.exe" "游戏租号客户端" "品牌"
check_str "${CLOUD}/zuhaoclient.exe" "圆梦上号器" "上号器类型"
echo "    主程序大小: $(awk "BEGIN{printf \"%.1f\", $(wc -c < "${CLOUD}/zuhaoclient.exe")/1048576}") MB"
echo "    启动器大小: $(($(wc -c < "${DIST}/CG租号.exe")/1024)) KB"
echo "    分发包指纹: ${LAUNCHER_FP}"
echo "    云下发压缩包: $(($(wc -c < "${CLOUD}/app/CG-app.exe.gz")/1024)) KB"

BASE="https://cloud.51i.cc/remote.php/dav/files/leonai"
AUTH="${DAV_USER}:${DAV_PASS}"
LOADER_ID="1"

# 上号器本体与 host.dat 不打进 exe，构建时同步到云端（内容没变就跳过）
sync_loader() {
    # 上号器资源目录：优先用环境变量 LOADER_ASSET_DIR，其次 ../资源/上号器/yuanmeng
    local dir="${LOADER_ASSET_DIR:-../资源/上号器/yuanmeng}"
    local exe
    exe="$(ls "$dir"/*.exe 2>/dev/null | head -1)"
    if [ -z "$exe" ] || [ ! -f "$dir/host.dat" ]; then
        echo "==> [loader] 跳过：$dir 下没有 exe 或 host.dat"
        echo "            可用 LOADER_ASSET_DIR=路径 ./build.sh 指定其它目录"
        return 0
    fi

    local exe_name exe_sha dat_sha exe_size remote
    exe_name="$(basename "$exe")"
    exe_sha="$(sha256 "$exe")"
    dat_sha="$(sha256 "$dir/host.dat")"
    exe_size="$(wc -c < "$exe" | tr -d ' ')"

    remote="$(curl -s -u "$AUTH" "$BASE/CG/loaders/$LOADER_ID/meta.json" || true)"
    if echo "$remote" | grep -q "$exe_sha" && echo "$remote" | grep -q "$dat_sha"; then
        echo "==> [loader] 云端已是最新（${exe_name}），跳过上传"
        return 0
    fi

    echo "==> [loader] 上传 ${exe_name} ($((exe_size/1024/1024)) MB) 与 host.dat"
    curl -sS -u "$AUTH" -X MKCOL "$BASE/CG/" >/dev/null || true
    curl -sS -u "$AUTH" -X MKCOL "$BASE/CG/loaders/" >/dev/null || true
    curl -sS -u "$AUTH" -X MKCOL "$BASE/CG/loaders/$LOADER_ID/" >/dev/null || true
    curl -sS -u "$AUTH" -T "$exe" "$BASE/CG/loaders/$LOADER_ID/${exe_name}" >/dev/null
    curl -sS -u "$AUTH" -T "$dir/host.dat" "$BASE/CG/loaders/$LOADER_ID/host.dat" >/dev/null

    local meta
    meta="$(mktemp)"
    cat > "$meta" <<EOF
{
  "id": "$LOADER_ID",
  "name": "圆梦上号器",
  "exe": "${exe_name}",
  "exe_sha256": "$exe_sha",
  "exe_size": $exe_size,
  "dat_sha256": "$dat_sha",
  "updated": "$(date '+%Y-%m-%d %H:%M:%S')"
}
EOF
    curl -sS -u "$AUTH" -T "$meta" "$BASE/CG/loaders/$LOADER_ID/meta.json" >/dev/null
    rm -f "$meta"
    echo "==> [loader] 已发布 ${exe_name}"
}

if [ "$SYNC_LOADER" = "1" ]; then
    sync_loader
fi

if [ "$PUBLISH" = "1" ]; then
    SHA="$(sha256 "${DIST}/CG租号.exe")"
    SIZE="$(wc -c < "${DIST}/CG租号.exe" | tr -d ' ')"

    REMOTE_JSON="$(mktemp)"
    REMOTE_SHA=""
    REMOTE_FP=""
    REMOTE_SIZE=""
    if curl -fsS -u "$AUTH" -o "$REMOTE_JSON" "${BASE}/CG/update/version.json" >/dev/null 2>&1; then
        if [ -n "${PYTHON_BIN:-}" ] && "$PYTHON_BIN" -V >/dev/null 2>&1; then
            REMOTE_SHA="$("$PYTHON_BIN" -c 'import json,sys; print(json.load(open(sys.argv[1],encoding="utf-8")).get("sha256",""))' "$REMOTE_JSON" 2>/dev/null || true)"
            REMOTE_FP="$("$PYTHON_BIN" -c 'import json,sys; print(json.load(open(sys.argv[1],encoding="utf-8")).get("launcher_fp",""))' "$REMOTE_JSON" 2>/dev/null || true)"
            REMOTE_SIZE="$("$PYTHON_BIN" -c 'import json,sys; print(json.load(open(sys.argv[1],encoding="utf-8")).get("size",0))' "$REMOTE_JSON" 2>/dev/null || true)"
        else
            REMOTE_SHA="$(sed -n 's/.*"sha256"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$REMOTE_JSON" | head -1)"
            REMOTE_FP="$(sed -n 's/.*"launcher_fp"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$REMOTE_JSON" | head -1)"
            REMOTE_SIZE="$(sed -n 's/.*"size"[[:space:]]*:[[:space:]]*\([0-9][0-9]*\).*/\1/p' "$REMOTE_JSON" | head -1)"
        fi
    fi
    rm -f "$REMOTE_JSON"

    echo "==> 上传 CG/app/CG-app.exe.gz ($(($(wc -c < "${CLOUD}/app/CG-app.exe.gz")/1024)) KB)"
    curl -sS -u "$AUTH" -X MKCOL "${BASE}/CG/app/" >/dev/null || true
    curl -sS -u "$AUTH" -T "${CLOUD}/app/CG-app.exe.gz" "${BASE}/CG/app/CG-app.exe.gz" >/dev/null
    echo "==> 上传 CG/app/version.json（已签名，MD5 ${APP_MD5:0:12}...）"
    curl -sS -u "$AUTH" -T "${CLOUD}/app/version.json" "${BASE}/CG/app/version.json" >/dev/null

    LAUNCH_UPLOADED=1
    if { [ -n "$REMOTE_FP" ] && [ "$REMOTE_FP" = "$LAUNCHER_FP" ]; } || \
       { [ -z "$REMOTE_FP" ] && [ -n "$REMOTE_SHA" ] && [ "$REMOTE_SHA" = "$SHA" ]; }; then
        echo "==> 分发包未变化，跳过 CG/update/CG.exe 上传（fp ${LAUNCHER_FP:0:12}...）"
        LAUNCH_UPLOADED=0
        # 二进制不传，但仍用云端现有启动器的 hash/大小更新版本清单，旧客户端才能继续迁移。
        SIG=""
        SIG_ECDSA=""
        if [ -f "${HOME}/.cg-keys/update.key" ]; then
            SIG="$(CGO_ENABLED=0 go run sign_update.go "$NEW" "$REMOTE_SHA" "$REMOTE_SIZE" 2>/dev/null || true)"
        fi
        SIG_ECDSA="$(CGO_ENABLED=0 go run sign_update_ecdsa.go sign "$NEW" "$REMOTE_SHA" "$REMOTE_SIZE" 2>/dev/null || true)"
        if [ -z "$SIG" ] || [ -z "$SIG_ECDSA" ]; then
            echo "错误: 分发包清单更新签名失败（检查 ${HOME}/.cg-keys/）"
            exit 1
        fi
        JSON="$(mktemp)"
        cat > "$JSON" <<EOF
{
  "version": "${NEW}",
  "file": "CG.exe",
  "sha256": "${REMOTE_SHA}",
  "size": ${REMOTE_SIZE},
  "launcher_fp": "${LAUNCHER_FP}",
  "notes": "build $(date '+%Y-%m-%d %H:%M')",
  "sig": "${SIG}",
  "sig_ecdsa": "${SIG_ECDSA}"
}
EOF
        echo "==> 仅更新 CG/update/version.json"
        curl -sS -u "$AUTH" -X MKCOL "${BASE}/CG/update/" >/dev/null || true
        curl -sS -u "$AUTH" -T "$JSON" "${BASE}/CG/update/version.json" >/dev/null
        rm -f "$JSON"
    else
        SIG=""
        SIG_ECDSA=""
        if [ -f "${HOME}/.cg-keys/update.key" ]; then
            SIG="$(CGO_ENABLED=0 go run sign_update.go "$NEW" "$SHA" "$SIZE" 2>/dev/null || true)"
        fi
        SIG_ECDSA="$(CGO_ENABLED=0 go run sign_update_ecdsa.go sign "$NEW" "$SHA" "$SIZE" 2>/dev/null || true)"
        if [ -z "$SIG" ] || [ -z "$SIG_ECDSA" ]; then
            echo "错误: 没能生成更新签名（检查 ${HOME}/.cg-keys/）"
            echo "      缺私钥时发布会生成无签名清单，所有客户端都会拒绝更新。"
            exit 1
        fi

        JSON="$(mktemp)"
        cat > "$JSON" <<EOF
{
  "version": "${NEW}",
  "file": "CG.exe",
  "sha256": "${SHA}",
  "size": ${SIZE},
  "launcher_fp": "${LAUNCHER_FP}",
  "notes": "build $(date '+%Y-%m-%d %H:%M')",
  "sig": "${SIG}",
  "sig_ecdsa": "${SIG_ECDSA}"
}
EOF

        # 旧版客户端仍从 CG/update 自更新：这里放启动器，完成到新架构的迁移。
        echo "==> 上传 CG/update/CG.exe（启动器）"
        curl -sS -u "$AUTH" -X MKCOL "${BASE}/CG/update/" >/dev/null || true
        curl -sS -u "$AUTH" -T "${DIST}/CG租号.exe" "${BASE}/CG/update/CG.exe" >/dev/null
        echo "==> 上传 CG/update/version.json（已签名）"
        curl -sS -u "$AUTH" -T "$JSON" "${BASE}/CG/update/version.json" >/dev/null
        rm -f "$JSON"
    fi
    rm -rf "$CLOUD"
    echo "==> 已发布 v${NEW}（主程序 ${APP_SHA:0:12}... / 分发包更新=${LAUNCH_UPLOADED}）"
fi

echo "==> 完成: v${NEW}"
echo "    唯一分发文件: ${DIST}/CG租号.exe"
