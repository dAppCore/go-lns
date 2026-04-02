// SPDX-License-Identifier: EUPL-1.2

package primitives

import (
	"bytes"
	"encoding/binary"

	core "dappco.re/go/core"
)

func normaliseCovenantIndex(index, length int) (int, bool) {
	if index < 0 {
		index += length
	}

	if index < 0 || index >= length {
		return 0, false
	}

	return index, true
}

func sizeVarint(n int) int {
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

func sizeVarBytes(n int) int {
	return sizeVarint(n) + n
}

// Len reports the number of covenant items.
func (c Covenant) Len() int {
	return len(c.Items)
}

// Clone returns a shallow copy of the covenant.
func (c Covenant) Clone() Covenant {
	return Covenant{
		Type:  c.Type,
		Items: c.ToArray(),
	}
}

// Inject copies another covenant into this one.
func (c *Covenant) Inject(other Covenant) *Covenant {
	c.Type = other.Type
	c.Items = other.ToArray()
	return c
}

// ToArray returns a shallow copy of the covenant item slice.
func (c Covenant) ToArray() [][]byte {
	if len(c.Items) == 0 {
		return nil
	}

	items := make([][]byte, len(c.Items))
	copy(items, c.Items)
	return items
}

// FromArray replaces the covenant items with a shallow copy of the provided slice.
func (c *Covenant) FromArray(items [][]byte) *Covenant {
	if len(items) == 0 {
		c.Items = nil
		return c
	}

	c.Items = make([][]byte, len(items))
	copy(c.Items, items)
	return c
}

// IndexOf reports the first matching item index or -1 if the item is absent.
func (c Covenant) IndexOf(data []byte) int {
	for i, item := range c.Items {
		if bytes.Equal(item, data) {
			return i
		}
	}

	return -1
}

// Get returns the item at the requested index.
func (c Covenant) Get(index int) ([]byte, error) {
	index, ok := normaliseCovenantIndex(index, len(c.Items))
	if !ok {
		return nil, core.E("primitives.Covenant.Get", "index out of range", nil)
	}

	return c.Items[index], nil
}

// Set replaces or appends an item at the requested index.
func (c *Covenant) Set(index int, item []byte) error {
	if index < 0 {
		index += len(c.Items)
	}

	if index < 0 || index > len(c.Items) {
		return core.E("primitives.Covenant.Set", "index out of range", nil)
	}

	if index == len(c.Items) {
		c.Items = append(c.Items, item)
		return nil
	}

	c.Items[index] = item
	return nil
}

// Push appends an item to the covenant.
func (c *Covenant) Push(item []byte) *Covenant {
	c.Items = append(c.Items, item)
	return c
}

// SetOpen configures the covenant as an OPEN operation.
func (c *Covenant) SetOpen(nameHash Hash, rawName []byte) *Covenant {
	c.Type = covenantTypeOpen
	c.Items = nil
	c.PushHash(nameHash)
	c.PushU32(0)
	c.Push(rawName)
	return c
}

// SetBid configures the covenant as a BID operation.
func (c *Covenant) SetBid(nameHash Hash, height uint32, rawName []byte, blind Hash) *Covenant {
	c.Type = covenantTypeBid
	c.Items = nil
	c.PushHash(nameHash)
	c.PushU32(height)
	c.Push(rawName)
	c.PushHash(blind)
	return c
}

// SetReveal configures the covenant as a REVEAL operation.
func (c *Covenant) SetReveal(nameHash Hash, height uint32, nonce Hash) *Covenant {
	c.Type = covenantTypeReveal
	c.Items = nil
	c.PushHash(nameHash)
	c.PushU32(height)
	c.PushHash(nonce)
	return c
}

// SetClaim configures the covenant as a CLAIM operation.
func (c *Covenant) SetClaim(nameHash Hash, height uint32, rawName []byte, flags uint8, commitHash Hash, commitHeight uint32) *Covenant {
	c.Type = covenantTypeClaim
	c.Items = nil
	c.PushHash(nameHash)
	c.PushU32(height)
	c.Push(rawName)
	c.PushU8(flags)
	c.PushHash(commitHash)
	c.PushU32(commitHeight)
	return c
}

// GetU8 returns the item at the requested index as a uint8.
func (c Covenant) GetU8(index int) (uint8, error) {
	item, err := c.Get(index)
	if err != nil {
		return 0, err
	}

	if len(item) != 1 {
		return 0, core.E("primitives.Covenant.GetU8", "item is not a uint8", nil)
	}

	return item[0], nil
}

// PushU8 appends a uint8 item.
func (c *Covenant) PushU8(num uint8) *Covenant {
	c.Items = append(c.Items, []byte{num})
	return c
}

// GetU32 returns the item at the requested index as a uint32.
func (c Covenant) GetU32(index int) (uint32, error) {
	item, err := c.Get(index)
	if err != nil {
		return 0, err
	}

	if len(item) != 4 {
		return 0, core.E("primitives.Covenant.GetU32", "item is not a uint32", nil)
	}

	return binary.LittleEndian.Uint32(item), nil
}

// PushU32 appends a uint32 item.
func (c *Covenant) PushU32(num uint32) *Covenant {
	item := make([]byte, 4)
	binary.LittleEndian.PutUint32(item, num)
	c.Items = append(c.Items, item)
	return c
}

// GetHash returns the item at the requested index as a 32-byte hash.
func (c Covenant) GetHash(index int) (Hash, error) {
	item, err := c.Get(index)
	if err != nil {
		return Hash{}, err
	}

	if len(item) != len(Hash{}) {
		return Hash{}, core.E("primitives.Covenant.GetHash", "item is not a hash", nil)
	}

	var hash Hash
	copy(hash[:], item)
	return hash, nil
}

// PushHash appends a 32-byte hash item.
func (c *Covenant) PushHash(hash Hash) *Covenant {
	item := make([]byte, len(hash))
	copy(item, hash[:])
	c.Items = append(c.Items, item)
	return c
}

// GetString returns the item at the requested index as a binary string.
func (c Covenant) GetString(index int) (string, error) {
	item, err := c.Get(index)
	if err != nil {
		return "", err
	}

	if len(item) == 0 || len(item) > 63 {
		return "", core.E("primitives.Covenant.GetString", "item is not a string", nil)
	}

	return string(item), nil
}

// PushString appends a binary string item.
func (c *Covenant) PushString(str string) *Covenant {
	c.Items = append(c.Items, []byte(str))
	return c
}

// GetSize returns the serialised size of all covenant items, excluding the
// type byte and item-count prefix.
func (c Covenant) GetSize() int {
	size := 0

	for _, item := range c.Items {
		size += sizeVarBytes(len(item))
	}

	return size
}

// GetVarSize returns the serialised covenant size including the type byte and
// item-count prefix.
func (c Covenant) GetVarSize() int {
	return 1 + sizeVarint(len(c.Items)) + c.GetSize()
}
