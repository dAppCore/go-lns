// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestOutpointJson_Outpoint_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpointJson_Outpoint_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpointJson_Outpoint_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpointJson_Outpoint_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).FromJSON(coverageValue[OutpointJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpointJson_Outpoint_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).FromJSON(coverageValue[OutpointJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpointJson_Outpoint_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).FromJSON(coverageValue[OutpointJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
