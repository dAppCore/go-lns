// SPDX-License-Identifier: EUPL-1.2

package net

import (
	core "dappco.re/go"
)

// fakeTransport is an in-memory Transport that records broadcasts so the gossip
// announce path can be exercised without a live go-p2p mesh.
type fakeTransport struct {
	broadcasts []GossipMessage
	sends      map[string][]GossipMessage
}

func newFakeTransport() *fakeTransport {
	return &fakeTransport{sends: make(map[string][]GossipMessage)}
}

func (f *fakeTransport) Broadcast(msg GossipMessage) core.Result {
	f.broadcasts = append(f.broadcasts, msg)
	return core.Ok(nil)
}

func (f *fakeTransport) Send(peerID string, msg GossipMessage) core.Result {
	f.sends[peerID] = append(f.sends[peerID], msg)
	return core.Ok(nil)
}

// TestRegistry_Register_Good registers a peer, scores it via the default
// Poindexter scorer, and looks it back up.
func TestRegistry_Register_Good(t *core.T) {
	reg := NewRegistry(RegistryOptions{NodeID: "self"})
	rec, err := reg.Register(PeerRecord{
		Name:      "charon",
		PublicKey: make([]byte, 33),
		Addresses: []string{"10.0.0.1:14038"},
		Cap:       []Capability{CapVPN},
	}, PeerVector{VerifiedChain: true, VerifiedHNS: true, Uptime: uptimeFullSeconds, ResponseTime: 10, ConsistencyN: 10})
	core.AssertNoError(t, err)
	core.AssertInDelta(t, 1.0, rec.TrustScore, 1e-9)

	got, ok := reg.Lookup("charon")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "charon", got.Name)
	core.AssertInDelta(t, 1.0, got.TrustScore, 1e-9)
}

// TestRegistry_Register_Bad rejects an invalid record.
func TestRegistry_Register_Bad(t *core.T) {
	reg := NewRegistry(RegistryOptions{})
	_, err := reg.Register(PeerRecord{Name: ""}, PeerVector{})
	core.AssertError(t, err)
	_, ok := reg.Lookup("")
	core.AssertFalse(t, ok)
}

// TestRegistry_WithCapability_Good returns capable peers ordered by descending
// trust score (RFC §11.3 — returned, not filtered).
func TestRegistry_WithCapability_Good(t *core.T) {
	reg := NewRegistry(RegistryOptions{})
	_, _ = reg.Register(PeerRecord{Name: "low", PublicKey: make([]byte, 33), Cap: []Capability{CapVPN}}, PeerVector{VerifiedChain: true})
	_, _ = reg.Register(PeerRecord{Name: "high", PublicKey: make([]byte, 33), Cap: []Capability{CapVPN}}, PeerVector{VerifiedChain: true, VerifiedHNS: true, Uptime: uptimeFullSeconds, ResponseTime: 5, ConsistencyN: 10})
	_, _ = reg.Register(PeerRecord{Name: "dns-only", PublicKey: make([]byte, 33), Cap: []Capability{CapDNS}}, PeerVector{VerifiedChain: true})

	vpn := reg.WithCapability(CapVPN)
	core.AssertEqual(t, 2, len(vpn))
	core.AssertEqual(t, "high", vpn[0].Name) // highest trust first
	core.AssertEqual(t, "low", vpn[1].Name)
}

// TestRegistry_Announce_Good broadcasts a signed gossip delta over the transport.
func TestRegistry_Announce_Good(t *core.T) {
	priv, pub := newTestKey(t)
	ft := newFakeTransport()
	reg := NewRegistry(RegistryOptions{NodeID: "self", Transport: ft})

	gossip, err := SignGossip(priv, "charon", []NameRecord{{Name: "charon", PublicKey: pub}})
	core.AssertNoError(t, err)

	res := reg.Announce(gossip)
	core.AssertTrue(t, res.OK)
	core.AssertEqual(t, 1, len(ft.broadcasts))
	core.AssertEqual(t, GossipMessageType, ft.broadcasts[0].Type)
	core.AssertEqual(t, "self", ft.broadcasts[0].From)

	// Round-trip the broadcast payload back into a Gossip.
	decoded, err := DecodeGossipMessage(ft.broadcasts[0])
	core.AssertNoError(t, err)
	core.AssertEqual(t, "charon", decoded.Peer)
}

// TestRegistry_Announce_Bad fails on a nil gossip and is a no-op success with no
// transport configured.
func TestRegistry_Announce_Bad(t *core.T) {
	reg := NewRegistry(RegistryOptions{NodeID: "self"}) // no transport
	core.AssertFalse(t, reg.Announce(nil).OK)
	// No transport: a valid gossip announce is a no-op success.
	priv, pub := newTestKey(t)
	gossip, _ := SignGossip(priv, "charon", []NameRecord{{Name: "charon", PublicKey: pub}})
	core.AssertTrue(t, reg.Announce(gossip).OK)
}

// TestRegistry_Receive_Good verifies and merges an inbound gossip via the
// registry surface.
func TestRegistry_Receive_Good(t *core.T) {
	priv, pub := newTestKey(t)
	reg := NewRegistry(RegistryOptions{NodeID: "self"})

	gossip, err := SignGossip(priv, "charon", []NameRecord{{Name: "charon", PublicKey: pub, TrustScore: 0.8}})
	core.AssertNoError(t, err)

	res, err := reg.Receive(gossip, GossipAuthority{Peer: "charon", PublicKey: pub, Address: "1.2.3.4:14038"})
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, res.Merged)

	rec, ok := reg.Cache().Get("charon")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, SourceGossip, rec.Source)
}

// TestRegistry_Remove_Ugly removes a peer and its gossip, tolerating an absent
// peer.
func TestRegistry_Remove_Ugly(t *core.T) {
	reg := NewRegistry(RegistryOptions{})
	reg.Remove("missing") // no-op
	_, _ = reg.Register(PeerRecord{Name: "charon", PublicKey: make([]byte, 33)}, PeerVector{})
	reg.Cache().Put(NameRecord{Name: "x", Source: SourceGossip, GossipPeer: "charon"})

	reg.Remove("charon")
	_, ok := reg.Lookup("charon")
	core.AssertFalse(t, ok)
	_, ok = reg.Cache().Get("x")
	core.AssertFalse(t, ok)
}

// TestRegistry_DecodeGossipMessage_Bad rejects a wrong-typed or oversized
// transport message (RFC §12 untrusted-length guard).
func TestRegistry_DecodeGossipMessage_Bad(t *core.T) {
	_, err := DecodeGossipMessage(GossipMessage{Type: "not-gossip"})
	core.AssertError(t, err)

	// Oversized record count in the encoded payload.
	records := make([]NameRecord, MaxGossipRecords+1)
	big := &Gossip{Peer: "charon", Records: records, Signature: make([]byte, GossipSignatureLen)}
	enc := core.JSONMarshal(big)
	core.AssertTrue(t, enc.OK)
	_, err = DecodeGossipMessage(GossipMessage{Type: GossipMessageType, Payload: enc.Value.([]byte)})
	core.AssertError(t, err)
}
