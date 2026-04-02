// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	_, ok := GetReservedName(name)
	return ok
}

// GetReservedName returns the reserved-name entry for the provided label.
//
// The input is lower-cased before validation and hashing, matching the JS
// reference helper.
func GetReservedName(name string) (ReservedName, bool) {
	name = strings.ToLower(name)

	if !VerifyString(name) {
		return ReservedName{}, false
	}

	hash, err := HashString(name)
	if err != nil {
		return ReservedName{}, false
	}

	return defaultReservedCatalog().Get(hash)
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
