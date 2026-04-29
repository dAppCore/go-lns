// SPDX-License-Identifier: EUPL-1.2

// Package covenant documents the locked-name catalog path for the Lethean Name
// System. Path as documentation: pkg/covenant/locked_lookup.go covers the
// locked reference table, while reserved_lookup.go covers the reserved table.
//
// Look up a locked label:
//
//	name, ok := covenant.GetLockedName("example")
package covenant

import (
	"encoding/hex"
	"runtime"
	"sort"
	"sync"

	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

// LockedName describes an entry from the JS locked-name reference data.
type LockedName struct {
	Name   string
	Hash   primitives.Hash
	Target string
	Root   bool
	Custom bool
}

// LockedCatalog exposes the locked-name reference table.
//
//	catalog := covenant.DefaultLockedCatalog()
type LockedCatalog struct {
	byHash map[primitives.Hash]LockedName
}

// LockedEntry pairs a locked-name hash with its catalog entry.
type LockedEntry struct {
	Hash  primitives.Hash
	Value LockedName
}

const lockedCatalogPrefixSize = 4

var (
	lockedCatalogOnce sync.Once
	lockedCatalog     *LockedCatalog
)

// HasLockedName reports whether the JS lockup table contains the name.
//
// Lookup is case-insensitive and accepts a canonical `.lthn` suffix in
// addition to bare labels.
//
//	ok := covenant.HasLockedName("example")
func HasLockedName(name string) bool {
	return DefaultLockedCatalog().HasByName(name)
}

// GetHasLockedName is an alias for HasLockedName.
//
//	ok := covenant.GetHasLockedName("example")
func GetHasLockedName(name string) bool {
	return HasLockedName(name)
}

// HasLockedByName is an alias for HasLockedName.
//
//	ok := covenant.HasLockedByName("example")
func HasLockedByName(name string) bool {
	return HasLockedName(name)
}

// GetHasLockedByName is an alias for HasLockedByName.
//
//	ok := covenant.GetHasLockedByName("example")
func GetHasLockedByName(name string) bool {
	return HasLockedByName(name)
}

// HasLockedString is an alias for HasLockedName.
//
//	ok := covenant.HasLockedString("example")
func HasLockedString(name string) bool {
	return HasLockedName(name)
}

// GetHasLockedString is an alias for HasLockedString.
//
//	ok := covenant.GetHasLockedString("example")
func GetHasLockedString(name string) bool {
	return HasLockedString(name)
}

// HasLockedBinary is an alias for HasLockedName.
//
//	ok := covenant.HasLockedBinary([]byte("example"))
func HasLockedBinary(name []byte) bool {
	return HasLockedName(string(name))
}

// GetHasLockedBinary is an alias for HasLockedBinary.
//
//	ok := covenant.GetHasLockedBinary([]byte("example"))
func GetHasLockedBinary(name []byte) bool {
	return HasLockedBinary(name)
}

// GetLockedName returns the locked-name entry for the provided label.
//
// The input is lower-cased before hashing, and canonical `.lthn` suffixes are
// stripped before lookup.
//
//	item, ok := covenant.GetLockedName("example")
func GetLockedName(name string) (LockedName, bool) {
	return DefaultLockedCatalog().GetByName(name)
}

// GetLockedByName is an alias for GetLockedName.
//
//	item, ok := covenant.GetLockedByName("example")
func GetLockedByName(name string) (LockedName, bool) {
	return GetLockedName(name)
}

// GetLockedString is an alias for GetLockedName.
//
//	item, ok := covenant.GetLockedString("example")
func GetLockedString(name string) (LockedName, bool) {
	return GetLockedName(name)
}

// GetLockedBinary is an alias for GetLockedName.
//
//	item, ok := covenant.GetLockedBinary([]byte("example"))
func GetLockedBinary(name []byte) (LockedName, bool) {
	return GetLockedName(string(name))
}

// HasLockedByString is an alias for HasLockedString.
//
//	ok := covenant.HasLockedByString("example")
func HasLockedByString(name string) bool {
	return HasLockedString(name)
}

// GetHasLockedByString is an alias for HasLockedByString.
//
//	ok := covenant.GetHasLockedByString("example")
func GetHasLockedByString(name string) bool {
	return HasLockedByString(name)
}

// HasLockedByBinary is an alias for HasLockedBinary.
//
//	ok := covenant.HasLockedByBinary([]byte("example"))
func HasLockedByBinary(name []byte) bool {
	return HasLockedBinary(name)
}

// GetHasLockedByBinary is an alias for HasLockedByBinary.
//
//	ok := covenant.GetHasLockedByBinary([]byte("example"))
func GetHasLockedByBinary(name []byte) bool {
	return HasLockedByBinary(name)
}

// GetLockedByString is an alias for GetLockedString.
//
//	item, ok := covenant.GetLockedByString("example")
func GetLockedByString(name string) (LockedName, bool) {
	return GetLockedString(name)
}

// GetLockedByBinary is an alias for GetLockedBinary.
//
//	item, ok := covenant.GetLockedByBinary([]byte("example"))
func GetLockedByBinary(name []byte) (LockedName, bool) {
	return GetLockedBinary(name)
}

// GetLockedHash returns the locked-name entry for a precomputed hash.
//
//	item, ok := covenant.GetLockedHash(hash)
func GetLockedHash(hash primitives.Hash) (LockedName, bool) {
	return defaultLockedCatalog().Get(hash)
}

// HasLockedHash reports whether a precomputed hash is present in the catalog.
//
//	ok := covenant.HasLockedHash(hash)
func HasLockedHash(hash primitives.Hash) bool {
	return DefaultLockedCatalog().Has(hash)
}

// GetHasLockedHash is an alias for HasLockedHash.
//
//	ok := covenant.GetHasLockedHash(hash)
func GetHasLockedHash(hash primitives.Hash) bool {
	return HasLockedHash(hash)
}

// GetLockedByHash is an alias for GetLockedHash.
//
//	item, ok := covenant.GetLockedByHash(hash)
func GetLockedByHash(hash primitives.Hash) (LockedName, bool) {
	return GetLockedHash(hash)
}

// HasLockedByHash is an alias for HasLockedHash.
//
//	ok := covenant.HasLockedByHash(hash)
func HasLockedByHash(hash primitives.Hash) bool {
	return HasLockedHash(hash)
}

// GetHasLockedByHash is an alias for HasLockedByHash.
//
//	ok := covenant.GetHasLockedByHash(hash)
func GetHasLockedByHash(hash primitives.Hash) bool {
	return HasLockedByHash(hash)
}

// DefaultLockedCatalog returns the lazily loaded locked-name catalog.
//
//	catalog := covenant.DefaultLockedCatalog()
func DefaultLockedCatalog() *LockedCatalog {
	return defaultLockedCatalog()
}

// GetLockedCatalog is an alias for DefaultLockedCatalog.
//
//	catalog := covenant.GetLockedCatalog()
func GetLockedCatalog() *LockedCatalog {
	return DefaultLockedCatalog()
}

// GetDefaultLockedCatalog is an alias for DefaultLockedCatalog.
//
//	catalog := covenant.GetDefaultLockedCatalog()
func GetDefaultLockedCatalog() *LockedCatalog {
	return DefaultLockedCatalog()
}

// Size reports the number of locked-name entries in the catalog.
//
//	size := catalog.Size()
func (c *LockedCatalog) Size() int {
	if c == nil {
		return 0
	}

	return len(c.byHash)
}

// GetSize is an alias for Size.
//
//	size := catalog.GetSize()
func (c *LockedCatalog) GetSize() int {
	return c.Size()
}

// PrefixSize reports the byte offset where the locked-name table entries begin.
//
//	offset := catalog.PrefixSize()
func (c *LockedCatalog) PrefixSize() int {
	return lockedCatalogPrefixSize
}

// GetPrefixSize is an alias for PrefixSize.
//
//	offset := catalog.GetPrefixSize()
func (c *LockedCatalog) GetPrefixSize() int {
	return c.PrefixSize()
}

// Has reports whether the catalog contains the hash.
//
//	ok := catalog.Has(hash)
func (c *LockedCatalog) Has(hash primitives.Hash) bool {
	if c == nil {
		return false
	}

	_, ok := c.byHash[hash]
	return ok
}

// HasByHash is an alias for Has.
//
//	ok := catalog.HasByHash(hash)
func (c *LockedCatalog) HasByHash(hash primitives.Hash) bool {
	return c.Has(hash)
}

// GetHasByHash is an alias for HasByHash.
//
//	ok := catalog.GetHasByHash(hash)
func (c *LockedCatalog) GetHasByHash(hash primitives.Hash) bool {
	return c.HasByHash(hash)
}

// Get returns the locked-name entry for the hash.
//
//	item, ok := catalog.Get(hash)
func (c *LockedCatalog) Get(hash primitives.Hash) (LockedName, bool) {
	if c == nil {
		return LockedName{}, false
	}

	item, ok := c.byHash[hash]
	return item, ok
}

// GetByHash is an alias for Get.
//
//	item, ok := catalog.GetByHash(hash)
func (c *LockedCatalog) GetByHash(hash primitives.Hash) (LockedName, bool) {
	return c.Get(hash)
}

// HasByName reports whether the catalog contains the provided label.
//
//	ok := catalog.HasByName("example")
func (c *LockedCatalog) HasByName(name string) bool {
	_, ok := c.GetByName(name)
	return ok
}

// GetHasByName is an alias for HasByName.
//
//	ok := catalog.GetHasByName("example")
func (c *LockedCatalog) GetHasByName(name string) bool {
	return c.HasByName(name)
}

// GetByName returns the locked-name entry for the provided label.
//
// The lookup lower-cases ASCII input, strips canonical `.lthn` suffixes, and
// hashes the normalized label directly.
//
//	item, ok := catalog.GetByName("example")
func (c *LockedCatalog) GetByName(name string) (LockedName, bool) {
	name, ok := asciiLowerName(name)
	if !ok {
		return LockedName{}, false
	}

	hash := primitives.Hash(core.SHA3_256([]byte(name)))
	return c.Get(hash)
}

// Entries returns the catalog contents as hash/item pairs.
//
//	entries := catalog.Entries()
func (c *LockedCatalog) Entries() []LockedEntry {
	if c == nil || len(c.byHash) == 0 {
		return nil
	}

	hashes := c.sortedHashes()
	entries := make([]LockedEntry, 0, len(hashes))

	for _, hash := range hashes {
		item, ok := c.byHash[hash]
		if !ok {
			continue
		}

		entries = append(entries, LockedEntry{Hash: hash, Value: item})
	}

	return entries
}

// GetEntries is an alias for Entries.
//
//	entries := catalog.GetEntries()
func (c *LockedCatalog) GetEntries() []LockedEntry {
	return c.Entries()
}

// Keys returns the locked-name hashes in deterministic order.
//
//	keys := catalog.Keys()
func (c *LockedCatalog) Keys() []primitives.Hash {
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
func (c *LockedCatalog) GetKeys() []primitives.Hash {
	return c.Keys()
}

// Values returns the locked-name entries in deterministic order.
//
//	values := catalog.Values()
func (c *LockedCatalog) Values() []LockedName {
	if c == nil || len(c.byHash) == 0 {
		return nil
	}

	hashes := c.sortedHashes()
	values := make([]LockedName, 0, len(hashes))

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
func (c *LockedCatalog) GetValues() []LockedName {
	return c.Values()
}

func (c *LockedCatalog) sortedHashes() []primitives.Hash {
	hashes := make([]primitives.Hash, 0, len(c.byHash))

	for hash := range c.byHash {
		hashes = append(hashes, hash)
	}

	sort.Slice(hashes, func(i, j int) bool {
		return compareBytes(hashes[i][:], hashes[j][:]) < 0
	})

	return hashes
}

func defaultLockedCatalog() *LockedCatalog {
	lockedCatalogOnce.Do(func() {
		catalog, err := loadLockedCatalog()
		if err != nil {
			lockedCatalog = &LockedCatalog{byHash: make(map[primitives.Hash]LockedName)}
			return
		}

		lockedCatalog = catalog
	})

	if lockedCatalog == nil {
		lockedCatalog = &LockedCatalog{byHash: make(map[primitives.Hash]LockedName)}
	}

	return lockedCatalog
}

func loadLockedCatalog() (*LockedCatalog, error) {
	data, err := readLockedData()
	if err != nil {
		return nil, err
	}

	var raw map[string]jsonraw
	if r := core.JSONUnmarshal(data, &raw); !r.OK {
		return nil, core.E("covenant.loadLockedCatalog", "failed to decode locked catalog", r.Value.(error))
	}

	catalog := &LockedCatalog{
		byHash: make(map[primitives.Hash]LockedName, len(raw)),
	}

	for key, value := range raw {
		if key == reservedZeroHash {
			continue
		}

		hash, err := decodeLockedHash(key)
		if err != nil {
			return nil, err
		}

		item, err := parseLockedEntry(hash, value)
		if err != nil {
			return nil, err
		}

		catalog.byHash[hash] = item
	}

	return catalog, nil
}

func parseLockedEntry(hash primitives.Hash, raw jsonraw) (LockedName, error) {
	var fields []jsonraw
	if r := core.JSONUnmarshal(raw, &fields); !r.OK {
		return LockedName{}, core.E("covenant.parseLockedEntry", "failed to decode locked entry", r.Value.(error))
	}

	if len(fields) != 2 {
		return LockedName{}, core.E("covenant.parseLockedEntry", "invalid entry size", nil)
	}

	var target string
	if r := core.JSONUnmarshal(fields[0], &target); !r.OK {
		return LockedName{}, core.E("covenant.parseLockedEntry", "failed to decode target", r.Value.(error))
	}

	parts := core.SplitN(target, ".", 2)
	if len(parts) != 2 || len(parts[0]) == 0 {
		return LockedName{}, core.E("covenant.parseLockedEntry", "invalid target", nil)
	}

	flags, err := parseLockedUint(fields[1])
	if err != nil {
		return LockedName{}, core.E("covenant.parseLockedEntry", "failed to decode flags", err)
	}

	item := LockedName{
		Name:   parts[0],
		Hash:   hash,
		Target: target,
		Root:   flags&1 != 0,
		Custom: flags&2 != 0,
	}

	return item, nil
}

func parseLockedUint(raw jsonraw) (uint64, error) {
	var num uint64
	if r := core.JSONUnmarshal(raw, &num); !r.OK {
		return 0, r.Value.(error)
	}

	return num, nil
}

func decodeLockedHash(key string) (primitives.Hash, error) {
	raw, err := hex.DecodeString(key)
	if err != nil {
		return primitives.Hash{}, core.E("covenant.decodeLockedHash", "failed to decode locked hash", err)
	}

	if len(raw) != len(primitives.Hash{}) {
		return primitives.Hash{}, core.E("covenant.decodeLockedHash", "invalid hash length", nil)
	}

	var hash primitives.Hash
	copy(hash[:], raw)
	return hash, nil
}

func readLockedData() ([]byte, error) {
	candidates := lockedDataPaths()
	fs := core.Fs{}

	for _, path := range candidates {
		if result := fs.Read(path); result.OK {
			return []byte(result.Value.(string)), nil
		}
	}

	return nil, core.E("covenant.readLockedData", "could not locate docs/js-covenants/lockup.json", nil)
}

func lockedDataPaths() []string {
	paths := []string{
		core.PathJoin("docs", "js-covenants", "lockup.json"),
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		base := core.PathDir(file)
		paths = append(paths,
			core.PathJoin(base, "..", "..", "docs", "js-covenants", "lockup.json"),
			core.PathJoin(base, "..", "..", "..", "docs", "js-covenants", "lockup.json"),
		)
	}

	return paths
}
