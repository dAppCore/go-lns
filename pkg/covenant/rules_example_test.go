// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleGetVerificationFlags() {
	fn := ax7ValueOf(GetVerificationFlags)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerificationFlagsByVal() {
	fn := ax7ValueOf(GetVerificationFlagsByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetRollout() {
	fn := ax7ValueOf(GetRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasRollout() {
	fn := ax7ValueOf(HasRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasRollout() {
	fn := ax7ValueOf(GetHasRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleIsReserved() {
	fn := ax7ValueOf(IsReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetIsReserved() {
	fn := ax7ValueOf(GetIsReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleIsLockedUp() {
	fn := ax7ValueOf(IsLockedUp)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetIsLockedUp() {
	fn := ax7ValueOf(GetIsLockedUp)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
