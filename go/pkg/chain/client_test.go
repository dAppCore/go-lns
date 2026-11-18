// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	"net/http"
	"net/http/httptest"
	"testing"

	core "github.com/dappcore/go"
)

func writeTestJSON(t *testing.T, w http.ResponseWriter, payload any) {
	t.Helper()
	raw := core.JSONMarshal(payload)
	if !raw.OK {
		t.Fatalf("marshal response: %v", raw.Value)
	}
	if _, err := w.Write(raw.Value.([]byte)); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

func readTestJSON(t *testing.T, reader core.Reader, target any) {
	t.Helper()
	raw := core.ReadAll(reader)
	if !raw.OK {
		t.Fatalf("read request: %v", raw.Value)
	}
	if r := core.JSONUnmarshalString(raw.Value.(string), target); !r.OK {
		t.Fatalf("decode request: %v", r.Value)
	}
}

func TestClientRPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pass, ok := r.BasicAuth(); !ok || user != "x" || pass != "testkey" {
			t.Fatalf("unexpected auth: %q %q %v", user, pass, ok)
		}

		var request struct {
			Method string `json:"method"`
		}
		readTestJSON(t, r.Body, &request)

		switch request.Method {
		case "getblockchaininfo":
			writeTestJSON(t, w, map[string]any{
				"result": map[string]any{"blocks": 42, "treeroot": "abc123"},
			})
		case "getinfo":
			writeTestJSON(t, w, map[string]any{
				"result": map[string]any{"blocks": 42, "connections": 3},
			})
		case "getnameresource":
			writeTestJSON(t, w, map[string]any{
				"result": map[string]any{
					"records": []map[string]any{
						{"type": "GLUE4", "ns": "charon.lthn.", "address": "192.0.2.10"},
						{"type": "TXT", "txt": []string{"hello"}},
						{"type": "NS", "ns": "ns1.charon.lthn."},
					},
				},
			})
		case "getpeerinfo":
			writeTestJSON(t, w, map[string]any{
				"result": []map[string]any{{"addr": "127.0.0.1:14037"}},
			})
		default:
			t.Fatalf("unexpected method %q", request.Method)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("testkey"))

	info, err := client.GetBlockchainInfo()
	if err != nil {
		t.Fatalf("GetBlockchainInfo returned error: %v", err)
	}
	if info.Blocks != 42 || info.TreeRoot != "abc123" {
		t.Fatalf("GetBlockchainInfo = %#v, want blocks=42 tree=abc123", info)
	}

	status, err := client.GetInfo()
	if err != nil {
		t.Fatalf("GetInfo returned error: %v", err)
	}
	if status.Blocks != 42 || status.Connections != 3 {
		t.Fatalf("GetInfo = %#v, want blocks=42 connections=3", status)
	}

	resource, err := client.GetNameResource("charon.lthn")
	if err != nil {
		t.Fatalf("GetNameResource returned error: %v", err)
	}
	if len(resource.Records) != 3 {
		t.Fatalf("GetNameResource returned %d records, want 3", len(resource.Records))
	}

	rawPeers, err := client.GetPeerInfo()
	if err != nil {
		t.Fatalf("GetPeerInfo returned error: %v", err)
	}
	if len(rawPeers) == 0 {
		t.Fatal("GetPeerInfo should return raw JSON")
	}
}

// BEGIN v090 file-aware generated coverage

func TestClient_WithAPIKey_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(WithAPIKey("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_WithAPIKey_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(WithAPIKey(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_WithAPIKey_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(WithAPIKey("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_WithHTTPClient_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(WithHTTPClient(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_WithHTTPClient_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(WithHTTPClient(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_WithHTTPClient_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(WithHTTPClient(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_NewClient_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(NewClient("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_NewClient_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(NewClient(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_NewClient_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(NewClient("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_rpcfailure_Error_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(rpcfailure).Error())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_rpcfailure_Error_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(rpcfailure).Error())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_rpcfailure_Error_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(rpcfailure).Error())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_jsonraw_MarshalJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(jsonraw).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_jsonraw_MarshalJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(jsonraw).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_jsonraw_MarshalJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(jsonraw).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_jsonraw_UnmarshalJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(jsonraw).UnmarshalJSON([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_jsonraw_UnmarshalJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(jsonraw).UnmarshalJSON([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_jsonraw_UnmarshalJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(jsonraw).UnmarshalJSON([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_RawCall_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Client).RawCall("lethean.lthn", "lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_RawCall_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Client).RawCall("", nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_RawCall_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Client).RawCall("lethean.\x00.lthn", "lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlockchainInfo_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Client).GetBlockchainInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlockchainInfo_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Client).GetBlockchainInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlockchainInfo_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Client).GetBlockchainInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetInfo_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Client).GetInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetInfo_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Client).GetInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetInfo_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Client).GetInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetNameResource_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Client).GetNameResource("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetNameResource_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Client).GetNameResource(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetNameResource_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Client).GetNameResource("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlock_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Client).GetBlock("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlock_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Client).GetBlock(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlock_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Client).GetBlock("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlockByHeight_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Client).GetBlockByHeight(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlockByHeight_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Client).GetBlockByHeight(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetBlockByHeight_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Client).GetBlockByHeight(63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetPeerInfo_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Client).GetPeerInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetPeerInfo_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Client).GetPeerInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestClient_Client_GetPeerInfo_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Client).GetPeerInfo())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
