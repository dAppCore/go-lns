// SPDX-License-Identifier: EUPL-1.2

package lns

import core "dappco.re/go"

func ExampleWithCore() {
	core.Println(coverageExample(WithCore(core.New())))
	// Output: 1
}

func ExampleGetWithCore() {
	core.Println(coverageExample(GetWithCore(core.New())))
	// Output: 1
}

func ExampleWithReservedCatalog() {
	core.Println(coverageExample(WithReservedCatalog(nil)))
	// Output: 1
}

func ExampleGetWithReservedCatalog() {
	core.Println(coverageExample(GetWithReservedCatalog(nil)))
	// Output: 1
}

func ExampleWithLockedCatalog() {
	core.Println(coverageExample(WithLockedCatalog(nil)))
	// Output: 1
}

func ExampleGetWithLockedCatalog() {
	core.Println(coverageExample(GetWithLockedCatalog(nil)))
	// Output: 1
}

func ExampleNewServiceWithOptions() {
	core.Println(coverageExample(NewServiceWithOptions()))
	// Output: 1
}

func ExampleGetNewServiceWithOptions() {
	core.Println(coverageExample(GetNewServiceWithOptions()))
	// Output: 1
}

func ExampleNewService() {
	core.Println(coverageExample(NewService(core.New())))
	// Output: 1
}

func ExampleGetNewService() {
	core.Println(coverageExample(GetNewService(core.New())))
	// Output: 1
}

func ExampleRegister() {
	core.Println(coverageExample(Register(core.New())))
	// Output: 1
}

func ExampleGetRegister() {
	core.Println(coverageExample(GetRegister(core.New())))
	// Output: 1
}

func ExampleRegisterWithOptions() {
	core.Println(coverageExample(RegisterWithOptions()))
	// Output: 1
}

func ExampleGetRegisterWithOptions() {
	core.Println(coverageExample(GetRegisterWithOptions()))
	// Output: 1
}

func ExampleGetServiceName() {
	core.Println(coverageExample(GetServiceName()))
	// Output: 1
}

func ExampleGetDefaultTTL() {
	core.Println(coverageExample(GetDefaultTTL()))
	// Output: 1
}

func ExampleGetDummy() {
	core.Println(coverageExample(GetDummy()))
	// Output: 1
}

func ExampleGetTypeMapRoot() {
	core.Println(coverageExample(GetTypeMapRoot()))
	// Output: 1
}

func ExampleGetTypeMapEmpty() {
	core.Println(coverageExample(GetTypeMapEmpty()))
	// Output: 1
}

func ExampleGetTypeMapNS() {
	core.Println(coverageExample(GetTypeMapNS()))
	// Output: 1
}

func ExampleGetTypeMapTXT() {
	core.Println(coverageExample(GetTypeMapTXT()))
	// Output: 1
}

func ExampleGetTypeMapA() {
	core.Println(coverageExample(GetTypeMapA()))
	// Output: 1
}

func ExampleGetTypeMapAAAA() {
	core.Println(coverageExample(GetTypeMapAAAA()))
	// Output: 1
}

func ExampleGetHSTypes() {
	core.Println(coverageExample(GetHSTypes()))
	// Output: 1
}

func ExampleGetHSTypesByVal() {
	core.Println(coverageExample(GetHSTypesByVal()))
	// Output: 1
}

func ExampleNewResource() {
	core.Println(coverageExample(NewResource()))
	// Output: 1
}

func ExampleGetNewResource() {
	core.Println(coverageExample(GetNewResource()))
	// Output: 1
}

func ExampleDecodeResource() {
	core.Println(coverageExample(DecodeResource([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetDecodeResource() {
	core.Println(coverageExample(GetDecodeResource([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleCreate() {
	core.Println(coverageExample(Create("lethean.lthn", "lethean.lthn", []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetCreate() {
	core.Println(coverageExample(GetCreate("lethean.lthn", "lethean.lthn", []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleNextName() {
	core.Println(coverageExample(NextName("lethean.lthn")))
	// Output: 1
}

func ExampleGetNextName() {
	core.Println(coverageExample(GetNextName("lethean.lthn")))
	// Output: 1
}

func ExamplePrevName() {
	core.Println(coverageExample(PrevName("lethean.lthn")))
	// Output: 1
}

func ExampleGetPrevName() {
	core.Println(coverageExample(GetPrevName("lethean.lthn")))
	// Output: 1
}

func ExampleService_NewClaim() {
	core.Println(coverageExample(new(Service).NewClaim()))
	// Output: 1
}

func ExampleService_GetNewClaim() {
	core.Println(coverageExample(new(Service).GetNewClaim()))
	// Output: 1
}

func ExampleService_NewOutpoint() {
	core.Println(coverageExample(new(Service).NewOutpoint()))
	// Output: 1
}

func ExampleService_GetNewOutpoint() {
	core.Println(coverageExample(new(Service).GetNewOutpoint()))
	// Output: 1
}

func ExampleService_NewNameView() {
	core.Println(coverageExample(new(Service).NewNameView()))
	// Output: 1
}

func ExampleService_GetNewNameView() {
	core.Println(coverageExample(new(Service).GetNewNameView()))
	// Output: 1
}

func ExampleService_Create() {
	core.Println(coverageExample(new(Service).Create("lethean.lthn", "lethean.lthn", []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_GetCreate() {
	core.Println(coverageExample(new(Service).GetCreate("lethean.lthn", "lethean.lthn", []byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_DefaultTTL() {
	core.Println(coverageExample(new(Service).DefaultTTL()))
	// Output: 1
}

func ExampleService_GetDefaultTTL() {
	core.Println(coverageExample(new(Service).GetDefaultTTL()))
	// Output: 1
}

func ExampleService_Dummy() {
	core.Println(coverageExample(new(Service).Dummy()))
	// Output: 1
}

func ExampleService_GetDummy() {
	core.Println(coverageExample(new(Service).GetDummy()))
	// Output: 1
}

func ExampleService_TypeMapRoot() {
	core.Println(coverageExample(new(Service).TypeMapRoot()))
	// Output: 1
}

func ExampleService_GetTypeMapRoot() {
	core.Println(coverageExample(new(Service).GetTypeMapRoot()))
	// Output: 1
}

func ExampleService_TypeMapEmpty() {
	core.Println(coverageExample(new(Service).TypeMapEmpty()))
	// Output: 1
}

func ExampleService_GetTypeMapEmpty() {
	core.Println(coverageExample(new(Service).GetTypeMapEmpty()))
	// Output: 1
}

func ExampleService_TypeMapNS() {
	core.Println(coverageExample(new(Service).TypeMapNS()))
	// Output: 1
}

func ExampleService_GetTypeMapNS() {
	core.Println(coverageExample(new(Service).GetTypeMapNS()))
	// Output: 1
}

func ExampleService_TypeMapTXT() {
	core.Println(coverageExample(new(Service).TypeMapTXT()))
	// Output: 1
}

func ExampleService_GetTypeMapTXT() {
	core.Println(coverageExample(new(Service).GetTypeMapTXT()))
	// Output: 1
}

func ExampleService_TypeMapA() {
	core.Println(coverageExample(new(Service).TypeMapA()))
	// Output: 1
}

func ExampleService_GetTypeMapA() {
	core.Println(coverageExample(new(Service).GetTypeMapA()))
	// Output: 1
}

func ExampleService_TypeMapAAAA() {
	core.Println(coverageExample(new(Service).TypeMapAAAA()))
	// Output: 1
}

func ExampleService_GetTypeMapAAAA() {
	core.Println(coverageExample(new(Service).GetTypeMapAAAA()))
	// Output: 1
}

func ExampleService_HSTypes() {
	core.Println(coverageExample(new(Service).HSTypes()))
	// Output: 1
}

func ExampleService_GetHSTypes() {
	core.Println(coverageExample(new(Service).GetHSTypes()))
	// Output: 1
}

func ExampleService_HSTypesByVal() {
	core.Println(coverageExample(new(Service).HSTypesByVal()))
	// Output: 1
}

func ExampleService_GetHSTypesByVal() {
	core.Println(coverageExample(new(Service).GetHSTypesByVal()))
	// Output: 1
}

func ExampleService_NextName() {
	core.Println(coverageExample(new(Service).NextName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetNextName() {
	core.Println(coverageExample(new(Service).GetNextName("lethean.lthn")))
	// Output: 1
}

func ExampleService_PrevName() {
	core.Println(coverageExample(new(Service).PrevName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetPrevName() {
	core.Println(coverageExample(new(Service).GetPrevName("lethean.lthn")))
	// Output: 1
}

func ExampleService_NewResource() {
	core.Println(coverageExample(new(Service).NewResource()))
	// Output: 1
}

func ExampleService_GetNewResource() {
	core.Println(coverageExample(new(Service).GetNewResource()))
	// Output: 1
}

func ExampleService_DecodeResource() {
	core.Println(coverageExample(new(Service).DecodeResource([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_GetDecodeResource() {
	core.Println(coverageExample(new(Service).GetDecodeResource([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_OnStartup() {
	core.Println(coverageExample(new(Service).OnStartup(core.Background())))
	// Output: 1
}

func ExampleService_GetOnStartup() {
	core.Println(coverageExample(new(Service).GetOnStartup(core.Background())))
	// Output: 1
}

func ExampleService_OnShutdown() {
	core.Println(coverageExample(new(Service).OnShutdown(core.Background())))
	// Output: 1
}

func ExampleService_GetOnShutdown() {
	core.Println(coverageExample(new(Service).GetOnShutdown(core.Background())))
	// Output: 1
}

func ExampleService_HandleIPCEvents() {
	core.Println(coverageExample(new(Service).HandleIPCEvents(core.New(), "lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHandleIPCEvents() {
	core.Println(coverageExample(new(Service).GetHandleIPCEvents(core.New(), "lethean.lthn")))
	// Output: 1
}

func ExampleLookupCatalogName() {
	core.Println(coverageExample(LookupCatalogName[any]("lethean.lthn", nil, nil)))
	// Output: 2
}

func ExampleGetLookupCatalogName() {
	core.Println(coverageExample(GetLookupCatalogName[any]("lethean.lthn", nil, nil)))
	// Output: 2
}

func ExampleHasLookupCatalogName() {
	core.Println(coverageExample(HasLookupCatalogName[any]("lethean.lthn", nil, nil)))
	// Output: 1
}

func ExampleGetHasLookupCatalogName() {
	core.Println(coverageExample(GetHasLookupCatalogName[any]("lethean.lthn", nil, nil)))
	// Output: 1
}

func ExampleService_LookupCatalogName() {
	core.Println(coverageExample(new(Service).LookupCatalogName("lethean.lthn", nil, nil)))
	// Output: 2
}

func ExampleService_GetLookupCatalogName() {
	core.Println(coverageExample(new(Service).GetLookupCatalogName("lethean.lthn", nil, nil)))
	// Output: 2
}

func ExampleService_HasLookupCatalogName() {
	core.Println(coverageExample(new(Service).HasLookupCatalogName("lethean.lthn", nil, nil)))
	// Output: 1
}

func ExampleService_GetHasLookupCatalogName() {
	core.Println(coverageExample(new(Service).GetHasLookupCatalogName("lethean.lthn", nil, nil)))
	// Output: 1
}

func ExampleResolve() {
	core.Println(coverageExample(Resolve("lethean.lthn")))
	// Output: 2
}

func ExampleGetResolve() {
	core.Println(coverageExample(GetResolve("lethean.lthn")))
	// Output: 2
}

func ExampleHash() {
	core.Println(coverageExample(Hash("lethean.lthn")))
	// Output: 2
}

func ExampleGetHash() {
	core.Println(coverageExample(GetHash("lethean.lthn")))
	// Output: 2
}

func ExampleVerify() {
	core.Println(coverageExample(Verify("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerify() {
	core.Println(coverageExample(GetVerify("lethean.lthn")))
	// Output: 1
}

func ExampleHashString() {
	core.Println(coverageExample(HashString("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashString() {
	core.Println(coverageExample(GetHashString("lethean.lthn")))
	// Output: 2
}

func ExampleHashBinary() {
	core.Println(coverageExample(HashBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetHashBinary() {
	core.Println(coverageExample(GetHashBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleResolveByString() {
	core.Println(coverageExample(ResolveByString("lethean.lthn")))
	// Output: 2
}

func ExampleGetResolveByString() {
	core.Println(coverageExample(GetResolveByString("lethean.lthn")))
	// Output: 2
}

func ExampleResolveByBinary() {
	core.Println(coverageExample(ResolveByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetResolveByBinary() {
	core.Println(coverageExample(GetResolveByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleHashByString() {
	core.Println(coverageExample(HashByString("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashByString() {
	core.Println(coverageExample(GetHashByString("lethean.lthn")))
	// Output: 2
}

func ExampleHashByBinary() {
	core.Println(coverageExample(HashByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetHashByBinary() {
	core.Println(coverageExample(GetHashByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleResolveString() {
	core.Println(coverageExample(ResolveString("lethean.lthn")))
	// Output: 2
}

func ExampleGetResolveString() {
	core.Println(coverageExample(GetResolveString("lethean.lthn")))
	// Output: 2
}

func ExampleResolveBinary() {
	core.Println(coverageExample(ResolveBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleGetResolveBinary() {
	core.Println(coverageExample(GetResolveBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleResolveName() {
	core.Println(coverageExample(ResolveName("lethean.lthn")))
	// Output: 2
}

func ExampleGetResolveName() {
	core.Println(coverageExample(GetResolveName("lethean.lthn")))
	// Output: 2
}

func ExampleHashName() {
	core.Println(coverageExample(HashName("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashName() {
	core.Println(coverageExample(GetHashName("lethean.lthn")))
	// Output: 2
}

func ExampleVerifyName() {
	core.Println(coverageExample(VerifyName("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyName() {
	core.Println(coverageExample(GetVerifyName("lethean.lthn")))
	// Output: 1
}

func ExampleResolveByName() {
	core.Println(coverageExample(ResolveByName("lethean.lthn")))
	// Output: 2
}

func ExampleGetResolveByName() {
	core.Println(coverageExample(GetResolveByName("lethean.lthn")))
	// Output: 2
}

func ExampleHashByName() {
	core.Println(coverageExample(HashByName("lethean.lthn")))
	// Output: 2
}

func ExampleGetHashByName() {
	core.Println(coverageExample(GetHashByName("lethean.lthn")))
	// Output: 2
}

func ExampleVerifyByName() {
	core.Println(coverageExample(VerifyByName("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyByName() {
	core.Println(coverageExample(GetVerifyByName("lethean.lthn")))
	// Output: 1
}

func ExampleVerifyString() {
	core.Println(coverageExample(VerifyString("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyString() {
	core.Println(coverageExample(GetVerifyString("lethean.lthn")))
	// Output: 1
}

func ExampleVerifyBinary() {
	core.Println(coverageExample(VerifyBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetVerifyBinary() {
	core.Println(coverageExample(GetVerifyBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleVerifyByString() {
	core.Println(coverageExample(VerifyByString("lethean.lthn")))
	// Output: 1
}

func ExampleGetVerifyByString() {
	core.Println(coverageExample(GetVerifyByString("lethean.lthn")))
	// Output: 1
}

func ExampleVerifyByBinary() {
	core.Println(coverageExample(VerifyByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetVerifyByBinary() {
	core.Println(coverageExample(GetVerifyByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleBlind() {
	core.Println(coverageExample(Blind(1, [32]byte{1})))
	// Output: 2
}

func ExampleGetBlind() {
	core.Println(coverageExample(GetBlind(1, [32]byte{1})))
	// Output: 2
}

func ExampleTypeName() {
	core.Println(coverageExample(TypeName(1)))
	// Output: 1
}

func ExampleGetBlacklist() {
	core.Println(coverageExample(GetBlacklist()))
	// Output: 1
}

func ExampleGetVerificationFlags() {
	core.Println(coverageExample(GetVerificationFlags()))
	// Output: 1
}

func ExampleGetVerificationFlagsByVal() {
	core.Println(coverageExample(GetVerificationFlagsByVal()))
	// Output: 1
}

func ExampleNewClaim() {
	core.Println(coverageExample(NewClaim()))
	// Output: 1
}

func ExampleGetNewClaim() {
	core.Println(coverageExample(GetNewClaim()))
	// Output: 1
}

func ExampleNewOutpoint() {
	core.Println(coverageExample(NewOutpoint()))
	// Output: 1
}

func ExampleGetNewOutpoint() {
	core.Println(coverageExample(GetNewOutpoint()))
	// Output: 1
}

func ExampleNewNameView() {
	core.Println(coverageExample(NewNameView()))
	// Output: 1
}

func ExampleGetNewNameView() {
	core.Println(coverageExample(GetNewNameView()))
	// Output: 1
}

func ExampleGetMaxNameSize() {
	core.Println(coverageExample(GetMaxNameSize()))
	// Output: 1
}

func ExampleGetMaxResourceSize() {
	core.Println(coverageExample(GetMaxResourceSize()))
	// Output: 1
}

func ExampleGetVerifyCovenantsNone() {
	core.Println(coverageExample(GetVerifyCovenantsNone()))
	// Output: 1
}

func ExampleGetVerifyCovenantsHardened() {
	core.Println(coverageExample(GetVerifyCovenantsHardened()))
	// Output: 1
}

func ExampleGetVerifyCovenantsLockup() {
	core.Println(coverageExample(GetVerifyCovenantsLockup()))
	// Output: 1
}

func ExampleGetMandatoryVerifyCovenantFlags() {
	core.Println(coverageExample(GetMandatoryVerifyCovenantFlags()))
	// Output: 1
}

func ExampleGetMaxCovenantSize() {
	core.Println(coverageExample(GetMaxCovenantSize()))
	// Output: 1
}

func ExampleGetCovenantMaxSize() {
	core.Println(coverageExample(GetCovenantMaxSize()))
	// Output: 1
}

func ExampleGetTypes() {
	core.Println(coverageExample(GetTypes()))
	// Output: 1
}

func ExampleGetTypesByVal() {
	core.Println(coverageExample(GetTypesByVal()))
	// Output: 1
}

func ExampleGetTypeName() {
	core.Println(coverageExample(GetTypeName(1)))
	// Output: 1
}

func ExampleIsName() {
	core.Println(coverageExample(IsName(1)))
	// Output: 1
}

func ExampleGetIsName() {
	core.Println(coverageExample(GetIsName(1)))
	// Output: 1
}

func ExampleIsKnown() {
	core.Println(coverageExample(IsKnown(1)))
	// Output: 1
}

func ExampleGetIsKnown() {
	core.Println(coverageExample(GetIsKnown(1)))
	// Output: 1
}

func ExampleIsLinked() {
	core.Println(coverageExample(IsLinked(1)))
	// Output: 1
}

func ExampleGetIsLinked() {
	core.Println(coverageExample(GetIsLinked(1)))
	// Output: 1
}

func ExampleService_MaxNameSize() {
	core.Println(coverageExample(new(Service).MaxNameSize()))
	// Output: 1
}

func ExampleService_GetMaxNameSize() {
	core.Println(coverageExample(new(Service).GetMaxNameSize()))
	// Output: 1
}

func ExampleService_MaxResourceSize() {
	core.Println(coverageExample(new(Service).MaxResourceSize()))
	// Output: 1
}

func ExampleService_GetMaxResourceSize() {
	core.Println(coverageExample(new(Service).GetMaxResourceSize()))
	// Output: 1
}

func ExampleService_VerifyCovenantsNone() {
	core.Println(coverageExample(new(Service).VerifyCovenantsNone()))
	// Output: 1
}

func ExampleService_GetVerifyCovenantsNone() {
	core.Println(coverageExample(new(Service).GetVerifyCovenantsNone()))
	// Output: 1
}

func ExampleService_VerifyCovenantsHardened() {
	core.Println(coverageExample(new(Service).VerifyCovenantsHardened()))
	// Output: 1
}

func ExampleService_GetVerifyCovenantsHardened() {
	core.Println(coverageExample(new(Service).GetVerifyCovenantsHardened()))
	// Output: 1
}

func ExampleService_VerifyCovenantsLockup() {
	core.Println(coverageExample(new(Service).VerifyCovenantsLockup()))
	// Output: 1
}

func ExampleService_GetVerifyCovenantsLockup() {
	core.Println(coverageExample(new(Service).GetVerifyCovenantsLockup()))
	// Output: 1
}

func ExampleService_MandatoryVerifyCovenantFlags() {
	core.Println(coverageExample(new(Service).MandatoryVerifyCovenantFlags()))
	// Output: 1
}

func ExampleService_GetMandatoryVerifyCovenantFlags() {
	core.Println(coverageExample(new(Service).GetMandatoryVerifyCovenantFlags()))
	// Output: 1
}

func ExampleService_MaxCovenantSize() {
	core.Println(coverageExample(new(Service).MaxCovenantSize()))
	// Output: 1
}

func ExampleService_GetMaxCovenantSize() {
	core.Println(coverageExample(new(Service).GetMaxCovenantSize()))
	// Output: 1
}

func ExampleService_CovenantMaxSize() {
	core.Println(coverageExample(new(Service).CovenantMaxSize()))
	// Output: 1
}

func ExampleService_GetCovenantMaxSize() {
	core.Println(coverageExample(new(Service).GetCovenantMaxSize()))
	// Output: 1
}

func ExampleGetRollout() {
	core.Println(coverageExample(GetRollout([32]byte{1}, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleHasRollout() {
	core.Println(coverageExample(HasRollout([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleGetHasRollout() {
	core.Println(coverageExample(GetHasRollout([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleIsReserved() {
	core.Println(coverageExample(IsReserved([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleGetIsReserved() {
	core.Println(coverageExample(GetIsReserved([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleIsLockedUp() {
	core.Println(coverageExample(IsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleGetIsLockedUp() {
	core.Println(coverageExample(GetIsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleGrindName() {
	core.Println(coverageExample(GrindName(1, 1, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleGetGrindName() {
	core.Println(coverageExample(GetGrindName(1, 1, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleCountOpens() {
	core.Println(coverageExample(CountOpens(coverageTransaction("good"))))
	// Output: 1
}

func ExampleGetCountOpens() {
	core.Println(coverageExample(GetCountOpens(coverageTransaction("good"))))
	// Output: 1
}

func ExampleCountUpdates() {
	core.Println(coverageExample(CountUpdates(coverageTransaction("good"))))
	// Output: 1
}

func ExampleGetCountUpdates() {
	core.Println(coverageExample(GetCountUpdates(coverageTransaction("good"))))
	// Output: 1
}

func ExampleCountRenewals() {
	core.Println(coverageExample(CountRenewals(coverageTransaction("good"))))
	// Output: 1
}

func ExampleGetCountRenewals() {
	core.Println(coverageExample(GetCountRenewals(coverageTransaction("good"))))
	// Output: 1
}

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

func ExampleVerifyCovenants() {
	core.Println(coverageExample(VerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good"))))
	// Output: 1
}

func ExampleGetVerifyCovenants() {
	core.Println(coverageExample(GetVerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good"))))
	// Output: 1
}

func ExampleService_ServiceName() {
	core.Println(coverageExample(new(Service).ServiceName()))
	// Output: 1
}

func ExampleService_GetServiceName() {
	core.Println(coverageExample(new(Service).GetServiceName()))
	// Output: 1
}

func ExampleService_Resolve() {
	core.Println(coverageExample(new(Service).Resolve("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetResolve() {
	core.Println(coverageExample(new(Service).GetResolve("lethean.lthn")))
	// Output: 2
}

func ExampleService_Hash() {
	core.Println(coverageExample(new(Service).Hash("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetHash() {
	core.Println(coverageExample(new(Service).GetHash("lethean.lthn")))
	// Output: 2
}

func ExampleService_Verify() {
	core.Println(coverageExample(new(Service).Verify("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetVerify() {
	core.Println(coverageExample(new(Service).GetVerify("lethean.lthn")))
	// Output: 1
}

func ExampleService_HashString() {
	core.Println(coverageExample(new(Service).HashString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetHashString() {
	core.Println(coverageExample(new(Service).GetHashString("lethean.lthn")))
	// Output: 2
}

func ExampleService_HashBinary() {
	core.Println(coverageExample(new(Service).HashBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_GetHashBinary() {
	core.Println(coverageExample(new(Service).GetHashBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_ResolveByString() {
	core.Println(coverageExample(new(Service).ResolveByString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetResolveByString() {
	core.Println(coverageExample(new(Service).GetResolveByString("lethean.lthn")))
	// Output: 2
}

func ExampleService_ResolveByBinary() {
	core.Println(coverageExample(new(Service).ResolveByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_GetResolveByBinary() {
	core.Println(coverageExample(new(Service).GetResolveByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_HashByString() {
	core.Println(coverageExample(new(Service).HashByString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetHashByString() {
	core.Println(coverageExample(new(Service).GetHashByString("lethean.lthn")))
	// Output: 2
}

func ExampleService_HashByBinary() {
	core.Println(coverageExample(new(Service).HashByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_GetHashByBinary() {
	core.Println(coverageExample(new(Service).GetHashByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_ResolveString() {
	core.Println(coverageExample(new(Service).ResolveString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetResolveString() {
	core.Println(coverageExample(new(Service).GetResolveString("lethean.lthn")))
	// Output: 2
}

func ExampleService_ResolveBinary() {
	core.Println(coverageExample(new(Service).ResolveBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_GetResolveBinary() {
	core.Println(coverageExample(new(Service).GetResolveBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_ResolveName() {
	core.Println(coverageExample(new(Service).ResolveName("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetResolveName() {
	core.Println(coverageExample(new(Service).GetResolveName("lethean.lthn")))
	// Output: 2
}

func ExampleService_HashName() {
	core.Println(coverageExample(new(Service).HashName("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetHashName() {
	core.Println(coverageExample(new(Service).GetHashName("lethean.lthn")))
	// Output: 2
}

func ExampleService_VerifyName() {
	core.Println(coverageExample(new(Service).VerifyName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetVerifyName() {
	core.Println(coverageExample(new(Service).GetVerifyName("lethean.lthn")))
	// Output: 1
}

func ExampleService_ResolveByName() {
	core.Println(coverageExample(new(Service).ResolveByName("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetResolveByName() {
	core.Println(coverageExample(new(Service).GetResolveByName("lethean.lthn")))
	// Output: 2
}

func ExampleService_HashByName() {
	core.Println(coverageExample(new(Service).HashByName("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetHashByName() {
	core.Println(coverageExample(new(Service).GetHashByName("lethean.lthn")))
	// Output: 2
}

func ExampleService_VerifyByName() {
	core.Println(coverageExample(new(Service).VerifyByName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetVerifyByName() {
	core.Println(coverageExample(new(Service).GetVerifyByName("lethean.lthn")))
	// Output: 1
}

func ExampleService_VerifyString() {
	core.Println(coverageExample(new(Service).VerifyString("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetVerifyString() {
	core.Println(coverageExample(new(Service).GetVerifyString("lethean.lthn")))
	// Output: 1
}

func ExampleService_VerifyBinary() {
	core.Println(coverageExample(new(Service).VerifyBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_GetVerifyBinary() {
	core.Println(coverageExample(new(Service).GetVerifyBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_VerifyByString() {
	core.Println(coverageExample(new(Service).VerifyByString("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetVerifyByString() {
	core.Println(coverageExample(new(Service).GetVerifyByString("lethean.lthn")))
	// Output: 1
}

func ExampleService_VerifyByBinary() {
	core.Println(coverageExample(new(Service).VerifyByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_GetVerifyByBinary() {
	core.Println(coverageExample(new(Service).GetVerifyByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_Blind() {
	core.Println(coverageExample(new(Service).Blind(1, [32]byte{1})))
	// Output: 2
}

func ExampleService_GetBlind() {
	core.Println(coverageExample(new(Service).GetBlind(1, [32]byte{1})))
	// Output: 2
}

func ExampleService_TypeName() {
	core.Println(coverageExample(new(Service).TypeName(1)))
	// Output: 1
}

func ExampleService_GetTypeName() {
	core.Println(coverageExample(new(Service).GetTypeName(1)))
	// Output: 1
}

func ExampleService_IsName() {
	core.Println(coverageExample(new(Service).IsName(1)))
	// Output: 1
}

func ExampleService_GetIsName() {
	core.Println(coverageExample(new(Service).GetIsName(1)))
	// Output: 1
}

func ExampleService_IsKnown() {
	core.Println(coverageExample(new(Service).IsKnown(1)))
	// Output: 1
}

func ExampleService_GetIsKnown() {
	core.Println(coverageExample(new(Service).GetIsKnown(1)))
	// Output: 1
}

func ExampleService_IsLinked() {
	core.Println(coverageExample(new(Service).IsLinked(1)))
	// Output: 1
}

func ExampleService_GetIsLinked() {
	core.Println(coverageExample(new(Service).GetIsLinked(1)))
	// Output: 1
}

func ExampleService_GetRollout() {
	core.Println(coverageExample(new(Service).GetRollout([32]byte{1}, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleService_HasRollout() {
	core.Println(coverageExample(new(Service).HasRollout([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleService_GetHasRollout() {
	core.Println(coverageExample(new(Service).GetHasRollout([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleService_IsReserved() {
	core.Println(coverageExample(new(Service).IsReserved([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleService_GetIsReserved() {
	core.Println(coverageExample(new(Service).GetIsReserved([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleService_IsLockedUp() {
	core.Println(coverageExample(new(Service).IsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleService_GetIsLockedUp() {
	core.Println(coverageExample(new(Service).GetIsLockedUp([32]byte{1}, 1, coverageValue[NameRules]("good"))))
	// Output: 1
}

func ExampleService_GrindName() {
	core.Println(coverageExample(new(Service).GrindName(1, 1, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleService_GetGrindName() {
	core.Println(coverageExample(new(Service).GetGrindName(1, 1, coverageValue[NameRules]("good"))))
	// Output: 2
}

func ExampleService_CountOpens() {
	core.Println(coverageExample(new(Service).CountOpens(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_GetCountOpens() {
	core.Println(coverageExample(new(Service).GetCountOpens(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_CountUpdates() {
	core.Println(coverageExample(new(Service).CountUpdates(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_GetCountUpdates() {
	core.Println(coverageExample(new(Service).GetCountUpdates(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_CountRenewals() {
	core.Println(coverageExample(new(Service).CountRenewals(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_GetCountRenewals() {
	core.Println(coverageExample(new(Service).GetCountRenewals(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_HasNames() {
	core.Println(coverageExample(new(Service).HasNames(coverageTransaction("good"), nil)))
	// Output: 1
}

func ExampleService_GetHasNames() {
	core.Println(coverageExample(new(Service).GetHasNames(coverageTransaction("good"), nil)))
	// Output: 1
}

func ExampleService_AddNames() {
	new(Service).AddNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleService_GetAddNames() {
	new(Service).GetAddNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleService_RemoveNames() {
	new(Service).RemoveNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleService_GetRemoveNames() {
	new(Service).GetRemoveNames(coverageTransaction("good"), nil)
	core.Println(coverageExample())
	// Output: 0
}

func ExampleService_HasSaneCovenants() {
	core.Println(coverageExample(new(Service).HasSaneCovenants(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_GetHasSaneCovenants() {
	core.Println(coverageExample(new(Service).GetHasSaneCovenants(coverageTransaction("good"))))
	// Output: 1
}

func ExampleService_VerifyCovenants() {
	core.Println(coverageExample(new(Service).VerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good"))))
	// Output: 1
}

func ExampleService_GetVerifyCovenants() {
	core.Println(coverageExample(new(Service).GetVerifyCovenants(coverageTransaction("good"), coverageValue[CoinView]("good"), 1, coverageValue[Network]("good"))))
	// Output: 1
}

func ExampleService_Types() {
	core.Println(coverageExample(new(Service).Types()))
	// Output: 1
}

func ExampleService_GetTypes() {
	core.Println(coverageExample(new(Service).GetTypes()))
	// Output: 1
}

func ExampleService_TypesByVal() {
	core.Println(coverageExample(new(Service).TypesByVal()))
	// Output: 1
}

func ExampleService_GetTypesByVal() {
	core.Println(coverageExample(new(Service).GetTypesByVal()))
	// Output: 1
}

func ExampleService_Blacklist() {
	core.Println(coverageExample(new(Service).Blacklist()))
	// Output: 1
}

func ExampleService_GetBlacklist() {
	core.Println(coverageExample(new(Service).GetBlacklist()))
	// Output: 1
}

func ExampleService_VerificationFlags() {
	core.Println(coverageExample(new(Service).VerificationFlags()))
	// Output: 1
}

func ExampleService_GetVerificationFlags() {
	core.Println(coverageExample(new(Service).GetVerificationFlags()))
	// Output: 1
}

func ExampleService_VerificationFlagsByVal() {
	core.Println(coverageExample(new(Service).VerificationFlagsByVal()))
	// Output: 1
}

func ExampleService_GetVerificationFlagsByVal() {
	core.Println(coverageExample(new(Service).GetVerificationFlagsByVal()))
	// Output: 1
}

func ExampleService_GetReserved() {
	core.Println(coverageExample(new(Service).GetReserved("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetReservedName() {
	core.Println(coverageExample(new(Service).GetReservedName("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetReservedString() {
	core.Println(coverageExample(new(Service).GetReservedString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetReservedBinary() {
	core.Println(coverageExample(new(Service).GetReservedBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_GetReservedByString() {
	core.Println(coverageExample(new(Service).GetReservedByString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetReservedByBinary() {
	core.Println(coverageExample(new(Service).GetReservedByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_HasReserved() {
	core.Println(coverageExample(new(Service).HasReserved("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasReservedString() {
	core.Println(coverageExample(new(Service).HasReservedString("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasReservedBinary() {
	core.Println(coverageExample(new(Service).HasReservedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_HasReservedByString() {
	core.Println(coverageExample(new(Service).HasReservedByString("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasReservedByBinary() {
	core.Println(coverageExample(new(Service).HasReservedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_GetReservedByName() {
	core.Println(coverageExample(new(Service).GetReservedByName("lethean.lthn")))
	// Output: 2
}

func ExampleService_HasReservedByName() {
	core.Println(coverageExample(new(Service).HasReservedByName("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasReservedName() {
	core.Println(coverageExample(new(Service).HasReservedName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetReservedHash() {
	core.Println(coverageExample(new(Service).GetReservedHash([32]byte{1})))
	// Output: 2
}

func ExampleService_HasReservedHash() {
	core.Println(coverageExample(new(Service).HasReservedHash([32]byte{1})))
	// Output: 1
}

func ExampleService_GetReservedByHash() {
	core.Println(coverageExample(new(Service).GetReservedByHash([32]byte{1})))
	// Output: 2
}

func ExampleService_HasReservedByHash() {
	core.Println(coverageExample(new(Service).HasReservedByHash([32]byte{1})))
	// Output: 1
}

func ExampleService_GetLocked() {
	core.Println(coverageExample(new(Service).GetLocked("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetLockedName() {
	core.Println(coverageExample(new(Service).GetLockedName("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetLockedString() {
	core.Println(coverageExample(new(Service).GetLockedString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetLockedBinary() {
	core.Println(coverageExample(new(Service).GetLockedBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_GetLockedByString() {
	core.Println(coverageExample(new(Service).GetLockedByString("lethean.lthn")))
	// Output: 2
}

func ExampleService_GetLockedByBinary() {
	core.Println(coverageExample(new(Service).GetLockedByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleService_HasLocked() {
	core.Println(coverageExample(new(Service).HasLocked("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasLockedString() {
	core.Println(coverageExample(new(Service).HasLockedString("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasLockedBinary() {
	core.Println(coverageExample(new(Service).HasLockedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_HasLockedByString() {
	core.Println(coverageExample(new(Service).HasLockedByString("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasLockedByBinary() {
	core.Println(coverageExample(new(Service).HasLockedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleService_GetLockedByName() {
	core.Println(coverageExample(new(Service).GetLockedByName("lethean.lthn")))
	// Output: 2
}

func ExampleService_HasLockedByName() {
	core.Println(coverageExample(new(Service).HasLockedByName("lethean.lthn")))
	// Output: 1
}

func ExampleService_HasLockedName() {
	core.Println(coverageExample(new(Service).HasLockedName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetLockedHash() {
	core.Println(coverageExample(new(Service).GetLockedHash([32]byte{1})))
	// Output: 2
}

func ExampleService_HasLockedHash() {
	core.Println(coverageExample(new(Service).HasLockedHash([32]byte{1})))
	// Output: 1
}

func ExampleService_GetLockedByHash() {
	core.Println(coverageExample(new(Service).GetLockedByHash([32]byte{1})))
	// Output: 2
}

func ExampleService_HasLockedByHash() {
	core.Println(coverageExample(new(Service).HasLockedByHash([32]byte{1})))
	// Output: 1
}

func ExampleService_ReservedCatalog() {
	core.Println(coverageExample(new(Service).ReservedCatalog()))
	// Output: 1
}

func ExampleService_Core() {
	core.Println(coverageExample(new(Service).Core()))
	// Output: 1
}

func ExampleService_GetCore() {
	core.Println(coverageExample(new(Service).GetCore()))
	// Output: 1
}

func ExampleService_GetReservedCatalog() {
	core.Println(coverageExample(new(Service).GetReservedCatalog()))
	// Output: 1
}

func ExampleService_DefaultReservedCatalog() {
	core.Println(coverageExample(new(Service).DefaultReservedCatalog()))
	// Output: 1
}

func ExampleService_GetDefaultReservedCatalog() {
	core.Println(coverageExample(new(Service).GetDefaultReservedCatalog()))
	// Output: 1
}

func ExampleService_HasReservedCatalog() {
	core.Println(coverageExample(new(Service).HasReservedCatalog()))
	// Output: 1
}

func ExampleService_GetHasReservedCatalog() {
	core.Println(coverageExample(new(Service).GetHasReservedCatalog()))
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

func ExampleService_ReservedSize() {
	core.Println(coverageExample(new(Service).ReservedSize()))
	// Output: 1
}

func ExampleService_GetReservedSize() {
	core.Println(coverageExample(new(Service).GetReservedSize()))
	// Output: 1
}

func ExampleService_ReservedNameValue() {
	core.Println(coverageExample(new(Service).ReservedNameValue()))
	// Output: 1
}

func ExampleService_GetReservedNameValue() {
	core.Println(coverageExample(new(Service).GetReservedNameValue()))
	// Output: 1
}

func ExampleService_ReservedRootValue() {
	core.Println(coverageExample(new(Service).ReservedRootValue()))
	// Output: 1
}

func ExampleService_GetReservedRootValue() {
	core.Println(coverageExample(new(Service).GetReservedRootValue()))
	// Output: 1
}

func ExampleService_ReservedTopValue() {
	core.Println(coverageExample(new(Service).ReservedTopValue()))
	// Output: 1
}

func ExampleService_GetReservedTopValue() {
	core.Println(coverageExample(new(Service).GetReservedTopValue()))
	// Output: 1
}

func ExampleService_NameValue() {
	core.Println(coverageExample(new(Service).NameValue()))
	// Output: 1
}

func ExampleService_GetNameValue() {
	core.Println(coverageExample(new(Service).GetNameValue()))
	// Output: 1
}

func ExampleService_RootValue() {
	core.Println(coverageExample(new(Service).RootValue()))
	// Output: 1
}

func ExampleService_GetRootValue() {
	core.Println(coverageExample(new(Service).GetRootValue()))
	// Output: 1
}

func ExampleService_TopValue() {
	core.Println(coverageExample(new(Service).TopValue()))
	// Output: 1
}

func ExampleService_GetTopValue() {
	core.Println(coverageExample(new(Service).GetTopValue()))
	// Output: 1
}

func ExampleService_PrefixSize() {
	core.Println(coverageExample(new(Service).PrefixSize()))
	// Output: 1
}

func ExampleService_GetPrefixSize() {
	core.Println(coverageExample(new(Service).GetPrefixSize()))
	// Output: 1
}

func ExampleService_ReservedEntries() {
	core.Println(coverageExample(new(Service).ReservedEntries()))
	// Output: 1
}

func ExampleService_GetReservedEntries() {
	core.Println(coverageExample(new(Service).GetReservedEntries()))
	// Output: 1
}

func ExampleService_ReservedKeys() {
	core.Println(coverageExample(new(Service).ReservedKeys()))
	// Output: 1
}

func ExampleService_GetReservedKeys() {
	core.Println(coverageExample(new(Service).GetReservedKeys()))
	// Output: 1
}

func ExampleService_ReservedValues() {
	core.Println(coverageExample(new(Service).ReservedValues()))
	// Output: 1
}

func ExampleService_GetReservedValues() {
	core.Println(coverageExample(new(Service).GetReservedValues()))
	// Output: 1
}

func ExampleService_LockedCatalog() {
	core.Println(coverageExample(new(Service).LockedCatalog()))
	// Output: 1
}

func ExampleService_GetLockedCatalog() {
	core.Println(coverageExample(new(Service).GetLockedCatalog()))
	// Output: 1
}

func ExampleService_HasLockedCatalog() {
	core.Println(coverageExample(new(Service).HasLockedCatalog()))
	// Output: 1
}

func ExampleService_GetHasLockedCatalog() {
	core.Println(coverageExample(new(Service).GetHasLockedCatalog()))
	// Output: 1
}

func ExampleService_DefaultLockedCatalog() {
	core.Println(coverageExample(new(Service).DefaultLockedCatalog()))
	// Output: 1
}

func ExampleService_GetDefaultLockedCatalog() {
	core.Println(coverageExample(new(Service).GetDefaultLockedCatalog()))
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

func ExampleHasReservedCatalog() {
	core.Println(coverageExample(HasReservedCatalog()))
	// Output: 1
}

func ExampleGetHasReservedCatalog() {
	core.Println(coverageExample(GetHasReservedCatalog()))
	// Output: 1
}

func ExampleHasLockedCatalog() {
	core.Println(coverageExample(HasLockedCatalog()))
	// Output: 1
}

func ExampleGetHasLockedCatalog() {
	core.Println(coverageExample(GetHasLockedCatalog()))
	// Output: 1
}

func ExampleGetDefaultLockedCatalog() {
	core.Println(coverageExample(GetDefaultLockedCatalog()))
	// Output: 1
}

func ExampleReservedCatalog() {
	core.Println(coverageExample(ReservedCatalog()))
	// Output: 1
}

func ExampleLockedCatalog() {
	core.Println(coverageExample(LockedCatalog()))
	// Output: 1
}

func ExampleReservedSize() {
	core.Println(coverageExample(ReservedSize()))
	// Output: 1
}

func ExampleGetReservedSize() {
	core.Println(coverageExample(GetReservedSize()))
	// Output: 1
}

func ExampleReservedNameValue() {
	core.Println(coverageExample(ReservedNameValue()))
	// Output: 1
}

func ExampleGetReservedNameValue() {
	core.Println(coverageExample(GetReservedNameValue()))
	// Output: 1
}

func ExampleReservedRootValue() {
	core.Println(coverageExample(ReservedRootValue()))
	// Output: 1
}

func ExampleGetReservedRootValue() {
	core.Println(coverageExample(GetReservedRootValue()))
	// Output: 1
}

func ExampleReservedTopValue() {
	core.Println(coverageExample(ReservedTopValue()))
	// Output: 1
}

func ExampleGetReservedTopValue() {
	core.Println(coverageExample(GetReservedTopValue()))
	// Output: 1
}

func ExampleNameValue() {
	core.Println(coverageExample(NameValue()))
	// Output: 1
}

func ExampleGetNameValue() {
	core.Println(coverageExample(GetNameValue()))
	// Output: 1
}

func ExampleRootValue() {
	core.Println(coverageExample(RootValue()))
	// Output: 1
}

func ExampleGetRootValue() {
	core.Println(coverageExample(GetRootValue()))
	// Output: 1
}

func ExampleTopValue() {
	core.Println(coverageExample(TopValue()))
	// Output: 1
}

func ExampleGetTopValue() {
	core.Println(coverageExample(GetTopValue()))
	// Output: 1
}

func ExamplePrefixSize() {
	core.Println(coverageExample(PrefixSize()))
	// Output: 1
}

func ExampleGetPrefixSize() {
	core.Println(coverageExample(GetPrefixSize()))
	// Output: 1
}

func ExampleReservedPrefixSize() {
	core.Println(coverageExample(ReservedPrefixSize()))
	// Output: 1
}

func ExampleGetReservedPrefixSize() {
	core.Println(coverageExample(GetReservedPrefixSize()))
	// Output: 1
}

func ExampleReservedEntries() {
	core.Println(coverageExample(ReservedEntries()))
	// Output: 1
}

func ExampleGetReservedEntries() {
	core.Println(coverageExample(GetReservedEntries()))
	// Output: 1
}

func ExampleReservedKeys() {
	core.Println(coverageExample(ReservedKeys()))
	// Output: 1
}

func ExampleGetReservedKeys() {
	core.Println(coverageExample(GetReservedKeys()))
	// Output: 1
}

func ExampleReservedValues() {
	core.Println(coverageExample(ReservedValues()))
	// Output: 1
}

func ExampleGetReservedValues() {
	core.Println(coverageExample(GetReservedValues()))
	// Output: 1
}

func ExampleLockedSize() {
	core.Println(coverageExample(LockedSize()))
	// Output: 1
}

func ExampleGetLockedSize() {
	core.Println(coverageExample(GetLockedSize()))
	// Output: 1
}

func ExampleLockedPrefixSize() {
	core.Println(coverageExample(LockedPrefixSize()))
	// Output: 1
}

func ExampleGetLockedPrefixSize() {
	core.Println(coverageExample(GetLockedPrefixSize()))
	// Output: 1
}

func ExampleLockedEntries() {
	core.Println(coverageExample(LockedEntries()))
	// Output: 1
}

func ExampleGetLockedEntries() {
	core.Println(coverageExample(GetLockedEntries()))
	// Output: 1
}

func ExampleLockedKeys() {
	core.Println(coverageExample(LockedKeys()))
	// Output: 1
}

func ExampleGetLockedKeys() {
	core.Println(coverageExample(GetLockedKeys()))
	// Output: 1
}

func ExampleLockedValues() {
	core.Println(coverageExample(LockedValues()))
	// Output: 1
}

func ExampleGetLockedValues() {
	core.Println(coverageExample(GetLockedValues()))
	// Output: 1
}

func ExampleGetReserved() {
	core.Println(coverageExample(GetReserved("lethean.lthn")))
	// Output: 2
}

func ExampleGetReservedName() {
	core.Println(coverageExample(GetReservedName("lethean.lthn")))
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

func ExampleGetReservedByString() {
	core.Println(coverageExample(GetReservedByString("lethean.lthn")))
	// Output: 2
}

func ExampleGetReservedByBinary() {
	core.Println(coverageExample(GetReservedByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleHasReserved() {
	core.Println(coverageExample(HasReserved("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasReserved() {
	core.Println(coverageExample(GetHasReserved("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasReserved() {
	core.Println(coverageExample(new(Service).GetHasReserved("lethean.lthn")))
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

func ExampleService_GetHasReservedString() {
	core.Println(coverageExample(new(Service).GetHasReservedString("lethean.lthn")))
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

func ExampleService_GetHasReservedBinary() {
	core.Println(coverageExample(new(Service).GetHasReservedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleHasReservedByString() {
	core.Println(coverageExample(HasReservedByString("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasReservedByString() {
	core.Println(coverageExample(GetHasReservedByString("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasReservedByString() {
	core.Println(coverageExample(new(Service).GetHasReservedByString("lethean.lthn")))
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

func ExampleService_GetHasReservedByBinary() {
	core.Println(coverageExample(new(Service).GetHasReservedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetReservedByName() {
	core.Println(coverageExample(GetReservedByName("lethean.lthn")))
	// Output: 2
}

func ExampleHasReservedByName() {
	core.Println(coverageExample(HasReservedByName("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasReservedByName() {
	core.Println(coverageExample(GetHasReservedByName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasReservedByName() {
	core.Println(coverageExample(new(Service).GetHasReservedByName("lethean.lthn")))
	// Output: 1
}

func ExampleHasReservedName() {
	core.Println(coverageExample(HasReservedName("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasReservedName() {
	core.Println(coverageExample(GetHasReservedName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasReservedName() {
	core.Println(coverageExample(new(Service).GetHasReservedName("lethean.lthn")))
	// Output: 1
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

func ExampleService_GetHasReservedHash() {
	core.Println(coverageExample(new(Service).GetHasReservedHash([32]byte{1})))
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

func ExampleService_GetHasReservedByHash() {
	core.Println(coverageExample(new(Service).GetHasReservedByHash([32]byte{1})))
	// Output: 1
}

func ExampleGetLocked() {
	core.Println(coverageExample(GetLocked("lethean.lthn")))
	// Output: 2
}

func ExampleGetLockedName() {
	core.Println(coverageExample(GetLockedName("lethean.lthn")))
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

func ExampleGetLockedByString() {
	core.Println(coverageExample(GetLockedByString("lethean.lthn")))
	// Output: 2
}

func ExampleGetLockedByBinary() {
	core.Println(coverageExample(GetLockedByBinary([]byte("lethean.lthn"))))
	// Output: 2
}

func ExampleHasLocked() {
	core.Println(coverageExample(HasLocked("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLocked() {
	core.Println(coverageExample(GetHasLocked("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasLocked() {
	core.Println(coverageExample(new(Service).GetHasLocked("lethean.lthn")))
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

func ExampleService_GetHasLockedString() {
	core.Println(coverageExample(new(Service).GetHasLockedString("lethean.lthn")))
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

func ExampleService_GetHasLockedBinary() {
	core.Println(coverageExample(new(Service).GetHasLockedBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleHasLockedByString() {
	core.Println(coverageExample(HasLockedByString("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLockedByString() {
	core.Println(coverageExample(GetHasLockedByString("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasLockedByString() {
	core.Println(coverageExample(new(Service).GetHasLockedByString("lethean.lthn")))
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

func ExampleService_GetHasLockedByBinary() {
	core.Println(coverageExample(new(Service).GetHasLockedByBinary([]byte("lethean.lthn"))))
	// Output: 1
}

func ExampleGetLockedByName() {
	core.Println(coverageExample(GetLockedByName("lethean.lthn")))
	// Output: 2
}

func ExampleHasLockedByName() {
	core.Println(coverageExample(HasLockedByName("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLockedByName() {
	core.Println(coverageExample(GetHasLockedByName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasLockedByName() {
	core.Println(coverageExample(new(Service).GetHasLockedByName("lethean.lthn")))
	// Output: 1
}

func ExampleHasLockedName() {
	core.Println(coverageExample(HasLockedName("lethean.lthn")))
	// Output: 1
}

func ExampleGetHasLockedName() {
	core.Println(coverageExample(GetHasLockedName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetHasLockedName() {
	core.Println(coverageExample(new(Service).GetHasLockedName("lethean.lthn")))
	// Output: 1
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

func ExampleService_GetHasLockedHash() {
	core.Println(coverageExample(new(Service).GetHasLockedHash([32]byte{1})))
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

func ExampleService_GetHasLockedByHash() {
	core.Println(coverageExample(new(Service).GetHasLockedByHash([32]byte{1})))
	// Output: 1
}

func ExampleService_LockedSize() {
	core.Println(coverageExample(new(Service).LockedSize()))
	// Output: 1
}

func ExampleService_GetLockedSize() {
	core.Println(coverageExample(new(Service).GetLockedSize()))
	// Output: 1
}

func ExampleService_LockedPrefixSize() {
	core.Println(coverageExample(new(Service).LockedPrefixSize()))
	// Output: 1
}

func ExampleService_GetLockedPrefixSize() {
	core.Println(coverageExample(new(Service).GetLockedPrefixSize()))
	// Output: 1
}

func ExampleService_ReservedPrefixSize() {
	core.Println(coverageExample(new(Service).ReservedPrefixSize()))
	// Output: 1
}

func ExampleService_GetReservedPrefixSize() {
	core.Println(coverageExample(new(Service).GetReservedPrefixSize()))
	// Output: 1
}

func ExampleService_LockedEntries() {
	core.Println(coverageExample(new(Service).LockedEntries()))
	// Output: 1
}

func ExampleService_GetLockedEntries() {
	core.Println(coverageExample(new(Service).GetLockedEntries()))
	// Output: 1
}

func ExampleService_LockedKeys() {
	core.Println(coverageExample(new(Service).LockedKeys()))
	// Output: 1
}

func ExampleService_GetLockedKeys() {
	core.Println(coverageExample(new(Service).GetLockedKeys()))
	// Output: 1
}

func ExampleService_LockedValues() {
	core.Println(coverageExample(new(Service).LockedValues()))
	// Output: 1
}

func ExampleService_GetLockedValues() {
	core.Println(coverageExample(new(Service).GetLockedValues()))
	// Output: 1
}
