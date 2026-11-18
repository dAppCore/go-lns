// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

// BEGIN v090 file-aware generated coverage

func TestCovenantJson_Covenant_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantJson_Covenant_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantJson_Covenant_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantJson_Covenant_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).FromJSON(coverageValue[CovenantJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantJson_Covenant_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).FromJSON(coverageValue[CovenantJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantJson_Covenant_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).FromJSON(coverageValue[CovenantJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
