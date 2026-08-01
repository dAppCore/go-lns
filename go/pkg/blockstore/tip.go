// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	core "dappco.re/go"

	"dappco.re/go/lns/pkg/primitives"
)

// Tip returns the hash and height of the current chain tip: the most-work block
// the store has been told about. It returns ErrNotFound when no block has been
// stored yet (an empty store has no tip). The tip is advanced automatically by
// PutBlock when a stored block extends the chain, and can be rewound explicitly
// with SetTip during a reorganisation.
//
//	tipHash, tipHeight, err := bs.Tip()
func (s *Store) Tip() (primitives.Hash, uint32, error) {
	if s == nil || s.scoped == nil {
		return primitives.Hash{}, 0, core.E("blockstore.Tip", "read tip", ErrNilStore)
	}

	hashExists, existResult := s.scoped.ExistsIn(groupMeta, metaTipHash)
	if existErr := existResult.Err(); existErr != nil {
		return primitives.Hash{}, 0, core.E("blockstore.Tip", "probe tip hash", existErr)
	}
	if !hashExists {
		return primitives.Hash{}, 0, core.E("blockstore.Tip", "read tip", ErrNotFound)
	}

	rawHash, err := s.getBytes(groupMeta, metaTipHash)
	if err != nil {
		return primitives.Hash{}, 0, core.E("blockstore.Tip", "read tip hash", err)
	}
	if len(rawHash) != len(primitives.Hash{}) {
		return primitives.Hash{}, 0, core.E("blockstore.Tip",
			core.Sprintf("tip hash is %d bytes, want 32", len(rawHash)), ErrCorrupt)
	}

	heightStr, heightResult := s.scoped.GetFrom(groupMeta, metaTipHeight)
	if heightErr := heightResult.Err(); heightErr != nil {
		return primitives.Hash{}, 0, core.E("blockstore.Tip", "read tip height", heightErr)
	}
	parsed := core.Atoi(core.Trim(heightStr))
	if !parsed.OK {
		return primitives.Hash{}, 0, core.E("blockstore.Tip", "parse tip height", ErrCorrupt)
	}
	heightInt, _ := parsed.Value.(int)
	if heightInt < 0 {
		return primitives.Hash{}, 0, core.E("blockstore.Tip",
			core.Sprintf("tip height %d is negative", heightInt), ErrCorrupt)
	}

	var hash primitives.Hash
	copy(hash[:], rawHash)
	return hash, uint32(heightInt), nil
}

// HasTip reports whether a tip has been recorded yet. It distinguishes an empty
// store (false, ok result) from a store failure (false, error).
//
//	if has, _ := bs.HasTip(); !has { /* fresh store, start from genesis */ }
func (s *Store) HasTip() (bool, error) {
	if s == nil || s.scoped == nil {
		return false, core.E("blockstore.HasTip", "probe tip", ErrNilStore)
	}

	exists, existsResult := s.scoped.ExistsIn(groupMeta, metaTipHash)
	if err := existsResult.Err(); err != nil {
		return false, core.E("blockstore.HasTip", "probe tip", err)
	}
	return exists, nil
}

// SetTip records hash and height as the chain tip unconditionally. PutBlock keeps
// the tip current for blocks that extend the chain, so the common path needs no
// explicit call; SetTip exists for the cases PutBlock cannot infer — rewinding
// the tip to a lower height when disconnecting blocks during a reorganisation,
// or pinning the tip to a specific block of a competing branch. It does not
// verify that a block is stored under hash: tip is a pointer, and the caller
// owns the invariant that it points at a stored block.
//
//	if err := bs.SetTip(forkHash, forkHeight); err != nil { return err }
func (s *Store) SetTip(hash primitives.Hash, height uint32) error {
	if s == nil || s.scoped == nil {
		return core.E("blockstore.SetTip", "write tip", ErrNilStore)
	}

	if err := s.putBytes(groupMeta, metaTipHash, hash[:]); err != nil {
		return core.E("blockstore.SetTip", "write tip hash", err)
	}
	if err := s.scoped.SetIn(groupMeta, metaTipHeight, core.Itoa(int(height))).Err(); err != nil {
		return core.E("blockstore.SetTip", "write tip height", err)
	}
	return nil
}

// advanceTip sets the tip to hash/height when the store either has no tip yet or
// the candidate height is at least the current tip height. Equal heights advance
// the pointer too, so re-storing the tip block (or replacing it with a competing
// block at the same height) leaves the tip pointing at the most recently stored
// block at the chain head. A block stored below the tip — backfilling history —
// does not move the tip backward.
func (s *Store) advanceTip(hash primitives.Hash, height uint32) error {
	hasTip, err := s.HasTip()
	if err != nil {
		return core.E("blockstore.advanceTip", "probe tip", err)
	}
	if hasTip {
		_, currentHeight, tipErr := s.Tip()
		if tipErr != nil {
			return core.E("blockstore.advanceTip", "read tip", tipErr)
		}
		if height < currentHeight {
			return nil
		}
	}
	return s.SetTip(hash, height)
}
