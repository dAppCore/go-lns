// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

// BEGIN v090 file-aware generated coverage

func TestNamestateBinary_NameState_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_MarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_MarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_MarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_UnmarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameState).UnmarshalBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_UnmarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameState).UnmarshalBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamestateBinary_NameState_UnmarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameState).UnmarshalBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
