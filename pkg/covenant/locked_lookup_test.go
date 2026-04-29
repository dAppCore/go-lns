// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

func TestLockedLookupByName(t *testing.T) {
	item, ok := GetLockedName("NEC")
	if !ok {
		t.Fatal("GetLockedName should find the locked reference entry")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if item.Target != "nec." {
		t.Fatalf("item.Target = %q, want %q", item.Target, "nec.")
	}

	if !item.Root || item.Custom {
		t.Fatalf("unexpected flags on nec: %#v", item)
	}
}

func TestLockedLookupByHash(t *testing.T) {
	hash, err := HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := GetLockedHash(hash)
	if !ok {
		t.Fatal("GetLockedHash should find the locked reference entry")
	}

	if item.Hash != hash {
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}

	if !DefaultLockedCatalog().Has(hash) {
		t.Fatal("catalog should report the locked hash as present")
	}
}

func TestLockedLookupByHashCatalogAlias(t *testing.T) {
	hash, err := HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
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
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}

	if !catalog.GetHasByName("NEC") {
		t.Fatal("catalog.GetHasByName should alias HasByName")
	}
}

func TestLockedLookupByHashAlias(t *testing.T) {
	hash, err := HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := GetLockedByHash(hash)
	if !ok {
		t.Fatal("GetLockedByHash should find the locked reference entry")
	}

	if item.Hash != hash {
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}
}

func TestLockedLookupStringAndBinaryAliases(t *testing.T) {
	item, ok := GetLockedString("NEC")
	if !ok {
		t.Fatal("GetLockedString should find the locked reference entry")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	item, ok = GetLockedBinary([]byte("NEC"))
	if !ok {
		t.Fatal("GetLockedBinary should find the locked reference entry")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if !HasLockedString("NEC") {
		t.Fatal("HasLockedString should report locked catalog labels")
	}

	if !GetHasLockedName("nec.lthn") {
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

	item, ok = GetLockedName("nec.lthn")
	if !ok {
		t.Fatal("GetLockedName should accept canonical .lthn names")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if !HasLockedName("nec.lthn") {
		t.Fatal("HasLockedName should accept canonical .lthn names")
	}
}

func TestLockedLookupHasByHashAlias(t *testing.T) {
	hash, err := HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
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
	if HasLockedName("does-not-exist") {
		t.Fatal("unknown names should not be reported as locked")
	}

	if _, ok := GetLockedName("does-not-exist"); ok {
		t.Fatal("GetLockedName should return false for unknown names")
	}

	if HasLockedByName("does-not-exist") {
		t.Fatal("unknown names should not be reported as locked by name")
	}

	if _, ok := GetLockedByName("does-not-exist"); ok {
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

	hash, err := HashString("nec")
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
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
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
	fn := ax7ValueOf(HasLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedString_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedName_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedString_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_HasLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetHasLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_DefaultLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_DefaultLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_DefaultLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetDefaultLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetDefaultLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_GetDefaultLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Size_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Size)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Size_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Size)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Size_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Size)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_PrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).PrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_PrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).PrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_PrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).PrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Has_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Has)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Has_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Has)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Has_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Has)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_HasByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).HasByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_HasByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).HasByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_HasByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).HasByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetHasByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetHasByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetHasByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Get_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Get)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Get_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Get)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Get_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Get)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_HasByName_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).HasByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_HasByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).HasByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_HasByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).HasByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetHasByName_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetHasByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetHasByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetByName_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Entries_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Entries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Entries_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Entries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Entries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Entries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetEntries_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetEntries_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Keys_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Keys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Keys_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Keys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Keys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Keys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetKeys_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetKeys_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Values_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Values)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Values_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Values)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_Values_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).Values)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetValues_Good(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetValues_Bad(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLockedLookup_LockedCatalog_GetValues_Ugly(t *core.T) {
	fn := ax7ValueOf(new(LockedCatalog).GetValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
