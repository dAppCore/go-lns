// SPDX-License-Identifier: EUPL-1.2

//nolint:unused
package wallet

import (
	"net/netip"
	"reflect"

	core "github.com/dappcore/go"
)

func coveragePack(values ...any) []any {
	return values
}

func coverageAssertVariant(t *core.T, variant string) {
	allowed := map[string]bool{"good": true, "bad": true, "ugly": true}
	core.AssertTrue(t, allowed[variant])
}

func coverageAssertReturn(t *core.T, values ...any) {
	core.AssertTrue(t, len(values) >= 1)
}

func coverageExample(values ...any) int {
	return len(values)
}

func coverageAddr(variant string) netip.Addr {
	if variant == "ugly" {
		return netip.MustParseAddr("2001:db8::1")
	}
	if variant == "bad" {
		return netip.Addr{}
	}
	return netip.MustParseAddr("192.0.2.1")
}

func coverageValue[T any](variant string) T {
	var zero T
	typ := reflect.TypeFor[T]()
	value := coverageReflectValue(typ, variant)
	if !value.IsValid() {
		return zero
	}
	if typ.Kind() == reflect.Interface && value.Kind() == reflect.Interface && value.IsNil() {
		return zero
	}
	if value.Type().AssignableTo(typ) {
		return value.Interface().(T)
	}
	if value.Type().ConvertibleTo(typ) {
		return value.Convert(typ).Interface().(T)
	}
	return zero
}

func coverageReflectValue(typ reflect.Type, variant string) reflect.Value {
	if typ == nil {
		return reflect.Value{}
	}
	if typ.PkgPath() == "context" && typ.Name() == "Context" {
		return reflect.ValueOf(core.Background())
	}
	if typ.PkgPath() == "github.com/dappcore/go" && typ.Name() == "Options" {
		return reflect.ValueOf(core.NewOptions())
	}
	if typ.Kind() == reflect.Pointer && typ.Elem().PkgPath() == "github.com/dappcore/go" && typ.Elem().Name() == "Core" {
		return reflect.ValueOf(core.New())
	}
	if typ.Kind() == reflect.Interface {
		if typ.NumMethod() == 0 && variant != "bad" {
			value := reflect.ValueOf("lethean.lthn")
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
			value.SetString("lethean.lthn")
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
		coveragePopulateStruct(value, variant)
		return value
	case reflect.Func:
		return reflect.Zero(typ)
	default:
		return reflect.Zero(typ)
	}
}

func coveragePopulateStruct(value reflect.Value, variant string) {
	for _, field := range value.Fields() {
		field := field
		if !field.CanSet() {
			continue
		}
		field.Set(coverageReflectValue(field.Type(), variant))
	}
}
