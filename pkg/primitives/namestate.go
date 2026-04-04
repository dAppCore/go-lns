// SPDX-License-Identifier: EUPL-1.2

package primitives

import "bytes"

// Delta returns the cached sparse delta, allocating it on demand.
func (ns *NameState) Delta() *NameDelta {
	if ns.delta == nil {
		ns.delta = &NameDelta{}
	}

	return ns.delta
}

// GetDelta is an alias for Delta.
func (ns *NameState) GetDelta() *NameDelta {
	return ns.Delta()
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

// GetInject is an alias for Inject.
func (ns *NameState) GetInject(other NameState) *NameState {
	return ns.Inject(other)
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

// GetClone is an alias for Clone.
func (ns NameState) GetClone() NameState {
	return ns.Clone()
}

// Clear discards any cached delta and returns the receiver.
func (ns *NameState) Clear() *NameState {
	ns.delta = nil
	return ns
}

// GetClear is an alias for Clear.
func (ns *NameState) GetClear() *NameState {
	return ns.Clear()
}

// HasDelta reports whether a cached delta exists and records changes.
func (ns NameState) HasDelta() bool {
	return ns.delta != nil && !ns.delta.IsNull()
}

// GetHasDelta is an alias for HasDelta.
func (ns NameState) GetHasDelta() bool {
	return ns.HasDelta()
}

func (ns *NameState) setHeight(height uint32) *NameState {
	if ns.Height == height {
		return ns
	}

	delta := ns.Delta()
	if delta.Height == nil {
		old := ns.Height
		delta.Height = &old
	}

	ns.Height = height
	return ns
}

func (ns *NameState) setRenewal(renewal uint32) *NameState {
	if ns.Renewal == renewal {
		return ns
	}

	delta := ns.Delta()
	if delta.Renewal == nil {
		old := ns.Renewal
		delta.Renewal = &old
	}

	ns.Renewal = renewal
	return ns
}

func outpointEqual(a, b Outpoint) bool {
	return a.TxHash == b.TxHash && a.Index == b.Index
}

func outpointIsNull(owner Outpoint) bool {
	return owner == (Outpoint{}) || owner.IsNull()
}

func (ns *NameState) setOwner(owner Outpoint) *NameState {
	if outpointEqual(ns.Owner, owner) {
		return ns
	}

	delta := ns.Delta()
	if delta.Owner == nil {
		old := ns.Owner
		delta.Owner = &old
	}

	ns.Owner = owner
	return ns
}

func (ns *NameState) setValue(value uint64) *NameState {
	if ns.Value == value {
		return ns
	}

	delta := ns.Delta()
	if delta.Value == nil {
		old := ns.Value
		delta.Value = &old
	}

	ns.Value = value
	return ns
}

func (ns *NameState) setHighest(highest uint64) *NameState {
	if ns.Highest == highest {
		return ns
	}

	delta := ns.Delta()
	if delta.Highest == nil {
		old := ns.Highest
		delta.Highest = &old
	}

	ns.Highest = highest
	return ns
}

func (ns *NameState) setData(data []byte) *NameState {
	if data == nil {
		data = []byte{}
	}

	if bytes.Equal(ns.Data, data) {
		return ns
	}

	delta := ns.Delta()
	if delta.Data == nil {
		old := ns.Data
		delta.Data = &old
	}

	ns.Data = data
	return ns
}

func (ns *NameState) setTransfer(transfer uint32) *NameState {
	if ns.Transfer == transfer {
		return ns
	}

	delta := ns.Delta()
	if delta.Transfer == nil {
		old := ns.Transfer
		delta.Transfer = &old
	}

	ns.Transfer = transfer
	return ns
}

func (ns *NameState) setRevoked(revoked uint32) *NameState {
	if ns.Revoked == revoked {
		return ns
	}

	delta := ns.Delta()
	if delta.Revoked == nil {
		old := ns.Revoked
		delta.Revoked = &old
	}

	ns.Revoked = revoked
	return ns
}

func (ns *NameState) setClaimed(claimed uint32) *NameState {
	if ns.Claimed == claimed {
		return ns
	}

	delta := ns.Delta()
	if delta.Claimed == nil {
		old := ns.Claimed
		delta.Claimed = &old
	}

	ns.Claimed = claimed
	return ns
}

func (ns *NameState) setRenewals(renewals uint32) *NameState {
	if ns.Renewals == renewals {
		return ns
	}

	delta := ns.Delta()
	if delta.Renewals == nil {
		old := ns.Renewals
		delta.Renewals = &old
	}

	ns.Renewals = renewals
	return ns
}

func (ns *NameState) setRegistered(registered bool) *NameState {
	if ns.Registered == registered {
		return ns
	}

	delta := ns.Delta()
	if delta.Registered == nil {
		old := ns.Registered
		delta.Registered = &old
	}

	ns.Registered = registered
	return ns
}

func (ns *NameState) setExpired(expired bool) *NameState {
	if ns.Expired == expired {
		return ns
	}

	delta := ns.Delta()
	if delta.Expired == nil {
		old := ns.Expired
		delta.Expired = &old
	}

	ns.Expired = expired
	return ns
}

func (ns *NameState) setWeak(weak bool) *NameState {
	if ns.Weak == weak {
		return ns
	}

	delta := ns.Delta()
	if delta.Weak == nil {
		old := ns.Weak
		delta.Weak = &old
	}

	ns.Weak = weak
	return ns
}

// Reset clears the chain-derived fields while keeping the name metadata.
func (ns *NameState) Reset(height uint32) *NameState {
	ns.setHeight(height)
	ns.setRenewal(height)
	ns.setOwner(NewOutpoint())
	ns.setValue(0)
	ns.setHighest(0)
	ns.setData(nil)
	ns.setTransfer(0)
	ns.setRevoked(0)
	ns.setClaimed(0)
	ns.setRenewals(0)
	ns.setRegistered(false)
	ns.setExpired(false)
	ns.setWeak(false)

	return ns
}

// GetReset is an alias for Reset.
func (ns *NameState) GetReset(height uint32) *NameState {
	return ns.Reset(height)
}

// Set assigns the raw name and resets the mutable chain state.
func (ns *NameState) Set(name []byte, height uint32) *NameState {
	ns.Name = name
	ns.Reset(height)
	return ns
}

// GetSet is an alias for Set.
func (ns *NameState) GetSet(name []byte, height uint32) *NameState {
	return ns.Set(name, height)
}

// MaybeExpire resets the state if it has expired, preserving the current data.
func (ns *NameState) MaybeExpire(height uint32, rules NameStateRules) bool {
	if !ns.IsExpired(height, rules) {
		return false
	}

	data := ns.Data
	ns.Reset(height)
	ns.setExpired(true)
	ns.setData(data)

	return true
}

// GetMaybeExpire is an alias for MaybeExpire.
func (ns *NameState) GetMaybeExpire(height uint32, rules NameStateRules) bool {
	return ns.MaybeExpire(height, rules)
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

// GetApplyState is an alias for ApplyState.
func (ns *NameState) GetApplyState(delta *NameDelta) *NameState {
	return ns.ApplyState(delta)
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

// GetIsNull is an alias for IsNull.
func (ns NameState) GetIsNull() bool {
	return ns.IsNull()
}
