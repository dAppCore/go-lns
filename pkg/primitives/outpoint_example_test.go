// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNewOutpoint() {
	fn := ax7ValueOf(NewOutpoint)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_Inject() {
	fn := ax7ValueOf(new(Outpoint).Inject)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_Clone() {
	fn := ax7ValueOf(new(Outpoint).Clone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_Equals() {
	fn := ax7ValueOf(new(Outpoint).Equals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_Compare() {
	fn := ax7ValueOf(new(Outpoint).Compare)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_GetSize() {
	fn := ax7ValueOf(new(Outpoint).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_MarshalBinary() {
	fn := ax7ValueOf(new(Outpoint).MarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleOutpoint_UnmarshalBinary() {
	fn := ax7ValueOf(new(Outpoint).UnmarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
