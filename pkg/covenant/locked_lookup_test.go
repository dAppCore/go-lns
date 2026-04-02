// SPDX-License-Identifier: EUPL-1.2

package covenant

import "testing"

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
}
