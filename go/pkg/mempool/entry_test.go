// SPDX-License-Identifier: EUPL-1.2

package mempool

import (
	"testing"
	"time"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// TestEntry_FeeRate_Good checks the fee rate is floor(fee*1000/size).
func TestEntry_FeeRate_Good(t *testing.T) {
	e := &Entry{Fee: 1000, Size: 250}
	if got := e.FeeRate(); got != 4000 {
		t.Fatalf("FeeRate = %d, want 4000", got) // 1000*1000/250
	}

	e2 := &Entry{Fee: 100, Size: 200}
	if got := e2.FeeRate(); got != 500 {
		t.Fatalf("FeeRate = %d, want 500", got)
	}
}

// TestEntry_FeeRate_Bad checks a higher-rate entry outranks a lower-rate one even
// when its absolute fee is smaller (rate, not raw fee, drives ordering).
func TestEntry_FeeRate_Bad(t *testing.T) {
	small := &Entry{Fee: 300, Size: 100}  // rate 3000
	bigAbs := &Entry{Fee: 500, Size: 500} // rate 1000, bigger absolute fee
	if small.FeeRate() <= bigAbs.FeeRate() {
		t.Fatal("smaller-but-denser fee should rate higher")
	}
}

// TestEntry_FeeRate_Ugly checks zero / negative size and nil receiver are safe.
func TestEntry_FeeRate_Ugly(t *testing.T) {
	if (&Entry{Fee: 1000, Size: 0}).FeeRate() != 0 {
		t.Fatal("zero-size entry should rate 0")
	}
	var nilEntry *Entry
	if nilEntry.FeeRate() != 0 {
		t.Fatal("nil entry should rate 0")
	}
}

// TestEntry_Age_Good checks age is now minus admission time.
func TestEntry_Age_Good(t *testing.T) {
	base := time.Unix(1000, 0)
	e := &Entry{Time: base}
	if got := e.Age(base.Add(90 * time.Minute)); got != 90*time.Minute {
		t.Fatalf("Age = %v, want 90m", got)
	}
}

// TestEntry_Age_Bad checks a clock before admission yields a negative age, not a panic.
func TestEntry_Age_Bad(t *testing.T) {
	base := time.Unix(2000, 0)
	e := &Entry{Time: base}
	if got := e.Age(base.Add(-time.Minute)); got >= 0 {
		t.Fatalf("Age = %v, want negative", got)
	}
}

// TestEntry_Age_Ugly checks the nil receiver returns zero.
func TestEntry_Age_Ugly(t *testing.T) {
	var nilEntry *Entry
	if nilEntry.Age(time.Now()) != 0 {
		t.Fatal("nil entry age should be 0")
	}
}

// TestEntry_computeFee_Good checks fee = inputs - outputs through a coin view.
func TestEntry_computeFee_Good(t *testing.T) {
	op := outpoint(0xD1, 0)
	set := NewSpentSet()
	set.Add(op, primitives.Output{Value: 5000})

	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: op}},
		Outputs: []primitives.Output{{Value: 4200}},
	}
	fee, err := computeFee(tx, set)
	if err != nil {
		t.Fatalf("computeFee: %v", err)
	}
	if fee != 800 {
		t.Fatalf("fee = %d, want 800", fee)
	}
}

// TestEntry_computeFee_Bad checks negative fees and missing inputs are rejected.
func TestEntry_computeFee_Bad(t *testing.T) {
	op := outpoint(0xD2, 0)
	set := NewSpentSet()
	set.Add(op, primitives.Output{Value: 1000})

	// Outputs exceed inputs.
	over := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: op}},
		Outputs: []primitives.Output{{Value: 2000}},
	}
	if _, err := computeFee(over, set); err == nil {
		t.Fatal("expected ErrNegativeFee for outputs > inputs")
	}

	// Input not in the view.
	missing := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: outpoint(0xD3, 0)}},
		Outputs: []primitives.Output{{Value: 1}},
	}
	if _, err := computeFee(missing, set); err == nil {
		t.Fatal("expected ErrMissingInputs for an unresolved input")
	}
}

// TestEntry_computeFee_Ugly checks a nil view is rejected, not dereferenced.
func TestEntry_computeFee_Ugly(t *testing.T) {
	tx := primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: outpoint(0xD4, 0)}},
		Outputs: []primitives.Output{{Value: 1}},
	}
	if _, err := computeFee(tx, nil); err == nil {
		t.Fatal("expected error for nil coin view")
	}
}
