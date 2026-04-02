// SPDX-License-Identifier: EUPL-1.2

package primitives

// Inject copies another name state into this one.
func (ns *NameState) Inject(other NameState) *NameState {
	*ns = other
	return ns
}

// Clone returns a shallow copy of the name state.
func (ns NameState) Clone() NameState {
	return ns
}

// IsNull reports whether the name state carries any non-default data.
func (ns NameState) IsNull() bool {
	return ns.Height == 0 &&
		ns.Renewal == 0 &&
		(ns.Owner == (Outpoint{}) || ns.Owner.IsNull()) &&
		ns.Value == 0 &&
		ns.Highest == 0 &&
		len(ns.Data) == 0 &&
		ns.Transfer == 0 &&
		ns.Revoked == 0 &&
		ns.Claimed == 0 &&
		ns.Renewals == 0 &&
		!ns.Registered &&
		!ns.Expired &&
		!ns.Weak
}
