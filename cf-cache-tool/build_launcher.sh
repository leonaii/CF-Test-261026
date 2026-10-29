#!/usr/bin/env bash
# 只重编译/发布 Rust 启动器，不重建主程序。
#
#   ./build_launcher.sh
#   ./build_launcher.sh --publish
#
# 主程序版本和 CG/app/CG-app.exe.gz 保持不变，只补 ECDSA 清单签名。
#
# 外观复制顺序（不可调换）：
#   链接 → stamp_exe.go -icon -version → UPX → stamp_exe.go -sig -time
#
# 1) 图标/版本信息走 UpdateResource，会重写 .rsrc，必须在 UPX 之前；先压缩再改会
#    破坏 UPX 的节布局。
# 2) 签名是往文件尾追加 PE 证书表并改写 DataDirectory[4]，之后任何写文件都会让
#    内容与签名彻底不符，只能放最后。
# 3) 文件时间要排在 writeChecksum 之后，写文件会把 mtime 刷成当前时间。

set -euo pipefail
cd "$(dirname "$0")"

PUBLISH=0
for arg in "$@"; do
    case "$arg" in
        --publish) PUBLISH=1 ;;
        -h|--help) sed -n '2,16p' "$0"; exit 0 ;;
        *) echo "未知参数: $arg"; exit 1 ;;
    esac
done

if [ -f ".build.env" ]; then
    set -a
    # shellcheck disable=SC1091
    . ./.build.env
    set +a
fi
if [ ! -f ".cg-share.env" ]; then
    echo "错误: 缺少 .cg-share.env，请先执行 ./build.sh 生成只读分享"
    exit 1
fi
set -a
# shellcheck disable=SC1091
. ./.cg-share.env
set +a
if [ -z "${DAV_USER:-}" ] || [ -z "${DAV_PASS:-}" ] || [ -z "${CG_SHARE_TOKEN:-}" ] || [ -z "${CG_SHARE_PASS:-}" ]; then
    echo "错误: 构建凭据不完整"
    exit 1
fi

VER="$(tr -d ' \t\r\n' < version.txt)"
DIST="../dist"
mkdir -p "$DIST"
rm -f "${DIST}/CG-app.exe" "${DIST}/CG-cli.exe" "${DIST}/CG-test.exe" "${DIST}/CG.exe"
rm -rf "${DIST}/app" "${DIST}/cloud"

sha256() {
    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$1" | cut -d' ' -f1
    else
        sha256sum "$1" | cut -d' ' -f1
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
UPX="$(find_upx)"

cleanup() { rm -f launcher/rsrc.syso; }
trap cleanup EXIT

echo "==> 生成图标/版本资源"
CGO_ENABLED=0 go run build_resources.go "$VER"

echo "==> 编译 Rust 启动器"
(
    cd launcher
    export CG_SHARE_TOKEN CG_SHARE_PASS CG_LAUNCHER_VERSION="$VER"
    export VCVARS64="${VCVARS64:-C:\\Program Files (x86)\\Microsoft Visual Studio\\2022\\BuildTools\\VC\\Auxiliary\\Build\\vcvars64.bat}"
    cmd.exe //c build.cmd
)
cp launcher/target/release/cg-launcher.exe "${DIST}/CG租号.exe"

# 图标 / 版本信息必须在 UPX 之前写入：UpdateResource 会重写 .rsrc，
# 先压缩再改会破坏 UPX 的节布局。签名只能放在最后（追加到文件尾）。
DONOR="exe/zuhaoclient.exe"
if [ -f "$DONOR" ]; then
    echo "==> 复制图标/版本信息（参考 $(basename "$DONOR")）"
    CGO_ENABLED=0 go run stamp_exe.go -from "$DONOR" -to "${DIST}/CG租号.exe" -icon -version
else
    echo "==> 跳过图标/版本复制：缺少 $DONOR"
fi

if [ -n "$UPX" ]; then
    echo "==> UPX 压缩"
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

SHA="$(sha256 "${DIST}/CG租号.exe")"
SIZE="$(wc -c < "${DIST}/CG租号.exe" | tr -d ' ')"
TOOLCHAIN="$(rustc -Vv 2>/dev/null | tr '\n' ' ' || true)"
SECRET_FP="$(printf '%s\n%s' "$CG_SHARE_TOKEN" "$CG_SHARE_PASS" | sha256_stdin)"
LAUNCHER_FP="$(
    {
        printf '%s\n%s\n%s\n' "$TOOLCHAIN" "$SECRET_FP" "$VER"
        for f in launcher/Cargo.toml launcher/Cargo.lock launcher/build.rs launcher/build.cmd \
                 launcher/src/main.rs launcher/src/crypto.rs launcher/src/cng.rs launcher/src/net.rs \
                 launcher/src/gzip.rs build_resources.go stamp_exe.go exe/zuhaoclient.exe assets/zuhaoclient.ico \
                 assets/icon_256.png assets/icon_64.png assets/icon_48.png; do
            sha256 "$f"
        done
    } | sha256_stdin
)"

echo "    启动器: ${DIST}/CG租号.exe"
echo "    大小: $((SIZE/1024)) KB"
echo "    SHA256: $SHA"
echo "    分发包指纹: $LAUNCHER_FP"

if [ "$PUBLISH" != "1" ]; then
    echo "==> 未发布（加 --publish 才上传云端）"
    exit 0
fi

BASE="https://cloud.51i.cc/remote.php/dav/files/leonai"
AUTH="${DAV_USER}:${DAV_PASS}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"; cleanup' EXIT

echo "==> 给现有主程序清单补 ECDSA 签名"
curl -fsS -u "$AUTH" -o "$TMP_DIR/app.json" "${BASE}/CG/app/version.json"
CGO_ENABLED=0 go run sign_update_ecdsa.go add "$TMP_DIR/app.json" "$TMP_DIR/app-signed.json"
curl -fsS -u "$AUTH" -T "$TMP_DIR/app-signed.json" "${BASE}/CG/app/version.json" >/dev/null

echo "==> 生成启动器更新清单"
SIG="$(CGO_ENABLED=0 go run sign_update.go "$VER" "$SHA" "$SIZE")"
SIG_ECDSA="$(CGO_ENABLED=0 go run sign_update_ecdsa.go sign "$VER" "$SHA" "$SIZE")"
cat > "$TMP_DIR/update.json" <<EOF
{
  "version": "${VER}",
  "file": "CG.exe",
  "sha256": "${SHA}",
  "size": ${SIZE},
  "launcher_fp": "${LAUNCHER_FP}",
  "notes": "launcher $(date '+%Y-%m-%d %H:%M')",
  "sig": "${SIG}",
  "sig_ecdsa": "${SIG_ECDSA}"
}
EOF

curl -fsS -u "$AUTH" -X MKCOL "${BASE}/CG/update/" >/dev/null 2>&1 || true
echo "==> 上传 CG/update/CG.exe"
curl -fsS -u "$AUTH" -T "${DIST}/CG租号.exe" "${BASE}/CG/update/CG.exe" >/dev/null
echo "==> 上传 CG/update/version.json"
curl -fsS -u "$AUTH" -T "$TMP_DIR/update.json" "${BASE}/CG/update/version.json" >/dev/null
echo "==> 启动器发布完成（主程序未改动）"
