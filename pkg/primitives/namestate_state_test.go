// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "dappco.re/go"
	"testing"
)

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

// BEGIN v090 file-aware generated coverage

func TestNamestateState_NameStateStatus_String_Good(t *core.T) {
	fn := ax7ValueOf(new(NameStateStatus).String)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameStateStatus_String_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameStateStatus).String)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameStateStatus_String_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameStateStatus).String)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_State_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).State)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_State_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).State)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_State_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).State)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetState_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetState)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetState_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetState)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetState_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetState)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsOpening_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsOpening)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsOpening_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsOpening)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsOpening_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsOpening)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsOpening_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsOpening)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsOpening_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsOpening)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsOpening_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsOpening)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsLocked_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsLocked_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsLocked_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsLocked_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsBidding_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsBidding)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsBidding_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsBidding)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsBidding_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsBidding)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsBidding_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsBidding)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsBidding_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsBidding)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsBidding_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsBidding)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsReveal_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsReveal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsReveal_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsReveal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsReveal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsReveal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsReveal_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsReveal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsReveal_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsReveal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsReveal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsReveal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsClosed_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsClosed)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsClosed_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsClosed)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsClosed_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsClosed)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsClosed_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsClosed)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsClosed_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsClosed)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsClosed_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsClosed)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsRevoked_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsRevoked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsRevoked_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsRevoked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsRevoked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsRevoked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsRevoked_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsRevoked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsRevoked_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsRevoked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsRevoked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsRevoked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsRedeemable_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsRedeemable)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsRedeemable_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsRedeemable)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsRedeemable_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsRedeemable)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsRedeemable_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsRedeemable)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsRedeemable_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsRedeemable)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsRedeemable_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsRedeemable)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsClaimable_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsClaimable)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsClaimable_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsClaimable)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsClaimable_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsClaimable)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsClaimable_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsClaimable)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsClaimable_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsClaimable)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsClaimable_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsClaimable)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsExpired_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsExpired)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsExpired_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsExpired)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_IsExpired_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsExpired)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsExpired_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsExpired)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsExpired_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsExpired)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateState_NameState_GetIsExpired_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsExpired)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
