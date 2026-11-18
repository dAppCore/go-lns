// SPDX-License-Identifier: EUPL-1.2

package net

import (
	core "github.com/dappcore/go"
)

// TestTrust_ScorePeer_Good gives a fully-verified, long-uptime, fast,
// consistent peer a high score (well above the RFC §11.3 routing thresholds).
func TestTrust_ScorePeer_Good(t *core.T) {
	score := DefaultScorer{}.ScorePeer(PeerVector{
		Uptime:        uptimeFullSeconds,
		VerifiedChain: true,
		VerifiedHNS:   true,
		ResponseTime:  10,
		ConsistencyN:  10,
	})
	core.AssertInDelta(t, 1.0, score, 1e-9)
	core.AssertFalse(t, IsStale(score))
}

// TestTrust_ScorePeer_Bad gives an unverified, zero-signal peer the floor score
// of zero, but still in range (the LNS returns it; the caller decides per §11.3).
func TestTrust_ScorePeer_Bad(t *core.T) {
	score := DefaultScorer{}.ScorePeer(PeerVector{})
	core.AssertInDelta(t, 0.0, score, 1e-9)
	core.AssertTrue(t, IsStale(score))
}

// TestTrust_ScorePeer_Ugly keeps a chain-verified-only peer in range and above
// the verification floor even with no other positive signals.
func TestTrust_ScorePeer_Ugly(t *core.T) {
	score := DefaultScorer{}.ScorePeer(PeerVector{VerifiedChain: true, ResponseTime: 9999})
	// weightChain (0.40) alone, latency clamps to 0 above the ceiling.
	core.AssertInDelta(t, weightChain, score, 1e-9)
	core.AssertGreaterOrEqual(t, score, 0.0)
	core.AssertLessOrEqual(t, score, 1.0)
}

// TestTrust_DecayScore_Good halves a score after ~69 hours (RFC §11.2).
func TestTrust_DecayScore_Good(t *core.T) {
	decayed := DecayScore(0.8, 69)
	// 0.8 * e^(-0.01*69) ≈ 0.8 * 0.5016 ≈ 0.401.
	core.AssertInDelta(t, 0.401, decayed, 0.01)
}

// TestTrust_DecayScore_Bad leaves the score unchanged for non-positive elapsed
// time.
func TestTrust_DecayScore_Bad(t *core.T) {
	core.AssertInDelta(t, 0.8, DecayScore(0.8, 0), 1e-9)
	core.AssertInDelta(t, 0.8, DecayScore(0.8, -5), 1e-9)
}

// TestTrust_DecayScore_Ugly drops a long-idle score below the stale threshold
// (RFC §11.2: a stale record triggers full re-verification).
func TestTrust_DecayScore_Ugly(t *core.T) {
	decayed := DecayScore(0.9, 1000)
	core.AssertLess(t, decayed, StaleThreshold)
	core.AssertTrue(t, IsStale(decayed))
}

// TestTrust_ApplyGossipPenalty_Good applies the 10% gossip penalty (RFC §12.1).
func TestTrust_ApplyGossipPenalty_Good(t *core.T) {
	core.AssertInDelta(t, 0.72, ApplyGossipPenalty(0.8), 1e-9)
	core.AssertInDelta(t, 0.0, ApplyGossipPenalty(0.0), 1e-9)
}

// TestTrust_ApplyGossipPenalty_Ugly keeps the result within the trust range even
// for an out-of-range input.
func TestTrust_ApplyGossipPenalty_Ugly(t *core.T) {
	core.AssertLessOrEqual(t, ApplyGossipPenalty(2.0), 1.0)
	core.AssertGreaterOrEqual(t, ApplyGossipPenalty(-1.0), 0.0)
}
