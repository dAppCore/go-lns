// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"time"

	core "github.com/dappcore/go"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// Entry is a transaction held in the pool together with the metadata the pool
// orders and evicts by: its id, serialised size, the fee it pays, and the time
// it was admitted. It mirrors the fields of an HSD MempoolEntry that the core
// pool needs (lib/mempool/mempoolentry.js), without the ancestor/descendant
// package accounting that this pool does not model.
//
//	entry := pool.GetEntry(txid)
//	rate := entry.FeeRate() // dollarydoo per 1000 bytes
type Entry struct {
	// TxID is the transaction id (non-witness BLAKE2b-256 digest).
	TxID TxID

	// Tx is the transaction itself.
	Tx primitives.Transaction

	// Size is the serialised byte size including witnesses.
	Size int

	// Fee is the fee in dollarydoos: sum(input values) - sum(output values).
	Fee uint64

	// Time is when the entry was admitted to the pool.
	Time time.Time
}

// newEntry builds an Entry for a transaction, computing its id, size, and fee
// from the spent coins resolved by view. It returns ErrMissingInputs if any
// input cannot be resolved and ErrNegativeFee if the outputs exceed the inputs.
func newEntry(tx primitives.Transaction, view CoinView, now time.Time) (*Entry, error) {
	id, err := ComputeTxID(tx)
	if err != nil {
		return nil, err
	}

	size, err := TxSize(tx)
	if err != nil {
		return nil, err
	}

	fee, err := computeFee(tx, view)
	if err != nil {
		return nil, err
	}

	return &Entry{
		TxID: id,
		Tx:   tx,
		Size: size,
		Fee:  fee,
		Time: now,
	}, nil
}

// computeFee returns sum(input values) - sum(output values), resolving each
// input's value through the coin view. It errors if an input is unresolved
// (ErrMissingInputs) or the fee would be negative (ErrNegativeFee). Coinbase
// transactions have no resolvable inputs and are rejected by the pool before
// this is reached.
func computeFee(tx primitives.Transaction, view CoinView) (uint64, error) {
	if view == nil {
		return 0, core.E("mempool.computeFee", "nil coin view", ErrMissingInputs)
	}

	var in uint64
	for i := range tx.Inputs {
		coin, ok := view.GetOutput(tx.Inputs[i].Prevout)
		if !ok {
			return 0, core.E("mempool.computeFee", "input not found in view", ErrMissingInputs)
		}
		// Guard the running input sum against overflow before adding.
		if coin.Value > ^uint64(0)-in {
			return 0, core.E("mempool.computeFee", "input value overflow", ErrNegativeFee)
		}
		in += coin.Value
	}

	var out uint64
	for i := range tx.Outputs {
		v := tx.Outputs[i].Value
		if v > ^uint64(0)-out {
			return 0, core.E("mempool.computeFee", "output value overflow", ErrNegativeFee)
		}
		out += v
	}

	if out > in {
		return 0, core.E("mempool.computeFee", "outputs exceed inputs", ErrNegativeFee)
	}

	return in - out, nil
}

// FeeRate returns the fee rate in dollarydoos per 1000 bytes, the value the pool
// orders block templates and eviction by. It is floor(fee * 1000 / size),
// matching policy.getRate in the HSD reference. A zero-size entry rates as zero.
//
//	rate := entry.FeeRate()
func (e *Entry) FeeRate() uint64 {
	if e == nil || e.Size <= 0 {
		return 0
	}
	return e.Fee * 1000 / uint64(e.Size)
}

// Age returns how long the entry has been in the pool relative to now.
//
//	if entry.Age(time.Now()) > time.Hour { ... }
func (e *Entry) Age(now time.Time) time.Duration {
	if e == nil {
		return 0
	}
	return now.Sub(e.Time)
}
