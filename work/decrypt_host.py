"""host.dat v2 解密：还原卡密服务器地址。

结构:
  0x00 magic  d8 3a 91 c7
  0x04 ver    02
  0x05 iv     16 bytes
  0x15 len    u16 (ciphertext length)
  0x17 ct     len bytes
  ...  hmac  32 bytes  (HMAC-SHA256(mac_key, file[0:0x17+len]))
明文 = nibble_swap( aes_cbc_dec(ct) ^ xor_mask )
"""
import hashlib
import hmac as hmac_mod
import struct

from pe import read_rva

T1 = struct.unpack("<4I", read_rva(0xFB780, 16))
T2 = struct.unpack("<4I", read_rva(0xFB790, 16))
T3 = struct.unpack("<4I", read_rva(0xFB770, 16))


def derive_seed() -> bytes:
    """复刻 RVA 0x8E300 中的 32 字节种子生成循环（混淆常量）。"""
    seed = bytearray(32)
    seed[0] = 0x92
    b, a, r8, t, esi = 1, 3, 2, 0xE4C89A02, 0x0D
    for i in range(1, 32):
        x = r8 & 3
        t_old = t
        r8 = a & 3
        t = T1[x]
        v = (T2[b & 3] >> ((b & 3) * 8)) ^ (t >> ((a & 3) * 8))
        v ^= T3[a & 3] >> ((~b & 3) * 8)
        v = (((b * b) ^ 0x5A) + esi + v)
        v ^= T2[(b - 1) & 3] >> 11
        v ^= t_old >> 3
        seed[i] = v & 0xFF
        a += 1
        b += 1
        esi += 0x0D
    return bytes(seed)


def derive_keys(seed: bytes):
    return (
        hashlib.sha256(seed + b"aes-host-v2").digest(),
        hashlib.sha256(seed + b"xor-host-v2").digest(),
        hashlib.sha256(seed + b"mac-host-v2").digest(),
    )


def nibble_swap(data: bytes) -> bytes:
    return bytes(((b << 4) & 0xF0) | (b >> 4) for b in data)


def aes_cbc_decrypt(key: bytes, iv: bytes, ct: bytes) -> bytes:
    try:
        from Crypto.Cipher import AES
        return AES.new(key, AES.MODE_CBC, iv).decrypt(ct)
    except ImportError:
        from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
        d = Cipher(algorithms.AES(key), modes.CBC(iv)).decryptor()
        return d.update(ct) + d.finalize()


def decrypt_host(path: str) -> str:
    raw = open(path, "rb").read()
    assert raw[:4] == b"\xd8\x3a\x91\xc7", "magic mismatch"
    assert raw[4] == 2, "version mismatch"
    iv = raw[5:0x15]
    length = struct.unpack_from("<H", raw, 0x15)[0]
    ct = raw[0x17:0x17 + length]
    tag = raw[0x17 + length:0x17 + length + 32]

    seed = derive_seed()
    aes_key, xor_key, mac_key = derive_keys(seed)

    calc = hmac_mod.new(mac_key, raw[:0x17 + length], hashlib.sha256).digest()
    print("seed      :", seed.hex())
    print("hmac ok   :", calc == tag)

    plain = aes_cbc_decrypt(aes_key, iv, ct)
    if plain and 1 <= plain[-1] <= 16:
        plain = plain[:-plain[-1]]

    out = bytearray(plain)
    for i in range(len(out)):
        out[i] ^= (xor_key[(i + 3) & 31] ^ (xor_key[i & 31] + i * 17)) & 0xFF
    return nibble_swap(bytes(out)).decode("utf-8", "replace")


if __name__ == "__main__":
    print("result    :", repr(decrypt_host("/Users/leonos/Downloads/cf/新1.17上号器/host.dat")))
