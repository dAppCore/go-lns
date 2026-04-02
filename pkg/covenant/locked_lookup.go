// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"bytes"
	"crypto/sha3"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

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
func HasLockedName(name string) bool {
	return DefaultLockedCatalog().HasByName(name)
}

// HasLockedByName is an alias for HasLockedName.
func HasLockedByName(name string) bool {
	return HasLockedName(name)
}

// GetLockedName returns the locked-name entry for the provided label.
func GetLockedName(name string) (LockedName, bool) {
	return DefaultLockedCatalog().GetByName(name)
}

// GetLockedByName is an alias for GetLockedName.
func GetLockedByName(name string) (LockedName, bool) {
	return GetLockedName(name)
}

// GetLockedHash returns the locked-name entry for a precomputed hash.
func GetLockedHash(hash primitives.Hash) (LockedName, bool) {
	return defaultLockedCatalog().Get(hash)
}

// GetLockedByHash is an alias for GetLockedHash.
func GetLockedByHash(hash primitives.Hash) (LockedName, bool) {
	return GetLockedHash(hash)
}

// DefaultLockedCatalog returns the lazily loaded locked-name catalog.
func DefaultLockedCatalog() *LockedCatalog {
	return defaultLockedCatalog()
}

// Size reports the number of locked-name entries in the catalog.
func (c *LockedCatalog) Size() int {
	if c == nil {
		return 0
	}

	return len(c.byHash)
}

// PrefixSize reports the byte offset where the locked-name table entries begin.
func (c *LockedCatalog) PrefixSize() int {
	return lockedCatalogPrefixSize
}

// Has reports whether the catalog contains the hash.
func (c *LockedCatalog) Has(hash primitives.Hash) bool {
	if c == nil {
		return false
	}

	_, ok := c.byHash[hash]
	return ok
}

// HasByHash is an alias for Has.
func (c *LockedCatalog) HasByHash(hash primitives.Hash) bool {
	return c.Has(hash)
}

// Get returns the locked-name entry for the hash.
func (c *LockedCatalog) Get(hash primitives.Hash) (LockedName, bool) {
	if c == nil {
		return LockedName{}, false
	}

	item, ok := c.byHash[hash]
	return item, ok
}

// GetByHash is an alias for Get.
func (c *LockedCatalog) GetByHash(hash primitives.Hash) (LockedName, bool) {
	return c.Get(hash)
}

// HasByName reports whether the catalog contains the provided label.
func (c *LockedCatalog) HasByName(name string) bool {
	_, ok := c.GetByName(name)
	return ok
}

// GetByName returns the locked-name entry for the provided label.
//
// The lookup lower-cases the input and hashes it directly, matching the JS
// reference helper's behavior.
func (c *LockedCatalog) GetByName(name string) (LockedName, bool) {
	name = strings.ToLower(name)
	if len(name) == 0 || len(name) > maxNameSize {
		return LockedName{}, false
	}

	hash := primitives.Hash(sha3.Sum256([]byte(name)))
	return c.Get(hash)
}

// Entries returns the catalog contents as hash/item pairs.
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

// Keys returns the locked-name hashes in deterministic order.
func (c *LockedCatalog) Keys() []primitives.Hash {
	if c == nil || len(c.byHash) == 0 {
		return nil
	}

	hashes := c.sortedHashes()
	keys := make([]primitives.Hash, len(hashes))
	copy(keys, hashes)
	return keys
}

// Values returns the locked-name entries in deterministic order.
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

func (c *LockedCatalog) sortedHashes() []primitives.Hash {
	hashes := make([]primitives.Hash, 0, len(c.byHash))

	for hash := range c.byHash {
		hashes = append(hashes, hash)
	}

	sort.Slice(hashes, func(i, j int) bool {
		return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
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

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("covenant.loadLockedCatalog: %w", err)
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

func parseLockedEntry(hash primitives.Hash, raw json.RawMessage) (LockedName, error) {
	var fields []json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return LockedName{}, fmt.Errorf("covenant.parseLockedEntry: %w", err)
	}

	if len(fields) != 2 {
		return LockedName{}, fmt.Errorf("covenant.parseLockedEntry: invalid entry size %d", len(fields))
	}

	var target string
	if err := json.Unmarshal(fields[0], &target); err != nil {
		return LockedName{}, fmt.Errorf("covenant.parseLockedEntry: %w", err)
	}

	dot := strings.IndexByte(target, '.')
	if dot <= 0 {
		return LockedName{}, fmt.Errorf("covenant.parseLockedEntry: invalid target %q", target)
	}

	flags, err := parseLockedUint(fields[1])
	if err != nil {
		return LockedName{}, fmt.Errorf("covenant.parseLockedEntry: %w", err)
	}

	item := LockedName{
		Name:   target[:dot],
		Hash:   hash,
		Target: target,
		Root:   flags&1 != 0,
		Custom: flags&2 != 0,
	}

	return item, nil
}

func parseLockedUint(raw json.RawMessage) (uint64, error) {
	var num uint64
	if err := json.Unmarshal(raw, &num); err != nil {
		return 0, err
	}

	return num, nil
}

func decodeLockedHash(key string) (primitives.Hash, error) {
	raw, err := hex.DecodeString(key)
	if err != nil {
		return primitives.Hash{}, fmt.Errorf("covenant.decodeLockedHash: %w", err)
	}

	if len(raw) != len(primitives.Hash{}) {
		return primitives.Hash{}, fmt.Errorf("covenant.decodeLockedHash: invalid hash length %d", len(raw))
	}

	var hash primitives.Hash
	copy(hash[:], raw)
	return hash, nil
}

func readLockedData() ([]byte, error) {
	candidates := lockedDataPaths()

	for _, path := range candidates {
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("covenant.readLockedData: could not locate docs/js-covenants/lockup.json")
}

func lockedDataPaths() []string {
	paths := []string{
		filepath.Join("docs", "js-covenants", "lockup.json"),
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		base := filepath.Dir(file)
		paths = append(paths,
			filepath.Join(base, "..", "..", "docs", "js-covenants", "lockup.json"),
			filepath.Join(base, "..", "..", "..", "docs", "js-covenants", "lockup.json"),
		)
	}

	return paths
}
