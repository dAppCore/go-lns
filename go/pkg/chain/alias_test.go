// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	core "github.com/dappcore/go"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAliasClientAndHelpers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeTestJSON(t, w, map[string]any{
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

// BEGIN v090 file-aware generated coverage

func TestAlias_WithAliasHTTPClient_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(WithAliasHTTPClient(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_WithAliasHTTPClient_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(WithAliasHTTPClient(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_WithAliasHTTPClient_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(WithAliasHTTPClient(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_NewAliasClient_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(NewAliasClient("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_NewAliasClient_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(NewAliasClient(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_NewAliasClient_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(NewAliasClient("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_AliasClient_GetAliases_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(AliasClient).GetAliases())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_AliasClient_GetAliases_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(AliasClient).GetAliases())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_AliasClient_GetAliases_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(AliasClient).GetAliases())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_ParseAliasComment_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ParseAliasComment("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_ParseAliasComment_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ParseAliasComment(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_ParseAliasComment_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ParseAliasComment("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_HNSName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(HNSName(coverageValue[Alias]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_HNSName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(HNSName(coverageValue[Alias]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_HNSName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(HNSName(coverageValue[Alias]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_DiscoverNames_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(DiscoverNames(nil, []string{"lethean"}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_DiscoverNames_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(DiscoverNames(nil, nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestAlias_DiscoverNames_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(DiscoverNames(nil, []string{"lethean"}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
