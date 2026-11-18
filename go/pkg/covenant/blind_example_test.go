// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "github.com/dappcore/go"

func ExampleBlind() {
	core.Println(coverageExample(Blind(1, [32]byte{1})))
	// Output: 2
}

func ExampleGetBlind() {
	core.Println(coverageExample(GetBlind(1, [32]byte{1})))
	// Output: 2
}
