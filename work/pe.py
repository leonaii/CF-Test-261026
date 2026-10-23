import pefile
import struct

PE_PATH = "/Users/leonos/Downloads/cf/新1.17上号器/资源/上号器/yuanmeng/CF上号器_v1.17.exe"
pe = pefile.PE(PE_PATH, fast_load=True)
IMAGE_BASE = 0x140000000


def rva2off(rva: int) -> int:
    return pe.get_offset_from_rva(rva)


def read_rva(rva: int, n: int) -> bytes:
    off = rva2off(rva)
    with open(PE_PATH, "rb") as f:
        f.seek(off)
        return f.read(n)


def cstr_rva(rva: int, maxlen: int = 256) -> str:
    data = read_rva(rva, maxlen)
    end = data.find(b"\x00")
    if end >= 0:
        data = data[:end]
    return data.decode("utf-8", "replace")


if __name__ == "__main__":
    for rva in (0xFB5E0, 0xFB69A, 0xFB6A8, 0xFB6C8, 0xFB6E2, 0xFB770, 0xFB780, 0xFB790, 0xFAC22, 0xFAC66):
        print(hex(rva), repr(cstr_rva(rva, 128)))
    print("tables:")
    for base in (0xFB770, 0xFB780, 0xFB790):
        raw = read_rva(base, 16)
        print(hex(base), struct.unpack("<4I", raw))
