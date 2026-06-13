// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	"errors"
	"testing"
)

// TestDelete_Good checks DeleteBlock removes the block, its height index, its
// header, and its undo data together.
func TestDelete_Good(t *testing.T) {
	bs := newTestStore(t)
	hash := hashOf(0x50)

	if err := bs.PutBlock(hash, 3, []byte("doomed")); err != nil {
		t.Fatalf("PutBlock: %v", err)
	}
	if err := bs.PutHeader(hash, []byte("doomed-header")); err != nil {
		t.Fatalf("PutHeader: %v", err)
	}
	if err := bs.PutUndo(hash, []byte("doomed-undo")); err != nil {
		t.Fatalf("PutUndo: %v", err)
	}

	if err := bs.DeleteBlock(hash); err != nil {
		t.Fatalf("DeleteBlock: %v", err)
	}

	if has, _ := bs.HasBlock(hash); has {
		t.Fatal("block present after delete")
	}
	if has, _ := bs.HasHeight(3); has {
		t.Fatal("height index present after delete")
	}
	if has, _ := bs.HasHeader(hash); has {
		t.Fatal("header present after delete")
	}
	if has, _ := bs.HasUndo(hash); has {
		t.Fatal("undo present after delete")
	}
}

// TestDelete_Bad checks deleting a hash that is not stored is a no-op success:
// the end state (block absent) is reached either way.
func TestDelete_Bad(t *testing.T) {
	bs := newTestStore(t)
	if err := bs.DeleteBlock(hashOf(0x99)); err != nil {
		t.Fatalf("DeleteBlock(absent) = %v, want nil", err)
	}
}

// TestDelete_Ugly checks the height index is only removed when it still points at
// the block being deleted. Storing a competing block at the same height moves the
// index to the competitor; deleting the original must leave the competitor's
// index entry intact.
func TestDelete_Ugly(t *testing.T) {
	bs := newTestStore(t)
	original, competitor := hashOf(0x60), hashOf(0x61)

	if err := bs.PutBlock(original, 7, []byte("original")); err != nil {
		t.Fatalf("PutBlock original: %v", err)
	}
	// A competing block at the same height takes over the height index.
	if err := bs.PutBlock(competitor, 7, []byte("competitor")); err != nil {
		t.Fatalf("PutBlock competitor: %v", err)
	}
	if hash, _ := bs.HashAtHeight(7); hash != competitor {
		t.Fatalf("HashAtHeight 7 = %x, want competitor %x", hash, competitor)
	}

	// Deleting the original must not disturb the competitor's index entry.
	if err := bs.DeleteBlock(original); err != nil {
		t.Fatalf("DeleteBlock original: %v", err)
	}
	if hash, err := bs.HashAtHeight(7); err != nil || hash != competitor {
		t.Fatalf("HashAtHeight 7 after delete = (%x, %v), want competitor", hash, err)
	}
	if has, _ := bs.HasBlock(competitor); !has {
		t.Fatal("competitor block removed by deleting original")
	}
}

// TestDelete_NilReceiver_Ugly checks DeleteBlock tolerates a nil receiver.
func TestDelete_NilReceiver_Ugly(t *testing.T) {
	var nilStore *Store
	if err := nilStore.DeleteBlock(hashOf(1)); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil DeleteBlock = %v, want ErrNilStore", err)
	}
}
