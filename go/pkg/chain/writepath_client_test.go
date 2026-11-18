// SPDX-License-Identifier: EUPL-1.2

package chain

import (
	"net/http"
	"net/http/httptest"
	"testing"

	core "github.com/dappcore/go"
)

// fakeRPC stands up an httptest server that dispatches on the JSON-RPC method
// and records the last params, so write-path RPC wrappers can be asserted
// without a live daemon.
type fakeRPC struct {
	server     *httptest.Server
	lastMethod string
	lastParams []any
	respond    func(method string) any
	failWith   string
}

func newFakeRPC(t *testing.T, respond func(method string) any) *fakeRPC {
	t.Helper()
	f := &fakeRPC{respond: respond}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := core.ReadAll(r.Body)
		if !raw.OK {
			t.Fatalf("read body: %v", raw.Value)
		}
		var req struct {
			Method string `json:"method"`
			Params []any  `json:"params"`
		}
		if res := core.JSONUnmarshalString(raw.Value.(string), &req); !res.OK {
			t.Fatalf("decode request: %v", res.Value)
		}
		f.lastMethod = req.Method
		f.lastParams = req.Params

		if f.failWith != "" {
			writeRPC(t, w, map[string]any{"error": map[string]any{"code": -1, "message": f.failWith}})
			return
		}
		writeRPC(t, w, map[string]any{"result": f.respond(req.Method)})
	}))
	t.Cleanup(f.server.Close)
	return f
}

func writeRPC(t *testing.T, w http.ResponseWriter, payload any) {
	t.Helper()
	raw := core.JSONMarshal(payload)
	if !raw.OK {
		t.Fatalf("marshal payload: %v", raw.Value)
	}
	if _, err := w.Write(raw.Value.([]byte)); err != nil {
		t.Fatalf("write payload: %v", err)
	}
}

// TestClient_SendRawTransaction_Good asserts the wrapper sends the hex under the
// sendrawtransaction method and returns the daemon txid.
func TestClient_SendRawTransaction_Good(t *testing.T) {
	rpc := newFakeRPC(t, func(string) any { return "deadbeeftxid" })
	client := NewClient(rpc.server.URL)

	txid, err := client.SendRawTransaction("00aa11bb")
	if err != nil {
		t.Fatalf("SendRawTransaction error: %v", err)
	}
	if txid != "deadbeeftxid" {
		t.Fatalf("txid = %q, want deadbeeftxid", txid)
	}
	if rpc.lastMethod != "sendrawtransaction" {
		t.Fatalf("method = %q, want sendrawtransaction", rpc.lastMethod)
	}
	if len(rpc.lastParams) != 1 || rpc.lastParams[0] != "00aa11bb" {
		t.Fatalf("params = %#v, want [00aa11bb]", rpc.lastParams)
	}
}

// TestClient_SendRawTransaction_Bad covers an empty hex and a daemon rejection.
func TestClient_SendRawTransaction_Bad(t *testing.T) {
	rpc := newFakeRPC(t, func(string) any { return "x" })
	client := NewClient(rpc.server.URL)

	if _, err := client.SendRawTransaction("   "); err == nil {
		t.Fatalf("SendRawTransaction accepted empty hex, want error")
	}

	rpc.failWith = "TX rejected"
	if _, err := client.SendRawTransaction("00aa"); err == nil {
		t.Fatalf("SendRawTransaction ignored daemon rejection, want error")
	}
}

// TestClient_SendRawTransaction_Ugly covers a daemon that returns the txid inside
// an object rather than as a bare string.
func TestClient_SendRawTransaction_Ugly(t *testing.T) {
	rpc := newFakeRPC(t, func(string) any { return map[string]any{"hash": "wrappedtxid"} })
	client := NewClient(rpc.server.URL)

	txid, err := client.SendRawTransaction("00aa")
	if err != nil {
		t.Fatalf("SendRawTransaction error: %v", err)
	}
	if txid != "wrappedtxid" {
		t.Fatalf("txid = %q, want wrappedtxid", txid)
	}
}

// TestClient_GetCoinsByAddress_Good decodes a coin set returned by the daemon.
func TestClient_GetCoinsByAddress_Good(t *testing.T) {
	rpc := newFakeRPC(t, func(string) any {
		return []map[string]any{
			{"version": 0, "height": 10, "value": 1234, "address": "ts1qxyz", "coinbase": false,
				"hash": "aa", "index": 1},
		}
	})
	client := NewClient(rpc.server.URL)

	coins, err := client.GetCoinsByAddress("ts1qxyz")
	if err != nil {
		t.Fatalf("GetCoinsByAddress error: %v", err)
	}
	if len(coins) != 1 || coins[0].Value != 1234 || coins[0].Index != 1 {
		t.Fatalf("coins = %#v, want one coin value 1234 index 1", coins)
	}
	if rpc.lastMethod != "getcoinsbyaddress" {
		t.Fatalf("method = %q, want getcoinsbyaddress", rpc.lastMethod)
	}
}

// TestClient_GetCoinsByAddress_Bad covers a missing address.
func TestClient_GetCoinsByAddress_Bad(t *testing.T) {
	rpc := newFakeRPC(t, func(string) any { return []any{} })
	client := NewClient(rpc.server.URL)
	if _, err := client.GetCoinsByAddress(""); err == nil {
		t.Fatalf("GetCoinsByAddress accepted empty address, want error")
	}
}

// TestClient_GetNameOwner_Good decodes the owner outpoint from the wrapped
// getnameinfo payload.
func TestClient_GetNameOwner_Good(t *testing.T) {
	rpc := newFakeRPC(t, func(string) any {
		return map[string]any{
			"info": map[string]any{
				"owner": map[string]any{"hash": "cafe", "index": 3},
			},
		}
	})
	client := NewClient(rpc.server.URL)

	owner, err := client.GetNameOwner("example")
	if err != nil {
		t.Fatalf("GetNameOwner error: %v", err)
	}
	if owner.Hash != "cafe" || owner.Index != 3 {
		t.Fatalf("owner = %#v, want hash cafe index 3", owner)
	}
	if rpc.lastMethod != "getnameinfo" {
		t.Fatalf("method = %q, want getnameinfo", rpc.lastMethod)
	}
}

// TestClient_GetNameOwner_Ugly covers a name with no owner (empty result).
func TestClient_GetNameOwner_Ugly(t *testing.T) {
	rpc := newFakeRPC(t, func(string) any { return map[string]any{"info": map[string]any{}} })
	client := NewClient(rpc.server.URL)

	owner, err := client.GetNameOwner("example")
	if err != nil {
		t.Fatalf("GetNameOwner error: %v", err)
	}
	if owner.Hash != "" {
		t.Fatalf("owner = %#v, want empty owner", owner)
	}
}
