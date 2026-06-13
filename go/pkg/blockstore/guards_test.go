// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	"errors"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/store"
)

// TestGuards_NilReceiver_Ugly checks every read and probe helper reports
// ErrNilStore on a nil receiver rather than panicking, so a misconstructed store
// fails loudly at the call site.
func TestGuards_NilReceiver_Ugly(t *testing.T) {
	var s *Store
	h := hashOf(0x01)

	if _, err := s.GetBlock(h); !errors.Is(err, ErrNilStore) {
		t.Errorf("GetBlock nil = %v", err)
	}
	if _, _, err := s.GetBlockWithHeight(h); !errors.Is(err, ErrNilStore) {
		t.Errorf("GetBlockWithHeight nil = %v", err)
	}
	if _, err := s.GetBlockByHeight(1); !errors.Is(err, ErrNilStore) {
		t.Errorf("GetBlockByHeight nil = %v", err)
	}
	if _, err := s.HashAtHeight(1); !errors.Is(err, ErrNilStore) {
		t.Errorf("HashAtHeight nil = %v", err)
	}
	if _, err := s.HasBlock(h); !errors.Is(err, ErrNilStore) {
		t.Errorf("HasBlock nil = %v", err)
	}
	if _, err := s.HasHeight(1); !errors.Is(err, ErrNilStore) {
		t.Errorf("HasHeight nil = %v", err)
	}
	if _, err := s.HasHeader(h); !errors.Is(err, ErrNilStore) {
		t.Errorf("HasHeader nil = %v", err)
	}
	if _, err := s.HasUndo(h); !errors.Is(err, ErrNilStore) {
		t.Errorf("HasUndo nil = %v", err)
	}
}

// TestGuards_CorruptRecordDelete_Ugly checks DeleteBlock still clears a block
// whose stored record is corrupt (cannot be decoded): the bad value is removed so
// it does not linger, even though its height cannot be matched in the index.
func TestGuards_CorruptRecordDelete_Ugly(t *testing.T) {
	storeInstance, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { _ = storeInstance.Close() })
	scoped := store.NewScoped(storeInstance, "corrupt-delete")
	bs, _ := New(scoped)

	hash := hashOf(0x70)
	// Valid hex, but a record too short to decode (header < 8 bytes).
	if err := scoped.SetIn(groupBlock, hashKey(hash), core.HexEncode([]byte{0x01})); err != nil {
		t.Fatalf("seed corrupt record: %v", err)
	}
	if has, _ := bs.HasBlock(hash); !has {
		t.Fatal("seed not visible")
	}

	if err := bs.DeleteBlock(hash); err != nil {
		t.Fatalf("DeleteBlock(corrupt) = %v, want nil", err)
	}
	if has, _ := bs.HasBlock(hash); has {
		t.Fatal("corrupt block still present after delete")
	}
}

// TestGuards_PutHeaderUndoOversize_Bad checks the oversize guards on PutHeader and
// PutUndo reject a blob larger than the maximum.
func TestGuards_PutHeaderUndoOversize_Bad(t *testing.T) {
	bs := newTestStore(t)
	oversize := make([]byte, maxBlockBytes+1)

	if err := bs.PutHeader(hashOf(0x71), oversize); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("PutHeader(oversize) = %v, want ErrCorrupt", err)
	}
	if err := bs.PutUndo(hashOf(0x72), oversize); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("PutUndo(oversize) = %v, want ErrCorrupt", err)
	}
}

// TestGuards_HashAtHeightCorrupt_Ugly checks HashAtHeight rejects an index entry
// whose stored value is not a 32-byte hash.
func TestGuards_HashAtHeightCorrupt_Ugly(t *testing.T) {
	storeInstance, err := store.New(":memory:")
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() { _ = storeInstance.Close() })
	scoped := store.NewScoped(storeInstance, "corrupt-index")
	bs, _ := New(scoped)

	// A 4-byte value where a 32-byte hash is expected.
	if err := scoped.SetIn(groupHeight, heightKey(3), core.HexEncode([]byte{1, 2, 3, 4})); err != nil {
		t.Fatalf("seed short index: %v", err)
	}
	if _, err := bs.HashAtHeight(3); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("HashAtHeight(short) = %v, want ErrCorrupt", err)
	}
}
