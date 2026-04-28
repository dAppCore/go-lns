// SPDX-License-Identifier: EUPL-1.2

package chain

import core "dappco.re/go"

func TestAX7_Error_Error_Good(t *core.T) {
	err := &rpcError{Code: -1, Message: "method failed"}
	got := err.Error()

	core.AssertContains(t, got, "-1")
	core.AssertContains(t, got, "method failed")
}

func TestAX7_Error_Error_Bad(t *core.T) {
	err := &rpcError{}
	got := err.Error()

	core.AssertContains(t, got, "0")
	core.AssertContains(t, got, "rpc error")
}

func TestAX7_Error_Error_Ugly(t *core.T) {
	err := &rpcError{Code: 42, Message: ""}
	got := err.Error()

	core.AssertContains(t, got, "42")
	core.AssertContains(t, got, "rpc error")
}
