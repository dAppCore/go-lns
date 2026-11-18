// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

func ExampleNameDelta_IsNull() {
	core.Println(coverageExample(new(NameDelta).IsNull()))
	// Output: 1
}

func ExampleNameDelta_Clone() {
	core.Println(coverageExample(new(NameDelta).Clone()))
	// Output: 1
}

func ExampleNameDelta_Inject() {
	core.Println(coverageExample(new(NameDelta).Inject(coverageValue[NameDelta]("good"))))
	// Output: 1
}

func ExampleNameDelta_Clear() {
	core.Println(coverageExample(new(NameDelta).Clear()))
	// Output: 1
}

func ExampleNameDelta_GetField() {
	core.Println(coverageExample(new(NameDelta).GetField()))
	// Output: 1
}

func ExampleNameDelta_GetSize() {
	core.Println(coverageExample(new(NameDelta).GetSize()))
	// Output: 1
}

func ExampleNameDelta_MarshalBinary() {
	core.Println(coverageExample(new(NameDelta).MarshalBinary()))
	// Output: 2
}

func ExampleNameDelta_UnmarshalBinary() {
	core.Println(coverageExample(new(NameDelta).UnmarshalBinary([]byte("lethean.lthn"))))
	// Output: 1
}
