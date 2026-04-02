// SPDX-License-Identifier: EUPL-1.2

// Package lns is the Lethean Name System — a CoreGO native conversion of the
// HSD (Handshake) JavaScript sidechain. It provides name auctions, DNS
// resolution, and covenant-based name management for .lthn domains.
//
// Register LNS as a Core service:
//
//	c := core.New(core.WithService(lns.Register))
//
// Sub-packages (primitives, covenant, dns) provide the domain types
// and resolution logic used by the service.
package lns

import (
	"strings"

	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
)

// ServiceName is the registered name for the LNS service group.
const ServiceName = "lns"

// Service is the LNS service that manages .lthn name resolution,
// covenant processing, and name-chain state queries.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns")
type Service struct {
	*core.ServiceRuntime[serviceOptions]
}

// serviceOptions holds configuration for the LNS service.
type serviceOptions struct{}

// Register constructs an LNS Service bound to the provided Core instance.
// Pass to core.WithService during Core construction:
//
//	c := core.New(core.WithService(lns.Register))
func Register(c *core.Core) core.Result {
	svc := &Service{
		ServiceRuntime: core.NewServiceRuntime(c, serviceOptions{}),
	}
	return core.Result{Value: svc, OK: true}
}

// Resolve returns the canonical hash for a validated .lthn name.
//
// This mirrors the subpackage resolver while keeping the service-level API
// ergonomic for callers that only have an lns.Service instance.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.Resolve("example.lthn")
func (s *Service) Resolve(name any) (primitives.Hash, error) {
	switch v := name.(type) {
	case string:
		return s.ResolveString(v)
	case []byte:
		return s.ResolveBinary(v)
	default:
		return primitives.Hash{}, core.E("lns.Service.Resolve", "invalid name type", nil)
	}
}

// Verify reports whether a .lthn name is valid after canonicalisation.
//
// This keeps the service-level API aligned with the DNS resolver helpers.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.Verify("example.lthn")
func (s *Service) Verify(name any) bool {
	switch v := name.(type) {
	case string:
		return s.VerifyString(v)
	case []byte:
		return s.VerifyBinary(v)
	default:
		return false
	}
}

// ResolveString returns the canonical hash for a .lthn name supplied as a string.
func (s *Service) ResolveString(name string) (primitives.Hash, error) {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return primitives.Hash{}, core.E("lns.Service.ResolveString", "invalid name", nil)
	}

	return covenant.HashString(normalized)
}

// ResolveBinary returns the canonical hash for a .lthn name supplied as bytes.
func (s *Service) ResolveBinary(name []byte) (primitives.Hash, error) {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return primitives.Hash{}, core.E("lns.Service.ResolveBinary", "invalid name", nil)
	}

	return covenant.HashString(normalized)
}

// ResolveName is an alias for Resolve.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.ResolveName("example.lthn")
func (s *Service) ResolveName(name any) (primitives.Hash, error) {
	return s.Resolve(name)
}

// VerifyName is an alias for Verify.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.VerifyName("example.lthn")
func (s *Service) VerifyName(name any) bool {
	return s.Verify(name)
}

// VerifyString reports whether a .lthn name supplied as a string is valid.
func (s *Service) VerifyString(name string) bool {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return false
	}

	return covenant.VerifyString(normalized)
}

// VerifyBinary reports whether a .lthn name supplied as bytes is valid.
func (s *Service) VerifyBinary(name []byte) bool {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return false
	}

	return covenant.VerifyString(normalized)
}

// GetReserved looks up a reserved .lthn name after canonicalisation.
//
// The lookup accepts the same inputs as Resolve and Verify, then returns the
// matching reserved-name entry from the reference catalog when present.
func (s *Service) GetReserved(name any) (covenant.ReservedName, bool) {
	hash, err := s.Resolve(name)
	if err != nil {
		return covenant.ReservedName{}, false
	}

	return covenant.GetReservedHash(hash)
}

// GetReservedName looks up a reserved name using the raw catalog label.
//
// The helper accepts the bare label first, then falls back to canonical
// `.lthn` resolution so existing callers can pass either form.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	item, ok := svc.GetReservedName("reserved")
func (s *Service) GetReservedName(name any) (covenant.ReservedName, bool) {
	if label, ok := catalogLabel(name); ok {
		if item, ok := covenant.GetReservedName(label); ok {
			return item, true
		}
	}

	return s.GetReserved(name)
}

// HasReserved reports whether a canonicalised .lthn name is reserved.
//
// This mirrors the covenant package lookup helper without exposing the hash
// plumbing to service-level callers.
func (s *Service) HasReserved(name any) bool {
	_, ok := s.GetReserved(name)
	return ok
}

// GetReservedByName is an alias for GetReservedName.
//
// This keeps the service-level name lookup helpers consistent: callers can
// pass either a raw catalog label or a canonical `.lthn` name.
func (s *Service) GetReservedByName(name any) (covenant.ReservedName, bool) {
	return s.GetReservedName(name)
}

// HasReservedByName is an alias for HasReservedName.
func (s *Service) HasReservedByName(name any) bool {
	return s.HasReservedName(name)
}

// HasReservedName reports whether a reserved name exists using the raw
// catalog label.
func (s *Service) HasReservedName(name any) bool {
	_, ok := s.GetReservedName(name)
	return ok
}

// GetReservedHash looks up a reserved .lthn name by its canonical hash.
//
// This mirrors the covenant package lookup helper for callers that already
// have the hashed form available.
func (s *Service) GetReservedHash(hash primitives.Hash) (covenant.ReservedName, bool) {
	return covenant.GetReservedHash(hash)
}

// HasReservedHash reports whether a canonical hash is reserved.
func (s *Service) HasReservedHash(hash primitives.Hash) bool {
	_, ok := s.GetReservedHash(hash)
	return ok
}

// GetReservedByHash is an alias for GetReservedHash.
func (s *Service) GetReservedByHash(hash primitives.Hash) (covenant.ReservedName, bool) {
	return s.GetReservedHash(hash)
}

// HasReservedByHash is an alias for HasReservedHash.
func (s *Service) HasReservedByHash(hash primitives.Hash) bool {
	return s.HasReservedHash(hash)
}

// GetLocked looks up a locked .lthn name after canonicalisation.
//
// The lookup accepts the same inputs as Resolve and Verify, then returns the
// matching locked-name entry from the reference catalog when present.
func (s *Service) GetLocked(name any) (covenant.LockedName, bool) {
	hash, err := s.Resolve(name)
	if err != nil {
		return covenant.LockedName{}, false
	}

	return covenant.GetLockedHash(hash)
}

// GetLockedName looks up a locked name using the raw catalog label.
//
// The helper accepts the bare label first, then falls back to canonical
// `.lthn` resolution so existing callers can pass either form.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	item, ok := svc.GetLockedName("nec")
func (s *Service) GetLockedName(name any) (covenant.LockedName, bool) {
	if label, ok := catalogLabel(name); ok {
		if item, ok := covenant.GetLockedName(label); ok {
			return item, true
		}
	}

	return s.GetLocked(name)
}

// HasLocked reports whether a canonicalised .lthn name is locked.
//
// This mirrors the covenant package lookup helper without exposing the hash
// plumbing to service-level callers.
func (s *Service) HasLocked(name any) bool {
	_, ok := s.GetLocked(name)
	return ok
}

// GetLockedByName is an alias for GetLockedName.
//
// This keeps the service-level name lookup helpers consistent: callers can
// pass either a raw catalog label or a canonical `.lthn` name.
func (s *Service) GetLockedByName(name any) (covenant.LockedName, bool) {
	return s.GetLockedName(name)
}

// HasLockedByName is an alias for HasLockedName.
func (s *Service) HasLockedByName(name any) bool {
	return s.HasLockedName(name)
}

// HasLockedName reports whether a locked name exists using the raw catalog
// label.
func (s *Service) HasLockedName(name any) bool {
	_, ok := s.GetLockedName(name)
	return ok
}

// GetLockedHash looks up a locked .lthn name by its canonical hash.
//
// This mirrors the covenant package lookup helper for callers that already
// have the hashed form available.
func (s *Service) GetLockedHash(hash primitives.Hash) (covenant.LockedName, bool) {
	return covenant.GetLockedHash(hash)
}

// HasLockedHash reports whether a canonical hash is locked.
func (s *Service) HasLockedHash(hash primitives.Hash) bool {
	_, ok := s.GetLockedHash(hash)
	return ok
}

// GetLockedByHash is an alias for GetLockedHash.
func (s *Service) GetLockedByHash(hash primitives.Hash) (covenant.LockedName, bool) {
	return s.GetLockedHash(hash)
}

// HasLockedByHash is an alias for HasLockedHash.
func (s *Service) HasLockedByHash(hash primitives.Hash) bool {
	return s.HasLockedHash(hash)
}

// ReservedCatalog exposes the reserved-name catalog used by the service.
//
// The returned catalog provides deterministic iteration helpers so callers can
// enumerate the reference set without reaching into the covenant package.
func (s *Service) ReservedCatalog() *covenant.ReservedCatalog {
	return covenant.DefaultReservedCatalog()
}

// ReservedSize reports the number of reserved-name entries available to the service.
func (s *Service) ReservedSize() int {
	return s.ReservedCatalog().Size()
}

// ReservedNameValue reports the base value used for reserved-name lookups.
func (s *Service) ReservedNameValue() uint64 {
	return s.ReservedCatalog().NameValue()
}

// ReservedRootValue reports the root-name increment used for reserved-name lookups.
func (s *Service) ReservedRootValue() uint64 {
	return s.ReservedCatalog().RootValue()
}

// ReservedTopValue reports the top-100 increment used for reserved-name lookups.
func (s *Service) ReservedTopValue() uint64 {
	return s.ReservedCatalog().TopValue()
}

// ReservedEntries exposes the reserved-name catalog entries used by the service.
func (s *Service) ReservedEntries() []covenant.ReservedEntry {
	return s.ReservedCatalog().Entries()
}

// ReservedKeys exposes the reserved-name hashes used by the service.
func (s *Service) ReservedKeys() []primitives.Hash {
	return s.ReservedCatalog().Keys()
}

// ReservedValues exposes the reserved-name entries used by the service.
func (s *Service) ReservedValues() []covenant.ReservedName {
	return s.ReservedCatalog().Values()
}

// LockedCatalog exposes the locked-name catalog used by the service.
//
// The returned catalog provides deterministic iteration helpers so callers can
// enumerate the reference set without reaching into the covenant package.
func (s *Service) LockedCatalog() *covenant.LockedCatalog {
	return covenant.DefaultLockedCatalog()
}

// LockedSize reports the number of locked-name entries available to the service.
func (s *Service) LockedSize() int {
	return s.LockedCatalog().Size()
}

// LockedPrefixSize reports the byte offset where locked-name entries begin in
// the reference table.
func (s *Service) LockedPrefixSize() int {
	return s.LockedCatalog().PrefixSize()
}

// ReservedPrefixSize reports the byte offset where reserved-name entries begin in
// the reference table.
func (s *Service) ReservedPrefixSize() int {
	return s.ReservedCatalog().PrefixSize()
}

// LockedEntries exposes the locked-name catalog entries used by the service.
func (s *Service) LockedEntries() []covenant.LockedEntry {
	return s.LockedCatalog().Entries()
}

// LockedKeys exposes the locked-name hashes used by the service.
func (s *Service) LockedKeys() []primitives.Hash {
	return s.LockedCatalog().Keys()
}

// LockedValues exposes the locked-name entries used by the service.
func (s *Service) LockedValues() []covenant.LockedName {
	return s.LockedCatalog().Values()
}

func catalogLabel(name any) (string, bool) {
	switch v := name.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	default:
		return "", false
	}
}

func canonicalizeName(name any) (string, bool) {
	var value string

	switch v := name.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	default:
		return "", false
	}

	value = strings.ToLower(value)

	switch {
	case strings.HasSuffix(value, ".lthn."):
		value = strings.TrimSuffix(value, ".lthn.")
	case strings.HasSuffix(value, ".lthn"):
		value = strings.TrimSuffix(value, ".lthn")
	case strings.HasSuffix(value, "."):
		value = strings.TrimSuffix(value, ".")
	}

	return value, true
}
