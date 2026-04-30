// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

func TestGetRolloutAndHasRollout(t *testing.T) {
	var hash primitives.Hash

	rules := NameRules{
		AuctionStart:    100,
		RolloutInterval: 7,
	}

	start, week := GetRollout(hash, rules)
	if start != 100 || week != 0 {
		t.Fatalf("GetRollout(zero) = (%d, %d), want (100, 0)", start, week)
	}

	if !HasRollout(hash, 100, rules) {
		t.Fatal("HasRollout should accept the rollout start height")
	}

	if !GetHasRollout(hash, 100, rules) {
		t.Fatal("GetHasRollout should alias HasRollout")
	}

	if HasRollout(hash, 99, rules) {
		t.Fatal("HasRollout should reject heights before rollout")
	}

	rules.NoRollout = true
	start, week = GetRollout(hash, rules)
	if start != 0 || week != 0 {
		t.Fatalf("GetRollout(no rollout) = (%d, %d), want (0, 0)", start, week)
	}

	if !HasRollout(hash, 0, rules) {
		t.Fatal("HasRollout should be disabled when rollout gating is off")
	}

	if !GetHasRollout(hash, 0, rules) {
		t.Fatal("GetHasRollout should alias HasRollout when rollout gating is off")
	}
}

func TestIsReserved(t *testing.T) {
	hash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	rules := NameRules{
		ClaimPeriod: 100,
	}

	if !IsReserved(hash, 99, rules) {
		t.Fatal("IsReserved should report reserved names before the claim period")
	}

	if IsReserved(hash, 100, rules) {
		t.Fatal("IsReserved should reject reserved names at or after the claim period")
	}

	if !GetIsReserved(hash, 99, rules) {
		t.Fatal("GetIsReserved should alias IsReserved")
	}

	rules.NoReserved = true
	if IsReserved(hash, 99, rules) {
		t.Fatal("IsReserved should be disabled when reserved-name protection is off")
	}
}

func TestIsLockedUp(t *testing.T) {
	rootHash, err := HashString(testLockedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	nonRootHash, err := HashString("openclassrooms")
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	rules := NameRules{
		ClaimPeriod:       100,
		AlexaLockupPeriod: 250,
	}

	if IsLockedUp(nonRootHash, 99, rules) {
		t.Fatal("IsLockedUp should be false before the claim period")
	}

	if !IsLockedUp(rootHash, 100, rules) {
		t.Fatal("IsLockedUp should keep root names locked after the claim period")
	}

	if !IsLockedUp(nonRootHash, 100, rules) {
		t.Fatal("IsLockedUp should lock non-root names until the Alexa threshold")
	}

	if IsLockedUp(nonRootHash, 250, rules) {
		t.Fatal("IsLockedUp should release non-root names after the Alexa threshold")
	}

	if !GetIsLockedUp(rootHash, 100, rules) {
		t.Fatal("GetIsLockedUp should alias IsLockedUp")
	}

	rules.NoReserved = true
	if IsLockedUp(rootHash, 100, rules) {
		t.Fatal("IsLockedUp should be disabled when reserved-name protection is off")
	}
}
