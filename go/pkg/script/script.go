// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"crypto/sha3"
	"encoding/binary"

	core "github.com/dappcore/go"
)

// Consensus limits for the script system. They mirror the HSD protocol
// constants (lib/protocol/consensus.js) consumed by the interpreter.
const (
	// MaxScriptSize is the maximum serialised length of a single script.
	MaxScriptSize = 10000
	// MaxScriptStack is the maximum combined size of the main and alt stacks.
	MaxScriptStack = 1000
	// MaxScriptPush is the maximum length of a pushed stack item.
	MaxScriptPush = 520
	// MaxScriptOps is the maximum number of executed non-push opcodes.
	MaxScriptOps = 201
	// MaxMultisigPubkeys is the maximum key count accepted by OP_CHECKMULTISIG.
	MaxMultisigPubkeys = 20
)

// pubkeyLen is the length of a compressed secp256k1 public key.
const pubkeyLen = 33

// sigLen is the length of a Handshake script signature: a 64-byte compact
// secp256k1 signature plus a trailing one-byte sighash type.
const sigLen = 65

// Script is a parsed transaction script: the raw bytes plus the decoded opcode
// sequence. On the name chain a Script is a witness redeem script (P2WSH) or the
// templated P2WPKH/multisig form a witness programme expands to. It mirrors
// HSD's Script (docs/js-script via lib/script/script.js).
//
//	s, err := script.DecodeScript(raw)
//	if s.IsPubkeyhash(false) { /* P2PKH redeem template */ }
type Script struct {
	// raw is the canonical serialised form; it is the source of truth for
	// hashing and re-encoding.
	raw []byte
	// code is the decoded opcode sequence.
	code []Opcode
}

// NewScript returns an empty script.
func NewScript() *Script {
	return &Script{raw: []byte{}, code: []Opcode{}}
}

// Code returns the decoded opcode sequence. The returned slice aliases the
// script; callers must not mutate it.
func (s *Script) Code() []Opcode {
	return s.code
}

// Encode returns the raw serialised script. The returned slice aliases the
// script; callers must not mutate it.
//
//	raw := redeem.Encode()
func (s *Script) Encode() []byte {
	return s.raw
}

// GetSize returns the serialised length of the script in bytes.
func (s *Script) GetSize() int {
	return len(s.raw)
}

// GetVarSize returns the length of the script when written with a varint length
// prefix (the form used inside the signature-hash preimage and keyrings).
func (s *Script) GetVarSize() int {
	return sizeVarint(len(s.raw)) + len(s.raw)
}

// Sha3 returns the SHA3-256 digest of the raw script. This is the witness
// programme hash for a P2WSH address (matching KeyRing.getScriptHash, which uses
// script.sha3()).
//
//	hash := redeem.Sha3()
func (s *Script) Sha3() []byte {
	sum := sha3.Sum256(s.raw)
	return sum[:]
}

// DecodeScript parses a raw script into its opcode sequence. A script that
// exceeds MaxScriptSize is rejected. Truncated push operations decode to
// parse-error opcodes (rejected at execution) rather than failing the decode, so
// witness data round-trips byte-for-byte. Push lengths are always bounded by the
// remaining buffer before any slice is taken.
//
//	s, err := script.DecodeScript(raw)
func DecodeScript(raw []byte) (*Script, error) {
	if len(raw) > MaxScriptSize {
		return nil, core.E("script.DecodeScript", "script exceeds maximum size", nil)
	}

	s := &Script{raw: append([]byte(nil), raw...)}
	rest := s.raw
	for len(rest) > 0 {
		op, n := readOpcode(rest)
		s.code = append(s.code, op)
		if n <= 0 {
			break
		}
		rest = rest[n:]
	}
	if s.code == nil {
		s.code = []Opcode{}
	}
	return s, nil
}

// FromOps builds a script from a sequence of opcodes, recomputing the raw form.
//
//	s, err := script.FromOps([]script.Opcode{script.OpFromOp(script.OP_1)})
func FromOps(ops []Opcode) (*Script, error) {
	s := &Script{code: append([]Opcode(nil), ops...)}
	if err := s.compile(); err != nil {
		return nil, err
	}
	return s, nil
}

// compile (re)serialises the opcode sequence into the raw form.
func (s *Script) compile() error {
	raw := make([]byte, 0, 32)
	for _, op := range s.code {
		var err error
		raw, err = op.appendTo(raw)
		if err != nil {
			return core.E("script.Script.compile", "encode opcode", err)
		}
	}
	if len(raw) > MaxScriptSize {
		return core.E("script.Script.compile", "script exceeds maximum size", nil)
	}
	s.raw = raw
	return nil
}

// FromPubkeyhash builds the standard pay-to-pubkey-hash redeem template for a
// 20-byte BLAKE2b-160 key hash:
//
//	OP_DUP OP_BLAKE160 <hash> OP_EQUALVERIFY OP_CHECKSIG
//
// This is the redeem script a version-0 / 20-byte witness programme expands to.
// Mirrors Script.fromPubkeyhash.
//
//	redeem := script.NewScript().FromPubkeyhash(hash160)
func (s *Script) FromPubkeyhash(hash []byte) *Script {
	raw := make([]byte, 25)
	raw[0] = OP_DUP
	raw[1] = OP_BLAKE160
	raw[2] = 0x14
	copy(raw[3:23], hash)
	raw[23] = OP_EQUALVERIFY
	raw[24] = OP_CHECKSIG

	s.raw = raw
	s.code = []Opcode{
		OpFromOp(OP_DUP),
		OpFromOp(OP_BLAKE160),
		OpFromPush(raw[3:23]),
		OpFromOp(OP_EQUALVERIFY),
		OpFromOp(OP_CHECKSIG),
	}
	return s
}

// FromPubkeyhash builds a pay-to-pubkey-hash redeem script. It validates the
// hash length (20 bytes).
//
//	redeem, err := script.FromPubkeyhash(hash160)
func FromPubkeyhash(hash []byte) (*Script, error) {
	if len(hash) != 20 {
		return nil, core.E("script.FromPubkeyhash", "pubkey hash must be 20 bytes", nil)
	}
	return NewScript().FromPubkeyhash(hash), nil
}

// FromPubkey builds a bare pay-to-pubkey redeem script: <pubkey> OP_CHECKSIG.
//
//	redeem, err := script.FromPubkey(compressedKey)
func FromPubkey(key []byte) (*Script, error) {
	if len(key) != pubkeyLen {
		return nil, core.E("script.FromPubkey", "public key must be 33 bytes", nil)
	}
	return FromOps([]Opcode{OpFromPush(key), OpFromOp(OP_CHECKSIG)})
}

// FromMultisig builds an m-of-n pay-to-multisig redeem script:
//
//	OP_<m> <key_1> ... <key_n> OP_<n> OP_CHECKMULTISIG
//
// Keys are sorted lexicographically (matching Script.fromMultisig). Each key
// must be a 33-byte compressed public key.
//
//	redeem, err := script.FromMultisig(2, 3, keys)
func FromMultisig(m, n int, keys [][]byte) (*Script, error) {
	if m < 1 || n < 1 || n > 15 || m > n {
		return nil, core.E("script.FromMultisig", "invalid m-of-n parameters", nil)
	}
	if len(keys) != n {
		return nil, core.E("script.FromMultisig", "key count does not match n", nil)
	}

	sorted := make([][]byte, len(keys))
	for i, k := range keys {
		if len(k) != pubkeyLen {
			return nil, core.E("script.FromMultisig", "public key must be 33 bytes", nil)
		}
		sorted[i] = k
	}
	sortKeys(sorted)

	ops := make([]Opcode, 0, n+3)
	ops = append(ops, OpFromSmall(m))
	for _, k := range sorted {
		ops = append(ops, OpFromPush(k))
	}
	ops = append(ops, OpFromSmall(n), OpFromOp(OP_CHECKMULTISIG))
	return FromOps(ops)
}

// IsPubkey reports whether the script is a bare pay-to-pubkey template
// (<33-byte key> OP_CHECKSIG). Mirrors Script.isPubkey.
func (s *Script) IsPubkey() bool {
	return len(s.raw) == 35 &&
		s.raw[0] == 33 &&
		s.raw[34] == OP_CHECKSIG
}

// IsPubkeyhash reports whether the script is the canonical P2PKH template.
// Mirrors Script.isPubkeyhash.
func (s *Script) IsPubkeyhash() bool {
	return len(s.raw) == 25 &&
		s.raw[0] == OP_DUP &&
		s.raw[1] == OP_BLAKE160 &&
		s.raw[2] == 0x14 &&
		s.raw[23] == OP_EQUALVERIFY &&
		s.raw[24] == OP_CHECKSIG
}

// GetPubkeyhash returns the 20-byte key hash if the script is a P2PKH template,
// or nil otherwise.
func (s *Script) GetPubkeyhash() []byte {
	if !s.IsPubkeyhash() {
		return nil
	}
	return append([]byte(nil), s.raw[3:23]...)
}

// IsMultisig reports whether the script is an m-of-n multisig template. Mirrors
// Script.isMultisig.
func (s *Script) IsMultisig() bool {
	if len(s.code) < 4 || len(s.code) > 19 {
		return false
	}
	if s.code[len(s.code)-1].Value != OP_CHECKMULTISIG {
		return false
	}
	m := s.code[0].ToSmall()
	if m < 1 {
		return false
	}
	n := s.code[len(s.code)-2].ToSmall()
	if n < 1 || m > n {
		return false
	}
	if len(s.code) != n+3 {
		return false
	}
	for i := 1; i < n+1; i++ {
		if len(s.code[i].Data) != pubkeyLen {
			return false
		}
	}
	return true
}

// GetSubscript returns a copy of the script from the given opcode index to the
// end, with OP_CODESEPARATOR opcodes removed. It is the code the signature hash
// commits to (mirrors Script.getSubscript + removeSeparators).
func (s *Script) GetSubscript(index int) (*Script, error) {
	if index < 0 {
		index = 0
	}
	if index > len(s.code) {
		index = len(s.code)
	}
	ops := make([]Opcode, 0, len(s.code)-index)
	for i := index; i < len(s.code); i++ {
		if s.code[i].Value == OP_CODESEPARATOR {
			continue
		}
		ops = append(ops, s.code[i])
	}
	return FromOps(ops)
}

/*
 * varint helpers — the Bitcoin/HSD CompactSize encoding, matching
 * primitives' (unexported) writer so witness/script wire forms interoperate.
 */

func sizeVarint(n int) int {
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

func appendVarint(dst []byte, n uint64) []byte {
	switch {
	case n < 0xfd:
		return append(dst, byte(n))
	case n <= 0xffff:
		return append(dst, 0xfd, byte(n), byte(n>>8))
	case n <= 0xffffffff:
		return append(dst, 0xfe, byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
	default:
		return append(dst,
			0xff,
			byte(n), byte(n>>8), byte(n>>16), byte(n>>24),
			byte(n>>32), byte(n>>40), byte(n>>48), byte(n>>56))
	}
}

func appendVarBytes(dst, src []byte) []byte {
	dst = appendVarint(dst, uint64(len(src)))
	return append(dst, src...)
}

func appendU32(dst []byte, v uint32) []byte {
	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], v)
	return append(dst, tmp[:]...)
}

func appendU64(dst []byte, v uint64) []byte {
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], v)
	return append(dst, tmp[:]...)
}

// sortKeys sorts public keys lexicographically in place (the ordering
// Script.fromMultisig applies).
func sortKeys(keys [][]byte) {
	// Insertion sort keeps the dependency surface minimal and n<=15 for multisig.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && compareBytes(keys[j-1], keys[j]) > 0; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
}

func compareBytes(left, right []byte) int {
	n := min(len(right), len(left))
	for i := 0; i < n; i++ {
		switch {
		case left[i] < right[i]:
			return -1
		case left[i] > right[i]:
			return 1
		}
	}
	switch {
	case len(left) < len(right):
		return -1
	case len(left) > len(right):
		return 1
	default:
		return 0
	}
}

func equalBytes(left, right []byte) bool {
	return compareBytes(left, right) == 0
}
