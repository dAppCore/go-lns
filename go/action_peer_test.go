// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"encoding/hex"
	"testing"

	core "github.com/dappcore/go"
	netpkg "github.com/dappcore/go-lns/go/pkg/net"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

const (
	testPeerCharon   = "charon"
	testPeerAddr     = "10.69.69.165:14038"
	testGossipRecvOK = "lns.gossip.receive returned %#v, want OK"
)

// peerTestService spins up an LNS service with all actions registered.
func peerTestService(t *testing.T) (*core.Core, *Service) {
	t.Helper()
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(t.Context()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}
	return c, svc
}

// peerTestKey returns a fresh secp256k1 keypair plus its hex-encoded compressed
// public key, the form the peer and gossip actions accept over options.
func peerTestKey(t *testing.T) (*secp256k1.PrivateKey, []byte, string) {
	t.Helper()
	priv, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		t.Fatalf("GeneratePrivateKey returned error: %v", err)
	}
	pub := priv.PubKey().SerializeCompressed()
	return priv, pub, hex.EncodeToString(pub)
}

// signedGossipJSON signs a single-record gossip delta for the named peer and
// returns its JSON encoding, the form the gossip actions accept over options.
func signedGossipJSON(t *testing.T, priv *secp256k1.PrivateKey, pub []byte, peer, recordName string) string {
	t.Helper()
	record := netpkg.NameRecord{
		Name:       recordName,
		PublicKey:  pub,
		Addresses:  []string{testPeerAddr},
		Cap:        []netpkg.Capability{netpkg.CapVPN, netpkg.CapDNS},
		TrustScore: 0.8,
		Source:     netpkg.SourceHSD,
	}
	gossip, err := netpkg.SignGossip(priv, peer, []netpkg.NameRecord{record})
	if err != nil {
		t.Fatalf("SignGossip returned error: %v", err)
	}
	return core.JSONMarshalString(gossip)
}

func TestServiceRegistersPeerAndGossipActions_Good(t *testing.T) {
	c, _ := peerTestService(t)

	for _, name := range []string{
		actionNamePeerRegister,
		"lns.peer.lookup",
		"lns.peer.list",
		actionNameGossipReceive,
		actionNameGossipAnnounce,
	} {
		if !c.Action(name).Exists() {
			t.Fatalf("expected action %q to be registered", name)
		}
	}
}

func TestPeerRegisterLookupList_Good(t *testing.T) {
	c, _ := peerTestService(t)
	_, _, pubHex := peerTestKey(t)

	registerResult := c.Action(actionNamePeerRegister).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "addresses", Value: []string{testPeerAddr}},
		core.Option{Key: "cap", Value: []string{netpkg.CapVPN}},
		core.Option{Key: "verified_chain", Value: true},
		core.Option{Key: "verified_hns", Value: true},
	))
	if !registerResult.OK {
		t.Fatalf("lns.peer.register returned %#v, want OK", registerResult)
	}
	registered, ok := registerResult.Value.(peerRecordResult)
	if !ok {
		t.Fatalf("lns.peer.register returned %T, want peerRecordResult", registerResult.Value)
	}
	if registered.Name != testPeerCharon || registered.TrustScore <= 0 {
		t.Fatalf("registered peer = %#v, want scored charon", registered)
	}

	lookupResult := c.Action("lns.peer.lookup").Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testPeerCharon},
	))
	if !lookupResult.OK {
		t.Fatalf("lns.peer.lookup returned %#v, want OK", lookupResult)
	}
	lookup, ok := lookupResult.Value.(peerLookupResult)
	if !ok {
		t.Fatalf("lns.peer.lookup returned %T, want peerLookupResult", lookupResult.Value)
	}
	if !lookup.Found || lookup.Peer == nil || lookup.Peer.Name != testPeerCharon {
		t.Fatalf("peer lookup = %#v, want found charon", lookup)
	}

	listResult := c.Action("lns.peer.list").Run(t.Context(), core.NewOptions(
		core.Option{Key: "cap", Value: netpkg.CapVPN},
	))
	if !listResult.OK {
		t.Fatalf("lns.peer.list returned %#v, want OK", listResult)
	}
	list, ok := listResult.Value.(peerListResult)
	if !ok {
		t.Fatalf("lns.peer.list returned %T, want peerListResult", listResult.Value)
	}
	if list.Count != 1 || len(list.Peers) != 1 || list.Peers[0].Name != testPeerCharon {
		t.Fatalf("peer list = %#v, want one charon with vpn cap", list)
	}
}

func TestPeerLookupMissing_Good(t *testing.T) {
	c, _ := peerTestService(t)

	result := c.Action("lns.peer.lookup").Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: "ghost"},
	))
	if !result.OK {
		t.Fatalf("lns.peer.lookup returned %#v, want OK", result)
	}
	lookup, ok := result.Value.(peerLookupResult)
	if !ok {
		t.Fatalf("lns.peer.lookup returned %T, want peerLookupResult", result.Value)
	}
	if lookup.Found || lookup.Peer != nil {
		t.Fatalf("missing lookup = %#v, want not found", lookup)
	}
}

func TestPeerRegisterRequiresKey_Bad(t *testing.T) {
	c, _ := peerTestService(t)

	result := c.Action(actionNamePeerRegister).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testPeerCharon},
	))
	if result.OK {
		t.Fatalf("lns.peer.register without public_key returned OK, want failure")
	}
}

func TestGossipReceiveMergesSignedDelta_Good(t *testing.T) {
	c, _ := peerTestService(t)
	priv, _, pubHex := peerTestKey(t)
	pub, _ := hex.DecodeString(pubHex)
	gossipJSON := signedGossipJSON(t, priv, pub, testPeerCharon, testPeerCharon)

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "address", Value: testPeerAddr},
		core.Option{Key: "gossip", Value: gossipJSON},
	))
	if !result.OK {
		t.Fatalf(testGossipRecvOK, result)
	}
	merge, ok := result.Value.(gossipReceiveResult)
	if !ok {
		t.Fatalf("lns.gossip.receive returned %T, want gossipReceiveResult", result.Value)
	}
	if merge.Merged != 1 || merge.CacheSize != 1 || merge.Peer != testPeerCharon {
		t.Fatalf("gossip receive = %#v, want one merged record", merge)
	}
}

func TestGossipReceiveViaMessageEnvelope_Good(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, pubHex := peerTestKey(t)
	gossip, err := netpkg.SignGossip(priv, testPeerCharon, []netpkg.NameRecord{{
		Name:       testPeerCharon,
		PublicKey:  pub,
		Addresses:  []string{testPeerAddr},
		Source:     netpkg.SourceHSD,
		TrustScore: 0.9,
	}})
	if err != nil {
		t.Fatalf("SignGossip returned error: %v", err)
	}
	payload := core.JSONMarshal(gossip)
	if !payload.OK {
		t.Fatalf("marshal gossip payload failed: %#v", payload)
	}
	rawPayload, _ := payload.Value.([]byte)

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "message", Value: map[string]any{
			"type":    netpkg.GossipMessageType,
			"from":    testPeerCharon,
			"payload": rawPayload,
		}},
	))
	if !result.OK {
		t.Fatalf(testGossipRecvOK, result)
	}
	merge, ok := result.Value.(gossipReceiveResult)
	if !ok {
		t.Fatalf("lns.gossip.receive returned %T, want gossipReceiveResult", result.Value)
	}
	if merge.Merged != 1 {
		t.Fatalf("gossip receive via envelope = %#v, want one merged record", merge)
	}
}

func TestGossipReceiveRejectsForgedSignature_Bad(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, pubHex := peerTestKey(t)
	gossip, err := netpkg.SignGossip(priv, testPeerCharon, []netpkg.NameRecord{{
		Name:      testPeerCharon,
		PublicKey: pub,
		Source:    netpkg.SourceHSD,
	}})
	if err != nil {
		t.Fatalf("SignGossip returned error: %v", err)
	}
	gossip.Signature[0] ^= 0xff
	forged := core.JSONMarshalString(gossip)

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "gossip", Value: forged},
	))
	if result.OK {
		t.Fatalf("lns.gossip.receive with forged signature returned OK, want failure")
	}
}

func TestGossipReceiveRequiresPublicKey_Bad(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, _ := peerTestKey(t)
	gossipJSON := signedGossipJSON(t, priv, pub, testPeerCharon, testPeerCharon)

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "gossip", Value: gossipJSON},
	))
	if result.OK {
		t.Fatalf("lns.gossip.receive without public_key returned OK, want failure")
	}
}

func TestGossipReceiveDropsUnauthorisedRecord_Ugly(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, pubHex := peerTestKey(t)
	// The peer signs a record for a name it is NOT authoritative for: its own
	// label is "charon" but the record names "stranger" (RFC §12.1 step 2).
	gossipJSON := signedGossipJSON(t, priv, pub, testPeerCharon, "stranger")

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "gossip", Value: gossipJSON},
	))
	if !result.OK {
		t.Fatalf(testGossipRecvOK, result)
	}
	merge, ok := result.Value.(gossipReceiveResult)
	if !ok {
		t.Fatalf("lns.gossip.receive returned %T, want gossipReceiveResult", result.Value)
	}
	if merge.Merged != 0 || merge.RejectedUnauthorised != 1 || merge.CacheSize != 0 {
		t.Fatalf("gossip receive = %#v, want one unauthorised rejection", merge)
	}
}

func TestGossipReceiveAllowsAuthorisedName_Good(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, pubHex := peerTestKey(t)
	// The peer relays a record for "stranger" but is explicitly authorised for
	// that name via authorized_names (RFC §12.1 step 2).
	gossipJSON := signedGossipJSON(t, priv, pub, testPeerCharon, "stranger")

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "authorized_names", Value: []string{"stranger"}},
		core.Option{Key: "gossip", Value: gossipJSON},
	))
	if !result.OK {
		t.Fatalf(testGossipRecvOK, result)
	}
	merge, ok := result.Value.(gossipReceiveResult)
	if !ok {
		t.Fatalf("lns.gossip.receive returned %T, want gossipReceiveResult", result.Value)
	}
	if merge.Merged != 1 {
		t.Fatalf("gossip receive = %#v, want authorised record merged", merge)
	}
}

func TestGossipReceiveRequiresPayload_Bad(t *testing.T) {
	c, _ := peerTestService(t)
	_, _, pubHex := peerTestKey(t)

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
	))
	if result.OK {
		t.Fatalf("lns.gossip.receive without gossip payload returned OK, want failure")
	}
}

func TestGossipReceiveViaPayloadBytes_Good(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, pubHex := peerTestKey(t)
	gossip, err := netpkg.SignGossip(priv, testPeerCharon, []netpkg.NameRecord{{
		Name:      testPeerCharon,
		PublicKey: pub,
		Source:    netpkg.SourceHSD,
	}})
	if err != nil {
		t.Fatalf("SignGossip returned error: %v", err)
	}
	encoded := core.JSONMarshal(gossip)
	if !encoded.OK {
		t.Fatalf("marshal gossip failed: %#v", encoded)
	}
	raw, _ := encoded.Value.([]byte)

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "payload", Value: raw},
	))
	if !result.OK {
		t.Fatalf(testGossipRecvOK, result)
	}
	if result.Value.(gossipReceiveResult).Merged != 1 {
		t.Fatalf("payload-bytes gossip receive = %#v, want one merged", result.Value)
	}
}

// TestPeerRegisterScoresFromStringVector_Good exercises the float and string
// coercion in the trust vector parser: options arriving as strings (the common
// shape over MCP/CLI) still produce a scored peer.
func TestPeerRegisterScoresFromStringVector_Good(t *testing.T) {
	c, _ := peerTestService(t)
	_, _, pubHex := peerTestKey(t)

	result := c.Action(actionNamePeerRegister).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: "scored"},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "uptime", Value: "86400"},
		core.Option{Key: "response_time", Value: "15.0"},
		core.Option{Key: "verified_chain", Value: "true"},
		core.Option{Key: "consistency", Value: "10"},
	))
	if !result.OK {
		t.Fatalf("lns.peer.register returned %#v, want OK", result)
	}
	registered, ok := result.Value.(peerRecordResult)
	if !ok {
		t.Fatalf("lns.peer.register returned %T, want peerRecordResult", result.Value)
	}
	if registered.TrustScore <= 0 {
		t.Fatalf("string-vector register = %#v, want positive trust score", registered)
	}
}

// TestGossipReceiveViaGossipObject_Good passes a decoded net.Gossip value
// directly (the in-process call shape) rather than its JSON encoding.
func TestGossipReceiveViaGossipObject_Good(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, pubHex := peerTestKey(t)
	gossip, err := netpkg.SignGossip(priv, testPeerCharon, []netpkg.NameRecord{{
		Name:      testPeerCharon,
		PublicKey: pub,
		Source:    netpkg.SourceHSD,
	}})
	if err != nil {
		t.Fatalf("SignGossip returned error: %v", err)
	}

	result := c.Action(actionNameGossipReceive).Run(t.Context(), core.NewOptions(
		core.Option{Key: "peer", Value: testPeerCharon},
		core.Option{Key: "public_key", Value: pubHex},
		core.Option{Key: "gossip", Value: gossip},
	))
	if !result.OK {
		t.Fatalf(testGossipRecvOK, result)
	}
	if result.Value.(gossipReceiveResult).Merged != 1 {
		t.Fatalf("gossip-object receive = %#v, want one merged", result.Value)
	}
}

func TestGossipAnnounceOfflineIsHonest_Ugly(t *testing.T) {
	c, _ := peerTestService(t)
	priv, pub, _ := peerTestKey(t)
	gossipJSON := signedGossipJSON(t, priv, pub, testPeerCharon, testPeerCharon)

	result := c.Action(actionNameGossipAnnounce).Run(t.Context(), core.NewOptions(
		core.Option{Key: "gossip", Value: gossipJSON},
	))
	if !result.OK {
		t.Fatalf("lns.gossip.announce returned %#v, want OK", result)
	}
	announce, ok := result.Value.(gossipAnnounceResult)
	if !ok {
		t.Fatalf("lns.gossip.announce returned %T, want gossipAnnounceResult", result.Value)
	}
	if announce.Broadcast {
		t.Fatalf("offline announce = %#v, want Broadcast=false", announce)
	}
	if len(announce.Notes) == 0 {
		t.Fatalf("offline announce = %#v, want an explanatory note", announce)
	}
}
