// SPDX-License-Identifier: EUPL-1.2

package primitives

import "testing"

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
