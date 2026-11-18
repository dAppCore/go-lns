// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	core "github.com/dappcore/go"
	"encoding/binary"
	"testing"

	"golang.org/x/crypto/blake2b"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

func TestBlind(t *testing.T) {
	var nonce primitives.Hash
	for i := range nonce {
		nonce[i] = byte(i)
	}

	got, err := Blind(0x1122334455667788, nonce)
	if err != nil {
		t.Fatalf("Blind returned error: %v", err)
	}

	var input [40]byte
	binary.LittleEndian.PutUint64(input[:8], 0x1122334455667788)
	copy(input[8:], nonce[:])

	want := blake2b.Sum256(input[:])
	if got != want {
		t.Fatalf("Blind returned %x, want %x", got, want)
	}
}

func TestGetBlind(t *testing.T) {
	var nonce primitives.Hash
	for i := range nonce {
		nonce[i] = byte(i)
	}

	got, err := GetBlind(0x1122334455667788, nonce)
	if err != nil {
		t.Fatalf("GetBlind returned error: %v", err)
	}

	want, err := Blind(0x1122334455667788, nonce)
	if err != nil {
		t.Fatalf("Blind returned error: %v", err)
	}

	if got != want {
		t.Fatalf("GetBlind returned %x, want %x", got, want)
	}
}

// BEGIN v090 file-aware generated coverage

func TestBlind_Blind_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(Blind(1, [32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBlind_Blind_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(Blind(0, [32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBlind_Blind_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(Blind(63, [32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBlind_GetBlind_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetBlind(1, [32]byte{1}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBlind_GetBlind_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetBlind(0, [32]byte{}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestBlind_GetBlind_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetBlind(63, [32]byte{0xff}))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
