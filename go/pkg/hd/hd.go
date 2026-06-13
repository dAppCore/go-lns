// SPDX-License-Identifier: EUPL-1.2

// Package hd implements BIP-32 hierarchical-deterministic key derivation and
// BIP-44 account path derivation for the Lethean Name System, matching the HSD
// (Handshake) wallet key scheme that this repository ports.
//
// # Key scheme
//
// The scheme is inherited from HSD (the upstream Handshake daemon, see
// docs/js-wallet/masterkey.js and docs/js-primitives/keyring.js):
//
//   - Curve: secp256k1. Public keys are the 33-byte compressed form
//     (secp256k1.PublicKey.SerializeCompressed).
//   - Derivation: standard BIP-32. A node carries a 32-byte private (or 33-byte
//     compressed public) key alongside a 32-byte chain code. Child derivation
//     uses HMAC-SHA512(chainCode, data); the hardened boundary is index
//     2^31 (HardenedKeyStart).
//   - Account paths: BIP-44, "m/44'/coin'/account'/branch/index". The purpose
//     (44'), coin type, and account are hardened; branch (0 = receive,
//     1 = change) and index are normal. See docs/js-wallet/path.js and
//     docs/js-wallet/account.js (deriveKey → deriveAccount(44, type, account)
//     → derive(branch) → derive(index)).
//   - Addresses: a version-0 witness programme whose 20-byte hash is the
//     BLAKE2b-160 digest of the 33-byte compressed public key, encoded as
//     primitives.Address (see docs/js-primitives/keyring.js getKeyHash). This
//     is the same primitives.Address type consumed by pkg/wallet, so keys
//     derived here produce wallet-compatible receive addresses.
//
// # Compatibility note
//
// The current pkg/wallet scaffold derives placeholder receive addresses with a
// SHA3-based hash and does not use a curve. This package implements the real
// HSD key scheme so callers (a future chain-backed wallet) can derive spendable
// secp256k1 keys whose primitives.Address matches the on-chain P2WPKH format.
// The address output type (primitives.Address, version 0, 20-byte hash) is
// identical to pkg/wallet's; the hash algorithm follows HSD (BLAKE2b-160) rather
// than the scaffold's SHA3 placeholder.
//
// # Example
//
//	master, _ := hd.NewMaster(seed)
//	acct, _ := hd.DeriveAccount(master, hd.CoinTypeLethean, 0)
//	key, _ := hd.DeriveAddressKey(acct, hd.BranchReceive, 0)
//	addr := key.Address()
package hd

import (
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// HardenedKeyStart is the index of the first hardened child key (2^31).
// Child indices at or above this value derive hardened keys.
//
//	hardened := hd.HardenedKeyStart + 0 // first hardened child
const HardenedKeyStart uint32 = 0x80000000

// Bip44Purpose is the BIP-44 purpose level value (44), used hardened as the
// first path element under the master key.
const Bip44Purpose uint32 = 44

// CoinTypeLethean is the BIP-44 coin type for the Lethean name chain. HSD/
// Handshake registered SLIP-44 coin type 5353; the Lethean sidechain inherits
// it via the network keyPrefix.coinType (docs/js-wallet/path.js).
const CoinTypeLethean uint32 = 5353

// Branch identifiers for BIP-44 derivation. Branch 0 is the external (receive)
// chain; branch 1 is the internal (change) chain.
const (
	// BranchReceive is the external address branch (incoming funds).
	BranchReceive uint32 = 0
	// BranchChange is the internal address branch (change outputs).
	BranchChange uint32 = 1
)

const (
	// SeedMinLen is the minimum allowed master-seed length in bytes (BIP-32
	// requires 128 bits of entropy).
	SeedMinLen = 16
	// SeedMaxLen is the maximum allowed master-seed length in bytes (BIP-32
	// allows up to 512 bits of entropy).
	SeedMaxLen = 64
	// keyLen is the length in bytes of a serialised private key or X-only
	// scalar.
	keyLen = 32
	// chainCodeLen is the length in bytes of a BIP-32 chain code.
	chainCodeLen = 32
	// compressedPubKeyLen is the length in bytes of a compressed secp256k1
	// public key.
	compressedPubKeyLen = 33
	// addressHashLen is the length in bytes of a version-0 witness pubkey hash
	// (BLAKE2b-160).
	addressHashLen = 20
)

// masterKeyHMACKey is the HMAC-SHA512 key used to derive the master node from a
// seed. BIP-32 fixes this to the ASCII string "Bitcoin seed"; HSD reuses the
// same constant (its hd/private.js seeds with the identical value), so Lethean
// master keys are derived with it for cross-implementation compatibility.
var masterKeyHMACKey = []byte("Bitcoin seed")

// curveOrder is the order (n) of the secp256k1 group. Child scalars are reduced
// modulo this value during BIP-32 derivation.
var curveOrder = secp256k1.Params().N

// ser32 serialises a uint32 as a 4-byte big-endian slice, as required by the
// BIP-32 child-derivation HMAC input.
//
//	data := hd.ser32(0) // []byte{0,0,0,0}
func ser32(i uint32) []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

// parse256 interprets a 32-byte big-endian slice as a non-negative integer, the
// inverse of the BIP-32 ser256 operation.
func parse256(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

// ser256 serialises a big integer into a fixed 32-byte big-endian slice,
// left-padding with zeroes. Values wider than 32 bytes are truncated to the low
// 32 bytes, which never occurs for in-range scalars.
func ser256(x *big.Int) []byte {
	out := make([]byte, keyLen)
	x.FillBytes(out)
	return out
}
