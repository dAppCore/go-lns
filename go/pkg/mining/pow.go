// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"math/big"

	core "dappco.re/go"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"

	"dappco.re/go/lns/pkg/primitives"
)

// shareHash computes the HSD share hash over a header's preheader. The share
// hash is the pre-mask PoW value; mining pools target it to award shares while
// hiding whether a share is also a full block.
//
//	left      = BLAKE2b-512(prehead)            // 64 bytes
//	right     = SHA3-256(prehead ‖ pad(8))      // 32 bytes
//	shareHash = BLAKE2b-256(left ‖ pad(32) ‖ right)
func shareHash(h primitives.BlockHeader) (primitives.Hash, error) {
	data := preHead(h)

	left, err := blake2b.New512(nil)
	if err != nil {
		return primitives.Hash{}, core.E("mining.shareHash", "init blake2b-512", err)
	}
	left.Write(data)
	leftSum := left.Sum(nil)

	sha := sha3.New256()
	sha.Write(data)
	sha.Write(padding(h, 8))
	rightSum := sha.Sum(nil)

	combined := make([]byte, 0, len(leftSum)+32+len(rightSum))
	combined = append(combined, leftSum...)
	combined = append(combined, padding(h, 32)...)
	combined = append(combined, rightSum...)

	return blake2b.Sum256(combined), nil
}

// PowHash computes the header proof-of-work hash: the share hash XORed with the
// 32-byte mask. This is the value compared against the difficulty target. With a
// zero mask (the common case outside pooled mining) it equals the share hash.
//
//	hash, err := mining.PowHash(block.Header)
func PowHash(header primitives.BlockHeader) (primitives.Hash, error) {
	hash, err := shareHash(header)
	if err != nil {
		return primitives.Hash{}, core.E("mining.PowHash", "share hash", err)
	}

	mask := fixedBytes(header.Mask, maskSize)
	for i := 0; i < HashSize; i++ {
		hash[i] ^= mask[i]
	}
	return hash, nil
}

// Solve searches for a nonce that makes the header's PoW hash meet the target,
// trying up to rounds values starting from the header's current nonce. It
// returns the winning nonce and true on success, or the last nonce tried and
// false if the range is exhausted. The supplied header is not mutated; the
// caller assigns the returned nonce on success.
//
// Solve mutates only the nonce. To search a wider space, roll the extra nonce
// (or timestamp) between calls. Keep rounds small in tests by mining at a
// trivial target so the search terminates quickly.
//
//	nonce, ok := mining.Solve(hdr, mining.TargetFromBits(hdr.Bits), 1<<20)
//	if ok { hdr.Nonce = nonce }
func Solve(header primitives.BlockHeader, target *big.Int, rounds uint32) (uint32, bool) {
	work := normalizeHeader(header)
	nonce := work.Nonce

	for i := uint32(0); i < rounds; i++ {
		work.Nonce = nonce

		hash, err := shareHash(work)
		if err != nil {
			return nonce, false
		}

		mask := work.Mask // already normalised to 32 bytes
		for k := 0; k < HashSize; k++ {
			hash[k] ^= mask[k]
		}

		if meetsTarget(hash, target) {
			return nonce, true
		}

		// Stop before wrapping past the uint32 range.
		if nonce == ^uint32(0) {
			return nonce, false
		}
		nonce++
	}

	return nonce, false
}

// meetsTarget reports whether a hash read big-endian is a valid solution for a
// positive, in-range target.
func meetsTarget(hash primitives.Hash, target *big.Int) bool {
	if target == nil || target.Sign() <= 0 || target.BitLen() > maxTargetBits {
		return false
	}
	value := new(big.Int).SetBytes(hash[:])
	return value.Cmp(target) <= 0
}
