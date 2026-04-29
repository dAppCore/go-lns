// SPDX-License-Identifier: EUPL-1.2

package chain

import core "dappco.re/go"

func ExampleWithAPIKey() {
	fn := ax7ValueOf(WithAPIKey)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWithHTTPClient() {
	fn := ax7ValueOf(WithHTTPClient)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewClient() {
	fn := ax7ValueOf(NewClient)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func Examplerpcfailure_Error() {
	fn := ax7ValueOf(new(rpcfailure).Error)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func Examplejsonraw_MarshalJSON() {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func Examplejsonraw_UnmarshalJSON() {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleClient_RawCall() {
	fn := ax7ValueOf(new(Client).RawCall)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleClient_GetBlockchainInfo() {
	fn := ax7ValueOf(new(Client).GetBlockchainInfo)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleClient_GetInfo() {
	fn := ax7ValueOf(new(Client).GetInfo)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleClient_GetNameResource() {
	fn := ax7ValueOf(new(Client).GetNameResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleClient_GetBlock() {
	fn := ax7ValueOf(new(Client).GetBlock)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleClient_GetBlockByHeight() {
	fn := ax7ValueOf(new(Client).GetBlockByHeight)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleClient_GetPeerInfo() {
	fn := ax7ValueOf(new(Client).GetPeerInfo)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}
