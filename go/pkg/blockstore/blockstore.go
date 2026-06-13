// SPDX-License-Identifier: EUPL-1.2

// Package blockstore is the local block-persistence layer for the Lethean Name
// System. It keeps the raw bytes of confirmed blocks on go-store and indexes
// them so the chain can be read back by hash or by height, the tip can be
// tracked, and the per-block undo (rollback) data needed to disconnect a block
// during a reorganisation can be retrieved.
//
// # Why it exists alongside pkg/chain
//
// pkg/chain is an RPC client: it fetches blocks from a remote sidechain daemon
// over JSON-RPC (chain.Client.GetBlock / GetBlockByHeight) and discovers main-
// chain aliases. It holds nothing locally — every read is a network round-trip
// and there is no durable copy of the chain on disk. blockstore is the other
// half: it owns the on-disk representation. A sync routine fetches a block with
// pkg/chain, validates it, and hands the raw bytes to blockstore.PutBlock; from
// then on the block is served locally by hash or height with no daemon.
//
// blockstore deliberately treats blocks as opaque bytes. It does not parse a
// block's transactions or header fields — that is the job of pkg/primitives and
// the chain-validation layer. blockstore's contract is durability and indexing:
// given the bytes plus the hash and height the caller has already computed, it
// stores them and gives them back, and it maintains the height->hash index and
// the tip pointer so the chain can be walked without re-deriving anything.
//
// # Storage model
//
// Everything is persisted through a go-store *store.ScopedStore. go-store is a
// string key/value store grouped into named groups, so binary payloads (blocks,
// headers, undo data) are hex-encoded on write and decoded on read, and the
// indices are stored as plain strings. The groups used are:
//
//	block   key=hex(hash)        value=hex(blockRecord{height,bytes})
//	height  key=hex(beUint32)    value=hex(hash)
//	header  key=hex(hash)        value=hex(headerBytes)
//	undo    key=hex(hash)        value=hex(undoBytes)
//	meta    key="tip.hash"       value=hex(hash)
//	        key="tip.height"     value=decimal height
//
// Height keys are the big-endian encoding of the uint32 height, hex-encoded, so
// the natural key order is also height order — a paged scan of the height group
// walks the chain from genesis upward.
//
// # Usage
//
//	scoped := store.NewScoped(storeInstance, "lns-blocks")
//	bs, err := blockstore.New(scoped)
//	if err != nil { return err }
//
//	if err := bs.PutBlock(hash, height, rawBlockBytes); err != nil { return err }
//	raw, err := bs.GetBlock(hash)            // by hash
//	raw, err = bs.GetBlockByHeight(height)    // by height (via the index)
//	if has, _ := bs.HasBlock(hash); has { /* already stored */ }
//	tipHash, tipHeight, _ := bs.Tip()
package blockstore

import (
	"encoding/binary"

	core "dappco.re/go"
	"dappco.re/go/store"

	"dappco.re/go/lns/pkg/primitives"
)

// Sentinel causes returned (wrapped by core.E) by the read and write helpers.
// Match them with errors.Is via core.E's Unwrap chain:
//
//	_, err := bs.GetBlock(hash)
//	if errors.Is(err, blockstore.ErrNotFound) { /* block is not stored */ }
var (
	// ErrNilStore is returned by New when the supplied scoped store is nil.
	ErrNilStore = core.NewError("blockstore: scoped store is nil")

	// ErrNotFound is returned by the read helpers when no record exists for the
	// requested hash or height.
	ErrNotFound = core.NewError("blockstore: block not found")

	// ErrEmptyBlock is returned by PutBlock when the raw block byte slice is
	// empty: an empty block is never a valid thing to persist and storing it
	// would make HasBlock report a block that cannot be served back meaningfully.
	ErrEmptyBlock = core.NewError("blockstore: block bytes are empty")

	// ErrCorrupt is returned when a stored value cannot be decoded: a hex value
	// that is not valid hex, or a block record whose framing is inconsistent.
	// It signals on-disk corruption (or a value written by something that is not
	// this package), not a missing key.
	ErrCorrupt = core.NewError("blockstore: stored value is corrupt")
)

// The go-store groups blockstore writes into. They are namespaced again by the
// ScopedStore's own namespace, so two blockstores on different namespaces of the
// same underlying store do not collide.
const (
	groupBlock  = "block"
	groupHeight = "height"
	groupHeader = "header"
	groupUndo   = "undo"
	groupMeta   = "meta"

	metaTipHash   = "tip.hash"
	metaTipHeight = "tip.height"
)

// maxBlockBytes bounds the size of a single decoded block record. The Handshake
// consensus block size limit this repository ports is well under a megabyte of
// serialised data; 8 MiB is a generous ceiling that still refuses a hostile or
// corrupt length prefix before it can drive a large allocation in decodeRecord.
const maxBlockBytes = 8 << 20

// Store persists confirmed blocks on a go-store scoped store and indexes them by
// hash and by height. The zero value is not usable; construct one with New.
//
//	bs, err := blockstore.New(store.NewScoped(storeInstance, "lns-blocks"))
type Store struct {
	scoped *store.ScopedStore
}

// New constructs a block store over an existing go-store scoped store. The
// scoped store carries the namespace under which every block, index, and meta
// key is written, so callers choose isolation by choosing the namespace.
//
//	scoped := store.NewScoped(storeInstance, "lns-blocks")
//	bs, err := blockstore.New(scoped)
func New(scoped *store.ScopedStore) (*Store, error) {
	if scoped == nil {
		return nil, core.E("blockstore.New", "construct block store", ErrNilStore)
	}
	return &Store{scoped: scoped}, nil
}

// blockRecord is the on-disk framing for a stored block: its height followed by
// the raw block bytes. The height travels with the block so GetBlock can report
// the height of a block fetched by hash without consulting the index, and so the
// index can be rebuilt from the block group alone if it is ever lost.
//
// Wire layout:
//
//	height(4 LE) ‖ length(4 LE) ‖ rawBlockBytes
//
// The explicit length prefix makes the framing self-checking: decodeRecord
// rejects a record whose declared length does not match the bytes that follow,
// and refuses a length beyond maxBlockBytes before allocating.
type blockRecord struct {
	height uint32
	raw    []byte
}

// encodeRecord serialises a block record to its wire form. It never fails for a
// well-formed in-memory record; the length prefix is the byte length of raw.
func encodeRecord(rec blockRecord) []byte {
	buf := make([]byte, 0, 8+len(rec.raw))
	var head [8]byte
	binary.LittleEndian.PutUint32(head[0:4], rec.height)
	binary.LittleEndian.PutUint32(head[4:8], uint32(len(rec.raw)))
	buf = append(buf, head[:]...)
	buf = append(buf, rec.raw...)
	return buf
}

// decodeRecord parses a block record produced by encodeRecord. It validates the
// framing before trusting any length: an 8-byte header must be present, the
// declared length must not exceed maxBlockBytes (so a corrupt prefix cannot
// drive an unbounded allocation), and the declared length must equal the number
// of trailing bytes (so truncation or trailing junk is caught rather than
// silently returning a short or over-long block).
func decodeRecord(data []byte) (blockRecord, error) {
	if len(data) < 8 {
		return blockRecord{}, core.E("blockstore.decodeRecord", "short record header", ErrCorrupt)
	}

	height := binary.LittleEndian.Uint32(data[0:4])
	length := binary.LittleEndian.Uint32(data[4:8])
	body := data[8:]

	if uint64(length) > maxBlockBytes {
		return blockRecord{}, core.E("blockstore.decodeRecord",
			core.Sprintf("declared length %d exceeds maximum %d", length, maxBlockBytes), ErrCorrupt)
	}
	if int(length) != len(body) {
		return blockRecord{}, core.E("blockstore.decodeRecord",
			core.Sprintf("declared length %d does not match %d trailing bytes", length, len(body)), ErrCorrupt)
	}

	raw := make([]byte, len(body))
	copy(raw, body)
	return blockRecord{height: height, raw: raw}, nil
}

// hashKey is the group key for a block, header, or undo entry: the lowercase hex
// of the 32-byte hash. Using hex keeps the key inside the printable-ASCII set
// the underlying string store expects.
func hashKey(hash primitives.Hash) string {
	return core.HexEncode(hash[:])
}

// heightKey is the group key for a height->hash index entry: the hex of the
// big-endian uint32 height. Big-endian so lexical key order matches numeric
// height order, which lets a paged scan of the height group walk the chain in
// order.
func heightKey(height uint32) string {
	var be [4]byte
	binary.BigEndian.PutUint32(be[:], height)
	return core.HexEncode(be[:])
}

// putBytes hex-encodes a binary value and writes it to a group. Centralising the
// encoding keeps every binary write in the package consistent and reversible by
// getBytes.
func (s *Store) putBytes(group, key string, value []byte) error {
	if err := s.scoped.SetIn(group, key, core.HexEncode(value)); err != nil {
		return core.E("blockstore.putBytes", core.Sprintf("write %s/%s", group, key), err)
	}
	return nil
}

// getBytes reads a hex-encoded binary value from a group. It returns ErrNotFound
// when the key is absent and ErrCorrupt when the stored string is not valid hex.
func (s *Store) getBytes(group, key string) ([]byte, error) {
	exists, err := s.scoped.ExistsIn(group, key)
	if err != nil {
		return nil, core.E("blockstore.getBytes", core.Sprintf("probe %s/%s", group, key), err)
	}
	if !exists {
		return nil, core.E("blockstore.getBytes", core.Sprintf("read %s/%s", group, key), ErrNotFound)
	}

	encoded, err := s.scoped.GetFrom(group, key)
	if err != nil {
		return nil, core.E("blockstore.getBytes", core.Sprintf("read %s/%s", group, key), err)
	}

	decoded := core.HexDecode(encoded)
	if !decoded.OK {
		return nil, core.E("blockstore.getBytes", core.Sprintf("decode %s/%s", group, key), ErrCorrupt)
	}
	raw, _ := decoded.Value.([]byte)
	return raw, nil
}
