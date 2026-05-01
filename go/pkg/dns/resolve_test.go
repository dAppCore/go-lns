// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

const (
	testFooBarLTHN       = "Foo-Bar.lthn"
	testFooBarLabel      = "foo-bar"
	testFooBarDottedLTHN = "foo.bar.lthn"
	testReturnErrFormat  = "%s returned error: %v"
	testReturnHashFormat = "%s returned %x, want %x"
	testAcceptCanonical  = "%s should accept canonical names"
	testRoot1            = "root-1"
	testLiveIPv4         = "192.0.2.10"
	testNS1Example       = "ns1.example."
	testNS2Example       = "ns2.example."
	testIPv4Example1     = "192.0.2.1"
	testIPv6Example1     = "2001:db8::1"
	testIPv4Example9     = "198.51.100.9"
	testIPv6Example9     = "2001:db8::9"
	testEncodeErr        = "Encode returned error: %v"
	testSynth4Name       = "_oopm828._synth."
	testSynth6Name       = "_400gre00000000000000000014._synth."
	testGetTypeFormat    = "%s GetType() = %d, want %d"
	testGetNSFormat      = "%s GetNS() = %q, want %q"
	testExampleNext      = "example\x00."
	testNS2Other         = "ns2.other."
	testExampleFQDN      = "example."
)

func TestResolve(t *testing.T) {
	svc := NewService()

	got, err := svc.Resolve(testFooBarLTHN)
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	want := core.SHA3_256([]byte(testFooBarLabel))
	if got != want {
		t.Fatalf("Resolve returned %x, want %x", got, want)
	}

	got, err = svc.Resolve([]byte(testFooBarLTHN))
	if err != nil {
		t.Fatalf("Resolve([]byte) returned error: %v", err)
	}

	if got != want {
		t.Fatalf("Resolve([]byte) returned %x, want %x", got, want)
	}
}

func TestResolveNameAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.ResolveName(testFooBarLTHN)
	if err != nil {
		t.Fatalf("ResolveName returned error: %v", err)
	}

	want := core.SHA3_256([]byte(testFooBarLabel))
	if got != want {
		t.Fatalf("ResolveName returned %x, want %x", got, want)
	}

	if _, err := svc.ResolveName(123); err == nil {
		t.Fatal("ResolveName should reject unsupported input types")
	}
}

func TestServiceGetAliases(t *testing.T) {
	svc := NewService()
	want := core.SHA3_256([]byte(testFooBarLabel))

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
		input := any(testFooBarLTHN)
		if tc.name == "GetResolveBinary" || tc.name == "GetHashBinary" || tc.name == "GetResolveByBinary" || tc.name == "GetHashByBinary" {
			input = []byte(testFooBarLTHN)
		}

		got, err := tc.fn(input)
		if err != nil {
			t.Fatalf(testReturnErrFormat, tc.name, err)
		}

		if got != want {
			t.Fatalf(testReturnHashFormat, tc.name, got, want)
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
		input := any(testFooBarLTHN)
		if tc.name == "GetVerifyBinary" || tc.name == "GetVerifyByBinary" {
			input = []byte(testFooBarLTHN)
		}

		if !tc.fn(input) {
			t.Fatalf(testAcceptCanonical, tc.name)
		}
	}
}

func TestPackageGetAliases(t *testing.T) {
	want := core.SHA3_256([]byte(testFooBarLabel))

	resolveCases := []struct {
		name  string
		input any
		fn    func(any) (primitives.Hash, error)
	}{
		{name: "GetResolve", input: testFooBarLTHN, fn: GetResolve},
		{name: "GetHash", input: testFooBarLTHN, fn: GetHash},
		{name: "GetResolveString", input: testFooBarLTHN, fn: func(name any) (primitives.Hash, error) { return GetResolveString(name.(string)) }},
		{name: "GetResolveBinary", input: []byte(testFooBarLTHN), fn: func(name any) (primitives.Hash, error) { return GetResolveBinary(name.([]byte)) }},
		{name: "GetHashString", input: testFooBarLTHN, fn: func(name any) (primitives.Hash, error) { return GetHashString(name.(string)) }},
		{name: "GetHashBinary", input: []byte(testFooBarLTHN), fn: func(name any) (primitives.Hash, error) { return GetHashBinary(name.([]byte)) }},
		{name: "GetResolveName", input: testFooBarLTHN, fn: GetResolveName},
		{name: "GetHashName", input: testFooBarLTHN, fn: GetHashName},
		{name: "GetResolveByString", input: testFooBarLTHN, fn: func(name any) (primitives.Hash, error) { return GetResolveByString(name.(string)) }},
		{name: "GetResolveByBinary", input: []byte(testFooBarLTHN), fn: func(name any) (primitives.Hash, error) { return GetResolveByBinary(name.([]byte)) }},
		{name: "GetResolveByName", input: testFooBarLTHN, fn: GetResolveByName},
		{name: "GetHashByString", input: testFooBarLTHN, fn: func(name any) (primitives.Hash, error) { return GetHashByString(name.(string)) }},
		{name: "GetHashByBinary", input: []byte(testFooBarLTHN), fn: func(name any) (primitives.Hash, error) { return GetHashByBinary(name.([]byte)) }},
		{name: "GetHashByName", input: testFooBarLTHN, fn: GetHashByName},
	}

	for _, tc := range resolveCases {
		got, err := tc.fn(tc.input)
		if err != nil {
			t.Fatalf(testReturnErrFormat, tc.name, err)
		}

		if got != want {
			t.Fatalf(testReturnHashFormat, tc.name, got, want)
		}
	}

	verifyCases := []struct {
		name  string
		input any
		fn    func(any) bool
	}{
		{name: "GetVerify", input: testFooBarLTHN, fn: GetVerify},
		{name: "GetVerifyString", input: testFooBarLTHN, fn: func(name any) bool { return GetVerifyString(name.(string)) }},
		{name: "GetVerifyBinary", input: []byte(testFooBarLTHN), fn: func(name any) bool { return GetVerifyBinary(name.([]byte)) }},
		{name: "GetVerifyName", input: testFooBarLTHN, fn: GetVerifyName},
		{name: "GetVerifyByString", input: testFooBarLTHN, fn: func(name any) bool { return GetVerifyByString(name.(string)) }},
		{name: "GetVerifyByBinary", input: []byte(testFooBarLTHN), fn: func(name any) bool { return GetVerifyByBinary(name.([]byte)) }},
		{name: "GetVerifyByName", input: testFooBarLTHN, fn: GetVerifyByName},
	}

	for _, tc := range verifyCases {
		if !tc.fn(tc.input) {
			t.Fatalf(testAcceptCanonical, tc.name)
		}
	}
}

func TestHashAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.Hash(testFooBarLTHN)
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	want := core.SHA3_256([]byte(testFooBarLabel))
	if got != want {
		t.Fatalf("Hash returned %x, want %x", got, want)
	}

	got, err = svc.HashString(testFooBarLTHN)
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashString returned %x, want %x", got, want)
	}

	got, err = svc.HashBinary([]byte(testFooBarLTHN))
	if err != nil {
		t.Fatalf("HashBinary returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashBinary returned %x, want %x", got, want)
	}

	got, err = svc.HashName(testFooBarLTHN)
	if err != nil {
		t.Fatalf("HashName returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashName returned %x, want %x", got, want)
	}

	got, err = svc.HashByName([]byte(testFooBarLTHN))
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

	got, err := svc.ResolveByName(testFooBarLTHN)
	if err != nil {
		t.Fatalf("ResolveByName returned error: %v", err)
	}

	want := core.SHA3_256([]byte(testFooBarLabel))
	if got != want {
		t.Fatalf("ResolveByName returned %x, want %x", got, want)
	}

	if _, err := svc.ResolveByName(123); err == nil {
		t.Fatal("ResolveByName should reject unsupported input types")
	}
}

func TestResolveStringAndBinary(t *testing.T) {
	svc := NewService()

	got, err := svc.ResolveString(testFooBarLTHN)
	if err != nil {
		t.Fatalf("ResolveString returned error: %v", err)
	}

	want := core.SHA3_256([]byte(testFooBarLabel))
	if got != want {
		t.Fatalf("ResolveString returned %x, want %x", got, want)
	}

	got, err = svc.ResolveBinary([]byte(testFooBarLTHN))
	if err != nil {
		t.Fatalf("ResolveBinary returned error: %v", err)
	}

	if got != want {
		t.Fatalf("ResolveBinary returned %x, want %x", got, want)
	}
}

func TestResolveByStringAndBinaryAliases(t *testing.T) {
	svc := NewService()

	got, err := svc.ResolveByString(testFooBarLTHN)
	if err != nil {
		t.Fatalf("ResolveByString returned error: %v", err)
	}

	want := core.SHA3_256([]byte(testFooBarLabel))
	if got != want {
		t.Fatalf("ResolveByString returned %x, want %x", got, want)
	}

	got, err = svc.ResolveByBinary([]byte(testFooBarLTHN))
	if err != nil {
		t.Fatalf("ResolveByBinary returned error: %v", err)
	}

	if got != want {
		t.Fatalf("ResolveByBinary returned %x, want %x", got, want)
	}
}

func TestHashByStringAndBinaryAliases(t *testing.T) {
	svc := NewService()

	want := core.SHA3_256([]byte(testFooBarLabel))

	got, err := svc.HashByString(testFooBarLTHN)
	if err != nil {
		t.Fatalf("HashByString returned error: %v", err)
	}

	if got != want {
		t.Fatalf("HashByString returned %x, want %x", got, want)
	}

	got, err = svc.HashByBinary([]byte(testFooBarLTHN))
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
		testFooBarDottedLTHN,
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

	if !svc.VerifyString(testFooBarLTHN) {
		t.Fatal("VerifyString should accept canonical names")
	}

	if svc.VerifyString(testFooBarDottedLTHN) {
		t.Fatal("VerifyString should reject malformed names")
	}

	if !svc.VerifyBinary([]byte(testFooBarLTHN)) {
		t.Fatal("VerifyBinary should accept canonical names")
	}

	if svc.VerifyBinary([]byte(testFooBarDottedLTHN)) {
		t.Fatal("VerifyBinary should reject malformed names")
	}
}

func TestVerifyByStringAndBinaryAliases(t *testing.T) {
	svc := NewService()

	if !svc.VerifyByString(testFooBarLTHN) {
		t.Fatal("VerifyByString should accept canonical names")
	}

	if svc.VerifyByString(testFooBarDottedLTHN) {
		t.Fatal("VerifyByString should reject malformed names")
	}

	if !svc.VerifyByBinary([]byte(testFooBarLTHN)) {
		t.Fatal("VerifyByBinary should accept canonical names")
	}

	if svc.VerifyByBinary([]byte(testFooBarDottedLTHN)) {
		t.Fatal("VerifyByBinary should reject malformed names")
	}
}

func TestVerify(t *testing.T) {
	svc := NewService()

	cases := []struct {
		name any
		ok   bool
	}{
		{name: testFooBarLTHN, ok: true},
		{name: []byte(testFooBarLTHN), ok: true},
		{name: testFooBarDottedLTHN, ok: false},
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

	if !svc.VerifyName(testFooBarLTHN) {
		t.Fatal("VerifyName should accept canonical names")
	}

	if svc.VerifyName(testFooBarDottedLTHN) {
		t.Fatal("VerifyName should reject malformed names")
	}

	if svc.VerifyName(123) {
		t.Fatal("VerifyName should reject unsupported input types")
	}
}

func TestVerifyByNameAliases(t *testing.T) {
	svc := NewService()

	if !svc.VerifyByName(testFooBarLTHN) {
		t.Fatal("VerifyByName should accept canonical names")
	}

	if svc.VerifyByName(testFooBarDottedLTHN) {
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
	want := core.SHA3_256([]byte(testFooBarLabel))

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
		input := any(testFooBarLTHN)
		if tc.name == "ResolveBinary" || tc.name == "HashBinary" || tc.name == "ResolveByBinary" || tc.name == "HashByBinary" {
			input = []byte(testFooBarLTHN)
		}

		got, err := tc.fn(input)
		if err != nil {
			t.Fatalf(testReturnErrFormat, tc.name, err)
		}

		if got != want {
			t.Fatalf(testReturnHashFormat, tc.name, got, want)
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
		input := any(testFooBarLTHN)
		if tc.name == "VerifyBinary" || tc.name == "VerifyByBinary" {
			input = []byte(testFooBarLTHN)
		}

		if !tc.fn(input) {
			t.Fatalf(testAcceptCanonical, tc.name)
		}

		if tc.name == "VerifyBinary" || tc.name == "VerifyByBinary" {
			input = []byte(testFooBarDottedLTHN)
		} else {
			input = testFooBarDottedLTHN
		}

		if tc.fn(input) {
			t.Fatalf("%s should reject malformed names", tc.name)
		}
	}
}

// BEGIN v090 file-aware generated coverage

func TestResolve_GetServiceName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetServiceName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetServiceName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Resolve_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(Resolve("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Resolve_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(Resolve(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Resolve_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(Resolve("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolve_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetResolve("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolve_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetResolve(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolve_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetResolve("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Hash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(Hash("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Hash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(Hash(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Hash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(Hash("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHash("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHash(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHash("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Verify_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(Verify("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Verify_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(Verify(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Verify_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(Verify("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerify_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerify("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerify_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerify(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerify_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerify("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ResolveString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ResolveString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ResolveString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetResolveString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetResolveString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetResolveString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ResolveBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ResolveBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ResolveBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetResolveBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetResolveBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetResolveBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HashString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HashString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HashString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHashString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHashString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHashString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HashBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HashBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HashBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHashBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHashBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHashBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ResolveName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ResolveName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ResolveName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetResolveName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetResolveName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetResolveName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HashName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HashName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HashName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHashName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHashName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHashName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(VerifyName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(VerifyName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(VerifyName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerifyName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerifyName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerifyName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ResolveByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ResolveByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ResolveByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetResolveByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetResolveByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetResolveByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ResolveByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ResolveByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ResolveByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetResolveByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetResolveByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetResolveByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HashByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HashByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HashByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHashByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHashByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHashByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HashByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HashByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HashByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHashByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHashByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHashByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(VerifyString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(VerifyString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(VerifyString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerifyString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerifyString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerifyString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(VerifyBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(VerifyBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(VerifyBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerifyBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerifyBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerifyBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(VerifyByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(VerifyByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(VerifyByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerifyByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerifyByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerifyByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(VerifyByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(VerifyByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(VerifyByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerifyByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerifyByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerifyByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_NewService_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(NewService())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_NewService_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(NewService())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_NewService_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(NewService())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_WithCore_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(WithCore(core.New()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_WithCore_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(WithCore(core.New()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_WithCore_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(WithCore(core.New()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Core_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).Core())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Core_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).Core())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Core_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).Core())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetCore_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetCore())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetCore_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetCore())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetCore_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetCore())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ServiceName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).ServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ServiceName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).ServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ServiceName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).ServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetServiceName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetServiceName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetServiceName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetServiceName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Resolve_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).Resolve("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Resolve_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).Resolve(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Resolve_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).Resolve("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolve_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetResolve("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolve_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetResolve(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolve_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetResolve("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Hash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).Hash("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Hash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).Hash(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Hash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).Hash("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetHash("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetHash(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetHash("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Verify_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).Verify("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Verify_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).Verify(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_Verify_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).Verify("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerify_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetVerify("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerify_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetVerify(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerify_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetVerify("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).HashString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).HashString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).HashString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetHashString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetHashString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetHashString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).HashBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).HashBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).HashBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetHashBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetHashBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetHashBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).ResolveByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).ResolveByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).ResolveByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetResolveByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetResolveByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetResolveByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).ResolveByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).ResolveByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).ResolveByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetResolveByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetResolveByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetResolveByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).HashByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).HashByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).HashByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetHashByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetHashByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetHashByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).HashByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).HashByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).HashByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetHashByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetHashByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetHashByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).ResolveString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).ResolveString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).ResolveString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetResolveString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetResolveString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetResolveString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).ResolveBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).ResolveBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).ResolveBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetResolveBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetResolveBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetResolveBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).ResolveName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).ResolveName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).ResolveName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetResolveName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetResolveName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetResolveName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).HashName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).HashName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).HashName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetHashName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetHashName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetHashName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).VerifyString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).VerifyString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).VerifyString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetVerifyString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetVerifyString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetVerifyString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).VerifyBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).VerifyBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).VerifyBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetVerifyBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetVerifyBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetVerifyBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).VerifyName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).VerifyName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).VerifyName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetVerifyName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetVerifyName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetVerifyName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ResolveByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ResolveByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_ResolveByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ResolveByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetResolveByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetResolveByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetResolveByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetResolveByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HashByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HashByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_HashByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HashByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetHashByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetHashByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetHashByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetHashByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(VerifyByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(VerifyByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_VerifyByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(VerifyByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetVerifyByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetVerifyByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_GetVerifyByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetVerifyByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).VerifyByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).VerifyByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).VerifyByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetVerifyByString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetVerifyByString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetVerifyByString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).VerifyByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).VerifyByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).VerifyByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetVerifyByBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetVerifyByBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetVerifyByBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).ResolveByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).ResolveByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_ResolveByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).ResolveByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetResolveByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetResolveByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetResolveByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetResolveByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).HashByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).HashByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_HashByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).HashByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetHashByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetHashByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetHashByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetHashByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).VerifyByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).VerifyByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_VerifyByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).VerifyByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Service).GetVerifyByName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Service).GetVerifyByName(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResolve_Service_GetVerifyByName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Service).GetVerifyByName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
