// SPDX-License-Identifier: EUPL-1.2

package primitives

func compareBytes(left, right []byte) int {
	n := len(left)
	if len(right) < n {
		n = len(right)
	}
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
