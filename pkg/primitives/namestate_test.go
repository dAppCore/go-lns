// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"testing"
)

func TestNameStateIsNull(t *testing.T) {
	if got := (NameState{}).IsNull(); !got {
		t.Fatal("zero name state should be null")
	}

	ns := NameState{
		Name:     []byte("foo"),
		NameHash: Hash{1, 2, 3},
	}
	if got := ns.IsNull(); !got {
		t.Fatal("name metadata alone should not make the state non-null")
	}

	ns.Height = 1
	if got := ns.IsNull(); got {
		t.Fatal("name state with chain data should not be null")
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

	if injected.delta != nil {
		t.Fatal("Inject should not copy cached delta state")
	}
}

func TestNameStateDeltaHelpers(t *testing.T) {
	var ns NameState

	if got := ns.HasDelta(); got {
		t.Fatal("zero name state should not report a delta")
	}

	delta := ns.Delta()
	if delta == nil {
		t.Fatal("Delta should allocate a delta on demand")
	}

	if got := ns.HasDelta(); got {
		t.Fatal("empty delta should not count as a change")
	}

	height := uint32(42)
	owner := Outpoint{TxHash: Hash{1, 2, 3}, Index: 7}
	data := []byte("payload")
	registered := true

	delta.Height = &height
	delta.Owner = &owner
	delta.Data = &data
	delta.Registered = &registered

	if got := ns.HasDelta(); !got {
		t.Fatal("populated delta should be reported")
	}

	ns.ApplyState(delta)

	if ns.Height != 42 || ns.Owner != owner || string(ns.Data) != "payload" || !ns.Registered {
		t.Fatal("ApplyState should merge populated delta fields")
	}

	if got := ns.Clear(); got != &ns {
		t.Fatal("Clear should return the receiver")
	}

	if ns.delta != nil {
		t.Fatal("Clear should discard the cached delta")
	}
}

func TestNameStateMaybeExpire(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:    0,
		LockupPeriod:    1,
		BiddingPeriod:   1,
		RevealPeriod:    1,
		ClaimPeriod:     6,
		RenewalWindow:   7,
		AuctionMaturity: 8,
	}

	ns := NameState{
		Height:     10,
		Renewal:    10,
		Owner:      Outpoint{TxHash: Hash{1, 2, 3}, Index: 4},
		Value:      12,
		Highest:    13,
		Data:       []byte("payload"),
		Transfer:   14,
		Revoked:    0,
		Claimed:    0,
		Renewals:   15,
		Registered: true,
		Expired:    false,
		Weak:       true,
	}

	if got := ns.MaybeExpire(13, rules); got {
		t.Fatal("MaybeExpire should not expire a closed name without a renewal timeout")
	}

	if ns.Expired {
		t.Fatal("MaybeExpire should leave the state unchanged when not expired")
	}

	ns.Renewal = 0
	if got := ns.MaybeExpire(13, rules); !got {
		t.Fatal("MaybeExpire should expire a name once the renewal window has elapsed")
	}

	if !ns.Expired {
		t.Fatal("MaybeExpire should mark the name expired")
	}

	if ns.Height != 13 || ns.Renewal != 13 {
		t.Fatalf("MaybeExpire should reset timestamps, got height=%d renewal=%d", ns.Height, ns.Renewal)
	}

	if !ns.Owner.IsNull() || ns.Value != 0 || ns.Highest != 0 || ns.Transfer != 0 || ns.Revoked != 0 || ns.Claimed != 0 || ns.Renewals != 0 || ns.Registered || ns.Weak {
		t.Fatal("MaybeExpire should clear mutable ownership and auction fields")
	}

	if string(ns.Data) != "payload" {
		t.Fatalf("MaybeExpire should preserve data, got %q", string(ns.Data))
	}
}
