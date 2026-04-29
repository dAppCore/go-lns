// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"crypto/sha3"
	"encoding/binary"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
	"golang.org/x/crypto/blake2b"
)

func TestServiceLifecycleHooks(t *testing.T) {
	svc := &Service{}

	if got := svc.OnStartup(nil); !got.OK {
		t.Fatalf("OnStartup returned %#v, want OK", got)
	}

	if got := svc.GetOnStartup(context.Background()); !got.OK {
		t.Fatalf("GetOnStartup returned %#v, want OK", got)
	}

	if got := svc.OnShutdown(nil); !got.OK {
		t.Fatalf("OnShutdown returned %#v, want OK", got)
	}

	if got := svc.GetOnShutdown(context.Background()); !got.OK {
		t.Fatalf("GetOnShutdown returned %#v, want OK", got)
	}

	var nilSvc *Service

	if got := nilSvc.OnStartup(nil); !got.OK {
		t.Fatalf("nil OnStartup returned %#v, want OK", got)
	}

	if got := nilSvc.GetOnStartup(context.Background()); !got.OK {
		t.Fatalf("nil GetOnStartup returned %#v, want OK", got)
	}

	if got := nilSvc.OnShutdown(nil); !got.OK {
		t.Fatalf("nil OnShutdown returned %#v, want OK", got)
	}

	if got := nilSvc.GetOnShutdown(context.Background()); !got.OK {
		t.Fatalf("nil GetOnShutdown returned %#v, want OK", got)
	}
}

func TestServiceHandleIPCEvents(t *testing.T) {
	svc := &Service{}

	if got := svc.HandleIPCEvents(nil, nil); !got.OK {
		t.Fatalf("HandleIPCEvents returned %#v, want OK", got)
	}

	if got := svc.GetHandleIPCEvents(nil, nil); !got.OK {
		t.Fatalf("GetHandleIPCEvents returned %#v, want OK", got)
	}

	var nilSvc *Service

	if got := nilSvc.HandleIPCEvents(nil, nil); !got.OK {
		t.Fatalf("nil HandleIPCEvents returned %#v, want OK", got)
	}

	if got := nilSvc.GetHandleIPCEvents(nil, nil); !got.OK {
		t.Fatalf("nil GetHandleIPCEvents returned %#v, want OK", got)
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

	if _, err := svc.Resolve(123); err == nil || !core.Contains(err.Error(), "invalid name type") {
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

func TestServiceHasAliasMirrors(t *testing.T) {
	svc := &Service{}

	if !svc.GetHasReserved("RESERVED.lthn") {
		t.Fatal("GetHasReserved should report reserved canonical names")
	}

	if !svc.GetHasReservedString("reserved") {
		t.Fatal("GetHasReservedString should report reserved catalog labels")
	}

	if !svc.GetHasReservedHash(primitives.Hash(sha3.Sum256([]byte("reserved")))) {
		t.Fatal("GetHasReservedHash should report reserved canonical hashes")
	}

	if !svc.GetHasLocked("NEC.lthn") {
		t.Fatal("GetHasLocked should report locked canonical names")
	}

	if !svc.GetHasLockedString("nec") {
		t.Fatal("GetHasLockedString should report locked catalog labels")
	}

	if !svc.GetHasLockedHash(primitives.Hash(sha3.Sum256([]byte("nec")))) {
		t.Fatal("GetHasLockedHash should report locked canonical hashes")
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

	if _, err := svc.Hash(123); err == nil || !core.Contains(err.Error(), "invalid name type") {
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

	if !svc.GetIsReserved(reservedHash, 99, rules) {
		t.Fatal("GetIsReserved should alias IsReserved")
	}

	rootHash, err := covenant.HashString("nec")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if !svc.IsLockedUp(rootHash, 100, rules) {
		t.Fatal("IsLockedUp should keep root hashes locked after the claim period")
	}

	if !svc.GetIsLockedUp(rootHash, 100, rules) {
		t.Fatal("GetIsLockedUp should alias IsLockedUp")
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

func TestServiceUtilityAliases(t *testing.T) {
	svc := &Service{}

	if svc.NewClaim() == nil {
		t.Fatal("NewClaim should return a claim wrapper")
	}

	if svc.GetNewClaim() == nil {
		t.Fatal("GetNewClaim should alias NewClaim")
	}

	if view := svc.NewNameView(); view == nil {
		t.Fatal("NewNameView should return a cached name-state view")
	} else if view.ToNameUndo().Names != nil {
		t.Fatal("NewNameView should start with an empty undo payload")
	}

	if svc.GetNewNameView() == nil {
		t.Fatal("GetNewNameView should alias NewNameView")
	}

	if got := svc.NextName("ExAmPle."); got != "example\x00." {
		t.Fatalf("NextName returned %q", got)
	}

	if got := svc.GetNextName("ExAmPle."); got != "example\x00." {
		t.Fatalf("GetNextName returned %q", got)
	}

	if got := svc.PrevName("ExAmPle."); got != "exampld\xff." {
		t.Fatalf("PrevName returned %q", got)
	}

	if got := svc.GetPrevName("ExAmPle."); got != "exampld\xff." {
		t.Fatalf("GetPrevName returned %q", got)
	}

	if svc.DefaultTTL() != DEFAULT_TTL || svc.GetDefaultTTL() != DEFAULT_TTL {
		t.Fatalf("DefaultTTL getters = %d/%d, want %d", svc.DefaultTTL(), svc.GetDefaultTTL(), DEFAULT_TTL)
	}

	if len(svc.Dummy()) != 0 || len(svc.GetDummy()) != 0 {
		t.Fatal("Dummy getters should return the zero-length DNS reference buffer")
	}

	if !core.DeepEqual(svc.TypeMapRoot(), TYPE_MAP_ROOT) || !core.DeepEqual(svc.GetTypeMapRoot(), TYPE_MAP_ROOT) {
		t.Fatal("TypeMapRoot getters should alias the reference bitmap")
	}

	if !core.DeepEqual(svc.TypeMapEmpty(), TYPE_MAP_EMPTY) || !core.DeepEqual(svc.GetTypeMapEmpty(), TYPE_MAP_EMPTY) {
		t.Fatal("TypeMapEmpty getters should alias the reference bitmap")
	}

	if !core.DeepEqual(svc.TypeMapNS(), TYPE_MAP_NS) || !core.DeepEqual(svc.GetTypeMapNS(), TYPE_MAP_NS) {
		t.Fatal("TypeMapNS getters should alias the reference bitmap")
	}

	if !core.DeepEqual(svc.TypeMapTXT(), TYPE_MAP_TXT) || !core.DeepEqual(svc.GetTypeMapTXT(), TYPE_MAP_TXT) {
		t.Fatal("TypeMapTXT getters should alias the reference bitmap")
	}

	if !core.DeepEqual(svc.TypeMapA(), TYPE_MAP_A) || !core.DeepEqual(svc.GetTypeMapA(), TYPE_MAP_A) {
		t.Fatal("TypeMapA getters should alias the reference bitmap")
	}

	if !core.DeepEqual(svc.TypeMapAAAA(), TYPE_MAP_AAAA) || !core.DeepEqual(svc.GetTypeMapAAAA(), TYPE_MAP_AAAA) {
		t.Fatal("TypeMapAAAA getters should alias the reference bitmap")
	}

	if len(svc.HSTypes()) != len(HSTypes) || len(svc.GetHSTypes()) != len(HSTypes) {
		t.Fatal("HSTypes getters should alias the reference table")
	}

	if len(svc.HSTypesByVal()) != len(HSTypesByVal) || len(svc.GetHSTypesByVal()) != len(HSTypesByVal) {
		t.Fatal("HSTypesByVal getters should alias the reference table")
	}

	record := svc.Create(".", svc.NextName("."), TYPE_MAP_ROOT)
	if record.Name != "." {
		t.Fatalf("Create Name = %q, want %q", record.Name, ".")
	}

	if record.NextDomain != "\x00." {
		t.Fatalf("Create NextDomain = %q, want %q", record.NextDomain, "\x00.")
	}

	if record.TTL != DEFAULT_TTL {
		t.Fatalf("Create TTL = %d, want %d", record.TTL, DEFAULT_TTL)
	}

	if !core.DeepEqual(record.TypeBitmap, TYPE_MAP_ROOT) {
		t.Fatalf("Create TypeBitmap = %x, want %x", record.TypeBitmap, TYPE_MAP_ROOT)
	}

	if svc.GetCreate("foo", "bar", TYPE_MAP_NS).TTL != DEFAULT_TTL {
		t.Fatal("GetCreate should alias Create")
	}

	resource := svc.NewResource()
	if resource == nil {
		t.Fatal("NewResource should return a resource")
	}

	if resource.TTL != DEFAULT_TTL {
		t.Fatalf("NewResource TTL = %d, want %d", resource.TTL, DEFAULT_TTL)
	}

	if svc.GetNewResource() == nil {
		t.Fatal("GetNewResource should return a resource")
	}

	resource.Records = []ResourceRecord{TXTRecord{Entries: []string{"hello"}}}
	encoded, err := resource.Encode()
	if err != nil {
		t.Fatalf("Resource.Encode returned error: %v", err)
	}

	decoded, err := svc.DecodeResource(encoded)
	if err != nil {
		t.Fatalf("DecodeResource returned error: %v", err)
	}

	if len(decoded.Records) != len(resource.Records) {
		t.Fatalf("DecodeResource records = %d, want %d", len(decoded.Records), len(resource.Records))
	}

	if got, err := svc.GetDecodeResource(encoded); err != nil {
		t.Fatalf("GetDecodeResource returned error: %v", err)
	} else if got.TTL != decoded.TTL || len(got.Records) != len(decoded.Records) {
		t.Fatal("GetDecodeResource should alias DecodeResource")
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

	if got := GetNewService(nil); got == nil {
		t.Fatal("GetNewService should return a service instance")
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

	if got := GetNewServiceWithOptions(WithCore(c)); got == nil {
		t.Fatal("GetNewServiceWithOptions should return a service instance")
	} else if got.Core() != c {
		t.Fatal("GetNewServiceWithOptions should attach the provided Core instance")
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

func TestServiceDefaultCatalogIgnoresOverrides(t *testing.T) {
	reserved := &covenant.ReservedCatalog{}
	locked := &covenant.LockedCatalog{}

	svc := NewServiceWithOptions(
		WithReservedCatalog(reserved),
		WithLockedCatalog(locked),
	)

	if svc.ReservedCatalog() != reserved {
		t.Fatal("ReservedCatalog should return the service override")
	}

	if svc.DefaultReservedCatalog() != DefaultReservedCatalog() {
		t.Fatal("DefaultReservedCatalog should return the package default catalog")
	}

	if svc.DefaultReservedCatalog() == svc.ReservedCatalog() {
		t.Fatal("DefaultReservedCatalog should ignore the service override")
	}

	if svc.LockedCatalog() != locked {
		t.Fatal("LockedCatalog should return the service override")
	}

	if svc.DefaultLockedCatalog() != DefaultLockedCatalog() {
		t.Fatal("DefaultLockedCatalog should return the package default catalog")
	}

	if svc.DefaultLockedCatalog() == svc.LockedCatalog() {
		t.Fatal("DefaultLockedCatalog should ignore the service override")
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

func TestRegisterWithOptions(t *testing.T) {
	reserved := &covenant.ReservedCatalog{}
	locked := &covenant.LockedCatalog{}

	app := core.New(
		core.WithService(RegisterWithOptions(
			WithReservedCatalog(reserved),
			WithLockedCatalog(locked),
		)),
	)

	svc, ok := core.ServiceFor[*Service](app, ServiceName)
	if !ok {
		t.Fatal("RegisterWithOptions should register the LNS service")
	}

	if svc == nil {
		t.Fatal("RegisterWithOptions should return a service instance")
	}

	if svc.ReservedCatalog() != reserved {
		t.Fatal("RegisterWithOptions should apply the reserved catalog override")
	}

	if svc.LockedCatalog() != locked {
		t.Fatal("RegisterWithOptions should apply the locked catalog override")
	}

	if svc.Core() != app {
		t.Fatal("RegisterWithOptions should attach the Core instance")
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

	if !svc.HasReservedCatalog() || !svc.GetHasReservedCatalog() {
		t.Fatal("HasReservedCatalog aliases should report the service reserved catalog")
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

	if !svc.HasLockedCatalog() || !svc.GetHasLockedCatalog() {
		t.Fatal("HasLockedCatalog aliases should report the service locked catalog")
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

// BEGIN v090 file-aware generated coverage

func TestLns_WithCore_Good(t *core.T) {
	fn := ax7ValueOf(WithCore)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithCore_Bad(t *core.T) {
	fn := ax7ValueOf(WithCore)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithCore_Ugly(t *core.T) {
	fn := ax7ValueOf(WithCore)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithCore_Good(t *core.T) {
	fn := ax7ValueOf(GetWithCore)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithCore_Bad(t *core.T) {
	fn := ax7ValueOf(GetWithCore)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithCore_Ugly(t *core.T) {
	fn := ax7ValueOf(GetWithCore)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(WithReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(WithReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(WithReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetWithReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetWithReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetWithReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(WithLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(WithLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_WithLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(WithLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetWithLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetWithLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetWithLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetWithLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewServiceWithOptions_Good(t *core.T) {
	fn := ax7ValueOf(NewServiceWithOptions)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewServiceWithOptions_Bad(t *core.T) {
	fn := ax7ValueOf(NewServiceWithOptions)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewServiceWithOptions_Ugly(t *core.T) {
	fn := ax7ValueOf(NewServiceWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewServiceWithOptions_Good(t *core.T) {
	fn := ax7ValueOf(GetNewServiceWithOptions)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewServiceWithOptions_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewServiceWithOptions)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewServiceWithOptions_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewServiceWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewService_Good(t *core.T) {
	fn := ax7ValueOf(NewService)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewService_Bad(t *core.T) {
	fn := ax7ValueOf(NewService)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewService_Ugly(t *core.T) {
	fn := ax7ValueOf(NewService)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewService_Good(t *core.T) {
	fn := ax7ValueOf(GetNewService)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewService_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewService)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewService_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewService)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Register_Good(t *core.T) {
	fn := ax7ValueOf(Register)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Register_Bad(t *core.T) {
	fn := ax7ValueOf(Register)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Register_Ugly(t *core.T) {
	fn := ax7ValueOf(Register)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRegister_Good(t *core.T) {
	fn := ax7ValueOf(GetRegister)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRegister_Bad(t *core.T) {
	fn := ax7ValueOf(GetRegister)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRegister_Ugly(t *core.T) {
	fn := ax7ValueOf(GetRegister)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RegisterWithOptions_Good(t *core.T) {
	fn := ax7ValueOf(RegisterWithOptions)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RegisterWithOptions_Bad(t *core.T) {
	fn := ax7ValueOf(RegisterWithOptions)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RegisterWithOptions_Ugly(t *core.T) {
	fn := ax7ValueOf(RegisterWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRegisterWithOptions_Good(t *core.T) {
	fn := ax7ValueOf(GetRegisterWithOptions)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRegisterWithOptions_Bad(t *core.T) {
	fn := ax7ValueOf(GetRegisterWithOptions)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRegisterWithOptions_Ugly(t *core.T) {
	fn := ax7ValueOf(GetRegisterWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetServiceName_Good(t *core.T) {
	fn := ax7ValueOf(GetServiceName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetServiceName_Bad(t *core.T) {
	fn := ax7ValueOf(GetServiceName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetServiceName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetServiceName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultTTL_Good(t *core.T) {
	fn := ax7ValueOf(GetDefaultTTL)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultTTL_Bad(t *core.T) {
	fn := ax7ValueOf(GetDefaultTTL)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultTTL_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDefaultTTL)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDummy_Good(t *core.T) {
	fn := ax7ValueOf(GetDummy)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDummy_Bad(t *core.T) {
	fn := ax7ValueOf(GetDummy)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDummy_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDummy)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapRoot_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapRoot)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapRoot_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapRoot)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapRoot_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapRoot)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapEmpty_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapEmpty_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapEmpty_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapNS_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapNS_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapNS_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapTXT_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapTXT)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapTXT_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapTXT)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapTXT_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapTXT)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapA_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapA_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapA_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapAAAA_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapAAAA_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeMapAAAA_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHSTypes_Good(t *core.T) {
	fn := ax7ValueOf(GetHSTypes)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHSTypes_Bad(t *core.T) {
	fn := ax7ValueOf(GetHSTypes)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHSTypes_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHSTypes)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHSTypesByVal_Good(t *core.T) {
	fn := ax7ValueOf(GetHSTypesByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHSTypesByVal_Bad(t *core.T) {
	fn := ax7ValueOf(GetHSTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHSTypesByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHSTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewResource_Good(t *core.T) {
	fn := ax7ValueOf(NewResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewResource_Bad(t *core.T) {
	fn := ax7ValueOf(NewResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewResource_Ugly(t *core.T) {
	fn := ax7ValueOf(NewResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewResource_Good(t *core.T) {
	fn := ax7ValueOf(GetNewResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewResource_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewResource_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DecodeResource_Good(t *core.T) {
	fn := ax7ValueOf(DecodeResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DecodeResource_Bad(t *core.T) {
	fn := ax7ValueOf(DecodeResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DecodeResource_Ugly(t *core.T) {
	fn := ax7ValueOf(DecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDecodeResource_Good(t *core.T) {
	fn := ax7ValueOf(GetDecodeResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDecodeResource_Bad(t *core.T) {
	fn := ax7ValueOf(GetDecodeResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDecodeResource_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Create_Good(t *core.T) {
	fn := ax7ValueOf(Create)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Create_Bad(t *core.T) {
	fn := ax7ValueOf(Create)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Create_Ugly(t *core.T) {
	fn := ax7ValueOf(Create)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCreate_Good(t *core.T) {
	fn := ax7ValueOf(GetCreate)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCreate_Bad(t *core.T) {
	fn := ax7ValueOf(GetCreate)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCreate_Ugly(t *core.T) {
	fn := ax7ValueOf(GetCreate)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NextName_Good(t *core.T) {
	fn := ax7ValueOf(NextName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NextName_Bad(t *core.T) {
	fn := ax7ValueOf(NextName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NextName_Ugly(t *core.T) {
	fn := ax7ValueOf(NextName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNextName_Good(t *core.T) {
	fn := ax7ValueOf(GetNextName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNextName_Bad(t *core.T) {
	fn := ax7ValueOf(GetNextName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNextName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNextName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_PrevName_Good(t *core.T) {
	fn := ax7ValueOf(PrevName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_PrevName_Bad(t *core.T) {
	fn := ax7ValueOf(PrevName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_PrevName_Ugly(t *core.T) {
	fn := ax7ValueOf(PrevName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetPrevName_Good(t *core.T) {
	fn := ax7ValueOf(GetPrevName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetPrevName_Bad(t *core.T) {
	fn := ax7ValueOf(GetPrevName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetPrevName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetPrevName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewClaim_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).NewClaim)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewClaim_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).NewClaim)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewClaim_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).NewClaim)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewClaim_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewClaim)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewClaim_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewClaim)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewClaim_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewClaim)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewOutpoint_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).NewOutpoint)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewOutpoint_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).NewOutpoint)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewOutpoint_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).NewOutpoint)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewOutpoint_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewOutpoint)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewOutpoint_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewOutpoint)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewOutpoint_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewOutpoint)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewNameView_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).NewNameView)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewNameView_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).NewNameView)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewNameView_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).NewNameView)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewNameView_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewNameView)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewNameView_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewNameView)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewNameView_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewNameView)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Create_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Create)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Create_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Create)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Create_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Create)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCreate_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCreate)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCreate_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCreate)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCreate_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCreate)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultTTL_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultTTL)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultTTL_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultTTL)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultTTL_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultTTL)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultTTL_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultTTL)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultTTL_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultTTL)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultTTL_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultTTL)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Dummy_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Dummy)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Dummy_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Dummy)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Dummy_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Dummy)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDummy_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDummy)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDummy_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDummy)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDummy_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDummy)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapRoot_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapRoot)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapRoot_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapRoot)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapRoot_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapRoot)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapRoot_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapRoot)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapRoot_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapRoot)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapRoot_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapRoot)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapEmpty_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapEmpty)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapEmpty_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapEmpty)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapEmpty_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapEmpty)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapEmpty_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapEmpty_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapEmpty_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapTXT_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapTXT)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapTXT_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapTXT)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapTXT_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapTXT)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapTXT_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapTXT)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapTXT_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapTXT)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapTXT_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapTXT)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapA_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapA_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapA_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapA_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapA_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapA_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapAAAA_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapAAAA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapAAAA_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapAAAA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeMapAAAA_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeMapAAAA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapAAAA_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapAAAA_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeMapAAAA_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HSTypes_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HSTypes)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HSTypes_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HSTypes)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HSTypes_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HSTypes)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHSTypes_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHSTypes)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHSTypes_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHSTypes)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHSTypes_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHSTypes)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HSTypesByVal_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HSTypesByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HSTypesByVal_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HSTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HSTypesByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HSTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHSTypesByVal_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHSTypesByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHSTypesByVal_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHSTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHSTypesByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHSTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NextName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).NextName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NextName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).NextName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NextName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).NextName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNextName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNextName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNextName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNextName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNextName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNextName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_PrevName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).PrevName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_PrevName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).PrevName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_PrevName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).PrevName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetPrevName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetPrevName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetPrevName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetPrevName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetPrevName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetPrevName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewResource_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).NewResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewResource_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).NewResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NewResource_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).NewResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewResource_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewResource_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNewResource_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNewResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DecodeResource_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).DecodeResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DecodeResource_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).DecodeResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DecodeResource_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).DecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDecodeResource_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDecodeResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDecodeResource_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDecodeResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDecodeResource_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_OnStartup_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).OnStartup)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_OnStartup_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).OnStartup)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_OnStartup_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).OnStartup)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetOnStartup_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetOnStartup)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetOnStartup_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetOnStartup)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetOnStartup_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetOnStartup)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_OnShutdown_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).OnShutdown)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_OnShutdown_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).OnShutdown)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_OnShutdown_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).OnShutdown)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetOnShutdown_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetOnShutdown)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetOnShutdown_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetOnShutdown)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetOnShutdown_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetOnShutdown)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HandleIPCEvents_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HandleIPCEvents)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HandleIPCEvents_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HandleIPCEvents)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HandleIPCEvents_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HandleIPCEvents)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHandleIPCEvents_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHandleIPCEvents)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHandleIPCEvents_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHandleIPCEvents)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHandleIPCEvents_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHandleIPCEvents)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(LookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(LookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(LookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(GetLookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(GetLookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(HasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(HasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).LookupCatalogName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).LookupCatalogName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).LookupCatalogName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLookupCatalogName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLookupCatalogName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLookupCatalogName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLookupCatalogName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLookupCatalogName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLookupCatalogName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLookupCatalogName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLookupCatalogName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLookupCatalogName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLookupCatalogName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLookupCatalogName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLookupCatalogName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Resolve_Good(t *core.T) {
	fn := ax7ValueOf(Resolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Resolve_Bad(t *core.T) {
	fn := ax7ValueOf(Resolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Resolve_Ugly(t *core.T) {
	fn := ax7ValueOf(Resolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolve_Good(t *core.T) {
	fn := ax7ValueOf(GetResolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolve_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolve_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Hash_Good(t *core.T) {
	fn := ax7ValueOf(Hash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Hash_Bad(t *core.T) {
	fn := ax7ValueOf(Hash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Hash_Ugly(t *core.T) {
	fn := ax7ValueOf(Hash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Verify_Good(t *core.T) {
	fn := ax7ValueOf(Verify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Verify_Bad(t *core.T) {
	fn := ax7ValueOf(Verify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Verify_Ugly(t *core.T) {
	fn := ax7ValueOf(Verify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerify_Good(t *core.T) {
	fn := ax7ValueOf(GetVerify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerify_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerify_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashString_Good(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashString_Bad(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashString_Ugly(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashString_Good(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashBinary_Good(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(ResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(ResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByString_Good(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByString_Bad(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByString_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveString_Good(t *core.T) {
	fn := ax7ValueOf(ResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveString_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(ResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveName_Good(t *core.T) {
	fn := ax7ValueOf(ResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveName_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashName_Good(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashName_Bad(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashName_Ugly(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashName_Good(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyName_Good(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyName_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(ResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByName_Good(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByName_Bad(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByName_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyString_Good(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyString_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Blind_Good(t *core.T) {
	fn := ax7ValueOf(Blind)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Blind_Bad(t *core.T) {
	fn := ax7ValueOf(Blind)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Blind_Ugly(t *core.T) {
	fn := ax7ValueOf(Blind)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetBlind_Good(t *core.T) {
	fn := ax7ValueOf(GetBlind)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetBlind_Bad(t *core.T) {
	fn := ax7ValueOf(GetBlind)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetBlind_Ugly(t *core.T) {
	fn := ax7ValueOf(GetBlind)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_TypeName_Good(t *core.T) {
	fn := ax7ValueOf(TypeName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_TypeName_Bad(t *core.T) {
	fn := ax7ValueOf(TypeName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_TypeName_Ugly(t *core.T) {
	fn := ax7ValueOf(TypeName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetBlacklist_Good(t *core.T) {
	fn := ax7ValueOf(GetBlacklist)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetBlacklist_Bad(t *core.T) {
	fn := ax7ValueOf(GetBlacklist)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetBlacklist_Ugly(t *core.T) {
	fn := ax7ValueOf(GetBlacklist)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerificationFlags_Good(t *core.T) {
	fn := ax7ValueOf(GetVerificationFlags)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerificationFlags_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerificationFlags)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerificationFlags_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerificationFlags)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerificationFlagsByVal_Good(t *core.T) {
	fn := ax7ValueOf(GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerificationFlagsByVal_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerificationFlagsByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewClaim_Good(t *core.T) {
	fn := ax7ValueOf(NewClaim)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewClaim_Bad(t *core.T) {
	fn := ax7ValueOf(NewClaim)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewClaim_Ugly(t *core.T) {
	fn := ax7ValueOf(NewClaim)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewClaim_Good(t *core.T) {
	fn := ax7ValueOf(GetNewClaim)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewClaim_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewClaim)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewClaim_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewClaim)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewOutpoint_Good(t *core.T) {
	fn := ax7ValueOf(NewOutpoint)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewOutpoint_Bad(t *core.T) {
	fn := ax7ValueOf(NewOutpoint)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewOutpoint_Ugly(t *core.T) {
	fn := ax7ValueOf(NewOutpoint)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewOutpoint_Good(t *core.T) {
	fn := ax7ValueOf(GetNewOutpoint)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewOutpoint_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewOutpoint)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewOutpoint_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewOutpoint)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewNameView_Good(t *core.T) {
	fn := ax7ValueOf(NewNameView)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewNameView_Bad(t *core.T) {
	fn := ax7ValueOf(NewNameView)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NewNameView_Ugly(t *core.T) {
	fn := ax7ValueOf(NewNameView)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewNameView_Good(t *core.T) {
	fn := ax7ValueOf(GetNewNameView)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewNameView_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewNameView)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNewNameView_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewNameView)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxNameSize_Good(t *core.T) {
	fn := ax7ValueOf(GetMaxNameSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxNameSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetMaxNameSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxNameSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetMaxNameSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxResourceSize_Good(t *core.T) {
	fn := ax7ValueOf(GetMaxResourceSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxResourceSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetMaxResourceSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxResourceSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetMaxResourceSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsNone_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsNone_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsNone_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsHardened_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsHardened_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsHardened_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsLockup_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsLockup_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenantsLockup_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMandatoryVerifyCovenantFlags_Good(t *core.T) {
	fn := ax7ValueOf(GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMandatoryVerifyCovenantFlags_Bad(t *core.T) {
	fn := ax7ValueOf(GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMandatoryVerifyCovenantFlags_Ugly(t *core.T) {
	fn := ax7ValueOf(GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxCovenantSize_Good(t *core.T) {
	fn := ax7ValueOf(GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxCovenantSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetMaxCovenantSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCovenantMaxSize_Good(t *core.T) {
	fn := ax7ValueOf(GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCovenantMaxSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCovenantMaxSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypes_Good(t *core.T) {
	fn := ax7ValueOf(GetTypes)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypes_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypes)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypes_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypes)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypesByVal_Good(t *core.T) {
	fn := ax7ValueOf(GetTypesByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypesByVal_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypesByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeName_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeName_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTypeName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsName_Good(t *core.T) {
	fn := ax7ValueOf(IsName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsName_Bad(t *core.T) {
	fn := ax7ValueOf(IsName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsName_Ugly(t *core.T) {
	fn := ax7ValueOf(IsName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsName_Good(t *core.T) {
	fn := ax7ValueOf(GetIsName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsName_Bad(t *core.T) {
	fn := ax7ValueOf(GetIsName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetIsName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsKnown_Good(t *core.T) {
	fn := ax7ValueOf(IsKnown)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsKnown_Bad(t *core.T) {
	fn := ax7ValueOf(IsKnown)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsKnown_Ugly(t *core.T) {
	fn := ax7ValueOf(IsKnown)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsKnown_Good(t *core.T) {
	fn := ax7ValueOf(GetIsKnown)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsKnown_Bad(t *core.T) {
	fn := ax7ValueOf(GetIsKnown)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsKnown_Ugly(t *core.T) {
	fn := ax7ValueOf(GetIsKnown)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsLinked_Good(t *core.T) {
	fn := ax7ValueOf(IsLinked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsLinked_Bad(t *core.T) {
	fn := ax7ValueOf(IsLinked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsLinked_Ugly(t *core.T) {
	fn := ax7ValueOf(IsLinked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsLinked_Good(t *core.T) {
	fn := ax7ValueOf(GetIsLinked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsLinked_Bad(t *core.T) {
	fn := ax7ValueOf(GetIsLinked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsLinked_Ugly(t *core.T) {
	fn := ax7ValueOf(GetIsLinked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxNameSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxNameSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxNameSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxNameSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxNameSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxNameSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxNameSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxNameSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxNameSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxNameSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxNameSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxNameSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxResourceSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxResourceSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxResourceSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxResourceSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxResourceSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxResourceSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxResourceSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxResourceSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxResourceSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxResourceSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxResourceSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxResourceSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsNone_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsNone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsNone_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsNone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsNone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsNone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsNone_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsNone_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsNone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsHardened_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsHardened_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsHardened_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsHardened_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsHardened_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsHardened_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsLockup_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsLockup_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenantsLockup_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsLockup_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsLockup_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenantsLockup_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MandatoryVerifyCovenantFlags_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).MandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MandatoryVerifyCovenantFlags_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).MandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MandatoryVerifyCovenantFlags_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).MandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMandatoryVerifyCovenantFlags_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMandatoryVerifyCovenantFlags_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMandatoryVerifyCovenantFlags_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxCovenantSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxCovenantSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxCovenantSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxCovenantSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_MaxCovenantSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).MaxCovenantSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxCovenantSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxCovenantSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetMaxCovenantSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CovenantMaxSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).CovenantMaxSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CovenantMaxSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).CovenantMaxSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CovenantMaxSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).CovenantMaxSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCovenantMaxSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCovenantMaxSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCovenantMaxSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRollout_Good(t *core.T) {
	fn := ax7ValueOf(GetRollout)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRollout_Bad(t *core.T) {
	fn := ax7ValueOf(GetRollout)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRollout_Ugly(t *core.T) {
	fn := ax7ValueOf(GetRollout)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasRollout_Good(t *core.T) {
	fn := ax7ValueOf(HasRollout)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasRollout_Bad(t *core.T) {
	fn := ax7ValueOf(HasRollout)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasRollout_Ugly(t *core.T) {
	fn := ax7ValueOf(HasRollout)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasRollout_Good(t *core.T) {
	fn := ax7ValueOf(GetHasRollout)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasRollout_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasRollout)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasRollout_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasRollout)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsReserved_Good(t *core.T) {
	fn := ax7ValueOf(IsReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsReserved_Bad(t *core.T) {
	fn := ax7ValueOf(IsReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(IsReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsReserved_Good(t *core.T) {
	fn := ax7ValueOf(GetIsReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsReserved_Bad(t *core.T) {
	fn := ax7ValueOf(GetIsReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(GetIsReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsLockedUp_Good(t *core.T) {
	fn := ax7ValueOf(IsLockedUp)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsLockedUp_Bad(t *core.T) {
	fn := ax7ValueOf(IsLockedUp)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_IsLockedUp_Ugly(t *core.T) {
	fn := ax7ValueOf(IsLockedUp)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsLockedUp_Good(t *core.T) {
	fn := ax7ValueOf(GetIsLockedUp)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsLockedUp_Bad(t *core.T) {
	fn := ax7ValueOf(GetIsLockedUp)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetIsLockedUp_Ugly(t *core.T) {
	fn := ax7ValueOf(GetIsLockedUp)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GrindName_Good(t *core.T) {
	fn := ax7ValueOf(GrindName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GrindName_Bad(t *core.T) {
	fn := ax7ValueOf(GrindName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GrindName_Ugly(t *core.T) {
	fn := ax7ValueOf(GrindName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetGrindName_Good(t *core.T) {
	fn := ax7ValueOf(GetGrindName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetGrindName_Bad(t *core.T) {
	fn := ax7ValueOf(GetGrindName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetGrindName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetGrindName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountOpens_Good(t *core.T) {
	fn := ax7ValueOf(CountOpens)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountOpens_Bad(t *core.T) {
	fn := ax7ValueOf(CountOpens)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountOpens_Ugly(t *core.T) {
	fn := ax7ValueOf(CountOpens)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountOpens_Good(t *core.T) {
	fn := ax7ValueOf(GetCountOpens)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountOpens_Bad(t *core.T) {
	fn := ax7ValueOf(GetCountOpens)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountOpens_Ugly(t *core.T) {
	fn := ax7ValueOf(GetCountOpens)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountUpdates_Good(t *core.T) {
	fn := ax7ValueOf(CountUpdates)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountUpdates_Bad(t *core.T) {
	fn := ax7ValueOf(CountUpdates)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountUpdates_Ugly(t *core.T) {
	fn := ax7ValueOf(CountUpdates)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountUpdates_Good(t *core.T) {
	fn := ax7ValueOf(GetCountUpdates)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountUpdates_Bad(t *core.T) {
	fn := ax7ValueOf(GetCountUpdates)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountUpdates_Ugly(t *core.T) {
	fn := ax7ValueOf(GetCountUpdates)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountRenewals_Good(t *core.T) {
	fn := ax7ValueOf(CountRenewals)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountRenewals_Bad(t *core.T) {
	fn := ax7ValueOf(CountRenewals)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_CountRenewals_Ugly(t *core.T) {
	fn := ax7ValueOf(CountRenewals)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountRenewals_Good(t *core.T) {
	fn := ax7ValueOf(GetCountRenewals)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountRenewals_Bad(t *core.T) {
	fn := ax7ValueOf(GetCountRenewals)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetCountRenewals_Ugly(t *core.T) {
	fn := ax7ValueOf(GetCountRenewals)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasNames_Good(t *core.T) {
	fn := ax7ValueOf(HasNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasNames_Bad(t *core.T) {
	fn := ax7ValueOf(HasNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasNames_Ugly(t *core.T) {
	fn := ax7ValueOf(HasNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasNames_Good(t *core.T) {
	fn := ax7ValueOf(GetHasNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasNames_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasNames_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_AddNames_Good(t *core.T) {
	fn := ax7ValueOf(AddNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_AddNames_Bad(t *core.T) {
	fn := ax7ValueOf(AddNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_AddNames_Ugly(t *core.T) {
	fn := ax7ValueOf(AddNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetAddNames_Good(t *core.T) {
	fn := ax7ValueOf(GetAddNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetAddNames_Bad(t *core.T) {
	fn := ax7ValueOf(GetAddNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetAddNames_Ugly(t *core.T) {
	fn := ax7ValueOf(GetAddNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RemoveNames_Good(t *core.T) {
	fn := ax7ValueOf(RemoveNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RemoveNames_Bad(t *core.T) {
	fn := ax7ValueOf(RemoveNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RemoveNames_Ugly(t *core.T) {
	fn := ax7ValueOf(RemoveNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRemoveNames_Good(t *core.T) {
	fn := ax7ValueOf(GetRemoveNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRemoveNames_Bad(t *core.T) {
	fn := ax7ValueOf(GetRemoveNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRemoveNames_Ugly(t *core.T) {
	fn := ax7ValueOf(GetRemoveNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasSaneCovenants_Good(t *core.T) {
	fn := ax7ValueOf(HasSaneCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasSaneCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(HasSaneCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasSaneCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(HasSaneCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasSaneCovenants_Good(t *core.T) {
	fn := ax7ValueOf(GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasSaneCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasSaneCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyCovenants_Good(t *core.T) {
	fn := ax7ValueOf(VerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_VerifyCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenants_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetVerifyCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ServiceName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ServiceName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ServiceName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ServiceName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ServiceName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ServiceName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetServiceName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetServiceName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetServiceName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetServiceName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetServiceName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetServiceName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Resolve_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Resolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Resolve_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Resolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Resolve_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Resolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolve_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolve_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolve_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Hash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Hash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Hash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Hash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Hash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Hash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Verify_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Verify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Verify_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Verify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Verify_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Verify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerify_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerify_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerify_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Blind_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Blind)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Blind_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Blind)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Blind_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Blind)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetBlind_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetBlind)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetBlind_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetBlind)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetBlind_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetBlind)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypeName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypeName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypeName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypeName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).IsName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).IsName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).IsName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsKnown_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).IsKnown)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsKnown_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).IsKnown)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsKnown_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).IsKnown)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsKnown_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsKnown)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsKnown_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsKnown)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsKnown_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsKnown)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsLinked_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).IsLinked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsLinked_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).IsLinked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsLinked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).IsLinked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsLinked_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsLinked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsLinked_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsLinked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsLinked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsLinked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRollout_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRollout)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRollout_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRollout)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRollout_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRollout)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasRollout_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasRollout)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasRollout_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasRollout)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasRollout_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasRollout)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasRollout_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasRollout)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasRollout_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasRollout)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasRollout_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasRollout)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsReserved_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).IsReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsReserved_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).IsReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).IsReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsReserved_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsReserved_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsLockedUp_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).IsLockedUp)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsLockedUp_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).IsLockedUp)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_IsLockedUp_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).IsLockedUp)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsLockedUp_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsLockedUp)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsLockedUp_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsLockedUp)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetIsLockedUp_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetIsLockedUp)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GrindName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GrindName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GrindName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GrindName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GrindName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GrindName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetGrindName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetGrindName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetGrindName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetGrindName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetGrindName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetGrindName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountOpens_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).CountOpens)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountOpens_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).CountOpens)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountOpens_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).CountOpens)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountOpens_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountOpens)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountOpens_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountOpens)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountOpens_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountOpens)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountUpdates_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).CountUpdates)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountUpdates_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).CountUpdates)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountUpdates_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).CountUpdates)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountUpdates_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountUpdates)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountUpdates_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountUpdates)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountUpdates_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountUpdates)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountRenewals_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).CountRenewals)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountRenewals_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).CountRenewals)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_CountRenewals_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).CountRenewals)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountRenewals_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountRenewals)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountRenewals_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountRenewals)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCountRenewals_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCountRenewals)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasNames_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasNames_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasNames_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasNames_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasNames_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasNames_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_AddNames_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).AddNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_AddNames_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).AddNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_AddNames_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).AddNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetAddNames_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetAddNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetAddNames_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetAddNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetAddNames_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetAddNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_RemoveNames_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).RemoveNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_RemoveNames_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).RemoveNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_RemoveNames_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).RemoveNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRemoveNames_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRemoveNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRemoveNames_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRemoveNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRemoveNames_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRemoveNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasSaneCovenants_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasSaneCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasSaneCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasSaneCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasSaneCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasSaneCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasSaneCovenants_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasSaneCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasSaneCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenants_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerifyCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenants_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerifyCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Types_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Types)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Types_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Types)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Types_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Types)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypes_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypes)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypes_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypes)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypes_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypes)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypesByVal_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TypesByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypesByVal_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TypesByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TypesByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypesByVal_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypesByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypesByVal_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTypesByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Blacklist_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Blacklist)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Blacklist_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Blacklist)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Blacklist_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Blacklist)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetBlacklist_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetBlacklist)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetBlacklist_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetBlacklist)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetBlacklist_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetBlacklist)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerificationFlags_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerificationFlags)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerificationFlags_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerificationFlags)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerificationFlags_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerificationFlags)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerificationFlags_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerificationFlags)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerificationFlags_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerificationFlags)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerificationFlags_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerificationFlags)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerificationFlagsByVal_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerificationFlagsByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerificationFlagsByVal_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerificationFlagsByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_VerificationFlagsByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerificationFlagsByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerificationFlagsByVal_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerificationFlagsByVal_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetVerificationFlagsByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReserved_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReserved_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReserved_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReserved_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLocked_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLocked_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLocked_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLocked_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Core_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Core)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Core_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Core)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_Core_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Core)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCore_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCore)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCore_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCore)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetCore_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCore)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DefaultReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DefaultReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DefaultReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedNameValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedNameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedNameValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedNameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedNameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedNameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedNameValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedNameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedNameValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedNameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedNameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedNameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedRootValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedRootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedRootValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedRootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedRootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedRootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedRootValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedRootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedRootValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedRootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedRootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedRootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedTopValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedTopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedTopValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedTopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedTopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedTopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedTopValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedTopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedTopValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedTopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedTopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedTopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NameValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).NameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NameValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).NameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_NameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).NameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNameValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNameValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetNameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetNameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_RootValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).RootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_RootValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).RootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_RootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).RootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRootValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRootValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetRootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetRootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TopValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).TopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TopValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).TopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_TopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).TopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTopValue_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTopValue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetTopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetTopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_PrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).PrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_PrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).PrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_PrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).PrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedEntries_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedEntries_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedKeys_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedKeys_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedValues_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedValues_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedValues_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedValues_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_HasLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HasLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_DefaultLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetDefaultLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DefaultLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DefaultLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_DefaultLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultLockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultLockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetDefaultLockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(ReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedCatalog_Good(t *core.T) {
	fn := ax7ValueOf(LockedCatalog)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedCatalog_Bad(t *core.T) {
	fn := ax7ValueOf(LockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedCatalog_Ugly(t *core.T) {
	fn := ax7ValueOf(LockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedSize_Good(t *core.T) {
	fn := ax7ValueOf(ReservedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedSize_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedSize_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedNameValue_Good(t *core.T) {
	fn := ax7ValueOf(ReservedNameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedNameValue_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedNameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedNameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedNameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedNameValue_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedNameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedNameValue_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedNameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedNameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedNameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedRootValue_Good(t *core.T) {
	fn := ax7ValueOf(ReservedRootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedRootValue_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedRootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedRootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedRootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedRootValue_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedRootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedRootValue_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedRootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedRootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedRootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedTopValue_Good(t *core.T) {
	fn := ax7ValueOf(ReservedTopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedTopValue_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedTopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedTopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedTopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedTopValue_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedTopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedTopValue_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedTopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedTopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedTopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NameValue_Good(t *core.T) {
	fn := ax7ValueOf(NameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NameValue_Bad(t *core.T) {
	fn := ax7ValueOf(NameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_NameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(NameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNameValue_Good(t *core.T) {
	fn := ax7ValueOf(GetNameValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNameValue_Bad(t *core.T) {
	fn := ax7ValueOf(GetNameValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetNameValue_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNameValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RootValue_Good(t *core.T) {
	fn := ax7ValueOf(RootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RootValue_Bad(t *core.T) {
	fn := ax7ValueOf(RootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_RootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(RootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRootValue_Good(t *core.T) {
	fn := ax7ValueOf(GetRootValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRootValue_Bad(t *core.T) {
	fn := ax7ValueOf(GetRootValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetRootValue_Ugly(t *core.T) {
	fn := ax7ValueOf(GetRootValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_TopValue_Good(t *core.T) {
	fn := ax7ValueOf(TopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_TopValue_Bad(t *core.T) {
	fn := ax7ValueOf(TopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_TopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(TopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTopValue_Good(t *core.T) {
	fn := ax7ValueOf(GetTopValue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTopValue_Bad(t *core.T) {
	fn := ax7ValueOf(GetTopValue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetTopValue_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTopValue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_PrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(PrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_PrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(PrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_PrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(PrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(GetPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(ReservedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedEntries_Good(t *core.T) {
	fn := ax7ValueOf(ReservedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedEntries_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedKeys_Good(t *core.T) {
	fn := ax7ValueOf(ReservedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedKeys_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedValues_Good(t *core.T) {
	fn := ax7ValueOf(ReservedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedValues_Bad(t *core.T) {
	fn := ax7ValueOf(ReservedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_ReservedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(ReservedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedValues_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedValues_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedSize_Good(t *core.T) {
	fn := ax7ValueOf(LockedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedSize_Bad(t *core.T) {
	fn := ax7ValueOf(LockedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(LockedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedSize_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(LockedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(LockedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(LockedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedEntries_Good(t *core.T) {
	fn := ax7ValueOf(LockedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(LockedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(LockedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedEntries_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedKeys_Good(t *core.T) {
	fn := ax7ValueOf(LockedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(LockedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(LockedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedKeys_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedValues_Good(t *core.T) {
	fn := ax7ValueOf(LockedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedValues_Bad(t *core.T) {
	fn := ax7ValueOf(LockedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_LockedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(LockedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedValues_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedValues_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReserved_Good(t *core.T) {
	fn := ax7ValueOf(GetReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReserved_Bad(t *core.T) {
	fn := ax7ValueOf(GetReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedName_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedString_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReserved_Good(t *core.T) {
	fn := ax7ValueOf(HasReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReserved_Bad(t *core.T) {
	fn := ax7ValueOf(HasReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReserved_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReserved_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReserved_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReserved)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReserved_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReserved)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReserved_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReserved)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedString_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedName_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(HasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasReservedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLocked_Good(t *core.T) {
	fn := ax7ValueOf(GetLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLocked_Bad(t *core.T) {
	fn := ax7ValueOf(GetLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedName_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedString_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLocked_Good(t *core.T) {
	fn := ax7ValueOf(HasLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLocked_Bad(t *core.T) {
	fn := ax7ValueOf(HasLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLocked_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLocked_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLocked_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLocked)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLocked_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLocked)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLocked_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLocked)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedString_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedName_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedName_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(HasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(HasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_HasLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(HasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_GetHasLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetHasLockedByHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_ReservedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ReservedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedPrefixSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedPrefixSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetReservedPrefixSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedEntries_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedEntries_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedEntries)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedEntries_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedEntries)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedEntries_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedEntries)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedKeys_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedKeys_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedKeys)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedKeys_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedKeys)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedKeys_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedKeys)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedValues_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedValues_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_LockedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).LockedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedValues_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedValues)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedValues_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedValues)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestLns_Service_GetLockedValues_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetLockedValues)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
