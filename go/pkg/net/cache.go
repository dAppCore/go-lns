// SPDX-License-Identifier: EUPL-1.2

package net

import (
	"iter"
	"slices"

	core "github.com/dappcore/go"
)

// Source identifies where a NameRecord came from (RFC §12.2). HSD-sourced
// records are authoritative; gossip-sourced records carry a trust penalty and
// are subject to eviction when HSD contradicts them.
type Source = string

// Source constants enumerate the provenance of a cached record.
const (
	// SourceHSD marks a record fetched from the HSD sidechain.
	SourceHSD Source = "hsd"
	// SourceChainAlias marks a record derived from a main-chain alias.
	SourceChainAlias Source = "chain-alias"
	// SourceGossip marks a record replicated from a peer over gossip.
	SourceGossip Source = "gossip"
	// SourceHardcoded marks a record baked into the binary.
	SourceHardcoded Source = "hardcoded"
)

// NameRecord is a cached, resolvable name record (RFC §12.2). Alongside the
// resolution payload it tracks provenance: the Source, the GossipPeer that
// supplied it (gossip only), and the TrustScore after any gossip penalty.
//
//	nr := net.NameRecord{Name: "charon", Source: net.SourceHSD}
//	if nr.IsGossip() { /* sourced from a peer, penalty applied */ }
type NameRecord struct {
	// Name is the bare name label without the .lthn suffix.
	Name string `json:"name"`
	// PublicKey is the resolved chain public key for the name.
	PublicKey []byte `json:"publicKey,omitempty"`
	// Addresses are the resolved endpoints for the name.
	Addresses []string `json:"addresses,omitempty"`
	// Cap lists the capabilities advertised for the name.
	Cap []Capability `json:"cap,omitempty"`
	// Source is the provenance of this record.
	Source Source `json:"source"`
	// GossipPeer is the peer label that supplied this record (gossip only).
	GossipPeer string `json:"gossipPeer,omitempty"`
	// GossipPeerAddress is the network address of the supplying peer (gossip only).
	GossipPeerAddress string `json:"gossipPeerAddress,omitempty"`
	// TrustScore is the score after any gossip penalty, 0.0–1.0.
	TrustScore float64 `json:"trustScore"`
}

// IsGossip reports whether the record was replicated from a peer.
func (r NameRecord) IsGossip() bool {
	return r.Source == SourceGossip
}

// Clone returns a deep copy so cache mutations cannot alias caller-held slices.
func (r NameRecord) Clone() NameRecord {
	return NameRecord{
		Name:              r.Name,
		PublicKey:         append([]byte(nil), r.PublicKey...),
		Addresses:         append([]string(nil), r.Addresses...),
		Cap:               append([]Capability(nil), r.Cap...),
		Source:            r.Source,
		GossipPeer:        r.GossipPeer,
		GossipPeerAddress: r.GossipPeerAddress,
		TrustScore:        r.TrustScore,
	}
}

// Cache is the in-memory local name-record store gossip merges into. It is
// safe for concurrent use. Records are keyed by bare name; a later HSD-sourced
// record always supersedes a gossip-sourced one for the same name, but a gossip
// record never overwrites an HSD record (RFC §12.3 staleness rules).
//
//	cache := net.NewCache()
//	cache.Put(net.NameRecord{Name: "charon", Source: net.SourceHSD})
//	rec, ok := cache.Get("charon")
type Cache struct {
	mu      core.RWMutex
	records map[string]NameRecord
}

// NewCache constructs an empty cache ready for use.
//
//	cache := net.NewCache()
func NewCache() *Cache {
	return &Cache{records: make(map[string]NameRecord)}
}

// Get returns a copy of the cached record for a name and whether it was present.
//
//	rec, ok := cache.Get("charon")
func (c *Cache) Get(name string) (NameRecord, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	rec, ok := c.records[name]
	if !ok {
		return NameRecord{}, false
	}
	return rec.Clone(), true
}

// Put stores a record unconditionally, replacing any existing entry for the
// name. Callers that must honour HSD-over-gossip precedence should use Merge on
// a Gossip message rather than Put directly.
//
//	cache.Put(net.NameRecord{Name: "charon", Source: net.SourceHSD})
func (c *Cache) Put(rec NameRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.records == nil {
		c.records = make(map[string]NameRecord)
	}
	c.records[rec.Name] = rec.Clone()
}

// Delete removes a record by name. Removing an absent name is a no-op.
//
//	cache.Delete("charon")
func (c *Cache) Delete(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.records, name)
}

// Len returns the number of cached records.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.records)
}

// Records returns an iterator over copies of every cached record, ordered by
// name so iteration is deterministic.
//
//	for rec := range cache.Records() { /* ... */ }
func (c *Cache) Records() iter.Seq[NameRecord] {
	return func(yield func(NameRecord) bool) {
		c.mu.RLock()
		names := make([]string, 0, len(c.records))
		for name := range c.records {
			names = append(names, name)
		}
		c.mu.RUnlock()
		slices.Sort(names)
		for _, name := range names {
			if rec, ok := c.Get(name); ok {
				if !yield(rec) {
					return
				}
			}
		}
	}
}

// EvictGossipFrom removes every gossip-sourced record supplied by the named
// peer. It is the §12.3 eviction path for a peer that has dropped below the
// reachability floor, and returns the number of records evicted.
//
//	n := cache.EvictGossipFrom("charon")
func (c *Cache) EvictGossipFrom(peer string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	evicted := 0
	for name, rec := range c.records {
		if rec.Source == SourceGossip && rec.GossipPeer == peer {
			delete(c.records, name)
			evicted++
		}
	}
	return evicted
}
