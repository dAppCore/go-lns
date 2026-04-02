// SPDX-License-Identifier: EUPL-1.2

package covenant

// asciiLowerName lower-cases ASCII labels without performing Unicode
// case-folding. The JS reference uses ASCII encoding for these catalog lookups,
// so non-ASCII input is rejected instead of being hashed differently.
func asciiLowerName(name string) (string, bool) {
	if len(name) == 0 || len(name) > maxNameSize {
		return "", false
	}

	buf := make([]byte, len(name))

	for i := 0; i < len(name); i++ {
		ch := name[i]

		if ch&0x80 != 0 {
			return "", false
		}

		if ch >= 'A' && ch <= 'Z' {
			ch += 'a' - 'A'
		}

		buf[i] = ch
	}

	return string(buf), true
}
