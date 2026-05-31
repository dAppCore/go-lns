// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

// coinbaseInput returns an input that spends the null outpoint.
func coinbaseInput() primitives.Input {
	return primitives.Input{Prevout: primitives.Outpoint{Index: ^uint32(0)}, Sequence: ^uint32(0)}
}

// faucetProofValue builds a faucet proof witness paying value with the given fee.
func faucetProofValue(addr []byte, version uint8, fee uint64) []byte {
	b := &proofBuilder{}
	b.u32(0).u8(0).u8(0).u8(0)
	b.varbytes([]byte{0x01}) // non-address key -> value == airdropReward
	b.u8(version)
	b.u8(uint8(len(addr)))
	b.buf = append(b.buf, addr...)
	b.varint(fee)
	b.varbytes(nil)
	return b.bytes()
}

func TestVerifyBehaviour_simpleTransfer_Good(t *testing.T) {
	// A non-coinbase tx with a single TypeNone output and a coin view that
	// resolves the input passes verification.
	prev := primitives.Outpoint{Index: 0}
	prev.TxHash[0] = 0x01
	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: prev, Witness: [][]byte{{0x00}}, Sequence: ^uint32(0)}},
		Outputs: []primitives.Output{{Value: 100, Covenant: primitives.Covenant{Type: uint8(TypeNone)}}},
	}
	view := mapView{prev: primitives.Output{Value: 100, Covenant: primitives.Covenant{Type: uint8(TypeNone)}}}
	core.AssertEqual(t, 0, VerifyCovenants(tx, view, 1, Network{}))
}

func TestVerifyBehaviour_simpleTransfer_Bad(t *testing.T) {
	// Missing coin view entry fails verification.
	prev := primitives.Outpoint{Index: 0}
	prev.TxHash[0] = 0x02
	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: prev, Sequence: ^uint32(0)}},
		Outputs: []primitives.Output{{Value: 100, Covenant: primitives.Covenant{Type: uint8(TypeNone)}}},
	}
	core.AssertEqual(t, -1, VerifyCovenants(tx, mapView{}, 1, Network{}))

	// A nil view on a non-coinbase tx fails.
	core.AssertEqual(t, -1, VerifyCovenants(tx, nil, 1, Network{}))
}

func TestVerifyBehaviour_simpleTransfer_Ugly(t *testing.T) {
	// A TypeClaim covenant on a non-coinbase tx is structurally rejected.
	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: primitives.Outpoint{Index: 1}, Sequence: ^uint32(0)}},
		Outputs: []primitives.Output{{Value: 1, Covenant: primitives.Covenant{Type: uint8(TypeClaim)}}},
	}
	core.AssertEqual(t, -1, VerifyCovenants(tx, mapView{}, 1, Network{}))
}

func TestVerifyBehaviour_coinbaseAirdrop_Good(t *testing.T) {
	addr := []byte{0xaa, 0xbb}
	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			coinbaseInput(),
			{Witness: [][]byte{faucetProofValue(addr, 0, 0)}, Sequence: ^uint32(0)},
		},
		Outputs: []primitives.Output{
			{Value: 0, Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			{
				Value:    airdropReward,
				Address:  primitives.Address{Version: 0, Hash: addr},
				Covenant: primitives.Covenant{Type: uint8(TypeNone)},
			},
		},
	}
	core.AssertTrue(t, HasSaneCovenants(tx))
	core.AssertEqual(t, 0, VerifyCovenants(tx, nil, 1, Network{}))
}

func TestVerifyBehaviour_coinbaseAirdrop_Bad(t *testing.T) {
	addr := []byte{0xaa, 0xbb}
	// Output value does not match value-fee -> fails.
	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			coinbaseInput(),
			{Witness: [][]byte{faucetProofValue(addr, 0, 0)}, Sequence: ^uint32(0)},
		},
		Outputs: []primitives.Output{
			{Value: 0, Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			{
				Value:    airdropReward - 1, // wrong value
				Address:  primitives.Address{Version: 0, Hash: addr},
				Covenant: primitives.Covenant{Type: uint8(TypeNone)},
			},
		},
	}
	core.AssertEqual(t, -1, VerifyCovenants(tx, nil, 1, Network{}))
}

func TestVerifyBehaviour_coinbaseAirdrop_Ugly(t *testing.T) {
	addr := []byte{0xaa, 0xbb}
	// Witness stack with the wrong number of items fails the sanity check.
	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			coinbaseInput(),
			{Witness: [][]byte{faucetProofValue(addr, 0, 0), {0x00}}, Sequence: ^uint32(0)},
		},
		Outputs: []primitives.Output{
			{Value: 0, Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			{
				Value:    airdropReward,
				Address:  primitives.Address{Version: 0, Hash: addr},
				Covenant: primitives.Covenant{Type: uint8(TypeNone)},
			},
		},
	}
	core.AssertFalse(t, HasSaneCovenants(tx))
	core.AssertEqual(t, -1, VerifyCovenants(tx, nil, 1, Network{}))
}

func TestVerifyBehaviour_HasSaneCovenants_coinbaseShape(t *testing.T) {
	// More inputs than outputs in a coinbase is rejected.
	tx := primitives.Transaction{
		Inputs: []primitives.Input{coinbaseInput(), coinbaseInput()},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
		},
	}
	core.AssertFalse(t, HasSaneCovenants(tx))

	// A coinbase TypeNone output carrying items is rejected.
	withItems := primitives.Transaction{
		Inputs: []primitives.Input{coinbaseInput()},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeNone), Items: [][]byte{{0x00}}}},
		},
	}
	core.AssertFalse(t, HasSaneCovenants(withItems))
}

// mapView is a CoinView backed by a map for verify tests.
type mapView map[primitives.Outpoint]primitives.Output

func (m mapView) GetOutput(o primitives.Outpoint) (primitives.Output, bool) {
	out, ok := m[o]
	return out, ok
}
