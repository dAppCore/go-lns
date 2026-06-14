// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"encoding/binary"

	core "dappco.re/go"
	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/primitives"
)

// Header field byte lengths and offsets, mirroring HSD lib/primitives/
// abstractblock.js. The preheader and subheader are each exactly one BLAKE2b
// block (128 bytes) when the commitment/padding form is hashed.
const (
	subheadSize = 128 // extraNonce(24)+reservedRoot(32)+witnessRoot(32)+merkleRoot(32)+version(4)+bits(4)
	preheadSize = 128 // nonce(4)+time(8)+pad(20)+prevBlock(32)+treeRoot(32)+commitHash(32)
	maskSize    = 32
)

// normalizeHeader returns a copy of the header with the variable-length byte
// fields (extra nonce, mask) padded or truncated to their fixed consensus sizes,
// so a caller may leave them nil. The original is not mutated.
func normalizeHeader(h primitives.BlockHeader) primitives.BlockHeader {
	h.ExtraNonce = fixedBytes(h.ExtraNonce, NonceSize)
	h.Mask = fixedBytes(h.Mask, maskSize)
	return h
}

// fixedBytes returns a length-n copy of src, left-aligned: shorter input is
// zero-padded on the right, longer input is truncated.
func fixedBytes(src []byte, n int) []byte {
	out := make([]byte, n)
	copy(out, src)
	return out
}

// padding derives n bytes of deterministic header padding from the previous
// block and tree root: pad[i] = prevBlock[i%32] XOR treeRoot[i%32]. It mirrors
// HSD AbstractBlock.padding and binds the PoW to the chain position.
func padding(h primitives.BlockHeader, n int) []byte {
	pad := make([]byte, n)
	for i := 0; i < n; i++ {
		pad[i] = h.PrevBlock[i%HashSize] ^ h.TreeRoot[i%HashSize]
	}
	return pad
}

// subHead serialises the subheader: the miner-mutable, less SPV-essential data,
// in HSD wire order, exactly one BLAKE2b block.
//
//	extraNonce(24) ‖ reservedRoot(32) ‖ witnessRoot(32) ‖ merkleRoot(32) ‖
//	version(4 LE) ‖ bits(4 LE)
func subHead(h primitives.BlockHeader) []byte {
	buf := make([]byte, 0, subheadSize)
	buf = append(buf, fixedBytes(h.ExtraNonce, NonceSize)...)
	buf = append(buf, h.ReservedRoot[:]...)
	buf = append(buf, h.WitnessRoot[:]...)
	buf = append(buf, h.MerkleRoot[:]...)
	buf = appendU32(buf, h.Version)
	buf = appendU32(buf, h.Bits)
	return buf
}

// subHash is BLAKE2b-256 of the subheader.
func subHash(h primitives.BlockHeader) primitives.Hash {
	return blake2b.Sum256(subHead(h))
}

// maskHash is BLAKE2b-256(prevBlock ‖ mask). HSD hashes the mask with the
// previous block so a pool may reuse a mask across blocks safely.
func maskHash(h primitives.BlockHeader) primitives.Hash {
	buf := make([]byte, 0, HashSize+maskSize)
	buf = append(buf, h.PrevBlock[:]...)
	buf = append(buf, fixedBytes(h.Mask, maskSize)...)
	return blake2b.Sum256(buf)
}

// commitHash is BLAKE2b-256(subHash ‖ maskHash). It folds all miner-malleable
// subheader data and the mask commitment into one 32-byte value carried in the
// preheader, so changing any subheader field costs the miner two hash rounds.
func commitHash(h primitives.BlockHeader) primitives.Hash {
	sub := subHash(h)
	mask := maskHash(h)
	buf := make([]byte, 0, 2*HashSize)
	buf = append(buf, sub[:]...)
	buf = append(buf, mask[:]...)
	return blake2b.Sum256(buf)
}

// preHead serialises the preheader: the SPV-essential data plus the nonce and
// timestamp, in HSD wire order, exactly one BLAKE2b block. This is the buffer
// the share hash is computed over.
//
//	nonce(4 LE) ‖ time(8 LE) ‖ pad(20) ‖ prevBlock(32) ‖ treeRoot(32) ‖
//	commitHash(32)
func preHead(h primitives.BlockHeader) []byte {
	commit := commitHash(h)
	buf := make([]byte, 0, preheadSize)
	buf = appendU32(buf, h.Nonce)
	buf = appendU64(buf, h.Time)
	buf = append(buf, padding(h, 20)...)
	buf = append(buf, h.PrevBlock[:]...)
	buf = append(buf, h.TreeRoot[:]...)
	buf = append(buf, commit[:]...)
	return buf
}

// MarshalHeader serialises a block header to its 236-byte wire form (HSD
// HEADER_SIZE): preheader fields, subheader fields, then the 32-byte mask. This
// is the on-the-wire header, distinct from the hashed preheader/subheader forms.
//
//	raw, err := mining.MarshalHeader(block.Header)
func MarshalHeader(h primitives.BlockHeader) ([]byte, error) {
	buf := make([]byte, 0, HeaderSize)

	// Preheader.
	buf = appendU32(buf, h.Nonce)
	buf = appendU64(buf, h.Time)
	buf = append(buf, h.PrevBlock[:]...)
	buf = append(buf, h.TreeRoot[:]...)

	// Subheader.
	buf = append(buf, fixedBytes(h.ExtraNonce, NonceSize)...)
	buf = append(buf, h.ReservedRoot[:]...)
	buf = append(buf, h.WitnessRoot[:]...)
	buf = append(buf, h.MerkleRoot[:]...)
	buf = appendU32(buf, h.Version)
	buf = appendU32(buf, h.Bits)

	// Mask.
	buf = append(buf, fixedBytes(h.Mask, maskSize)...)

	if len(buf) != HeaderSize {
		return nil, core.E("mining.MarshalHeader",
			core.Sprintf("encoded header is %d bytes, want %d", len(buf), HeaderSize), nil)
	}
	return buf, nil
}

// appendU32 appends a little-endian uint32.
func appendU32(dst []byte, v uint32) []byte {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], v)
	return append(dst, tmp[:]...)
}

// appendU64 appends a little-endian uint64.
func appendU64(dst []byte, v uint64) []byte {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], v)
	return append(dst, tmp[:]...)
}

// appendVarint appends a Bitcoin/HSD compact-size integer.
func appendVarint(dst []byte, n uint64) []byte {
	switch {
	case n < 0xfd:
		return append(dst, byte(n))
	case n <= 0xffff:
		return append(dst, 0xfd, byte(n), byte(n>>8))
	case n <= 0xffffffff:
		return append(dst, 0xfe, byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
	default:
		return append(dst,
			0xff,
			byte(n), byte(n>>8), byte(n>>16), byte(n>>24),
			byte(n>>32), byte(n>>40), byte(n>>48), byte(n>>56))
	}
}
