// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleGrindName() {
	fn := ax7ValueOf(GrindName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetGrindName() {
	fn := ax7ValueOf(GetGrindName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCountOpens() {
	fn := ax7ValueOf(CountOpens)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCountOpens() {
	fn := ax7ValueOf(GetCountOpens)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCountUpdates() {
	fn := ax7ValueOf(CountUpdates)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCountUpdates() {
	fn := ax7ValueOf(GetCountUpdates)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCountRenewals() {
	fn := ax7ValueOf(CountRenewals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCountRenewals() {
	fn := ax7ValueOf(GetCountRenewals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
