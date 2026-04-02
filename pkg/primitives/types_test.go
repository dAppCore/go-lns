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
		isNone         bool
		isClaim        bool
		isOpen         bool
		isBid          bool
		isReveal       bool
		isRedeem       bool
		isRegister     bool
		isUpdate       bool
		isRenew        bool
		isTransfer     bool
		isFinalize     bool
		isRevoke       bool
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
			isNone:         true,
			isKnown:        true,
			isUnknown:      false,
			isName:         false,
			isDustworthy:   true,
			isNonspendable: false,
			isLinked:       false,
		},
		{
			name:           "claim",
			covenant:       Covenant{Type: covenantTypeClaim},
			isClaim:        true,
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       false,
		},
		{
			name:           "open",
			covenant:       Covenant{Type: covenantTypeOpen},
			isOpen:         true,
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
			isBid:          true,
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
			isReveal:       true,
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       true,
		},
		{
			name:           "redeem",
			covenant:       Covenant{Type: covenantTypeRedeem},
			isRedeem:       true,
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: false,
			isLinked:       true,
		},
		{
			name:           "register",
			covenant:       Covenant{Type: covenantTypeRegister},
			isRegister:     true,
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       true,
		},
		{
			name:           "update",
			covenant:       Covenant{Type: covenantTypeUpdate},
			isUpdate:       true,
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       true,
		},
		{
			name:           "renew",
			covenant:       Covenant{Type: covenantTypeRenew},
			isRenew:        true,
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       true,
		},
		{
			name:           "transfer",
			covenant:       Covenant{Type: covenantTypeTransfer},
			isTransfer:     true,
			isKnown:        true,
			isUnknown:      false,
			isName:         true,
			isDustworthy:   false,
			isNonspendable: true,
			isLinked:       true,
		},
		{
			name:           "finalize",
			covenant:       Covenant{Type: covenantTypeFinalize},
			isFinalize:     true,
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
			isRevoke:       true,
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

		if got := tc.covenant.IsNone(); got != tc.isNone {
			t.Fatalf("%s: IsNone() = %v, want %v", tc.name, got, tc.isNone)
		}

		if got := tc.covenant.IsClaim(); got != tc.isClaim {
			t.Fatalf("%s: IsClaim() = %v, want %v", tc.name, got, tc.isClaim)
		}

		if got := tc.covenant.IsOpen(); got != tc.isOpen {
			t.Fatalf("%s: IsOpen() = %v, want %v", tc.name, got, tc.isOpen)
		}

		if got := tc.covenant.IsBid(); got != tc.isBid {
			t.Fatalf("%s: IsBid() = %v, want %v", tc.name, got, tc.isBid)
		}

		if got := tc.covenant.IsReveal(); got != tc.isReveal {
			t.Fatalf("%s: IsReveal() = %v, want %v", tc.name, got, tc.isReveal)
		}

		if got := tc.covenant.IsRedeem(); got != tc.isRedeem {
			t.Fatalf("%s: IsRedeem() = %v, want %v", tc.name, got, tc.isRedeem)
		}

		if got := tc.covenant.IsRegister(); got != tc.isRegister {
			t.Fatalf("%s: IsRegister() = %v, want %v", tc.name, got, tc.isRegister)
		}

		if got := tc.covenant.IsUpdate(); got != tc.isUpdate {
			t.Fatalf("%s: IsUpdate() = %v, want %v", tc.name, got, tc.isUpdate)
		}

		if got := tc.covenant.IsRenew(); got != tc.isRenew {
			t.Fatalf("%s: IsRenew() = %v, want %v", tc.name, got, tc.isRenew)
		}

		if got := tc.covenant.IsTransfer(); got != tc.isTransfer {
			t.Fatalf("%s: IsTransfer() = %v, want %v", tc.name, got, tc.isTransfer)
		}

		if got := tc.covenant.IsFinalize(); got != tc.isFinalize {
			t.Fatalf("%s: IsFinalize() = %v, want %v", tc.name, got, tc.isFinalize)
		}

		if got := tc.covenant.IsRevoke(); got != tc.isRevoke {
			t.Fatalf("%s: IsRevoke() = %v, want %v", tc.name, got, tc.isRevoke)
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
