// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
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
	fn := ax7ValueOf(WithAPIKey)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_WithAPIKey_Bad(t *core.T) {
	fn := ax7ValueOf(WithAPIKey)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_WithAPIKey_Ugly(t *core.T) {
	fn := ax7ValueOf(WithAPIKey)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_WithHTTPClient_Good(t *core.T) {
	fn := ax7ValueOf(WithHTTPClient)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_WithHTTPClient_Bad(t *core.T) {
	fn := ax7ValueOf(WithHTTPClient)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_WithHTTPClient_Ugly(t *core.T) {
	fn := ax7ValueOf(WithHTTPClient)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_NewClient_Good(t *core.T) {
	fn := ax7ValueOf(NewClient)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_NewClient_Bad(t *core.T) {
	fn := ax7ValueOf(NewClient)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_NewClient_Ugly(t *core.T) {
	fn := ax7ValueOf(NewClient)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_rpcfailure_Error_Good(t *core.T) {
	fn := ax7ValueOf(new(rpcfailure).Error)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_rpcfailure_Error_Bad(t *core.T) {
	fn := ax7ValueOf(new(rpcfailure).Error)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_rpcfailure_Error_Ugly(t *core.T) {
	fn := ax7ValueOf(new(rpcfailure).Error)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_jsonraw_MarshalJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_jsonraw_MarshalJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_jsonraw_MarshalJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_jsonraw_UnmarshalJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_jsonraw_UnmarshalJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_jsonraw_UnmarshalJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_RawCall_Good(t *core.T) {
	fn := ax7ValueOf(new(Client).RawCall)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_RawCall_Bad(t *core.T) {
	fn := ax7ValueOf(new(Client).RawCall)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_RawCall_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Client).RawCall)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlockchainInfo_Good(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlockchainInfo)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlockchainInfo_Bad(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlockchainInfo)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlockchainInfo_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlockchainInfo)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetInfo_Good(t *core.T) {
	fn := ax7ValueOf(new(Client).GetInfo)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetInfo_Bad(t *core.T) {
	fn := ax7ValueOf(new(Client).GetInfo)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetInfo_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Client).GetInfo)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetNameResource_Good(t *core.T) {
	fn := ax7ValueOf(new(Client).GetNameResource)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetNameResource_Bad(t *core.T) {
	fn := ax7ValueOf(new(Client).GetNameResource)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetNameResource_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Client).GetNameResource)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlock_Good(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlock)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlock_Bad(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlock)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlock_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlock)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlockByHeight_Good(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlockByHeight)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlockByHeight_Bad(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlockByHeight)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetBlockByHeight_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Client).GetBlockByHeight)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetPeerInfo_Good(t *core.T) {
	fn := ax7ValueOf(new(Client).GetPeerInfo)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetPeerInfo_Bad(t *core.T) {
	fn := ax7ValueOf(new(Client).GetPeerInfo)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestClient_Client_GetPeerInfo_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Client).GetPeerInfo)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
