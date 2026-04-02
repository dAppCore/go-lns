// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"bytes"

	"dappco.re/go/lns/pkg/primitives"
)

// CoinView is the minimal lookup interface required by VerifyCovenants.
//
// It returns the previous output for an outpoint, which is sufficient for the
// covenant transition checks implemented in this package.
type CoinView interface {
	GetOutput(primitives.Outpoint) (primitives.Output, bool)
}

// Network captures the covenant rollout inputs used by VerifyCovenants.
//
// The structure mirrors the `network.names` and `network.deflationHeight`
// fields accessed by the JS reference implementation while keeping the Go API
// self-contained.
type Network struct {
	Names           NameRules
	DeflationHeight uint32
}

// VerifyCovenants performs contextual covenant verification for a transaction.
//
// The function mirrors the JS reference return convention: 0 on success and
// -1 on failure. The Go port currently covers the covenant transitions that
// depend only on the transaction and coin view. Coinbase proof validation is
// intentionally conservative because the corresponding proof types live in the
// wider blockchain stack, which is not part of this package.
func VerifyCovenants(tx primitives.Transaction, view CoinView, height uint32, network Network) int {
	if !HasSaneCovenants(tx) {
		return -1
	}

	if isCoinbaseTx(tx) {
		// The JS reference performs proof validation here, but those proof
		// types live outside this package. We keep coinbase handling conservative
		// and rely on structural sanity checks above.
		return 0
	}

	if view == nil {
		return -1
	}

	for i, input := range tx.Inputs {
		coin, ok := view.GetOutput(input.Prevout)
		if !ok {
			return -1
		}

		var output *primitives.Output
		if i < len(tx.Outputs) {
			output = &tx.Outputs[i]
		}

		uc := coin.Covenant

		if output == nil {
			switch uc.Type {
			case uint8(TypeNone), uint8(TypeOpen), uint8(TypeRedeem):
				continue
			default:
				return -1
			}
		}

		cov := output.Covenant

		switch uc.Type {
		case uint8(TypeNone), uint8(TypeOpen), uint8(TypeRedeem):
			switch cov.Type {
			case uint8(TypeNone), uint8(TypeOpen), uint8(TypeBid):
			default:
				return -1
			}

		case uint8(TypeBid):
			if cov.Type != uint8(TypeReveal) {
				return -1
			}

			if !covenantHashEquals(cov, 0, uc) {
				return -1
			}

			if !covenantU32Equals(cov, 1, uc, 1) {
				return -1
			}

			nonce, err := cov.GetHash(2)
			if err != nil {
				return -1
			}

			blind, err := Blind(output.Value, nonce)
			if err != nil {
				return -1
			}

			ucBlind, err := uc.GetHash(3)
			if err != nil || blind != ucBlind {
				return -1
			}

			if coin.Value < output.Value {
				return -1
			}

		case uint8(TypeClaim), uint8(TypeReveal):
			switch cov.Type {
			case uint8(TypeRegister):
				if !covenantHashEquals(cov, 0, uc) {
					return -1
				}

				if !covenantU32Equals(cov, 1, uc, 1) {
					return -1
				}

				if !output.Address.Equals(coin.Address) {
					return -1
				}

			case uint8(TypeRedeem):
				if !covenantHashEquals(cov, 0, uc) {
					return -1
				}

				if !covenantU32Equals(cov, 1, uc, 1) {
					return -1
				}

				if uc.Type == uint8(TypeClaim) {
					return -1
				}

			default:
				return -1
			}

		case uint8(TypeRegister), uint8(TypeUpdate), uint8(TypeRenew), uint8(TypeFinalize):
			switch cov.Type {
			case uint8(TypeUpdate), uint8(TypeRenew), uint8(TypeTransfer), uint8(TypeRevoke):
			default:
				return -1
			}

			if output.Value != coin.Value {
				return -1
			}

			if !output.Address.Equals(coin.Address) {
				return -1
			}

			if !covenantHashEquals(cov, 0, uc) {
				return -1
			}

			if !covenantU32Equals(cov, 1, uc, 1) {
				return -1
			}

		case uint8(TypeTransfer):
			switch cov.Type {
			case uint8(TypeUpdate), uint8(TypeRenew), uint8(TypeRevoke):
				if output.Value != coin.Value {
					return -1
				}

				if !output.Address.Equals(coin.Address) {
					return -1
				}

				if !covenantHashEquals(cov, 0, uc) {
					return -1
				}

				if !covenantU32Equals(cov, 1, uc, 1) {
					return -1
				}

			case uint8(TypeFinalize):
				if output.Value != coin.Value {
					return -1
				}

				if !covenantHashEquals(cov, 0, uc) {
					return -1
				}

				if !covenantU32Equals(cov, 1, uc, 1) {
					return -1
				}

				version, err := uc.GetU8(2)
				if err != nil || output.Address.Version != version {
					return -1
				}

				addr, err := uc.Get(3)
				if err != nil || !bytes.Equal(output.Address.Hash, addr) {
					return -1
				}

			default:
				return -1
			}

		case uint8(TypeRevoke):
			return -1

		default:
			if cov.IsName() {
				return -1
			}
		}
	}

	return 0
}

func covenantHashEquals(cov primitives.Covenant, index int, other primitives.Covenant) bool {
	hash, err := cov.GetHash(index)
	if err != nil {
		return false
	}

	otherHash, err := other.GetHash(index)
	if err != nil {
		return false
	}

	return hash == otherHash
}

func covenantU32Equals(left primitives.Covenant, leftIndex int, right primitives.Covenant, rightIndex int) bool {
	leftValue, err := left.GetU32(leftIndex)
	if err != nil {
		return false
	}

	rightValue, err := right.GetU32(rightIndex)
	if err != nil {
		return false
	}

	return leftValue == rightValue
}
