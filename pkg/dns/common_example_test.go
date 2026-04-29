// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleGetDummy() {
	fn := ax7ValueOf(GetDummy)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetDefaultTTL() {
	fn := ax7ValueOf(GetDefaultTTL)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapRoot() {
	fn := ax7ValueOf(GetTypeMapRoot)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapEmpty() {
	fn := ax7ValueOf(GetTypeMapEmpty)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapNS() {
	fn := ax7ValueOf(GetTypeMapNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapTXT() {
	fn := ax7ValueOf(GetTypeMapTXT)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapA() {
	fn := ax7ValueOf(GetTypeMapA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapAAAA() {
	fn := ax7ValueOf(GetTypeMapAAAA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHSTypes() {
	fn := ax7ValueOf(GetHSTypes)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHSTypesByVal() {
	fn := ax7ValueOf(GetHSTypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
