// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func Examplejsonraw_MarshalJSON() {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func Examplejsonraw_UnmarshalJSON() {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
