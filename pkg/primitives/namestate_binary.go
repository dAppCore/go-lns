// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"encoding/binary"

	core "dappco.re/go"
)

const (
	// MaxNameStateSize is the maximum serialized size of a name state,
	// mirroring the JS reference's `NameState.MAX_SIZE` constant.
	MaxNameStateSize = 668

	// NameStateMaxSize is an alias for MaxNameStateSize.
	NameStateMaxSize = MaxNameStateSize
)

// GetSize returns the serialized size of the name state.
func (ns NameState) GetSize() int {
	size := 1 + len(ns.Name) + 2 + len(ns.Data) + 4 + 4 + 2

	if !outpointIsNull(ns.Owner) {
		size += len(ns.Owner.TxHash)
		size += sizeVarint64(uint64(ns.Owner.Index))
	}

	if ns.Value != 0 {
		size += sizeVarint64(ns.Value)
	}

	if ns.Highest != 0 {
		size += sizeVarint64(ns.Highest)
	}

	if ns.Transfer != 0 {
		size += 4
	}

	if ns.Revoked != 0 {
		size += 4
	}

	if ns.Claimed != 0 {
		size += 4
	}

	if ns.Renewals != 0 {
		size += sizeVarint64(uint64(ns.Renewals))
	}

	return size
}

// MarshalBinary serializes the name state using the compact JS reference layout.
//
// The encoding does not include NameHash, which is derived from Name by the
// caller and is not part of the on-wire payload.
func (ns NameState) MarshalBinary() ([]byte, error) {
	if len(ns.Name) > 0xff {
		return nil, core.Errorf("primitives.NameState.MarshalBinary: name too long")
	}

	if len(ns.Data) > 0xffff {
		return nil, core.Errorf("primitives.NameState.MarshalBinary: data too long")
	}

	buf := make([]byte, 0, ns.GetSize())

	buf = append(buf, byte(len(ns.Name)))
	buf = append(buf, ns.Name...)

	var tmp [10]byte
	binary.LittleEndian.PutUint16(tmp[:2], uint16(len(ns.Data)))
	buf = append(buf, tmp[:2]...)
	buf = append(buf, ns.Data...)

	binary.LittleEndian.PutUint32(tmp[:4], ns.Height)
	buf = append(buf, tmp[:4]...)

	binary.LittleEndian.PutUint32(tmp[:4], ns.Renewal)
	buf = append(buf, tmp[:4]...)

	var field uint16
	if !outpointIsNull(ns.Owner) {
		field |= 1 << 0
	}
	if ns.Value != 0 {
		field |= 1 << 1
	}
	if ns.Highest != 0 {
		field |= 1 << 2
	}
	if ns.Transfer != 0 {
		field |= 1 << 3
	}
	if ns.Revoked != 0 {
		field |= 1 << 4
	}
	if ns.Claimed != 0 {
		field |= 1 << 5
	}
	if ns.Renewals != 0 {
		field |= 1 << 6
	}
	if ns.Registered {
		field |= 1 << 7
	}
	if ns.Expired {
		field |= 1 << 8
	}
	if ns.Weak {
		field |= 1 << 9
	}

	binary.LittleEndian.PutUint16(tmp[:2], field)
	buf = append(buf, tmp[:2]...)

	if !outpointIsNull(ns.Owner) {
		buf = append(buf, ns.Owner.TxHash[:]...)
		buf = appendVarint(buf, uint64(ns.Owner.Index))
	}

	if ns.Value != 0 {
		buf = appendVarint(buf, ns.Value)
	}

	if ns.Highest != 0 {
		buf = appendVarint(buf, ns.Highest)
	}

	if ns.Transfer != 0 {
		binary.LittleEndian.PutUint32(tmp[:4], ns.Transfer)
		buf = append(buf, tmp[:4]...)
	}

	if ns.Revoked != 0 {
		binary.LittleEndian.PutUint32(tmp[:4], ns.Revoked)
		buf = append(buf, tmp[:4]...)
	}

	if ns.Claimed != 0 {
		binary.LittleEndian.PutUint32(tmp[:4], ns.Claimed)
		buf = append(buf, tmp[:4]...)
	}

	if ns.Renewals != 0 {
		buf = appendVarint(buf, uint64(ns.Renewals))
	}

	return buf, nil
}

// UnmarshalBinary decodes the compact name-state layout used by the JS reference.
func (ns *NameState) UnmarshalBinary(data []byte) error {
	*ns = NameState{}

	if len(data) < 1 {
		return core.Errorf("primitives.NameState.UnmarshalBinary: short buffer")
	}

	nameLen := int(data[0])
	data = data[1:]

	if len(data) < nameLen+2+4+4+2 {
		return core.Errorf("primitives.NameState.UnmarshalBinary: short buffer")
	}

	ns.Name = append([]byte(nil), data[:nameLen]...)
	data = data[nameLen:]

	dataLen := int(binary.LittleEndian.Uint16(data[:2]))
	data = data[2:]

	if len(data) < dataLen+4+4+2 {
		return core.Errorf("primitives.NameState.UnmarshalBinary: short buffer")
	}

	ns.Data = append([]byte(nil), data[:dataLen]...)
	data = data[dataLen:]

	ns.Height = binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	ns.Renewal = binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	field := binary.LittleEndian.Uint16(data[:2])
	data = data[2:]

	if field&(1<<0) != 0 {
		if len(data) < len(ns.Owner.TxHash) {
			return core.Errorf("primitives.NameState.UnmarshalBinary: short buffer")
		}

		copy(ns.Owner.TxHash[:], data[:len(ns.Owner.TxHash)])
		data = data[len(ns.Owner.TxHash):]

		index, n, err := readVarint(data)
		if err != nil {
			return core.Errorf("primitives.NameState.UnmarshalBinary: %w", err)
		}

		ns.Owner.Index = uint32(index)
		data = data[n:]
	}

	if field&(1<<1) != 0 {
		value, n, err := readVarint(data)
		if err != nil {
			return core.Errorf("primitives.NameState.UnmarshalBinary: %w", err)
		}

		ns.Value = value
		data = data[n:]
	}

	if field&(1<<2) != 0 {
		highest, n, err := readVarint(data)
		if err != nil {
			return core.Errorf("primitives.NameState.UnmarshalBinary: %w", err)
		}

		ns.Highest = highest
		data = data[n:]
	}

	if field&(1<<3) != 0 {
		if len(data) < 4 {
			return core.Errorf("primitives.NameState.UnmarshalBinary: short buffer")
		}

		ns.Transfer = binary.LittleEndian.Uint32(data[:4])
		data = data[4:]
	}

	if field&(1<<4) != 0 {
		if len(data) < 4 {
			return core.Errorf("primitives.NameState.UnmarshalBinary: short buffer")
		}

		ns.Revoked = binary.LittleEndian.Uint32(data[:4])
		data = data[4:]
	}

	if field&(1<<5) != 0 {
		if len(data) < 4 {
			return core.Errorf("primitives.NameState.UnmarshalBinary: short buffer")
		}

		ns.Claimed = binary.LittleEndian.Uint32(data[:4])
		data = data[4:]
	}

	if field&(1<<6) != 0 {
		renewals, n, err := readVarint(data)
		if err != nil {
			return core.Errorf("primitives.NameState.UnmarshalBinary: %w", err)
		}

		ns.Renewals = uint32(renewals)
		data = data[n:]
	}

	if field&(1<<7) != 0 {
		ns.Registered = true
	}

	if field&(1<<8) != 0 {
		ns.Expired = true
	}

	if field&(1<<9) != 0 {
		ns.Weak = true
	}

	if len(data) != 0 {
		return core.Errorf("primitives.NameState.UnmarshalBinary: trailing data")
	}

	return nil
}
