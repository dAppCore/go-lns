// SPDX-License-Identifier: EUPL-1.2

package primitives

import "testing"

func TestNameStateStateTransitions(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:  10,
		LockupPeriod:  20,
		BiddingPeriod: 30,
		RevealPeriod:  40,
	}

	ns := NameState{Height: 100}

	cases := []struct {
		name   string
		height uint32
		state  NameStateStatus
	}{
		{name: "opening", height: 100, state: NameStateOpening},
		{name: "bidding", height: 111, state: NameStateBidding},
		{name: "reveal", height: 141, state: NameStateReveal},
		{name: "closed", height: 181, state: NameStateClosed},
	}

	for _, tc := range cases {
		if got := ns.State(tc.height, rules); got != tc.state {
			t.Fatalf("%s: State() = %s, want %s", tc.name, got, tc.state)
		}
	}
}

func TestNameStateStateClaimAndRevoke(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:  10,
		LockupPeriod:  20,
		BiddingPeriod: 30,
		RevealPeriod:  40,
	}

	claimed := NameState{Height: 100, Claimed: 1}
	if got := claimed.State(110, rules); got != NameStateLocked {
		t.Fatalf("claimed name at lockup height = %s, want LOCKED", got)
	}

	if got := claimed.State(120, rules); got != NameStateClosed {
		t.Fatalf("claimed name after lockup = %s, want CLOSED", got)
	}

	revoked := NameState{Height: 100, Revoked: 1}
	if got := revoked.State(110, rules); got != NameStateRevoked {
		t.Fatalf("revoked name = %s, want REVOKED", got)
	}
}

func TestNameStatePhasePredicates(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:  10,
		LockupPeriod:  20,
		BiddingPeriod: 30,
		RevealPeriod:  40,
	}

	ns := NameState{Height: 100}

	if !ns.IsOpening(100, rules) {
		t.Fatal("IsOpening should report the opening phase")
	}

	if !ns.IsBidding(111, rules) {
		t.Fatal("IsBidding should report the bidding phase")
	}

	if !ns.IsReveal(141, rules) {
		t.Fatal("IsReveal should report the reveal phase")
	}

	if !ns.IsClosed(181, rules) {
		t.Fatal("IsClosed should report the closed phase")
	}

	if !ns.IsRedeemable(181, rules) {
		t.Fatal("IsRedeemable should report closed or later phases")
	}

	if ns.IsRevoked(181, rules) {
		t.Fatal("IsRevoked should be false for a normal name")
	}
}

func TestNameStateClaimableAndExpired(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:    10,
		LockupPeriod:    20,
		BiddingPeriod:   30,
		RevealPeriod:    40,
		ClaimPeriod:     140,
		RenewalWindow:   50,
		AuctionMaturity: 25,
	}

	claimed := NameState{
		Height:  100,
		Renewal: 100,
		Claimed: 1,
	}

	if !claimed.IsClaimable(130, rules) {
		t.Fatal("IsClaimable should report claimed names before the claim cutoff")
	}

	rules.NoReserved = true
	if claimed.IsClaimable(130, rules) {
		t.Fatal("IsClaimable should be disabled when reserved names are not allowed")
	}
	rules.NoReserved = false

	if claimed.IsExpired(130, rules) {
		t.Fatal("IsExpired should stay false while the name is still within the claim window")
	}

	if claimed.IsExpired(140, rules) {
		t.Fatal("IsExpired should stay false while the name is still within the renewal window")
	}

	if !claimed.IsExpired(160, rules) {
		t.Fatal("IsExpired should become true after the renewal window elapses")
	}

	revoked := NameState{
		Height:  100,
		Revoked: 120,
	}

	if revoked.IsExpired(140, rules) {
		t.Fatal("revoked names should stay pending until auction maturity")
	}

	if !revoked.IsExpired(150, rules) {
		t.Fatal("revoked names should expire after auction maturity")
	}

	closedNoOwner := NameState{
		Height:  100,
		Renewal: 100,
	}

	if !closedNoOwner.IsExpired(190, rules) {
		t.Fatal("closed names without an owner should expire once the renewal window ends")
	}
}
