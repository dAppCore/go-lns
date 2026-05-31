// SPDX-License-Identifier: EUPL-1.2

package wallet

import (
	"testing"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

func TestWalletBehaviour_nilReceiver(t *testing.T) {
	var w *Wallet

	// Every accessor on a nil wallet returns the zero value without panicking.
	core.AssertEqual(t, "", w.ExportSeed())
	core.AssertEqual(t, "", w.SeedID())
	core.AssertEqual(t, Snapshot{}, w.Snapshot())
	core.AssertEqual(t, primitives.Address{}, w.ReceiveAddress())
	core.AssertEqual(t, primitives.Address{}, w.NextReceiveAddress())
	core.AssertNil(t, w.History())

	confirmed, locked := w.Balance()
	core.AssertEqual(t, uint64(0), confirmed)
	core.AssertEqual(t, uint64(0), locked)

	// Mutators on a nil wallet are no-ops.
	core.AssertNotPanics(t, func() { w.Credit(1, false, 0, primitives.Hash{}, "") })
	core.AssertNotPanics(t, func() { w.Debit(1, false, 0, primitives.Hash{}, "") })

	// PendingDebit returns the entry even for a nil wallet but stores nothing.
	entry := w.PendingDebit(5, primitives.Hash{}, "draft")
	core.AssertEqual(t, EntryDebit, entry.Kind)
	core.AssertTrue(t, entry.Pending)
}

func TestWalletBehaviour_Balance_lockedAndPending(t *testing.T) {
	w, err := Create()
	core.AssertNoError(t, err)

	w.Credit(100, false, 1, primitives.Hash{}, "income")
	w.Credit(50, true, 1, primitives.Hash{}, "locked income")
	w.Debit(30, false, 2, primitives.Hash{}, "spend")
	w.PendingDebit(999, primitives.Hash{}, "draft") // ignored by Balance

	confirmed, locked := w.Balance()
	core.AssertEqual(t, uint64(70), confirmed)
	core.AssertEqual(t, uint64(50), locked)

	// History is newest-first and includes the pending draft.
	history := w.History()
	core.AssertLen(t, history, 4)
	core.AssertTrue(t, history[0].Pending)
}

func TestWalletBehaviour_applySignedAmount_underflow(t *testing.T) {
	// Debiting more than the current balance floors at zero, never wraps.
	core.AssertEqual(t, uint64(0), applySignedAmount(10, -1, 20))
	core.AssertEqual(t, uint64(5), applySignedAmount(10, -1, 5))
	core.AssertEqual(t, uint64(30), applySignedAmount(10, 1, 20))
	// An exact-match debit also floors at zero.
	core.AssertEqual(t, uint64(0), applySignedAmount(10, -1, 10))
}

func TestWalletBehaviour_Balance_debitFloor(t *testing.T) {
	w, err := Create()
	core.AssertNoError(t, err)

	// A debit larger than credits must not underflow the confirmed total.
	w.Credit(10, false, 1, primitives.Hash{}, "")
	w.Debit(40, false, 2, primitives.Hash{}, "")
	confirmed, _ := w.Balance()
	core.AssertEqual(t, uint64(0), confirmed)
}
