// SPDX-License-Identifier: EUPL-1.2

package primitives

// NameDelta tracks a sparse update to a name state.
//
// The zero value means "no changes recorded". Pointer fields let the delta
// distinguish between "not set" and "set to the zero value".
type NameDelta struct {
	Height     *uint32
	Renewal    *uint32
	Owner      *Outpoint
	Value      *uint64
	Highest    *uint64
	Data       *[]byte
	Transfer   *uint32
	Revoked    *uint32
	Claimed    *uint32
	Renewals   *uint32
	Registered *bool
	Expired    *bool
	Weak       *bool
}

// IsNull reports whether the delta contains any recorded changes.
func (d NameDelta) IsNull() bool {
	return d.Height == nil &&
		d.Renewal == nil &&
		d.Owner == nil &&
		d.Value == nil &&
		d.Highest == nil &&
		d.Data == nil &&
		d.Transfer == nil &&
		d.Revoked == nil &&
		d.Claimed == nil &&
		d.Renewals == nil &&
		d.Registered == nil &&
		d.Expired == nil &&
		d.Weak == nil
}

// Clone returns a shallow copy of the delta.
func (d NameDelta) Clone() NameDelta {
	return d
}

// Inject copies another delta into this one.
func (d *NameDelta) Inject(other NameDelta) *NameDelta {
	*d = other
	return d
}

// Clear resets the delta to its zero value.
func (d *NameDelta) Clear() *NameDelta {
	*d = NameDelta{}
	return d
}
