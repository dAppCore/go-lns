// SPDX-License-Identifier: EUPL-1.2

package nameutil

import core "github.com/dappcore/go"

// CatalogLabel converts string and byte inputs into a raw catalog label.
//
// It accepts the same input types as the public service helpers and preserves
// the bytes exactly, without applying canonical name suffix trimming.
func CatalogLabel(name any) (string, bool) {
	switch v := name.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	default:
		return "", false
	}
}

// Canonicalize normalizes string and byte inputs into a lowercased LNS label.
//
// The helper rejects non-ASCII input, lowercases ASCII letters, and strips the
// optional trailing `.lthn` or trailing dot suffix used by service callers.
func Canonicalize(name any) (string, bool) {
	var value []byte

	switch v := name.(type) {
	case string:
		value = []byte(v)
	case []byte:
		value = append([]byte(nil), v...)
	default:
		return "", false
	}

	if len(value) == 0 {
		return "", false
	}

	for i := range value {
		if value[i]&0x80 != 0 {
			return "", false
		}

		if value[i] >= 'A' && value[i] <= 'Z' {
			value[i] += 'a' - 'A'
		}
	}

	str := string(value)

	switch {
	case core.HasSuffix(str, ".lthn."):
		str = core.TrimSuffix(str, ".lthn.")
	case core.HasSuffix(str, ".lthn"):
		str = core.TrimSuffix(str, ".lthn")
	case core.HasSuffix(str, "."):
		str = core.TrimSuffix(str, ".")
	}

	if len(str) == 0 {
		return "", false
	}

	return str, true
}
