// SPDX-License-Identifier: EUPL-1.2

package net

import (
	core "dappco.re/go"
)

// Capability is a peer service capability advertised through a PeerRecord.
// Capabilities gate which peers a caller will route particular traffic to.
type Capability = string

// Capability constants enumerate the service roles a peer may advertise.
// They mirror the alias-comment cap field (RFC §16.3.3) and the PeerRecord
// Cap list (RFC §10.1).
const (
	// CapVPN marks a peer that terminates VPN tunnels.
	CapVPN Capability = "vpn"
	// CapDNS marks a peer that serves the LNS DNS zone.
	CapDNS Capability = "dns"
	// CapProxy marks a peer that relays proxied traffic.
	CapProxy Capability = "proxy"
	// CapExit marks a peer that acts as an exit node.
	CapExit Capability = "exit"
)

// Bound constants cap untrusted-length fields so a hostile peer cannot drive
// an unbounded allocation when a record is decoded from the wire.
const (
	// MaxNameLen is the longest accepted bare peer label (RFC name limit is 63).
	MaxNameLen = 63
	// MaxPublicKeyLen is the longest accepted chain public key (33-byte
	// compressed secp256k1 plus headroom for alternative encodings).
	MaxPublicKeyLen = 64
	// MaxAddresses is the most endpoint strings accepted on one record.
	MaxAddresses = 16
	// MaxAddressLen is the longest accepted single endpoint string.
	MaxAddressLen = 256
	// MaxCaps is the most capability tokens accepted on one record.
	MaxCaps = 16
	// MaxCapLen is the longest accepted single capability token.
	MaxCapLen = 32
)

// PeerRecord identifies a network node (RFC §10.1, "Peer Names"). It resolves a
// bare label to the peer's chain public key plus reachable addresses, the
// capabilities the peer advertises, and a Poindexter-derived trust score in the
// range 0.0–1.0.
//
//	rec := net.PeerRecord{
//	    Name:       "charon",
//	    PublicKey:  pub, // 33-byte compressed secp256k1 chain key
//	    Addresses:  []string{"10.69.69.165:14038"},
//	    Cap:        []net.Capability{net.CapVPN, net.CapDNS},
//	    TrustScore: 0.82,
//	}
type PeerRecord struct {
	// Name is the bare peer label without the .lthn suffix ("charon").
	Name string `json:"name"`
	// PublicKey is the peer's chain public key (compressed secp256k1) used to
	// verify gossip the peer signs.
	PublicKey []byte `json:"publicKey"`
	// Addresses are the reachable endpoints (host:port or multiaddr).
	Addresses []string `json:"addresses"`
	// Cap lists the capabilities the peer advertises.
	Cap []Capability `json:"cap"`
	// TrustScore is the Poindexter-derived score, 0.0–1.0.
	TrustScore float64 `json:"trustScore"`
}

// HasCapability reports whether the peer advertises the given capability.
//
//	if rec.HasCapability(net.CapExit) { /* eligible as an exit node */ }
func (r PeerRecord) HasCapability(cap Capability) bool {
	for _, c := range r.Cap {
		if c == cap {
			return true
		}
	}
	return false
}

// Validate reports whether the record's fields are within the bounded limits.
// It rejects an empty name, an over-long or empty public key, and any field
// that exceeds the wire bounds, so a record decoded from an untrusted peer
// cannot smuggle oversized data into the cache.
//
//	if err := rec.Validate(); err != nil { /* reject */ }
func (r PeerRecord) Validate() error {
	if r.Name == "" {
		return core.E("net.PeerRecord.Validate", "name is required", nil)
	}
	if len(r.Name) > MaxNameLen {
		return core.Errorf("net.PeerRecord.Validate: name length %d exceeds %d", len(r.Name), MaxNameLen)
	}
	if len(r.PublicKey) == 0 {
		return core.E("net.PeerRecord.Validate", "public key is required", nil)
	}
	if len(r.PublicKey) > MaxPublicKeyLen {
		return core.Errorf("net.PeerRecord.Validate: public key length %d exceeds %d", len(r.PublicKey), MaxPublicKeyLen)
	}
	if len(r.Addresses) > MaxAddresses {
		return core.Errorf("net.PeerRecord.Validate: address count %d exceeds %d", len(r.Addresses), MaxAddresses)
	}
	for _, addr := range r.Addresses {
		if len(addr) > MaxAddressLen {
			return core.Errorf("net.PeerRecord.Validate: address length %d exceeds %d", len(addr), MaxAddressLen)
		}
	}
	if len(r.Cap) > MaxCaps {
		return core.Errorf("net.PeerRecord.Validate: capability count %d exceeds %d", len(r.Cap), MaxCaps)
	}
	for _, c := range r.Cap {
		if len(c) > MaxCapLen {
			return core.Errorf("net.PeerRecord.Validate: capability length %d exceeds %d", len(c), MaxCapLen)
		}
	}
	return nil
}

// AsNameRecord projects the peer record into a NameRecord for gossip transport.
// The result is HSD-sourced by default; the gossip merge path re-stamps the
// source and trust penalty for records that arrive over the wire.
//
//	nr := rec.AsNameRecord()
func (r PeerRecord) AsNameRecord() NameRecord {
	return NameRecord{
		Name:       r.Name,
		PublicKey:  append([]byte(nil), r.PublicKey...),
		Addresses:  append([]string(nil), r.Addresses...),
		Cap:        append([]Capability(nil), r.Cap...),
		TrustScore: r.TrustScore,
		Source:     SourceHSD,
	}
}

// Clone returns a deep copy of the record so callers cannot mutate slices held
// elsewhere (the registry and cache hand out copies for this reason).
func (r PeerRecord) Clone() PeerRecord {
	return PeerRecord{
		Name:       r.Name,
		PublicKey:  append([]byte(nil), r.PublicKey...),
		Addresses:  append([]string(nil), r.Addresses...),
		Cap:        append([]Capability(nil), r.Cap...),
		TrustScore: r.TrustScore,
	}
}
