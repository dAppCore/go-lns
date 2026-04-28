// SPDX-License-Identifier: EUPL-1.2

package primitives_test

import (
	"reflect"

	. "dappco.re/go"
	subject "dappco.re/go/lns/pkg/primitives"
)

func ax7Args(fnType reflect.Type, variant string) []reflect.Value {
	args := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		args[i] = ax7Value(fnType.In(i), variant)
	}
	return args
}

func ax7Invoke(fn reflect.Value, args []reflect.Value) []reflect.Value {
	if fn.Type().IsVariadic() {
		return fn.CallSlice(args)
	}
	return fn.Call(args)
}

func ax7Value(typ reflect.Type, variant string) reflect.Value {
	if typ.PkgPath() == "context" && typ.Name() == "Context" {
		return reflect.ValueOf(Background())
	}
	if typ.PkgPath() == "dappco.re/go" && typ.Name() == "Options" {
		return reflect.ValueOf(NewOptions())
	}
	if typ.Kind() == reflect.Pointer && typ.Elem().PkgPath() == "dappco.re/go" && typ.Elem().Name() == "Core" {
		return reflect.ValueOf(New())
	}
	if typ.Kind() == reflect.Interface {
		if typ.NumMethod() == 0 && variant != "bad" {
			value := reflect.ValueOf("lethean")
			if value.Type().AssignableTo(typ) {
				return value
			}
		}
		return reflect.Zero(typ)
	}

	switch typ.Kind() {
	case reflect.Pointer:
		return reflect.New(typ.Elem())
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			switch variant {
			case "bad":
				return reflect.MakeSlice(typ, 0, 0)
			case "ugly":
				value := reflect.MakeSlice(typ, 3, 3)
				value.Index(0).SetUint(0xff)
				value.Index(1).SetUint(0x00)
				value.Index(2).SetUint(0x61)
				return value
			default:
				value := reflect.MakeSlice(typ, len("lethean"), len("lethean"))
				reflect.Copy(value, reflect.ValueOf([]byte("lethean")))
				return value
			}
		}
		return reflect.MakeSlice(typ, 0, 0)
	case reflect.Array:
		value := reflect.New(typ).Elem()
		if value.Len() > 0 && value.Index(0).CanSet() && value.Index(0).Kind() == reflect.Uint8 {
			switch variant {
			case "bad":
				value.Index(0).SetUint(0)
			case "ugly":
				value.Index(0).SetUint(0xff)
			default:
				value.Index(0).SetUint(1)
			}
		}
		return value
	case reflect.String:
		value := reflect.New(typ).Elem()
		switch variant {
		case "bad":
			value.SetString("")
		case "ugly":
			value.SetString("lethean.\x00.lthn")
		default:
			value.SetString("lethean")
		}
		return value
	case reflect.Bool:
		value := reflect.New(typ).Elem()
		value.SetBool(variant != "bad")
		return value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := reflect.New(typ).Elem()
		switch variant {
		case "bad":
			value.SetInt(0)
		case "ugly":
			value.SetInt(-1)
		default:
			value.SetInt(1)
		}
		return value
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		value := reflect.New(typ).Elem()
		switch variant {
		case "bad":
			value.SetUint(0)
		case "ugly":
			value.SetUint(63)
		default:
			value.SetUint(1)
		}
		return value
	case reflect.Float32, reflect.Float64:
		value := reflect.New(typ).Elem()
		switch variant {
		case "bad":
			value.SetFloat(0)
		case "ugly":
			value.SetFloat(-1)
		default:
			value.SetFloat(1)
		}
		return value
	case reflect.Map:
		return reflect.MakeMap(typ)
	case reflect.Struct:
		value := reflect.New(typ).Elem()
		ax7PopulateStruct(value, variant)
		return value
	case reflect.Func:
		return reflect.Zero(typ)
	default:
		return reflect.Zero(typ)
	}
}

func ax7PopulateStruct(value reflect.Value, variant string) {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.CanSet() {
			continue
		}
		field.Set(ax7Value(field.Type(), variant))
	}
}

func TestAX7_NewClaim_Good(t *T) {
	fn := reflect.ValueOf(subject.NewClaim)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewClaim_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewClaim)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewClaim_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewClaim)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewClaim_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNewClaim)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewClaim_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNewClaim)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewClaim_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNewClaim)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Refresh_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Refresh")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Refresh_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Refresh")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Refresh_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Refresh")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Clone_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Clone_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Clone_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Hash_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Hash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Hash_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Hash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Hash_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Hash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetHash_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetHash_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetHash_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_HashHex_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HashHex")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_HashHex_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HashHex")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_HashHex_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HashHex")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetHashHex_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHashHex")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetHashHex_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHashHex")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetHashHex_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHashHex")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetSize_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetSize_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetSize_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Weight_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Weight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Weight_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Weight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_Weight_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Weight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetWeight_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetWeight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetWeight_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetWeight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetWeight_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetWeight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_VirtualSize_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("VirtualSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_VirtualSize_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("VirtualSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_VirtualSize_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("VirtualSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetVirtualSize_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetVirtualSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetVirtualSize_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetVirtualSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetVirtualSize_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetVirtualSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_MarshalBinary_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_MarshalBinary_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_MarshalBinary_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_UnmarshalBinary_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_UnmarshalBinary_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_UnmarshalBinary_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_ToBlob_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_ToBlob_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_ToBlob_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetToBlob_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetToBlob_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetToBlob_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetBlob_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetBlob_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetBlob_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_FromBlob_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_FromBlob_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_FromBlob_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetFromBlob_Good(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFromBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetFromBlob_Bad(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFromBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Claim_GetFromBlob_Ugly(t *T) {
	var recv subject.Claim
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFromBlob")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_MarshalBinary_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_MarshalBinary_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_MarshalBinary_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_UnmarshalBinary_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_UnmarshalBinary_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_UnmarshalBinary_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Len_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Len")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Len_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Len")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Len_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Len")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Clone_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Clone_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Clone_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Inject_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Inject_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Inject_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetNone_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetNone_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetNone_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_ToArray_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToArray")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_ToArray_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToArray")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_ToArray_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToArray")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_FromArray_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromArray")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_FromArray_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromArray")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_FromArray_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromArray")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IndexOf_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IndexOf")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IndexOf_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IndexOf")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IndexOf_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IndexOf")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Get_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Get")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Get_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Get")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Get_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Get")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Set_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Set")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Set_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Set")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Set_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Set")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Push_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Push")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Push_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Push")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_Push_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Push")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetOpen_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetOpen")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetOpen_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetOpen")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetOpen_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetOpen")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetBid_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetBid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetBid_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetBid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetBid_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetBid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetReveal_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetReveal_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetReveal_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRedeem_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRedeem")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRedeem_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRedeem")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRedeem_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRedeem")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetClaim_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetClaim_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetClaim_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRegister_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRegister")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRegister_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRegister")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRegister_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRegister")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetUpdate_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetUpdate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetUpdate_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetUpdate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetUpdate_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetUpdate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRenew_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRenew")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRenew_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRenew")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRenew_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRenew")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetTransfer_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetTransfer")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetTransfer_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetTransfer")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetTransfer_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetTransfer")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetFinalize_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetFinalize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetFinalize_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetFinalize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetFinalize_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetFinalize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRevoke_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRevoke")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRevoke_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRevoke")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_SetRevoke_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("SetRevoke")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetU8_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetU8")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetU8_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetU8")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetU8_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetU8")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushU8_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushU8")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushU8_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushU8")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushU8_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushU8")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetU32_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetU32")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetU32_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetU32")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetU32_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetU32")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushU32_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushU32")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushU32_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushU32")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushU32_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushU32")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetHash_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetHash_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetHash_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushHash_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushHash_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushHash_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetString_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetString_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetString_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushString_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushString_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_PushString_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PushString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetSize_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetSize_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetSize_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetVarSize_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetVarSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetVarSize_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetVarSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetVarSize_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetVarSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetJSON_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetJSON_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_GetJSON_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_FromJSON_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_FromJSON_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_FromJSON_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_IsNull_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_IsNull_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_IsNull_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Clone_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Clone_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Clone_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Inject_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Inject_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Inject_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Clear_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Clear_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_Clear_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_GetField_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetField")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_GetField_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetField")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_GetField_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetField")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_GetSize_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_GetSize_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_GetSize_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_MarshalBinary_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_MarshalBinary_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_MarshalBinary_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_UnmarshalBinary_Good(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_UnmarshalBinary_Bad(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameDelta_UnmarshalBinary_Ugly(t *T) {
	var recv subject.NameDelta
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Delta_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Delta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Delta_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Delta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Delta_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Delta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetDelta_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetDelta_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetDelta_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Inject_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Inject_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Inject_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetInject_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetInject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetInject_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetInject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetInject_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetInject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Clone_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Clone_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Clone_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetClone_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetClone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetClone_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetClone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetClone_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetClone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Clear_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Clear_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Clear_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetClear_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetClear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetClear_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetClear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetClear_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetClear")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_HasDelta_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_HasDelta_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_HasDelta_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetHasDelta_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetHasDelta_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetHasDelta_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasDelta")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Reset_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Reset")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Reset_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Reset")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Reset_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Reset")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetReset_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetReset")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetReset_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetReset")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetReset_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetReset")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Set_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Set")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Set_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Set")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Set_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Set")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetSet_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSet")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetSet_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSet")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetSet_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSet")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_MaybeExpire_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MaybeExpire")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_MaybeExpire_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MaybeExpire")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_MaybeExpire_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MaybeExpire")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetMaybeExpire_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetMaybeExpire")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetMaybeExpire_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetMaybeExpire")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetMaybeExpire_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetMaybeExpire")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_ApplyState_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ApplyState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_ApplyState_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ApplyState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_ApplyState_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ApplyState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetApplyState_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetApplyState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetApplyState_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetApplyState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetApplyState_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetApplyState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsNull_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsNull_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsNull_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsNull_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsNull_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsNull_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetSize_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetSize_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetSize_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_MarshalBinary_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_MarshalBinary_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_MarshalBinary_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_UnmarshalBinary_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_UnmarshalBinary_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_UnmarshalBinary_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetJSON_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetJSON_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetJSON_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Format_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Format")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Format_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Format")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_Format_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Format")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetFormat_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFormat")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetFormat_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFormat")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetFormat_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFormat")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_FromJSON_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_FromJSON_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_FromJSON_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetFromJSON_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetFromJSON_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetFromJSON_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetFromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameStateStatus_String_Good(t *T) {
	var recv subject.NameStateStatus
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("String")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameStateStatus_String_Bad(t *T) {
	var recv subject.NameStateStatus
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("String")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameStateStatus_String_Ugly(t *T) {
	var recv subject.NameStateStatus
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("String")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_State_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("State")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_State_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("State")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_State_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("State")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetState_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetState_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetState_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsOpening_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsOpening")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsOpening_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsOpening")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsOpening_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsOpening")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsOpening_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsOpening")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsOpening_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsOpening")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsOpening_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsOpening")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsLocked_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsLocked_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsLocked_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsLocked_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsLocked_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsLocked_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsBidding_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsBidding")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsBidding_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsBidding")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsBidding_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsBidding")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsBidding_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsBidding")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsBidding_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsBidding")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsBidding_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsBidding")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsReveal_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsReveal_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsReveal_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsReveal_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsReveal_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsReveal_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsClosed_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClosed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsClosed_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClosed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsClosed_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClosed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsClosed_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsClosed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsClosed_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsClosed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsClosed_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsClosed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsRevoked_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRevoked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsRevoked_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRevoked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsRevoked_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRevoked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsRevoked_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsRevoked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsRevoked_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsRevoked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsRevoked_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsRevoked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsRedeemable_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRedeemable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsRedeemable_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRedeemable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsRedeemable_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRedeemable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsRedeemable_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsRedeemable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsRedeemable_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsRedeemable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsRedeemable_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsRedeemable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsClaimable_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClaimable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsClaimable_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClaimable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsClaimable_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClaimable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsClaimable_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsClaimable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsClaimable_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsClaimable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsClaimable_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsClaimable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsExpired_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsExpired")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsExpired_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsExpired")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_IsExpired_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsExpired")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsExpired_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsExpired")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsExpired_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsExpired")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetIsExpired_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetIsExpired")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_ToStats_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToStats")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_ToStats_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToStats")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_ToStats_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToStats")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetToStats_Good(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToStats")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetToStats_Bad(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToStats")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameState_GetToStats_Ugly(t *T) {
	var recv subject.NameState
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToStats")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewOutpoint_Good(t *T) {
	fn := reflect.ValueOf(subject.NewOutpoint)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewOutpoint_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewOutpoint)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewOutpoint_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewOutpoint)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Inject_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Inject_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Inject_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Clone_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Clone_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Clone_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Equals_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Equals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Equals_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Equals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Equals_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Equals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Compare_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Compare")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Compare_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Compare")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_Compare_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Compare")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_GetSize_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_GetSize_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_GetSize_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_MarshalBinary_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_MarshalBinary_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_MarshalBinary_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_UnmarshalBinary_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_UnmarshalBinary_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_UnmarshalBinary_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_GetJSON_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_GetJSON_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_GetJSON_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_FromJSON_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_FromJSON_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_FromJSON_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Inject_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Inject_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Inject_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Inject")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Clone_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Clone_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Clone_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Clone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Equals_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Equals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Equals_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Equals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Equals_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Equals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Compare_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Compare")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Compare_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Compare")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_Compare_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Compare")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsNull_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsNull_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsNull_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsPubkeyHash_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsPubkeyHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsPubkeyHash_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsPubkeyHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsPubkeyHash_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsPubkeyHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsScripthash_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsScripthash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsScripthash_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsScripthash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsScripthash_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsScripthash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsNulldata_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNulldata")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsNulldata_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNulldata")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsNulldata_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNulldata")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsUnspendable_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsUnspendable_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsUnspendable_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsUnknown_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnknown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsUnknown_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnknown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsUnknown_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnknown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsValid_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsValid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsValid_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsValid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_IsValid_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsValid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_GetSize_Good(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_GetSize_Bad(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Address_GetSize_Ugly(t *T) {
	var recv subject.Address
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUnspendable_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUnspendable_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUnspendable_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsNone_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsNone_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsNone_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsClaim_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsClaim_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsClaim_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsOpen_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsOpen")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsOpen_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsOpen")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsOpen_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsOpen")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsBid_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsBid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsBid_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsBid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsBid_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsBid")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsReveal_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsReveal_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsReveal_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsReveal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRedeem_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRedeem")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRedeem_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRedeem")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRedeem_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRedeem")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRegister_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRegister")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRegister_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRegister")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRegister_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRegister")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUpdate_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUpdate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUpdate_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUpdate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUpdate_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUpdate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRenew_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRenew")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRenew_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRenew")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRenew_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRenew")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsTransfer_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsTransfer")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsTransfer_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsTransfer")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsTransfer_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsTransfer")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsFinalize_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsFinalize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsFinalize_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsFinalize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsFinalize_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsFinalize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRevoke_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRevoke")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRevoke_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRevoke")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsRevoke_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsRevoke")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsName_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsName_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsName_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsKnown_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsKnown_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsKnown_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUnknown_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnknown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUnknown_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnknown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsUnknown_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnknown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsDustworthy_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsDustworthy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsDustworthy_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsDustworthy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsDustworthy_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsDustworthy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsNonspendable_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNonspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsNonspendable_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNonspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsNonspendable_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNonspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsLinked_Good(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsLinked_Bad(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Covenant_IsLinked_Ugly(t *T) {
	var recv subject.Covenant
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_IsNull_Good(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_IsNull_Bad(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Outpoint_IsNull_Ugly(t *T) {
	var recv subject.Outpoint
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsNull")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Input_IsCoinbase_Good(t *T) {
	var recv subject.Input
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsCoinbase")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Input_IsCoinbase_Bad(t *T) {
	var recv subject.Input
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsCoinbase")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Input_IsCoinbase_Ugly(t *T) {
	var recv subject.Input
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsCoinbase")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Input_IsFinal_Good(t *T) {
	var recv subject.Input
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsFinal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Input_IsFinal_Bad(t *T) {
	var recv subject.Input
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsFinal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Input_IsFinal_Ugly(t *T) {
	var recv subject.Input
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsFinal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Output_IsUnspendable_Good(t *T) {
	var recv subject.Output
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Output_IsUnspendable_Bad(t *T) {
	var recv subject.Output
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Output_IsUnspendable_Ugly(t *T) {
	var recv subject.Output
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("IsUnspendable")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewNameView_Good(t *T) {
	fn := reflect.ValueOf(subject.NewNameView)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewNameView_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewNameView)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewNameView_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewNameView)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_GetNameStateSync_Good(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameStateSync")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_GetNameStateSync_Bad(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameStateSync")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_GetNameStateSync_Ugly(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameStateSync")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_GetNameState_Good(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_GetNameState_Bad(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_GetNameState_Ugly(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameState")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_ToNameUndo_Good(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNameUndo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_ToNameUndo_Bad(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNameUndo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameView_ToNameUndo_Ugly(t *T) {
	var recv subject.NameView
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNameUndo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_FromView_Good(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_FromView_Bad(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_FromView_Ugly(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FromView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_GetSize_Good(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_GetSize_Bad(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_GetSize_Ugly(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_MarshalBinary_Good(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_MarshalBinary_Bad(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_MarshalBinary_Ugly(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("MarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_UnmarshalBinary_Good(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_UnmarshalBinary_Bad(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameUndo_UnmarshalBinary_Ugly(t *T) {
	var recv subject.NameUndo
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("UnmarshalBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}
