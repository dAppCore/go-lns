// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleCovenant_MarshalBinary() {
	fn := ax7ValueOf(new(Covenant).MarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCovenant_UnmarshalBinary() {
	fn := ax7ValueOf(new(Covenant).UnmarshalBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
