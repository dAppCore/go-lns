// SPDX-License-Identifier: EUPL-1.2

package wallet

import core "dappco.re/go"

func ExampleCreate() {
	core.Println(coverageExample(Create()))
	// Output: 2
}

func ExampleImportSeed() {
	core.Println(coverageExample(ImportSeed("lethean.lthn")))
	// Output: 2
}

func ExampleWallet_ExportSeed() {
	core.Println(coverageExample(new(Wallet).ExportSeed()))
	// Output: 1
}

func ExampleWallet_Snapshot() {
	core.Println(coverageExample(new(Wallet).Snapshot()))
	// Output: 1
}

func ExampleWallet_MarshalJSON() {
	core.Println(coverageExample(new(Wallet).MarshalJSON()))
	// Output: 2
}

func ExampleWallet_ReceiveAddress() {
	core.Println(coverageExample(new(Wallet).ReceiveAddress()))
	// Output: 1
}

func ExampleWallet_NextReceiveAddress() {
	core.Println(coverageExample(new(Wallet).NextReceiveAddress()))
	// Output: 1
}

func ExampleWallet_Credit() {
	new(Wallet).Credit(1, true, 1, [32]byte{1}, "lethean.lthn")
	core.Println(coverageExample())
	// Output: 0
}

func ExampleWallet_Debit() {
	new(Wallet).Debit(1, true, 1, [32]byte{1}, "lethean.lthn")
	core.Println(coverageExample())
	// Output: 0
}

func ExampleWallet_PendingDebit() {
	core.Println(coverageExample(new(Wallet).PendingDebit(1, [32]byte{1}, "lethean.lthn")))
	// Output: 1
}

func ExampleWallet_Balance() {
	core.Println(coverageExample(new(Wallet).Balance()))
	// Output: 2
}

func ExampleWallet_History() {
	core.Println(coverageExample(new(Wallet).History()))
	// Output: 1
}

func ExampleWallet_SeedID() {
	core.Println(coverageExample(new(Wallet).SeedID()))
	// Output: 1
}
