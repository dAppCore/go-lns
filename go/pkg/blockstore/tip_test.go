// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	"errors"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/store"
)

// TestTip_Good checks the tip advances as higher blocks are stored and reports
// the hash and height of the most-work block.
func TestTip_Good(t *testing.T) {
	bs := newTestStore(t)
	g, one, two := hashOf(0x00), hashOf(0x01), hashOf(0x02)

	if err := bs.PutBlock(g, 0, []byte("genesis")); err != nil {
		t.Fatalf("PutBlock genesis: %v", err)
	}
	if err := bs.PutBlock(one, 1, []byte("one")); err != nil {
		t.Fatalf("PutBlock one: %v", err)
	}
	if err := bs.PutBlock(two, 2, []byte("two")); err != nil {
		t.Fatalf("PutBlock two: %v", err)
	}

	tipHash, tipHeight, err := bs.Tip()
	if err != nil {
		t.Fatalf("Tip: %v", err)
	}
	if tipHash != two || tipHeight != 2 {
		t.Fatalf("Tip = (%x, %d), want (%x, 2)", tipHash, tipHeight, two)
	}
}

// TestTip_Bad checks an empty store has no tip and HasTip reports false without
// erroring, while Tip reports ErrNotFound.
func TestTip_Bad(t *testing.T) {
	bs := newTestStore(t)

	if has, err := bs.HasTip(); err != nil || has {
		t.Fatalf("HasTip on empty = (%v, %v), want (false, nil)", has, err)
	}
	if _, _, err := bs.Tip(); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Tip on empty = %v, want ErrNotFound", err)
	}
}

// TestTip_Ugly checks two things: backfilling a lower block does not move the tip
// backward, and SetTip can explicitly rewind the tip during a reorganisation.
func TestTip_Ugly(t *testing.T) {
	bs := newTestStore(t)
	low, high := hashOf(0x05), hashOf(0x09)

	if err := bs.PutBlock(high, 9, []byte("high")); err != nil {
		t.Fatalf("PutBlock high: %v", err)
	}
	// Backfill an earlier block: tip must stay at height 9.
	if err := bs.PutBlock(low, 5, []byte("low")); err != nil {
		t.Fatalf("PutBlock low: %v", err)
	}
	if _, height, _ := bs.Tip(); height != 9 {
		t.Fatalf("tip moved to %d after backfill, want 9", height)
	}

	// Explicit rewind for a reorg.
	if err := bs.SetTip(low, 5); err != nil {
		t.Fatalf("SetTip rewind: %v", err)
	}
	if h, height, _ := bs.Tip(); h != low || height != 5 {
		t.Fatalf("after SetTip = (%x, %d), want (%x, 5)", h, height, low)
	}
}

// TestTip_Corrupt_Ugly checks a corrupt tip-height value is reported as
// ErrCorrupt rather than parsed into a bogus height.
func TestTip_Corrupt_Ugly(t *testing.T) {
	storeInstance, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { _ = storeInstance.Close() })
	scoped := store.NewScoped(storeInstance, "corrupt-tip")
	bs, _ := New(scoped)

	hash := hashOf(0x42)
	if err := scoped.SetIn(groupMeta, metaTipHash, core.HexEncode(hash[:])); err != nil {
		t.Fatalf("seed tip hash: %v", err)
	}
	if err := scoped.SetIn(groupMeta, metaTipHeight, "not-a-number"); err != nil {
		t.Fatalf("seed tip height: %v", err)
	}

	if _, _, err := bs.Tip(); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("Tip(corrupt height) = %v, want ErrCorrupt", err)
	}
}

// TestTip_NilReceiver_Ugly checks the tip helpers tolerate a nil receiver.
func TestTip_NilReceiver_Ugly(t *testing.T) {
	var nilStore *Store
	if _, _, err := nilStore.Tip(); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil Tip = %v, want ErrNilStore", err)
	}
	if err := nilStore.SetTip(hashOf(1), 1); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil SetTip = %v, want ErrNilStore", err)
	}
	if _, err := nilStore.HasTip(); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil HasTip = %v, want ErrNilStore", err)
	}
}
