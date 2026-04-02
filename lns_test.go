// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"crypto/sha3"
	"encoding/binary"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
	"golang.org/x/crypto/blake2b"
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

func TestServiceResolveStringAndBinary(t *testing.T) {
	svc := &Service{}

	got, err := svc.ResolveString("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveString returned error: %v", err)
	}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))
	if got != want {
		t.Fatalf("ResolveString returned %x, want %x", got, want)
	}

	got, err = svc.ResolveBinary([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("ResolveBinary returned error: %v", err)
	}

	if got != want {
		t.Fatalf("ResolveBinary returned %x, want %x", got, want)
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
		"Kfoo.lthn",
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
		{name: "Kfoo.lthn", ok: false},
		{name: 123, ok: false},
	}

	for _, tc := range cases {
		if got := svc.Verify(tc.name); got != tc.ok {
			t.Fatalf("Verify(%v) = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestServiceVerifyStringAndBinary(t *testing.T) {
	svc := &Service{}

	if !svc.VerifyString("Foo-Bar.lthn") {
		t.Fatal("VerifyString should accept canonical names")
	}

	if svc.VerifyString("foo.bar.lthn") {
		t.Fatal("VerifyString should reject malformed names")
	}

	if !svc.VerifyBinary([]byte("Foo-Bar.lthn")) {
		t.Fatal("VerifyBinary should accept canonical names")
	}

	if svc.VerifyBinary([]byte("foo.bar.lthn")) {
		t.Fatal("VerifyBinary should reject malformed names")
	}

	if svc.VerifyString("Kfoo.lthn") {
		t.Fatal("VerifyString should reject non-ASCII names")
	}

	if svc.VerifyBinary([]byte("Kfoo.lthn")) {
		t.Fatal("VerifyBinary should reject non-ASCII names")
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

func TestServiceHashAliases(t *testing.T) {
	svc := &Service{}

	got, err := svc.Hash("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))
	if got != want {
		t.Fatalf("Hash returned %x, want %x", got, want)
	}

	got, err = svc.HashString("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashString returned %x, want %x", got, want)
	}

	got, err = svc.HashBinary([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("HashBinary returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashBinary returned %x, want %x", got, want)
	}

	got, err = svc.HashName("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("HashName returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashName returned %x, want %x", got, want)
	}

	got, err = svc.HashByName([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("HashByName returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashByName returned %x, want %x", got, want)
	}

	if _, err := svc.Hash(123); err == nil {
		t.Fatal("Hash should reject unsupported input types")
	}
}

func TestServiceRuleTables(t *testing.T) {
	svc := &Service{}

	if len(svc.Types()) == 0 || len(svc.TypesByVal()) == 0 {
		t.Fatal("service rule tables should not be empty")
	}

	if len(svc.GetTypes()) != len(svc.Types()) {
		t.Fatal("GetTypes should alias Types")
	}

	if len(svc.GetTypesByVal()) != len(svc.TypesByVal()) {
		t.Fatal("GetTypesByVal should alias TypesByVal")
	}

	if _, ok := svc.Types()["BID"]; !ok {
		t.Fatal("Types should expose BID")
	}

	if got, ok := svc.TypesByVal()[covenant.TypeBid]; !ok || got != "BID" {
		t.Fatalf("TypesByVal[TypeBid] = %q, want %q", got, "BID")
	}

	if len(svc.Blacklist()) == 0 {
		t.Fatal("Blacklist should not be empty")
	}

	if len(svc.GetBlacklist()) != len(svc.Blacklist()) {
		t.Fatal("GetBlacklist should alias Blacklist")
	}

	if _, ok := svc.Blacklist()["test"]; !ok {
		t.Fatal("Blacklist should expose test")
	}
}

func TestServiceVerificationFlagTables(t *testing.T) {
	svc := &Service{}

	if len(svc.VerificationFlags()) != 3 {
		t.Fatalf("VerificationFlags has %d entries, want 3", len(svc.VerificationFlags()))
	}

	if len(svc.GetVerificationFlags()) != len(svc.VerificationFlags()) {
		t.Fatal("GetVerificationFlags should alias VerificationFlags")
	}

	if len(svc.VerificationFlagsByVal()) != len(svc.VerificationFlags()) {
		t.Fatalf("VerificationFlagsByVal has %d entries, want %d", len(svc.VerificationFlagsByVal()), len(svc.VerificationFlags()))
	}

	if len(svc.GetVerificationFlagsByVal()) != len(svc.VerificationFlagsByVal()) {
		t.Fatal("GetVerificationFlagsByVal should alias VerificationFlagsByVal")
	}

	if got, ok := svc.VerificationFlags()["VERIFY_COVENANTS_HARDENED"]; !ok || got != covenant.VerifyCovenantsHardened {
		t.Fatalf("VerificationFlags[HARDENED] = %d, want %d", got, covenant.VerifyCovenantsHardened)
	}

	if got, ok := svc.VerificationFlagsByVal()[covenant.VerifyCovenantsLockup]; !ok || got != "VERIFY_COVENANTS_LOCKUP" {
		t.Fatalf("VerificationFlagsByVal[LOCKUP] = %q, want %q", got, "VERIFY_COVENANTS_LOCKUP")
	}
}

func TestServiceHashByStringAndBinaryAliases(t *testing.T) {
	svc := &Service{}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))

	got, err := svc.HashByString("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("HashByString returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashByString returned %x, want %x", got, want)
	}

	got, err = svc.HashByBinary([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("HashByBinary returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashByBinary returned %x, want %x", got, want)
	}
}

func TestServiceResolveByNameAliases(t *testing.T) {
	svc := &Service{}

	got, err := svc.ResolveByName("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveByName returned error: %v", err)
	}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))
	if got != want {
		t.Fatalf("ResolveByName returned %x, want %x", got, want)
	}

	if _, err := svc.ResolveByName(123); err == nil {
		t.Fatal("ResolveByName should reject unsupported input types")
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

func TestServiceResolveByStringAndBinaryAliases(t *testing.T) {
	svc := &Service{}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))

	got, err := svc.ResolveByString("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveByString returned error: %v", err)
	}

	if got != want {
		t.Fatalf("ResolveByString returned %x, want %x", got, want)
	}

	got, err = svc.ResolveByBinary([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("ResolveByBinary returned error: %v", err)
	}

	if got != want {
		t.Fatalf("ResolveByBinary returned %x, want %x", got, want)
	}
}

func TestServiceVerifyByStringAndBinaryAliases(t *testing.T) {
	svc := &Service{}

	if !svc.VerifyByString("Foo-Bar.lthn") {
		t.Fatal("VerifyByString should accept canonical names")
	}

	if svc.VerifyByString("foo.bar.lthn") {
		t.Fatal("VerifyByString should reject malformed names")
	}

	if !svc.VerifyByBinary([]byte("Foo-Bar.lthn")) {
		t.Fatal("VerifyByBinary should accept canonical names")
	}

	if svc.VerifyByBinary([]byte("foo.bar.lthn")) {
		t.Fatal("VerifyByBinary should reject malformed names")
	}
}

func TestServiceVerifyByNameAliases(t *testing.T) {
	svc := &Service{}

	if !svc.VerifyByName("Foo-Bar.lthn") {
		t.Fatal("VerifyByName should accept canonical names")
	}

	if svc.VerifyByName("foo.bar.lthn") {
		t.Fatal("VerifyByName should reject malformed names")
	}

	if svc.VerifyByName(123) {
		t.Fatal("VerifyByName should reject unsupported input types")
	}
}

func TestServiceBlind(t *testing.T) {
	svc := &Service{}

	var nonce primitives.Hash
	for i := range nonce {
		nonce[i] = byte(i)
	}

	got, err := svc.Blind(0x1122334455667788, nonce)
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
}

func TestServiceTypeName(t *testing.T) {
	svc := &Service{}

	if got := svc.TypeName(covenant.TypeBid); got != "BID" {
		t.Fatalf("TypeName(TypeBid) = %q, want %q", got, "BID")
	}

	if got := svc.TypeName(covenant.CovenantType(99)); got != "UNKNOWN" {
		t.Fatalf("TypeName(99) = %q, want %q", got, "UNKNOWN")
	}
}

func TestServiceGetBlindAndGetTypeNameAliases(t *testing.T) {
	svc := &Service{}

	var nonce primitives.Hash
	got, err := svc.GetBlind(1, nonce)
	if err != nil {
		t.Fatalf("GetBlind returned error: %v", err)
	}

	want, err := svc.Blind(1, nonce)
	if err != nil {
		t.Fatalf("Blind returned error: %v", err)
	}

	if got != want {
		t.Fatalf("GetBlind returned %x, want %x", got, want)
	}

	if got := svc.GetTypeName(covenant.TypeBid); got != "BID" {
		t.Fatalf("GetTypeName(TypeBid) = %q, want %q", got, "BID")
	}
}

func TestServiceCoreAccessors(t *testing.T) {
	runtime := core.NewServiceRuntime(nil, serviceOptions{})
	svc := &Service{ServiceRuntime: runtime}

	if svc.Core() != nil {
		t.Fatal("Core should return nil when the service has no Core")
	}

	if svc.GetCore() != nil {
		t.Fatal("GetCore should alias Core")
	}

	if got := svc.ServiceName(); got != ServiceName {
		t.Fatalf("ServiceName() = %q, want %q", got, ServiceName)
	}

	if got := svc.GetServiceName(); got != ServiceName {
		t.Fatalf("GetServiceName() = %q, want %q", got, ServiceName)
	}
}

func TestServiceCoreNilSafety(t *testing.T) {
	var svc *Service

	if svc.Core() != nil {
		t.Fatal("Core should return nil for a nil receiver")
	}

	if svc.GetCore() != nil {
		t.Fatal("GetCore should return nil for a nil receiver")
	}

	if got := (&Service{}).Core(); got != nil {
		t.Fatal("Core should return nil for a zero-value service")
	}
}

func TestNewService(t *testing.T) {
	svc := NewService(nil)
	if svc == nil {
		t.Fatal("NewService should return a service instance")
	}

	if svc.Core() != nil {
		t.Fatal("NewService(nil) should preserve a nil Core")
	}

	if got := svc.ServiceName(); got != ServiceName {
		t.Fatalf("NewService ServiceName() = %q, want %q", got, ServiceName)
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

func TestServiceGetReservedStringAndBinary(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetReservedString("RESERVED")
	if !ok {
		t.Fatal("GetReservedString should find the reserved catalog label")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	item, ok = svc.GetReservedBinary([]byte("reserved"))
	if !ok {
		t.Fatal("GetReservedBinary should find the reserved catalog label")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if !svc.HasReservedString("reserved") {
		t.Fatal("HasReservedString should report reserved catalog labels")
	}

	if !svc.HasReservedBinary([]byte("RESERVED")) {
		t.Fatal("HasReservedBinary should report reserved catalog labels")
	}

	if _, ok := svc.GetReservedString("RESERVED.lthn"); !ok {
		t.Fatal("GetReservedString should also accept canonical reserved names")
	}

	if _, ok := svc.GetReservedBinary([]byte("reserved.lthn")); !ok {
		t.Fatal("GetReservedBinary should also accept canonical reserved names")
	}

	if !svc.HasReservedString("reserved.lthn") {
		t.Fatal("HasReservedString should report canonical reserved names")
	}

	if !svc.HasReservedBinary([]byte("RESERVED.lthn")) {
		t.Fatal("HasReservedBinary should report canonical reserved names")
	}
}

func TestServiceGetReservedByStringAndBinary(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetReservedByString("RESERVED")
	if !ok {
		t.Fatal("GetReservedByString should find the reserved catalog label")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	item, ok = svc.GetReservedByBinary([]byte("reserved"))
	if !ok {
		t.Fatal("GetReservedByBinary should find the reserved catalog label")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if !svc.HasReservedByString("reserved") {
		t.Fatal("HasReservedByString should report reserved catalog labels")
	}

	if !svc.HasReservedByBinary([]byte("RESERVED")) {
		t.Fatal("HasReservedByBinary should report reserved catalog labels")
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

func TestServiceGetLockedStringAndBinary(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetLockedString("NEC")
	if !ok {
		t.Fatal("GetLockedString should find the locked catalog label")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	item, ok = svc.GetLockedBinary([]byte("nec"))
	if !ok {
		t.Fatal("GetLockedBinary should find the locked catalog label")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if !svc.HasLockedString("nec") {
		t.Fatal("HasLockedString should report locked catalog labels")
	}

	if !svc.HasLockedBinary([]byte("NEC")) {
		t.Fatal("HasLockedBinary should report locked catalog labels")
	}

	if _, ok := svc.GetLockedString("NEC.lthn"); !ok {
		t.Fatal("GetLockedString should also accept canonical locked names")
	}

	if _, ok := svc.GetLockedBinary([]byte("nec.lthn")); !ok {
		t.Fatal("GetLockedBinary should also accept canonical locked names")
	}

	if !svc.HasLockedString("nec.lthn") {
		t.Fatal("HasLockedString should report canonical locked names")
	}

	if !svc.HasLockedBinary([]byte("NEC.lthn")) {
		t.Fatal("HasLockedBinary should report canonical locked names")
	}
}

func TestServiceGetLockedByStringAndBinary(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetLockedByString("NEC")
	if !ok {
		t.Fatal("GetLockedByString should find the locked catalog label")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	item, ok = svc.GetLockedByBinary([]byte("nec"))
	if !ok {
		t.Fatal("GetLockedByBinary should find the locked catalog label")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if !svc.HasLockedByString("nec") {
		t.Fatal("HasLockedByString should report locked catalog labels")
	}

	if !svc.HasLockedByBinary([]byte("NEC")) {
		t.Fatal("HasLockedByBinary should report locked catalog labels")
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
	getCatalog := svc.GetReservedCatalog()
	defaultCatalog := DefaultReservedCatalog()
	packageGetCatalog := GetReservedCatalog()

	if catalog == nil {
		t.Fatal("ReservedCatalog should return a catalog")
	}

	if getCatalog == nil {
		t.Fatal("GetReservedCatalog should return a catalog")
	}

	if getCatalog != catalog {
		t.Fatal("GetReservedCatalog should alias ReservedCatalog")
	}

	if defaultCatalog != catalog {
		t.Fatal("DefaultReservedCatalog should alias the covenant default catalog")
	}

	if packageGetCatalog != catalog {
		t.Fatal("GetReservedCatalog should alias DefaultReservedCatalog")
	}

	if got := svc.ReservedSize(); got != catalog.Size() {
		t.Fatalf("ReservedSize() = %d, want %d", got, catalog.Size())
	}

	if got := svc.GetReservedSize(); got != catalog.Size() {
		t.Fatalf("GetReservedSize() = %d, want %d", got, catalog.Size())
	}

	if got := svc.ReservedNameValue(); got != catalog.NameValue() {
		t.Fatalf("ReservedNameValue() = %d, want %d", got, catalog.NameValue())
	}

	if got := svc.GetReservedNameValue(); got != catalog.NameValue() {
		t.Fatalf("GetReservedNameValue() = %d, want %d", got, catalog.NameValue())
	}

	if got := svc.ReservedRootValue(); got != catalog.RootValue() {
		t.Fatalf("ReservedRootValue() = %d, want %d", got, catalog.RootValue())
	}

	if got := svc.GetReservedRootValue(); got != catalog.RootValue() {
		t.Fatalf("GetReservedRootValue() = %d, want %d", got, catalog.RootValue())
	}

	if got := svc.ReservedTopValue(); got != catalog.TopValue() {
		t.Fatalf("ReservedTopValue() = %d, want %d", got, catalog.TopValue())
	}

	if got := svc.GetReservedTopValue(); got != catalog.TopValue() {
		t.Fatalf("GetReservedTopValue() = %d, want %d", got, catalog.TopValue())
	}

	if got := svc.ReservedPrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("ReservedPrefixSize() = %d, want %d", got, catalog.PrefixSize())
	}

	if got := svc.GetReservedPrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("GetReservedPrefixSize() = %d, want %d", got, catalog.PrefixSize())
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

	if len(svc.GetReservedEntries()) == 0 {
		t.Fatal("GetReservedEntries should expose catalog entries")
	}

	if len(svc.ReservedKeys()) == 0 {
		t.Fatal("ReservedKeys should expose catalog keys")
	}

	if len(svc.GetReservedKeys()) == 0 {
		t.Fatal("GetReservedKeys should expose catalog keys")
	}

	if len(svc.ReservedValues()) == 0 {
		t.Fatal("ReservedValues should expose catalog values")
	}

	if len(svc.GetReservedValues()) == 0 {
		t.Fatal("GetReservedValues should expose catalog values")
	}
}

func TestServiceLockedCatalog(t *testing.T) {
	svc := &Service{}
	catalog := svc.LockedCatalog()
	getCatalog := svc.GetLockedCatalog()
	defaultCatalog := DefaultLockedCatalog()
	packageGetCatalog := GetLockedCatalog()

	if catalog == nil {
		t.Fatal("LockedCatalog should return a catalog")
	}

	if getCatalog == nil {
		t.Fatal("GetLockedCatalog should return a catalog")
	}

	if getCatalog != catalog {
		t.Fatal("GetLockedCatalog should alias LockedCatalog")
	}

	if defaultCatalog != catalog {
		t.Fatal("DefaultLockedCatalog should alias the covenant default catalog")
	}

	if packageGetCatalog != catalog {
		t.Fatal("GetLockedCatalog should alias DefaultLockedCatalog")
	}

	if got := svc.LockedSize(); got != catalog.Size() {
		t.Fatalf("LockedSize() = %d, want %d", got, catalog.Size())
	}

	if got := svc.GetLockedSize(); got != catalog.Size() {
		t.Fatalf("GetLockedSize() = %d, want %d", got, catalog.Size())
	}

	if got := svc.LockedPrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("LockedPrefixSize() = %d, want %d", got, catalog.PrefixSize())
	}

	if got := svc.GetLockedPrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("GetLockedPrefixSize() = %d, want %d", got, catalog.PrefixSize())
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

	if len(svc.GetLockedEntries()) == 0 {
		t.Fatal("GetLockedEntries should expose catalog entries")
	}

	if len(svc.LockedKeys()) == 0 {
		t.Fatal("LockedKeys should expose catalog keys")
	}

	if len(svc.GetLockedKeys()) == 0 {
		t.Fatal("GetLockedKeys should expose catalog keys")
	}

	if len(svc.LockedValues()) == 0 {
		t.Fatal("LockedValues should expose catalog values")
	}

	if len(svc.GetLockedValues()) == 0 {
		t.Fatal("GetLockedValues should expose catalog values")
	}
}
