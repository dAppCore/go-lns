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

const reservedZeroHash = "0000000000000000000000000000000000000000000000000000000000000000"

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
type ReservedCatalog struct {
	byHash map[primitives.Hash]ReservedName
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
// Lookup is case-insensitive and matches the reference table's lower-casing
// behaviour.
func HasReservedName(name string) bool {
	return DefaultReservedCatalog().HasByName(name)
}

// GetReservedName returns the reserved-name entry for the provided label.
//
// The input is lower-cased before hashing, matching the JS reference helper.
func GetReservedName(name string) (ReservedName, bool) {
	return DefaultReservedCatalog().GetByName(name)
}

// GetReservedHash returns the reserved-name entry for a precomputed hash.
func GetReservedHash(hash primitives.Hash) (ReservedName, bool) {
	return defaultReservedCatalog().Get(hash)
}

// DefaultReservedCatalog returns the lazily loaded reserved-name catalog.
func DefaultReservedCatalog() *ReservedCatalog {
	return defaultReservedCatalog()
}

// Has reports whether the catalog contains the hash.
func (c *ReservedCatalog) Has(hash primitives.Hash) bool {
	if c == nil {
		return false
	}

	_, ok := c.byHash[hash]
	return ok
}

// Get returns the reserved-name entry for the hash.
func (c *ReservedCatalog) Get(hash primitives.Hash) (ReservedName, bool) {
	if c == nil {
		return ReservedName{}, false
	}

	item, ok := c.byHash[hash]
	return item, ok
}

// HasByName reports whether the catalog contains the provided label.
func (c *ReservedCatalog) HasByName(name string) bool {
	_, ok := c.GetByName(name)
	return ok
}

// GetByName returns the reserved-name entry for the provided label.
//
// The lookup lower-cases the input and hashes it directly, matching the JS
// reference helper's behavior.
func (c *ReservedCatalog) GetByName(name string) (ReservedName, bool) {
	name = strings.ToLower(name)
	if len(name) == 0 || len(name) > maxNameSize {
		return ReservedName{}, false
	}

	hash := primitives.Hash(sha3.Sum256([]byte(name)))
	return c.Get(hash)
}

// Entries returns the catalog contents as hash/item pairs.
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

// Keys returns the reserved-name hashes in deterministic order.
func (c *ReservedCatalog) Keys() []primitives.Hash {
	if c == nil || len(c.byHash) == 0 {
		return nil
	}

	hashes := c.sortedHashes()
	keys := make([]primitives.Hash, len(hashes))
	copy(keys, hashes)
	return keys
}

// Values returns the reserved-name entries in deterministic order.
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
		return nil, fmt.Errorf("covenant.loadReservedCatalog: %w", err)
	}

	catalog := &ReservedCatalog{
		byHash: make(map[primitives.Hash]ReservedName, len(raw)),
	}

	meta, err := parseReservedMeta(raw[reservedZeroHash])
	if err != nil {
		return nil, err
	}

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
		return reservedCatalogMeta{}, fmt.Errorf("covenant.parseReservedMeta: %w", err)
	}

	if len(values) != 4 {
		return reservedCatalogMeta{}, fmt.Errorf("covenant.parseReservedMeta: invalid metadata size %d", len(values))
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
		return ReservedName{}, fmt.Errorf("covenant.parseReservedEntry: %w", err)
	}

	if len(fields) < 2 || len(fields) > 3 {
		return ReservedName{}, fmt.Errorf("covenant.parseReservedEntry: invalid entry size %d", len(fields))
	}

	var target string
	if err := json.Unmarshal(fields[0], &target); err != nil {
		return ReservedName{}, fmt.Errorf("covenant.parseReservedEntry: %w", err)
	}

	dot := strings.IndexByte(target, '.')
	if dot <= 0 {
		return ReservedName{}, fmt.Errorf("covenant.parseReservedEntry: invalid target %q", target)
	}

	flags, err := parseReservedUint(fields[1])
	if err != nil {
		return ReservedName{}, fmt.Errorf("covenant.parseReservedEntry: %w", err)
	}

	var custom uint64
	if len(fields) == 3 {
		custom, err = parseReservedUint(fields[2])
		if err != nil {
			return ReservedName{}, fmt.Errorf("covenant.parseReservedEntry: %w", err)
		}
	}

	item := ReservedName{
		Name:   target[:dot],
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
		return primitives.Hash{}, fmt.Errorf("covenant.decodeReservedHash: %w", err)
	}

	if len(raw) != len(primitives.Hash{}) {
		return primitives.Hash{}, fmt.Errorf("covenant.decodeReservedHash: invalid hash length %d", len(raw))
	}

	var hash primitives.Hash
	copy(hash[:], raw)
	return hash, nil
}

func readReservedData() ([]byte, error) {
	candidates := reservedDataPaths()

	for _, path := range candidates {
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("covenant.readReservedData: could not locate docs/js-covenants/names.json")
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
