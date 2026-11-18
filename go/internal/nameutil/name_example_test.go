// SPDX-License-Identifier: EUPL-1.2

package nameutil

import core "github.com/dappcore/go"

func ExampleCatalogLabel() {
	core.Println(coverageExample(CatalogLabel("lethean.lthn")))
	// Output: 2
}

func ExampleCanonicalize() {
	core.Println(coverageExample(Canonicalize("lethean.lthn")))
	// Output: 2
}
