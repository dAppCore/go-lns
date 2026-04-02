// SPDX-License-Identifier: EUPL-1.2

package covenant

import "bytes"

// asciiLowerName lower-cases ASCII labels without performing Unicode
// case-folding. The helper also strips an optional canonical `.lthn` suffix so
// the catalog lookups accept the same normalized forms as the top-level LNS API.
func asciiLowerName(name string) (string, bool) {
	if len(name) == 0 {
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

	switch {
	case bytes.HasSuffix(buf, []byte(".lthn.")):
		buf = buf[:len(buf)-6]
	case bytes.HasSuffix(buf, []byte(".lthn")):
		buf = buf[:len(buf)-5]
	case bytes.HasSuffix(buf, []byte(".")):
		buf = buf[:len(buf)-1]
	}

	if len(buf) == 0 || len(buf) > maxNameSize {
		return "", false
	}

	return string(buf), true
}
