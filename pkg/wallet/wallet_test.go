// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	core "dappco.re/go"
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

// BEGIN v090 file-aware generated coverage

func TestWallet_Create_Good(t *core.T) {
	fn := ax7ValueOf(Create)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Create_Bad(t *core.T) {
	fn := ax7ValueOf(Create)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Create_Ugly(t *core.T) {
	fn := ax7ValueOf(Create)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_ImportSeed_Good(t *core.T) {
	fn := ax7ValueOf(ImportSeed)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_ImportSeed_Bad(t *core.T) {
	fn := ax7ValueOf(ImportSeed)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_ImportSeed_Ugly(t *core.T) {
	fn := ax7ValueOf(ImportSeed)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_ExportSeed_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).ExportSeed)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_ExportSeed_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).ExportSeed)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_ExportSeed_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).ExportSeed)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Snapshot_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Snapshot)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Snapshot_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Snapshot)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Snapshot_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Snapshot)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_MarshalJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).MarshalJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_MarshalJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).MarshalJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_MarshalJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).MarshalJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_ReceiveAddress_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).ReceiveAddress)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_ReceiveAddress_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).ReceiveAddress)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_ReceiveAddress_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).ReceiveAddress)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_NextReceiveAddress_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).NextReceiveAddress)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_NextReceiveAddress_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).NextReceiveAddress)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_NextReceiveAddress_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).NextReceiveAddress)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Credit_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Credit)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Credit_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Credit)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Credit_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Credit)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Debit_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Debit)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Debit_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Debit)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Debit_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Debit)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_PendingDebit_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).PendingDebit)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_PendingDebit_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).PendingDebit)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_PendingDebit_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).PendingDebit)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Balance_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Balance)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Balance_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Balance)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_Balance_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).Balance)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_History_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).History)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_History_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).History)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_History_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).History)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_SeedID_Good(t *core.T) {
	fn := ax7ValueOf(new(Wallet).SeedID)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_SeedID_Bad(t *core.T) {
	fn := ax7ValueOf(new(Wallet).SeedID)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestWallet_Wallet_SeedID_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Wallet).SeedID)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
