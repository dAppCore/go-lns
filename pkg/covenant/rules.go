// SPDX-License-Identifier: EUPL-1.2

package covenant

// NameFlags mirrors the covenant verification flags used by the JS reference
// rules module.
type NameFlags uint8

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
