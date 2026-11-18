// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"encoding/binary"
	"testing"

	core "github.com/dappcore/go"
)

// proofBuilder assembles airdrop-proof wire bytes for decode tests.
type proofBuilder struct {
	buf []byte
}

func (b *proofBuilder) u8(v uint8) *proofBuilder { b.buf = append(b.buf, v); return b }

func (b *proofBuilder) u32(v uint32) *proofBuilder {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], v)
	b.buf = append(b.buf, tmp[:]...)
	return b
}

func (b *proofBuilder) hash() *proofBuilder { b.buf = append(b.buf, make([]byte, 32)...); return b }

func (b *proofBuilder) varint(n uint64) *proofBuilder {
	switch {
	case n < 0xfd:
		b.buf = append(b.buf, byte(n))
	case n <= 0xffff:
		var t [2]byte
		binary.LittleEndian.PutUint16(t[:], uint16(n))
		b.buf = append(b.buf, 0xfd)
		b.buf = append(b.buf, t[:]...)
	case n <= 0xffffffff:
		var t [4]byte
		binary.LittleEndian.PutUint32(t[:], uint32(n))
		b.buf = append(b.buf, 0xfe)
		b.buf = append(b.buf, t[:]...)
	default:
		var t [8]byte
		binary.LittleEndian.PutUint64(t[:], n)
		b.buf = append(b.buf, 0xff)
		b.buf = append(b.buf, t[:]...)
	}
	return b
}

func (b *proofBuilder) varbytes(data []byte) *proofBuilder {
	b.varint(uint64(len(data)))
	b.buf = append(b.buf, data...)
	return b
}

func (b *proofBuilder) bytes() []byte { return b.buf }

// faucetProof builds a minimal valid faucet (non-address-key) proof.
func faucetProof() []byte {
	b := &proofBuilder{}
	b.u32(0)                          // index
	b.u8(0)                           // proof count
	b.u8(0)                           // subindex
	b.u8(0)                           // subproof total
	b.varbytes([]byte{0x01})          // key (non-address: first byte != airdropKeyAddress)
	b.u8(0)                           // version
	b.u8(2)                           // address length
	b.buf = append(b.buf, 0xaa, 0xbb) // address
	b.varint(0)                       // fee
	b.varbytes(nil)                   // signature
	return b.bytes()
}

func TestAirdropProofBehaviour_decode_Good(t *testing.T) {
	proof, err := decodeAirdropProof(faucetProof())
	core.AssertNoError(t, err)
	core.AssertNotNil(t, proof)
	core.AssertEqual(t, uint32(0), proof.index)
	core.AssertLen(t, proof.address, 2)
	core.AssertTrue(t, proof.isSane())
	// A non-address key falls back to the fixed airdrop reward value.
	core.AssertEqual(t, uint64(airdropReward), proof.getValue())
}

func TestAirdropProofBehaviour_decode_Bad(t *testing.T) {
	// Oversized buffer is rejected up front.
	big := make([]byte, maxAirdropProofLen+1)
	_, err := decodeAirdropProof(big)
	core.AssertError(t, err)

	// Truncated buffer (just the index) fails on the proof-length read.
	short := (&proofBuilder{}).u32(0).bytes()
	_, err = decodeAirdropProof(short)
	core.AssertError(t, err)

	// Empty key is explicitly rejected.
	b := &proofBuilder{}
	b.u32(0).u8(0).u8(0).u8(0).varbytes(nil)
	_, err = decodeAirdropProof(b.bytes())
	core.AssertError(t, err)
}

func TestAirdropProofBehaviour_decode_Ugly(t *testing.T) {
	// Proof depth beyond airdropDepth is rejected.
	b := &proofBuilder{}
	b.u32(0).u8(airdropDepth + 1)
	_, err := decodeAirdropProof(b.bytes())
	core.AssertError(t, err)

	// Invalid address length (1) is rejected.
	bad := &proofBuilder{}
	bad.u32(0).u8(0).u8(0).u8(0).varbytes([]byte{0x01}).u8(0).u8(1)
	bad.buf = append(bad.buf, 0xaa)
	_, err = decodeAirdropProof(bad.bytes())
	core.AssertError(t, err)

	// Trailing data after a complete proof is rejected.
	trailing := append(faucetProof(), 0xff)
	_, err = decodeAirdropProof(trailing)
	core.AssertError(t, err)
}

func TestAirdropProofBehaviour_addressKey(t *testing.T) {
	// Build a key with the address marker: [marker, version, addrLen, addr..., value(8), trailer]
	addr := []byte{0xaa, 0xbb}
	key := []byte{airdropKeyAddress, 0x00, byte(len(addr))}
	key = append(key, addr...)
	var value [8]byte
	binary.LittleEndian.PutUint64(value[:], 12345)
	key = append(key, value[:]...)
	key = append(key, 0x00) // trailer byte

	b := &proofBuilder{}
	b.u32(0).u8(0).u8(0).u8(0).varbytes(key).u8(0).u8(byte(len(addr)))
	b.buf = append(b.buf, addr...)
	b.varint(0).varbytes(nil)

	proof, err := decodeAirdropProof(b.bytes())
	core.AssertNoError(t, err)
	core.AssertTrue(t, proof.isAddress())
	core.AssertEqual(t, uint64(12345), proof.getValue())
	core.AssertTrue(t, proof.isSane())

	parsed, ok := proof.getKey()
	core.AssertTrue(t, ok)
	core.AssertEqual(t, uint64(12345), parsed.value)
}

func TestAirdropProofBehaviour_isSane_Bad(t *testing.T) {
	// Fee greater than the value is rejected.
	proof := &airdropProof{
		key:     []byte{0x01},
		version: 0,
		address: []byte{0xaa, 0xbb},
		fee:     airdropReward + 1,
	}
	core.AssertFalse(t, proof.isSane())

	// Version above 31 is rejected.
	bad := &airdropProof{key: []byte{0x01}, version: 32, address: []byte{0xaa, 0xbb}}
	core.AssertFalse(t, bad.isSane())

	// Empty key is rejected.
	core.AssertFalse(t, (&airdropProof{}).isSane())
}

func TestAirdropProofBehaviour_size_Good(t *testing.T) {
	proof := &airdropProof{
		proof:     [][]byte{make([]byte, 32)},
		subproof:  [][]byte{make([]byte, 32)},
		key:       []byte{0x01},
		address:   []byte{0xaa, 0xbb},
		signature: []byte{0x01, 0x02},
	}
	// 4(index)+1(count)+32(proof)+1(subindex)+1(total)+32(subproof)
	// +2(varbytes key:1+1)+1(version)+1(addrlen)+2(address)
	// +1(varint fee)+3(varbytes sig:1+2) = 81
	core.AssertEqual(t, 81, proof.size())
}

func TestAirdropProofBehaviour_varint_helpers(t *testing.T) {
	core.AssertEqual(t, 1, sizeProofVarint(0xfc))
	core.AssertEqual(t, 3, sizeProofVarint(0xffff))
	core.AssertEqual(t, 5, sizeProofVarint(0xffffffff))
	core.AssertEqual(t, 9, sizeProofVarint(0x1_0000_0000))

	core.AssertEqual(t, 1+5, sizeProofVarBytes(5))

	// readProofVarint covers each prefix branch.
	v, n, err := readProofVarint([]byte{0x10})
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint64(0x10), v)
	core.AssertEqual(t, 1, n)

	v, n, err = readProofVarint([]byte{0xfd, 0x34, 0x12})
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint64(0x1234), v)
	core.AssertEqual(t, 3, n)

	v, n, err = readProofVarint([]byte{0xfe, 0x78, 0x56, 0x34, 0x12})
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint64(0x12345678), v)
	core.AssertEqual(t, 5, n)

	v, n, err = readProofVarint([]byte{0xff, 1, 0, 0, 0, 0, 0, 0, 0})
	core.AssertNoError(t, err)
	core.AssertEqual(t, uint64(1), v)
	core.AssertEqual(t, 9, n)
}

func TestAirdropProofBehaviour_varint_short(t *testing.T) {
	_, _, err := readProofVarint(nil)
	core.AssertError(t, err)

	_, _, err = readProofVarint([]byte{0xfd, 0x00})
	core.AssertError(t, err)

	_, _, err = readProofVarint([]byte{0xfe, 0x00})
	core.AssertError(t, err)

	_, _, err = readProofVarint([]byte{0xff, 0x00})
	core.AssertError(t, err)
}
