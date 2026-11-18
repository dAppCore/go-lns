// SPDX-License-Identifier: EUPL-1.2

// Package script implements the Lethean Name System transaction script system,
// ported from the HSD (Handshake) reference implementation at
// itns-sidechain/lib/script/ (mirrored upstream at handshake-org/hsd). It
// provides the opcode set, the Script encode/decode layer, the witness stack
// interpreter, and the witness-program verification entry point that full
// transaction validation calls per input.
//
// # Model
//
// The chain is segregated-witness only: an output commits to a witness
// programme (primitives.Address — a version plus a 20- or 32-byte hash) rather
// than a script, and the spending input carries a witness stack
// (primitives.Input.Witness, a [][]byte). There is no scriptSig and no
// scriptPubKey. Verify reconstructs the redeem script from the address:
//
//   - Version 0, 20-byte hash: pay-to-witness-pubkey-hash (P2WPKH). The redeem
//     script is the canonical OP_DUP OP_BLAKE160 <hash> OP_EQUALVERIFY OP_CHECKSIG
//     template and the witness stack must be exactly [signature, pubkey].
//   - Version 0, 32-byte hash: pay-to-witness-script-hash (P2WSH). The last
//     witness item is the redeem script; its SHA3-256 digest must equal the
//     address hash. The remaining items are the input stack.
//   - Version 31: nulldata — unspendable (OP_RETURN error).
//   - Other versions: reserved upgradable witness programmes.
//
// # Hashing
//
// Handshake replaces the Bitcoin hash functions: P2PKH templates use BLAKE2b-160
// (OP_BLAKE160), script hashes use SHA3-256 (OP_SHA3, and Script.Sha3 for the
// P2WSH commitment). OP_KECCAK is the pre-standard Keccak-256. The legacy
// OP_HASH160/OP_HASH256/OP_SHA256/OP_SHA1/OP_RIPEMD160 opcodes remain available
// for custom redeem scripts. Signatures are 64-byte compact secp256k1 (R‖S, low-S)
// with a trailing 1-byte sighash type, and the signature hash is the BLAKE2b
// segwit-style commitment in sighash.go.
//
// # Example
//
//	redeem := script.NewScript().FromPubkeyhash(pubkeyHash160)
//	raw := redeem.Encode()
//	parsed, _ := script.DecodeScript(raw)
package script

import (
	"encoding/binary"

	core "github.com/dappcore/go"
)

// Opcode values for the Lethean name-chain script system. The set mirrors HSD
// (docs/js-script via lib/script/common.js): the Bitcoin-derived opcodes plus
// the Handshake crypto additions (OP_BLAKE160, OP_BLAKE256, OP_SHA3, OP_KECCAK)
// and OP_TYPE.
//
//	op := script.OP_CHECKSIG
const (
	// Push value opcodes.
	OP_0 = 0x00

	OP_PUSHDATA1 = 0x4c
	OP_PUSHDATA2 = 0x4d
	OP_PUSHDATA4 = 0x4e

	OP_1NEGATE = 0x4f

	OP_RESERVED = 0x50

	OP_1  = 0x51
	OP_2  = 0x52
	OP_3  = 0x53
	OP_4  = 0x54
	OP_5  = 0x55
	OP_6  = 0x56
	OP_7  = 0x57
	OP_8  = 0x58
	OP_9  = 0x59
	OP_10 = 0x5a
	OP_11 = 0x5b
	OP_12 = 0x5c
	OP_13 = 0x5d
	OP_14 = 0x5e
	OP_15 = 0x5f
	OP_16 = 0x60

	// Control opcodes.
	OP_NOP      = 0x61
	OP_VER      = 0x62
	OP_IF       = 0x63
	OP_NOTIF    = 0x64
	OP_VERIF    = 0x65
	OP_VERNOTIF = 0x66
	OP_ELSE     = 0x67
	OP_ENDIF    = 0x68
	OP_VERIFY   = 0x69
	OP_RETURN   = 0x6a

	// Stack opcodes.
	OP_TOALTSTACK   = 0x6b
	OP_FROMALTSTACK = 0x6c
	OP_2DROP        = 0x6d
	OP_2DUP         = 0x6e
	OP_3DUP         = 0x6f
	OP_2OVER        = 0x70
	OP_2ROT         = 0x71
	OP_2SWAP        = 0x72
	OP_IFDUP        = 0x73
	OP_DEPTH        = 0x74
	OP_DROP         = 0x75
	OP_DUP          = 0x76
	OP_NIP          = 0x77
	OP_OVER         = 0x78
	OP_PICK         = 0x79
	OP_ROLL         = 0x7a
	OP_ROT          = 0x7b
	OP_SWAP         = 0x7c
	OP_TUCK         = 0x7d

	// Splice opcodes (all disabled in consensus).
	OP_CAT    = 0x7e
	OP_SUBSTR = 0x7f
	OP_LEFT   = 0x80
	OP_RIGHT  = 0x81
	OP_SIZE   = 0x82

	// Bitwise opcodes.
	OP_INVERT      = 0x83
	OP_AND         = 0x84
	OP_OR          = 0x85
	OP_XOR         = 0x86
	OP_EQUAL       = 0x87
	OP_EQUALVERIFY = 0x88
	OP_RESERVED1   = 0x89
	OP_RESERVED2   = 0x8a

	// Numeric opcodes.
	OP_1ADD               = 0x8b
	OP_1SUB               = 0x8c
	OP_2MUL               = 0x8d
	OP_2DIV               = 0x8e
	OP_NEGATE             = 0x8f
	OP_ABS                = 0x90
	OP_NOT                = 0x91
	OP_0NOTEQUAL          = 0x92
	OP_ADD                = 0x93
	OP_SUB                = 0x94
	OP_MUL                = 0x95
	OP_DIV                = 0x96
	OP_MOD                = 0x97
	OP_LSHIFT             = 0x98
	OP_RSHIFT             = 0x99
	OP_BOOLAND            = 0x9a
	OP_BOOLOR             = 0x9b
	OP_NUMEQUAL           = 0x9c
	OP_NUMEQUALVERIFY     = 0x9d
	OP_NUMNOTEQUAL        = 0x9e
	OP_LESSTHAN           = 0x9f
	OP_GREATERTHAN        = 0xa0
	OP_LESSTHANOREQUAL    = 0xa1
	OP_GREATERTHANOREQUAL = 0xa2
	OP_MIN                = 0xa3
	OP_MAX                = 0xa4
	OP_WITHIN             = 0xa5

	// Crypto opcodes.
	OP_RIPEMD160           = 0xa6
	OP_SHA1                = 0xa7
	OP_SHA256              = 0xa8
	OP_HASH160             = 0xa9
	OP_HASH256             = 0xaa
	OP_CODESEPARATOR       = 0xab
	OP_CHECKSIG            = 0xac
	OP_CHECKSIGVERIFY      = 0xad
	OP_CHECKMULTISIG       = 0xae
	OP_CHECKMULTISIGVERIFY = 0xaf

	// Expansion opcodes.
	OP_NOP1                = 0xb0
	OP_CHECKLOCKTIMEVERIFY = 0xb1
	OP_CHECKSEQUENCEVERIFY = 0xb2
	OP_NOP4                = 0xb3
	OP_NOP5                = 0xb4
	OP_NOP6                = 0xb5
	OP_NOP7                = 0xb6
	OP_NOP8                = 0xb7
	OP_NOP9                = 0xb8
	OP_NOP10               = 0xb9

	// Handshake crypto additions.
	OP_BLAKE160 = 0xc0
	OP_BLAKE256 = 0xc1
	OP_SHA3     = 0xc2
	OP_KECCAK   = 0xc3

	// Introspection.
	OP_TYPE = 0xd0

	// Sentinel.
	OP_INVALIDOPCODE = 0xff
)

// opParseError is the sentinel opcode value used to mark a push-data opcode that
// could not be decoded (truncated length or body). It mirrors the JS PARSE_ERROR
// opcode whose value is -1; in Go it surfaces as a value that round-trips
// through Opcode.Value and is rejected by the interpreter as BAD_OPCODE.
const opParseError = -1

// opcodeNames maps opcode values to their canonical symbolic names. It is used
// for diagnostics (Opcode.Symbol, ScriptError messages).
var opcodeNames = map[int]string{
	OP_0: "OP_0", OP_PUSHDATA1: "OP_PUSHDATA1", OP_PUSHDATA2: "OP_PUSHDATA2",
	OP_PUSHDATA4: "OP_PUSHDATA4", OP_1NEGATE: "OP_1NEGATE", OP_RESERVED: "OP_RESERVED",
	OP_1: "OP_1", OP_2: "OP_2", OP_3: "OP_3", OP_4: "OP_4", OP_5: "OP_5",
	OP_6: "OP_6", OP_7: "OP_7", OP_8: "OP_8", OP_9: "OP_9", OP_10: "OP_10",
	OP_11: "OP_11", OP_12: "OP_12", OP_13: "OP_13", OP_14: "OP_14", OP_15: "OP_15",
	OP_16: "OP_16", OP_NOP: "OP_NOP", OP_VER: "OP_VER", OP_IF: "OP_IF",
	OP_NOTIF: "OP_NOTIF", OP_VERIF: "OP_VERIF", OP_VERNOTIF: "OP_VERNOTIF",
	OP_ELSE: "OP_ELSE", OP_ENDIF: "OP_ENDIF", OP_VERIFY: "OP_VERIFY",
	OP_RETURN: "OP_RETURN", OP_TOALTSTACK: "OP_TOALTSTACK",
	OP_FROMALTSTACK: "OP_FROMALTSTACK", OP_2DROP: "OP_2DROP", OP_2DUP: "OP_2DUP",
	OP_3DUP: "OP_3DUP", OP_2OVER: "OP_2OVER", OP_2ROT: "OP_2ROT",
	OP_2SWAP: "OP_2SWAP", OP_IFDUP: "OP_IFDUP", OP_DEPTH: "OP_DEPTH",
	OP_DROP: "OP_DROP", OP_DUP: "OP_DUP", OP_NIP: "OP_NIP", OP_OVER: "OP_OVER",
	OP_PICK: "OP_PICK", OP_ROLL: "OP_ROLL", OP_ROT: "OP_ROT", OP_SWAP: "OP_SWAP",
	OP_TUCK: "OP_TUCK", OP_CAT: "OP_CAT", OP_SUBSTR: "OP_SUBSTR", OP_LEFT: "OP_LEFT",
	OP_RIGHT: "OP_RIGHT", OP_SIZE: "OP_SIZE", OP_INVERT: "OP_INVERT", OP_AND: "OP_AND",
	OP_OR: "OP_OR", OP_XOR: "OP_XOR", OP_EQUAL: "OP_EQUAL",
	OP_EQUALVERIFY: "OP_EQUALVERIFY", OP_RESERVED1: "OP_RESERVED1",
	OP_RESERVED2: "OP_RESERVED2", OP_1ADD: "OP_1ADD", OP_1SUB: "OP_1SUB",
	OP_2MUL: "OP_2MUL", OP_2DIV: "OP_2DIV", OP_NEGATE: "OP_NEGATE", OP_ABS: "OP_ABS",
	OP_NOT: "OP_NOT", OP_0NOTEQUAL: "OP_0NOTEQUAL", OP_ADD: "OP_ADD", OP_SUB: "OP_SUB",
	OP_MUL: "OP_MUL", OP_DIV: "OP_DIV", OP_MOD: "OP_MOD", OP_LSHIFT: "OP_LSHIFT",
	OP_RSHIFT: "OP_RSHIFT", OP_BOOLAND: "OP_BOOLAND", OP_BOOLOR: "OP_BOOLOR",
	OP_NUMEQUAL: "OP_NUMEQUAL", OP_NUMEQUALVERIFY: "OP_NUMEQUALVERIFY",
	OP_NUMNOTEQUAL: "OP_NUMNOTEQUAL", OP_LESSTHAN: "OP_LESSTHAN",
	OP_GREATERTHAN: "OP_GREATERTHAN", OP_LESSTHANOREQUAL: "OP_LESSTHANOREQUAL",
	OP_GREATERTHANOREQUAL: "OP_GREATERTHANOREQUAL", OP_MIN: "OP_MIN", OP_MAX: "OP_MAX",
	OP_WITHIN: "OP_WITHIN", OP_RIPEMD160: "OP_RIPEMD160", OP_SHA1: "OP_SHA1",
	OP_SHA256: "OP_SHA256", OP_HASH160: "OP_HASH160", OP_HASH256: "OP_HASH256",
	OP_CODESEPARATOR: "OP_CODESEPARATOR", OP_CHECKSIG: "OP_CHECKSIG",
	OP_CHECKSIGVERIFY: "OP_CHECKSIGVERIFY", OP_CHECKMULTISIG: "OP_CHECKMULTISIG",
	OP_CHECKMULTISIGVERIFY: "OP_CHECKMULTISIGVERIFY", OP_NOP1: "OP_NOP1",
	OP_CHECKLOCKTIMEVERIFY: "OP_CHECKLOCKTIMEVERIFY",
	OP_CHECKSEQUENCEVERIFY: "OP_CHECKSEQUENCEVERIFY", OP_NOP4: "OP_NOP4",
	OP_NOP5: "OP_NOP5", OP_NOP6: "OP_NOP6", OP_NOP7: "OP_NOP7", OP_NOP8: "OP_NOP8",
	OP_NOP9: "OP_NOP9", OP_NOP10: "OP_NOP10", OP_BLAKE160: "OP_BLAKE160",
	OP_BLAKE256: "OP_BLAKE256", OP_SHA3: "OP_SHA3", OP_KECCAK: "OP_KECCAK",
	OP_TYPE: "OP_TYPE", OP_INVALIDOPCODE: "OP_INVALIDOPCODE",
}

// OpcodeName returns the canonical symbolic name for an opcode value, or an
// "OP_UNKNOWN(0xNN)" form for unrecognised values.
//
//	name := script.OpcodeName(script.OP_CHECKSIG) // "OP_CHECKSIG"
func OpcodeName(value int) string {
	if value == opParseError {
		return "OP_INVALIDOPCODE"
	}
	if name, ok := opcodeNames[value]; ok {
		return name
	}
	return core.Sprintf("OP_UNKNOWN(0x%02x)", value&0xff)
}

// Opcode is a single parsed script element: an opcode value and, for push
// operations, the pushed data. A pure opcode (no push) has Data == nil.
//
//	op := script.OpFromData([]byte{0x01, 0x02})
//	_ = op.Value // OP_PUSHDATA byte or the direct length 0x00..0x4b
type Opcode struct {
	// Value is the opcode byte (0x00..0xff) or opParseError (-1) for an opcode
	// that failed to decode.
	Value int

	// Data is the pushed bytes for a push opcode, or nil for a non-push opcode.
	Data []byte
}

// OpFromOp returns a bare opcode with no push data.
//
//	op := script.OpFromOp(script.OP_DUP)
func OpFromOp(value int) Opcode {
	return Opcode{Value: value}
}

// OpFromData returns the minimal push opcode for the given data, encoding small
// integers (0, 1..16, -1) as their single-byte opcodes the way HSD's
// Opcode.fromData does.
//
//	op := script.OpFromData([]byte{0x05}) // encodes as OP_5
func OpFromData(data []byte) Opcode {
	if len(data) == 1 {
		if data[0] == 0x81 {
			return OpFromOp(OP_1NEGATE)
		}
		if data[0] >= 1 && data[0] <= 16 {
			return OpFromOp(int(data[0]) + 0x50)
		}
	}
	return OpFromPush(data)
}

// OpFromPush returns a push opcode for the given data, always using a pushdata
// form (never collapsing to a small-integer opcode). Mirrors Opcode.fromPush.
//
//	op := script.OpFromPush([]byte{0x05}) // direct push of length 1
func OpFromPush(data []byte) Opcode {
	switch {
	case len(data) == 0:
		return OpFromOp(OP_0)
	case len(data) <= 0x4b:
		return Opcode{Value: len(data), Data: append([]byte(nil), data...)}
	case len(data) <= 0xff:
		return Opcode{Value: OP_PUSHDATA1, Data: append([]byte(nil), data...)}
	case len(data) <= 0xffff:
		return Opcode{Value: OP_PUSHDATA2, Data: append([]byte(nil), data...)}
	default:
		return Opcode{Value: OP_PUSHDATA4, Data: append([]byte(nil), data...)}
	}
}

// OpFromSmall returns the single-byte opcode for a small integer 0..16.
//
//	op := script.OpFromSmall(3) // OP_3
func OpFromSmall(num int) Opcode {
	if num <= 0 {
		return OpFromOp(OP_0)
	}
	if num > 16 {
		num = 16
	}
	return OpFromOp(num + 0x50)
}

// IsParseError reports whether the opcode is the decode sentinel.
func (o Opcode) IsParseError() bool {
	return o.Value == opParseError
}

// IsBranch reports whether the opcode is a conditional (OP_IF..OP_ENDIF).
func (o Opcode) IsBranch() bool {
	return o.Value >= OP_IF && o.Value <= OP_ENDIF
}

// IsDisabled reports whether the opcode is permanently disabled by consensus.
// Disabled opcodes fail a script even inside an unexecuted branch.
func (o Opcode) IsDisabled() bool {
	switch o.Value {
	case OP_CAT, OP_SUBSTR, OP_LEFT, OP_RIGHT,
		OP_INVERT, OP_AND, OP_OR, OP_XOR,
		OP_2MUL, OP_2DIV, OP_MUL, OP_DIV, OP_MOD,
		OP_LSHIFT, OP_RSHIFT:
		return true
	}
	return false
}

// IsMinimal reports whether a push opcode uses the smallest possible encoding
// for its data (the MINIMALDATA rule). Non-push opcodes are always minimal.
func (o Opcode) IsMinimal() bool {
	if o.Data == nil {
		return true
	}
	if len(o.Data) == 1 {
		if o.Data[0] == 0x81 {
			return false
		}
		if o.Data[0] >= 1 && o.Data[0] <= 16 {
			return false
		}
	}
	if len(o.Data) <= 0x4b {
		return o.Value == len(o.Data)
	}
	if len(o.Data) <= 0xff {
		return o.Value == OP_PUSHDATA1
	}
	if len(o.Data) <= 0xffff {
		return o.Value == OP_PUSHDATA2
	}
	return o.Value == OP_PUSHDATA4
}

// ToSmall returns the integer value of a small-integer opcode (0, 1..16) or -1
// if the opcode is not one. Mirrors Opcode.toSmall.
func (o Opcode) ToSmall() int {
	if o.Value == OP_0 {
		return 0
	}
	if o.Value >= OP_1 && o.Value <= OP_16 {
		return o.Value - 0x50
	}
	return -1
}

// ToPush returns the data this opcode pushes onto the stack: the push bytes for
// a pushdata op, or the small-integer encoding for OP_0/OP_1NEGATE/OP_1..OP_16.
// Returns nil for opcodes that do not push. Mirrors Opcode.toPush.
func (o Opcode) ToPush() []byte {
	if o.Value == OP_0 {
		return []byte{}
	}
	if o.Value == OP_1NEGATE {
		return []byte{0x81}
	}
	if o.Value >= OP_1 && o.Value <= OP_16 {
		return []byte{byte(o.Value - 0x50)}
	}
	return o.Data
}

// Symbol returns the opcode's symbolic name (see OpcodeName).
func (o Opcode) Symbol() string {
	return OpcodeName(o.Value)
}

// Size returns the serialised size of the opcode in bytes.
func (o Opcode) Size() int {
	if o.Data == nil {
		return 1
	}
	switch o.Value {
	case OP_PUSHDATA1:
		return 2 + len(o.Data)
	case OP_PUSHDATA2:
		return 3 + len(o.Data)
	case OP_PUSHDATA4:
		return 5 + len(o.Data)
	default:
		return 1 + len(o.Data)
	}
}

// appendTo writes the opcode's wire encoding onto dst and returns the extended
// slice. A parse-error opcode cannot be serialised and returns an error.
func (o Opcode) appendTo(dst []byte) ([]byte, error) {
	if o.Value == opParseError {
		return nil, core.E("script.Opcode.appendTo", "cannot reserialise a parse error", nil)
	}
	if o.Data == nil {
		return append(dst, byte(o.Value)), nil
	}
	switch o.Value {
	case OP_PUSHDATA1:
		dst = append(dst, OP_PUSHDATA1, byte(len(o.Data)))
	case OP_PUSHDATA2:
		var tmp [2]byte
		binary.LittleEndian.PutUint16(tmp[:], uint16(len(o.Data)))
		dst = append(dst, OP_PUSHDATA2, tmp[0], tmp[1])
	case OP_PUSHDATA4:
		var tmp [4]byte
		binary.LittleEndian.PutUint32(tmp[:], uint32(len(o.Data)))
		dst = append(dst, OP_PUSHDATA4, tmp[0], tmp[1], tmp[2], tmp[3])
	default:
		// Direct push: the opcode value is the data length.
		dst = append(dst, byte(o.Value))
	}
	return append(dst, o.Data...), nil
}

// readOpcode decodes a single opcode from src and returns it with the number of
// bytes consumed. A truncated push length or body yields a parse-error opcode
// that consumes the remaining bytes (matching the JS reader's seek-to-end
// behaviour), so the caller still makes forward progress and the interpreter
// rejects the script. The push length is taken from the encoding itself and is
// always bounded by the remaining buffer before any slice is taken, so an
// untrusted length cannot drive an out-of-range read or oversized allocation.
func readOpcode(src []byte) (Opcode, int) {
	if len(src) == 0 {
		return Opcode{Value: opParseError}, 0
	}

	value := int(src[0])
	rest := src[1:]

	switch {
	case value == OP_PUSHDATA1:
		if len(rest) < 1 {
			return Opcode{Value: opParseError}, len(src)
		}
		size := int(rest[0])
		body := rest[1:]
		if len(body) < size {
			return Opcode{Value: opParseError}, len(src)
		}
		return Opcode{Value: value, Data: append([]byte(nil), body[:size]...)}, 1 + 1 + size

	case value == OP_PUSHDATA2:
		if len(rest) < 2 {
			return Opcode{Value: opParseError}, len(src)
		}
		size := int(binary.LittleEndian.Uint16(rest[:2]))
		body := rest[2:]
		if len(body) < size {
			return Opcode{Value: opParseError}, len(src)
		}
		return Opcode{Value: value, Data: append([]byte(nil), body[:size]...)}, 1 + 2 + size

	case value == OP_PUSHDATA4:
		if len(rest) < 4 {
			return Opcode{Value: opParseError}, len(src)
		}
		size := int(binary.LittleEndian.Uint32(rest[:4]))
		body := rest[4:]
		// size is bounded by the buffer below; on 32-bit a >2GiB length would be
		// negative after the int conversion, which the bound check rejects.
		if size < 0 || len(body) < size {
			return Opcode{Value: opParseError}, len(src)
		}
		return Opcode{Value: value, Data: append([]byte(nil), body[:size]...)}, 1 + 4 + size

	case value >= 0x01 && value <= 0x4b:
		// Direct push: the opcode value is the data length.
		if len(rest) < value {
			return Opcode{Value: opParseError}, len(src)
		}
		return Opcode{Value: value, Data: append([]byte(nil), rest[:value]...)}, 1 + value

	default:
		return Opcode{Value: value}, 1
	}
}
