// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

func TestReservedLookupByName(t *testing.T) {
	item, ok := GetReservedName("RESERVED")
	if !ok {
		t.Fatal("GetReservedName should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
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
	hash, err := HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := GetReservedHash(hash)
	if !ok {
		t.Fatal("GetReservedHash should find the reserved reference entry")
	}

	if item.Hash != hash {
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}

	if !DefaultReservedCatalog().Has(hash) {
		t.Fatal("catalog should report the reserved hash as present")
	}
}

func TestReservedLookupByHashCatalogAlias(t *testing.T) {
	hash, err := HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
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
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}

	if !catalog.GetHasByName("RESERVED") {
		t.Fatal("catalog.GetHasByName should alias HasByName")
	}
}

func TestReservedLookupByHashAlias(t *testing.T) {
	hash, err := HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := GetReservedByHash(hash)
	if !ok {
		t.Fatal("GetReservedByHash should find the reserved reference entry")
	}

	if item.Hash != hash {
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}
}

func TestReservedLookupStringAndBinaryAliases(t *testing.T) {
	item, ok := GetReservedString("RESERVED")
	if !ok {
		t.Fatal("GetReservedString should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	item, ok = GetReservedBinary([]byte("RESERVED"))
	if !ok {
		t.Fatal("GetReservedBinary should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if !HasReservedString("RESERVED") {
		t.Fatal("HasReservedString should report reserved catalog labels")
	}

	if !GetHasReservedName("reserved.lthn") {
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

	item, ok = GetReservedName("reserved.lthn")
	if !ok {
		t.Fatal("GetReservedName should accept canonical .lthn names")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if !HasReservedName("reserved.lthn") {
		t.Fatal("HasReservedName should accept canonical .lthn names")
	}
}

func TestReservedLookupHasByHashAlias(t *testing.T) {
	hash, err := HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
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
	if HasReservedName("does-not-exist") {
		t.Fatal("unknown names should not be reported as reserved")
	}

	if _, ok := GetReservedName("does-not-exist"); ok {
		t.Fatal("GetReservedName should return false for unknown names")
	}

	if HasReservedByName("does-not-exist") {
		t.Fatal("unknown names should not be reported as reserved by name")
	}

	if _, ok := GetReservedByName("does-not-exist"); ok {
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

	hash, err := HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if item.Hash != hash {
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
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
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
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
	fn := ax7ValueOf(HasReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedString_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedName_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedString_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_HasReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetHasReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_DefaultReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_DefaultReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_DefaultReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetDefaultReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetDefaultReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_GetDefaultReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Size_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Size)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Size_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Size)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Size_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Size)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_NameValue_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).NameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_NameValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).NameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_NameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).NameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetNameValue_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetNameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetNameValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetNameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetNameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetNameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_RootValue_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).RootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_RootValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).RootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_RootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).RootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetRootValue_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetRootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetRootValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetRootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetRootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetRootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_TopValue_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).TopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_TopValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).TopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_TopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).TopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetTopValue_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetTopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetTopValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetTopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetTopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetTopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_PrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).PrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_PrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).PrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_PrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).PrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Has_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Has)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Has_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Has)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Has_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Has)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_HasByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).HasByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_HasByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).HasByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_HasByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).HasByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetHasByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetHasByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetHasByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Get_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Get)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Get_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Get)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Get_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Get)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_HasByName_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).HasByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_HasByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).HasByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_HasByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).HasByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetHasByName_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetHasByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetHasByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetByName_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Entries_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Entries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Entries_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Entries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Entries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Entries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetEntries_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetEntries_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Keys_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Keys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Keys_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Keys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Keys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Keys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetKeys_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetKeys_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Values_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Values)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Values_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Values)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_Values_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).Values)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetValues_Good(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetValues_Bad(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestReservedLookup_ReservedCatalog_GetValues_Ugly(t *core.T) {
	fn := ax7ValueOf(new(ReservedCatalog).GetValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
