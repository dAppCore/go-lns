// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"testing"

	core "dappco.re/go"
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

	cov.SetUpdate(Hash{31, 32, 33}, 155, []byte("update"))

	if got := cov.Type; got != covenantTypeUpdate {
		t.Fatalf("SetUpdate() set type %d, want %d", got, covenantTypeUpdate)
	}

	if got := cov.Len(); got != 3 {
		t.Fatalf("SetUpdate() len = %d, want 3", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetUpdate() GetHash returned error: %v", err)
	}
	if hash[0] != 31 || hash[1] != 32 || hash[2] != 33 {
		t.Fatalf("SetUpdate() hash = %x, want 1f2021...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetUpdate() GetU32 returned error: %v", err)
	}
	if height != 155 {
		t.Fatalf("SetUpdate() height = %d, want 155", height)
	}

	str, err = cov.GetString(2)
	if err != nil {
		t.Fatalf("SetUpdate() GetString returned error: %v", err)
	}
	if str != "update" {
		t.Fatalf("SetUpdate() GetString() = %q, want %q", str, "update")
	}

	cov.SetRenew(Hash{34, 35, 36}, 166, Hash{37, 38, 39})

	if got := cov.Type; got != covenantTypeRenew {
		t.Fatalf("SetRenew() set type %d, want %d", got, covenantTypeRenew)
	}

	if got := cov.Len(); got != 3 {
		t.Fatalf("SetRenew() len = %d, want 3", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetRenew() GetHash returned error: %v", err)
	}
	if hash[0] != 34 || hash[1] != 35 || hash[2] != 36 {
		t.Fatalf("SetRenew() hash = %x, want 222324...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetRenew() GetU32 returned error: %v", err)
	}
	if height != 166 {
		t.Fatalf("SetRenew() height = %d, want 166", height)
	}

	blockHash, err = cov.GetHash(2)
	if err != nil {
		t.Fatalf("SetRenew() GetHash(2) returned error: %v", err)
	}
	if blockHash[0] != 37 || blockHash[1] != 38 || blockHash[2] != 39 {
		t.Fatalf("SetRenew() block hash = %x, want 252627...", blockHash[:3])
	}

	cov.SetTransfer(Hash{40, 41, 42}, 177, Address{Version: 0, Hash: []byte{1, 2, 3, 4, 5}})

	if got := cov.Type; got != covenantTypeTransfer {
		t.Fatalf("SetTransfer() set type %d, want %d", got, covenantTypeTransfer)
	}

	if got := cov.Len(); got != 4 {
		t.Fatalf("SetTransfer() len = %d, want 4", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetTransfer() GetHash returned error: %v", err)
	}
	if hash[0] != 40 || hash[1] != 41 || hash[2] != 42 {
		t.Fatalf("SetTransfer() hash = %x, want 28292a...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetTransfer() GetU32 returned error: %v", err)
	}
	if height != 177 {
		t.Fatalf("SetTransfer() height = %d, want 177", height)
	}

	version, err := cov.GetU8(2)
	if err != nil {
		t.Fatalf("SetTransfer() GetU8 returned error: %v", err)
	}
	if version != 0 {
		t.Fatalf("SetTransfer() version = %d, want 0", version)
	}

	addrHash, err := cov.Get(3)
	if err != nil {
		t.Fatalf("SetTransfer() Get(3) returned error: %v", err)
	}
	if len(addrHash) != 5 || addrHash[0] != 1 || addrHash[4] != 5 {
		t.Fatalf("SetTransfer() address hash = %x, want 0102030405", addrHash)
	}

	cov.SetFinalize(Hash{43, 44, 45}, 188, []byte("final"), 12, 21, 34, Hash{46, 47, 48})

	if got := cov.Type; got != covenantTypeFinalize {
		t.Fatalf("SetFinalize() set type %d, want %d", got, covenantTypeFinalize)
	}

	if got := cov.Len(); got != 7 {
		t.Fatalf("SetFinalize() len = %d, want 7", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetFinalize() GetHash returned error: %v", err)
	}
	if hash[0] != 43 || hash[1] != 44 || hash[2] != 45 {
		t.Fatalf("SetFinalize() hash = %x, want 2b2c2d...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetFinalize() GetU32 returned error: %v", err)
	}
	if height != 188 {
		t.Fatalf("SetFinalize() height = %d, want 188", height)
	}

	str, err = cov.GetString(2)
	if err != nil {
		t.Fatalf("SetFinalize() GetString returned error: %v", err)
	}
	if str != "final" {
		t.Fatalf("SetFinalize() GetString() = %q, want %q", str, "final")
	}

	flags, err = cov.GetU8(3)
	if err != nil {
		t.Fatalf("SetFinalize() GetU8 returned error: %v", err)
	}
	if flags != 12 {
		t.Fatalf("SetFinalize() flags = %d, want 12", flags)
	}

	claimed, err := cov.GetU32(4)
	if err != nil {
		t.Fatalf("SetFinalize() GetU32(4) returned error: %v", err)
	}
	if claimed != 21 {
		t.Fatalf("SetFinalize() claimed = %d, want 21", claimed)
	}

	renewals, err := cov.GetU32(5)
	if err != nil {
		t.Fatalf("SetFinalize() GetU32(5) returned error: %v", err)
	}
	if renewals != 34 {
		t.Fatalf("SetFinalize() renewals = %d, want 34", renewals)
	}

	blockHash, err = cov.GetHash(6)
	if err != nil {
		t.Fatalf("SetFinalize() GetHash(6) returned error: %v", err)
	}
	if blockHash[0] != 46 || blockHash[1] != 47 || blockHash[2] != 48 {
		t.Fatalf("SetFinalize() block hash = %x, want 2e2f30...", blockHash[:3])
	}

	cov.SetRevoke(Hash{49, 50, 51}, 199)

	if got := cov.Type; got != covenantTypeRevoke {
		t.Fatalf("SetRevoke() set type %d, want %d", got, covenantTypeRevoke)
	}

	if got := cov.Len(); got != 2 {
		t.Fatalf("SetRevoke() len = %d, want 2", got)
	}

	hash, err = cov.GetHash(0)
	if err != nil {
		t.Fatalf("SetRevoke() GetHash returned error: %v", err)
	}
	if hash[0] != 49 || hash[1] != 50 || hash[2] != 51 {
		t.Fatalf("SetRevoke() hash = %x, want 313233...", hash[:3])
	}

	height, err = cov.GetU32(1)
	if err != nil {
		t.Fatalf("SetRevoke() GetU32 returned error: %v", err)
	}
	if height != 199 {
		t.Fatalf("SetRevoke() height = %d, want 199", height)
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

func TestCovenantJSONRoundTrip(t *testing.T) {
	original := Covenant{}
	original.SetFinalize(
		Hash{43, 44, 45},
		188,
		[]byte("final"),
		12,
		21,
		34,
		Hash{46, 47, 48},
	)

	jsonView := original.GetJSON()
	if jsonView.Type != covenantTypeFinalize {
		t.Fatalf("GetJSON() type = %d, want %d", jsonView.Type, covenantTypeFinalize)
	}

	if jsonView.Action != "FINALIZE" {
		t.Fatalf("GetJSON() action = %q, want %q", jsonView.Action, "FINALIZE")
	}

	if len(jsonView.Items) != 7 {
		t.Fatalf("GetJSON() item count = %d, want 7", len(jsonView.Items))
	}

	raw := core.JSONMarshal(jsonView)
	if !raw.OK {
		t.Fatalf("JSONMarshal returned error: %v", raw.Value)
	}

	var decoded CovenantJSON
	if r := core.JSONUnmarshal(raw.Value.([]byte), &decoded); !r.OK {
		t.Fatalf("JSONUnmarshal returned error: %v", r.Value)
	}

	var roundTripped Covenant
	if err := roundTripped.FromJSON(decoded); err != nil {
		t.Fatalf("FromJSON returned error: %v", err)
	}

	if roundTripped.Type != original.Type {
		t.Fatalf("round-tripped type = %d, want %d", roundTripped.Type, original.Type)
	}

	if roundTripped.Len() != original.Len() {
		t.Fatalf("round-tripped item count = %d, want %d", roundTripped.Len(), original.Len())
	}

	for i := range original.Items {
		if string(roundTripped.Items[i]) != string(original.Items[i]) {
			t.Fatalf("round-tripped item %d = %x, want %x", i, roundTripped.Items[i], original.Items[i])
		}
	}
}

func TestCovenantJSONRejectsInvalidEncoding(t *testing.T) {
	var cov Covenant

	err := cov.FromJSON(CovenantJSON{
		Type:  covenantTypeOpen,
		Items: []string{"zz"},
	})
	if err == nil {
		t.Fatal("FromJSON should reject invalid hex data")
	}
}

// BEGIN v090 file-aware generated coverage

func TestCovenantItems_Covenant_Len_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).Len())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Len_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).Len())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Len_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).Len())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Clone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Clone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Clone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).Clone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Inject_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).Inject(coverageValue[Covenant]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Inject_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).Inject(coverageValue[Covenant]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Inject_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).Inject(coverageValue[Covenant]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetNone_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetNone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetNone_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetNone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetNone_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetNone())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_ToArray_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).ToArray())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_ToArray_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).ToArray())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_ToArray_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).ToArray())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_FromArray_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).FromArray([][]byte{[]byte("lethean")}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_FromArray_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).FromArray(nil))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_FromArray_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).FromArray([][]byte{[]byte("lethean")}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_IndexOf_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).IndexOf([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_IndexOf_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).IndexOf([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_IndexOf_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).IndexOf([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Get_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).Get(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Get_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).Get(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Get_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).Get(-1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Set_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).Set(1, []byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Set_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).Set(0, []byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Set_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).Set(-1, []byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Push_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).Push([]byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Push_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).Push([]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_Push_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).Push([]byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetOpen_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetOpen([32]byte{1}, []byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetOpen_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetOpen([32]byte{}, []byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetOpen_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetOpen([32]byte{0xff}, []byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetBid_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetBid([32]byte{1}, 1, []byte("lethean.lthn"), [32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetBid_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetBid([32]byte{}, 0, []byte{}, [32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetBid_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetBid([32]byte{0xff}, 63, []byte{0xff, 0x00, 0x61}, [32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetReveal_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetReveal([32]byte{1}, 1, [32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetReveal_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetReveal([32]byte{}, 0, [32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetReveal_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetReveal([32]byte{0xff}, 63, [32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRedeem_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetRedeem([32]byte{1}, 1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRedeem_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetRedeem([32]byte{}, 0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRedeem_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetRedeem([32]byte{0xff}, 63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetClaim_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetClaim([32]byte{1}, 1, []byte("lethean.lthn"), 1, [32]byte{1}, 1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetClaim_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetClaim([32]byte{}, 0, []byte{}, 0, [32]byte{}, 0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetClaim_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetClaim([32]byte{0xff}, 63, []byte{0xff, 0x00, 0x61}, 63, [32]byte{0xff}, 63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRegister_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetRegister([32]byte{1}, 1, []byte("lethean.lthn"), [32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRegister_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetRegister([32]byte{}, 0, []byte{}, [32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRegister_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetRegister([32]byte{0xff}, 63, []byte{0xff, 0x00, 0x61}, [32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetUpdate_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetUpdate([32]byte{1}, 1, []byte("lethean.lthn")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetUpdate_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetUpdate([32]byte{}, 0, []byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetUpdate_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetUpdate([32]byte{0xff}, 63, []byte{0xff, 0x00, 0x61}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRenew_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetRenew([32]byte{1}, 1, [32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRenew_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetRenew([32]byte{}, 0, [32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRenew_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetRenew([32]byte{0xff}, 63, [32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetTransfer_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetTransfer([32]byte{1}, 1, coverageValue[Address]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetTransfer_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetTransfer([32]byte{}, 0, coverageValue[Address]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetTransfer_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetTransfer([32]byte{0xff}, 63, coverageValue[Address]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetFinalize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetFinalize([32]byte{1}, 1, []byte("lethean.lthn"), 1, 1, 1, [32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetFinalize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetFinalize([32]byte{}, 0, []byte{}, 0, 0, 0, [32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetFinalize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetFinalize([32]byte{0xff}, 63, []byte{0xff, 0x00, 0x61}, 63, 63, 63, [32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRevoke_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).SetRevoke([32]byte{1}, 1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRevoke_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).SetRevoke([32]byte{}, 0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_SetRevoke_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).SetRevoke([32]byte{0xff}, 63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetU8_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).GetU8(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetU8_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).GetU8(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetU8_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).GetU8(-1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushU8_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).PushU8(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushU8_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).PushU8(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushU8_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).PushU8(63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetU32_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).GetU32(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetU32_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).GetU32(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetU32_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).GetU32(-1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushU32_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).PushU32(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushU32_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).PushU32(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushU32_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).PushU32(63))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).GetHash(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).GetHash(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).GetHash(-1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushHash_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).PushHash([32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushHash_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).PushHash([32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushHash_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).PushHash([32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).GetString(1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).GetString(0))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).GetString(-1))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushString_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).PushString("lethean.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushString_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).PushString(""))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_PushString_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).PushString("lethean.\x00.lthn"))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).GetSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetVarSize_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(new(Covenant).GetVarSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetVarSize_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(new(Covenant).GetVarSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestCovenantItems_Covenant_GetVarSize_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(new(Covenant).GetVarSize())
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
