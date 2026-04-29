// SPDX-License-Identifier: EUPL-1.2

package nameutil

import core "dappco.re/go"

func ExampleCatalogLabel() {
	fn := ax7ValueOf(CatalogLabel)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCanonicalize() {
	fn := ax7ValueOf(Canonicalize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}
