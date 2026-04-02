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
