// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleGetServiceName() {
	core.Println(coverageExample(GetServiceName()))
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

func ExampleNewService() {
	core.Println(coverageExample(NewService()))
	// Output: 1
}

func ExampleWithCore() {
	core.Println(coverageExample(WithCore(core.New())))
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

func ExampleService_VerifyName() {
	core.Println(coverageExample(new(Service).VerifyName("lethean.lthn")))
	// Output: 1
}

func ExampleService_GetVerifyName() {
	core.Println(coverageExample(new(Service).GetVerifyName("lethean.lthn")))
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
