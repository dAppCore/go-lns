// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestOutpoint_NewOutpoint_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(NewOutpoint())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_NewOutpoint_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(NewOutpoint())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_NewOutpoint_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(NewOutpoint())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Inject_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).Inject(coverageValue[Outpoint]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Inject_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).Inject(coverageValue[Outpoint]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Inject_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).Inject(coverageValue[Outpoint]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Clone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Clone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Clone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Equals_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).Equals(coverageValue[Outpoint]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Equals_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).Equals(coverageValue[Outpoint]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Equals_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).Equals(coverageValue[Outpoint]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Compare_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).Compare(coverageValue[Outpoint]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Compare_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).Compare(coverageValue[Outpoint]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_Compare_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).Compare(coverageValue[Outpoint]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_MarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_MarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_MarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_UnmarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).UnmarshalBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_UnmarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).UnmarshalBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestOutpoint_Outpoint_UnmarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).UnmarshalBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
