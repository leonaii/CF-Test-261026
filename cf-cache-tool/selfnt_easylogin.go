package main

// NT 登录：SsoKeyExchange → SsoNTLoginEasyLogin → QQConnectLogin.auth。
// EasyLogin 复刻 sub_1400586F0，FirstAuth 复刻 sub_14005B890/sub_14005A240。

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// selfNTEasyDevice 保留旧默认设备类型；QQ 主路径使用 selfNTEasyDeviceType(..., 5)。
func selfNTEasyDevice(payload []byte) []byte {
	return selfNTEasyDeviceType(payload, 4)
}

func selfNTEasyDeviceType(payload []byte, deviceType uint64) []byte {
	out := make([]byte, 0, 128+len(payload))
	str := func(field int, s string) {
		out = selfNTAppendVarint(out, uint64(field)<<3|2)
		out = selfNTAppendBytes(out, []byte(s))
	}
	num := func(field int, v uint64) {
		out = selfNTAppendVarint(out, uint64(field)<<3)
		out = selfNTAppendVarint(out, v)
	}
	bytesField := func(field int, b []byte) {
		out = selfNTAppendVarint(out, uint64(field)<<3|2)
		out = selfNTAppendBytes(out, b)
	}
	str(1, "Windows")
	str(2, "Administrator")
	num(3, deviceType)
	bytesField(4, payload)
	str(5, "6.1.7601")
	str(6, "MS-20220623RAMD")
	return out
}

// selfNTSealTickets 把两张票封成载荷：0x01 ‖ GCM(protobuf{1:票A,2:票B}) ‖ 0x00。
//
// key 长度：sub_14005CFA0 末尾那句是
//
//	sub_140030760(pbSecret, 0x20u, ...)   ← 明确 0x20 = 32 字节
//
// 同一函数前面的第一处 GCM 用的是 a4[5]..a4[6]（即 32 字节共享密钥），
// 所以这里取 shared 的前 32 字节；不足 32 时退回 16。
func selfNTSealTickets(shared []byte, ticketA string, ticketB []byte) ([]byte, error) {
	keyLen := 32
	if len(shared) < 32 {
		keyLen = 16
	}
	if len(shared) < keyLen {
		return nil, fmt.Errorf("共享密钥太短: %d", len(shared))
	}
	inner := selfNTAppendVarint(nil, 0x0A)
	inner = selfNTAppendBytes(inner, []byte(ticketA))
	inner = selfNTAppendVarint(inner, 0x12)
	inner = selfNTAppendBytes(inner, ticketB)

	blk, err := aes.NewCipher(shared[:keyLen])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blk)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize()) // 12
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	sealed := gcm.Seal(nil, nonce, inner, nil) // 密文 ‖ tag(16)

	out := make([]byte, 0, 1+len(nonce)+len(sealed)+1)
	out = append(out, 0x01)
	out = append(out, nonce...)
	out = append(out, sealed...)
	out = append(out, 0x00)
	return out, nil
}

// selfNTOpenTickets 是 selfNTSealTickets 的逆操作（自检/解回包用）。
func selfNTOpenTickets(shared, payload []byte) (string, []byte, error) {
	keyLen := 32
	if len(shared) < 32 {
		keyLen = 16
	}
	if len(shared) < keyLen {
		return "", nil, errors.New("共享密钥太短")
	}
	if len(payload) < 1+12+16+1 || payload[0] != 0x01 {
		return "", nil, errors.New("载荷头不对")
	}
	body := payload[1 : len(payload)-1]
	nonce, ct := body[:12], body[12:]
	blk, err := aes.NewCipher(shared[:keyLen])
	if err != nil {
		return "", nil, err
	}
	gcm, err := cipher.NewGCM(blk)
	if err != nil {
		return "", nil, err
	}
	inner, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", nil, err
	}
	// protobuf{1:string,2:bytes}
	readField := func(b []byte, want int) ([]byte, []byte, error) {
		tag := b[0]
		if int(tag>>3) != want {
			return nil, nil, fmt.Errorf("字段顺序不对: %#x", tag)
		}
		l := 0
		shift := uint(0)
		i := 1
		for {
			c := b[i]
			i++
			l |= int(c&0x7f) << shift
			if c < 0x80 {
				break
			}
			shift += 7
		}
		if i+l > len(b) {
			return nil, nil, errors.New("长度越界")
		}
		return b[i : i+l], b[i+l:], nil
	}
	a, rest, err := readField(inner, 1)
	if err != nil {
		return "", nil, err
	}
	bb, _, err := readField(rest, 2)
	if err != nil {
		return "", nil, err
	}
	return string(a), bb, nil
}

// selfNTEasyFingerprint 复刻 sub_140057020 的 a12=0 分支。
// H1=MD5(seed) 写入明文；TEA key=MD5(H1‖BE32(0)‖BE32(uin))。
func selfNTEasyFingerprint(uin uint64, computerID []byte, timestamp uint32, seed []byte) ([]byte, error) {
	r4 := make([]byte, 4)
	r16 := make([]byte, 16)
	if _, err := rand.Read(r4); err != nil {
		return nil, err
	}
	if _, err := rand.Read(r16); err != nil {
		return nil, err
	}
	uinText := strconv.FormatUint(uin, 10)
	seedHash := md5.Sum(seed)
	keyInput := make([]byte, 0, 24)
	keyInput = append(keyInput, seedHash[:]...)
	keyInput = append(keyInput, selfNTBE32Bytes(0)...)
	keyInput = append(keyInput, selfNTBE32Bytes(uint32(uin))...)
	key := md5.Sum(keyInput)
	plain := make([]byte, 0, 160+len(computerID)+len(uinText))
	plain = append(plain, 0, 4) // sub_14002C3D0 是 BE16(4)，不是四个零
	plain = append(plain, r4...)
	plain = append(plain, selfNTBE32Bytes(0)...)
	plain = append(plain, selfNTBE32Bytes(1600001614)...)
	plain = append(plain, selfNTBE32Bytes(8001)...)
	plain = append(plain, selfNTBE32Bytes(0)...)
	plain = append(plain, selfNTBE32Bytes(uint32(uin))...)
	plain = append(plain, selfNTBE32Bytes(timestamp)...)
	plain = append(plain, selfNTBE32Bytes(0)...)
	plain = append(plain, 0)
	plain = append(plain, seedHash[:]...)
	plain = append(plain, r16...)
	plain = append(plain, selfNTBE32Bytes(0)...)
	plain = append(plain, 0)
	plain = append(plain, computerID...)
	plain = append(plain, selfNTBE32Bytes(1)...)
	plain = append(plain, selfNTBE32Bytes(0)...)
	plain = append(plain, byte(len(uinText)>>8), byte(len(uinText))) // BE16 长度
	plain = append(plain, uinText...)
	return selfNTTEAEncryptKey(selfNTTEAWrap(plain), key[:]), nil
}

// selfNTEasyLoginBody 复刻 sub_1400586F0：
// 外层 {1:TicketB, 3:GCM(TicketA, inner), 4:1}；inner.field2.field1 是 Token0106。
// 与 PasswordLogin 不同，这里不包含现场生成的设备指纹。
func selfNTEasyLoginBody(uin uint64, computerID, token0106, ticketA, ticketB []byte) ([]byte, error) {
	if len(ticketA) != 32 {
		return nil, fmt.Errorf("EasyLogin TicketA 应为 32 字节，实际 %d", len(ticketA))
	}
	if len(ticketB) == 0 {
		return nil, errors.New("EasyLogin TicketB 为空")
	}
	if len(token0106) == 0 {
		return nil, errors.New("EasyLogin Token0106 为空")
	}
	cid := selfNTComputerID(uin, computerID)
	bytesField := func(out []byte, field int, data []byte) []byte {
		out = selfNTAppendVarint(out, uint64(field)<<3|2)
		return selfNTAppendBytes(out, data)
	}
	varintField := func(out []byte, field int, v uint64) []byte {
		out = selfNTAppendVarint(out, uint64(field)<<3)
		return selfNTAppendVarint(out, v)
	}

	uinMsg := bytesField(nil, 1, []byte(strconv.FormatUint(uin, 10)))
	device := selfNTEasyDeviceType(cid, 4)
	client := bytesField(nil, 1, []byte("1.0"))
	client = varintField(client, 2, 1600001614)
	client = bytesField(client, 3, []byte("com.tencent.wegame"))
	client = bytesField(client, 4, []byte("nt_pc_sdk"))
	client = bytesField(client, 5, []byte("8.9.8"))

	login := bytesField(nil, 1, uinMsg)
	login = bytesField(login, 2, device)
	login = bytesField(login, 3, client)
	tokenMsg := bytesField(nil, 1, token0106)
	inner := bytesField(nil, 1, login)
	inner = bytesField(inner, 2, tokenMsg)

	sealed, err := selfNTGCMSeal(ticketA, inner)
	if err != nil {
		return nil, err
	}
	out := bytesField(nil, 1, ticketB)
	out = bytesField(out, 3, sealed)
	out = varintField(out, 4, 1)
	return out, nil
}

func selfNTEasyLoginFrame(uin uint64, computerID, token0106, ticketA, ticketB []byte, ts int64) ([]byte, error) {
	if ts <= 0 {
		ts = time.Now().Unix()
	}
	body, err := selfNTEasyLoginBody(uin, computerID, token0106, ticketA, ticketB)
	if err != nil {
		return nil, err
	}
	cid := selfNTComputerID(uin, computerID)
	return selfNTCallFrameWithA5(0x2712, "trpc.login.ecdh.EcdhService.SsoNTLoginEasyLogin", body, []byte(hex.EncodeToString(cid)), strconv.FormatInt(ts+1, 10)), nil
}

type selfNTSession struct {
	Token0106 []byte
	TGT       []byte
	D2        []byte
	D2K       []byte
}

// selfNTParseEasyLoginResponse 解开 field3，并读取 LoginBody.field1 的会话材料。
func selfNTParseEasyLoginResponse(plain, ticketA []byte) (*selfNTSession, error) {
	payload, err := selfNTExtractPayload(plain)
	if err != nil {
		return nil, err
	}
	if len(payload) == 0 {
		return nil, errors.New("EasyLogin 回包无业务数据")
	}
	enc, err := selfNTKexBytesField(payload, 3)
	if err != nil {
		return nil, err
	}
	dec, err := selfNTGCMOpen(ticketA, enc)
	if err != nil {
		return nil, fmt.Errorf("EasyLogin 回包 GCM 解密失败: %w", err)
	}
	head, err := selfNTKexBytesField(dec, 1)
	if err != nil {
		return nil, err
	}
	if status, _ := selfNTKexBytesField(head, 4); len(status) > 0 {
		var code uint64
		if len(status) < 2 || status[0] != 0x08 {
			return nil, errors.New("EasyLogin 状态字段格式错误")
		}
		off := 1
		v, ok := selfNTKexReadVarint(status, &off)
		if !ok {
			return nil, errors.New("EasyLogin 状态码截断")
		}
		code = v
		if code != 0 {
			msg, _ := selfNTKexBytesField(status, 3)
			if len(msg) == 0 {
				msg, _ = selfNTKexBytesField(status, 2)
			}
			if len(msg) == 0 {
				msg = []byte("未知错误")
			}
			return nil, fmt.Errorf("EasyLogin 失败(%d): %s", code, msg)
		}
	}
	body, err := selfNTKexBytesField(dec, 2)
	if err != nil {
		return nil, err
	}
	sm, err := selfNTKexBytesField(body, 1)
	if err != nil {
		return nil, err
	}
	s := &selfNTSession{}
	s.Token0106, _ = selfNTKexBytesField(sm, 3)
	if s.TGT, err = selfNTKexBytesField(sm, 4); err != nil {
		return nil, fmt.Errorf("EasyLogin 缺少 tgt: %w", err)
	}
	if s.D2, err = selfNTKexBytesField(sm, 5); err != nil {
		return nil, fmt.Errorf("EasyLogin 缺少 d2: %w", err)
	}
	if s.D2K, err = selfNTKexBytesField(sm, 6); err != nil {
		return nil, fmt.Errorf("EasyLogin 缺少 d2k: %w", err)
	}
	if len(s.TGT) == 0 || len(s.D2) < 16 || len(s.D2K) < 16 {
		return nil, errors.New("EasyLogin 未返回完整会话")
	}
	return s, nil
}

// EasyLogin 主流程 sub_140069AF0 在 0x14006C6EC 传 dword_1400DA09C。
// 这里是客户端版本号，不是 QQ 主流程 sub_14006FC00 使用的小 appid 表。
var selfNTFirstAuthAppIDs = [...]uint32{1600001614}

// selfNTFirstAuthFrame 生成 cmd=10003 / QQConnectLogin.auth，TEA key=d2k[:16]。
// 外层为 BE32(12)‖01‖LP(d2)‖00‖LP(uin)‖TEA，见 sub_14005A240 尾部。
func selfNTFirstAuthFrame(uin uint64, guid []byte, connSeq, uid string, tgt, d2, d2k []byte, appID uint32) ([]byte, error) {
	if len(tgt) == 0 || len(d2) == 0 || len(d2k) < 16 {
		return nil, errors.New("FIRST_AUTH 会话材料不完整")
	}
	bytesField := func(out []byte, field int, data []byte) []byte {
		out = selfNTAppendVarint(out, uint64(field)<<3|2)
		return selfNTAppendBytes(out, data)
	}
	varintField := func(out []byte, field int, v uint64) []byte {
		out = selfNTAppendVarint(out, uint64(field)<<3)
		return selfNTAppendVarint(out, v)
	}
	login := varintField(nil, 1, 8)
	login = bytesField(login, 2, []byte(hex.EncodeToString(tgt)))
	login = varintField(login, 3, uint64(appID))
	payload := varintField(nil, 1, 1)
	payload = varintField(payload, 2, 101478283)
	payload = varintField(payload, 3, 1010)
	payload = bytesField(payload, 4, []byte("win"))
	payload = bytesField(payload, 5, []byte("nt_pc_sdk"))
	payload = bytesField(payload, 6, []byte("1.8"))
	payload = bytesField(payload, 7, []byte("8.9.8"))
	payload = varintField(payload, 20, 8)
	payload = varintField(payload, 21, 0)
	payload = bytesField(payload, 65, login)

	hdr := selfNTHeader(payload, []byte(selfNTRandHeader15()), []byte(connSeq))
	if uid != "" {
		hdr = selfNTAppendVarint(hdr, 16<<3|2)
		hdr = selfNTAppendBytes(hdr, []byte(uid))
	}
	const svc = "QQConnectLogin.auth"
	body := make([]byte, 0, 256+len(tgt)+len(guid)+len(hdr))
	body = append(body, selfNTBE32Bytes(10003)...)
	body = append(body, selfNTBE32Bytes(selfNTCmdMagic)...)
	body = append(body, selfNTBE32Bytes(selfNTLocale)...)
	body = append(body, 0x02)
	body = append(body, make([]byte, 11)...)
	body = append(body, selfNTLenPrefixed(tgt)...)
	body = append(body, selfNTLenPrefixed([]byte(svc))...)
	body = append(body, selfNTLenPrefixed(nil)...)
	body = append(body, selfNTLenPrefixed(guid)...)
	body = append(body, selfNTLenPrefixed(nil)...)
	body = append(body, 0x00, 0x07)
	body = append(body, selfNTFrameMagic5...)
	body = append(body, selfNTLenPrefixed(hdr)...)
	plain := selfNTLenPrefixed(body)
	plain = append(plain, selfNTLenPrefixed(payload)...)
	sealed := selfNTTEAEncryptKey(selfNTTEAWrap(plain), d2k[:16])
	outer := make([]byte, 0, 16+len(d2)+len(sealed))
	outer = append(outer, selfNTBE32Bytes(12)...)
	outer = append(outer, 0x01)
	outer = append(outer, selfNTLenPrefixed(d2)...)
	outer = append(outer, 0x00)
	outer = append(outer, selfNTLenPrefixed([]byte(strconv.FormatUint(uin, 10)))...)
	outer = append(outer, sealed...)
	return selfNTLenPrefixed(outer), nil
}

func selfNTHex32(b []byte) bool {
	if len(b) != 32 {
		return false
	}
	for _, c := range b {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// selfNTParseFirstAuthResponse 对应 sub_1400568C0：优先取 field4 的 32 位 hex sig。
func selfNTParseFirstAuthResponse(plain []byte) (string, error) {
	if payload, err := selfNTExtractPayload(plain); err == nil && len(payload) > 0 {
		if b, _ := selfNTKexBytesField(payload, 4); selfNTHex32(b) {
			return strings.ToUpper(string(b)), nil
		}
		if msg, _ := selfNTKexBytesField(payload, 2); len(msg) > 0 {
			return "", fmt.Errorf("FIRST_AUTH 返回错误: %s", msg)
		}
	}
	for i := 0; i+32 <= len(plain); i++ {
		if selfNTHex32(plain[i : i+32]) {
			return strings.ToUpper(string(plain[i : i+32])), nil
		}
	}
	return "", errors.New("FIRST_AUTH 回包中未找到 32 位 hex sig")
}

type selfNTLoginResult struct {
	Sig       string
	Token0106 []byte
}

// selfNTLogin 用 NT-only 票据完成 KEX、EasyLogin 和 FirstAuth，返回 sig。
func selfNTLogin(uin uint64, machineID, token0106 []byte) (*selfNTLoginResult, error) {
	ts := time.Now().Unix()
	conn, _, err := selfNTConnectSSO(5 * time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	kexFrame, kp, _, err := selfNTBuildKexFrame(uin, machineID, ts)
	if err != nil {
		return nil, err
	}
	_, plain, err := selfNTCall(conn, kexFrame, 12*time.Second)
	if err != nil {
		return nil, err
	}
	kex, err := selfNTParseKexResponse(plain, kp)
	if err != nil {
		return nil, err
	}
	easyFrame, err := selfNTEasyLoginFrame(uin, machineID, token0106, kex.TicketA, kex.TicketB, ts)
	if err != nil {
		return nil, err
	}
	_, plain, err = selfNTCall(conn, easyFrame, 15*time.Second)
	if err != nil {
		return nil, err
	}
	sess, err := selfNTParseEasyLoginResponse(plain, kex.TicketA)
	if err != nil {
		return nil, err
	}

	guid := selfNTComputerID(uin, machineID)
	connSeq := strconv.FormatInt(ts+1, 10)
	uid := strconv.FormatUint(uin, 10)
	resultToken := sess.Token0106
	if len(resultToken) == 0 {
		resultToken = token0106
	}
	var last error
	for _, appID := range selfNTFirstAuthAppIDs {
		frame, err := selfNTFirstAuthFrame(uin, []byte(hex.EncodeToString(guid)), connSeq, uid, sess.TGT, sess.D2, sess.D2K, appID)
		if err != nil {
			return nil, err
		}
		_, plain, err = selfNTCallKey(conn, frame, sess.D2K, 15*time.Second)
		if err != nil {
			last = err
			continue
		}
		sig, err := selfNTParseFirstAuthResponse(plain)
		if err == nil {
			return &selfNTLoginResult{Sig: sig, Token0106: resultToken}, nil
		}
		last = err
	}
	return nil, fmt.Errorf("FIRST_AUTH 全部 appid 均未取得 sig: %w", last)
}
