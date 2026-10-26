#!/usr/bin/env bash
# CG 打包脚本
#
#   ./build.sh              自动把版本号 +1（1.0.1 -> 1.0.2）并编译
#   ./build.sh 1.2.0        指定版本号编译
#   ./build.sh --publish    自动 +1 并上传到云端更新目录
#   ./build.sh 1.2.0 --publish
#   ./build.sh --no-loader  跳过上号器资源同步（离线打包时用）
#
# 产物：../dist/CG.exe、../dist/CG-cli.exe
# 云端：CG/update/CG.exe + CG/update/version.json

set -euo pipefail
cd "$(dirname "$0")"

VER_FILE="version.txt"
PUBLISH=0
SYNC_LOADER=1
NEW=""

for arg in "$@"; do
    case "$arg" in
        --publish) PUBLISH=1 ;;
        --no-loader) SYNC_LOADER=0 ;;
        -h|--help) sed -n '2,12p' "$0"; exit 0 ;;
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

LDFLAGS="-s -w -X main.toolVersion=${NEW}"
DIST="../dist"
mkdir -p "$DIST"

# macOS 只有 shasum，Git Bash/Linux 只有 sha256sum
sha256() {
    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$1" | cut -d' ' -f1
    else
        sha256sum "$1" | cut -d' ' -f1
    fi
}

echo "==> 编译 CG.exe (windowsgui)"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
    -ldflags "-H=windowsgui ${LDFLAGS}" -o "${DIST}/CG.exe" .

echo "==> 编译 CG-cli.exe (console)"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
    -ldflags "${LDFLAGS}" -o "${DIST}/CG-cli.exe" .

echo "==> 校验内嵌资源"
if command -v python3 >/dev/null 2>&1 && python3 -V >/dev/null 2>&1; then
python3 - "$DIST/CG.exe" "$NEW" <<'PY'
import sys
path, ver = sys.argv[1], sys.argv[2]
data = open(path, "rb").read()
checks = {
    "版本号": ver.encode(),
    "上号器类型": "圆梦科级上号器".encode(),
    "卡网": b"517fk.cn",
}
for name, needle in checks.items():
    print(f"    {name}: {'OK' if needle in data else '缺失'}")
print(f"    大小: {len(data)/1024/1024:.1f} MB")
PY
else
    echo "    跳过（本机 python3 不可用）"
fi

BASE="https://cloud.51i.cc/remote.php/dav/files/leonai"
AUTH="leonai:QGdWe-otHHx-35wCj-MtCKM-owZkK"
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
  "name": "圆梦科级上号器",
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
    JSON="$(mktemp)"
    cat > "$JSON" <<EOF
{
  "version": "${NEW}",
  "file": "CG.exe",
  "sha256": "${SHA}",
  "size": ${SIZE},
  "notes": "build $(date '+%Y-%m-%d %H:%M')"
}
EOF

    echo "==> 上传 CG/update/CG.exe"
    curl -sS -u "$AUTH" -X MKCOL "${BASE}/CG/update/" >/dev/null || true
    curl -sS -u "$AUTH" -T "${DIST}/CG.exe" "${BASE}/CG/update/CG.exe" >/dev/null
    echo "==> 上传 CG/update/version.json"
    curl -sS -u "$AUTH" -T "$JSON" "${BASE}/CG/update/version.json" >/dev/null
    rm -f "$JSON"
    echo "==> 已发布 v${NEW}（sha256 ${SHA:0:12}...）"
fi

echo "==> 完成: ${DIST}/CG.exe v${NEW}"
