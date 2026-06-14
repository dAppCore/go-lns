// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"encoding/hex"
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

func witnessTestTx() primitives.Transaction {
	return primitives.Transaction{
		Version: 0,
		Inputs: []primitives.Input{
			{
				Prevout:  primitives.Outpoint{TxHash: primitives.Hash{1, 2, 3}, Index: 0},
				Sequence: 0xffffffff,
				Witness:  [][]byte{{0xaa, 0xbb}, {0xcc}},
			},
		},
		Outputs: []primitives.Output{
			{Value: 1000, Address: primitives.Address{Version: 0, Hash: make([]byte, 20)}},
		},
		LockTime: 0,
	}
}

// TestSerializeWitness_Good checks the witness serialisation begins with the
// non-witness serialisation and appends the witness stack bytes.
func TestSerializeWitness_Good(t *testing.T) {
	tx := witnessTestTx()

	full, err := SerializeWithWitness(tx)
	if err != nil {
		t.Fatalf("SerializeWithWitness error: %v", err)
	}
	base, err := SerializeNoWitness(tx)
	if err != nil {
		t.Fatalf("SerializeNoWitness error: %v", err)
	}

	if len(full) <= len(base) {
		t.Fatalf("witness serialisation (%d) not longer than base (%d)", len(full), len(base))
	}
	if hex.EncodeToString(full[:len(base)]) != hex.EncodeToString(base) {
		t.Fatalf("witness serialisation does not start with the base serialisation")
	}

	// The witness tail: varint(2 items) ‖ varint(2)+aabb ‖ varint(1)+cc.
	tail := full[len(base):]
	want := []byte{0x02, 0x02, 0xaa, 0xbb, 0x01, 0xcc}
	if hex.EncodeToString(tail) != hex.EncodeToString(want) {
		t.Fatalf("witness tail = %x, want %x", tail, want)
	}
}

// TestSerializeWitness_Bad covers the empty-witness case: a tx whose input has no
// witness items still serialises, with a zero item-count for that input.
func TestSerializeWitness_Bad(t *testing.T) {
	tx := witnessTestTx()
	tx.Inputs[0].Witness = nil

	full, err := SerializeWithWitness(tx)
	if err != nil {
		t.Fatalf("SerializeWithWitness error: %v", err)
	}
	base, err := SerializeNoWitness(tx)
	if err != nil {
		t.Fatalf("SerializeNoWitness error: %v", err)
	}
	if len(full) != len(base)+1 {
		t.Fatalf("empty-witness serialisation = %d bytes, want base+1 (%d)", len(full), len(base)+1)
	}
	if full[len(base)] != 0x00 {
		t.Fatalf("empty witness item count = %d, want 0", full[len(base)])
	}
}

// TestSerializeWitness_Ugly confirms the txid is witness-independent: changing
// the witness leaves SerializeNoWitness (and thus the txid) unchanged while
// SerializeWithWitness differs.
func TestSerializeWitness_Ugly(t *testing.T) {
	a := witnessTestTx()
	b := witnessTestTx()
	b.Inputs[0].Witness = [][]byte{{0x11}}

	baseA, _ := SerializeNoWitness(a)
	baseB, _ := SerializeNoWitness(b)
	if hex.EncodeToString(baseA) != hex.EncodeToString(baseB) {
		t.Fatalf("non-witness serialisation changed with the witness")
	}

	fullA, _ := SerializeWithWitness(a)
	fullB, _ := SerializeWithWitness(b)
	if hex.EncodeToString(fullA) == hex.EncodeToString(fullB) {
		t.Fatalf("witness serialisation did not change with the witness")
	}
}
