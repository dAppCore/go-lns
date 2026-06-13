// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"sort"
	"sync"
	"time"

	core "dappco.re/go"

	"dappco.re/go/lns/pkg/primitives"
	"dappco.re/go/lns/pkg/script"
)

// HSD policy defaults (lib/protocol/policy.js), reused so the pool's limits match
// the chain this repository ports.
const (
	// DefaultMaxSize is the default pool byte budget: 100 MB.
	DefaultMaxSize = 100 * 1000 * 1000

	// DefaultMaxCount is the default cap on the number of pooled transactions.
	// HSD bounds the pool by bytes alone; the count cap is an additional guard so
	// a flood of tiny transactions cannot exhaust map memory.
	DefaultMaxCount = 100000

	// DefaultExpiry is how long a transaction may sit in the pool before it is
	// evicted as stale: 72 hours (policy.MEMPOOL_EXPIRY_TIME).
	DefaultExpiry = 72 * time.Hour

	// evictionThresholdNum/Den set the size the pool drains down to once it
	// overflows: 90% of MaxSize (HSD: maxSize - maxSize/10), so eviction is
	// amortised rather than run on every single Add.
	evictionThresholdNum = 9
	evictionThresholdDen = 10
)

// Config tunes a Pool's limits. The zero value is not usable; start from
// DefaultConfig and override fields as needed.
//
//	cfg := mempool.DefaultConfig()
//	cfg.MaxSize = 8 * 1000 * 1000
//	pool := mempool.New(cfg)
type Config struct {
	// MaxSize is the pool byte budget. When the total serialised size of pooled
	// transactions exceeds it, lowest-fee-rate transactions are evicted.
	MaxSize int

	// MaxCount is the maximum number of pooled transactions. Zero means no count
	// limit (bytes still apply).
	MaxCount int

	// Expiry is how long a transaction may remain pooled. Zero disables expiry.
	Expiry time.Duration

	// VerifyFlags are the script verification flags applied to each input.
	// Defaults to script.StandardVerifyFlags (the policy/mempool flag set).
	VerifyFlags int

	// now returns the current time. Tests inject a clock here; production leaves
	// it nil and the pool uses time.Now.
	now func() time.Time
}

// DefaultConfig returns the standard pool configuration: 100 MB, 72-hour expiry,
// standard (policy) script flags.
//
//	pool := mempool.New(mempool.DefaultConfig())
func DefaultConfig() Config {
	return Config{
		MaxSize:     DefaultMaxSize,
		MaxCount:    DefaultMaxCount,
		Expiry:      DefaultExpiry,
		VerifyFlags: script.StandardVerifyFlags,
	}
}

// Pool is a local pool of validated, unconfirmed transactions. It is safe for
// concurrent use. Transactions are indexed by id and by every outpoint they
// claim, so duplicate-detection, double-spend-detection, and removal are all
// O(1) in the number of inputs.
//
//	pool := mempool.New(mempool.DefaultConfig())
//	id, err := pool.Add(tx, chainView)
type Pool struct {
	cfg Config

	mu      sync.RWMutex
	entries map[TxID]*Entry
	// claimed maps each spent outpoint to the id of the pool transaction that
	// claims it, giving O(1) double-spend detection within the pool.
	claimed map[outpointKey]TxID
	bytes   int
}

// New creates a pool from a configuration. Missing limits fall back to their
// defaults so a partially-filled Config is still usable.
//
//	pool := mempool.New(mempool.Config{MaxSize: 4 * 1000 * 1000})
func New(cfg Config) *Pool {
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = DefaultMaxSize
	}
	if cfg.Expiry < 0 {
		cfg.Expiry = 0
	}
	if cfg.VerifyFlags == 0 {
		cfg.VerifyFlags = script.StandardVerifyFlags
	}

	return &Pool{
		cfg:     cfg,
		entries: make(map[TxID]*Entry),
		claimed: make(map[outpointKey]TxID),
	}
}

// Add validates a transaction and admits it to the pool, returning its id.
//
// Validation, in order:
//
//  1. Structure: at least one input and one output, not a coinbase, no input
//     spent twice within the transaction.
//  2. Duplicate: the id is not already pooled (ErrDuplicate).
//  3. Inputs available: every input resolves to a pool output or an unspent
//     chain coin (ErrMissingInputs); none is already claimed by another pool
//     transaction or spent on chain (ErrDoubleSpend).
//  4. Fee: outputs do not exceed inputs (ErrNegativeFee).
//  5. Script: every input's witness verifies against the spent output
//     (ErrScript).
//  6. Capacity: the transaction is not larger than the whole pool (ErrTooLarge).
//
// On success the transaction is indexed and the pool is trimmed to its limits
// (expired and lowest-fee-rate transactions evicted). chainView is the confirmed
// UTXO set; pass nil only if every input is satisfied by another pooled
// transaction.
//
//	id, err := pool.Add(tx, chainView)
//	if errors.Is(err, mempool.ErrDoubleSpend) { /* conflict */ }
func (p *Pool) Add(tx primitives.Transaction, chainView CoinView) (TxID, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := checkStructure(tx); err != nil {
		return TxID{}, err
	}

	id, err := ComputeTxID(tx)
	if err != nil {
		return TxID{}, err
	}

	if _, ok := p.entries[id]; ok {
		return TxID{}, core.E("mempool.Add", "duplicate transaction", ErrDuplicate)
	}

	// Build the viewpoint: pool outputs layered over the confirmed chain.
	view := PoolView{pool: p, chain: chainView}

	// Reject conflicts and unresolved inputs before doing the expensive script
	// checks. An input must resolve in the view (chain coin or pool output) and
	// must not already be claimed by another pool transaction.
	for i := range tx.Inputs {
		prevout := tx.Inputs[i].Prevout
		key := keyForOutpoint(prevout)

		if _, ok := p.claimed[key]; ok {
			return TxID{}, core.E("mempool.Add", "input already claimed in pool", ErrDoubleSpend)
		}

		if _, ok := view.GetOutput(prevout); !ok {
			return TxID{}, core.E("mempool.Add", "input not available", ErrMissingInputs)
		}
	}

	entry, err := newEntry(tx, view, p.clock())
	if err != nil {
		return TxID{}, err
	}

	if entry.Size > p.cfg.MaxSize {
		return TxID{}, core.E("mempool.Add", "transaction exceeds pool capacity", ErrTooLarge)
	}

	if err := p.verifyScripts(tx, view); err != nil {
		return TxID{}, err
	}

	p.insert(entry)
	p.limit(id)

	return id, nil
}

// verifyScripts runs witness verification for every non-coinbase input against
// the output it spends, using the pool's configured flags. The first failing
// input returns ErrScript wrapping the underlying script error.
func (p *Pool) verifyScripts(tx primitives.Transaction, view CoinView) error {
	for i := range tx.Inputs {
		coin, ok := view.GetOutput(tx.Inputs[i].Prevout)
		if !ok {
			return core.E("mempool.verifyScripts", "input not found", ErrMissingInputs)
		}
		if err := script.VerifyInput(coin, tx, i, p.cfg.VerifyFlags); err != nil {
			// Surface ErrScript as the cause so callers can match it with
			// errors.Is, while the message preserves the underlying script error.
			return core.E("mempool.verifyScripts", core.Sprintf("input %d: %s", i, err.Error()), ErrScript)
		}
	}
	return nil
}

// insert adds an entry to the indexes and updates the byte total. The caller
// holds the write lock and has already confirmed the inputs are unclaimed.
func (p *Pool) insert(entry *Entry) {
	p.entries[entry.TxID] = entry
	p.bytes += entry.Size
	for i := range entry.Tx.Inputs {
		p.claimed[keyForOutpoint(entry.Tx.Inputs[i].Prevout)] = entry.TxID
	}
}

// removeLocked deletes an entry by id, freeing its claimed outpoints and byte
// count. It reports whether the entry existed. The caller holds the write lock.
func (p *Pool) removeLocked(id TxID) bool {
	entry, ok := p.entries[id]
	if !ok {
		return false
	}

	delete(p.entries, id)
	p.bytes -= entry.Size
	for i := range entry.Tx.Inputs {
		key := keyForOutpoint(entry.Tx.Inputs[i].Prevout)
		// Only release the claim if it still points at this entry; a later
		// transaction never overwrites an existing claim, but this keeps removal
		// self-consistent.
		if owner, ok := p.claimed[key]; ok && owner == id {
			delete(p.claimed, key)
		}
	}
	return true
}

// Remove deletes a transaction from the pool by id (for example after it is
// confirmed in a block) and reports whether it was present.
//
//	removed := pool.Remove(id)
func (p *Pool) Remove(id TxID) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.removeLocked(id)
}

// limit trims the pool to its limits: it first drops expired transactions, then,
// if still over the byte or count budget, evicts the lowest-fee-rate
// transactions down to the eviction threshold. added is the id just inserted; it
// is never evicted by its own Add call (HSD keeps the just-added entry when it
// can). The caller holds the write lock.
func (p *Pool) limit(added TxID) {
	p.expireLocked()

	if !p.overLimit() {
		return
	}

	// Order all entries except the just-added one by ascending fee rate (ties
	// broken by oldest first), then evict from the front until under threshold.
	victims := make([]*Entry, 0, len(p.entries))
	for id, entry := range p.entries {
		if id == added {
			continue
		}
		victims = append(victims, entry)
	}

	sort.Slice(victims, func(i, j int) bool {
		ri, rj := victims[i].FeeRate(), victims[j].FeeRate()
		if ri != rj {
			return ri < rj
		}
		return victims[i].Time.Before(victims[j].Time)
	})

	threshold := p.cfg.MaxSize * evictionThresholdNum / evictionThresholdDen
	for _, entry := range victims {
		if !p.overThreshold(threshold) {
			break
		}
		p.removeLocked(entry.TxID)
	}

	// If the added transaction alone still pushes us over (every other entry was
	// evicted), drop it too so the invariant holds.
	if p.overLimit() {
		p.removeLocked(added)
	}
}

// expireLocked evicts transactions older than the configured expiry. The caller
// holds the write lock. A zero expiry disables the sweep.
func (p *Pool) expireLocked() {
	if p.cfg.Expiry <= 0 {
		return
	}

	now := p.clock()
	for id, entry := range p.entries {
		if now.Sub(entry.Time) >= p.cfg.Expiry {
			p.removeLocked(id)
		}
	}
}

// overLimit reports whether the pool exceeds either the byte or count budget.
func (p *Pool) overLimit() bool {
	if p.bytes > p.cfg.MaxSize {
		return true
	}
	if p.cfg.MaxCount > 0 && len(p.entries) > p.cfg.MaxCount {
		return true
	}
	return false
}

// overThreshold reports whether the pool is still above the drain-down target
// for either budget.
func (p *Pool) overThreshold(byteThreshold int) bool {
	if p.bytes > byteThreshold {
		return true
	}
	if p.cfg.MaxCount > 0 && len(p.entries) > p.cfg.MaxCount {
		return true
	}
	return false
}

// Expire removes stale transactions immediately and returns the number evicted.
// It is exposed so a caller can run expiry on a timer independent of Add.
//
//	dropped := pool.Expire()
func (p *Pool) Expire() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	before := len(p.entries)
	p.expireLocked()
	return before - len(p.entries)
}

// Has reports whether a transaction id is in the pool.
//
//	if pool.Has(id) { ... }
func (p *Pool) Has(id TxID) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.entries[id]
	return ok
}

// Get returns the transaction for an id and whether it was found.
//
//	tx, ok := pool.Get(id)
func (p *Pool) Get(id TxID) (primitives.Transaction, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	entry, ok := p.entries[id]
	if !ok {
		return primitives.Transaction{}, false
	}
	return entry.Tx, true
}

// GetEntry returns the full pool entry for an id (transaction plus fee, size,
// and admission time) and whether it was found.
//
//	entry, ok := pool.GetEntry(id)
func (p *Pool) GetEntry(id TxID) (*Entry, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	entry, ok := p.entries[id]
	return entry, ok
}

// Size returns the number of transactions in the pool.
//
//	n := pool.Size()
func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.entries)
}

// Bytes returns the total serialised size of all pooled transactions.
//
//	used := pool.Bytes()
func (p *Pool) Bytes() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.bytes
}

// IsSpent reports whether an outpoint is already claimed by a pool transaction.
//
//	if pool.IsSpent(outpoint) { /* would double-spend */ }
func (p *Pool) IsSpent(outpoint primitives.Outpoint) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, ok := p.claimed[keyForOutpoint(outpoint)]
	return ok
}

// Snapshot returns the pooled transactions ordered by descending fee rate
// (highest-paying first), ties broken by oldest first. It is the order a miner
// fills a block from. The slice is a copy; mutating it does not affect the pool.
//
//	for _, tx := range pool.Snapshot() { /* add to block until full */ }
func (p *Pool) Snapshot() []primitives.Transaction {
	return p.GetForBlock(0)
}

// GetForBlock returns up to limit transactions ordered by descending fee rate
// for a block template; a limit of zero or less returns all of them. Ordering
// matches Snapshot. The result is a copy.
//
//	template := pool.GetForBlock(2000) // up to 2000 highest-fee-rate txs
func (p *Pool) GetForBlock(limit int) []primitives.Transaction {
	p.mu.RLock()
	defer p.mu.RUnlock()

	entries := make([]*Entry, 0, len(p.entries))
	for _, entry := range p.entries {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		ri, rj := entries[i].FeeRate(), entries[j].FeeRate()
		if ri != rj {
			return ri > rj
		}
		return entries[i].Time.Before(entries[j].Time)
	})

	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}

	out := make([]primitives.Transaction, len(entries))
	for i := range entries {
		out[i] = copyTx(entries[i].Tx)
	}
	return out
}

// copyTx returns a transaction whose input, output, and witness slices are
// freshly allocated, so a caller may edit the returned transaction (reorder or
// reprice outputs, rebuild a witness) without mutating the entry the pool still
// holds. It backs the copy promise GetForBlock and Snapshot make: the entries
// are stored by value but share their slice backing arrays, so a shallow struct
// copy would alias them.
func copyTx(tx primitives.Transaction) primitives.Transaction {
	dup := tx // copies the scalar fields (version, locktime) directly

	if tx.Inputs != nil {
		dup.Inputs = make([]primitives.Input, len(tx.Inputs))
		copy(dup.Inputs, tx.Inputs)
		for i := range tx.Inputs {
			if w := tx.Inputs[i].Witness; w != nil {
				stack := make([][]byte, len(w))
				copy(stack, w)
				dup.Inputs[i].Witness = stack
			}
		}
	}

	if tx.Outputs != nil {
		dup.Outputs = make([]primitives.Output, len(tx.Outputs))
		copy(dup.Outputs, tx.Outputs)
	}

	return dup
}

// outputFor resolves an outpoint to an output created by a pool transaction. It
// is used by PoolView to let one pooled transaction spend another's output. The
// caller (via PoolView) may hold the read lock; this only reads the entry map.
func (p *Pool) outputFor(outpoint primitives.Outpoint) (primitives.Output, bool) {
	entry, ok := p.entries[outpoint.TxHash]
	if !ok {
		return primitives.Output{}, false
	}
	if int(outpoint.Index) >= len(entry.Tx.Outputs) {
		return primitives.Output{}, false
	}
	return entry.Tx.Outputs[outpoint.Index], true
}

// clock returns the current time using the injected clock or time.Now.
func (p *Pool) clock() time.Time {
	if p.cfg.now != nil {
		return p.cfg.now()
	}
	return time.Now()
}

// checkStructure performs context-free structural validation: a transaction must
// have at least one input and one output, must not be a coinbase, and must not
// spend the same outpoint twice.
func checkStructure(tx primitives.Transaction) error {
	if len(tx.Inputs) == 0 || len(tx.Outputs) == 0 {
		return core.E("mempool.checkStructure", "empty inputs or outputs", ErrEmptyTx)
	}

	for i := range tx.Inputs {
		if tx.Inputs[i].IsCoinbase() {
			return core.E("mempool.checkStructure", "coinbase input", ErrCoinbase)
		}
	}

	seen := make(map[outpointKey]struct{}, len(tx.Inputs))
	for i := range tx.Inputs {
		key := keyForOutpoint(tx.Inputs[i].Prevout)
		if _, ok := seen[key]; ok {
			return core.E("mempool.checkStructure", "duplicate input", ErrDuplicateInput)
		}
		seen[key] = struct{}{}
	}

	return nil
}
