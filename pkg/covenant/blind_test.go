// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"encoding/binary"
	"testing"

	"golang.org/x/crypto/blake2b"

	"dappco.re/go/lns/pkg/primitives"
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
