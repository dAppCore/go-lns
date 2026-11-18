// SPDX-License-Identifier: EUPL-1.2

// Package net is the LNS peer-networking layer (actions prefix lns.peer.*).
//
// It builds the Lethean Name System peer and gossip-replication layer on top
// of the go-p2p transport (github.com/dappcore/go-p2p/go), the Poindexter trust scorer
// (forge.lthn.ai/Snider/Poindexter), and the LNS chain primitives in
// github.com/dappcore/go-lns/go/pkg/primitives. It does NOT reimplement transport, peer
// discovery scoring, or identity — those come from go-p2p.
//
// The three responsibilities are:
//
//   - PeerRecord — the resolvable record for a peer name (RFC §10.1): a bare
//     label, the peer's chain public key, network addresses, capabilities, and
//     a Poindexter-derived trust score.
//
//   - Registry — a registry/discovery surface over go-p2p's node.PeerRegistry.
//     Peers are added, looked up, and ranked; trust scores are recomputed on
//     write from a pluggable ScorePeer hook (Poindexter by default).
//
//   - Gossip — record replication (RFC §12). On receive the layer verifies the
//     gossip message is signed by the originating peer's chain key, accepts
//     only records the peer is authoritative for, and merges them into a local
//     Cache with a gossip=true flag, the origin peer address, and a trust
//     penalty applied relative to HSD-sourced records.
//
// Build a peer record, score it, and merge a signed gossip delta:
//
//	rec := net.PeerRecord{Name: "charon", PublicKey: pub, Addresses: []string{"10.0.0.1:14038"}}
//	rec.TrustScore = net.DefaultScorer{}.ScorePeer(net.PeerVector{VerifiedChain: true})
//	cache := net.NewCache()
//	gossip, _ := net.SignGossip(priv, "charon", []net.NameRecord{rec.AsNameRecord()})
//	res := gossip.Merge(cache, net.GossipAuthority{Peer: "charon", PublicKey: pub})
package net
