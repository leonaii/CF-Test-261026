#!/usr/bin/env bash
# CG 打包脚本
#
#   ./build.sh              自动把版本号 +1（1.0.1 -> 1.0.2）并编译
#   ./build.sh 1.2.0        指定版本号编译
#   ./build.sh --publish    自动 +1 并上传到云端更新目录
#   ./build.sh 1.2.0 --publish
#   ./build.sh --no-loader  跳过上号器资源同步（离线打包时用）
#
# 凭据：.build.env 里的 DAV_USER / DAV_PASS 只用于构建机上传和创建只读分享，
#       不进入任何客户端二进制。客户端只拿只读分享 token / 口令。
#
# 产物：../dist/CG.exe（唯一对外分发的启动器）
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
        -h|--help) sed -n '2,16p' "$0"; exit 0 ;;
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
LAUNCHER_SECRET="launcher/zz_build_secret.go"
LAUNCHER_RSRC="launcher/rsrc.syso"
cat > "$SECRET_FILE" <<EOF
//go:build cg_secret

package main

func buildToolVersion() string    { return "${NEW}" }
func buildCloudShareToken() string { return deob("${TOKEN_ENC}") }
func buildCloudSharePass() string  { return deob("${SHARE_ENC}") }
EOF
cat > "$LAUNCHER_SECRET" <<EOF
//go:build cg_secret

package main

func buildShareToken() string { return deob("${TOKEN_ENC}") }
func buildSharePass() string  { return deob("${SHARE_ENC}") }
EOF
cp rsrc.syso "$LAUNCHER_RSRC"
cleanup() { rm -f "$SECRET_FILE" "$LAUNCHER_SECRET" "$LAUNCHER_RSRC"; }
trap cleanup EXIT

LDFLAGS="-s -w"
DIST="../dist"
CLOUD="${DIST}/cloud"
mkdir -p "$DIST" "$CLOUD/app"
# 旧版本遗留的独立主程序/CLI/测试启动器不再分发。
rm -f "${DIST}/CG-app.exe" "${DIST}/CG-cli.exe" "${DIST}/CG-test.exe"
rm -rf "${DIST}/app"

# ---------- 编译器 ----------
# 不用 garble：它会把依赖包（lxn/walk declarative）的未导出字段名一起混淆，
# 而 declarative 靠反射按字段名赋值，混淆后启动即 panic（实测 v0.18.0）。
# 凭据改用 obfstr.go 的 base64+XOR 混淆，strings 同样搜不到明文。
BUILD_TOOL="go"

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

echo "==> 编译主程序（云端下发）"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $BUILD_TOOL build -trimpath -tags cg_secret \
    -ldflags "-H=windowsgui ${LDFLAGS}" -o "${CLOUD}/CG-app.exe" .

echo "==> 编译 CG.exe (唯一分发启动器)"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $BUILD_TOOL build -trimpath -tags cg_secret \
    -ldflags "-H=windowsgui ${LDFLAGS}" -o "${DIST}/CG.exe" ./launcher

echo "==> 生成云下发主程序包"
APP_SHA="$(sha256 "${CLOUD}/CG-app.exe")"
APP_MD5="$(md5_of "${CLOUD}/CG-app.exe")"
APP_SIZE="$(wc -c < "${CLOUD}/CG-app.exe" | tr -d ' ')"
APP_SIG="$(CGO_ENABLED=0 go run sign_update.go "$NEW" "$APP_SHA" "$APP_SIZE" 2>/dev/null || true)"
if [ -z "$APP_SIG" ]; then
    echo "错误: 主程序清单签名失败（检查 ${HOME}/.cg-keys/update.key）"
    exit 1
fi
gzip -9 -n -c "${CLOUD}/CG-app.exe" > "${CLOUD}/app/CG-app.exe.gz"
cat > "${CLOUD}/app/version.json" <<EOF
{
  "version": "${NEW}",
  "file": "CG-app.exe.gz",
  "sha256": "${APP_SHA}",
  "md5": "${APP_MD5}",
  "size": ${APP_SIZE},
  "notes": "build $(date '+%Y-%m-%d %H:%M')",
  "sig": "${APP_SIG}"
}
EOF

echo "==> 校验内嵌资源"
PYTHON_BIN="${PYTHON_BIN:-}"
if [ -z "$PYTHON_BIN" ]; then
    for c in python3 python; do
        if command -v "$c" >/dev/null 2>&1 && "$c" -V >/dev/null 2>&1; then
            PYTHON_BIN="$c"
            break
        fi
    done
fi
if [ -n "$PYTHON_BIN" ] && "$PYTHON_BIN" -V >/dev/null 2>&1; then
"$PYTHON_BIN" - "$CLOUD/CG-app.exe" "$NEW" <<'PY'
import sys
path, ver = sys.argv[1], sys.argv[2]
data = open(path, "rb").read()
checks = {
    "版本号": ver.encode(),
    "品牌": "CG租号".encode(),
    "上号器类型": "圆梦上号器".encode(),
}
for name, needle in checks.items():
    print(f"    {name}: {'OK' if needle in data else '缺失'}")
print(f"    主程序大小: {len(data)/1024/1024:.1f} MB")
PY
else
    echo "    跳过（本机 python 不可用）"
fi
echo "    启动器大小: $(($(wc -c < "${DIST}/CG.exe")/1024)) KB"
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
    SHA="$(sha256 "${DIST}/CG.exe")"
    SIZE="$(wc -c < "${DIST}/CG.exe" | tr -d ' ')"

    REMOTE_JSON="$(mktemp)"
    REMOTE_SHA=""
    if curl -fsS -u "$AUTH" -o "$REMOTE_JSON" "${BASE}/CG/update/version.json" >/dev/null 2>&1; then
        if [ -n "${PYTHON_BIN:-}" ] && "$PYTHON_BIN" -V >/dev/null 2>&1; then
            REMOTE_SHA="$("$PYTHON_BIN" -c 'import json,sys; print(json.load(open(sys.argv[1],encoding="utf-8")).get("sha256",""))' "$REMOTE_JSON" 2>/dev/null || true)"
        else
            REMOTE_SHA="$(sed -n 's/.*"sha256"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$REMOTE_JSON" | head -1)"
        fi
    fi
    rm -f "$REMOTE_JSON"

    echo "==> 上传 CG/app/CG-app.exe.gz ($(($(wc -c < "${CLOUD}/app/CG-app.exe.gz")/1024)) KB)"
    curl -sS -u "$AUTH" -X MKCOL "${BASE}/CG/app/" >/dev/null || true
    curl -sS -u "$AUTH" -T "${CLOUD}/app/CG-app.exe.gz" "${BASE}/CG/app/CG-app.exe.gz" >/dev/null
    echo "==> 上传 CG/app/version.json（已签名，MD5 ${APP_MD5:0:12}...）"
    curl -sS -u "$AUTH" -T "${CLOUD}/app/version.json" "${BASE}/CG/app/version.json" >/dev/null

    LAUNCH_UPLOADED=1
    if [ -n "$REMOTE_SHA" ] && [ "$REMOTE_SHA" = "$SHA" ]; then
        echo "==> 分发包未变化，跳过 CG/update 上传"
        LAUNCH_UPLOADED=0
    else
        SIG=""
        if [ -f "${HOME}/.cg-keys/update.key" ]; then
            SIG="$(CGO_ENABLED=0 go run sign_update.go "$NEW" "$SHA" "$SIZE" 2>/dev/null || true)"
        fi
        if [ -z "$SIG" ]; then
            echo "错误: 没能生成更新签名（检查 ${HOME}/.cg-keys/update.key）"
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
  "notes": "build $(date '+%Y-%m-%d %H:%M')",
  "sig": "${SIG}"
}
EOF

        # 旧版客户端仍从 CG/update 自更新：这里放启动器，完成到新架构的迁移。
        echo "==> 上传 CG/update/CG.exe（启动器）"
        curl -sS -u "$AUTH" -X MKCOL "${BASE}/CG/update/" >/dev/null || true
        curl -sS -u "$AUTH" -T "${DIST}/CG.exe" "${BASE}/CG/update/CG.exe" >/dev/null
        echo "==> 上传 CG/update/version.json（已签名）"
        curl -sS -u "$AUTH" -T "$JSON" "${BASE}/CG/update/version.json" >/dev/null
        rm -f "$JSON"
    fi
    rm -rf "$CLOUD"
    echo "==> 已发布 v${NEW}（主程序 ${APP_SHA:0:12}... / 分发包更新=${LAUNCH_UPLOADED}）"
fi

echo "==> 完成: v${NEW}"
echo "    唯一分发文件: ${DIST}/CG.exe"
