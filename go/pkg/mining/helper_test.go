// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"testing"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// hashOf builds a deterministic 32-byte hash from a seed byte, for test
// transactions and roots.
func hashOf(seed byte) primitives.Hash {
	var h primitives.Hash
	for i := range h {
		h[i] = seed
	}
	return h
}

// coinbaseTx builds a minimal but structurally valid coinbase: one input
// spending the null outpoint, one output. The locktime carries the height the
// way HSD coinbases do.
func coinbaseTx(height uint32, value uint64) primitives.Transaction {
	return primitives.Transaction{
		Version: 0,
		Inputs: []primitives.Input{{
			Prevout:  primitives.Outpoint{TxHash: primitives.Hash{}, Index: ^uint32(0)},
			Sequence: ^uint32(0),
			Witness:  [][]byte{{0x00}},
		}},
		Outputs: []primitives.Output{{
			Value:   value,
			Address: primitives.Address{Version: 0, Hash: make([]byte, 20)},
		}},
		LockTime: height,
	}
}

// spendTx builds a non-coinbase transaction spending a distinct prior outpoint,
// seeded so its txid differs per seed. The witness length lets a test vary the
// witness footprint.
func spendTx(seed byte, value uint64, witnessLen int) primitives.Transaction {
	addr := primitives.Address{Version: 0, Hash: make([]byte, 20)}
	addr.Hash[0] = seed

	witness := make([]byte, witnessLen)
	for i := range witness {
		witness[i] = seed
	}

	return primitives.Transaction{
		Version: 0,
		Inputs: []primitives.Input{{
			Prevout:  primitives.Outpoint{TxHash: hashOf(seed), Index: 0},
			Sequence: 0xffffffff,
			Witness:  [][]byte{witness},
		}},
		Outputs: []primitives.Output{{
			Value:   value,
			Address: addr,
		}},
		LockTime: 0,
	}
}

// fixedSource is a TemplateSource that returns a preset transaction list in the
// order given, standing in for a fee-ordered mempool snapshot.
type fixedSource struct {
	txs []primitives.Transaction
}

func (s fixedSource) GetForBlock(limit int) []primitives.Transaction {
	if limit > 0 && limit < len(s.txs) {
		out := make([]primitives.Transaction, limit)
		copy(out, s.txs)
		return out
	}
	out := make([]primitives.Transaction, len(s.txs))
	copy(out, s.txs)
	return out
}

// easyBits is a compact target so wide that essentially any hash meets it, so
// Solve terminates on the first nonce. Exponent 0x21 (33) with mantissa 0x00ffff
// gives a 256-bit-class target near the maximum.
const easyBits uint32 = 0x2100ffff

// hardBits is a compact target so tight (exponent 0x03, tiny mantissa) that a
// random hash will not meet it, used to exercise rejection without a long
// search.
const hardBits uint32 = 0x03000001

// mustSolve fails the test if the header cannot be solved within rounds.
func mustSolve(t *testing.T, hdr primitives.BlockHeader, rounds uint32) primitives.BlockHeader {
	t.Helper()
	nonce, ok := Solve(hdr, TargetFromBits(hdr.Bits), rounds)
	if !ok {
		t.Fatalf("Solve failed to find a nonce within %d rounds at bits %#x", rounds, hdr.Bits)
	}
	hdr.Nonce = nonce
	return hdr
}
