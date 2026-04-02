// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"testing"
)

func TestNameStateIsNull(t *testing.T) {
	if got := (NameState{}).IsNull(); !got {
		t.Fatal("zero name state should be null")
	}

	ns := NameState{Name: []byte("foo")}
	if got := ns.IsNull(); got {
		t.Fatal("name state with data should not be null")
	}
}

func TestNameStateCloneAndInject(t *testing.T) {
	original := NameState{
		Name:       []byte("foo"),
		NameHash:   Hash{1, 2, 3},
		Height:     10,
		Renewal:    11,
		Owner:      Outpoint{TxHash: Hash{4, 5, 6}, Index: 7},
		Value:      8,
		Highest:    9,
		Data:       []byte{0xaa, 0xbb},
		Transfer:   12,
		Revoked:    13,
		Claimed:    14,
		Renewals:   15,
		Registered: true,
		Expired:    true,
		Weak:       true,
	}

	clone := original.Clone()
	if clone.NameHash != original.NameHash ||
		clone.Height != original.Height ||
		clone.Renewal != original.Renewal ||
		clone.Owner != original.Owner ||
		clone.Value != original.Value ||
		clone.Highest != original.Highest ||
		clone.Transfer != original.Transfer ||
		clone.Revoked != original.Revoked ||
		clone.Claimed != original.Claimed ||
		clone.Renewals != original.Renewals ||
		clone.Registered != original.Registered ||
		clone.Expired != original.Expired ||
		clone.Weak != original.Weak {
		t.Fatal("Clone should copy the full name state")
	}

	if len(clone.Name) == 0 || len(clone.Data) == 0 {
		t.Fatal("test fixture must include non-empty slices")
	}

	if &clone.Name[0] != &original.Name[0] {
		t.Fatal("Clone should be shallow for name bytes")
	}

	if &clone.Data[0] != &original.Data[0] {
		t.Fatal("Clone should be shallow for data bytes")
	}

	var injected NameState
	if got := injected.Inject(original); got != &injected {
		t.Fatal("Inject should return the receiver")
	}

	if injected.NameHash != original.NameHash ||
		injected.Height != original.Height ||
		injected.Renewal != original.Renewal ||
		injected.Owner != original.Owner ||
		injected.Value != original.Value ||
		injected.Highest != original.Highest ||
		injected.Transfer != original.Transfer ||
		injected.Revoked != original.Revoked ||
		injected.Claimed != original.Claimed ||
		injected.Renewals != original.Renewals ||
		injected.Registered != original.Registered ||
		injected.Expired != original.Expired ||
		injected.Weak != original.Weak {
		t.Fatal("Inject should copy the full name state")
	}
}
