// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"encoding/hex"

	core "dappco.re/go"
	netpkg "dappco.re/go/lns/pkg/net"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// exampleService returns a started LNS service with all actions registered.
func exampleService() (*core.Core, *Service) {
	c := core.New()
	svc := NewService(c)
	svc.OnStartup(context.Background())
	return c, svc
}

// ExampleService_peerRegister registers a peer over the lns.peer.register action
// and looks it back up via lns.peer.lookup.
func ExampleService_peerRegister() {
	c, _ := exampleService()
	priv, _ := secp256k1.GeneratePrivateKey()
	pubHex := hex.EncodeToString(priv.PubKey().SerializeCompressed())

	c.Action(actionNamePeerRegister).Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon"},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "cap", Value: []string{netpkg.CapVPN}},
	))

	lookup := c.Action("lns.peer.lookup").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon"},
	))
	core.Println(lookup.Value.(peerLookupResult).Found)
	// Output: true
}

// ExampleService_gossipReceive verifies and merges a signed gossip delta from a
// peer that is authoritative for the record name (RFC §12.1).
func ExampleService_gossipReceive() {
	c, _ := exampleService()
	priv, _ := secp256k1.GeneratePrivateKey()
	pub := priv.PubKey().SerializeCompressed()
	pubHex := hex.EncodeToString(pub)

	gossip, _ := netpkg.SignGossip(priv, "charon", []netpkg.NameRecord{{
		Name:      "charon",
		PublicKey: pub,
		Source:    netpkg.SourceHSD,
	}})

	result := c.Action(actionNameGossipReceive).Run(context.Background(), core.NewOptions(
		core.Option{Key: "peer", Value: "charon"},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "gossip", Value: core.JSONMarshalString(gossip)},
	))
	core.Println(result.Value.(gossipReceiveResult).Merged)
	// Output: 1
}

// ExampleService_gossipAnnounce shows the offline broadcast path: with no
// transport configured the signed delta is accepted locally but not broadcast.
func ExampleService_gossipAnnounce() {
	c, _ := exampleService()
	priv, _ := secp256k1.GeneratePrivateKey()
	pub := priv.PubKey().SerializeCompressed()

	gossip, _ := netpkg.SignGossip(priv, "charon", []netpkg.NameRecord{{
		Name:      "charon",
		PublicKey: pub,
		Source:    netpkg.SourceHSD,
	}})

	result := c.Action(actionNameGossipAnnounce).Run(context.Background(), core.NewOptions(
		core.Option{Key: "gossip", Value: core.JSONMarshalString(gossip)},
	))
	core.Println(result.Value.(gossipAnnounceResult).Broadcast)
	// Output: false
}
