// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"crypto/sha3"

	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/primitives"
)

// MaxNameSize is the maximum size of a name label in bytes.
const MaxNameSize = 63

// MaxResourceSize is the maximum size of a DNS resource payload in bytes.
const MaxResourceSize = 512

const maxNameSize = MaxNameSize

var blacklist = map[string]struct{}{
	"example":   {},
	"invalid":   {},
	"local":     {},
	"localhost": {},
	"test":      {},
}

// VerifyString reports whether a domain name meets the covenant rules.
func VerifyString(name string) bool {
	if len(name) == 0 || len(name) > maxNameSize {
		return false
	}

	for i := 0; i < len(name); i++ {
		ch := name[i]

		if ch&0x80 != 0 {
			return false
		}

		switch {
		case ch >= '0' && ch <= '9':
		case ch >= 'A' && ch <= 'Z':
			return false
		case ch >= 'a' && ch <= 'z':
		case ch == '-' || ch == '_':
			if i == 0 || i == len(name)-1 {
				return false
			}
		default:
			return false
		}
	}

	_, blocked := blacklist[name]
	return !blocked
}

// VerifyBinary reports whether a binary domain name meets the covenant rules.
func VerifyBinary(name []byte) bool {
	if len(name) == 0 || len(name) > maxNameSize {
		return false
	}

	for i := 0; i < len(name); i++ {
		ch := name[i]

		if ch&0x80 != 0 {
			return false
		}

		switch {
		case ch >= '0' && ch <= '9':
		case ch >= 'A' && ch <= 'Z':
			return false
		case ch >= 'a' && ch <= 'z':
		case ch == '-' || ch == '_':
			if i == 0 || i == len(name)-1 {
				return false
			}
		default:
			return false
		}
	}

	_, blocked := blacklist[string(name)]
	return !blocked
}

// VerifyName reports whether a name meets the covenant rules.
//
// It accepts either a string or a byte slice to mirror the JS reference API.
func VerifyName(name any) bool {
	switch v := name.(type) {
	case string:
		return VerifyString(v)
	case []byte:
		return VerifyBinary(v)
	default:
		return false
	}
}

// HashString hashes a validated domain name.
func HashString(name string) (primitives.Hash, error) {
	if !VerifyString(name) {
		return primitives.Hash{}, core.E("covenant.HashString", "invalid name", nil)
	}

	sum := sha3.Sum256([]byte(name))
	return primitives.Hash(sum), nil
}

// HashBinary hashes a validated domain name.
func HashBinary(name []byte) (primitives.Hash, error) {
	if !VerifyBinary(name) {
		return primitives.Hash{}, core.E("covenant.HashBinary", "invalid name", nil)
	}

	sum := sha3.Sum256(name)
	return primitives.Hash(sum), nil
}

// HashName hashes a validated domain name.
//
// It accepts either a string or a byte slice to mirror the JS reference API.
func HashName(name any) (primitives.Hash, error) {
	switch v := name.(type) {
	case string:
		return HashString(v)
	case []byte:
		return HashBinary(v)
	default:
		return primitives.Hash{}, core.E("covenant.HashName", "invalid name type", nil)
	}
}
