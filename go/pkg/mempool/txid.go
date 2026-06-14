// SPDX-License-Identifier: EUPL-1.2

// Package mempool implements the local pool of unconfirmed transactions for the
// Lethean Name System. It mirrors the behaviour of the HSD (Handshake) mempool
// (lib/mempool/mempool.js) that this repository ports, restricted to the
// well-defined core: structural validation, witness/script verification,
// double-spend detection against the pool and the chain, fee-rate ordered block
// templates, and size/expiry eviction.
//
// # What it does
//
// A Pool accepts a primitives.Transaction, validates it (structure + script +
// double-spend), indexes it by transaction id, tracks every outpoint it claims
// so conflicts are O(1), and exposes the pool ordered by fee rate so a miner can
// build a block template. When the pool exceeds its byte or count budget the
// lowest-fee-rate transactions are evicted; transactions older than the expiry
// window are dropped first.
//
//	pool := mempool.New(mempool.DefaultConfig())
//	_, err := pool.Add(tx, chainView)   // chainView is the confirmed UTXO set
//	if pool.Has(txid) { ... }
//	template := pool.GetForBlock(0)     // fee-rate ordered snapshot
//
// # Scope
//
// The package implements the single-transaction core of the HSD mempool. The
// following HSD features are deliberately NOT implemented here, because they are
// large and the LNS RFC (code/core/go/lns/RFC.md §2) specifies only that
// pkg/mempool is "the mempool" without fixing their semantics. They are left to
// a later, separately-specified pass rather than stubbed:
//
//   - Replace-by-fee (RBF): a conflicting higher-fee transaction is rejected,
//     not swapped in. Add reports ErrDoubleSpend.
//   - Orphan handling: a transaction whose inputs are neither in the pool nor in
//     the supplied chain view is rejected (ErrMissingInputs), not parked in an
//     orphan set awaiting its parents.
//   - Ancestor/descendant fee packages: fee rate is computed per transaction.
//     Child-pays-for-parent packaging and the ancestor-count limit are not
//     modelled.
//   - Claim/airdrop name de-duplication (the claimNames index) and minimum
//     relay-fee / free-transaction priority policy.
//
// Chained spends ARE supported: a transaction may spend an output created by
// another transaction already in the pool (PoolView resolves pool outputs), it
// is simply not grouped into a fee package for eviction.
package mempool

import (
	"encoding/binary"

	core "dappco.re/go"
	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/primitives"
)

// TxID is the 32-byte identifier of a transaction: the BLAKE2b-256 digest of its
// non-witness serialisation. Witness data is excluded so that re-signing or
// witness malleation cannot change the id, matching the HSD segwit-style txid
// (the "hash" in lib/primitives/tx.js, distinct from the witness hash).
//
//	id := mempool.ComputeTxID(tx)
//	if id == (mempool.TxID{}) { /* zero id, only for an empty/invalid tx */ }
type TxID = primitives.Hash

// ComputeTxID returns the transaction id: BLAKE2b-256 over the canonical
// non-witness serialisation produced by SerializeNoWitness. Two transactions
// with identical inputs (outpoints + sequences), outputs, version, and locktime
// share an id regardless of their witnesses.
//
//	id, err := mempool.ComputeTxID(tx)
func ComputeTxID(tx primitives.Transaction) (TxID, error) {
	enc, err := SerializeNoWitness(tx)
	if err != nil {
		return TxID{}, core.E("mempool.ComputeTxID", "serialise transaction", err)
	}

	digest := blake2b.Sum256(enc)
	return TxID(digest), nil
}

// SerializeNoWitness encodes a transaction without its input witnesses, in the
// HSD wire order used for the transaction id:
//
//	version(4 LE) ‖ varint(inputCount) ‖ inputs ‖ varint(outputCount) ‖
//	outputs ‖ locktime(4 LE)
//
// Each input is prevout(36) ‖ sequence(4 LE). Each output is reused from the
// shared output encoding: value(8 LE) ‖ version(1) ‖ hashLen(1) ‖ hash ‖
// covenant. The witness stack is omitted entirely.
//
//	raw, err := mempool.SerializeNoWitness(tx)
func SerializeNoWitness(tx primitives.Transaction) ([]byte, error) {
	buf := make([]byte, 0, estimateSerializedSize(tx))

	buf = appendU32(buf, tx.Version)

	buf = appendVarint(buf, uint64(len(tx.Inputs)))
	for i := range tx.Inputs {
		in := tx.Inputs[i]

		prevout, err := in.Prevout.MarshalBinary()
		if err != nil {
			return nil, core.E("mempool.SerializeNoWitness", "marshal prevout", err)
		}
		buf = append(buf, prevout...)
		buf = appendU32(buf, in.Sequence)
	}

	buf = appendVarint(buf, uint64(len(tx.Outputs)))
	for i := range tx.Outputs {
		out, err := encodeOutput(tx.Outputs[i])
		if err != nil {
			return nil, core.E("mempool.SerializeNoWitness", "marshal output", err)
		}
		buf = append(buf, out...)
	}

	buf = appendU32(buf, tx.LockTime)

	return buf, nil
}

// SerializeWithWitness encodes a transaction in the full HSD witness wire form:
// the base (non-witness) serialisation produced by SerializeNoWitness followed
// by each input's witness stack, in input order:
//
//	version(4 LE) ‖ varint(inCount) ‖ inputs ‖ varint(outCount) ‖ outputs ‖
//	locktime(4 LE) ‖ for each input: varint(itemCount) ‖ (varint(len) ‖ item)…
//
// This is the byte string a fully-signed transaction is hex-encoded from for the
// daemon's sendrawtransaction RPC. The base segment is byte-identical to the
// txid preimage, so the broadcast bytes and the txid stay consistent. It mirrors
// the block-body transaction layout (pkg/mining marshalTxWitness) so a
// transaction serialises identically whether it is broadcast or mined.
//
//	raw, err := mempool.SerializeWithWitness(tx)
//	txhex := core.HexEncode(raw)
func SerializeWithWitness(tx primitives.Transaction) ([]byte, error) {
	base, err := SerializeNoWitness(tx)
	if err != nil {
		return nil, core.E("mempool.SerializeWithWitness", "serialise base", err)
	}

	buf := make([]byte, 0, len(base)+witnessDataSize(tx))
	buf = append(buf, base...)
	for i := range tx.Inputs {
		witness := tx.Inputs[i].Witness
		buf = appendVarint(buf, uint64(len(witness)))
		for _, item := range witness {
			buf = appendVarint(buf, uint64(len(item)))
			buf = append(buf, item...)
		}
	}
	return buf, nil
}

// witnessDataSize returns the serialised size of all input witness stacks, used
// to pre-size the SerializeWithWitness buffer.
func witnessDataSize(tx primitives.Transaction) int {
	size := 0
	for i := range tx.Inputs {
		size += witnessSize(tx.Inputs[i].Witness)
	}
	return size
}

// TxSize returns the serialised byte size of a transaction including witnesses.
// It is the size the pool charges against its byte budget and divides fees by to
// derive a fee rate. The witness contributes its full length (this pool does not
// apply a witness-scale discount; the LNS RFC does not specify one).
//
//	size, err := mempool.TxSize(tx)
func TxSize(tx primitives.Transaction) (int, error) {
	base, err := SerializeNoWitness(tx)
	if err != nil {
		return 0, core.E("mempool.TxSize", "serialise base", err)
	}

	size := len(base)
	for i := range tx.Inputs {
		size += witnessSize(tx.Inputs[i].Witness)
	}

	return size, nil
}

// witnessSize returns the serialised size of one input's witness stack:
// varint(itemCount) followed by each item as varint(len) ‖ bytes.
func witnessSize(witness [][]byte) int {
	size := varintSize(uint64(len(witness)))
	for _, item := range witness {
		size += varintSize(uint64(len(item))) + len(item)
	}
	return size
}

// encodeOutput serialises a transaction output in the shared HSD wire form:
//
//	value(8 LE) ‖ address.version(1) ‖ address.hashLen(1) ‖ address.hash ‖ covenant
//
// It mirrors script.encodeOutput (which is unexported) and reuses
// primitives.Covenant.MarshalBinary for the covenant tail so the byte layout is
// identical to the one the signature hash commits to.
func encodeOutput(out primitives.Output) ([]byte, error) {
	cov, err := out.Covenant.MarshalBinary()
	if err != nil {
		return nil, core.E("mempool.encodeOutput", "marshal covenant", err)
	}
	if len(out.Address.Hash) > 0xff {
		return nil, core.E("mempool.encodeOutput", "address hash too long", nil)
	}

	buf := make([]byte, 0, 8+2+len(out.Address.Hash)+len(cov))
	buf = appendU64(buf, out.Value)
	buf = append(buf, out.Address.Version, byte(len(out.Address.Hash)))
	buf = append(buf, out.Address.Hash...)
	buf = append(buf, cov...)
	return buf, nil
}

// estimateSerializedSize returns an upper bound on the non-witness size used to
// pre-size the serialisation buffer. It is a hint, not a guarantee.
func estimateSerializedSize(tx primitives.Transaction) int {
	size := 4 + 4 + 9 + 9 // version + locktime + two varint count prefixes (max)
	size += len(tx.Inputs) * (36 + 4)
	for i := range tx.Outputs {
		size += 8 + 2 + len(tx.Outputs[i].Address.Hash) + tx.Outputs[i].Covenant.GetVarSize()
	}
	return size
}

// appendU32 appends a little-endian uint32.
func appendU32(dst []byte, v uint32) []byte {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], v)
	return append(dst, tmp[:]...)
}

// appendU64 appends a little-endian uint64.
func appendU64(dst []byte, v uint64) []byte {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], v)
	return append(dst, tmp[:]...)
}

// appendVarint appends a Bitcoin/HSD compact-size integer.
func appendVarint(dst []byte, n uint64) []byte {
	switch {
	case n < 0xfd:
		return append(dst, byte(n))
	case n <= 0xffff:
		return append(dst, 0xfd, byte(n), byte(n>>8))
	case n <= 0xffffffff:
		return append(dst, 0xfe, byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
	default:
		return append(dst,
			0xff,
			byte(n), byte(n>>8), byte(n>>16), byte(n>>24),
			byte(n>>32), byte(n>>40), byte(n>>48), byte(n>>56))
	}
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
