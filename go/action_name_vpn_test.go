// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
)

func TestNameLinkedActionsReturnDraftPlans(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}

	resource := NewResource()
	resource.Records = []ResourceRecord{
		TXTRecord{Entries: []string{"hello"}},
	}

	updateResult := c.Action(actionNameUpdate).Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: testCharonName},
		core.Option{Key: "height", Value: 125},
		core.Option{Key: "resource", Value: resource},
	))
	if !updateResult.OK {
		t.Fatalf("lns.name.update returned %#v, want OK", updateResult)
	}

	updatePlan, ok := updateResult.Value.(nameActionPlanResult)
	if !ok {
		t.Fatalf("lns.name.update returned %T, want nameActionPlanResult", updateResult.Value)
	}
	if len(updatePlan.Steps) != 1 || updatePlan.Steps[0].CovenantType != TypeUpdate || updatePlan.ResourceSize == 0 {
		t.Fatalf("update plan = %#v, want one UPDATE step with resource bytes", updatePlan)
	}
	if updatePlan.Steps[0].Draft == nil || len(updatePlan.Steps[0].Draft.Transaction.Outputs) != 1 {
		t.Fatalf("update draft = %#v, want one-output draft transaction", updatePlan.Steps[0].Draft)
	}

	renewResult := c.Action(actionNameRenew).Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: testCharonName},
		core.Option{Key: "height", Value: 125},
		core.Option{Key: "block_hash", Value: "2222222222222222222222222222222222222222222222222222222222222222"},
	))
	if !renewResult.OK {
		t.Fatalf("lns.name.renew returned %#v, want OK", renewResult)
	}

	renewPlan, ok := renewResult.Value.(nameActionPlanResult)
	if !ok {
		t.Fatalf("lns.name.renew returned %T, want nameActionPlanResult", renewResult.Value)
	}
	if len(renewPlan.Steps) != 1 || renewPlan.Steps[0].CovenantType != TypeRenew {
		t.Fatalf("renew plan = %#v, want one RENEW step", renewPlan)
	}
	if renewPlan.Steps[0].Draft == nil || renewPlan.Steps[0].Draft.Fee == 0 {
		t.Fatalf("renew draft = %#v, want fee placeholder", renewPlan.Steps[0].Draft)
	}

	transferResult := c.Action(actionNameTransfer).Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: testCharonName},
		core.Option{Key: "height", Value: 125},
		core.Option{Key: "address", Value: Address{Version: 0, Hash: []byte{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		}}},
	))
	if !transferResult.OK {
		t.Fatalf("lns.name.transfer returned %#v, want OK", transferResult)
	}

	transferPlan, ok := transferResult.Value.(nameActionPlanResult)
	if !ok {
		t.Fatalf("lns.name.transfer returned %T, want nameActionPlanResult", transferResult.Value)
	}
	if len(transferPlan.Steps) != 1 || transferPlan.Steps[0].CovenantType != TypeTransfer || transferPlan.Address == nil {
		t.Fatalf("transfer plan = %#v, want one TRANSFER step and address", transferPlan)
	}
	if transferPlan.Steps[0].Draft == nil || transferPlan.Steps[0].Draft.PrimaryOutputValue != 0 {
		t.Fatalf("transfer draft = %#v, want zero-value placeholder output by default", transferPlan.Steps[0].Draft)
	}

	revokeResult := c.Action(actionNameRevoke).Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: testCharonName},
		core.Option{Key: "height", Value: 125},
	))
	if !revokeResult.OK {
		t.Fatalf("lns.name.revoke returned %#v, want OK", revokeResult)
	}

	revokePlan, ok := revokeResult.Value.(nameActionPlanResult)
	if !ok {
		t.Fatalf("lns.name.revoke returned %T, want nameActionPlanResult", revokeResult.Value)
	}
	if len(revokePlan.Steps) != 1 || revokePlan.Steps[0].CovenantType != TypeRevoke {
		t.Fatalf("revoke plan = %#v, want one REVOKE step", revokePlan)
	}
}

func TestVPNActionsDiscoverAndTrackSession(t *testing.T) {
	aliasServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeActionTestJSON(t, w, map[string]any{
			"result": map[string]any{
				"aliases": []map[string]any{
					{"alias": "charon", "address": "LTHN...", "comment": "v=lthn1;type=gateway;cap=vpn,dns;hns=charon.lthn"},
				},
			},
		})
	}))
	defer aliasServer.Close()

	chainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
		}
		readActionTestJSON(t, r.Body, &request)

		switch request.Method {
		case "getblockchaininfo":
			writeActionTestJSON(t, w, map[string]any{
				"result": map[string]any{"blocks": 50, "treeroot": "root-1"},
			})
		case "getnameresource":
			writeActionTestJSON(t, w, map[string]any{
				"result": map[string]any{
					"records": []map[string]any{
						{"type": "GLUE4", "ns": testCharonFQDN, "address": testLiveIPv4},
						{"type": "TXT", "txt": []string{"vpn=wireguard", "endpoint=198.51.100.20:51820"}},
					},
				},
			})
		default:
			t.Fatalf("unexpected method %q", request.Method)
		}
	}))
	defer chainServer.Close()

	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}

	options := core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "alias_url", Value: aliasServer.URL},
	)

	discoverResult := c.Action(actionNameVPNDiscover).Run(context.Background(), options)
	if !discoverResult.OK {
		t.Fatalf("lns.vpn.discover returned %#v, want OK", discoverResult)
	}

	discovered, ok := discoverResult.Value.(vpnDiscoverResult)
	if !ok {
		t.Fatalf("lns.vpn.discover returned %T, want vpnDiscoverResult", discoverResult.Value)
	}
	if len(discovered.Gateways) != 1 || discovered.Gateways[0].Name != "charon" {
		t.Fatalf("vpn discover = %#v, want one charon gateway", discovered)
	}

	connectResult := c.Action(actionNameVPNConnect).Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "alias_url", Value: aliasServer.URL},
		core.Option{Key: "name", Value: testCharonName},
	))
	if !connectResult.OK {
		t.Fatalf("lns.vpn.connect returned %#v, want OK", connectResult)
	}

	connected, ok := connectResult.Value.(vpnStatusResult)
	if !ok {
		t.Fatalf("lns.vpn.connect returned %T, want vpnStatusResult", connectResult.Value)
	}
	if !connected.Connected || connected.Gateway == nil || connected.Gateway.Name != "charon" {
		t.Fatalf("vpn connect = %#v, want connected charon gateway", connected)
	}

	statusResult := c.Action("lns.vpn.status").Run(context.Background(), core.NewOptions())
	if !statusResult.OK {
		t.Fatalf("lns.vpn.status returned %#v, want OK", statusResult)
	}

	status, ok := statusResult.Value.(vpnStatusResult)
	if !ok {
		t.Fatalf("lns.vpn.status returned %T, want vpnStatusResult", statusResult.Value)
	}
	if !status.Connected || status.Gateway == nil || status.Gateway.Name != "charon" {
		t.Fatalf("vpn status = %#v, want connected charon gateway", status)
	}

	disconnectResult := c.Action("lns.vpn.disconnect").Run(context.Background(), core.NewOptions())
	if !disconnectResult.OK {
		t.Fatalf("lns.vpn.disconnect returned %#v, want OK", disconnectResult)
	}

	disconnected, ok := disconnectResult.Value.(vpnStatusResult)
	if !ok {
		t.Fatalf("lns.vpn.disconnect returned %T, want vpnStatusResult", disconnectResult.Value)
	}
	if disconnected.Connected {
		t.Fatalf("vpn disconnect = %#v, want disconnected state", disconnected)
	}
}
