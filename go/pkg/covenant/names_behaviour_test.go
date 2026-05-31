// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

// nameTx returns a transaction with a single TypeTransfer name covenant
// whose first item is the given 32-byte hash.
func nameTx(hash primitives.Hash) primitives.Transaction {
	return primitives.Transaction{
		Outputs: []primitives.Output{{
			Covenant: primitives.Covenant{
				Type:  uint8(TypeTransfer),
				Items: [][]byte{append([]byte(nil), hash[:]...)},
			},
		}},
	}
}

func TestNamesBehaviour_AddHasRemove_Good(t *testing.T) {
	var hash primitives.Hash
	hash[0] = 0x42
	tx := nameTx(hash)
	set := map[primitives.Hash]struct{}{}

	core.AssertFalse(t, HasNames(tx, set))
	AddNames(tx, set)
	core.AssertTrue(t, HasNames(tx, set))
	core.AssertTrue(t, GetHasNames(tx, set))

	RemoveNames(tx, set)
	core.AssertFalse(t, HasNames(tx, set))
	core.AssertLen(t, set, 0)
}

func TestNamesBehaviour_AddHasRemove_Bad(t *testing.T) {
	var hash primitives.Hash
	hash[0] = 0x01
	tx := nameTx(hash)

	// All name set operations no-op on a nil set without panicking.
	core.AssertFalse(t, HasNames(tx, nil))
	core.AssertNotPanics(t, func() { AddNames(tx, nil) })
	core.AssertNotPanics(t, func() { RemoveNames(tx, nil) })
	core.AssertNotPanics(t, func() { GetAddNames(tx, nil) })
	core.AssertNotPanics(t, func() { GetRemoveNames(tx, nil) })
}

func TestNamesBehaviour_AddHasRemove_Ugly(t *testing.T) {
	// A TypeNone output is not a name covenant and is ignored.
	plain := primitives.Transaction{
		Outputs: []primitives.Output{{Covenant: primitives.Covenant{Type: uint8(TypeNone)}}},
	}
	set := map[primitives.Hash]struct{}{}
	AddNames(plain, set)
	core.AssertLen(t, set, 0)

	// A TypeBid name covenant is not in the tracked transition set, so it is
	// skipped by Add/Has/Remove even though IsName() is true.
	var hash primitives.Hash
	hash[0] = 0x99
	bid := primitives.Transaction{
		Outputs: []primitives.Output{{
			Covenant: primitives.Covenant{
				Type:  uint8(TypeBid),
				Items: [][]byte{append([]byte(nil), hash[:]...)},
			},
		}},
	}
	AddNames(bid, set)
	core.AssertLen(t, set, 0)
	core.AssertFalse(t, HasNames(bid, set))
}

func TestNamesBehaviour_nameSetHelpers(t *testing.T) {
	var hash primitives.Hash
	hash[0] = 0x07

	// nil-set guards return safe zero behaviour.
	core.AssertFalse(t, nameSetHas(nil, hash))
	core.AssertNotPanics(t, func() { nameSetAdd(nil, hash) })
	core.AssertNotPanics(t, func() { nameSetDelete(nil, hash) })

	set := map[primitives.Hash]struct{}{}
	core.AssertFalse(t, nameSetHas(set, hash))
	nameSetAdd(set, hash)
	core.AssertTrue(t, nameSetHas(set, hash))
	nameSetDelete(set, hash)
	core.AssertFalse(t, nameSetHas(set, hash))
}
