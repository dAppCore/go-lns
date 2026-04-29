// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

func ExampleHasLockedName() {
	fn := ax7ValueOf(HasLockedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedName() {
	fn := ax7ValueOf(GetHasLockedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasLockedByName() {
	fn := ax7ValueOf(HasLockedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedByName() {
	fn := ax7ValueOf(GetHasLockedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasLockedString() {
	fn := ax7ValueOf(HasLockedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedString() {
	fn := ax7ValueOf(GetHasLockedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasLockedBinary() {
	fn := ax7ValueOf(HasLockedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedBinary() {
	fn := ax7ValueOf(GetHasLockedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedName() {
	fn := ax7ValueOf(GetLockedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetLockedByName() {
	fn := ax7ValueOf(GetLockedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetLockedString() {
	fn := ax7ValueOf(GetLockedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetLockedBinary() {
	fn := ax7ValueOf(GetLockedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasLockedByString() {
	fn := ax7ValueOf(HasLockedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedByString() {
	fn := ax7ValueOf(GetHasLockedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasLockedByBinary() {
	fn := ax7ValueOf(HasLockedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedByBinary() {
	fn := ax7ValueOf(GetHasLockedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedByString() {
	fn := ax7ValueOf(GetLockedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetLockedByBinary() {
	fn := ax7ValueOf(GetLockedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetLockedHash() {
	fn := ax7ValueOf(GetLockedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasLockedHash() {
	fn := ax7ValueOf(HasLockedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedHash() {
	fn := ax7ValueOf(GetHasLockedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedByHash() {
	fn := ax7ValueOf(GetLockedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasLockedByHash() {
	fn := ax7ValueOf(HasLockedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedByHash() {
	fn := ax7ValueOf(GetHasLockedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDefaultLockedCatalog() {
	fn := ax7ValueOf(DefaultLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedCatalog() {
	fn := ax7ValueOf(GetLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetDefaultLockedCatalog() {
	fn := ax7ValueOf(GetDefaultLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_Size() {
	fn := ax7ValueOf(new(LockedCatalog).Size)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetSize() {
	fn := ax7ValueOf(new(LockedCatalog).GetSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_PrefixSize() {
	fn := ax7ValueOf(new(LockedCatalog).PrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetPrefixSize() {
	fn := ax7ValueOf(new(LockedCatalog).GetPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_Has() {
	fn := ax7ValueOf(new(LockedCatalog).Has)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_HasByHash() {
	fn := ax7ValueOf(new(LockedCatalog).HasByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetHasByHash() {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_Get() {
	fn := ax7ValueOf(new(LockedCatalog).Get)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleLockedCatalog_GetByHash() {
	fn := ax7ValueOf(new(LockedCatalog).GetByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleLockedCatalog_HasByName() {
	fn := ax7ValueOf(new(LockedCatalog).HasByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetHasByName() {
	fn := ax7ValueOf(new(LockedCatalog).GetHasByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetByName() {
	fn := ax7ValueOf(new(LockedCatalog).GetByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleLockedCatalog_Entries() {
	fn := ax7ValueOf(new(LockedCatalog).Entries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetEntries() {
	fn := ax7ValueOf(new(LockedCatalog).GetEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_Keys() {
	fn := ax7ValueOf(new(LockedCatalog).Keys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetKeys() {
	fn := ax7ValueOf(new(LockedCatalog).GetKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_Values() {
	fn := ax7ValueOf(new(LockedCatalog).Values)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog_GetValues() {
	fn := ax7ValueOf(new(LockedCatalog).GetValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
