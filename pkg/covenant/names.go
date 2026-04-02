// SPDX-License-Identifier: EUPL-1.2

package covenant

import "dappco.re/go/lns/pkg/primitives"

func nameHashFromCovenant(c primitives.Covenant) (primitives.Hash, bool) {
	hash, err := c.GetHash(0)
	if err != nil {
		return primitives.Hash{}, false
	}

	return hash, true
}

func nameSetHas(set map[primitives.Hash]struct{}, hash primitives.Hash) bool {
	if set == nil {
		return false
	}

	_, ok := set[hash]
	return ok
}

func nameSetAdd(set map[primitives.Hash]struct{}, hash primitives.Hash) {
	if set == nil {
		return
	}

	set[hash] = struct{}{}
}

func nameSetDelete(set map[primitives.Hash]struct{}, hash primitives.Hash) {
	if set == nil {
		return
	}

	delete(set, hash)
}

// HasNames reports whether any name covenant in the transaction references a
// hash that is already present in the set.
func HasNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) bool {
	for _, output := range tx.Outputs {
		if !output.Covenant.IsName() {
			continue
		}

		hash, ok := nameHashFromCovenant(output.Covenant)
		if !ok {
			continue
		}

		switch CovenantType(output.Covenant.Type) {
		case TypeClaim, TypeOpen, TypeRegister, TypeUpdate, TypeRenew, TypeTransfer, TypeFinalize, TypeRevoke:
			if nameSetHas(set, hash) {
				return true
			}
		}
	}

	return false
}

// AddNames adds the transaction's name covenants to the set.
func AddNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) {
	for _, output := range tx.Outputs {
		if !output.Covenant.IsName() {
			continue
		}

		hash, ok := nameHashFromCovenant(output.Covenant)
		if !ok {
			continue
		}

		switch CovenantType(output.Covenant.Type) {
		case TypeClaim, TypeOpen, TypeRegister, TypeUpdate, TypeRenew, TypeTransfer, TypeFinalize, TypeRevoke:
			nameSetAdd(set, hash)
		}
	}
}

// RemoveNames removes the transaction's name covenants from the set.
func RemoveNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) {
	for _, output := range tx.Outputs {
		if !output.Covenant.IsName() {
			continue
		}

		hash, ok := nameHashFromCovenant(output.Covenant)
		if !ok {
			continue
		}

		switch CovenantType(output.Covenant.Type) {
		case TypeClaim, TypeOpen, TypeRegister, TypeUpdate, TypeRenew, TypeTransfer, TypeFinalize, TypeRevoke:
			nameSetDelete(set, hash)
		}
	}
}

func isCoinbaseTx(tx primitives.Transaction) bool {
	return len(tx.Inputs) > 0 && tx.Inputs[0].IsCoinbase()
}

func validNameHashAndLabel(c primitives.Covenant, hashIndex, nameIndex int) bool {
	hash, err := c.GetHash(hashIndex)
	if err != nil {
		return false
	}

	name, err := c.GetString(nameIndex)
	if err != nil {
		return false
	}

	key, err := HashString(name)
	if err != nil {
		return false
	}

	return key == hash
}

// HasSaneCovenants performs structural covenant validation on a transaction.
//
// The helper mirrors the JS rules.js sanity checks that only depend on the
// transaction contents. Contextual witness and proof validation is handled by
// higher-level consensus logic.
func HasSaneCovenants(tx primitives.Transaction) bool {
	if isCoinbaseTx(tx) {
		if len(tx.Inputs) > len(tx.Outputs) {
			return false
		}

		for i, output := range tx.Outputs {
			cov := output.Covenant

			switch CovenantType(cov.Type) {
			case TypeNone:
				if cov.Len() != 0 {
					return false
				}

				if i > 0 && i < len(tx.Inputs) {
					input := tx.Inputs[i]
					if len(input.Witness) != 1 {
						return false
					}
				}

			case TypeClaim:
				if i == 0 || i >= len(tx.Inputs) {
					return false
				}

				if cov.Len() != 6 {
					return false
				}

				if len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) == 0 || len(cov.Items[2]) > MaxNameSize || len(cov.Items[3]) != 1 || len(cov.Items[4]) != 32 || len(cov.Items[5]) != 4 {
					return false
				}

				if !VerifyString(string(cov.Items[2])) {
					return false
				}

				var hash primitives.Hash
				copy(hash[:], cov.Items[0])

				if !HasReservedHash(hash) {
					return false
				}

				if !validNameHashAndLabel(cov, 0, 2) {
					return false
				}

			default:
				return false
			}
		}

		return true
	}

	for i, output := range tx.Outputs {
		cov := output.Covenant

		switch CovenantType(cov.Type) {
		case TypeNone:
			if cov.Len() != 0 {
				return false
			}

		case TypeClaim:
			return false

		case TypeOpen:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 3 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) == 0 || len(cov.Items[2]) > MaxNameSize {
				return false
			}

			if !VerifyString(string(cov.Items[2])) || !validNameHashAndLabel(cov, 0, 2) {
				return false
			}

			if h, err := cov.GetU32(1); err != nil || h != 0 {
				return false
			}

		case TypeBid:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 4 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) == 0 || len(cov.Items[2]) > MaxNameSize || len(cov.Items[3]) != 32 {
				return false
			}

			if !VerifyString(string(cov.Items[2])) || !validNameHashAndLabel(cov, 0, 2) {
				return false
			}

		case TypeReveal:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 3 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) != 32 {
				return false
			}

		case TypeRedeem:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 2 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 {
				return false
			}

		case TypeRegister:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 4 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) > MaxResourceSize || len(cov.Items[3]) != 32 {
				return false
			}

		case TypeUpdate:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 3 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) > MaxResourceSize {
				return false
			}

		case TypeRenew:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 3 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) != 32 {
				return false
			}

		case TypeTransfer:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 4 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) != 1 {
				return false
			}

			version := cov.Items[2][0]
			hash := cov.Items[3]

			if version > 31 || len(hash) < 2 || len(hash) > 40 {
				return false
			}

		case TypeFinalize:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 7 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 || len(cov.Items[2]) == 0 || len(cov.Items[2]) > MaxNameSize || len(cov.Items[3]) != 1 || len(cov.Items[4]) != 4 || len(cov.Items[5]) != 4 || len(cov.Items[6]) != 32 {
				return false
			}

			if !VerifyString(string(cov.Items[2])) || !validNameHashAndLabel(cov, 0, 2) {
				return false
			}

		case TypeRevoke:
			if i >= len(tx.Inputs) {
				return false
			}

			if cov.Len() != 2 || len(cov.Items[0]) != 32 || len(cov.Items[1]) != 4 {
				return false
			}

		default:
			if cov.Len() > 255 || cov.GetSize() > MaxCovenantSize {
				return false
			}
		}
	}

	return true
}
