// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	core "dappco.re/go"

	"dappco.re/go/lns/pkg/primitives"
)

// PutUndo stores the undo (rollback) data for a block hash. Undo data is the set
// of coins a block spent, captured so the block can be disconnected during a
// reorganisation: to rewind a block you restore the coins it consumed, and those
// coins are exactly what the undo blob records. blockstore keeps the undo blob
// opaque; the chain layer decides its format. An empty undo blob is rejected
// with ErrEmptyBlock — a block with no spends should simply have no undo entry
// rather than an empty one, so HasUndo stays meaningful.
//
//	if err := bs.PutUndo(hash, undoBytes); err != nil { return err }
func (s *Store) PutUndo(hash primitives.Hash, undo []byte) error {
	if s == nil || s.scoped == nil {
		return core.E("blockstore.PutUndo", "store undo", ErrNilStore)
	}
	if len(undo) == 0 {
		return core.E("blockstore.PutUndo", "store undo", ErrEmptyBlock)
	}
	if len(undo) > maxBlockBytes {
		return core.E("blockstore.PutUndo",
			core.Sprintf("undo of %d bytes exceeds maximum %d", len(undo), maxBlockBytes), ErrCorrupt)
	}

	if err := s.putBytes(groupUndo, hashKey(hash), undo); err != nil {
		return core.E("blockstore.PutUndo", "write undo", err)
	}
	return nil
}

// GetUndo returns the undo data stored under hash. It returns ErrNotFound when no
// undo data is stored for that hash (a block with no spends, or one stored
// without undo) and ErrCorrupt when the stored value is not valid hex.
//
//	undo, err := bs.GetUndo(hash)
func (s *Store) GetUndo(hash primitives.Hash) ([]byte, error) {
	if s == nil || s.scoped == nil {
		return nil, core.E("blockstore.GetUndo", "read undo", ErrNilStore)
	}

	undo, err := s.getBytes(groupUndo, hashKey(hash))
	if err != nil {
		return nil, core.E("blockstore.GetUndo", "read undo", err)
	}
	return undo, nil
}

// HasUndo reports whether undo data is stored under hash.
//
//	if has, _ := bs.HasUndo(hash); has { /* can disconnect this block */ }
func (s *Store) HasUndo(hash primitives.Hash) (bool, error) {
	if s == nil || s.scoped == nil {
		return false, core.E("blockstore.HasUndo", "probe undo", ErrNilStore)
	}

	exists, existsResult := s.scoped.ExistsIn(groupUndo, hashKey(hash))
	if err := existsResult.Err(); err != nil {
		return false, core.E("blockstore.HasUndo", "probe undo", err)
	}
	return exists, nil
}
