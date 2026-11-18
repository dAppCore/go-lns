// SPDX-License-Identifier: EUPL-1.2

package chain

import core "github.com/dappcore/go"

func ExampleWithAliasHTTPClient() {
	core.Println(coverageExample(WithAliasHTTPClient(nil)))
	// Output: 1
}

func ExampleNewAliasClient() {
	core.Println(coverageExample(NewAliasClient("lethean.lthn")))
	// Output: 1
}

func ExampleAliasClient_GetAliases() {
	core.Println(coverageExample(new(AliasClient).GetAliases()))
	// Output: 2
}

func ExampleParseAliasComment() {
	core.Println(coverageExample(ParseAliasComment("lethean.lthn")))
	// Output: 1
}

func ExampleHNSName() {
	core.Println(coverageExample(HNSName(coverageValue[Alias]("good"))))
	// Output: 1
}

func ExampleDiscoverNames() {
	core.Println(coverageExample(DiscoverNames(nil, []string{"lethean"})))
	// Output: 1
}
