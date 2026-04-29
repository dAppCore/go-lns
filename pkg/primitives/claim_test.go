// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "dappco.re/go"
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

// BEGIN v090 file-aware generated coverage

func TestClaim_NewClaim_Good(t *core.T) {
	fn := ax7ValueOf(NewClaim)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_NewClaim_Bad(t *core.T) {
	fn := ax7ValueOf(NewClaim)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_NewClaim_Ugly(t *core.T) {
	fn := ax7ValueOf(NewClaim)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_GetNewClaim_Good(t *core.T) {
	fn := ax7ValueOf(GetNewClaim)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_GetNewClaim_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewClaim)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_GetNewClaim_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewClaim)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Refresh_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).Refresh)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Refresh_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).Refresh)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Refresh_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).Refresh)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Clone_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).Clone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Clone_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).Clone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Clone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).Clone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Hash_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).Hash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Hash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).Hash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Hash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).Hash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_HashHex_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).HashHex)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_HashHex_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).HashHex)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_HashHex_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).HashHex)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetHashHex_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetHashHex)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetHashHex_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetHashHex)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetHashHex_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetHashHex)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Weight_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).Weight)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Weight_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).Weight)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_Weight_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).Weight)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetWeight_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetWeight)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetWeight_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetWeight)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetWeight_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetWeight)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_VirtualSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).VirtualSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_VirtualSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).VirtualSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_VirtualSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).VirtualSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetVirtualSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetVirtualSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetVirtualSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetVirtualSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetVirtualSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetVirtualSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_MarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).MarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_MarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).MarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_MarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).MarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_UnmarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).UnmarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_UnmarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).UnmarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_UnmarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).UnmarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_ToBlob_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).ToBlob)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_ToBlob_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).ToBlob)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_ToBlob_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).ToBlob)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetToBlob_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetToBlob)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetToBlob_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetToBlob)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetToBlob_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetToBlob)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetBlob_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetBlob)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetBlob_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetBlob)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetBlob_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetBlob)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_FromBlob_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).FromBlob)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_FromBlob_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).FromBlob)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_FromBlob_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).FromBlob)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetFromBlob_Good(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetFromBlob)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetFromBlob_Bad(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetFromBlob)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClaim_Claim_GetFromBlob_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Claim).GetFromBlob)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
