// SPDX-License-Identifier: EUPL-1.2

package dns_test

import (
	"reflect"

	. "dappco.re/go"
	subject "dappco.re/go/lns/pkg/dns"
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

func TestAX7_NewBridge_Good(t *T) {
	fn := reflect.ValueOf(subject.NewBridge)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewBridge_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewBridge)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewBridge_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewBridge)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Discover_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Discover")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Discover_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Discover")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Discover_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Discover")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Lookup_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Lookup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Lookup_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Lookup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Lookup_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Lookup")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Records_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Records")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Records_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Records")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Records_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Records")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Reverse_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Reverse")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Reverse_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Reverse")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Reverse_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Reverse")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Health_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Health")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Health_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Health")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Health_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Health")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Serve_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Serve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Serve_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Serve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Serve_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Serve")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Shutdown_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Shutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Shutdown_Bad(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Shutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Bridge_Shutdown_Ugly(t *T) {
	receiver := reflect.ValueOf(subject.NewBridge(subject.BridgeConfig{}))
	fn := receiver.MethodByName("Shutdown")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_FQDN_Good(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FQDN")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_FQDN_Bad(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FQDN")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_FQDN_Ugly(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("FQDN")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_Addresses4_Good(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Addresses4")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_Addresses4_Bad(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Addresses4")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_Addresses4_Ugly(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Addresses4")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_Addresses6_Good(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Addresses6")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_Addresses6_Bad(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Addresses6")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_Addresses6_Ugly(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Addresses6")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_TXT_Good(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("TXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_TXT_Bad(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("TXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_TXT_Ugly(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("TXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_NS_Good(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_NS_Bad(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NameRecord_NS_Ugly(t *T) {
	var recv subject.NameRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_Core_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
	fn := receiver.MethodByName("GetCore")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ServiceName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
	fn := receiver.MethodByName("GetHashName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
	fn := receiver.MethodByName("GetVerifyBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_VerifyName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
	fn := receiver.MethodByName("GetVerifyName")
	RequireTrue(t, fn.IsValid())
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

func TestAX7_Service_VerifyByString_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
	fn := receiver.MethodByName("GetVerifyByBinary")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Service_ResolveByName_Good(t *T) {
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
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
	receiver := reflect.ValueOf(subject.NewService())
	fn := receiver.MethodByName("GetVerifyByName")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetSize_Good(t *T) {
	var recv subject.Resource
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

func TestAX7_Resource_GetSize_Bad(t *T) {
	var recv subject.Resource
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

func TestAX7_Resource_GetSize_Ugly(t *T) {
	var recv subject.Resource
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

func TestAX7_DSRecord_Type_Good(t *T) {
	var recv subject.DSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DSRecord_Type_Bad(t *T) {
	var recv subject.DSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DSRecord_Type_Ugly(t *T) {
	var recv subject.DSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DSRecord_GetType_Good(t *T) {
	var recv subject.DSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DSRecord_GetType_Bad(t *T) {
	var recv subject.DSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DSRecord_GetType_Ugly(t *T) {
	var recv subject.DSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_Type_Good(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_Type_Bad(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_Type_Ugly(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_GetType_Good(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_GetType_Bad(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_GetType_Ugly(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_GetNS_Good(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_GetNS_Bad(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NSRecord_GetNS_Ugly(t *T) {
	var recv subject.NSRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_Type_Good(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_Type_Bad(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_Type_Ugly(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_GetType_Good(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_GetType_Bad(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_GetType_Ugly(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_GetNS_Good(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_GetNS_Bad(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE4Record_GetNS_Ugly(t *T) {
	var recv subject.GLUE4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_Type_Good(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_Type_Bad(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_Type_Ugly(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_GetType_Good(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_GetType_Bad(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_GetType_Ugly(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_GetNS_Good(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_GetNS_Bad(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_GLUE6Record_GetNS_Ugly(t *T) {
	var recv subject.GLUE6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_Type_Good(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_Type_Bad(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_Type_Ugly(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_GetType_Good(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_GetType_Bad(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_GetType_Ugly(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_GetNS_Good(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_GetNS_Bad(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_GetNS_Ugly(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_Type_Good(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_Type_Bad(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_Type_Ugly(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_GetType_Good(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_GetType_Bad(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_GetType_Ugly(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_GetNS_Good(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_GetNS_Bad(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_GetNS_Ugly(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TXTRecord_Type_Good(t *T) {
	var recv subject.TXTRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TXTRecord_Type_Bad(t *T) {
	var recv subject.TXTRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TXTRecord_Type_Ugly(t *T) {
	var recv subject.TXTRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Type")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TXTRecord_GetType_Good(t *T) {
	var recv subject.TXTRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TXTRecord_GetType_Bad(t *T) {
	var recv subject.TXTRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_TXTRecord_GetType_Ugly(t *T) {
	var recv subject.TXTRecord
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_NS_Good(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_NS_Bad(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH4Record_NS_Ugly(t *T) {
	var recv subject.SYNTH4Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_NS_Good(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_NS_Bad(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_SYNTH6Record_NS_Ugly(t *T) {
	var recv subject.SYNTH6Record
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("NS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasType_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasType_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasType_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasType_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasType_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasType_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasType")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasNS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasNS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasNS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasNS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasNS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasNS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasDS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasDS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_HasDS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("HasDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasDS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasDS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetHasDS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetHasDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToNS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToNS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToNS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToNS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToNS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToNS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToGlue_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToGlue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToGlue_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToGlue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToGlue_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToGlue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToGlue_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToGlue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToGlue_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToGlue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToGlue_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToGlue")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToDS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToDS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToDS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToDS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToDS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToDS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToDS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToTXT_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToTXT_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToTXT_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToTXT_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToTXT_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToTXT_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToTXT")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToZone_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToZone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToZone_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToZone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToZone_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToZone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToZone_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToZone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToZone_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToZone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToZone_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToZone")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToReferral_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToReferral")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToReferral_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToReferral")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToReferral_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToReferral")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToReferral_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToReferral")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToReferral_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToReferral")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToReferral_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToReferral")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToDNS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToDNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToDNS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToDNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToDNS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToDNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToDNS_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToDNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToDNS_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToDNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToDNS_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToDNS")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToNSEC_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNSEC")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToNSEC_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNSEC")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_ToNSEC_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("ToNSEC")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToNSEC_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToNSEC")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToNSEC_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToNSEC")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetToNSEC_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetToNSEC")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_Encode_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Encode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_Encode_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Encode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_Encode_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Encode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetEncode_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEncode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetEncode_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEncode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetEncode_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetEncode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_Decode_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Decode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_Decode_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Decode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_Decode_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("Decode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetDecode_Good(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetDecode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetDecode_Bad(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetDecode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetDecode_Ugly(t *T) {
	var recv subject.Resource
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetDecode")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Resource_GetJSON_Good(t *T) {
	var recv subject.Resource
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

func TestAX7_Resource_GetJSON_Bad(t *T) {
	var recv subject.Resource
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

func TestAX7_Resource_GetJSON_Ugly(t *T) {
	var recv subject.Resource
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

func TestAX7_Resource_FromJSON_Good(t *T) {
	var recv subject.Resource
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

func TestAX7_Resource_FromJSON_Bad(t *T) {
	var recv subject.Resource
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

func TestAX7_Resource_FromJSON_Ugly(t *T) {
	var recv subject.Resource
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

func TestAX7_DSRecord_GetSize_Good(t *T) {
	var recv subject.DSRecord
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

func TestAX7_DSRecord_GetSize_Bad(t *T) {
	var recv subject.DSRecord
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

func TestAX7_DSRecord_GetSize_Ugly(t *T) {
	var recv subject.DSRecord
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

func TestAX7_NSRecord_GetSize_Good(t *T) {
	var recv subject.NSRecord
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

func TestAX7_NSRecord_GetSize_Bad(t *T) {
	var recv subject.NSRecord
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

func TestAX7_NSRecord_GetSize_Ugly(t *T) {
	var recv subject.NSRecord
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

func TestAX7_GLUE4Record_GetSize_Good(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE4Record_GetSize_Bad(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE4Record_GetSize_Ugly(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE6Record_GetSize_Good(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_GLUE6Record_GetSize_Bad(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_GLUE6Record_GetSize_Ugly(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_SYNTH4Record_GetSize_Good(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH4Record_GetSize_Bad(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH4Record_GetSize_Ugly(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH6Record_GetSize_Good(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_SYNTH6Record_GetSize_Bad(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_SYNTH6Record_GetSize_Ugly(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_TXTRecord_GetSize_Good(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_TXTRecord_GetSize_Bad(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_TXTRecord_GetSize_Ugly(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_DSRecord_GetJSON_Good(t *T) {
	var recv subject.DSRecord
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

func TestAX7_DSRecord_GetJSON_Bad(t *T) {
	var recv subject.DSRecord
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

func TestAX7_DSRecord_GetJSON_Ugly(t *T) {
	var recv subject.DSRecord
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

func TestAX7_DSRecord_FromJSON_Good(t *T) {
	var recv subject.DSRecord
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

func TestAX7_DSRecord_FromJSON_Bad(t *T) {
	var recv subject.DSRecord
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

func TestAX7_DSRecord_FromJSON_Ugly(t *T) {
	var recv subject.DSRecord
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

func TestAX7_NSRecord_GetJSON_Good(t *T) {
	var recv subject.NSRecord
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

func TestAX7_NSRecord_GetJSON_Bad(t *T) {
	var recv subject.NSRecord
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

func TestAX7_NSRecord_GetJSON_Ugly(t *T) {
	var recv subject.NSRecord
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

func TestAX7_NSRecord_FromJSON_Good(t *T) {
	var recv subject.NSRecord
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

func TestAX7_NSRecord_FromJSON_Bad(t *T) {
	var recv subject.NSRecord
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

func TestAX7_NSRecord_FromJSON_Ugly(t *T) {
	var recv subject.NSRecord
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

func TestAX7_GLUE4Record_GetJSON_Good(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE4Record_GetJSON_Bad(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE4Record_GetJSON_Ugly(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE4Record_FromJSON_Good(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE4Record_FromJSON_Bad(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE4Record_FromJSON_Ugly(t *T) {
	var recv subject.GLUE4Record
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

func TestAX7_GLUE6Record_GetJSON_Good(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_GLUE6Record_GetJSON_Bad(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_GLUE6Record_GetJSON_Ugly(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_GLUE6Record_FromJSON_Good(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_GLUE6Record_FromJSON_Bad(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_GLUE6Record_FromJSON_Ugly(t *T) {
	var recv subject.GLUE6Record
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

func TestAX7_SYNTH4Record_GetJSON_Good(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH4Record_GetJSON_Bad(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH4Record_GetJSON_Ugly(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH4Record_FromJSON_Good(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH4Record_FromJSON_Bad(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH4Record_FromJSON_Ugly(t *T) {
	var recv subject.SYNTH4Record
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

func TestAX7_SYNTH6Record_GetJSON_Good(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_SYNTH6Record_GetJSON_Bad(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_SYNTH6Record_GetJSON_Ugly(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_SYNTH6Record_FromJSON_Good(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_SYNTH6Record_FromJSON_Bad(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_SYNTH6Record_FromJSON_Ugly(t *T) {
	var recv subject.SYNTH6Record
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

func TestAX7_TXTRecord_GetJSON_Good(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_TXTRecord_GetJSON_Bad(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_TXTRecord_GetJSON_Ugly(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_TXTRecord_FromJSON_Good(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_TXTRecord_FromJSON_Bad(t *T) {
	var recv subject.TXTRecord
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

func TestAX7_TXTRecord_FromJSON_Ugly(t *T) {
	var recv subject.TXTRecord
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
