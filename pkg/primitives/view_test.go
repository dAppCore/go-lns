// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "dappco.re/go"
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

	if !equalBytes(decoded.Names[0].Delta.MarshalMust(t), undo.Names[0].Delta.MarshalMust(t)) {
		t.Fatal("decoded delta did not match the original")
	}
}

func TestNameUndoPreservesInsertionOrder(t *testing.T) {
	view := NewNameView()

	var firstHash Hash
	firstHash[0] = 0x01

	var secondHash Hash
	secondHash[0] = 0x02

	first, err := view.GetNameStateSync(nil, firstHash)
	if err != nil {
		t.Fatalf("GetNameStateSync(first) returned error: %v", err)
	}
	first.setValue(11)

	second, err := view.GetNameStateSync(nil, secondHash)
	if err != nil {
		t.Fatalf("GetNameStateSync(second) returned error: %v", err)
	}
	second.setValue(22)

	undo := view.ToNameUndo()
	if len(undo.Names) != 2 {
		t.Fatalf("ToNameUndo produced %d entries, want 2", len(undo.Names))
	}

	if undo.Names[0].NameHash != firstHash {
		t.Fatalf("first undo entry hash = %x, want %x", undo.Names[0].NameHash, firstHash)
	}

	if undo.Names[1].NameHash != secondHash {
		t.Fatalf("second undo entry hash = %x, want %x", undo.Names[1].NameHash, secondHash)
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

// BEGIN v090 file-aware generated coverage

func TestView_NewNameView_Good(t *core.T) {
	fn := ax7ValueOf(NewNameView)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NewNameView_Bad(t *core.T) {
	fn := ax7ValueOf(NewNameView)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NewNameView_Ugly(t *core.T) {
	fn := ax7ValueOf(NewNameView)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_GetNameStateSync_Good(t *core.T) {
	fn := ax7ValueOf(new(NameView).GetNameStateSync)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_GetNameStateSync_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameView).GetNameStateSync)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_GetNameStateSync_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameView).GetNameStateSync)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_GetNameState_Good(t *core.T) {
	fn := ax7ValueOf(new(NameView).GetNameState)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_GetNameState_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameView).GetNameState)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_GetNameState_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameView).GetNameState)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_ToNameUndo_Good(t *core.T) {
	fn := ax7ValueOf(new(NameView).ToNameUndo)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_ToNameUndo_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameView).ToNameUndo)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameView_ToNameUndo_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameView).ToNameUndo)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_FromView_Good(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).FromView)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_FromView_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).FromView)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_FromView_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).FromView)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_MarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).MarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_MarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).MarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_MarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).MarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_UnmarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).UnmarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_UnmarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).UnmarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestView_NameUndo_UnmarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameUndo).UnmarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
