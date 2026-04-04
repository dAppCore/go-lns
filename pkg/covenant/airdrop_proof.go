// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"encoding/binary"

	core "dappco.re/go/core"
)

const (
	airdropReward      = 4246994314
	airdropLeaves      = 216199
	airdropDepth       = 18
	airdropSubdepth    = 3
	airdropSubleaves   = 8
	faucetDepth        = 11
	faucetLeaves       = 1358
	maxAirdropProofLen = 3400
	airdropKeyAddress  = 4
)

type airdropProof struct {
	index     uint32
	proof     [][]byte
	subindex  uint8
	subproof  [][]byte
	key       []byte
	version   uint8
	address   []byte
	fee       uint64
	signature []byte
}

type airdropKey struct {
	version uint8
	address []byte
	value   uint64
}

func decodeAirdropProof(data []byte) (*airdropProof, error) {
	if len(data) > maxAirdropProofLen {
		return nil, core.E("covenant.decodeAirdropProof", "proof too large", nil)
	}

	proof := &airdropProof{}
	var n int
	var err error

	if proof.index, n, err = readProofU32(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode index", err)
	}
	data = data[n:]

	var count uint8
	if count, n, err = readProofU8(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode proof length", err)
	}
	data = data[n:]

	if int(count) > airdropDepth {
		return nil, core.E("covenant.decodeAirdropProof", "invalid proof depth", nil)
	}

	proof.proof = make([][]byte, 0, count)
	for i := 0; i < int(count); i++ {
		var hash []byte
		if hash, n, err = readProofBytes(data, 32); err != nil {
			return nil, core.E("covenant.decodeAirdropProof", "failed to decode proof hash", err)
		}
		data = data[n:]
		proof.proof = append(proof.proof, hash)
	}

	if proof.subindex, n, err = readProofU8(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode subindex", err)
	}
	data = data[n:]

	var total uint8
	if total, n, err = readProofU8(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode subproof length", err)
	}
	data = data[n:]

	if int(total) > airdropSubdepth {
		return nil, core.E("covenant.decodeAirdropProof", "invalid subproof depth", nil)
	}

	proof.subproof = make([][]byte, 0, total)
	for i := 0; i < int(total); i++ {
		var hash []byte
		if hash, n, err = readProofBytes(data, 32); err != nil {
			return nil, core.E("covenant.decodeAirdropProof", "failed to decode subproof hash", err)
		}
		data = data[n:]
		proof.subproof = append(proof.subproof, hash)
	}

	if proof.key, n, err = readProofVarBytes(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode key", err)
	}
	data = data[n:]

	if len(proof.key) == 0 {
		return nil, core.E("covenant.decodeAirdropProof", "missing key", nil)
	}

	if proof.version, n, err = readProofU8(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode version", err)
	}
	data = data[n:]

	var size uint8
	if size, n, err = readProofU8(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode address length", err)
	}
	data = data[n:]

	if size < 2 || size > 40 {
		return nil, core.E("covenant.decodeAirdropProof", "invalid address length", nil)
	}

	if proof.address, n, err = readProofBytes(data, int(size)); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode address", err)
	}
	data = data[n:]

	if proof.fee, n, err = readProofVarint(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode fee", err)
	}
	data = data[n:]

	if proof.signature, n, err = readProofVarBytes(data); err != nil {
		return nil, core.E("covenant.decodeAirdropProof", "failed to decode signature", err)
	}
	data = data[n:]

	if len(data) != 0 {
		return nil, core.E("covenant.decodeAirdropProof", "trailing data", nil)
	}

	return proof, nil
}

func (p *airdropProof) isAddress() bool {
	return len(p.key) > 0 && p.key[0] == airdropKeyAddress
}

func (p *airdropProof) getValue() uint64 {
	if !p.isAddress() {
		return airdropReward
	}

	key, ok := p.getKey()
	if !ok {
		return 0
	}

	return key.value
}

func (p *airdropProof) getKey() (airdropKey, bool) {
	if !p.isAddress() {
		return airdropKey{}, false
	}

	if len(p.key) < 1+1+1+8+1 {
		return airdropKey{}, false
	}

	key := airdropKey{}
	key.version = p.key[1]

	addrLen := int(p.key[2])
	if addrLen < 2 || addrLen > 40 {
		return airdropKey{}, false
	}

	want := 1 + 1 + 1 + addrLen + 8 + 1
	if len(p.key) != want {
		return airdropKey{}, false
	}

	key.address = append([]byte(nil), p.key[3:3+addrLen]...)
	key.value = binary.LittleEndian.Uint64(p.key[3+addrLen : 3+addrLen+8])
	return key, true
}

func (p *airdropProof) isSane() bool {
	if len(p.key) == 0 {
		return false
	}

	if p.version > 31 {
		return false
	}

	if len(p.address) < 2 || len(p.address) > 40 {
		return false
	}

	if p.isAddress() {
		if _, ok := p.getKey(); !ok {
			return false
		}
	}

	value := p.getValue()

	if p.fee > value {
		return false
	}

	if p.isAddress() {
		if len(p.subproof) != 0 {
			return false
		}

		if p.subindex != 0 {
			return false
		}

		if len(p.proof) > faucetDepth {
			return false
		}

		if p.index >= faucetLeaves {
			return false
		}

		return true
	}

	if len(p.subproof) > airdropSubdepth {
		return false
	}

	if p.subindex >= airdropSubleaves {
		return false
	}

	if len(p.proof) > airdropDepth {
		return false
	}

	if p.index >= airdropLeaves {
		return false
	}

	if p.size() > maxAirdropProofLen {
		return false
	}

	return true
}

func (p *airdropProof) size() int {
	size := 0
	size += 4
	size += 1
	size += len(p.proof) * 32
	size += 1
	size += 1
	size += len(p.subproof) * 32
	size += sizeProofVarBytes(len(p.key))
	size += 1
	size += 1
	size += len(p.address)
	size += sizeProofVarint(p.fee)
	size += sizeProofVarBytes(len(p.signature))
	return size
}

func readProofU8(src []byte) (uint8, int, error) {
	if len(src) < 1 {
		return 0, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
	}

	return src[0], 1, nil
}

func readProofU32(src []byte) (uint32, int, error) {
	if len(src) < 4 {
		return 0, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
	}

	return binary.LittleEndian.Uint32(src[:4]), 4, nil
}

func readProofBytes(src []byte, n int) ([]byte, int, error) {
	if n < 0 || len(src) < n {
		return nil, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
	}

	buf := make([]byte, n)
	copy(buf, src[:n])
	return buf, n, nil
}

func readProofVarint(src []byte) (uint64, int, error) {
	if len(src) == 0 {
		return 0, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
	}

	switch prefix := src[0]; {
	case prefix < 0xfd:
		return uint64(prefix), 1, nil
	case prefix == 0xfd:
		if len(src) < 3 {
			return 0, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
		}
		return uint64(binary.LittleEndian.Uint16(src[1:3])), 3, nil
	case prefix == 0xfe:
		if len(src) < 5 {
			return 0, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
		}
		return uint64(binary.LittleEndian.Uint32(src[1:5])), 5, nil
	default:
		if len(src) < 9 {
			return 0, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
		}
		return binary.LittleEndian.Uint64(src[1:9]), 9, nil
	}
}

func readProofVarBytes(src []byte) ([]byte, int, error) {
	n, consumed, err := readProofVarint(src)
	if err != nil {
		return nil, 0, err
	}

	if uint64(len(src[consumed:])) < n {
		return nil, 0, core.E("covenant.decodeAirdropProof", "short buffer", nil)
	}

	end := consumed + int(n)
	buf := make([]byte, int(n))
	copy(buf, src[consumed:end])
	return buf, end, nil
}

func sizeProofVarint(n uint64) int {
	switch {
	case n < 0xfd:
		return 1
	case n <= 0xffff:
		return 3
	case n <= 0xffffffff:
		return 5
	default:
		return 9
	}
}

func sizeProofVarBytes(n int) int {
	return sizeProofVarint(uint64(n)) + n
}
