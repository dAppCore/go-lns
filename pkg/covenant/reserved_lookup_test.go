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
