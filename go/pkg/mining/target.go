// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"math/big"

	core "dappco.re/go"

	"dappco.re/go/lns/pkg/primitives"
)

// maxTargetBits is the bit ceiling for a valid target. A target wider than 256
// bits cannot be met by any 32-byte hash and is rejected, matching HSD
// consensus.verifyPOW.
const maxTargetBits = 256

// TargetFromBits decodes the compact difficulty representation (the header Bits
// field) into a 256-bit target. The compact form packs an exponent in the top
// byte and a 23-bit mantissa, with bit 23 reserved as a sign flag. It mirrors
// HSD consensus.fromCompact.
//
//	target := mining.TargetFromBits(0x1d00ffff)
//	if target.Sign() == 0 { /* zero bits: no valid PoW */ }
func TargetFromBits(bits uint32) *big.Int {
	if bits == 0 {
		return big.NewInt(0)
	}

	exponent := bits >> 24
	negative := (bits >> 23) & 1
	mantissa := bits & 0x7fffff

	target := new(big.Int)
	if exponent <= 3 {
		mantissa >>= 8 * (3 - exponent)
		target.SetUint64(uint64(mantissa))
	} else {
		target.SetUint64(uint64(mantissa))
		target.Lsh(target, uint(8*(exponent-3)))
	}

	if negative == 1 {
		target.Neg(target)
	}

	return target
}

// BitsFromTarget encodes a 256-bit target into the compact representation stored
// in the header Bits field. It is the inverse of TargetFromBits and mirrors HSD
// consensus.toCompact, including the mantissa normalisation that keeps bit 23
// clear (so it is never mistaken for the sign flag).
//
//	bits := mining.BitsFromTarget(target)
func BitsFromTarget(target *big.Int) uint32 {
	if target.Sign() == 0 {
		return 0
	}

	work := new(big.Int).Abs(target)
	exponent := uint32(len(work.Bytes()))

	var mantissa uint32
	if exponent <= 3 {
		// Shift the value up into the mantissa position.
		mantissa = uint32(work.Uint64())
		mantissa <<= 8 * (3 - exponent)
	} else {
		shifted := new(big.Int).Rsh(work, uint(8*(exponent-3)))
		mantissa = uint32(shifted.Uint64())
	}

	// If the top mantissa bit is set it would collide with the sign flag, so
	// shift down one byte and bump the exponent.
	if mantissa&0x800000 != 0 {
		mantissa >>= 8
		exponent++
	}

	compact := (exponent << 24) | mantissa
	if target.Sign() < 0 {
		compact |= 0x800000
	}

	return compact
}

// CheckTarget reports whether a 32-byte PoW hash satisfies the target decoded
// from bits: the hash, read big-endian, must be a positive target and less than
// or equal to it. A zero, negative, or over-wide target can never be satisfied.
// It mirrors HSD consensus.verifyPOW.
//
//	if mining.CheckTarget(powHash, header.Bits) { /* block solves */ }
func CheckTarget(hash primitives.Hash, bits uint32) bool {
	target := TargetFromBits(bits)
	if target.Sign() <= 0 {
		return false
	}
	if target.BitLen() > maxTargetBits {
		return false
	}

	value := new(big.Int).SetBytes(hash[:])
	return value.Cmp(target) <= 0
}

// VerifyPOW recomputes the header PoW hash and reports whether it meets the
// header's own difficulty target. It is the one-call header validity check a
// chain validator runs on an inbound block.
//
//	ok, err := mining.VerifyPOW(block.Header)
func VerifyPOW(header primitives.BlockHeader) (bool, error) {
	hash, err := PowHash(header)
	if err != nil {
		return false, core.E("mining.VerifyPOW", "compute pow hash", err)
	}
	return CheckTarget(hash, header.Bits), nil
}
