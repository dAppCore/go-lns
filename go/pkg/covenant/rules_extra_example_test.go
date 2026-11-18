// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "github.com/dappcore/go"

func ExampleGrindName() {
	core.Println(coverageExample(GrindName(1, 1, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleGetGrindName() {
	core.Println(coverageExample(GetGrindName(1, 1, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleCountOpens() {
	core.Println(coverageExample(CountOpens(coverageTransaction("good"))))
	// Output: 1
}

func ExampleGetCountOpens() {
	core.Println(coverageExample(GetCountOpens(coverageTransaction("good"))))
	// Output: 1
}

func ExampleCountUpdates() {
	core.Println(coverageExample(CountUpdates(coverageTransaction("good"))))
	// Output: 1
}

func ExampleGetCountUpdates() {
	core.Println(coverageExample(GetCountUpdates(coverageTransaction("good"))))
	// Output: 1
}

func ExampleCountRenewals() {
	core.Println(coverageExample(CountRenewals(coverageTransaction("good"))))
	// Output: 1
}

func ExampleGetCountRenewals() {
	core.Println(coverageExample(GetCountRenewals(coverageTransaction("good"))))
	// Output: 1
}
