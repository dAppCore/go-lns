// SPDX-License-Identifier: EUPL-1.2

package primitives

// NameStateStatus represents the auction/registration phase for a name.
type NameStateStatus uint8

const (
	// NameStateOpening is the initial pre-bid phase.
	NameStateOpening NameStateStatus = iota

	// NameStateLocked is the post-claim lockup phase.
	NameStateLocked

	// NameStateBidding is the sealed-bid auction phase.
	NameStateBidding

	// NameStateReveal is the bid-reveal phase.
	NameStateReveal

	// NameStateClosed is the post-auction registration phase.
	NameStateClosed

	// NameStateRevoked is the revocation state.
	NameStateRevoked
)

// String returns the human-readable phase name.
func (s NameStateStatus) String() string {
	switch s {
	case NameStateOpening:
		return "OPENING"
	case NameStateLocked:
		return "LOCKED"
	case NameStateBidding:
		return "BIDDING"
	case NameStateReveal:
		return "REVEAL"
	case NameStateClosed:
		return "CLOSED"
	case NameStateRevoked:
		return "REVOKED"
	default:
		return "UNKNOWN"
	}
}

// NameStateRules captures the period lengths needed to evaluate a name's phase.
//
// The structure keeps the logic independent from any larger network object,
// while still matching the phase transitions from the JS reference.
type NameStateRules struct {
	TreeInterval  uint32
	LockupPeriod  uint32
	BiddingPeriod uint32
	RevealPeriod  uint32
}

// State returns the current phase for the name at the provided height.
func (ns NameState) State(height uint32, rules NameStateRules) NameStateStatus {
	if ns.Revoked != 0 {
		return NameStateRevoked
	}

	if ns.Claimed != 0 {
		if height < ns.Height+rules.LockupPeriod {
			return NameStateLocked
		}
		return NameStateClosed
	}

	openPeriod := rules.TreeInterval + 1
	if height < ns.Height+openPeriod {
		return NameStateOpening
	}

	if height < ns.Height+openPeriod+rules.BiddingPeriod {
		return NameStateBidding
	}

	if height < ns.Height+openPeriod+rules.BiddingPeriod+rules.RevealPeriod {
		return NameStateReveal
	}

	return NameStateClosed
}

// IsOpening reports whether the name is currently in the opening phase.
func (ns NameState) IsOpening(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateOpening
}

// IsLocked reports whether the name is currently in the claim lockup phase.
func (ns NameState) IsLocked(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateLocked
}

// IsBidding reports whether the name is currently in the bidding phase.
func (ns NameState) IsBidding(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateBidding
}

// IsReveal reports whether the name is currently in the reveal phase.
func (ns NameState) IsReveal(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateReveal
}

// IsClosed reports whether the name is currently in the closed phase.
func (ns NameState) IsClosed(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateClosed
}

// IsRevoked reports whether the name has been revoked.
func (ns NameState) IsRevoked(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateRevoked
}

// IsRedeemable reports whether the name has reached the closed state.
func (ns NameState) IsRedeemable(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) >= NameStateClosed
}
