#!/usr/bin/env bash
# CG 打包脚本
#
#   ./build.sh              自动把版本号 +1（1.0.1 -> 1.0.2）并编译
#   ./build.sh 1.2.0        指定版本号编译
#   ./build.sh --publish    自动 +1 并上传到云端更新目录
#   ./build.sh 1.2.0 --publish
#
# 产物：../dist/CG.exe、../dist/CG-cli.exe
# 云端：CG/update/CG.exe + CG/update/version.json

set -euo pipefail
cd "$(dirname "$0")"

VER_FILE="version.txt"
PUBLISH=0
NEW=""

for arg in "$@"; do
    case "$arg" in
        --publish) PUBLISH=1 ;;
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

echo "==> 编译 CG.exe (windowsgui)"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
    -ldflags "-H=windowsgui ${LDFLAGS}" -o "${DIST}/CG.exe" .

echo "==> 编译 CG-cli.exe (console)"
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath \
    -ldflags "${LDFLAGS}" -o "${DIST}/CG-cli.exe" .

echo "==> 校验内嵌资源"
python3 - "$DIST/CG.exe" "$NEW" <<'PY'
import sys
path, ver = sys.argv[1], sys.argv[2]
data = open(path, "rb").read()
checks = {
    "版本号": ver.encode(),
    "上号器": "圆梦科级上号器".encode(),
    "卡网": b"517fk.cn",
}
for name, needle in checks.items():
    print(f"    {name}: {'OK' if needle in data else '缺失'}")
print(f"    大小: {len(data)/1024/1024:.1f} MB")
PY

if [ "$PUBLISH" = "1" ]; then
    BASE="https://cloud.51i.cc/remote.php/dav/files/leonai"
    AUTH="leonai:QGdWe-otHHx-35wCj-MtCKM-owZkK"
    SHA="$(shasum -a 256 "${DIST}/CG.exe" | cut -d' ' -f1)"
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
