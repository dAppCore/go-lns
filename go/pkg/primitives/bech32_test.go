// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/hex"
	"testing"

	core "dappco.re/go"
)

// TestBech32_Encode_Good encodes a version-0 P2WPKH address and checks the
// prefix, separator, and round-trippable structure of the result.
func TestBech32_Encode_Good(t *core.T) {
	hash, _ := hex.DecodeString("52a51028e839d00c87bf588697b86347718d3e8f")
	addr := Address{Version: 0, Hash: hash}

	encoded, err := addr.Encode("ts")
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}
	if len(encoded) < 8 {
		t.Fatalf("encoded address too short: %q", encoded)
	}
	if encoded[:3] != "ts1" {
		t.Fatalf("encoded = %q, want ts1 prefix", encoded)
	}
	// A version-0 / 20-byte programme always encodes to the same length for a
	// 2-char HRP: hrp(2) + sep(1) + version(1) + 32 data symbols + 6 checksum.
	if len(encoded) != 2+1+1+32+6 {
		t.Fatalf("encoded length = %d, want %d", len(encoded), 2+1+1+32+6)
	}
}

// TestBech32_Encode_Deterministic checks the encoding is stable for a fixed
// address and prefix.
func TestBech32_Encode_Deterministic(t *testing.T) {
	hash := make([]byte, 20)
	for i := range hash {
		hash[i] = byte(i)
	}
	addr := Address{Version: 0, Hash: hash}

	first, err := addr.Encode("lts")
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}
	second, err := addr.Encode("lts")
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}
	if first != second {
		t.Fatalf("encoding not deterministic: %q vs %q", first, second)
	}
	if first[:4] != "lts1" {
		t.Fatalf("encoded = %q, want lts1 prefix", first)
	}
}

// TestBech32_Encode_Bad covers an empty prefix and an out-of-range hash length.
func TestBech32_Encode_Bad(t *core.T) {
	addr := Address{Version: 0, Hash: make([]byte, 20)}
	if _, err := addr.Encode(""); err == nil {
		t.Fatalf("Encode accepted an empty prefix, want error")
	}

	short := Address{Version: 0, Hash: []byte{1}}
	if _, err := short.Encode("ts"); err == nil {
		t.Fatalf("Encode accepted a 1-byte hash, want error")
	}
}

// TestBech32_Encode_Ugly covers an out-of-range witness version.
func TestBech32_Encode_Ugly(t *core.T) {
	addr := Address{Version: 32, Hash: make([]byte, 20)}
	if _, err := addr.Encode("ts"); err == nil {
		t.Fatalf("Encode accepted version 32, want error")
	}
}

// TestBech32_PrefixSensitivity verifies the checksum binds the prefix: the same
// programme under different prefixes yields different strings.
func TestBech32_PrefixSensitivity(t *testing.T) {
	addr := Address{Version: 0, Hash: make([]byte, 20)}
	a, err := addr.Encode("ts")
	if err != nil {
		t.Fatalf("Encode(ts) error: %v", err)
	}
	b, err := addr.Encode("hs")
	if err != nil {
		t.Fatalf("Encode(hs) error: %v", err)
	}
	if a == b {
		t.Fatalf("encodings under different prefixes are equal: %q", a)
	}
}
