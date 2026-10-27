package main

// 用 Windows CNG 复现样本 sub_1400B3960 的流程，验证 selfNTSharedSecret 逐字节一致：
//   BCryptImportKeyPair(ECK2 私钥) + BCryptImportKeyPair(ECK1 对端公钥)
//   -> BCryptSecretAgreement -> BCryptDeriveKey("TRUNCATE", 32B) -> 反转

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"runtime"
	"syscall"
	"testing"
	"unsafe"
)

var (
	tBcrypt        = syscall.NewLazyDLL("bcrypt.dll")
	tOpenAlg       = tBcrypt.NewProc("BCryptOpenAlgorithmProvider")
	tGenKeyPair    = tBcrypt.NewProc("BCryptGenerateKeyPair")
	tFinalizeKey   = tBcrypt.NewProc("BCryptFinalizeKeyPair")
	tImportKeyPair = tBcrypt.NewProc("BCryptImportKeyPair")
	tSecretAgree   = tBcrypt.NewProc("BCryptSecretAgreement")
	tDeriveKey     = tBcrypt.NewProc("BCryptDeriveKey")
	tDestroyKey    = tBcrypt.NewProc("BCryptDestroyKey")
	tDestroySecret = tBcrypt.NewProc("BCryptDestroySecret")
	tCloseAlg      = tBcrypt.NewProc("BCryptCloseAlgorithmProvider")
)

func tWide(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}

func cngRawSecret(priv *ecdh.PrivateKey, peerPoint []byte) ([]byte, error) {
	var hAlg, hPriv, hPeer, hSecret syscall.Handle
	if r, _, _ := tOpenAlg.Call(uintptr(unsafe.Pointer(&hAlg)), uintptr(unsafe.Pointer(tWide("ECDH_P256"))), 0, 0); int32(r) < 0 {
		return nil, syscall.Errno(uintptr(r))
	}
	defer tCloseAlg.Call(uintptr(hAlg))

	pub := priv.PublicKey().Bytes()
	privBlob := append([]byte{'E', 'C', 'K', '2', 32, 0, 0, 0}, pub[1:]...)
	privBlob = append(privBlob, priv.Bytes()...)
	if r, _, _ := tImportKeyPair.Call(uintptr(hAlg), 0, uintptr(unsafe.Pointer(tWide("ECCPRIVATEBLOB"))),
		uintptr(unsafe.Pointer(&hPriv)), uintptr(unsafe.Pointer(&privBlob[0])), uintptr(len(privBlob)), 0); int32(r) < 0 {
		return nil, syscall.Errno(uintptr(r))
	}
	defer tDestroyKey.Call(uintptr(hPriv))

	peerBlob := append([]byte{'E', 'C', 'K', '1', 32, 0, 0, 0}, peerPoint[1:]...)
	if r, _, _ := tImportKeyPair.Call(uintptr(hAlg), 0, uintptr(unsafe.Pointer(tWide("ECCPUBLICBLOB"))),
		uintptr(unsafe.Pointer(&hPeer)), uintptr(unsafe.Pointer(&peerBlob[0])), uintptr(len(peerBlob)), 0); int32(r) < 0 {
		return nil, syscall.Errno(uintptr(r))
	}
	defer tDestroyKey.Call(uintptr(hPeer))

	if r, _, _ := tSecretAgree.Call(uintptr(hPriv), uintptr(hPeer), uintptr(unsafe.Pointer(&hSecret)), 0); int32(r) < 0 {
		return nil, syscall.Errno(uintptr(r))
	}
	defer tDestroySecret.Call(uintptr(hSecret))

	out := make([]byte, 32)
	cb := uint32(32)
	if r, _, _ := tDeriveKey.Call(uintptr(hSecret), uintptr(unsafe.Pointer(tWide("TRUNCATE"))), 0,
		uintptr(unsafe.Pointer(&out[0])), 32, uintptr(unsafe.Pointer(&cb)), 0); int32(r) < 0 {
		return nil, syscall.Errno(uintptr(r))
	}
	return out[:cb], nil
}

// TestSelfNTSharedSecretMatchesCNG 需要能加载 bcrypt.dll（Windows）。
func TestSelfNTSharedSecretMatchesCNG(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("需要 Windows CNG")
	}
	for i := 0; i < 3; i++ {
		ka, err := selfNTNewKeyPair()
		if err != nil {
			t.Fatal(err)
		}
		kb, err := ecdh.P256().GenerateKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		got, err := ka.SharedSecret(kb.PublicKey().Bytes())
		if err != nil {
			t.Fatal(err)
		}
		raw, err := cngRawSecret(ka.priv, kb.PublicKey().Bytes())
		if err != nil {
			t.Fatalf("CNG: %v", err)
		}
		want := make([]byte, len(raw))
		for j := range raw {
			want[j] = raw[len(raw)-1-j]
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("round %d 不一致\n got=%x\nwant=%x", i, got, want)
		}
	}
}
