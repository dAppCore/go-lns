// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleGetVerificationFlags() {
	core.Println(coverageExample(GetVerificationFlags()))
	// Output: 1
}

func ExampleGetVerificationFlagsByVal() {
	core.Println(coverageExample(GetVerificationFlagsByVal()))
	// Output: 1
}

func ExampleGetRollout() {
	core.Println(coverageExample(GetRollout([32]byte{1}, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleHasRollout() {
	core.Println(coverageExample(HasRollout([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleGetHasRollout() {
	core.Println(coverageExample(GetHasRollout([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleIsReserved() {
	core.Println(coverageExample(IsReserved([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleGetIsReserved() {
	core.Println(coverageExample(GetIsReserved([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleIsLockedUp() {
	core.Println(coverageExample(IsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleGetIsLockedUp() {
	core.Println(coverageExample(GetIsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}
