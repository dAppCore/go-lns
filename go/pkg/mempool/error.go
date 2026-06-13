// SPDX-License-Identifier: EUPL-1.2

package mempool

import core "dappco.re/go"

// Sentinel causes returned (wrapped by core.E) when Add rejects a transaction.
// Match them with errors.Is via core.E's Unwrap chain:
//
//	if errors.Is(err, mempool.ErrDoubleSpend) { /* conflicting input */ }
var (
	// ErrDuplicate is returned when the transaction is already in the pool.
	ErrDuplicate = core.NewError("transaction already in pool")

	// ErrEmptyTx is returned when the transaction has no inputs or no outputs.
	ErrEmptyTx = core.NewError("transaction has no inputs or outputs")

	// ErrCoinbase is returned when a coinbase transaction is offered to the
	// pool. Coinbases are created by miners in blocks, never relayed.
	ErrCoinbase = core.NewError("coinbase transaction not allowed in pool")

	// ErrDuplicateInput is returned when a transaction spends the same outpoint
	// twice within itself.
	ErrDuplicateInput = core.NewError("transaction spends an outpoint twice")

	// ErrMissingInputs is returned when an input is neither an unspent chain coin
	// nor an output of a pool transaction (an orphan; not parked by this pool).
	ErrMissingInputs = core.NewError("transaction inputs not available")

	// ErrDoubleSpend is returned when an input is already claimed by another
	// transaction in the pool, or already spent on chain.
	ErrDoubleSpend = core.NewError("transaction conflicts with a pooled or spent input")

	// ErrNegativeFee is returned when output value exceeds input value.
	ErrNegativeFee = core.NewError("transaction outputs exceed inputs")

	// ErrScript is returned when an input fails witness/script verification.
	ErrScript = core.NewError("transaction input failed script verification")

	// ErrTooLarge is returned when a single transaction is larger than the whole
	// pool byte budget and so could never be admitted.
	ErrTooLarge = core.NewError("transaction larger than pool capacity")
)
