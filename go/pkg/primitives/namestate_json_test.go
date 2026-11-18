// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

// BEGIN v090 file-aware generated coverage

func TestNamestateJson_NameState_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetJSON(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetJSON(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetJSON(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_Format_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).Format(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_Format_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).Format(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_Format_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).Format(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetFormat_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetFormat(1, coverageValue[NameStateRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetFormat_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetFormat(0, coverageValue[NameStateRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetFormat_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetFormat(63, coverageValue[NameStateRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).FromJSON(coverageValue[NameStateJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).FromJSON(coverageValue[NameStateJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).FromJSON(coverageValue[NameStateJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetFromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetFromJSON(coverageValue[NameStateJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetFromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetFromJSON(coverageValue[NameStateJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateJson_NameState_GetFromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetFromJSON(coverageValue[NameStateJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
