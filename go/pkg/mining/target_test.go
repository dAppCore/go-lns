// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"math/big"
	"testing"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// TestTarget_TargetFromBits_Good decodes well-formed compact targets and checks
// the bit width matches the encoded exponent.
func TestTarget_TargetFromBits_Good(t *testing.T) {
	cases := []struct {
		bits   uint32
		bitLen int
	}{
		{0x1d00ffff, 224}, // classic Bitcoin-genesis-style target
		{0x2100ffff, 256}, // very wide target
		{0x03000001, 1},   // tiny target (value 1)
	}
	for _, c := range cases {
		target := TargetFromBits(c.bits)
		if target.Sign() <= 0 {
			t.Fatalf("bits %#x: target is non-positive", c.bits)
		}
		if target.BitLen() != c.bitLen {
			t.Fatalf("bits %#x: bitLen = %d, want %d", c.bits, target.BitLen(), c.bitLen)
		}
	}
}

// TestTarget_TargetFromBits_Bad rejects the zero compact form (no valid PoW
// target).
func TestTarget_TargetFromBits_Bad(t *testing.T) {
	target := TargetFromBits(0)
	if target.Sign() != 0 {
		t.Fatalf("zero bits decoded to a non-zero target: %s", target)
	}
}

// TestTarget_TargetFromBits_Ugly decodes the negative-flag form and confirms the
// sign bit is honoured, so a malformed target cannot pass as positive.
func TestTarget_TargetFromBits_Ugly(t *testing.T) {
	// Set the sign flag (bit 23) on a small mantissa.
	bits := uint32(0x03_800001)
	target := TargetFromBits(bits)
	if target.Sign() >= 0 {
		t.Fatalf("sign-flagged bits %#x decoded to non-negative target %s", bits, target)
	}
}

// TestTarget_BitsFromTarget_Good round-trips compact targets through decode and
// re-encode.
func TestTarget_BitsFromTarget_Good(t *testing.T) {
	for _, bits := range []uint32{0x1d00ffff, 0x1b0404cb, 0x2100ffff} {
		target := TargetFromBits(bits)
		got := BitsFromTarget(target)
		if got != bits {
			t.Fatalf("round-trip bits %#x -> target -> %#x", bits, got)
		}
	}
}

// TestTarget_BitsFromTarget_Bad encodes zero to the zero compact form.
func TestTarget_BitsFromTarget_Bad(t *testing.T) {
	if got := BitsFromTarget(big.NewInt(0)); got != 0 {
		t.Fatalf("BitsFromTarget(0) = %#x, want 0", got)
	}
}

// TestTarget_BitsFromTarget_Ugly encodes a value whose top byte has bit 23 set,
// forcing the mantissa normalisation (shift down, bump exponent) so the sign
// flag is not accidentally set.
func TestTarget_BitsFromTarget_Ugly(t *testing.T) {
	// 0x800000 has the high mantissa bit set at exponent 3.
	target := big.NewInt(0x800000)
	bits := BitsFromTarget(target)
	// Decoding must recover the same value (normalisation is lossless here).
	if TargetFromBits(bits).Cmp(target) != 0 {
		t.Fatalf("normalised round-trip lost value: bits %#x decodes to %s, want %s",
			bits, TargetFromBits(bits), target)
	}
	// And the sign flag must be clear.
	if bits&0x800000 != 0 {
		t.Fatalf("normalisation set the sign flag: bits %#x", bits)
	}
}

// TestTarget_CheckTarget_Good accepts a hash below the target.
func TestTarget_CheckTarget_Good(t *testing.T) {
	// A near-max target accepts essentially any hash.
	var low primitives.Hash
	low[31] = 0x01
	if !CheckTarget(low, 0x2100ffff) {
		t.Fatal("a tiny hash failed a near-max target")
	}
}

// TestTarget_CheckTarget_Bad rejects a hash above the target.
func TestTarget_CheckTarget_Bad(t *testing.T) {
	// A maximal hash against a tiny (value 1) target must fail.
	var high primitives.Hash
	for i := range high {
		high[i] = 0xff
	}
	if CheckTarget(high, 0x03000001) {
		t.Fatal("a maximal hash passed a tiny target")
	}
}

// TestTarget_CheckTarget_Ugly rejects against a zero target, which can never be
// met.
func TestTarget_CheckTarget_Ugly(t *testing.T) {
	var zero primitives.Hash
	if CheckTarget(zero, 0) {
		t.Fatal("zero target was satisfiable")
	}
}
