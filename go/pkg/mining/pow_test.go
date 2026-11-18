// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"testing"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// baseHeader builds a header with non-zero roots and a chosen difficulty, ready
// to mine.
func baseHeader(bits uint32) primitives.BlockHeader {
	return primitives.BlockHeader{
		Version:      1,
		PrevBlock:    hashOf(0xaa),
		MerkleRoot:   hashOf(0xbb),
		WitnessRoot:  hashOf(0xcc),
		TreeRoot:     hashOf(0xdd),
		ReservedRoot: primitives.Hash{},
		Time:         1_700_000_000,
		Bits:         bits,
	}
}

// TestPow_Solve_Good mines at a trivial difficulty and confirms the found nonce
// verifies, the whole round-trip closing through VerifyPOW.
func TestPow_Solve_Good(t *testing.T) {
	hdr := baseHeader(easyBits)
	solved := mustSolve(t, hdr, 1<<16)

	ok, err := VerifyPOW(solved)
	if err != nil {
		t.Fatalf("VerifyPOW: %v", err)
	}
	if !ok {
		t.Fatal("solved header did not verify against its own target")
	}
}

// TestPow_Solve_Bad rejects a header whose PoW hash does not meet a tight
// target. Nonce 0 against a value-1 target will (overwhelmingly) not solve.
func TestPow_Solve_Bad(t *testing.T) {
	hdr := baseHeader(hardBits)
	hdr.Nonce = 0

	ok, err := VerifyPOW(hdr)
	if err != nil {
		t.Fatalf("VerifyPOW: %v", err)
	}
	if ok {
		t.Fatal("an unmined header met a tight target")
	}

	// A bounded search at this difficulty must also fail.
	if _, solved := Solve(hdr, TargetFromBits(hdr.Bits), 256); solved {
		t.Fatal("Solve claimed success against a value-1 target within 256 rounds")
	}
}

// TestPow_Solve_Ugly confirms PowHash is deterministic and that Solve does not
// mutate the caller's header.
func TestPow_Solve_Ugly(t *testing.T) {
	hdr := baseHeader(easyBits)

	first, err := PowHash(hdr)
	if err != nil {
		t.Fatalf("PowHash: %v", err)
	}
	again, err := PowHash(hdr)
	if err != nil {
		t.Fatalf("PowHash again: %v", err)
	}
	if first != again {
		t.Fatal("PowHash is not deterministic")
	}

	before := hdr.Nonce
	_, _ = Solve(hdr, TargetFromBits(hdr.Bits), 4)
	if hdr.Nonce != before {
		t.Fatalf("Solve mutated the caller's header nonce: %d -> %d", before, hdr.Nonce)
	}
}

// TestPow_MaskChangesHash confirms the mask is folded into the PoW hash. The
// mask influences the result by two paths in the HSD scheme: it feeds the
// commitment hash (via maskHash, changing the share hash) and it is XORed over
// the final share hash. Either way, a non-zero mask must change the PoW hash.
func TestPow_MaskChangesHash(t *testing.T) {
	hdr := baseHeader(easyBits)

	zeroMaskHash, err := PowHash(hdr)
	if err != nil {
		t.Fatalf("PowHash zero mask: %v", err)
	}

	hdr.Mask = make([]byte, maskSize)
	hdr.Mask[0] = 0xff
	maskedHash, err := PowHash(hdr)
	if err != nil {
		t.Fatalf("PowHash with mask: %v", err)
	}

	if zeroMaskHash == maskedHash {
		t.Fatal("mask did not affect the PoW hash")
	}
}

// TestPow_MarshalHeader_Good confirms the header serialises to exactly
// HeaderSize bytes with variable fields left nil (padded internally).
func TestPow_MarshalHeader_Good(t *testing.T) {
	hdr := baseHeader(easyBits)
	raw, err := MarshalHeader(hdr)
	if err != nil {
		t.Fatalf("MarshalHeader: %v", err)
	}
	if len(raw) != HeaderSize {
		t.Fatalf("header length = %d, want %d", len(raw), HeaderSize)
	}
}
