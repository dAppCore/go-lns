// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNewOutpoint() {
	core.Println(coverageExample(NewOutpoint()))
	// Output: 1
}

func ExampleOutpoint_Inject() {
	core.Println(coverageExample(new(Outpoint).Inject(coverageValue[Outpoint]("good"))))
	// Output: 1
}

func ExampleOutpoint_Clone() {
	core.Println(coverageExample(new(Outpoint).Clone()))
	// Output: 1
}

func ExampleOutpoint_Equals() {
	core.Println(coverageExample(new(Outpoint).Equals(coverageValue[Outpoint]("good"))))
	// Output: 1
}

func ExampleOutpoint_Compare() {
	core.Println(coverageExample(new(Outpoint).Compare(coverageValue[Outpoint]("good"))))
	// Output: 1
}

func ExampleOutpoint_GetSize() {
	core.Println(coverageExample(new(Outpoint).GetSize()))
	// Output: 1
}

func ExampleOutpoint_MarshalBinary() {
	core.Println(coverageExample(new(Outpoint).MarshalBinary()))
	// Output: 2
}

func ExampleOutpoint_UnmarshalBinary() {
	core.Println(coverageExample(new(Outpoint).UnmarshalBinary([]byte("lethean.lthn"))))
	// Output: 1
}
