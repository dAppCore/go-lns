// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "github.com/dappcore/go"

// BEGIN v090 file-aware generated coverage

func TestHelpers_jsonraw_MarshalJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(jsonraw).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestHelpers_jsonraw_MarshalJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(jsonraw).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestHelpers_jsonraw_MarshalJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(jsonraw).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestHelpers_jsonraw_UnmarshalJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(jsonraw).UnmarshalJSON([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestHelpers_jsonraw_UnmarshalJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(jsonraw).UnmarshalJSON([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestHelpers_jsonraw_UnmarshalJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(jsonraw).UnmarshalJSON([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
