// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

func TestCountHelpers(t *testing.T) {
	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeClaim)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeUpdate)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeTransfer)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRevoke)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRegister)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRenew)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeFinalize)}},
		},
	}

	if got := CountOpens(tx); got != 1 {
		t.Fatalf("CountOpens() = %d, want 1", got)
	}

	if got := CountUpdates(tx); got != 5 {
		t.Fatalf("CountUpdates() = %d, want 5", got)
	}

	if got := CountRenewals(tx); got != 3 {
		t.Fatalf("CountRenewals() = %d, want 3", got)
	}
}

func TestGrindName(t *testing.T) {
	rules := NameRules{
		AuctionStart:    100,
		RolloutInterval: 7,
	}

	name, err := GrindName(8, 100, rules)
	if err != nil {
		t.Fatalf("GrindName returned error: %v", err)
	}

	if len(name) != 8 {
		t.Fatalf("GrindName returned %q with length %d, want 8", name, len(name))
	}

	if !VerifyString(name) {
		t.Fatalf("GrindName returned invalid name %q", name)
	}

	hash, err := HashString(name)
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if !HasRollout(hash, 100, rules) {
		t.Fatalf("GrindName returned %q that does not satisfy rollout", name)
	}

	if HasReservedHash(hash) {
		t.Fatalf("GrindName returned reserved name %q", name)
	}

	if _, err := GrindName(0, 100, rules); err == nil {
		t.Fatal("GrindName should reject zero-length names")
	}
}

func TestNameSetHelpers(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], []byte{0, 0, 0, 0}, []byte("example-name")}}},
			{Covenant: primitives.Covenant{Type: uint8(TypeBid), Items: [][]byte{hash[:], []byte{1, 0, 0, 0}, []byte("example-name"), make([]byte, 32)}}},
		},
	}

	set := map[primitives.Hash]struct{}{
		hash: {},
	}

	if !HasNames(tx, set) {
		t.Fatal("HasNames should report a matching name hash")
	}

	RemoveNames(tx, set)

	if len(set) != 0 {
		t.Fatalf("RemoveNames left %d entries, want 0", len(set))
	}

	AddNames(tx, set)

	if len(set) != 1 {
		t.Fatalf("AddNames left %d entries, want 1", len(set))
	}
}

func TestHasSaneCovenants(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	valid := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], []byte{0, 0, 0, 0}, []byte("example-name")}}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRegister), Items: [][]byte{hash[:], []byte{1, 0, 0, 0}, []byte("record"), make([]byte, 32)}}},
		},
	}

	if !HasSaneCovenants(valid) {
		t.Fatal("HasSaneCovenants should accept structurally valid covenants")
	}

	invalid := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], []byte{1, 0, 0, 0}, []byte("example-name")}}},
		},
	}

	if HasSaneCovenants(invalid) {
		t.Fatal("HasSaneCovenants should reject non-zero open heights")
	}
}
