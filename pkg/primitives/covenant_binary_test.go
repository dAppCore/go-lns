// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "dappco.re/go"
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

// BEGIN v090 file-aware generated coverage

func TestCovenantBinary_Covenant_MarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Covenant).MarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCovenantBinary_Covenant_MarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Covenant).MarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCovenantBinary_Covenant_MarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Covenant).MarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCovenantBinary_Covenant_UnmarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Covenant).UnmarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCovenantBinary_Covenant_UnmarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Covenant).UnmarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCovenantBinary_Covenant_UnmarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Covenant).UnmarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
