// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleOutpoint_GetJSON() {
	core.Println(coverageExample(new(Outpoint).GetJSON()))
	// Output: 1
}

func ExampleOutpoint_FromJSON() {
	core.Println(coverageExample(new(Outpoint).FromJSON(coverageValue[OutpointJSON]("good"))))
	// Output: 1
}
