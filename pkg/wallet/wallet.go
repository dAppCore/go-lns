// SPDX-License-Identifier: EUPL-1.2

// Package wallet provides a lightweight deterministic wallet scaffold for the
// LNS sidechain. It focuses on seed management, receive-address derivation, and
// local balance/history tracking until full chain-backed spend construction
// lands.
package wallet

import (
	"crypto/rand"
	"crypto/sha3"
	"encoding/hex"
	"slices"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

// EntryKind identifies a wallet history record.
type EntryKind string

const (
	// EntryCredit records an incoming amount.
	EntryCredit EntryKind = "credit"
	// EntryDebit records an outgoing amount.
	EntryDebit EntryKind = "debit"
)

// Entry describes a local wallet history item.
type Entry struct {
	Kind    EntryKind       `json:"kind"`
	Amount  uint64          `json:"amount"`
	Locked  bool            `json:"locked"`
	Pending bool            `json:"pending,omitempty"`
	Height  uint32          `json:"height"`
	Hash    primitives.Hash `json:"hash"`
	Note    string          `json:"note,omitempty"`
}

// Snapshot is a serializable wallet view.
type Snapshot struct {
	SeedHex      string             `json:"seedHex"`
	ReceiveIndex uint32             `json:"receiveIndex"`
	Address      primitives.Address `json:"address"`
}

// Wallet holds a deterministic seed and local ledger entries.
type Wallet struct {
	seed         [32]byte
	receiveIndex uint32
	history      []Entry
}

// Create constructs a wallet from fresh random seed material.
func Create() (*Wallet, error) {
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		return nil, core.E("wallet.Create", "generate seed", err)
	}
	return &Wallet{seed: seed}, nil
}

// ImportSeed reconstructs a wallet from a hex-encoded 32-byte seed.
func ImportSeed(seedHex string) (*Wallet, error) {
	raw, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, core.E("wallet.ImportSeed", "decode seed", err)
	}
	if len(raw) != 32 {
		return nil, core.E("wallet.ImportSeed", "seed must be 32 bytes", nil)
	}

	var seed [32]byte
	copy(seed[:], raw)
	return &Wallet{seed: seed}, nil
}

// ExportSeed returns the hex-encoded seed.
func (w *Wallet) ExportSeed() string {
	if w == nil {
		return ""
	}
	return hex.EncodeToString(w.seed[:])
}

// Snapshot returns a serializable view of the wallet.
func (w *Wallet) Snapshot() Snapshot {
	if w == nil {
		return Snapshot{}
	}
	return Snapshot{
		SeedHex:      w.ExportSeed(),
		ReceiveIndex: w.receiveIndex,
		Address:      w.ReceiveAddress(),
	}
}

// MarshalJSON emits the wallet snapshot for API responses.
func (w *Wallet) MarshalJSON() ([]byte, error) {
	r := core.JSONMarshal(w.Snapshot())
	if !r.OK {
		return nil, r.Value.(error)
	}
	return r.Value.([]byte), nil
}

// ReceiveAddress returns the current deterministic receive address.
func (w *Wallet) ReceiveAddress() primitives.Address {
	if w == nil {
		return primitives.Address{}
	}
	hash := sha3.Sum256(append(w.seed[:], byte(w.receiveIndex>>24), byte(w.receiveIndex>>16), byte(w.receiveIndex>>8), byte(w.receiveIndex)))
	return primitives.Address{
		Version: 0,
		Hash:    append([]byte(nil), hash[:20]...),
	}
}

// NextReceiveAddress increments the receive index and returns the next address.
func (w *Wallet) NextReceiveAddress() primitives.Address {
	if w == nil {
		return primitives.Address{}
	}
	w.receiveIndex++
	return w.ReceiveAddress()
}

// Credit records an incoming amount in local history.
func (w *Wallet) Credit(amount uint64, locked bool, height uint32, txHash primitives.Hash, note string) {
	if w == nil {
		return
	}
	w.history = append(w.history, Entry{
		Kind:   EntryCredit,
		Amount: amount,
		Locked: locked,
		Height: height,
		Hash:   txHash,
		Note:   note,
	})
}

// Debit records an outgoing amount in local history.
func (w *Wallet) Debit(amount uint64, locked bool, height uint32, txHash primitives.Hash, note string) {
	if w == nil {
		return
	}
	w.history = append(w.history, Entry{
		Kind:   EntryDebit,
		Amount: amount,
		Locked: locked,
		Height: height,
		Hash:   txHash,
		Note:   note,
	})
}

// PendingDebit records a local draft outgoing transfer without mutating
// confirmed or locked balances.
func (w *Wallet) PendingDebit(amount uint64, txHash primitives.Hash, note string) Entry {
	entry := Entry{
		Kind:    EntryDebit,
		Amount:  amount,
		Pending: true,
		Hash:    txHash,
		Note:    note,
	}
	if w == nil {
		return entry
	}

	w.history = append(w.history, entry)
	return entry
}

// Balance returns confirmed and locked totals from the local ledger.
func (w *Wallet) Balance() (confirmed uint64, locked uint64) {
	if w == nil {
		return 0, 0
	}

	for _, entry := range w.history {
		if entry.Pending {
			continue
		}

		sign := int64(1)
		if entry.Kind == EntryDebit {
			sign = -1
		}

		if entry.Locked {
			locked = applySignedAmount(locked, sign, entry.Amount)
			continue
		}

		confirmed = applySignedAmount(confirmed, sign, entry.Amount)
	}

	return confirmed, locked
}

func applySignedAmount(current uint64, sign int64, amount uint64) uint64 {
	if sign >= 0 {
		return current + amount
	}
	if amount >= current {
		return 0
	}
	return current - amount
}

// History returns the local history entries ordered by newest first.
func (w *Wallet) History() []Entry {
	if w == nil {
		return nil
	}
	history := append([]Entry(nil), w.history...)
	slices.Reverse(history)
	return history
}

// SeedID returns a short deterministic identifier for the wallet seed.
func (w *Wallet) SeedID() string {
	if w == nil {
		return ""
	}
	hash := sha3.Sum256(w.seed[:])
	return core.Sprintf("%x", hash[:6])
}
