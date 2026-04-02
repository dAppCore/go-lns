// SPDX-License-Identifier: EUPL-1.2

package primitives

import "math"

// NameStateStats describes the derived timing information for a name.
//
// The fields mirror the timing summary produced by the JS reference
// implementation's `NameState.toStats()` helper.
type NameStateStats struct {
	OpenPeriodStart          uint32
	OpenPeriodEnd            uint32
	BlocksUntilBidding       uint32
	HoursUntilBidding        float64
	LockupPeriodStart        uint32
	LockupPeriodEnd          uint32
	BlocksUntilClosed        uint32
	HoursUntilClosed         float64
	BidPeriodStart           uint32
	BidPeriodEnd             uint32
	BlocksUntilReveal        uint32
	HoursUntilReveal         float64
	RevealPeriodStart        uint32
	RevealPeriodEnd          uint32
	BlocksUntilClose         uint32
	HoursUntilClose          float64
	RenewalPeriodStart       uint32
	RenewalPeriodEnd         uint32
	BlocksUntilExpire        uint32
	DaysUntilExpire          float64
	TransferLockupStart      uint32
	TransferLockupEnd        uint32
	BlocksUntilValidFinalize uint32
	HoursUntilValidFinalize  float64
	RevokePeriodStart        uint32
	RevokePeriodEnd          uint32
	BlocksUntilReopen        uint32
	HoursUntilReopen         float64
	BlocksSinceExpired       uint32
}

func blocksToHours(blocks uint32, blocksPerHour float64) float64 {
	if blocksPerHour == 0 {
		return 0
	}

	return math.Round((float64(blocks)/blocksPerHour)*100) / 100
}

func blocksToDays(blocks uint32, blocksPerDay float64) float64 {
	if blocksPerDay == 0 {
		return 0
	}

	return math.Round((float64(blocks)/blocksPerDay)*100) / 100
}

func maxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}

	return b
}

// ToStats returns a derived timing summary for the name at the provided height.
//
// The result mirrors the JS reference implementation's `toStats()` helper.
// Hour and day estimates are rounded to two decimal places to match that API.
// It returns nil when the name is expired and has no owner reference.
func (ns NameState) ToStats(height uint32, rules NameStateRules) *NameStateStats {
	if outpointIsNull(ns.Owner) && ns.IsExpired(height, rules) {
		return nil
	}

	stats := &NameStateStats{}
	blocksPerDay := float64(rules.BlocksPerDay)
	blocksPerHour := blocksPerDay / 24
	openPeriod := rules.TreeInterval + 1

	state := ns.State(height, rules)
	expired := ns.IsExpired(height, rules)

	if expired && state != NameStateRevoked {
		blocksSinceExpired := height - (ns.Renewal + rules.RenewalWindow)
		stats.BlocksSinceExpired = blocksSinceExpired
		return stats
	}

	switch state {
	case NameStateOpening:
		start := ns.Height
		end := ns.Height + openPeriod
		blocks := end - height

		stats.OpenPeriodStart = start
		stats.OpenPeriodEnd = end
		stats.BlocksUntilBidding = blocks
		stats.HoursUntilBidding = blocksToHours(blocks, blocksPerHour)
	case NameStateLocked:
		start := ns.Height
		end := ns.Height + rules.LockupPeriod
		blocks := end - height

		stats.LockupPeriodStart = start
		stats.LockupPeriodEnd = end
		stats.BlocksUntilClosed = blocks
		stats.HoursUntilClosed = blocksToHours(blocks, blocksPerHour)
	case NameStateBidding:
		start := ns.Height + openPeriod
		end := start + rules.BiddingPeriod
		blocks := end - height

		stats.BidPeriodStart = start
		stats.BidPeriodEnd = end
		stats.BlocksUntilReveal = blocks
		stats.HoursUntilReveal = blocksToHours(blocks, blocksPerHour)
	case NameStateReveal:
		start := ns.Height + openPeriod + rules.BiddingPeriod
		end := start + rules.RevealPeriod
		blocks := end - height

		stats.RevealPeriodStart = start
		stats.RevealPeriodEnd = end
		stats.BlocksUntilClose = blocks
		stats.HoursUntilClose = blocksToHours(blocks, blocksPerHour)
	case NameStateClosed:
		start := ns.Renewal
		normalEnd := start + rules.RenewalWindow
		end := normalEnd
		if ns.Claimed != 0 {
			end = maxUint32(rules.ClaimPeriod, normalEnd)
		}

		blocks := end - height

		stats.RenewalPeriodStart = start
		stats.RenewalPeriodEnd = end
		stats.BlocksUntilExpire = blocks
		stats.DaysUntilExpire = blocksToDays(blocks, blocksPerDay)

		if ns.Transfer != 0 {
			start := ns.Transfer
			end := start + rules.TransferLockup
			blocks := end - height

			stats.TransferLockupStart = start
			stats.TransferLockupEnd = end
			stats.BlocksUntilValidFinalize = blocks
			stats.HoursUntilValidFinalize = blocksToHours(blocks, blocksPerHour)
		}
	case NameStateRevoked:
		start := ns.Revoked
		end := start + rules.AuctionMaturity
		blocks := end - height

		stats.RevokePeriodStart = start
		stats.RevokePeriodEnd = end
		stats.BlocksUntilReopen = blocks
		stats.HoursUntilReopen = blocksToHours(blocks, blocksPerHour)
	}

	return stats
}
