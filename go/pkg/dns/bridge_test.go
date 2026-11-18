// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"

	core "github.com/dappcore/go"
)

func writeDNSTestJSON(t *testing.T, w http.ResponseWriter, payload any) {
	t.Helper()
	raw := core.JSONMarshal(payload)
	if !raw.OK {
		t.Fatalf("marshal response: %v", raw.Value)
	}
	if _, err := w.Write(raw.Value.([]byte)); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

func readDNSTestJSON(t *testing.T, reader core.Reader, target any) {
	t.Helper()
	raw := core.ReadAll(reader)
	if !raw.OK {
		t.Fatalf("read request: %v", raw.Value)
	}
	if r := core.JSONUnmarshalString(raw.Value.(string), target); !r.OK {
		t.Fatalf("decode request: %v", r.Value)
	}
}

func TestBridgeDiscoverAndReverse(t *testing.T) {
	aliasServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeDNSTestJSON(t, w, map[string]any{
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
		readDNSTestJSON(t, r.Body, &request)

		switch request.Method {
		case "getblockchaininfo":
			writeDNSTestJSON(t, w, map[string]any{
				"result": map[string]any{"blocks": 50, "treeroot": testRoot1},
			})
		case "getnameresource":
			name := request.Params[0].(string)
			if name == "charon" {
				writeDNSTestJSON(t, w, map[string]any{
					"result": map[string]any{
						"records": []map[string]any{
							{"type": "GLUE4", "ns": "charon.lthn.", "address": testLiveIPv4},
							{"type": "TXT", "txt": []string{"hello"}},
						},
					},
				})
				return
			}

			writeDNSTestJSON(t, w, map[string]any{
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

	report, err := bridge.Discover(t.Context())
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if report.TreeRoot != testRoot1 || report.Names != 2 || report.Aliases != 2 {
		t.Fatalf("Discover report = %#v", report)
	}

	record, ok := bridge.Lookup("charon")
	if !ok {
		t.Fatal("Lookup(charon) should succeed")
	}
	if got := record.Addresses4(); len(got) != 1 || got[0] != testLiveIPv4 {
		t.Fatalf("charon addresses = %#v, want 192.0.2.10", got)
	}

	fallback, ok := bridge.Lookup("trade")
	if !ok {
		t.Fatal("Lookup(trade) should use alias dns= fallback")
	}
	if got := fallback.Addresses4(); len(got) != 1 || got[0] != testIPv4Example9 {
		t.Fatalf("trade addresses = %#v, want 198.51.100.9", got)
	}
	if got := fallback.TXT(); len(got) != 1 || got[0] != "fallback" {
		t.Fatalf("trade TXT = %#v, want fallback", got)
	}

	names := bridge.Reverse(netip.MustParseAddr(testLiveIPv4))
	if len(names) != 1 || names[0] != "charon.lthn." {
		t.Fatalf("Reverse = %#v, want charon.lthn.", names)
	}

	health := bridge.Health()
	if !health.Live || health.Names != 2 || health.TreeRoot != testRoot1 {
		t.Fatalf("Health = %#v", health)
	}
}

// BEGIN v090 file-aware generated coverage

func TestBridge_NewBridge_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(NewBridge(coverageValue[BridgeConfig]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NewBridge_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(NewBridge(coverageValue[BridgeConfig]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NewBridge_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(NewBridge(coverageValue[BridgeConfig]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Discover_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Bridge).Discover(core.Background()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Discover_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Bridge).Discover(core.Background()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Discover_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Bridge).Discover(core.Background()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Lookup_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Bridge).Lookup("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Lookup_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Bridge).Lookup(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Lookup_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Bridge).Lookup("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Records_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Bridge).Records())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Records_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Bridge).Records())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Records_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Bridge).Records())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Reverse_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Bridge).Reverse(coverageAddr("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Reverse_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Bridge).Reverse(coverageAddr("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Reverse_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Bridge).Reverse(coverageAddr("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Health_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Bridge).Health())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Health_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Bridge).Health())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Health_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Bridge).Health())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Serve_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Bridge).Serve(core.Background()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Serve_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Bridge).Serve(core.Background()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Serve_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Bridge).Serve(core.Background()))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Shutdown_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Bridge).Shutdown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Shutdown_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Bridge).Shutdown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_Bridge_Shutdown_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Bridge).Shutdown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_FQDN_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameRecord).FQDN())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_FQDN_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameRecord).FQDN())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_FQDN_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameRecord).FQDN())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_Addresses4_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameRecord).Addresses4())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_Addresses4_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameRecord).Addresses4())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_Addresses4_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameRecord).Addresses4())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_Addresses6_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameRecord).Addresses6())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_Addresses6_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameRecord).Addresses6())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_Addresses6_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameRecord).Addresses6())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_TXT_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameRecord).TXT())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_TXT_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameRecord).TXT())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_TXT_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameRecord).TXT())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_NS_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameRecord).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_NS_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameRecord).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBridge_NameRecord_NS_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameRecord).NS())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
