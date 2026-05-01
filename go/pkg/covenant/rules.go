// SPDX-License-Identifier: EUPL-1.2

package covenant

import "dappco.re/go/lns/pkg/primitives"

// NameFlags mirrors the covenant verification flags used by the JS reference
// rules module.
type NameFlags uint8

// NameRules captures the network name-timing parameters used by the JS
// covenant rules helpers.
//
// The structure keeps the rollout and lockup checks independent from any
// broader network object while preserving the reference behaviour.
type NameRules struct {
	NoRollout         bool
	RolloutInterval   uint32
	AuctionStart      uint32
	ClaimPeriod       uint32
	NoReserved        bool
	AlexaLockupPeriod uint32
}

const (
	// VerifyCovenantsNone disables covenant verification flags.
	VerifyCovenantsNone NameFlags = 0

	// VerifyCovenantsHardened is enabled when the covenant hardening soft fork
	// activates.
	VerifyCovenantsHardened NameFlags = 1 << 0

	// VerifyCovenantsLockup is enabled when the ICANN lockup soft fork
	// activates.
	VerifyCovenantsLockup NameFlags = 1 << 1

	// MandatoryVerifyCovenantFlags is the default covenant verification mask.
	MandatoryVerifyCovenantFlags NameFlags = VerifyCovenantsNone

	// MaxCovenantSize is the maximum serialized covenant size mirrored from the
	// JS rules reference.
	MaxCovenantSize = 1 + 32 + 1 + 4 + 2 + MaxResourceSize + 1 + 32

	// CovenantMaxSize is an alias for MaxCovenantSize.
	CovenantMaxSize = MaxCovenantSize
)

// VerificationFlags mirrors the JS rules.nameFlags lookup table.
//
// The map is exported so callers can discover the named verification masks
// without hard-coding the string keys used by the reference implementation.
var VerificationFlags = map[string]NameFlags{
	"VERIFY_COVENANTS_NONE":     VerifyCovenantsNone,
	"VERIFY_COVENANTS_HARDENED": VerifyCovenantsHardened,
	"VERIFY_COVENANTS_LOCKUP":   VerifyCovenantsLockup,
}

// VerificationFlagsByVal mirrors the JS nameFlags lookup table in reverse.
var VerificationFlagsByVal = map[NameFlags]string{
	VerifyCovenantsNone:     "VERIFY_COVENANTS_NONE",
	VerifyCovenantsHardened: "VERIFY_COVENANTS_HARDENED",
	VerifyCovenantsLockup:   "VERIFY_COVENANTS_LOCKUP",
}

// GetVerificationFlags returns the covenant verification flag lookup table.
func GetVerificationFlags() map[string]NameFlags {
	return VerificationFlags
}

// GetVerificationFlagsByVal returns the reverse covenant verification table.
func GetVerificationFlagsByVal() map[NameFlags]string {
	return VerificationFlagsByVal
}

// GetRollout returns the rollout start height and rollout week for a name hash.
//
// When rollout gating is disabled, it returns zero values so HasRollout reports
// true for any height.
func GetRollout(nameHash primitives.Hash, rules NameRules) (uint32, uint32) {
	if rules.NoRollout {
		return 0, 0
	}

	week := modBuffer(nameHash[:], 52)
	height := week * rules.RolloutInterval
	return rules.AuctionStart + height, week
}

// HasRollout reports whether a name has passed its rollout start height.
func HasRollout(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	start, _ := GetRollout(nameHash, rules)
	return height >= start
}

// GetHasRollout is an alias for HasRollout.
//
//	ok := covenant.GetHasRollout(nameHash, height, rules)
func GetHasRollout(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return HasRollout(nameHash, height, rules)
}

// IsReserved reports whether a name is reserved at the provided height.
//
// Names are only reserved before the claim period and only when the network
// has reserved-name protection enabled.
func IsReserved(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	if rules.NoReserved {
		return false
	}

	if height >= rules.ClaimPeriod {
		return false
	}

	return HasReservedHash(nameHash)
}

// GetIsReserved is an alias for IsReserved.
//
//	ok := covenant.GetIsReserved(nameHash, height, rules)
func GetIsReserved(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return IsReserved(nameHash, height, rules)
}

// IsLockedUp reports whether a name is locked after the claim period.
//
// ICANN root names are always locked; non-root names remain locked until the
// Alexa lockup threshold has elapsed.
func IsLockedUp(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	if rules.NoReserved {
		return false
	}

	if height < rules.ClaimPeriod {
		return false
	}

	item, ok := GetLockedHash(nameHash)
	if !ok {
		return false
	}

	if item.Root {
		return true
	}

	return height < rules.AlexaLockupPeriod
}

// GetIsLockedUp is an alias for IsLockedUp.
//
//	ok := covenant.GetIsLockedUp(nameHash, height, rules)
func GetIsLockedUp(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return IsLockedUp(nameHash, height, rules)
}

func modBuffer(buf []byte, num uint32) uint32 {
	if num == 0 {
		return 0
	}

	p := uint32(256 % num)
	var acc uint32

	for _, b := range buf {
		acc = (p*acc + uint32(b)) % num
	}

	return acc
}
