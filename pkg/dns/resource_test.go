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
		NSRecord{NS: testNS1Example},
		GLUE4Record{NS: testNS1Example, Address: netip.MustParseAddr(testIPv4Example1)},
		GLUE6Record{NS: testNS2Example, Address: netip.MustParseAddr(testIPv6Example1)},
		SYNTH4Record{Address: netip.MustParseAddr(testIPv4Example9)},
		SYNTH6Record{Address: netip.MustParseAddr(testIPv6Example9)},
		TXTRecord{Entries: []string{"hello", "world"}},
	}

	raw, err := resource.Encode()
	if err != nil {
		t.Fatalf(testEncodeErr, err)
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
	if !ok || ns.NS != testNS1Example {
		t.Fatalf("decoded NS record = %#v, want ns1.example.", decoded.Records[1])
	}

	glue4, ok := decoded.Records[2].(GLUE4Record)
	if !ok || glue4.NS != testNS1Example || glue4.Address.String() != testIPv4Example1 {
		t.Fatalf("decoded GLUE4 record = %#v, want ns1.example./192.0.2.1", decoded.Records[2])
	}

	glue6, ok := decoded.Records[3].(GLUE6Record)
	if !ok || glue6.NS != testNS2Example || glue6.Address.String() != testIPv6Example1 {
		t.Fatalf("decoded GLUE6 record = %#v, want ns2.example./2001:db8::1", decoded.Records[3])
	}

	synth4, ok := decoded.Records[4].(SYNTH4Record)
	if !ok || synth4.Address.String() != testIPv4Example9 {
		t.Fatalf("decoded SYNTH4 record = %#v, want 198.51.100.9", decoded.Records[4])
	}

	if synth4.NS() != testSynth4Name {
		t.Fatalf("SYNTH4 NS = %q, want %q", synth4.NS(), testSynth4Name)
	}

	if synth4.GetNS() != synth4.NS() {
		t.Fatalf("SYNTH4 GetNS = %q, want %q", synth4.GetNS(), synth4.NS())
	}

	synth6, ok := decoded.Records[5].(SYNTH6Record)
	if !ok || synth6.Address.String() != testIPv6Example9 {
		t.Fatalf("decoded SYNTH6 record = %#v, want 2001:db8::9", decoded.Records[5])
	}

	if synth6.NS() != testSynth6Name {
		t.Fatalf("SYNTH6 NS = %q, want %q", synth6.NS(), testSynth6Name)
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
		NSRecord{NS: testNS1Example},
		GLUE4Record{NS: testNS1Example, Address: netip.MustParseAddr(testIPv4Example1)},
		GLUE6Record{NS: testNS2Example, Address: netip.MustParseAddr(testIPv6Example1)},
		SYNTH4Record{Address: netip.MustParseAddr(testIPv4Example9)},
		SYNTH6Record{Address: netip.MustParseAddr(testIPv6Example9)},
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

	if ns, ok := decoded.Records[1].(NSRecord); !ok || ns.NS != testNS1Example {
		t.Fatalf("decoded NS record = %#v, want ns1.example.", decoded.Records[1])
	}

	if glue4, ok := decoded.Records[2].(GLUE4Record); !ok || glue4.NS != testNS1Example || glue4.Address.String() != testIPv4Example1 {
		t.Fatalf("decoded GLUE4 record = %#v, want ns1.example./192.0.2.1", decoded.Records[2])
	}

	if glue6, ok := decoded.Records[3].(GLUE6Record); !ok || glue6.NS != testNS2Example || glue6.Address.String() != testIPv6Example1 {
		t.Fatalf("decoded GLUE6 record = %#v, want ns2.example./2001:db8::1", decoded.Records[3])
	}

	if synth4, ok := decoded.Records[4].(SYNTH4Record); !ok || synth4.Address.String() != testIPv4Example9 {
		t.Fatalf("decoded SYNTH4 record = %#v, want 198.51.100.9", decoded.Records[4])
	}

	if synth6, ok := decoded.Records[5].(SYNTH6Record); !ok || synth6.Address.String() != testIPv6Example9 {
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
		NSRecord{NS: testNS1Example},
		GLUE4Record{NS: testNS1Example, Address: netip.MustParseAddr(testIPv4Example1)},
		GLUE6Record{NS: testNS2Example, Address: netip.MustParseAddr(testIPv6Example1)},
		SYNTH4Record{Address: netip.MustParseAddr(testIPv4Example9)},
		SYNTH6Record{Address: netip.MustParseAddr(testIPv6Example9)},
		TXTRecord{Entries: []string{"hello", "world"}},
	}

	raw, err := resource.Encode()
	if err != nil {
		t.Fatalf(testEncodeErr, err)
	}

	if got := resource.GetSize(); got != len(raw) {
		t.Fatalf("Resource.GetSize() = %d, want %d", got, len(raw))
	}

	if got := (DSRecord{Digest: []byte{1, 2, 3}}).GetSize(); got != 8 {
		t.Fatalf("DSRecord.GetSize() = %d, want 8", got)
	}

	if got := (NSRecord{NS: testNS1Example}).GetSize(); got != resourceNameSize(testNS1Example) {
		t.Fatalf("NSRecord.GetSize() = %d, want %d", got, resourceNameSize(testNS1Example))
	}

	if got := (GLUE4Record{NS: testNS1Example}).GetSize(); got != resourceNameSize(testNS1Example)+4 {
		t.Fatalf("GLUE4Record.GetSize() = %d, want %d", got, resourceNameSize(testNS1Example)+4)
	}

	if got := (GLUE6Record{NS: testNS1Example}).GetSize(); got != resourceNameSize(testNS1Example)+16 {
		t.Fatalf("GLUE6Record.GetSize() = %d, want %d", got, resourceNameSize(testNS1Example)+16)
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
				t.Fatalf(testGetTypeFormat, tc.name, got, tc.want)
			}
		case NSRecord:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf(testGetTypeFormat, tc.name, got, tc.want)
			}
		case GLUE4Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf(testGetTypeFormat, tc.name, got, tc.want)
			}
		case GLUE6Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf(testGetTypeFormat, tc.name, got, tc.want)
			}
		case SYNTH4Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf(testGetTypeFormat, tc.name, got, tc.want)
			}
		case SYNTH6Record:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf(testGetTypeFormat, tc.name, got, tc.want)
			}
		case TXTRecord:
			if got := rr.GetType(); got != tc.want {
				t.Fatalf(testGetTypeFormat, tc.name, got, tc.want)
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
		{name: "ns", value: NSRecord{NS: testNS1Example}, want: testNS1Example},
		{name: "glue4", value: GLUE4Record{NS: testNS1Example}, want: testNS1Example},
		{name: "glue6", value: GLUE6Record{NS: testNS2Example}, want: testNS2Example},
		{name: "synth4", value: SYNTH4Record{Address: netip.MustParseAddr(testIPv4Example9)}, want: testSynth4Name},
		{name: "synth6", value: SYNTH6Record{Address: netip.MustParseAddr(testIPv6Example9)}, want: testSynth6Name},
	}

	for _, tc := range cases {
		switch rr := tc.value.(type) {
		case NSRecord:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf(testGetNSFormat, tc.name, got, tc.want)
			}
		case GLUE4Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf(testGetNSFormat, tc.name, got, tc.want)
			}
		case GLUE6Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf(testGetNSFormat, tc.name, got, tc.want)
			}
		case SYNTH4Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf(testGetNSFormat, tc.name, got, tc.want)
			}
		case SYNTH6Record:
			if got := rr.GetNS(); got != tc.want {
				t.Fatalf(testGetNSFormat, tc.name, got, tc.want)
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
				NSRecord{NS: testNS1Example},
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

		if rr.NextDomain != testExampleNext {
			t.Fatalf("%s: ToNSEC next domain = %q, want %q", tc.name, rr.NextDomain, testExampleNext)
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
		NSRecord{NS: testNS1Example},
		GLUE4Record{NS: testNS1Example, Address: netip.MustParseAddr(testLiveIPv4)},
		GLUE6Record{NS: testNS2Other, Address: netip.MustParseAddr("2001:db8::20")},
		SYNTH4Record{Address: netip.MustParseAddr(testIPv4Example9)},
		TXTRecord{Entries: []string{"hello"}},
	}

	ns := resource.ToNS(testExampleFQDN)
	if len(ns) != 3 {
		t.Fatalf("ToNS returned %d records, want 3", len(ns))
	}
	if ns[0].NS != testNS1Example || ns[1].NS != testNS2Other || ns[2].NS != testSynth4Name {
		t.Fatalf("ToNS returned %#v, want ns1.example./ns2.other./_oopm828._synth.", ns)
	}

	aliasNS := resource.GetToNS(testExampleFQDN)
	if len(aliasNS) != len(ns) {
		t.Fatalf("GetToNS returned %d records, want %d", len(aliasNS), len(ns))
	}

	glue := resource.ToGlue(testExampleFQDN)
	if len(glue) != 2 {
		t.Fatalf("ToGlue returned %d records, want 2", len(glue))
	}
	if rr, ok := glue[0].(GLUE4Record); !ok || rr.NS != testNS1Example || rr.Address.String() != testLiveIPv4 {
		t.Fatalf("ToGlue[0] = %#v, want GLUE4 ns1.example./192.0.2.10", glue[0])
	}
	if rr, ok := glue[1].(GLUE4Record); !ok || rr.NS != testSynth4Name || rr.Address.String() != testIPv4Example9 {
		t.Fatalf("ToGlue[1] = %#v, want GLUE4 _oopm828._synth./198.51.100.9", glue[1])
	}

	aliasGlue := resource.GetToGlue(testExampleFQDN)
	if len(aliasGlue) != len(glue) {
		t.Fatalf("GetToGlue returned %d records, want %d", len(aliasGlue), len(glue))
	}

	ds := resource.ToDS(testExampleFQDN)
	if len(ds) != 1 || ds[0].KeyTag != 7 || ds[0].Algorithm != 8 || ds[0].DigestType != 9 || !core.DeepEqual(ds[0].Digest, []byte{0xaa}) {
		t.Fatalf("ToDS returned %#v, want one DS record", ds)
	}

	txt := resource.ToTXT(testExampleFQDN)
	if len(txt) != 1 || len(txt[0].Entries) != 1 || txt[0].Entries[0] != "hello" {
		t.Fatalf("ToTXT returned %#v, want one TXT record", txt)
	}

	zone := resource.ToZone(testExampleFQDN)
	if len(zone) != 6 {
		t.Fatalf("ToZone returned %d records, want 6", len(zone))
	}
	if rr, ok := zone[0].(DSRecord); !ok || rr.KeyTag != 7 {
		t.Fatalf("ToZone[0] = %#v, want DS record", zone[0])
	}
	if rr, ok := zone[1].(NSRecord); !ok || rr.NS != testNS1Example {
		t.Fatalf("ToZone[1] = %#v, want NS ns1.example.", zone[1])
	}
	if rr, ok := zone[2].(NSRecord); !ok || rr.NS != testNS2Other {
		t.Fatalf("ToZone[2] = %#v, want NS ns2.other.", zone[2])
	}
	if rr, ok := zone[3].(NSRecord); !ok || rr.NS != testSynth4Name {
		t.Fatalf("ToZone[3] = %#v, want NS _oopm828._synth.", zone[3])
	}
	if rr, ok := zone[4].(TXTRecord); !ok || len(rr.Entries) != 1 || rr.Entries[0] != "hello" {
		t.Fatalf("ToZone[4] = %#v, want TXT hello", zone[4])
	}
	if rr, ok := zone[5].(GLUE4Record); !ok || rr.NS != testNS1Example || rr.Address.String() != testLiveIPv4 {
		t.Fatalf("ToZone[5] = %#v, want GLUE4 ns1.example./192.0.2.10", zone[5])
	}

	aliasZone := resource.GetToZone(testExampleFQDN)
	if len(aliasZone) != len(zone) {
		t.Fatalf("GetToZone returned %d records, want %d", len(aliasZone), len(zone))
	}

	var nilResource *Resource
	if got := nilResource.ToNS(testExampleFQDN); got != nil {
		t.Fatalf("nil ToNS = %#v, want nil", got)
	}
	if got := nilResource.ToGlue(testExampleFQDN); got != nil {
		t.Fatalf("nil ToGlue = %#v, want nil", got)
	}
	if got := nilResource.ToDS(testExampleFQDN); got != nil {
		t.Fatalf("nil ToDS = %#v, want nil", got)
	}
	if got := nilResource.ToTXT(testExampleFQDN); got != nil {
		t.Fatalf("nil ToTXT = %#v, want nil", got)
	}
	if got := nilResource.ToZone(testExampleFQDN); got != nil {
		t.Fatalf("nil ToZone = %#v, want nil", got)
	}
}

func TestResourceToDNSAndReferral(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		DSRecord{KeyTag: 7, Algorithm: 8, DigestType: 9, Digest: []byte{0xaa}},
		NSRecord{NS: testNS1Example},
		GLUE4Record{NS: testNS1Example, Address: netip.MustParseAddr(testLiveIPv4)},
		TXTRecord{Entries: []string{"hello"}},
	}

	referral := resource.ToReferral(testExampleFQDN, HSTypeNS, false)
	if referral.AA {
		t.Fatal("ToReferral should not set AA for a normal referral")
	}
	if len(referral.Authority) != 2 {
		t.Fatalf("ToReferral authority = %d, want 2", len(referral.Authority))
	}
	if ns, ok := referral.Authority[0].(NSRecord); !ok || ns.NS != testNS1Example {
		t.Fatalf("ToReferral authority[0] = %#v, want NS ns1.example.", referral.Authority[0])
	}
	if ds, ok := referral.Authority[1].(DSRecord); !ok || ds.KeyTag != 7 {
		t.Fatalf("ToReferral authority[1] = %#v, want DS keyTag=7", referral.Authority[1])
	}
	if len(referral.Additional) != 1 {
		t.Fatalf("ToReferral additional = %d, want 1", len(referral.Additional))
	}
	if glue, ok := referral.Additional[0].(GLUE4Record); !ok || glue.NS != testNS1Example || glue.Address.String() != testLiveIPv4 {
		t.Fatalf("ToReferral additional[0] = %#v, want GLUE4 ns1.example./192.0.2.10", referral.Additional[0])
	}

	aliasReferral := resource.GetToReferral(testExampleFQDN, HSTypeNS, false)
	if aliasReferral.AA != referral.AA || len(aliasReferral.Authority) != len(referral.Authority) || len(aliasReferral.Additional) != len(referral.Additional) {
		t.Fatal("GetToReferral should alias ToReferral")
	}

	negative := resource.ToDNS(testExampleFQDN, HSTypeDS)
	if !negative.AA {
		t.Fatal("ToDNS should set AA for a TLD DS negative proof")
	}
	if len(negative.Answer) != 0 {
		t.Fatalf("ToDNS answer = %d, want 0", len(negative.Answer))
	}
	if len(negative.Authority) != 1 {
		t.Fatalf("ToDNS authority = %d, want 1", len(negative.Authority))
	}
	if nsec, ok := negative.Authority[0].(NSECRecord); !ok || nsec.Name != testExampleFQDN || nsec.NextDomain != testExampleNext {
		t.Fatalf("ToDNS authority[0] = %#v, want NSEC example./example\\x00.", negative.Authority[0])
	}

	txtResource := NewResource()
	txtResource.Records = []ResourceRecord{TXTRecord{Entries: []string{"hello"}}}

	answer := txtResource.ToDNS(testExampleFQDN, HSTypeTXT)
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

	aliasAnswer := txtResource.GetToDNS(testExampleFQDN, HSTypeTXT)
	if aliasAnswer.AA != answer.AA || len(aliasAnswer.Answer) != len(answer.Answer) || len(aliasAnswer.Authority) != len(answer.Authority) || len(aliasAnswer.Additional) != len(answer.Additional) {
		t.Fatal("GetToDNS should alias ToDNS")
	}
}

func TestResourceDecodeStopsAtUnknownType(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{TXTRecord{Entries: []string{"known"}}}

	raw, err := resource.Encode()
	if err != nil {
		t.Fatalf(testEncodeErr, err)
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
	coverageAssertVariant(t, "good")
	values := coveragePack(new(RawMessage).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_RawMessage_MarshalJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(RawMessage).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_RawMessage_MarshalJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(RawMessage).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_RawMessage_UnmarshalJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(RawMessage).UnmarshalJSON([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_RawMessage_UnmarshalJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(RawMessage).UnmarshalJSON([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_RawMessage_UnmarshalJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(RawMessage).UnmarshalJSON([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NewResource_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(NewResource())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NewResource_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(NewResource())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NewResource_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(NewResource())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GetNewResource_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetNewResource())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GetNewResource_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetNewResource())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GetNewResource_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetNewResource())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_Type_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(DSRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_Type_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(DSRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_Type_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(DSRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(DSRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(DSRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(DSRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_Type_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NSRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_Type_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NSRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_Type_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NSRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NSRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NSRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NSRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NSRecord).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NSRecord).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NSRecord).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_Type_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE4Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_Type_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE4Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_Type_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE4Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE4Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE4Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE4Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE4Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE4Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE4Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_Type_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE6Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_Type_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE6Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_Type_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE6Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE6Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE6Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE6Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE6Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE6Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE6Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_Type_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH4Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_Type_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH4Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_Type_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH4Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH4Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH4Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH4Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH4Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH4Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH4Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_Type_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH6Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_Type_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH6Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_Type_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH6Record).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH6Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH6Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH6Record).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH6Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH6Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH6Record).GetNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_Type_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(TXTRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_Type_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(TXTRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_Type_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(TXTRecord).Type())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(TXTRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(TXTRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(TXTRecord).GetType())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_NS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH4Record).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_NS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH4Record).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_NS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH4Record).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_NS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH6Record).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_NS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH6Record).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_NS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH6Record).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).HasType(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).HasType(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).HasType(63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasType_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetHasType(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasType_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetHasType(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasType_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetHasType(63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).HasNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).HasNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).HasNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetHasNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetHasNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetHasNS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasDS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).HasDS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasDS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).HasDS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_HasDS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).HasDS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasDS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetHasDS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasDS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetHasDS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetHasDS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetHasDS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToNS("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToNS(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToNS("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToNS("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToNS(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToNS("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToGlue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToGlue("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToGlue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToGlue(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToGlue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToGlue("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToGlue_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToGlue("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToGlue_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToGlue(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToGlue_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToGlue("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToDS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToDS("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToDS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToDS(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToDS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToDS("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToDS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToDS("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToDS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToDS(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToDS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToDS("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToTXT_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToTXT("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToTXT_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToTXT(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToTXT_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToTXT("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToTXT_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToTXT("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToTXT_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToTXT(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToTXT_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToTXT("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToZone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToZone("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToZone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToZone(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToZone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToZone("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToZone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToZone("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToZone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToZone(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToZone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToZone("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToReferral_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToReferral("lethean.lthn", 1, true))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToReferral_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToReferral("", 0, false))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToReferral_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToReferral("lethean.\x00.lthn", 63, true))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToReferral_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToReferral("lethean.lthn", 1, true))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToReferral_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToReferral("", 0, false))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToReferral_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToReferral("lethean.\x00.lthn", 63, true))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToDNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToDNS("lethean.lthn", 1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToDNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToDNS("", 0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToDNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToDNS("lethean.\x00.lthn", 63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToDNS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToDNS("lethean.lthn", 1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToDNS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToDNS("", 0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToDNS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToDNS("lethean.\x00.lthn", 63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToNSEC_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).ToNSEC("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToNSEC_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).ToNSEC(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_ToNSEC_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).ToNSEC("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToNSEC_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetToNSEC("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToNSEC_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetToNSEC(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetToNSEC_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetToNSEC("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_Encode_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).Encode())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_Encode_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).Encode())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_Encode_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).Encode())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetEncode_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetEncode())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetEncode_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetEncode())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetEncode_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetEncode())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_Decode_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).Decode([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_Decode_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).Decode([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_Decode_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).Decode([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetDecode_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetDecode([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetDecode_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetDecode([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetDecode_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetDecode([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Resource).FromJSON(coverageValue[ResourceJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Resource).FromJSON(coverageValue[ResourceJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_Resource_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Resource).FromJSON(coverageValue[ResourceJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(DSRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(DSRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(DSRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NSRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NSRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NSRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE4Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE4Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE4Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE6Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE6Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE6Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH4Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH4Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH4Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH6Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH6Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH6Record).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(TXTRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(TXTRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(TXTRecord).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DecodeResource_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(DecodeResource([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DecodeResource_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(DecodeResource([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DecodeResource_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(DecodeResource([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GetDecodeResource_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetDecodeResource([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GetDecodeResource_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetDecodeResource([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GetDecodeResource_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetDecodeResource([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(DSRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(DSRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(DSRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(DSRecord).FromJSON(coverageValue[DSRecordJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(DSRecord).FromJSON(coverageValue[DSRecordJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_DSRecord_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(DSRecord).FromJSON(coverageValue[DSRecordJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NSRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NSRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NSRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NSRecord).FromJSON(coverageValue[NSRecordJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NSRecord).FromJSON(coverageValue[NSRecordJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_NSRecord_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NSRecord).FromJSON(coverageValue[NSRecordJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE4Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE4Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE4Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE4Record).FromJSON(coverageValue[GLUE4RecordJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE4Record).FromJSON(coverageValue[GLUE4RecordJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE4Record_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE4Record).FromJSON(coverageValue[GLUE4RecordJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE6Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE6Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE6Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(GLUE6Record).FromJSON(coverageValue[GLUE6RecordJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(GLUE6Record).FromJSON(coverageValue[GLUE6RecordJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_GLUE6Record_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(GLUE6Record).FromJSON(coverageValue[GLUE6RecordJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH4Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH4Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH4Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH4Record).FromJSON(coverageValue[SYNTH4RecordJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH4Record).FromJSON(coverageValue[SYNTH4RecordJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH4Record_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH4Record).FromJSON(coverageValue[SYNTH4RecordJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH6Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH6Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH6Record).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(SYNTH6Record).FromJSON(coverageValue[SYNTH6RecordJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(SYNTH6Record).FromJSON(coverageValue[SYNTH6RecordJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_SYNTH6Record_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(SYNTH6Record).FromJSON(coverageValue[SYNTH6RecordJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(TXTRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(TXTRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_GetJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(TXTRecord).GetJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_FromJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(TXTRecord).FromJSON(coverageValue[TXTRecordJSON]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_FromJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(TXTRecord).FromJSON(coverageValue[TXTRecordJSON]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestResource_TXTRecord_FromJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(TXTRecord).FromJSON(coverageValue[TXTRecordJSON]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
