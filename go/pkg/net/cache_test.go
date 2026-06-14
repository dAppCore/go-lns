// SPDX-License-Identifier: EUPL-1.2

package net

import (
	"slices"

	core "dappco.re/go"
)

// TestCache_PutGet_Good stores and retrieves a record, returning a deep copy.
func TestCache_PutGet_Good(t *core.T) {
	cache := NewCache()
	cache.Put(NameRecord{Name: "charon", PublicKey: []byte{1, 2}, Source: SourceHSD})

	rec, ok := cache.Get("charon")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "charon", rec.Name)

	// Mutating the returned copy must not affect the cached record.
	rec.PublicKey[0] = 9
	again, _ := cache.Get("charon")
	core.AssertEqual(t, byte(1), again.PublicKey[0])
}

// TestCache_Get_Bad returns not-found for an absent name.
func TestCache_Get_Bad(t *core.T) {
	cache := NewCache()
	_, ok := cache.Get("missing")
	core.AssertFalse(t, ok)
	core.AssertEqual(t, 0, cache.Len())
}

// TestCache_Delete_Ugly deletes an absent name without panicking, then deletes a
// present one.
func TestCache_Delete_Ugly(t *core.T) {
	cache := NewCache()
	cache.Delete("missing") // no-op
	cache.Put(NameRecord{Name: "charon", Source: SourceGossip})
	cache.Delete("charon")
	_, ok := cache.Get("charon")
	core.AssertFalse(t, ok)
}

// TestCache_Records_Good iterates records in deterministic name order.
func TestCache_Records_Good(t *core.T) {
	cache := NewCache()
	cache.Put(NameRecord{Name: "zeta", Source: SourceHSD})
	cache.Put(NameRecord{Name: "alpha", Source: SourceHSD})
	cache.Put(NameRecord{Name: "mu", Source: SourceHSD})

	var names []string
	for rec := range cache.Records() {
		names = append(names, rec.Name)
	}
	core.AssertTrue(t, slices.Equal([]string{"alpha", "mu", "zeta"}, names))
}

// TestCache_EvictGossipFrom_Good evicts only the gossip records from a given
// peer, leaving HSD records and other peers' gossip intact (RFC §12.3).
func TestCache_EvictGossipFrom_Good(t *core.T) {
	cache := NewCache()
	cache.Put(NameRecord{Name: "a", Source: SourceGossip, GossipPeer: "charon"})
	cache.Put(NameRecord{Name: "b", Source: SourceGossip, GossipPeer: "charon"})
	cache.Put(NameRecord{Name: "c", Source: SourceGossip, GossipPeer: "other"})
	cache.Put(NameRecord{Name: "d", Source: SourceHSD})

	evicted := cache.EvictGossipFrom("charon")
	core.AssertEqual(t, 2, evicted)
	core.AssertEqual(t, 2, cache.Len())
	_, ok := cache.Get("c")
	core.AssertTrue(t, ok)
	_, ok = cache.Get("d")
	core.AssertTrue(t, ok)
}

// TestCache_NameRecord_IsGossip_Ugly distinguishes gossip from non-gossip
// sources.
func TestCache_NameRecord_IsGossip_Ugly(t *core.T) {
	core.AssertTrue(t, NameRecord{Source: SourceGossip}.IsGossip())
	core.AssertFalse(t, NameRecord{Source: SourceHSD}.IsGossip())
	core.AssertFalse(t, NameRecord{Source: SourceHardcoded}.IsGossip())
	core.AssertFalse(t, NameRecord{}.IsGossip())
}
