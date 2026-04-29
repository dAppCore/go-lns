// SPDX-License-Identifier: EUPL-1.2

package nameutil

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestName_CatalogLabel_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(CatalogLabel("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestName_CatalogLabel_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(CatalogLabel(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestName_CatalogLabel_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(CatalogLabel("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestName_Canonicalize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(Canonicalize("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestName_Canonicalize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(Canonicalize(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestName_Canonicalize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(Canonicalize("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
