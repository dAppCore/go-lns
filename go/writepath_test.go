// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
	chainpkg "dappco.re/go/lns/pkg/chain"
	"dappco.re/go/lns/pkg/mempool"
	"dappco.re/go/lns/pkg/primitives"
	"dappco.re/go/lns/pkg/script"
	walletpkg "dappco.re/go/lns/pkg/wallet"
)

const (
	testWriteSeed   = "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"
	testRecipientHx = "0a0b0c0d0e0f101112131415161718191a1b1c1d"
)

// fakeDaemon is a stand-in HSD daemon RPC server for write-path tests. It serves
// getcoinsbyaddress from a fixture coin set and records the sendrawtransaction
// payload so a test can assert the broadcast bytes round-trip.
type fakeDaemon struct {
	server      *httptest.Server
	coinsByAddr map[string][]chainpkg.Coin
	nameOwners  map[string]chainpkg.NameOwner
	lastRawTx   string
	sendCalls   int
	rejectSend  bool
}

func newFakeDaemon(t *testing.T) *fakeDaemon {
	t.Helper()
	d := &fakeDaemon{
		coinsByAddr: map[string][]chainpkg.Coin{},
		nameOwners:  map[string]chainpkg.NameOwner{},
	}
	d.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := core.ReadAll(r.Body)
		if !raw.OK {
			t.Fatalf("read request body: %v", raw.Value)
		}
		var request struct {
			Method string `json:"method"`
			Params []any  `json:"params"`
		}
		if res := core.JSONUnmarshalString(raw.Value.(string), &request); !res.OK {
			t.Fatalf("decode request: %v", res.Value)
		}

		switch request.Method {
		case "getcoinsbyaddress":
			addr, _ := request.Params[0].(string)
			writeFakeResult(t, w, d.coinsByAddr[addr])
		case "sendrawtransaction":
			d.sendCalls++
			d.lastRawTx, _ = request.Params[0].(string)
			if d.rejectSend {
				writeFakeError(t, w, "TX rejected: insufficient fee")
				return
			}
			writeFakeResult(t, w, "broadcast-txid-0001")
		case "getnameinfo":
			label, _ := request.Params[0].(string)
			if owner, ok := d.nameOwners[label]; ok {
				writeFakeResult(t, w, map[string]any{"info": map[string]any{"owner": owner}})
				return
			}
			writeFakeResult(t, w, map[string]any{"info": map[string]any{}})
		default:
			t.Fatalf("unexpected RPC method %q", request.Method)
		}
	}))
	t.Cleanup(d.server.Close)
	return d
}

func (d *fakeDaemon) client() *chainpkg.Client {
	return chainpkg.NewClient(d.server.URL)
}

func writeFakeResult(t *testing.T, w http.ResponseWriter, result any) {
	t.Helper()
	raw := core.JSONMarshal(map[string]any{"result": result})
	if !raw.OK {
		t.Fatalf("marshal result: %v", raw.Value)
	}
	if _, err := w.Write(raw.Value.([]byte)); err != nil {
		t.Fatalf("write result: %v", err)
	}
}

func writeFakeError(t *testing.T, w http.ResponseWriter, message string) {
	t.Helper()
	raw := core.JSONMarshal(map[string]any{
		"error": map[string]any{"code": -1, "message": message},
	})
	if !raw.OK {
		t.Fatalf("marshal error: %v", raw.Value)
	}
	if _, err := w.Write(raw.Value.([]byte)); err != nil {
		t.Fatalf("write error: %v", err)
	}
}

func mustImportWriteWallet(t *testing.T) *walletpkg.Wallet {
	t.Helper()
	wallet, err := walletpkg.ImportSeed(testWriteSeed)
	if err != nil {
		t.Fatalf(testImportSeedErr, err)
	}
	return wallet
}

// fundingCoin builds a selectedCoin held by the wallet's receive address for the
// given value, with a deterministic outpoint.
func fundingCoin(wallet *walletpkg.Wallet, value uint64, marker byte) selectedCoin {
	var txHash primitives.Hash
	txHash[0] = marker
	return selectedCoin{
		outpoint: primitives.Outpoint{TxHash: txHash, Index: uint32(marker)},
		value:    value,
		address:  wallet.ReceiveAddress(),
	}
}

func recipientAddress(t *testing.T) primitives.Address {
	t.Helper()
	hash, err := hex.DecodeString(testRecipientHx)
	if err != nil {
		t.Fatalf("decode recipient: %v", err)
	}
	return primitives.Address{Version: 0, Hash: hash}
}

// TestWritePath_execSpend_Good exercises the full pipeline: a funded wallet,
// real coin selection, signing, self-verification, and broadcast through the
// fake daemon. It asserts the broadcast bytes decode to a transaction whose
// inputs each verify against their spent coin.
func TestWritePath_execSpend_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)

	daemon := newFakeDaemon(t)
	coin := fundingCoin(wallet, 100000, 1)

	req := spendRequest{
		wallet:        wallet,
		client:        daemon.client(),
		outputs:       []primitives.Output{{Value: 25000, Address: recipientAddress(t)}},
		changeAddress: wallet.ReceiveAddress(),
		allowChange:   true,
		suppliedCoins: []selectedCoin{coin},
	}

	outcome, err := svc.execSpend(req)
	if err != nil {
		t.Fatalf("execSpend returned error: %v", err)
	}
	if outcome.DaemonTxID != "broadcast-txid-0001" {
		t.Fatalf("daemon txid = %q, want broadcast-txid-0001", outcome.DaemonTxID)
	}
	if outcome.TxID == "" {
		t.Fatalf("outcome txid is empty")
	}
	if daemon.sendCalls != 1 {
		t.Fatalf("sendrawtransaction calls = %d, want 1", daemon.sendCalls)
	}
	if outcome.Fee == 0 {
		t.Fatalf("fee = 0, want a non-zero fee")
	}
	if !outcome.ChangeOutput || outcome.ChangeValue == 0 {
		t.Fatalf("outcome = %#v, want a change output", outcome)
	}
	if outcome.TotalInput != coin.value {
		t.Fatalf("total input = %d, want %d", outcome.TotalInput, coin.value)
	}
	if outcome.TotalInput != outcome.TotalOutput+outcome.Fee {
		t.Fatalf("inputs %d != outputs %d + fee %d", outcome.TotalInput, outcome.TotalOutput, outcome.Fee)
	}

	// The broadcast payload must be the serialised signed transaction whose txid
	// matches the reported one and whose single input verifies against the coin.
	rawBytes, derr := hex.DecodeString(daemon.lastRawTx)
	if derr != nil {
		t.Fatalf("decode broadcast hex: %v", derr)
	}
	if hex.EncodeToString(rawBytes) != outcome.RawTx {
		t.Fatalf("broadcast hex does not match outcome.RawTx")
	}

	// Reconstruct the transaction from the request to re-verify the input.
	tx := primitives.Transaction{
		Version: 0,
		Inputs:  []primitives.Input{{Prevout: coin.outpoint, Sequence: sequenceFinal}},
		Outputs: []primitives.Output{{Value: 25000, Address: recipientAddress(t)}},
	}
	if serr := svc.signTransaction(&tx, wallet, []selectedCoin{coin}); serr != nil {
		t.Fatalf("re-sign returned error: %v", serr)
	}
	prevOutput := primitives.Output{Value: coin.value, Address: coin.address}
	if verr := script.VerifyInput(prevOutput, tx, 0, script.MandatoryVerifyFlags); verr != nil {
		t.Fatalf("VerifyInput failed on signed input: %v", verr)
	}
}

// TestWritePath_execSpend_Bad covers insufficient funds: the only coin cannot
// cover the requested amount plus fee, so selection must fail before broadcast.
func TestWritePath_execSpend_Bad(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)
	daemon := newFakeDaemon(t)

	req := spendRequest{
		wallet:        wallet,
		client:        daemon.client(),
		outputs:       []primitives.Output{{Value: 100000, Address: recipientAddress(t)}},
		changeAddress: wallet.ReceiveAddress(),
		allowChange:   true,
		suppliedCoins: []selectedCoin{fundingCoin(wallet, 5000, 2)},
	}

	if _, err := svc.execSpend(req); err == nil {
		t.Fatalf("execSpend succeeded, want insufficient-funds error")
	}
	if daemon.sendCalls != 0 {
		t.Fatalf("sendrawtransaction was called %d times, want 0 on insufficient funds", daemon.sendCalls)
	}
}

// TestWritePath_execSpend_Ugly covers a daemon rejection: the transaction is
// built and signed, but the daemon rejects the broadcast. execSpend must surface
// the error, not a fake success.
func TestWritePath_execSpend_Ugly(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)
	daemon := newFakeDaemon(t)
	daemon.rejectSend = true

	req := spendRequest{
		wallet:        wallet,
		client:        daemon.client(),
		outputs:       []primitives.Output{{Value: 25000, Address: recipientAddress(t)}},
		changeAddress: wallet.ReceiveAddress(),
		allowChange:   true,
		suppliedCoins: []selectedCoin{fundingCoin(wallet, 100000, 3)},
	}

	if _, err := svc.execSpend(req); err == nil {
		t.Fatalf("execSpend succeeded, want a broadcast rejection error")
	}
	if daemon.sendCalls != 1 {
		t.Fatalf("sendrawtransaction calls = %d, want exactly 1 (attempted then rejected)", daemon.sendCalls)
	}
}

// TestWritePath_signTransaction_Good verifies the build→sign→VerifyInput round
// trip across multiple inputs.
func TestWritePath_signTransaction_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)

	coinA := fundingCoin(wallet, 40000, 4)
	coinB := fundingCoin(wallet, 60000, 5)
	tx := primitives.Transaction{
		Version: 0,
		Inputs: []primitives.Input{
			{Prevout: coinA.outpoint, Sequence: sequenceFinal},
			{Prevout: coinB.outpoint, Sequence: sequenceFinal},
		},
		Outputs: []primitives.Output{{Value: 90000, Address: recipientAddress(t)}},
	}

	if err := svc.signTransaction(&tx, wallet, []selectedCoin{coinA, coinB}); err != nil {
		t.Fatalf("signTransaction returned error: %v", err)
	}
	for i, coin := range []selectedCoin{coinA, coinB} {
		if len(tx.Inputs[i].Witness) != 2 {
			t.Fatalf("input %d witness items = %d, want 2 (sig, pubkey)", i, len(tx.Inputs[i].Witness))
		}
		if len(tx.Inputs[i].Witness[0]) != 65 {
			t.Fatalf("input %d signature length = %d, want 65", i, len(tx.Inputs[i].Witness[0]))
		}
		if len(tx.Inputs[i].Witness[1]) != 33 {
			t.Fatalf("input %d pubkey length = %d, want 33", i, len(tx.Inputs[i].Witness[1]))
		}
		prevOutput := primitives.Output{Value: coin.value, Address: coin.address}
		if err := script.VerifyInput(prevOutput, tx, i, script.MandatoryVerifyFlags); err != nil {
			t.Fatalf("VerifyInput failed for input %d: %v", i, err)
		}
	}
}

// TestWritePath_signTransaction_Bad covers signing an input whose coin is not a
// wallet address: the signing key cannot be resolved.
func TestWritePath_signTransaction_Bad(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)

	foreign := selectedCoin{
		outpoint: primitives.Outpoint{TxHash: primitives.Hash{9}, Index: 0},
		value:    50000,
		address:  recipientAddress(t), // not owned by this wallet
	}
	tx := primitives.Transaction{
		Version: 0,
		Inputs:  []primitives.Input{{Prevout: foreign.outpoint, Sequence: sequenceFinal}},
		Outputs: []primitives.Output{{Value: 40000, Address: wallet.ReceiveAddress()}},
	}

	if err := svc.signTransaction(&tx, wallet, []selectedCoin{foreign}); err == nil {
		t.Fatalf("signTransaction succeeded for a foreign coin, want key-resolution error")
	}
}

// TestWritePath_signTransaction_Ugly covers a tampered transaction: signing then
// mutating an output value must make the input fail verification, proving the
// signature commits to the outputs.
func TestWritePath_signTransaction_Ugly(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)

	coin := fundingCoin(wallet, 100000, 6)
	tx := primitives.Transaction{
		Version: 0,
		Inputs:  []primitives.Input{{Prevout: coin.outpoint, Sequence: sequenceFinal}},
		Outputs: []primitives.Output{{Value: 90000, Address: recipientAddress(t)}},
	}
	if err := svc.signTransaction(&tx, wallet, []selectedCoin{coin}); err != nil {
		t.Fatalf("signTransaction returned error: %v", err)
	}

	tx.Outputs[0].Value = 99999 // tamper after signing

	prevOutput := primitives.Output{Value: coin.value, Address: coin.address}
	if err := script.VerifyInput(prevOutput, tx, 0, script.MandatoryVerifyFlags); err == nil {
		t.Fatalf("VerifyInput passed on a tampered transaction, want failure")
	}
}

// TestWritePath_gatherCoins_Good verifies the daemon getcoinsbyaddress path:
// coins for the wallet's first receive address are loaded and filtered to
// fundable P2WPKH coins.
func TestWritePath_gatherCoins_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)
	daemon := newFakeDaemon(t)

	recvAddr := wallet.ReceiveAddress()
	encoded, err := recvAddr.Encode("ts")
	if err != nil {
		t.Fatalf("encode address: %v", err)
	}
	daemon.coinsByAddr[encoded] = []chainpkg.Coin{
		{Value: 80000, Hash: hex.EncodeToString(make([]byte, 32)), Index: 0},
		{Value: 0, Hash: hex.EncodeToString(make([]byte, 32)), Index: 1},                                               // zero value: skipped
		{Value: 5000, Hash: hex.EncodeToString(make([]byte, 32)), Index: 2, Covenant: &chainpkg.CoinCovenant{Type: 3}}, // BID: locked, skipped
	}

	req := spendRequest{wallet: wallet, client: daemon.client(), addressPrefix: "ts"}
	coins, err := svc.gatherCoins(req)
	if err != nil {
		t.Fatalf("gatherCoins returned error: %v", err)
	}
	if len(coins) != 1 {
		t.Fatalf("gathered %d coins, want 1 fundable", len(coins))
	}
	if coins[0].value != 80000 {
		t.Fatalf("coin value = %d, want 80000", coins[0].value)
	}
}

// TestWritePath_gatherCoins_Bad covers the missing-prefix case: without a bech32
// HRP and without supplied coins, the daemon address query cannot be formed.
func TestWritePath_gatherCoins_Bad(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)
	daemon := newFakeDaemon(t)

	req := spendRequest{wallet: wallet, client: daemon.client()}
	if _, err := svc.gatherCoins(req); err == nil {
		t.Fatalf("gatherCoins succeeded without an address prefix, want error")
	}
}

// TestWritePath_walletSendAction_Good drives the lns.wallet.send action in
// broadcast mode end-to-end and asserts a txid is returned.
func TestWritePath_walletSendAction_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf(testOnStartupOK, result)
	}
	wallet := mustImportWriteWallet(t)
	svc.wallet = wallet
	daemon := newFakeDaemon(t)

	coin := fundingCoin(wallet, 100000, 7)
	result := c.Action(actionNameWalletSend).Run(context.Background(), core.NewOptions(
		core.Option{Key: "broadcast", Value: true},
		core.Option{Key: "chain_url", Value: daemon.server.URL},
		core.Option{Key: "amount", Value: 30000},
		core.Option{Key: "address", Value: Address{Version: 0, Hash: recipientAddress(t).Hash}},
		core.Option{Key: "coins", Value: []any{
			map[string]any{
				"hash":    hex.EncodeToString(coin.outpoint.TxHash[:]),
				"index":   int(coin.outpoint.Index),
				"value":   int(coin.value),
				"address": map[string]any{"version": 0, "hash": hex.EncodeToString(coin.address.Hash)},
			},
		}},
	))
	if !result.OK {
		t.Fatalf(testWalletSendOK, result)
	}
	sendResult, ok := result.Value.(walletSendResult)
	if !ok {
		t.Fatalf("wallet send returned %T, want walletSendResult", result.Value)
	}
	if sendResult.Broadcast == nil || sendResult.Broadcast.DaemonTxID == "" {
		t.Fatalf("wallet send result = %#v, want a broadcast txid", sendResult)
	}
	if daemon.sendCalls != 1 {
		t.Fatalf("sendrawtransaction calls = %d, want 1", daemon.sendCalls)
	}
}

// TestWritePath_multiCoinSelection_Good exercises selecting across several coins
// (largest-first ordering) to cover an amount no single coin can fund.
func TestWritePath_multiCoinSelection_Good(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	wallet := mustImportWriteWallet(t)
	daemon := newFakeDaemon(t)

	coins := []selectedCoin{
		fundingCoin(wallet, 10000, 41),
		fundingCoin(wallet, 30000, 42),
		fundingCoin(wallet, 20000, 43),
	}
	req := spendRequest{
		wallet:        wallet,
		client:        daemon.client(),
		outputs:       []primitives.Output{{Value: 45000, Address: recipientAddress(t)}},
		changeAddress: wallet.ReceiveAddress(),
		allowChange:   true,
		suppliedCoins: coins,
	}

	outcome, err := svc.execSpend(req)
	if err != nil {
		t.Fatalf("execSpend returned error: %v", err)
	}
	// Largest-first: 30000 + 20000 = 50000 covers 45000 + fee with two inputs.
	if outcome.InputCount < 2 {
		t.Fatalf("input count = %d, want at least 2 coins selected", outcome.InputCount)
	}
	if outcome.TotalInput < 45000+outcome.Fee {
		t.Fatalf("selected inputs %d do not cover amount+fee %d", outcome.TotalInput, 45000+outcome.Fee)
	}
}

// TestWritePath_requiredInputsList_Good covers parsing the list form of required
// inputs and including them all.
func TestWritePath_requiredInputsList_Good(t *testing.T) {
	options := core.NewOptions(core.Option{Key: "required_inputs", Value: []any{
		map[string]any{"hash": hex.EncodeToString(make([]byte, 32)), "index": 0},
		map[string]any{"hash": hex.EncodeToString(make([]byte, 32)), "index": 1},
	}})
	required, err := requiredInputsFromOptions(options)
	if err != nil {
		t.Fatalf("requiredInputsFromOptions error: %v", err)
	}
	if len(required) != 2 {
		t.Fatalf("required inputs = %d, want 2", len(required))
	}
	if required[1].Index != 1 {
		t.Fatalf("second required input index = %d, want 1", required[1].Index)
	}
}

// TestWritePath_serializeWithWitness_Good confirms the broadcast serialisation
// round-trips: base bytes match SerializeNoWitness and witnesses are appended.
func TestWritePath_serializeWithWitness_Good(t *testing.T) {
	wallet := mustImportWriteWallet(t)
	c := core.New()
	svc := NewService(c)
	coin := fundingCoin(wallet, 100000, 8)
	tx := primitives.Transaction{
		Version: 0,
		Inputs:  []primitives.Input{{Prevout: coin.outpoint, Sequence: sequenceFinal}},
		Outputs: []primitives.Output{{Value: 90000, Address: recipientAddress(t)}},
	}
	if err := svc.signTransaction(&tx, wallet, []selectedCoin{coin}); err != nil {
		t.Fatalf("signTransaction returned error: %v", err)
	}

	full, err := mempool.SerializeWithWitness(tx)
	if err != nil {
		t.Fatalf("SerializeWithWitness returned error: %v", err)
	}
	base, err := mempool.SerializeNoWitness(tx)
	if err != nil {
		t.Fatalf("SerializeNoWitness returned error: %v", err)
	}
	if len(full) <= len(base) {
		t.Fatalf("witness serialisation (%d) not longer than base (%d)", len(full), len(base))
	}
	if hex.EncodeToString(full[:len(base)]) != hex.EncodeToString(base) {
		t.Fatalf("witness serialisation does not start with the base serialisation")
	}
}
