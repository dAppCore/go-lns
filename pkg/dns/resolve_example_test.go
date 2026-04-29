// SPDX-License-Identifier: EUPL-1.2

package dns

import core "dappco.re/go"

func ExampleGetServiceName() {
	fn := ax7ValueOf(GetServiceName)
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

func ExampleNewService() {
	fn := ax7ValueOf(NewService)
	out := ax7Invoke(fn, ax7Args(fn.Type(), "good"))
	core.Println(len(out))
	// Output: 1
}

func ExampleWithCore() {
	fn := ax7ValueOf(WithCore)
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
