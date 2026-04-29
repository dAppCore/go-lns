// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleGetBlacklist() {
	core.Println(coverageExample(GetBlacklist()))
	// Output: 1
}

func ExampleVerifyString() {
	core.Println(coverageExample(VerifyString("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyString() {
	core.Println(coverageExample(GetVerifyString("lethean.lthn")))
	// Output: 1
}

func ExampleVerifyBinary() {
	core.Println(coverageExample(VerifyBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetVerifyBinary() {
	core.Println(coverageExample(GetVerifyBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleVerifyName() {
	core.Println(coverageExample(VerifyName("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyName() {
	core.Println(coverageExample(GetVerifyName("lethean.lthn")))
	// Output: 1
}

func ExampleVerifyByName() {
	core.Println(coverageExample(VerifyByName("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyByName() {
	core.Println(coverageExample(GetVerifyByName("lethean.lthn")))
	// Output: 1
}

func ExampleVerifyByString() {
	core.Println(coverageExample(VerifyByString("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyByString() {
	core.Println(coverageExample(GetVerifyByString("lethean.lthn")))
	// Output: 1
}

func ExampleVerifyByBinary() {
	core.Println(coverageExample(VerifyByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetVerifyByBinary() {
	core.Println(coverageExample(GetVerifyByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleHashString() {
	core.Println(coverageExample(HashString("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashString() {
	core.Println(coverageExample(GetHashString("lethean.lthn")))
	// Output: 2
}

func ExampleHashByString() {
	core.Println(coverageExample(HashByString("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashByString() {
	core.Println(coverageExample(GetHashByString("lethean.lthn")))
	// Output: 2
}

func ExampleHashBinary() {
	core.Println(coverageExample(HashBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetHashBinary() {
	core.Println(coverageExample(GetHashBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleHashByBinary() {
	core.Println(coverageExample(HashByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetHashByBinary() {
	core.Println(coverageExample(GetHashByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleHashName() {
	core.Println(coverageExample(HashName("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashName() {
	core.Println(coverageExample(GetHashName("lethean.lthn")))
	// Output: 2
}

func ExampleHashByName() {
	core.Println(coverageExample(HashByName("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashByName() {
	core.Println(coverageExample(GetHashByName("lethean.lthn")))
	// Output: 2
}
