// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"bytes"
	"strings"
	"testing"
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

	if !bytes.Equal(rr.TypeBitmap, bitmap) {
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
		{name: "max label", in: strings.Repeat("a", 63), want: strings.Repeat("a", 62) + "b."},
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
