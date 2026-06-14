// SPDX-License-Identifier: EUPL-1.2

package net_test

import (
	core "dappco.re/go"
	"dappco.re/go/lns/pkg/net"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// exampleKey returns a fixed secp256k1 key so example output is deterministic.
func exampleKey() (*secp256k1.PrivateKey, []byte) {
	var b [32]byte
	b[31] = 0x2a // fixed non-zero scalar
	priv := secp256k1.PrivKeyFromBytes(b[:])
	return priv, priv.PubKey().SerializeCompressed()
}

// ExampleSignGossip_merge shows the full sign → verify → merge flow: a peer
// signs a record it owns and the record lands in the cache with the gossip flag.
func ExampleSignGossip_merge() {
	priv, pub := exampleKey()
	cache := net.NewCache()

	gossip, _ := net.SignGossip(priv, "charon", []net.NameRecord{{
		Name:       "charon",
		PublicKey:  pub,
		TrustScore: 0.8,
	}})

	auth := net.GossipAuthority{Peer: "charon", PublicKey: pub, Address: "10.0.0.1:14038"}
	res, _ := gossip.Merge(cache, auth)

	rec, _ := cache.Get("charon")
	core.Println(core.Sprintf("merged=%d source=%s peer=%s", res.Merged, rec.Source, rec.GossipPeer))
	// Output: merged=1 source=gossip peer=charon
}

// ExampleGossip_Merge_unauthorised shows a peer's foreign record being rejected.
func ExampleGossip_Merge_unauthorised() {
	priv, pub := exampleKey()
	cache := net.NewCache()

	gossip, _ := net.SignGossip(priv, "charon", []net.NameRecord{{Name: "someone-else", PublicKey: pub}})
	res, _ := gossip.Merge(cache, net.GossipAuthority{Peer: "charon", PublicKey: pub})

	core.Println(core.Sprintf("merged=%d rejected=%d", res.Merged, res.RejectedUnauthorised))
	// Output: merged=0 rejected=1
}

// ExampleDefaultScorer_ScorePeer shows the Poindexter-backed trust score for a
// fully-verified peer.
func ExampleDefaultScorer_ScorePeer() {
	score := net.DefaultScorer{}.ScorePeer(net.PeerVector{
		VerifiedChain: true,
		VerifiedHNS:   true,
		Uptime:        86400,
		ResponseTime:  10,
		ConsistencyN:  10,
	})
	core.Println(core.Sprintf("score=%.2f", score))
	// Output: score=1.00
}

// ExampleRegistry_WithCapability shows capable peers returned ranked by trust.
func ExampleRegistry_WithCapability() {
	reg := net.NewRegistry(net.RegistryOptions{})
	_, _ = reg.Register(net.PeerRecord{Name: "vpn-node", PublicKey: make([]byte, 33), Cap: []net.Capability{net.CapVPN}},
		net.PeerVector{VerifiedChain: true})
	candidates := reg.WithCapability(net.CapVPN)
	core.Println(core.Sprintf("vpn peers=%d first=%s", len(candidates), candidates[0].Name))
	// Output: vpn peers=1 first=vpn-node
}
