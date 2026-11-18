// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

func ExampleNameState_ToStats() {
	core.Println(coverageExample(new(NameState).ToStats(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetToStats() {
	core.Println(coverageExample(new(NameState).GetToStats(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}
