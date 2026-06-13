// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// MarshalBinary serializes the covenant using the compact JS reference layout.
//
// The encoding is:
//   - 1 byte covenant type
//   - varint item count
//   - repeated varbytes items
func (c Covenant) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 0, c.GetVarSize())

	buf = append(buf, byte(c.Type))
	buf = appendVarint(buf, uint64(len(c.Items)))

	for _, item := range c.Items {
		buf = appendVarBytes(buf, item)
	}

	return buf, nil
}

// UnmarshalBinary decodes the compact covenant layout used by the JS reference.
func (c *Covenant) UnmarshalBinary(data []byte) error {
	*c = Covenant{}

	if len(data) < 1 {
		return core.Errorf("primitives.Covenant.UnmarshalBinary: short buffer")
	}

	c.Type = uint8(data[0])
	data = data[1:]

	count, consumed, err := readVarint(data)
	if err != nil {
		return core.Errorf("primitives.Covenant.UnmarshalBinary: %w", err)
	}
	data = data[consumed:]

	// Each item is a varbytes entry that consumes at least one byte (its varint
	// length prefix). Reject a count the remaining buffer cannot hold before
	// allocating, so an untrusted length prefix cannot drive an unbounded make().
	if count > uint64(len(data)) {
		return core.Errorf("primitives.Covenant.UnmarshalBinary: item count %d exceeds buffer capacity", count)
	}

	c.Items = make([][]byte, 0, int(count))
	for i := uint64(0); i < count; i++ {
		item, n, err := readVarBytes(data)
		if err != nil {
			return core.Errorf("primitives.Covenant.UnmarshalBinary: %w", err)
		}
		c.Items = append(c.Items, item)
		data = data[n:]
	}

	if len(data) != 0 {
		return core.Errorf("primitives.Covenant.UnmarshalBinary: trailing data")
	}

	return nil
}
