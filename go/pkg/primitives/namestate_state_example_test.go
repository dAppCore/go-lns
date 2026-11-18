// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

func ExampleNameStateStatus_String() {
	core.Println(coverageExample(new(NameStateStatus).String()))
	// Output: 1
}

func ExampleNameState_State() {
	core.Println(coverageExample(new(NameState).State(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetState() {
	core.Println(coverageExample(new(NameState).GetState(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsOpening() {
	core.Println(coverageExample(new(NameState).IsOpening(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsOpening() {
	core.Println(coverageExample(new(NameState).GetIsOpening(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsLocked() {
	core.Println(coverageExample(new(NameState).IsLocked(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsLocked() {
	core.Println(coverageExample(new(NameState).GetIsLocked(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsBidding() {
	core.Println(coverageExample(new(NameState).IsBidding(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsBidding() {
	core.Println(coverageExample(new(NameState).GetIsBidding(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsReveal() {
	core.Println(coverageExample(new(NameState).IsReveal(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsReveal() {
	core.Println(coverageExample(new(NameState).GetIsReveal(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsClosed() {
	core.Println(coverageExample(new(NameState).IsClosed(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsClosed() {
	core.Println(coverageExample(new(NameState).GetIsClosed(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsRevoked() {
	core.Println(coverageExample(new(NameState).IsRevoked(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsRevoked() {
	core.Println(coverageExample(new(NameState).GetIsRevoked(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsRedeemable() {
	core.Println(coverageExample(new(NameState).IsRedeemable(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsRedeemable() {
	core.Println(coverageExample(new(NameState).GetIsRedeemable(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsClaimable() {
	core.Println(coverageExample(new(NameState).IsClaimable(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsClaimable() {
	core.Println(coverageExample(new(NameState).GetIsClaimable(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_IsExpired() {
	core.Println(coverageExample(new(NameState).IsExpired(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetIsExpired() {
	core.Println(coverageExample(new(NameState).GetIsExpired(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}
