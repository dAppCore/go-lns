// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"crypto/sha3"
	"encoding/binary"
	"strings"
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

	if _, err := Resolve(123); err == nil || !strings.Contains(err.Error(), "invalid name type") {
		t.Fatalf("Resolve should reject unsupported input types with an invalid name type error, got %v", err)
	}

	if _, err := Hash(123); err == nil || !strings.Contains(err.Error(), "invalid name type") {
		t.Fatalf("Hash should reject unsupported input types with an invalid name type error, got %v", err)
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

func TestPackageDefaultCatalogAliases(t *testing.T) {
	if DefaultReservedCatalog() != ReservedCatalog() {
		t.Fatal("DefaultReservedCatalog should alias ReservedCatalog")
	}

	if GetReservedCatalog() != ReservedCatalog() {
		t.Fatal("GetReservedCatalog should alias ReservedCatalog")
	}

	if GetDefaultReservedCatalog() != DefaultReservedCatalog() {
		t.Fatal("GetDefaultReservedCatalog should alias DefaultReservedCatalog")
	}

	if DefaultLockedCatalog() != LockedCatalog() {
		t.Fatal("DefaultLockedCatalog should alias LockedCatalog")
	}

	if GetLockedCatalog() != LockedCatalog() {
		t.Fatal("GetLockedCatalog should alias LockedCatalog")
	}

	if GetDefaultLockedCatalog() != DefaultLockedCatalog() {
		t.Fatal("GetDefaultLockedCatalog should alias DefaultLockedCatalog")
	}
}

func TestPackageTypeTables(t *testing.T) {
	if len(Types) == 0 || len(TypesByVal) == 0 {
		t.Fatal("type lookup tables should not be empty")
	}

	if len(GetTypes()) != len(Types) {
		t.Fatal("GetTypes should alias Types")
	}

	if len(GetTypesByVal()) != len(TypesByVal) {
		t.Fatal("GetTypesByVal should alias TypesByVal")
	}

	if got, ok := Types["BID"]; !ok || got != covenant.TypeBid {
		t.Fatalf("Types[\"BID\"] = %d, want %d", got, covenant.TypeBid)
	}

	if got, ok := TypesByVal[covenant.TypeBid]; !ok || got != "BID" {
		t.Fatalf("TypesByVal[TypeBid] = %q, want %q", got, "BID")
	}

	if got, ok := GetTypes()["BID"]; !ok || got != covenant.TypeBid {
		t.Fatalf("GetTypes()[\"BID\"] = %d, want %d", got, covenant.TypeBid)
	}

	if got, ok := GetTypesByVal()[covenant.TypeBid]; !ok || got != "BID" {
		t.Fatalf("GetTypesByVal()[TypeBid] = %q, want %q", got, "BID")
	}

	if CovenantType(TypeBid) != covenant.TypeBid {
		t.Fatalf("CovenantType(TypeBid) = %d, want %d", CovenantType(TypeBid), covenant.TypeBid)
	}
}

func TestPackageBlacklistAlias(t *testing.T) {
	if len(Blacklist) == 0 {
		t.Fatal("Blacklist should not be empty")
	}

	if len(GetBlacklist()) != len(Blacklist) {
		t.Fatal("GetBlacklist should alias Blacklist")
	}

	if _, ok := GetBlacklist()["test"]; !ok {
		t.Fatal("GetBlacklist should expose the covenant verifier blacklist")
	}
}

func TestPackageVerificationFlagTables(t *testing.T) {
	if len(VerificationFlags) != 3 {
		t.Fatalf("VerificationFlags has %d entries, want 3", len(VerificationFlags))
	}

	if len(GetVerificationFlags()) != len(VerificationFlags) {
		t.Fatal("GetVerificationFlags should alias VerificationFlags")
	}

	if len(VerificationFlagsByVal) != len(VerificationFlags) {
		t.Fatalf("VerificationFlagsByVal has %d entries, want %d", len(VerificationFlagsByVal), len(VerificationFlags))
	}

	if len(GetVerificationFlagsByVal()) != len(VerificationFlagsByVal) {
		t.Fatal("GetVerificationFlagsByVal should alias VerificationFlagsByVal")
	}

	if got, ok := VerificationFlags["VERIFY_COVENANTS_LOCKUP"]; !ok || got != covenant.VerifyCovenantsLockup {
		t.Fatalf("VerificationFlags[LOCKUP] = %d, want %d", got, covenant.VerifyCovenantsLockup)
	}

	if got, ok := VerificationFlagsByVal[covenant.VerifyCovenantsNone]; !ok || got != "VERIFY_COVENANTS_NONE" {
		t.Fatalf("VerificationFlagsByVal[None] = %q, want %q", got, "VERIFY_COVENANTS_NONE")
	}
}

func TestPackageConstantGetters(t *testing.T) {
	if GetMaxNameSize() != MaxNameSize {
		t.Fatalf("GetMaxNameSize() = %d, want %d", GetMaxNameSize(), MaxNameSize)
	}

	if GetMaxResourceSize() != MaxResourceSize {
		t.Fatalf("GetMaxResourceSize() = %d, want %d", GetMaxResourceSize(), MaxResourceSize)
	}

	if GetVerifyCovenantsNone() != VerifyCovenantsNone {
		t.Fatalf("GetVerifyCovenantsNone() = %d, want %d", GetVerifyCovenantsNone(), VerifyCovenantsNone)
	}

	if GetVerifyCovenantsHardened() != VerifyCovenantsHardened {
		t.Fatalf("GetVerifyCovenantsHardened() = %d, want %d", GetVerifyCovenantsHardened(), VerifyCovenantsHardened)
	}

	if GetVerifyCovenantsLockup() != VerifyCovenantsLockup {
		t.Fatalf("GetVerifyCovenantsLockup() = %d, want %d", GetVerifyCovenantsLockup(), VerifyCovenantsLockup)
	}

	if GetMandatoryVerifyCovenantFlags() != MandatoryVerifyCovenantFlags {
		t.Fatalf("GetMandatoryVerifyCovenantFlags() = %d, want %d", GetMandatoryVerifyCovenantFlags(), MandatoryVerifyCovenantFlags)
	}

	if GetMaxCovenantSize() != MaxCovenantSize {
		t.Fatalf("GetMaxCovenantSize() = %d, want %d", GetMaxCovenantSize(), MaxCovenantSize)
	}

	if GetCovenantMaxSize() != CovenantMaxSize {
		t.Fatalf("GetCovenantMaxSize() = %d, want %d", GetCovenantMaxSize(), CovenantMaxSize)
	}
}

func TestPackageNameSetHelpers(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeOpen), Items: [][]byte{hash[:], []byte{0, 0, 0, 0}, []byte("example-name")}}},
		},
	}

	set := map[primitives.Hash]struct{}{
		hash: {},
	}

	if !HasNames(tx, set) {
		t.Fatal("HasNames should report a matching name hash")
	}

	RemoveNames(tx, set)

	if len(set) != 0 {
		t.Fatalf("RemoveNames left %d entries, want 0", len(set))
	}

	AddNames(tx, set)

	if len(set) != 1 {
		t.Fatalf("AddNames left %d entries, want 1", len(set))
	}
}

func TestPackageHasSaneCovenants(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	valid := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: primitives.Hash{1}, Index: 0}},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeOpen), Items: [][]byte{hash[:], []byte{0, 0, 0, 0}, []byte("example-name")}}},
		},
	}

	if !HasSaneCovenants(valid) {
		t.Fatal("HasSaneCovenants should accept structurally valid covenants")
	}
}

func TestPackageGrindAndCountHelpers(t *testing.T) {
	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeOpen)}},
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeClaim)}},
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeUpdate)}},
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeTransfer)}},
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeRevoke)}},
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeRegister)}},
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeRenew)}},
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeFinalize)}},
		},
	}

	if got := CountOpens(tx); got != 1 {
		t.Fatalf("CountOpens() = %d, want 1", got)
	}

	if got := CountUpdates(tx); got != 5 {
		t.Fatalf("CountUpdates() = %d, want 5", got)
	}

	if got := CountRenewals(tx); got != 3 {
		t.Fatalf("CountRenewals() = %d, want 3", got)
	}

	rules := NameRules{
		AuctionStart:    100,
		RolloutInterval: 7,
	}

	name, err := GrindName(8, 100, rules)
	if err != nil {
		t.Fatalf("GrindName returned error: %v", err)
	}

	alias, err := GetGrindName(8, 100, rules)
	if err != nil {
		t.Fatalf("GetGrindName returned error: %v", err)
	}

	if len(name) != 8 {
		t.Fatalf("GrindName returned %q with length %d, want 8", name, len(name))
	}

	if len(alias) != 8 {
		t.Fatalf("GetGrindName returned %q with length %d, want 8", alias, len(alias))
	}

	hash, err := covenant.HashString(name)
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if !HasRollout(hash, 100, rules) {
		t.Fatalf("GrindName returned %q that does not satisfy rollout", name)
	}
}

func TestPackageRolloutHelpers(t *testing.T) {
	var hash primitives.Hash

	rules := NameRules{
		AuctionStart:      100,
		RolloutInterval:   7,
		ClaimPeriod:       100,
		AlexaLockupPeriod: 250,
	}

	start, week := GetRollout(hash, rules)
	if start != 100 || week != 0 {
		t.Fatalf("GetRollout(zero) = (%d, %d), want (100, 0)", start, week)
	}

	if !HasRollout(hash, 100, rules) {
		t.Fatal("HasRollout should accept the rollout start height")
	}

	hash, err := covenant.HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if !IsReserved(hash, 99, rules) {
		t.Fatal("IsReserved should report reserved hashes before the claim period")
	}

	rootHash, err := covenant.HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if !IsLockedUp(rootHash, 100, rules) {
		t.Fatal("IsLockedUp should keep root hashes locked after the claim period")
	}
}

type packageTestCoinView struct {
	coins map[primitives.Outpoint]primitives.Output
}

func (v packageTestCoinView) GetOutput(prevout primitives.Outpoint) (primitives.Output, bool) {
	coin, ok := v.coins[prevout]
	return coin, ok
}

func TestPackageVerifyCovenants(t *testing.T) {
	hash, err := covenant.HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	var prevHash primitives.Hash
	prevHash[0] = 3

	var finalHash primitives.Hash
	finalHash[0] = 4
	committedAddress := append([]byte(nil), finalHash[:]...)
	outputAddress := append([]byte(nil), finalHash[:]...)

	transferCovenant := primitives.Covenant{
		Type: uint8(covenant.TypeTransfer),
		Items: [][]byte{
			hash[:],
			make([]byte, 4),
			[]byte{0},
			committedAddress,
		},
	}
	binary.LittleEndian.PutUint32(transferCovenant.Items[1], 100)

	finalizeCovenant := primitives.Covenant{
		Type: uint8(covenant.TypeFinalize),
		Items: [][]byte{
			hash[:],
			make([]byte, 4),
			[]byte("example-name"),
			[]byte{0},
			make([]byte, 4),
			make([]byte, 4),
			make([]byte, 32),
		},
	}
	binary.LittleEndian.PutUint32(finalizeCovenant.Items[1], 100)
	binary.LittleEndian.PutUint32(finalizeCovenant.Items[4], 1)
	binary.LittleEndian.PutUint32(finalizeCovenant.Items[5], 2)

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: prevHash, Index: 0}},
		},
		Outputs: []primitives.Output{
			{
				Value: 1000,
				Address: primitives.Address{
					Version: 0,
					Hash:    outputAddress,
				},
				Covenant: finalizeCovenant,
			},
		},
	}

	view := packageTestCoinView{
		coins: map[primitives.Outpoint]primitives.Output{
			primitives.Outpoint{TxHash: prevHash, Index: 0}: primitives.Output{
				Value: 1000,
				Address: primitives.Address{
					Version: 0,
					Hash:    []byte{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
				},
				Covenant: transferCovenant,
			},
		},
	}

	if got := VerifyCovenants(tx, view, 100, covenant.Network{}); got != 0 {
		t.Fatalf("VerifyCovenants returned %d, want 0", got)
	}

	tx.Outputs[0].Address.Hash[0] ^= 1
	if got := GetVerifyCovenants(tx, view, 100, covenant.Network{}); got != -1 {
		t.Fatalf("GetVerifyCovenants returned %d for an invalid finalize address, want -1", got)
	}
}

func TestPackageTypeConstants(t *testing.T) {
	cases := []struct {
		name string
		got  covenant.CovenantType
		want covenant.CovenantType
	}{
		{name: "TypeNone", got: TypeNone, want: covenant.TypeNone},
		{name: "TypeClaim", got: TypeClaim, want: covenant.TypeClaim},
		{name: "TypeOpen", got: TypeOpen, want: covenant.TypeOpen},
		{name: "TypeBid", got: TypeBid, want: covenant.TypeBid},
		{name: "TypeReveal", got: TypeReveal, want: covenant.TypeReveal},
		{name: "TypeRedeem", got: TypeRedeem, want: covenant.TypeRedeem},
		{name: "TypeRegister", got: TypeRegister, want: covenant.TypeRegister},
		{name: "TypeUpdate", got: TypeUpdate, want: covenant.TypeUpdate},
		{name: "TypeRenew", got: TypeRenew, want: covenant.TypeRenew},
		{name: "TypeTransfer", got: TypeTransfer, want: covenant.TypeTransfer},
		{name: "TypeFinalize", got: TypeFinalize, want: covenant.TypeFinalize},
		{name: "TypeRevoke", got: TypeRevoke, want: covenant.TypeRevoke},
	}

	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("%s = %d, want %d", tc.name, tc.got, tc.want)
		}
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

	if Reserved != reservedCatalog {
		t.Fatal("Reserved should alias ReservedCatalog")
	}

	if Locked != lockedCatalog {
		t.Fatal("Locked should alias LockedCatalog")
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

	if GetDefaultReservedCatalog() != reservedCatalog {
		t.Fatal("GetDefaultReservedCatalog should match ReservedCatalog")
	}

	if GetDefaultLockedCatalog() != lockedCatalog {
		t.Fatal("GetDefaultLockedCatalog should match LockedCatalog")
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

	if !HasReservedHash(reservedCatalog.Keys()[0]) || !HasReservedByHash(reservedCatalog.Keys()[0]) {
		t.Fatal("HasReservedHash aliases should report reserved entries")
	}

	if _, ok := GetLockedByHash(lockedCatalog.Keys()[0]); !ok {
		t.Fatal("GetLockedByHash should return a locked entry")
	}

	if !HasLockedHash(lockedCatalog.Keys()[0]) || !HasLockedByHash(lockedCatalog.Keys()[0]) {
		t.Fatal("HasLockedHash aliases should report locked entries")
	}
}
