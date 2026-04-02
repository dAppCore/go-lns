// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"crypto/rand"
	"math/big"

	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/primitives"
)

const grindAttempts = 500000

var grindNameChars = []byte("abcdefghijklmnopqrstuvwxyz")

func randomName(size int) (string, error) {
	if size <= 0 || size > MaxNameSize {
		return "", core.E("covenant.GrindName", "invalid size", nil)
	}

	buf := make([]byte, size)

	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(grindNameChars))))
		if err != nil {
			return "", core.E("covenant.GrindName", "random source failure", err)
		}

		buf[i] = grindNameChars[n.Int64()]
	}

	return string(buf), nil
}

// GrindName searches for an available name that satisfies the current rollout
// and reservation rules. It mirrors the JS helper used by reference tests.
func GrindName(size int, height uint32, rules NameRules) (string, error) {
	if size <= 0 || size > MaxNameSize {
		return "", core.E("covenant.GrindName", "invalid size", nil)
	}

	if height < rules.AuctionStart {
		height = rules.AuctionStart
	}

	for i := 0; i < grindAttempts; i++ {
		name, err := randomName(size)
		if err != nil {
			return "", err
		}

		if !VerifyString(name) {
			continue
		}

		hash, err := HashString(name)
		if err != nil {
			continue
		}

		if !HasRollout(hash, height, rules) {
			continue
		}

		if HasReservedHash(hash) {
			continue
		}

		return name, nil
	}

	return "", core.E("covenant.GrindName", "could not find available name", nil)
}

// CountOpens counts OPEN covenants in the transaction outputs.
func CountOpens(tx primitives.Transaction) int {
	total := 0

	for _, output := range tx.Outputs {
		if output.Covenant.IsOpen() {
			total++
		}
	}

	return total
}

// CountUpdates counts UPDATE-style covenants in the transaction outputs.
func CountUpdates(tx primitives.Transaction) int {
	total := 0

	for _, output := range tx.Outputs {
		switch CovenantType(output.Covenant.Type) {
		case TypeClaim, TypeOpen, TypeUpdate, TypeTransfer, TypeRevoke:
			total++
		}
	}

	return total
}

// CountRenewals counts REGISTER, RENEW, and FINALIZE covenants in the transaction outputs.
func CountRenewals(tx primitives.Transaction) int {
	total := 0

	for _, output := range tx.Outputs {
		switch CovenantType(output.Covenant.Type) {
		case TypeRegister, TypeRenew, TypeFinalize:
			total++
		}
	}

	return total
}
