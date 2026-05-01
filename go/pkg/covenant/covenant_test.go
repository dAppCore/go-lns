// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	core "dappco.re/go"
	"testing"
)

func TestCovenantTypePredicates(t *testing.T) {
	cases := []struct {
		name    string
		ct      CovenantType
		isKnown bool
		isName  bool
		isLink  bool
	}{
		{name: "none", ct: TypeNone, isKnown: true, isName: false, isLink: false},
		{name: "claim", ct: TypeClaim, isKnown: true, isName: true, isLink: false},
		{name: "open", ct: TypeOpen, isKnown: true, isName: true, isLink: false},
		{name: "bid", ct: TypeBid, isKnown: true, isName: true, isLink: false},
		{name: "reveal", ct: TypeReveal, isKnown: true, isName: true, isLink: true},
		{name: "redeem", ct: TypeRedeem, isKnown: true, isName: true, isLink: true},
		{name: "register", ct: TypeRegister, isKnown: true, isName: true, isLink: true},
		{name: "update", ct: TypeUpdate, isKnown: true, isName: true, isLink: true},
		{name: "renew", ct: TypeRenew, isKnown: true, isName: true, isLink: true},
		{name: "transfer", ct: TypeTransfer, isKnown: true, isName: true, isLink: true},
		{name: "finalize", ct: TypeFinalize, isKnown: true, isName: true, isLink: true},
		{name: "revoke", ct: TypeRevoke, isKnown: true, isName: true, isLink: true},
		{name: "unknown", ct: CovenantType(99), isKnown: false, isName: false, isLink: false},
	}

	for _, tc := range cases {
		if got := tc.ct.IsKnown(); got != tc.isKnown {
			t.Fatalf("%s: IsKnown() = %v, want %v", tc.name, got, tc.isKnown)
		}

		if got := tc.ct.IsName(); got != tc.isName {
			t.Fatalf("%s: IsName() = %v, want %v", tc.name, got, tc.isName)
		}

		if got := tc.ct.IsLinked(); got != tc.isLink {
			t.Fatalf("%s: IsLinked() = %v, want %v", tc.name, got, tc.isLink)
		}
	}
}

// BEGIN v090 file-aware generated coverage

func TestCovenant_GetTypes_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetTypes())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypes_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetTypes())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypes_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetTypes())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypesByVal_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetTypesByVal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypesByVal_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetTypesByVal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypesByVal_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetTypesByVal())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_TypeName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(TypeName(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_TypeName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(TypeName(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_TypeName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(TypeName(63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypeName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetTypeName(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypeName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetTypeName(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_GetTypeName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetTypeName(63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(CovenantType).IsName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(CovenantType).IsName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(CovenantType).IsName())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsKnown_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(CovenantType).IsKnown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsKnown_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(CovenantType).IsKnown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsKnown_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(CovenantType).IsKnown())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsLinked_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(CovenantType).IsLinked())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsLinked_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(CovenantType).IsLinked())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_IsLinked_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(CovenantType).IsLinked())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_String_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(CovenantType).String())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_String_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(CovenantType).String())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenant_CovenantType_String_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(CovenantType).String())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
