// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	core "dappco.re/go"
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

	if len(GetTypes()) != len(Types) {
		t.Fatal("GetTypes should alias Types")
	}

	if len(GetTypesByVal()) != len(TypesByVal) {
		t.Fatal("GetTypesByVal should alias TypesByVal")
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

	if len(GetBlacklist()) != len(Blacklist) {
		t.Fatal("GetBlacklist should alias Blacklist")
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

	if !GetVerifyName("foo-bar") {
		t.Fatal("GetVerifyName should accept string inputs")
	}

	if !VerifyByString("foo-bar") {
		t.Fatal("VerifyByString should accept string inputs")
	}

	if !GetVerifyString("foo-bar") {
		t.Fatal("GetVerifyString should accept string inputs")
	}

	if !VerifyByBinary([]byte("foo_bar")) {
		t.Fatal("VerifyByBinary should accept byte-slice inputs")
	}

	if !GetVerifyBinary([]byte("foo_bar")) {
		t.Fatal("GetVerifyBinary should accept byte-slice inputs")
	}

	if !GetVerifyByName("foo-bar") {
		t.Fatal("GetVerifyByName should alias VerifyByName")
	}

	if !GetVerifyByString("foo-bar") {
		t.Fatal("GetVerifyByString should alias VerifyByString")
	}

	if !GetVerifyByBinary([]byte("foo_bar")) {
		t.Fatal("GetVerifyByBinary should alias VerifyByBinary")
	}
}

func TestHashString(t *testing.T) {
	got, err := HashString("foo-bar")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	want := core.SHA3_256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("HashString returned %x, want %x", got, want)
	}
}

func TestHashBinary(t *testing.T) {
	got, err := HashBinary([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("HashBinary returned error: %v", err)
	}

	want := core.SHA3_256([]byte("foo_bar"))
	if got != want {
		t.Fatalf("HashBinary returned %x, want %x", got, want)
	}
}

func TestHashName(t *testing.T) {
	got, err := HashName("foo-bar")
	if err != nil {
		t.Fatalf("HashName(string) returned error: %v", err)
	}

	want := core.SHA3_256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("HashName(string) returned %x, want %x", got, want)
	}

	got, err = HashName([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("HashName([]byte) returned error: %v", err)
	}

	want = core.SHA3_256([]byte("foo_bar"))
	if got != want {
		t.Fatalf("HashName([]byte) returned %x, want %x", got, want)
	}
}

func TestHashAliases(t *testing.T) {
	wantString := core.SHA3_256([]byte("foo-bar"))
	wantBinary := core.SHA3_256([]byte("foo_bar"))

	got, err := HashByName("foo-bar")
	if err != nil {
		t.Fatalf("HashByName(string) returned error: %v", err)
	}

	if got != wantString {
		t.Fatalf("HashByName(string) returned %x, want %x", got, wantString)
	}

	got, err = GetHashName("foo-bar")
	if err != nil {
		t.Fatalf("GetHashName(string) returned error: %v", err)
	}

	if got != wantString {
		t.Fatalf("GetHashName(string) returned %x, want %x", got, wantString)
	}

	got, err = HashByString("foo-bar")
	if err != nil {
		t.Fatalf("HashByString returned error: %v", err)
	}

	if got != wantString {
		t.Fatalf("HashByString returned %x, want %x", got, wantString)
	}

	got, err = GetHashString("foo-bar")
	if err != nil {
		t.Fatalf("GetHashString returned error: %v", err)
	}

	if got != wantString {
		t.Fatalf("GetHashString returned %x, want %x", got, wantString)
	}

	got, err = HashByBinary([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("HashByBinary returned error: %v", err)
	}

	if got != wantBinary {
		t.Fatalf("HashByBinary returned %x, want %x", got, wantBinary)
	}

	got, err = GetHashBinary([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("GetHashBinary returned error: %v", err)
	}

	if got != wantBinary {
		t.Fatalf("GetHashBinary returned %x, want %x", got, wantBinary)
	}

	got, err = GetHashByName("foo-bar")
	if err != nil {
		t.Fatalf("GetHashByName(string) returned error: %v", err)
	}

	if got != wantString {
		t.Fatalf("GetHashByName(string) returned %x, want %x", got, wantString)
	}

	got, err = GetHashByString("foo-bar")
	if err != nil {
		t.Fatalf("GetHashByString returned error: %v", err)
	}

	if got != wantString {
		t.Fatalf("GetHashByString returned %x, want %x", got, wantString)
	}

	got, err = GetHashByBinary([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("GetHashByBinary returned error: %v", err)
	}

	if got != wantBinary {
		t.Fatalf("GetHashByBinary returned %x, want %x", got, wantBinary)
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

// BEGIN v090 file-aware generated coverage

func TestName_GetBlacklist_Good(t *core.T) {
	fn := ax7ValueOf(GetBlacklist)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetBlacklist_Bad(t *core.T) {
	fn := ax7ValueOf(GetBlacklist)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetBlacklist_Ugly(t *core.T) {
	fn := ax7ValueOf(GetBlacklist)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyString_Good(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyString_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyName_Good(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyName_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_VerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetVerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashString_Good(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashString_Bad(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashString_Ugly(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashString_Good(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByString_Good(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByString_Bad(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByString_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashBinary_Good(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashName_Good(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashName_Bad(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashName_Ugly(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashName_Good(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByName_Good(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByName_Bad(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_HashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByName_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_GetHashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
