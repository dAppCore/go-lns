// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleHasReservedName() {
	core.Println(coverageExample(HasReservedName("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasReservedName() {
	core.Println(coverageExample(GetHasReservedName("lethean.lthn")))
	// Output: 1
}

func ExampleHasReservedByName() {
	core.Println(coverageExample(HasReservedByName("lethean.lthn")))
	// Output: 1
}

func ExampleHasReservedString() {
	core.Println(coverageExample(HasReservedString("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasReservedString() {
	core.Println(coverageExample(GetHasReservedString("lethean.lthn")))
	// Output: 1
}

func ExampleHasReservedBinary() {
	core.Println(coverageExample(HasReservedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetHasReservedBinary() {
	core.Println(coverageExample(GetHasReservedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetReservedName() {
	core.Println(coverageExample(GetReservedName("lethean.lthn")))
	// Output: 2
}

func ExampleGetReservedByName() {
	core.Println(coverageExample(GetReservedByName("lethean.lthn")))
	// Output: 2
}

func ExampleGetReservedString() {
	core.Println(coverageExample(GetReservedString("lethean.lthn")))
	// Output: 2
}

func ExampleGetReservedBinary() {
	core.Println(coverageExample(GetReservedBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleHasReservedByString() {
	core.Println(coverageExample(HasReservedByString("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasReservedByString() {
	core.Println(coverageExample(GetHasReservedByString("lethean.lthn")))
	// Output: 1
}

func ExampleHasReservedByBinary() {
	core.Println(coverageExample(HasReservedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetHasReservedByBinary() {
	core.Println(coverageExample(GetHasReservedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetReservedByString() {
	core.Println(coverageExample(GetReservedByString("lethean.lthn")))
	// Output: 2
}

func ExampleGetReservedByBinary() {
	core.Println(coverageExample(GetReservedByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetReservedHash() {
	core.Println(coverageExample(GetReservedHash([32]byte{1})))
	// Output: 2
}

func ExampleHasReservedHash() {
	core.Println(coverageExample(HasReservedHash([32]byte{1})))
	// Output: 1
}

func ExampleGetHasReservedHash() {
	core.Println(coverageExample(GetHasReservedHash([32]byte{1})))
	// Output: 1
}

func ExampleGetReservedByHash() {
	core.Println(coverageExample(GetReservedByHash([32]byte{1})))
	// Output: 2
}

func ExampleHasReservedByHash() {
	core.Println(coverageExample(HasReservedByHash([32]byte{1})))
	// Output: 1
}

func ExampleGetHasReservedByHash() {
	core.Println(coverageExample(GetHasReservedByHash([32]byte{1})))
	// Output: 1
}

func ExampleGetHasReservedByName() {
	core.Println(coverageExample(GetHasReservedByName("lethean.lthn")))
	// Output: 1
}

func ExampleDefaultReservedCatalog() {
	core.Println(coverageExample(DefaultReservedCatalog()))
	// Output: 1
}

func ExampleGetReservedCatalog() {
	core.Println(coverageExample(GetReservedCatalog()))
	// Output: 1
}

func ExampleGetDefaultReservedCatalog() {
	core.Println(coverageExample(GetDefaultReservedCatalog()))
	// Output: 1
}

func ExampleReservedCatalog_Size() {
	core.Println(coverageExample(new(ReservedCatalog).Size()))
	// Output: 1
}

func ExampleReservedCatalog_GetSize() {
	core.Println(coverageExample(new(ReservedCatalog).GetSize()))
	// Output: 1
}

func ExampleReservedCatalog_NameValue() {
	core.Println(coverageExample(new(ReservedCatalog).NameValue()))
	// Output: 1
}

func ExampleReservedCatalog_GetNameValue() {
	core.Println(coverageExample(new(ReservedCatalog).GetNameValue()))
	// Output: 1
}

func ExampleReservedCatalog_RootValue() {
	core.Println(coverageExample(new(ReservedCatalog).RootValue()))
	// Output: 1
}

func ExampleReservedCatalog_GetRootValue() {
	core.Println(coverageExample(new(ReservedCatalog).GetRootValue()))
	// Output: 1
}

func ExampleReservedCatalog_TopValue() {
	core.Println(coverageExample(new(ReservedCatalog).TopValue()))
	// Output: 1
}

func ExampleReservedCatalog_GetTopValue() {
	core.Println(coverageExample(new(ReservedCatalog).GetTopValue()))
	// Output: 1
}

func ExampleReservedCatalog_PrefixSize() {
	core.Println(coverageExample(new(ReservedCatalog).PrefixSize()))
	// Output: 1
}

func ExampleReservedCatalog_GetPrefixSize() {
	core.Println(coverageExample(new(ReservedCatalog).GetPrefixSize()))
	// Output: 1
}

func ExampleReservedCatalog_Has() {
	core.Println(coverageExample(new(ReservedCatalog).Has([32]byte{1})))
	// Output: 1
}

func ExampleReservedCatalog_HasByHash() {
	core.Println(coverageExample(new(ReservedCatalog).HasByHash([32]byte{1})))
	// Output: 1
}

func ExampleReservedCatalog_GetHasByHash() {
	core.Println(coverageExample(new(ReservedCatalog).GetHasByHash([32]byte{1})))
	// Output: 1
}

func ExampleReservedCatalog_Get() {
	core.Println(coverageExample(new(ReservedCatalog).Get([32]byte{1})))
	// Output: 2
}

func ExampleReservedCatalog_GetByHash() {
	core.Println(coverageExample(new(ReservedCatalog).GetByHash([32]byte{1})))
	// Output: 2
}

func ExampleReservedCatalog_HasByName() {
	core.Println(coverageExample(new(ReservedCatalog).HasByName("lethean.lthn")))
	// Output: 1
}

func ExampleReservedCatalog_GetHasByName() {
	core.Println(coverageExample(new(ReservedCatalog).GetHasByName("lethean.lthn")))
	// Output: 1
}

func ExampleReservedCatalog_GetByName() {
	core.Println(coverageExample(new(ReservedCatalog).GetByName("lethean.lthn")))
	// Output: 2
}

func ExampleReservedCatalog_Entries() {
	core.Println(coverageExample(new(ReservedCatalog).Entries()))
	// Output: 1
}

func ExampleReservedCatalog_GetEntries() {
	core.Println(coverageExample(new(ReservedCatalog).GetEntries()))
	// Output: 1
}

func ExampleReservedCatalog_Keys() {
	core.Println(coverageExample(new(ReservedCatalog).Keys()))
	// Output: 1
}

func ExampleReservedCatalog_GetKeys() {
	core.Println(coverageExample(new(ReservedCatalog).GetKeys()))
	// Output: 1
}

func ExampleReservedCatalog_Values() {
	core.Println(coverageExample(new(ReservedCatalog).Values()))
	// Output: 1
}

func ExampleReservedCatalog_GetValues() {
	core.Println(coverageExample(new(ReservedCatalog).GetValues()))
	// Output: 1
}
