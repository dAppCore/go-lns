// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"

	core "dappco.re/go/core"
	dnspkg "dappco.re/go/lns/pkg/dns"
)

func TestLiveDNSActionsDiscoverServeAndResolve(t *testing.T) {
	aliasServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
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
			Params []any  `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		switch request.Method {
		case "getblockchaininfo":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{"blocks": 50, "treeroot": "root-1"},
			})
		case "getnameresource":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"records": []map[string]any{
						{"type": "GLUE4", "ns": "charon.lthn.", "address": "192.0.2.10"},
						{"type": "TXT", "txt": []string{"hello"}},
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
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	options := core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "alias_url", Value: aliasServer.URL},
		core.Option{Key: "dns_addr", Value: "127.0.0.1:0"},
	)

	discoverResult := c.Action("lns.dns.discover").Run(context.Background(), options)
	if !discoverResult.OK {
		t.Fatalf("lns.dns.discover returned %#v, want OK", discoverResult)
	}

	report, ok := discoverResult.Value.(dnspkg.SyncReport)
	if !ok {
		t.Fatalf("lns.dns.discover returned %T, want dns.SyncReport", discoverResult.Value)
	}
	if report.Names != 1 || report.TreeRoot != "root-1" {
		t.Fatalf("discover report = %#v, want one cached name at root-1", report)
	}

	resolveResult := c.Action("lns.dns.resolve").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon.lthn"},
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "alias_url", Value: aliasServer.URL},
	))
	if !resolveResult.OK {
		t.Fatalf("lns.dns.resolve returned %#v, want OK", resolveResult)
	}

	resolve, ok := resolveResult.Value.(dnsResolveResult)
	if !ok {
		t.Fatalf("lns.dns.resolve returned %T, want dnsResolveResult", resolveResult.Value)
	}
	if len(resolve.Records4) != 1 || len(resolve.Addresses) != 1 || resolve.Addresses[0] != "192.0.2.10" {
		t.Fatalf("resolve = %#v, want one IPv4 address from live bridge", resolve)
	}

	reverseResult := c.Action("lns.dns.reverse").Run(context.Background(), core.NewOptions(
		core.Option{Key: "ip", Value: "192.0.2.10"},
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "alias_url", Value: aliasServer.URL},
	))
	if !reverseResult.OK {
		t.Fatalf("lns.dns.reverse returned %#v, want OK", reverseResult)
	}

	reverse, ok := reverseResult.Value.(dnsReverseResult)
	if !ok {
		t.Fatalf("lns.dns.reverse returned %T, want dnsReverseResult", reverseResult.Value)
	}
	if len(reverse.Names) != 1 || reverse.Names[0] != "charon.lthn." {
		t.Fatalf("reverse = %#v, want charon.lthn.", reverse)
	}

	serveResult := c.Action("lns.dns.serve").Run(context.Background(), options)
	if !serveResult.OK {
		t.Fatalf("lns.dns.serve returned %#v, want OK", serveResult)
	}

	health, ok := serveResult.Value.(dnspkg.BridgeHealth)
	if !ok {
		t.Fatalf("lns.dns.serve returned %T, want dns.BridgeHealth", serveResult.Value)
	}
	if !health.Live || !health.Serving || health.Names != 1 {
		t.Fatalf("serve health = %#v, want live serving bridge with one name", health)
	}

	healthResult := c.Action("lns.dns.health").Run(context.Background(), options)
	if !healthResult.OK {
		t.Fatalf("lns.dns.health returned %#v, want OK", healthResult)
	}

	dnsHealth, ok := healthResult.Value.(dnsHealthResult)
	if !ok {
		t.Fatalf("lns.dns.health returned %T, want dnsHealthResult", healthResult.Value)
	}
	if !dnsHealth.Live || dnsHealth.Bridge == nil || !dnsHealth.Bridge.Serving {
		t.Fatalf("dns health = %#v, want live serving bridge state", dnsHealth)
	}

	if result := svc.OnShutdown(context.Background()); !result.OK {
		t.Fatalf("OnShutdown returned %#v, want OK", result)
	}
	if bridge := dnsHealth.Bridge; bridge != nil && bridge.DNSAddr == "" {
		t.Fatalf("bridge health = %#v, want DNS listen address", bridge)
	}
}

func TestDNSReverseWithExplicitResourceStillWorks(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	resource := NewResource()
	resource.Records = []ResourceRecord{
		GLUE4Record{NS: "charon.lthn.", Address: netip.MustParseAddr("192.0.2.10")},
	}

	result := c.Action("lns.dns.reverse").Run(context.Background(), core.NewOptions(
		core.Option{Key: "name", Value: "charon.lthn"},
		core.Option{Key: "resource", Value: resource},
		core.Option{Key: "ip", Value: "192.0.2.10"},
	))
	if !result.OK {
		t.Fatalf("lns.dns.reverse returned %#v, want OK", result)
	}
}
