// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

func ExampleNameState_GetJSON() {
	core.Println(coverageExample(new(NameState).GetJSON(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_Format() {
	core.Println(coverageExample(new(NameState).Format(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetFormat() {
	core.Println(coverageExample(new(NameState).GetFormat(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_FromJSON() {
	core.Println(coverageExample(new(NameState).FromJSON(coverageValue[NameStateJSON]("good"))))
	// Output: 1
}

func ExampleNameState_GetFromJSON() {
	core.Println(coverageExample(new(NameState).GetFromJSON(coverageValue[NameStateJSON]("good"))))
	// Output: 1
}
