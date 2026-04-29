// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleCovenant_Len() {
	fn := ax7ValueOf(new(Covenant).Len)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_Clone() {
	fn := ax7ValueOf(new(Covenant).Clone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_Inject() {
	fn := ax7ValueOf(new(Covenant).Inject)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetNone() {
	fn := ax7ValueOf(new(Covenant).SetNone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_ToArray() {
	fn := ax7ValueOf(new(Covenant).ToArray)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_FromArray() {
	fn := ax7ValueOf(new(Covenant).FromArray)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IndexOf() {
	fn := ax7ValueOf(new(Covenant).IndexOf)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_Get() {
	fn := ax7ValueOf(new(Covenant).Get)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCovenant_Set() {
	fn := ax7ValueOf(new(Covenant).Set)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_Push() {
	fn := ax7ValueOf(new(Covenant).Push)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetOpen() {
	fn := ax7ValueOf(new(Covenant).SetOpen)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetBid() {
	fn := ax7ValueOf(new(Covenant).SetBid)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetReveal() {
	fn := ax7ValueOf(new(Covenant).SetReveal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetRedeem() {
	fn := ax7ValueOf(new(Covenant).SetRedeem)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetClaim() {
	fn := ax7ValueOf(new(Covenant).SetClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetRegister() {
	fn := ax7ValueOf(new(Covenant).SetRegister)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetUpdate() {
	fn := ax7ValueOf(new(Covenant).SetUpdate)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetRenew() {
	fn := ax7ValueOf(new(Covenant).SetRenew)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetTransfer() {
	fn := ax7ValueOf(new(Covenant).SetTransfer)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetFinalize() {
	fn := ax7ValueOf(new(Covenant).SetFinalize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_SetRevoke() {
	fn := ax7ValueOf(new(Covenant).SetRevoke)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_GetU8() {
	fn := ax7ValueOf(new(Covenant).GetU8)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCovenant_PushU8() {
	fn := ax7ValueOf(new(Covenant).PushU8)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_GetU32() {
	fn := ax7ValueOf(new(Covenant).GetU32)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCovenant_PushU32() {
	fn := ax7ValueOf(new(Covenant).PushU32)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_GetHash() {
	fn := ax7ValueOf(new(Covenant).GetHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCovenant_PushHash() {
	fn := ax7ValueOf(new(Covenant).PushHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_GetString() {
	fn := ax7ValueOf(new(Covenant).GetString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCovenant_PushString() {
	fn := ax7ValueOf(new(Covenant).PushString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_GetSize() {
	fn := ax7ValueOf(new(Covenant).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_GetVarSize() {
	fn := ax7ValueOf(new(Covenant).GetVarSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
