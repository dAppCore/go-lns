// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"bytes"
	"testing"
)

func TestCommonConstants(t *testing.T) {
	if DEFAULT_TTL != 21600 {
		t.Fatalf("DEFAULT_TTL = %d, want 21600", DEFAULT_TTL)
	}

	if len(DUMMY) != 0 {
		t.Fatalf("DUMMY length = %d, want 0", len(DUMMY))
	}

	cases := []struct {
		name string
		got  []byte
		want []byte
	}{
		{name: "TYPE_MAP_ROOT", got: TYPE_MAP_ROOT, want: mustHex("000722000000000380")},
		{name: "TYPE_MAP_EMPTY", got: TYPE_MAP_EMPTY, want: mustHex("0006000000000003")},
		{name: "TYPE_MAP_NS", got: TYPE_MAP_NS, want: mustHex("0006200000000003")},
		{name: "TYPE_MAP_TXT", got: TYPE_MAP_TXT, want: mustHex("0006000080000003")},
		{name: "TYPE_MAP_A", got: TYPE_MAP_A, want: mustHex("0006400000000003")},
		{name: "TYPE_MAP_AAAA", got: TYPE_MAP_AAAA, want: mustHex("0006000000080003")},
	}

	for _, tc := range cases {
		if !bytes.Equal(tc.got, tc.want) {
			t.Fatalf("%s = %x, want %x", tc.name, tc.got, tc.want)
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
