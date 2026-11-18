// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	"testing"

	"github.com/dappcore/go-lns/go/pkg/mempool"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// templateOpts builds a baseline TemplateOptions with a valid coinbase and easy
// difficulty, overridden per test.
func templateOpts(src TemplateSource) TemplateOptions {
	return TemplateOptions{
		Mempool:   src,
		Coinbase:  coinbaseTx(100, 2000),
		PrevBlock: hashOf(0xaa),
		TreeRoot:  hashOf(0xdd),
		Height:    100,
		Version:   1,
		Bits:      easyBits,
		Time:      1_700_000_000,
	}
}

// TestTemplate_NewTemplate_Good selects fee-ordered transactions, places the
// coinbase first, and fills the header roots.
func TestTemplate_NewTemplate_Good(t *testing.T) {
	src := fixedSource{txs: []primitives.Transaction{
		spendTx(0x10, 1000, 4),
		spendTx(0x20, 900, 4),
		spendTx(0x30, 800, 4),
	}}

	tmpl, err := NewTemplate(templateOpts(src))
	if err != nil {
		t.Fatalf("NewTemplate: %v", err)
	}

	blk := tmpl.Block()
	if blk.TxCount() != 4 {
		t.Fatalf("block has %d txs, want 4 (coinbase + 3)", blk.TxCount())
	}

	cb, ok := blk.Coinbase()
	if !ok || !isCoinbase(cb) {
		t.Fatal("first transaction is not a coinbase")
	}
	if tmpl.Included != 3 {
		t.Fatalf("Included = %d, want 3", tmpl.Included)
	}

	// Roots must be filled and consistent with the assembled set.
	wantMerkle, err := MerkleRoot(blk.Txs)
	if err != nil {
		t.Fatalf("MerkleRoot: %v", err)
	}
	if blk.Header.MerkleRoot != wantMerkle {
		t.Fatal("header merkle root does not match the assembled transactions")
	}
	if blk.Header.MerkleRoot == (primitives.Hash{}) {
		t.Fatal("header merkle root is zero")
	}
}

// TestTemplate_NewTemplate_Bad rejects assembly without a valid coinbase.
func TestTemplate_NewTemplate_Bad(t *testing.T) {
	opts := templateOpts(fixedSource{})
	opts.Coinbase = spendTx(0x10, 1000, 4) // not a coinbase (spends a real outpoint)

	if _, err := NewTemplate(opts); err == nil {
		t.Fatal("NewTemplate accepted a non-coinbase as the coinbase")
	}
}

// TestTemplate_NewTemplate_Ugly assembles from an empty mempool: the block holds
// only the coinbase, and the roots are still well-formed.
func TestTemplate_NewTemplate_Ugly(t *testing.T) {
	tmpl, err := NewTemplate(templateOpts(nil))
	if err != nil {
		t.Fatalf("NewTemplate empty: %v", err)
	}

	blk := tmpl.Block()
	if blk.TxCount() != 1 {
		t.Fatalf("empty-mempool block has %d txs, want 1 (coinbase only)", blk.TxCount())
	}
	if tmpl.Included != 0 {
		t.Fatalf("Included = %d, want 0", tmpl.Included)
	}
	if blk.Header.MerkleRoot == (primitives.Hash{}) {
		t.Fatal("coinbase-only merkle root is zero")
	}
}

// TestTemplate_SizeLimit skips transactions that would breach the base-size cap,
// confirming selection respects limits rather than blindly draining the pool.
func TestTemplate_SizeLimit(t *testing.T) {
	// Three small txs; a base cap that admits the coinbase plus exactly one.
	src := fixedSource{txs: []primitives.Transaction{
		spendTx(0x10, 1000, 4),
		spendTx(0x20, 900, 4),
		spendTx(0x30, 800, 4),
	}}

	opts := templateOpts(src)

	// Measure the footprint so the cap is meaningful: header + coinbase base +
	// one spend base + the count varint.
	cbBase, _, err := txSizes(opts.Coinbase)
	if err != nil {
		t.Fatalf("txSizes coinbase: %v", err)
	}
	spendBase, _, err := txSizes(src.txs[0])
	if err != nil {
		t.Fatalf("txSizes spend: %v", err)
	}
	// Allow the coinbase and one spend, but not two.
	opts.MaxSize = HeaderSize + cbBase + spendBase + 1

	tmpl, err := NewTemplate(opts)
	if err != nil {
		t.Fatalf("NewTemplate: %v", err)
	}
	if tmpl.Included != 1 {
		t.Fatalf("Included = %d, want 1 under the size cap", tmpl.Included)
	}
	if got := tmpl.Block().TxCount(); got != 2 {
		t.Fatalf("block has %d txs, want 2 (coinbase + 1)", got)
	}
}

// TestTemplate_SkipsExtraCoinbase ignores a second coinbase offered by the
// mempool, since a block may have only one.
func TestTemplate_SkipsExtraCoinbase(t *testing.T) {
	src := fixedSource{txs: []primitives.Transaction{
		coinbaseTx(100, 5000), // a second coinbase — must be skipped
		spendTx(0x10, 1000, 4),
	}}

	tmpl, err := NewTemplate(templateOpts(src))
	if err != nil {
		t.Fatalf("NewTemplate: %v", err)
	}
	if tmpl.Included != 1 {
		t.Fatalf("Included = %d, want 1 (the spend only)", tmpl.Included)
	}
	blk := tmpl.Block()
	for i := 1; i < blk.TxCount(); i++ {
		if isCoinbase(blk.Txs[i]) {
			t.Fatalf("transaction %d is a second coinbase", i)
		}
	}
}

// poolImplementsSource is a compile-time assertion that the real mempool.Pool
// satisfies TemplateSource, so a production template can be fed directly from the
// pool. The pool's own tests cover the fee-rate ordering GetForBlock produces;
// this package's contract is to preserve that order, which the next test checks.
var _ TemplateSource = (*mempool.Pool)(nil)

// TestTemplate_PreservesSourceOrder confirms the template emits transactions in
// exactly the order the source returns them (the mempool returns them fee-rate
// ordered), so block position tracks the source ranking.
func TestTemplate_PreservesSourceOrder(t *testing.T) {
	// The source presents three distinct spends in a fixed order, standing in for
	// a fee-rate-ordered mempool snapshot.
	ordered := []primitives.Transaction{
		spendTx(0x51, 1000, 4),
		spendTx(0x52, 900, 4),
		spendTx(0x53, 800, 4),
	}
	tmpl, err := NewTemplate(templateOpts(fixedSource{txs: ordered}))
	if err != nil {
		t.Fatalf("NewTemplate: %v", err)
	}

	blk := tmpl.Block()
	if blk.TxCount() != 4 {
		t.Fatalf("block has %d txs, want 4 (coinbase + 3)", blk.TxCount())
	}

	for i, want := range ordered {
		wantID, _ := mempool.ComputeTxID(want)
		gotID, _ := mempool.ComputeTxID(blk.Txs[i+1]) // +1 to skip the coinbase
		if gotID != wantID {
			t.Fatalf("position %d: source order not preserved", i)
		}
	}
}

// TestTemplate_MineMarshalVerify is the full block-production round-trip: build a
// template, mine it at trivial difficulty, marshal the block, and verify both
// the PoW and that the marshalled bytes are non-trivial.
func TestTemplate_MineMarshalVerify(t *testing.T) {
	src := fixedSource{txs: []primitives.Transaction{
		spendTx(0x10, 1000, 4),
		spendTx(0x20, 900, 8),
	}}

	tmpl, err := NewTemplate(templateOpts(src))
	if err != nil {
		t.Fatalf("NewTemplate: %v", err)
	}

	blk := tmpl.Block()
	blk.Header = mustSolve(t, blk.Header, 1<<16)

	ok, err := VerifyPOW(blk.Header)
	if err != nil {
		t.Fatalf("VerifyPOW: %v", err)
	}
	if !ok {
		t.Fatal("mined block failed PoW verification")
	}

	raw, err := blk.Marshal()
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if len(raw) <= HeaderSize {
		t.Fatalf("marshalled block is %d bytes, expected more than the %d-byte header",
			len(raw), HeaderSize)
	}

	hash, err := blk.Hash()
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if hash == (primitives.Hash{}) {
		t.Fatal("block hash is zero")
	}
}
