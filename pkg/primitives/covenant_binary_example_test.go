// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleCovenant_MarshalBinary() {
	core.Println(coverageExample(new(Covenant).MarshalBinary()))
	// Output: 2
}

func ExampleCovenant_UnmarshalBinary() {
	core.Println(coverageExample(new(Covenant).UnmarshalBinary([]byte("lethean.lthn"))))
	// Output: 1
}
