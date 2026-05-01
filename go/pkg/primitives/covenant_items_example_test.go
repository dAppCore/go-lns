// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleCovenant_Len() {
	core.Println(coverageExample(new(Covenant).Len()))
	// Output: 1
}

func ExampleCovenant_Clone() {
	core.Println(coverageExample(new(Covenant).Clone()))
	// Output: 1
}

func ExampleCovenant_Inject() {
	core.Println(coverageExample(new(Covenant).Inject(coverageValue[Covenant]("good"))))
	// Output: 1
}

func ExampleCovenant_SetNone() {
	core.Println(coverageExample(new(Covenant).SetNone()))
	// Output: 1
}

func ExampleCovenant_ToArray() {
	core.Println(coverageExample(new(Covenant).ToArray()))
	// Output: 1
}

func ExampleCovenant_FromArray() {
	core.Println(coverageExample(new(Covenant).FromArray([][]byte{[]byte("lethean")})))
	// Output: 1
}

func ExampleCovenant_IndexOf() {
	core.Println(coverageExample(new(Covenant).IndexOf([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleCovenant_Get() {
	core.Println(coverageExample(new(Covenant).Get(1)))
	// Output: 2
}

func ExampleCovenant_Set() {
	core.Println(coverageExample(new(Covenant).Set(1, []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleCovenant_Push() {
	core.Println(coverageExample(new(Covenant).Push([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleCovenant_SetOpen() {
	core.Println(coverageExample(new(Covenant).SetOpen([32]byte{1}, []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleCovenant_SetBid() {
	core.Println(coverageExample(new(Covenant).SetBid([32]byte{1}, 1, []byte("lethean.lthn"), [32]byte{1})))
	// Output: 1
}

func ExampleCovenant_SetReveal() {
	core.Println(coverageExample(new(Covenant).SetReveal([32]byte{1}, 1, [32]byte{1})))
	// Output: 1
}

func ExampleCovenant_SetRedeem() {
	core.Println(coverageExample(new(Covenant).SetRedeem([32]byte{1}, 1)))
	// Output: 1
}

func ExampleCovenant_SetClaim() {
	core.Println(coverageExample(new(Covenant).SetClaim([32]byte{1}, 1, []byte("lethean.lthn"), 1, [32]byte{1}, 1)))
	// Output: 1
}

func ExampleCovenant_SetRegister() {
	core.Println(coverageExample(new(Covenant).SetRegister([32]byte{1}, 1, []byte("lethean.lthn"), [32]byte{1})))
	// Output: 1
}

func ExampleCovenant_SetUpdate() {
	core.Println(coverageExample(new(Covenant).SetUpdate([32]byte{1}, 1, []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleCovenant_SetRenew() {
	core.Println(coverageExample(new(Covenant).SetRenew([32]byte{1}, 1, [32]byte{1})))
	// Output: 1
}

func ExampleCovenant_SetTransfer() {
	core.Println(coverageExample(new(Covenant).SetTransfer([32]byte{1}, 1, coverageValue[Address]("good"))))
	// Output: 1
}

func ExampleCovenant_SetFinalize() {
	core.Println(coverageExample(new(Covenant).SetFinalize([32]byte{1}, 1, []byte("lethean.lthn"), 1, 1, 1, [32]byte{1})))
	// Output: 1
}

func ExampleCovenant_SetRevoke() {
	core.Println(coverageExample(new(Covenant).SetRevoke([32]byte{1}, 1)))
	// Output: 1
}

func ExampleCovenant_GetU8() {
	core.Println(coverageExample(new(Covenant).GetU8(1)))
	// Output: 2
}

func ExampleCovenant_PushU8() {
	core.Println(coverageExample(new(Covenant).PushU8(1)))
	// Output: 1
}

func ExampleCovenant_GetU32() {
	core.Println(coverageExample(new(Covenant).GetU32(1)))
	// Output: 2
}

func ExampleCovenant_PushU32() {
	core.Println(coverageExample(new(Covenant).PushU32(1)))
	// Output: 1
}

func ExampleCovenant_GetHash() {
	core.Println(coverageExample(new(Covenant).GetHash(1)))
	// Output: 2
}

func ExampleCovenant_PushHash() {
	core.Println(coverageExample(new(Covenant).PushHash([32]byte{1})))
	// Output: 1
}

func ExampleCovenant_GetString() {
	core.Println(coverageExample(new(Covenant).GetString(1)))
	// Output: 2
}

func ExampleCovenant_PushString() {
	core.Println(coverageExample(new(Covenant).PushString("lethean.lthn")))
	// Output: 1
}

func ExampleCovenant_GetSize() {
	core.Println(coverageExample(new(Covenant).GetSize()))
	// Output: 1
}

func ExampleCovenant_GetVarSize() {
	core.Println(coverageExample(new(Covenant).GetVarSize()))
	// Output: 1
}
