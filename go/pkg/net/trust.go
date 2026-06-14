// SPDX-License-Identifier: EUPL-1.2

package net

import (
	"math"

	poindexter "forge.lthn.ai/Snider/Poindexter"
)

// GossipTrustPenalty is the multiplicative penalty applied to a gossip-sourced
// record's trust score relative to the same record from HSD (RFC §12.1 step 4:
// a 10% penalty, i.e. retain 90%).
const GossipTrustPenalty = 0.9

// DecayLambda is the per-hour exponential decay constant for trust scores
// (RFC §11.2). At λ = 0.01 a score halves after roughly 69 hours without
// re-verification.
const DecayLambda = 0.01

// StaleThreshold is the trust score below which a record is considered stale and
// the next resolution triggers full re-verification (RFC §11.2).
const StaleThreshold = 0.2

// PeerVector holds the observable inputs to a peer trust score (RFC §11.1).
// Each field is a signal Poindexter folds into the composite score.
//
//	v := net.PeerVector{Uptime: 86400, VerifiedChain: true, VerifiedHNS: true}
type PeerVector struct {
	// Uptime is the observed uptime in seconds.
	Uptime float64
	// VerifiedChain reports whether the peer's alias is confirmed on chain.
	VerifiedChain bool
	// VerifiedHNS reports whether the peer's HNS TXT proof is confirmed.
	VerifiedHNS bool
	// ResponseTime is the most recent response latency in milliseconds.
	ResponseTime float64
	// ConsistencyN is the number of times the peer resolved the same address.
	ConsistencyN int
}

// ScorePeer derives a trust score in 0.0–1.0 from a PeerVector. It is the seam
// the LNS uses to plug a trust scorer in; the default implementation delegates
// to Poindexter (RFC §11). A caller may supply an alternative implementation,
// but the LNS ships with Poindexter wired in by default.
type ScorePeer interface {
	// ScorePeer returns the trust score for the supplied vector, clamped 0.0–1.0.
	ScorePeer(v PeerVector) float64
}

// DefaultScorer is the Poindexter-backed ScorePeer implementation. It maps the
// PeerVector signals onto Poindexter Factor weights and folds them with
// poindexter.WeightedScore, then clamps the result to 0.0–1.0.
//
//	score := net.DefaultScorer{}.ScorePeer(net.PeerVector{VerifiedChain: true})
type DefaultScorer struct{}

// Compile-time assertion that DefaultScorer satisfies ScorePeer.
var _ ScorePeer = DefaultScorer{}

// Scoring weights. Verification dominates: a chain-confirmed, HNS-confirmed peer
// floors well above the routing thresholds in RFC §11.3 even before uptime and
// latency contribute.
const (
	weightChain       = 0.40
	weightHNS         = 0.25
	weightUptime      = 0.15
	weightLatency     = 0.10
	weightConsistency = 0.10

	// uptimeFullSeconds is the uptime at which the uptime factor saturates (24h).
	uptimeFullSeconds = 86400.0
	// latencyFloorMs is the latency at or below which the latency factor is 1.0.
	latencyFloorMs = 20.0
	// latencyCeilMs is the latency at or above which the latency factor is 0.0.
	latencyCeilMs = 2000.0
	// consistencyFull is the consistency count at which the factor saturates.
	consistencyFull = 10.0
)

// ScorePeer implements the ScorePeer interface via Poindexter.
func (DefaultScorer) ScorePeer(v PeerVector) float64 {
	factors := []poindexter.Factor{
		{Value: boolFactor(v.VerifiedChain), Weight: weightChain},
		{Value: boolFactor(v.VerifiedHNS), Weight: weightHNS},
		{Value: poindexter.Ratio(clampFloat(v.Uptime, 0, uptimeFullSeconds), uptimeFullSeconds), Weight: weightUptime},
		{Value: latencyFactor(v.ResponseTime), Weight: weightLatency},
		{Value: poindexter.Ratio(clampFloat(float64(v.ConsistencyN), 0, consistencyFull), consistencyFull), Weight: weightConsistency},
	}
	return clampScore(poindexter.WeightedScore(factors))
}

// DecayScore applies the §11.2 exponential decay to a stored score given the
// hours elapsed since its last verification: currentScore = stored * e^(-λ·h).
// Elapsed time at or below zero returns the stored score unchanged.
//
//	current := net.DecayScore(0.8, 69) // ~0.40
func DecayScore(storedScore, hoursElapsed float64) float64 {
	if hoursElapsed <= 0 {
		return clampScore(storedScore)
	}
	return clampScore(storedScore * math.Exp(-DecayLambda*hoursElapsed))
}

// IsStale reports whether a score has decayed below the stale threshold and a
// full re-verification is required (RFC §11.2).
//
//	if net.IsStale(net.DecayScore(stored, hours)) { /* re-verify */ }
func IsStale(score float64) bool {
	return score < StaleThreshold
}

// ApplyGossipPenalty returns the score a gossip-sourced record carries relative
// to its HSD-sourced equivalent (RFC §12.1 step 4).
//
//	penalised := net.ApplyGossipPenalty(0.8) // 0.72
func ApplyGossipPenalty(score float64) float64 {
	return clampScore(score * GossipTrustPenalty)
}

// boolFactor maps a boolean signal to a 0.0/1.0 factor value.
func boolFactor(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// latencyFactor maps a measured latency in milliseconds to a 1.0 (fast) → 0.0
// (slow) factor, linearly between the floor and ceiling. A non-positive value is
// treated as unmeasured and scores 0.0 — an unresponsive or never-probed peer
// must not earn the best latency factor by default.
func latencyFactor(ms float64) float64 {
	if ms <= 0 {
		return 0.0
	}
	if ms <= latencyFloorMs {
		return 1.0
	}
	if ms >= latencyCeilMs {
		return 0.0
	}
	return 1.0 - (ms-latencyFloorMs)/(latencyCeilMs-latencyFloorMs)
}

// clampFloat constrains v to the inclusive [lo, hi] range.
func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// clampScore constrains a score to the canonical 0.0–1.0 trust range.
func clampScore(v float64) float64 {
	return clampFloat(v, 0.0, 1.0)
}
