// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleHasNames() {
	fn := ax7ValueOf(HasNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasNames() {
	fn := ax7ValueOf(GetHasNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddNames() {
	fn := ax7ValueOf(AddNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleGetAddNames() {
	fn := ax7ValueOf(GetAddNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleRemoveNames() {
	fn := ax7ValueOf(RemoveNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleGetRemoveNames() {
	fn := ax7ValueOf(GetRemoveNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleHasSaneCovenants() {
	fn := ax7ValueOf(HasSaneCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasSaneCovenants() {
	fn := ax7ValueOf(GetHasSaneCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
