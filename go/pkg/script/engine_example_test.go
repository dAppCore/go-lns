// SPDX-License-Identifier: EUPL-1.2

package script_test

import (
	core "dappco.re/go"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/primitives"
	"dappco.re/go/lns/pkg/script"
)

// ExampleVerifyInput demonstrates a complete pay-to-witness-pubkey-hash spend:
// derive a key, build the address, sign the input, attach the witness, and
// verify it.
func ExampleVerifyInput() {
	// Deterministic key for a reproducible example.
	var seed [32]byte
	for i := range seed {
		seed[i] = byte(i + 1)
	}
	priv := secp256k1.PrivKeyFromBytes(seed[:])
	pub := priv.PubKey().SerializeCompressed()

	// P2WPKH programme hash = BLAKE2b-160 of the compressed public key.
	h, _ := blake2b.New(20, nil)
	_, _ = h.Write(pub)
	addr := primitives.Address{Version: 0, Hash: h.Sum(nil)}

	// A transaction spending one input to one output.
	tx := primitives.Transaction{
		Version:  0,
		Inputs:   []primitives.Input{{Prevout: primitives.Outpoint{Index: 0}, Sequence: 0xffffffff}},
		Outputs:  []primitives.Output{{Value: 1000, Address: addr}},
		LockTime: 0,
	}

	value := uint64(2000)
	redeem := script.NewScript().FromPubkeyhash(addr.Hash)
	sig, _ := script.Sign(tx, 0, redeem, value, priv, script.SighashAll)
	tx.Inputs[0].Witness = [][]byte{sig, pub}

	prevOut := primitives.Output{Value: value, Address: addr}
	err := script.VerifyInput(prevOut, tx, 0, script.MandatoryVerifyFlags)
	core.Println(err == nil)
	// Output: true
}

// ExampleScript_Execute runs a small arithmetic script and reports the result of
// 2 + 3 left on the stack.
func ExampleScript_Execute() {
	redeem, _ := script.FromOps([]script.Opcode{
		script.OpFromOp(script.OP_2),
		script.OpFromOp(script.OP_3),
		script.OpFromOp(script.OP_ADD),
	})
	stack := script.NewStack()
	_ = redeem.Execute(stack, script.VerifyNone, nil)
	n, _ := stack.Num(-1, false, 4)
	core.Println(int(n))
	// Output: 5
}
