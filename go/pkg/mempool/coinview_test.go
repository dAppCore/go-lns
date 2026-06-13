// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

// TestCoinview_SpentSet_Good checks add, lookup, and spend semantics.
func TestCoinview_SpentSet_Good(t *testing.T) {
	set := NewSpentSet()
	op := outpoint(0xA1, 0)
	out := primitives.Output{Value: 1000, Address: primitives.Address{Version: 0, Hash: make([]byte, 20)}}

	set.Add(op, out)
	if set.Len() != 1 {
		t.Fatalf("Len = %d, want 1", set.Len())
	}

	got, ok := set.GetOutput(op)
	if !ok {
		t.Fatal("output not found after Add")
	}
	if got.Value != 1000 {
		t.Fatalf("value = %d, want 1000", got.Value)
	}

	if !set.Spend(op) {
		t.Fatal("Spend returned false for a present outpoint")
	}
	if _, ok := set.GetOutput(op); ok {
		t.Fatal("output still present after Spend")
	}
}

// TestCoinview_SpentSet_Bad checks lookups for absent / re-spent outpoints.
func TestCoinview_SpentSet_Bad(t *testing.T) {
	set := NewSpentSet()
	op := outpoint(0xA2, 7)

	if _, ok := set.GetOutput(op); ok {
		t.Fatal("found an output that was never added")
	}
	if set.Spend(op) {
		t.Fatal("Spend returned true for an absent outpoint")
	}
}

// TestCoinview_SpentSet_Ugly checks the index distinguishes outpoints by index,
// not just transaction hash.
func TestCoinview_SpentSet_Ugly(t *testing.T) {
	set := NewSpentSet()
	var h primitives.Hash
	for i := range h {
		h[i] = 0xCC
	}
	op0 := primitives.Outpoint{TxHash: h, Index: 0}
	op1 := primitives.Outpoint{TxHash: h, Index: 1}

	set.Add(op0, primitives.Output{Value: 1})
	set.Add(op1, primitives.Output{Value: 2})

	a, _ := set.GetOutput(op0)
	b, _ := set.GetOutput(op1)
	if a.Value == b.Value {
		t.Fatal("same-hash outpoints with different indexes collided")
	}
}

// TestCoinview_PoolView_Good checks pool outputs layer over the chain view.
func TestCoinview_PoolView_Good(t *testing.T) {
	pool := New(DefaultConfig())

	// Fund and admit a parent transaction; its output becomes spendable in-pool.
	parentOp := outpoint(0xB1, 0)
	addr, priv, chain := fundedCoin(t, 10, 5000, parentOp)
	parent := signedSpend(t, priv, parentOp, 5000, 4000, addr)
	parentID, err := pool.Add(parent, chain)
	if err != nil {
		t.Fatalf("Add parent: %v", err)
	}

	view := PoolView{pool: pool, chain: chain}

	// The parent's output (index 0) resolves through the pool, not the chain.
	childPrev := primitives.Outpoint{TxHash: parentID, Index: 0}
	got, ok := view.GetOutput(childPrev)
	if !ok {
		t.Fatal("pool output not resolvable through PoolView")
	}
	if got.Value != 4000 {
		t.Fatalf("pool output value = %d, want 4000", got.Value)
	}

	// A chain coin still resolves through the fallback.
	if _, ok := view.GetOutput(parentOp); !ok {
		// parentOp was spent on chain? No — SpentSet still holds it (pool does not
		// mutate the chain view). It must resolve.
		t.Fatal("chain coin not resolvable through PoolView")
	}
}

// TestCoinview_PoolView_Bad checks unknown outpoints resolve to not-found.
func TestCoinview_PoolView_Bad(t *testing.T) {
	pool := New(DefaultConfig())
	view := PoolView{pool: pool, chain: NewSpentSet()}

	if _, ok := view.GetOutput(outpoint(0xB2, 0)); ok {
		t.Fatal("resolved an outpoint that exists nowhere")
	}
}

// TestCoinview_PoolView_Ugly checks a nil chain view does not panic.
func TestCoinview_PoolView_Ugly(t *testing.T) {
	view := PoolView{pool: nil, chain: nil}
	if _, ok := view.GetOutput(outpoint(0xB3, 0)); ok {
		t.Fatal("nil view resolved an outpoint")
	}
}
