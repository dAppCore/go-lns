// SPDX-License-Identifier: EUPL-1.2

package net

import (
	core "dappco.re/go"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// newTestKey returns a deterministic-enough secp256k1 keypair for tests, along
// with its 33-byte compressed public key bytes.
func newTestKey(t *core.T) (*secp256k1.PrivateKey, []byte) {
	t.Helper()
	priv, err := secp256k1.GeneratePrivateKey()
	core.AssertNoError(t, err)
	return priv, priv.PubKey().SerializeCompressed()
}

// sampleRecord builds a valid name record for "charon".
func sampleRecord(pub []byte) NameRecord {
	return NameRecord{
		Name:       "charon",
		PublicKey:  pub,
		Addresses:  []string{"10.69.69.165:14038"},
		Cap:        []Capability{CapVPN, CapDNS},
		TrustScore: 0.8,
		Source:     SourceHSD,
	}
}

// TestGossip_Merge_Good covers the happy path: a peer signs a delta for a name
// it owns, the signature verifies, and the record merges into the cache with the
// gossip flag, origin peer/address, and the trust penalty applied (RFC §12.1).
func TestGossip_Merge_Good(t *core.T) {
	priv, pub := newTestKey(t)
	cache := NewCache()

	gossip, err := SignGossip(priv, "charon", []NameRecord{sampleRecord(pub)})
	core.AssertNoError(t, err)

	auth := GossipAuthority{Peer: "charon", PublicKey: pub, Address: "10.69.69.165:14038"}
	res, err := gossip.Merge(cache, auth)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, res.Merged)
	core.AssertEqual(t, 0, res.RejectedUnauthorised)

	merged, ok := cache.Get("charon")
	core.AssertTrue(t, ok)
	core.AssertTrue(t, merged.IsGossip())
	core.AssertEqual(t, SourceGossip, merged.Source)
	core.AssertEqual(t, "charon", merged.GossipPeer)
	core.AssertEqual(t, "10.69.69.165:14038", merged.GossipPeerAddress)
	// 0.8 * 0.9 = 0.72 after the gossip penalty.
	core.AssertInDelta(t, 0.72, merged.TrustScore, 1e-9)
}

// TestGossip_Merge_Bad rejects a message whose signature does not verify against
// the authority's public key (RFC §12.1 step 1). A tampered signature must not
// merge anything.
func TestGossip_Merge_Bad(t *core.T) {
	priv, pub := newTestKey(t)
	cache := NewCache()

	gossip, err := SignGossip(priv, "charon", []NameRecord{sampleRecord(pub)})
	core.AssertNoError(t, err)

	// Forge: flip a byte of the signature.
	gossip.Signature[0] ^= 0xff

	auth := GossipAuthority{Peer: "charon", PublicKey: pub}
	res, err := gossip.Merge(cache, auth)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, res.Merged)
	core.AssertEqual(t, 0, cache.Len())
}

// TestGossip_Merge_Bad_WrongKey rejects a validly-signed message when verified
// against a different peer's key — a peer cannot impersonate another's records.
func TestGossip_Merge_Bad_WrongKey(t *core.T) {
	priv, _ := newTestKey(t)
	_, otherPub := newTestKey(t)
	cache := NewCache()

	gossip, err := SignGossip(priv, "charon", []NameRecord{sampleRecord(otherPub)})
	core.AssertNoError(t, err)

	auth := GossipAuthority{Peer: "charon", PublicKey: otherPub}
	_, err = gossip.Merge(cache, auth)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, cache.Len())
}

// TestGossip_Merge_Ugly rejects records the peer is not authoritative for
// (RFC §12.1 step 2): a peer signing a name that is neither its own nor in its
// authorised set has that record dropped, while a record it does own merges.
func TestGossip_Merge_Ugly(t *core.T) {
	priv, pub := newTestKey(t)
	cache := NewCache()

	own := NameRecord{Name: "charon", PublicKey: pub, TrustScore: 0.5}
	foreign := NameRecord{Name: "someone-else", PublicKey: pub, TrustScore: 0.9}

	gossip, err := SignGossip(priv, "charon", []NameRecord{own, foreign})
	core.AssertNoError(t, err)

	auth := GossipAuthority{Peer: "charon", PublicKey: pub}
	res, err := gossip.Merge(cache, auth)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, res.Merged)
	core.AssertEqual(t, 1, res.RejectedUnauthorised)

	_, ok := cache.Get("charon")
	core.AssertTrue(t, ok)
	_, ok = cache.Get("someone-else")
	core.AssertFalse(t, ok)
}

// TestGossip_Merge_AuthorizedNames accepts a non-own name the peer has observed
// from HSD and listed in AuthorizedNames (RFC §12.1 step 2).
func TestGossip_Merge_AuthorizedNames(t *core.T) {
	priv, pub := newTestKey(t)
	cache := NewCache()

	relayed := NameRecord{Name: "trade", PublicKey: pub, TrustScore: 0.6}
	gossip, err := SignGossip(priv, "charon", []NameRecord{relayed})
	core.AssertNoError(t, err)

	auth := GossipAuthority{Peer: "charon", PublicKey: pub, AuthorizedNames: []string{"trade"}}
	res, err := gossip.Merge(cache, auth)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, res.Merged)

	rec, ok := cache.Get("trade")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "charon", rec.GossipPeer)
}

// TestGossip_Merge_HSDPrecedence confirms gossip never overwrites an
// authoritative HSD-sourced record for the same name (RFC §12.3).
func TestGossip_Merge_HSDPrecedence(t *core.T) {
	priv, pub := newTestKey(t)
	cache := NewCache()

	// Pre-seed an HSD record.
	cache.Put(NameRecord{Name: "charon", PublicKey: pub, TrustScore: 0.95, Source: SourceHSD})

	gossip, err := SignGossip(priv, "charon", []NameRecord{sampleRecord(pub)})
	core.AssertNoError(t, err)

	auth := GossipAuthority{Peer: "charon", PublicKey: pub}
	res, err := gossip.Merge(cache, auth)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 0, res.Merged)
	core.AssertEqual(t, 1, res.SkippedHSD)

	rec, _ := cache.Get("charon")
	core.AssertEqual(t, SourceHSD, rec.Source)
	core.AssertInDelta(t, 0.95, rec.TrustScore, 1e-9)
}

// TestGossip_Merge_Oversized rejects a message whose record count exceeds the
// wire bound, without touching the cache (RFC §12 untrusted-length guard).
func TestGossip_Merge_Oversized(t *core.T) {
	_, pub := newTestKey(t)
	cache := NewCache()

	records := make([]NameRecord, MaxGossipRecords+1)
	for i := range records {
		records[i] = NameRecord{Name: "n", PublicKey: pub}
	}
	gossip := &Gossip{Peer: "charon", Records: records, Signature: make([]byte, GossipSignatureLen)}

	auth := GossipAuthority{Peer: "charon", PublicKey: pub}
	res, err := gossip.Merge(cache, auth)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, res.Merged)
	core.AssertEqual(t, 0, cache.Len())
}

// TestGossip_Merge_OversizedRecord drops a single record whose field exceeds the
// bounded limit while merging the valid records around it.
func TestGossip_Merge_OversizedRecord(t *core.T) {
	priv, pub := newTestKey(t)
	cache := NewCache()

	tooLong := make([]byte, 0, MaxNameLen+10)
	for i := 0; i < MaxNameLen+5; i++ {
		tooLong = append(tooLong, 'a')
	}
	bad := NameRecord{Name: string(tooLong), PublicKey: pub}
	// Authorise the over-long name so the rejection is on bounds, not authority.
	gossip, err := SignGossip(priv, "charon", []NameRecord{bad})
	core.AssertNoError(t, err)

	auth := GossipAuthority{Peer: "charon", PublicKey: pub, AuthorizedNames: []string{string(tooLong)}}
	res, err := gossip.Merge(cache, auth)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 0, res.Merged)
	core.AssertEqual(t, 1, res.RejectedInvalid)
}

// TestGossip_SignGossip_Bad rejects signing with a nil key or an oversized batch.
func TestGossip_SignGossip_Bad(t *core.T) {
	_, err := SignGossip(nil, "charon", nil)
	core.AssertError(t, err)

	priv, pub := newTestKey(t)
	records := make([]NameRecord, MaxGossipRecords+1)
	for i := range records {
		records[i] = NameRecord{Name: "n", PublicKey: pub}
	}
	_, err = SignGossip(priv, "charon", records)
	core.AssertError(t, err)

	_, err = SignGossip(priv, "", nil)
	core.AssertError(t, err)
}

// TestGossip_VerifySignature_Ugly rejects a malformed signature length and a
// malformed public key without panicking.
func TestGossip_VerifySignature_Ugly(t *core.T) {
	priv, pub := newTestKey(t)
	gossip, err := SignGossip(priv, "charon", []NameRecord{{Name: "charon", PublicKey: pub}})
	core.AssertNoError(t, err)

	// Wrong-length signature.
	short := &Gossip{Peer: gossip.Peer, Records: gossip.Records, Signature: []byte{0x01, 0x02}}
	core.AssertError(t, short.VerifySignature(pub))

	// Malformed public key.
	core.AssertError(t, gossip.VerifySignature([]byte{0x00}))

	// Nil receiver.
	var nilGossip *Gossip
	core.AssertError(t, nilGossip.VerifySignature(pub))
}

// TestGossip_PeerMismatch rejects a message whose peer label does not match the
// authority being checked.
func TestGossip_PeerMismatch(t *core.T) {
	priv, pub := newTestKey(t)
	cache := NewCache()
	gossip, err := SignGossip(priv, "charon", []NameRecord{{Name: "charon", PublicKey: pub}})
	core.AssertNoError(t, err)

	auth := GossipAuthority{Peer: "different", PublicKey: pub}
	_, err = gossip.Merge(cache, auth)
	core.AssertError(t, err)
}
