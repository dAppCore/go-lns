// SPDX-License-Identifier: EUPL-1.2

package hd

import (
	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

const (
	opDeriveAccount    = "hd.DeriveAccount"
	opDeriveAddressKey = "hd.DeriveAddressKey"
	opDerivePath       = "hd.Derive"
	opDeriveAddress    = "hd.DeriveAddress"
)

// DeriveAccount derives a BIP-44 account node from a master key:
// "m/44'/coinType'/account'". All three levels are hardened, so the master must
// be a private node. The returned node is the account extended key from which
// receive and change address keys are derived.
//
//	acct, err := hd.DeriveAccount(master, hd.CoinTypeLethean, 0)
func DeriveAccount(master *ExtendedKey, coinType, account uint32) (*ExtendedKey, error) {
	if master == nil {
		return nil, core.E(opDeriveAccount, "nil master key", nil)
	}
	if !master.IsPrivate() {
		return nil, core.E(opDeriveAccount, "account derivation requires a private master key", nil)
	}

	purpose, err := master.Child(HardenedKeyStart + Bip44Purpose)
	if err != nil {
		return nil, core.E(opDeriveAccount, "derive purpose level", err)
	}

	coin, err := purpose.Child(HardenedKeyStart + coinType)
	if err != nil {
		return nil, core.E(opDeriveAccount, "derive coin-type level", err)
	}

	acct, err := coin.Child(HardenedKeyStart + account)
	if err != nil {
		return nil, core.E(opDeriveAccount, "derive account level", err)
	}
	return acct, nil
}

// DeriveAddressKey derives an address key from an account node:
// "<account>/branch/index". Branch and index are normal (non-hardened), so the
// account node may be private or a neutered public node. Branch must be
// BranchReceive (0) or BranchChange (1).
//
//	key, err := hd.DeriveAddressKey(acct, hd.BranchReceive, 0)
func DeriveAddressKey(account *ExtendedKey, branch, index uint32) (*ExtendedKey, error) {
	if account == nil {
		return nil, core.E(opDeriveAddressKey, "nil account key", nil)
	}
	if branch != BranchReceive && branch != BranchChange {
		return nil, core.E(opDeriveAddressKey, "branch must be 0 (receive) or 1 (change)", nil)
	}
	if index >= HardenedKeyStart {
		return nil, core.E(opDeriveAddressKey, "address index must be non-hardened", nil)
	}

	branchKey, err := account.Child(branch)
	if err != nil {
		return nil, core.E(opDeriveAddressKey, "derive branch level", err)
	}

	addrKey, err := branchKey.Child(index)
	if err != nil {
		return nil, core.E(opDeriveAddressKey, "derive index level", err)
	}
	return addrKey, nil
}

// Path describes a BIP-44 derivation path relative to the master key. It mirrors
// the HSD wallet path object (docs/js-wallet/path.js): coin type, account,
// branch, and index, with purpose fixed to 44.
//
//	p := hd.Path{CoinType: hd.CoinTypeLethean, Account: 0, Branch: hd.BranchReceive, Index: 3}
type Path struct {
	// CoinType is the BIP-44 coin type (hardened in the derived path).
	CoinType uint32
	// Account is the BIP-44 account index (hardened in the derived path).
	Account uint32
	// Branch is the address branch: 0 receive, 1 change (non-hardened).
	Branch uint32
	// Index is the address index within the branch (non-hardened).
	Index uint32
}

// Derive derives the full BIP-44 address key for the path from a master key:
// "m/44'/CoinType'/Account'/Branch/Index". It combines DeriveAccount and
// DeriveAddressKey.
//
//	key, err := hd.Derive(master, hd.Path{CoinType: hd.CoinTypeLethean, Index: 0})
func Derive(master *ExtendedKey, p Path) (*ExtendedKey, error) {
	if master == nil {
		return nil, core.E(opDerivePath, "nil master key", nil)
	}

	acct, err := DeriveAccount(master, p.CoinType, p.Account)
	if err != nil {
		return nil, core.E(opDerivePath, "derive account", err)
	}

	key, err := DeriveAddressKey(acct, p.Branch, p.Index)
	if err != nil {
		return nil, core.E(opDerivePath, "derive address key", err)
	}
	return key, nil
}

// DeriveAddress derives the version-0 witness address for a BIP-44 path from a
// master key. It is shorthand for Derive followed by Address.
//
//	addr, err := hd.DeriveAddress(master, hd.Path{CoinType: hd.CoinTypeLethean, Index: 0})
func DeriveAddress(master *ExtendedKey, p Path) (primitives.Address, error) {
	key, err := Derive(master, p)
	if err != nil {
		return primitives.Address{}, core.E(opDeriveAddress, "derive key for address", err)
	}
	return key.Address(), nil
}
