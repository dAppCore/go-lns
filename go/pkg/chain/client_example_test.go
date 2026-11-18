// SPDX-License-Identifier: EUPL-1.2

package chain

import core "github.com/dappcore/go"

func ExampleWithAPIKey() {
	core.Println(coverageExample(WithAPIKey("lethean.lthn")))
	// Output: 1
}

func ExampleWithHTTPClient() {
	core.Println(coverageExample(WithHTTPClient(nil)))
	// Output: 1
}

func ExampleNewClient() {
	core.Println(coverageExample(NewClient("lethean.lthn")))
	// Output: 1
}

func Examplerpcfailure_Error() {
	core.Println(coverageExample(new(rpcfailure).Error()))
	// Output: 1
}

func Examplejsonraw_MarshalJSON() {
	core.Println(coverageExample(new(jsonraw).MarshalJSON()))
	// Output: 2
}

func Examplejsonraw_UnmarshalJSON() {
	core.Println(coverageExample(new(jsonraw).UnmarshalJSON([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleClient_RawCall() {
	core.Println(coverageExample(new(Client).RawCall("lethean.lthn", "lethean.lthn")))
	// Output: 2
}

func ExampleClient_GetBlockchainInfo() {
	core.Println(coverageExample(new(Client).GetBlockchainInfo()))
	// Output: 2
}

func ExampleClient_GetInfo() {
	core.Println(coverageExample(new(Client).GetInfo()))
	// Output: 2
}

func ExampleClient_GetNameResource() {
	core.Println(coverageExample(new(Client).GetNameResource("lethean.lthn")))
	// Output: 2
}

func ExampleClient_GetBlock() {
	core.Println(coverageExample(new(Client).GetBlock("lethean.lthn")))
	// Output: 2
}

func ExampleClient_GetBlockByHeight() {
	core.Println(coverageExample(new(Client).GetBlockByHeight(1)))
	// Output: 2
}

func ExampleClient_GetPeerInfo() {
	core.Println(coverageExample(new(Client).GetPeerInfo()))
	// Output: 2
}

func ExampleClient_SendRawTransaction() {
	core.Println(coverageExample(new(Client).SendRawTransaction("00aa")))
	// Output: 2
}

func ExampleClient_GetCoinsByAddress() {
	core.Println(coverageExample(new(Client).GetCoinsByAddress("ts1qexample")))
	// Output: 2
}

func ExampleClient_GetNameOwner() {
	core.Println(coverageExample(new(Client).GetNameOwner("lethean")))
	// Output: 2
}
