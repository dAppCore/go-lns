// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	core "dappco.re/go"

	"dappco.re/go/lns/pkg/primitives"
)

// PutHeader stores the serialised header bytes for a block hash. Headers are
// kept in their own group so a header-only client, or a routine walking the
// chain to check proof-of-work, can read headers without paying to decode whole
// blocks. blockstore treats the header as opaque bytes; the caller serialises a
// primitives.BlockHeader and hands the bytes over. An empty header is rejected
// with ErrEmptyBlock.
//
//	if err := bs.PutHeader(hash, headerBytes); err != nil { return err }
func (s *Store) PutHeader(hash primitives.Hash, header []byte) error {
	if s == nil || s.scoped == nil {
		return core.E("blockstore.PutHeader", "store header", ErrNilStore)
	}
	if len(header) == 0 {
		return core.E("blockstore.PutHeader", "store header", ErrEmptyBlock)
	}
	if len(header) > maxBlockBytes {
		return core.E("blockstore.PutHeader",
			core.Sprintf("header of %d bytes exceeds maximum %d", len(header), maxBlockBytes), ErrCorrupt)
	}

	if err := s.putBytes(groupHeader, hashKey(hash), header); err != nil {
		return core.E("blockstore.PutHeader", "write header", err)
	}
	return nil
}

// GetHeader returns the serialised header bytes stored under hash. It returns
// ErrNotFound when no header is stored for that hash and ErrCorrupt when the
// stored value is not valid hex.
//
//	header, err := bs.GetHeader(hash)
func (s *Store) GetHeader(hash primitives.Hash) ([]byte, error) {
	if s == nil || s.scoped == nil {
		return nil, core.E("blockstore.GetHeader", "read header", ErrNilStore)
	}

	header, err := s.getBytes(groupHeader, hashKey(hash))
	if err != nil {
		return nil, core.E("blockstore.GetHeader", "read header", err)
	}
	return header, nil
}

// HasHeader reports whether a header is stored under hash.
//
//	if has, _ := bs.HasHeader(hash); has { /* header already indexed */ }
func (s *Store) HasHeader(hash primitives.Hash) (bool, error) {
	if s == nil || s.scoped == nil {
		return false, core.E("blockstore.HasHeader", "probe header", ErrNilStore)
	}

	exists, err := s.scoped.ExistsIn(groupHeader, hashKey(hash))
	if err != nil {
		return false, core.E("blockstore.HasHeader", "probe header", err)
	}
	return exists, nil
}
