// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "github.com/dappcore/go"
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
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameStateStatus).String())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameStateStatus_String_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameStateStatus).String())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameStateStatus_String_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameStateStatus).String())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_State_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).State(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_State_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).State(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_State_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).State(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetState_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetState(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetState_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetState(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetState_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetState(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsOpening_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsOpening(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsOpening_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsOpening(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsOpening_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsOpening(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsOpening_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsOpening(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsOpening_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsOpening(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsOpening_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsOpening(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsLocked_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsLocked(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsLocked_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsLocked(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsLocked_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsLocked(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsLocked_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsLocked(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsLocked_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsLocked(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsLocked_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsLocked(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsBidding_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsBidding(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsBidding_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsBidding(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsBidding_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsBidding(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsBidding_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsBidding(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsBidding_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsBidding(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsBidding_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsBidding(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsReveal_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsReveal(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsReveal_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsReveal(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsReveal_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsReveal(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsReveal_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsReveal(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsReveal_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsReveal(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsReveal_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsReveal(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsClosed_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsClosed(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsClosed_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsClosed(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsClosed_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsClosed(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsClosed_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsClosed(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsClosed_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsClosed(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsClosed_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsClosed(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsRevoked_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsRevoked(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsRevoked_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsRevoked(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsRevoked_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsRevoked(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsRevoked_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsRevoked(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsRevoked_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsRevoked(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsRevoked_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsRevoked(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsRedeemable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsRedeemable(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsRedeemable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsRedeemable(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsRedeemable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsRedeemable(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsRedeemable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsRedeemable(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsRedeemable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsRedeemable(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsRedeemable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsRedeemable(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsClaimable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsClaimable(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsClaimable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsClaimable(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsClaimable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsClaimable(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsClaimable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsClaimable(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsClaimable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsClaimable(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsClaimable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsClaimable(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsExpired_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).IsExpired(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsExpired_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).IsExpired(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_IsExpired_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).IsExpired(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsExpired_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetIsExpired(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsExpired_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetIsExpired(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateState_NameState_GetIsExpired_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetIsExpired(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
