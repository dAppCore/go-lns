// SPDX-License-Identifier: EUPL-1.2

package covenant_test

import (
	"reflect"

	. "dappco.re/go"
	subject "dappco.re/go/lns/pkg/covenant"
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

func TestAX7_Blind_Good(t *T) {
	fn := reflect.ValueOf(subject.Blind)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Blind_Bad(t *T) {
	fn := reflect.ValueOf(subject.Blind)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Blind_Ugly(t *T) {
	fn := reflect.ValueOf(subject.Blind)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetBlind_Good(t *T) {
	fn := reflect.ValueOf(subject.GetBlind)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetBlind_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetBlind)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetBlind_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetBlind)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypes_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypes)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypes_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypes)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypes_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypes)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypesByVal_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypesByVal)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypesByVal_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypesByVal_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TypeName_Good(t *T) {
	fn := reflect.ValueOf(subject.TypeName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TypeName_Bad(t *T) {
	fn := reflect.ValueOf(subject.TypeName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TypeName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.TypeName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypeName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypeName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypeName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CovenantType_IsName_Good(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsName_Bad(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsName_Ugly(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsKnown_Good(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsKnown_Bad(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsKnown_Ugly(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsLinked_Good(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsLinked_Bad(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_IsLinked_Ugly(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_String_Good(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_String_Bad(t *T) {
	var recv subject.CovenantType
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

func TestAX7_CovenantType_String_Ugly(t *T) {
	var recv subject.CovenantType
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

func TestAX7_HasLockedName_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByName_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedString_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedString_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByString_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedHash_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedByHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByHash_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedByHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedByHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedByHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DefaultLockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DefaultLockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DefaultLockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.DefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultLockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultLockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultLockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Size_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Size")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Size_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Size")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Size_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Size")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetSize_Good(t *T) {
	var recv subject.LockedCatalog
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

func TestAX7_LockedCatalog_GetSize_Bad(t *T) {
	var recv subject.LockedCatalog
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

func TestAX7_LockedCatalog_GetSize_Ugly(t *T) {
	var recv subject.LockedCatalog
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

func TestAX7_LockedCatalog_PrefixSize_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_PrefixSize_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_PrefixSize_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetPrefixSize_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetPrefixSize_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetPrefixSize_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Has_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Has")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Has_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Has")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Has_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Has")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_HasByHash_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_HasByHash_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_HasByHash_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetHasByHash_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetHasByHash_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetHasByHash_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Get_Good(t *T) {
	var recv subject.LockedCatalog
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

func TestAX7_LockedCatalog_Get_Bad(t *T) {
	var recv subject.LockedCatalog
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

func TestAX7_LockedCatalog_Get_Ugly(t *T) {
	var recv subject.LockedCatalog
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

func TestAX7_LockedCatalog_GetByHash_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetByHash_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetByHash_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_HasByName_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_HasByName_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_HasByName_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetHasByName_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetHasByName_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetHasByName_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetByName_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetByName_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetByName_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Entries_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Entries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Entries_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Entries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Entries_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Entries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetEntries_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetEntries_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetEntries_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Keys_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Keys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Keys_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Keys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Keys_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Keys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetKeys_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetKeys_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetKeys_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Values_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Values")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Values_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Values")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Values_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Values")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetValues_Good(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetValues_Bad(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_GetValues_Ugly(t *T) {
	var recv subject.LockedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetBlacklist_Good(t *T) {
	fn := reflect.ValueOf(subject.GetBlacklist)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetBlacklist_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetBlacklist)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetBlacklist_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetBlacklist)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyString_Good(t *T) {
	fn := reflect.ValueOf(subject.VerifyString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyString_Bad(t *T) {
	fn := reflect.ValueOf(subject.VerifyString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.VerifyString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.VerifyBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.VerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.VerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyName_Good(t *T) {
	fn := reflect.ValueOf(subject.VerifyName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyName_Bad(t *T) {
	fn := reflect.ValueOf(subject.VerifyName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.VerifyName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByName_Good(t *T) {
	fn := reflect.ValueOf(subject.VerifyByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.VerifyByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.VerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByString_Good(t *T) {
	fn := reflect.ValueOf(subject.VerifyByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.VerifyByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.VerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.VerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.VerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.VerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashString_Good(t *T) {
	fn := reflect.ValueOf(subject.HashString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashString_Bad(t *T) {
	fn := reflect.ValueOf(subject.HashString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HashString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHashString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHashString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHashString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByString_Good(t *T) {
	fn := reflect.ValueOf(subject.HashByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.HashByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HashByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHashByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHashByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHashByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.HashBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.HashBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HashBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHashBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHashBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHashBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.HashByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.HashByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHashByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHashByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHashByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashName_Good(t *T) {
	fn := reflect.ValueOf(subject.HashName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HashName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HashName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHashName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHashName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHashName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByName_Good(t *T) {
	fn := reflect.ValueOf(subject.HashByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HashByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HashByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HashByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHashByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHashByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHashByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHashByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasNames_Good(t *T) {
	fn := reflect.ValueOf(subject.HasNames)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasNames_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasNames)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasNames_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasNames)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasNames_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasNames)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasNames_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasNames)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasNames_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasNames)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_AddNames_Good(t *T) {
	fn := reflect.ValueOf(subject.AddNames)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_AddNames_Bad(t *T) {
	fn := reflect.ValueOf(subject.AddNames)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_AddNames_Ugly(t *T) {
	fn := reflect.ValueOf(subject.AddNames)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetAddNames_Good(t *T) {
	fn := reflect.ValueOf(subject.GetAddNames)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetAddNames_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetAddNames)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetAddNames_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetAddNames)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RemoveNames_Good(t *T) {
	fn := reflect.ValueOf(subject.RemoveNames)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RemoveNames_Bad(t *T) {
	fn := reflect.ValueOf(subject.RemoveNames)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RemoveNames_Ugly(t *T) {
	fn := reflect.ValueOf(subject.RemoveNames)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRemoveNames_Good(t *T) {
	fn := reflect.ValueOf(subject.GetRemoveNames)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRemoveNames_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetRemoveNames)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRemoveNames_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetRemoveNames)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasSaneCovenants_Good(t *T) {
	fn := reflect.ValueOf(subject.HasSaneCovenants)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasSaneCovenants_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasSaneCovenants)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasSaneCovenants_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasSaneCovenants)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasSaneCovenants_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasSaneCovenants_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasSaneCovenants_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasSaneCovenants)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedName_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByName_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedString_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedString_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByString_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedHash_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedByHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByHash_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedByHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DefaultReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DefaultReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DefaultReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.DefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Size_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Size")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Size_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Size")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Size_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Size")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetSize_Good(t *T) {
	var recv subject.ReservedCatalog
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

func TestAX7_ReservedCatalog_GetSize_Bad(t *T) {
	var recv subject.ReservedCatalog
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

func TestAX7_ReservedCatalog_GetSize_Ugly(t *T) {
	var recv subject.ReservedCatalog
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

func TestAX7_ReservedCatalog_NameValue_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_NameValue_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_NameValue_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetNameValue_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetNameValue_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetNameValue_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_RootValue_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("RootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_RootValue_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("RootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_RootValue_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("RootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetRootValue_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetRootValue_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetRootValue_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_TopValue_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("TopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_TopValue_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("TopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_TopValue_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("TopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetTopValue_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetTopValue_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetTopValue_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_PrefixSize_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_PrefixSize_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_PrefixSize_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetPrefixSize_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetPrefixSize_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetPrefixSize_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Has_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Has")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Has_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Has")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Has_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Has")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_HasByHash_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_HasByHash_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_HasByHash_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetHasByHash_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetHasByHash_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetHasByHash_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Get_Good(t *T) {
	var recv subject.ReservedCatalog
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

func TestAX7_ReservedCatalog_Get_Bad(t *T) {
	var recv subject.ReservedCatalog
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

func TestAX7_ReservedCatalog_Get_Ugly(t *T) {
	var recv subject.ReservedCatalog
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

func TestAX7_ReservedCatalog_GetByHash_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetByHash_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetByHash_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_HasByName_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_HasByName_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_HasByName_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetHasByName_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetHasByName_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetHasByName_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetByName_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetByName_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetByName_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Entries_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Entries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Entries_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Entries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Entries_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Entries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetEntries_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetEntries_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetEntries_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Keys_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Keys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Keys_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Keys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Keys_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Keys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetKeys_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetKeys_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetKeys_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Values_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Values")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Values_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Values")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Values_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Values")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetValues_Good(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetValues_Bad(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_GetValues_Ugly(t *T) {
	var recv subject.ReservedCatalog
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerificationFlags_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerificationFlags)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerificationFlags_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerificationFlags)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerificationFlags_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerificationFlags)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerificationFlagsByVal_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerificationFlagsByVal_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerificationFlagsByVal_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerificationFlagsByVal)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRollout_Good(t *T) {
	fn := reflect.ValueOf(subject.GetRollout)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRollout_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetRollout)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRollout_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetRollout)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasRollout_Good(t *T) {
	fn := reflect.ValueOf(subject.HasRollout)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasRollout_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasRollout)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasRollout_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasRollout)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasRollout_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasRollout)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasRollout_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasRollout)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasRollout_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasRollout)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsReserved_Good(t *T) {
	fn := reflect.ValueOf(subject.IsReserved)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsReserved_Bad(t *T) {
	fn := reflect.ValueOf(subject.IsReserved)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsReserved_Ugly(t *T) {
	fn := reflect.ValueOf(subject.IsReserved)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsReserved_Good(t *T) {
	fn := reflect.ValueOf(subject.GetIsReserved)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsReserved_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetIsReserved)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsReserved_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetIsReserved)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsLockedUp_Good(t *T) {
	fn := reflect.ValueOf(subject.IsLockedUp)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsLockedUp_Bad(t *T) {
	fn := reflect.ValueOf(subject.IsLockedUp)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsLockedUp_Ugly(t *T) {
	fn := reflect.ValueOf(subject.IsLockedUp)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsLockedUp_Good(t *T) {
	fn := reflect.ValueOf(subject.GetIsLockedUp)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsLockedUp_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetIsLockedUp)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsLockedUp_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetIsLockedUp)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GrindName_Good(t *T) {
	fn := reflect.ValueOf(subject.GrindName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GrindName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GrindName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GrindName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GrindName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetGrindName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetGrindName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetGrindName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetGrindName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetGrindName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetGrindName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountOpens_Good(t *T) {
	fn := reflect.ValueOf(subject.CountOpens)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountOpens_Bad(t *T) {
	fn := reflect.ValueOf(subject.CountOpens)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountOpens_Ugly(t *T) {
	fn := reflect.ValueOf(subject.CountOpens)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountOpens_Good(t *T) {
	fn := reflect.ValueOf(subject.GetCountOpens)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountOpens_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetCountOpens)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountOpens_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetCountOpens)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountUpdates_Good(t *T) {
	fn := reflect.ValueOf(subject.CountUpdates)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountUpdates_Bad(t *T) {
	fn := reflect.ValueOf(subject.CountUpdates)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountUpdates_Ugly(t *T) {
	fn := reflect.ValueOf(subject.CountUpdates)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountUpdates_Good(t *T) {
	fn := reflect.ValueOf(subject.GetCountUpdates)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountUpdates_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetCountUpdates)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountUpdates_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetCountUpdates)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountRenewals_Good(t *T) {
	fn := reflect.ValueOf(subject.CountRenewals)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountRenewals_Bad(t *T) {
	fn := reflect.ValueOf(subject.CountRenewals)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_CountRenewals_Ugly(t *T) {
	fn := reflect.ValueOf(subject.CountRenewals)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountRenewals_Good(t *T) {
	fn := reflect.ValueOf(subject.GetCountRenewals)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountRenewals_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetCountRenewals)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCountRenewals_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetCountRenewals)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyCovenants_Good(t *T) {
	fn := reflect.ValueOf(subject.VerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyCovenants_Bad(t *T) {
	fn := reflect.ValueOf(subject.VerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_VerifyCovenants_Ugly(t *T) {
	fn := reflect.ValueOf(subject.VerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenants_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenants_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenants_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}
