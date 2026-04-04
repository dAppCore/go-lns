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
