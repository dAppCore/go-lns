// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	"bytes"
	"errors"
	"testing"
)

// TestHeader_Good checks a header round-trips by hash and HasHeader reports it.
func TestHeader_Good(t *testing.T) {
	bs := newTestStore(t)
	hash := hashOf(0x10)
	header := []byte("serialised-block-header-bytes")

	if err := bs.PutHeader(hash, header); err != nil {
		t.Fatalf("PutHeader: %v", err)
	}
	got, err := bs.GetHeader(hash)
	if err != nil {
		t.Fatalf("GetHeader: %v", err)
	}
	if !bytes.Equal(got, header) {
		t.Fatalf("GetHeader = %q, want %q", got, header)
	}
	if has, _ := bs.HasHeader(hash); !has {
		t.Fatal("HasHeader false after PutHeader")
	}
}

// TestHeader_Bad checks an empty header is rejected and a missing header reads as
// ErrNotFound.
func TestHeader_Bad(t *testing.T) {
	bs := newTestStore(t)
	hash := hashOf(0x20)

	if err := bs.PutHeader(hash, nil); !errors.Is(err, ErrEmptyBlock) {
		t.Fatalf("PutHeader(nil) = %v, want ErrEmptyBlock", err)
	}
	if _, err := bs.GetHeader(hashOf(0xee)); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetHeader(absent) = %v, want ErrNotFound", err)
	}
	if has, _ := bs.HasHeader(hashOf(0xee)); has {
		t.Fatal("HasHeader true for absent header")
	}
}

// TestHeader_Ugly checks the header helpers tolerate a nil receiver.
func TestHeader_Ugly(t *testing.T) {
	var nilStore *Store
	if err := nilStore.PutHeader(hashOf(1), []byte("h")); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil PutHeader = %v, want ErrNilStore", err)
	}
	if _, err := nilStore.GetHeader(hashOf(1)); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil GetHeader = %v, want ErrNilStore", err)
	}
}

// TestUndo_Good checks undo data round-trips by hash and HasUndo reports it.
func TestUndo_Good(t *testing.T) {
	bs := newTestStore(t)
	hash := hashOf(0x30)
	undo := []byte("undo-coins-spent-by-this-block")

	if err := bs.PutUndo(hash, undo); err != nil {
		t.Fatalf("PutUndo: %v", err)
	}
	got, err := bs.GetUndo(hash)
	if err != nil {
		t.Fatalf("GetUndo: %v", err)
	}
	if !bytes.Equal(got, undo) {
		t.Fatalf("GetUndo = %q, want %q", got, undo)
	}
	if has, _ := bs.HasUndo(hash); !has {
		t.Fatal("HasUndo false after PutUndo")
	}
}

// TestUndo_Bad checks an empty undo blob is rejected and a missing one reads as
// ErrNotFound.
func TestUndo_Bad(t *testing.T) {
	bs := newTestStore(t)

	if err := bs.PutUndo(hashOf(0x40), []byte{}); !errors.Is(err, ErrEmptyBlock) {
		t.Fatalf("PutUndo(empty) = %v, want ErrEmptyBlock", err)
	}
	if _, err := bs.GetUndo(hashOf(0xab)); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetUndo(absent) = %v, want ErrNotFound", err)
	}
	if has, _ := bs.HasUndo(hashOf(0xab)); has {
		t.Fatal("HasUndo true for absent undo")
	}
}

// TestUndo_Ugly checks the undo helpers tolerate a nil receiver.
func TestUndo_Ugly(t *testing.T) {
	var nilStore *Store
	if err := nilStore.PutUndo(hashOf(1), []byte("u")); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil PutUndo = %v, want ErrNilStore", err)
	}
	if _, err := nilStore.GetUndo(hashOf(1)); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil GetUndo = %v, want ErrNilStore", err)
	}
}
