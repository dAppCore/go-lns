// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/binary"
	"testing"

	"golang.org/x/crypto/blake2b"
)

func TestClaimHelpers(t *testing.T) {
	var claim Claim
	blob := []byte("ownership-proof")

	claim.FromBlob(blob)

	if got := claim.GetBlob(); string(got) != string(blob) {
		t.Fatalf("GetBlob() = %q, want %q", got, blob)
	}

	wantHash := blake2b.Sum256(blob)
	if got := claim.Hash(); got != wantHash {
		t.Fatalf("Hash() = %x, want %x", got, wantHash)
	}

	if got := claim.GetHash(); got != wantHash {
		t.Fatalf("GetHash() = %x, want %x", got, wantHash)
	}

	if got := claim.HashHex(); got != claim.GetHashHex() {
		t.Fatalf("HashHex() and GetHashHex() should agree, got %q and %q", got, claim.GetHashHex())
	}

	if got := claim.GetSize(); got != 2+len(blob) {
		t.Fatalf("GetSize() = %d, want %d", got, 2+len(blob))
	}

	if got := claim.Weight(); got != claim.GetSize() {
		t.Fatalf("Weight() = %d, want %d", got, claim.GetSize())
	}

	if got := claim.GetWeight(); got != claim.Weight() {
		t.Fatalf("GetWeight() = %d, want %d", got, claim.Weight())
	}

	if got := claim.VirtualSize(); got != (claim.Weight()+3)/4 {
		t.Fatalf("VirtualSize() = %d, want %d", got, (claim.Weight()+3)/4)
	}

	if got := claim.GetVirtualSize(); got != claim.VirtualSize() {
		t.Fatalf("GetVirtualSize() = %d, want %d", got, claim.VirtualSize())
	}

	if got := claim.ToBlob(); string(got) != string(blob) {
		t.Fatalf("ToBlob() = %q, want %q", got, blob)
	}

	if got := claim.GetToBlob(); string(got) != string(blob) {
		t.Fatalf("GetToBlob() = %q, want %q", got, blob)
	}

	raw, err := claim.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	if size := binary.LittleEndian.Uint16(raw[:2]); int(size) != len(blob) {
		t.Fatalf("MarshalBinary length prefix = %d, want %d", size, len(blob))
	}

	var decoded Claim
	if err := decoded.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if got := decoded.GetBlob(); string(got) != string(blob) {
		t.Fatalf("decoded blob = %q, want %q", got, blob)
	}

	clone := claim.Clone()
	if got := clone.GetBlob(); string(got) != string(blob) {
		t.Fatalf("Clone() = %q, want %q", got, blob)
	}

	if NewClaim() == nil || GetNewClaim() == nil {
		t.Fatal("NewClaim helpers should return a claim wrapper")
	}
}
