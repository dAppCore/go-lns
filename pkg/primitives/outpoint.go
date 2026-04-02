// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// NewOutpoint returns the null outpoint used by coinbase inputs and empty
// name-state owner references.
func NewOutpoint() Outpoint {
	return Outpoint{Index: ^uint32(0)}
}

// Inject copies another outpoint into this one.
func (o *Outpoint) Inject(other Outpoint) *Outpoint {
	o.TxHash = other.TxHash
	o.Index = other.Index
	return o
}

// Clone returns a shallow copy of the outpoint.
func (o Outpoint) Clone() Outpoint {
	var clone Outpoint
	return *clone.Inject(o)
}

// Equals reports whether two outpoints refer to the same transaction output.
func (o Outpoint) Equals(other Outpoint) bool {
	return o.TxHash == other.TxHash && o.Index == other.Index
}

// Compare orders two outpoints by hash and then by output index.
func (o Outpoint) Compare(other Outpoint) int {
	if cmp := bytes.Compare(o.TxHash[:], other.TxHash[:]); cmp != 0 {
		return cmp
	}

	switch {
	case o.Index < other.Index:
		return -1
	case o.Index > other.Index:
		return 1
	default:
		return 0
	}
}

// GetSize returns the serialized size of the outpoint.
func (o Outpoint) GetSize() int {
	return len(o.TxHash) + 4
}

// MarshalBinary serializes the outpoint into the reference wire format.
func (o Outpoint) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 0, o.GetSize())
	buf = append(buf, o.TxHash[:]...)

	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], o.Index)
	buf = append(buf, tmp[:]...)
	return buf, nil
}

// UnmarshalBinary decodes the outpoint from the reference wire format.
func (o *Outpoint) UnmarshalBinary(data []byte) error {
	if len(data) != o.GetSize() {
		return fmt.Errorf("primitives.Outpoint.UnmarshalBinary: invalid length")
	}

	copy(o.TxHash[:], data[:len(o.TxHash)])
	o.Index = binary.LittleEndian.Uint32(data[len(o.TxHash):])
	return nil
}
