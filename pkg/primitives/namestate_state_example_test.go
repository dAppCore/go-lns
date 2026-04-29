// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNameStateStatus_String() {
	fn := ax7ValueOf(new(NameStateStatus).String)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_State() {
	fn := ax7ValueOf(new(NameState).State)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetState() {
	fn := ax7ValueOf(new(NameState).GetState)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsOpening() {
	fn := ax7ValueOf(new(NameState).IsOpening)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsOpening() {
	fn := ax7ValueOf(new(NameState).GetIsOpening)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsLocked() {
	fn := ax7ValueOf(new(NameState).IsLocked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsLocked() {
	fn := ax7ValueOf(new(NameState).GetIsLocked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsBidding() {
	fn := ax7ValueOf(new(NameState).IsBidding)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsBidding() {
	fn := ax7ValueOf(new(NameState).GetIsBidding)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsReveal() {
	fn := ax7ValueOf(new(NameState).IsReveal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsReveal() {
	fn := ax7ValueOf(new(NameState).GetIsReveal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsClosed() {
	fn := ax7ValueOf(new(NameState).IsClosed)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsClosed() {
	fn := ax7ValueOf(new(NameState).GetIsClosed)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsRevoked() {
	fn := ax7ValueOf(new(NameState).IsRevoked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsRevoked() {
	fn := ax7ValueOf(new(NameState).GetIsRevoked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsRedeemable() {
	fn := ax7ValueOf(new(NameState).IsRedeemable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsRedeemable() {
	fn := ax7ValueOf(new(NameState).GetIsRedeemable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsClaimable() {
	fn := ax7ValueOf(new(NameState).IsClaimable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsClaimable() {
	fn := ax7ValueOf(new(NameState).GetIsClaimable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_IsExpired() {
	fn := ax7ValueOf(new(NameState).IsExpired)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetIsExpired() {
	fn := ax7ValueOf(new(NameState).GetIsExpired)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
