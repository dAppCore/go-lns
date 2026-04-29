// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"math"
	"reflect"
	"testing"

	core "dappco.re/go"
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

func TestNameStateGetAliases(t *testing.T) {
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

	ns := NameState{
		Name:       []byte("example"),
		NameHash:   Hash{1, 2, 3, 4},
		Height:     100,
		Renewal:    150,
		Owner:      Outpoint{TxHash: Hash{5, 6, 7, 8}, Index: 9},
		Value:      10,
		Highest:    11,
		Data:       []byte{0xaa, 0xbb},
		Transfer:   180,
		Claimed:    1,
		Renewals:   12,
		Registered: true,
		Weak:       true,
	}

	if ns.GetDelta() == nil {
		t.Fatal("GetDelta should allocate a delta on demand")
	}

	if ns.GetHasDelta() {
		t.Fatal("GetHasDelta should stay false until the delta carries changes")
	}

	ns.GetDelta().Height = new(uint32)
	if !ns.GetHasDelta() {
		t.Fatal("GetHasDelta should report a populated delta")
	}

	if got := ns.GetClone(); got.NameHash != ns.NameHash || got.Renewal != ns.Renewal {
		t.Fatalf("GetClone returned %+v, want %+v", got, ns)
	}

	if got := new(NameState).GetInject(ns); got == nil || got.NameHash != ns.NameHash {
		t.Fatalf("GetInject returned %+v, want %+v", got, ns)
	}

	if got := new(NameState).GetSet([]byte("alias"), 200); got == nil || got.Height != 200 || string(got.Name) != "alias" {
		t.Fatalf("GetSet returned %+v", got)
	}

	if got := new(NameState).GetReset(200); got == nil || got.Height != 200 || got.Renewal != 200 {
		t.Fatalf("GetReset returned %+v", got)
	}

	if got := new(NameState).GetClear(); got == nil {
		t.Fatal("GetClear should return the receiver")
	}

	if got := (NameState{}).GetIsNull(); !got {
		t.Fatal("GetIsNull should report the zero-value state shape")
	}

	if ns.GetIsNull() {
		t.Fatal("GetIsNull should reject populated states")
	}

	if got := ns.GetState(100, rules); got != NameStateLocked {
		t.Fatalf("GetState returned %s, want LOCKED", got)
	}

	if !ns.GetIsLocked(100, rules) || !ns.GetIsClaimable(120, rules) {
		t.Fatal("phase predicate aliases should match the underlying helpers")
	}

	if ns.GetIsOpening(100, rules) || ns.GetIsBidding(140, rules) || ns.GetIsReveal(170, rules) {
		t.Fatal("phase predicate aliases should reflect the current state")
	}

	if !ns.GetIsClosed(200, rules) || ns.GetIsRevoked(200, rules) || !ns.GetIsRedeemable(200, rules) {
		t.Fatal("state aliases should report closed names correctly")
	}

	expired := ns.GetClone()
	if !expired.GetMaybeExpire(220, rules) {
		t.Fatal("GetMaybeExpire should expire the closed name after the renewal window")
	}

	if stats := ns.GetClone().GetToStats(188, rules); stats == nil || stats.RenewalPeriodStart != 150 {
		t.Fatalf("GetToStats returned %+v", stats)
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

	if alias := original.GetFormat(188, rules); !reflect.DeepEqual(alias, jsonView) {
		t.Fatalf("GetFormat returned %+v, want %+v", alias, jsonView)
	}

	raw := core.JSONMarshal(jsonView)
	if !raw.OK {
		t.Fatalf("JSONMarshal returned error: %v", raw.Value)
	}

	var decoded NameStateJSON
	if r := core.JSONUnmarshal(raw.Value.([]byte), &decoded); !r.OK {
		t.Fatalf("JSONUnmarshal returned error: %v", r.Value)
	}

	var roundTripped NameState
	if err := roundTripped.GetFromJSON(decoded); err != nil {
		t.Fatalf("GetFromJSON returned error: %v", err)
	}

	if !equalBytes(roundTripped.Name, original.Name) {
		t.Fatalf("round-tripped name = %q, want %q", roundTripped.Name, original.Name)
	}

	if roundTripped.NameHash != original.NameHash {
		t.Fatalf("round-tripped name hash = %x, want %x", roundTripped.NameHash, original.NameHash)
	}

	if roundTripped.Owner != original.Owner {
		t.Fatalf("round-tripped owner = %+v, want %+v", roundTripped.Owner, original.Owner)
	}

	if !equalBytes(roundTripped.Data, original.Data) {
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

	if !equalBytes(decoded.Name, original.Name) {
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

	if !equalBytes(decoded.Data, original.Data) {
		t.Fatalf("decoded data = %x, want %x", decoded.Data, original.Data)
	}
}

func TestNameStateMaxSize(t *testing.T) {
	if MaxNameStateSize != 668 {
		t.Fatalf("MaxNameStateSize = %d, want 668", MaxNameStateSize)
	}

	if NameStateMaxSize != MaxNameStateSize {
		t.Fatalf("NameStateMaxSize = %d, want %d", NameStateMaxSize, MaxNameStateSize)
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

// BEGIN v090 file-aware generated coverage

func TestNamestate_NameState_Delta_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).Delta)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Delta_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).Delta)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Delta_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).Delta)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetDelta_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetDelta)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetDelta_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetDelta)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetDelta_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetDelta)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Inject_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).Inject)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Inject_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).Inject)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Inject_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).Inject)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetInject_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetInject)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetInject_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetInject)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetInject_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetInject)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Clone_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).Clone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Clone_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).Clone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Clone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).Clone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetClone_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetClone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetClone_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetClone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetClone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetClone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Clear_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).Clear)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Clear_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).Clear)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Clear_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).Clear)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetClear_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetClear)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetClear_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetClear)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetClear_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetClear)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_HasDelta_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).HasDelta)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_HasDelta_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).HasDelta)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_HasDelta_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).HasDelta)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetHasDelta_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetHasDelta)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetHasDelta_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetHasDelta)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetHasDelta_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetHasDelta)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Reset_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).Reset)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Reset_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).Reset)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Reset_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).Reset)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetReset_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetReset)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetReset_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetReset)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetReset_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetReset)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Set_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).Set)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Set_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).Set)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_Set_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).Set)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetSet_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetSet)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetSet_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetSet)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetSet_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetSet)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_MaybeExpire_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).MaybeExpire)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_MaybeExpire_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).MaybeExpire)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_MaybeExpire_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).MaybeExpire)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetMaybeExpire_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetMaybeExpire)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetMaybeExpire_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetMaybeExpire)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetMaybeExpire_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetMaybeExpire)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_ApplyState_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).ApplyState)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_ApplyState_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).ApplyState)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_ApplyState_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).ApplyState)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetApplyState_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetApplyState)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetApplyState_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetApplyState)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetApplyState_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetApplyState)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_IsNull_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsNull)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_IsNull_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsNull)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_IsNull_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).IsNull)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetIsNull_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsNull)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetIsNull_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsNull)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestate_NameState_GetIsNull_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetIsNull)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
