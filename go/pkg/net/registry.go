// SPDX-License-Identifier: EUPL-1.2

package net

import (
	"iter"
	"slices"

	core "dappco.re/go"
)

// GossipMessageType is the transport message type carrying an LNS gossip delta
// across the go-p2p mesh (RFC §12.1). It namespaces LNS gossip away from the
// transport's own message types.
const GossipMessageType = "lns_gossip"

// GossipMessage is the transport envelope for a signed gossip delta. It is the
// minimal shape the gossip layer exchanges with the underlying transport: a
// type tag, the sender's node ID, and the JSON-encoded Gossip payload. The
// go-p2p node.Message envelope maps onto this shape, so a thin adapter binds the
// two once the go-p2p node package is published (see Transport).
//
//	msg := net.GossipMessage{Type: net.GossipMessageType, From: "self", Payload: raw}
type GossipMessage struct {
	// Type is the message type tag (GossipMessageType for LNS gossip).
	Type string `json:"type"`
	// From is the sending node's identifier.
	From string `json:"from"`
	// To is the recipient node ID, empty for a broadcast.
	To string `json:"to,omitempty"`
	// Payload is the JSON-encoded Gossip message.
	Payload []byte `json:"payload"`
}

// Transport is the seam the registry uses to announce and receive gossip over
// the network. It is deliberately the minimal subset of the go-p2p node
// transport the gossip layer needs (broadcast/send), so the LNS gossip layer is
// built ON go-p2p without importing a go-p2p package that is not yet published
// in a tagged release. A thin adapter over *node.Transport satisfies this
// interface; tests supply an in-memory fake.
//
// Deferred wiring: the live adapter binds to dappco.re/go/p2p/node.Transport
// (Broadcast/Send/OnMessage) once go-p2p publishes a tag that contains the node
// package. Until then the registry runs offline (announcements are no-ops) or
// against the test fake. See the package report for the deferred-wiring note.
type Transport interface {
	// Broadcast sends a gossip message to every connected peer.
	Broadcast(msg GossipMessage) core.Result
	// Send delivers a gossip message to a single peer by node ID.
	Send(peerID string, msg GossipMessage) core.Result
}

// Registry is the LNS peer registry and gossip surface (actions prefix
// lns.peer.*). It owns the registered PeerRecord set (keyed by bare name,
// ranked by trust score), the local NameRecord Cache that gossip merges into, a
// ScorePeer hook (Poindexter by default) used to score discovered peers, and an
// optional Transport that carries gossip across the go-p2p mesh.
//
// Register a peer, then look it up:
//
//	reg := net.NewRegistry(net.RegistryOptions{NodeID: "self"})
//	reg.Register(net.PeerRecord{Name: "charon", PublicKey: pub}, vec)
type Registry struct {
	nodeID    string
	cache     *Cache
	scorer    ScorePeer
	transport Transport

	mu      core.RWMutex
	records map[string]PeerRecord // bare name -> peer record
}

// RegistryOptions configures a Registry. All fields are optional: an empty
// options value yields a usable in-memory registry with the Poindexter-backed
// default scorer and a fresh cache.
//
//	opts := net.RegistryOptions{NodeID: "self", Transport: t}
type RegistryOptions struct {
	// NodeID is this node's identifier, stamped as the sender of announcements.
	NodeID string
	// Cache is the name-record cache gossip merges into. When nil a fresh cache
	// is created.
	Cache *Cache
	// Scorer derives peer trust scores. When nil DefaultScorer (Poindexter) is
	// used.
	Scorer ScorePeer
	// Transport carries gossip across the mesh. When nil, announcements are a
	// no-op (offline mode).
	Transport Transport
}

// NewRegistry constructs a Registry from the supplied options, defaulting the
// cache and scorer where omitted.
//
//	reg := net.NewRegistry(net.RegistryOptions{})
func NewRegistry(opts RegistryOptions) *Registry {
	cache := opts.Cache
	if cache == nil {
		cache = NewCache()
	}
	var scorer ScorePeer = opts.Scorer
	if scorer == nil {
		scorer = DefaultScorer{}
	}
	return &Registry{
		nodeID:    opts.NodeID,
		cache:     cache,
		scorer:    scorer,
		transport: opts.Transport,
		records:   make(map[string]PeerRecord),
	}
}

// Cache returns the registry's name-record cache.
func (r *Registry) Cache() *Cache { return r.cache }

// Register validates and stores a peer record, recomputing its trust score from
// the supplied vector via the configured scorer (Poindexter by default).
//
//	rec, err := reg.Register(net.PeerRecord{Name: "charon", PublicKey: pub}, vec)
func (r *Registry) Register(rec PeerRecord, vec PeerVector) (PeerRecord, error) {
	if err := rec.Validate(); err != nil {
		return PeerRecord{}, err
	}
	rec.TrustScore = r.scorer.ScorePeer(vec)

	r.mu.Lock()
	r.records[rec.Name] = rec.Clone()
	r.mu.Unlock()

	return rec.Clone(), nil
}

// Lookup returns a copy of the peer record for a bare name and whether it was
// found (RFC §10.1 resolution).
//
//	rec, ok := reg.Lookup("charon")
func (r *Registry) Lookup(name string) (PeerRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.records[name]
	if !ok {
		return PeerRecord{}, false
	}
	return rec.Clone(), true
}

// Remove deletes a peer record and its cached gossip. Removing an absent peer
// is a no-op.
//
//	reg.Remove("charon")
func (r *Registry) Remove(name string) {
	r.mu.Lock()
	delete(r.records, name)
	r.mu.Unlock()
	r.cache.EvictGossipFrom(name)
}

// Records returns an iterator over copies of every registered peer record,
// ordered by name for deterministic iteration.
//
//	for rec := range reg.Records() { /* ... */ }
func (r *Registry) Records() iter.Seq[PeerRecord] {
	return func(yield func(PeerRecord) bool) {
		r.mu.RLock()
		names := make([]string, 0, len(r.records))
		for name := range r.records {
			names = append(names, name)
		}
		r.mu.RUnlock()
		slices.Sort(names)
		for _, name := range names {
			if rec, ok := r.Lookup(name); ok {
				if !yield(rec) {
					return
				}
			}
		}
	}
}

// WithCapability returns the registered peers advertising the given capability,
// ordered by descending trust score so the strongest candidates come first. The
// LNS does not filter by score (RFC §11.3) — callers apply the routing minimum.
//
//	candidates := reg.WithCapability(net.CapVPN)
func (r *Registry) WithCapability(cap Capability) []PeerRecord {
	var out []PeerRecord
	for rec := range r.Records() {
		if rec.HasCapability(cap) {
			out = append(out, rec)
		}
	}
	slices.SortStableFunc(out, func(a, b PeerRecord) int {
		switch {
		case a.TrustScore > b.TrustScore:
			return -1
		case a.TrustScore < b.TrustScore:
			return 1
		default:
			return 0
		}
	})
	return out
}

// Announce broadcasts a signed gossip delta to connected peers over the
// transport (RFC §12.1 lns.gossip.announce). With no transport configured it is
// a no-op success. It fails if the gossip payload cannot be encoded.
//
//	res := reg.Announce(gossip)
func (r *Registry) Announce(gossip *Gossip) core.Result {
	if gossip == nil {
		return core.Fail(core.E("net.Registry.Announce", "nil gossip", nil))
	}
	if r.transport == nil {
		return core.Ok(nil)
	}
	encoded := core.JSONMarshal(gossip)
	if !encoded.OK {
		return encoded
	}
	payload, _ := encoded.Value.([]byte)
	return r.transport.Broadcast(GossipMessage{
		Type:    GossipMessageType,
		From:    r.nodeID,
		Payload: payload,
	})
}

// Receive verifies and merges an inbound gossip message into the cache against
// the supplied authority (RFC §12.1 lns.gossip.receive). It returns the merge
// outcome so callers can record acceptance/rejection counts.
//
//	res, err := reg.Receive(gossip, auth)
func (r *Registry) Receive(gossip *Gossip, auth GossipAuthority) (MergeResult, error) {
	if gossip == nil {
		return MergeResult{}, core.E("net.Registry.Receive", "nil gossip", nil)
	}
	return gossip.Merge(r.cache, auth)
}

// DecodeGossipMessage decodes the Gossip payload from a transport message,
// rejecting a message whose type is not GossipMessageType and bounding the
// decoded record count (RFC §12 untrusted-length guard). It is the inbound
// counterpart to Announce: a transport adapter calls it on each received
// message before handing the Gossip to Receive.
//
//	gossip, err := net.DecodeGossipMessage(msg)
func DecodeGossipMessage(msg GossipMessage) (*Gossip, error) {
	if msg.Type != GossipMessageType {
		return nil, core.Errorf("net.DecodeGossipMessage: unexpected message type %q", msg.Type)
	}
	var gossip Gossip
	if res := core.JSONUnmarshal(msg.Payload, &gossip); !res.OK {
		err, _ := res.Value.(error)
		return nil, core.E("net.DecodeGossipMessage", "decode gossip payload", err)
	}
	if len(gossip.Records) > MaxGossipRecords {
		return nil, core.Errorf("net.DecodeGossipMessage: record count %d exceeds %d", len(gossip.Records), MaxGossipRecords)
	}
	return &gossip, nil
}
