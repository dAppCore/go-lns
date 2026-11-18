// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "github.com/dappcore/go"
	"testing"
)

func TestCovenantBinaryRoundTrip(t *testing.T) {
	original := Covenant{}
	original.SetFinalize(
		Hash{43, 44, 45},
		188,
		[]byte("final"),
		12,
		21,
		34,
		Hash{46, 47, 48},
	)

	encoded, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	if got := original.GetVarSize(); got != len(encoded) {
		t.Fatalf("GetVarSize() = %d, want encoded length %d", got, len(encoded))
	}

	var decoded Covenant
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if decoded.Type != original.Type {
		t.Fatalf("decoded type = %d, want %d", decoded.Type, original.Type)
	}

	if decoded.Len() != original.Len() {
		t.Fatalf("decoded item count = %d, want %d", decoded.Len(), original.Len())
	}

	for i := range original.Items {
		if string(decoded.Items[i]) != string(original.Items[i]) {
			t.Fatalf("decoded item %d = %x, want %x", i, decoded.Items[i], original.Items[i])
		}
	}
}

func TestCovenantBinaryZeroValue(t *testing.T) {
	var original Covenant

	encoded, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	if got := original.GetVarSize(); got != len(encoded) {
		t.Fatalf("GetVarSize() = %d, want encoded length %d", got, len(encoded))
	}

	var decoded Covenant
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if decoded.Type != 0 || decoded.Len() != 0 {
		t.Fatalf("decoded zero value = %+v, want zero value", decoded)
	}
}

// TestCovenantUnmarshalBinaryRejectsOversizedItemCount ensures the varint item
// count is bounded against the remaining buffer before allocation, so a crafted
// prefix cannot drive an unbounded make([][]byte, 0, count) (DoS). Each item
// needs at least one byte (its varint length prefix), so a count larger than
// the remaining buffer is malformed.
func TestCovenantUnmarshalBinaryRejectsOversizedItemCount(t *testing.T) {
	// Type byte, then varint 0xfe + uint32(0xffffffff) = ~4.29 billion items,
	// with no item bytes following.
	payload := []byte{0x00, 0xfe, 0xff, 0xff, 0xff, 0xff}

	var decoded Covenant
	if err := decoded.UnmarshalBinary(payload); err == nil {
		t.Fatal("UnmarshalBinary accepted an oversized item count, want error")
	}

	if decoded.Len() != 0 {
		t.Fatalf("UnmarshalBinary allocated %d items on a rejected payload", decoded.Len())
	}
}

// TestCovenantUnmarshalBinaryValidStillRoundTrips guards the bound against
// rejecting legitimate covenants: a real encoded covenant must still decode.
func TestCovenantUnmarshalBinaryValidStillRoundTrips(t *testing.T) {
	original := Covenant{}
	original.SetOpen(Hash{1, 2, 3}, []byte("example"))

	encoded, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	var decoded Covenant
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("UnmarshalBinary rejected a valid payload: %v", err)
	}

	if decoded.Type != original.Type || decoded.Len() != original.Len() {
		t.Fatalf("decoded covenant = %+v, want type %d with %d items", decoded, original.Type, original.Len())
	}
}

// BEGIN v090 file-aware generated coverage

func TestCovenantBinary_Covenant_MarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantBinary_Covenant_MarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantBinary_Covenant_MarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantBinary_Covenant_UnmarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).UnmarshalBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantBinary_Covenant_UnmarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).UnmarshalBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantBinary_Covenant_UnmarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).UnmarshalBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
