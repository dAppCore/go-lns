// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"crypto/sha3"
	"encoding/binary"
	"strings"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
	"golang.org/x/crypto/blake2b"
)

func TestServiceLifecycleHooks(t *testing.T) {
	svc := &Service{}

	if got := svc.OnStartup(nil); !got.OK {
		t.Fatalf("OnStartup returned %#v, want OK", got)
	}

	if got := svc.OnShutdown(nil); !got.OK {
		t.Fatalf("OnShutdown returned %#v, want OK", got)
	}
}

func TestServiceHandleIPCEvents(t *testing.T) {
	svc := &Service{}

	if got := svc.HandleIPCEvents(nil, nil); !got.OK {
		t.Fatalf("HandleIPCEvents returned %#v, want OK", got)
	}
}

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

	got, err = svc.GetResolve("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("GetResolve returned error: %v", err)
	}

	if got != want {
		t.Fatalf("GetResolve returned %x, want %x", got, want)
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

func TestServiceGetResolverAliases(t *testing.T) {
	svc := &Service{}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))

	resolveCases := []struct {
		name string
		fn   func(any) (primitives.Hash, error)
	}{
		{name: "GetResolveString", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveString(name.(string)) }},
		{name: "GetResolveBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveBinary(name.([]byte)) }},
		{name: "GetResolveName", fn: svc.GetResolveName},
		{name: "GetResolveByName", fn: svc.GetResolveByName},
		{name: "GetResolveByString", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveByString(name.(string)) }},
		{name: "GetResolveByBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveByBinary(name.([]byte)) }},
		{name: "GetHashString", fn: func(name any) (primitives.Hash, error) { return svc.GetHashString(name.(string)) }},
		{name: "GetHashBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetHashBinary(name.([]byte)) }},
		{name: "GetHashName", fn: svc.GetHashName},
		{name: "GetHashByName", fn: svc.GetHashByName},
		{name: "GetHashByString", fn: func(name any) (primitives.Hash, error) { return svc.GetHashByString(name.(string)) }},
		{name: "GetHashByBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetHashByBinary(name.([]byte)) }},
	}

	for _, tc := range resolveCases {
		var input any = "Foo-Bar.lthn"
		switch tc.name {
		case "GetResolveBinary", "GetResolveByBinary", "GetHashBinary", "GetHashByBinary":
			input = []byte("Foo-Bar.lthn")
		}

		got, err := tc.fn(input)
		if err != nil {
			t.Fatalf("%s returned error: %v", tc.name, err)
		}

		if got != want {
			t.Fatalf("%s returned %x, want %x", tc.name, got, want)
		}
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

	if _, err := svc.Resolve(123); err == nil || !strings.Contains(err.Error(), "invalid name type") {
		t.Fatalf("Resolve should reject unsupported input types with an invalid name type error, got %v", err)
	}
}

func TestHashOperationLabels(t *testing.T) {
	if _, err := Hash(123); err == nil {
		t.Fatal("Hash should reject unsupported input types")
	} else if got := core.Operation(err); got != "lns.Hash" {
		t.Fatalf("Hash error operation = %q, want %q", got, "lns.Hash")
	}

	if _, err := HashString("foo-"); err == nil {
		t.Fatal("HashString should reject malformed names")
	} else if got := core.Operation(err); got != "lns.HashString" {
		t.Fatalf("HashString error operation = %q, want %q", got, "lns.HashString")
	}

	if _, err := HashBinary([]byte("foo-")); err == nil {
		t.Fatal("HashBinary should reject malformed names")
	} else if got := core.Operation(err); got != "lns.HashBinary" {
		t.Fatalf("HashBinary error operation = %q, want %q", got, "lns.HashBinary")
	}

	svc := &Service{}

	if _, err := svc.Hash(123); err == nil {
		t.Fatal("Service.Hash should reject unsupported input types")
	} else if got := core.Operation(err); got != "lns.Service.Hash" {
		t.Fatalf("Service.Hash error operation = %q, want %q", got, "lns.Service.Hash")
	}

	if _, err := svc.HashString("foo-"); err == nil {
		t.Fatal("Service.HashString should reject malformed names")
	} else if got := core.Operation(err); got != "lns.Service.HashString" {
		t.Fatalf("Service.HashString error operation = %q, want %q", got, "lns.Service.HashString")
	}

	if _, err := svc.HashBinary([]byte("foo-")); err == nil {
		t.Fatal("Service.HashBinary should reject malformed names")
	} else if got := core.Operation(err); got != "lns.Service.HashBinary" {
		t.Fatalf("Service.HashBinary error operation = %q, want %q", got, "lns.Service.HashBinary")
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

func TestServiceRolloutAlias(t *testing.T) {
	svc := &Service{}
	var hash primitives.Hash

	rules := NameRules{
		AuctionStart:    100,
		RolloutInterval: 7,
	}

	if !svc.HasRollout(hash, 100, rules) {
		t.Fatal("HasRollout should accept the rollout start height")
	}

	if !svc.GetHasRollout(hash, 100, rules) {
		t.Fatal("GetHasRollout should alias HasRollout")
	}
}

func TestServiceGetVerifyAliases(t *testing.T) {
	svc := &Service{}

	verifyCases := []struct {
		name  string
		input any
		fn    func(any) bool
	}{
		{name: "GetVerifyName", input: "Foo-Bar.lthn", fn: svc.GetVerifyName},
		{name: "GetVerifyByName", input: "Foo-Bar.lthn", fn: svc.GetVerifyByName},
		{name: "GetVerifyString", input: "Foo-Bar.lthn", fn: func(name any) bool { return svc.GetVerifyString(name.(string)) }},
		{name: "GetVerifyBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) bool { return svc.GetVerifyBinary(name.([]byte)) }},
		{name: "GetVerifyByString", input: "Foo-Bar.lthn", fn: func(name any) bool { return svc.GetVerifyByString(name.(string)) }},
		{name: "GetVerifyByBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) bool { return svc.GetVerifyByBinary(name.([]byte)) }},
	}

	for _, tc := range verifyCases {
		if !tc.fn(tc.input) {
			t.Fatalf("%s should accept canonical names", tc.name)
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

	got, err = svc.GetHash("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("GetHash returned error: %v", err)
	}

	if got != want {
		t.Fatalf("GetHash returned %x, want %x", got, want)
	}

	got, err = svc.HashByName([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("HashByName returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashByName returned %x, want %x", got, want)
	}

	if _, err := svc.Hash(123); err == nil || !strings.Contains(err.Error(), "invalid name type") {
		t.Fatalf("Hash should reject unsupported input types with an invalid name type error, got %v", err)
	}
}

func TestServiceGetHashAndVerifyAliases(t *testing.T) {
	svc := &Service{}

	want := primitives.Hash(sha3.Sum256([]byte("foo-bar")))

	got, err := svc.GetHash("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("GetHash returned error: %v", err)
	}

	if got != want {
		t.Fatalf("GetHash returned %x, want %x", got, want)
	}

	if !svc.GetVerify("Foo-Bar.lthn") {
		t.Fatal("GetVerify should accept canonical names")
	}

	if svc.GetVerify("foo.bar.lthn") {
		t.Fatal("GetVerify should reject malformed names")
	}

	if svc.GetVerify(123) {
		t.Fatal("GetVerify should reject unsupported input types")
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

	if svc.DefaultReservedCatalog() != svc.ReservedCatalog() {
		t.Fatal("DefaultReservedCatalog should alias ReservedCatalog")
	}

	if svc.GetDefaultReservedCatalog() != svc.ReservedCatalog() {
		t.Fatal("GetDefaultReservedCatalog should alias ReservedCatalog")
	}

	if svc.DefaultLockedCatalog() != svc.LockedCatalog() {
		t.Fatal("DefaultLockedCatalog should alias LockedCatalog")
	}

	if svc.GetDefaultLockedCatalog() != svc.LockedCatalog() {
		t.Fatal("GetDefaultLockedCatalog should alias LockedCatalog")
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

func TestServiceConstantGetters(t *testing.T) {
	svc := &Service{}

	if svc.MaxNameSize() != MaxNameSize || svc.GetMaxNameSize() != MaxNameSize {
		t.Fatalf("MaxNameSize getters = %d/%d, want %d", svc.MaxNameSize(), svc.GetMaxNameSize(), MaxNameSize)
	}

	if svc.MaxResourceSize() != MaxResourceSize || svc.GetMaxResourceSize() != MaxResourceSize {
		t.Fatalf("MaxResourceSize getters = %d/%d, want %d", svc.MaxResourceSize(), svc.GetMaxResourceSize(), MaxResourceSize)
	}

	if svc.VerifyCovenantsNone() != VerifyCovenantsNone || svc.GetVerifyCovenantsNone() != VerifyCovenantsNone {
		t.Fatalf("VerifyCovenantsNone getters = %d/%d, want %d", svc.VerifyCovenantsNone(), svc.GetVerifyCovenantsNone(), VerifyCovenantsNone)
	}

	if svc.VerifyCovenantsHardened() != VerifyCovenantsHardened || svc.GetVerifyCovenantsHardened() != VerifyCovenantsHardened {
		t.Fatalf("VerifyCovenantsHardened getters = %d/%d, want %d", svc.VerifyCovenantsHardened(), svc.GetVerifyCovenantsHardened(), VerifyCovenantsHardened)
	}

	if svc.VerifyCovenantsLockup() != VerifyCovenantsLockup || svc.GetVerifyCovenantsLockup() != VerifyCovenantsLockup {
		t.Fatalf("VerifyCovenantsLockup getters = %d/%d, want %d", svc.VerifyCovenantsLockup(), svc.GetVerifyCovenantsLockup(), VerifyCovenantsLockup)
	}

	if svc.MandatoryVerifyCovenantFlags() != MandatoryVerifyCovenantFlags || svc.GetMandatoryVerifyCovenantFlags() != MandatoryVerifyCovenantFlags {
		t.Fatalf("MandatoryVerifyCovenantFlags getters = %d/%d, want %d", svc.MandatoryVerifyCovenantFlags(), svc.GetMandatoryVerifyCovenantFlags(), MandatoryVerifyCovenantFlags)
	}

	if svc.MaxCovenantSize() != MaxCovenantSize || svc.GetMaxCovenantSize() != MaxCovenantSize {
		t.Fatalf("MaxCovenantSize getters = %d/%d, want %d", svc.MaxCovenantSize(), svc.GetMaxCovenantSize(), MaxCovenantSize)
	}

	if svc.CovenantMaxSize() != CovenantMaxSize || svc.GetCovenantMaxSize() != CovenantMaxSize {
		t.Fatalf("CovenantMaxSize getters = %d/%d, want %d", svc.CovenantMaxSize(), svc.GetCovenantMaxSize(), CovenantMaxSize)
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

func TestServiceCovenantTypePredicates(t *testing.T) {
	svc := &Service{}

	if !svc.IsName(covenant.TypeOpen) || !svc.GetIsName(covenant.TypeOpen) {
		t.Fatal("IsName should report open covenants")
	}

	if svc.IsName(covenant.TypeNone) || svc.GetIsName(covenant.TypeNone) {
		t.Fatal("IsName should reject non-name covenants")
	}

	if !svc.IsKnown(covenant.TypeBid) || !svc.GetIsKnown(covenant.TypeBid) {
		t.Fatal("IsKnown should report recognized covenants")
	}

	if svc.IsKnown(covenant.CovenantType(99)) || svc.GetIsKnown(covenant.CovenantType(99)) {
		t.Fatal("IsKnown should reject unknown covenants")
	}

	if !svc.IsLinked(covenant.TypeReveal) || !svc.GetIsLinked(covenant.TypeReveal) {
		t.Fatal("IsLinked should report linked covenants")
	}

	if svc.IsLinked(covenant.TypeOpen) || svc.GetIsLinked(covenant.TypeOpen) {
		t.Fatal("IsLinked should reject unlinked covenants")
	}
}

func TestServiceRolloutHelpers(t *testing.T) {
	svc := &Service{}

	var hash primitives.Hash
	rules := NameRules{
		AuctionStart:      100,
		RolloutInterval:   7,
		ClaimPeriod:       100,
		AlexaLockupPeriod: 250,
	}

	start, week := svc.GetRollout(hash, rules)
	if start != 100 || week != 0 {
		t.Fatalf("GetRollout(zero) = (%d, %d), want (100, 0)", start, week)
	}

	if !svc.HasRollout(hash, 100, rules) {
		t.Fatal("HasRollout should accept the rollout start height")
	}

	reservedHash, err := covenant.HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if !svc.IsReserved(reservedHash, 99, rules) {
		t.Fatal("IsReserved should report reserved hashes before the claim period")
	}

	rootHash, err := covenant.HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if !svc.IsLockedUp(rootHash, 100, rules) {
		t.Fatal("IsLockedUp should keep root hashes locked after the claim period")
	}
}

func TestServiceGrindAndCountHelpers(t *testing.T) {
	svc := &Service{}

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

	if got := svc.CountOpens(tx); got != 1 {
		t.Fatalf("CountOpens() = %d, want 1", got)
	}

	if got := svc.GetCountOpens(tx); got != 1 {
		t.Fatalf("GetCountOpens() = %d, want 1", got)
	}

	if got := svc.CountUpdates(tx); got != 5 {
		t.Fatalf("CountUpdates() = %d, want 5", got)
	}

	if got := svc.GetCountUpdates(tx); got != 5 {
		t.Fatalf("GetCountUpdates() = %d, want 5", got)
	}

	if got := svc.CountRenewals(tx); got != 3 {
		t.Fatalf("CountRenewals() = %d, want 3", got)
	}

	if got := svc.GetCountRenewals(tx); got != 3 {
		t.Fatalf("GetCountRenewals() = %d, want 3", got)
	}

	rules := NameRules{
		AuctionStart:    100,
		RolloutInterval: 7,
	}

	name, err := svc.GrindName(8, 100, rules)
	if err != nil {
		t.Fatalf("GrindName returned error: %v", err)
	}

	alias, err := svc.GetGrindName(8, 100, rules)
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

	if !svc.HasRollout(hash, 100, rules) {
		t.Fatalf("GrindName returned %q that does not satisfy rollout", name)
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

func TestServiceNilReceiverCatalogAccessors(t *testing.T) {
	var svc *Service

	if svc.ReservedCatalog() != Reserved {
		t.Fatal("ReservedCatalog should fall back to the package catalog on a nil receiver")
	}

	if svc.GetReservedCatalog() != Reserved {
		t.Fatal("GetReservedCatalog should fall back to the package catalog on a nil receiver")
	}

	if svc.LockedCatalog() != Locked {
		t.Fatal("LockedCatalog should fall back to the package catalog on a nil receiver")
	}

	if svc.GetLockedCatalog() != Locked {
		t.Fatal("GetLockedCatalog should fall back to the package catalog on a nil receiver")
	}

	if svc.ReservedSize() != Reserved.Size() {
		t.Fatal("ReservedSize should be nil-safe")
	}

	if svc.LockedSize() != Locked.Size() {
		t.Fatal("LockedSize should be nil-safe")
	}

	if len(svc.ReservedEntries()) != len(Reserved.Entries()) {
		t.Fatal("ReservedEntries should be nil-safe")
	}

	if len(svc.LockedEntries()) != len(Locked.Entries()) {
		t.Fatal("LockedEntries should be nil-safe")
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

func TestNewServiceWithOptions(t *testing.T) {
	c := core.New()
	svc := NewServiceWithOptions(WithCore(c))
	if svc == nil {
		t.Fatal("NewServiceWithOptions should return a service instance")
	}

	if svc.Core() != c {
		t.Fatal("WithCore should attach the provided Core instance")
	}

	if got := svc.ServiceName(); got != ServiceName {
		t.Fatalf("NewServiceWithOptions ServiceName() = %q, want %q", got, ServiceName)
	}
}

func TestNewServiceWithNilCoreOption(t *testing.T) {
	svc := NewServiceWithOptions(WithCore(nil))
	if svc == nil {
		t.Fatal("NewServiceWithOptions should return a service instance")
	}

	if svc.Core() != nil {
		t.Fatal("WithCore(nil) should preserve a nil Core")
	}
}

func TestNewServiceWithCatalogOverrides(t *testing.T) {
	reserved := covenant.DefaultReservedCatalog()
	locked := covenant.DefaultLockedCatalog()

	svc := NewServiceWithOptions(
		WithReservedCatalog(reserved),
		WithLockedCatalog(locked),
	)

	if svc.ReservedCatalog() != reserved {
		t.Fatal("WithReservedCatalog should install the provided reserved catalog")
	}

	if svc.LockedCatalog() != locked {
		t.Fatal("WithLockedCatalog should install the provided locked catalog")
	}
}

func TestCatalogOptionIgnoresNilCatalogs(t *testing.T) {
	svc := NewServiceWithOptions(
		WithReservedCatalog(nil),
		WithLockedCatalog(nil),
	)

	if svc.ReservedCatalog() != ReservedCatalog() {
		t.Fatal("WithReservedCatalog(nil) should not replace the default reserved catalog")
	}

	if svc.LockedCatalog() != LockedCatalog() {
		t.Fatal("WithLockedCatalog(nil) should not replace the default locked catalog")
	}
}

func TestRegisterRejectsNilCore(t *testing.T) {
	res := Register(nil)
	if res.OK {
		t.Fatal("Register(nil) should fail")
	}

	if _, ok := res.Value.(error); !ok {
		t.Fatalf("Register(nil) should return an error result, got %T", res.Value)
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

	if catalog != Reserved {
		t.Fatal("ReservedCatalog should return the exported Reserved instance")
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

	if got := svc.NameValue(); got != catalog.NameValue() {
		t.Fatalf("NameValue() = %d, want %d", got, catalog.NameValue())
	}

	if got := svc.GetNameValue(); got != catalog.NameValue() {
		t.Fatalf("GetNameValue() = %d, want %d", got, catalog.NameValue())
	}

	if got := svc.RootValue(); got != catalog.RootValue() {
		t.Fatalf("RootValue() = %d, want %d", got, catalog.RootValue())
	}

	if got := svc.GetRootValue(); got != catalog.RootValue() {
		t.Fatalf("GetRootValue() = %d, want %d", got, catalog.RootValue())
	}

	if got := svc.TopValue(); got != catalog.TopValue() {
		t.Fatalf("TopValue() = %d, want %d", got, catalog.TopValue())
	}

	if got := svc.GetTopValue(); got != catalog.TopValue() {
		t.Fatalf("GetTopValue() = %d, want %d", got, catalog.TopValue())
	}

	if got := svc.PrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("PrefixSize() = %d, want %d", got, catalog.PrefixSize())
	}

	if got := svc.GetPrefixSize(); got != catalog.PrefixSize() {
		t.Fatalf("GetPrefixSize() = %d, want %d", got, catalog.PrefixSize())
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

	if catalog != Locked {
		t.Fatal("LockedCatalog should return the exported Locked instance")
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

func TestServiceNameSetHelpers(t *testing.T) {
	svc := &Service{}

	hash, err := svc.HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: primitives.Hash{1}, Index: 0}},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(covenant.TypeOpen), Items: [][]byte{hash[:], []byte{0, 0, 0, 0}, []byte("example-name")}}},
		},
	}

	set := map[primitives.Hash]struct{}{
		hash: {},
	}

	if !svc.HasNames(tx, set) {
		t.Fatal("HasNames should report a matching name hash")
	}

	if !svc.GetHasNames(tx, set) {
		t.Fatal("GetHasNames should report a matching name hash")
	}

	svc.RemoveNames(tx, set)
	svc.GetRemoveNames(tx, set)

	if len(set) != 0 {
		t.Fatalf("RemoveNames/GetRemoveNames left %d entries, want 0", len(set))
	}

	svc.AddNames(tx, set)
	svc.GetAddNames(tx, set)

	if len(set) != 1 {
		t.Fatalf("AddNames/GetAddNames left %d entries, want 1", len(set))
	}
}

func TestServiceVerifyCovenants(t *testing.T) {
	svc := &Service{}

	hash, err := covenant.HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	var prevHash primitives.Hash
	prevHash[0] = 5

	var finalHash primitives.Hash
	finalHash[0] = 6
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

	if !svc.HasSaneCovenants(tx) {
		t.Fatal("HasSaneCovenants should accept structurally valid covenants")
	}

	if !svc.GetHasSaneCovenants(tx) {
		t.Fatal("GetHasSaneCovenants should accept structurally valid covenants")
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

	if got := svc.VerifyCovenants(tx, view, 100, Network{}); got != 0 {
		t.Fatalf("VerifyCovenants returned %d, want 0", got)
	}

	tx.Outputs[0].Address.Hash[0] ^= 1
	if got := svc.GetVerifyCovenants(tx, view, 100, Network{}); got != -1 {
		t.Fatalf("GetVerifyCovenants returned %d for an invalid finalize address, want -1", got)
	}
}
