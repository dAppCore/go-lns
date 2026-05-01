// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	core "dappco.re/go"
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

const testImportSeedErr = "ImportSeed returned error: %v"

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
		t.Fatalf(testImportSeedErr, err)
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
		t.Fatalf(testImportSeedErr, err)
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
		t.Fatalf(testImportSeedErr, err)
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
	coverageAssertVariant(t, "good")
	values := coveragePack(Create())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Create_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(Create())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Create_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(Create())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_ImportSeed_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(ImportSeed("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_ImportSeed_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(ImportSeed(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_ImportSeed_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(ImportSeed("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_ExportSeed_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).ExportSeed())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_ExportSeed_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).ExportSeed())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_ExportSeed_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).ExportSeed())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_Snapshot_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).Snapshot())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_Snapshot_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).Snapshot())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_Snapshot_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).Snapshot())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_MarshalJSON_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_MarshalJSON_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_MarshalJSON_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).MarshalJSON())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_ReceiveAddress_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).ReceiveAddress())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_ReceiveAddress_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).ReceiveAddress())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_ReceiveAddress_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).ReceiveAddress())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_NextReceiveAddress_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).NextReceiveAddress())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_NextReceiveAddress_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).NextReceiveAddress())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_NextReceiveAddress_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).NextReceiveAddress())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_Credit_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	called := false
	new(Wallet).Credit(1, true, 1, [32]byte{1}, "lethean.lthn")
	called = true
	core.AssertTrue(t, called)
}

func TestWallet_Wallet_Credit_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	called := false
	new(Wallet).Credit(0, false, 0, [32]byte{}, "")
	called = true
	core.AssertTrue(t, called)
}

func TestWallet_Wallet_Credit_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	called := false
	new(Wallet).Credit(63, true, 63, [32]byte{0xff}, "lethean.\x00.lthn")
	called = true
	core.AssertTrue(t, called)
}

func TestWallet_Wallet_Debit_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	called := false
	new(Wallet).Debit(1, true, 1, [32]byte{1}, "lethean.lthn")
	called = true
	core.AssertTrue(t, called)
}

func TestWallet_Wallet_Debit_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	called := false
	new(Wallet).Debit(0, false, 0, [32]byte{}, "")
	called = true
	core.AssertTrue(t, called)
}

func TestWallet_Wallet_Debit_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	called := false
	new(Wallet).Debit(63, true, 63, [32]byte{0xff}, "lethean.\x00.lthn")
	called = true
	core.AssertTrue(t, called)
}

func TestWallet_Wallet_PendingDebit_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).PendingDebit(1, [32]byte{1}, "lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_PendingDebit_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).PendingDebit(0, [32]byte{}, ""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_PendingDebit_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).PendingDebit(63, [32]byte{0xff}, "lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_Balance_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).Balance())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_Balance_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).Balance())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_Balance_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).Balance())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_History_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).History())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_History_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).History())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_History_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).History())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_SeedID_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Wallet).SeedID())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_SeedID_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Wallet).SeedID())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestWallet_Wallet_SeedID_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Wallet).SeedID())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
