// SPDX-License-Identifier: EUPL-1.2

package hd

import (
	"testing"

	core "github.com/dappcore/go"
)

// deriveManual walks m/44'/coin'/account'/branch/index one Child call at a time,
// providing an independent cross-check for the DeriveAccount/DeriveAddressKey
// helpers.
func deriveManual(t *testing.T, master *ExtendedKey, coin, account, branch, index uint32) *ExtendedKey {
	t.Helper()
	purpose, err := master.Child(HardenedKeyStart + Bip44Purpose)
	core.AssertNoError(t, err)
	coinKey, err := purpose.Child(HardenedKeyStart + coin)
	core.AssertNoError(t, err)
	acct, err := coinKey.Child(HardenedKeyStart + account)
	core.AssertNoError(t, err)
	branchKey, err := acct.Child(branch)
	core.AssertNoError(t, err)
	addrKey, err := branchKey.Child(index)
	core.AssertNoError(t, err)
	return addrKey
}

// TestDerive_MatchesManualWalk confirms Derive == step-by-step Child derivation
// for a Lethean BIP-44 path.
func TestDerive_MatchesManualWalk(t *testing.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)

	want := deriveManual(t, master, CoinTypeLethean, 0, BranchReceive, 0)
	wantPriv, err := want.PrivateKeyBytes()
	core.AssertNoError(t, err)

	got, err := Derive(master, Path{CoinType: CoinTypeLethean, Account: 0, Branch: BranchReceive, Index: 0})
	core.AssertNoError(t, err)
	gotPriv, err := got.PrivateKeyBytes()
	core.AssertNoError(t, err)

	core.AssertEqual(t, core.HexEncode(wantPriv), core.HexEncode(gotPriv))
	core.AssertEqual(t, uint8(5), got.Depth())
}

// TestDeriveAccount_PublicCanDeriveAddresses confirms a neutered account key can
// derive normal address public keys identical to the private path — the
// watch-only wallet property.
func TestDeriveAccount_PublicCanDeriveAddresses(t *testing.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)

	acct, err := DeriveAccount(master, CoinTypeLethean, 0)
	core.AssertNoError(t, err)

	privAddrKey, err := DeriveAddressKey(acct, BranchReceive, 5)
	core.AssertNoError(t, err)
	privPub, err := privAddrKey.PublicKeyBytes()
	core.AssertNoError(t, err)

	watchAcct, err := acct.Neuter()
	core.AssertNoError(t, err)
	pubAddrKey, err := DeriveAddressKey(watchAcct, BranchReceive, 5)
	core.AssertNoError(t, err)
	pubPub, err := pubAddrKey.PublicKeyBytes()
	core.AssertNoError(t, err)

	core.AssertEqual(t, core.HexEncode(privPub), core.HexEncode(pubPub))
}

// TestDeriveAddress_WalletCompatible confirms derived addresses use the
// primitives.Address version-0 P2WPKH form consumed by pkg/wallet.
func TestDeriveAddress_WalletCompatible(t *testing.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)

	addr, err := DeriveAddress(master, Path{CoinType: CoinTypeLethean, Account: 0, Branch: BranchReceive, Index: 0})
	core.AssertNoError(t, err)
	core.AssertTrue(t, addr.IsPubkeyHash())
	core.AssertTrue(t, addr.IsValid())
	core.AssertEqual(t, uint8(0), addr.Version)
	core.AssertEqual(t, addressHashLen, len(addr.Hash))
}

// TestDeriveAddressKey_ChangeBranchDiffers confirms receive and change branches
// derive distinct keys.
func TestDeriveAddressKey_ChangeBranchDiffers(t *testing.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	acct, err := DeriveAccount(master, CoinTypeLethean, 0)
	core.AssertNoError(t, err)

	recv, err := DeriveAddressKey(acct, BranchReceive, 0)
	core.AssertNoError(t, err)
	chg, err := DeriveAddressKey(acct, BranchChange, 0)
	core.AssertNoError(t, err)

	rp, err := recv.PrivateKeyBytes()
	core.AssertNoError(t, err)
	cp, err := chg.PrivateKeyBytes()
	core.AssertNoError(t, err)
	core.AssertTrue(t, core.HexEncode(rp) != core.HexEncode(cp))
}

// --- DeriveAccount Good/Bad/Ugly ---

func TestDerive_DeriveAccount_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	acct, err := DeriveAccount(master, CoinTypeLethean, 0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint8(3), acct.Depth())
	core.AssertTrue(t, acct.IsPrivate())
}

func TestDerive_DeriveAccount_Bad(t *core.T) {
	_, err := DeriveAccount(nil, CoinTypeLethean, 0)
	core.AssertError(t, err)
}

func TestDerive_DeriveAccount_Ugly(t *core.T) {
	// A neutered master cannot derive a hardened account.
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	_, err = DeriveAccount(pub, CoinTypeLethean, 0)
	core.AssertError(t, err)
}

// --- DeriveAddressKey Good/Bad/Ugly ---

func TestDerive_DeriveAddressKey_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	acct, err := DeriveAccount(master, CoinTypeLethean, 0)
	core.AssertNoError(t, err)
	key, err := DeriveAddressKey(acct, BranchReceive, 0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint8(5), key.Depth())
}

func TestDerive_DeriveAddressKey_Bad(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	acct, err := DeriveAccount(master, CoinTypeLethean, 0)
	core.AssertNoError(t, err)
	// Invalid branch (not 0 or 1) must fail.
	_, err = DeriveAddressKey(acct, 2, 0)
	core.AssertError(t, err)
}

func TestDerive_DeriveAddressKey_Ugly(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	acct, err := DeriveAccount(master, CoinTypeLethean, 0)
	core.AssertNoError(t, err)
	// Hardened address index is rejected (BIP-44 address level is normal).
	_, err = DeriveAddressKey(acct, BranchReceive, HardenedKeyStart+1)
	core.AssertError(t, err)
}

// --- Derive Good/Bad/Ugly ---

func TestDerive_Derive_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	key, err := Derive(master, Path{CoinType: CoinTypeLethean, Account: 0, Branch: BranchReceive, Index: 0})
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint8(5), key.Depth())
}

func TestDerive_Derive_Bad(t *core.T) {
	_, err := Derive(nil, Path{CoinType: CoinTypeLethean})
	core.AssertError(t, err)
}

func TestDerive_Derive_Ugly(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	// Invalid branch in the path propagates as an error.
	_, err = Derive(master, Path{CoinType: CoinTypeLethean, Branch: 9})
	core.AssertError(t, err)
}

// --- DeriveAddress Good/Bad/Ugly ---

func TestDerive_DeriveAddress_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	addr, err := DeriveAddress(master, Path{CoinType: CoinTypeLethean, Account: 0, Branch: BranchReceive, Index: 0})
	core.AssertNoError(t, err)
	core.AssertTrue(t, addr.IsPubkeyHash())
}

func TestDerive_DeriveAddress_Bad(t *core.T) {
	_, err := DeriveAddress(nil, Path{CoinType: CoinTypeLethean})
	core.AssertError(t, err)
}

func TestDerive_DeriveAddress_Ugly(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	_, err = DeriveAddress(master, Path{CoinType: CoinTypeLethean, Branch: 7})
	core.AssertError(t, err)
}
