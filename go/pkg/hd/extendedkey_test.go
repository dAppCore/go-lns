// SPDX-License-Identifier: EUPL-1.2

package hd

import (
	"testing"

	core "github.com/dappcore/go"
)

// assertKeyMatches verifies a derived node against known chain-code, private-key
// and public-key vectors.
func assertKeyMatches(t *testing.T, k *ExtendedKey, chainHex, privHex, pubHex string) {
	t.Helper()

	core.AssertEqual(t, chainHex, core.HexEncode(k.ChainCode()))

	priv, err := k.PrivateKeyBytes()
	core.AssertNoError(t, err)
	core.AssertEqual(t, privHex, core.HexEncode(priv))

	pub, err := k.PublicKeyBytes()
	core.AssertNoError(t, err)
	core.AssertEqual(t, pubHex, core.HexEncode(pub))
}

// TestExtendedKey_BIP32Vector1 walks BIP-32 Test Vector 1 end to end, asserting
// every chain code and key against the published known-answer values. This is
// the primary correctness gate for the derivation maths.
func TestExtendedKey_BIP32Vector1(t *testing.T) {
	seed := mustHex(t, vector1SeedHex)

	master, err := NewMaster(seed)
	core.AssertNoError(t, err)
	assertKeyMatches(t, master, vector1MasterChainHex, vector1MasterPrivHex, vector1MasterPubHex)

	m0h, err := master.Child(HardenedKeyStart + 0)
	core.AssertNoError(t, err)
	assertKeyMatches(t, m0h, vector1M0hChainHex, vector1M0hPrivHex, vector1M0hPubHex)

	m0h1, err := m0h.Child(1)
	core.AssertNoError(t, err)
	assertKeyMatches(t, m0h1, vector1M0h1ChainHex, vector1M0h1PrivHex, vector1M0h1PubHex)

	m0h12h, err := m0h1.Child(HardenedKeyStart + 2)
	core.AssertNoError(t, err)
	assertKeyMatches(t, m0h12h, vector1M0h12hChainHex, vector1M0h12hPrivHex, vector1M0h12hPubHex)
}

// TestExtendedKey_PublicDerivationMatchesPrivate confirms that deriving a normal
// child via a neutered (public-only) parent yields the same public key as
// deriving it privately. This exercises the EC point-addition path.
func TestExtendedKey_PublicDerivationMatchesPrivate(t *testing.T) {
	seed := mustHex(t, vector1SeedHex)
	master, err := NewMaster(seed)
	core.AssertNoError(t, err)

	m0h, err := master.Child(HardenedKeyStart + 0)
	core.AssertNoError(t, err)

	// Private path: m/0'/1 public key.
	privChild, err := m0h.Child(1)
	core.AssertNoError(t, err)
	privPub, err := privChild.PublicKeyBytes()
	core.AssertNoError(t, err)

	// Public path: neuter m/0', then derive normal child 1.
	pubParent, err := m0h.Neuter()
	core.AssertNoError(t, err)
	core.AssertTrue(t, !pubParent.IsPrivate())

	pubChild, err := pubParent.Child(1)
	core.AssertNoError(t, err)
	pubChildPub, err := pubChild.PublicKeyBytes()
	core.AssertNoError(t, err)

	core.AssertEqual(t, core.HexEncode(privPub), core.HexEncode(pubChildPub))
	core.AssertEqual(t, vector1M0h1PubHex, core.HexEncode(pubChildPub))
}

// TestExtendedKey_Determinism confirms the same seed always derives the same key.
func TestExtendedKey_Determinism(t *testing.T) {
	seed := mustHex(t, vector1SeedHex)

	first, err := NewMaster(seed)
	core.AssertNoError(t, err)
	second, err := NewMaster(seed)
	core.AssertNoError(t, err)

	p1, err := first.PrivateKeyBytes()
	core.AssertNoError(t, err)
	p2, err := second.PrivateKeyBytes()
	core.AssertNoError(t, err)
	core.AssertEqual(t, core.HexEncode(p1), core.HexEncode(p2))
}

// TestExtendedKey_HardenedDivergesFromNormal confirms hardened and normal
// children at the same numeric index produce different keys.
func TestExtendedKey_HardenedDivergesFromNormal(t *testing.T) {
	seed := mustHex(t, vector1SeedHex)
	master, err := NewMaster(seed)
	core.AssertNoError(t, err)

	normal, err := master.Child(0)
	core.AssertNoError(t, err)
	hardened, err := master.Child(HardenedKeyStart + 0)
	core.AssertNoError(t, err)

	np, err := normal.PrivateKeyBytes()
	core.AssertNoError(t, err)
	hp, err := hardened.PrivateKeyBytes()
	core.AssertNoError(t, err)

	core.AssertTrue(t, core.HexEncode(np) != core.HexEncode(hp))
}

// --- NewMaster Good/Bad/Ugly ---

func TestExtendedKey_NewMaster_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	core.AssertTrue(t, master.IsPrivate())
}

func TestExtendedKey_NewMaster_Bad(t *core.T) {
	// Seed too short (below SeedMinLen) must error.
	_, err := NewMaster([]byte{1, 2, 3})
	core.AssertError(t, err)
}

func TestExtendedKey_NewMaster_Ugly(t *core.T) {
	// Seed too long (above SeedMaxLen) must error.
	_, err := NewMaster(make([]byte, SeedMaxLen+1))
	core.AssertError(t, err)
}

// --- Child Good/Bad/Ugly ---

func TestExtendedKey_Child_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	child, err := master.Child(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint8(1), child.Depth())
	core.AssertEqual(t, uint32(0), child.ChildNumber())
}

func TestExtendedKey_Child_Bad(t *core.T) {
	// Hardened derivation from a public-only node must fail.
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	_, err = pub.Child(HardenedKeyStart + 0)
	core.AssertError(t, err)
}

func TestExtendedKey_Child_Ugly(t *core.T) {
	// Nil receiver must fail rather than panic.
	var k *ExtendedKey
	_, err := k.Child(0)
	core.AssertError(t, err)
}

// --- Neuter Good/Bad/Ugly ---

func TestExtendedKey_Neuter_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	core.AssertTrue(t, !pub.IsPrivate())

	gotPub, err := pub.PublicKeyBytes()
	core.AssertNoError(t, err)
	core.AssertEqual(t, vector1MasterPubHex, core.HexEncode(gotPub))
}

func TestExtendedKey_Neuter_Bad(t *core.T) {
	var k *ExtendedKey
	_, err := k.Neuter()
	core.AssertError(t, err)
}

func TestExtendedKey_Neuter_Ugly(t *core.T) {
	// Neutering an already-public node yields an equivalent public node.
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	pub2, err := pub.Neuter()
	core.AssertNoError(t, err)
	core.AssertTrue(t, !pub2.IsPrivate())
}

// --- PrivateKeyBytes Good/Bad/Ugly ---

func TestExtendedKey_PrivateKeyBytes_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	priv, err := master.PrivateKeyBytes()
	core.AssertNoError(t, err)
	core.AssertEqual(t, vector1MasterPrivHex, core.HexEncode(priv))
}

func TestExtendedKey_PrivateKeyBytes_Bad(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	_, err = pub.PrivateKeyBytes()
	core.AssertError(t, err)
}

func TestExtendedKey_PrivateKeyBytes_Ugly(t *core.T) {
	var k *ExtendedKey
	_, err := k.PrivateKeyBytes()
	core.AssertError(t, err)
}

// --- PublicKeyBytes Good/Bad/Ugly ---

func TestExtendedKey_PublicKeyBytes_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.PublicKeyBytes()
	core.AssertNoError(t, err)
	core.AssertEqual(t, compressedPubKeyLen, len(pub))
	core.AssertEqual(t, vector1MasterPubHex, core.HexEncode(pub))
}

func TestExtendedKey_PublicKeyBytes_Bad(t *core.T) {
	var k *ExtendedKey
	_, err := k.PublicKeyBytes()
	core.AssertError(t, err)
}

func TestExtendedKey_PublicKeyBytes_Ugly(t *core.T) {
	// A neutered node still yields its compressed public key.
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	got, err := pub.PublicKeyBytes()
	core.AssertNoError(t, err)
	core.AssertEqual(t, vector1MasterPubHex, core.HexEncode(got))
}

// --- Address Good/Bad/Ugly ---

func TestExtendedKey_Address_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	addr := master.Address()
	core.AssertEqual(t, uint8(0), addr.Version)
	core.AssertEqual(t, addressHashLen, len(addr.Hash))
	core.AssertTrue(t, addr.IsPubkeyHash())
}

func TestExtendedKey_Address_Bad(t *core.T) {
	var k *ExtendedKey
	addr := k.Address()
	// Nil receiver yields the zero address (no hash).
	core.AssertEqual(t, 0, len(addr.Hash))
}

func TestExtendedKey_Address_Ugly(t *core.T) {
	// Address is deterministic: same key, same address hash, and the public
	// node yields the same address as the private node.
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	core.AssertEqual(t, core.HexEncode(master.Address().Hash), core.HexEncode(pub.Address().Hash))
}

// --- IsPrivate / Depth / ChildNumber / ChainCode Good/Bad/Ugly ---

func TestExtendedKey_IsPrivate_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	core.AssertTrue(t, master.IsPrivate())
}

func TestExtendedKey_IsPrivate_Bad(t *core.T) {
	var k *ExtendedKey
	core.AssertTrue(t, !k.IsPrivate())
}

func TestExtendedKey_IsPrivate_Ugly(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	pub, err := master.Neuter()
	core.AssertNoError(t, err)
	core.AssertTrue(t, !pub.IsPrivate())
}

func TestExtendedKey_Depth_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint8(0), master.Depth())
	child, err := master.Child(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint8(1), child.Depth())
}

func TestExtendedKey_Depth_Bad(t *core.T) {
	var k *ExtendedKey
	core.AssertEqual(t, uint8(0), k.Depth())
}

func TestExtendedKey_Depth_Ugly(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	d2, err := master.Child(HardenedKeyStart)
	core.AssertNoError(t, err)
	d3, err := d2.Child(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint8(2), d3.Depth())
}

func TestExtendedKey_ChildNumber_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	child, err := master.Child(7)
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint32(7), child.ChildNumber())
}

func TestExtendedKey_ChildNumber_Bad(t *core.T) {
	var k *ExtendedKey
	core.AssertEqual(t, uint32(0), k.ChildNumber())
}

func TestExtendedKey_ChildNumber_Ugly(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	child, err := master.Child(HardenedKeyStart + 9)
	core.AssertNoError(t, err)
	core.AssertEqual(t, HardenedKeyStart+9, child.ChildNumber())
}

func TestExtendedKey_ChainCode_Good(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	core.AssertEqual(t, vector1MasterChainHex, core.HexEncode(master.ChainCode()))
}

func TestExtendedKey_ChainCode_Bad(t *core.T) {
	var k *ExtendedKey
	core.AssertTrue(t, k.ChainCode() == nil)
}

func TestExtendedKey_ChainCode_Ugly(t *core.T) {
	master, err := NewMaster(mustHex(t, vector1SeedHex))
	core.AssertNoError(t, err)
	// Mutating the returned chain code must not affect internal state.
	cc := master.ChainCode()
	if len(cc) > 0 {
		cc[0] ^= 0xff
	}
	core.AssertEqual(t, vector1MasterChainHex, core.HexEncode(master.ChainCode()))
}
