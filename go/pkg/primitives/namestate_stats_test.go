// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

// BEGIN v090 file-aware generated coverage

func TestNamestateStats_NameState_ToStats_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).ToStats(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateStats_NameState_ToStats_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).ToStats(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateStats_NameState_ToStats_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).ToStats(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateStats_NameState_GetToStats_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetToStats(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateStats_NameState_GetToStats_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetToStats(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateStats_NameState_GetToStats_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetToStats(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
