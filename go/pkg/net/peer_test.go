// SPDX-License-Identifier: EUPL-1.2

package net

import (
	core "github.com/dappcore/go"
)

// TestPeerRecord_Validate_Good accepts a well-formed peer record.
func TestPeerRecord_Validate_Good(t *core.T) {
	rec := PeerRecord{
		Name:      "charon",
		PublicKey: make([]byte, 33),
		Addresses: []string{"10.0.0.1:14038"},
		Cap:       []Capability{CapVPN, CapDNS},
	}
	core.AssertNoError(t, rec.Validate())
	core.AssertTrue(t, rec.HasCapability(CapVPN))
	core.AssertFalse(t, rec.HasCapability(CapExit))
}

// TestPeerRecord_Validate_Bad rejects an empty name and an empty public key.
func TestPeerRecord_Validate_Bad(t *core.T) {
	core.AssertError(t, PeerRecord{PublicKey: make([]byte, 33)}.Validate())
	core.AssertError(t, PeerRecord{Name: "charon"}.Validate())
}

// TestPeerRecord_Validate_Ugly rejects every oversized field (RFC §12 untrusted
// length guard reused for record bounds).
func TestPeerRecord_Validate_Ugly(t *core.T) {
	long := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = 'x'
		}
		return string(b)
	}

	core.AssertError(t, PeerRecord{Name: long(MaxNameLen + 1), PublicKey: make([]byte, 33)}.Validate())
	core.AssertError(t, PeerRecord{Name: "c", PublicKey: make([]byte, MaxPublicKeyLen+1)}.Validate())

	addrs := make([]string, MaxAddresses+1)
	core.AssertError(t, PeerRecord{Name: "c", PublicKey: make([]byte, 33), Addresses: addrs}.Validate())

	core.AssertError(t, PeerRecord{Name: "c", PublicKey: make([]byte, 33), Addresses: []string{long(MaxAddressLen + 1)}}.Validate())

	caps := make([]Capability, MaxCaps+1)
	core.AssertError(t, PeerRecord{Name: "c", PublicKey: make([]byte, 33), Cap: caps}.Validate())

	core.AssertError(t, PeerRecord{Name: "c", PublicKey: make([]byte, 33), Cap: []Capability{long(MaxCapLen + 1)}}.Validate())
}

// TestPeerRecord_AsNameRecord_Good projects a peer record into a name record and
// confirms the projection is HSD-sourced and deep-copied.
func TestPeerRecord_AsNameRecord_Good(t *core.T) {
	rec := PeerRecord{
		Name:       "charon",
		PublicKey:  []byte{1, 2, 3},
		Addresses:  []string{"a"},
		Cap:        []Capability{CapVPN},
		TrustScore: 0.7,
	}
	nr := rec.AsNameRecord()
	core.AssertEqual(t, "charon", nr.Name)
	core.AssertEqual(t, SourceHSD, nr.Source)
	core.AssertInDelta(t, 0.7, nr.TrustScore, 1e-9)

	// Mutating the projection must not touch the source record.
	nr.PublicKey[0] = 9
	core.AssertEqual(t, byte(1), rec.PublicKey[0])
}

// TestPeerRecord_Clone_Good confirms Clone deep-copies all slices.
func TestPeerRecord_Clone_Good(t *core.T) {
	rec := PeerRecord{Name: "c", PublicKey: []byte{1}, Addresses: []string{"a"}, Cap: []Capability{CapDNS}}
	clone := rec.Clone()
	clone.PublicKey[0] = 9
	clone.Addresses[0] = "b"
	clone.Cap[0] = CapExit
	core.AssertEqual(t, byte(1), rec.PublicKey[0])
	core.AssertEqual(t, "a", rec.Addresses[0])
	core.AssertEqual(t, CapDNS, rec.Cap[0])
}
