// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"crypto/sha3"
	"testing"

	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
)

func TestServiceResolve(t *testing.T) {
	svc := &Service{}

	got, err := svc.Resolve("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("Resolve returned %x, want %x", got, want)
	}

	got, err = svc.Resolve([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("Resolve([]byte) returned error: %v", err)
	}

	if got != want {
		t.Fatalf("Resolve([]byte) returned %x, want %x", got, want)
	}
}

func TestServiceResolveRejectsInvalidNames(t *testing.T) {
	svc := &Service{}

	cases := []string{
		"",
		" foo-bar.lthn",
		"foo-bar.lthn ",
		"foo.bar.lthn",
		"foo..lthn",
		"foo.lthn..",
		"foo-",
		"test.lthn",
	}

	for _, name := range cases {
		if _, err := svc.Resolve(name); err == nil {
			t.Fatalf("Resolve(%q) should reject the name", name)
		}
	}

	if _, err := svc.Resolve(123); err == nil {
		t.Fatal("Resolve should reject unsupported input types")
	}
}

func TestServiceVerify(t *testing.T) {
	svc := &Service{}

	cases := []struct {
		name any
		ok   bool
	}{
		{name: "Foo-Bar.lthn", ok: true},
		{name: []byte("Foo-Bar.lthn"), ok: true},
		{name: "foo.bar.lthn", ok: false},
		{name: "foo-", ok: false},
		{name: 123, ok: false},
	}

	for _, tc := range cases {
		if got := svc.Verify(tc.name); got != tc.ok {
			t.Fatalf("Verify(%v) = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestServiceResolveNameAliases(t *testing.T) {
	svc := &Service{}

	got, err := svc.ResolveName("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveName returned error: %v", err)
	}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))
	if got != want {
		t.Fatalf("ResolveName returned %x, want %x", got, want)
	}

	if _, err := svc.ResolveName(123); err == nil {
		t.Fatal("ResolveName should reject unsupported input types")
	}
}

func TestServiceVerifyNameAliases(t *testing.T) {
	svc := &Service{}

	if !svc.VerifyName("Foo-Bar.lthn") {
		t.Fatal("VerifyName should accept canonical names")
	}

	if svc.VerifyName("foo.bar.lthn") {
		t.Fatal("VerifyName should reject malformed names")
	}

	if svc.VerifyName(123) {
		t.Fatal("VerifyName should reject unsupported input types")
	}
}

func TestServiceGetReserved(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetReserved("RESERVED.lthn")
	if !ok {
		t.Fatal("GetReserved should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if item.Target != "reserved.com." {
		t.Fatalf("item.Target = %q, want %q", item.Target, "reserved.com.")
	}

	hash, err := covenant.HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if item.Hash != hash {
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}

	if _, ok := svc.GetReserved([]byte("reserved.lthn")); !ok {
		t.Fatal("GetReserved([]byte) should find the reserved reference entry")
	}

	if _, ok := svc.GetReserved(123); ok {
		t.Fatal("GetReserved should reject unsupported input types")
	}
}

func TestServiceHasReserved(t *testing.T) {
	svc := &Service{}

	if !svc.HasReserved("RESERVED.lthn") {
		t.Fatal("HasReserved should report canonical reserved names")
	}

	if !svc.HasReserved([]byte("reserved.lthn")) {
		t.Fatal("HasReserved([]byte) should report canonical reserved names")
	}

	if svc.HasReserved("not-reserved.lthn") {
		t.Fatal("HasReserved should return false for unknown names")
	}

	if svc.HasReserved(123) {
		t.Fatal("HasReserved should reject unsupported input types")
	}
}

func TestServiceGetReservedByName(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetReservedByName("RESERVED")
	if !ok {
		t.Fatal("GetReservedByName should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if _, ok := svc.GetReservedByName([]byte("reserved")); !ok {
		t.Fatal("GetReservedByName([]byte) should find the reserved reference entry")
	}

	if _, ok := svc.GetReservedByName("RESERVED.lthn"); !ok {
		t.Fatal("GetReservedByName should also accept canonical reserved names")
	}

	if _, ok := svc.GetReservedByName(123); ok {
		t.Fatal("GetReservedByName should reject unsupported input types")
	}
}

func TestServiceHasReservedByName(t *testing.T) {
	svc := &Service{}

	if !svc.HasReservedByName("RESERVED") {
		t.Fatal("HasReservedByName should report reserved catalog labels")
	}

	if !svc.HasReservedByName([]byte("reserved")) {
		t.Fatal("HasReservedByName([]byte) should report reserved catalog labels")
	}

	if !svc.HasReservedByName("reserved.lthn") {
		t.Fatal("HasReservedByName should report canonical reserved names")
	}

	if svc.HasReservedByName("not-reserved") {
		t.Fatal("HasReservedByName should return false for unknown labels")
	}

	if svc.HasReservedByName(123) {
		t.Fatal("HasReservedByName should reject unsupported input types")
	}
}

func TestServiceReservedNameAliases(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetReservedName("RESERVED")
	if !ok {
		t.Fatal("GetReservedName should find the reserved catalog label")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	item, ok = svc.GetReservedName("RESERVED.lthn")
	if !ok {
		t.Fatal("GetReservedName should also accept canonical reserved names")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if !svc.HasReservedName("reserved") {
		t.Fatal("HasReservedName should report reserved catalog labels")
	}

	if !svc.HasReservedName("reserved.lthn") {
		t.Fatal("HasReservedName should report canonical reserved names")
	}
}

func TestServiceGetReservedHash(t *testing.T) {
	svc := &Service{}

	hash, err := covenant.HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := svc.GetReservedHash(hash)
	if !ok {
		t.Fatal("GetReservedHash should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if !svc.HasReservedHash(hash) {
		t.Fatal("HasReservedHash should report canonical reserved hashes")
	}

	if svc.HasReservedHash(primitives.Hash{}) {
		t.Fatal("HasReservedHash should return false for unknown hashes")
	}
}

func TestServiceGetReservedByHash(t *testing.T) {
	svc := &Service{}

	hash, err := covenant.HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := svc.GetReservedByHash(hash)
	if !ok {
		t.Fatal("GetReservedByHash should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if !svc.HasReservedByHash(hash) {
		t.Fatal("HasReservedByHash should report canonical reserved hashes")
	}

	if svc.HasReservedByHash(primitives.Hash{}) {
		t.Fatal("HasReservedByHash should return false for unknown hashes")
	}
}

func TestServiceGetLocked(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetLocked("NEC.lthn")
	if !ok {
		t.Fatal("GetLocked should find the locked reference entry")
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

	if _, ok := svc.GetLocked([]byte("nec.lthn")); !ok {
		t.Fatal("GetLocked([]byte) should find the locked reference entry")
	}

	if _, ok := svc.GetLocked(123); ok {
		t.Fatal("GetLocked should reject unsupported input types")
	}
}

func TestServiceHasLocked(t *testing.T) {
	svc := &Service{}

	if !svc.HasLocked("NEC.lthn") {
		t.Fatal("HasLocked should report canonical locked names")
	}

	if !svc.HasLocked([]byte("nec.lthn")) {
		t.Fatal("HasLocked([]byte) should report canonical locked names")
	}

	if svc.HasLocked("not-locked.lthn") {
		t.Fatal("HasLocked should return false for unknown names")
	}

	if svc.HasLocked(123) {
		t.Fatal("HasLocked should reject unsupported input types")
	}
}

func TestServiceGetLockedByName(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetLockedByName("NEC")
	if !ok {
		t.Fatal("GetLockedByName should find the locked reference entry")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if _, ok := svc.GetLockedByName([]byte("nec")); !ok {
		t.Fatal("GetLockedByName([]byte) should find the locked reference entry")
	}

	if _, ok := svc.GetLockedByName("NEC.lthn"); !ok {
		t.Fatal("GetLockedByName should also accept canonical locked names")
	}

	if _, ok := svc.GetLockedByName(123); ok {
		t.Fatal("GetLockedByName should reject unsupported input types")
	}
}

func TestServiceHasLockedByName(t *testing.T) {
	svc := &Service{}

	if !svc.HasLockedByName("NEC") {
		t.Fatal("HasLockedByName should report locked catalog labels")
	}

	if !svc.HasLockedByName([]byte("nec")) {
		t.Fatal("HasLockedByName([]byte) should report locked catalog labels")
	}

	if !svc.HasLockedByName("nec.lthn") {
		t.Fatal("HasLockedByName should report canonical locked names")
	}

	if svc.HasLockedByName("not-locked") {
		t.Fatal("HasLockedByName should return false for unknown labels")
	}

	if svc.HasLockedByName(123) {
		t.Fatal("HasLockedByName should reject unsupported input types")
	}
}

func TestServiceLockedNameAliases(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetLockedName("NEC")
	if !ok {
		t.Fatal("GetLockedName should find the locked catalog label")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	item, ok = svc.GetLockedName("NEC.lthn")
	if !ok {
		t.Fatal("GetLockedName should also accept canonical locked names")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if !svc.HasLockedName("nec") {
		t.Fatal("HasLockedName should report locked catalog labels")
	}

	if !svc.HasLockedName("nec.lthn") {
		t.Fatal("HasLockedName should report canonical locked names")
	}
}

func TestServiceGetLockedHash(t *testing.T) {
	svc := &Service{}

	hash, err := covenant.HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := svc.GetLockedHash(hash)
	if !ok {
		t.Fatal("GetLockedHash should find the locked reference entry")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if !svc.HasLockedHash(hash) {
		t.Fatal("HasLockedHash should report canonical locked hashes")
	}

	if svc.HasLockedHash(primitives.Hash{}) {
		t.Fatal("HasLockedHash should return false for unknown hashes")
	}
}

func TestServiceGetLockedByHash(t *testing.T) {
	svc := &Service{}

	hash, err := covenant.HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	item, ok := svc.GetLockedByHash(hash)
	if !ok {
		t.Fatal("GetLockedByHash should find the locked reference entry")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if !svc.HasLockedByHash(hash) {
		t.Fatal("HasLockedByHash should report canonical locked hashes")
	}

	if svc.HasLockedByHash(primitives.Hash{}) {
		t.Fatal("HasLockedByHash should return false for unknown hashes")
	}
}

func TestServiceReservedCatalog(t *testing.T) {
	svc := &Service{}
	catalog := svc.ReservedCatalog()

	if catalog == nil {
		t.Fatal("ReservedCatalog should return a catalog")
	}

	if got := svc.ReservedSize(); got != catalog.Size() {
		t.Fatalf("ReservedSize() = %d, want %d", got, catalog.Size())
	}

	if got := svc.ReservedNameValue(); got != catalog.NameValue() {
		t.Fatalf("ReservedNameValue() = %d, want %d", got, catalog.NameValue())
	}

	if got := svc.ReservedRootValue(); got != catalog.RootValue() {
		t.Fatalf("ReservedRootValue() = %d, want %d", got, catalog.RootValue())
	}

	if got := svc.ReservedTopValue(); got != catalog.TopValue() {
		t.Fatalf("ReservedTopValue() = %d, want %d", got, catalog.TopValue())
	}

	if got := svc.ReservedPrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("ReservedPrefixSize() = %d, want %d", got, catalog.PrefixSize())
	}

	if !catalog.HasByName("RESERVED") {
		t.Fatal("ReservedCatalog should expose the reserved catalog helpers")
	}

	if len(catalog.Keys()) == 0 {
		t.Fatal("ReservedCatalog should expose hash lookup aliases")
	}

	if !catalog.HasByHash(catalog.Keys()[0]) {
		t.Fatal("ReservedCatalog should expose hash lookup aliases")
	}

	if _, ok := catalog.GetByHash(catalog.Keys()[0]); !ok {
		t.Fatal("ReservedCatalog should expose hash lookup aliases")
	}

	if len(catalog.Entries()) == 0 {
		t.Fatal("ReservedCatalog should expose catalog entries")
	}

	if len(svc.ReservedEntries()) == 0 {
		t.Fatal("ReservedEntries should expose catalog entries")
	}

	if len(svc.ReservedKeys()) == 0 {
		t.Fatal("ReservedKeys should expose catalog keys")
	}

	if len(svc.ReservedValues()) == 0 {
		t.Fatal("ReservedValues should expose catalog values")
	}
}

func TestServiceLockedCatalog(t *testing.T) {
	svc := &Service{}
	catalog := svc.LockedCatalog()

	if catalog == nil {
		t.Fatal("LockedCatalog should return a catalog")
	}

	if got := svc.LockedSize(); got != catalog.Size() {
		t.Fatalf("LockedSize() = %d, want %d", got, catalog.Size())
	}

	if got := svc.LockedPrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("LockedPrefixSize() = %d, want %d", got, catalog.PrefixSize())
	}

	if !catalog.HasByName("NEC") {
		t.Fatal("LockedCatalog should expose the locked catalog helpers")
	}

	if len(catalog.Keys()) == 0 {
		t.Fatal("LockedCatalog should expose hash lookup aliases")
	}

	if !catalog.HasByHash(catalog.Keys()[0]) {
		t.Fatal("LockedCatalog should expose hash lookup aliases")
	}

	if _, ok := catalog.GetByHash(catalog.Keys()[0]); !ok {
		t.Fatal("LockedCatalog should expose hash lookup aliases")
	}

	if len(catalog.Entries()) == 0 {
		t.Fatal("LockedCatalog should expose catalog entries")
	}

	if len(svc.LockedEntries()) == 0 {
		t.Fatal("LockedEntries should expose catalog entries")
	}

	if len(svc.LockedKeys()) == 0 {
		t.Fatal("LockedKeys should expose catalog keys")
	}

	if len(svc.LockedValues()) == 0 {
		t.Fatal("LockedValues should expose catalog values")
	}
}
