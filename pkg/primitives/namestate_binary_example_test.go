// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNameState_GetSize() {
	core.Println(coverageExample(new(NameState).GetSize()))
	// Output: 1
}

func ExampleNameState_MarshalBinary() {
	core.Println(coverageExample(new(NameState).MarshalBinary()))
	// Output: 2
}

func ExampleNameState_UnmarshalBinary() {
	core.Println(coverageExample(new(NameState).UnmarshalBinary([]byte("lethean.lthn"))))
	// Output: 1
}
