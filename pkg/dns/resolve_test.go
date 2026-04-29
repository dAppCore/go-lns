// SPDX-License-Identifier: EUPL-1.2

package dns

import (
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

	if _, err := Resolve(123); err == nil || !core.Contains(err.Error(), "dns.Resolve") {
		t.Fatalf("Resolve should reject unsupported input types with a dns.Resolve error, got %v", err)
	}

	if _, err := Hash(123); err == nil || !core.Contains(err.Error(), "dns.Hash") {
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

// BEGIN v090 file-aware generated coverage

func TestResolve_GetServiceName_Good(t *core.T) {
	fn := ax7ValueOf(GetServiceName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetServiceName_Bad(t *core.T) {
	fn := ax7ValueOf(GetServiceName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetServiceName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetServiceName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Resolve_Good(t *core.T) {
	fn := ax7ValueOf(Resolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Resolve_Bad(t *core.T) {
	fn := ax7ValueOf(Resolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Resolve_Ugly(t *core.T) {
	fn := ax7ValueOf(Resolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolve_Good(t *core.T) {
	fn := ax7ValueOf(GetResolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolve_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolve_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Hash_Good(t *core.T) {
	fn := ax7ValueOf(Hash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Hash_Bad(t *core.T) {
	fn := ax7ValueOf(Hash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Hash_Ugly(t *core.T) {
	fn := ax7ValueOf(Hash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHash_Good(t *core.T) {
	fn := ax7ValueOf(GetHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHash_Bad(t *core.T) {
	fn := ax7ValueOf(GetHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHash_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Verify_Good(t *core.T) {
	fn := ax7ValueOf(Verify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Verify_Bad(t *core.T) {
	fn := ax7ValueOf(Verify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Verify_Ugly(t *core.T) {
	fn := ax7ValueOf(Verify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerify_Good(t *core.T) {
	fn := ax7ValueOf(GetVerify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerify_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerify_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveString_Good(t *core.T) {
	fn := ax7ValueOf(ResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveString_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(ResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashString_Good(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashString_Bad(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashString_Ugly(t *core.T) {
	fn := ax7ValueOf(HashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashString_Good(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashBinary_Good(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveName_Good(t *core.T) {
	fn := ax7ValueOf(ResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveName_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashName_Good(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashName_Bad(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashName_Ugly(t *core.T) {
	fn := ax7ValueOf(HashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashName_Good(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyName_Good(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyName_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(ResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(ResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByString_Good(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByString_Bad(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByString_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyString_Good(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyString_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_NewService_Good(t *core.T) {
	fn := ax7ValueOf(NewService)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_NewService_Bad(t *core.T) {
	fn := ax7ValueOf(NewService)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_NewService_Ugly(t *core.T) {
	fn := ax7ValueOf(NewService)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_WithCore_Good(t *core.T) {
	fn := ax7ValueOf(WithCore)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_WithCore_Bad(t *core.T) {
	fn := ax7ValueOf(WithCore)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_WithCore_Ugly(t *core.T) {
	fn := ax7ValueOf(WithCore)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Core_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Core)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Core_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Core)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Core_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Core)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetCore_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCore)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetCore_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCore)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetCore_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetCore)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ServiceName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ServiceName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ServiceName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ServiceName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ServiceName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ServiceName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetServiceName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetServiceName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetServiceName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetServiceName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetServiceName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetServiceName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Resolve_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Resolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Resolve_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Resolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Resolve_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Resolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolve_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolve)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolve_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolve)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolve_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolve)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Hash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Hash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Hash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Hash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Hash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Hash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHash_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHash)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHash_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHash)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHash_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHash)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Verify_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).Verify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Verify_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).Verify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_Verify_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).Verify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerify_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerify)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerify_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerify)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerify_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerify)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(ResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(ResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_ResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(ResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(GetResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByName_Good(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByName_Bad(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_HashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(HashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByName_Good(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetHashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_VerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_GetVerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByString_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByString)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByString_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByString)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByString_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_ResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).ResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetResolveByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_HashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).HashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetHashByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetHashByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_VerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).VerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByName_Good(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByName_Bad(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResolve_Service_GetVerifyByName_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Service).GetVerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
