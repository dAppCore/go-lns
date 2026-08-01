// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"bytes"
	"crypto/sha3"
	"testing"
)

func TestScript_FromPubkeyhash_Good(t *testing.T) {
	hash := bytes.Repeat([]byte{0xab}, 20)
	s, err := FromPubkeyhash(hash)
	if err != nil {
		t.Fatalf("FromPubkeyhash: %v", err)
	}
	want := []byte{OP_DUP, OP_BLAKE160, 0x14}
	want = append(want, hash...)
	want = append(want, OP_EQUALVERIFY, OP_CHECKSIG)
	if !bytes.Equal(s.Encode(), want) {
		t.Fatalf("P2PKH raw = %x\nwant %x", s.Encode(), want)
	}
	if !s.IsPubkeyhash() {
		t.Fatal("IsPubkeyhash should be true")
	}
	if !bytes.Equal(s.GetPubkeyhash(), hash) {
		t.Fatal("GetPubkeyhash mismatch")
	}
}

func TestScript_FromPubkeyhash_Bad(t *testing.T) {
	if _, err := FromPubkeyhash([]byte{0x01, 0x02}); err == nil {
		t.Fatal("expected error for short pubkey hash")
	}
}

func TestScript_FromMultisig_Good(t *testing.T) {
	keys := [][]byte{
		append([]byte{0x02}, bytes.Repeat([]byte{0x03}, 32)...),
		append([]byte{0x02}, bytes.Repeat([]byte{0x01}, 32)...),
		append([]byte{0x02}, bytes.Repeat([]byte{0x02}, 32)...),
	}
	s, err := FromMultisig(2, 3, keys)
	if err != nil {
		t.Fatalf("FromMultisig: %v", err)
	}
	if !s.IsMultisig() {
		t.Fatal("IsMultisig should be true")
	}
	code := s.Code()
	if code[0].Value != OP_2 || code[len(code)-2].Value != OP_3 {
		t.Fatalf("multisig m/n opcodes wrong: %s %s", code[0].Symbol(), code[len(code)-2].Symbol())
	}
	if code[len(code)-1].Value != OP_CHECKMULTISIG {
		t.Fatal("multisig must end with OP_CHECKMULTISIG")
	}
	// Keys must be sorted lexicographically: 0x01.., 0x02.., 0x03..
	if code[1].Data[1] != 0x01 || code[2].Data[1] != 0x02 || code[3].Data[1] != 0x03 {
		t.Fatalf("keys not sorted: %x %x %x", code[1].Data[1:2], code[2].Data[1:2], code[3].Data[1:2])
	}
}

func TestScript_FromMultisig_Bad(t *testing.T) {
	key := append([]byte{0x02}, bytes.Repeat([]byte{0x01}, 32)...)
	// m > n.
	if _, err := FromMultisig(3, 2, [][]byte{key, key}); err == nil {
		t.Fatal("expected error for m > n")
	}
	// key count mismatch.
	if _, err := FromMultisig(1, 2, [][]byte{key}); err == nil {
		t.Fatal("expected error for key count != n")
	}
	// bad key length.
	if _, err := FromMultisig(1, 1, [][]byte{{0x02, 0x03}}); err == nil {
		t.Fatal("expected error for short key")
	}
}

func TestScript_Sha3_Good(t *testing.T) {
	s, _ := FromPubkey(append([]byte{0x02}, bytes.Repeat([]byte{0x05}, 32)...))
	want := sha3.Sum256(s.Encode())
	if !bytes.Equal(s.Sha3(), want[:]) {
		t.Fatalf("Sha3 mismatch")
	}
}

func TestScript_Decode_Bad(t *testing.T) {
	// A script larger than MaxScriptSize is rejected.
	big := make([]byte, MaxScriptSize+1)
	if _, err := DecodeScript(big); err == nil {
		t.Fatal("expected error decoding oversized script")
	}
}

func TestScript_Decode_Ugly(t *testing.T) {
	// A truncated push inside the script decodes to a parse-error opcode rather
	// than failing, and the raw bytes round-trip exactly.
	raw := []byte{OP_DUP, OP_PUSHDATA1, 0x05, 0xaa} // truncated push
	s, err := DecodeScript(raw)
	if err != nil {
		t.Fatalf("DecodeScript: %v", err)
	}
	if !bytes.Equal(s.Encode(), raw) {
		t.Fatalf("raw not preserved: %x", s.Encode())
	}
	last := s.Code()[len(s.Code())-1]
	if !last.IsParseError() {
		t.Fatalf("expected trailing parse error, got %s", last.Symbol())
	}
}

func TestScript_GetSubscript_Good(t *testing.T) {
	s, _ := FromOps([]Opcode{
		OpFromOp(OP_DUP),
		OpFromOp(OP_CODESEPARATOR),
		OpFromOp(OP_CHECKSIG),
	})
	sub, err := s.GetSubscript(0)
	if err != nil {
		t.Fatalf("GetSubscript: %v", err)
	}
	for _, op := range sub.Code() {
		if op.Value == OP_CODESEPARATOR {
			t.Fatal("subscript must not contain OP_CODESEPARATOR")
		}
	}
	if len(sub.Code()) != 2 {
		t.Fatalf("subscript len = %d, want 2", len(sub.Code()))
	}
}

func TestScript_GetVarSize_Good(t *testing.T) {
	s, _ := FromOps([]Opcode{OpFromOp(OP_1)})
	if s.GetVarSize() != 1+s.GetSize() {
		t.Fatalf("GetVarSize = %d, want %d", s.GetVarSize(), 1+s.GetSize())
	}
}
