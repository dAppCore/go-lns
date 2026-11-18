// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	core "github.com/dappcore/go"

	"github.com/dappcore/go-lns/go/pkg/mempool"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// TemplateSource is the subset of mempool.Pool a template needs: a fee-ordered
// snapshot of candidate transactions. *mempool.Pool satisfies it. Accepting the
// interface keeps the package testable without a full pool.
//
//	var src mining.TemplateSource = pool
type TemplateSource interface {
	// GetForBlock returns up to limit transactions ordered by descending fee
	// rate; a non-positive limit returns all of them.
	GetForBlock(limit int) []primitives.Transaction
}

// TemplateOptions configures template assembly.
//
//	opts := mining.TemplateOptions{Mempool: pool, PrevBlock: tip, Height: h,
//	    Bits: bits, Coinbase: cb, Time: now}
type TemplateOptions struct {
	// Mempool supplies fee-ordered candidate transactions. If nil, the template
	// contains only the coinbase.
	Mempool TemplateSource

	// Coinbase is the block's coinbase transaction. It MUST be a coinbase (its
	// first input spends the null outpoint). Assembly does not build it.
	Coinbase primitives.Transaction

	// PrevBlock is the hash of the block this one extends.
	PrevBlock primitives.Hash

	// TreeRoot is the name/UTXO Urkel tree root at this height. Assembly does
	// not recompute it.
	TreeRoot primitives.Hash

	// ReservedRoot is the reserved root field (zero unless a soft fork uses it).
	ReservedRoot primitives.Hash

	// Height is the block height, used only for the template's record; it is not
	// written into the header (HSD carries height in the coinbase locktime).
	Height uint32

	// Version is the block version. Zero is accepted.
	Version uint32

	// Bits is the compact difficulty target to mine against. Assembly does not
	// retarget.
	Bits uint32

	// Time is the block timestamp in Unix seconds.
	Time uint64

	// Mask is the optional 32-byte PoW mask for pooled mining. Nil means no
	// mask (the PoW hash equals the share hash).
	Mask []byte

	// ExtraNonce is the optional 24-byte starting extra nonce. Nil means zero.
	ExtraNonce []byte

	// MaxWeight caps the assembled block weight. Zero uses MaxBlockWeight.
	MaxWeight int

	// MaxSize caps the assembled base (non-witness) block size in bytes. Zero
	// uses MaxBlockSize.
	MaxSize int
}

// Template is an assembled, unsolved block plus the accounting from assembly: it
// pairs the candidate Block with the height it targets and the totals a caller
// uses to size the coinbase reward (fees) and report capacity (weight, size).
//
//	tmpl, err := mining.NewTemplate(opts)
//	blk := tmpl.Block()
type Template struct {
	// block is the assembled candidate (coinbase first), roots filled, nonce
	// unsolved.
	block Block

	// Height is the block height this template targets.
	Height uint32

	// Weight is the assembled block weight.
	Weight int

	// BaseSize is the assembled base (non-witness) size in bytes.
	BaseSize int

	// Included is the number of non-coinbase transactions selected from the
	// mempool. Assembly has no chain view, so it cannot total absolute fees;
	// callers that need an exact coinbase reward compute it from their own view.
	Included int
}

// Block returns a copy of the assembled candidate block. The header roots are
// filled and the timestamp is set; the nonce and extra nonce are left for the
// miner to search.
//
//	blk := tmpl.Block()
func (t *Template) Block() Block {
	out := Block{Header: t.block.Header}
	out.Txs = make([]primitives.Transaction, len(t.block.Txs))
	copy(out.Txs, t.block.Txs)
	return out
}

// NewTemplate assembles a block template: it places the coinbase first, then
// drains fee-ordered transactions from the mempool while they fit within the
// size and weight limits, computes the merkle and witness roots over the chosen
// set, and builds the header. The returned template is ready to mine with Solve.
//
// Selection is greedy by fee rate (the mempool's GetForBlock order): each
// candidate is added if it keeps both the base size at or under MaxSize and the
// weight at or under MaxWeight, otherwise it is skipped and the next is tried.
// The coinbase is always included and counts against both limits.
//
//	tmpl, err := mining.NewTemplate(mining.TemplateOptions{
//	    Mempool: pool, PrevBlock: tip, Height: h, Bits: bits,
//	    Coinbase: cb, Time: now,
//	})
func NewTemplate(opts TemplateOptions) (*Template, error) {
	if !isCoinbase(opts.Coinbase) {
		return nil, core.E("mining.NewTemplate", "coinbase transaction is required and must spend the null outpoint", nil)
	}

	maxSize := opts.MaxSize
	if maxSize <= 0 {
		maxSize = MaxBlockSize
	}
	maxWeight := opts.MaxWeight
	if maxWeight <= 0 {
		maxWeight = MaxBlockWeight
	}

	txs := []primitives.Transaction{opts.Coinbase}

	cbBase, cbWitness, err := txSizes(opts.Coinbase)
	if err != nil {
		return nil, core.E("mining.NewTemplate", "size coinbase", err)
	}

	// The block body framing: header + the transaction-count varint. The count
	// is recomputed as transactions are added so the varint width is honest.
	baseSize := HeaderSize + cbBase
	witnessSize := cbWitness
	included := 0

	var candidates []primitives.Transaction
	if opts.Mempool != nil {
		candidates = opts.Mempool.GetForBlock(0)
	}

	for i := range candidates {
		tx := candidates[i]
		if isCoinbase(tx) {
			// A second coinbase is invalid; skip it.
			continue
		}

		base, witness, err := txSizes(tx)
		if err != nil {
			return nil, core.E("mining.NewTemplate", "size transaction", err)
		}

		newBase := baseSize + base
		newWitness := witnessSize + witness
		if newBase > maxSize {
			continue
		}
		if blockWeight(newBase, newWitness) > maxWeight {
			continue
		}

		txs = append(txs, tx)
		baseSize = newBase
		witnessSize = newWitness
		included++
	}

	// Account for the transaction-count varint in the framing now that the count
	// is known.
	baseSize += varintSize(uint64(len(txs)))

	merkleRoot, err := MerkleRoot(txs)
	if err != nil {
		return nil, core.E("mining.NewTemplate", "merkle root", err)
	}
	witnessRoot, err := WitnessRoot(txs)
	if err != nil {
		return nil, core.E("mining.NewTemplate", "witness root", err)
	}

	header := primitives.BlockHeader{
		Version:      opts.Version,
		PrevBlock:    opts.PrevBlock,
		MerkleRoot:   merkleRoot,
		WitnessRoot:  witnessRoot,
		TreeRoot:     opts.TreeRoot,
		ReservedRoot: opts.ReservedRoot,
		Time:         opts.Time,
		Bits:         opts.Bits,
		Nonce:        0,
		ExtraNonce:   fixedBytes(opts.ExtraNonce, NonceSize),
		Mask:         fixedBytes(opts.Mask, maskSize),
	}

	return &Template{
		block:    Block{Header: header, Txs: txs},
		Height:   opts.Height,
		Weight:   blockWeight(baseSize, witnessSize),
		BaseSize: baseSize,
		Included: included,
	}, nil
}

// isCoinbase reports whether a transaction is a coinbase: it has at least one
// input and the first spends the null outpoint.
func isCoinbase(tx primitives.Transaction) bool {
	return len(tx.Inputs) > 0 && tx.Inputs[0].IsCoinbase()
}

// txSizes returns a transaction's base (non-witness) and witness byte sizes,
// matching the HSD split used for weight.
func txSizes(tx primitives.Transaction) (base, witness int, err error) {
	baseBytes, err := mempool.SerializeNoWitness(tx)
	if err != nil {
		return 0, 0, core.E("mining.txSizes", "serialise base", err)
	}
	base = len(baseBytes)
	witness = len(serializeWitnessData(tx))
	return base, witness, nil
}

// blockWeight computes block weight from base and witness sizes:
// weight = base*(WitnessScaleFactor-1) + (base + witness). Witness bytes are
// discounted, matching HSD getWeight.
func blockWeight(base, witness int) int {
	return base*(WitnessScaleFactor-1) + base + witness
}

// varintSize returns the encoded length of a compact-size integer.
func varintSize(n uint64) int {
	switch {
	case n < 0xfd:
		return 1
	case n <= 0xffff:
		return 3
	case n <= 0xffffffff:
		return 5
	default:
		return 9
	}
}
