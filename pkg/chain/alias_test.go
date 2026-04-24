// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAliasClientAndHelpers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"aliases": []map[string]any{
					{"alias": "charon", "address": "LTHN...", "comment": "v=lthn1;hns=charon.lthn"},
					{"alias": "trade", "address": "LTHN...", "comment": "v=lthn1;type=gateway"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewAliasClient(server.URL)
	aliases, err := client.GetAliases()
	if err != nil {
		t.Fatalf("GetAliases returned error: %v", err)
	}

	if len(aliases) != 2 {
		t.Fatalf("GetAliases returned %d aliases, want 2", len(aliases))
	}

	if got := HNSName(aliases[0]); got != "charon" {
		t.Fatalf("HNSName(alias[0]) = %q, want %q", got, "charon")
	}

	if got := HNSName(aliases[1]); got != "trade" {
		t.Fatalf("HNSName(alias[1]) = %q, want %q", got, "trade")
	}

	comment := ParseAliasComment("v=lthn1;type=gateway;cap=vpn,dns")
	if comment["type"] != "gateway" || comment["cap"] != "vpn,dns" {
		t.Fatalf("ParseAliasComment returned %#v", comment)
	}

	names := DiscoverNames(aliases, []string{"fallback"})
	if !reflect.DeepEqual(names, []string{"charon", "trade"}) {
		t.Fatalf("DiscoverNames = %#v, want %#v", names, []string{"charon", "trade"})
	}

	fallback := DiscoverNames(nil, []string{"fallback"})
	if !reflect.DeepEqual(fallback, []string{"fallback"}) {
		t.Fatalf("DiscoverNames fallback = %#v, want %#v", fallback, []string{"fallback"})
	}
}
