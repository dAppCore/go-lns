// SPDX-License-Identifier: EUPL-1.2

// Package primitives provides the core data types for the Lethean Name System.
// These types mirror the HSD (Handshake) sidechain primitives, converted from
// the JavaScript reference implementation at itns-sidechain/lib/primitives/.
//
// Create an address from a public key hash:
//
//	addr := primitives.Address{Version: 0, Hash: pubkeyHash}
//
// Build a transaction with inputs and outputs:
//
//	tx := primitives.Transaction{Version: 0, LockTime: 0}
//	tx.Inputs = append(tx.Inputs, input)
//	tx.Outputs = append(tx.Outputs, output)
package primitives

// Hash is a 32-byte hash used throughout the name system for block references,
// transaction identifiers, and name hashes.
//
//	var h primitives.Hash  // zero-value is a valid empty hash
type Hash [32]byte

// Address represents a bech32-encoded witness address on the Lethean name chain.
// Version 0 addresses use either 20-byte (pubkey hash) or 32-byte (script hash)
// programmes. Version 31 is reserved for nulldata (unspendable) outputs.
//
//	addr := primitives.Address{Version: 0, Hash: blake2bDigest}
//	if addr.IsPubkeyHash() { /* 20-byte witness programme */ }
type Address struct {
	// Version is the witness programme version (0-31).
	// Version 0 = standard witness, version 31 = nulldata.
	Version uint8

	// Hash is the witness programme hash.
	// 20 bytes for pubkey hash, 32 bytes for script hash at version 0.
	Hash []byte
}

// Inject copies another address into this one.
func (a *Address) Inject(other Address) *Address {
	a.Version = other.Version
	a.Hash = other.Hash
	return a
}

// Clone returns a shallow copy of the address.
func (a Address) Clone() Address {
	var clone Address
	return *clone.Inject(a)
}

// Equals reports whether two addresses have the same version and hash.
func (a Address) Equals(other Address) bool {
	if a.Version != other.Version {
		return false
	}

	if len(a.Hash) != len(other.Hash) {
		return false
	}

	for i := range a.Hash {
		if a.Hash[i] != other.Hash[i] {
			return false
		}
	}

	return true
}

// Compare orders two addresses by version and then by hash bytes.
func (a Address) Compare(other Address) int {
	if a.Version < other.Version {
		return -1
	}

	if a.Version > other.Version {
		return 1
	}

	min := len(a.Hash)
	if len(other.Hash) < min {
		min = len(other.Hash)
	}

	for i := 0; i < min; i++ {
		if a.Hash[i] < other.Hash[i] {
			return -1
		}
		if a.Hash[i] > other.Hash[i] {
			return 1
		}
	}

	switch {
	case len(a.Hash) < len(other.Hash):
		return -1
	case len(a.Hash) > len(other.Hash):
		return 1
	default:
		return 0
	}
}

// IsNull reports whether the address hash is entirely zero bytes.
func (a Address) IsNull() bool {
	for _, b := range a.Hash {
		if b != 0 {
			return false
		}
	}

	return true
}

// IsPubkeyHash reports whether the address is a version 0 P2WPKH programme.
func (a Address) IsPubkeyHash() bool {
	return a.Version == 0 && len(a.Hash) == 20
}

// IsScripthash reports whether the address is a version 0 P2WSH programme.
func (a Address) IsScripthash() bool {
	return a.Version == 0 && len(a.Hash) == 32
}

// IsNulldata reports whether the address is an unspendable nulldata programme.
func (a Address) IsNulldata() bool {
	return a.Version == 31
}

// IsUnspendable reports whether the address cannot be spent.
func (a Address) IsUnspendable() bool {
	return a.IsNulldata()
}

// IsUnknown reports whether the address is a non-standard witness programme.
func (a Address) IsUnknown() bool {
	switch a.Version {
	case 0:
		return len(a.Hash) != 20 && len(a.Hash) != 32
	case 31:
		return false
	default:
		return true
	}
}

// IsValid reports whether the address is a syntactically valid witness
// programme.
func (a Address) IsValid() bool {
	if a.Version > 31 {
		return false
	}

	if len(a.Hash) < 2 || len(a.Hash) > 40 {
		return false
	}

	if a.Version == 0 {
		return len(a.Hash) == 20 || len(a.Hash) == 32
	}

	return true
}

// GetSize returns the size of the address when serialised.
func (a Address) GetSize() int {
	return 1 + 1 + len(a.Hash)
}

// Outpoint references a specific output from a previous transaction.
// Used by transaction inputs to identify the coin being spent.
//
//	ref := primitives.Outpoint{TxHash: prevTxHash, Index: 0}
type Outpoint struct {
	// TxHash is the hash of the transaction containing the referenced output.
	TxHash Hash

	// Index is the zero-based position of the output within that transaction.
	Index uint32
}

// Input represents a transaction input that spends a previous output.
// Each input references an outpoint and carries a witness for validation.
//
//	in := primitives.Input{Prevout: ref, Sequence: 0xffffffff}
type Input struct {
	// Prevout identifies the output being spent.
	Prevout Outpoint

	// Witness contains the witness stack items (signatures, public keys).
	Witness [][]byte

	// Sequence is the input sequence number, used for relative timelocks.
	Sequence uint32
}

// Output represents a transaction output: a value sent to an address,
// optionally carrying a covenant for name operations.
//
//	out := primitives.Output{Value: 1000, Address: addr, Covenant: cov}
type Output struct {
	// Value is the amount in the smallest denomination (dollarydoos).
	Value uint64

	// Address is the destination witness address.
	Address Address

	// Covenant carries name-system operation data (NONE for plain transfers).
	Covenant Covenant
}

// Covenant attaches name-system semantics to a transaction output.
// The type field determines the operation (OPEN, BID, REVEAL, REGISTER, etc.)
// and the items carry operation-specific data (name hashes, heights, records).
//
//	cov := primitives.Covenant{Type: covenant.TypeOpen, Items: [][]byte{nameHash, heightBytes, rawName}}
type Covenant struct {
	// Type is the covenant operation type (see pkg/covenant for constants).
	Type uint8

	// Items is the ordered list of data payloads for this covenant.
	// Contents vary by type (name hash, height, resource data, blind, etc.).
	Items [][]byte
}

const (
	covenantTypeNone     = 0
	covenantTypeClaim    = 1
	covenantTypeOpen     = 2
	covenantTypeBid      = 3
	covenantTypeReveal   = 4
	covenantTypeRedeem   = 5
	covenantTypeRegister = 6
	covenantTypeUpdate   = 7
	covenantTypeRenew    = 8
	covenantTypeTransfer = 9
	covenantTypeFinalize = 10
	covenantTypeRevoke   = 11 // Matches the JS covenant.TypeRevoke value.
)

// IsUnspendable reports whether the covenant burns the output.
func (c Covenant) IsUnspendable() bool {
	return c.Type == covenantTypeRevoke
}

// IsNone reports whether the covenant is a plain transfer.
func (c Covenant) IsNone() bool {
	return c.Type == covenantTypeNone
}

// IsClaim reports whether the covenant is a claim.
func (c Covenant) IsClaim() bool {
	return c.Type == covenantTypeClaim
}

// IsOpen reports whether the covenant is an open.
func (c Covenant) IsOpen() bool {
	return c.Type == covenantTypeOpen
}

// IsBid reports whether the covenant is a bid.
func (c Covenant) IsBid() bool {
	return c.Type == covenantTypeBid
}

// IsReveal reports whether the covenant is a reveal.
func (c Covenant) IsReveal() bool {
	return c.Type == covenantTypeReveal
}

// IsRedeem reports whether the covenant is a redeem.
func (c Covenant) IsRedeem() bool {
	return c.Type == covenantTypeRedeem
}

// IsRegister reports whether the covenant is a register.
func (c Covenant) IsRegister() bool {
	return c.Type == covenantTypeRegister
}

// IsUpdate reports whether the covenant is an update.
func (c Covenant) IsUpdate() bool {
	return c.Type == covenantTypeUpdate
}

// IsRenew reports whether the covenant is a renewal.
func (c Covenant) IsRenew() bool {
	return c.Type == covenantTypeRenew
}

// IsTransfer reports whether the covenant is a transfer.
func (c Covenant) IsTransfer() bool {
	return c.Type == covenantTypeTransfer
}

// IsFinalize reports whether the covenant is a finalise.
func (c Covenant) IsFinalize() bool {
	return c.Type == covenantTypeFinalize
}

// IsRevoke reports whether the covenant is a revocation.
func (c Covenant) IsRevoke() bool {
	return c.Type == covenantTypeRevoke
}

// IsName reports whether the covenant carries name-system semantics.
func (c Covenant) IsName() bool {
	return c.Type >= covenantTypeClaim && c.Type <= covenantTypeRevoke
}

// IsKnown reports whether the covenant type is recognised.
func (c Covenant) IsKnown() bool {
	return c.Type <= covenantTypeRevoke
}

// IsUnknown reports whether the covenant type is not recognised.
func (c Covenant) IsUnknown() bool {
	return c.Type > covenantTypeRevoke
}

// IsDustworthy reports whether the covenant is subject to dust policy.
func (c Covenant) IsDustworthy() bool {
	switch c.Type {
	case covenantTypeNone, covenantTypeBid:
		return true
	default:
		return c.Type > covenantTypeRevoke
	}
}

// IsNonspendable reports whether the covenant should be treated as locked.
func (c Covenant) IsNonspendable() bool {
	switch c.Type {
	case covenantTypeNone, covenantTypeOpen, covenantTypeRedeem:
		return false
	default:
		return true
	}
}

// IsLinked reports whether the covenant is part of the linked covenant set.
func (c Covenant) IsLinked() bool {
	return c.Type >= covenantTypeReveal && c.Type <= covenantTypeRevoke
}

// Transaction is a complete name-chain transaction containing inputs that
// spend previous outputs and outputs that create new coins or name operations.
//
//	tx := primitives.Transaction{Version: 0, LockTime: 0}
type Transaction struct {
	// Version is the transaction format version.
	Version uint32

	// Inputs is the list of inputs spending previous outputs.
	Inputs []Input

	// Outputs is the list of new outputs created by this transaction.
	Outputs []Output

	// LockTime is the earliest time or block height at which this
	// transaction may be included in a block.
	LockTime uint32
}

// BlockHeader contains the fields present in every LNS block header.
// The header is split into a preheader (essential for SPV) and a subheader
// (miner-mutable data), following the HSD dual-hash PoW design.
//
//	hdr := primitives.BlockHeader{Version: 1, Time: uint64(time.Now().Unix())}
type BlockHeader struct {
	// Version determines which consensus rules apply to this block.
	Version uint32

	// PrevBlock is the hash of the previous block in the chain.
	PrevBlock Hash

	// MerkleRoot is the root hash of the transaction merkle tree.
	MerkleRoot Hash

	// WitnessRoot is the root hash of the witness merkle tree.
	WitnessRoot Hash

	// TreeRoot is the root hash of the name/UTXO Urkel tree.
	TreeRoot Hash

	// ReservedRoot is reserved for future use (must be zero for now).
	ReservedRoot Hash

	// Time is the Unix epoch timestamp (seconds) when the block was created.
	Time uint64

	// Bits encodes the target difficulty in compact form.
	Bits uint32

	// Nonce is iterated by the miner to find a valid PoW solution.
	Nonce uint32

	// ExtraNonce provides additional entropy for miners (24 bytes in HSD).
	ExtraNonce []byte

	// Mask is the PoW mask used for pool block-withholding protection (32 bytes).
	Mask []byte
}

// NameState tracks the current state of a registered name on the chain.
// Names progress through auction phases (OPENING -> BIDDING -> REVEAL -> CLOSED)
// and can be transferred, renewed, or revoked.
//
//	ns := primitives.NameState{Name: []byte("example"), Height: 1000}
//	if ns.IsExpired() { /* name can be re-auctioned */ }
type NameState struct {
	// Name is the raw name bytes (1-63 characters).
	Name []byte

	// NameHash is the blake2b hash of the name used for lookups.
	NameHash Hash

	// Height is the block height at which the current state was established.
	Height uint32

	// Renewal is the block height of the last renewal.
	Renewal uint32

	// Owner references the outpoint that currently owns this name.
	Owner Outpoint

	// Value is the locked value (winning bid amount).
	Value uint64

	// Highest is the highest bid seen during the auction.
	Highest uint64

	// Data is the DNS resource data attached to this name.
	Data []byte

	// Transfer is the block height at which a transfer was initiated (0 if none).
	Transfer uint32

	// Revoked is the block height at which the name was revoked (0 if not revoked).
	Revoked uint32

	// Claimed is the number of times this name has been claimed.
	Claimed uint32

	// Renewals is the total number of renewals for this name.
	Renewals uint32

	// Registered indicates whether the name has completed the REGISTER covenant.
	Registered bool

	// Expired indicates whether the name has passed its renewal window.
	Expired bool

	// Weak indicates a weak claim (imported from ICANN without full proof).
	Weak bool

	// delta caches the in-memory state delta used by undo/redo flows.
	// It is not part of the persistent name-state payload.
	delta *NameDelta
}

// IsNull reports whether the outpoint is the null outpoint used by coinbase
// inputs.
func (o Outpoint) IsNull() bool {
	return o.TxHash == (Hash{}) && o.Index == ^uint32(0)
}

// IsCoinbase reports whether the input spends the null outpoint.
func (i Input) IsCoinbase() bool {
	return i.Prevout.IsNull()
}

// IsFinal reports whether the input sequence disables locktime semantics.
func (i Input) IsFinal() bool {
	return i.Sequence == ^uint32(0)
}

// IsUnspendable reports whether the output cannot be spent.
func (o Output) IsUnspendable() bool {
	return o.Address.IsUnspendable() || o.Covenant.IsUnspendable()
}
