// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
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

// Blacklist mirrors the JS rules.blacklist set.
//
// The map is exported so callers can inspect the reserved labels used by the
// covenant verifier without duplicating the list.
var Blacklist = blacklist

// GetBlacklist returns the covenant verifier blacklist.
func GetBlacklist() map[string]struct{} {
	return Blacklist
}

// VerifyString reports whether a domain name meets the covenant rules.
//
//	ok := covenant.VerifyString("example")
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

	_, blocked := Blacklist[name]
	return !blocked
}

// GetVerifyString is an alias for VerifyString.
//
//	ok := covenant.GetVerifyString("example")
func GetVerifyString(name string) bool {
	return VerifyString(name)
}

// VerifyBinary reports whether a binary domain name meets the covenant rules.
//
//	ok := covenant.VerifyBinary([]byte("example"))
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

	_, blocked := Blacklist[string(name)]
	return !blocked
}

// GetVerifyBinary is an alias for VerifyBinary.
//
//	ok := covenant.GetVerifyBinary([]byte("example"))
func GetVerifyBinary(name []byte) bool {
	return VerifyBinary(name)
}

// VerifyName reports whether a name meets the covenant rules.
//
// It accepts either a string or a byte slice to mirror the JS reference API.
//
//	ok := covenant.VerifyName("example")
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

// GetVerifyName is an alias for VerifyName.
//
//	ok := covenant.GetVerifyName("example")
func GetVerifyName(name any) bool {
	return VerifyName(name)
}

// VerifyByName is an alias for VerifyName.
//
//	ok := covenant.VerifyByName("example")
func VerifyByName(name any) bool {
	return VerifyName(name)
}

// GetVerifyByName is an alias for VerifyByName.
//
//	ok := covenant.GetVerifyByName("example")
func GetVerifyByName(name any) bool {
	return VerifyByName(name)
}

// VerifyByString is an alias for VerifyString.
//
//	ok := covenant.VerifyByString("example")
func VerifyByString(name string) bool {
	return VerifyString(name)
}

// GetVerifyByString is an alias for VerifyByString.
//
//	ok := covenant.GetVerifyByString("example")
func GetVerifyByString(name string) bool {
	return VerifyByString(name)
}

// VerifyByBinary is an alias for VerifyBinary.
//
//	ok := covenant.VerifyByBinary([]byte("example"))
func VerifyByBinary(name []byte) bool {
	return VerifyBinary(name)
}

// GetVerifyByBinary is an alias for VerifyByBinary.
//
//	ok := covenant.GetVerifyByBinary([]byte("example"))
func GetVerifyByBinary(name []byte) bool {
	return VerifyByBinary(name)
}

// HashString hashes a validated domain name.
//
//	hash, err := covenant.HashString("example")
func HashString(name string) (primitives.Hash, error) {
	if !VerifyString(name) {
		return primitives.Hash{}, core.E("covenant.HashString", "invalid name", nil)
	}

	sum := core.SHA3Keccak256([]byte(name))
	return primitives.Hash(sum), nil
}

// GetHashString is an alias for HashString.
//
//	hash, err := covenant.GetHashString("example")
func GetHashString(name string) (primitives.Hash, error) {
	return HashString(name)
}

// HashByString is an alias for HashString.
//
//	hash, err := covenant.HashByString("example")
func HashByString(name string) (primitives.Hash, error) {
	return HashString(name)
}

// GetHashByString is an alias for HashByString.
//
//	hash, err := covenant.GetHashByString("example")
func GetHashByString(name string) (primitives.Hash, error) {
	return HashByString(name)
}

// HashBinary hashes a validated domain name.
//
//	hash, err := covenant.HashBinary([]byte("example"))
func HashBinary(name []byte) (primitives.Hash, error) {
	if !VerifyBinary(name) {
		return primitives.Hash{}, core.E("covenant.HashBinary", "invalid name", nil)
	}

	sum := core.SHA3Keccak256(name)
	return primitives.Hash(sum), nil
}

// GetHashBinary is an alias for HashBinary.
//
//	hash, err := covenant.GetHashBinary([]byte("example"))
func GetHashBinary(name []byte) (primitives.Hash, error) {
	return HashBinary(name)
}

// HashByBinary is an alias for HashBinary.
//
//	hash, err := covenant.HashByBinary([]byte("example"))
func HashByBinary(name []byte) (primitives.Hash, error) {
	return HashBinary(name)
}

// GetHashByBinary is an alias for HashByBinary.
//
//	hash, err := covenant.GetHashByBinary([]byte("example"))
func GetHashByBinary(name []byte) (primitives.Hash, error) {
	return HashByBinary(name)
}

// HashName hashes a validated domain name.
//
// It accepts either a string or a byte slice to mirror the JS reference API.
//
//	hash, err := covenant.HashName("example")
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

// GetHashName is an alias for HashName.
//
//	hash, err := covenant.GetHashName("example")
func GetHashName(name any) (primitives.Hash, error) {
	return HashName(name)
}

// HashByName is an alias for HashName.
//
//	hash, err := covenant.HashByName("example")
func HashByName(name any) (primitives.Hash, error) {
	return HashName(name)
}

// GetHashByName is an alias for HashByName.
//
//	hash, err := covenant.GetHashByName("example")
func GetHashByName(name any) (primitives.Hash, error) {
	return HashByName(name)
}
