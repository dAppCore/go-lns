// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	core "dappco.re/go"
	dnspkg "dappco.re/go/lns/pkg/dns"
)

// txtResource builds a Resource carrying the given TXT entries.
func txtResource(t *core.T, entries ...string) *Resource {
	items := make([]any, 0, len(entries))
	for _, e := range entries {
		items = append(items, e)
	}
	resource, err := resourceFromAny(map[string]any{
		"records": []any{map[string]any{"type": "TXT", "txt": items}},
	})
	core.AssertNoError(t, err)
	return resource
}

func TestActionVPNBehaviour_vpnMetadataFromTXT(t *core.T) {
	// Both ';' and ',' separate metadata fields, so "cap=vpn,relay" yields
	// cap=vpn plus a bare "relay" field that is dropped (no '=').
	meta := vpnMetadataFromTXT([]string{"vpn=1;type=gateway", "cap=vpn,relay"})
	core.AssertEqual(t, "1", meta["vpn"])
	core.AssertEqual(t, "gateway", meta["type"])
	core.AssertEqual(t, "vpn", meta["cap"])

	// Entries without an = sign are skipped.
	none := vpnMetadataFromTXT([]string{"just text", "  "})
	core.AssertLen(t, none, 0)
}

func TestActionVPNBehaviour_cutKeyValue(t *core.T) {
	k, v, ok := cutKeyValue("type=gateway")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "type", k)
	core.AssertEqual(t, "gateway", v)

	_, _, ok = cutKeyValue("novalue")
	core.AssertFalse(t, ok)
}

func TestActionVPNBehaviour_splitMetadataFields(t *core.T) {
	parts := splitMetadataFields("a=1;b=2,c=3")
	core.AssertLen(t, parts, 3)
}

func TestActionVPNBehaviour_vpnGatewayFromRecord_Good(t *core.T) {
	record := dnspkg.NameRecord{
		Name:     "charon",
		Resource: txtResource(t, "vpn=1;type=gateway;cap=vpn"),
		Source:   "resource",
	}
	gateway, ok := vpnGatewayFromRecord(record)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "charon", gateway.Name)
	core.AssertEqual(t, "gateway", gateway.Type)
	core.AssertContains(t, gateway.Capabilities, "vpn")
}

func TestActionVPNBehaviour_vpnGatewayFromRecord_Bad(t *core.T) {
	// A record with no VPN markers is not a gateway.
	record := dnspkg.NameRecord{
		Name:     "trade",
		Resource: txtResource(t, "v=spf1"),
		Source:   "resource",
	}
	_, ok := vpnGatewayFromRecord(record)
	core.AssertFalse(t, ok)
}

func TestActionVPNBehaviour_vpnGatewayFromRecord_Ugly(t *core.T) {
	// txtLooksVPN heuristic alone (no explicit metadata) still flags a gateway.
	record := dnspkg.NameRecord{
		Name:     "gateway",
		Resource: txtResource(t, "wg=base64publickey"),
		Source:   "resource",
	}
	gateway, ok := vpnGatewayFromRecord(record)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, "gateway", gateway.Name)
}

func TestActionVPNBehaviour_vpnRecordsFromResourceMap(t *core.T) {
	records, err := vpnRecordsFromResourceMap(map[string]any{
		"charon.lthn": resourceTXTJSON,
	})
	core.AssertNoError(t, err)
	core.AssertLen(t, records, 1)
	// Zone keys are reduced to bare labels.
	core.AssertEqual(t, "charon", records[0].Name)

	// An unsupported value type propagates the error.
	_, err = vpnRecordsFromResourceMap(42)
	core.AssertError(t, err)
}

func TestActionVPNBehaviour_vpnRecordFromResource(t *core.T) {
	opts := core.NewOptions()
	opts.Set("name", "charon.lthn")
	records, err := vpnRecordFromResource(resourceTXTJSON, opts)
	core.AssertNoError(t, err)
	core.AssertLen(t, records, 1)
	core.AssertEqual(t, "charon", records[0].Name)

	// Missing name is rejected.
	_, err = vpnRecordFromResource(resourceTXTJSON, core.NewOptions())
	core.AssertError(t, err)

	// An invalid resource payload propagates the error.
	_, err = vpnRecordFromResource("{broken", opts)
	core.AssertError(t, err)
}
