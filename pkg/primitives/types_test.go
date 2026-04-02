// SPDX-License-Identifier: EUPL-1.2

package primitives

import "testing"

func TestOutpointIsNull(t *testing.T) {
	zero := Outpoint{Index: ^uint32(0)}
	if !zero.IsNull() {
		t.Fatal("null outpoint should be recognised")
	}

	nonZero := Outpoint{Index: 0}
	if nonZero.IsNull() {
		t.Fatal("non-zero outpoint should not be null")
	}
}

func TestInputIsCoinbase(t *testing.T) {
	coinbase := Input{Prevout: Outpoint{Index: ^uint32(0)}}
	if !coinbase.IsCoinbase() {
		t.Fatal("input with a null prevout should be coinbase")
	}

	regular := Input{Prevout: Outpoint{Index: 0}}
	if regular.IsCoinbase() {
		t.Fatal("input with a real prevout should not be coinbase")
	}
}

func TestInputIsFinal(t *testing.T) {
	final := Input{Sequence: ^uint32(0)}
	if !final.IsFinal() {
		t.Fatal("input with max sequence should be final")
	}

	nonFinal := Input{Sequence: 0}
	if nonFinal.IsFinal() {
		t.Fatal("input with a lower sequence should not be final")
	}
}

func TestCovenantIsUnspendable(t *testing.T) {
	if got := (Covenant{}).IsUnspendable(); got {
		t.Fatal("zero covenant should be spendable")
	}

	if got := (Covenant{Type: covenantTypeRevoke}).IsUnspendable(); !got {
		t.Fatal("revoke covenant should be unspendable")
	}
}

func TestCovenantPredicates(t *testing.T) {
	cases := []struct {
		name           string
		covenant       Covenant
		isKnown        bool
		isUnknown      bool
		isName         bool
		isDustworthy   bool
		isNonspendable bool
		isLinked       bool
	}{
		{
			name:           "none",
			covenant:       Covenant{Type: covenantTypeNone},
			isKnown:        true,
			isUnknown:      false,
			isName:         false,
			isDustworthy:   true,
			isNonspendable: false,
			isLinked:       false,
		},
		{
			name:           "open",
			covenant:       Covenant{Type: covenantTypeOpen},
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: false,
			isLinked:       false,
		},
		{
			name:           "bid",
			covenant:       Covenant{Type: covenantTypeBid},
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   true,
			isNonspendable: true,
			isLinked:       false,
		},
		{
			name:           "reveal",
			covenant:       Covenant{Type: covenantTypeReveal},
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       true,
		},
		{
			name:           "revoke",
			covenant:       Covenant{Type: covenantTypeRevoke},
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       true,
		},
		{
			name:           "unknown",
			covenant:       Covenant{Type: 99},
			isKnown:        false,
			isUnknown:      true,
			isName:         false,
			isDustworthy:   true,
			isNonspendable: true,
			isLinked:       false,
		},
	}

	for _, tc := range cases {
		if got := tc.covenant.IsKnown(); got != tc.isKnown {
			t.Fatalf("%s: IsKnown() = %v, want %v", tc.name, got, tc.isKnown)
		}

		if got := tc.covenant.IsUnknown(); got != tc.isUnknown {
			t.Fatalf("%s: IsUnknown() = %v, want %v", tc.name, got, tc.isUnknown)
		}

		if got := tc.covenant.IsName(); got != tc.isName {
			t.Fatalf("%s: IsName() = %v, want %v", tc.name, got, tc.isName)
		}

		if got := tc.covenant.IsDustworthy(); got != tc.isDustworthy {
			t.Fatalf("%s: IsDustworthy() = %v, want %v", tc.name, got, tc.isDustworthy)
		}

		if got := tc.covenant.IsNonspendable(); got != tc.isNonspendable {
			t.Fatalf("%s: IsNonspendable() = %v, want %v", tc.name, got, tc.isNonspendable)
		}

		if got := tc.covenant.IsLinked(); got != tc.isLinked {
			t.Fatalf("%s: IsLinked() = %v, want %v", tc.name, got, tc.isLinked)
		}
	}
}

func TestOutputIsUnspendable(t *testing.T) {
	if got := (Output{}).IsUnspendable(); got {
		t.Fatal("zero output should be spendable")
	}

	if got := (Output{Address: Address{Version: 31}}).IsUnspendable(); !got {
		t.Fatal("nulldata output should be unspendable")
	}

	if got := (Output{Covenant: Covenant{Type: covenantTypeRevoke}}).IsUnspendable(); !got {
		t.Fatal("revoke covenant output should be unspendable")
	}
}
