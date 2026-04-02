// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"bytes"
	"testing"
)

type stubNameStateDB struct {
	state *NameState
	err   error
}

func (db stubNameStateDB) GetNameState(hash Hash) (*NameState, error) {
	if db.err != nil {
		return nil, db.err
	}

	if db.state == nil {
		return nil, nil
	}

	cp := db.state.Clone()
	cp.NameHash = hash
	return &cp, nil
}

func TestNameViewGetNameState(t *testing.T) {
	var hash Hash
	hash[0] = 0x42

	view := NewNameView()
	db := stubNameStateDB{
		state: &NameState{Name: []byte("example"), Height: 7},
	}

	ns, err := view.GetNameStateSync(db, hash)
	if err != nil {
		t.Fatalf("GetNameStateSync returned error: %v", err)
	}

	if ns == nil {
		t.Fatal("GetNameStateSync returned nil state")
	}

	if ns.NameHash != hash {
		t.Fatalf("GetNameStateSync returned hash %x, want %x", ns.NameHash, hash)
	}

	ns2, err := view.GetNameState(db, hash)
	if err != nil {
		t.Fatalf("GetNameState returned error: %v", err)
	}

	if ns2 != ns {
		t.Fatal("GetNameState should return the cached state")
	}

	missingHash := Hash{0x99}
	empty, err := view.GetNameState(nil, missingHash)
	if err != nil {
		t.Fatalf("GetNameState(nil) returned error: %v", err)
	}

	if empty == nil || empty.NameHash != missingHash {
		t.Fatal("GetNameState(nil) should create an empty state with the requested hash")
	}
}

func TestNameUndoRoundTrip(t *testing.T) {
	var hash Hash
	hash[0] = 0x51

	view := NewNameView()
	ns := &NameState{NameHash: hash, Name: []byte("example"), Height: 10}
	ns.setValue(123)
	view.names[hash] = ns

	undo := view.ToNameUndo()
	if len(undo.Names) != 1 {
		t.Fatalf("ToNameUndo produced %d entries, want 1", len(undo.Names))
	}

	encoded, err := undo.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	if got := undo.GetSize(); got != len(encoded) {
		t.Fatalf("GetSize() = %d, want %d", got, len(encoded))
	}

	var decoded NameUndo
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if len(decoded.Names) != 1 {
		t.Fatalf("decoded undo has %d entries, want 1", len(decoded.Names))
	}

	if decoded.Names[0].NameHash != hash {
		t.Fatalf("decoded hash %x, want %x", decoded.Names[0].NameHash, hash)
	}

	if !bytes.Equal(decoded.Names[0].Delta.MarshalMust(t), undo.Names[0].Delta.MarshalMust(t)) {
		t.Fatal("decoded delta did not match the original")
	}
}

func (d NameDelta) MarshalMust(t *testing.T) []byte {
	t.Helper()

	data, err := d.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	return data
}
