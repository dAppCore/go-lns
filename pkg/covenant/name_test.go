// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"crypto/sha3"
	"testing"
)

func TestNameSizeLimits(t *testing.T) {
	if MaxNameSize != 63 {
		t.Fatalf("MaxNameSize = %d, want 63", MaxNameSize)
	}

	if MaxResourceSize != 512 {
		t.Fatalf("MaxResourceSize = %d, want 512", MaxResourceSize)
	}

	if MaxCovenantSize != 585 {
		t.Fatalf("MaxCovenantSize = %d, want 585", MaxCovenantSize)
	}

	if CovenantMaxSize != MaxCovenantSize {
		t.Fatalf("CovenantMaxSize = %d, want %d", CovenantMaxSize, MaxCovenantSize)
	}
}

func TestNameFlags(t *testing.T) {
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
}

func TestVerificationFlagTables(t *testing.T) {
	if len(VerificationFlags) != 3 {
		t.Fatalf("VerificationFlags has %d entries, want 3", len(VerificationFlags))
	}

	if len(VerificationFlagsByVal) != len(VerificationFlags) {
		t.Fatalf("VerificationFlagsByVal has %d entries, want %d", len(VerificationFlagsByVal), len(VerificationFlags))
	}

	if got, ok := VerificationFlags["VERIFY_COVENANTS_HARDENED"]; !ok || got != VerifyCovenantsHardened {
		t.Fatalf("VerificationFlags[HARDENED] = %d, want %d", got, VerifyCovenantsHardened)
	}

	if got, ok := VerificationFlagsByVal[VerifyCovenantsLockup]; !ok || got != "VERIFY_COVENANTS_LOCKUP" {
		t.Fatalf("VerificationFlagsByVal[LOCKUP] = %q, want %q", got, "VERIFY_COVENANTS_LOCKUP")
	}

	if GetVerificationFlags()["VERIFY_COVENANTS_NONE"] != VerifyCovenantsNone {
		t.Fatal("GetVerificationFlags should alias VerificationFlags")
	}

	if GetVerificationFlagsByVal()[VerifyCovenantsHardened] != "VERIFY_COVENANTS_HARDENED" {
		t.Fatal("GetVerificationFlagsByVal should alias VerificationFlagsByVal")
	}
}

func TestTypeMaps(t *testing.T) {
	if len(Types) != len(TypesByVal) {
		t.Fatalf("Types and TypesByVal should expose the same number of entries: %d vs %d", len(Types), len(TypesByVal))
	}

	for label, typ := range Types {
		if got, ok := TypesByVal[typ]; !ok || got != label {
			t.Fatalf("TypesByVal[%d] = %q, want %q", typ, got, label)
		}
	}

	if got, ok := Types["BID"]; !ok || got != TypeBid {
		t.Fatalf("Types[\"BID\"] = %d, want %d", got, TypeBid)
	}
}

func TestBlacklist(t *testing.T) {
	for _, name := range []string{"example", "invalid", "local", "localhost", "test"} {
		if _, ok := Blacklist[name]; !ok {
			t.Fatalf("Blacklist should include %q", name)
		}
	}
}

func TestGetTypeName(t *testing.T) {
	if got := GetTypeName(TypeBid); got != "BID" {
		t.Fatalf("GetTypeName(TypeBid) = %q, want %q", got, "BID")
	}

	if got := GetTypeName(CovenantType(99)); got != "UNKNOWN" {
		t.Fatalf("GetTypeName(99) = %q, want %q", got, "UNKNOWN")
	}
}

func TestVerifyString(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{name: "example", ok: false},
		{name: "invalid", ok: false},
		{name: "localhost", ok: false},
		{name: "test", ok: false},
		{name: "", ok: false},
		{name: "foo", ok: true},
		{name: "foo-bar", ok: true},
		{name: "foo_bar", ok: true},
		{name: "-foo", ok: false},
		{name: "foo-", ok: false},
		{name: "_foo", ok: false},
		{name: "foo_", ok: false},
		{name: "Foo", ok: false},
		{name: "f\x7fo", ok: false},
	}

	for _, tc := range cases {
		if got := VerifyString(tc.name); got != tc.ok {
			t.Fatalf("VerifyString(%q) = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestVerifyBinary(t *testing.T) {
	cases := []struct {
		name []byte
		ok   bool
	}{
		{name: nil, ok: false},
		{name: []byte("example"), ok: false},
		{name: []byte("foo"), ok: true},
		{name: []byte("foo-bar"), ok: true},
		{name: []byte("foo_bar"), ok: true},
		{name: []byte("-foo"), ok: false},
		{name: []byte("foo-"), ok: false},
		{name: []byte("Foo"), ok: false},
		{name: []byte{0xff}, ok: false},
	}

	for _, tc := range cases {
		if got := VerifyBinary(tc.name); got != tc.ok {
			t.Fatalf("VerifyBinary(%q) = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestVerifyName(t *testing.T) {
	if !VerifyName("foo-bar") {
		t.Fatal("VerifyName should accept string inputs")
	}

	if !VerifyName([]byte("foo_bar")) {
		t.Fatal("VerifyName should accept byte-slice inputs")
	}

	if VerifyName(123) {
		t.Fatal("VerifyName should reject unsupported input types")
	}
}

func TestVerifyAliases(t *testing.T) {
	if !VerifyByName("foo-bar") {
		t.Fatal("VerifyByName should accept string inputs")
	}

	if !VerifyByString("foo-bar") {
		t.Fatal("VerifyByString should accept string inputs")
	}

	if !VerifyByBinary([]byte("foo_bar")) {
		t.Fatal("VerifyByBinary should accept byte-slice inputs")
	}
}

func TestHashString(t *testing.T) {
	got, err := HashString("foo-bar")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("HashString returned %x, want %x", got, want)
	}
}

func TestHashBinary(t *testing.T) {
	got, err := HashBinary([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("HashBinary returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo_bar"))
	if got != want {
		t.Fatalf("HashBinary returned %x, want %x", got, want)
	}
}

func TestHashName(t *testing.T) {
	got, err := HashName("foo-bar")
	if err != nil {
		t.Fatalf("HashName(string) returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("HashName(string) returned %x, want %x", got, want)
	}

	got, err = HashName([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("HashName([]byte) returned error: %v", err)
	}

	want = sha3.Sum256([]byte("foo_bar"))
	if got != want {
		t.Fatalf("HashName([]byte) returned %x, want %x", got, want)
	}
}

func TestHashAliases(t *testing.T) {
	got, err := HashByName("foo-bar")
	if err != nil {
		t.Fatalf("HashByName(string) returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("HashByName(string) returned %x, want %x", got, want)
	}

	got, err = HashByString("foo-bar")
	if err != nil {
		t.Fatalf("HashByString returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashByString returned %x, want %x", got, want)
	}

	got, err = HashByBinary([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("HashByBinary returned error: %v", err)
	}

	want = sha3.Sum256([]byte("foo_bar"))
	if got != want {
		t.Fatalf("HashByBinary returned %x, want %x", got, want)
	}
}

func TestHashRejectsInvalidName(t *testing.T) {
	if _, err := HashString("Foo"); err == nil {
		t.Fatal("HashString should reject invalid names")
	}

	if _, err := HashBinary([]byte("foo-")); err == nil {
		t.Fatal("HashBinary should reject invalid names")
	}

	if _, err := HashName(123); err == nil {
		t.Fatal("HashName should reject unsupported input types")
	}
}
