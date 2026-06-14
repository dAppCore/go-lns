// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	core "dappco.re/go"
	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/mempool"
	"dappco.re/go/lns/pkg/primitives"
)

// Merkle node domain-separation prefixes, mirroring bcrypto/lib/mrkl. Leaves and
// internal nodes are tagged so a leaf hash can never be reinterpreted as an
// internal node (the CVE-2012-2459 class of malleability).
var (
	merkleLeafPrefix     = []byte{0x00}
	merkleInternalPrefix = []byte{0x01}
)

// MerkleRoot computes the transaction merkle root over a block's transaction
// ids, in the order given. The id of each transaction is its BLAKE2b-256
// non-witness digest (mempool.ComputeTxID), so this is the root the header
// MerkleRoot field commits to. An empty slice yields the sentinel hash
// BLAKE2b-256("") rather than the zero hash, matching HSD.
//
//	root, err := mining.MerkleRoot(block.Txs)
func MerkleRoot(txs []primitives.Transaction) (primitives.Hash, error) {
	leaves := make([]primitives.Hash, len(txs))
	for i := range txs {
		id, err := mempool.ComputeTxID(txs[i])
		if err != nil {
			return primitives.Hash{}, core.E("mining.MerkleRoot", "compute tx id", err)
		}
		leaves[i] = id
	}
	return merkleRootFromLeaves(leaves), nil
}

// WitnessRoot computes the witness merkle root over each transaction's witness
// hash, in the order given. It is the root the header WitnessRoot field commits
// to. An empty slice yields the sentinel hash.
//
//	root, err := mining.WitnessRoot(block.Txs)
func WitnessRoot(txs []primitives.Transaction) (primitives.Hash, error) {
	leaves := make([]primitives.Hash, len(txs))
	for i := range txs {
		wh, err := WitnessHash(txs[i])
		if err != nil {
			return primitives.Hash{}, core.E("mining.WitnessRoot", "compute witness hash", err)
		}
		leaves[i] = wh
	}
	return merkleRootFromLeaves(leaves), nil
}

// WitnessHash returns the witness transaction id (wtxid): the BLAKE2b-256 digest
// of the non-witness txid concatenated with the BLAKE2b-256 digest of the
// witness data. It mirrors HSD tx.witnessHash, where left = HASH(base-data),
// right = HASH(witness-data), and wtxid = HASH(left ‖ right).
//
//	wh, err := mining.WitnessHash(tx)
func WitnessHash(tx primitives.Transaction) (primitives.Hash, error) {
	left, err := mempool.ComputeTxID(tx)
	if err != nil {
		return primitives.Hash{}, core.E("mining.WitnessHash", "compute base id", err)
	}

	right := blake2b.Sum256(serializeWitnessData(tx))

	var buf [2 * HashSize]byte
	copy(buf[:HashSize], left[:])
	copy(buf[HashSize:], right[:])

	return blake2b.Sum256(buf[:]), nil
}

// serializeWitnessData encodes only the witness stacks, in input order, in the
// HSD trailing-witness layout: for each input, varint(itemCount) then each item
// as varint(len) ‖ bytes. This is the "witness data" segment HSD hashes for the
// right branch of the wtxid.
func serializeWitnessData(tx primitives.Transaction) []byte {
	buf := make([]byte, 0, 64)
	for i := range tx.Inputs {
		witness := tx.Inputs[i].Witness
		buf = appendVarint(buf, uint64(len(witness)))
		for _, item := range witness {
			buf = appendVarint(buf, uint64(len(item)))
			buf = append(buf, item...)
		}
	}
	return buf
}

// merkleRootFromLeaves reduces leaf hashes to a single root using the tagged,
// sentinel-padded tree from bcrypto/lib/mrkl: each leaf is HASH(0x00 ‖ id), each
// internal node is HASH(0x01 ‖ left ‖ right), and an odd node at a level pairs
// with the sentinel HASH(""). No leaf count yields the sentinel.
func merkleRootFromLeaves(leaves []primitives.Hash) primitives.Hash {
	sentinel := blake2b.Sum256(nil)

	if len(leaves) == 0 {
		return sentinel
	}

	level := make([]primitives.Hash, len(leaves))
	for i := range leaves {
		level[i] = hashMerkleLeaf(leaves[i])
	}

	for len(level) > 1 {
		next := make([]primitives.Hash, 0, (len(level)+1)/2)
		for j := 0; j < len(level); j += 2 {
			left := level[j]
			right := sentinel
			if j+1 < len(level) {
				right = level[j+1]
			}
			next = append(next, hashMerkleInternal(left, right))
		}
		level = next
	}

	return level[0]
}

// hashMerkleLeaf hashes a leaf with the leaf domain prefix.
func hashMerkleLeaf(data primitives.Hash) primitives.Hash {
	buf := make([]byte, 0, 1+HashSize)
	buf = append(buf, merkleLeafPrefix...)
	buf = append(buf, data[:]...)
	return blake2b.Sum256(buf)
}

// hashMerkleInternal hashes two child nodes with the internal domain prefix.
func hashMerkleInternal(left, right primitives.Hash) primitives.Hash {
	buf := make([]byte, 0, 1+2*HashSize)
	buf = append(buf, merkleInternalPrefix...)
	buf = append(buf, left[:]...)
	buf = append(buf, right[:]...)
	return blake2b.Sum256(buf)
}
