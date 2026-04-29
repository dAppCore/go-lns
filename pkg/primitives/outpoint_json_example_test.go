// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleOutpoint_GetJSON() {
	fn := ax7ValueOf(new(Outpoint).GetJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleOutpoint_FromJSON() {
	fn := ax7ValueOf(new(Outpoint).FromJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
