// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"context"
	"encoding/hex"
	"testing"

	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/primitives"
	walletpkg "dappco.re/go/lns/pkg/wallet"
)

func TestWalletSendReturnsDraftPlan(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	wallet, err := walletpkg.ImportSeed("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatalf("ImportSeed returned error: %v", err)
	}
	wallet.Credit(1000, false, 10, primitives.Hash{1}, "funding")
	svc.wallet = wallet

	result := c.Action("lns.wallet.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "amount", Value: 200},
		core.Option{Key: "address", Value: Address{Version: 0, Hash: []byte{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		}}},
	))
	if !result.OK {
		t.Fatalf("lns.wallet.send returned %#v, want OK", result)
	}

	plan, ok := result.Value.(walletSendPlanResult)
	if !ok {
		t.Fatalf("lns.wallet.send returned %T, want walletSendPlanResult", result.Value)
	}

	if !plan.SufficientFunds || plan.Confirmed != 1000 || plan.Amount != 200 {
		t.Fatalf("wallet send plan = %#v, want sufficient balance for 200/1000", plan)
	}
	if plan.Draft == nil || len(plan.Draft.Transaction.Outputs) != 2 {
		t.Fatalf("wallet send draft = %#v, want recipient and change outputs", plan.Draft)
	}
	if plan.Draft.Fee == 0 || !plan.Draft.ChangeOutputPresent || plan.ChangeAddress == nil {
		t.Fatalf("wallet send draft = %#v, want placeholder fee and change output", plan.Draft)
	}
	if plan.DraftID == "" {
		t.Fatalf("wallet send plan = %#v, want deterministic draft ID", plan)
	}

	second := c.Action("lns.wallet.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "amount", Value: 200},
		core.Option{Key: "address", Value: Address{Version: 0, Hash: []byte{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		}}},
	))
	if !second.OK {
		t.Fatalf("second lns.wallet.send returned %#v, want OK", second)
	}

	secondPlan, ok := second.Value.(walletSendPlanResult)
	if !ok {
		t.Fatalf("second lns.wallet.send returned %T, want walletSendPlanResult", second.Value)
	}
	if secondPlan.DraftID != plan.DraftID {
		t.Fatalf("draft ID = %q, second draft ID = %q, want deterministic match", plan.DraftID, secondPlan.DraftID)
	}
}

func TestWalletSendCanRecordPendingHistoryEntry(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	wallet, err := walletpkg.ImportSeed("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatalf("ImportSeed returned error: %v", err)
	}
	wallet.Credit(1000, false, 10, primitives.Hash{1}, "funding")
	svc.wallet = wallet

	result := c.Action("lns.wallet.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "amount", Value: 200},
		core.Option{Key: "record_history", Value: true},
		core.Option{Key: "note", Value: "queued send"},
		core.Option{Key: "address", Value: Address{Version: 0, Hash: []byte{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		}}},
	))
	if !result.OK {
		t.Fatalf("lns.wallet.send returned %#v, want OK", result)
	}

	plan, ok := result.Value.(walletSendPlanResult)
	if !ok {
		t.Fatalf("lns.wallet.send returned %T, want walletSendPlanResult", result.Value)
	}
	if !plan.Recorded || plan.PendingEntry == nil || !plan.PendingEntry.Pending {
		t.Fatalf("wallet send plan = %#v, want recorded pending entry", plan)
	}
	if plan.PendingEntry.Note != "queued send" {
		t.Fatalf("pending entry = %#v, want note %q", plan.PendingEntry, "queued send")
	}
	if got := hex.EncodeToString(plan.PendingEntry.Hash[:]); got != plan.DraftID {
		t.Fatalf("pending entry hash = %q, want draft ID %q", got, plan.DraftID)
	}

	confirmed, locked := wallet.Balance()
	if confirmed != 1000 || locked != 0 {
		t.Fatalf("Balance = (%d, %d), want (1000, 0) after pending record", confirmed, locked)
	}

	history := wallet.History()
	if len(history) == 0 || !history[0].Pending {
		t.Fatalf("History = %#v, want pending entry first", history)
	}
}

func TestWalletSendCanSuppressAutoChange(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	wallet, err := walletpkg.ImportSeed("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatalf("ImportSeed returned error: %v", err)
	}
	wallet.Credit(1000, false, 10, primitives.Hash{1}, "funding")
	svc.wallet = wallet

	result := c.Action("lns.wallet.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "amount", Value: 200},
		core.Option{Key: "no_change", Value: true},
		core.Option{Key: "address", Value: Address{Version: 0, Hash: []byte{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		}}},
	))
	if !result.OK {
		t.Fatalf("lns.wallet.send returned %#v, want OK", result)
	}

	plan, ok := result.Value.(walletSendPlanResult)
	if !ok {
		t.Fatalf("lns.wallet.send returned %T, want walletSendPlanResult", result.Value)
	}
	if plan.Draft == nil || len(plan.Draft.Transaction.Outputs) != 1 {
		t.Fatalf("wallet send draft = %#v, want recipient-only draft output", plan.Draft)
	}
	if plan.Draft.ChangeOutputPresent || plan.ChangeAddress != nil {
		t.Fatalf("wallet send draft = %#v, want no auto-change output", plan.Draft)
	}
}

func TestWalletSendRequiresWallet(t *testing.T) {
	c := core.New()
	svc := NewService(c)
	if result := svc.OnStartup(context.Background()); !result.OK {
		t.Fatalf("OnStartup returned %#v, want OK", result)
	}

	result := c.Action("lns.wallet.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "amount", Value: 200},
		core.Option{Key: "address", Value: Address{Version: 0, Hash: []byte{
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
			2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
		}}},
	))
	if result.OK {
		t.Fatal("lns.wallet.send should require an initialized wallet")
	}

	err, ok := result.Value.(error)
	if !ok {
		t.Fatalf("lns.wallet.send returned %T, want error", result.Value)
	}
	if got := core.Operation(err); got != "lns.wallet.send" {
		t.Fatalf("wallet.send error operation = %q, want %q", got, "lns.wallet.send")
	}
}
