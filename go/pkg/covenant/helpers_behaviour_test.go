// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"testing"

	core "github.com/dappcore/go"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

func TestHelpersBehaviour_compareBytes(t *testing.T) {
	core.AssertEqual(t, 0, compareBytes([]byte{1, 2}, []byte{1, 2}))
	core.AssertEqual(t, -1, compareBytes([]byte{1, 1}, []byte{1, 2}))
	core.AssertEqual(t, 1, compareBytes([]byte{1, 3}, []byte{1, 2}))
	// Shorter prefix sorts before its extension.
	core.AssertEqual(t, -1, compareBytes([]byte{1}, []byte{1, 2}))
	core.AssertEqual(t, 1, compareBytes([]byte{1, 2}, []byte{1}))

	core.AssertTrue(t, equalBytes([]byte{9}, []byte{9}))
	core.AssertFalse(t, equalBytes([]byte{9}, []byte{8}))
}

func TestHelpersBehaviour_jsonraw(t *testing.T) {
	var nilRaw jsonraw
	encoded, err := nilRaw.MarshalJSON()
	core.AssertNoError(t, err)
	core.AssertEqual(t, "null", string(encoded))

	populated := jsonraw(`{"k":1}`)
	encoded, err = populated.MarshalJSON()
	core.AssertNoError(t, err)
	core.AssertEqual(t, `{"k":1}`, string(encoded))

	var target jsonraw
	core.AssertNoError(t, target.UnmarshalJSON([]byte(`[1,2]`)))
	core.AssertEqual(t, `[1,2]`, string(target))

	var nilPtr *jsonraw
	core.AssertError(t, nilPtr.UnmarshalJSON([]byte("1")))
}

func TestHelpersBehaviour_covenantHashEquals(t *testing.T) {
	var hash primitives.Hash
	hash[0] = 0x42

	a := primitives.Covenant{}
	a.PushHash(hash)
	b := primitives.Covenant{}
	b.PushHash(hash)
	core.AssertTrue(t, covenantHashEquals(a, 0, b))

	var other primitives.Hash
	other[0] = 0x99
	c := primitives.Covenant{}
	c.PushHash(other)
	core.AssertFalse(t, covenantHashEquals(a, 0, c))

	// A missing/short item on either side yields false.
	empty := primitives.Covenant{}
	core.AssertFalse(t, covenantHashEquals(a, 0, empty))
	core.AssertFalse(t, covenantHashEquals(empty, 0, a))
}

func TestHelpersBehaviour_covenantU32Equals(t *testing.T) {
	a := primitives.Covenant{}
	a.PushU32(7)
	b := primitives.Covenant{}
	b.PushU32(7)
	core.AssertTrue(t, covenantU32Equals(a, 0, b, 0))

	c := primitives.Covenant{}
	c.PushU32(8)
	core.AssertFalse(t, covenantU32Equals(a, 0, c, 0))

	// A missing item on either side yields false.
	empty := primitives.Covenant{}
	core.AssertFalse(t, covenantU32Equals(a, 0, empty, 0))
	core.AssertFalse(t, covenantU32Equals(empty, 0, a, 0))
}

func TestHelpersBehaviour_reservedCatalog_nilGuards(t *testing.T) {
	var nilCatalog *ReservedCatalog
	core.AssertEqual(t, 0, nilCatalog.Size())
	core.AssertEqual(t, uint64(0), nilCatalog.NameValue())
	core.AssertEqual(t, uint64(0), nilCatalog.RootValue())
	core.AssertEqual(t, uint64(0), nilCatalog.TopValue())

	// The default catalog is populated and exposes positive base values.
	catalog := DefaultReservedCatalog()
	core.AssertGreaterOrEqual(t, catalog.Size(), 0)
	core.AssertEqual(t, catalog.Size(), catalog.GetSize())
	core.AssertEqual(t, catalog.NameValue(), catalog.GetNameValue())
}

func TestHelpersBehaviour_lockedCatalog_nilGuards(t *testing.T) {
	var nilCatalog *LockedCatalog
	core.AssertEqual(t, 0, nilCatalog.Size())
	core.AssertFalse(t, nilCatalog.Has(primitives.Hash{}))

	catalog := DefaultLockedCatalog()
	core.AssertGreaterOrEqual(t, catalog.Size(), 0)
	core.AssertEqual(t, catalog.Size(), catalog.GetSize())
}
