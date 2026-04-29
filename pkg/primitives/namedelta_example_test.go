// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNameDelta_IsNull() {
	fn := ax7ValueOf(new(NameDelta).IsNull)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameDelta_Clone() {
	fn := ax7ValueOf(new(NameDelta).Clone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameDelta_Inject() {
	fn := ax7ValueOf(new(NameDelta).Inject)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameDelta_Clear() {
	fn := ax7ValueOf(new(NameDelta).Clear)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameDelta_GetField() {
	fn := ax7ValueOf(new(NameDelta).GetField)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameDelta_GetSize() {
	fn := ax7ValueOf(new(NameDelta).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameDelta_MarshalBinary() {
	fn := ax7ValueOf(new(NameDelta).MarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleNameDelta_UnmarshalBinary() {
	fn := ax7ValueOf(new(NameDelta).UnmarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
