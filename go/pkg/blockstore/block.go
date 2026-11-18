// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	core "github.com/dappcore/go"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// PutBlock stores a confirmed block under its hash, records the height->hash
// index entry, and advances the tip if this block extends it. The caller has
// already computed hash and height (and validated the block); blockstore treats
// raw as opaque bytes. Storing the same hash twice overwrites the previous
// record with the new height and bytes, which is what a caller re-validating a
// block at a corrected height wants. An empty raw slice is rejected with
// ErrEmptyBlock.
//
//	if err := bs.PutBlock(hash, height, rawBlockBytes); err != nil { return err }
func (s *Store) PutBlock(hash primitives.Hash, height uint32, raw []byte) error {
	if s == nil || s.scoped == nil {
		return core.E("blockstore.PutBlock", "store block", ErrNilStore)
	}
	if len(raw) == 0 {
		return core.E("blockstore.PutBlock", "store block", ErrEmptyBlock)
	}
	if len(raw) > maxBlockBytes {
		return core.E("blockstore.PutBlock",
			core.Sprintf("block of %d bytes exceeds maximum %d", len(raw), maxBlockBytes), ErrCorrupt)
	}

	record := encodeRecord(blockRecord{height: height, raw: raw})
	if err := s.putBytes(groupBlock, hashKey(hash), record); err != nil {
		return core.E("blockstore.PutBlock", "write block record", err)
	}

	if err := s.putBytes(groupHeight, heightKey(height), hash[:]); err != nil {
		return core.E("blockstore.PutBlock", "write height index", err)
	}

	if err := s.advanceTip(hash, height); err != nil {
		return core.E("blockstore.PutBlock", "advance tip", err)
	}

	return nil
}

// GetBlock returns the raw bytes of the block stored under hash. It returns
// ErrNotFound when no block is stored for that hash and ErrCorrupt when the
// stored record cannot be decoded.
//
//	raw, err := bs.GetBlock(hash)
func (s *Store) GetBlock(hash primitives.Hash) ([]byte, error) {
	if s == nil || s.scoped == nil {
		return nil, core.E("blockstore.GetBlock", "read block", ErrNilStore)
	}

	encoded, err := s.getBytes(groupBlock, hashKey(hash))
	if err != nil {
		return nil, core.E("blockstore.GetBlock", "read block record", err)
	}

	record, err := decodeRecord(encoded)
	if err != nil {
		return nil, core.E("blockstore.GetBlock", "decode block record", err)
	}
	return record.raw, nil
}

// GetBlockWithHeight returns the raw bytes of the block stored under hash
// together with the height recorded alongside it. It is the same read as
// GetBlock but also surfaces the height the block was stored at, without a
// second lookup against the index.
//
//	raw, height, err := bs.GetBlockWithHeight(hash)
func (s *Store) GetBlockWithHeight(hash primitives.Hash) ([]byte, uint32, error) {
	if s == nil || s.scoped == nil {
		return nil, 0, core.E("blockstore.GetBlockWithHeight", "read block", ErrNilStore)
	}

	encoded, err := s.getBytes(groupBlock, hashKey(hash))
	if err != nil {
		return nil, 0, core.E("blockstore.GetBlockWithHeight", "read block record", err)
	}

	record, err := decodeRecord(encoded)
	if err != nil {
		return nil, 0, core.E("blockstore.GetBlockWithHeight", "decode block record", err)
	}
	return record.raw, record.height, nil
}

// HashAtHeight returns the hash recorded at a height in the index. It returns
// ErrNotFound when no block is indexed at that height and ErrCorrupt when the
// stored hash is not a 32-byte value.
//
//	hash, err := bs.HashAtHeight(height)
func (s *Store) HashAtHeight(height uint32) (primitives.Hash, error) {
	if s == nil || s.scoped == nil {
		return primitives.Hash{}, core.E("blockstore.HashAtHeight", "read height index", ErrNilStore)
	}

	raw, err := s.getBytes(groupHeight, heightKey(height))
	if err != nil {
		return primitives.Hash{}, core.E("blockstore.HashAtHeight", "read height index", err)
	}
	if len(raw) != len(primitives.Hash{}) {
		return primitives.Hash{}, core.E("blockstore.HashAtHeight",
			core.Sprintf("indexed hash is %d bytes, want 32", len(raw)), ErrCorrupt)
	}

	var hash primitives.Hash
	copy(hash[:], raw)
	return hash, nil
}

// GetBlockByHeight returns the raw bytes of the block indexed at height. It
// resolves the height to a hash through the index and then reads the block,
// returning ErrNotFound when either the index entry or the block is absent.
//
//	raw, err := bs.GetBlockByHeight(height)
func (s *Store) GetBlockByHeight(height uint32) ([]byte, error) {
	if s == nil || s.scoped == nil {
		return nil, core.E("blockstore.GetBlockByHeight", "read block", ErrNilStore)
	}

	hash, err := s.HashAtHeight(height)
	if err != nil {
		return nil, core.E("blockstore.GetBlockByHeight", "resolve height", err)
	}

	raw, err := s.GetBlock(hash)
	if err != nil {
		return nil, core.E("blockstore.GetBlockByHeight", "read indexed block", err)
	}
	return raw, nil
}

// HasBlock reports whether a block is stored under hash. It distinguishes a
// missing block (false, ok result) from a store failure (false, error).
//
//	if has, _ := bs.HasBlock(hash); has { /* skip refetch */ }
func (s *Store) HasBlock(hash primitives.Hash) (bool, error) {
	if s == nil || s.scoped == nil {
		return false, core.E("blockstore.HasBlock", "probe block", ErrNilStore)
	}

	exists, existsResult := s.scoped.ExistsIn(groupBlock, hashKey(hash))
	if err := existsResult.Err(); err != nil {
		return false, core.E("blockstore.HasBlock", "probe block", err)
	}
	return exists, nil
}

// HasHeight reports whether a block is indexed at height.
//
//	if has, _ := bs.HasHeight(height); has { /* height already filled */ }
func (s *Store) HasHeight(height uint32) (bool, error) {
	if s == nil || s.scoped == nil {
		return false, core.E("blockstore.HasHeight", "probe height", ErrNilStore)
	}

	exists, existsResult := s.scoped.ExistsIn(groupHeight, heightKey(height))
	if err := existsResult.Err(); err != nil {
		return false, core.E("blockstore.HasHeight", "probe height", err)
	}
	return exists, nil
}

// DeleteBlock removes a block, its height-index entry, its header, and its undo
// data. It is the disconnect half of a reorganisation: the caller has already
// rewound the tip with SetTip (or will), and DeleteBlock drops the now-orphaned
// data. Deleting a hash that is not stored is not an error — the end state, the
// block being absent, is the same either way. The height-index entry is only
// removed if it still points at this hash, so deleting a block that a competing
// block has already replaced at that height leaves the winner's index intact.
//
//	if err := bs.DeleteBlock(hash); err != nil { return err }
func (s *Store) DeleteBlock(hash primitives.Hash) error {
	if s == nil || s.scoped == nil {
		return core.E("blockstore.DeleteBlock", "delete block", ErrNilStore)
	}

	// Recover the height from the record (if present) so the matching index
	// entry can be removed. A missing record means nothing to delete.
	encoded, err := s.getBytes(groupBlock, hashKey(hash))
	if err != nil {
		if core.Is(err, ErrNotFound) {
			return nil
		}
		return core.E("blockstore.DeleteBlock", "read block record", err)
	}

	record, decErr := decodeRecord(encoded)
	if decErr != nil {
		// The record is corrupt; still drop the block, header, and undo so the
		// bad value does not linger, but the height index cannot be matched.
		if delErr := s.deleteBlockBytes(hash); delErr != nil {
			return core.E("blockstore.DeleteBlock", "delete corrupt block", delErr)
		}
		return nil
	}

	if delErr := s.deleteHeightIfMatches(record.height, hash); delErr != nil {
		return core.E("blockstore.DeleteBlock", "delete height index", delErr)
	}
	if delErr := s.deleteBlockBytes(hash); delErr != nil {
		return core.E("blockstore.DeleteBlock", "delete block data", delErr)
	}
	return nil
}

// deleteBlockBytes removes the block record, header, and undo entries for a hash.
// Each group is independent, so it deletes from all three regardless of whether
// every group held an entry.
func (s *Store) deleteBlockBytes(hash primitives.Hash) error {
	key := hashKey(hash)
	if err := s.scoped.Delete(groupBlock, key).Err(); err != nil {
		return core.E("blockstore.deleteBlockBytes", "delete block", err)
	}
	if err := s.scoped.Delete(groupHeader, key).Err(); err != nil {
		return core.E("blockstore.deleteBlockBytes", "delete header", err)
	}
	if err := s.scoped.Delete(groupUndo, key).Err(); err != nil {
		return core.E("blockstore.deleteBlockBytes", "delete undo", err)
	}
	return nil
}

// deleteHeightIfMatches removes the height-index entry only when it still maps to
// the given hash, so a block superseded at its height by a competing block does
// not delete the competitor's index entry when it is later pruned.
func (s *Store) deleteHeightIfMatches(height uint32, hash primitives.Hash) error {
	indexed, err := s.HashAtHeight(height)
	if err != nil {
		if core.Is(err, ErrNotFound) {
			return nil
		}
		return core.E("blockstore.deleteHeightIfMatches", "read height index", err)
	}
	if indexed != hash {
		return nil
	}
	if delErr := s.scoped.Delete(groupHeight, heightKey(height)).Err(); delErr != nil {
		return core.E("blockstore.deleteHeightIfMatches", "delete height index", delErr)
	}
	return nil
}
