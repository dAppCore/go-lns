// SPDX-License-Identifier: EUPL-1.2

package lns

import core "dappco.re/go"

func ExampleWithCore() {
	fn := ax7ValueOf(WithCore)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetWithCore() {
	fn := ax7ValueOf(GetWithCore)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWithReservedCatalog() {
	fn := ax7ValueOf(WithReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetWithReservedCatalog() {
	fn := ax7ValueOf(GetWithReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWithLockedCatalog() {
	fn := ax7ValueOf(WithLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetWithLockedCatalog() {
	fn := ax7ValueOf(GetWithLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewServiceWithOptions() {
	fn := ax7ValueOf(NewServiceWithOptions)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewServiceWithOptions() {
	fn := ax7ValueOf(GetNewServiceWithOptions)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewService() {
	fn := ax7ValueOf(NewService)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewService() {
	fn := ax7ValueOf(GetNewService)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleRegister() {
	fn := ax7ValueOf(Register)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetRegister() {
	fn := ax7ValueOf(GetRegister)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleRegisterWithOptions() {
	fn := ax7ValueOf(RegisterWithOptions)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetRegisterWithOptions() {
	fn := ax7ValueOf(GetRegisterWithOptions)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetServiceName() {
	fn := ax7ValueOf(GetServiceName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetDefaultTTL() {
	fn := ax7ValueOf(GetDefaultTTL)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetDummy() {
	fn := ax7ValueOf(GetDummy)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapRoot() {
	fn := ax7ValueOf(GetTypeMapRoot)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapEmpty() {
	fn := ax7ValueOf(GetTypeMapEmpty)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapNS() {
	fn := ax7ValueOf(GetTypeMapNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapTXT() {
	fn := ax7ValueOf(GetTypeMapTXT)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapA() {
	fn := ax7ValueOf(GetTypeMapA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeMapAAAA() {
	fn := ax7ValueOf(GetTypeMapAAAA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHSTypes() {
	fn := ax7ValueOf(GetHSTypes)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHSTypesByVal() {
	fn := ax7ValueOf(GetHSTypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewResource() {
	fn := ax7ValueOf(NewResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewResource() {
	fn := ax7ValueOf(GetNewResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleDecodeResource() {
	fn := ax7ValueOf(DecodeResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetDecodeResource() {
	fn := ax7ValueOf(GetDecodeResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCreate() {
	fn := ax7ValueOf(Create)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCreate() {
	fn := ax7ValueOf(GetCreate)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNextName() {
	fn := ax7ValueOf(NextName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNextName() {
	fn := ax7ValueOf(GetNextName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExamplePrevName() {
	fn := ax7ValueOf(PrevName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetPrevName() {
	fn := ax7ValueOf(GetPrevName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_NewClaim() {
	fn := ax7ValueOf(new(Service).NewClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetNewClaim() {
	fn := ax7ValueOf(new(Service).GetNewClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_NewOutpoint() {
	fn := ax7ValueOf(new(Service).NewOutpoint)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetNewOutpoint() {
	fn := ax7ValueOf(new(Service).GetNewOutpoint)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_NewNameView() {
	fn := ax7ValueOf(new(Service).NewNameView)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetNewNameView() {
	fn := ax7ValueOf(new(Service).GetNewNameView)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_Create() {
	fn := ax7ValueOf(new(Service).Create)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetCreate() {
	fn := ax7ValueOf(new(Service).GetCreate)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_DefaultTTL() {
	fn := ax7ValueOf(new(Service).DefaultTTL)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetDefaultTTL() {
	fn := ax7ValueOf(new(Service).GetDefaultTTL)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_Dummy() {
	fn := ax7ValueOf(new(Service).Dummy)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetDummy() {
	fn := ax7ValueOf(new(Service).GetDummy)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TypeMapRoot() {
	fn := ax7ValueOf(new(Service).TypeMapRoot)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypeMapRoot() {
	fn := ax7ValueOf(new(Service).GetTypeMapRoot)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TypeMapEmpty() {
	fn := ax7ValueOf(new(Service).TypeMapEmpty)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypeMapEmpty() {
	fn := ax7ValueOf(new(Service).GetTypeMapEmpty)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TypeMapNS() {
	fn := ax7ValueOf(new(Service).TypeMapNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypeMapNS() {
	fn := ax7ValueOf(new(Service).GetTypeMapNS)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TypeMapTXT() {
	fn := ax7ValueOf(new(Service).TypeMapTXT)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypeMapTXT() {
	fn := ax7ValueOf(new(Service).GetTypeMapTXT)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TypeMapA() {
	fn := ax7ValueOf(new(Service).TypeMapA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypeMapA() {
	fn := ax7ValueOf(new(Service).GetTypeMapA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TypeMapAAAA() {
	fn := ax7ValueOf(new(Service).TypeMapAAAA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypeMapAAAA() {
	fn := ax7ValueOf(new(Service).GetTypeMapAAAA)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HSTypes() {
	fn := ax7ValueOf(new(Service).HSTypes)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHSTypes() {
	fn := ax7ValueOf(new(Service).GetHSTypes)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HSTypesByVal() {
	fn := ax7ValueOf(new(Service).HSTypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHSTypesByVal() {
	fn := ax7ValueOf(new(Service).GetHSTypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_NextName() {
	fn := ax7ValueOf(new(Service).NextName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetNextName() {
	fn := ax7ValueOf(new(Service).GetNextName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_PrevName() {
	fn := ax7ValueOf(new(Service).PrevName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetPrevName() {
	fn := ax7ValueOf(new(Service).GetPrevName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_NewResource() {
	fn := ax7ValueOf(new(Service).NewResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetNewResource() {
	fn := ax7ValueOf(new(Service).GetNewResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_DecodeResource() {
	fn := ax7ValueOf(new(Service).DecodeResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetDecodeResource() {
	fn := ax7ValueOf(new(Service).GetDecodeResource)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_OnStartup() {
	fn := ax7ValueOf(new(Service).OnStartup)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetOnStartup() {
	fn := ax7ValueOf(new(Service).GetOnStartup)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_OnShutdown() {
	fn := ax7ValueOf(new(Service).OnShutdown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetOnShutdown() {
	fn := ax7ValueOf(new(Service).GetOnShutdown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HandleIPCEvents() {
	fn := ax7ValueOf(new(Service).HandleIPCEvents)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHandleIPCEvents() {
	fn := ax7ValueOf(new(Service).GetHandleIPCEvents)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLookupCatalogName() {
	fn := ax7ValueOf(LookupCatalogName[any])
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetLookupCatalogName() {
	fn := ax7ValueOf(GetLookupCatalogName[any])
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasLookupCatalogName() {
	fn := ax7ValueOf(HasLookupCatalogName[any])
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLookupCatalogName() {
	fn := ax7ValueOf(GetHasLookupCatalogName[any])
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_LookupCatalogName() {
	fn := ax7ValueOf(new(Service).LookupCatalogName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetLookupCatalogName() {
	fn := ax7ValueOf(new(Service).GetLookupCatalogName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasLookupCatalogName() {
	fn := ax7ValueOf(new(Service).HasLookupCatalogName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasLookupCatalogName() {
	fn := ax7ValueOf(new(Service).GetHasLookupCatalogName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResolve() {
	fn := ax7ValueOf(Resolve)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetResolve() {
	fn := ax7ValueOf(GetResolve)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHash() {
	fn := ax7ValueOf(Hash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetHash() {
	fn := ax7ValueOf(GetHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleVerify() {
	fn := ax7ValueOf(Verify)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerify() {
	fn := ax7ValueOf(GetVerify)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHashString() {
	fn := ax7ValueOf(HashString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetHashString() {
	fn := ax7ValueOf(GetHashString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHashBinary() {
	fn := ax7ValueOf(HashBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetHashBinary() {
	fn := ax7ValueOf(GetHashBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleResolveByString() {
	fn := ax7ValueOf(ResolveByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetResolveByString() {
	fn := ax7ValueOf(GetResolveByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleResolveByBinary() {
	fn := ax7ValueOf(ResolveByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetResolveByBinary() {
	fn := ax7ValueOf(GetResolveByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHashByString() {
	fn := ax7ValueOf(HashByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetHashByString() {
	fn := ax7ValueOf(GetHashByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHashByBinary() {
	fn := ax7ValueOf(HashByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetHashByBinary() {
	fn := ax7ValueOf(GetHashByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleResolveString() {
	fn := ax7ValueOf(ResolveString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetResolveString() {
	fn := ax7ValueOf(GetResolveString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleResolveBinary() {
	fn := ax7ValueOf(ResolveBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetResolveBinary() {
	fn := ax7ValueOf(GetResolveBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleResolveName() {
	fn := ax7ValueOf(ResolveName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetResolveName() {
	fn := ax7ValueOf(GetResolveName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHashName() {
	fn := ax7ValueOf(HashName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetHashName() {
	fn := ax7ValueOf(GetHashName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleVerifyName() {
	fn := ax7ValueOf(VerifyName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyName() {
	fn := ax7ValueOf(GetVerifyName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleResolveByName() {
	fn := ax7ValueOf(ResolveByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetResolveByName() {
	fn := ax7ValueOf(GetResolveByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHashByName() {
	fn := ax7ValueOf(HashByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetHashByName() {
	fn := ax7ValueOf(GetHashByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleVerifyByName() {
	fn := ax7ValueOf(VerifyByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyByName() {
	fn := ax7ValueOf(GetVerifyByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleVerifyString() {
	fn := ax7ValueOf(VerifyString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyString() {
	fn := ax7ValueOf(GetVerifyString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleVerifyBinary() {
	fn := ax7ValueOf(VerifyBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyBinary() {
	fn := ax7ValueOf(GetVerifyBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleVerifyByString() {
	fn := ax7ValueOf(VerifyByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyByString() {
	fn := ax7ValueOf(GetVerifyByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleVerifyByBinary() {
	fn := ax7ValueOf(VerifyByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyByBinary() {
	fn := ax7ValueOf(GetVerifyByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleBlind() {
	fn := ax7ValueOf(Blind)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetBlind() {
	fn := ax7ValueOf(GetBlind)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleTypeName() {
	fn := ax7ValueOf(TypeName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetBlacklist() {
	fn := ax7ValueOf(GetBlacklist)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerificationFlags() {
	fn := ax7ValueOf(GetVerificationFlags)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerificationFlagsByVal() {
	fn := ax7ValueOf(GetVerificationFlagsByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewClaim() {
	fn := ax7ValueOf(NewClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewClaim() {
	fn := ax7ValueOf(GetNewClaim)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewOutpoint() {
	fn := ax7ValueOf(NewOutpoint)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewOutpoint() {
	fn := ax7ValueOf(GetNewOutpoint)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNewNameView() {
	fn := ax7ValueOf(NewNameView)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNewNameView() {
	fn := ax7ValueOf(GetNewNameView)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetMaxNameSize() {
	fn := ax7ValueOf(GetMaxNameSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetMaxResourceSize() {
	fn := ax7ValueOf(GetMaxResourceSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyCovenantsNone() {
	fn := ax7ValueOf(GetVerifyCovenantsNone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyCovenantsHardened() {
	fn := ax7ValueOf(GetVerifyCovenantsHardened)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyCovenantsLockup() {
	fn := ax7ValueOf(GetVerifyCovenantsLockup)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetMandatoryVerifyCovenantFlags() {
	fn := ax7ValueOf(GetMandatoryVerifyCovenantFlags)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetMaxCovenantSize() {
	fn := ax7ValueOf(GetMaxCovenantSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCovenantMaxSize() {
	fn := ax7ValueOf(GetCovenantMaxSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypes() {
	fn := ax7ValueOf(GetTypes)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypesByVal() {
	fn := ax7ValueOf(GetTypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTypeName() {
	fn := ax7ValueOf(GetTypeName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleIsName() {
	fn := ax7ValueOf(IsName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetIsName() {
	fn := ax7ValueOf(GetIsName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleIsKnown() {
	fn := ax7ValueOf(IsKnown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetIsKnown() {
	fn := ax7ValueOf(GetIsKnown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleIsLinked() {
	fn := ax7ValueOf(IsLinked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetIsLinked() {
	fn := ax7ValueOf(GetIsLinked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_MaxNameSize() {
	fn := ax7ValueOf(new(Service).MaxNameSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetMaxNameSize() {
	fn := ax7ValueOf(new(Service).GetMaxNameSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_MaxResourceSize() {
	fn := ax7ValueOf(new(Service).MaxResourceSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetMaxResourceSize() {
	fn := ax7ValueOf(new(Service).GetMaxResourceSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyCovenantsNone() {
	fn := ax7ValueOf(new(Service).VerifyCovenantsNone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyCovenantsNone() {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsNone)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyCovenantsHardened() {
	fn := ax7ValueOf(new(Service).VerifyCovenantsHardened)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyCovenantsHardened() {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsHardened)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyCovenantsLockup() {
	fn := ax7ValueOf(new(Service).VerifyCovenantsLockup)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyCovenantsLockup() {
	fn := ax7ValueOf(new(Service).GetVerifyCovenantsLockup)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_MandatoryVerifyCovenantFlags() {
	fn := ax7ValueOf(new(Service).MandatoryVerifyCovenantFlags)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetMandatoryVerifyCovenantFlags() {
	fn := ax7ValueOf(new(Service).GetMandatoryVerifyCovenantFlags)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_MaxCovenantSize() {
	fn := ax7ValueOf(new(Service).MaxCovenantSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetMaxCovenantSize() {
	fn := ax7ValueOf(new(Service).GetMaxCovenantSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_CovenantMaxSize() {
	fn := ax7ValueOf(new(Service).CovenantMaxSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetCovenantMaxSize() {
	fn := ax7ValueOf(new(Service).GetCovenantMaxSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetRollout() {
	fn := ax7ValueOf(GetRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasRollout() {
	fn := ax7ValueOf(HasRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasRollout() {
	fn := ax7ValueOf(GetHasRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleIsReserved() {
	fn := ax7ValueOf(IsReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetIsReserved() {
	fn := ax7ValueOf(GetIsReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleIsLockedUp() {
	fn := ax7ValueOf(IsLockedUp)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetIsLockedUp() {
	fn := ax7ValueOf(GetIsLockedUp)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGrindName() {
	fn := ax7ValueOf(GrindName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetGrindName() {
	fn := ax7ValueOf(GetGrindName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleCountOpens() {
	fn := ax7ValueOf(CountOpens)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCountOpens() {
	fn := ax7ValueOf(GetCountOpens)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCountUpdates() {
	fn := ax7ValueOf(CountUpdates)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCountUpdates() {
	fn := ax7ValueOf(GetCountUpdates)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleCountRenewals() {
	fn := ax7ValueOf(CountRenewals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetCountRenewals() {
	fn := ax7ValueOf(GetCountRenewals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasNames() {
	fn := ax7ValueOf(HasNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasNames() {
	fn := ax7ValueOf(GetHasNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleAddNames() {
	fn := ax7ValueOf(AddNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleGetAddNames() {
	fn := ax7ValueOf(GetAddNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleRemoveNames() {
	fn := ax7ValueOf(RemoveNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleGetRemoveNames() {
	fn := ax7ValueOf(GetRemoveNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleHasSaneCovenants() {
	fn := ax7ValueOf(HasSaneCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasSaneCovenants() {
	fn := ax7ValueOf(GetHasSaneCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleVerifyCovenants() {
	fn := ax7ValueOf(VerifyCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetVerifyCovenants() {
	fn := ax7ValueOf(GetVerifyCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ServiceName() {
	fn := ax7ValueOf(new(Service).ServiceName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetServiceName() {
	fn := ax7ValueOf(new(Service).GetServiceName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_Resolve() {
	fn := ax7ValueOf(new(Service).Resolve)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetResolve() {
	fn := ax7ValueOf(new(Service).GetResolve)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_Hash() {
	fn := ax7ValueOf(new(Service).Hash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetHash() {
	fn := ax7ValueOf(new(Service).GetHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_Verify() {
	fn := ax7ValueOf(new(Service).Verify)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerify() {
	fn := ax7ValueOf(new(Service).GetVerify)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HashString() {
	fn := ax7ValueOf(new(Service).HashString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetHashString() {
	fn := ax7ValueOf(new(Service).GetHashString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HashBinary() {
	fn := ax7ValueOf(new(Service).HashBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetHashBinary() {
	fn := ax7ValueOf(new(Service).GetHashBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_ResolveByString() {
	fn := ax7ValueOf(new(Service).ResolveByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetResolveByString() {
	fn := ax7ValueOf(new(Service).GetResolveByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_ResolveByBinary() {
	fn := ax7ValueOf(new(Service).ResolveByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetResolveByBinary() {
	fn := ax7ValueOf(new(Service).GetResolveByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HashByString() {
	fn := ax7ValueOf(new(Service).HashByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetHashByString() {
	fn := ax7ValueOf(new(Service).GetHashByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HashByBinary() {
	fn := ax7ValueOf(new(Service).HashByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetHashByBinary() {
	fn := ax7ValueOf(new(Service).GetHashByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_ResolveString() {
	fn := ax7ValueOf(new(Service).ResolveString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetResolveString() {
	fn := ax7ValueOf(new(Service).GetResolveString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_ResolveBinary() {
	fn := ax7ValueOf(new(Service).ResolveBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetResolveBinary() {
	fn := ax7ValueOf(new(Service).GetResolveBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_ResolveName() {
	fn := ax7ValueOf(new(Service).ResolveName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetResolveName() {
	fn := ax7ValueOf(new(Service).GetResolveName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HashName() {
	fn := ax7ValueOf(new(Service).HashName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetHashName() {
	fn := ax7ValueOf(new(Service).GetHashName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_VerifyName() {
	fn := ax7ValueOf(new(Service).VerifyName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyName() {
	fn := ax7ValueOf(new(Service).GetVerifyName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ResolveByName() {
	fn := ax7ValueOf(new(Service).ResolveByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetResolveByName() {
	fn := ax7ValueOf(new(Service).GetResolveByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HashByName() {
	fn := ax7ValueOf(new(Service).HashByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetHashByName() {
	fn := ax7ValueOf(new(Service).GetHashByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_VerifyByName() {
	fn := ax7ValueOf(new(Service).VerifyByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyByName() {
	fn := ax7ValueOf(new(Service).GetVerifyByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyString() {
	fn := ax7ValueOf(new(Service).VerifyString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyString() {
	fn := ax7ValueOf(new(Service).GetVerifyString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyBinary() {
	fn := ax7ValueOf(new(Service).VerifyBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyBinary() {
	fn := ax7ValueOf(new(Service).GetVerifyBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyByString() {
	fn := ax7ValueOf(new(Service).VerifyByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyByString() {
	fn := ax7ValueOf(new(Service).GetVerifyByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyByBinary() {
	fn := ax7ValueOf(new(Service).VerifyByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyByBinary() {
	fn := ax7ValueOf(new(Service).GetVerifyByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_Blind() {
	fn := ax7ValueOf(new(Service).Blind)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetBlind() {
	fn := ax7ValueOf(new(Service).GetBlind)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_TypeName() {
	fn := ax7ValueOf(new(Service).TypeName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypeName() {
	fn := ax7ValueOf(new(Service).GetTypeName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_IsName() {
	fn := ax7ValueOf(new(Service).IsName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetIsName() {
	fn := ax7ValueOf(new(Service).GetIsName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_IsKnown() {
	fn := ax7ValueOf(new(Service).IsKnown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetIsKnown() {
	fn := ax7ValueOf(new(Service).GetIsKnown)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_IsLinked() {
	fn := ax7ValueOf(new(Service).IsLinked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetIsLinked() {
	fn := ax7ValueOf(new(Service).GetIsLinked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetRollout() {
	fn := ax7ValueOf(new(Service).GetRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasRollout() {
	fn := ax7ValueOf(new(Service).HasRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasRollout() {
	fn := ax7ValueOf(new(Service).GetHasRollout)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_IsReserved() {
	fn := ax7ValueOf(new(Service).IsReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetIsReserved() {
	fn := ax7ValueOf(new(Service).GetIsReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_IsLockedUp() {
	fn := ax7ValueOf(new(Service).IsLockedUp)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetIsLockedUp() {
	fn := ax7ValueOf(new(Service).GetIsLockedUp)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GrindName() {
	fn := ax7ValueOf(new(Service).GrindName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetGrindName() {
	fn := ax7ValueOf(new(Service).GetGrindName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_CountOpens() {
	fn := ax7ValueOf(new(Service).CountOpens)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetCountOpens() {
	fn := ax7ValueOf(new(Service).GetCountOpens)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_CountUpdates() {
	fn := ax7ValueOf(new(Service).CountUpdates)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetCountUpdates() {
	fn := ax7ValueOf(new(Service).GetCountUpdates)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_CountRenewals() {
	fn := ax7ValueOf(new(Service).CountRenewals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetCountRenewals() {
	fn := ax7ValueOf(new(Service).GetCountRenewals)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasNames() {
	fn := ax7ValueOf(new(Service).HasNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasNames() {
	fn := ax7ValueOf(new(Service).GetHasNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_AddNames() {
	fn := ax7ValueOf(new(Service).AddNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleService_GetAddNames() {
	fn := ax7ValueOf(new(Service).GetAddNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleService_RemoveNames() {
	fn := ax7ValueOf(new(Service).RemoveNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleService_GetRemoveNames() {
	fn := ax7ValueOf(new(Service).GetRemoveNames)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 0
}

func ExampleService_HasSaneCovenants() {
	fn := ax7ValueOf(new(Service).HasSaneCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasSaneCovenants() {
	fn := ax7ValueOf(new(Service).GetHasSaneCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerifyCovenants() {
	fn := ax7ValueOf(new(Service).VerifyCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerifyCovenants() {
	fn := ax7ValueOf(new(Service).GetVerifyCovenants)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_Types() {
	fn := ax7ValueOf(new(Service).Types)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypes() {
	fn := ax7ValueOf(new(Service).GetTypes)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TypesByVal() {
	fn := ax7ValueOf(new(Service).TypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTypesByVal() {
	fn := ax7ValueOf(new(Service).GetTypesByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_Blacklist() {
	fn := ax7ValueOf(new(Service).Blacklist)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetBlacklist() {
	fn := ax7ValueOf(new(Service).GetBlacklist)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerificationFlags() {
	fn := ax7ValueOf(new(Service).VerificationFlags)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerificationFlags() {
	fn := ax7ValueOf(new(Service).GetVerificationFlags)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_VerificationFlagsByVal() {
	fn := ax7ValueOf(new(Service).VerificationFlagsByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetVerificationFlagsByVal() {
	fn := ax7ValueOf(new(Service).GetVerificationFlagsByVal)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReserved() {
	fn := ax7ValueOf(new(Service).GetReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetReservedName() {
	fn := ax7ValueOf(new(Service).GetReservedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetReservedString() {
	fn := ax7ValueOf(new(Service).GetReservedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetReservedBinary() {
	fn := ax7ValueOf(new(Service).GetReservedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetReservedByString() {
	fn := ax7ValueOf(new(Service).GetReservedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetReservedByBinary() {
	fn := ax7ValueOf(new(Service).GetReservedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasReserved() {
	fn := ax7ValueOf(new(Service).HasReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasReservedString() {
	fn := ax7ValueOf(new(Service).HasReservedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasReservedBinary() {
	fn := ax7ValueOf(new(Service).HasReservedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasReservedByString() {
	fn := ax7ValueOf(new(Service).HasReservedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasReservedByBinary() {
	fn := ax7ValueOf(new(Service).HasReservedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedByName() {
	fn := ax7ValueOf(new(Service).GetReservedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasReservedByName() {
	fn := ax7ValueOf(new(Service).HasReservedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasReservedName() {
	fn := ax7ValueOf(new(Service).HasReservedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedHash() {
	fn := ax7ValueOf(new(Service).GetReservedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasReservedHash() {
	fn := ax7ValueOf(new(Service).HasReservedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedByHash() {
	fn := ax7ValueOf(new(Service).GetReservedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasReservedByHash() {
	fn := ax7ValueOf(new(Service).HasReservedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLocked() {
	fn := ax7ValueOf(new(Service).GetLocked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetLockedName() {
	fn := ax7ValueOf(new(Service).GetLockedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetLockedString() {
	fn := ax7ValueOf(new(Service).GetLockedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetLockedBinary() {
	fn := ax7ValueOf(new(Service).GetLockedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetLockedByString() {
	fn := ax7ValueOf(new(Service).GetLockedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_GetLockedByBinary() {
	fn := ax7ValueOf(new(Service).GetLockedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasLocked() {
	fn := ax7ValueOf(new(Service).HasLocked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasLockedString() {
	fn := ax7ValueOf(new(Service).HasLockedString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasLockedBinary() {
	fn := ax7ValueOf(new(Service).HasLockedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasLockedByString() {
	fn := ax7ValueOf(new(Service).HasLockedByString)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasLockedByBinary() {
	fn := ax7ValueOf(new(Service).HasLockedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedByName() {
	fn := ax7ValueOf(new(Service).GetLockedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasLockedByName() {
	fn := ax7ValueOf(new(Service).HasLockedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasLockedName() {
	fn := ax7ValueOf(new(Service).HasLockedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedHash() {
	fn := ax7ValueOf(new(Service).GetLockedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasLockedHash() {
	fn := ax7ValueOf(new(Service).HasLockedHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedByHash() {
	fn := ax7ValueOf(new(Service).GetLockedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleService_HasLockedByHash() {
	fn := ax7ValueOf(new(Service).HasLockedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedCatalog() {
	fn := ax7ValueOf(new(Service).ReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_Core() {
	fn := ax7ValueOf(new(Service).Core)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetCore() {
	fn := ax7ValueOf(new(Service).GetCore)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedCatalog() {
	fn := ax7ValueOf(new(Service).GetReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_DefaultReservedCatalog() {
	fn := ax7ValueOf(new(Service).DefaultReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetDefaultReservedCatalog() {
	fn := ax7ValueOf(new(Service).GetDefaultReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasReservedCatalog() {
	fn := ax7ValueOf(new(Service).HasReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasReservedCatalog() {
	fn := ax7ValueOf(new(Service).GetHasReservedCatalog)
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

func ExampleService_ReservedSize() {
	fn := ax7ValueOf(new(Service).ReservedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedSize() {
	fn := ax7ValueOf(new(Service).GetReservedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedNameValue() {
	fn := ax7ValueOf(new(Service).ReservedNameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedNameValue() {
	fn := ax7ValueOf(new(Service).GetReservedNameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedRootValue() {
	fn := ax7ValueOf(new(Service).ReservedRootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedRootValue() {
	fn := ax7ValueOf(new(Service).GetReservedRootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedTopValue() {
	fn := ax7ValueOf(new(Service).ReservedTopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedTopValue() {
	fn := ax7ValueOf(new(Service).GetReservedTopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_NameValue() {
	fn := ax7ValueOf(new(Service).NameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetNameValue() {
	fn := ax7ValueOf(new(Service).GetNameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_RootValue() {
	fn := ax7ValueOf(new(Service).RootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetRootValue() {
	fn := ax7ValueOf(new(Service).GetRootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_TopValue() {
	fn := ax7ValueOf(new(Service).TopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetTopValue() {
	fn := ax7ValueOf(new(Service).GetTopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_PrefixSize() {
	fn := ax7ValueOf(new(Service).PrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetPrefixSize() {
	fn := ax7ValueOf(new(Service).GetPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedEntries() {
	fn := ax7ValueOf(new(Service).ReservedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedEntries() {
	fn := ax7ValueOf(new(Service).GetReservedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedKeys() {
	fn := ax7ValueOf(new(Service).ReservedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedKeys() {
	fn := ax7ValueOf(new(Service).GetReservedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedValues() {
	fn := ax7ValueOf(new(Service).ReservedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedValues() {
	fn := ax7ValueOf(new(Service).GetReservedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_LockedCatalog() {
	fn := ax7ValueOf(new(Service).LockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedCatalog() {
	fn := ax7ValueOf(new(Service).GetLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_HasLockedCatalog() {
	fn := ax7ValueOf(new(Service).HasLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasLockedCatalog() {
	fn := ax7ValueOf(new(Service).GetHasLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_DefaultLockedCatalog() {
	fn := ax7ValueOf(new(Service).DefaultLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetDefaultLockedCatalog() {
	fn := ax7ValueOf(new(Service).GetDefaultLockedCatalog)
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

func ExampleHasReservedCatalog() {
	fn := ax7ValueOf(HasReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReservedCatalog() {
	fn := ax7ValueOf(GetHasReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleHasLockedCatalog() {
	fn := ax7ValueOf(HasLockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLockedCatalog() {
	fn := ax7ValueOf(GetHasLockedCatalog)
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

func ExampleReservedCatalog() {
	fn := ax7ValueOf(ReservedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedCatalog() {
	fn := ax7ValueOf(LockedCatalog)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedSize() {
	fn := ax7ValueOf(ReservedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedSize() {
	fn := ax7ValueOf(GetReservedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedNameValue() {
	fn := ax7ValueOf(ReservedNameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedNameValue() {
	fn := ax7ValueOf(GetReservedNameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedRootValue() {
	fn := ax7ValueOf(ReservedRootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedRootValue() {
	fn := ax7ValueOf(GetReservedRootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedTopValue() {
	fn := ax7ValueOf(ReservedTopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedTopValue() {
	fn := ax7ValueOf(GetReservedTopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleNameValue() {
	fn := ax7ValueOf(NameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetNameValue() {
	fn := ax7ValueOf(GetNameValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleRootValue() {
	fn := ax7ValueOf(RootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetRootValue() {
	fn := ax7ValueOf(GetRootValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleTopValue() {
	fn := ax7ValueOf(TopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetTopValue() {
	fn := ax7ValueOf(GetTopValue)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExamplePrefixSize() {
	fn := ax7ValueOf(PrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetPrefixSize() {
	fn := ax7ValueOf(GetPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedPrefixSize() {
	fn := ax7ValueOf(ReservedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedPrefixSize() {
	fn := ax7ValueOf(GetReservedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedEntries() {
	fn := ax7ValueOf(ReservedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedEntries() {
	fn := ax7ValueOf(GetReservedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedKeys() {
	fn := ax7ValueOf(ReservedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedKeys() {
	fn := ax7ValueOf(GetReservedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleReservedValues() {
	fn := ax7ValueOf(ReservedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedValues() {
	fn := ax7ValueOf(GetReservedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedSize() {
	fn := ax7ValueOf(LockedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedSize() {
	fn := ax7ValueOf(GetLockedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedPrefixSize() {
	fn := ax7ValueOf(LockedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedPrefixSize() {
	fn := ax7ValueOf(GetLockedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedEntries() {
	fn := ax7ValueOf(LockedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedEntries() {
	fn := ax7ValueOf(GetLockedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedKeys() {
	fn := ax7ValueOf(LockedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedKeys() {
	fn := ax7ValueOf(GetLockedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleLockedValues() {
	fn := ax7ValueOf(LockedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedValues() {
	fn := ax7ValueOf(GetLockedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReserved() {
	fn := ax7ValueOf(GetReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetReservedName() {
	fn := ax7ValueOf(GetReservedName)
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

func ExampleHasReserved() {
	fn := ax7ValueOf(HasReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasReserved() {
	fn := ax7ValueOf(GetHasReserved)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasReserved() {
	fn := ax7ValueOf(new(Service).GetHasReserved)
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

func ExampleService_GetHasReservedString() {
	fn := ax7ValueOf(new(Service).GetHasReservedString)
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

func ExampleService_GetHasReservedBinary() {
	fn := ax7ValueOf(new(Service).GetHasReservedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
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

func ExampleService_GetHasReservedByString() {
	fn := ax7ValueOf(new(Service).GetHasReservedByString)
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

func ExampleService_GetHasReservedByBinary() {
	fn := ax7ValueOf(new(Service).GetHasReservedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetReservedByName() {
	fn := ax7ValueOf(GetReservedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleHasReservedByName() {
	fn := ax7ValueOf(HasReservedByName)
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

func ExampleService_GetHasReservedByName() {
	fn := ax7ValueOf(new(Service).GetHasReservedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

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

func ExampleService_GetHasReservedName() {
	fn := ax7ValueOf(new(Service).GetHasReservedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
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

func ExampleService_GetHasReservedHash() {
	fn := ax7ValueOf(new(Service).GetHasReservedHash)
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

func ExampleService_GetHasReservedByHash() {
	fn := ax7ValueOf(new(Service).GetHasReservedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLocked() {
	fn := ax7ValueOf(GetLocked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
}

func ExampleGetLockedName() {
	fn := ax7ValueOf(GetLockedName)
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

func ExampleHasLocked() {
	fn := ax7ValueOf(HasLocked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetHasLocked() {
	fn := ax7ValueOf(GetHasLocked)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetHasLocked() {
	fn := ax7ValueOf(new(Service).GetHasLocked)
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

func ExampleService_GetHasLockedString() {
	fn := ax7ValueOf(new(Service).GetHasLockedString)
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

func ExampleService_GetHasLockedBinary() {
	fn := ax7ValueOf(new(Service).GetHasLockedBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
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

func ExampleService_GetHasLockedByString() {
	fn := ax7ValueOf(new(Service).GetHasLockedByString)
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

func ExampleService_GetHasLockedByBinary() {
	fn := ax7ValueOf(new(Service).GetHasLockedByBinary)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleGetLockedByName() {
	fn := ax7ValueOf(GetLockedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 2
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

func ExampleService_GetHasLockedByName() {
	fn := ax7ValueOf(new(Service).GetHasLockedByName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

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

func ExampleService_GetHasLockedName() {
	fn := ax7ValueOf(new(Service).GetHasLockedName)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
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

func ExampleService_GetHasLockedHash() {
	fn := ax7ValueOf(new(Service).GetHasLockedHash)
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

func ExampleService_GetHasLockedByHash() {
	fn := ax7ValueOf(new(Service).GetHasLockedByHash)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_LockedSize() {
	fn := ax7ValueOf(new(Service).LockedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedSize() {
	fn := ax7ValueOf(new(Service).GetLockedSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_LockedPrefixSize() {
	fn := ax7ValueOf(new(Service).LockedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedPrefixSize() {
	fn := ax7ValueOf(new(Service).GetLockedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_ReservedPrefixSize() {
	fn := ax7ValueOf(new(Service).ReservedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetReservedPrefixSize() {
	fn := ax7ValueOf(new(Service).GetReservedPrefixSize)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_LockedEntries() {
	fn := ax7ValueOf(new(Service).LockedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedEntries() {
	fn := ax7ValueOf(new(Service).GetLockedEntries)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_LockedKeys() {
	fn := ax7ValueOf(new(Service).LockedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedKeys() {
	fn := ax7ValueOf(new(Service).GetLockedKeys)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_LockedValues() {
	fn := ax7ValueOf(new(Service).LockedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleService_GetLockedValues() {
	fn := ax7ValueOf(new(Service).GetLockedValues)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}
