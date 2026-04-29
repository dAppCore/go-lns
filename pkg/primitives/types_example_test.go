// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleAddress_Inject() {
	fn := ax7ValueOf(new(Address).Inject)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_Clone() {
	fn := ax7ValueOf(new(Address).Clone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_Equals() {
	fn := ax7ValueOf(new(Address).Equals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_Compare() {
	fn := ax7ValueOf(new(Address).Compare)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_IsNull() {
	fn := ax7ValueOf(new(Address).IsNull)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_IsPubkeyHash() {
	fn := ax7ValueOf(new(Address).IsPubkeyHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_IsScripthash() {
	fn := ax7ValueOf(new(Address).IsScripthash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_IsNulldata() {
	fn := ax7ValueOf(new(Address).IsNulldata)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_IsUnspendable() {
	fn := ax7ValueOf(new(Address).IsUnspendable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_IsUnknown() {
	fn := ax7ValueOf(new(Address).IsUnknown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_IsValid() {
	fn := ax7ValueOf(new(Address).IsValid)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddress_GetSize() {
	fn := ax7ValueOf(new(Address).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsUnspendable() {
	fn := ax7ValueOf(new(Covenant).IsUnspendable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsNone() {
	fn := ax7ValueOf(new(Covenant).IsNone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsClaim() {
	fn := ax7ValueOf(new(Covenant).IsClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsOpen() {
	fn := ax7ValueOf(new(Covenant).IsOpen)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsBid() {
	fn := ax7ValueOf(new(Covenant).IsBid)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsReveal() {
	fn := ax7ValueOf(new(Covenant).IsReveal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsRedeem() {
	fn := ax7ValueOf(new(Covenant).IsRedeem)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsRegister() {
	fn := ax7ValueOf(new(Covenant).IsRegister)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsUpdate() {
	fn := ax7ValueOf(new(Covenant).IsUpdate)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsRenew() {
	fn := ax7ValueOf(new(Covenant).IsRenew)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsTransfer() {
	fn := ax7ValueOf(new(Covenant).IsTransfer)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsFinalize() {
	fn := ax7ValueOf(new(Covenant).IsFinalize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsRevoke() {
	fn := ax7ValueOf(new(Covenant).IsRevoke)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsName() {
	fn := ax7ValueOf(new(Covenant).IsName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsKnown() {
	fn := ax7ValueOf(new(Covenant).IsKnown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsUnknown() {
	fn := ax7ValueOf(new(Covenant).IsUnknown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsDustworthy() {
	fn := ax7ValueOf(new(Covenant).IsDustworthy)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsNonspendable() {
	fn := ax7ValueOf(new(Covenant).IsNonspendable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_IsLinked() {
	fn := ax7ValueOf(new(Covenant).IsLinked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_IsNull() {
	fn := ax7ValueOf(new(Outpoint).IsNull)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleInput_IsCoinbase() {
	fn := ax7ValueOf(new(Input).IsCoinbase)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleInput_IsFinal() {
	fn := ax7ValueOf(new(Input).IsFinal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutput_IsUnspendable() {
	fn := ax7ValueOf(new(Output).IsUnspendable)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
