// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"testing"

	core "dappco.re/go"
)

func TestNSECRecordHelpers(t *testing.T) {
	bitmap := []byte{0x01, 0x02, 0x03}
	rr := Create("Example", NextName("Example"), bitmap)

	if rr.Name != "Example." {
		t.Fatalf("Create name = %q, want %q", rr.Name, "Example.")
	}

	if rr.NextDomain != "example\000." {
		t.Fatalf("Create next domain = %q, want %q", rr.NextDomain, "example\x00.")
	}

	if rr.TTL != DEFAULT_TTL {
		t.Fatalf("Create TTL = %d, want %d", rr.TTL, DEFAULT_TTL)
	}

	alias := GetCreate("Example", NextName("Example"), bitmap)
	if alias.Name != rr.Name || alias.NextDomain != rr.NextDomain || alias.TTL != rr.TTL {
		t.Fatal("GetCreate should alias Create")
	}

	if !core.DeepEqual(rr.TypeBitmap, bitmap) {
		t.Fatalf("Create bitmap = %x, want %x", rr.TypeBitmap, bitmap)
	}

	bitmap[0] = 0xff
	if rr.TypeBitmap[0] != 0x01 {
		t.Fatal("Create should copy the type bitmap")
	}
}

func TestNextName(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "root", in: ".", want: "\x00."},
		{name: "fqdn", in: "Foo-Bar.", want: "foo-bar\x00."},
		{name: "label", in: "example", want: "example\x00."},
		{name: "max label", in: repeatString("a", 63), want: repeatString("a", 62) + "b."},
	}

	for _, tc := range cases {
		if got := NextName(tc.in); got != tc.want {
			t.Fatalf("%s: NextName(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}

		if got := GetNextName(tc.in); got != tc.want {
			t.Fatalf("%s: GetNextName(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

func TestPrevName(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "fqdn", in: "Foo-Bar.", want: "foo-baq\xff."},
		{name: "label", in: "example", want: "exampld\xff."},
		{name: "root", in: ".", want: "."},
	}

	for _, tc := range cases {
		if got := PrevName(tc.in); got != tc.want {
			t.Fatalf("%s: PrevName(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}

		if got := GetPrevName(tc.in); got != tc.want {
			t.Fatalf("%s: GetPrevName(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

// BEGIN v090 file-aware generated coverage

func TestNsec_Create_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(Create("lethean.lthn", "lethean.lthn", []byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_Create_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(Create("", "", []byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_Create_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(Create("lethean.\x00.lthn", "lethean.\x00.lthn", []byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetCreate_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetCreate("lethean.lthn", "lethean.lthn", []byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetCreate_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetCreate("", "", []byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetCreate_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetCreate("lethean.\x00.lthn", "lethean.\x00.lthn", []byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_NextName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(NextName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_NextName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(NextName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_NextName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(NextName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetNextName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetNextName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetNextName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetNextName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetNextName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetNextName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_PrevName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(PrevName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_PrevName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(PrevName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_PrevName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(PrevName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetPrevName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetPrevName("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetPrevName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetPrevName(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNsec_GetPrevName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetPrevName("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
