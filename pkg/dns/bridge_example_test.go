// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleNewBridge() {
	core.Println(coverageExample(NewBridge(coverageValue[BridgeConfig]("good"))))
	// Output: 1
}

func ExampleBridge_Discover() {
	core.Println(coverageExample(new(Bridge).Discover(core.Background())))
	// Output: 2
}

func ExampleBridge_Lookup() {
	core.Println(coverageExample(new(Bridge).Lookup("lethean.lthn")))
	// Output: 2
}

func ExampleBridge_Records() {
	core.Println(coverageExample(new(Bridge).Records()))
	// Output: 1
}

func ExampleBridge_Reverse() {
	core.Println(coverageExample(new(Bridge).Reverse(coverageAddr("good"))))
	// Output: 1
}

func ExampleBridge_Health() {
	core.Println(coverageExample(new(Bridge).Health()))
	// Output: 1
}

func ExampleBridge_Serve() {
	core.Println(coverageExample(new(Bridge).Serve(core.Background())))
	// Output: 1
}

func ExampleBridge_Shutdown() {
	core.Println(coverageExample(new(Bridge).Shutdown()))
	// Output: 1
}

func ExampleNameRecord_FQDN() {
	core.Println(coverageExample(new(NameRecord).FQDN()))
	// Output: 1
}

func ExampleNameRecord_Addresses4() {
	core.Println(coverageExample(new(NameRecord).Addresses4()))
	// Output: 1
}

func ExampleNameRecord_Addresses6() {
	core.Println(coverageExample(new(NameRecord).Addresses6()))
	// Output: 1
}

func ExampleNameRecord_TXT() {
	core.Println(coverageExample(new(NameRecord).TXT()))
	// Output: 1
}

func ExampleNameRecord_NS() {
	core.Println(coverageExample(new(NameRecord).NS()))
	// Output: 1
}
