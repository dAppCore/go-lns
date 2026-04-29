// SPDX-License-Identifier: EUPL-1.2

package chain

import core "dappco.re/go"

func ExampleWithAliasHTTPClient() {
	fn := ax7ValueOf(WithAliasHTTPClient)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewAliasClient() {
	fn := ax7ValueOf(NewAliasClient)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAliasClient_GetAliases() {
	fn := ax7ValueOf(new(AliasClient).GetAliases)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleParseAliasComment() {
	fn := ax7ValueOf(ParseAliasComment)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHNSName() {
	fn := ax7ValueOf(HNSName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDiscoverNames() {
	fn := ax7ValueOf(DiscoverNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
