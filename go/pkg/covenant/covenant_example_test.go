// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleGetTypes() {
	core.Println(coverageExample(GetTypes()))
	// Output: 1
}

func ExampleGetTypesByVal() {
	core.Println(coverageExample(GetTypesByVal()))
	// Output: 1
}

func ExampleTypeName() {
	core.Println(coverageExample(TypeName(1)))
	// Output: 1
}

func ExampleGetTypeName() {
	core.Println(coverageExample(GetTypeName(1)))
	// Output: 1
}

func ExampleCovenantType_IsName() {
	core.Println(coverageExample(new(CovenantType).IsName()))
	// Output: 1
}

func ExampleCovenantType_IsKnown() {
	core.Println(coverageExample(new(CovenantType).IsKnown()))
	// Output: 1
}

func ExampleCovenantType_IsLinked() {
	core.Println(coverageExample(new(CovenantType).IsLinked()))
	// Output: 1
}

func ExampleCovenantType_String() {
	core.Println(coverageExample(new(CovenantType).String()))
	// Output: 1
}
