// SPDX-License-Identifier: EUPL-1.2

package covenant

import "testing"

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

