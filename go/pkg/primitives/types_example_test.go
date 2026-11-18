// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

func ExampleAddress_Inject() {
	core.Println(coverageExample(new(Address).Inject(coverageValue[Address]("good"))))
	// Output: 1
}

func ExampleAddress_Clone() {
	core.Println(coverageExample(new(Address).Clone()))
	// Output: 1
}

func ExampleAddress_Equals() {
	core.Println(coverageExample(new(Address).Equals(coverageValue[Address]("good"))))
	// Output: 1
}

func ExampleAddress_Compare() {
	core.Println(coverageExample(new(Address).Compare(coverageValue[Address]("good"))))
	// Output: 1
}

func ExampleAddress_IsNull() {
	core.Println(coverageExample(new(Address).IsNull()))
	// Output: 1
}

func ExampleAddress_IsPubkeyHash() {
	core.Println(coverageExample(new(Address).IsPubkeyHash()))
	// Output: 1
}

func ExampleAddress_IsScripthash() {
	core.Println(coverageExample(new(Address).IsScripthash()))
	// Output: 1
}

func ExampleAddress_IsNulldata() {
	core.Println(coverageExample(new(Address).IsNulldata()))
	// Output: 1
}

func ExampleAddress_IsUnspendable() {
	core.Println(coverageExample(new(Address).IsUnspendable()))
	// Output: 1
}

func ExampleAddress_IsUnknown() {
	core.Println(coverageExample(new(Address).IsUnknown()))
	// Output: 1
}

func ExampleAddress_IsValid() {
	core.Println(coverageExample(new(Address).IsValid()))
	// Output: 1
}

func ExampleAddress_GetSize() {
	core.Println(coverageExample(new(Address).GetSize()))
	// Output: 1
}

func ExampleCovenant_IsUnspendable() {
	core.Println(coverageExample(new(Covenant).IsUnspendable()))
	// Output: 1
}

func ExampleCovenant_IsNone() {
	core.Println(coverageExample(new(Covenant).IsNone()))
	// Output: 1
}

func ExampleCovenant_IsClaim() {
	core.Println(coverageExample(new(Covenant).IsClaim()))
	// Output: 1
}

func ExampleCovenant_IsOpen() {
	core.Println(coverageExample(new(Covenant).IsOpen()))
	// Output: 1
}

func ExampleCovenant_IsBid() {
	core.Println(coverageExample(new(Covenant).IsBid()))
	// Output: 1
}

func ExampleCovenant_IsReveal() {
	core.Println(coverageExample(new(Covenant).IsReveal()))
	// Output: 1
}

func ExampleCovenant_IsRedeem() {
	core.Println(coverageExample(new(Covenant).IsRedeem()))
	// Output: 1
}

func ExampleCovenant_IsRegister() {
	core.Println(coverageExample(new(Covenant).IsRegister()))
	// Output: 1
}

func ExampleCovenant_IsUpdate() {
	core.Println(coverageExample(new(Covenant).IsUpdate()))
	// Output: 1
}

func ExampleCovenant_IsRenew() {
	core.Println(coverageExample(new(Covenant).IsRenew()))
	// Output: 1
}

func ExampleCovenant_IsTransfer() {
	core.Println(coverageExample(new(Covenant).IsTransfer()))
	// Output: 1
}

func ExampleCovenant_IsFinalize() {
	core.Println(coverageExample(new(Covenant).IsFinalize()))
	// Output: 1
}

func ExampleCovenant_IsRevoke() {
	core.Println(coverageExample(new(Covenant).IsRevoke()))
	// Output: 1
}

func ExampleCovenant_IsName() {
	core.Println(coverageExample(new(Covenant).IsName()))
	// Output: 1
}

func ExampleCovenant_IsKnown() {
	core.Println(coverageExample(new(Covenant).IsKnown()))
	// Output: 1
}

func ExampleCovenant_IsUnknown() {
	core.Println(coverageExample(new(Covenant).IsUnknown()))
	// Output: 1
}

func ExampleCovenant_IsDustworthy() {
	core.Println(coverageExample(new(Covenant).IsDustworthy()))
	// Output: 1
}

func ExampleCovenant_IsNonspendable() {
	core.Println(coverageExample(new(Covenant).IsNonspendable()))
	// Output: 1
}

func ExampleCovenant_IsLinked() {
	core.Println(coverageExample(new(Covenant).IsLinked()))
	// Output: 1
}

func ExampleOutpoint_IsNull() {
	core.Println(coverageExample(new(Outpoint).IsNull()))
	// Output: 1
}

func ExampleInput_IsCoinbase() {
	core.Println(coverageExample(new(Input).IsCoinbase()))
	// Output: 1
}

func ExampleInput_IsFinal() {
	core.Println(coverageExample(new(Input).IsFinal()))
	// Output: 1
}

func ExampleOutput_IsUnspendable() {
	core.Println(coverageExample(new(Output).IsUnspendable()))
	// Output: 1
}
