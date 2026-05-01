// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNewClaim() {
	core.Println(coverageExample(NewClaim()))
	// Output: 1
}

func ExampleGetNewClaim() {
	core.Println(coverageExample(GetNewClaim()))
	// Output: 1
}

func ExampleClaim_Refresh() {
	core.Println(coverageExample(new(Claim).Refresh()))
	// Output: 1
}

func ExampleClaim_Clone() {
	core.Println(coverageExample(new(Claim).Clone()))
	// Output: 1
}

func ExampleClaim_Hash() {
	core.Println(coverageExample(new(Claim).Hash()))
	// Output: 1
}

func ExampleClaim_GetHash() {
	core.Println(coverageExample(new(Claim).GetHash()))
	// Output: 1
}

func ExampleClaim_HashHex() {
	core.Println(coverageExample(new(Claim).HashHex()))
	// Output: 1
}

func ExampleClaim_GetHashHex() {
	core.Println(coverageExample(new(Claim).GetHashHex()))
	// Output: 1
}

func ExampleClaim_GetSize() {
	core.Println(coverageExample(new(Claim).GetSize()))
	// Output: 1
}

func ExampleClaim_Weight() {
	core.Println(coverageExample(new(Claim).Weight()))
	// Output: 1
}

func ExampleClaim_GetWeight() {
	core.Println(coverageExample(new(Claim).GetWeight()))
	// Output: 1
}

func ExampleClaim_VirtualSize() {
	core.Println(coverageExample(new(Claim).VirtualSize()))
	// Output: 1
}

func ExampleClaim_GetVirtualSize() {
	core.Println(coverageExample(new(Claim).GetVirtualSize()))
	// Output: 1
}

func ExampleClaim_MarshalBinary() {
	core.Println(coverageExample(new(Claim).MarshalBinary()))
	// Output: 2
}

func ExampleClaim_UnmarshalBinary() {
	core.Println(coverageExample(new(Claim).UnmarshalBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleClaim_ToBlob() {
	core.Println(coverageExample(new(Claim).ToBlob()))
	// Output: 1
}

func ExampleClaim_GetToBlob() {
	core.Println(coverageExample(new(Claim).GetToBlob()))
	// Output: 1
}

func ExampleClaim_GetBlob() {
	core.Println(coverageExample(new(Claim).GetBlob()))
	// Output: 1
}

func ExampleClaim_FromBlob() {
	core.Println(coverageExample(new(Claim).FromBlob([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleClaim_GetFromBlob() {
	core.Println(coverageExample(new(Claim).GetFromBlob([]byte("lethean.lthn"))))
	// Output: 1
}
