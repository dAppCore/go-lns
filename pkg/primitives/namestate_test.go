// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"bytes"
	"encoding/json"
	"math"
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

func TestNameStateToStats(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:    10,
		LockupPeriod:    20,
		BiddingPeriod:   30,
		RevealPeriod:    40,
		RenewalWindow:   50,
		AuctionMaturity: 25,
		BlocksPerDay:    250,
		TransferLockup:  15,
	}

	opening := NameState{Height: 100}
	stats := opening.ToStats(100, rules)
	if stats == nil {
		t.Fatal("ToStats should return stats for an active name")
	}

	if stats.OpenPeriodStart != 100 || stats.OpenPeriodEnd != 111 || stats.BlocksUntilBidding != 11 {
		t.Fatalf("unexpected opening stats: %+v", stats)
	}

	if stats.HoursUntilBidding != 1.06 {
		t.Fatalf("unexpected opening hours: got %v, want 1.06", stats.HoursUntilBidding)
	}

	closed := NameState{
		Height:   100,
		Renewal:  150,
		Owner:    Outpoint{TxHash: Hash{1, 2, 3}, Index: 4},
		Transfer: 180,
	}

	stats = closed.ToStats(188, rules)
	if stats == nil {
		t.Fatal("ToStats should return stats for a closed name")
	}

	if stats.RenewalPeriodStart != 150 || stats.RenewalPeriodEnd != 200 || stats.BlocksUntilExpire != 12 {
		t.Fatalf("unexpected closed stats: %+v", stats)
	}

	if math.Abs(stats.DaysUntilExpire-0.05) > 1e-9 {
		t.Fatalf("unexpected expire days: got %v, want 0.05", stats.DaysUntilExpire)
	}

	if stats.TransferLockupStart != 180 || stats.TransferLockupEnd != 195 || stats.BlocksUntilValidFinalize != 7 {
		t.Fatalf("unexpected transfer stats: %+v", stats)
	}

	if math.Abs(stats.HoursUntilValidFinalize-0.67) > 1e-9 {
		t.Fatalf("unexpected transfer hours: got %v, want 0.67", stats.HoursUntilValidFinalize)
	}

	expired := NameState{
		Height:  100,
		Renewal: 120,
		Owner:   Outpoint{TxHash: Hash{1, 2, 3}, Index: 4},
	}

	stats = expired.ToStats(210, rules)
	if stats == nil {
		t.Fatal("ToStats should return stats for an expired name with an owner")
	}

	if stats.BlocksSinceExpired != 40 {
		t.Fatalf("unexpected expired stats: %+v", stats)
	}

	if got := (NameState{
		Height:  100,
		Renewal: 120,
		Owner:   Outpoint{Index: ^uint32(0)},
	}).ToStats(210, rules); got != nil {
		t.Fatal("ToStats should return nil for an expired name without an owner")
	}
}

func TestNameStateJSONRoundTrip(t *testing.T) {
	rules := NameStateRules{
		TreeInterval:    10,
		LockupPeriod:    20,
		BiddingPeriod:   30,
		RevealPeriod:    40,
		ClaimPeriod:     140,
		RenewalWindow:   50,
		AuctionMaturity: 25,
		BlocksPerDay:    250,
		TransferLockup:  15,
	}

	original := NameState{
		Name:       []byte("example"),
		NameHash:   Hash{1, 2, 3, 4},
		Height:     100,
		Renewal:    150,
		Owner:      Outpoint{TxHash: Hash{5, 6, 7, 8}, Index: 9},
		Value:      10,
		Highest:    11,
		Data:       []byte{0xaa, 0xbb},
		Transfer:   180,
		Revoked:    0,
		Claimed:    0,
		Renewals:   12,
		Registered: true,
		Expired:    false,
		Weak:       true,
	}

	jsonView := original.GetJSON(188, rules)
	if jsonView.State != "CLOSED" {
		t.Fatalf("GetJSON should include the derived state, got %q", jsonView.State)
	}

	if jsonView.Stats == nil {
		t.Fatal("GetJSON should include stats for an active name")
	}

	if jsonView.Stats.BlocksUntilExpire != 12 {
		t.Fatalf("GetJSON should expose stats, got %+v", jsonView.Stats)
	}

	raw, err := json.Marshal(jsonView)
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	var decoded NameStateJSON
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("json.Unmarshal returned error: %v", err)
	}

	var roundTripped NameState
	if err := roundTripped.FromJSON(decoded); err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}

	if !bytes.Equal(roundTripped.Name, original.Name) {
		t.Fatalf("round-tripped name = %q, want %q", roundTripped.Name, original.Name)
	}

	if roundTripped.NameHash != original.NameHash {
		t.Fatalf("round-tripped name hash = %x, want %x", roundTripped.NameHash, original.NameHash)
	}

	if roundTripped.Owner != original.Owner {
		t.Fatalf("round-tripped owner = %+v, want %+v", roundTripped.Owner, original.Owner)
	}

	if !bytes.Equal(roundTripped.Data, original.Data) {
		t.Fatalf("round-tripped data = %x, want %x", roundTripped.Data, original.Data)
	}

	if roundTripped.Transfer != original.Transfer ||
		roundTripped.Registered != original.Registered ||
		roundTripped.Weak != original.Weak {
		t.Fatalf("round-tripped state did not preserve fields: %+v", roundTripped)
	}
}

func TestNameStateBinaryRoundTrip(t *testing.T) {
	original := NameState{
		Name:       []byte("example"),
		NameHash:   Hash{1, 2, 3, 4},
		Height:     100,
		Renewal:    150,
		Owner:      Outpoint{TxHash: Hash{5, 6, 7, 8}, Index: 9},
		Value:      10,
		Highest:    11,
		Data:       []byte{0xaa, 0xbb},
		Transfer:   180,
		Revoked:    190,
		Claimed:    2,
		Renewals:   12,
		Registered: true,
		Expired:    true,
		Weak:       true,
	}

	encoded, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	if got := original.GetSize(); got != len(encoded) {
		t.Fatalf("GetSize() = %d, want encoded length %d", got, len(encoded))
	}

	var decoded NameState
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if !bytes.Equal(decoded.Name, original.Name) {
		t.Fatalf("decoded name = %q, want %q", decoded.Name, original.Name)
	}

	if decoded.NameHash != (Hash{}) {
		t.Fatalf("decoded name hash = %x, want zero value", decoded.NameHash)
	}

	if decoded.Height != original.Height ||
		decoded.Renewal != original.Renewal ||
		decoded.Owner != original.Owner ||
		decoded.Value != original.Value ||
		decoded.Highest != original.Highest ||
		decoded.Transfer != original.Transfer ||
		decoded.Revoked != original.Revoked ||
		decoded.Claimed != original.Claimed ||
		decoded.Renewals != original.Renewals ||
		decoded.Registered != original.Registered ||
		decoded.Expired != original.Expired ||
		decoded.Weak != original.Weak {
		t.Fatalf("decoded binary state did not preserve fields: %+v", decoded)
	}

	if !bytes.Equal(decoded.Data, original.Data) {
		t.Fatalf("decoded data = %x, want %x", decoded.Data, original.Data)
	}
}

func TestNameStateBinaryZeroValue(t *testing.T) {
	var original NameState

	encoded, err := original.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	if got := original.GetSize(); got != len(encoded) {
		t.Fatalf("GetSize() = %d, want encoded length %d", got, len(encoded))
	}

	var decoded NameState
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if len(decoded.Name) != 0 ||
		decoded.NameHash != (Hash{}) ||
		decoded.Height != 0 ||
		decoded.Renewal != 0 ||
		decoded.Owner != (Outpoint{}) ||
		decoded.Value != 0 ||
		decoded.Highest != 0 ||
		len(decoded.Data) != 0 ||
		decoded.Transfer != 0 ||
		decoded.Revoked != 0 ||
		decoded.Claimed != 0 ||
		decoded.Renewals != 0 ||
		decoded.Registered ||
		decoded.Expired ||
		decoded.Weak {
		t.Fatalf("decoded zero value = %+v, want zero value", decoded)
	}
}
