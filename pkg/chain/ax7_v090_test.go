// SPDX-License-Identifier: EUPL-1.2

package chain_test

import (
	"reflect"

	. "dappco.re/go"
	subject "dappco.re/go/lns/pkg/chain"
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

func TestAX7_WithAliasHTTPClient_Good(t *T) {
	fn := reflect.ValueOf(subject.WithAliasHTTPClient)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithAliasHTTPClient_Bad(t *T) {
	fn := reflect.ValueOf(subject.WithAliasHTTPClient)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithAliasHTTPClient_Ugly(t *T) {
	fn := reflect.ValueOf(subject.WithAliasHTTPClient)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewAliasClient_Good(t *T) {
	fn := reflect.ValueOf(subject.NewAliasClient)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewAliasClient_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewAliasClient)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewAliasClient_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewAliasClient)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_AliasClient_GetAliases_Good(t *T) {
	var recv subject.AliasClient
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetAliases")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_AliasClient_GetAliases_Bad(t *T) {
	var recv subject.AliasClient
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetAliases")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_AliasClient_GetAliases_Ugly(t *T) {
	var recv subject.AliasClient
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetAliases")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ParseAliasComment_Good(t *T) {
	fn := reflect.ValueOf(subject.ParseAliasComment)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ParseAliasComment_Bad(t *T) {
	fn := reflect.ValueOf(subject.ParseAliasComment)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ParseAliasComment_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ParseAliasComment)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HNSName_Good(t *T) {
	fn := reflect.ValueOf(subject.HNSName)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HNSName_Bad(t *T) {
	fn := reflect.ValueOf(subject.HNSName)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_HNSName_Ugly(t *T) {
	fn := reflect.ValueOf(subject.HNSName)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DiscoverNames_Good(t *T) {
	fn := reflect.ValueOf(subject.DiscoverNames)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DiscoverNames_Bad(t *T) {
	fn := reflect.ValueOf(subject.DiscoverNames)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_DiscoverNames_Ugly(t *T) {
	fn := reflect.ValueOf(subject.DiscoverNames)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithAPIKey_Good(t *T) {
	fn := reflect.ValueOf(subject.WithAPIKey)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithAPIKey_Bad(t *T) {
	fn := reflect.ValueOf(subject.WithAPIKey)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithAPIKey_Ugly(t *T) {
	fn := reflect.ValueOf(subject.WithAPIKey)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithHTTPClient_Good(t *T) {
	fn := reflect.ValueOf(subject.WithHTTPClient)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithHTTPClient_Bad(t *T) {
	fn := reflect.ValueOf(subject.WithHTTPClient)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_WithHTTPClient_Ugly(t *T) {
	fn := reflect.ValueOf(subject.WithHTTPClient)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewClient_Good(t *T) {
	fn := reflect.ValueOf(subject.NewClient)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewClient_Bad(t *T) {
	fn := reflect.ValueOf(subject.NewClient)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_NewClient_Ugly(t *T) {
	fn := reflect.ValueOf(subject.NewClient)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_RawCall_Good(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("RawCall")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_RawCall_Bad(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("RawCall")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_RawCall_Ugly(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("RawCall")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlockchainInfo_Good(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlockchainInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlockchainInfo_Bad(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlockchainInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlockchainInfo_Ugly(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlockchainInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetInfo_Good(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetInfo_Bad(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetInfo_Ugly(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetNameResource_Good(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetNameResource_Bad(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetNameResource_Ugly(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetNameResource")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlock_Good(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlock")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlock_Bad(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlock")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlock_Ugly(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlock")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlockByHeight_Good(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlockByHeight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlockByHeight_Bad(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlockByHeight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetBlockByHeight_Ugly(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetBlockByHeight")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetPeerInfo_Good(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPeerInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetPeerInfo_Bad(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPeerInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Client_GetPeerInfo_Ugly(t *T) {
	var recv subject.Client
	receiver := reflect.ValueOf(&recv)
	fn := receiver.MethodByName("GetPeerInfo")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}
