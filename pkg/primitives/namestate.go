// SPDX-License-Identifier: EUPL-1.2

package primitives

// Delta returns the cached sparse delta, allocating it on demand.
func (ns *NameState) Delta() *NameDelta {
	if ns.delta == nil {
		ns.delta = &NameDelta{}
	}

	return ns.delta
}

// Inject copies another name state into this one.
func (ns *NameState) Inject(other NameState) *NameState {
	ns.Name = other.Name
	ns.NameHash = other.NameHash
	ns.Height = other.Height
	ns.Renewal = other.Renewal
	ns.Owner = other.Owner
	ns.Value = other.Value
	ns.Highest = other.Highest
	ns.Data = other.Data
	ns.Transfer = other.Transfer
	ns.Revoked = other.Revoked
	ns.Claimed = other.Claimed
	ns.Renewals = other.Renewals
	ns.Registered = other.Registered
	ns.Expired = other.Expired
	ns.Weak = other.Weak
	return ns
}

// Clone returns a shallow copy of the name state.
func (ns NameState) Clone() NameState {
	return NameState{
		Name:       ns.Name,
		NameHash:   ns.NameHash,
		Height:     ns.Height,
		Renewal:    ns.Renewal,
		Owner:      ns.Owner,
		Value:      ns.Value,
		Highest:    ns.Highest,
		Data:       ns.Data,
		Transfer:   ns.Transfer,
		Revoked:    ns.Revoked,
		Claimed:    ns.Claimed,
		Renewals:   ns.Renewals,
		Registered: ns.Registered,
		Expired:    ns.Expired,
		Weak:       ns.Weak,
	}
}

// Clear discards any cached delta and returns the receiver.
func (ns *NameState) Clear() *NameState {
	ns.delta = nil
	return ns
}

// HasDelta reports whether a cached delta exists and records changes.
func (ns NameState) HasDelta() bool {
	return ns.delta != nil && !ns.delta.IsNull()
}

// ApplyState merges a sparse delta into the name state.
func (ns *NameState) ApplyState(delta *NameDelta) *NameState {
	if delta == nil {
		return ns
	}

	if delta.Height != nil {
		ns.Height = *delta.Height
	}

	if delta.Renewal != nil {
		ns.Renewal = *delta.Renewal
	}

	if delta.Owner != nil {
		ns.Owner = *delta.Owner
	}

	if delta.Value != nil {
		ns.Value = *delta.Value
	}

	if delta.Highest != nil {
		ns.Highest = *delta.Highest
	}

	if delta.Data != nil {
		ns.Data = *delta.Data
	}

	if delta.Transfer != nil {
		ns.Transfer = *delta.Transfer
	}

	if delta.Revoked != nil {
		ns.Revoked = *delta.Revoked
	}

	if delta.Claimed != nil {
		ns.Claimed = *delta.Claimed
	}

	if delta.Renewals != nil {
		ns.Renewals = *delta.Renewals
	}

	if delta.Registered != nil {
		ns.Registered = *delta.Registered
	}

	if delta.Expired != nil {
		ns.Expired = *delta.Expired
	}

	if delta.Weak != nil {
		ns.Weak = *delta.Weak
	}

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
