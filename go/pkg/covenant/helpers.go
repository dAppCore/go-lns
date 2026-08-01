// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

type jsonraw []byte

func (m jsonraw) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return []byte(m), nil
}

func (m *jsonraw) UnmarshalJSON(data []byte) error {
	if m == nil {
		return core.NewError("covenant.jsonraw: nil receiver")
	}
	*m = append((*m)[:0], data...)
	return nil
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
