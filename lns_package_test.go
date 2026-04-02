// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"crypto/sha3"
	"encoding/binary"
	"testing"

	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
)

func TestPackageResolveAndVerifyAliases(t *testing.T) {
	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))

	resolveCases := []struct {
		name string
		fn   func(any) (primitives.Hash, error)
	}{
		{name: "Resolve", fn: Resolve},
		{name: "Hash", fn: Hash},
		{name: "ResolveName", fn: ResolveName},
		{name: "HashName", fn: HashName},
		{name: "ResolveByName", fn: ResolveByName},
		{name: "HashByName", fn: HashByName},
	}

	for _, tc := range resolveCases {
		got, err := tc.fn("Foo-Bar.lthn")
		if err != nil {
			t.Fatalf("%s returned error: %v", tc.name, err)
		}

		if got != want {
			t.Fatalf("%s returned %x, want %x", tc.name, got, want)
		}
	}

	if got, err := ResolveString("Foo-Bar.lthn"); err != nil {
		t.Fatalf("ResolveString returned error: %v", err)
	} else if got != want {
		t.Fatalf("ResolveString returned %x, want %x", got, want)
	}

	if got, err := ResolveBinary([]byte("Foo-Bar.lthn")); err != nil {
		t.Fatalf("ResolveBinary returned error: %v", err)
	} else if got != want {
		t.Fatalf("ResolveBinary returned %x, want %x", got, want)
	}

	if got, err := HashString("Foo-Bar.lthn"); err != nil {
		t.Fatalf("HashString returned error: %v", err)
	} else if got != want {
		t.Fatalf("HashString returned %x, want %x", got, want)
	}

	if got, err := HashBinary([]byte("Foo-Bar.lthn")); err != nil {
		t.Fatalf("HashBinary returned error: %v", err)
	} else if got != want {
		t.Fatalf("HashBinary returned %x, want %x", got, want)
	}

	verifyCases := []struct {
		name string
		fn   func(any) bool
	}{
		{name: "Verify", fn: Verify},
		{name: "VerifyName", fn: VerifyName},
		{name: "VerifyByName", fn: VerifyByName},
	}

	for _, tc := range verifyCases {
		if !tc.fn("Foo-Bar.lthn") {
			t.Fatalf("%s should accept canonical names", tc.name)
		}

		if tc.fn("foo.bar.lthn") {
			t.Fatalf("%s should reject malformed names", tc.name)
		}
	}

	if !VerifyString("Foo-Bar.lthn") {
		t.Fatal("VerifyString should accept canonical names")
	}

	if VerifyString("foo.bar.lthn") {
		t.Fatal("VerifyString should reject malformed names")
	}

	if !VerifyBinary([]byte("Foo-Bar.lthn")) {
		t.Fatal("VerifyBinary should accept canonical names")
	}

	if VerifyBinary([]byte("foo.bar.lthn")) {
		t.Fatal("VerifyBinary should reject malformed names")
	}

	if !VerifyByString("Foo-Bar.lthn") {
		t.Fatal("VerifyByString should accept canonical names")
	}

	if !VerifyByBinary([]byte("Foo-Bar.lthn")) {
		t.Fatal("VerifyByBinary should accept canonical names")
	}
}

func TestPackageBlindAndTypeAliases(t *testing.T) {
	var nonce primitives.Hash
	for i := range nonce {
		nonce[i] = byte(i)
	}

	got, err := Blind(0x1122334455667788, nonce)
	if err != nil {
		t.Fatalf("Blind returned error: %v", err)
	}

	var input [40]byte
	binary.LittleEndian.PutUint64(input[:8], 0x1122334455667788)
	copy(input[8:], nonce[:])

	want := blake2b.Sum256(input[:])
	if got != want {
		t.Fatalf("Blind returned %x, want %x", got, want)
	}

	got, err = GetBlind(0x1122334455667788, nonce)
	if err != nil {
		t.Fatalf("GetBlind returned error: %v", err)
	}

	if got != want {
		t.Fatalf("GetBlind returned %x, want %x", got, want)
	}

	if got := TypeName(covenant.TypeBid); got != "BID" {
		t.Fatalf("TypeName(TypeBid) = %q, want %q", got, "BID")
	}

	if got := GetTypeName(covenant.CovenantType(99)); got != "UNKNOWN" {
		t.Fatalf("GetTypeName(99) = %q, want %q", got, "UNKNOWN")
	}
}

func TestPackageServiceNameAlias(t *testing.T) {
	if got := GetServiceName(); got != ServiceName {
		t.Fatalf("GetServiceName() = %q, want %q", got, ServiceName)
	}
}

func TestPackageTypeTables(t *testing.T) {
	if len(Types) == 0 || len(TypesByVal) == 0 {
		t.Fatal("type lookup tables should not be empty")
	}

	if got, ok := Types["BID"]; !ok || got != covenant.TypeBid {
		t.Fatalf("Types[\"BID\"] = %d, want %d", got, covenant.TypeBid)
	}

	if got, ok := TypesByVal[covenant.TypeBid]; !ok || got != "BID" {
		t.Fatalf("TypesByVal[TypeBid] = %q, want %q", got, "BID")
	}

	if _, ok := Blacklist["test"]; !ok {
		t.Fatal("Blacklist should expose the covenant verifier blacklist")
	}

	if MaxNameSize != 63 {
		t.Fatalf("MaxNameSize = %d, want 63", MaxNameSize)
	}

	if MaxResourceSize != 512 {
		t.Fatalf("MaxResourceSize = %d, want 512", MaxResourceSize)
	}

	if VerifyCovenantsNone != 0 {
		t.Fatalf("VerifyCovenantsNone = %d, want 0", VerifyCovenantsNone)
	}

	if VerifyCovenantsHardened != 1<<0 {
		t.Fatalf("VerifyCovenantsHardened = %d, want %d", VerifyCovenantsHardened, 1<<0)
	}

	if VerifyCovenantsLockup != 1<<1 {
		t.Fatalf("VerifyCovenantsLockup = %d, want %d", VerifyCovenantsLockup, 1<<1)
	}

	if MandatoryVerifyCovenantFlags != VerifyCovenantsNone {
		t.Fatalf("MandatoryVerifyCovenantFlags = %d, want %d", MandatoryVerifyCovenantFlags, VerifyCovenantsNone)
	}

	if MaxCovenantSize != 585 {
		t.Fatalf("MaxCovenantSize = %d, want 585", MaxCovenantSize)
	}

	if CovenantMaxSize != MaxCovenantSize {
		t.Fatalf("CovenantMaxSize = %d, want %d", CovenantMaxSize, MaxCovenantSize)
	}
}

func TestPackageCatalogAliases(t *testing.T) {
	reservedCatalog := ReservedCatalog()
	lockedCatalog := LockedCatalog()

	if reservedCatalog == nil {
		t.Fatal("ReservedCatalog should return a catalog")
	}

	if lockedCatalog == nil {
		t.Fatal("LockedCatalog should return a catalog")
	}

	if GetReservedCatalog() != reservedCatalog {
		t.Fatal("GetReservedCatalog should alias ReservedCatalog")
	}

	if GetLockedCatalog() != lockedCatalog {
		t.Fatal("GetLockedCatalog should alias LockedCatalog")
	}

	if DefaultReservedCatalog() != reservedCatalog {
		t.Fatal("DefaultReservedCatalog should match ReservedCatalog")
	}

	if DefaultLockedCatalog() != lockedCatalog {
		t.Fatal("DefaultLockedCatalog should match LockedCatalog")
	}

	if ReservedSize() != reservedCatalog.Size() || GetReservedSize() != reservedCatalog.Size() {
		t.Fatal("ReservedSize aliases should match the catalog size")
	}

	if LockedSize() != lockedCatalog.Size() || GetLockedSize() != lockedCatalog.Size() {
		t.Fatal("LockedSize aliases should match the catalog size")
	}

	if ReservedPrefixSize() != reservedCatalog.PrefixSize() || GetReservedPrefixSize() != reservedCatalog.PrefixSize() {
		t.Fatal("ReservedPrefixSize aliases should match the catalog prefix")
	}

	if LockedPrefixSize() != lockedCatalog.PrefixSize() || GetLockedPrefixSize() != lockedCatalog.PrefixSize() {
		t.Fatal("LockedPrefixSize aliases should match the catalog prefix")
	}

	if len(ReservedEntries()) == 0 || len(GetReservedEntries()) == 0 {
		t.Fatal("ReservedEntries aliases should expose catalog entries")
	}

	if len(ReservedKeys()) == 0 || len(GetReservedKeys()) == 0 {
		t.Fatal("ReservedKeys aliases should expose catalog keys")
	}

	if len(ReservedValues()) == 0 || len(GetReservedValues()) == 0 {
		t.Fatal("ReservedValues aliases should expose catalog values")
	}

	if len(LockedEntries()) == 0 || len(GetLockedEntries()) == 0 {
		t.Fatal("LockedEntries aliases should expose catalog entries")
	}

	if len(LockedKeys()) == 0 || len(GetLockedKeys()) == 0 {
		t.Fatal("LockedKeys aliases should expose catalog keys")
	}

	if len(LockedValues()) == 0 || len(GetLockedValues()) == 0 {
		t.Fatal("LockedValues aliases should expose catalog values")
	}

	if _, ok := GetReserved("RESERVED.lthn"); !ok {
		t.Fatal("GetReserved should find the reserved reference entry")
	}

	if !HasReserved("reserved.lthn") || !HasReservedByName("reserved") || !HasReservedName("reserved") {
		t.Fatal("reserved lookup aliases should report reserved names")
	}

	if _, ok := GetReservedByString("RESERVED"); !ok {
		t.Fatal("GetReservedByString should find the reserved catalog label")
	}

	if _, ok := GetReservedByBinary([]byte("reserved")); !ok {
		t.Fatal("GetReservedByBinary should find the reserved catalog label")
	}

	if _, ok := GetLocked("NEC.lthn"); !ok {
		t.Fatal("GetLocked should find the locked reference entry")
	}

	if !HasLocked("nec.lthn") || !HasLockedByName("nec") || !HasLockedName("nec") {
		t.Fatal("locked lookup aliases should report locked names")
	}

	if _, ok := GetLockedByString("NEC"); !ok {
		t.Fatal("GetLockedByString should find the locked catalog label")
	}

	if _, ok := GetLockedByBinary([]byte("nec")); !ok {
		t.Fatal("GetLockedByBinary should find the locked catalog label")
	}

	if _, ok := GetReservedByHash(reservedCatalog.Keys()[0]); !ok {
		t.Fatal("GetReservedByHash should return a reserved entry")
	}

	if _, ok := GetLockedByHash(lockedCatalog.Keys()[0]); !ok {
		t.Fatal("GetLockedByHash should return a locked entry")
	}
}
