// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleHasReservedName() {
	fn := ax7ValueOf(HasReservedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedName() {
	fn := ax7ValueOf(GetHasReservedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasReservedByName() {
	fn := ax7ValueOf(HasReservedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasReservedString() {
	fn := ax7ValueOf(HasReservedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedString() {
	fn := ax7ValueOf(GetHasReservedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasReservedBinary() {
	fn := ax7ValueOf(HasReservedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedBinary() {
	fn := ax7ValueOf(GetHasReservedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedName() {
	fn := ax7ValueOf(GetReservedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetReservedByName() {
	fn := ax7ValueOf(GetReservedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetReservedString() {
	fn := ax7ValueOf(GetReservedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetReservedBinary() {
	fn := ax7ValueOf(GetReservedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasReservedByString() {
	fn := ax7ValueOf(HasReservedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedByString() {
	fn := ax7ValueOf(GetHasReservedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasReservedByBinary() {
	fn := ax7ValueOf(HasReservedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedByBinary() {
	fn := ax7ValueOf(GetHasReservedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedByString() {
	fn := ax7ValueOf(GetReservedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetReservedByBinary() {
	fn := ax7ValueOf(GetReservedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetReservedHash() {
	fn := ax7ValueOf(GetReservedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasReservedHash() {
	fn := ax7ValueOf(HasReservedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedHash() {
	fn := ax7ValueOf(GetHasReservedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedByHash() {
	fn := ax7ValueOf(GetReservedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasReservedByHash() {
	fn := ax7ValueOf(HasReservedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedByHash() {
	fn := ax7ValueOf(GetHasReservedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedByName() {
	fn := ax7ValueOf(GetHasReservedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDefaultReservedCatalog() {
	fn := ax7ValueOf(DefaultReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedCatalog() {
	fn := ax7ValueOf(GetReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetDefaultReservedCatalog() {
	fn := ax7ValueOf(GetDefaultReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_Size() {
	fn := ax7ValueOf(new(ReservedCatalog).Size)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetSize() {
	fn := ax7ValueOf(new(ReservedCatalog).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_NameValue() {
	fn := ax7ValueOf(new(ReservedCatalog).NameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetNameValue() {
	fn := ax7ValueOf(new(ReservedCatalog).GetNameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_RootValue() {
	fn := ax7ValueOf(new(ReservedCatalog).RootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetRootValue() {
	fn := ax7ValueOf(new(ReservedCatalog).GetRootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_TopValue() {
	fn := ax7ValueOf(new(ReservedCatalog).TopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetTopValue() {
	fn := ax7ValueOf(new(ReservedCatalog).GetTopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_PrefixSize() {
	fn := ax7ValueOf(new(ReservedCatalog).PrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetPrefixSize() {
	fn := ax7ValueOf(new(ReservedCatalog).GetPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_Has() {
	fn := ax7ValueOf(new(ReservedCatalog).Has)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_HasByHash() {
	fn := ax7ValueOf(new(ReservedCatalog).HasByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetHasByHash() {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_Get() {
	fn := ax7ValueOf(new(ReservedCatalog).Get)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleReservedCatalog_GetByHash() {
	fn := ax7ValueOf(new(ReservedCatalog).GetByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleReservedCatalog_HasByName() {
	fn := ax7ValueOf(new(ReservedCatalog).HasByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetHasByName() {
	fn := ax7ValueOf(new(ReservedCatalog).GetHasByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetByName() {
	fn := ax7ValueOf(new(ReservedCatalog).GetByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleReservedCatalog_Entries() {
	fn := ax7ValueOf(new(ReservedCatalog).Entries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetEntries() {
	fn := ax7ValueOf(new(ReservedCatalog).GetEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_Keys() {
	fn := ax7ValueOf(new(ReservedCatalog).Keys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetKeys() {
	fn := ax7ValueOf(new(ReservedCatalog).GetKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_Values() {
	fn := ax7ValueOf(new(ReservedCatalog).Values)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedCatalog_GetValues() {
	fn := ax7ValueOf(new(ReservedCatalog).GetValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
