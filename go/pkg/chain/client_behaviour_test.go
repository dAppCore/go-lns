// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
)

func TestClientBehaviour_RawCall_Good(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		writeTestJSON(t, w, map[string]any{"result": map[string]any{"blocks": 7, "treeroot": "abc"}})
	}))
	defer server.Close()

	client := NewClient(server.URL, WithAPIKey("secret"))
	info, err := client.GetBlockchainInfo()
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint64(7), info.Blocks)
	core.AssertEqual(t, "abc", info.TreeRoot)
	// Basic auth header is set when an API key is configured.
	core.AssertNotEqual(t, "", gotAuth)
}

func TestClientBehaviour_RawCall_Bad(t *testing.T) {
	// Nil-url client is rejected.
	_, err := NewClient("").RawCall("getinfo", nil)
	core.AssertError(t, err)

	// Non-200 status is surfaced.
	notFound := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer notFound.Close()
	_, err = NewClient(notFound.URL).RawCall("getinfo", nil)
	core.AssertError(t, err)

	// RPC-level error object is surfaced.
	rpcErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeTestJSON(t, w, map[string]any{"error": map[string]any{"code": -1, "message": "boom"}})
	}))
	defer rpcErr.Close()
	_, err = NewClient(rpcErr.URL).RawCall("getinfo", nil)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "boom")
}

func TestClientBehaviour_RawCall_Ugly(t *testing.T) {
	// Malformed JSON body fails the decode step.
	badJSON := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("{not json")); err != nil {
			t.Fatalf("write: %v", err)
		}
	}))
	defer badJSON.Close()
	_, err := NewClient(badJSON.URL).RawCall("getinfo", nil)
	core.AssertError(t, err)
}

func TestClientBehaviour_call_decodeError(t *testing.T) {
	// A result that is not the expected struct shape fails the typed decode.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeTestJSON(t, w, map[string]any{"result": "a string, not an object"})
	}))
	defer server.Close()
	_, err := NewClient(server.URL).GetBlockchainInfo()
	core.AssertError(t, err)
}

func TestClientBehaviour_options_nilGuards(t *testing.T) {
	// Option functions tolerate nil clients and nil transports.
	core.AssertNotPanics(t, func() { WithAPIKey("x")(nil) })
	core.AssertNotPanics(t, func() { WithHTTPClient(nil)(NewClient("u")) })
	core.AssertNotPanics(t, func() { WithHTTPClient(&http.Client{})(nil) })
	core.AssertNotPanics(t, func() { WithAliasHTTPClient(nil)(NewAliasClient("u")) })
	core.AssertNotPanics(t, func() { WithAliasHTTPClient(&http.Client{})(nil) })
}

func TestClientBehaviour_jsonraw_roundtrip(t *testing.T) {
	var nilRaw jsonraw
	encoded, err := nilRaw.MarshalJSON()
	core.AssertNoError(t, err)
	core.AssertEqual(t, "null", string(encoded))

	var target jsonraw
	core.AssertNoError(t, target.UnmarshalJSON([]byte(`{"k":1}`)))
	core.AssertEqual(t, `{"k":1}`, string(target))

	// A nil receiver is rejected.
	var nilPtr *jsonraw
	core.AssertError(t, nilPtr.UnmarshalJSON([]byte("1")))
}

func TestClientBehaviour_rpcfailure_Error(t *testing.T) {
	failure := &rpcfailure{Code: -32601, Message: "method not found"}
	core.AssertContains(t, failure.Error(), "method not found")
	core.AssertContains(t, failure.Error(), "-32601")
}

func TestAliasBehaviour_NewAliasClient_url(t *testing.T) {
	// Trailing slashes are stripped and the json_rpc path is appended.
	client := NewAliasClient("http://host:1234/")
	core.AssertEqual(t, "http://host:1234/json_rpc", client.url)

	// An already-suffixed url is left intact.
	suffixed := NewAliasClient("http://host:1234/json_rpc")
	core.AssertEqual(t, "http://host:1234/json_rpc", suffixed.url)

	// Empty url stays empty.
	core.AssertEqual(t, "", NewAliasClient("").url)
}

func TestAliasBehaviour_GetAliases_Good(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeTestJSON(t, w, map[string]any{
			"result": map[string]any{
				"aliases": []map[string]any{
					{"alias": "charon", "address": "LTHN1", "comment": "v=lthn1"},
				},
			},
		})
	}))
	defer server.Close()

	aliases, err := NewAliasClient(server.URL).GetAliases()
	core.AssertNoError(t, err)
	core.AssertLen(t, aliases, 1)
	core.AssertEqual(t, "charon", aliases[0].Name)
}

func TestAliasBehaviour_GetAliases_Bad(t *testing.T) {
	// Missing url is rejected.
	_, err := (&AliasClient{}).GetAliases()
	core.AssertError(t, err)

	// Non-200 status is surfaced.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()
	_, err = NewAliasClient(server.URL).GetAliases()
	core.AssertError(t, err)
}

func TestAliasBehaviour_GetAliases_Ugly(t *testing.T) {
	// Malformed JSON fails the decode step.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("{broken")); err != nil {
			t.Fatalf("write: %v", err)
		}
	}))
	defer server.Close()
	_, err := NewAliasClient(server.URL).GetAliases()
	core.AssertError(t, err)
}
