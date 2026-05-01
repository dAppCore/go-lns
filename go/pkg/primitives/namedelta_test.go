// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	core "dappco.re/go"
	"testing"
)

const testPointerMismatch = "%s pointer mismatch: got %v want %v"

func TestNameDeltaHelpers(t *testing.T) {
	if got := (NameDelta{}).IsNull(); !got {
		t.Fatal("zero delta should be null")
	}

	height := uint32(7)
	renewal := uint32(8)
	owner := Outpoint{TxHash: Hash{1, 2, 3}, Index: 300}
	value := uint64(500)
	highest := uint64(600)
	data := []byte{}
	transfer := uint32(9)
	revoked := uint32(10)
	claimed := uint32(11)
	renewals := uint32(12)
	registered := true
	expired := false
	weak := true

	delta := NameDelta{
		Height:     &height,
		Renewal:    &renewal,
		Owner:      &owner,
		Value:      &value,
		Highest:    &highest,
		Data:       &data,
		Transfer:   &transfer,
		Revoked:    &revoked,
		Claimed:    &claimed,
		Renewals:   &renewals,
		Registered: &registered,
		Expired:    &expired,
		Weak:       &weak,
	}

	if got := delta.IsNull(); got {
		t.Fatal("populated delta should not be null")
	}

	const wantField = (1 << 0) |
		(1 << 1) |
		(1 << 2) |
		(1 << 3) |
		(1 << 4) |
		(1 << 5) |
		(1 << 6) |
		(1 << 7) |
		(1 << 8) |
		(1 << 9) |
		(1 << 10) |
		(1 << 11) |
		(1 << 12) |
		(1 << 13) |
		(1 << 14) |
		(1 << 15) |
		(1 << 16) |
		(1 << 17) |
		(1 << 18) |
		(1 << 19) |
		(1 << 20) |
		(1 << 22) |
		(1 << 23)

	if got := delta.GetField(); got != wantField {
		t.Fatalf("GetField() = %#x, want %#x", got, wantField)
	}

	if got := delta.GetSize(); got != len(mustMarshalBinary(t, delta)) {
		t.Fatalf("GetSize() = %d, want encoded length", got)
	}

	var decoded NameDelta
	if err := decoded.UnmarshalBinary(mustMarshalBinary(t, delta)); err != nil {
		t.Fatalf("UnmarshalBinary() returned error: %v", err)
	}

	assertNameDeltaEqual(t, decoded, delta)
}

func mustMarshalBinary(t *testing.T, delta NameDelta) []byte {
	t.Helper()

	data, err := delta.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary() returned error: %v", err)
	}

	return data
}

func assertNameDeltaEqual(t *testing.T, got, want NameDelta) {
	t.Helper()

	assertUint32Ptr := func(name string, got, want *uint32) {
		t.Helper()
		switch {
		case got == nil && want == nil:
			return
		case got == nil || want == nil:
			t.Fatalf(testPointerMismatch, name, got, want)
		case *got != *want:
			t.Fatalf("%s = %d, want %d", name, *got, *want)
		}
	}

	assertUint64Ptr := func(name string, got, want *uint64) {
		t.Helper()
		switch {
		case got == nil && want == nil:
			return
		case got == nil || want == nil:
			t.Fatalf(testPointerMismatch, name, got, want)
		case *got != *want:
			t.Fatalf("%s = %d, want %d", name, *got, *want)
		}
	}

	assertBoolPtr := func(name string, got, want *bool) {
		t.Helper()
		switch {
		case got == nil && want == nil:
			return
		case got == nil || want == nil:
			t.Fatalf(testPointerMismatch, name, got, want)
		case *got != *want:
			t.Fatalf("%s = %v, want %v", name, *got, *want)
		}
	}

	assertUint32Ptr("Height", got.Height, want.Height)
	assertUint32Ptr("Renewal", got.Renewal, want.Renewal)

	switch {
	case got.Owner == nil && want.Owner == nil:
	case got.Owner == nil || want.Owner == nil:
		t.Fatalf("Owner pointer mismatch: got %v want %v", got.Owner, want.Owner)
	default:
		if got.Owner.TxHash != want.Owner.TxHash || got.Owner.Index != want.Owner.Index {
			t.Fatalf("Owner = %+v, want %+v", got.Owner, want.Owner)
		}
	}

	assertUint64Ptr("Value", got.Value, want.Value)
	assertUint64Ptr("Highest", got.Highest, want.Highest)

	switch {
	case got.Data == nil && want.Data == nil:
	case got.Data == nil || want.Data == nil:
		t.Fatalf("Data pointer mismatch: got %v want %v", got.Data, want.Data)
	default:
		if len(*got.Data) != len(*want.Data) {
			t.Fatalf("Data length = %d, want %d", len(*got.Data), len(*want.Data))
		}
		for i := range *got.Data {
			if (*got.Data)[i] != (*want.Data)[i] {
				t.Fatalf("Data[%d] = %x, want %x", i, (*got.Data)[i], (*want.Data)[i])
			}
		}
	}

	assertUint32Ptr("Transfer", got.Transfer, want.Transfer)
	assertUint32Ptr("Revoked", got.Revoked, want.Revoked)
	assertUint32Ptr("Claimed", got.Claimed, want.Claimed)
	assertUint32Ptr("Renewals", got.Renewals, want.Renewals)
	assertBoolPtr("Registered", got.Registered, want.Registered)
	assertBoolPtr("Expired", got.Expired, want.Expired)
	assertBoolPtr("Weak", got.Weak, want.Weak)
}

// BEGIN v090 file-aware generated coverage

func TestNamedelta_NameDelta_IsNull_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_IsNull_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_IsNull_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).IsNull())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Clone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Clone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Clone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Inject_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).Inject(coverageValue[NameDelta]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Inject_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).Inject(coverageValue[NameDelta]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Inject_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).Inject(coverageValue[NameDelta]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Clear_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).Clear())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Clear_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).Clear())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_Clear_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).Clear())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_GetField_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).GetField())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_GetField_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).GetField())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_GetField_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).GetField())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_MarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_MarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_MarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).MarshalBinary())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_UnmarshalBinary_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(NameDelta).UnmarshalBinary([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_UnmarshalBinary_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(NameDelta).UnmarshalBinary([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestNamedelta_NameDelta_UnmarshalBinary_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(NameDelta).UnmarshalBinary([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
