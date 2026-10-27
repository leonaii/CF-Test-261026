package main

// NT 登录第二层：SsoKeyExchange 请求体。
//
// 逐字节复刻样本 sub_14004FB70（IDA 反编译 + 汇编确认）：
//   1. sub_1400B3570 生成本机临时 P-256 密钥对；
//   2. sub_1400B3960 与硬编码服务端公钥（.rdata 0x1400F97A0，65 字节 0x04||X||Y）
//      做 ECDH，得到 32 字节，反转后取前 16 字节当 AES key；
//   3. 明文是 {"mid":"<mid>","timestamp":<unix 秒>}（0x14004FE0B 的 "{\"mid\":"）；
//   4. AES-128-CBC + PKCS7，IV 是 16 字节全零（0x14004FED3 lea r8 到被 pxor 清零的栈缓冲）。
//
// 已确认不是 GCM：sub_1400402D0 里明确 BCryptSetProperty(ChainingModeCBC) +
// BCryptEncrypt(..., 1 = BCRYPT_BLOCK_PADDING)。GCM（nonce 12 + tag 16）是
// sub_140030760，用在握手之后的登录包上。

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// selfNTServerPub 是样本里写死的服务端 P-256 公钥（未压缩点）。
// 来源：.rdata 0x1400F97A0，由 sub_14004FC50 的 r9 传入 sub_1400B3960（a5=65）。
const selfNTServerPubHex = "04c0bc66a3434f824638eb7ada5d1c481f8ebac83825f6c34289a85f18918b5167" +
	"013107266ecb9985d9d1b6c281e55fbfb34be6b96146eb1307a4c2b374840542"

func selfNTServerPub() ([]byte, error) {
	b, err := hex.DecodeString(selfNTServerPubHex)
	if err != nil {
		return nil, fmt.Errorf("内置服务端公钥损坏: %w", err)
	}
	if len(b) != 65 || b[0] != 4 {
		return nil, errors.New("内置服务端公钥不是 65 字节未压缩点")
	}
	return b, nil
}

// selfNTSSOServerPubHex 是 NT SSO 的 SsoKeyExchange 专用服务端公钥。
// 来源：sub_1400B3570 @0x1400B3777，sub_1400B3960(..., &unk_1400F9840, 65)。
// 注意：它与 mid_list 使用的 selfNTServerPubHex（0x1400F97A0）不是同一把。
const selfNTSSOServerPubHex = "049d1423332735980edabe7e9ea451b3395b6f35250db8fc56f25889f628cb" +
	"ae3e8e73077914071eeebc108f4e0170057792bb17aa303af652313d17c1ac815e79"

func selfNTSSOServerPub() ([]byte, error) {
	b, err := hex.DecodeString(selfNTSSOServerPubHex)
	if err != nil {
		return nil, fmt.Errorf("SSO 服务端公钥损坏: %w", err)
	}
	if len(b) != 65 || b[0] != 4 {
		return nil, errors.New("SSO 服务端公钥不是 65 字节未压缩点")
	}
	return b, nil
}

// selfNTKexPlain 构造 SsoKeyExchange 的明文：
// `{"mid":"<mid>","timestamp":<unix>}`（样本 0x14004FE0B 起按字段拼接）。
func selfNTKexPlain(mid string, ts int64) []byte {
	var b bytes.Buffer
	b.WriteString(`{"mid":"`)
	b.WriteString(mid)
	b.WriteString(`","timestamp":`)
	b.WriteString(strconv.FormatInt(ts, 10))
	b.WriteString(`}`)
	return b.Bytes()
}

// selfNTPKCS7 补 PadSize 字节的 PKCS7（AES 块 16）。
func selfNTPKCS7(p []byte) []byte {
	n := aes.BlockSize - len(p)%aes.BlockSize
	out := make([]byte, len(p)+n)
	copy(out, p)
	for i := len(p); i < len(out); i++ {
		out[i] = byte(n)
	}
	return out
}

// selfNTKexPayload 生成 sub_14005CFA0 的 plain 段。
//
// 调用点证据：
//
//	0x14006A3CF/3DA  sub_1400C7C80(rbx, uin) → a2 是 QQ 号十进制文本
//	0x14006A3EF      a3 = [rsp+180h]          → ComputerID
//	0x14006A3DF      a5 = r14d                → 成功路径上 r14 来自 _time64
//
// 因此 plain = clientPub ‖ BE32(1) ‖ GCM(shared, {1:uinText,2:computerID})
//
//	‖ BE32(0) ‖ BE32(timestamp)。
//
// 开头的 65 字节公钥由 0x14005D0B0～0x14005D143 确认，不能省略。
func selfNTKexPayload(shared, uinText, computerID, clientPub []byte, timestamp uint32) (plain, innerGCM []byte, err error) {
	inner := selfNTAppendVarint(nil, 0x0A)
	inner = selfNTAppendBytes(inner, uinText)
	inner = selfNTAppendVarint(inner, 0x12)
	inner = selfNTAppendBytes(inner, computerID)

	innerGCM, err = selfNTGCMSeal(shared, inner)
	if err != nil {
		return nil, nil, err
	}
	plain = make([]byte, 0, len(clientPub)+12+len(innerGCM))
	plain = append(plain, clientPub...)
	plain = append(plain, selfNTBE32Bytes(1)...)
	plain = append(plain, innerGCM...)
	plain = append(plain, selfNTBE32Bytes(0)...)
	plain = append(plain, selfNTBE32Bytes(timestamp)...)
	return plain, innerGCM, nil
}

// selfNTGCMSeal 是 sub_140030760 的输出布局：nonce(12) ‖ 密文 ‖ tag(16)。
func selfNTGCMSeal(key, plain []byte) ([]byte, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(blk)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, g.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	sealed := g.Seal(nil, nonce, plain, nil)
	out := make([]byte, 0, len(nonce)+len(sealed))
	out = append(out, nonce...)
	out = append(out, sealed...)
	return out, nil
}

// selfNTGCMOpen 是 sub_14003FA50：nonce(12) ‖ 密文 ‖ tag(16)，无 AAD。
func selfNTGCMOpen(key, sealed []byte) ([]byte, error) {
	if len(sealed) < 28 {
		return nil, fmt.Errorf("GCM 数据太短: %d", len(sealed))
	}
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	g, err := cipher.NewGCM(blk)
	if err != nil {
		return nil, err
	}
	return g.Open(nil, sealed[:g.NonceSize()], sealed[g.NonceSize():], nil)
}

// selfNTKexBody 完整复刻 sub_14005CFA0 的外层 protobuf。
//
//	ctx+0x10 = 本机 P-256 公钥（65 字节）
//	ctx+0x28 = ECDH 共享密钥 ptr，+0x30 = 长度
//
//	out = {1:clientPub, 2:1, 3:innerGCM, 4:timestamp, 5:GCM(signKey,SHA256(plain))}
//
// 注意两层 GCM 的 key 不同：plain 内层用 ECDH 共享密钥；
// 字段 5 用全局 pbSecret（0x14005D30C `lea rcx,pbSecret`，仅这一处引用）。
var selfNTKexSignKey = []byte{
	0xe2, 0x73, 0x3b, 0xf4, 0x03, 0x14, 0x99, 0x13,
	0xcb, 0xf8, 0x0c, 0x7a, 0x95, 0x16, 0x8b, 0xd4,
	0xca, 0x69, 0x35, 0xee, 0x53, 0xcd, 0x39, 0x76,
	0x4b, 0xee, 0xbe, 0x2e, 0x00, 0x7e, 0x3a, 0xee,
}

func selfNTKexBody(pbSecret, uinText, computerID, clientPub []byte, timestamp uint32) ([]byte, error) {
	plain, innerGCM, err := selfNTKexPayload(pbSecret, uinText, computerID, clientPub, timestamp)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(plain)
	g2, err := selfNTGCMSeal(selfNTKexSignKey, sum[:])
	if err != nil {
		return nil, err
	}

	out := selfNTAppendVarint(nil, 0x0A)
	out = selfNTAppendBytes(out, clientPub)
	out = selfNTAppendVarint(out, 0x10) // field 2: varint 1
	out = selfNTAppendVarint(out, 1)
	out = selfNTAppendVarint(out, 0x1A) // field 3: innerGCM
	out = selfNTAppendBytes(out, innerGCM)
	out = selfNTAppendVarint(out, 0x20) // field 4: timestamp
	out = selfNTAppendVarint(out, uint64(timestamp))
	out = selfNTAppendVarint(out, 0x2A) // field 5: GCM(SHA256(plain))
	out = selfNTAppendBytes(out, g2)
	return out, nil
}

// selfNTComputerID 复刻 sub_140038D90：
// 调用方给出的值恰好 16 字节时直接使用，否则计算
// MD5(十进制 uin + "ComputerID")；0x1400F72EA 的算法名是 "MD5"。
func selfNTComputerID(uin uint64, supplied []byte) []byte {
	if len(supplied) == 16 {
		return append([]byte(nil), supplied...)
	}
	sum := md5.Sum([]byte(strconv.FormatUint(uin, 10) + "ComputerID"))
	return sum[:]
}

// selfNTBuildKex 生成一次完整的 SsoKeyExchange 请求体。
// 同时返回本机密钥对与共享密钥，后续 EasyLogin 要用同一份。
func selfNTBuildKex(uin uint64, computerID []byte, ts int64) (body []byte, kp *selfNTKeyPair, secret []byte, err error) {
	if ts <= 0 {
		ts = time.Now().Unix()
	}
	cid := selfNTComputerID(uin, computerID)
	kp, err = selfNTNewKeyPair()
	if err != nil {
		return nil, nil, nil, err
	}
	peer, err := selfNTSSOServerPub()
	if err != nil {
		return nil, nil, nil, err
	}
	secret, err = selfNTSharedSecret(kp.priv, peer)
	if err != nil {
		return nil, nil, nil, err
	}
	// selfNTSharedSecret 保留 mid_list 所需的反转形式；NT SSO 的 GCM
	// 使用未反转的 ECDH X 坐标。真实 SsoKeyExchange 已验证此方向。
	selfNTReverseBytes(secret)
	body, err = selfNTKexBody(secret, []byte(strconv.FormatUint(uin, 10)), cid, kp.Uncompressed(), uint32(ts))
	return body, kp, secret, err
}

// selfNTBuildKexFrame 直接给出可发送的 trpc 帧（cmd=0x2711）。
func selfNTBuildKexFrame(uin uint64, computerID []byte, ts int64) (frame []byte, kp *selfNTKeyPair, secret []byte, err error) {
	if ts <= 0 {
		ts = time.Now().Unix()
	}
	body, kp, secret, err := selfNTBuildKex(uin, computerID, ts)
	if err != nil {
		return nil, nil, nil, err
	}
	cid := selfNTComputerID(uin, computerID)
	frame = selfNTCallFrameWithA5(0x2711, "trpc.login.ecdh.EcdhService.SsoKeyExchange", body, []byte(hex.EncodeToString(cid)), strconv.FormatInt(ts+1, 10))
	return frame, kp, secret, nil
}

func selfNTReverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func selfNTKexReadVarint(b []byte, off *int) (uint64, bool) {
	var v uint64
	for i := 0; i < 10 && *off < len(b); i++ {
		c := b[*off]
		*off = *off + 1
		v |= uint64(c&0x7f) << (7 * i)
		if c < 0x80 {
			return v, true
		}
	}
	return 0, false
}

// selfNTKexBytesField 读取一个顶层 protobuf bytes 字段。
func selfNTKexBytesField(b []byte, want uint64) ([]byte, error) {
	for off := 0; off < len(b); {
		tag, ok := selfNTKexReadVarint(b, &off)
		if !ok {
			return nil, errors.New("protobuf tag 截断")
		}
		field, wire := tag>>3, tag&7
		switch wire {
		case 0:
			if _, ok := selfNTKexReadVarint(b, &off); !ok {
				return nil, errors.New("protobuf varint 截断")
			}
		case 1:
			off += 8
		case 2:
			n, ok := selfNTKexReadVarint(b, &off)
			if !ok || n > uint64(len(b)-off) {
				return nil, errors.New("protobuf bytes 长度越界")
			}
			if field == want {
				return append([]byte(nil), b[off:off+int(n)]...), nil
			}
			off += int(n)
		case 5:
			off += 4
		default:
			return nil, fmt.Errorf("不支持的 protobuf wire=%d", wire)
		}
		if off > len(b) {
			return nil, errors.New("protobuf 字段越界")
		}
	}
	return nil, fmt.Errorf("protobuf 缺少字段 %d", want)
}

type selfNTKexResult struct {
	TicketA       []byte
	TicketB       []byte
	SessionSecret []byte
}

// selfNTParseKexResponse 复刻 sub_1400564C0：
// field3 服务端公钥 → 第二次 ECDH → GCM 解 field1 → 内层 field1/field2。
func selfNTParseKexResponse(plain []byte, kp *selfNTKeyPair) (*selfNTKexResult, error) {
	if kp == nil || kp.priv == nil {
		return nil, errors.New("KEX 本机密钥对为空")
	}
	payload, err := selfNTExtractPayload(plain)
	if err != nil {
		return nil, err
	}
	if len(payload) == 0 {
		return nil, errors.New("KEX 回包无业务数据")
	}
	enc, err := selfNTKexBytesField(payload, 1)
	if err != nil {
		return nil, err
	}
	serverPub, err := selfNTKexBytesField(payload, 3)
	if err != nil {
		return nil, err
	}
	secret, err := selfNTSharedSecret(kp.priv, serverPub)
	if err != nil {
		return nil, fmt.Errorf("KEX 回包 ECDH 失败: %w", err)
	}
	selfNTReverseBytes(secret)
	decrypted, err := selfNTGCMOpen(secret, enc)
	if err != nil {
		return nil, fmt.Errorf("KEX 回包 GCM 解密失败: %w", err)
	}
	a, err := selfNTKexBytesField(decrypted, 1)
	if err != nil {
		return nil, fmt.Errorf("KEX 票据 A: %w", err)
	}
	b, err := selfNTKexBytesField(decrypted, 2)
	if err != nil {
		return nil, fmt.Errorf("KEX 票据 B: %w", err)
	}
	return &selfNTKexResult{TicketA: a, TicketB: b, SessionSecret: secret}, nil
}

// selfNTSealKex 用本机密钥对 + 内置服务端公钥加密 KeyExchange 明文。
// 返回 16 字节 AES key（会话后续要用）和 CBC 密文。
func selfNTSealKex(kp *selfNTKeyPair, plain []byte) (key, ct []byte, err error) {
	if kp == nil {
		return nil, nil, errors.New("密钥对为空")
	}
	pub, err := selfNTServerPub()
	if err != nil {
		return nil, nil, err
	}
	shared, err := kp.SharedSecret(pub)
	if err != nil {
		return nil, nil, err
	}
	if len(shared) < 16 {
		return nil, nil, fmt.Errorf("共享密钥长度异常: %d", len(shared))
	}
	key = shared[:16]
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	iv := make([]byte, aes.BlockSize) // 样本里就是全零
	buf := selfNTPKCS7(plain)
	out := make([]byte, len(buf))
	cipher.NewCBCEncrypter(blk, iv).CryptBlocks(out, buf)
	return key, out, nil
}

// selfNTOpenKex 是 selfNTSealKex 的逆操作，用于自检/解回包。
func selfNTOpenKex(key, ct []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, fmt.Errorf("key 长度应为 16，实际 %d", len(key))
	}
	if len(ct) == 0 || len(ct)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("密文长度不是 16 的倍数: %d", len(ct))
	}
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	out := make([]byte, len(ct))
	cipher.NewCBCDecrypter(blk, iv).CryptBlocks(out, ct)
	// 去 PKCS7
	n := int(out[len(out)-1])
	if n <= 0 || n > aes.BlockSize || n > len(out) {
		return nil, errors.New("PKCS7 填充不合法")
	}
	for _, c := range out[len(out)-n:] {
		if int(c) != n {
			return nil, errors.New("PKCS7 填充内容不一致")
		}
	}
	return out[:len(out)-n], nil
}
