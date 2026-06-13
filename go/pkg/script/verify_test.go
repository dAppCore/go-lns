// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/primitives"
)

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

// blake160Of returns the BLAKE2b-160 digest of a compressed public key, i.e. the
// P2WPKH programme hash.
func blake160Of(t *testing.T, pub []byte) []byte {
	t.Helper()
	h, err := blake2b.New(20, nil)
	if err != nil {
		t.Fatalf("blake2b: %v", err)
	}
	_, _ = h.Write(pub)
	return h.Sum(nil)
}

// spendTx builds a single-input single-output transaction spending to addr.
func spendTx(addr primitives.Address, outValue uint64) primitives.Transaction {
	return primitives.Transaction{
		Version:  0,
		Inputs:   []primitives.Input{{Prevout: primitives.Outpoint{Index: 0}, Sequence: 0xffffffff}},
		Outputs:  []primitives.Output{{Value: outValue, Address: addr}},
		LockTime: 0,
	}
}

func TestVerify_P2WPKH_Good(t *testing.T) {
	priv, pub := testKey(t, 1)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}

	tx := spendTx(addr, 1000)
	value := uint64(2000)
	redeem := NewScript().FromPubkeyhash(pkh)

	sig, err := Sign(tx, 0, redeem, value, priv, SighashAll)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	tx.Inputs[0].Witness = [][]byte{sig, pub}

	prevOut := primitives.Output{Value: value, Address: addr}
	if err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags); err != nil {
		t.Fatalf("VerifyInput: %v", err)
	}
}

func TestVerify_P2WPKH_Bad(t *testing.T) {
	priv, pub := testKey(t, 2)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	value := uint64(2000)
	redeem := NewScript().FromPubkeyhash(pkh)

	sig, err := Sign(tx, 0, redeem, value, priv, SighashAll)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	// Wrong public key for the committed hash -> EQUALVERIFY failure.
	_, otherPub := testKey(t, 99)
	tx.Inputs[0].Witness = [][]byte{sig, otherPub}
	prevOut := primitives.Output{Value: value, Address: addr}
	if err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags); err == nil {
		t.Fatal("expected verification failure with wrong pubkey")
	}

	// Tampered signature -> EVAL_FALSE / NULLFAIL.
	bad := append([]byte(nil), sig...)
	bad[5] ^= 0xff
	tx.Inputs[0].Witness = [][]byte{bad, pub}
	if err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags); err == nil {
		t.Fatal("expected verification failure with tampered sig")
	}
}

func TestVerify_P2WPKH_Ugly(t *testing.T) {
	_, pub := testKey(t, 3)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	prevOut := primitives.Output{Value: 2000, Address: addr}

	// P2WPKH requires exactly two witness items.
	tx.Inputs[0].Witness = [][]byte{pub}
	err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags)
	assertScriptError(t, err, errWitnessProgramMismatch)

	tx.Inputs[0].Witness = [][]byte{pub, pub, pub}
	err = VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags)
	assertScriptError(t, err, errWitnessProgramMismatch)
}

func TestVerify_P2WSH_Good(t *testing.T) {
	priv, pub := testKey(t, 4)
	// Redeem script: <pubkey> OP_CHECKSIG (bare P2PK inside a witness script).
	redeem, err := FromPubkey(pub)
	if err != nil {
		t.Fatalf("FromPubkey: %v", err)
	}
	addr := primitives.Address{Version: 0, Hash: redeem.Sha3()}
	tx := spendTx(addr, 500)
	value := uint64(1000)

	sig, err := Sign(tx, 0, redeem, value, priv, SighashAll)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	// Witness stack: [sig, redeemScript].
	tx.Inputs[0].Witness = [][]byte{sig, redeem.Encode()}

	prevOut := primitives.Output{Value: value, Address: addr}
	if err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags); err != nil {
		t.Fatalf("VerifyInput P2WSH: %v", err)
	}
}

func TestVerify_P2WSH_Multisig_Good(t *testing.T) {
	priv0, pub0 := testKey(t, 5)
	priv1, pub1 := testKey(t, 6)
	keys := [][]byte{pub0, pub1}
	redeem, err := FromMultisig(2, 2, keys)
	if err != nil {
		t.Fatalf("FromMultisig: %v", err)
	}
	addr := primitives.Address{Version: 0, Hash: redeem.Sha3()}
	tx := spendTx(addr, 500)
	value := uint64(1000)

	// Signatures must be ordered to match the (sorted) key order in the script.
	sorted := [][]byte{pub0, pub1}
	sortKeys(sorted)
	privFor := map[string]*secp256k1.PrivateKey{
		string(pub0): priv0,
		string(pub1): priv1,
	}
	var sigs [][]byte
	for _, k := range sorted {
		sig, err := Sign(tx, 0, redeem, value, privFor[string(k)], SighashAll)
		if err != nil {
			t.Fatalf("Sign: %v", err)
		}
		sigs = append(sigs, sig)
	}
	// Witness: [null-dummy, sig0, sig1, redeemScript].
	witness := [][]byte{{}}
	witness = append(witness, sigs...)
	witness = append(witness, redeem.Encode())
	tx.Inputs[0].Witness = witness

	prevOut := primitives.Output{Value: value, Address: addr}
	if err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags); err != nil {
		t.Fatalf("VerifyInput P2WSH multisig: %v", err)
	}
}

func TestVerify_P2WSH_Bad(t *testing.T) {
	priv, pub := testKey(t, 7)
	redeem, _ := FromPubkey(pub)
	addr := primitives.Address{Version: 0, Hash: redeem.Sha3()}
	tx := spendTx(addr, 500)
	value := uint64(1000)
	sig, _ := Sign(tx, 0, redeem, value, priv, SighashAll)

	prevOut := primitives.Output{Value: value, Address: addr}

	// Redeem script whose hash does not match the programme.
	wrong, _ := FromPubkey(mustOtherPub(t))
	tx.Inputs[0].Witness = [][]byte{sig, wrong.Encode()}
	err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags)
	assertScriptError(t, err, errWitnessProgramMismatch)

	// Empty witness for a P2WSH programme.
	tx.Inputs[0].Witness = [][]byte{}
	err = VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags)
	assertScriptError(t, err, errWitnessProgramWitnessEmpty)
}

func TestVerify_Unspendable_Ugly(t *testing.T) {
	// Version 31 (nulldata) is unspendable.
	addr := primitives.Address{Version: 31, Hash: make([]byte, 20)}
	tx := spendTx(addr, 0)
	tx.Inputs[0].Witness = [][]byte{{0x01}}
	prevOut := primitives.Output{Value: 100, Address: addr}
	err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags)
	assertScriptError(t, err, errOpReturn)
}

func TestVerify_WrongLength_Ugly(t *testing.T) {
	// Version 0 with a non-standard (10-byte) programme is wrong-length.
	addr := primitives.Address{Version: 0, Hash: make([]byte, 10)}
	tx := spendTx(addr, 0)
	tx.Inputs[0].Witness = [][]byte{{0x01}, {0x02}}
	prevOut := primitives.Output{Value: 100, Address: addr}
	err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags)
	assertScriptError(t, err, errWitnessProgramWrongLength)
}

func TestVerify_UpgradableWitness_Ugly(t *testing.T) {
	// An unknown witness version is accepted by default but rejected under the
	// discourage flag.
	addr := primitives.Address{Version: 1, Hash: make([]byte, 32)}
	tx := spendTx(addr, 0)
	tx.Inputs[0].Witness = [][]byte{{0x01}}
	prevOut := primitives.Output{Value: 100, Address: addr}

	if err := VerifyInput(prevOut, tx, 0, VerifyNone); err != nil {
		t.Fatalf("unknown witness version should pass without discourage flag: %v", err)
	}
	err := VerifyInput(prevOut, tx, 0, VerifyDiscourageUpgradableWitnessProgram)
	assertScriptError(t, err, errDiscourageUpgradableWitness)
}

func TestVerify_StackTooLarge_Bad(t *testing.T) {
	addr := primitives.Address{Version: 0, Hash: make([]byte, 32)}
	tx := spendTx(addr, 0)
	witness := make([][]byte, MaxScriptStack+1)
	for i := range witness {
		witness[i] = []byte{0x01}
	}
	tx.Inputs[0].Witness = witness
	prevOut := primitives.Output{Value: 100, Address: addr}
	err := VerifyInput(prevOut, tx, 0, MandatoryVerifyFlags)
	assertScriptError(t, err, errStackSize)
}

func mustOtherPub(t *testing.T) []byte {
	t.Helper()
	_, pub := testKey(t, 200)
	return pub
}
