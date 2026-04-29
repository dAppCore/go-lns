// SPDX-License-Identifier: EUPL-1.2

package wallet

import core "dappco.re/go"

func ExampleCreate() {
	fn := ax7ValueOf(Create)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleImportSeed() {
	fn := ax7ValueOf(ImportSeed)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleWallet_ExportSeed() {
	fn := ax7ValueOf(new(Wallet).ExportSeed)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWallet_Snapshot() {
	fn := ax7ValueOf(new(Wallet).Snapshot)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWallet_MarshalJSON() {
	fn := ax7ValueOf(new(Wallet).MarshalJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleWallet_ReceiveAddress() {
	fn := ax7ValueOf(new(Wallet).ReceiveAddress)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWallet_NextReceiveAddress() {
	fn := ax7ValueOf(new(Wallet).NextReceiveAddress)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWallet_Credit() {
	fn := ax7ValueOf(new(Wallet).Credit)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleWallet_Debit() {
	fn := ax7ValueOf(new(Wallet).Debit)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleWallet_PendingDebit() {
	fn := ax7ValueOf(new(Wallet).PendingDebit)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWallet_Balance() {
	fn := ax7ValueOf(new(Wallet).Balance)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleWallet_History() {
	fn := ax7ValueOf(new(Wallet).History)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWallet_SeedID() {
	fn := ax7ValueOf(new(Wallet).SeedID)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
