// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"testing"

	"golang.org/x/crypto/blake2b"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// TestMerkle_MerkleRoot_Good computes a stable, non-zero root over several
// transactions and confirms recomputation is deterministic.
func TestMerkle_MerkleRoot_Good(t *testing.T) {
	txs := []primitives.Transaction{
		coinbaseTx(1, 2000),
		spendTx(0x10, 1000, 4),
		spendTx(0x20, 900, 8),
	}

	root, err := MerkleRoot(txs)
	if err != nil {
		t.Fatalf("MerkleRoot: %v", err)
	}
	if root == (primitives.Hash{}) {
		t.Fatal("merkle root is the zero hash for a non-empty block")
	}

	again, err := MerkleRoot(txs)
	if err != nil {
		t.Fatalf("MerkleRoot again: %v", err)
	}
	if root != again {
		t.Fatal("merkle root is not deterministic")
	}
}

// TestMerkle_MerkleRoot_Bad confirms an empty transaction set yields the
// sentinel hash BLAKE2b-256(""), not the zero hash.
func TestMerkle_MerkleRoot_Bad(t *testing.T) {
	root, err := MerkleRoot(nil)
	if err != nil {
		t.Fatalf("MerkleRoot(nil): %v", err)
	}
	sentinel := primitives.Hash(blake2b.Sum256(nil))
	if root != sentinel {
		t.Fatalf("empty merkle root = %x, want sentinel %x", root, sentinel)
	}
	if root == (primitives.Hash{}) {
		t.Fatal("empty merkle root must not be the zero hash")
	}
}

// TestMerkle_MerkleRoot_Ugly confirms ordering matters: swapping two
// transactions changes the root (the tree is order-sensitive).
func TestMerkle_MerkleRoot_Ugly(t *testing.T) {
	a := spendTx(0x10, 1000, 4)
	b := spendTx(0x20, 900, 8)

	root1, err := MerkleRoot([]primitives.Transaction{a, b})
	if err != nil {
		t.Fatalf("MerkleRoot ab: %v", err)
	}
	root2, err := MerkleRoot([]primitives.Transaction{b, a})
	if err != nil {
		t.Fatalf("MerkleRoot ba: %v", err)
	}
	if root1 == root2 {
		t.Fatal("merkle root is order-insensitive; expected a different root after swap")
	}
}

// TestMerkle_WitnessRoot_Good confirms the witness root over a witnessed set is
// non-zero and differs from the merkle (txid) root, since wtxids commit to the
// witness data.
func TestMerkle_WitnessRoot_Good(t *testing.T) {
	txs := []primitives.Transaction{
		coinbaseTx(1, 2000),
		spendTx(0x10, 1000, 16),
	}

	wroot, err := WitnessRoot(txs)
	if err != nil {
		t.Fatalf("WitnessRoot: %v", err)
	}
	mroot, err := MerkleRoot(txs)
	if err != nil {
		t.Fatalf("MerkleRoot: %v", err)
	}
	if wroot == (primitives.Hash{}) {
		t.Fatal("witness root is zero for a witnessed block")
	}
	if wroot == mroot {
		t.Fatal("witness root equals merkle root despite non-empty witnesses")
	}
}

// TestMerkle_WitnessHash_Bad confirms two transactions identical except for the
// witness share a txid but differ in witness hash.
func TestMerkle_WitnessHash_Bad(t *testing.T) {
	tx := spendTx(0x33, 1000, 4)

	clone := tx
	clone.Inputs = make([]primitives.Input, len(tx.Inputs))
	copy(clone.Inputs, tx.Inputs)
	clone.Inputs[0].Witness = [][]byte{{0x99, 0x98, 0x97}} // different witness, same prevout

	wh1, err := WitnessHash(tx)
	if err != nil {
		t.Fatalf("WitnessHash tx: %v", err)
	}
	wh2, err := WitnessHash(clone)
	if err != nil {
		t.Fatalf("WitnessHash clone: %v", err)
	}
	if wh1 == wh2 {
		t.Fatal("witness hash ignored the witness change")
	}
}

// TestMerkle_WitnessRoot_Ugly confirms a single-transaction witness root equals
// the leaf hash of that transaction's witness hash (no internal node is formed).
func TestMerkle_WitnessRoot_Ugly(t *testing.T) {
	tx := coinbaseTx(1, 2000)

	wroot, err := WitnessRoot([]primitives.Transaction{tx})
	if err != nil {
		t.Fatalf("WitnessRoot: %v", err)
	}

	wh, err := WitnessHash(tx)
	if err != nil {
		t.Fatalf("WitnessHash: %v", err)
	}
	want := hashMerkleLeaf(wh)
	if wroot != want {
		t.Fatalf("single-tx witness root = %x, want leaf hash %x", wroot, want)
	}
}
