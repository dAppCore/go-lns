// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "github.com/dappcore/go"

func Examplejsonraw_MarshalJSON() {
	core.Println(coverageExample(new(jsonraw).MarshalJSON()))
	// Output: 2
}

func Examplejsonraw_UnmarshalJSON() {
	core.Println(coverageExample(new(jsonraw).UnmarshalJSON([]byte("lethean.lthn"))))
	// Output: 1
}
