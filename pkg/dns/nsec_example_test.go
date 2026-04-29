// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleCreate() {
	core.Println(coverageExample(Create("lethean.lthn", "lethean.lthn", []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetCreate() {
	core.Println(coverageExample(GetCreate("lethean.lthn", "lethean.lthn", []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleNextName() {
	core.Println(coverageExample(NextName("lethean.lthn")))
	// Output: 1
}

func ExampleGetNextName() {
	core.Println(coverageExample(GetNextName("lethean.lthn")))
	// Output: 1
}

func ExamplePrevName() {
	core.Println(coverageExample(PrevName("lethean.lthn")))
	// Output: 1
}

func ExampleGetPrevName() {
	core.Println(coverageExample(GetPrevName("lethean.lthn")))
	// Output: 1
}
