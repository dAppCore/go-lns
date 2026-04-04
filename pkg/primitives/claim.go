// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	core "dappco.re/go/core"
	"golang.org/x/crypto/blake2b"
)

const maxClaimSize = 10000

// Claim wraps a DNSSEC ownership proof blob used by name claims.
//
// The structure mirrors the JS reference wrapper: it stores the raw proof
// bytes, hashes them with blake2b, and serializes them with a compact length
// prefix.
type Claim struct {
	Blob []byte
}

// NewClaim constructs an empty claim wrapper.
func NewClaim() *Claim {
	return &Claim{}
}

// GetNewClaim is an alias for NewClaim.
func GetNewClaim() *Claim {
	return NewClaim()
}

// Refresh is kept for parity with the JS reference wrapper.
func (c *Claim) Refresh() *Claim {
	return c
}

// Clone returns a shallow copy of the claim wrapper.
func (c Claim) Clone() Claim {
	return Claim{Blob: c.GetBlob()}
}

// Hash returns the blake2b digest of the claim blob.
func (c Claim) Hash() Hash {
	return blake2b.Sum256(c.Blob)
}

// GetHash is an alias for Hash.
func (c Claim) GetHash() Hash {
	return c.Hash()
}

// HashHex returns the claim hash encoded as lowercase hexadecimal.
func (c Claim) HashHex() string {
	sum := c.Hash()
	return hex.EncodeToString(sum[:])
}

// GetHashHex is an alias for HashHex.
func (c Claim) GetHashHex() string {
	return c.HashHex()
}

// GetSize returns the serialized size of the claim wrapper.
func (c Claim) GetSize() int {
	return 2 + len(c.Blob)
}

// Weight returns the claim weight in weight units.
//
// The claim wrapper is witness-like data, so its weight matches its serialized
// size in this implementation.
func (c Claim) Weight() int {
	return c.GetSize()
}

// GetWeight is an alias for Weight.
func (c Claim) GetWeight() int {
	return c.Weight()
}

// VirtualSize returns the compact size used for fee estimation.
//
// This mirrors the JS reference helper by rounding weight units up to the next
// virtual byte.
func (c Claim) VirtualSize() int {
	return (c.Weight() + 3) / 4
}

// GetVirtualSize is an alias for VirtualSize.
func (c Claim) GetVirtualSize() int {
	return c.VirtualSize()
}

// MarshalBinary serializes the claim wrapper.
func (c Claim) MarshalBinary() ([]byte, error) {
	if len(c.Blob) > maxClaimSize {
		return nil, core.E("primitives.Claim.MarshalBinary", "claim too large", nil)
	}

	if len(c.Blob) > 0xffff {
		return nil, core.E("primitives.Claim.MarshalBinary", "claim exceeds wire size", nil)
	}

	buf := make([]byte, 2+len(c.Blob))
	binary.LittleEndian.PutUint16(buf[:2], uint16(len(c.Blob)))
	copy(buf[2:], c.Blob)
	return buf, nil
}

// UnmarshalBinary decodes the claim wrapper.
func (c *Claim) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("primitives.Claim.UnmarshalBinary: short buffer")
	}

	size := binary.LittleEndian.Uint16(data[:2])
	if int(size) > maxClaimSize {
		return fmt.Errorf("primitives.Claim.UnmarshalBinary: claim too large")
	}

	if len(data) != 2+int(size) {
		return fmt.Errorf("primitives.Claim.UnmarshalBinary: trailing data")
	}

	c.Blob = append(c.Blob[:0], data[2:]...)
	return nil
}

// ToBlob returns a copy of the raw claim blob.
func (c Claim) ToBlob() []byte {
	return c.GetBlob()
}

// GetToBlob is an alias for ToBlob.
func (c Claim) GetToBlob() []byte {
	return c.ToBlob()
}

// GetBlob returns a copy of the raw claim blob.
func (c Claim) GetBlob() []byte {
	if len(c.Blob) == 0 {
		return nil
	}

	blob := make([]byte, len(c.Blob))
	copy(blob, c.Blob)
	return blob
}

// FromBlob replaces the raw claim blob.
func (c *Claim) FromBlob(blob []byte) *Claim {
	if c == nil {
		return nil
	}

	c.Blob = append(c.Blob[:0], blob...)
	return c
}

// GetFromBlob is an alias for FromBlob.
func (c *Claim) GetFromBlob(blob []byte) *Claim {
	return c.FromBlob(blob)
}
