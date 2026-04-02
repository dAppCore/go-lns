// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"crypto/sha3"
	"testing"
)

func TestVerifyString(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{name: "example", ok: false},
		{name: "invalid", ok: false},
		{name: "localhost", ok: false},
		{name: "test", ok: false},
		{name: "", ok: false},
		{name: "foo", ok: true},
		{name: "foo-bar", ok: true},
		{name: "foo_bar", ok: true},
		{name: "-foo", ok: false},
		{name: "foo-", ok: false},
		{name: "_foo", ok: false},
		{name: "foo_", ok: false},
		{name: "Foo", ok: false},
		{name: "f\x7fo", ok: false},
	}

	for _, tc := range cases {
		if got := VerifyString(tc.name); got != tc.ok {
			t.Fatalf("VerifyString(%q) = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestVerifyBinary(t *testing.T) {
	cases := []struct {
		name []byte
		ok   bool
	}{
		{name: nil, ok: false},
		{name: []byte("example"), ok: false},
		{name: []byte("foo"), ok: true},
		{name: []byte("foo-bar"), ok: true},
		{name: []byte("foo_bar"), ok: true},
		{name: []byte("-foo"), ok: false},
		{name: []byte("foo-"), ok: false},
		{name: []byte("Foo"), ok: false},
		{name: []byte{0xff}, ok: false},
	}

	for _, tc := range cases {
		if got := VerifyBinary(tc.name); got != tc.ok {
			t.Fatalf("VerifyBinary(%q) = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestHashString(t *testing.T) {
	got, err := HashString("foo-bar")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("HashString returned %x, want %x", got, want)
	}
}

func TestHashBinary(t *testing.T) {
	got, err := HashBinary([]byte("foo_bar"))
	if err != nil {
		t.Fatalf("HashBinary returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo_bar"))
	if got != want {
		t.Fatalf("HashBinary returned %x, want %x", got, want)
	}
}

func TestHashRejectsInvalidName(t *testing.T) {
	if _, err := HashString("Foo"); err == nil {
		t.Fatal("HashString should reject invalid names")
	}

	if _, err := HashBinary([]byte("foo-")); err == nil {
		t.Fatal("HashBinary should reject invalid names")
	}
}
