// SPDX-License-Identifier: EUPL-1.2

package blockstore

import (
	"bytes"
	"errors"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/store"

	"dappco.re/go/lns/pkg/primitives"
)

// newTestStore builds a block store over a real in-memory go-store. It uses a
// :memory: SQLite database, so the test exercises the genuine store path with no
// mocks; the store is closed when the test ends.
func newTestStore(t *testing.T) *Store {
	t.Helper()

	storeResult := store.New(":memory:")
	if err := storeResult.Err(); err != nil {
		t.Fatalf("store.New(:memory:): %v", err)
	}
	storeInstance := storeResult.Value.(*store.Store)
	t.Cleanup(func() { _ = storeInstance.Close() })

	scoped := store.NewScoped(storeInstance, "lns-blocks-test")
	if scoped == nil {
		t.Fatal("store.NewScoped returned nil")
	}

	bs, err := New(scoped)
	if err != nil {
		t.Fatalf("blockstore.New: %v", err)
	}
	return bs
}

// hashOf builds a deterministic non-zero test hash from a single seed byte.
func hashOf(seed byte) primitives.Hash {
	var h primitives.Hash
	for i := range h {
		h[i] = seed
	}
	return h
}

// TestBlockstore_PutBlock_Good checks a stored block round-trips by hash and by
// height, and that the height recorded alongside the block is returned.
func TestBlockstore_PutBlock_Good(t *testing.T) {
	bs := newTestStore(t)
	hash := hashOf(0x11)
	raw := []byte("the quick brown fox block payload")

	if err := bs.PutBlock(hash, 42, raw); err != nil {
		t.Fatalf("PutBlock: %v", err)
	}

	got, err := bs.GetBlock(hash)
	if err != nil {
		t.Fatalf("GetBlock: %v", err)
	}
	if !bytes.Equal(got, raw) {
		t.Fatalf("GetBlock = %q, want %q", got, raw)
	}

	gotByHeight, err := bs.GetBlockByHeight(42)
	if err != nil {
		t.Fatalf("GetBlockByHeight: %v", err)
	}
	if !bytes.Equal(gotByHeight, raw) {
		t.Fatalf("GetBlockByHeight = %q, want %q", gotByHeight, raw)
	}

	gotRaw, gotHeight, err := bs.GetBlockWithHeight(hash)
	if err != nil {
		t.Fatalf("GetBlockWithHeight: %v", err)
	}
	if gotHeight != 42 || !bytes.Equal(gotRaw, raw) {
		t.Fatalf("GetBlockWithHeight = (%q, %d), want (%q, 42)", gotRaw, gotHeight, raw)
	}
}

// TestBlockstore_PutBlock_Bad checks the empty-block and oversized guards reject
// the write rather than persisting something HasBlock would then report.
func TestBlockstore_PutBlock_Bad(t *testing.T) {
	bs := newTestStore(t)
	hash := hashOf(0x22)

	if err := bs.PutBlock(hash, 1, nil); !errors.Is(err, ErrEmptyBlock) {
		t.Fatalf("PutBlock(nil) error = %v, want ErrEmptyBlock", err)
	}
	if err := bs.PutBlock(hash, 1, []byte{}); !errors.Is(err, ErrEmptyBlock) {
		t.Fatalf("PutBlock(empty) error = %v, want ErrEmptyBlock", err)
	}

	oversize := make([]byte, maxBlockBytes+1)
	if err := bs.PutBlock(hash, 1, oversize); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("PutBlock(oversize) error = %v, want ErrCorrupt", err)
	}

	// A rejected write must not have created any state.
	if has, _ := bs.HasBlock(hash); has {
		t.Fatal("HasBlock true after a rejected PutBlock")
	}
}

// TestBlockstore_PutBlock_Ugly checks a nil-receiver store and a nil internal
// store are reported, not panicked, on the write path.
func TestBlockstore_PutBlock_Ugly(t *testing.T) {
	var nilStore *Store
	if err := nilStore.PutBlock(hashOf(1), 1, []byte("x")); !errors.Is(err, ErrNilStore) {
		t.Fatalf("nil receiver PutBlock = %v, want ErrNilStore", err)
	}

	empty := &Store{}
	if err := empty.PutBlock(hashOf(1), 1, []byte("x")); !errors.Is(err, ErrNilStore) {
		t.Fatalf("empty store PutBlock = %v, want ErrNilStore", err)
	}
}

// TestBlockstore_GetBlock_Good checks two distinct blocks stored at different
// heights are kept apart and each resolves independently by hash and height.
func TestBlockstore_GetBlock_Good(t *testing.T) {
	bs := newTestStore(t)
	a, b := hashOf(0xa1), hashOf(0xb2)

	if err := bs.PutBlock(a, 10, []byte("block-a")); err != nil {
		t.Fatalf("PutBlock a: %v", err)
	}
	if err := bs.PutBlock(b, 11, []byte("block-b")); err != nil {
		t.Fatalf("PutBlock b: %v", err)
	}

	if got, _ := bs.GetBlock(a); string(got) != "block-a" {
		t.Fatalf("GetBlock a = %q", got)
	}
	if got, _ := bs.GetBlockByHeight(11); string(got) != "block-b" {
		t.Fatalf("GetBlockByHeight 11 = %q", got)
	}
	if hash, _ := bs.HashAtHeight(10); hash != a {
		t.Fatalf("HashAtHeight 10 = %x, want %x", hash, a)
	}
}

// TestBlockstore_GetBlock_Bad checks reads for an absent hash and an absent
// height both report ErrNotFound, distinct from a store error.
func TestBlockstore_GetBlock_Bad(t *testing.T) {
	bs := newTestStore(t)

	if _, err := bs.GetBlock(hashOf(0xde)); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetBlock(absent) = %v, want ErrNotFound", err)
	}
	if _, err := bs.GetBlockByHeight(9999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetBlockByHeight(absent) = %v, want ErrNotFound", err)
	}
	if _, err := bs.HashAtHeight(9999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("HashAtHeight(absent) = %v, want ErrNotFound", err)
	}
	if has, err := bs.HasBlock(hashOf(0xde)); err != nil || has {
		t.Fatalf("HasBlock(absent) = (%v, %v), want (false, nil)", has, err)
	}
}

// TestBlockstore_GetBlock_Ugly checks a block record whose stored value is
// corrupt is rejected with ErrCorrupt instead of returning garbage bytes. The
// corruption is injected by writing a deliberately malformed value under the
// exact key GetBlock will read, through the same scoped store.
func TestBlockstore_GetBlock_Ugly(t *testing.T) {
	storeResult := store.New(":memory:")
	if err := storeResult.Err(); err != nil {
		t.Fatalf("store.New: %v", err)
	}
	storeInstance := storeResult.Value.(*store.Store)
	t.Cleanup(func() { _ = storeInstance.Close() })
	scoped := store.NewScoped(storeInstance, "corrupt-ns")
	bs, err := New(scoped)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	hash := hashOf(0x77)

	// Not valid hex: HexDecode fails -> ErrCorrupt.
	if err := scoped.SetIn(groupBlock, hashKey(hash), "zznothex").Err(); err != nil {
		t.Fatalf("seed non-hex: %v", err)
	}
	if _, err := bs.GetBlock(hash); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("GetBlock(non-hex) = %v, want ErrCorrupt", err)
	}

	// Valid hex but a record header shorter than 8 bytes: decodeRecord fails.
	if err := scoped.SetIn(groupBlock, hashKey(hash), core.HexEncode([]byte{0x01, 0x02})).Err(); err != nil {
		t.Fatalf("seed short record: %v", err)
	}
	if _, err := bs.GetBlock(hash); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("GetBlock(short record) = %v, want ErrCorrupt", err)
	}
}

// TestBlockstore_decodeRecord_Good checks a record encoded by encodeRecord
// decodes back to the same height and bytes.
func TestBlockstore_decodeRecord_Good(t *testing.T) {
	raw := []byte("payload-bytes-of-a-block")
	rec, err := decodeRecord(encodeRecord(blockRecord{height: 7, raw: raw}))
	if err != nil {
		t.Fatalf("decodeRecord: %v", err)
	}
	if rec.height != 7 || !bytes.Equal(rec.raw, raw) {
		t.Fatalf("decodeRecord = (%d, %q), want (7, %q)", rec.height, rec.raw, raw)
	}
}

// TestBlockstore_decodeRecord_Bad checks a record whose declared length does not
// match the trailing bytes is rejected, so truncation or padding is caught.
func TestBlockstore_decodeRecord_Bad(t *testing.T) {
	good := encodeRecord(blockRecord{height: 1, raw: []byte("abcdef")})

	// Drop the last body byte: declared length now exceeds the trailing bytes.
	truncated := good[:len(good)-1]
	if _, err := decodeRecord(truncated); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("decodeRecord(truncated) = %v, want ErrCorrupt", err)
	}

	// Append a junk byte: trailing bytes now exceed the declared length.
	padded := append(append([]byte{}, good...), 0xff)
	if _, err := decodeRecord(padded); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("decodeRecord(padded) = %v, want ErrCorrupt", err)
	}
}

// TestBlockstore_decodeRecord_Ugly checks the bounded-decode guard: a record
// whose 8-byte header declares a length beyond maxBlockBytes is rejected before
// the body is even examined, so a hostile length prefix cannot drive a large
// allocation. The body is left tiny on purpose.
func TestBlockstore_decodeRecord_Ugly(t *testing.T) {
	var header [8]byte
	// height = 0, length = maxBlockBytes + 1, with a 1-byte body. Mask each byte
	// so the constant expression does not overflow byte at compile time.
	const declared = uint32(maxBlockBytes + 1)
	header[4] = byte(declared & 0xff)
	header[5] = byte((declared >> 8) & 0xff)
	header[6] = byte((declared >> 16) & 0xff)
	header[7] = byte((declared >> 24) & 0xff)
	hostile := append(header[:], 0x00)

	if _, err := decodeRecord(hostile); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("decodeRecord(hostile length) = %v, want ErrCorrupt", err)
	}

	// A header shorter than 8 bytes is also rejected.
	if _, err := decodeRecord([]byte{0x00, 0x01, 0x02}); !errors.Is(err, ErrCorrupt) {
		t.Fatalf("decodeRecord(short) = %v, want ErrCorrupt", err)
	}
}

// TestBlockstore_New_Bad checks constructing over a nil scoped store fails with
// ErrNilStore rather than yielding a store that panics on first use.
func TestBlockstore_New_Bad(t *testing.T) {
	if _, err := New(nil); !errors.Is(err, ErrNilStore) {
		t.Fatalf("New(nil) = %v, want ErrNilStore", err)
	}
}

// TestBlockstore_Overwrite_Good checks re-storing the same hash overwrites the
// previous height and bytes, which is what re-validating a block at a corrected
// height needs.
func TestBlockstore_Overwrite_Good(t *testing.T) {
	bs := newTestStore(t)
	hash := hashOf(0x33)

	if err := bs.PutBlock(hash, 5, []byte("first")); err != nil {
		t.Fatalf("PutBlock first: %v", err)
	}
	if err := bs.PutBlock(hash, 6, []byte("second-longer")); err != nil {
		t.Fatalf("PutBlock second: %v", err)
	}

	raw, height, err := bs.GetBlockWithHeight(hash)
	if err != nil {
		t.Fatalf("GetBlockWithHeight: %v", err)
	}
	if height != 6 || string(raw) != "second-longer" {
		t.Fatalf("after overwrite = (%q, %d), want (second-longer, 6)", raw, height)
	}
}
