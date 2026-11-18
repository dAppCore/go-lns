// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"testing"
	"time"

	core "github.com/dappcore/go"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// TestPool_Add_Good admits a valid, fully-signed transaction and indexes it.
func TestPool_Add_Good(t *testing.T) {
	pool := New(DefaultConfig())

	op := outpoint(0x01, 0)
	addr, priv, chain := fundedCoin(t, 1, 5000, op)
	tx := signedSpend(t, priv, op, 5000, 4000, addr)

	id, err := pool.Add(tx, chain)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if !pool.Has(id) {
		t.Fatal("Has returned false after Add")
	}
	if pool.Size() != 1 {
		t.Fatalf("Size = %d, want 1", pool.Size())
	}
	if pool.Bytes() <= 0 {
		t.Fatal("Bytes should be positive after Add")
	}

	got, ok := pool.Get(id)
	if !ok {
		t.Fatal("Get returned false after Add")
	}
	if len(got.Outputs) != 1 || got.Outputs[0].Value != 4000 {
		t.Fatal("Get returned the wrong transaction")
	}

	entry, ok := pool.GetEntry(id)
	if !ok {
		t.Fatal("GetEntry returned false")
	}
	if entry.Fee != 1000 { // 5000 in - 4000 out
		t.Fatalf("entry fee = %d, want 1000", entry.Fee)
	}

	// The spent outpoint is now claimed.
	if !pool.IsSpent(op) {
		t.Fatal("input outpoint not marked spent after Add")
	}
}

// TestPool_Add_Bad rejects structurally invalid and duplicate transactions.
func TestPool_Add_Bad(t *testing.T) {
	pool := New(DefaultConfig())

	// Empty transaction.
	if _, err := pool.Add(primitives.Transaction{}, NewSpentSet()); !core.Is(err, ErrEmptyTx) {
		t.Fatalf("expected ErrEmptyTx, got %v", err)
	}

	// Coinbase (null prevout) is not relayed.
	cb := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: primitives.NewOutpoint()}},
		Outputs: []primitives.Output{{Value: 1}},
	}
	if _, err := pool.Add(cb, NewSpentSet()); !core.Is(err, ErrCoinbase) {
		t.Fatalf("expected ErrCoinbase, got %v", err)
	}

	// Same outpoint spent twice within one transaction.
	op := outpoint(0x02, 0)
	dup := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: op, Sequence: 0xffffffff},
			{Prevout: op, Sequence: 0xffffffff},
		},
		Outputs: []primitives.Output{{Value: 1}},
	}
	if _, err := pool.Add(dup, NewSpentSet()); !core.Is(err, ErrDuplicateInput) {
		t.Fatalf("expected ErrDuplicateInput, got %v", err)
	}

	// Duplicate admission of the same valid transaction.
	addr, priv, chain := fundedCoin(t, 2, 5000, outpoint(0x03, 0))
	tx := signedSpend(t, priv, outpoint(0x03, 0), 5000, 4000, addr)
	if _, err := pool.Add(tx, chain); err != nil {
		t.Fatalf("first Add: %v", err)
	}
	if _, err := pool.Add(tx, chain); !core.Is(err, ErrDuplicate) {
		t.Fatalf("expected ErrDuplicate, got %v", err)
	}
}

// TestPool_Add_Ugly rejects double-spends (pool + chain), missing inputs, bad
// scripts, negative fees, and over-capacity transactions.
func TestPool_Add_Ugly(t *testing.T) {
	// --- Double-spend within the pool: two txs spending the same chain coin. ---
	pool := New(DefaultConfig())
	op := outpoint(0x10, 0)
	addr, priv, chain := fundedCoin(t, 3, 5000, op)

	txA := signedSpend(t, priv, op, 5000, 4000, addr)
	txB := signedSpend(t, priv, op, 5000, 3000, addr) // same input, different output

	if _, err := pool.Add(txA, chain); err != nil {
		t.Fatalf("Add txA: %v", err)
	}
	if _, err := pool.Add(txB, chain); !core.Is(err, ErrDoubleSpend) {
		t.Fatalf("expected ErrDoubleSpend (pool conflict), got %v", err)
	}

	// --- Double-spend against the chain: the coin is already spent on chain. ---
	pool2 := New(DefaultConfig())
	op2 := outpoint(0x11, 0)
	addr2, priv2, chain2 := fundedCoin(t, 4, 5000, op2)
	tx2 := signedSpend(t, priv2, op2, 5000, 4000, addr2)
	chain2.Spend(op2) // coin no longer unspent on chain
	if _, err := pool2.Add(tx2, chain2); !core.Is(err, ErrMissingInputs) {
		t.Fatalf("expected ErrMissingInputs for a chain-spent coin, got %v", err)
	}

	// --- Invalid script: valid structure, wrong witness pubkey. ---
	pool3 := New(DefaultConfig())
	op3 := outpoint(0x12, 0)
	addr3, priv3, chain3 := fundedCoin(t, 5, 5000, op3)
	badScript := signedSpend(t, priv3, op3, 5000, 4000, addr3)
	_, wrongPub := testKey(t, 222)
	badScript.Inputs[0].Witness = [][]byte{badScript.Inputs[0].Witness[0], wrongPub}
	if _, err := pool3.Add(badScript, chain3); !core.Is(err, ErrScript) {
		t.Fatalf("expected ErrScript for a bad witness, got %v", err)
	}

	// --- Negative fee: outputs exceed the spent coin. ---
	pool4 := New(DefaultConfig())
	op4 := outpoint(0x13, 0)
	addr4, priv4, chain4 := fundedCoin(t, 6, 1000, op4)
	negFee := signedSpend(t, priv4, op4, 1000, 2000, addr4) // out 2000 > in 1000
	if _, err := pool4.Add(negFee, chain4); !core.Is(err, ErrNegativeFee) {
		t.Fatalf("expected ErrNegativeFee, got %v", err)
	}

	// --- Over capacity: a transaction larger than the whole pool. ---
	cfg := DefaultConfig()
	cfg.MaxSize = 50 // bytes; smaller than any real tx
	pool5 := New(cfg)
	op5 := outpoint(0x14, 0)
	addr5, priv5, chain5 := fundedCoin(t, 7, 5000, op5)
	big := signedSpend(t, priv5, op5, 5000, 4000, addr5)
	if _, err := pool5.Add(big, chain5); !core.Is(err, ErrTooLarge) {
		t.Fatalf("expected ErrTooLarge, got %v", err)
	}
}

// TestPool_GetForBlock_Good orders transactions by descending fee rate.
func TestPool_GetForBlock_Good(t *testing.T) {
	pool := New(DefaultConfig())

	// Three independent transactions with increasing fee (same size → same order
	// by fee rate). Fees: 200, 600, 1000 → rates ascending in that order.
	type spec struct {
		seed byte
		tag  byte
		fee  uint64
	}
	specs := []spec{{20, 0x20, 200}, {21, 0x21, 600}, {22, 0x22, 1000}}
	ids := make([]TxID, len(specs))
	for i, s := range specs {
		op := outpoint(s.tag, 0)
		addr, priv, chain := fundedCoin(t, s.seed, 5000, op)
		tx := signedSpend(t, priv, op, 5000, 5000-s.fee, addr)
		id, err := pool.Add(tx, chain)
		if err != nil {
			t.Fatalf("Add %d: %v", i, err)
		}
		ids[i] = id
	}

	ordered := pool.GetForBlock(0)
	if len(ordered) != 3 {
		t.Fatalf("GetForBlock returned %d txs, want 3", len(ordered))
	}

	// Highest fee (1000) must come first, lowest (200) last.
	firstID, _ := ComputeTxID(ordered[0])
	lastID, _ := ComputeTxID(ordered[len(ordered)-1])
	if firstID != ids[2] {
		t.Fatal("highest-fee-rate transaction is not first")
	}
	if lastID != ids[0] {
		t.Fatal("lowest-fee-rate transaction is not last")
	}

	// Limit is honoured.
	top := pool.GetForBlock(1)
	if len(top) != 1 {
		t.Fatalf("GetForBlock(1) returned %d, want 1", len(top))
	}
	topID, _ := ComputeTxID(top[0])
	if topID != ids[2] {
		t.Fatal("GetForBlock(1) did not return the highest-fee-rate transaction")
	}
}

// TestPool_GetForBlock_Bad checks an empty pool yields an empty template.
func TestPool_GetForBlock_Bad(t *testing.T) {
	pool := New(DefaultConfig())
	if got := pool.GetForBlock(10); len(got) != 0 {
		t.Fatalf("empty pool returned %d txs", len(got))
	}
	if got := pool.Snapshot(); len(got) != 0 {
		t.Fatalf("empty pool snapshot returned %d txs", len(got))
	}
}

// TestPool_GetForBlock_Ugly checks a limit larger than the pool returns all, and
// the returned slice is a copy that does not alias pool state.
func TestPool_GetForBlock_Ugly(t *testing.T) {
	pool := New(DefaultConfig())
	op := outpoint(0x30, 0)
	addr, priv, chain := fundedCoin(t, 30, 5000, op)
	tx := signedSpend(t, priv, op, 5000, 4000, addr)
	if _, err := pool.Add(tx, chain); err != nil {
		t.Fatalf("Add: %v", err)
	}

	got := pool.GetForBlock(1000)
	if len(got) != 1 {
		t.Fatalf("over-large limit returned %d, want 1", len(got))
	}

	// Mutating the returned copy must not corrupt the pool's stored tx.
	got[0].Outputs[0].Value = 999999
	stored, _ := pool.Get(tx.Inputs[0].Prevout.TxHash) // wrong id on purpose: ignore
	_ = stored
	id, _ := ComputeTxID(tx)
	again, _ := pool.Get(id)
	if again.Outputs[0].Value != 4000 {
		t.Fatal("mutating the snapshot corrupted pool state")
	}
}

// TestPool_Remove_Good removes a transaction and frees its claimed inputs.
func TestPool_Remove_Good(t *testing.T) {
	pool := New(DefaultConfig())
	op := outpoint(0x40, 0)
	addr, priv, chain := fundedCoin(t, 40, 5000, op)
	tx := signedSpend(t, priv, op, 5000, 4000, addr)
	id, err := pool.Add(tx, chain)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if !pool.Remove(id) {
		t.Fatal("Remove returned false for a present tx")
	}
	if pool.Has(id) {
		t.Fatal("tx still present after Remove")
	}
	if pool.IsSpent(op) {
		t.Fatal("input still claimed after Remove")
	}
	if pool.Bytes() != 0 {
		t.Fatalf("Bytes = %d after removing the only tx, want 0", pool.Bytes())
	}

	// After removal the same input may be respent (e.g. by a replacement).
	if _, err := pool.Add(tx, chain); err != nil {
		t.Fatalf("re-Add after Remove: %v", err)
	}
}

// TestPool_Remove_Bad checks removing an absent id is a no-op.
func TestPool_Remove_Bad(t *testing.T) {
	pool := New(DefaultConfig())
	if pool.Remove(outpoint(0x41, 0).TxHash) {
		t.Fatal("Remove returned true for an absent id")
	}
}

// TestPool_Remove_Ugly checks removal of one of several txs leaves the rest intact.
func TestPool_Remove_Ugly(t *testing.T) {
	pool := New(DefaultConfig())
	var keep TxID
	var drop TxID
	for i, tag := range []byte{0x50, 0x51} {
		op := outpoint(tag, 0)
		addr, priv, chain := fundedCoin(t, byte(50+i), 5000, op)
		tx := signedSpend(t, priv, op, 5000, 4000, addr)
		id, err := pool.Add(tx, chain)
		if err != nil {
			t.Fatalf("Add %d: %v", i, err)
		}
		if i == 0 {
			drop = id
		} else {
			keep = id
		}
	}

	pool.Remove(drop)
	if pool.Has(drop) {
		t.Fatal("dropped tx still present")
	}
	if !pool.Has(keep) {
		t.Fatal("kept tx was removed")
	}
	if pool.Size() != 1 {
		t.Fatalf("Size = %d after removing one of two, want 1", pool.Size())
	}
}

// TestPool_Expire_Good evicts transactions older than the expiry window.
func TestPool_Expire_Good(t *testing.T) {
	now := time.Unix(100000, 0)
	cfg := DefaultConfig()
	cfg.Expiry = time.Hour
	cfg.now = func() time.Time { return now }
	pool := New(cfg)

	op := outpoint(0x60, 0)
	addr, priv, chain := fundedCoin(t, 60, 5000, op)
	tx := signedSpend(t, priv, op, 5000, 4000, addr)
	id, err := pool.Add(tx, chain)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	// Advance the clock past the expiry window.
	now = now.Add(2 * time.Hour)
	dropped := pool.Expire()
	if dropped != 1 {
		t.Fatalf("Expire dropped %d, want 1", dropped)
	}
	if pool.Has(id) {
		t.Fatal("expired tx still present")
	}
}

// TestPool_Expire_Bad checks fresh transactions are not expired.
func TestPool_Expire_Bad(t *testing.T) {
	now := time.Unix(200000, 0)
	cfg := DefaultConfig()
	cfg.Expiry = time.Hour
	cfg.now = func() time.Time { return now }
	pool := New(cfg)

	op := outpoint(0x61, 0)
	addr, priv, chain := fundedCoin(t, 61, 5000, op)
	tx := signedSpend(t, priv, op, 5000, 4000, addr)
	if _, err := pool.Add(tx, chain); err != nil {
		t.Fatalf("Add: %v", err)
	}

	now = now.Add(30 * time.Minute) // still within the window
	if dropped := pool.Expire(); dropped != 0 {
		t.Fatalf("Expire dropped %d fresh txs, want 0", dropped)
	}
}

// TestPool_Expire_Ugly checks a zero expiry disables the sweep entirely.
func TestPool_Expire_Ugly(t *testing.T) {
	now := time.Unix(300000, 0)
	cfg := DefaultConfig()
	cfg.Expiry = 0 // disabled
	cfg.now = func() time.Time { return now }
	pool := New(cfg)

	op := outpoint(0x62, 0)
	addr, priv, chain := fundedCoin(t, 62, 5000, op)
	tx := signedSpend(t, priv, op, 5000, 4000, addr)
	if _, err := pool.Add(tx, chain); err != nil {
		t.Fatalf("Add: %v", err)
	}

	now = now.Add(10000 * time.Hour)
	if dropped := pool.Expire(); dropped != 0 {
		t.Fatalf("zero-expiry pool dropped %d, want 0", dropped)
	}
	if pool.Size() != 1 {
		t.Fatal("zero-expiry pool lost its transaction")
	}
}

// TestPool_Add_EvictsLowestFee checks overflow evicts the lowest-fee-rate tx and
// keeps the just-added higher-fee one. This exercises Add's limit() path.
func TestPool_Add_EvictsLowestFee(t *testing.T) {
	// Size a pool to hold roughly two of these ~200-byte transactions, forcing the
	// third to trigger eviction down to the 90% threshold.
	op0 := outpoint(0x70, 0)
	addr0, priv0, chain0 := fundedCoin(t, 70, 5000, op0)
	probe := signedSpend(t, priv0, op0, 5000, 4900, addr0)
	txSize, err := TxSize(probe)
	if err != nil {
		t.Fatalf("TxSize: %v", err)
	}

	cfg := DefaultConfig()
	cfg.MaxSize = txSize*2 + txSize/2 // room for 2 txs, the 3rd overflows
	pool := New(cfg)

	// Admit two low-fee transactions.
	lowFeeIDs := make([]TxID, 0, 2)
	for i, tag := range []byte{0x71, 0x72} {
		op := outpoint(tag, 0)
		addr, priv, chain := fundedCoin(t, byte(71+i), 5000, op)
		tx := signedSpend(t, priv, op, 5000, 4900, addr) // fee 100
		id, err := pool.Add(tx, chain)
		if err != nil {
			t.Fatalf("Add low %d: %v", i, err)
		}
		lowFeeIDs = append(lowFeeIDs, id)
	}

	// Admit a high-fee transaction that pushes the pool over its budget.
	opHi := outpoint(0x73, 0)
	addrHi, privHi, chainHi := fundedCoin(t, 73, 5000, opHi)
	hi := signedSpend(t, privHi, opHi, 5000, 1000, addrHi) // fee 4000 (much denser)
	hiID, err := pool.Add(hi, chainHi)
	if err != nil {
		t.Fatalf("Add high: %v", err)
	}

	// The high-fee tx must survive; at least one low-fee tx must have been evicted.
	if !pool.Has(hiID) {
		t.Fatal("high-fee transaction was evicted")
	}
	if pool.Bytes() > cfg.MaxSize {
		t.Fatalf("pool still over budget: %d > %d", pool.Bytes(), cfg.MaxSize)
	}
	evicted := 0
	for _, id := range lowFeeIDs {
		if !pool.Has(id) {
			evicted++
		}
	}
	if evicted == 0 {
		t.Fatal("no low-fee transaction was evicted on overflow")
	}
	_ = chain0
}

// TestPool_Add_ChainedSpend checks a tx may spend a pooled tx's output.
func TestPool_Add_ChainedSpend(t *testing.T) {
	pool := New(DefaultConfig())

	// Parent spends a chain coin to an address we also control via priv2.
	parentOp := outpoint(0x80, 0)
	_, priv1, chain := fundedCoin(t, 80, 10000, parentOp)
	// Send the parent's output to priv2's P2WPKH address.
	priv2, pub2 := testKey(t, 81)
	childAddr := primitives.Address{Version: 0, Hash: blake160(t, pub2)}
	parent := signedSpend(t, priv1, parentOp, 10000, 9000, childAddr)
	parentID, err := pool.Add(parent, chain)
	if err != nil {
		t.Fatalf("Add parent: %v", err)
	}

	// Child spends the parent's output (index 0, value 9000) — resolved in-pool.
	childPrev := primitives.Outpoint{TxHash: parentID, Index: 0}
	child := signedSpend(t, priv2, childPrev, 9000, 8000, childAddr)
	if _, err := pool.Add(child, chain); err != nil {
		t.Fatalf("Add child (chained spend): %v", err)
	}
	if pool.Size() != 2 {
		t.Fatalf("Size = %d, want 2", pool.Size())
	}
}
