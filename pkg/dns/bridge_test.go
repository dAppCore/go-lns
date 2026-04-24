// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
)

func TestBridgeDiscoverAndReverse(t *testing.T) {
	aliasServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"aliases": []map[string]any{
					{"alias": "charon", "address": "LTHN...", "comment": "v=lthn1;hns=charon.lthn"},
					{"alias": "trade", "address": "LTHN...", "comment": "v=lthn1;dns=A:@:198.51.100.9|TXT:@:fallback"},
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
			name := request.Params[0].(string)
			if name == "charon" {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"result": map[string]any{
						"records": []map[string]any{
							{"type": "GLUE4", "ns": "charon.lthn.", "address": "192.0.2.10"},
							{"type": "TXT", "txt": []string{"hello"}},
						},
					},
				})
				return
			}

			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{"records": []map[string]any{}},
			})
		default:
			t.Fatalf("unexpected method %q", request.Method)
		}
	}))
	defer chainServer.Close()

	bridge := NewBridge(BridgeConfig{
		AliasURL: aliasServer.URL,
		ChainURL: chainServer.URL,
		APIKey:   "testkey",
	})

	report, err := bridge.Discover(context.Background())
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if report.TreeRoot != "root-1" || report.Names != 2 || report.Aliases != 2 {
		t.Fatalf("Discover report = %#v", report)
	}

	record, ok := bridge.Lookup("charon")
	if !ok {
		t.Fatal("Lookup(charon) should succeed")
	}
	if got := record.Addresses4(); len(got) != 1 || got[0] != "192.0.2.10" {
		t.Fatalf("charon addresses = %#v, want 192.0.2.10", got)
	}

	fallback, ok := bridge.Lookup("trade")
	if !ok {
		t.Fatal("Lookup(trade) should use alias dns= fallback")
	}
	if got := fallback.Addresses4(); len(got) != 1 || got[0] != "198.51.100.9" {
		t.Fatalf("trade addresses = %#v, want 198.51.100.9", got)
	}
	if got := fallback.TXT(); len(got) != 1 || got[0] != "fallback" {
		t.Fatalf("trade TXT = %#v, want fallback", got)
	}

	names := bridge.Reverse(netip.MustParseAddr("192.0.2.10"))
	if len(names) != 1 || names[0] != "charon.lthn." {
		t.Fatalf("Reverse = %#v, want charon.lthn.", names)
	}

	health := bridge.Health()
	if !health.Live || health.Names != 2 || health.TreeRoot != "root-1" {
		t.Fatalf("Health = %#v", health)
	}
}
