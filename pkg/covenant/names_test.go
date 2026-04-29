// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestNames_HasNames_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasNames(coverageTransaction("good"), nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_HasNames_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasNames(coverageTransaction("bad"), nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_HasNames_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasNames(coverageTransaction("ugly"), nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_GetHasNames_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasNames(coverageTransaction("good"), nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_GetHasNames_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasNames(coverageTransaction("bad"), nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_GetHasNames_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasNames(coverageTransaction("ugly"), nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_AddNames_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	called := false
	AddNames(coverageTransaction("good"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_AddNames_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	called := false
	AddNames(coverageTransaction("bad"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_AddNames_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	called := false
	AddNames(coverageTransaction("ugly"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_GetAddNames_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	called := false
	GetAddNames(coverageTransaction("good"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_GetAddNames_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	called := false
	GetAddNames(coverageTransaction("bad"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_GetAddNames_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	called := false
	GetAddNames(coverageTransaction("ugly"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_RemoveNames_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	called := false
	RemoveNames(coverageTransaction("good"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_RemoveNames_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	called := false
	RemoveNames(coverageTransaction("bad"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_RemoveNames_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	called := false
	RemoveNames(coverageTransaction("ugly"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_GetRemoveNames_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	called := false
	GetRemoveNames(coverageTransaction("good"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_GetRemoveNames_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	called := false
	GetRemoveNames(coverageTransaction("bad"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_GetRemoveNames_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	called := false
	GetRemoveNames(coverageTransaction("ugly"), nil)
	called = true
	core.AssertTrue(t, called)
}

func TestNames_HasSaneCovenants_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasSaneCovenants(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_HasSaneCovenants_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasSaneCovenants(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_HasSaneCovenants_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasSaneCovenants(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_GetHasSaneCovenants_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasSaneCovenants(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_GetHasSaneCovenants_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasSaneCovenants(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNames_GetHasSaneCovenants_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasSaneCovenants(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
