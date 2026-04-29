// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNewClaim() {
	fn := ax7ValueOf(NewClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewClaim() {
	fn := ax7ValueOf(GetNewClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_Refresh() {
	fn := ax7ValueOf(new(Claim).Refresh)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_Clone() {
	fn := ax7ValueOf(new(Claim).Clone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_Hash() {
	fn := ax7ValueOf(new(Claim).Hash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetHash() {
	fn := ax7ValueOf(new(Claim).GetHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_HashHex() {
	fn := ax7ValueOf(new(Claim).HashHex)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetHashHex() {
	fn := ax7ValueOf(new(Claim).GetHashHex)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetSize() {
	fn := ax7ValueOf(new(Claim).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_Weight() {
	fn := ax7ValueOf(new(Claim).Weight)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetWeight() {
	fn := ax7ValueOf(new(Claim).GetWeight)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_VirtualSize() {
	fn := ax7ValueOf(new(Claim).VirtualSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetVirtualSize() {
	fn := ax7ValueOf(new(Claim).GetVirtualSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_MarshalBinary() {
	fn := ax7ValueOf(new(Claim).MarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleClaim_UnmarshalBinary() {
	fn := ax7ValueOf(new(Claim).UnmarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_ToBlob() {
	fn := ax7ValueOf(new(Claim).ToBlob)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetToBlob() {
	fn := ax7ValueOf(new(Claim).GetToBlob)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetBlob() {
	fn := ax7ValueOf(new(Claim).GetBlob)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_FromBlob() {
	fn := ax7ValueOf(new(Claim).FromBlob)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClaim_GetFromBlob() {
	fn := ax7ValueOf(new(Claim).GetFromBlob)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
