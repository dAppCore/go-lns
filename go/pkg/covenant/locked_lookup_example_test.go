// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleHasLockedName() {
	core.Println(coverageExample(HasLockedName("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLockedName() {
	core.Println(coverageExample(GetHasLockedName("lethean.lthn")))
	// Output: 1
}

func ExampleHasLockedByName() {
	core.Println(coverageExample(HasLockedByName("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLockedByName() {
	core.Println(coverageExample(GetHasLockedByName("lethean.lthn")))
	// Output: 1
}

func ExampleHasLockedString() {
	core.Println(coverageExample(HasLockedString("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLockedString() {
	core.Println(coverageExample(GetHasLockedString("lethean.lthn")))
	// Output: 1
}

func ExampleHasLockedBinary() {
	core.Println(coverageExample(HasLockedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetHasLockedBinary() {
	core.Println(coverageExample(GetHasLockedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetLockedName() {
	core.Println(coverageExample(GetLockedName("lethean.lthn")))
	// Output: 2
}

func ExampleGetLockedByName() {
	core.Println(coverageExample(GetLockedByName("lethean.lthn")))
	// Output: 2
}

func ExampleGetLockedString() {
	core.Println(coverageExample(GetLockedString("lethean.lthn")))
	// Output: 2
}

func ExampleGetLockedBinary() {
	core.Println(coverageExample(GetLockedBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleHasLockedByString() {
	core.Println(coverageExample(HasLockedByString("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLockedByString() {
	core.Println(coverageExample(GetHasLockedByString("lethean.lthn")))
	// Output: 1
}

func ExampleHasLockedByBinary() {
	core.Println(coverageExample(HasLockedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetHasLockedByBinary() {
	core.Println(coverageExample(GetHasLockedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetLockedByString() {
	core.Println(coverageExample(GetLockedByString("lethean.lthn")))
	// Output: 2
}

func ExampleGetLockedByBinary() {
	core.Println(coverageExample(GetLockedByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetLockedHash() {
	core.Println(coverageExample(GetLockedHash([32]byte{1})))
	// Output: 2
}

func ExampleHasLockedHash() {
	core.Println(coverageExample(HasLockedHash([32]byte{1})))
	// Output: 1
}

func ExampleGetHasLockedHash() {
	core.Println(coverageExample(GetHasLockedHash([32]byte{1})))
	// Output: 1
}

func ExampleGetLockedByHash() {
	core.Println(coverageExample(GetLockedByHash([32]byte{1})))
	// Output: 2
}

func ExampleHasLockedByHash() {
	core.Println(coverageExample(HasLockedByHash([32]byte{1})))
	// Output: 1
}

func ExampleGetHasLockedByHash() {
	core.Println(coverageExample(GetHasLockedByHash([32]byte{1})))
	// Output: 1
}

func ExampleDefaultLockedCatalog() {
	core.Println(coverageExample(DefaultLockedCatalog()))
	// Output: 1
}

func ExampleGetLockedCatalog() {
	core.Println(coverageExample(GetLockedCatalog()))
	// Output: 1
}

func ExampleGetDefaultLockedCatalog() {
	core.Println(coverageExample(GetDefaultLockedCatalog()))
	// Output: 1
}

func ExampleLockedCatalog_Size() {
	core.Println(coverageExample(new(LockedCatalog).Size()))
	// Output: 1
}

func ExampleLockedCatalog_GetSize() {
	core.Println(coverageExample(new(LockedCatalog).GetSize()))
	// Output: 1
}

func ExampleLockedCatalog_PrefixSize() {
	core.Println(coverageExample(new(LockedCatalog).PrefixSize()))
	// Output: 1
}

func ExampleLockedCatalog_GetPrefixSize() {
	core.Println(coverageExample(new(LockedCatalog).GetPrefixSize()))
	// Output: 1
}

func ExampleLockedCatalog_Has() {
	core.Println(coverageExample(new(LockedCatalog).Has([32]byte{1})))
	// Output: 1
}

func ExampleLockedCatalog_HasByHash() {
	core.Println(coverageExample(new(LockedCatalog).HasByHash([32]byte{1})))
	// Output: 1
}

func ExampleLockedCatalog_GetHasByHash() {
	core.Println(coverageExample(new(LockedCatalog).GetHasByHash([32]byte{1})))
	// Output: 1
}

func ExampleLockedCatalog_Get() {
	core.Println(coverageExample(new(LockedCatalog).Get([32]byte{1})))
	// Output: 2
}

func ExampleLockedCatalog_GetByHash() {
	core.Println(coverageExample(new(LockedCatalog).GetByHash([32]byte{1})))
	// Output: 2
}

func ExampleLockedCatalog_HasByName() {
	core.Println(coverageExample(new(LockedCatalog).HasByName("lethean.lthn")))
	// Output: 1
}

func ExampleLockedCatalog_GetHasByName() {
	core.Println(coverageExample(new(LockedCatalog).GetHasByName("lethean.lthn")))
	// Output: 1
}

func ExampleLockedCatalog_GetByName() {
	core.Println(coverageExample(new(LockedCatalog).GetByName("lethean.lthn")))
	// Output: 2
}

func ExampleLockedCatalog_Entries() {
	core.Println(coverageExample(new(LockedCatalog).Entries()))
	// Output: 1
}

func ExampleLockedCatalog_GetEntries() {
	core.Println(coverageExample(new(LockedCatalog).GetEntries()))
	// Output: 1
}

func ExampleLockedCatalog_Keys() {
	core.Println(coverageExample(new(LockedCatalog).Keys()))
	// Output: 1
}

func ExampleLockedCatalog_GetKeys() {
	core.Println(coverageExample(new(LockedCatalog).GetKeys()))
	// Output: 1
}

func ExampleLockedCatalog_Values() {
	core.Println(coverageExample(new(LockedCatalog).Values()))
	// Output: 1
}

func ExampleLockedCatalog_GetValues() {
	core.Println(coverageExample(new(LockedCatalog).GetValues()))
	// Output: 1
}
