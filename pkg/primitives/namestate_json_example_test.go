// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNameState_GetJSON() {
	fn := ax7ValueOf(new(NameState).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_Format() {
	fn := ax7ValueOf(new(NameState).Format)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetFormat() {
	fn := ax7ValueOf(new(NameState).GetFormat)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_FromJSON() {
	fn := ax7ValueOf(new(NameState).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameState_GetFromJSON() {
	fn := ax7ValueOf(new(NameState).GetFromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
