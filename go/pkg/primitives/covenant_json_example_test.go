// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

func ExampleCovenant_GetJSON() {
	core.Println(coverageExample(new(Covenant).GetJSON()))
	// Output: 1
}

func ExampleCovenant_FromJSON() {
	core.Println(coverageExample(new(Covenant).FromJSON(coverageValue[CovenantJSON]("good"))))
	// Output: 1
}
