// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"testing"

	core "github.com/dappcore/go"
)

func TestCovenantJSONBehaviour_actionNames(t *testing.T) {
	cases := []struct {
		typ  uint8
		want string
	}{
		{covenantTypeNone, "NONE"},
		{covenantTypeClaim, "CLAIM"},
		{covenantTypeOpen, "OPEN"},
		{covenantTypeBid, "BID"},
		{covenantTypeReveal, "REVEAL"},
		{covenantTypeRedeem, "REDEEM"},
		{covenantTypeRegister, "REGISTER"},
		{covenantTypeUpdate, "UPDATE"},
		{covenantTypeRenew, "RENEW"},
		{covenantTypeTransfer, "TRANSFER"},
		{covenantTypeFinalize, "FINALIZE"},
		{covenantTypeRevoke, "REVOKE"},
		{200, "UNKNOWN"},
	}
	for _, tc := range cases {
		cov := Covenant{Type: tc.typ}
		core.AssertEqual(t, tc.want, cov.GetJSON().Action)
	}
}

func TestCovenantJSONBehaviour_itemsHex(t *testing.T) {
	cov := Covenant{Type: covenantTypeOpen, Items: [][]byte{{0xde, 0xad}}}
	jsonView := cov.GetJSON()
	core.AssertEqual(t, "OPEN", jsonView.Action)
	core.AssertLen(t, jsonView.Items, 1)
	core.AssertEqual(t, "dead", jsonView.Items[0])
}
