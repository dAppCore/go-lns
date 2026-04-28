// SPDX-License-Identifier: EUPL-1.2

// Package covenant documents the reserved-name catalog path for the Lethean
// Name System. Path as documentation: pkg/covenant/reserved_lookup.go covers
// the reserved reference table, while locked_lookup.go covers the locked one.
//
// Look up a reserved label:
//
//	name, ok := covenant.GetReservedName("example")
package covenant

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"runtime"
	"sort"
	"sync"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

const reservedZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"
const reservedCatalogPrefixSize = 28

// ReservedName describes a name entry from the JS reserved-name reference data.
//
// The name is the leftmost label used for lookup, while Target preserves the
// full DNS name from the reference table.
type ReservedName struct {
	Name   string
	Hash   primitives.Hash
	Target string
	Value  uint64
	Root   bool
	Top100 bool
	Custom bool
	Zero   bool
}

// ReservedCatalog exposes the reserved-name reference table.
//
//	catalog := covenant.DefaultReservedCatalog()
type ReservedCatalog struct {
	byHash    map[primitives.Hash]ReservedName
	nameValue uint64
	rootValue uint64
	topValue  uint64
}

// ReservedEntry pairs a reserved-name hash with its catalog entry.
type ReservedEntry struct {
	Hash  primitives.Hash
	Value ReservedName
}

type reservedCatalogMeta struct {
	nameValue uint64
	rootValue uint64
	topValue  uint64
}

var (
	reservedCatalogOnce sync.Once
	reservedCatalog     *ReservedCatalog
)

// HasReservedName reports whether the JS reserved-name table contains the name.
//
// Lookup is case-insensitive and accepts a canonical `.lthn` suffix in
// addition to bare labels.
//
//	ok := covenant.HasReservedName("example")
func HasReservedName(name string) bool {
	return DefaultReservedCatalog().HasByName(name)
}

// GetHasReservedName is an alias for HasReservedName.
//
//	ok := covenant.GetHasReservedName("example")
func GetHasReservedName(name string) bool {
	return HasReservedName(name)
}

// HasReservedByName is an alias for HasReservedName.
//
//	ok := covenant.HasReservedByName("example")
func HasReservedByName(name string) bool {
	return HasReservedName(name)
}

// HasReservedString is an alias for HasReservedName.
//
//	ok := covenant.HasReservedString("example")
func HasReservedString(name string) bool {
	return HasReservedName(name)
}

// GetHasReservedString is an alias for HasReservedString.
//
//	ok := covenant.GetHasReservedString("example")
func GetHasReservedString(name string) bool {
	return HasReservedString(name)
}

// HasReservedBinary is an alias for HasReservedName.
//
//	ok := covenant.HasReservedBinary([]byte("example"))
func HasReservedBinary(name []byte) bool {
	return HasReservedName(string(name))
}

// GetHasReservedBinary is an alias for HasReservedBinary.
//
//	ok := covenant.GetHasReservedBinary([]byte("example"))
func GetHasReservedBinary(name []byte) bool {
	return HasReservedBinary(name)
}

// GetReservedName returns the reserved-name entry for the provided label.
//
// The input is lower-cased before hashing, and canonical `.lthn` suffixes are
// stripped before lookup.
//
//	item, ok := covenant.GetReservedName("example")
func GetReservedName(name string) (ReservedName, bool) {
	return DefaultReservedCatalog().GetByName(name)
}

// GetReservedByName is an alias for GetReservedName.
//
//	item, ok := covenant.GetReservedByName("example")
func GetReservedByName(name string) (ReservedName, bool) {
	return GetReservedName(name)
}

// GetReservedString is an alias for GetReservedName.
//
//	item, ok := covenant.GetReservedString("example")
func GetReservedString(name string) (ReservedName, bool) {
	return GetReservedName(name)
}

// GetReservedBinary is an alias for GetReservedName.
//
//	item, ok := covenant.GetReservedBinary([]byte("example"))
func GetReservedBinary(name []byte) (ReservedName, bool) {
	return GetReservedName(string(name))
}

// HasReservedByString is an alias for HasReservedString.
//
//	ok := covenant.HasReservedByString("example")
func HasReservedByString(name string) bool {
	return HasReservedString(name)
}

// GetHasReservedByString is an alias for HasReservedByString.
//
//	ok := covenant.GetHasReservedByString("example")
func GetHasReservedByString(name string) bool {
	return HasReservedByString(name)
}

// HasReservedByBinary is an alias for HasReservedBinary.
//
//	ok := covenant.HasReservedByBinary([]byte("example"))
func HasReservedByBinary(name []byte) bool {
	return HasReservedBinary(name)
}

// GetHasReservedByBinary is an alias for HasReservedByBinary.
//
//	ok := covenant.GetHasReservedByBinary([]byte("example"))
func GetHasReservedByBinary(name []byte) bool {
	return HasReservedByBinary(name)
}

// GetReservedByString is an alias for GetReservedString.
//
//	item, ok := covenant.GetReservedByString("example")
func GetReservedByString(name string) (ReservedName, bool) {
	return GetReservedString(name)
}

// GetReservedByBinary is an alias for GetReservedBinary.
//
//	item, ok := covenant.GetReservedByBinary([]byte("example"))
func GetReservedByBinary(name []byte) (ReservedName, bool) {
	return GetReservedBinary(name)
}

// GetReservedHash returns the reserved-name entry for a precomputed hash.
//
//	item, ok := covenant.GetReservedHash(hash)
func GetReservedHash(hash primitives.Hash) (ReservedName, bool) {
	return defaultReservedCatalog().Get(hash)
}

// HasReservedHash reports whether a precomputed hash is present in the catalog.
//
//	ok := covenant.HasReservedHash(hash)
func HasReservedHash(hash primitives.Hash) bool {
	return DefaultReservedCatalog().Has(hash)
}

// GetHasReservedHash is an alias for HasReservedHash.
//
//	ok := covenant.GetHasReservedHash(hash)
func GetHasReservedHash(hash primitives.Hash) bool {
	return HasReservedHash(hash)
}

// GetReservedByHash is an alias for GetReservedHash.
//
//	item, ok := covenant.GetReservedByHash(hash)
func GetReservedByHash(hash primitives.Hash) (ReservedName, bool) {
	return GetReservedHash(hash)
}

// HasReservedByHash is an alias for HasReservedHash.
//
//	ok := covenant.HasReservedByHash(hash)
func HasReservedByHash(hash primitives.Hash) bool {
	return HasReservedHash(hash)
}

// GetHasReservedByHash is an alias for HasReservedByHash.
//
//	ok := covenant.GetHasReservedByHash(hash)
func GetHasReservedByHash(hash primitives.Hash) bool {
	return HasReservedByHash(hash)
}

// GetHasReservedByName is an alias for HasReservedByName.
//
//	ok := covenant.GetHasReservedByName("example")
func GetHasReservedByName(name string) bool {
	return HasReservedByName(name)
}

// DefaultReservedCatalog returns the lazily loaded reserved-name catalog.
//
//	catalog := covenant.DefaultReservedCatalog()
func DefaultReservedCatalog() *ReservedCatalog {
	return defaultReservedCatalog()
}

// GetReservedCatalog is an alias for DefaultReservedCatalog.
//
//	catalog := covenant.GetReservedCatalog()
func GetReservedCatalog() *ReservedCatalog {
	return DefaultReservedCatalog()
}

// GetDefaultReservedCatalog is an alias for DefaultReservedCatalog.
//
//	catalog := covenant.GetDefaultReservedCatalog()
func GetDefaultReservedCatalog() *ReservedCatalog {
	return DefaultReservedCatalog()
}

// Size reports the number of reserved-name entries in the catalog.
//
//	size := catalog.Size()
func (c *ReservedCatalog) Size() int {
	if c == nil {
		return 0
	}

	return len(c.byHash)
}

// GetSize is an alias for Size.
//
//	size := catalog.GetSize()
func (c *ReservedCatalog) GetSize() int {
	return c.Size()
}

// NameValue reports the base value used for reserved-name lookups.
//
//	base := catalog.NameValue()
func (c *ReservedCatalog) NameValue() uint64 {
	if c == nil {
		return 0
	}

	return c.nameValue
}

// GetNameValue is an alias for NameValue.
//
//	base := catalog.GetNameValue()
func (c *ReservedCatalog) GetNameValue() uint64 {
	return c.NameValue()
}

// RootValue reports the root-name increment used for reserved-name lookups.
//
//	root := catalog.RootValue()
func (c *ReservedCatalog) RootValue() uint64 {
	if c == nil {
		return 0
	}

	return c.rootValue
}

// GetRootValue is an alias for RootValue.
//
//	root := catalog.GetRootValue()
func (c *ReservedCatalog) GetRootValue() uint64 {
	return c.RootValue()
}

// TopValue reports the top-100 increment used for reserved-name lookups.
//
//	top := catalog.TopValue()
func (c *ReservedCatalog) TopValue() uint64 {
	if c == nil {
		return 0
	}

	return c.topValue
}

// GetTopValue is an alias for TopValue.
//
//	top := catalog.GetTopValue()
func (c *ReservedCatalog) GetTopValue() uint64 {
	return c.TopValue()
}

// PrefixSize reports the byte offset where the reserved-name table entries begin.
//
//	offset := catalog.PrefixSize()
func (c *ReservedCatalog) PrefixSize() int {
	return reservedCatalogPrefixSize
}

// GetPrefixSize is an alias for PrefixSize.
//
//	offset := catalog.GetPrefixSize()
func (c *ReservedCatalog) GetPrefixSize() int {
	return c.PrefixSize()
}

// Has reports whether the catalog contains the hash.
//
//	ok := catalog.Has(hash)
func (c *ReservedCatalog) Has(hash primitives.Hash) bool {
	if c == nil {
		return false
	}

	_, ok := c.byHash[hash]
	return ok
}

// HasByHash is an alias for Has.
//
//	ok := catalog.HasByHash(hash)
func (c *ReservedCatalog) HasByHash(hash primitives.Hash) bool {
	return c.Has(hash)
}

// GetHasByHash is an alias for HasByHash.
//
//	ok := catalog.GetHasByHash(hash)
func (c *ReservedCatalog) GetHasByHash(hash primitives.Hash) bool {
	return c.HasByHash(hash)
}

// Get returns the reserved-name entry for the hash.
//
//	item, ok := catalog.Get(hash)
func (c *ReservedCatalog) Get(hash primitives.Hash) (ReservedName, bool) {
	if c == nil {
		return ReservedName{}, false
	}

	item, ok := c.byHash[hash]
	return item, ok
}

// GetByHash is an alias for Get.
//
//	item, ok := catalog.GetByHash(hash)
func (c *ReservedCatalog) GetByHash(hash primitives.Hash) (ReservedName, bool) {
	return c.Get(hash)
}

// HasByName reports whether the catalog contains the provided label.
//
//	ok := catalog.HasByName("example")
func (c *ReservedCatalog) HasByName(name string) bool {
	_, ok := c.GetByName(name)
	return ok
}

// GetHasByName is an alias for HasByName.
//
//	ok := catalog.GetHasByName("example")
func (c *ReservedCatalog) GetHasByName(name string) bool {
	return c.HasByName(name)
}

// GetByName returns the reserved-name entry for the provided label.
//
// The lookup lower-cases ASCII input, strips canonical `.lthn` suffixes, and
// hashes the normalized label directly.
//
//	item, ok := catalog.GetByName("example")
func (c *ReservedCatalog) GetByName(name string) (ReservedName, bool) {
	name, ok := asciiLowerName(name)
	if !ok {
		return ReservedName{}, false
	}

	hash := primitives.Hash(core.SHA3_256([]byte(name)))
	return c.Get(hash)
}

// Entries returns the catalog contents as hash/item pairs.
//
//	entries := catalog.Entries()
func (c *ReservedCatalog) Entries() []ReservedEntry {
	if c == nil || len(c.byHash) == 0 {
		return nil
	}

	hashes := c.sortedHashes()
	entries := make([]ReservedEntry, 0, len(hashes))

	for _, hash := range hashes {
		item, ok := c.byHash[hash]
		if !ok {
			continue
		}

		entries = append(entries, ReservedEntry{Hash: hash, Value: item})
	}

	return entries
}

// GetEntries is an alias for Entries.
//
//	entries := catalog.GetEntries()
func (c *ReservedCatalog) GetEntries() []ReservedEntry {
	return c.Entries()
}

// Keys returns the reserved-name hashes in deterministic order.
//
//	keys := catalog.Keys()
func (c *ReservedCatalog) Keys() []primitives.Hash {
	if c == nil || len(c.byHash) == 0 {
		return nil
	}

	hashes := c.sortedHashes()
	keys := make([]primitives.Hash, len(hashes))
	copy(keys, hashes)
	return keys
}

// GetKeys is an alias for Keys.
//
//	keys := catalog.GetKeys()
func (c *ReservedCatalog) GetKeys() []primitives.Hash {
	return c.Keys()
}

// Values returns the reserved-name entries in deterministic order.
//
//	values := catalog.Values()
func (c *ReservedCatalog) Values() []ReservedName {
	if c == nil || len(c.byHash) == 0 {
		return nil
	}

	hashes := c.sortedHashes()
	values := make([]ReservedName, 0, len(hashes))

	for _, hash := range hashes {
		item, ok := c.byHash[hash]
		if !ok {
			continue
		}

		values = append(values, item)
	}

	return values
}

// GetValues is an alias for Values.
//
//	values := catalog.GetValues()
func (c *ReservedCatalog) GetValues() []ReservedName {
	return c.Values()
}

func (c *ReservedCatalog) sortedHashes() []primitives.Hash {
	hashes := make([]primitives.Hash, 0, len(c.byHash))

	for hash := range c.byHash {
		hashes = append(hashes, hash)
	}

	sort.Slice(hashes, func(i, j int) bool {
		return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
	})

	return hashes
}

func defaultReservedCatalog() *ReservedCatalog {
	reservedCatalogOnce.Do(func() {
		catalog, err := loadReservedCatalog()
		if err != nil {
			reservedCatalog = &ReservedCatalog{byHash: make(map[primitives.Hash]ReservedName)}
			return
		}

		reservedCatalog = catalog
	})

	if reservedCatalog == nil {
		reservedCatalog = &ReservedCatalog{byHash: make(map[primitives.Hash]ReservedName)}
	}

	return reservedCatalog
}

func loadReservedCatalog() (*ReservedCatalog, error) {
	data, err := readReservedData()
	if err != nil {
		return nil, err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, core.E("covenant.loadReservedCatalog", "failed to decode reserved catalog", err)
	}

	catalog := &ReservedCatalog{
		byHash: make(map[primitives.Hash]ReservedName, len(raw)),
	}

	meta, err := parseReservedMeta(raw[reservedZeroHash])
	if err != nil {
		return nil, err
	}

	catalog.nameValue = meta.nameValue
	catalog.rootValue = meta.rootValue
	catalog.topValue = meta.topValue

	for key, value := range raw {
		if key == reservedZeroHash {
			continue
		}

		hash, err := decodeReservedHash(key)
		if err != nil {
			return nil, err
		}

		item, err := parseReservedEntry(hash, value, meta)
		if err != nil {
			return nil, err
		}

		catalog.byHash[hash] = item
	}

	return catalog, nil
}

func parseReservedMeta(raw json.RawMessage) (reservedCatalogMeta, error) {
	var values []uint64
	if err := json.Unmarshal(raw, &values); err != nil {
		return reservedCatalogMeta{}, core.E("covenant.parseReservedMeta", "failed to decode reserved metadata", err)
	}

	if len(values) != 4 {
		return reservedCatalogMeta{}, core.E("covenant.parseReservedMeta", "invalid metadata size", nil)
	}

	return reservedCatalogMeta{
		nameValue: values[1],
		rootValue: values[2],
		topValue:  values[3],
	}, nil
}

func parseReservedEntry(hash primitives.Hash, raw json.RawMessage, meta reservedCatalogMeta) (ReservedName, error) {
	var fields []json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return ReservedName{}, core.E("covenant.parseReservedEntry", "failed to decode reserved entry", err)
	}

	if len(fields) < 2 || len(fields) > 3 {
		return ReservedName{}, core.E("covenant.parseReservedEntry", "invalid entry size", nil)
	}

	var target string
	if err := json.Unmarshal(fields[0], &target); err != nil {
		return ReservedName{}, core.E("covenant.parseReservedEntry", "failed to decode target", err)
	}

	parts := core.SplitN(target, ".", 2)
	if len(parts) != 2 || len(parts[0]) == 0 {
		return ReservedName{}, core.E("covenant.parseReservedEntry", "invalid target", nil)
	}

	flags, err := parseReservedUint(fields[1])
	if err != nil {
		return ReservedName{}, core.E("covenant.parseReservedEntry", "failed to decode flags", err)
	}

	var custom uint64
	if len(fields) == 3 {
		custom, err = parseReservedUint(fields[2])
		if err != nil {
			return ReservedName{}, core.E("covenant.parseReservedEntry", "failed to decode custom value", err)
		}
	}

	item := ReservedName{
		Name:   parts[0],
		Hash:   hash,
		Target: target,
		Root:   flags&1 != 0,
		Top100: flags&2 != 0,
		Custom: flags&4 != 0,
		Zero:   flags&8 != 0,
	}

	item.Value = meta.nameValue

	if item.Root {
		item.Value += meta.rootValue
	}

	if item.Top100 {
		item.Value += meta.topValue
	}

	if item.Custom {
		item.Value += custom
	}

	if item.Zero {
		item.Value = 0
	}

	return item, nil
}

func parseReservedUint(raw json.RawMessage) (uint64, error) {
	var num uint64
	if err := json.Unmarshal(raw, &num); err != nil {
		return 0, err
	}

	return num, nil
}

func decodeReservedHash(key string) (primitives.Hash, error) {
	raw, err := hex.DecodeString(key)
	if err != nil {
		return primitives.Hash{}, core.E("covenant.decodeReservedHash", "failed to decode reserved hash", err)
	}

	if len(raw) != len(primitives.Hash{}) {
		return primitives.Hash{}, core.E("covenant.decodeReservedHash", "invalid hash length", nil)
	}

	var hash primitives.Hash
	copy(hash[:], raw)
	return hash, nil
}

func readReservedData() ([]byte, error) {
	candidates := reservedDataPaths()
	fs := core.Fs{}

	for _, path := range candidates {
		if result := fs.Read(path); result.OK {
			return []byte(result.Value.(string)), nil
		}
	}

	return nil, core.E("covenant.readReservedData", "could not locate docs/js-covenants/names.json", nil)
}

func reservedDataPaths() []string {
	paths := []string{
		filepath.Join("docs", "js-covenants", "names.json"),
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		base := filepath.Dir(file)
		paths = append(paths,
			filepath.Join(base, "..", "..", "docs", "js-covenants", "names.json"),
			filepath.Join(base, "..", "..", "..", "docs", "js-covenants", "names.json"),
		)
	}

	return paths
}
