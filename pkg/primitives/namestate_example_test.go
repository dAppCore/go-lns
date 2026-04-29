// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNameState_Delta() {
	fn := ax7ValueOf(new(NameState).Delta)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetDelta() {
	fn := ax7ValueOf(new(NameState).GetDelta)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_Inject() {
	fn := ax7ValueOf(new(NameState).Inject)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetInject() {
	fn := ax7ValueOf(new(NameState).GetInject)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_Clone() {
	fn := ax7ValueOf(new(NameState).Clone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetClone() {
	fn := ax7ValueOf(new(NameState).GetClone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_Clear() {
	fn := ax7ValueOf(new(NameState).Clear)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetClear() {
	fn := ax7ValueOf(new(NameState).GetClear)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_HasDelta() {
	fn := ax7ValueOf(new(NameState).HasDelta)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetHasDelta() {
	fn := ax7ValueOf(new(NameState).GetHasDelta)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_Reset() {
	fn := ax7ValueOf(new(NameState).Reset)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetReset() {
	fn := ax7ValueOf(new(NameState).GetReset)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_Set() {
	fn := ax7ValueOf(new(NameState).Set)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetSet() {
	fn := ax7ValueOf(new(NameState).GetSet)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_MaybeExpire() {
	fn := ax7ValueOf(new(NameState).MaybeExpire)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetMaybeExpire() {
	fn := ax7ValueOf(new(NameState).GetMaybeExpire)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_ApplyState() {
	fn := ax7ValueOf(new(NameState).ApplyState)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetApplyState() {
	fn := ax7ValueOf(new(NameState).GetApplyState)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsNull() {
	fn := ax7ValueOf(new(NameState).IsNull)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsNull() {
	fn := ax7ValueOf(new(NameState).GetIsNull)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
