// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/binary"

	core "dappco.re/go"
)

// NameDelta tracks a sparse update to a name state.
//
// The zero value means "no changes recorded". Pointer fields let the delta
// distinguish between "not set" and "set to the zero value".
type NameDelta struct {
	Height     *uint32
	Renewal    *uint32
	Owner      *Outpoint
	Value      *uint64
	Highest    *uint64
	Data       *[]byte
	Transfer   *uint32
	Revoked    *uint32
	Claimed    *uint32
	Renewals   *uint32
	Registered *bool
	Expired    *bool
	Weak       *bool
}

// IsNull reports whether the delta contains any recorded changes.
func (d NameDelta) IsNull() bool {
	return d.Height == nil &&
		d.Renewal == nil &&
		d.Owner == nil &&
		d.Value == nil &&
		d.Highest == nil &&
		d.Data == nil &&
		d.Transfer == nil &&
		d.Revoked == nil &&
		d.Claimed == nil &&
		d.Renewals == nil &&
		d.Registered == nil &&
		d.Expired == nil &&
		d.Weak == nil
}

// Clone returns a shallow copy of the delta.
func (d NameDelta) Clone() NameDelta {
	return d
}

// Inject copies another delta into this one.
func (d *NameDelta) Inject(other NameDelta) *NameDelta {
	*d = other
	return d
}

// Clear resets the delta to its zero value.
func (d *NameDelta) Clear() *NameDelta {
	*d = NameDelta{}
	return d
}

// GetField returns the sparse field mask used by the on-wire encoding.
func (d NameDelta) GetField() uint32 {
	var field uint32

	if d.Height != nil {
		field |= 1 << 0
	}

	if d.Renewal != nil {
		field |= 1 << 1
	}

	if d.Owner != nil {
		field |= 1 << 2

		if !d.Owner.IsNull() {
			field |= 1 << 3
		}
	}

	if d.Value != nil {
		field |= 1 << 4

		if *d.Value != 0 {
			field |= 1 << 5
		}
	}

	if d.Highest != nil {
		field |= 1 << 6

		if *d.Highest != 0 {
			field |= 1 << 7
		}
	}

	if d.Data != nil {
		field |= 1 << 8
		field |= 1 << 9
	}

	if d.Transfer != nil {
		field |= 1 << 10

		if *d.Transfer != 0 {
			field |= 1 << 11
		}
	}

	if d.Revoked != nil {
		field |= 1 << 12

		if *d.Revoked != 0 {
			field |= 1 << 13
		}
	}

	if d.Claimed != nil {
		field |= 1 << 14

		if *d.Claimed != 0 {
			field |= 1 << 15
		}
	}

	if d.Renewals != nil {
		field |= 1 << 16

		if *d.Renewals != 0 {
			field |= 1 << 17
		}
	}

	if d.Registered != nil {
		field |= 1 << 18

		if *d.Registered {
			field |= 1 << 19
		}
	}

	if d.Expired != nil {
		field |= 1 << 20

		if *d.Expired {
			field |= 1 << 21
		}
	}

	if d.Weak != nil {
		field |= 1 << 22

		if *d.Weak {
			field |= 1 << 23
		}
	}

	return field
}

// GetSize returns the serialized size of the delta payload.
func (d NameDelta) GetSize() int {
	size := 4

	if d.Height != nil {
		size += 4
	}

	if d.Renewal != nil {
		size += 4
	}

	if d.Owner != nil && !d.Owner.IsNull() {
		size += len(d.Owner.TxHash)
		size += sizeVarint64(uint64(d.Owner.Index))
	}

	if d.Value != nil && *d.Value != 0 {
		size += sizeVarint64(*d.Value)
	}

	if d.Highest != nil && *d.Highest != 0 {
		size += sizeVarint64(*d.Highest)
	}

	if d.Data != nil {
		size += sizeVarBytes(len(*d.Data))
	}

	if d.Transfer != nil && *d.Transfer != 0 {
		size += 4
	}

	if d.Revoked != nil && *d.Revoked != 0 {
		size += 4
	}

	if d.Claimed != nil && *d.Claimed != 0 {
		size += 4
	}

	if d.Renewals != nil && *d.Renewals != 0 {
		size += sizeVarint64(uint64(*d.Renewals))
	}

	return size
}

// MarshalBinary serialises the sparse delta to the same compact format used by
// the JS reference implementation.
func (d NameDelta) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 0, d.GetSize())

	var tmp [10]byte
	binary.LittleEndian.PutUint32(tmp[:4], d.GetField())
	buf = append(buf, tmp[:4]...)

	if d.Height != nil {
		binary.LittleEndian.PutUint32(tmp[:4], *d.Height)
		buf = append(buf, tmp[:4]...)
	}

	if d.Renewal != nil {
		binary.LittleEndian.PutUint32(tmp[:4], *d.Renewal)
		buf = append(buf, tmp[:4]...)
	}

	if d.Owner != nil && !d.Owner.IsNull() {
		buf = append(buf, d.Owner.TxHash[:]...)
		buf = appendVarint(buf, uint64(d.Owner.Index))
	}

	if d.Value != nil && *d.Value != 0 {
		buf = appendVarint(buf, *d.Value)
	}

	if d.Highest != nil && *d.Highest != 0 {
		buf = appendVarint(buf, *d.Highest)
	}

	if d.Data != nil {
		buf = appendVarBytes(buf, *d.Data)
	}

	if d.Transfer != nil && *d.Transfer != 0 {
		binary.LittleEndian.PutUint32(tmp[:4], *d.Transfer)
		buf = append(buf, tmp[:4]...)
	}

	if d.Revoked != nil && *d.Revoked != 0 {
		binary.LittleEndian.PutUint32(tmp[:4], *d.Revoked)
		buf = append(buf, tmp[:4]...)
	}

	if d.Claimed != nil && *d.Claimed != 0 {
		binary.LittleEndian.PutUint32(tmp[:4], *d.Claimed)
		buf = append(buf, tmp[:4]...)
	}

	if d.Renewals != nil && *d.Renewals != 0 {
		buf = appendVarint(buf, uint64(*d.Renewals))
	}

	return buf, nil
}

// UnmarshalBinary decodes a sparse delta from the compact wire format.
func (d *NameDelta) UnmarshalBinary(data []byte) error {
	*d = NameDelta{}

	if len(data) < 4 {
		return core.Errorf("primitives.NameDelta.UnmarshalBinary: short buffer")
	}

	field := binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	if field&(1<<0) != 0 {
		if len(data) < 4 {
			return core.Errorf("primitives.NameDelta.UnmarshalBinary: short buffer")
		}
		height := binary.LittleEndian.Uint32(data[:4])
		d.Height = &height
		data = data[4:]
	}

	if field&(1<<1) != 0 {
		if len(data) < 4 {
			return core.Errorf("primitives.NameDelta.UnmarshalBinary: short buffer")
		}
		renewal := binary.LittleEndian.Uint32(data[:4])
		d.Renewal = &renewal
		data = data[4:]
	}

	if field&(1<<2) != 0 {
		owner := Outpoint{}

		if field&(1<<3) != 0 {
			if len(data) < len(owner.TxHash) {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: short buffer")
			}
			copy(owner.TxHash[:], data[:len(owner.TxHash)])
			data = data[len(owner.TxHash):]

			index, n, err := readVarint(data)
			if err != nil {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: %w", err)
			}
			owner.Index = uint32(index)
			data = data[n:]
		}

		d.Owner = &owner
	}

	if field&(1<<4) != 0 {
		value := uint64(0)
		if field&(1<<5) != 0 {
			var n int
			var err error
			value, n, err = readVarint(data)
			if err != nil {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: %w", err)
			}
			data = data[n:]
		}
		d.Value = &value
	}

	if field&(1<<6) != 0 {
		highest := uint64(0)
		if field&(1<<7) != 0 {
			var n int
			var err error
			highest, n, err = readVarint(data)
			if err != nil {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: %w", err)
			}
			data = data[n:]
		}
		d.Highest = &highest
	}

	if field&(1<<8) != 0 {
		payload := []byte{}
		if field&(1<<9) != 0 {
			var n int
			var err error
			payload, n, err = readVarBytes(data)
			if err != nil {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: %w", err)
			}
			data = data[n:]
		}
		d.Data = &payload
	}

	if field&(1<<10) != 0 {
		transfer := uint32(0)
		if field&(1<<11) != 0 {
			if len(data) < 4 {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: short buffer")
			}
			transfer = binary.LittleEndian.Uint32(data[:4])
			data = data[4:]
		}
		d.Transfer = &transfer
	}

	if field&(1<<12) != 0 {
		revoked := uint32(0)
		if field&(1<<13) != 0 {
			if len(data) < 4 {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: short buffer")
			}
			revoked = binary.LittleEndian.Uint32(data[:4])
			data = data[4:]
		}
		d.Revoked = &revoked
	}

	if field&(1<<14) != 0 {
		claimed := uint32(0)
		if field&(1<<15) != 0 {
			if len(data) < 4 {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: short buffer")
			}
			claimed = binary.LittleEndian.Uint32(data[:4])
			data = data[4:]
		}
		d.Claimed = &claimed
	}

	if field&(1<<16) != 0 {
		renewals := uint64(0)
		if field&(1<<17) != 0 {
			var n int
			var err error
			renewals, n, err = readVarint(data)
			if err != nil {
				return core.Errorf("primitives.NameDelta.UnmarshalBinary: %w", err)
			}
			data = data[n:]
		}
		value := uint32(renewals)
		d.Renewals = &value
	}

	if field&(1<<18) != 0 {
		registered := false
		if field&(1<<19) != 0 {
			registered = true
		}
		d.Registered = &registered
	}

	if field&(1<<20) != 0 {
		expired := false
		if field&(1<<21) != 0 {
			expired = true
		}
		d.Expired = &expired
	}

	if field&(1<<22) != 0 {
		weak := false
		if field&(1<<23) != 0 {
			weak = true
		}
		d.Weak = &weak
	}

	if len(data) != 0 {
		return core.Errorf("primitives.NameDelta.UnmarshalBinary: trailing data")
	}

	return nil
}

func appendVarint(dst []byte, n uint64) []byte {
	switch {
	case n < 0xfd:
		return append(dst, byte(n))
	case n <= 0xffff:
		return append(dst, 0xfd, byte(n), byte(n>>8))
	case n <= 0xffffffff:
		return append(dst, 0xfe, byte(n), byte(n>>8), byte(n>>16), byte(n>>24))
	default:
		return append(dst,
			0xff,
			byte(n), byte(n>>8), byte(n>>16), byte(n>>24),
			byte(n>>32), byte(n>>40), byte(n>>48), byte(n>>56))
	}
}

func appendVarBytes(dst, src []byte) []byte {
	dst = appendVarint(dst, uint64(len(src)))
	return append(dst, src...)
}

func sizeVarint64(n uint64) int {
	switch {
	case n < 0xfd:
		return 1
	case n <= 0xffff:
		return 3
	case n <= 0xffffffff:
		return 5
	default:
		return 9
	}
}

func readVarint(src []byte) (uint64, int, error) {
	if len(src) == 0 {
		return 0, 0, core.Errorf("short buffer")
	}

	prefix := src[0]

	switch {
	case prefix < 0xfd:
		return uint64(prefix), 1, nil
	case prefix == 0xfd:
		if len(src) < 3 {
			return 0, 0, core.Errorf("short buffer")
		}
		return uint64(binary.LittleEndian.Uint16(src[1:3])), 3, nil
	case prefix == 0xfe:
		if len(src) < 5 {
			return 0, 0, core.Errorf("short buffer")
		}
		return uint64(binary.LittleEndian.Uint32(src[1:5])), 5, nil
	default:
		if len(src) < 9 {
			return 0, 0, core.Errorf("short buffer")
		}
		return binary.LittleEndian.Uint64(src[1:9]), 9, nil
	}
}

func readVarBytes(src []byte) ([]byte, int, error) {
	n, consumed, err := readVarint(src)
	if err != nil {
		return nil, 0, err
	}

	if uint64(len(src[consumed:])) < n {
		return nil, 0, core.Errorf("short buffer")
	}

	end := consumed + int(n)
	return append([]byte(nil), src[consumed:end]...), end, nil
}
