// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/binary"
	"fmt"
)

// NameStateGetter is the minimal lookup interface used by NameView.
type NameStateGetter interface {
	GetNameState(Hash) (*NameState, error)
}

// NameView caches name states during state-tree processing.
type NameView struct {
	names map[Hash]*NameState
}

// NewNameView creates an empty name-state view.
func NewNameView() *NameView {
	return &NameView{
		names: make(map[Hash]*NameState),
	}
}

func (v *NameView) ensureCache() {
	if v.names == nil {
		v.names = make(map[Hash]*NameState)
	}
}

// GetNameStateSync returns the cached name state or loads it from db.
//
// When the database has no entry for the requested hash, a new empty state is
// created with NameHash populated, mirroring the JS reference behavior.
func (v *NameView) GetNameStateSync(db NameStateGetter, nameHash Hash) (*NameState, error) {
	v.ensureCache()

	if ns, ok := v.names[nameHash]; ok {
		return ns, nil
	}

	if db != nil {
		ns, err := db.GetNameState(nameHash)
		if err != nil {
			return nil, err
		}

		if ns != nil {
			v.names[nameHash] = ns
			return ns, nil
		}
	}

	ns := &NameState{NameHash: nameHash}
	v.names[nameHash] = ns
	return ns, nil
}

// GetNameState is an alias for GetNameStateSync.
func (v *NameView) GetNameState(db NameStateGetter, nameHash Hash) (*NameState, error) {
	return v.GetNameStateSync(db, nameHash)
}

// ToNameUndo converts the cached view into a sparse undo record.
func (v *NameView) ToNameUndo() NameUndo {
	undo := NameUndo{}
	if v == nil {
		return undo
	}
	return undo.FromView(v)
}

// NameUndoEntry pairs a name hash with its sparse state delta.
type NameUndoEntry struct {
	NameHash Hash
	Delta    NameDelta
}

// NameUndo stores the deltas required to roll a name view back.
type NameUndo struct {
	Names []NameUndoEntry
}

// FromView populates the undo record from a cached view.
func (u *NameUndo) FromView(view *NameView) NameUndo {
	if view == nil {
		*u = NameUndo{}
		return *u
	}

	u.Names = u.Names[:0]
	for hash, ns := range view.names {
		if ns == nil || !ns.HasDelta() {
			continue
		}

		u.Names = append(u.Names, NameUndoEntry{
			NameHash: hash,
			Delta:    *ns.Delta(),
		})
	}

	return *u
}

// GetSize returns the serialized size of the undo payload.
func (u NameUndo) GetSize() int {
	size := 4
	for _, entry := range u.Names {
		size += len(entry.NameHash)
		size += entry.Delta.GetSize()
	}
	return size
}

// MarshalBinary serializes the undo payload.
func (u NameUndo) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 0, u.GetSize())

	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], uint32(len(u.Names)))
	buf = append(buf, tmp[:]...)
	for _, entry := range u.Names {
		buf = append(buf, entry.NameHash[:]...)
		delta, err := entry.Delta.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("primitives.NameUndo.MarshalBinary: %w", err)
		}
		buf = append(buf, delta...)
	}

	return buf, nil
}

// UnmarshalBinary decodes the undo payload.
func (u *NameUndo) UnmarshalBinary(data []byte) error {
	*u = NameUndo{}

	if len(data) < 4 {
		return fmt.Errorf("primitives.NameUndo.UnmarshalBinary: short buffer")
	}

	count := binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	u.Names = make([]NameUndoEntry, 0, int(count))
	for i := uint32(0); i < count; i++ {
		if len(data) < len(Hash{}) {
			return fmt.Errorf("primitives.NameUndo.UnmarshalBinary: short buffer")
		}

		var hash Hash
		copy(hash[:], data[:len(hash)])
		data = data[len(hash):]

		delta, consumed, err := decodeNameDeltaPrefix(data)
		if err != nil {
			return fmt.Errorf("primitives.NameUndo.UnmarshalBinary: %w", err)
		}

		u.Names = append(u.Names, NameUndoEntry{
			NameHash: hash,
			Delta:    delta,
		})
		data = data[consumed:]
	}

	if len(data) != 0 {
		return fmt.Errorf("primitives.NameUndo.UnmarshalBinary: trailing data")
	}

	return nil
}

func decodeNameDeltaPrefix(data []byte) (NameDelta, int, error) {
	for n := 4; n <= len(data); n++ {
		var delta NameDelta
		if err := delta.UnmarshalBinary(data[:n]); err == nil {
			return delta, n, nil
		}
	}

	return NameDelta{}, 0, fmt.Errorf("short buffer")
}
