// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleCovenant_GetJSON() {
	fn := ax7ValueOf(new(Covenant).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCovenant_FromJSON() {
	fn := ax7ValueOf(new(Covenant).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
