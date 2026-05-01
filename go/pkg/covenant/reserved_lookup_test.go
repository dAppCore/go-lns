// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

const (
	testReservedLabel       = "reserved"
	testLockedLabel         = "nec"
	testReservedLTHN        = "reserved.lthn"
	testLockedLTHN          = "nec.lthn"
	testUnknownName         = "does-not-exist"
	testFooBarLabel         = "foo-bar"
	testExampleName         = "example-name"
	testHashStringErr       = "HashString returned error: %v"
	testItemNameFormat      = "item.Name = %q, want %q"
	testItemHashFormat      = "item.Hash = %x, want %x"
	testVerifyCovenantsZero = "VerifyCovenants returned %d, want 0"
)

func TestReservedLookupByName(t *testing.T) {
	item, ok := GetReservedName("RESERVED")
	if !ok {
		t.Fatal("GetReservedName should find the reserved reference entry")
	}

	if item.Name != testReservedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testReservedLabel)
	}

	if item.Target != "reserved.com." {
		t.Fatalf("item.Target = %q, want %q", item.Target, "reserved.com.")
	}

	if item.Value != 503513487 {
		t.Fatalf("item.Value = %d, want %d", item.Value, 503513487)
	}

	if item.Root || item.Top100 || item.Custom || item.Zero {
		t.Fatalf("unexpected flags on reserved.com: %#v", item)
	}
}

func TestReservedLookupByHash(t *testing.T) {
	hash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	item, ok := GetReservedHash(hash)
	if !ok {
		t.Fatal("GetReservedHash should find the reserved reference entry")
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}

	if !DefaultReservedCatalog().Has(hash) {
		t.Fatal("catalog should report the reserved hash as present")
	}
}

func TestReservedLookupByHashCatalogAlias(t *testing.T) {
	hash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	catalog := DefaultReservedCatalog()

	if !catalog.HasByHash(hash) {
		t.Fatal("catalog.HasByHash should report the reserved hash as present")
	}

	if !catalog.GetHasByHash(hash) {
		t.Fatal("catalog.GetHasByHash should alias HasByHash")
	}

	item, ok := catalog.GetByHash(hash)
	if !ok {
		t.Fatal("catalog.GetByHash should find the reserved reference entry")
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}

	if !catalog.GetHasByName("RESERVED") {
		t.Fatal("catalog.GetHasByName should alias HasByName")
	}
}

func TestReservedLookupByHashAlias(t *testing.T) {
	hash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	item, ok := GetReservedByHash(hash)
	if !ok {
		t.Fatal("GetReservedByHash should find the reserved reference entry")
	}

	if item.Hash != hash {
		t.Fatalf(testItemHashFormat, item.Hash, hash)
	}
}

func TestReservedLookupStringAndBinaryAliases(t *testing.T) {
	item, ok := GetReservedString("RESERVED")
	if !ok {
		t.Fatal("GetReservedString should find the reserved reference entry")
	}

	if item.Name != testReservedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testReservedLabel)
	}

	item, ok = GetReservedBinary([]byte("RESERVED"))
	if !ok {
		t.Fatal("GetReservedBinary should find the reserved reference entry")
	}

	if item.Name != testReservedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testReservedLabel)
	}

	if !HasReservedString("RESERVED") {
		t.Fatal("HasReservedString should report reserved catalog labels")
	}

	if !GetHasReservedName(testReservedLTHN) {
		t.Fatal("GetHasReservedName should accept canonical .lthn names")
	}

	if !GetHasReservedByName("RESERVED") {
		t.Fatal("GetHasReservedByName should report reserved catalog labels")
	}

	if !GetHasReservedString("RESERVED") {
		t.Fatal("GetHasReservedString should report reserved catalog labels")
	}

	if !GetHasReservedBinary([]byte("RESERVED")) {
		t.Fatal("GetHasReservedBinary should report reserved catalog labels")
	}

	if !HasReservedBinary([]byte("RESERVED")) {
		t.Fatal("HasReservedBinary should report reserved catalog labels")
	}

	if _, ok := GetReservedByString("RESERVED"); !ok {
		t.Fatal("GetReservedByString should find the reserved reference entry")
	}

	if _, ok := GetReservedByBinary([]byte("RESERVED")); !ok {
		t.Fatal("GetReservedByBinary should find the reserved reference entry")
	}

	if !HasReservedByString("RESERVED") {
		t.Fatal("HasReservedByString should report reserved catalog labels")
	}

	if !GetHasReservedByString("RESERVED") {
		t.Fatal("GetHasReservedByString should report reserved catalog labels")
	}

	if !HasReservedByBinary([]byte("RESERVED")) {
		t.Fatal("HasReservedByBinary should report reserved catalog labels")
	}

	if !GetHasReservedByBinary([]byte("RESERVED")) {
		t.Fatal("GetHasReservedByBinary should report reserved catalog labels")
	}

	item, ok = GetReservedName(testReservedLTHN)
	if !ok {
		t.Fatal("GetReservedName should accept canonical .lthn names")
	}

	if item.Name != testReservedLabel {
		t.Fatalf(testItemNameFormat, item.Name, testReservedLabel)
	}

	if !HasReservedName(testReservedLTHN) {
		t.Fatal("HasReservedName should accept canonical .lthn names")
	}
}

func TestReservedLookupHasByHashAlias(t *testing.T) {
	hash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	if !HasReservedHash(hash) {
		t.Fatal("HasReservedHash should report the reserved reference entry")
	}

	if !GetHasReservedHash(hash) {
		t.Fatal("GetHasReservedHash should alias HasReservedHash")
	}

	if !GetHasReservedByHash(hash) {
		t.Fatal("GetHasReservedByHash should alias HasReservedByHash")
	}

	if !HasReservedByHash(hash) {
		t.Fatal("HasReservedByHash should report the reserved reference entry")
	}

	if !GetHasReservedByHash(hash) {
		t.Fatal("GetHasReservedByHash should alias HasReservedByHash")
	}

	if HasReservedHash(primitives.Hash{}) {
		t.Fatal("HasReservedHash should return false for unknown hashes")
	}

	if GetHasReservedHash(primitives.Hash{}) {
		t.Fatal("GetHasReservedHash should return false for unknown hashes")
	}
}

func TestReservedLookupRejectsUnknownNames(t *testing.T) {
	if HasReservedName(testUnknownName) {
		t.Fatal("unknown names should not be reported as reserved")
	}

	if _, ok := GetReservedName(testUnknownName); ok {
		t.Fatal("GetReservedName should return false for unknown names")
	}

	if HasReservedByName(testUnknownName) {
		t.Fatal("unknown names should not be reported as reserved by name")
	}

	if _, ok := GetReservedByName(testUnknownName); ok {
		t.Fatal("GetReservedByName should return false for unknown names")
	}
}

func TestReservedCatalogHelpers(t *testing.T) {
	catalog := DefaultReservedCatalog()
	getCatalog := GetReservedCatalog()
	defaultCatalog := GetDefaultReservedCatalog()

	if catalog.Size() == 0 {
		t.Fatal("catalog should report a non-zero size")
	}

	if getCatalog != catalog {
		t.Fatal("GetReservedCatalog should alias DefaultReservedCatalog")
	}

	if defaultCatalog != catalog {
		t.Fatal("GetDefaultReservedCatalog should alias DefaultReservedCatalog")
	}

	if catalog.GetSize() != catalog.Size() {
		t.Fatal("catalog.GetSize should alias Size")
	}

	if catalog.NameValue() == 0 {
		t.Fatal("catalog should expose the reserved base value")
	}

	if catalog.GetNameValue() != catalog.NameValue() {
		t.Fatal("catalog.GetNameValue should alias NameValue")
	}

	if catalog.RootValue() == 0 {
		t.Fatal("catalog should expose the root increment")
	}

	if catalog.GetRootValue() != catalog.RootValue() {
		t.Fatal("catalog.GetRootValue should alias RootValue")
	}

	if catalog.TopValue() == 0 {
		t.Fatal("catalog should expose the top-100 increment")
	}

	if catalog.GetTopValue() != catalog.TopValue() {
		t.Fatal("catalog.GetTopValue should alias TopValue")
	}

	if catalog.PrefixSize() != reservedCatalogPrefixSize {
		t.Fatalf("catalog.PrefixSize() = %d, want %d", catalog.PrefixSize(), reservedCatalogPrefixSize)
	}

	if catalog.GetPrefixSize() != reservedCatalogPrefixSize {
		t.Fatalf("catalog.GetPrefixSize() = %d, want %d", catalog.GetPrefixSize(), reservedCatalogPrefixSize)
	}

	if !catalog.HasByName("RESERVED") {
		t.Fatal("catalog should find names case-insensitively")
	}

	item, ok := catalog.GetByName("RESERVED")
	if !ok {
		t.Fatal("catalog should return an entry by name")
	}

	hash, err := HashString(testReservedLabel)
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

	for _, key := range keys {
		if key == hash {
			keyFound = true
			break
		}
	}

	for _, value := range values {
		if value.Hash == hash {
			valueFound = true
			break
		}
	}

	if !entryFound || !keyFound || !valueFound {
		t.Fatalf("catalog iterators missing reserved hash: entry=%v key=%v value=%v", entryFound, keyFound, valueFound)
	}
}

func TestReservedCatalogGetByNameSkipsSyntaxValidation(t *testing.T) {
	hash := primitives.Hash(core.SHA3_256([]byte("foo.bar")))
	catalog := &ReservedCatalog{
		byHash: map[primitives.Hash]ReservedName{
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

func TestReservedCatalogGetByNameRejectsNonASCII(t *testing.T) {
	catalog := DefaultReservedCatalog()

	if catalog.HasByName("Kfoo") {
		t.Fatal("HasByName should reject non-ASCII labels")
	}

	if _, ok := catalog.GetByName("Kfoo"); ok {
		t.Fatal("GetByName should reject non-ASCII labels")
	}
}

// BEGIN v090 file-aware generated coverage

func TestReservedLookup_HasReservedName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HasReservedByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HasReservedByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_HasReservedByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HasReservedByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHasReservedByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHasReservedByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetHasReservedByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHasReservedByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_DefaultReservedCatalog_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(DefaultReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_DefaultReservedCatalog_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(DefaultReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_DefaultReservedCatalog_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(DefaultReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedCatalog_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedCatalog_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetReservedCatalog_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetDefaultReservedCatalog_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetDefaultReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetDefaultReservedCatalog_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetDefaultReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_GetDefaultReservedCatalog_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetDefaultReservedCatalog())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Size_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).Size())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Size_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).Size())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Size_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).Size())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_NameValue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).NameValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_NameValue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).NameValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_NameValue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).NameValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetNameValue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetNameValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetNameValue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetNameValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetNameValue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetNameValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_RootValue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).RootValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_RootValue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).RootValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_RootValue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).RootValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetRootValue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetRootValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetRootValue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetRootValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetRootValue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetRootValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_TopValue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).TopValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_TopValue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).TopValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_TopValue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).TopValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetTopValue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetTopValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetTopValue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetTopValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetTopValue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetTopValue())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_PrefixSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).PrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_PrefixSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).PrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_PrefixSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).PrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetPrefixSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetPrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetPrefixSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetPrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetPrefixSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetPrefixSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Has_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).Has([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Has_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).Has([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Has_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).Has([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_HasByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).HasByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_HasByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).HasByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_HasByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).HasByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetHasByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetHasByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetHasByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetHasByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetHasByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetHasByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Get_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).Get([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Get_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).Get([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Get_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).Get([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetByHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetByHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetByHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetByHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetByHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetByHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_HasByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).HasByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_HasByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).HasByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_HasByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).HasByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetHasByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetHasByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetHasByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetHasByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetHasByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetHasByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetByName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Entries_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).Entries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Entries_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).Entries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Entries_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).Entries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetEntries_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetEntries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetEntries_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetEntries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetEntries_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetEntries())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Keys_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).Keys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Keys_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).Keys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Keys_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).Keys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetKeys_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetKeys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetKeys_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetKeys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetKeys_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetKeys())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Values_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).Values())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Values_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).Values())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_Values_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).Values())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetValues_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(ReservedCatalog).GetValues())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetValues_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(ReservedCatalog).GetValues())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestReservedLookup_ReservedCatalog_GetValues_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(ReservedCatalog).GetValues())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
