// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"encoding/hex"
	"testing"

	core "github.com/dappcore/go"
	chainpkg "github.com/dappcore/go-lns/go/pkg/chain"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

const testRegisterName = "charontest"

// coinOption renders a caller-supplied coin (held by the wallet receive address)
// as the option map the action layer parses.
func coinOption(wallet interface {
	ReceiveAddress() primitives.Address
}, value uint64, marker byte) map[string]any {
	var txHash primitives.Hash
	txHash[0] = marker
	addr := wallet.ReceiveAddress()
	return map[string]any{
		"hash":    hex.EncodeToString(txHash[:]),
		"index":   int(marker),
		"value":   int(value),
		"address": map[string]any{"version": 0, "hash": hex.EncodeToString(addr.Hash)},
	}
}

// TestNameBroadcast_OpenPhase_Good drives lns.name.register in broadcast mode for
// the OPEN phase: an unlinked covenant funded purely from the wallet, signed and
// broadcast through the fake daemon.
func TestNameBroadcast_OpenPhase_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(t.Context()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}
	wallet := mustImportWriteWallet(t)
	svc.wallet = wallet
	daemon := newFakeDaemon(t)

	result := c.Action(actionNameRegister).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testRegisterName},
		core.Option{Key: "broadcast", Value: true},
		core.Option{Key: "phase", Value: "open"},
		core.Option{Key: "chain_url", Value: daemon.server.URL},
		core.Option{Key: "open_value", Value: 0},
		core.Option{Key: "coins", Value: []any{coinOption(wallet, 100000, 11)}},
	))
	if !result.OK {
		t.Fatalf("lns.name.register (open broadcast) returned %#v, want OK", result)
	}
	res, ok := result.Value.(nameBroadcastResult)
	if !ok {
		t.Fatalf("register broadcast returned %T, want nameBroadcastResult", result.Value)
	}
	if res.Broadcast == nil || res.Broadcast.DaemonTxID == "" {
		t.Fatalf("register broadcast = %#v, want a daemon txid", res)
	}
	if res.Covenant != TypeOpen {
		t.Fatalf("covenant type = %v, want OPEN", res.Covenant)
	}
	if daemon.sendCalls != 1 {
		t.Fatalf("sendrawtransaction calls = %d, want 1", daemon.sendCalls)
	}
}

// TestNameBroadcast_RenewLinked_Good drives lns.name.renew in broadcast mode: a
// linked covenant that must spend the name's owner coin, supplied explicitly and
// also offered as a funding coin so the write-path can resolve its value.
func TestNameBroadcast_RenewLinked_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(t.Context()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}
	wallet := mustImportWriteWallet(t)
	svc.wallet = wallet
	daemon := newFakeDaemon(t)

	// The owner coin is a wallet coin (so it can be signed) and is named as the
	// required input.
	var ownerHash primitives.Hash
	ownerHash[0] = 21
	ownerCoin := coinOption(wallet, 50000, 21)

	result := c.Action(actionNameRenew).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testRegisterName},
		core.Option{Key: "broadcast", Value: true},
		core.Option{Key: "chain_url", Value: daemon.server.URL},
		core.Option{Key: "renew_value", Value: 50000},
		core.Option{Key: "owner_outpoint", Value: map[string]any{
			"hash": hex.EncodeToString(ownerHash[:]), "index": 21,
		}},
		core.Option{Key: "coins", Value: []any{ownerCoin, coinOption(wallet, 100000, 22)}},
	))
	if !result.OK {
		t.Fatalf("lns.name.renew (broadcast) returned %#v, want OK", result)
	}
	res, ok := result.Value.(nameBroadcastResult)
	if !ok {
		t.Fatalf("renew broadcast returned %T, want nameBroadcastResult", result.Value)
	}
	if res.Broadcast == nil || res.Broadcast.InputCount < 2 {
		t.Fatalf("renew broadcast = %#v, want owner coin plus funding input", res.Broadcast)
	}
	if res.Covenant != TypeRenew {
		t.Fatalf("covenant type = %v, want RENEW", res.Covenant)
	}
}

// TestNameBroadcast_RenewMissingOwner_Bad covers a linked covenant broadcast with
// no owner coin available: it must fail rather than broadcast an invalid
// transaction. The daemon getnameinfo returns no owner.
func TestNameBroadcast_RenewMissingOwner_Bad(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(t.Context()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}
	wallet := mustImportWriteWallet(t)
	svc.wallet = wallet
	daemon := newFakeDaemon(t) // getnameinfo returns nothing for an unknown name

	result := c.Action(actionNameRenew).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testRegisterName},
		core.Option{Key: "broadcast", Value: true},
		core.Option{Key: "chain_url", Value: daemon.server.URL},
		core.Option{Key: "renew_value", Value: 50000},
		core.Option{Key: "coins", Value: []any{coinOption(wallet, 100000, 23)}},
	))
	if result.OK {
		t.Fatalf("renew broadcast succeeded with no owner coin, want failure")
	}
	if daemon.sendCalls != 0 {
		t.Fatalf("sendrawtransaction called %d times, want 0 when owner is unresolved", daemon.sendCalls)
	}
}

// TestNameBroadcast_RenewOwnerFromDaemon_Good resolves the name owner coin from
// the daemon's getnameinfo (rather than an explicit owner_outpoint) and
// broadcasts the renewal.
func TestNameBroadcast_RenewOwnerFromDaemon_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(t.Context()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}
	wallet := mustImportWriteWallet(t)
	svc.wallet = wallet
	daemon := newFakeDaemon(t)

	// Register the owner coin in the daemon name info; the same coin is a wallet
	// coin so it can be signed.
	var ownerHash primitives.Hash
	ownerHash[0] = 31
	daemon.nameOwners[testRegisterName] = chainpkg.NameOwner{Hash: hex.EncodeToString(ownerHash[:]), Index: 31}

	result := c.Action(actionNameRenew).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testRegisterName},
		core.Option{Key: "broadcast", Value: true},
		core.Option{Key: "chain_url", Value: daemon.server.URL},
		core.Option{Key: "renew_value", Value: 40000},
		core.Option{Key: "coins", Value: []any{coinOption(wallet, 40000, 31), coinOption(wallet, 100000, 32)}},
	))
	if !result.OK {
		t.Fatalf("renew (owner from daemon) returned %#v, want OK", result)
	}
	res, ok := result.Value.(nameBroadcastResult)
	if !ok {
		t.Fatalf("renew broadcast returned %T, want nameBroadcastResult", result.Value)
	}
	if res.Broadcast == nil || res.Broadcast.InputCount < 2 {
		t.Fatalf("renew broadcast = %#v, want owner coin plus funding input", res.Broadcast)
	}
}

// TestNameBroadcast_RegisterBadPhase_Ugly covers an invalid phase selector.
func TestNameBroadcast_RegisterBadPhase_Ugly(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(t.Context()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}
	wallet := mustImportWriteWallet(t)
	svc.wallet = wallet
	daemon := newFakeDaemon(t)

	result := c.Action(actionNameRegister).Run(t.Context(), core.NewOptions(
		core.Option{Key: "name", Value: testRegisterName},
		core.Option{Key: "broadcast", Value: true},
		core.Option{Key: "phase", Value: "nonsense"},
		core.Option{Key: "chain_url", Value: daemon.server.URL},
		core.Option{Key: "coins", Value: []any{coinOption(wallet, 100000, 24)}},
	))
	if result.OK {
		t.Fatalf("register broadcast accepted an invalid phase, want failure")
	}
}
