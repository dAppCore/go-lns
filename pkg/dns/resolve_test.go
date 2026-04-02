// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"crypto/sha3"
	"testing"

	core "dappco.re/go/core"
)

func TestResolve(t *testing.T) {
	svc := NewService()

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

func TestResolveNameAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.ResolveName("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveName returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("ResolveName returned %x, want %x", got, want)
	}

	if _, err := svc.ResolveName(123); err == nil {
		t.Fatal("ResolveName should reject unsupported input types")
	}
}

func TestHashAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.Hash("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
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

func TestResolveByNameAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.ResolveByName("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveByName returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("ResolveByName returned %x, want %x", got, want)
	}

	if _, err := svc.ResolveByName(123); err == nil {
		t.Fatal("ResolveByName should reject unsupported input types")
	}
}

func TestResolveStringAndBinary(t *testing.T) {
	svc := NewService()

	got, err := svc.ResolveString("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveString returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
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

func TestResolveByStringAndBinaryAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.ResolveByString("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("ResolveByString returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
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

func TestHashByStringAndBinaryAliases(t *testing.T) {
	svc := NewService()

	want := sha3.Sum256([]byte("foo-bar"))

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

func TestResolveRejectsInvalidNames(t *testing.T) {
	svc := NewService()

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

func TestVerifyStringAndBinary(t *testing.T) {
	svc := NewService()

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
}

func TestVerifyByStringAndBinaryAliases(t *testing.T) {
	svc := NewService()

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

func TestVerify(t *testing.T) {
	svc := NewService()

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

func TestVerifyNameAliases(t *testing.T) {
	svc := NewService()

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

func TestVerifyByNameAliases(t *testing.T) {
	svc := NewService()

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

func TestCoreAccessors(t *testing.T) {
	c := &core.Core{}
	svc := NewService(WithCore(c))

	if svc.Core() != c {
		t.Fatal("Core should return the configured Core instance")
	}

	if svc.GetCore() != c {
		t.Fatal("GetCore should alias Core")
	}
}

func TestServiceName(t *testing.T) {
	svc := NewService()

	if got := svc.GetServiceName(); got != ServiceName {
		t.Fatalf("GetServiceName() = %q, want %q", got, ServiceName)
	}
}
