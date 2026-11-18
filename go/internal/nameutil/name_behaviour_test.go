// SPDX-License-Identifier: EUPL-1.2

package nameutil

import core "github.com/dappcore/go"

func TestNameBehaviour_CatalogLabel_Good(t *core.T) {
	label, ok := CatalogLabel("lethean.lthn")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "lethean.lthn", label)

	fromBytes, okBytes := CatalogLabel([]byte("lethean.lthn"))
	core.AssertTrue(t, okBytes)
	core.AssertEqual(t, "lethean.lthn", fromBytes)
}

func TestNameBehaviour_CatalogLabel_Bad(t *core.T) {
	label, ok := CatalogLabel(nil)
	core.AssertFalse(t, ok)
	core.AssertEqual(t, "", label)

	notString, okInt := CatalogLabel(42)
	core.AssertFalse(t, okInt)
	core.AssertEqual(t, "", notString)
}

func TestNameBehaviour_CatalogLabel_Ugly(t *core.T) {
	// CatalogLabel preserves bytes exactly, including embedded NUL and uppercase.
	raw := "Lethean.\x00.LTHN"
	label, ok := CatalogLabel(raw)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, raw, label)
}

func TestNameBehaviour_Canonicalize_Good(t *core.T) {
	// Lowercasing plus .lthn suffix trim.
	label, ok := Canonicalize("Lethean.LTHN")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "lethean", label)

	// Trailing-dot fully-qualified .lthn. form.
	fqdn, okFqdn := Canonicalize("lethean.lthn.")
	core.AssertTrue(t, okFqdn)
	core.AssertEqual(t, "lethean", fqdn)

	// Bare trailing dot is stripped.
	trailing, okTrailing := Canonicalize("lethean.")
	core.AssertTrue(t, okTrailing)
	core.AssertEqual(t, "lethean", trailing)

	// Byte input is accepted and lowercased.
	fromBytes, okBytes := Canonicalize([]byte("LETHEAN"))
	core.AssertTrue(t, okBytes)
	core.AssertEqual(t, "lethean", fromBytes)
}

func TestNameBehaviour_Canonicalize_Bad(t *core.T) {
	// Non string/byte input is rejected.
	label, ok := Canonicalize(42)
	core.AssertFalse(t, ok)
	core.AssertEqual(t, "", label)

	// Empty input is rejected.
	empty, okEmpty := Canonicalize("")
	core.AssertFalse(t, okEmpty)
	core.AssertEqual(t, "", empty)

	// A bare suffix that trims to nothing is rejected.
	suffixOnly, okSuffix := Canonicalize(".lthn")
	core.AssertFalse(t, okSuffix)
	core.AssertEqual(t, "", suffixOnly)

	// A lone dot trims to nothing and is rejected.
	dot, okDot := Canonicalize(".")
	core.AssertFalse(t, okDot)
	core.AssertEqual(t, "", dot)
}

func TestNameBehaviour_Canonicalize_Ugly(t *core.T) {
	// Non-ASCII (high bit set) input is rejected.
	label, ok := Canonicalize("léthean")
	core.AssertFalse(t, ok)
	core.AssertEqual(t, "", label)

	// Embedded NUL is below 0x80 so it is preserved, not rejected.
	withNul, okNul := Canonicalize("lethean\x00")
	core.AssertTrue(t, okNul)
	core.AssertEqual(t, "lethean\x00", withNul)
}
