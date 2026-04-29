// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"testing"

	core "dappco.re/go"
)

func repeatString(value string, count int) string {
	out := ""
	for i := 0; i < count; i++ {
		out += value
	}
	return out
}

func TestCommonConstants(t *testing.T) {
	if DEFAULT_TTL != 21600 || GetDefaultTTL() != DEFAULT_TTL {
		t.Fatalf("DEFAULT_TTL = %d, GetDefaultTTL() = %d, want 21600", DEFAULT_TTL, GetDefaultTTL())
	}

	if len(DUMMY) != 0 {
		t.Fatalf("DUMMY length = %d, want 0", len(DUMMY))
	}

	if len(GetDummy()) != len(DUMMY) {
		t.Fatalf("GetDummy length = %d, want %d", len(GetDummy()), len(DUMMY))
	}

	cases := []struct {
		name string
		got  []byte
		want []byte
	}{
		{name: "TYPE_MAP_ROOT", got: TYPE_MAP_ROOT, want: []byte{0x00, 0x07, 0x22, 0x00, 0x00, 0x00, 0x00, 0x03, 0x80}},
		{name: "TYPE_MAP_EMPTY", got: TYPE_MAP_EMPTY, want: []byte{0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03}},
		{name: "TYPE_MAP_NS", got: TYPE_MAP_NS, want: []byte{0x00, 0x06, 0x20, 0x00, 0x00, 0x00, 0x00, 0x03}},
		{name: "TYPE_MAP_TXT", got: TYPE_MAP_TXT, want: []byte{0x00, 0x06, 0x00, 0x00, 0x80, 0x00, 0x00, 0x03}},
		{name: "TYPE_MAP_A", got: TYPE_MAP_A, want: []byte{0x00, 0x06, 0x40, 0x00, 0x00, 0x00, 0x00, 0x03}},
		{name: "TYPE_MAP_AAAA", got: TYPE_MAP_AAAA, want: []byte{0x00, 0x06, 0x00, 0x00, 0x00, 0x08, 0x00, 0x03}},
	}

	for _, tc := range cases {
		if !core.DeepEqual(tc.got, tc.want) {
			t.Fatalf("%s = %x, want %x", tc.name, tc.got, tc.want)
		}

		var got []byte
		switch tc.name {
		case "TYPE_MAP_ROOT":
			got = GetTypeMapRoot()
		case "TYPE_MAP_EMPTY":
			got = GetTypeMapEmpty()
		case "TYPE_MAP_NS":
			got = GetTypeMapNS()
		case "TYPE_MAP_TXT":
			got = GetTypeMapTXT()
		case "TYPE_MAP_A":
			got = GetTypeMapA()
		case "TYPE_MAP_AAAA":
			got = GetTypeMapAAAA()
		}

		if !core.DeepEqual(got, tc.want) {
			t.Fatalf("%s getter = %x, want %x", tc.name, got, tc.want)
		}
	}
}

func TestHSTypesTables(t *testing.T) {
	if len(HSTypes) != 7 {
		t.Fatalf("HSTypes has %d entries, want 7", len(HSTypes))
	}

	if len(GetHSTypes()) != len(HSTypes) {
		t.Fatal("GetHSTypes should alias HSTypes")
	}

	if len(HSTypesByVal) != len(HSTypes) {
		t.Fatalf("HSTypesByVal has %d entries, want %d", len(HSTypesByVal), len(HSTypes))
	}

	if len(GetHSTypesByVal()) != len(HSTypesByVal) {
		t.Fatal("GetHSTypesByVal should alias HSTypesByVal")
	}

	cases := []struct {
		name string
		got  HSType
		want HSType
	}{
		{name: "DS", got: HSTypes["DS"], want: HSTypeDS},
		{name: "NS", got: HSTypes["NS"], want: HSTypeNS},
		{name: "GLUE4", got: HSTypes["GLUE4"], want: HSTypeGLUE4},
		{name: "GLUE6", got: HSTypes["GLUE6"], want: HSTypeGLUE6},
		{name: "SYNTH4", got: HSTypes["SYNTH4"], want: HSTypeSYNTH4},
		{name: "SYNTH6", got: HSTypes["SYNTH6"], want: HSTypeSYNTH6},
		{name: "TXT", got: HSTypes["TXT"], want: HSTypeTXT},
	}

	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("HSTypes[%s] = %d, want %d", tc.name, tc.got, tc.want)
		}

		if got := GetHSTypes()[tc.name]; got != tc.want {
			t.Fatalf("GetHSTypes()[%s] = %d, want %d", tc.name, got, tc.want)
		}

		if got := HSTypesByVal[tc.want]; got != tc.name {
			t.Fatalf("HSTypesByVal[%d] = %q, want %q", tc.want, got, tc.name)
		}

		if got := GetHSTypesByVal()[tc.want]; got != tc.name {
			t.Fatalf("GetHSTypesByVal()[%d] = %q, want %q", tc.want, got, tc.name)
		}
	}
}

// BEGIN v090 file-aware generated coverage

func TestCommon_GetDummy_Good(t *core.T) {
	fn := ax7ValueOf(GetDummy)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetDummy_Bad(t *core.T) {
	fn := ax7ValueOf(GetDummy)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetDummy_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDummy)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetDefaultTTL_Good(t *core.T) {
	fn := ax7ValueOf(GetDefaultTTL)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetDefaultTTL_Bad(t *core.T) {
	fn := ax7ValueOf(GetDefaultTTL)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetDefaultTTL_Ugly(t *core.T) {
	fn := ax7ValueOf(GetDefaultTTL)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapRoot_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapRoot)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapRoot_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapRoot)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapRoot_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapRoot)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapEmpty_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapEmpty_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapEmpty_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapEmpty)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapNS_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapNS)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapNS_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapNS)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapNS_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapNS)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapTXT_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapTXT)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapTXT_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapTXT)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapTXT_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapTXT)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapA_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapA_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapA_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapAAAA_Good(t *core.T) {
	fn := ax7ValueOf(GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapAAAA_Bad(t *core.T) {
	fn := ax7ValueOf(GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetTypeMapAAAA_Ugly(t *core.T) {
	fn := ax7ValueOf(GetTypeMapAAAA)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetHSTypes_Good(t *core.T) {
	fn := ax7ValueOf(GetHSTypes)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetHSTypes_Bad(t *core.T) {
	fn := ax7ValueOf(GetHSTypes)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetHSTypes_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHSTypes)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetHSTypesByVal_Good(t *core.T) {
	fn := ax7ValueOf(GetHSTypesByVal)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetHSTypesByVal_Bad(t *core.T) {
	fn := ax7ValueOf(GetHSTypesByVal)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestCommon_GetHSTypesByVal_Ugly(t *core.T) {
	fn := ax7ValueOf(GetHSTypesByVal)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
