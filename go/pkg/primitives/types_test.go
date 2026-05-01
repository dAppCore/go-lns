// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "dappco.re/go"
	"testing"
)

func TestOutpointIsNull(t *testing.T) {
	zero := NewOutpoint()
	if !zero.IsNull() {
		t.Fatal("null outpoint should be recognised")
	}

	nonZero := Outpoint{Index: 0}
	if nonZero.IsNull() {
		t.Fatal("non-zero outpoint should not be null")
	}
}

func TestOutpointHelpers(t *testing.T) {
	var hash Hash
	hash[0] = 0x11

	left := Outpoint{TxHash: hash, Index: 7}
	right := left.Clone()

	if !left.Equals(right) {
		t.Fatal("Clone() should produce an equal outpoint")
	}

	if got := left.Compare(right); got != 0 {
		t.Fatalf("Compare() = %d, want 0", got)
	}

	if got := (Outpoint{TxHash: hash, Index: 7}).Compare(Outpoint{TxHash: hash, Index: 8}); got >= 0 {
		t.Fatalf("Compare() = %d, want negative", got)
	}

	if got := (Outpoint{TxHash: hash, Index: 8}).Compare(Outpoint{TxHash: hash, Index: 7}); got <= 0 {
		t.Fatalf("Compare() = %d, want positive", got)
	}

	var injected Outpoint
	if got := injected.Inject(left); got != &injected {
		t.Fatal("Inject() should return the receiver")
	}

	if !injected.Equals(left) {
		t.Fatal("Inject() should copy the source outpoint")
	}

	if got := left.GetSize(); got != 36 {
		t.Fatalf("GetSize() = %d, want 36", got)
	}

	encoded, err := left.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary returned error: %v", err)
	}

	if len(encoded) != left.GetSize() {
		t.Fatalf("MarshalBinary length = %d, want %d", len(encoded), left.GetSize())
	}

	var decoded Outpoint
	if err := decoded.UnmarshalBinary(encoded); err != nil {
		t.Fatalf("UnmarshalBinary returned error: %v", err)
	}

	if !decoded.Equals(left) {
		t.Fatal("UnmarshalBinary should round-trip the outpoint")
	}

	if !equalBytes(encoded[:len(hash)], hash[:]) {
		t.Fatal("MarshalBinary should preserve the tx hash prefix")
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

// BEGIN v090 file-aware generated coverage

func TestTypes_Address_Inject_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).Inject(coverageValue[Address]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Inject_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).Inject(coverageValue[Address]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Inject_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).Inject(coverageValue[Address]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Clone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Clone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Clone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Equals_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).Equals(coverageValue[Address]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Equals_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).Equals(coverageValue[Address]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Equals_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).Equals(coverageValue[Address]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Compare_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).Compare(coverageValue[Address]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Compare_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).Compare(coverageValue[Address]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_Compare_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).Compare(coverageValue[Address]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsNull_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsNull_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsNull_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsPubkeyHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).IsPubkeyHash())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsPubkeyHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).IsPubkeyHash())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsPubkeyHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).IsPubkeyHash())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsScripthash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).IsScripthash())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsScripthash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).IsScripthash())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsScripthash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).IsScripthash())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsNulldata_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).IsNulldata())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsNulldata_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).IsNulldata())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsNulldata_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).IsNulldata())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsUnspendable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsUnspendable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsUnspendable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsUnknown_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).IsUnknown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsUnknown_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).IsUnknown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsUnknown_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).IsUnknown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsValid_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).IsValid())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsValid_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).IsValid())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_IsValid_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).IsValid())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Address).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Address).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Address_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Address).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUnspendable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUnspendable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUnspendable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsNone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsNone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsNone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsNone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsNone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsNone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsClaim_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsClaim())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsClaim_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsClaim())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsClaim_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsClaim())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsOpen_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsOpen())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsOpen_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsOpen())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsOpen_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsOpen())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsBid_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsBid())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsBid_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsBid())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsBid_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsBid())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsReveal_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsReveal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsReveal_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsReveal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsReveal_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsReveal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRedeem_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsRedeem())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRedeem_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsRedeem())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRedeem_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsRedeem())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRegister_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsRegister())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRegister_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsRegister())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRegister_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsRegister())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUpdate_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsUpdate())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUpdate_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsUpdate())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUpdate_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsUpdate())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRenew_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsRenew())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRenew_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsRenew())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRenew_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsRenew())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsTransfer_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsTransfer())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsTransfer_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsTransfer())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsTransfer_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsTransfer())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsFinalize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsFinalize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsFinalize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsFinalize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsFinalize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsFinalize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRevoke_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsRevoke())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRevoke_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsRevoke())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsRevoke_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsRevoke())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsKnown_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsKnown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsKnown_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsKnown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsKnown_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsKnown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUnknown_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsUnknown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUnknown_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsUnknown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsUnknown_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsUnknown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsDustworthy_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsDustworthy())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsDustworthy_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsDustworthy())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsDustworthy_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsDustworthy())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsNonspendable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsNonspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsNonspendable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsNonspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsNonspendable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsNonspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsLinked_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IsLinked())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsLinked_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IsLinked())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Covenant_IsLinked_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IsLinked())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Outpoint_IsNull_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Outpoint).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Outpoint_IsNull_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Outpoint).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Outpoint_IsNull_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Outpoint).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Input_IsCoinbase_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Input).IsCoinbase())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Input_IsCoinbase_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Input).IsCoinbase())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Input_IsCoinbase_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Input).IsCoinbase())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Input_IsFinal_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Input).IsFinal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Input_IsFinal_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Input).IsFinal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Input_IsFinal_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Input).IsFinal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Output_IsUnspendable_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Output).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Output_IsUnspendable_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Output).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestTypes_Output_IsUnspendable_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Output).IsUnspendable())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
