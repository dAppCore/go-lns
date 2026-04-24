// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

func TestWalletCreateImportExportRoundTrip(t *testing.T) {
	first, err := Create()
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	seed := first.ExportSeed()
	if len(seed) != 64 {
		t.Fatalf("ExportSeed length = %d, want 64", len(seed))
	}

	second, err := ImportSeed(seed)
	if err != nil {
		t.Fatalf("ImportSeed returned error: %v", err)
	}

	if got, want := second.ExportSeed(), seed; got != want {
		t.Fatalf("round-trip seed = %q, want %q", got, want)
	}

	if got, want := second.ReceiveAddress(), first.ReceiveAddress(); !got.Equals(want) {
		t.Fatalf("ReceiveAddress = %#v, want %#v", got, want)
	}
}

func TestWalletBalanceHistoryAndReceive(t *testing.T) {
	wallet, err := ImportSeed("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatalf("ImportSeed returned error: %v", err)
	}

	firstAddress := wallet.ReceiveAddress()
	nextAddress := wallet.NextReceiveAddress()
	if firstAddress.Equals(nextAddress) {
		t.Fatal("NextReceiveAddress should advance the derivation index")
	}

	wallet.Credit(100, false, 10, primitives.Hash{1}, "credit")
	wallet.Credit(25, true, 11, primitives.Hash{2}, "locked")
	wallet.Debit(30, false, 12, primitives.Hash{3}, "debit")

	confirmed, locked := wallet.Balance()
	if confirmed != 70 || locked != 25 {
		t.Fatalf("Balance = (%d, %d), want (70, 25)", confirmed, locked)
	}

	history := wallet.History()
	if len(history) != 3 {
		t.Fatalf("History returned %d entries, want 3", len(history))
	}
	if history[0].Hash != (primitives.Hash{3}) {
		t.Fatalf("History[0] = %#v, want newest debit first", history[0])
	}
}

func TestWalletPendingDebitDoesNotAffectBalance(t *testing.T) {
	wallet, err := ImportSeed("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	if err != nil {
		t.Fatalf("ImportSeed returned error: %v", err)
	}

	wallet.Credit(100, false, 10, primitives.Hash{1}, "credit")
	entry := wallet.PendingDebit(40, primitives.Hash{9}, "draft")

	confirmed, locked := wallet.Balance()
	if confirmed != 100 || locked != 0 {
		t.Fatalf("Balance = (%d, %d), want (100, 0)", confirmed, locked)
	}
	if !entry.Pending || entry.Kind != EntryDebit {
		t.Fatalf("PendingDebit() = %#v, want pending debit entry", entry)
	}

	history := wallet.History()
	if len(history) != 2 || !history[0].Pending || history[0].Hash != (primitives.Hash{9}) {
		t.Fatalf("History = %#v, want pending draft entry first", history)
	}
}
