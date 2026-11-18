// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "github.com/dappcore/go"

func ExampleVerifyCovenants() {
	core.Println(coverageExample(VerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good"))))
	// Output: 1
}

func ExampleGetVerifyCovenants() {
	core.Println(coverageExample(GetVerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good"))))
	// Output: 1
}
