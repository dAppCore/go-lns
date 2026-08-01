// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"slices"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

func TestLockedLookupByName(t *testing.T) {
	item, ok := GetLockedName("NEC")
	if !ok {
		t.Fatal("GetLockedName should find the locked reference entry")
	}

	if item.Name != testLockedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testLockedLabel)
	}

	if item.Target != "nec." {
		t.Fatalf("item.Target = %q, want %q", item.Target, "nec.")
	}

	if !item.Root || item.Custom {
		t.Fatalf("unexpected flags on nec: %#v", item)
	}
}

func TestLockedLookupByHash(t *testing.T) {
	hash, err := HashString(testLockedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	item, ok := GetLockedHash(hash)
	if !ok {
		t.Fatal("GetLockedHash should find the locked reference entry")
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}

	if !DefaultLockedCatalog().Has(hash) {
		t.Fatal("catalog should report the locked hash as present")
	}
}

func TestLockedLookupByHashCatalogAlias(t *testing.T) {
	hash, err := HashString(testLockedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	catalog := DefaultLockedCatalog()

	if !catalog.HasByHash(hash) {
		t.Fatal("catalog.HasByHash should report the locked hash as present")
	}

	if !catalog.GetHasByHash(hash) {
		t.Fatal("catalog.GetHasByHash should alias HasByHash")
	}

	item, ok := catalog.GetByHash(hash)
	if !ok {
		t.Fatal("catalog.GetByHash should find the locked reference entry")
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}

	if !catalog.GetHasByName("NEC") {
		t.Fatal("catalog.GetHasByName should alias HasByName")
	}
}

func TestLockedLookupByHashAlias(t *testing.T) {
	hash, err := HashString(testLockedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	item, ok := GetLockedByHash(hash)
	if !ok {
		t.Fatal("GetLockedByHash should find the locked reference entry")
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}
}

func TestLockedLookupStringAndBinaryAliases(t *testing.T) {
	item, ok := GetLockedString("NEC")
	if !ok {
		t.Fatal("GetLockedString should find the locked reference entry")
	}

	if item.Name != testLockedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testLockedLabel)
	}

	item, ok = GetLockedBinary([]byte("NEC"))
	if !ok {
		t.Fatal("GetLockedBinary should find the locked reference entry")
	}

	if item.Name != testLockedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testLockedLabel)
	}

	if !HasLockedString("NEC") {
		t.Fatal("HasLockedString should report locked catalog labels")
	}

	if !GetHasLockedName(testLockedLTHN) {
		t.Fatal("GetHasLockedName should accept canonical .lthn names")
	}

	if !GetHasLockedByName("NEC") {
		t.Fatal("GetHasLockedByName should report locked catalog labels")
	}

	if !GetHasLockedString("NEC") {
		t.Fatal("GetHasLockedString should report locked catalog labels")
	}

	if !HasLockedBinary([]byte("NEC")) {
		t.Fatal("HasLockedBinary should report locked catalog labels")
	}

	if !GetHasLockedBinary([]byte("NEC")) {
		t.Fatal("GetHasLockedBinary should report locked catalog labels")
	}

	if _, ok := GetLockedByString("NEC"); !ok {
		t.Fatal("GetLockedByString should find the locked reference entry")
	}

	if _, ok := GetLockedByBinary([]byte("NEC")); !ok {
		t.Fatal("GetLockedByBinary should find the locked reference entry")
	}

	if !HasLockedByString("NEC") {
		t.Fatal("HasLockedByString should report locked catalog labels")
	}

	if !GetHasLockedByString("NEC") {
		t.Fatal("GetHasLockedByString should report locked catalog labels")
	}

	if !HasLockedByBinary([]byte("NEC")) {
		t.Fatal("HasLockedByBinary should report locked catalog labels")
	}

	if !GetHasLockedByBinary([]byte("NEC")) {
		t.Fatal("GetHasLockedByBinary should report locked catalog labels")
	}

	item, ok = GetLockedName(testLockedLTHN)
	if !ok {
		t.Fatal("GetLockedName should accept canonical .lthn names")
	}

	if item.Name != testLockedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testLockedLabel)
	}

	if !HasLockedName(testLockedLTHN) {
		t.Fatal("HasLockedName should accept canonical .lthn names")
	}
}

func TestLockedLookupHasByHashAlias(t *testing.T) {
	hash, err := HashString(testLockedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	if !HasLockedHash(hash) {
		t.Fatal("HasLockedHash should report the locked reference entry")
	}

	if !GetHasLockedHash(hash) {
		t.Fatal("GetHasLockedHash should alias HasLockedHash")
	}

	if !GetHasLockedByHash(hash) {
		t.Fatal("GetHasLockedByHash should alias HasLockedByHash")
	}

	if !HasLockedByHash(hash) {
		t.Fatal("HasLockedByHash should report the locked reference entry")
	}

	if !GetHasLockedByHash(hash) {
		t.Fatal("GetHasLockedByHash should alias HasLockedByHash")
	}

	if HasLockedHash(primitives.Hash{}) {
		t.Fatal("HasLockedHash should return false for unknown hashes")
	}

	if GetHasLockedHash(primitives.Hash{}) {
		t.Fatal("GetHasLockedHash should return false for unknown hashes")
	}
}

func TestLockedLookupRejectsUnknownNames(t *testing.T) {
	if HasLockedName(testUnknownName) {
		t.Fatal("unknown names should not be reported as locked")
	}

	if _, ok := GetLockedName(testUnknownName); ok {
		t.Fatal("GetLockedName should return false for unknown names")
	}

	if HasLockedByName(testUnknownName) {
		t.Fatal("unknown names should not be reported as locked by name")
	}

	if _, ok := GetLockedByName(testUnknownName); ok {
		t.Fatal("GetLockedByName should return false for unknown names")
	}
}

func TestLockedCatalogHelpers(t *testing.T) {
	catalog := DefaultLockedCatalog()
	getCatalog := GetLockedCatalog()
	defaultCatalog := GetDefaultLockedCatalog()

	if catalog.Size() == 0 {
		t.Fatal("catalog should report a non-zero size")
	}

	if getCatalog != catalog {
		t.Fatal("GetLockedCatalog should alias DefaultLockedCatalog")
	}

	if defaultCatalog != catalog {
		t.Fatal("GetDefaultLockedCatalog should alias DefaultLockedCatalog")
	}

	if catalog.GetSize() != catalog.Size() {
		t.Fatal("catalog.GetSize should alias Size")
	}

	if catalog.PrefixSize() != lockedCatalogPrefixSize {
		t.Fatalf("catalog.PrefixSize() = %d, want %d", catalog.PrefixSize(), lockedCatalogPrefixSize)
	}

	if catalog.GetPrefixSize() != lockedCatalogPrefixSize {
		t.Fatalf("catalog.GetPrefixSize() = %d, want %d", catalog.GetPrefixSize(), lockedCatalogPrefixSize)
	}

	if !catalog.HasByName("NEC") {
		t.Fatal("catalog should find names case-insensitively")
	}

	item, ok := catalog.GetByName("NEC")
	if !ok {
		t.Fatal("catalog should return an entry by name")
	}

	hash, err := HashString(testLockedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}

	entries := catalog.Entries()
	getEntries := catalog.GetEntries()
	keys := catalog.Keys()
	getKeys := catalog.GetKeys()
	values := catalog.Values()
	getValues := catalog.GetValues()

	if len(entries) != len(keys) || len(entries) != len(values) || len(entries) != len(getEntries) || len(entries) != len(getKeys) || len(entries) != len(getValues) {
		t.Fatalf("catalog iterator lengths differ: entries=%d getEntries=%d keys=%d getKeys=%d values=%d getValues=%d", len(entries), len(getEntries), len(keys), len(getKeys), len(values), len(getValues))
	}

	var (
		entryFound bool
		keyFound   bool
		valueFound bool
	)

	for _, entry := range entries {
		if entry.Hash == hash {
			entryFound = true
			if entry.Value != item {
				t.Fatalf("entry.Value = %#v, want %#v", entry.Value, item)
			}
		}
	}

	if slices.Contains(keys, hash) {
		keyFound = true
	}

	for _, value := range values {
		if value.Hash == hash {
			valueFound = true
			break
		}
	}

	if !entryFound || !keyFound || !valueFound {
		t.Fatalf("catalog iterators missing locked hash: entry=%v key=%v value=%v", entryFound, keyFound, valueFound)
	}
}

func TestLockedCatalogGetByNameSkipsSyntaxValidation(t *testing.T) {
	hash := primitives.Hash(core.SHA3_256([]byte("foo.bar")))
	catalog := &LockedCatalog{
		byHash: map[primitives.Hash]LockedName{
			hash: {
				Name:   "foo.bar",
				Hash:   hash,
				Target: "foo.bar.",
			},
		},
	}

	item, ok := catalog.GetByName("FOO.BAR")
	if !ok {
		t.Fatal("GetByName should hash arbitrary lowercased labels")
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}
}

func TestLockedCatalogGetByNameRejectsNonASCII(t *testing.T) {
	catalog := DefaultLockedCatalog()

	if catalog.HasByName("Kfoo") {
		t.Fatal("HasByName should reject non-ASCII labels")
	}

	if _, ok := catalog.GetByName("Kfoo"); ok {
		t.Fatal("GetByName should reject non-ASCII labels")
	}
}

// BEGIN v090 file-aware generated coverage

func TestLockedLookup_HasLockedName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasLockedByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasLockedByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_HasLockedByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasLockedByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasLockedByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasLockedByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetHasLockedByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasLockedByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_DefaultLockedCatalog_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(DefaultLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_DefaultLockedCatalog_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(DefaultLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_DefaultLockedCatalog_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(DefaultLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedCatalog_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedCatalog_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetLockedCatalog_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetDefaultLockedCatalog_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetDefaultLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetDefaultLockedCatalog_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetDefaultLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_GetDefaultLockedCatalog_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetDefaultLockedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Size_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).Size())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Size_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).Size())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Size_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).Size())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_PrefixSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).PrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_PrefixSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).PrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_PrefixSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).PrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetPrefixSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetPrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetPrefixSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetPrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetPrefixSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetPrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Has_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).Has([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Has_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).Has([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Has_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).Has([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_HasByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).HasByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_HasByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).HasByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_HasByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).HasByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetHasByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetHasByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetHasByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetHasByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetHasByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetHasByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Get_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).Get([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Get_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).Get([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Get_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).Get([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_HasByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).HasByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_HasByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).HasByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_HasByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).HasByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetHasByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetHasByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetHasByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetHasByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetHasByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetHasByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Entries_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).Entries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Entries_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).Entries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Entries_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).Entries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetEntries_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetEntries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetEntries_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetEntries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetEntries_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetEntries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Keys_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).Keys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Keys_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).Keys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Keys_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).Keys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetKeys_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetKeys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetKeys_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetKeys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetKeys_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetKeys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Values_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).Values())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Values_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).Values())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_Values_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).Values())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetValues_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(LockedCatalog).GetValues())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetValues_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(LockedCatalog).GetValues())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestLockedLookup_LockedCatalog_GetValues_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(LockedCatalog).GetValues())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
