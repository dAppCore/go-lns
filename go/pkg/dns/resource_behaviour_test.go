// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"net/netip"
	"testing"

	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
)

func TestResourceBehaviour_ToGlue_pointerAndSynth(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		&GLUE4Record{NS: "ns1.charon.lthn.", Address: netip.MustParseAddr("192.0.2.1")},
		&GLUE6Record{NS: "ns2.charon.lthn.", Address: netip.MustParseAddr("2001:db8::1")},
		SYNTH4Record{Address: netip.MustParseAddr("198.51.100.9")},
		SYNTH6Record{Address: netip.MustParseAddr("2001:db8::9")},
		(*GLUE4Record)(nil), // nil pointer is skipped
	}

	glue := resource.ToGlue("charon.lthn.")
	// 2 in-zone glue (pointer variants) + 2 synth-derived glue records.
	core.AssertLen(t, glue, 4)
	core.AssertEqual(t, glue, resource.GetToGlue("charon.lthn."))

	// A nil receiver returns nil.
	core.AssertNil(t, (*Resource)(nil).ToGlue("charon.lthn."))
}

func TestResourceBehaviour_ToDS_ToTXT_pointer(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		&DSRecord{KeyTag: 1, Digest: []byte{0x01}},
		&TXTRecord{Entries: []string{"hello"}},
		(*DSRecord)(nil),
		(*TXTRecord)(nil),
	}

	core.AssertLen(t, resource.ToDS(""), 1)
	core.AssertLen(t, resource.ToTXT(""), 1)
	core.AssertNil(t, (*Resource)(nil).ToDS(""))
	core.AssertNil(t, (*Resource)(nil).ToTXT(""))
}

func TestResourceBehaviour_ToZone_collapsesNS(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		NSRecord{NS: "ns1.charon.lthn."},
		NSRecord{NS: "ns1.charon.lthn."}, // duplicate target, collapsed
		GLUE4Record{NS: "ns1.charon.lthn.", Address: netip.MustParseAddr("192.0.2.1")},
	}

	zone := resource.ToZone("charon.lthn.")
	core.AssertNotEqual(t, 0, len(zone))
	core.AssertNil(t, (*Resource)(nil).ToZone("charon.lthn."))
}

func TestResourceBehaviour_ToReferral(t *testing.T) {
	resource := NewResource()
	resource.Records = []ResourceRecord{
		NSRecord{NS: "ns1.charon.lthn."},
		GLUE4Record{NS: "ns1.charon.lthn.", Address: netip.MustParseAddr("192.0.2.1")},
	}

	msg := resource.ToReferral("charon.lthn.", HSTypeNS, false)
	core.AssertNotEqual(t, 0, len(msg.Authority))
	core.AssertNotEqual(t, 0, len(msg.Additional))

	// A TLD DS referral is a bad referral and falls through to the NSEC path.
	bad := resource.ToReferral("charon.lthn.", HSTypeDS, true)
	core.AssertTrue(t, bad.AA)
}

func TestResourceBehaviour_resourceFromChainResource_Good(t *testing.T) {
	chainResource := &chainpkg.NameResource{
		Records: []chainpkg.Record{
			{Type: "GLUE4", NS: "ns1.charon", Address: "192.0.2.1"},
			{Type: "GLUE6", NS: "ns2.charon", Address: "2001:db8::1"},
			{Type: "TXT", TXT: []string{"hello"}},
			{Type: "NS", NS: "ns3.charon"},
			{Type: "DS", KeyTag: 1, Algorithm: 2, DigestType: 3, Digest: "aabb"},
		},
	}

	resource, err := resourceFromChainResource(chainResource)
	core.AssertNoError(t, err)
	core.AssertLen(t, resource.Records, 5)

	// A nil chain resource yields an empty resource without error.
	empty, err := resourceFromChainResource(nil)
	core.AssertNoError(t, err)
	core.AssertLen(t, empty.Records, 0)
}

func TestResourceBehaviour_resourceFromChainResource_Bad(t *testing.T) {
	// An unparsable GLUE4 address is rejected.
	_, err := resourceFromChainResource(&chainpkg.NameResource{
		Records: []chainpkg.Record{{Type: "GLUE4", NS: "ns1", Address: "not-an-ip"}},
	})
	core.AssertError(t, err)

	// An invalid DS digest hex is rejected.
	_, err = resourceFromChainResource(&chainpkg.NameResource{
		Records: []chainpkg.Record{{Type: "DS", Digest: "zz"}},
	})
	core.AssertError(t, err)
}
