// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/primitives"
	"dappco.re/go/lns/pkg/script"
)

// These helpers build real, fully-signed P2WPKH transactions so the pool's
// validation exercises the genuine script.VerifyInput path rather than mocks,
// matching the testnet-as-harness philosophy in the LNS RFC §15. They follow the
// same construction the script package's own tests use (pkg/script/verify_test.go).

// testKey derives a deterministic secp256k1 key pair from a seed byte.
func testKey(t *testing.T, seed byte) (*secp256k1.PrivateKey, []byte) {
	t.Helper()
	var raw [32]byte
	for i := range raw {
		raw[i] = byte(int(seed)*13 + i + 1)
	}
	priv := secp256k1.PrivKeyFromBytes(raw[:])
	return priv, priv.PubKey().SerializeCompressed()
}

// blake160 returns the BLAKE2b-160 digest of a compressed public key: the
// P2WPKH programme hash.
func blake160(t *testing.T, pub []byte) []byte {
	t.Helper()
	h, err := blake2b.New(20, nil)
	if err != nil {
		t.Fatalf("blake2b: %v", err)
	}
	_, _ = h.Write(pub)
	return h.Sum(nil)
}

// fundedCoin returns a P2WPKH address controlled by the seed key, a chain coin
// set holding one output of the given value at the given outpoint, and the
// private key. The outpoint is what a spending transaction references.
func fundedCoin(t *testing.T, seed byte, value uint64, outpoint primitives.Outpoint) (primitives.Address, *secp256k1.PrivateKey, *SpentSet) {
	t.Helper()
	priv, pub := testKey(t, seed)
	addr := primitives.Address{Version: 0, Hash: blake160(t, pub)}

	set := NewSpentSet()
	set.Add(outpoint, primitives.Output{Value: value, Address: addr})
	return addr, priv, set
}

// signedSpend builds a single-input transaction spending `outpoint` (value
// prevValue, controlled by priv) to `to` for `outValue`, leaving prevValue-outValue
// as the fee. The witness is a valid P2WPKH signature.
func signedSpend(t *testing.T, priv *secp256k1.PrivateKey, outpoint primitives.Outpoint, prevValue, outValue uint64, to primitives.Address) primitives.Transaction {
	t.Helper()

	pub := priv.PubKey().SerializeCompressed()
	pkh := blake160(t, pub)

	tx := primitives.Transaction{
		Version:  0,
		Inputs:   []primitives.Input{{Prevout: outpoint, Sequence: 0xffffffff}},
		Outputs:  []primitives.Output{{Value: outValue, Address: to}},
		LockTime: 0,
	}

	redeem := script.NewScript().FromPubkeyhash(pkh)
	sig, err := script.Sign(tx, 0, redeem, prevValue, priv, script.SighashAll)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	tx.Inputs[0].Witness = [][]byte{sig, pub}
	return tx
}

// outpoint returns an outpoint with a recognisable hash derived from a tag byte.
func outpoint(tag byte, index uint32) primitives.Outpoint {
	var h primitives.Hash
	for i := range h {
		h[i] = tag
	}
	return primitives.Outpoint{TxHash: h, Index: index}
}
