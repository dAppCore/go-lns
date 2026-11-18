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

	core "github.com/dappcore/go"
	"github.com/dappcore/go-lns/go/pkg/hd"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// defaultAddressGap is the forward look-ahead applied when scanning derived
// addresses for owned coins or signing keys. It mirrors the conventional BIP-44
// gap limit so funds sent to a freshly-derived (not-yet-recorded) address are
// still discovered.
const defaultAddressGap uint32 = 20

// errNilWallet is the message used when an operation is attempted on a nil
// wallet receiver.
const errNilWallet = "wallet is not initialized"

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
//
// The address is derived through pkg/hd using the real HSD key scheme: the
// 32-byte seed seeds a BIP-32 master key, and the receive index selects the
// BIP-44 external-chain address "m/44'/5353'/0'/0/receiveIndex". The result is a
// version-0 P2WPKH primitives.Address whose 20-byte hash is the BLAKE2b-160
// digest of the compressed secp256k1 public key. Derivation is deterministic in
// the seed and index, so importing the same seed reproduces the same address.
func (w *Wallet) ReceiveAddress() primitives.Address {
	if w == nil {
		return primitives.Address{}
	}
	return w.deriveReceiveAddress(w.receiveIndex)
}

// deriveReceiveAddress derives the external-chain (receive) address at the given
// index via pkg/hd. A derivation failure (cryptographically improbable for a
// 32-byte seed) yields the zero Address, matching the nil-receiver contract.
func (w *Wallet) deriveReceiveAddress(index uint32) primitives.Address {
	master, err := hd.NewMaster(w.seed[:])
	if err != nil {
		return primitives.Address{}
	}

	addr, err := hd.DeriveAddress(master, hd.Path{
		CoinType: hd.CoinTypeLethean,
		Account:  0,
		Branch:   hd.BranchReceive,
		Index:    index,
	})
	if err != nil {
		return primitives.Address{}
	}
	return addr
}

// NextReceiveAddress increments the receive index and returns the next address.
func (w *Wallet) NextReceiveAddress() primitives.Address {
	if w == nil {
		return primitives.Address{}
	}
	w.receiveIndex++
	return w.ReceiveAddress()
}

// KeyRef pairs a derived address with its derivation coordinates so a caller can
// match a coin (by its holding address) to the key that signs it. Branch is 0
// for receive addresses and 1 for change addresses.
type KeyRef struct {
	// Address is the version-0 P2WPKH address derived at Branch/Index.
	Address primitives.Address
	// Branch is the BIP-44 branch: 0 receive, 1 change.
	Branch uint32
	// Index is the address index within the branch.
	Index uint32
}

// SpendableAddresses enumerates the wallet's owned addresses across the receive
// and change branches, covering indices [0, receiveIndex] plus a forward gap. It
// is the address set whose coins the wallet may spend: a caller queries the
// daemon for coins on these addresses, selects inputs, and resolves each input's
// signing key with SigningKey. The gap extends the scan past the current receive
// index so funds sent to a not-yet-used address (within the gap) are still
// found; a zero gap falls back to a small default.
//
//	refs := wallet.SpendableAddresses(0)
func (w *Wallet) SpendableAddresses(gap uint32) []KeyRef {
	if w == nil {
		return nil
	}
	if gap == 0 {
		gap = defaultAddressGap
	}

	master, err := hd.NewMaster(w.seed[:])
	if err != nil {
		return nil
	}

	last := w.receiveIndex + gap
	refs := make([]KeyRef, 0, int(last+1)*2)
	for _, branch := range []uint32{hd.BranchReceive, hd.BranchChange} {
		for index := uint32(0); index <= last; index++ {
			addr, err := hd.DeriveAddress(master, hd.Path{
				CoinType: hd.CoinTypeLethean,
				Account:  0,
				Branch:   branch,
				Index:    index,
			})
			if err != nil {
				continue
			}
			refs = append(refs, KeyRef{Address: addr, Branch: branch, Index: index})
		}
	}
	return refs
}

// SigningKey derives the private signing key for one of the wallet's addresses,
// scanning the receive and change branches up to receiveIndex plus a gap to find
// the matching derivation. It returns the 32-byte secp256k1 scalar and the
// 33-byte compressed public key, which together produce the P2WPKH witness for a
// coin held at addr. It fails if addr is not a derivable wallet address within
// the scanned range.
//
//	priv, pub, err := wallet.SigningKey(coinAddress, 0)
func (w *Wallet) SigningKey(addr primitives.Address, gap uint32) (priv []byte, pub []byte, err error) {
	if w == nil {
		return nil, nil, core.E("wallet.SigningKey", errNilWallet, nil)
	}
	if !addr.IsPubkeyHash() {
		return nil, nil, core.E("wallet.SigningKey", "address is not a version-0 pubkey hash", nil)
	}
	if gap == 0 {
		gap = defaultAddressGap
	}

	master, err := hd.NewMaster(w.seed[:])
	if err != nil {
		return nil, nil, core.E("wallet.SigningKey", "derive master key", err)
	}

	last := w.receiveIndex + gap
	for _, branch := range []uint32{hd.BranchReceive, hd.BranchChange} {
		for index := uint32(0); index <= last; index++ {
			key, derr := hd.Derive(master, hd.Path{
				CoinType: hd.CoinTypeLethean,
				Account:  0,
				Branch:   branch,
				Index:    index,
			})
			if derr != nil {
				continue
			}
			if !key.Address().Equals(addr) {
				continue
			}
			privBytes, perr := key.PrivateKeyBytes()
			if perr != nil {
				return nil, nil, core.E("wallet.SigningKey", "extract private key", perr)
			}
			pubBytes, perr := key.PublicKeyBytes()
			if perr != nil {
				return nil, nil, core.E("wallet.SigningKey", "extract public key", perr)
			}
			return privBytes, pubBytes, nil
		}
	}

	return nil, nil, core.E("wallet.SigningKey", "address is not owned by this wallet", nil)
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
