// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNewNameView() {
	fn := ax7ValueOf(NewNameView)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameView_GetNameStateSync() {
	fn := ax7ValueOf(new(NameView).GetNameStateSync)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleNameView_GetNameState() {
	fn := ax7ValueOf(new(NameView).GetNameState)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleNameView_ToNameUndo() {
	fn := ax7ValueOf(new(NameView).ToNameUndo)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameUndo_FromView() {
	fn := ax7ValueOf(new(NameUndo).FromView)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameUndo_GetSize() {
	fn := ax7ValueOf(new(NameUndo).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameUndo_MarshalBinary() {
	fn := ax7ValueOf(new(NameUndo).MarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleNameUndo_UnmarshalBinary() {
	fn := ax7ValueOf(new(NameUndo).UnmarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
