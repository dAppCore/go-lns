// SPDX-License-Identifier: EUPL-1.2

package net

import (
	"crypto/sha256"
	"slices"

	core "dappco.re/go"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
)

// Gossip wire bounds. They cap untrusted-length fields so a hostile peer cannot
// drive an unbounded allocation when a gossip message is decoded.
const (
	// MaxGossipRecords is the most records accepted in one gossip message.
	MaxGossipRecords = 256
	// GossipSignatureLen is the fixed length of the compact secp256k1
	// signature carried by a gossip message (64-byte R‖S).
	GossipSignatureLen = 64
)

// Gossip is a signed batch of name records replicated from one peer (RFC §12.1).
// Peer is the bare label of the originating peer; Records is the delta being
// announced; Signature is a compact secp256k1 signature over the canonical
// digest of (Peer, Records), produced with the peer's chain key.
//
//	g, _ := net.SignGossip(priv, "charon", []net.NameRecord{rec})
//	if err := g.VerifySignature(pub); err != nil { /* reject */ }
type Gossip struct {
	// Peer is the bare label of the announcing peer.
	Peer string `json:"peer"`
	// Records is the batch of name records being replicated.
	Records []NameRecord `json:"records"`
	// Signature is the 64-byte compact secp256k1 signature over the digest.
	Signature []byte `json:"signature"`
}

// GossipAuthority describes the peer a gossip message is being checked against:
// its bare label, its verified chain public key, and the network address the
// message arrived from. The merge path stamps Address onto accepted records.
//
//	auth := net.GossipAuthority{Peer: "charon", PublicKey: pub, Address: "10.0.0.1:14038"}
type GossipAuthority struct {
	// Peer is the bare label the peer is authoritative for (its own name).
	Peer string
	// PublicKey is the peer's chain public key used to verify the signature.
	PublicKey []byte
	// Address is the network address the gossip arrived from (stamped on merge).
	Address string
	// AuthorizedNames are additional names the peer has observed from HSD and is
	// therefore allowed to relay (RFC §12.1 step 2). The peer's own name is
	// always authorised regardless of this list.
	AuthorizedNames []string
}

// isAuthoritative reports whether the authority may relay a record for the given
// name: a peer is always authoritative for its own name, plus any name it has
// observed from HSD and listed in AuthorizedNames (RFC §12.1 step 2).
func (a GossipAuthority) isAuthoritative(name string) bool {
	if name == a.Peer {
		return true
	}
	return slices.Contains(a.AuthorizedNames, name)
}

// gossipDigest computes the canonical signing digest for a gossip message: the
// SHA-256 over peer-label bytes followed by each record's stable field bytes.
// Signing and verification both derive the digest the same way so the wire
// Signature is reproducible from the message contents alone.
func gossipDigest(peer string, records []NameRecord) []byte {
	h := sha256.New()
	h.Write([]byte(peer))
	h.Write([]byte{0})
	for _, rec := range records {
		h.Write([]byte(rec.Name))
		h.Write([]byte{0})
		h.Write(rec.PublicKey)
		h.Write([]byte{0})
		for _, addr := range rec.Addresses {
			h.Write([]byte(addr))
			h.Write([]byte{0})
		}
		h.Write([]byte{1})
		for _, c := range rec.Cap {
			h.Write([]byte(c))
			h.Write([]byte{0})
		}
		h.Write([]byte{1})
	}
	return h.Sum(nil)
}

// SignGossip produces a signed gossip message for a batch of records, signing
// the canonical digest with the peer's chain private key (compact secp256k1
// R‖S). It fails on a nil key, an empty peer label, or a batch that exceeds the
// wire bound.
//
//	g, err := net.SignGossip(priv, "charon", []net.NameRecord{rec})
func SignGossip(priv *secp256k1.PrivateKey, peer string, records []NameRecord) (*Gossip, error) {
	if priv == nil {
		return nil, core.E("net.SignGossip", "nil private key", nil)
	}
	if peer == "" {
		return nil, core.E("net.SignGossip", "peer label is required", nil)
	}
	if len(records) > MaxGossipRecords {
		return nil, core.Errorf("net.SignGossip: record count %d exceeds %d", len(records), MaxGossipRecords)
	}

	digest := gossipDigest(peer, records)
	sig := ecdsa.Sign(priv, digest)

	cloned := make([]NameRecord, len(records))
	for i := range records {
		cloned[i] = records[i].Clone()
	}

	return &Gossip{
		Peer:      peer,
		Records:   cloned,
		Signature: encodeCompactSig(sig),
	}, nil
}

// VerifySignature checks that the gossip message was signed by the holder of the
// supplied chain public key (RFC §12.1 step 1). It fails on a malformed key, a
// signature of the wrong length, or a signature that does not verify.
//
//	if err := g.VerifySignature(peerPubKey); err != nil { /* reject */ }
func (g *Gossip) VerifySignature(pub []byte) error {
	if g == nil {
		return core.E("net.Gossip.VerifySignature", "nil gossip", nil)
	}
	if len(g.Signature) != GossipSignatureLen {
		return core.Errorf("net.Gossip.VerifySignature: signature length %d is not %d", len(g.Signature), GossipSignatureLen)
	}
	pk, err := secp256k1.ParsePubKey(pub)
	if err != nil {
		return core.E("net.Gossip.VerifySignature", "parse public key", err)
	}
	parsed, err := decodeCompactSig(g.Signature)
	if err != nil {
		return core.E("net.Gossip.VerifySignature", "decode signature", err)
	}
	digest := gossipDigest(g.Peer, g.Records)
	if !parsed.Verify(digest, pk) {
		return core.E("net.Gossip.VerifySignature", "signature does not verify", nil)
	}
	return nil
}

// MergeResult reports the outcome of merging a gossip message into a cache.
type MergeResult struct {
	// Merged is the number of records accepted and written to the cache.
	Merged int
	// RejectedUnauthorised is the number of records dropped because the peer was
	// not authoritative for the name (RFC §12.1 step 2).
	RejectedUnauthorised int
	// RejectedInvalid is the number of records dropped because they failed
	// bounded validation.
	RejectedInvalid int
	// SkippedHSD is the number of records skipped because an authoritative
	// HSD-sourced record already holds the name (RFC §12.3 precedence).
	SkippedHSD int
}

// Merge verifies and merges the gossip message into the cache (RFC §12.1). It:
//
//  1. verifies the message is signed by the authority's chain key;
//  2. accepts only records the authority is authoritative for;
//  3. merges each accepted record with Source=gossip, the origin peer label and
//     address, and a 10% trust penalty applied to its score;
//  4. never overwrites an authoritative HSD-sourced record for the same name.
//
// A signature, authority, or bound failure rejects the whole message (step 1/2
// are message-level); per-record authority and validation failures drop only
// the offending record and are counted in the returned MergeResult.
//
//	res, err := g.Merge(cache, auth)
func (g *Gossip) Merge(cache *Cache, auth GossipAuthority) (MergeResult, error) {
	if g == nil {
		return MergeResult{}, core.E("net.Gossip.Merge", "nil gossip", nil)
	}
	if cache == nil {
		return MergeResult{}, core.E("net.Gossip.Merge", "nil cache", nil)
	}
	if len(g.Records) > MaxGossipRecords {
		return MergeResult{}, core.Errorf("net.Gossip.Merge: record count %d exceeds %d", len(g.Records), MaxGossipRecords)
	}
	if g.Peer != auth.Peer {
		return MergeResult{}, core.E("net.Gossip.Merge", "gossip peer does not match authority", nil)
	}
	if err := g.VerifySignature(auth.PublicKey); err != nil {
		return MergeResult{}, err
	}

	var result MergeResult
	for _, rec := range g.Records {
		if !auth.isAuthoritative(rec.Name) {
			result.RejectedUnauthorised++
			continue
		}
		if err := rec.validateForMerge(); err != nil {
			result.RejectedInvalid++
			continue
		}
		// HSD-sourced records are authoritative and are never displaced by
		// gossip (RFC §12.3). A gossip record only refreshes another gossip
		// record or fills an empty slot.
		if existing, ok := cache.Get(rec.Name); ok && existing.Source != SourceGossip {
			result.SkippedHSD++
			continue
		}

		merged := rec.Clone()
		merged.Source = SourceGossip
		merged.GossipPeer = auth.Peer
		merged.GossipPeerAddress = auth.Address
		merged.TrustScore = ApplyGossipPenalty(rec.TrustScore)
		cache.Put(merged)
		result.Merged++
	}
	return result, nil
}

// validateForMerge enforces the bounded limits on a record arriving over gossip,
// reusing the PeerRecord bounds. An over-long name, address, or capability is
// rejected so untrusted gossip cannot smuggle oversized data into the cache.
func (r NameRecord) validateForMerge() error {
	if r.Name == "" {
		return core.E("net.NameRecord.validateForMerge", "name is required", nil)
	}
	if len(r.Name) > MaxNameLen {
		return core.Errorf("net.NameRecord.validateForMerge: name length %d exceeds %d", len(r.Name), MaxNameLen)
	}
	if len(r.PublicKey) > MaxPublicKeyLen {
		return core.Errorf("net.NameRecord.validateForMerge: public key length %d exceeds %d", len(r.PublicKey), MaxPublicKeyLen)
	}
	if len(r.Addresses) > MaxAddresses {
		return core.Errorf("net.NameRecord.validateForMerge: address count %d exceeds %d", len(r.Addresses), MaxAddresses)
	}
	for _, addr := range r.Addresses {
		if len(addr) > MaxAddressLen {
			return core.Errorf("net.NameRecord.validateForMerge: address length %d exceeds %d", len(addr), MaxAddressLen)
		}
	}
	if len(r.Cap) > MaxCaps {
		return core.Errorf("net.NameRecord.validateForMerge: capability count %d exceeds %d", len(r.Cap), MaxCaps)
	}
	for _, c := range r.Cap {
		if len(c) > MaxCapLen {
			return core.Errorf("net.NameRecord.validateForMerge: capability length %d exceeds %d", len(c), MaxCapLen)
		}
	}
	return nil
}

// encodeCompactSig serialises an ECDSA signature as the 64-byte R‖S form used
// across the LNS chain (big-endian 32-byte R, big-endian 32-byte S).
func encodeCompactSig(sig *ecdsa.Signature) []byte {
	r := sig.R()
	s := sig.S()
	rb := r.Bytes()
	sb := s.Bytes()
	out := make([]byte, 0, GossipSignatureLen)
	out = append(out, rb[:]...)
	out = append(out, sb[:]...)
	return out
}

// decodeCompactSig parses a 64-byte R‖S signature into an ecdsa.Signature. It
// rejects a wrong-length body, an R or S that overflows the group order, and a
// zero R or S.
func decodeCompactSig(body []byte) (*ecdsa.Signature, error) {
	if len(body) != GossipSignatureLen {
		return nil, core.Errorf("net.decodeCompactSig: signature must be %d bytes", GossipSignatureLen)
	}
	var r, s secp256k1.ModNScalar
	var rb, sb [32]byte
	copy(rb[:], body[:32])
	copy(sb[:], body[32:])
	if r.SetBytes(&rb) != 0 {
		return nil, core.E("net.decodeCompactSig", "signature R overflows group order", nil)
	}
	if s.SetBytes(&sb) != 0 {
		return nil, core.E("net.decodeCompactSig", "signature S overflows group order", nil)
	}
	if r.IsZero() || s.IsZero() {
		return nil, core.E("net.decodeCompactSig", "signature R or S is zero", nil)
	}
	return ecdsa.NewSignature(&r, &s), nil
}
