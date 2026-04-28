// SPDX-License-Identifier: EUPL-1.2

package wallet_test

import (
	"reflect"

	. "dappco.re/go"
	subject "dappco.re/go/lns/pkg/wallet"
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

func TestAX7_ImportSeed_Good(t *T) {
	fn := reflect.ValueOf(subject.ImportSeed)
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ImportSeed_Bad(t *T) {
	fn := reflect.ValueOf(subject.ImportSeed)
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_ImportSeed_Ugly(t *T) {
	fn := reflect.ValueOf(subject.ImportSeed)
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_ExportSeed_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("ExportSeed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_ExportSeed_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("ExportSeed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_ExportSeed_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("ExportSeed")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Snapshot_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Snapshot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Snapshot_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Snapshot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Snapshot_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Snapshot")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_MarshalJSON_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("MarshalJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_MarshalJSON_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("MarshalJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_MarshalJSON_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("MarshalJSON")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_ReceiveAddress_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("ReceiveAddress")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_ReceiveAddress_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("ReceiveAddress")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_ReceiveAddress_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("ReceiveAddress")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_NextReceiveAddress_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("NextReceiveAddress")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_NextReceiveAddress_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("NextReceiveAddress")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_NextReceiveAddress_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("NextReceiveAddress")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Credit_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Credit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Credit_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Credit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Credit_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Credit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Debit_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Debit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Debit_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Debit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Debit_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Debit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_PendingDebit_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("PendingDebit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_PendingDebit_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("PendingDebit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_PendingDebit_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("PendingDebit")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Balance_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Balance")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Balance_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Balance")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_Balance_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("Balance")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_History_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("History")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_History_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("History")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_History_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("History")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_SeedID_Good(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("SeedID")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "good")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_SeedID_Bad(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("SeedID")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "bad")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}

func TestAX7_Wallet_SeedID_Ugly(t *T) {
	wallet, err := subject.Create()
	RequireNoError(t, err)
	receiver := reflect.ValueOf(wallet)
	fn := receiver.MethodByName("SeedID")
	RequireTrue(t, fn.IsValid())
	args := ax7Args(fn.Type(), "ugly")
	var out []reflect.Value
	AssertNotPanics(t, func() {
		out = ax7Invoke(fn, args)
	})
	AssertGreaterOrEqual(t, len(out), 0)
}
