// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/covenant"
	"net/netip"
)

func TestResourceEncodeDecodeRoundTrip(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		DSRecord{
			KeyTag:     1,
			Algorithm:  2,
			DigestType: 3,
			Digest:     []byte{0xaa, 0xbb, 0xcc},
		},
		NSRecord{NS: "ns1.example."},
		GLUE4Record{NS: "ns1.example.", Address: netip.MustParseAddr("192.0.2.1")},
		GLUE6Record{NS: "ns2.example.", Address: netip.MustParseAddr("2001:db8::1")},
		SYNTH4Record{Address: netip.MustParseAddr("198.51.100.9")},
		SYNTH6Record{Address: netip.MustParseAddr("2001:db8::9")},
		TXTRecord{Entries: []string{"hello", "world"}},
	}

	raw, err := resource.Encode()
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	if len(raw) == 0 || raw[0] != 0 {
		t.Fatalf("Encode version byte = %v, want 0-prefixed payload", raw)
	}

	decoded, err := DecodeResource(raw)
	if err != nil {
		t.Fatalf("DecodeResource returned error: %v", err)
	}

	if decoded.TTL != DEFAULT_TTL {
		t.Fatalf("decoded TTL = %d, want %d", decoded.TTL, DEFAULT_TTL)
	}

	if len(decoded.Records) != len(resource.Records) {
		t.Fatalf("decoded %d records, want %d", len(decoded.Records), len(resource.Records))
	}

	ds, ok := decoded.Records[0].(DSRecord)
	if !ok {
		t.Fatalf("decoded first record type = %T, want DSRecord", decoded.Records[0])
	}

	if ds.KeyTag != 1 || ds.Algorithm != 2 || ds.DigestType != 3 || !core.DeepEqual(ds.Digest, []byte{0xaa, 0xbb, 0xcc}) {
		t.Fatalf("decoded DS record = %#v, want keyTag=1 algorithm=2 digestType=3 digest=aabbcc", ds)
	}

	ns, ok := decoded.Records[1].(NSRecord)
	if !ok || ns.NS != "ns1.example." {
		t.Fatalf("decoded NS record = %#v, want ns1.example.", decoded.Records[1])
	}

	glue4, ok := decoded.Records[2].(GLUE4Record)
	if !ok || glue4.NS != "ns1.example." || glue4.Address.String() != "192.0.2.1" {
		t.Fatalf("decoded GLUE4 record = %#v, want ns1.example./192.0.2.1", decoded.Records[2])
	}

	glue6, ok := decoded.Records[3].(GLUE6Record)
	if !ok || glue6.NS != "ns2.example." || glue6.Address.String() != "2001:db8::1" {
		t.Fatalf("decoded GLUE6 record = %#v, want ns2.example./2001:db8::1", decoded.Records[3])
	}

	synth4, ok := decoded.Records[4].(SYNTH4Record)
	if !ok || synth4.Address.String() != "198.51.100.9" {
		t.Fatalf("decoded SYNTH4 record = %#v, want 198.51.100.9", decoded.Records[4])
	}

	if synth4.NS() != "_oopm828._synth." {
		t.Fatalf("SYNTH4 NS = %q, want %q", synth4.NS(), "_oopm828._synth.")
	}

	if synth4.GetNS() != synth4.NS() {
		t.Fatalf("SYNTH4 GetNS = %q, want %q", synth4.GetNS(), synth4.NS())
	}

	synth6, ok := decoded.Records[5].(SYNTH6Record)
	if !ok || synth6.Address.String() != "2001:db8::9" {
		t.Fatalf("decoded SYNTH6 record = %#v, want 2001:db8::9", decoded.Records[5])
	}

	if synth6.NS() != "_400gre00000000000000000014._synth." {
		t.Fatalf("SYNTH6 NS = %q, want %q", synth6.NS(), "_400gre00000000000000000014._synth.")
	}

	if synth6.GetNS() != synth6.NS() {
		t.Fatalf("SYNTH6 GetNS = %q, want %q", synth6.GetNS(), synth6.NS())
	}

	txt, ok := decoded.Records[6].(TXTRecord)
	if !ok || len(txt.Entries) != 2 || txt.Entries[0] != "hello" || txt.Entries[1] != "world" {
		t.Fatalf("decoded TXT record = %#v, want [hello world]", decoded.Records[6])
	}
}

func TestResourceJSONRoundTrip(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		DSRecord{
			KeyTag:     1,
			Algorithm:  2,
			DigestType: 3,
			Digest:     []byte{0xaa, 0xbb, 0xcc},
		},
		NSRecord{NS: "ns1.example."},
		GLUE4Record{NS: "ns1.example.", Address: netip.MustParseAddr("192.0.2.1")},
		GLUE6Record{NS: "ns2.example.", Address: netip.MustParseAddr("2001:db8::1")},
		SYNTH4Record{Address: netip.MustParseAddr("198.51.100.9")},
		SYNTH6Record{Address: netip.MustParseAddr("2001:db8::9")},
		TXTRecord{Entries: []string{"hello", "world"}},
	}

	jsonView := resource.GetJSON()
	if len(jsonView.Records) != len(resource.Records) {
		t.Fatalf("GetJSON returned %d records, want %d", len(jsonView.Records), len(resource.Records))
	}

	var decoded Resource
	if err := decoded.FromJSON(jsonView); err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}

	if decoded.TTL != DEFAULT_TTL {
		t.Fatalf("decoded TTL = %d, want %d", decoded.TTL, DEFAULT_TTL)
	}

	if len(decoded.Records) != len(resource.Records) {
		t.Fatalf("decoded %d records, want %d", len(decoded.Records), len(resource.Records))
	}

	if ds, ok := decoded.Records[0].(DSRecord); !ok || ds.KeyTag != 1 || ds.Algorithm != 2 || ds.DigestType != 3 || !core.DeepEqual(ds.Digest, []byte{0xaa, 0xbb, 0xcc}) {
		t.Fatalf("decoded DS record = %#v, want keyTag=1 algorithm=2 digestType=3 digest=aabbcc", decoded.Records[0])
	}

	if ns, ok := decoded.Records[1].(NSRecord); !ok || ns.NS != "ns1.example." {
		t.Fatalf("decoded NS record = %#v, want ns1.example.", decoded.Records[1])
	}

	if glue4, ok := decoded.Records[2].(GLUE4Record); !ok || glue4.NS != "ns1.example." || glue4.Address.String() != "192.0.2.1" {
		t.Fatalf("decoded GLUE4 record = %#v, want ns1.example./192.0.2.1", decoded.Records[2])
	}

	if glue6, ok := decoded.Records[3].(GLUE6Record); !ok || glue6.NS != "ns2.example." || glue6.Address.String() != "2001:db8::1" {
		t.Fatalf("decoded GLUE6 record = %#v, want ns2.example./2001:db8::1", decoded.Records[3])
	}

	if synth4, ok := decoded.Records[4].(SYNTH4Record); !ok || synth4.Address.String() != "198.51.100.9" {
		t.Fatalf("decoded SYNTH4 record = %#v, want 198.51.100.9", decoded.Records[4])
	}

	if synth6, ok := decoded.Records[5].(SYNTH6Record); !ok || synth6.Address.String() != "2001:db8::9" {
		t.Fatalf("decoded SYNTH6 record = %#v, want 2001:db8::9", decoded.Records[5])
	}

	if txt, ok := decoded.Records[6].(TXTRecord); !ok || len(txt.Entries) != 2 || txt.Entries[0] != "hello" || txt.Entries[1] != "world" {
		t.Fatalf("decoded TXT record = %#v, want [hello world]", decoded.Records[6])
	}
}

func TestResourceSizeHelpers(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		DSRecord{
			KeyTag:     1,
			Algorithm:  2,
			DigestType: 3,
			Digest:     []byte{0xaa, 0xbb, 0xcc},
		},
		NSRecord{NS: "ns1.example."},
		GLUE4Record{NS: "ns1.example.", Address: netip.MustParseAddr("192.0.2.1")},
		GLUE6Record{NS: "ns2.example.", Address: netip.MustParseAddr("2001:db8::1")},
		SYNTH4Record{Address: netip.MustParseAddr("198.51.100.9")},
		SYNTH6Record{Address: netip.MustParseAddr("2001:db8::9")},
		TXTRecord{Entries: []string{"hello", "world"}},
	}

	raw, err := resource.Encode()
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	if got := resource.GetSize(); got != len(raw) {
		t.Fatalf("Resource.GetSize() = %d, want %d", got, len(raw))
	}

	if got := (DSRecord{Digest: []byte{1, 2, 3}}).GetSize(); got != 8 {
		t.Fatalf("DSRecord.GetSize() = %d, want 8", got)
	}

	if got := (NSRecord{NS: "ns1.example."}).GetSize(); got != resourceNameSize("ns1.example.") {
		t.Fatalf("NSRecord.GetSize() = %d, want %d", got, resourceNameSize("ns1.example."))
	}

	if got := (GLUE4Record{NS: "ns1.example."}).GetSize(); got != resourceNameSize("ns1.example.")+4 {
		t.Fatalf("GLUE4Record.GetSize() = %d, want %d", got, resourceNameSize("ns1.example.")+4)
	}

	if got := (GLUE6Record{NS: "ns1.example."}).GetSize(); got != resourceNameSize("ns1.example.")+16 {
		t.Fatalf("GLUE6Record.GetSize() = %d, want %d", got, resourceNameSize("ns1.example.")+16)
	}

	if got := (SYNTH4Record{}).GetSize(); got != 4 {
		t.Fatalf("SYNTH4Record.GetSize() = %d, want 4", got)
	}

	if got := (SYNTH6Record{}).GetSize(); got != 16 {
		t.Fatalf("SYNTH6Record.GetSize() = %d, want 16", got)
	}

	if got := (TXTRecord{Entries: []string{"hello", "world"}}).GetSize(); got != 13 {
		t.Fatalf("TXTRecord.GetSize() = %d, want 13", got)
	}
}

func TestResourceTypeHelpers(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		TXTRecord{Entries: []string{"txt"}},
		SYNTH4Record{Address: netip.MustParseAddr("192.0.2.25")},
		DSRecord{Digest: []byte{1}},
	}

	if !resource.HasType(HSTypeTXT) || !resource.GetHasType(HSTypeTXT) {
		t.Fatal("HasType should report TXT records")
	}

	if resource.HasType(HSTypeGLUE6) {
		t.Fatal("HasType should reject absent record types")
	}

	if !resource.HasNS() || !resource.GetHasNS() {
		t.Fatal("HasNS should treat SYNTH records as NS-capable")
	}

	if !resource.HasDS() || !resource.GetHasDS() {
		t.Fatal("HasDS should report DS records")
	}
}

func TestResourceRecordTypeAliases(t *testing.T) {
	cases := []struct {
		name  string
		value ResourceRecord
		want  HSType
	}{
		{name: "ds", value: DSRecord{}, want: HSTypeDS},
		{name: "ns", value: NSRecord{}, want: HSTypeNS},
		{name: "glue4", value: GLUE4Record{}, want: HSTypeGLUE4},
		{name: "glue6", value: GLUE6Record{}, want: HSTypeGLUE6},
		{name: "synth4", value: SYNTH4Record{}, want: HSTypeSYNTH4},
		{name: "synth6", value: SYNTH6Record{}, want: HSTypeSYNTH6},
		{name: "txt", value: TXTRecord{}, want: HSTypeTXT},
	}

	for _, tc := range cases {
		switch rr := tc.value.(type) {
		case DSRecord:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf("%s GetType() = %d, want %d", tc.name, got, tc.want)
			}
		case NSRecord:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf("%s GetType() = %d, want %d", tc.name, got, tc.want)
			}
		case GLUE4Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf("%s GetType() = %d, want %d", tc.name, got, tc.want)
			}
		case GLUE6Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf("%s GetType() = %d, want %d", tc.name, got, tc.want)
			}
		case SYNTH4Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf("%s GetType() = %d, want %d", tc.name, got, tc.want)
			}
		case SYNTH6Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf("%s GetType() = %d, want %d", tc.name, got, tc.want)
			}
		case TXTRecord:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf("%s GetType() = %d, want %d", tc.name, got, tc.want)
			}
		default:
			t.Fatalf("%s: unexpected record type %T", tc.name, tc.value)
		}
	}
}

func TestResourceRecordNSAliases(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  string
	}{
		{name: "ns", value: NSRecord{NS: "ns1.example."}, want: "ns1.example."},
		{name: "glue4", value: GLUE4Record{NS: "ns1.example."}, want: "ns1.example."},
		{name: "glue6", value: GLUE6Record{NS: "ns2.example."}, want: "ns2.example."},
		{name: "synth4", value: SYNTH4Record{Address: netip.MustParseAddr("198.51.100.9")}, want: "_oopm828._synth."},
		{name: "synth6", value: SYNTH6Record{Address: netip.MustParseAddr("2001:db8::9")}, want: "_400gre00000000000000000014._synth."},
	}

	for _, tc := range cases {
		switch rr := tc.value.(type) {
		case NSRecord:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf("%s GetNS() = %q, want %q", tc.name, got, tc.want)
			}
		case GLUE4Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf("%s GetNS() = %q, want %q", tc.name, got, tc.want)
			}
		case GLUE6Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf("%s GetNS() = %q, want %q", tc.name, got, tc.want)
			}
		case SYNTH4Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf("%s GetNS() = %q, want %q", tc.name, got, tc.want)
			}
		case SYNTH6Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf("%s GetNS() = %q, want %q", tc.name, got, tc.want)
			}
		default:
			t.Fatalf("%s: unexpected record type %T", tc.name, tc.value)
		}
	}
}

func TestResourceToNSEC(t *testing.T) {
	cases := []struct {
		name   string
		record []ResourceRecord
		want   []byte
	}{
		{
			name: "empty",
			want: TYPE_MAP_EMPTY,
		},
		{
			name: "txt",
			record: []ResourceRecord{
				TXTRecord{Entries: []string{"txt"}},
			},
			want: TYPE_MAP_TXT,
		},
		{
			name: "ns",
			record: []ResourceRecord{
				TXTRecord{Entries: []string{"txt"}},
				NSRecord{NS: "ns1.example."},
			},
			want: TYPE_MAP_NS,
		},
	}

	for _, tc := range cases {
		resource := NewResource()
		resource.Records = append(resource.Records, tc.record...)

		rr := resource.ToNSEC("Example")
		if rr.Name != "Example." {
			t.Fatalf("%s: ToNSEC name = %q, want %q", tc.name, rr.Name, "Example.")
		}

		if rr.NextDomain != "example\x00." {
			t.Fatalf("%s: ToNSEC next domain = %q, want %q", tc.name, rr.NextDomain, "example\x00.")
		}

		if rr.TTL != DEFAULT_TTL {
			t.Fatalf("%s: ToNSEC TTL = %d, want %d", tc.name, rr.TTL, DEFAULT_TTL)
		}

		if !core.DeepEqual(rr.TypeBitmap, tc.want) {
			t.Fatalf("%s: ToNSEC bitmap = %x, want %x", tc.name, rr.TypeBitmap, tc.want)
		}

		alias := resource.GetToNSEC("Example")
		if alias.Name != rr.Name || alias.NextDomain != rr.NextDomain || alias.TTL != rr.TTL {
			t.Fatalf("%s: GetToNSEC should alias ToNSEC, got %+v want %+v", tc.name, alias, rr)
		}
	}

	var nilResource *Resource
	rr := nilResource.ToNSEC("Example")
	if !core.DeepEqual(rr.TypeBitmap, TYPE_MAP_EMPTY) {
		t.Fatalf("nil resource ToNSEC bitmap = %x, want %x", rr.TypeBitmap, TYPE_MAP_EMPTY)
	}
}

func TestResourceProjectionHelpers(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		DSRecord{KeyTag: 7, Algorithm: 8, DigestType: 9, Digest: []byte{0xaa}},
		NSRecord{NS: "ns1.example."},
		GLUE4Record{NS: "ns1.example.", Address: netip.MustParseAddr("192.0.2.10")},
		GLUE6Record{NS: "ns2.other.", Address: netip.MustParseAddr("2001:db8::20")},
		SYNTH4Record{Address: netip.MustParseAddr("198.51.100.9")},
		TXTRecord{Entries: []string{"hello"}},
	}

	ns := resource.ToNS("example.")
	if len(ns) != 3 {
		t.Fatalf("ToNS returned %d records, want 3", len(ns))
	}
	if ns[0].NS != "ns1.example." || ns[1].NS != "ns2.other." || ns[2].NS != "_oopm828._synth." {
		t.Fatalf("ToNS returned %#v, want ns1.example./ns2.other./_oopm828._synth.", ns)
	}

	aliasNS := resource.GetToNS("example.")
	if len(aliasNS) != len(ns) {
		t.Fatalf("GetToNS returned %d records, want %d", len(aliasNS), len(ns))
	}

	glue := resource.ToGlue("example.")
	if len(glue) != 2 {
		t.Fatalf("ToGlue returned %d records, want 2", len(glue))
	}
	if rr, ok := glue[0].(GLUE4Record); !ok || rr.NS != "ns1.example." || rr.Address.String() != "192.0.2.10" {
		t.Fatalf("ToGlue[0] = %#v, want GLUE4 ns1.example./192.0.2.10", glue[0])
	}
	if rr, ok := glue[1].(GLUE4Record); !ok || rr.NS != "_oopm828._synth." || rr.Address.String() != "198.51.100.9" {
		t.Fatalf("ToGlue[1] = %#v, want GLUE4 _oopm828._synth./198.51.100.9", glue[1])
	}

	aliasGlue := resource.GetToGlue("example.")
	if len(aliasGlue) != len(glue) {
		t.Fatalf("GetToGlue returned %d records, want %d", len(aliasGlue), len(glue))
	}

	ds := resource.ToDS("example.")
	if len(ds) != 1 || ds[0].KeyTag != 7 || ds[0].Algorithm != 8 || ds[0].DigestType != 9 || !core.DeepEqual(ds[0].Digest, []byte{0xaa}) {
		t.Fatalf("ToDS returned %#v, want one DS record", ds)
	}

	txt := resource.ToTXT("example.")
	if len(txt) != 1 || len(txt[0].Entries) != 1 || txt[0].Entries[0] != "hello" {
		t.Fatalf("ToTXT returned %#v, want one TXT record", txt)
	}

	zone := resource.ToZone("example.")
	if len(zone) != 6 {
		t.Fatalf("ToZone returned %d records, want 6", len(zone))
	}
	if rr, ok := zone[0].(DSRecord); !ok || rr.KeyTag != 7 {
		t.Fatalf("ToZone[0] = %#v, want DS record", zone[0])
	}
	if rr, ok := zone[1].(NSRecord); !ok || rr.NS != "ns1.example." {
		t.Fatalf("ToZone[1] = %#v, want NS ns1.example.", zone[1])
	}
	if rr, ok := zone[2].(NSRecord); !ok || rr.NS != "ns2.other." {
		t.Fatalf("ToZone[2] = %#v, want NS ns2.other.", zone[2])
	}
	if rr, ok := zone[3].(NSRecord); !ok || rr.NS != "_oopm828._synth." {
		t.Fatalf("ToZone[3] = %#v, want NS _oopm828._synth.", zone[3])
	}
	if rr, ok := zone[4].(TXTRecord); !ok || len(rr.Entries) != 1 || rr.Entries[0] != "hello" {
		t.Fatalf("ToZone[4] = %#v, want TXT hello", zone[4])
	}
	if rr, ok := zone[5].(GLUE4Record); !ok || rr.NS != "ns1.example." || rr.Address.String() != "192.0.2.10" {
		t.Fatalf("ToZone[5] = %#v, want GLUE4 ns1.example./192.0.2.10", zone[5])
	}

	aliasZone := resource.GetToZone("example.")
	if len(aliasZone) != len(zone) {
		t.Fatalf("GetToZone returned %d records, want %d", len(aliasZone), len(zone))
	}

	var nilResource *Resource
	if got := nilResource.ToNS("example."); got != nil {
		t.Fatalf("nil ToNS = %#v, want nil", got)
	}
	if got := nilResource.ToGlue("example."); got != nil {
		t.Fatalf("nil ToGlue = %#v, want nil", got)
	}
	if got := nilResource.ToDS("example."); got != nil {
		t.Fatalf("nil ToDS = %#v, want nil", got)
	}
	if got := nilResource.ToTXT("example."); got != nil {
		t.Fatalf("nil ToTXT = %#v, want nil", got)
	}
	if got := nilResource.ToZone("example."); got != nil {
		t.Fatalf("nil ToZone = %#v, want nil", got)
	}
}

func TestResourceToDNSAndReferral(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		DSRecord{KeyTag: 7, Algorithm: 8, DigestType: 9, Digest: []byte{0xaa}},
		NSRecord{NS: "ns1.example."},
		GLUE4Record{NS: "ns1.example.", Address: netip.MustParseAddr("192.0.2.10")},
		TXTRecord{Entries: []string{"hello"}},
	}

	referral := resource.ToReferral("example.", HSTypeNS, false)
	if referral.AA {
		t.Fatal("ToReferral should not set AA for a normal referral")
	}
	if len(referral.Authority) != 2 {
		t.Fatalf("ToReferral authority = %d, want 2", len(referral.Authority))
	}
	if ns, ok := referral.Authority[0].(NSRecord); !ok || ns.NS != "ns1.example." {
		t.Fatalf("ToReferral authority[0] = %#v, want NS ns1.example.", referral.Authority[0])
	}
	if ds, ok := referral.Authority[1].(DSRecord); !ok || ds.KeyTag != 7 {
		t.Fatalf("ToReferral authority[1] = %#v, want DS keyTag=7", referral.Authority[1])
	}
	if len(referral.Additional) != 1 {
		t.Fatalf("ToReferral additional = %d, want 1", len(referral.Additional))
	}
	if glue, ok := referral.Additional[0].(GLUE4Record); !ok || glue.NS != "ns1.example." || glue.Address.String() != "192.0.2.10" {
		t.Fatalf("ToReferral additional[0] = %#v, want GLUE4 ns1.example./192.0.2.10", referral.Additional[0])
	}

	aliasReferral := resource.GetToReferral("example.", HSTypeNS, false)
	if aliasReferral.AA != referral.AA || len(aliasReferral.Authority) != len(referral.Authority) || len(aliasReferral.Additional) != len(referral.Additional) {
		t.Fatal("GetToReferral should alias ToReferral")
	}

	negative := resource.ToDNS("example.", HSTypeDS)
	if !negative.AA {
		t.Fatal("ToDNS should set AA for a TLD DS negative proof")
	}
	if len(negative.Answer) != 0 {
		t.Fatalf("ToDNS answer = %d, want 0", len(negative.Answer))
	}
	if len(negative.Authority) != 1 {
		t.Fatalf("ToDNS authority = %d, want 1", len(negative.Authority))
	}
	if nsec, ok := negative.Authority[0].(NSECRecord); !ok || nsec.Name != "example." || nsec.NextDomain != "example\x00." {
		t.Fatalf("ToDNS authority[0] = %#v, want NSEC example./example\\x00.", negative.Authority[0])
	}

	txtResource := NewResource()
	txtResource.Records = []ResourceRecord{TXTRecord{Entries: []string{"hello"}}}

	answer := txtResource.ToDNS("example.", HSTypeTXT)
	if !answer.AA {
		t.Fatal("ToDNS should set AA for an in-zone TXT answer")
	}
	if len(answer.Answer) != 1 {
		t.Fatalf("ToDNS answer = %d, want 1", len(answer.Answer))
	}
	if txt, ok := answer.Answer[0].(TXTRecord); !ok || len(txt.Entries) != 1 || txt.Entries[0] != "hello" {
		t.Fatalf("ToDNS answer[0] = %#v, want TXT hello", answer.Answer[0])
	}
	if len(answer.Authority) != 0 || len(answer.Additional) != 0 {
		t.Fatalf("ToDNS TXT answer should not populate authority/additional, got %+v", answer)
	}

	aliasAnswer := txtResource.GetToDNS("example.", HSTypeTXT)
	if aliasAnswer.AA != answer.AA || len(aliasAnswer.Answer) != len(answer.Answer) || len(aliasAnswer.Authority) != len(answer.Authority) || len(aliasAnswer.Additional) != len(answer.Additional) {
		t.Fatal("GetToDNS should alias ToDNS")
	}
}

func TestResourceDecodeStopsAtUnknownType(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{TXTRecord{Entries: []string{"known"}}}

	raw, err := resource.Encode()
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	raw = append(raw, 0xff, 0x01, 0x02, 0x03)
	decoded, err := DecodeResource(raw)
	if err != nil {
		t.Fatalf("DecodeResource returned error: %v", err)
	}

	if len(decoded.Records) != 1 {
		t.Fatalf("decoded %d records, want 1", len(decoded.Records))
	}
}

func TestResourceEncodeRejectsOversizePayload(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		TXTRecord{Entries: []string{repeatString("a", covenant.MaxResourceSize)}},
	}

	if _, err := resource.Encode(); err == nil {
		t.Fatal("Encode should reject resource payloads larger than MaxResourceSize")
	}
}

func TestResourceDecodeRejectsInvalidPayloads(t *testing.T) {
	resource := NewResource()

	if err := resource.Decode(nil); err == nil {
		t.Fatal("Decode should reject empty payloads")
	}

	if err := resource.Decode([]byte{1}); err == nil {
		t.Fatal("Decode should reject unknown versions")
	}

	if err := resource.Decode(append([]byte{0, byte(HSTypeTXT), 1, 10}, []byte("short")...)); err == nil {
		t.Fatal("Decode should reject truncated TXT entries")
	}
}

// BEGIN v090 file-aware generated coverage

func TestResource_RawMessage_MarshalJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(RawMessage).MarshalJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_RawMessage_MarshalJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(RawMessage).MarshalJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_RawMessage_MarshalJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(RawMessage).MarshalJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_RawMessage_UnmarshalJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(RawMessage).UnmarshalJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_RawMessage_UnmarshalJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(RawMessage).UnmarshalJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_RawMessage_UnmarshalJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(RawMessage).UnmarshalJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NewResource_Good(t *core.T) {
	fn := ax7ValueOf(NewResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NewResource_Bad(t *core.T) {
	fn := ax7ValueOf(NewResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NewResource_Ugly(t *core.T) {
	fn := ax7ValueOf(NewResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GetNewResource_Good(t *core.T) {
	fn := ax7ValueOf(GetNewResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GetNewResource_Bad(t *core.T) {
	fn := ax7ValueOf(GetNewResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GetNewResource_Ugly(t *core.T) {
	fn := ax7ValueOf(GetNewResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_Type_Good(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).Type)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_Type_Bad(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).Type)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_Type_Ugly(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).Type)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetType_Good(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetType_Bad(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_Type_Good(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).Type)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_Type_Bad(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).Type)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_Type_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).Type)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetType_Good(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetType_Bad(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetNS_Good(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_Type_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).Type)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_Type_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).Type)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_Type_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).Type)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetType_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetType_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetNS_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_Type_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).Type)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_Type_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).Type)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_Type_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).Type)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetType_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetType_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetNS_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_Type_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).Type)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_Type_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).Type)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_Type_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).Type)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetType_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetType_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetNS_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_Type_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).Type)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_Type_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).Type)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_Type_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).Type)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetType_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetType_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetNS_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_Type_Good(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).Type)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_Type_Bad(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).Type)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_Type_Ugly(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).Type)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetType_Good(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetType_Bad(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_NS_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).NS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_NS_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).NS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_NS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).NS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_NS_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).NS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_NS_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).NS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_NS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).NS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasType_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasType_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasType_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasType)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasType_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasType)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasType_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasType)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasDS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasDS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasDS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasDS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_HasDS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).HasDS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasDS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasDS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasDS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasDS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetHasDS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetHasDS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToGlue_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToGlue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToGlue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToGlue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToGlue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToGlue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToGlue_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToGlue)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToGlue_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToGlue)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToGlue_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToGlue)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToDS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToDS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToDS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToDS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToDS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToDS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToDS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToDS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToDS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToDS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToDS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToDS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToTXT_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToTXT)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToTXT_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToTXT)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToTXT_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToTXT)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToTXT_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToTXT)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToTXT_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToTXT)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToTXT_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToTXT)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToZone_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToZone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToZone_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToZone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToZone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToZone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToZone_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToZone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToZone_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToZone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToZone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToZone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToReferral_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToReferral)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToReferral_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToReferral)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToReferral_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToReferral)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToReferral_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToReferral)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToReferral_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToReferral)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToReferral_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToReferral)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToDNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToDNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToDNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToDNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToDNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToDNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToDNS_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToDNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToDNS_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToDNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToDNS_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToDNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToNSEC_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToNSEC)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToNSEC_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToNSEC)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_ToNSEC_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).ToNSEC)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToNSEC_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToNSEC)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToNSEC_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToNSEC)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetToNSEC_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetToNSEC)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_Encode_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).Encode)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_Encode_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).Encode)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_Encode_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).Encode)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetEncode_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetEncode)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetEncode_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetEncode)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetEncode_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetEncode)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_Decode_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).Decode)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_Decode_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).Decode)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_Decode_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).Decode)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetDecode_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetDecode)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetDecode_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetDecode)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetDecode_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetDecode)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(Resource).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(Resource).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_Resource_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Resource).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DecodeResource_Good(t *core.T) {
	fn := ax7ValueOf(DecodeResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DecodeResource_Bad(t *core.T) {
	fn := ax7ValueOf(DecodeResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DecodeResource_Ugly(t *core.T) {
	fn := ax7ValueOf(DecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GetDecodeResource_Good(t *core.T) {
	fn := ax7ValueOf(GetDecodeResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GetDecodeResource_Bad(t *core.T) {
	fn := ax7ValueOf(GetDecodeResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GetDecodeResource_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_DSRecord_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(DSRecord).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_NSRecord_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NSRecord).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE4Record_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE4Record).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_GLUE6Record_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(GLUE6Record).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH4Record_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH4Record).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_SYNTH6Record_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(SYNTH6Record).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestResource_TXTRecord_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(TXTRecord).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
