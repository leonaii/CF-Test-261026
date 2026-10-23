#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""CF验号器 v1.3 接口直连工具。

逆向结论（CF验号器_v1.3.exe）：
  1. host.dat 解密得到服务器地址，默认 http://111.170.141.190:8790
  2. 启动后 POST /api/card/boot，body 为空 JSON {}，用于探测服务器与封禁状态
  3. 查询时 POST /api/query/status，body {"code":"<编号或卡密>"}
     - 卡密：16 位大写英文
     - 编号：8 位小写英文
  4. 无签名、无时间戳、无加密信封，UA 不校验，仅要求 Content-Type: application/json

用法：
  python3 cf_check.py --code abcdefgh
  python3 cf_check.py --code ABCDEFGHIJKLMNOP
  python3 cf_check.py --boot
  python3 cf_check.py --code abcdefgh --host http://127.0.0.1:8790
"""

import argparse
import json
import sys
import urllib.error
import urllib.request

DEFAULT_HOST = "http://111.170.141.190:8790"
USER_AGENT = "CfCheck/1.3"
TIMEOUT_SECONDS = 12

# 响应字段与中文含义，用于结果友好展示
FIELD_LABELS = {
    "ok": "请求成功",
    "banned": "是否封禁",
    "zaixian": "在线状态",
    "dongjie": "冻结状态",
    "renlian": "人脸状态",
    "lingpai": "令牌状态",
    "message": "提示信息",
    "code": "编号",
    "qq": "QQ",
    "zone": "区服",
    "note": "备注",
}


class CheckError(Exception):
    """接口调用异常。"""


def post_json(host: str, path: str, payload: dict) -> tuple:
    """发送 JSON POST 请求，返回 (http_status, 响应文本)。

    :param host: 服务器根地址，例如 http://111.170.141.190:8790
    :param path: 接口路径，例如 /api/query/status
    :param payload: 请求体字典
    :return: (状态码, 响应正文)
    """
    url = host.rstrip("/") + path
    body = json.dumps(payload, ensure_ascii=False, separators=(",", ":")).encode("utf-8")
    request = urllib.request.Request(
        url,
        data=body,
        method="POST",
        headers={
            "Content-Type": "application/json",
            "User-Agent": USER_AGENT,
        },
    )
    try:
        with urllib.request.urlopen(request, timeout=TIMEOUT_SECONDS) as response:
            return response.status, response.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as error:
        # 400 等业务错误同样带 JSON 响应体，需要读出来展示
        return error.code, error.read().decode("utf-8", "replace")
    except urllib.error.URLError as error:
        raise CheckError("连接失败：%s（检查服务器地址或网络）" % error.reason)


def parse_result(text: str) -> dict:
    """解析响应 JSON，失败时抛出异常。"""
    try:
        data = json.loads(text)
    except json.JSONDecodeError:
        raise CheckError("响应不是合法 JSON：%s" % text[:200])
    if not isinstance(data, dict):
        raise CheckError("响应结构异常：%s" % text[:200])
    return data


def print_result(status: int, data: dict) -> None:
    """打印结构化结果。"""
    print("HTTP %d" % status)
    print(json.dumps(data, ensure_ascii=False, indent=2))
    print("-" * 40)
    for key, label in FIELD_LABELS.items():
        if key in data:
            print("%-8s %s" % (label, data[key]))


def check_boot(host: str) -> int:
    """调用 /api/card/boot，探测服务器存活与封禁状态。"""
    status, text = post_json(host, "/api/card/boot", {})
    data = parse_result(text)
    print("[boot] %s/api/card/boot" % host.rstrip("/"))
    print_result(status, data)
    return 0 if data.get("ok") else 1


def check_code(host: str, code: str) -> int:
    """调用 /api/query/status 查询编号或卡密。"""
    status, text = post_json(host, "/api/query/status", {"code": code})
    data = parse_result(text)
    print("[query] %s/api/query/status code=%s" % (host.rstrip("/"), code))
    print_result(status, data)
    return 0 if data.get("ok") else 1


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(description="CF验号器 v1.3 接口直连工具")
    parser.add_argument("--code", help="编号（8 位小写英文）或卡密（16 位大写英文）")
    parser.add_argument("--host", default=DEFAULT_HOST, help="服务器地址，默认 %s" % DEFAULT_HOST)
    parser.add_argument("--boot", action="store_true", help="只调用 /api/card/boot 探测服务器")
    parser.add_argument("--json", action="store_true", help="只输出原始 JSON，便于其他程序调用")
    return parser


def main() -> int:
    args = build_parser().parse_args()
    try:
        if args.boot:
            return check_boot(args.host)
        if not args.code:
            print("请提供 --code 或 --boot", file=sys.stderr)
            return 2
        return check_code(args.host, args.code)
    except CheckError as error:
        print("错误：%s" % error, file=sys.stderr)
        return 3


if __name__ == "__main__":
    sys.exit(main())
