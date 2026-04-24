// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientRPC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pass, ok := r.BasicAuth(); !ok || user != "x" || pass != "testkey" {
			t.Fatalf("unexpected auth: %q %q %v", user, pass, ok)
		}

		var request struct {
			Method string `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		switch request.Method {
		case "getblockchaininfo":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{"blocks": 42, "treeroot": "abc123"},
			})
		case "getinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{"blocks": 42, "connections": 3},
			})
		case "getnameresource":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"records": []map[string]any{
						{"type": "GLUE4", "ns": "charon.lthn.", "address": "192.0.2.10"},
						{"type": "TXT", "txt": []string{"hello"}},
						{"type": "NS", "ns": "ns1.charon.lthn."},
					},
				},
			})
		case "getpeerinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{
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
