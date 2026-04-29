// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleCreate() {
	fn := ax7ValueOf(Create)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCreate() {
	fn := ax7ValueOf(GetCreate)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNextName() {
	fn := ax7ValueOf(NextName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNextName() {
	fn := ax7ValueOf(GetNextName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExamplePrevName() {
	fn := ax7ValueOf(PrevName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetPrevName() {
	fn := ax7ValueOf(GetPrevName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
