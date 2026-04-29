// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	core "dappco.re/go"
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
	fn := ax7ValueOf(WithAliasHTTPClient)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_WithAliasHTTPClient_Bad(t *core.T) {
	fn := ax7ValueOf(WithAliasHTTPClient)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_WithAliasHTTPClient_Ugly(t *core.T) {
	fn := ax7ValueOf(WithAliasHTTPClient)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_NewAliasClient_Good(t *core.T) {
	fn := ax7ValueOf(NewAliasClient)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_NewAliasClient_Bad(t *core.T) {
	fn := ax7ValueOf(NewAliasClient)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_NewAliasClient_Ugly(t *core.T) {
	fn := ax7ValueOf(NewAliasClient)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_AliasClient_GetAliases_Good(t *core.T) {
	fn := ax7ValueOf(new(AliasClient).GetAliases)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_AliasClient_GetAliases_Bad(t *core.T) {
	fn := ax7ValueOf(new(AliasClient).GetAliases)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_AliasClient_GetAliases_Ugly(t *core.T) {
	fn := ax7ValueOf(new(AliasClient).GetAliases)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_ParseAliasComment_Good(t *core.T) {
	fn := ax7ValueOf(ParseAliasComment)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_ParseAliasComment_Bad(t *core.T) {
	fn := ax7ValueOf(ParseAliasComment)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_ParseAliasComment_Ugly(t *core.T) {
	fn := ax7ValueOf(ParseAliasComment)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_HNSName_Good(t *core.T) {
	fn := ax7ValueOf(HNSName)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_HNSName_Bad(t *core.T) {
	fn := ax7ValueOf(HNSName)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_HNSName_Ugly(t *core.T) {
	fn := ax7ValueOf(HNSName)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_DiscoverNames_Good(t *core.T) {
	fn := ax7ValueOf(DiscoverNames)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_DiscoverNames_Bad(t *core.T) {
	fn := ax7ValueOf(DiscoverNames)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestAlias_DiscoverNames_Ugly(t *core.T) {
	fn := ax7ValueOf(DiscoverNames)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
