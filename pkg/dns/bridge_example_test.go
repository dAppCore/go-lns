// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleNewBridge() {
	fn := ax7ValueOf(NewBridge)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleBridge_Discover() {
	fn := ax7ValueOf(new(Bridge).Discover)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleBridge_Lookup() {
	fn := ax7ValueOf(new(Bridge).Lookup)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleBridge_Records() {
	fn := ax7ValueOf(new(Bridge).Records)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleBridge_Reverse() {
	fn := ax7ValueOf(new(Bridge).Reverse)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleBridge_Health() {
	fn := ax7ValueOf(new(Bridge).Health)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleBridge_Serve() {
	fn := ax7ValueOf(new(Bridge).Serve)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleBridge_Shutdown() {
	fn := ax7ValueOf(new(Bridge).Shutdown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameRecord_FQDN() {
	fn := ax7ValueOf(new(NameRecord).FQDN)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameRecord_Addresses4() {
	fn := ax7ValueOf(new(NameRecord).Addresses4)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameRecord_Addresses6() {
	fn := ax7ValueOf(new(NameRecord).Addresses6)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameRecord_TXT() {
	fn := ax7ValueOf(new(NameRecord).TXT)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameRecord_NS() {
	fn := ax7ValueOf(new(NameRecord).NS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
