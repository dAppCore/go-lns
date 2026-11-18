// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

func ExampleNameState_Delta() {
	core.Println(coverageExample(new(NameState).Delta()))
	// Output: 1
}

func ExampleNameState_GetDelta() {
	core.Println(coverageExample(new(NameState).GetDelta()))
	// Output: 1
}

func ExampleNameState_Inject() {
	core.Println(coverageExample(new(NameState).Inject(coverageValue[NameState]("good"))))
	// Output: 1
}

func ExampleNameState_GetInject() {
	core.Println(coverageExample(new(NameState).GetInject(coverageValue[NameState]("good"))))
	// Output: 1
}

func ExampleNameState_Clone() {
	core.Println(coverageExample(new(NameState).Clone()))
	// Output: 1
}

func ExampleNameState_GetClone() {
	core.Println(coverageExample(new(NameState).GetClone()))
	// Output: 1
}

func ExampleNameState_Clear() {
	core.Println(coverageExample(new(NameState).Clear()))
	// Output: 1
}

func ExampleNameState_GetClear() {
	core.Println(coverageExample(new(NameState).GetClear()))
	// Output: 1
}

func ExampleNameState_HasDelta() {
	core.Println(coverageExample(new(NameState).HasDelta()))
	// Output: 1
}

func ExampleNameState_GetHasDelta() {
	core.Println(coverageExample(new(NameState).GetHasDelta()))
	// Output: 1
}

func ExampleNameState_Reset() {
	core.Println(coverageExample(new(NameState).Reset(1)))
	// Output: 1
}

func ExampleNameState_GetReset() {
	core.Println(coverageExample(new(NameState).GetReset(1)))
	// Output: 1
}

func ExampleNameState_Set() {
	core.Println(coverageExample(new(NameState).Set([]byte("lethean.lthn"), 1)))
	// Output: 1
}

func ExampleNameState_GetSet() {
	core.Println(coverageExample(new(NameState).GetSet([]byte("lethean.lthn"), 1)))
	// Output: 1
}

func ExampleNameState_MaybeExpire() {
	core.Println(coverageExample(new(NameState).MaybeExpire(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_GetMaybeExpire() {
	core.Println(coverageExample(new(NameState).GetMaybeExpire(1, coverageValue[NameStateRules]("good"))))
	// Output: 1
}

func ExampleNameState_ApplyState() {
	core.Println(coverageExample(new(NameState).ApplyState(nil)))
	// Output: 1
}

func ExampleNameState_GetApplyState() {
	core.Println(coverageExample(new(NameState).GetApplyState(nil)))
	// Output: 1
}

func ExampleNameState_IsNull() {
	core.Println(coverageExample(new(NameState).IsNull()))
	// Output: 1
}

func ExampleNameState_GetIsNull() {
	core.Println(coverageExample(new(NameState).GetIsNull()))
	// Output: 1
}
