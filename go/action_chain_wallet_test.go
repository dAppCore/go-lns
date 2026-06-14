// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
	dnspkg "dappco.re/go/lns/pkg/dns"
	"dappco.re/go/lns/pkg/primitives"
	walletpkg "dappco.re/go/lns/pkg/wallet"
)

func TestChainHeightReturnsBlockchainInfo(t *testing.T) {
	chainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
		}
		readActionTestJSON(t, r.Body, &request)
		if request.Method != "getblockchaininfo" {
			t.Fatalf("unexpected method %q", request.Method)
		}
		writeActionTestJSON(t, w, map[string]any{
			"result": map[string]any{"blocks": 100, "treeroot": "abc123"},
		})
	}))
	defer chainServer.Close()

	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.height").Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
	))
	if !result.OK {
		t.Fatalf("lns.chain.height returned %#v, want OK", result)
	}

	height, ok := result.Value.(chainHeightResult)
	if !ok {
		t.Fatalf("lns.chain.height returned %T, want chainHeightResult", result.Value)
	}
	if height.Height != 100 || height.TreeRoot != "abc123" {
		t.Fatalf("chain height = %#v, want {Height:100 TreeRoot:abc123}", height)
	}
}

func TestChainHeightMissingURL(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.height").Run(context.Background(), core.NewOptions())
	if result.OK {
		t.Fatalf("lns.chain.height succeeded without chain_url, want error")
	}
}

func TestChainStatusReturnsCombinedInfo(t *testing.T) {
	chainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
		}
		readActionTestJSON(t, r.Body, &request)
		switch request.Method {
		case "getinfo":
			writeActionTestJSON(t, w, map[string]any{
				"result": map[string]any{"version": "v1", "network": "main"},
			})
		case "getblockchaininfo":
			writeActionTestJSON(t, w, map[string]any{
				"result": map[string]any{"blocks": 42, "treeroot": "deadbeef"},
			})
		default:
			t.Fatalf("unexpected method %q", request.Method)
		}
	}))
	defer chainServer.Close()

	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.status").Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
	))
	if !result.OK {
		t.Fatalf("lns.chain.status returned %#v, want OK", result)
	}

	status, ok := result.Value.(chainStatusResult)
	if !ok {
		t.Fatalf("lns.chain.status returned %T, want chainStatusResult", result.Value)
	}
	if status.BlockchainInfo == nil || status.BlockchainInfo.Blocks != 42 {
		t.Fatalf("chain status = %#v, want blockchain info with blocks=42", status)
	}
	if status.Info == nil || status.Info.Version != "v1" {
		t.Fatalf("chain status = %#v, want info with version=v1", status)
	}
}

func TestChainPeersReturnsPeerInfo(t *testing.T) {
	chainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
		}
		readActionTestJSON(t, r.Body, &request)
		if request.Method != "getpeerinfo" {
			t.Fatalf("unexpected method %q", request.Method)
		}
		writeActionTestJSON(t, w, map[string]any{
			"result": []any{map[string]any{"addr": "1.2.3.4", "version": "1.0"}},
		})
	}))
	defer chainServer.Close()

	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.peers").Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
	))
	if !result.OK {
		t.Fatalf("lns.chain.peers returned %#v, want OK", result)
	}
}

func TestChainBlockByHash(t *testing.T) {
	chainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
			Params []any  `json:"params"`
		}
		readActionTestJSON(t, r.Body, &request)
		if request.Method != "getblock" {
			t.Fatalf("unexpected method %q", request.Method)
		}
		writeActionTestJSON(t, w, map[string]any{
			"result": map[string]any{"hash": "abc123", "height": 50},
		})
	}))
	defer chainServer.Close()

	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.block").Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "hash", Value: "abc123"},
	))
	if !result.OK {
		t.Fatalf("lns.chain.block returned %#v, want OK", result)
	}
}

func TestChainBlockByHeight(t *testing.T) {
	chainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
			Params []any  `json:"params"`
		}
		readActionTestJSON(t, r.Body, &request)
		if request.Method != "getblockbyheight" {
			t.Fatalf("unexpected method %q", request.Method)
		}
		writeActionTestJSON(t, w, map[string]any{
			"result": map[string]any{"hash": "abc123", "height": 50},
		})
	}))
	defer chainServer.Close()

	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.block").Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "height", Value: uint32(50)},
	))
	if !result.OK {
		t.Fatalf("lns.chain.block returned %#v, want OK", result)
	}
}

func TestChainBlockMissingHashAndHeight(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.block").Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: "http://localhost:0"},
	))
	if result.OK {
		t.Fatalf("lns.chain.block succeeded without hash/height, want error")
	}
}

func TestWalletBalanceRequiresWallet(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.wallet.balance").Run(context.Background(), core.NewOptions())
	if result.OK {
		t.Fatalf("lns.wallet.balance succeeded without wallet, want error")
	}
}

func TestWalletCreateImportExportReceive(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	if result := c.Action("lns.wallet.export").Run(context.Background(), core.NewOptions()); result.OK {
		t.Fatal("lns.wallet.export succeeded without wallet, want error")
	}
	if result := c.Action("lns.wallet.receive").Run(context.Background(), core.NewOptions()); result.OK {
		t.Fatal("lns.wallet.receive succeeded without wallet, want error")
	}

	created := c.Action("lns.wallet.create").Run(context.Background(), core.NewOptions())
	if !created.OK {
		t.Fatalf("lns.wallet.create returned %#v, want OK", created)
	}

	exported := c.Action("lns.wallet.export").Run(context.Background(), core.NewOptions())
	if !exported.OK {
		t.Fatalf("lns.wallet.export returned %#v, want OK", exported)
	}

	current := c.Action("lns.wallet.receive").Run(context.Background(), core.NewOptions())
	if !current.OK {
		t.Fatalf("lns.wallet.receive current returned %#v, want OK", current)
	}
	currentReceive, ok := current.Value.(walletReceiveResult)
	if !ok {
		t.Fatalf("lns.wallet.receive returned %T, want walletReceiveResult", current.Value)
	}

	next := c.Action("lns.wallet.receive").Run(context.Background(), core.NewOptions(
		core.Option{Key: "next", Value: true},
	))
	if !next.OK {
		t.Fatalf("lns.wallet.receive next returned %#v, want OK", next)
	}
	nextReceive, ok := next.Value.(walletReceiveResult)
	if !ok {
		t.Fatalf("lns.wallet.receive next returned %T, want walletReceiveResult", next.Value)
	}
	if nextReceive.Index <= currentReceive.Index {
		t.Fatalf("receive index did not advance: current=%d next=%d", currentReceive.Index, nextReceive.Index)
	}

	if result := c.Action("lns.wallet.import").Run(context.Background(), core.NewOptions()); result.OK {
		t.Fatal("lns.wallet.import succeeded without seed, want error")
	}

	imported := c.Action("lns.wallet.import").Run(context.Background(), core.NewOptions(
		core.Option{Key: "seed", Value: "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"},
	))
	if !imported.OK {
		t.Fatalf("lns.wallet.import returned %#v, want OK", imported)
	}
}

func TestWalletBalanceReturnsBalances(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	wallet, err := walletpkg.ImportSeed("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatalf("ImportSeed: %v", err)
	}
	wallet.Credit(500, true, 10, primitives.Hash{1}, "funding")
	svc.wallet = wallet

	result := c.Action("lns.wallet.balance").Run(context.Background(), core.NewOptions())
	if !result.OK {
		t.Fatalf("lns.wallet.balance returned %#v, want OK", result)
	}

	balance, ok := result.Value.(walletBalanceResult)
	if !ok {
		t.Fatalf("lns.wallet.balance returned %T, want walletBalanceResult", result.Value)
	}
	if balance.Locked != 500 {
		t.Fatalf("wallet balance locked = %d, want 500", balance.Locked)
	}
}

func TestWalletHistoryRequiresWallet(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.wallet.history").Run(context.Background(), core.NewOptions())
	if result.OK {
		t.Fatalf("lns.wallet.history succeeded without wallet, want error")
	}
}

func TestWalletHistoryReturnsEntries(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	wallet, err := walletpkg.ImportSeed("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatalf("ImportSeed: %v", err)
	}
	wallet.Credit(100, false, 10, primitives.Hash{2}, "test deposit")
	svc.wallet = wallet

	result := c.Action("lns.wallet.history").Run(context.Background(), core.NewOptions())
	if !result.OK {
		t.Fatalf("lns.wallet.history returned %#v, want OK", result)
	}

	history, ok := result.Value.([]walletpkg.Entry)
	if !ok {
		t.Fatalf("lns.wallet.history returned %T, want []wallet.Entry", result.Value)
	}
	if len(history) != 1 {
		t.Fatalf("wallet history length = %d, want 1", len(history))
	}
}

func TestChainSyncRequiresChainURL(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.sync").Run(context.Background(), core.NewOptions())
	if result.OK {
		t.Fatalf("lns.chain.sync succeeded without chain_url, want error")
	}
}

func TestChainSyncDiscoversKnownName(t *testing.T) {
	chainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Method string `json:"method"`
			Params []any  `json:"params"`
		}
		readActionTestJSON(t, r.Body, &request)

		switch request.Method {
		case "getblockchaininfo":
			writeActionTestJSON(t, w, map[string]any{
				"result": map[string]any{"blocks": 100, "treeroot": "root-sync"},
			})
		case "getnameresource":
			if len(request.Params) != 1 || request.Params[0] != "charon" {
				t.Fatalf("unexpected getnameresource params %#v", request.Params)
			}
			writeActionTestJSON(t, w, map[string]any{
				"result": map[string]any{
					"records": []map[string]any{{"type": "TXT", "txt": []string{"synced"}}},
				},
			})
		default:
			t.Fatalf("unexpected method %q", request.Method)
		}
	}))
	defer chainServer.Close()

	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup: %v", result)
	}

	result := c.Action("lns.chain.sync").Run(context.Background(), core.NewOptions(
		core.Option{Key: "chain_url", Value: chainServer.URL},
		core.Option{Key: "known_names", Value: []string{"charon"}},
	))
	if !result.OK {
		t.Fatalf("lns.chain.sync returned %#v, want OK", result)
	}

	report, ok := result.Value.(dnspkg.SyncReport)
	if !ok {
		t.Fatalf("lns.chain.sync returned %T, want dns.SyncReport", result.Value)
	}
	if report.TreeRoot != "root-sync" || report.Names != 1 || len(report.Updated) != 1 || report.Updated[0] != "charon" {
		t.Fatalf("sync report = %#v, want charon update at root-sync", report)
	}
}
