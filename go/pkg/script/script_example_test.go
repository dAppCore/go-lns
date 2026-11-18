// SPDX-License-Identifier: EUPL-1.2

package script_test

import (
	"bytes"

	core "github.com/dappcore/go"

	"github.com/dappcore/go-lns/go/pkg/script"
)

// ExampleFromPubkeyhash builds the standard pay-to-witness-pubkey-hash redeem
// template and prints its opcode count.
func ExampleFromPubkeyhash() {
	hash := bytes.Repeat([]byte{0xab}, 20)
	redeem, _ := script.FromPubkeyhash(hash)
	core.Println(len(redeem.Code()))
	// Output: 5
}

// ExampleScript_Sha3 shows that the P2WSH programme hash is the SHA3-256 of the
// redeem script and is 32 bytes wide.
func ExampleScript_Sha3() {
	key := append([]byte{0x02}, bytes.Repeat([]byte{0x05}, 32)...)
	redeem, _ := script.FromPubkey(key)
	core.Println(len(redeem.Sha3()))
	// Output: 32
}

// ExampleDecodeScript round-trips a raw script and reports that the bytes are
// preserved exactly.
func ExampleDecodeScript() {
	raw := []byte{script.OP_DUP, script.OP_BLAKE160}
	s, _ := script.DecodeScript(raw)
	core.Println(bytes.Equal(s.Encode(), raw))
	// Output: true
}
