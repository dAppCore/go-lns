// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"crypto/sha3"
	"testing"

	"dappco.re/go/lns/pkg/covenant"
)

func TestServiceResolve(t *testing.T) {
	svc := &Service{}

	got, err := svc.Resolve("Foo-Bar.lthn")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}

	want := sha3.Sum256([]byte("foo-bar"))
	if got != want {
		t.Fatalf("Resolve returned %x, want %x", got, want)
	}

	got, err = svc.Resolve([]byte("Foo-Bar.lthn"))
	if err != nil {
		t.Fatalf("Resolve([]byte) returned error: %v", err)
	}

	if got != want {
		t.Fatalf("Resolve([]byte) returned %x, want %x", got, want)
	}
}

func TestServiceResolveRejectsInvalidNames(t *testing.T) {
	svc := &Service{}

	cases := []string{
		"",
		" foo-bar.lthn",
		"foo-bar.lthn ",
		"foo.bar.lthn",
		"foo..lthn",
		"foo.lthn..",
		"foo-",
		"test.lthn",
	}

	for _, name := range cases {
		if _, err := svc.Resolve(name); err == nil {
			t.Fatalf("Resolve(%q) should reject the name", name)
		}
	}

	if _, err := svc.Resolve(123); err == nil {
		t.Fatal("Resolve should reject unsupported input types")
	}
}

func TestServiceVerify(t *testing.T) {
	svc := &Service{}

	cases := []struct {
		name any
		ok   bool
	}{
		{name: "Foo-Bar.lthn", ok: true},
		{name: []byte("Foo-Bar.lthn"), ok: true},
		{name: "foo.bar.lthn", ok: false},
		{name: "foo-", ok: false},
		{name: 123, ok: false},
	}

	for _, tc := range cases {
		if got := svc.Verify(tc.name); got != tc.ok {
			t.Fatalf("Verify(%v) = %v, want %v", tc.name, got, tc.ok)
		}
	}
}

func TestServiceGetReserved(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetReserved("RESERVED.lthn")
	if !ok {
		t.Fatal("GetReserved should find the reserved reference entry")
	}

	if item.Name != "reserved" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "reserved")
	}

	if item.Target != "reserved.com." {
		t.Fatalf("item.Target = %q, want %q", item.Target, "reserved.com.")
	}

	hash, err := covenant.HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	if item.Hash != hash {
		t.Fatalf("item.Hash = %x, want %x", item.Hash, hash)
	}

	if _, ok := svc.GetReserved([]byte("reserved.lthn")); !ok {
		t.Fatal("GetReserved([]byte) should find the reserved reference entry")
	}

	if _, ok := svc.GetReserved(123); ok {
		t.Fatal("GetReserved should reject unsupported input types")
	}
}

func TestServiceGetLocked(t *testing.T) {
	svc := &Service{}

	item, ok := svc.GetLocked("NEC.lthn")
	if !ok {
		t.Fatal("GetLocked should find the locked reference entry")
	}

	if item.Name != "nec" {
		t.Fatalf("item.Name = %q, want %q", item.Name, "nec")
	}

	if item.Target != "nec." {
		t.Fatalf("item.Target = %q, want %q", item.Target, "nec.")
	}

	if !item.Root || item.Custom {
		t.Fatalf("unexpected flags on nec: %#v", item)
	}

	if _, ok := svc.GetLocked([]byte("nec.lthn")); !ok {
		t.Fatal("GetLocked([]byte) should find the locked reference entry")
	}

	if _, ok := svc.GetLocked(123); ok {
		t.Fatal("GetLocked should reject unsupported input types")
	}
}
