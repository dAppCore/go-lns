// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"crypto/sha3"
	"testing"

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

	if catalog.Size() == 0 {
		t.Fatal("catalog should report a non-zero size")
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
	keys := catalog.Keys()
	values := catalog.Values()

	if len(entries) != len(keys) || len(entries) != len(values) {
		t.Fatalf("catalog iterator lengths differ: entries=%d keys=%d values=%d", len(entries), len(keys), len(values))
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
	hash := primitives.Hash(sha3.Sum256([]byte("foo.bar")))
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
