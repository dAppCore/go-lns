// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

// TestTxid_ComputeTxID_Good checks the id is deterministic and witness-independent.
func TestTxid_ComputeTxID_Good(t *testing.T) {
	op := outpoint(0x11, 0)
	addr, priv, _ := fundedCoin(t, 1, 2000, op)
	tx := signedSpend(t, priv, op, 2000, 1500, addr)

	id, err := ComputeTxID(tx)
	if err != nil {
		t.Fatalf("ComputeTxID: %v", err)
	}
	if id == (primitives.Hash{}) {
		t.Fatal("id is zero for a non-empty transaction")
	}

	// Recomputing yields the same id.
	again, err := ComputeTxID(tx)
	if err != nil {
		t.Fatalf("ComputeTxID again: %v", err)
	}
	if id != again {
		t.Fatal("ComputeTxID is not deterministic")
	}

	// Mutating only the witness must NOT change the id (non-witness serialisation).
	clone := tx
	clone.Inputs = []primitives.Input{{Prevout: op, Sequence: 0xffffffff, Witness: [][]byte{{0x00}}}}
	clone.Outputs = tx.Outputs
	witnessOnly, err := ComputeTxID(clone)
	if err != nil {
		t.Fatalf("ComputeTxID witness-only: %v", err)
	}
	if witnessOnly != id {
		t.Fatal("changing the witness changed the txid")
	}
}

// TestTxid_ComputeTxID_Bad checks distinct transactions get distinct ids.
func TestTxid_ComputeTxID_Bad(t *testing.T) {
	op := outpoint(0x22, 0)
	addr, priv, _ := fundedCoin(t, 2, 2000, op)

	a := signedSpend(t, priv, op, 2000, 1500, addr)
	b := signedSpend(t, priv, op, 2000, 1400, addr) // different output value

	ida, err := ComputeTxID(a)
	if err != nil {
		t.Fatalf("ComputeTxID a: %v", err)
	}
	idb, err := ComputeTxID(b)
	if err != nil {
		t.Fatalf("ComputeTxID b: %v", err)
	}
	if ida == idb {
		t.Fatal("transactions differing in output value share an id")
	}
}

// TestTxid_ComputeTxID_Ugly checks degenerate transactions still serialise.
func TestTxid_ComputeTxID_Ugly(t *testing.T) {
	// Empty transaction: no inputs, no outputs. Still hashes (to a fixed digest),
	// the pool rejects it elsewhere via ErrEmptyTx.
	var empty primitives.Transaction
	id, err := ComputeTxID(empty)
	if err != nil {
		t.Fatalf("ComputeTxID empty: %v", err)
	}
	again, _ := ComputeTxID(empty)
	if id != again {
		t.Fatal("empty-tx id not deterministic")
	}

	// An over-long address hash (>255 bytes) cannot be serialised: error, no panic.
	bad := primitives.Transaction{
		Outputs: []primitives.Output{{Address: primitives.Address{Hash: make([]byte, 256)}}},
	}
	if _, err := ComputeTxID(bad); err == nil {
		t.Fatal("expected error for over-long address hash")
	}
}

// TestTxid_SerializeNoWitness_Good checks the encoding round-trips through size.
func TestTxid_SerializeNoWitness_Good(t *testing.T) {
	op := outpoint(0x33, 1)
	addr, priv, _ := fundedCoin(t, 3, 5000, op)
	tx := signedSpend(t, priv, op, 5000, 4000, addr)

	raw, err := SerializeNoWitness(tx)
	if err != nil {
		t.Fatalf("SerializeNoWitness: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("empty serialisation")
	}

	// version(4) + varint(1) + [prevout(36)+seq(4)] + varint(1) + output + locktime(4)
	want := 4 + 1 + 40 + 1 + (8 + 2 + len(addr.Hash) + 2) + 4
	if len(raw) != want {
		t.Fatalf("serialised size = %d, want %d", len(raw), want)
	}
}

// TestTxid_SerializeNoWitness_Bad checks an unserialisable output errors.
func TestTxid_SerializeNoWitness_Bad(t *testing.T) {
	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: outpoint(0x44, 0)}},
		Outputs: []primitives.Output{{Address: primitives.Address{Hash: make([]byte, 300)}}},
	}
	if _, err := SerializeNoWitness(tx); err == nil {
		t.Fatal("expected error for over-long address hash")
	}
}

// TestTxid_SerializeNoWitness_Ugly checks the empty transaction encodes minimally.
func TestTxid_SerializeNoWitness_Ugly(t *testing.T) {
	var empty primitives.Transaction
	raw, err := SerializeNoWitness(empty)
	if err != nil {
		t.Fatalf("SerializeNoWitness empty: %v", err)
	}
	// version(4) + varint(0) + varint(0) + locktime(4) = 10 bytes.
	if len(raw) != 10 {
		t.Fatalf("empty serialisation = %d bytes, want 10", len(raw))
	}
}

// TestTxid_TxSize_Good checks the witness-inclusive size exceeds the base size.
func TestTxid_TxSize_Good(t *testing.T) {
	op := outpoint(0x55, 0)
	addr, priv, _ := fundedCoin(t, 4, 3000, op)
	tx := signedSpend(t, priv, op, 3000, 2500, addr)

	base, err := SerializeNoWitness(tx)
	if err != nil {
		t.Fatalf("SerializeNoWitness: %v", err)
	}
	size, err := TxSize(tx)
	if err != nil {
		t.Fatalf("TxSize: %v", err)
	}
	if size <= len(base) {
		t.Fatalf("TxSize %d should exceed base %d (witness adds bytes)", size, len(base))
	}
}

// TestTxid_TxSize_Bad checks an unserialisable transaction errors.
func TestTxid_TxSize_Bad(t *testing.T) {
	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: outpoint(0x66, 0)}},
		Outputs: []primitives.Output{{Address: primitives.Address{Hash: make([]byte, 300)}}},
	}
	if _, err := TxSize(tx); err == nil {
		t.Fatal("expected error sizing an unserialisable transaction")
	}
}

// TestTxid_TxSize_Ugly checks a no-witness transaction sizes as base + empty stacks.
func TestTxid_TxSize_Ugly(t *testing.T) {
	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: outpoint(0x77, 0), Sequence: 0xffffffff}},
		Outputs: []primitives.Output{{Value: 1, Address: primitives.Address{Version: 0, Hash: make([]byte, 20)}}},
	}
	base, _ := SerializeNoWitness(tx)
	size, err := TxSize(tx)
	if err != nil {
		t.Fatalf("TxSize: %v", err)
	}
	// One input with an empty witness stack adds exactly one varint(0) byte.
	if size != len(base)+1 {
		t.Fatalf("TxSize = %d, want base+1 = %d", size, len(base)+1)
	}
}
