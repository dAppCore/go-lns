// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleGetDummy() {
	core.Println(coverageExample(GetDummy()))
	// Output: 1
}

func ExampleGetDefaultTTL() {
	core.Println(coverageExample(GetDefaultTTL()))
	// Output: 1
}

func ExampleGetTypeMapRoot() {
	core.Println(coverageExample(GetTypeMapRoot()))
	// Output: 1
}

func ExampleGetTypeMapEmpty() {
	core.Println(coverageExample(GetTypeMapEmpty()))
	// Output: 1
}

func ExampleGetTypeMapNS() {
	core.Println(coverageExample(GetTypeMapNS()))
	// Output: 1
}

func ExampleGetTypeMapTXT() {
	core.Println(coverageExample(GetTypeMapTXT()))
	// Output: 1
}

func ExampleGetTypeMapA() {
	core.Println(coverageExample(GetTypeMapA()))
	// Output: 1
}

func ExampleGetTypeMapAAAA() {
	core.Println(coverageExample(GetTypeMapAAAA()))
	// Output: 1
}

func ExampleGetHSTypes() {
	core.Println(coverageExample(GetHSTypes()))
	// Output: 1
}

func ExampleGetHSTypesByVal() {
	core.Println(coverageExample(GetHSTypesByVal()))
	// Output: 1
}
