// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleAddress_Encode() {
	addr := Address{Version: 0, Hash: make([]byte, 20)}
	core.Println(coverageExample(addr.Encode("ts")))
	// Output: 2
}
