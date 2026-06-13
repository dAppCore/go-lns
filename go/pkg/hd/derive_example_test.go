// SPDX-License-Identifier: EUPL-1.2

package hd

import core "dappco.re/go"

func ExampleDeriveAccount() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	acct, _ := DeriveAccount(master, CoinTypeLethean, 0)
	core.Println(acct.Depth())
	// Output: 3
}

func ExampleDeriveAddressKey() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	acct, _ := DeriveAccount(master, CoinTypeLethean, 0)
	key, _ := DeriveAddressKey(acct, BranchReceive, 0)
	core.Println(key.Depth())
	// Output: 5
}

func ExampleDerive() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	key, _ := Derive(master, Path{CoinType: CoinTypeLethean, Account: 0, Branch: BranchReceive, Index: 0})
	core.Println(key.Depth())
	// Output: 5
}

func ExampleDeriveAddress() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)
	master, _ := NewMaster(seed)
	addr, _ := DeriveAddress(master, Path{CoinType: CoinTypeLethean, Account: 0, Branch: BranchReceive, Index: 0})
	core.Println(addr.IsPubkeyHash())
	// Output: true
}
