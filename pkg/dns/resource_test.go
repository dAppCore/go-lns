// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"bytes"
	"strings"
	"testing"

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

	if ds.KeyTag != 1 || ds.Algorithm != 2 || ds.DigestType != 3 || !bytes.Equal(ds.Digest, []byte{0xaa, 0xbb, 0xcc}) {
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

	synth6, ok := decoded.Records[5].(SYNTH6Record)
	if !ok || synth6.Address.String() != "2001:db8::9" {
		t.Fatalf("decoded SYNTH6 record = %#v, want 2001:db8::9", decoded.Records[5])
	}

	if synth6.NS() != "_400gre00000000000000000014._synth." {
		t.Fatalf("SYNTH6 NS = %q, want %q", synth6.NS(), "_400gre00000000000000000014._synth.")
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

	if ds, ok := decoded.Records[0].(DSRecord); !ok || ds.KeyTag != 1 || ds.Algorithm != 2 || ds.DigestType != 3 || !bytes.Equal(ds.Digest, []byte{0xaa, 0xbb, 0xcc}) {
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

		if !bytes.Equal(rr.TypeBitmap, tc.want) {
			t.Fatalf("%s: ToNSEC bitmap = %x, want %x", tc.name, rr.TypeBitmap, tc.want)
		}

		alias := resource.GetToNSEC("Example")
		if alias.Name != rr.Name || alias.NextDomain != rr.NextDomain || alias.TTL != rr.TTL {
			t.Fatalf("%s: GetToNSEC should alias ToNSEC, got %+v want %+v", tc.name, alias, rr)
		}
	}

	var nilResource *Resource
	rr := nilResource.ToNSEC("Example")
	if !bytes.Equal(rr.TypeBitmap, TYPE_MAP_EMPTY) {
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
	if len(ds) != 1 || ds[0].KeyTag != 7 || ds[0].Algorithm != 8 || ds[0].DigestType != 9 || !bytes.Equal(ds[0].Digest, []byte{0xaa}) {
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
		TXTRecord{Entries: []string{strings.Repeat("a", covenant.MaxResourceSize)}},
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
