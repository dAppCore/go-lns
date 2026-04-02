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

	cov.Type = 99
	cov.Items = [][]byte{{1, 2, 3}}
	cov.SetNone()
	if got := cov.Type; got != covenantTypeNone {
		t.Fatalf("SetNone() set type %d, want %d", got, covenantTypeNone)
	}

	if got := cov.Len(); got != 0 {
		t.Fatalf("SetNone() len = %d, want 0", got)
	}

	cov.SetOpen(Hash{9, 8, 7}, []byte("bar"))

	if got := cov.Type; got != covenantTypeOpen {
		t.Fatalf("SetOpen() set type %d, want %d", got, covenantTypeOpen)
	}

	if got := cov.Len(); got != 3 {
		t.Fatalf("SetOpen() len = %d, want 3", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetOpen() GetHash returned error: %v", err)
	}
	if hash[0] != 9 || hash[1] != 8 || hash[2] != 7 {
		t.Fatalf("SetOpen() hash = %x, want 090807...", hash[:3])
	}

	height, err := cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetOpen() GetU32 returned error: %v", err)
	}
	if height != 0 {
		t.Fatalf("SetOpen() height = %d, want 0", height)
	}

	str, err = cov.GetString(2)
	if err != nil {
		t.Fatalf("SetOpen() GetString returned error: %v", err)
	}
	if str != "bar" {
		t.Fatalf("SetOpen() GetString() = %q, want %q", str, "bar")
	}

	cov.SetBid(Hash{4, 5, 6}, 42, []byte("baz"), Hash{7, 8, 9})

	if got := cov.Type; got != covenantTypeBid {
		t.Fatalf("SetBid() set type %d, want %d", got, covenantTypeBid)
	}

	if got := cov.Len(); got != 4 {
		t.Fatalf("SetBid() len = %d, want 4", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetBid() GetHash returned error: %v", err)
	}
	if hash[0] != 4 || hash[1] != 5 || hash[2] != 6 {
		t.Fatalf("SetBid() hash = %x, want 040506...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetBid() GetU32 returned error: %v", err)
	}
	if height != 42 {
		t.Fatalf("SetBid() height = %d, want 42", height)
	}

	str, err = cov.GetString(2)
	if err != nil {
		t.Fatalf("SetBid() GetString returned error: %v", err)
	}
	if str != "baz" {
		t.Fatalf("SetBid() GetString() = %q, want %q", str, "baz")
	}

	blind, err := cov.GetHash(3)
	if err != nil {
		t.Fatalf("SetBid() GetHash(3) returned error: %v", err)
	}
	if blind[0] != 7 || blind[1] != 8 || blind[2] != 9 {
		t.Fatalf("SetBid() blind = %x, want 070809...", blind[:3])
	}

	cov.SetReveal(Hash{10, 11, 12}, 77, Hash{13, 14, 15})

	if got := cov.Type; got != covenantTypeReveal {
		t.Fatalf("SetReveal() set type %d, want %d", got, covenantTypeReveal)
	}

	if got := cov.Len(); got != 3 {
		t.Fatalf("SetReveal() len = %d, want 3", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetReveal() GetHash returned error: %v", err)
	}
	if hash[0] != 10 || hash[1] != 11 || hash[2] != 12 {
		t.Fatalf("SetReveal() hash = %x, want 0a0b0c...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetReveal() GetU32 returned error: %v", err)
	}
	if height != 77 {
		t.Fatalf("SetReveal() height = %d, want 77", height)
	}

	nonce, err := cov.GetHash(2)
	if err != nil {
		t.Fatalf("SetReveal() GetHash(2) returned error: %v", err)
	}
	if nonce[0] != 13 || nonce[1] != 14 || nonce[2] != 15 {
		t.Fatalf("SetReveal() nonce = %x, want 0d0e0f...", nonce[:3])
	}

	cov.SetRedeem(Hash{22, 23, 24}, 123)

	if got := cov.Type; got != covenantTypeRedeem {
		t.Fatalf("SetRedeem() set type %d, want %d", got, covenantTypeRedeem)
	}

	if got := cov.Len(); got != 2 {
		t.Fatalf("SetRedeem() len = %d, want 2", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetRedeem() GetHash returned error: %v", err)
	}
	if hash[0] != 22 || hash[1] != 23 || hash[2] != 24 {
		t.Fatalf("SetRedeem() hash = %x, want 161718...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetRedeem() GetU32 returned error: %v", err)
	}
	if height != 123 {
		t.Fatalf("SetRedeem() height = %d, want 123", height)
	}

	cov.SetClaim(Hash{16, 17, 18}, 88, []byte("claim"), 9, Hash{19, 20, 21}, 99)

	if got := cov.Type; got != covenantTypeClaim {
		t.Fatalf("SetClaim() set type %d, want %d", got, covenantTypeClaim)
	}

	if got := cov.Len(); got != 6 {
		t.Fatalf("SetClaim() len = %d, want 6", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetClaim() GetHash returned error: %v", err)
	}
	if hash[0] != 16 || hash[1] != 17 || hash[2] != 18 {
		t.Fatalf("SetClaim() hash = %x, want 101112...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetClaim() GetU32 returned error: %v", err)
	}
	if height != 88 {
		t.Fatalf("SetClaim() height = %d, want 88", height)
	}

	str, err = cov.GetString(2)
	if err != nil {
		t.Fatalf("SetClaim() GetString returned error: %v", err)
	}
	if str != "claim" {
		t.Fatalf("SetClaim() GetString() = %q, want %q", str, "claim")
	}

	flags, err := cov.GetU8(3)
	if err != nil {
		t.Fatalf("SetClaim() GetU8 returned error: %v", err)
	}
	if flags != 9 {
		t.Fatalf("SetClaim() flags = %d, want 9", flags)
	}

	commitHash, err := cov.GetHash(4)
	if err != nil {
		t.Fatalf("SetClaim() GetHash(4) returned error: %v", err)
	}
	if commitHash[0] != 19 || commitHash[1] != 20 || commitHash[2] != 21 {
		t.Fatalf("SetClaim() commit hash = %x, want 131415...", commitHash[:3])
	}

	commitHeight, err := cov.GetU32(5)
	if err != nil {
		t.Fatalf("SetClaim() GetU32(5) returned error: %v", err)
	}
	if commitHeight != 99 {
		t.Fatalf("SetClaim() commit height = %d, want 99", commitHeight)
	}

	cov.SetRegister(Hash{25, 26, 27}, 144, []byte("record"), Hash{28, 29, 30})

	if got := cov.Type; got != covenantTypeRegister {
		t.Fatalf("SetRegister() set type %d, want %d", got, covenantTypeRegister)
	}

	if got := cov.Len(); got != 4 {
		t.Fatalf("SetRegister() len = %d, want 4", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetRegister() GetHash returned error: %v", err)
	}
	if hash[0] != 25 || hash[1] != 26 || hash[2] != 27 {
		t.Fatalf("SetRegister() hash = %x, want 191a1b...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetRegister() GetU32 returned error: %v", err)
	}
	if height != 144 {
		t.Fatalf("SetRegister() height = %d, want 144", height)
	}

	str, err = cov.GetString(2)
	if err != nil {
		t.Fatalf("SetRegister() GetString returned error: %v", err)
	}
	if str != "record" {
		t.Fatalf("SetRegister() GetString() = %q, want %q", str, "record")
	}

	blockHash, err := cov.GetHash(3)
	if err != nil {
		t.Fatalf("SetRegister() GetHash(3) returned error: %v", err)
	}
	if blockHash[0] != 28 || blockHash[1] != 29 || blockHash[2] != 30 {
		t.Fatalf("SetRegister() block hash = %x, want 1c1d1e...", blockHash[:3])
	}

	cov = Covenant{}
	cov.PushU8(7).
		PushU32(0x11223344).
		PushHash(Hash{1, 2, 3}).
		PushString("foo")

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
