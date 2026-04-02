// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"testing"

	core "dappco.re/go/core"
)

func TestCovenantItemHelpers(t *testing.T) {
	var cov Covenant

	cov.PushU8(7).
		PushU32(0x11223344).
		PushHash(Hash{1, 2, 3}).
		PushString("foo")

	if got := cov.Len(); got != 4 {
		t.Fatalf("Len() = %d, want 4", got)
	}

	if got := cov.IndexOf([]byte("foo")); got != 3 {
		t.Fatalf("IndexOf() = %d, want 3", got)
	}

	u8, err := cov.GetU8(0)
	if err != nil {
		t.Fatalf("GetU8() returned error: %v", err)
	}
	if u8 != 7 {
		t.Fatalf("GetU8() = %d, want 7", u8)
	}

	u32, err := cov.GetU32(1)
	if err != nil {
		t.Fatalf("GetU32() returned error: %v", err)
	}
	if u32 != 0x11223344 {
		t.Fatalf("GetU32() = %#x, want %#x", u32, 0x11223344)
	}

	hash, err := cov.GetHash(2)
	if err != nil {
		t.Fatalf("GetHash() returned error: %v", err)
	}
	if hash[0] != 1 || hash[1] != 2 || hash[2] != 3 {
		t.Fatalf("GetHash() = %x, want 010203...", hash[:3])
	}

	str, err := cov.GetString(-1)
	if err != nil {
		t.Fatalf("GetString() returned error: %v", err)
	}
	if str != "foo" {
		t.Fatalf("GetString() = %q, want %q", str, "foo")
	}

	if got := cov.GetSize(); got != 44 {
		t.Fatalf("GetSize() = %d, want 44", got)
	}

	if got := cov.GetVarSize(); got != 46 {
		t.Fatalf("GetVarSize() = %d, want 46", got)
	}

	if err := cov.Set(-1, []byte("bar")); err != nil {
		t.Fatalf("Set() returned error: %v", err)
	}

	str, err = cov.GetString(3)
	if err != nil {
		t.Fatalf("GetString() after Set returned error: %v", err)
	}
	if str != "bar" {
		t.Fatalf("GetString() after Set = %q, want %q", str, "bar")
	}

	clone := cov.Clone()
	if clone.Type != cov.Type {
		t.Fatalf("Clone() copied wrong type: %d, want %d", clone.Type, cov.Type)
	}
	if clone.Len() != cov.Len() {
		t.Fatalf("Clone() copied wrong item count: %d, want %d", clone.Len(), cov.Len())
	}
}

func TestCovenantItemErrors(t *testing.T) {
	var cov Covenant

	if _, err := cov.Get(0); err == nil {
		t.Fatal("Get() should reject an out-of-range index")
	} else if got := core.Operation(err); got != "primitives.Covenant.Get" {
		t.Fatalf("Get() error operation = %q, want %q", got, "primitives.Covenant.Get")
	}

	if err := cov.Set(1, []byte("x")); err == nil {
		t.Fatal("Set() should reject an out-of-range index")
	} else if got := core.Operation(err); got != "primitives.Covenant.Set" {
		t.Fatalf("Set() error operation = %q, want %q", got, "primitives.Covenant.Set")
	}

	cov.Push([]byte{1, 2})

	if _, err := cov.GetU8(0); err == nil {
		t.Fatal("GetU8() should reject a non-uint8 item")
	}

	if _, err := cov.GetU32(0); err == nil {
		t.Fatal("GetU32() should reject a non-uint32 item")
	}

	if _, err := cov.GetHash(0); err == nil {
		t.Fatal("GetHash() should reject a non-hash item")
	}

	cov.Push(make([]byte, 64))

	if _, err := cov.GetString(1); err == nil {
		t.Fatal("GetString() should reject a non-string item")
	}
}
