// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "github.com/dappcore/go"

func ExampleHasNames() {
	core.Println(coverageExample(HasNames(coverageTransaction("good"), nil)))
	// Output: 1
}

func ExampleGetHasNames() {
	core.Println(coverageExample(GetHasNames(coverageTransaction("good"), nil)))
	// Output: 1
}

func ExampleAddNames() {
	AddNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleGetAddNames() {
	GetAddNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleRemoveNames() {
	RemoveNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleGetRemoveNames() {
	GetRemoveNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleHasSaneCovenants() {
	core.Println(coverageExample(HasSaneCovenants(coverageTransaction("good"))))
	// Output: 1
}

func ExampleGetHasSaneCovenants() {
	core.Println(coverageExample(GetHasSaneCovenants(coverageTransaction("good"))))
	// Output: 1
}
