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
	TreeInterval    uint32
	LockupPeriod    uint32
	BiddingPeriod   uint32
	RevealPeriod    uint32
	ClaimPeriod     uint32
	RenewalWindow   uint32
	AuctionMaturity uint32
	BlocksPerDay    uint32
	TransferLockup  uint32
	NoReserved      bool
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

// GetState is an alias for State.
func (ns NameState) GetState(height uint32, rules NameStateRules) NameStateStatus {
	return ns.State(height, rules)
}

// IsOpening reports whether the name is currently in the opening phase.
func (ns NameState) IsOpening(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateOpening
}

// GetIsOpening is an alias for IsOpening.
func (ns NameState) GetIsOpening(height uint32, rules NameStateRules) bool {
	return ns.IsOpening(height, rules)
}

// IsLocked reports whether the name is currently in the claim lockup phase.
func (ns NameState) IsLocked(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateLocked
}

// GetIsLocked is an alias for IsLocked.
func (ns NameState) GetIsLocked(height uint32, rules NameStateRules) bool {
	return ns.IsLocked(height, rules)
}

// IsBidding reports whether the name is currently in the bidding phase.
func (ns NameState) IsBidding(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateBidding
}

// GetIsBidding is an alias for IsBidding.
func (ns NameState) GetIsBidding(height uint32, rules NameStateRules) bool {
	return ns.IsBidding(height, rules)
}

// IsReveal reports whether the name is currently in the reveal phase.
func (ns NameState) IsReveal(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateReveal
}

// GetIsReveal is an alias for IsReveal.
func (ns NameState) GetIsReveal(height uint32, rules NameStateRules) bool {
	return ns.IsReveal(height, rules)
}

// IsClosed reports whether the name is currently in the closed phase.
func (ns NameState) IsClosed(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateClosed
}

// GetIsClosed is an alias for IsClosed.
func (ns NameState) GetIsClosed(height uint32, rules NameStateRules) bool {
	return ns.IsClosed(height, rules)
}

// IsRevoked reports whether the name has been revoked.
func (ns NameState) IsRevoked(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) == NameStateRevoked
}

// GetIsRevoked is an alias for IsRevoked.
func (ns NameState) GetIsRevoked(height uint32, rules NameStateRules) bool {
	return ns.IsRevoked(height, rules)
}

// IsRedeemable reports whether the name has reached the closed state.
func (ns NameState) IsRedeemable(height uint32, rules NameStateRules) bool {
	return ns.State(height, rules) >= NameStateClosed
}

// GetIsRedeemable is an alias for IsRedeemable.
func (ns NameState) GetIsRedeemable(height uint32, rules NameStateRules) bool {
	return ns.IsRedeemable(height, rules)
}

// IsClaimable reports whether the name is still within the claim window.
func (ns NameState) IsClaimable(height uint32, rules NameStateRules) bool {
	if ns.Claimed == 0 || rules.NoReserved {
		return false
	}

	return height < rules.ClaimPeriod
}

// GetIsClaimable is an alias for IsClaimable.
func (ns NameState) GetIsClaimable(height uint32, rules NameStateRules) bool {
	return ns.IsClaimable(height, rules)
}

// IsExpired reports whether the name should be treated as expired.
func (ns NameState) IsExpired(height uint32, rules NameStateRules) bool {
	if ns.Revoked != 0 {
		return height >= ns.Revoked+rules.AuctionMaturity
	}

	// A name can only expire after it has reached the closed state.
	if !ns.IsClosed(height, rules) {
		return false
	}

	// Claimed names can only expire once the claim period has ended.
	if ns.IsClaimable(height, rules) {
		return false
	}

	// If we have not been renewed in time, the name expires.
	if height >= ns.Renewal+rules.RenewalWindow {
		return true
	}

	// If nobody revealed their bids, restart the auction.
	return ns.Owner.IsNull()
}

// GetIsExpired is an alias for IsExpired.
func (ns NameState) GetIsExpired(height uint32, rules NameStateRules) bool {
	return ns.IsExpired(height, rules)
}
