// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"bytes"
	"testing"
)

func TestOpcode_OpFromData_Good(t *testing.T) {
	// Small integers collapse to their dedicated opcodes.
	cases := []struct {
		in   []byte
		want int
	}{
		{[]byte{}, OP_0},
		{[]byte{0x01}, OP_1},
		{[]byte{0x10}, OP_16},
		{[]byte{0x81}, OP_1NEGATE},
	}
	for _, c := range cases {
		op := OpFromData(c.in)
		if op.Value != c.want {
			t.Fatalf("OpFromData(%x).Value = %#x, want %#x", c.in, op.Value, c.want)
		}
		if op.Data != nil {
			t.Fatalf("OpFromData(%x) should be a bare opcode, got data %x", c.in, op.Data)
		}
	}
}

func TestOpcode_OpFromData_Bad(t *testing.T) {
	// A two-byte value cannot collapse and must be a direct push of length 2.
	op := OpFromData([]byte{0x17, 0x00})
	if op.Value != 2 {
		t.Fatalf("OpFromData(2 bytes).Value = %d, want 2 (direct push)", op.Value)
	}
	if !bytes.Equal(op.Data, []byte{0x17, 0x00}) {
		t.Fatalf("OpFromData data = %x", op.Data)
	}
}

func TestOpcode_OpFromData_Ugly(t *testing.T) {
	// 0x4b bytes is the largest direct push; 0x4c forces OP_PUSHDATA1.
	op := OpFromPush(bytes.Repeat([]byte{0xaa}, 0x4c))
	if op.Value != OP_PUSHDATA1 {
		t.Fatalf("76-byte push opcode = %#x, want OP_PUSHDATA1", op.Value)
	}
	big := OpFromPush(bytes.Repeat([]byte{0xbb}, 0x100))
	if big.Value != OP_PUSHDATA2 {
		t.Fatalf("256-byte push opcode = %#x, want OP_PUSHDATA2", big.Value)
	}
}

func TestOpcode_RoundTrip_Good(t *testing.T) {
	// Encode/decode round-trip for representative push and non-push opcodes.
	ops := []Opcode{
		OpFromOp(OP_DUP),
		OpFromPush([]byte{0x01, 0x02, 0x03}),
		OpFromPush(bytes.Repeat([]byte{0x07}, 200)),   // OP_PUSHDATA1
		OpFromPush(bytes.Repeat([]byte{0x09}, 70000)), // OP_PUSHDATA4 (> 0xffff)
		OpFromOp(OP_CHECKSIG),
	}
	for _, op := range ops {
		raw, err := op.appendTo(nil)
		if err != nil {
			t.Fatalf("appendTo(%s): %v", op.Symbol(), err)
		}
		got, n := readOpcode(raw)
		if n != len(raw) {
			t.Fatalf("readOpcode consumed %d of %d for %s", n, len(raw), op.Symbol())
		}
		if got.Value != op.Value {
			t.Fatalf("round-trip value %#x != %#x", got.Value, op.Value)
		}
		if !bytes.Equal(got.Data, op.Data) {
			t.Fatalf("round-trip data mismatch for %s", op.Symbol())
		}
	}
}

func TestOpcode_ReadOpcode_Bad(t *testing.T) {
	// Truncated OP_PUSHDATA1 (claims 5 bytes, supplies 2) must yield a parse error
	// and consume the whole remaining buffer.
	raw := []byte{OP_PUSHDATA1, 0x05, 0xaa, 0xbb}
	op, n := readOpcode(raw)
	if !op.IsParseError() {
		t.Fatalf("expected parse error, got %s", op.Symbol())
	}
	if n != len(raw) {
		t.Fatalf("parse error should consume all %d bytes, consumed %d", len(raw), n)
	}
}

func TestOpcode_ReadOpcode_Ugly(t *testing.T) {
	// A direct-push opcode whose declared length exceeds the buffer (untrusted
	// length guard) must not panic and must report a parse error.
	raw := []byte{0x20, 0x01, 0x02} // claims 32 bytes, has 2
	op, n := readOpcode(raw)
	if !op.IsParseError() {
		t.Fatalf("expected parse error for over-long direct push, got %s", op.Symbol())
	}
	if n != len(raw) {
		t.Fatalf("consumed %d, want %d", n, len(raw))
	}

	// PUSHDATA4 with a 5GiB length must be rejected without allocating.
	huge := []byte{OP_PUSHDATA4, 0xff, 0xff, 0xff, 0xff, 0x00}
	hop, _ := readOpcode(huge)
	if !hop.IsParseError() {
		t.Fatalf("expected parse error for oversized PUSHDATA4, got %s", hop.Symbol())
	}
}

func TestOpcode_IsDisabled_Good(t *testing.T) {
	for _, v := range []int{OP_CAT, OP_SUBSTR, OP_MUL, OP_DIV, OP_LSHIFT, OP_2MUL} {
		if !OpFromOp(v).IsDisabled() {
			t.Fatalf("%s should be disabled", OpcodeName(v))
		}
	}
	for _, v := range []int{OP_DUP, OP_ADD, OP_CHECKSIG} {
		if OpFromOp(v).IsDisabled() {
			t.Fatalf("%s should not be disabled", OpcodeName(v))
		}
	}
}

func TestOpcode_IsMinimal_Bad(t *testing.T) {
	// A push of {0x01} via a direct push of length 1 is NOT minimal: it should be
	// OP_1. The minimal-data rule rejects it.
	nonMinimal := Opcode{Value: 1, Data: []byte{0x01}}
	if nonMinimal.IsMinimal() {
		t.Fatalf("direct push of small int 1 must be non-minimal")
	}
	// A genuine 1-byte push of 0x17 IS minimal.
	minimal := Opcode{Value: 1, Data: []byte{0x17}}
	if !minimal.IsMinimal() {
		t.Fatalf("direct push of 0x17 should be minimal")
	}
}

func TestOpcode_OpcodeName_Ugly(t *testing.T) {
	if got := OpcodeName(OP_SHA3); got != "OP_SHA3" {
		t.Fatalf("OP_SHA3 name = %q", got)
	}
	if got := OpcodeName(opParseError); got != "OP_INVALIDOPCODE" {
		t.Fatalf("parse-error name = %q", got)
	}
	if got := OpcodeName(0xee); got == "" {
		t.Fatalf("unknown opcode name should be non-empty")
	}
}

func TestOpcode_Helpers_Good(t *testing.T) {
	// Symbol for a bare opcode.
	if OpFromOp(OP_DUP).Symbol() != "OP_DUP" {
		t.Fatalf("Symbol = %q", OpFromOp(OP_DUP).Symbol())
	}
	// Size of a direct push (1 length byte already encoded as the value) is
	// 1 + len.
	push := OpFromPush([]byte{0x01, 0x02, 0x03})
	if push.Size() != 1+3 {
		t.Fatalf("direct push size = %d, want 4", push.Size())
	}
	// Size of a bare opcode is 1.
	if OpFromOp(OP_CHECKSIG).Size() != 1 {
		t.Fatalf("bare opcode size = %d, want 1", OpFromOp(OP_CHECKSIG).Size())
	}
	// PUSHDATA1 size includes the 1-byte length prefix.
	pd1 := OpFromPush(make([]byte, 100))
	if pd1.Value != OP_PUSHDATA1 || pd1.Size() != 2+100 {
		t.Fatalf("PUSHDATA1 size = %d, want 102", pd1.Size())
	}
}

func TestOpcode_ToPush_Good(t *testing.T) {
	// Small-integer opcodes report their canonical push bytes.
	if got := OpFromOp(OP_0).ToPush(); len(got) != 0 {
		t.Fatalf("OP_0 push = %x, want empty", got)
	}
	if got := OpFromOp(OP_1NEGATE).ToPush(); len(got) != 1 || got[0] != 0x81 {
		t.Fatalf("OP_1NEGATE push = %x, want 81", got)
	}
	if got := OpFromOp(OP_5).ToPush(); len(got) != 1 || got[0] != 0x05 {
		t.Fatalf("OP_5 push = %x, want 05", got)
	}
	// A pushdata opcode reports its data.
	push := OpFromPush([]byte{0xaa, 0xbb})
	if got := push.ToPush(); len(got) != 2 {
		t.Fatalf("push ToPush len = %d, want 2", len(got))
	}
	// ToSmall round-trip.
	if OpFromOp(OP_7).ToSmall() != 7 {
		t.Fatalf("OP_7 ToSmall = %d", OpFromOp(OP_7).ToSmall())
	}
	if OpFromOp(OP_DUP).ToSmall() != -1 {
		t.Fatalf("OP_DUP ToSmall = %d, want -1", OpFromOp(OP_DUP).ToSmall())
	}
}
