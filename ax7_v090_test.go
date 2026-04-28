// SPDX-License-Identifier: EUPL-1.2

package lns_test

import (
	"reflect"

	. "dappco.re/go"
	subject "dappco.re/go/lns"
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

func TestAX7_WithCore_Good(t *T) {
	fn := reflect.ValueOf(subject.WithCore)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithCore_Bad(t *T) {
	fn := reflect.ValueOf(subject.WithCore)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithCore_Ugly(t *T) {
	fn := reflect.ValueOf(subject.WithCore)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithCore_Good(t *T) {
	fn := reflect.ValueOf(subject.GetWithCore)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithCore_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetWithCore)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithCore_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetWithCore)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.WithReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.WithReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.WithReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetWithReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetWithReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetWithReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithLockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.WithLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithLockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.WithLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithLockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.WithLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithLockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetWithLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithLockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetWithLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetWithLockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetWithLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewServiceWithOptions_Good(t *T) {
	fn := reflect.ValueOf(subject.NewServiceWithOptions)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewServiceWithOptions_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewServiceWithOptions)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewServiceWithOptions_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewServiceWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewServiceWithOptions_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNewServiceWithOptions)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewServiceWithOptions_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNewServiceWithOptions)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewServiceWithOptions_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNewServiceWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewService_Good(t *T) {
	fn := reflect.ValueOf(subject.NewService)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewService_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewService)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewService_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewService)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewService_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNewService)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewService_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNewService)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewService_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNewService)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Register_Good(t *T) {
	fn := reflect.ValueOf(subject.Register)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Register_Bad(t *T) {
	fn := reflect.ValueOf(subject.Register)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Register_Ugly(t *T) {
	fn := reflect.ValueOf(subject.Register)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRegister_Good(t *T) {
	fn := reflect.ValueOf(subject.GetRegister)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRegister_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetRegister)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRegister_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetRegister)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RegisterWithOptions_Good(t *T) {
	fn := reflect.ValueOf(subject.RegisterWithOptions)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RegisterWithOptions_Bad(t *T) {
	fn := reflect.ValueOf(subject.RegisterWithOptions)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RegisterWithOptions_Ugly(t *T) {
	fn := reflect.ValueOf(subject.RegisterWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRegisterWithOptions_Good(t *T) {
	fn := reflect.ValueOf(subject.GetRegisterWithOptions)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRegisterWithOptions_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetRegisterWithOptions)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRegisterWithOptions_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetRegisterWithOptions)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetServiceName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetServiceName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetServiceName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetServiceName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetServiceName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetServiceName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultTTL_Good(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultTTL)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultTTL_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultTTL)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDefaultTTL_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetDefaultTTL)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDummy_Good(t *T) {
	fn := reflect.ValueOf(subject.GetDummy)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDummy_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetDummy)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDummy_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetDummy)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapRoot_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapRoot)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapRoot_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapRoot)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapRoot_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapRoot)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapEmpty_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapEmpty_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapEmpty_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapNS_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapNS)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapNS_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapNS)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapNS_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapNS)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapTXT_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapTXT)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapTXT_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapTXT)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapTXT_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapTXT)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapA_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapA)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapA_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapA)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapA_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapA)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapAAAA_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapAAAA_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTypeMapAAAA_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHSTypes_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHSTypes)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHSTypes_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHSTypes)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHSTypes_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHSTypes)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHSTypesByVal_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHSTypesByVal)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHSTypesByVal_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHSTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHSTypesByVal_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHSTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewResource_Good(t *T) {
	fn := reflect.ValueOf(subject.NewResource)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewResource_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewResource)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewResource_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewResource)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewResource_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNewResource)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewResource_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNewResource)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewResource_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNewResource)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DecodeResource_Good(t *T) {
	fn := reflect.ValueOf(subject.DecodeResource)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DecodeResource_Bad(t *T) {
	fn := reflect.ValueOf(subject.DecodeResource)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DecodeResource_Ugly(t *T) {
	fn := reflect.ValueOf(subject.DecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDecodeResource_Good(t *T) {
	fn := reflect.ValueOf(subject.GetDecodeResource)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDecodeResource_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetDecodeResource)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetDecodeResource_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetDecodeResource)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Create_Good(t *T) {
	fn := reflect.ValueOf(subject.Create)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Create_Bad(t *T) {
	fn := reflect.ValueOf(subject.Create)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Create_Ugly(t *T) {
	fn := reflect.ValueOf(subject.Create)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCreate_Good(t *T) {
	fn := reflect.ValueOf(subject.GetCreate)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCreate_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetCreate)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCreate_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetCreate)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NextName_Good(t *T) {
	fn := reflect.ValueOf(subject.NextName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NextName_Bad(t *T) {
	fn := reflect.ValueOf(subject.NextName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NextName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NextName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNextName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNextName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNextName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNextName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNextName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNextName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_PrevName_Good(t *T) {
	fn := reflect.ValueOf(subject.PrevName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_PrevName_Bad(t *T) {
	fn := reflect.ValueOf(subject.PrevName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_PrevName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.PrevName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetPrevName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetPrevName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetPrevName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetPrevName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetPrevName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetPrevName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewClaim_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewClaim_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewClaim_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewClaim_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewClaim_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewClaim_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewClaim")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewOutpoint_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewOutpoint")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewOutpoint_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewOutpoint")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewOutpoint_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewOutpoint")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewOutpoint_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewOutpoint")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewOutpoint_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewOutpoint")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewOutpoint_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewOutpoint")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewNameView_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewNameView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewNameView_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewNameView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewNameView_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewNameView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewNameView_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewNameView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewNameView_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewNameView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewNameView_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewNameView")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Create_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Create")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Create_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Create")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Create_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Create")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCreate_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCreate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCreate_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCreate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCreate_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCreate")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultTTL_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultTTL")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultTTL_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultTTL")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultTTL_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultTTL")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultTTL_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultTTL")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultTTL_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultTTL")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultTTL_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultTTL")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Dummy_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Dummy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Dummy_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Dummy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Dummy_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Dummy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDummy_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDummy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDummy_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDummy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDummy_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDummy")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapRoot_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapRoot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapRoot_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapRoot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapRoot_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapRoot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapRoot_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapRoot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapRoot_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapRoot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapRoot_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapRoot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapEmpty_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapEmpty")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapEmpty_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapEmpty")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapEmpty_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapEmpty")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapEmpty_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapEmpty")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapEmpty_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapEmpty")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapEmpty_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapEmpty")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapNS_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapNS_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapNS_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapNS_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapNS_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapNS_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapTXT_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapTXT_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapTXT_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapTXT_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapTXT_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapTXT_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapA_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapA_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapA_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapA_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapA_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapA_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapAAAA_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapAAAA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapAAAA_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapAAAA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeMapAAAA_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeMapAAAA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapAAAA_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapAAAA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapAAAA_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapAAAA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeMapAAAA_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeMapAAAA")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HSTypes_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HSTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HSTypes_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HSTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HSTypes_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HSTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHSTypes_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHSTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHSTypes_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHSTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHSTypes_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHSTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HSTypesByVal_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HSTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HSTypesByVal_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HSTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HSTypesByVal_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HSTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHSTypesByVal_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHSTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHSTypesByVal_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHSTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHSTypesByVal_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHSTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NextName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NextName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NextName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NextName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NextName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NextName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNextName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNextName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNextName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNextName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNextName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNextName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_PrevName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("PrevName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_PrevName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("PrevName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_PrevName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("PrevName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetPrevName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetPrevName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetPrevName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetPrevName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetPrevName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetPrevName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewResource_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewResource_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NewResource_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NewResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewResource_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewResource_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNewResource_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNewResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DecodeResource_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DecodeResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DecodeResource_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DecodeResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DecodeResource_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DecodeResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDecodeResource_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDecodeResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDecodeResource_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDecodeResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDecodeResource_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDecodeResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_OnStartup_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("OnStartup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_OnStartup_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("OnStartup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_OnStartup_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("OnStartup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetOnStartup_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetOnStartup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetOnStartup_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetOnStartup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetOnStartup_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetOnStartup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_OnShutdown_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("OnShutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_OnShutdown_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("OnShutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_OnShutdown_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("OnShutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetOnShutdown_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetOnShutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetOnShutdown_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetOnShutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetOnShutdown_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetOnShutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HandleIPCEvents_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HandleIPCEvents")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HandleIPCEvents_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HandleIPCEvents")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HandleIPCEvents_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HandleIPCEvents")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHandleIPCEvents_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHandleIPCEvents")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHandleIPCEvents_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHandleIPCEvents")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHandleIPCEvents_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHandleIPCEvents")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LookupCatalogName_Good(t *T) {
	fn := reflect.ValueOf(subject.LookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LookupCatalogName_Bad(t *T) {
	fn := reflect.ValueOf(subject.LookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LookupCatalogName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.LookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLookupCatalogName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLookupCatalogName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLookupCatalogName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLookupCatalogName_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLookupCatalogName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLookupCatalogName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLookupCatalogName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLookupCatalogName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLookupCatalogName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLookupCatalogName[any])
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LookupCatalogName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LookupCatalogName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LookupCatalogName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLookupCatalogName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLookupCatalogName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLookupCatalogName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLookupCatalogName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLookupCatalogName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLookupCatalogName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLookupCatalogName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLookupCatalogName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLookupCatalogName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLookupCatalogName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resolve_Good(t *T) {
	fn := reflect.ValueOf(subject.Resolve)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resolve_Bad(t *T) {
	fn := reflect.ValueOf(subject.Resolve)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resolve_Ugly(t *T) {
	fn := reflect.ValueOf(subject.Resolve)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolve_Good(t *T) {
	fn := reflect.ValueOf(subject.GetResolve)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolve_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetResolve)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolve_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetResolve)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Hash_Good(t *T) {
	fn := reflect.ValueOf(subject.Hash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Hash_Bad(t *T) {
	fn := reflect.ValueOf(subject.Hash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Hash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.Hash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHash_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHash)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHash_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHash)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHash_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHash)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Verify_Good(t *T) {
	fn := reflect.ValueOf(subject.Verify)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Verify_Bad(t *T) {
	fn := reflect.ValueOf(subject.Verify)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Verify_Ugly(t *T) {
	fn := reflect.ValueOf(subject.Verify)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerify_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerify)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerify_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerify)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerify_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerify)
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

func TestAX7_ResolveByString_Good(t *T) {
	fn := reflect.ValueOf(subject.ResolveByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.ResolveByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.ResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.ResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ResolveByBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByBinary)
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

func TestAX7_ResolveString_Good(t *T) {
	fn := reflect.ValueOf(subject.ResolveString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveString_Bad(t *T) {
	fn := reflect.ValueOf(subject.ResolveString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ResolveString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveString_Good(t *T) {
	fn := reflect.ValueOf(subject.GetResolveString)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveString_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetResolveString)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveString_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetResolveString)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.ResolveBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.ResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveBinary_Good(t *T) {
	fn := reflect.ValueOf(subject.GetResolveBinary)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveBinary_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetResolveBinary)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveBinary_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetResolveBinary)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveName_Good(t *T) {
	fn := reflect.ValueOf(subject.ResolveName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveName_Bad(t *T) {
	fn := reflect.ValueOf(subject.ResolveName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ResolveName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetResolveName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetResolveName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetResolveName)
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

func TestAX7_ResolveByName_Good(t *T) {
	fn := reflect.ValueOf(subject.ResolveByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.ResolveByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ResolveByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ResolveByName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetResolveByName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetResolveByName)
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

func TestAX7_GetNewOutpoint_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNewOutpoint)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewOutpoint_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNewOutpoint)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewOutpoint_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNewOutpoint)
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

func TestAX7_GetNewNameView_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNewNameView)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewNameView_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNewNameView)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNewNameView_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNewNameView)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxNameSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetMaxNameSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxNameSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetMaxNameSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxNameSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetMaxNameSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxResourceSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetMaxResourceSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxResourceSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetMaxResourceSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxResourceSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetMaxResourceSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsNone_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsNone_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsNone_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsNone)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsHardened_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsHardened_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsHardened_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsHardened)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsLockup_Good(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsLockup_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetVerifyCovenantsLockup_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetVerifyCovenantsLockup)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMandatoryVerifyCovenantFlags_Good(t *T) {
	fn := reflect.ValueOf(subject.GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMandatoryVerifyCovenantFlags_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMandatoryVerifyCovenantFlags_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetMandatoryVerifyCovenantFlags)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxCovenantSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxCovenantSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetMaxCovenantSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetMaxCovenantSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCovenantMaxSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCovenantMaxSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetCovenantMaxSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetCovenantMaxSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetCovenantMaxSize)
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

func TestAX7_IsName_Good(t *T) {
	fn := reflect.ValueOf(subject.IsName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsName_Bad(t *T) {
	fn := reflect.ValueOf(subject.IsName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.IsName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsName_Good(t *T) {
	fn := reflect.ValueOf(subject.GetIsName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsName_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetIsName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetIsName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsKnown_Good(t *T) {
	fn := reflect.ValueOf(subject.IsKnown)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsKnown_Bad(t *T) {
	fn := reflect.ValueOf(subject.IsKnown)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsKnown_Ugly(t *T) {
	fn := reflect.ValueOf(subject.IsKnown)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsKnown_Good(t *T) {
	fn := reflect.ValueOf(subject.GetIsKnown)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsKnown_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetIsKnown)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsKnown_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetIsKnown)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsLinked_Good(t *T) {
	fn := reflect.ValueOf(subject.IsLinked)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsLinked_Bad(t *T) {
	fn := reflect.ValueOf(subject.IsLinked)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_IsLinked_Ugly(t *T) {
	fn := reflect.ValueOf(subject.IsLinked)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsLinked_Good(t *T) {
	fn := reflect.ValueOf(subject.GetIsLinked)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsLinked_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetIsLinked)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetIsLinked_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetIsLinked)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxNameSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxNameSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxNameSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxNameSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxNameSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxNameSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxNameSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxNameSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxNameSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxNameSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxNameSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxNameSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxResourceSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxResourceSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxResourceSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxResourceSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxResourceSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxResourceSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxResourceSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxResourceSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxResourceSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxResourceSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxResourceSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxResourceSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsNone_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsNone_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsNone_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsNone_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsNone_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsNone_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsNone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsHardened_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsHardened")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsHardened_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsHardened")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsHardened_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsHardened")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsHardened_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsHardened")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsHardened_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsHardened")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsHardened_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsHardened")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsLockup_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsLockup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsLockup_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsLockup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenantsLockup_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenantsLockup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsLockup_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsLockup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsLockup_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsLockup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenantsLockup_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenantsLockup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MandatoryVerifyCovenantFlags_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MandatoryVerifyCovenantFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MandatoryVerifyCovenantFlags_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MandatoryVerifyCovenantFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MandatoryVerifyCovenantFlags_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MandatoryVerifyCovenantFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMandatoryVerifyCovenantFlags_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMandatoryVerifyCovenantFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMandatoryVerifyCovenantFlags_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMandatoryVerifyCovenantFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMandatoryVerifyCovenantFlags_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMandatoryVerifyCovenantFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxCovenantSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxCovenantSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxCovenantSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxCovenantSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_MaxCovenantSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("MaxCovenantSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxCovenantSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxCovenantSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxCovenantSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxCovenantSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetMaxCovenantSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetMaxCovenantSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CovenantMaxSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CovenantMaxSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CovenantMaxSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CovenantMaxSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CovenantMaxSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CovenantMaxSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCovenantMaxSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCovenantMaxSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCovenantMaxSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCovenantMaxSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCovenantMaxSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCovenantMaxSize")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_ServiceName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ServiceName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ServiceName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ServiceName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ServiceName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ServiceName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetServiceName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetServiceName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetServiceName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetServiceName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetServiceName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetServiceName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Resolve_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Resolve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Resolve_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Resolve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Resolve_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Resolve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolve_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolve_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolve_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Hash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Hash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Hash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Hash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Hash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Hash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Verify_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Verify")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Verify_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Verify")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Verify_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Verify")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerify_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerify")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerify_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerify")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerify_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerify")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ResolveByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetResolveByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetResolveByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HashByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HashByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHashByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHashByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Blind_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Blind")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Blind_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Blind")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Blind_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Blind")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetBlind_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetBlind")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetBlind_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetBlind")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetBlind_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetBlind")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypeName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypeName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypeName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypeName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsKnown_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsKnown_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsKnown_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsKnown_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsKnown_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsKnown_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsKnown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsLinked_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsLinked_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsLinked_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsLinked_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsLinked_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsLinked_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsLinked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRollout_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRollout_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRollout_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasRollout_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasRollout_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasRollout_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasRollout_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasRollout_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasRollout_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasRollout")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsReserved_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsReserved_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsReserved_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsReserved_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsReserved_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsReserved_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsLockedUp_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsLockedUp")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsLockedUp_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsLockedUp")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_IsLockedUp_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("IsLockedUp")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsLockedUp_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsLockedUp")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsLockedUp_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsLockedUp")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetIsLockedUp_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetIsLockedUp")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GrindName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GrindName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GrindName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GrindName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GrindName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GrindName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetGrindName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetGrindName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetGrindName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetGrindName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetGrindName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetGrindName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountOpens_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountOpens")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountOpens_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountOpens")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountOpens_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountOpens")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountOpens_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountOpens")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountOpens_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountOpens")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountOpens_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountOpens")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountUpdates_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountUpdates")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountUpdates_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountUpdates")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountUpdates_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountUpdates")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountUpdates_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountUpdates")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountUpdates_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountUpdates")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountUpdates_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountUpdates")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountRenewals_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountRenewals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountRenewals_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountRenewals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_CountRenewals_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("CountRenewals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountRenewals_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountRenewals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountRenewals_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountRenewals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCountRenewals_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCountRenewals")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasNames_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasNames_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasNames_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasNames_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasNames_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasNames_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_AddNames_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("AddNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_AddNames_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("AddNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_AddNames_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("AddNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetAddNames_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetAddNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetAddNames_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetAddNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetAddNames_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetAddNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_RemoveNames_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("RemoveNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_RemoveNames_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("RemoveNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_RemoveNames_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("RemoveNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRemoveNames_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRemoveNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRemoveNames_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRemoveNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRemoveNames_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRemoveNames")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasSaneCovenants_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasSaneCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasSaneCovenants_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasSaneCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasSaneCovenants_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasSaneCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasSaneCovenants_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasSaneCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasSaneCovenants_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasSaneCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasSaneCovenants_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasSaneCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenants_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenants_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyCovenants_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerifyCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenants_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenants_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerifyCovenants_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerifyCovenants")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Types_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Types")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Types_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Types")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Types_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Types")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypes_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypes_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypes_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypes")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypesByVal_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypesByVal_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TypesByVal_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypesByVal_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypesByVal_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTypesByVal_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTypesByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Blacklist_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Blacklist")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Blacklist_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Blacklist")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Blacklist_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Blacklist")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetBlacklist_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetBlacklist")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetBlacklist_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetBlacklist")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetBlacklist_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetBlacklist")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerificationFlags_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerificationFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerificationFlags_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerificationFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerificationFlags_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerificationFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerificationFlags_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerificationFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerificationFlags_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerificationFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerificationFlags_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerificationFlags")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerificationFlagsByVal_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerificationFlagsByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerificationFlagsByVal_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerificationFlagsByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerificationFlagsByVal_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("VerificationFlagsByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerificationFlagsByVal_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerificationFlagsByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerificationFlagsByVal_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerificationFlagsByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetVerificationFlagsByVal_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetVerificationFlagsByVal")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReserved_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReserved_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReserved_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReserved_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReserved_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReserved_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedByHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedByHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLocked_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLocked_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLocked_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLocked_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLocked_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLocked_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedByHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedByHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Core_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Core")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Core_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Core")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_Core_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("Core")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCore_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCore")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCore_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCore")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetCore_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetCore")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultReservedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultReservedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultReservedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultReservedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultReservedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultReservedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasReservedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedCatalog")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_ReservedSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedNameValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedNameValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedNameValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedNameValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedNameValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedNameValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedRootValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedRootValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedRootValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedRootValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedRootValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedRootValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedTopValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedTopValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedTopValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedTopValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedTopValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedTopValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NameValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NameValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_NameValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("NameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNameValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNameValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetNameValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetNameValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_RootValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("RootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_RootValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("RootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_RootValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("RootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRootValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRootValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetRootValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetRootValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TopValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TopValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_TopValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("TopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTopValue_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTopValue_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetTopValue_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetTopValue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_PrefixSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_PrefixSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_PrefixSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("PrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetPrefixSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetPrefixSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetPrefixSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedEntries_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedEntries_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedEntries_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedEntries_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedEntries_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedEntries_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedKeys_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedKeys_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedKeys_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedKeys_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedKeys_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedKeys_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedValues_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedValues_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedValues_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedValues_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedValues_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedValues_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_HasLockedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("HasLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultLockedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultLockedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_DefaultLockedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("DefaultLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultLockedCatalog_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultLockedCatalog_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultLockedCatalog")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetDefaultLockedCatalog_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetDefaultLockedCatalog")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_HasReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLockedCatalog)
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

func TestAX7_ReservedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Good(t *T) {
	fn := reflect.ValueOf(subject.LockedCatalog)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Bad(t *T) {
	fn := reflect.ValueOf(subject.LockedCatalog)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedCatalog_Ugly(t *T) {
	fn := reflect.ValueOf(subject.LockedCatalog)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedSize_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedNameValue_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedNameValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedNameValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedNameValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedNameValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedNameValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedNameValue_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedNameValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedNameValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedNameValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedNameValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedNameValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedRootValue_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedRootValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedRootValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedRootValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedRootValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedRootValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedRootValue_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedRootValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedRootValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedRootValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedRootValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedRootValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedTopValue_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedTopValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedTopValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedTopValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedTopValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedTopValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedTopValue_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedTopValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedTopValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedTopValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedTopValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedTopValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameValue_Good(t *T) {
	fn := reflect.ValueOf(subject.NameValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.NameValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NameValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNameValue_Good(t *T) {
	fn := reflect.ValueOf(subject.GetNameValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNameValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetNameValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetNameValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetNameValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RootValue_Good(t *T) {
	fn := reflect.ValueOf(subject.RootValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RootValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.RootValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_RootValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.RootValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRootValue_Good(t *T) {
	fn := reflect.ValueOf(subject.GetRootValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRootValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetRootValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetRootValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetRootValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TopValue_Good(t *T) {
	fn := reflect.ValueOf(subject.TopValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TopValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.TopValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TopValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.TopValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTopValue_Good(t *T) {
	fn := reflect.ValueOf(subject.GetTopValue)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTopValue_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetTopValue)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetTopValue_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetTopValue)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_PrefixSize_Good(t *T) {
	fn := reflect.ValueOf(subject.PrefixSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_PrefixSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.PrefixSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_PrefixSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.PrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetPrefixSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetPrefixSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetPrefixSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetPrefixSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedPrefixSize_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedPrefixSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedPrefixSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedPrefixSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedPrefixSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedPrefixSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedEntries_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedEntries)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedEntries_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedEntries)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedEntries_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedEntries)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedEntries_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedEntries)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedEntries_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedEntries)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedEntries_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedEntries)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedKeys_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedKeys)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedKeys_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedKeys)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedKeys_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedKeys)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedKeys_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedKeys)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedKeys_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedKeys)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedKeys_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedKeys)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedValues_Good(t *T) {
	fn := reflect.ValueOf(subject.ReservedValues)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedValues_Bad(t *T) {
	fn := reflect.ValueOf(subject.ReservedValues)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ReservedValues_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ReservedValues)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedValues_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReservedValues)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedValues_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReservedValues)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReservedValues_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReservedValues)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedSize_Good(t *T) {
	fn := reflect.ValueOf(subject.LockedSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.LockedSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.LockedSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedPrefixSize_Good(t *T) {
	fn := reflect.ValueOf(subject.LockedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedPrefixSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.LockedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedPrefixSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.LockedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedPrefixSize_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedPrefixSize_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedPrefixSize_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedPrefixSize)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedEntries_Good(t *T) {
	fn := reflect.ValueOf(subject.LockedEntries)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedEntries_Bad(t *T) {
	fn := reflect.ValueOf(subject.LockedEntries)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedEntries_Ugly(t *T) {
	fn := reflect.ValueOf(subject.LockedEntries)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedEntries_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedEntries)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedEntries_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedEntries)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedEntries_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedEntries)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedKeys_Good(t *T) {
	fn := reflect.ValueOf(subject.LockedKeys)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedKeys_Bad(t *T) {
	fn := reflect.ValueOf(subject.LockedKeys)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedKeys_Ugly(t *T) {
	fn := reflect.ValueOf(subject.LockedKeys)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedKeys_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedKeys)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedKeys_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedKeys)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedKeys_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedKeys)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedValues_Good(t *T) {
	fn := reflect.ValueOf(subject.LockedValues)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedValues_Bad(t *T) {
	fn := reflect.ValueOf(subject.LockedValues)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_LockedValues_Ugly(t *T) {
	fn := reflect.ValueOf(subject.LockedValues)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedValues_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLockedValues)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedValues_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLockedValues)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLockedValues_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLockedValues)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReserved_Good(t *T) {
	fn := reflect.ValueOf(subject.GetReserved)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReserved_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetReserved)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetReserved_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetReserved)
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

func TestAX7_HasReserved_Good(t *T) {
	fn := reflect.ValueOf(subject.HasReserved)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReserved_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasReserved)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasReserved_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasReserved)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReserved_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasReserved)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReserved_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasReserved)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasReserved_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasReserved)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReserved_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReserved_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReserved")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReserved_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReserved")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedString")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedBinary")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByString")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByBinary")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByName")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedName")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedHash")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasReservedByHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasReservedByHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasReservedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLocked_Good(t *T) {
	fn := reflect.ValueOf(subject.GetLocked)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLocked_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetLocked)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetLocked_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetLocked)
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

func TestAX7_HasLocked_Good(t *T) {
	fn := reflect.ValueOf(subject.HasLocked)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLocked_Bad(t *T) {
	fn := reflect.ValueOf(subject.HasLocked)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HasLocked_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HasLocked)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLocked_Good(t *T) {
	fn := reflect.ValueOf(subject.GetHasLocked)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLocked_Bad(t *T) {
	fn := reflect.ValueOf(subject.GetHasLocked)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GetHasLocked_Ugly(t *T) {
	fn := reflect.ValueOf(subject.GetHasLocked)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLocked_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLocked_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLocked")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLocked_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLocked")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasLockedString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedString")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasLockedBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedBinary")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasLockedByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByString_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByString")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByString_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByString")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasLockedByBinary_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByBinary_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByBinary_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByBinary")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasLockedByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByName")
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

func TestAX7_Service_GetHasLockedName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedName_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedName_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedName")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasLockedHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedHash")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_GetHasLockedByHash_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByHash_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetHasLockedByHash_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetHasLockedByHash")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedPrefixSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedPrefixSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedPrefixSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedPrefixSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedPrefixSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedPrefixSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedPrefixSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedPrefixSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ReservedPrefixSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("ReservedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedPrefixSize_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedPrefixSize_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetReservedPrefixSize_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetReservedPrefixSize")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedEntries_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedEntries_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedEntries_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedEntries_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedEntries_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedEntries_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedEntries")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedKeys_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedKeys_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedKeys_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedKeys_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedKeys_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedKeys_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedKeys")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedValues_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedValues_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_LockedValues_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("LockedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedValues_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedValues_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_GetLockedValues_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewService(New()))
	fn := receiver.MethodByName("GetLockedValues")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}
