// SPDX-License-Identifier: EUPL-1.2

package script

import core "github.com/dappcore/go"

// scriptNumMaxBytes is the default maximum length (in bytes) of a number read
// from the stack by arithmetic opcodes. It matches the HSD/Bitcoin 4-byte limit
// (numbers up to ±2^31-1), keeping arithmetic within int32 range.
const scriptNumMaxBytes = 4

// ScriptNum is a signed integer in the script number encoding: little-endian
// magnitude with a sign bit in the most-significant byte, zero encoded as the
// empty slice. It mirrors HSD's ScriptNum (docs/js-script via
// lib/script/scriptnum.js) which is a 64-bit value; arithmetic opcodes bound
// inputs to four bytes via decodeScriptNum.
//
//	n, _ := script.DecodeScriptNum([]byte{0x05}, true, 4)
//	_ = n.Int() // 5
type ScriptNum int64

// Int returns the value as an int (saturating to int32 bounds the way
// ScriptNum.getInt does, so callers using it for counts/indices stay in range).
func (n ScriptNum) Int() int {
	const (
		min = -(1 << 31)
		max = (1 << 31) - 1
	)
	switch {
	case int64(n) < min:
		return min
	case int64(n) > max:
		return max
	default:
		return int(n)
	}
}

// Encode serialises the number into the minimal script-number form.
//
//	raw := script.ScriptNum(-1).Encode() // []byte{0x81}
func (n ScriptNum) Encode() []byte {
	if n == 0 {
		return []byte{}
	}

	neg := n < 0
	abs := uint64(n)
	if neg {
		abs = uint64(-n)
	}

	var out []byte
	for abs > 0 {
		out = append(out, byte(abs&0xff))
		abs >>= 8
	}

	// If the most-significant byte already has its high bit set, append a
	// sign-only byte; otherwise fold the sign into the top byte.
	if out[len(out)-1]&0x80 != 0 {
		if neg {
			out = append(out, 0x80)
		} else {
			out = append(out, 0x00)
		}
	} else if neg {
		out[len(out)-1] |= 0x80
	}

	return out
}

// IsMinimalScriptNum reports whether data is the most compact encoding of its
// value. Mirrors ScriptNum.isMinimal: a trailing zero-magnitude byte is only
// permitted when needed to carry the sign bit.
//
//	ok := script.IsMinimalScriptNum([]byte{0x05}) // true
func IsMinimalScriptNum(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	if data[len(data)-1]&0x7f == 0 {
		if len(data) == 1 {
			return false
		}
		if data[len(data)-2]&0x80 == 0 {
			return false
		}
	}
	return true
}

// DecodeScriptNum decodes a script number, optionally enforcing minimal encoding
// and a maximum byte length. It mirrors ScriptNum.decode and is the entry point
// the interpreter uses for arithmetic operands.
//
//	n, err := script.DecodeScriptNum([]byte{0x05}, true, 4)
func DecodeScriptNum(data []byte, minimal bool, limit int) (ScriptNum, error) {
	if limit > 0 && len(data) > limit {
		return 0, core.E("script.DecodeScriptNum", "script number overflow", nil)
	}
	if minimal && !IsMinimalScriptNum(data) {
		return 0, core.E("script.DecodeScriptNum", "non-minimal script number", nil)
	}
	return decodeScriptNumUnsafe(data), nil
}

// decodeScriptNumUnsafe decodes without length/minimal checks. The caller bounds
// the input length; values are limited to 8 magnitude bytes so the shift stays
// within int64.
func decodeScriptNumUnsafe(data []byte) ScriptNum {
	if len(data) == 0 {
		return 0
	}

	var result int64
	n := min(len(data), 8)
	for i := 0; i < n; i++ {
		result |= int64(data[i]) << (8 * i)
	}

	// If the sign bit is set, clear it and negate.
	if data[len(data)-1]&0x80 != 0 {
		mask := int64(0x80) << (8 * (len(data) - 1))
		result &^= mask
		return ScriptNum(-result)
	}

	return ScriptNum(result)
}

// castToBool reports whether a stack item is truthy. A value is false if it is
// empty or all zero bytes except for an optional negative-zero sign byte.
// Mirrors Stack.getBool / the bitcoind CastToBool rule.
//
//	ok := script.castToBool([]byte{0x01}) // true
func castToBool(data []byte) bool {
	for i := range data {
		if data[i] != 0 {
			// Negative zero (0x80 in the top byte, rest zero) is still false.
			if i == len(data)-1 && data[i] == 0x80 {
				return false
			}
			return true
		}
	}
	return false
}
