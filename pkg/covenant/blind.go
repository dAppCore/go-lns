// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"encoding/binary"

	core "dappco.re/go/core"
	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/primitives"
)

// Blind creates the sealed-bid commitment used by the JS covenant rules.
//
// The encoding matches the reference implementation:
//   - 8-byte little-endian bid value
//   - 32-byte nonce
//
// The nonce is accepted as a fixed-size hash to keep the Go API aligned with
// the rest of the name-system primitives.
//
//	bid, err := covenant.Blind(1000, nonce)
func Blind(value uint64, nonce primitives.Hash) (primitives.Hash, error) {
	var data [40]byte
	binary.LittleEndian.PutUint64(data[:8], value)
	copy(data[8:], nonce[:])

	sum, err := blake2b.New256(nil)
	if err != nil {
		return primitives.Hash{}, core.E("covenant.Blind", "failed to initialise blake2b hash", err)
	}

	if _, err := sum.Write(data[:]); err != nil {
		return primitives.Hash{}, core.E("covenant.Blind", "failed to hash blind data", err)
	}

	var blind primitives.Hash
	copy(blind[:], sum.Sum(nil))
	return blind, nil
}

// GetBlind is an alias for Blind.
//
//	bid, err := covenant.GetBlind(1000, nonce)
func GetBlind(value uint64, nonce primitives.Hash) (primitives.Hash, error) {
	return Blind(value, nonce)
}
