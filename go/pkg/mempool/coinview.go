// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"dappco.re/go/lns/pkg/primitives"
)

// CoinView is the confirmed-chain UTXO lookup the pool validates spends against.
// It returns the output an outpoint refers to and whether that output is
// currently unspent (present in the chain's coin set).
//
// It is intentionally the same shape as covenant.CoinView so a single chain coin
// view satisfies both the covenant rules and the mempool. A spent or unknown
// outpoint returns ok=false; the pool treats that as "not spendable from the
// chain" and falls back to its own pool outputs (see PoolView).
//
//	coin, ok := view.GetOutput(input.Prevout)
type CoinView interface {
	// GetOutput returns the unspent output for an outpoint and whether it exists.
	GetOutput(primitives.Outpoint) (primitives.Output, bool)
}

// PoolView layers the pool's own unconfirmed outputs over a confirmed CoinView.
// It lets a transaction spend an output created by another transaction already
// in the pool: GetOutput first checks the pool's transactions, then falls back
// to the chain view. It is the viewpoint the pool builds for script and covenant
// verification, mirroring HSD mempool._getSpentView.
//
//	view := pool.SpentView(chainView)
//	coin, ok := view.GetOutput(outpoint)
type PoolView struct {
	pool  *Pool
	chain CoinView
}

// GetOutput resolves an outpoint against pool transactions first, then the
// confirmed chain view. An outpoint already spent by another pool transaction is
// still resolvable here (it returns the underlying output); double-spend
// rejection is the pool's job via its claimed-outpoint index, not the view's.
func (v PoolView) GetOutput(outpoint primitives.Outpoint) (primitives.Output, bool) {
	if v.pool != nil {
		if out, ok := v.pool.outputFor(outpoint); ok {
			return out, true
		}
	}

	if v.chain != nil {
		return v.chain.GetOutput(outpoint)
	}

	return primitives.Output{}, false
}

// SpentSet is an in-memory CoinView backed by a map of outpoints to outputs. It
// is the simplest confirmed-chain stand-in: callers (and tests) register the
// outputs a transaction is allowed to spend, and Spend removes one to model
// confirmation or a competing spend.
//
//	set := mempool.NewSpentSet()
//	set.Add(outpoint, output)
//	coin, ok := set.GetOutput(outpoint) // ok == true until Spend(outpoint)
type SpentSet struct {
	coins map[outpointKey]primitives.Output
}

// NewSpentSet returns an empty coin set.
func NewSpentSet() *SpentSet {
	return &SpentSet{coins: make(map[outpointKey]primitives.Output)}
}

// Add registers an unspent output for an outpoint, replacing any existing entry.
//
//	set.Add(primitives.Outpoint{TxHash: h, Index: 0}, out)
func (s *SpentSet) Add(outpoint primitives.Outpoint, out primitives.Output) {
	if s.coins == nil {
		s.coins = make(map[outpointKey]primitives.Output)
	}
	s.coins[keyForOutpoint(outpoint)] = out
}

// Spend removes an outpoint from the set, modelling its confirmation in a block
// or a competing spend. It reports whether the outpoint was present.
//
//	set.Spend(outpoint)
func (s *SpentSet) Spend(outpoint primitives.Outpoint) bool {
	key := keyForOutpoint(outpoint)
	if _, ok := s.coins[key]; !ok {
		return false
	}
	delete(s.coins, key)
	return true
}

// GetOutput implements CoinView.
func (s *SpentSet) GetOutput(outpoint primitives.Outpoint) (primitives.Output, bool) {
	out, ok := s.coins[keyForOutpoint(outpoint)]
	return out, ok
}

// Len returns the number of unspent outputs tracked.
func (s *SpentSet) Len() int {
	return len(s.coins)
}

// outpointKey is a comparable map key for an outpoint: the 32-byte transaction
// hash followed by the 4-byte index. Using an array (not a string) keeps the
// index allocation-free.
type outpointKey [36]byte

// keyForOutpoint derives the map key for an outpoint.
func keyForOutpoint(o primitives.Outpoint) outpointKey {
	var key outpointKey
	copy(key[:32], o.TxHash[:])
	key[32] = byte(o.Index)
	key[33] = byte(o.Index >> 8)
	key[34] = byte(o.Index >> 16)
	key[35] = byte(o.Index >> 24)
	return key
}
