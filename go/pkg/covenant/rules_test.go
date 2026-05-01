// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestRules_GetVerificationFlags_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerificationFlags())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetVerificationFlags_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerificationFlags())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetVerificationFlags_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerificationFlags())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetVerificationFlagsByVal_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerificationFlagsByVal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetVerificationFlagsByVal_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerificationFlagsByVal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetVerificationFlagsByVal_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerificationFlagsByVal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetRollout_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetRollout([32]byte{1}, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetRollout_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetRollout([32]byte{}, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetRollout_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetRollout([32]byte{0xff}, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_HasRollout_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasRollout([32]byte{1}, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_HasRollout_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasRollout([32]byte{}, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_HasRollout_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasRollout([32]byte{0xff}, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetHasRollout_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasRollout([32]byte{1}, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetHasRollout_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasRollout([32]byte{}, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetHasRollout_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasRollout([32]byte{0xff}, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_IsReserved_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(IsReserved([32]byte{1}, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_IsReserved_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(IsReserved([32]byte{}, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_IsReserved_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(IsReserved([32]byte{0xff}, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetIsReserved_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetIsReserved([32]byte{1}, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetIsReserved_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetIsReserved([32]byte{}, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetIsReserved_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetIsReserved([32]byte{0xff}, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_IsLockedUp_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(IsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_IsLockedUp_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(IsLockedUp([32]byte{}, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_IsLockedUp_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(IsLockedUp([32]byte{0xff}, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetIsLockedUp_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetIsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetIsLockedUp_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetIsLockedUp([32]byte{}, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRules_GetIsLockedUp_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetIsLockedUp([32]byte{0xff}, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
