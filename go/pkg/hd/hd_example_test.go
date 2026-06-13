// SPDX-License-Identifier: EUPL-1.2

package hd_test

import (
	core "dappco.re/go"
	"dappco.re/go/lns/pkg/hd"
)

// Example derives a Lethean BIP-44 receive address from a master seed,
// demonstrating the end-to-end package flow.
func Example() {
	seed := core.HexDecode("000102030405060708090a0b0c0d0e0f").Value.([]byte)

	master, err := hd.NewMaster(seed)
	if err != nil {
		core.Println("error")
		return
	}

	addr, err := hd.DeriveAddress(master, hd.Path{
		CoinType: hd.CoinTypeLethean,
		Account:  0,
		Branch:   hd.BranchReceive,
		Index:    0,
	})
	if err != nil {
		core.Println("error")
		return
	}

	core.Println(addr.IsPubkeyHash())
	// Output: true
}
