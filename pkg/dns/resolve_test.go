// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"strings"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

func TestResolve(t *testing.T) {
	svc := NewService()

	got, err := svc.Resolve("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	want := core.SHA3_256([]byte("foo-bar"))
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

	want := core.SHA3_256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("ResolveName returned %x, want %x", got, want)
	}

	if _, err := svc.ResolveName(123); err == nil {
		t.Fatal("ResolveName should reject unsupported input types")
	}
}

func TestServiceGetAliases(t *testing.T) {
	svc := NewService()
	want := core.SHA3_256([]byte("foo-bar"))

	resolveCases := []struct {
		name string
		fn   func(any) (primitives.Hash, error)
	}{
		{name: "GetResolve", fn: svc.GetResolve},
		{name: "GetHash", fn: svc.GetHash},
		{name: "GetResolveString", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveString(name.(string)) }},
		{name: "GetResolveBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveBinary(name.([]byte)) }},
		{name: "GetHashString", fn: func(name any) (primitives.Hash, error) { return svc.GetHashString(name.(string)) }},
		{name: "GetHashBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetHashBinary(name.([]byte)) }},
		{name: "GetResolveName", fn: svc.GetResolveName},
		{name: "GetHashName", fn: svc.GetHashName},
		{name: "GetResolveByString", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveByString(name.(string)) }},
		{name: "GetResolveByBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetResolveByBinary(name.([]byte)) }},
		{name: "GetResolveByName", fn: svc.GetResolveByName},
		{name: "GetHashByString", fn: func(name any) (primitives.Hash, error) { return svc.GetHashByString(name.(string)) }},
		{name: "GetHashByBinary", fn: func(name any) (primitives.Hash, error) { return svc.GetHashByBinary(name.([]byte)) }},
		{name: "GetHashByName", fn: svc.GetHashByName},
	}

	for _, tc := range resolveCases {
		input := any("Foo-Bar.lthn")
		if tc.name == "GetResolveBinary" || tc.name == "GetHashBinary" || tc.name == "GetResolveByBinary" || tc.name == "GetHashByBinary" {
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

	verifyCases := []struct {
		name string
		fn   func(any) bool
	}{
		{name: "GetVerify", fn: svc.GetVerify},
		{name: "GetVerifyString", fn: func(name any) bool { return svc.GetVerifyString(name.(string)) }},
		{name: "GetVerifyBinary", fn: func(name any) bool { return svc.GetVerifyBinary(name.([]byte)) }},
		{name: "GetVerifyName", fn: svc.GetVerifyName},
		{name: "GetVerifyByString", fn: func(name any) bool { return svc.GetVerifyByString(name.(string)) }},
		{name: "GetVerifyByBinary", fn: func(name any) bool { return svc.GetVerifyByBinary(name.([]byte)) }},
		{name: "GetVerifyByName", fn: svc.GetVerifyByName},
	}

	for _, tc := range verifyCases {
		input := any("Foo-Bar.lthn")
		if tc.name == "GetVerifyBinary" || tc.name == "GetVerifyByBinary" {
			input = []byte("Foo-Bar.lthn")
		}

		if !tc.fn(input) {
			t.Fatalf("%s should accept canonical names", tc.name)
		}
	}
}

func TestPackageGetAliases(t *testing.T) {
	want := core.SHA3_256([]byte("foo-bar"))

	resolveCases := []struct {
		name  string
		input any
		fn    func(any) (primitives.Hash, error)
	}{
		{name: "GetResolve", input: "Foo-Bar.lthn", fn: GetResolve},
		{name: "GetHash", input: "Foo-Bar.lthn", fn: GetHash},
		{name: "GetResolveString", input: "Foo-Bar.lthn", fn: func(name any) (primitives.Hash, error) { return GetResolveString(name.(string)) }},
		{name: "GetResolveBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) (primitives.Hash, error) { return GetResolveBinary(name.([]byte)) }},
		{name: "GetHashString", input: "Foo-Bar.lthn", fn: func(name any) (primitives.Hash, error) { return GetHashString(name.(string)) }},
		{name: "GetHashBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) (primitives.Hash, error) { return GetHashBinary(name.([]byte)) }},
		{name: "GetResolveName", input: "Foo-Bar.lthn", fn: GetResolveName},
		{name: "GetHashName", input: "Foo-Bar.lthn", fn: GetHashName},
		{name: "GetResolveByString", input: "Foo-Bar.lthn", fn: func(name any) (primitives.Hash, error) { return GetResolveByString(name.(string)) }},
		{name: "GetResolveByBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) (primitives.Hash, error) { return GetResolveByBinary(name.([]byte)) }},
		{name: "GetResolveByName", input: "Foo-Bar.lthn", fn: GetResolveByName},
		{name: "GetHashByString", input: "Foo-Bar.lthn", fn: func(name any) (primitives.Hash, error) { return GetHashByString(name.(string)) }},
		{name: "GetHashByBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) (primitives.Hash, error) { return GetHashByBinary(name.([]byte)) }},
		{name: "GetHashByName", input: "Foo-Bar.lthn", fn: GetHashByName},
	}

	for _, tc := range resolveCases {
		got, err := tc.fn(tc.input)
		if err != nil {
			t.Fatalf("%s returned error: %v", tc.name, err)
		}

		if got != want {
			t.Fatalf("%s returned %x, want %x", tc.name, got, want)
		}
	}

	verifyCases := []struct {
		name  string
		input any
		fn    func(any) bool
	}{
		{name: "GetVerify", input: "Foo-Bar.lthn", fn: GetVerify},
		{name: "GetVerifyString", input: "Foo-Bar.lthn", fn: func(name any) bool { return GetVerifyString(name.(string)) }},
		{name: "GetVerifyBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) bool { return GetVerifyBinary(name.([]byte)) }},
		{name: "GetVerifyName", input: "Foo-Bar.lthn", fn: GetVerifyName},
		{name: "GetVerifyByString", input: "Foo-Bar.lthn", fn: func(name any) bool { return GetVerifyByString(name.(string)) }},
		{name: "GetVerifyByBinary", input: []byte("Foo-Bar.lthn"), fn: func(name any) bool { return GetVerifyByBinary(name.([]byte)) }},
		{name: "GetVerifyByName", input: "Foo-Bar.lthn", fn: GetVerifyByName},
	}

	for _, tc := range verifyCases {
		if !tc.fn(tc.input) {
			t.Fatalf("%s should accept canonical names", tc.name)
		}
	}
}

func TestHashAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.Hash("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	want := core.SHA3_256([]byte("foo-bar"))
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

	want := core.SHA3_256([]byte("foo-bar"))
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

	want := core.SHA3_256([]byte("foo-bar"))
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

	want := core.SHA3_256([]byte("foo-bar"))
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

	want := core.SHA3_256([]byte("foo-bar"))

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

func TestNewServiceNilOptions(t *testing.T) {
	svc := NewService(nil, WithCore(nil))
	if svc == nil {
		t.Fatal("NewService should return a service even when options are nil")
	}

	if svc.Core() != nil {
		t.Fatal("NewService(nil, WithCore(nil)) should preserve a nil Core")
	}
}

func TestWithCoreNilSafety(t *testing.T) {
	var svc *Service
	WithCore(nil)(svc)
	WithCore(&core.Core{})(nil)
}

func TestServiceName(t *testing.T) {
	svc := NewService()

	if got := svc.ServiceName(); got != ServiceName {
		t.Fatalf("ServiceName() = %q, want %q", got, ServiceName)
	}

	if got := svc.GetServiceName(); got != ServiceName {
		t.Fatalf("GetServiceName() = %q, want %q", got, ServiceName)
	}

	if got := GetServiceName(); got != ServiceName {
		t.Fatalf("package GetServiceName() = %q, want %q", got, ServiceName)
	}
}

func TestPackageLevelResolveAndVerifyAliases(t *testing.T) {
	want := core.SHA3_256([]byte("foo-bar"))

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
		{name: "ResolveString", fn: func(name any) (primitives.Hash, error) { return ResolveString(name.(string)) }},
		{name: "ResolveBinary", fn: func(name any) (primitives.Hash, error) { return ResolveBinary(name.([]byte)) }},
		{name: "HashString", fn: func(name any) (primitives.Hash, error) { return HashString(name.(string)) }},
		{name: "HashBinary", fn: func(name any) (primitives.Hash, error) { return HashBinary(name.([]byte)) }},
		{name: "ResolveByString", fn: func(name any) (primitives.Hash, error) { return ResolveByString(name.(string)) }},
		{name: "ResolveByBinary", fn: func(name any) (primitives.Hash, error) { return ResolveByBinary(name.([]byte)) }},
		{name: "HashByString", fn: func(name any) (primitives.Hash, error) { return HashByString(name.(string)) }},
		{name: "HashByBinary", fn: func(name any) (primitives.Hash, error) { return HashByBinary(name.([]byte)) }},
	}

	for _, tc := range resolveCases {
		input := any("Foo-Bar.lthn")
		if tc.name == "ResolveBinary" || tc.name == "HashBinary" || tc.name == "ResolveByBinary" || tc.name == "HashByBinary" {
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

	if _, err := Resolve(123); err == nil || !strings.Contains(err.Error(), "dns.Resolve") {
		t.Fatalf("Resolve should reject unsupported input types with a dns.Resolve error, got %v", err)
	}

	if _, err := Hash(123); err == nil || !strings.Contains(err.Error(), "dns.Hash") {
		t.Fatalf("Hash should reject unsupported input types with a dns.Hash error, got %v", err)
	}

	verifyCases := []struct {
		name string
		fn   func(any) bool
	}{
		{name: "Verify", fn: Verify},
		{name: "VerifyName", fn: VerifyName},
		{name: "VerifyByName", fn: VerifyByName},
		{name: "VerifyString", fn: func(name any) bool { return VerifyString(name.(string)) }},
		{name: "VerifyBinary", fn: func(name any) bool { return VerifyBinary(name.([]byte)) }},
		{name: "VerifyByString", fn: func(name any) bool { return VerifyByString(name.(string)) }},
		{name: "VerifyByBinary", fn: func(name any) bool { return VerifyByBinary(name.([]byte)) }},
	}

	for _, tc := range verifyCases {
		input := any("Foo-Bar.lthn")
		if tc.name == "VerifyBinary" || tc.name == "VerifyByBinary" {
			input = []byte("Foo-Bar.lthn")
		}

		if !tc.fn(input) {
			t.Fatalf("%s should accept canonical names", tc.name)
		}

		if tc.name == "VerifyBinary" || tc.name == "VerifyByBinary" {
			input = []byte("foo.bar.lthn")
		} else {
			input = "foo.bar.lthn"
		}

		if tc.fn(input) {
			t.Fatalf("%s should reject malformed names", tc.name)
		}
	}
}
