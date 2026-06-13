// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"bytes"
	"testing"
)

// knownScriptNums are canonical encode/decode vectors matching the Bitcoin/HSD
// CScriptNum serialisation.
var knownScriptNums = []struct {
	value ScriptNum
	enc   []byte
}{
	{0, []byte{}},
	{1, []byte{0x01}},
	{-1, []byte{0x81}},
	{127, []byte{0x7f}},
	{128, []byte{0x80, 0x00}},
	{-128, []byte{0x80, 0x80}},
	{255, []byte{0xff, 0x00}},
	{256, []byte{0x00, 0x01}},
	{0x1234, []byte{0x34, 0x12}},
	{-0x1234, []byte{0x34, 0x92}},
	{0x7fffffff, []byte{0xff, 0xff, 0xff, 0x7f}},
}

func TestScriptNum_Encode_Good(t *testing.T) {
	for _, c := range knownScriptNums {
		got := c.value.Encode()
		if !bytes.Equal(got, c.enc) {
			t.Fatalf("Encode(%d) = %x, want %x", c.value, got, c.enc)
		}
	}
}

func TestScriptNum_Decode_Good(t *testing.T) {
	for _, c := range knownScriptNums {
		got, err := DecodeScriptNum(c.enc, true, 0)
		if err != nil {
			t.Fatalf("DecodeScriptNum(%x): %v", c.enc, err)
		}
		if got != c.value {
			t.Fatalf("DecodeScriptNum(%x) = %d, want %d", c.enc, got, c.value)
		}
	}
}

func TestScriptNum_Decode_Bad(t *testing.T) {
	// Non-minimal encodings are rejected when minimal is required.
	for _, enc := range [][]byte{
		{0x00},       // redundant zero
		{0x01, 0x00}, // trailing zero with high bit clear
		{0x80},       // negative zero
	} {
		if IsMinimalScriptNum(enc) {
			t.Fatalf("%x should be non-minimal", enc)
		}
		if _, err := DecodeScriptNum(enc, true, 4); err == nil {
			t.Fatalf("expected non-minimal rejection for %x", enc)
		}
	}
}

func TestScriptNum_Decode_Ugly(t *testing.T) {
	// Exceeding the length limit is an overflow error even for a valid number.
	enc := ScriptNum(0x1_0000_0000).Encode() // 5 bytes
	if _, err := DecodeScriptNum(enc, false, 4); err == nil {
		t.Fatalf("expected overflow rejection for 5-byte number under 4-byte limit")
	}
	// But it decodes fine without the limit.
	got, err := DecodeScriptNum(enc, false, 0)
	if err != nil {
		t.Fatalf("DecodeScriptNum no-limit: %v", err)
	}
	if got != 0x1_0000_0000 {
		t.Fatalf("got %d", got)
	}
}

func TestScriptNum_RoundTrip_Good(t *testing.T) {
	for _, v := range []ScriptNum{0, 1, -1, 42, -42, 1000, -1000, 0x7fffffff, -0x7fffffff} {
		enc := v.Encode()
		back, err := DecodeScriptNum(enc, false, 0)
		if err != nil {
			t.Fatalf("decode %d: %v", v, err)
		}
		if back != v {
			t.Fatalf("round-trip %d -> %x -> %d", v, enc, back)
		}
	}
}

func TestScriptNum_CastToBool_Good(t *testing.T) {
	if castToBool([]byte{}) {
		t.Fatal("empty must be false")
	}
	if castToBool([]byte{0x00, 0x00}) {
		t.Fatal("all-zero must be false")
	}
	if castToBool([]byte{0x80}) {
		t.Fatal("negative zero must be false")
	}
	if !castToBool([]byte{0x01}) {
		t.Fatal("0x01 must be true")
	}
	if !castToBool([]byte{0x00, 0x01}) {
		t.Fatal("non-zero high byte must be true")
	}
	// A leading non-zero byte makes the value true even if the top byte is 0x80
	// (this is not negative zero — there is magnitude below it).
	if !castToBool([]byte{0x01, 0x80}) {
		t.Fatal("value with magnitude and 0x80 top byte must be true")
	}
}

func TestScriptNum_Int_Ugly(t *testing.T) {
	// Saturation to int32 bounds.
	if ScriptNum(1<<40).Int() != (1<<31)-1 {
		t.Fatal("large positive should saturate to int32 max")
	}
	if ScriptNum(-(1 << 40)).Int() != -(1 << 31) {
		t.Fatal("large negative should saturate to int32 min")
	}
}
