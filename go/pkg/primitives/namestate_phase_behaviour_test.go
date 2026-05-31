// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"testing"

	core "dappco.re/go"
)

func TestNameStatePhaseBehaviour_String(t *testing.T) {
	cases := []struct {
		status NameStateStatus
		want   string
	}{
		{NameStateOpening, "OPENING"},
		{NameStateLocked, "LOCKED"},
		{NameStateBidding, "BIDDING"},
		{NameStateReveal, "REVEAL"},
		{NameStateClosed, "CLOSED"},
		{NameStateRevoked, "REVOKED"},
		{NameStateStatus(99), "UNKNOWN"},
	}
	for _, tc := range cases {
		core.AssertEqual(t, tc.want, tc.status.String())
	}
}

func TestNameStatePhaseBehaviour_State_auctionPhases(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:  4, // openPeriod = 5
		BiddingPeriod: 10,
		RevealPeriod:  10,
		LockupPeriod:  20,
	}
	ns := NameState{Height: 100}

	// Within openPeriod (height < 100+5).
	core.AssertEqual(t, NameStateOpening, ns.State(102, rules))
	core.AssertTrue(t, ns.IsOpening(102, rules))

	// Bidding window: 105 <= height < 115.
	core.AssertEqual(t, NameStateBidding, ns.State(110, rules))
	core.AssertTrue(t, ns.IsBidding(110, rules))

	// Reveal window: 115 <= height < 125.
	core.AssertEqual(t, NameStateReveal, ns.State(120, rules))
	core.AssertTrue(t, ns.IsReveal(120, rules))

	// Closed: height >= 125.
	core.AssertEqual(t, NameStateClosed, ns.State(130, rules))
	core.AssertTrue(t, ns.IsClosed(130, rules))
	core.AssertTrue(t, ns.IsRedeemable(130, rules))
}

func TestNameStatePhaseBehaviour_State_claimedAndRevoked(t *testing.T) {
	rules := NameStateRules{LockupPeriod: 20}

	// Revoked dominates regardless of other fields.
	revoked := NameState{Height: 100, Revoked: 50}
	core.AssertEqual(t, NameStateRevoked, revoked.State(200, rules))
	core.AssertTrue(t, revoked.IsRevoked(200, rules))

	// Claimed name in its lockup window.
	claimedLocked := NameState{Height: 100, Claimed: 1}
	core.AssertEqual(t, NameStateLocked, claimedLocked.State(110, rules))
	core.AssertTrue(t, claimedLocked.IsLocked(110, rules))

	// Claimed name past its lockup is closed.
	core.AssertEqual(t, NameStateClosed, claimedLocked.State(130, rules))
}
