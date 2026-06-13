// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"bytes"
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

func TestSighash_Deterministic_Good(t *testing.T) {
	_, pub := testKey(t, 10)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	redeem := NewScript().FromPubkeyhash(pkh)

	h1, err := SignatureHash(tx, 0, redeem, 2000, SighashAll)
	if err != nil {
		t.Fatalf("SignatureHash: %v", err)
	}
	h2, err := SignatureHash(tx, 0, redeem, 2000, SighashAll)
	if err != nil {
		t.Fatalf("SignatureHash: %v", err)
	}
	if !bytes.Equal(h1, h2) {
		t.Fatal("signature hash must be deterministic")
	}
	if len(h1) != 32 {
		t.Fatalf("signature hash length = %d, want 32", len(h1))
	}
}

func TestSighash_TypeSensitivity_Good(t *testing.T) {
	_, pub := testKey(t, 11)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	redeem := NewScript().FromPubkeyhash(pkh)

	all, _ := SignatureHash(tx, 0, redeem, 2000, SighashAll)
	none, _ := SignatureHash(tx, 0, redeem, 2000, SighashNone)
	single, _ := SignatureHash(tx, 0, redeem, 2000, SighashSingle)

	if bytes.Equal(all, none) || bytes.Equal(all, single) || bytes.Equal(none, single) {
		t.Fatal("different sighash types must produce different digests")
	}
}

func TestSighash_ValueSensitivity_Good(t *testing.T) {
	_, pub := testKey(t, 12)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	redeem := NewScript().FromPubkeyhash(pkh)

	a, _ := SignatureHash(tx, 0, redeem, 2000, SighashAll)
	b, _ := SignatureHash(tx, 0, redeem, 2001, SighashAll)
	if bytes.Equal(a, b) {
		t.Fatal("changing the spent value must change the digest")
	}
}

func TestSighash_SignChecksig_Good(t *testing.T) {
	priv, pub := testKey(t, 13)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	redeem := NewScript().FromPubkeyhash(pkh)
	value := uint64(2000)

	sig, err := Sign(tx, 0, redeem, value, priv, SighashAll)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if len(sig) != sigLen {
		t.Fatalf("sig length = %d, want %d", len(sig), sigLen)
	}
	ok, err := Checksig(tx, 0, redeem, value, sig, pub)
	if err != nil {
		t.Fatalf("Checksig: %v", err)
	}
	if !ok {
		t.Fatal("Checksig should accept a freshly produced signature")
	}
}

func TestSighash_Checksig_Bad(t *testing.T) {
	priv, pub := testKey(t, 14)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	redeem := NewScript().FromPubkeyhash(pkh)
	value := uint64(2000)
	sig, _ := Sign(tx, 0, redeem, value, priv, SighashAll)

	// Verify against a different key -> false.
	_, otherPub := testKey(t, 77)
	ok, err := Checksig(tx, 0, redeem, value, sig, otherPub)
	if err != nil {
		t.Fatalf("Checksig: %v", err)
	}
	if ok {
		t.Fatal("Checksig must reject a signature under the wrong key")
	}

	// Empty signature -> false, no error.
	ok, err = Checksig(tx, 0, redeem, value, nil, pub)
	if err != nil {
		t.Fatalf("Checksig empty: %v", err)
	}
	if ok {
		t.Fatal("empty signature must not verify")
	}
}

func TestSighash_Ugly(t *testing.T) {
	_, pub := testKey(t, 15)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	redeem := NewScript().FromPubkeyhash(pkh)

	// Out-of-range input index.
	if _, err := SignatureHash(tx, 5, redeem, 2000, SighashAll); err == nil {
		t.Fatal("expected error for out-of-range index")
	}
	// Nil previous script.
	if _, err := SignatureHash(tx, 0, nil, 2000, SighashAll); err == nil {
		t.Fatal("expected error for nil prev script")
	}
}

func TestSighash_DecodeCompactSig_Bad(t *testing.T) {
	// Zero R/S is rejected.
	if _, err := decodeCompactSig(make([]byte, 64)); err == nil {
		t.Fatal("zero signature should be rejected")
	}
	// Wrong length is rejected.
	if _, err := decodeCompactSig(make([]byte, 10)); err == nil {
		t.Fatal("short signature should be rejected")
	}
}

func TestSighash_EncodeOutput_Good(t *testing.T) {
	out := primitives.Output{
		Value:   0x1122334455667788,
		Address: primitives.Address{Version: 0, Hash: bytes.Repeat([]byte{0xaa}, 20)},
	}
	enc, err := encodeOutput(out)
	if err != nil {
		t.Fatalf("encodeOutput: %v", err)
	}
	// value(8) + version(1) + len(1) + hash(20) + covenant(type 1 + count varint 1).
	want := 8 + 1 + 1 + 20 + 2
	if len(enc) != want {
		t.Fatalf("encoded output length = %d, want %d", len(enc), want)
	}
	// First 8 bytes are the little-endian value.
	if enc[0] != 0x88 || enc[7] != 0x11 {
		t.Fatalf("value not little-endian: %x", enc[:8])
	}
	// version byte then hash length.
	if enc[8] != 0x00 || enc[9] != 20 {
		t.Fatalf("address header wrong: %x", enc[8:10])
	}
}
