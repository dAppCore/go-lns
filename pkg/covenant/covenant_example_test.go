// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleGetTypes() {
	fn := ax7ValueOf(GetTypes)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypesByVal() {
	fn := ax7ValueOf(GetTypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleTypeName() {
	fn := ax7ValueOf(TypeName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeName() {
	fn := ax7ValueOf(GetTypeName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenantType_IsName() {
	fn := ax7ValueOf(new(CovenantType).IsName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenantType_IsKnown() {
	fn := ax7ValueOf(new(CovenantType).IsKnown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenantType_IsLinked() {
	fn := ax7ValueOf(new(CovenantType).IsLinked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenantType_String() {
	fn := ax7ValueOf(new(CovenantType).String)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
