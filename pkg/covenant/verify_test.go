// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestVerify_VerifyCovenants_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(VerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestVerify_VerifyCovenants_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(VerifyCovenants(coverageTransaction("bad"), coverageValue[CoinView]("bad"), 0, coverageValue[Network]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestVerify_VerifyCovenants_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(VerifyCovenants(coverageTransaction("ugly"), coverageValue[CoinView]("ugly"), 63, coverageValue[Network]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestVerify_GetVerifyCovenants_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestVerify_GetVerifyCovenants_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerifyCovenants(coverageTransaction("bad"), coverageValue[CoinView]("bad"), 0, coverageValue[Network]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestVerify_GetVerifyCovenants_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerifyCovenants(coverageTransaction("ugly"), coverageValue[CoinView]("ugly"), 63, coverageValue[Network]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
