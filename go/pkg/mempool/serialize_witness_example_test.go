// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	core "github.com/dappcore/go"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

func ExampleSerializeWithWitness() {
	tx := primitives.Transaction{
		Version: 0,
		Inputs: []primitives.Input{
			{
				Prevout:  primitives.Outpoint{Index: 0},
				Sequence: 0xffffffff,
				Witness:  [][]byte{{0xaa}},
			},
		},
		Outputs: []primitives.Output{
			{Value: 1, Address: primitives.Address{Version: 0, Hash: make([]byte, 20)}},
		},
	}

	raw, err := SerializeWithWitness(tx)
	core.Println(err == nil && len(raw) > 0)
	// Output: true
}
