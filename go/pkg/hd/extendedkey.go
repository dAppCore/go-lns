// SPDX-License-Identifier: EUPL-1.2

package hd

import (
	"math/big"

	core "github.com/dappcore/go"
	"github.com/dappcore/go-lns/go/pkg/primitives"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/blake2b"
)

const (
	opNewMaster        = "hd.NewMaster"
	opChild            = "hd.ExtendedKey.Child"
	opNeuter           = "hd.ExtendedKey.Neuter"
	opPublicKeyBytes   = "hd.ExtendedKey.PublicKeyBytes"
	opPrivateKeyBytes  = "hd.ExtendedKey.PrivateKeyBytes"
	opAddress          = "hd.ExtendedKey.Address"
	opExtendedKeyValid = "hd.ExtendedKey"
)

// ExtendedKey is a BIP-32 hierarchical-deterministic node. It pairs a key with
// its chain code and tracks its position in the tree (depth and child number).
// A node is either private (key is a 32-byte scalar) or public-only (key is a
// 33-byte compressed point); call IsPrivate to distinguish them.
//
//	master, _ := hd.NewMaster(seed)
//	if master.IsPrivate() { /* can derive hardened children */ }
type ExtendedKey struct {
	// key holds the raw key material: a 32-byte private scalar when isPrivate
	// is true, otherwise a 33-byte compressed public key.
	key []byte
	// chainCode is the 32-byte BIP-32 chain code mixed into child derivation.
	chainCode []byte
	// depth is the number of derivations from the master node (0 for master).
	depth uint8
	// childNumber is the index used to derive this node from its parent.
	childNumber uint32
	// isPrivate reports whether key holds private material.
	isPrivate bool
}

// NewMaster derives the master extended key from a binary seed using
// HMAC-SHA512 with the fixed BIP-32 key. The seed must be between SeedMinLen and
// SeedMaxLen bytes. It fails if the derived scalar is zero or not less than the
// curve order, which is cryptographically improbable for a random seed.
//
//	master, err := hd.NewMaster(seed)
func NewMaster(seed []byte) (*ExtendedKey, error) {
	if len(seed) < SeedMinLen || len(seed) > SeedMaxLen {
		return nil, core.E(opNewMaster, "seed length out of range", nil)
	}

	r := core.HMAC("sha512", masterKeyHMACKey, seed)
	if !r.OK {
		return nil, core.E(opNewMaster, "hmac-sha512 master derivation failed", r.Value.(error))
	}
	sum := r.Value.([]byte)

	il := sum[:keyLen]
	ir := sum[keyLen:]

	secret := parse256(il)
	if secret.Sign() == 0 || secret.Cmp(curveOrder) >= 0 {
		return nil, core.E(opNewMaster, "invalid master key (out of range)", nil)
	}

	return &ExtendedKey{
		key:         ser256(secret),
		chainCode:   cloneBytes(ir),
		depth:       0,
		childNumber: 0,
		isPrivate:   true,
	}, nil
}

// IsPrivate reports whether the node carries private key material and can
// therefore derive hardened children and sign.
func (k *ExtendedKey) IsPrivate() bool {
	return k != nil && k.isPrivate
}

// Depth returns the node's distance from the master key (0 for the master).
func (k *ExtendedKey) Depth() uint8 {
	if k == nil {
		return 0
	}
	return k.depth
}

// ChildNumber returns the index this node was derived with from its parent.
func (k *ExtendedKey) ChildNumber() uint32 {
	if k == nil {
		return 0
	}
	return k.childNumber
}

// ChainCode returns a copy of the node's 32-byte chain code.
func (k *ExtendedKey) ChainCode() []byte {
	if k == nil {
		return nil
	}
	return cloneBytes(k.chainCode)
}

// Child derives the child extended key at the given index. Indices at or above
// HardenedKeyStart derive hardened children, which require a private parent.
// Normal children may be derived from either a private or a public node. It
// fails on a hardened request against a public node, or in the improbable event
// that the derived scalar is invalid (the BIP-32 retry case).
//
//	hardened, err := parent.Child(hd.HardenedKeyStart + 0)
//	normal, err := parent.Child(5)
func (k *ExtendedKey) Child(index uint32) (*ExtendedKey, error) {
	if k == nil {
		return nil, core.E(opChild, "nil extended key", nil)
	}
	if k.depth == 255 {
		return nil, core.E(opChild, "maximum derivation depth reached", nil)
	}

	hardened := index >= HardenedKeyStart
	if hardened && !k.isPrivate {
		return nil, core.E(opChild, "cannot derive hardened child from public key", nil)
	}

	data := make([]byte, 0, compressedPubKeyLen+4)
	if hardened {
		// 0x00 || ser256(parentKey) || ser32(index)
		data = append(data, 0x00)
		data = append(data, k.key...)
	} else {
		// serP(parentPubKey) || ser32(index)
		pub, err := k.publicKeyBytes()
		if err != nil {
			return nil, core.E(opChild, "derive parent public key", err)
		}
		data = append(data, pub...)
	}
	data = append(data, ser32(index)...)

	r := core.HMAC("sha512", k.chainCode, data)
	if !r.OK {
		return nil, core.E(opChild, "hmac-sha512 child derivation failed", r.Value.(error))
	}
	sum := r.Value.([]byte)

	il := sum[:keyLen]
	ir := sum[keyLen:]

	ilNum := parse256(il)
	if ilNum.Cmp(curveOrder) >= 0 {
		// I_L >= n: per BIP-32 this index is invalid; the caller should try
		// the next index.
		return nil, core.E(opChild, "derived scalar not less than curve order; try next index", nil)
	}

	child := &ExtendedKey{
		chainCode:   cloneBytes(ir),
		depth:       k.depth + 1,
		childNumber: index,
		isPrivate:   k.isPrivate,
	}

	if k.isPrivate {
		// childKey = (I_L + parentKey) mod n
		parentNum := parse256(k.key)
		childNum := new(big.Int).Add(ilNum, parentNum)
		childNum.Mod(childNum, curveOrder)
		if childNum.Sign() == 0 {
			return nil, core.E(opChild, "derived private key is zero; try next index", nil)
		}
		child.key = ser256(childNum)
		return child, nil
	}

	// Public parent: childPub = point(I_L) + parentPub.
	childPub, err := addPublicScalar(k.key, il)
	if err != nil {
		return nil, core.E(opChild, "derive public child point", err)
	}
	child.key = childPub
	return child, nil
}

// Neuter returns the public-only counterpart of the node, discarding private
// material while preserving the chain code and tree position. Neutering a node
// that is already public returns an equivalent copy.
//
//	pub, err := acct.Neuter()
func (k *ExtendedKey) Neuter() (*ExtendedKey, error) {
	if k == nil {
		return nil, core.E(opNeuter, "nil extended key", nil)
	}
	if !k.isPrivate {
		return &ExtendedKey{
			key:         cloneBytes(k.key),
			chainCode:   cloneBytes(k.chainCode),
			depth:       k.depth,
			childNumber: k.childNumber,
			isPrivate:   false,
		}, nil
	}

	pub, err := k.publicKeyBytes()
	if err != nil {
		return nil, core.E(opNeuter, "derive public key", err)
	}
	return &ExtendedKey{
		key:         pub,
		chainCode:   cloneBytes(k.chainCode),
		depth:       k.depth,
		childNumber: k.childNumber,
		isPrivate:   false,
	}, nil
}

// PrivateKeyBytes returns a copy of the 32-byte private scalar. It fails for a
// public-only node.
//
//	priv, err := key.PrivateKeyBytes()
func (k *ExtendedKey) PrivateKeyBytes() ([]byte, error) {
	if k == nil || !k.isPrivate {
		return nil, core.E(opPrivateKeyBytes, "no private key available", nil)
	}
	return cloneBytes(k.key), nil
}

// PublicKeyBytes returns the 33-byte compressed secp256k1 public key for the
// node, deriving it from the private scalar when necessary.
//
//	pub, err := key.PublicKeyBytes()
func (k *ExtendedKey) PublicKeyBytes() ([]byte, error) {
	if k == nil {
		return nil, core.E(opPublicKeyBytes, "nil extended key", nil)
	}
	return k.publicKeyBytes()
}

// publicKeyBytes is the unexported helper backing PublicKeyBytes and the
// child-derivation public-key path.
func (k *ExtendedKey) publicKeyBytes() ([]byte, error) {
	if !k.isPrivate {
		return cloneBytes(k.key), nil
	}
	priv := secp256k1.PrivKeyFromBytes(k.key)
	return priv.PubKey().SerializeCompressed(), nil
}

// Address returns the version-0 witness address for the node: a 20-byte
// BLAKE2b-160 digest of the compressed public key, wrapped as
// primitives.Address. This matches HSD's keyring getKeyAddress and produces a
// hash compatible with pkg/wallet's primitives.Address receive address type.
//
//	addr := key.Address()
func (k *ExtendedKey) Address() primitives.Address {
	if k == nil {
		return primitives.Address{}
	}
	pub, err := k.publicKeyBytes()
	if err != nil {
		return primitives.Address{}
	}
	return primitives.Address{Version: 0, Hash: addressHash(pub)}
}

// addressHash computes the BLAKE2b-160 digest of a compressed public key, the
// HSD pubkey-hash used for version-0 P2WPKH witness programmes.
func addressHash(pub []byte) []byte {
	h, err := blake2b.New(addressHashLen, nil)
	if err != nil {
		return nil
	}
	h.Write(pub)
	return h.Sum(nil)
}

// addPublicScalar computes serP(point(scalar) + parentPub) for normal public
// child derivation. It parses the parent compressed key, performs the EC point
// addition, and returns the compressed sum. It fails if the parent key is
// invalid or the sum is the point at infinity (the BIP-32 retry case).
func addPublicScalar(parentPub, scalar []byte) ([]byte, error) {
	parent, err := secp256k1.ParsePubKey(parentPub)
	if err != nil {
		return nil, core.E("hd.addPublicScalar", "parse parent public key", err)
	}

	var ilScalar secp256k1.ModNScalar
	if overflow := ilScalar.SetByteSlice(scalar); overflow {
		return nil, core.E("hd.addPublicScalar", "scalar overflows curve order", nil)
	}
	if ilScalar.IsZero() {
		return nil, core.E("hd.addPublicScalar", "scalar is zero; try next index", nil)
	}

	var ilPoint secp256k1.JacobianPoint
	secp256k1.ScalarBaseMultNonConst(&ilScalar, &ilPoint)

	var parentPoint secp256k1.JacobianPoint
	parent.AsJacobian(&parentPoint)

	var sum secp256k1.JacobianPoint
	secp256k1.AddNonConst(&ilPoint, &parentPoint, &sum)
	if (sum.X.IsZero() && sum.Y.IsZero()) || sum.Z.IsZero() {
		return nil, core.E("hd.addPublicScalar", "derived point at infinity; try next index", nil)
	}

	sum.ToAffine()
	childPub := secp256k1.NewPublicKey(&sum.X, &sum.Y)
	return childPub.SerializeCompressed(), nil
}

// cloneBytes returns a copy of b so callers cannot mutate internal key state.
func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}
