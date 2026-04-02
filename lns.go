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
	core "dappco.re/go/core"
	"dappco.re/go/lns/internal/nameutil"
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

// Hash returns the canonical hash for a validated .lthn name.
//
// Hash is the direct service-level alias for Resolve.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.Hash("example.lthn")
func (s *Service) Hash(name any) (primitives.Hash, error) {
	return s.Resolve(name)
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

// HashString returns the canonical hash for a .lthn name supplied as a string.
//
//	ok := svc.HashString("example.lthn")
func (s *Service) HashString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// HashBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	ok := svc.HashBinary([]byte("example.lthn"))
func (s *Service) HashBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// ResolveByString is an alias for ResolveString.
//
//	ok := svc.ResolveByString("example.lthn")
func (s *Service) ResolveByString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// ResolveByBinary is an alias for ResolveBinary.
//
//	ok := svc.ResolveByBinary([]byte("example.lthn"))
func (s *Service) ResolveByBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// HashByString is an alias for HashString.
//
//	ok := svc.HashByString("example.lthn")
func (s *Service) HashByString(name string) (primitives.Hash, error) {
	return s.HashString(name)
}

// HashByBinary is an alias for HashBinary.
//
//	ok := svc.HashByBinary([]byte("example.lthn"))
func (s *Service) HashByBinary(name []byte) (primitives.Hash, error) {
	return s.HashBinary(name)
}

// ResolveString returns the canonical hash for a .lthn name supplied as a string.
func (s *Service) ResolveString(name string) (primitives.Hash, error) {
	normalized, ok := nameutil.Canonicalize(name)
	if !ok {
		return primitives.Hash{}, core.E("lns.Service.ResolveString", "invalid name", nil)
	}

	return covenant.HashString(normalized)
}

// ResolveBinary returns the canonical hash for a .lthn name supplied as bytes.
func (s *Service) ResolveBinary(name []byte) (primitives.Hash, error) {
	normalized, ok := nameutil.Canonicalize(name)
	if !ok {
		return primitives.Hash{}, core.E("lns.Service.ResolveBinary", "invalid name", nil)
	}

	return covenant.HashBinary([]byte(normalized))
}

// ResolveName is an alias for Resolve.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.ResolveName("example.lthn")
func (s *Service) ResolveName(name any) (primitives.Hash, error) {
	return s.Resolve(name)
}

// HashName is an alias for Hash.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.HashName("example.lthn")
func (s *Service) HashName(name any) (primitives.Hash, error) {
	return s.Hash(name)
}

// VerifyName is an alias for Verify.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.VerifyName("example.lthn")
func (s *Service) VerifyName(name any) bool {
	return s.Verify(name)
}

// ResolveByName is an alias for ResolveName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.ResolveByName("example.lthn")
func (s *Service) ResolveByName(name any) (primitives.Hash, error) {
	return s.ResolveName(name)
}

// HashByName is an alias for HashName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.HashByName("example.lthn")
func (s *Service) HashByName(name any) (primitives.Hash, error) {
	return s.HashName(name)
}

// VerifyByName is an alias for VerifyName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.VerifyByName("example.lthn")
func (s *Service) VerifyByName(name any) bool {
	return s.VerifyName(name)
}

// VerifyString reports whether a .lthn name supplied as a string is valid.
//
//	ok := svc.VerifyString("example.lthn")
func (s *Service) VerifyString(name string) bool {
	normalized, ok := nameutil.Canonicalize(name)
	if !ok {
		return false
	}

	return covenant.VerifyString(normalized)
}

// VerifyBinary reports whether a .lthn name supplied as bytes is valid.
//
//	ok := svc.VerifyBinary([]byte("example.lthn"))
func (s *Service) VerifyBinary(name []byte) bool {
	normalized, ok := nameutil.Canonicalize(name)
	if !ok {
		return false
	}

	return covenant.VerifyBinary([]byte(normalized))
}

// VerifyByString is an alias for VerifyString.
//
//	ok := svc.VerifyByString("example.lthn")
func (s *Service) VerifyByString(name string) bool {
	return s.VerifyString(name)
}

// VerifyByBinary is an alias for VerifyBinary.
//
//	ok := svc.VerifyByBinary([]byte("example.lthn"))
func (s *Service) VerifyByBinary(name []byte) bool {
	return s.VerifyBinary(name)
}

// GetReserved looks up a reserved .lthn name after canonicalisation.
//
// The lookup accepts the same inputs as Resolve and Verify, then returns the
// matching reserved-name entry from the reference catalog when present.
//
//	item, ok := svc.GetReserved("reserved.lthn")
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
	if label, ok := nameutil.CatalogLabel(name); ok {
		if item, ok := covenant.GetReservedName(label); ok {
			return item, true
		}
	}

	return s.GetReserved(name)
}

// GetReservedString looks up a reserved name supplied as a string.
//
// The helper accepts either a raw catalog label or a canonical `.lthn` name.
//
//	item, ok := svc.GetReservedString("reserved")
func (s *Service) GetReservedString(name string) (covenant.ReservedName, bool) {
	return s.GetReservedName(name)
}

// GetReservedBinary looks up a reserved name supplied as bytes.
//
//	item, ok := svc.GetReservedBinary([]byte("reserved"))
func (s *Service) GetReservedBinary(name []byte) (covenant.ReservedName, bool) {
	return s.GetReservedName(name)
}

// GetReservedByString is an alias for GetReservedString.
//
//	item, ok := svc.GetReservedByString("reserved")
func (s *Service) GetReservedByString(name string) (covenant.ReservedName, bool) {
	return s.GetReservedString(name)
}

// GetReservedByBinary is an alias for GetReservedBinary.
//
//	item, ok := svc.GetReservedByBinary([]byte("reserved"))
func (s *Service) GetReservedByBinary(name []byte) (covenant.ReservedName, bool) {
	return s.GetReservedBinary(name)
}

// HasReserved reports whether a canonicalised .lthn name is reserved.
//
// This mirrors the covenant package lookup helper without exposing the hash
// plumbing to service-level callers.
//
//	ok := svc.HasReserved("reserved.lthn")
func (s *Service) HasReserved(name any) bool {
	_, ok := s.GetReserved(name)
	return ok
}

// HasReservedString reports whether a reserved name supplied as a string exists.
//
//	ok := svc.HasReservedString("reserved")
func (s *Service) HasReservedString(name string) bool {
	_, ok := s.GetReservedName(name)
	return ok
}

// HasReservedBinary reports whether a reserved name supplied as bytes exists.
//
//	ok := svc.HasReservedBinary([]byte("reserved"))
func (s *Service) HasReservedBinary(name []byte) bool {
	_, ok := s.GetReservedName(name)
	return ok
}

// HasReservedByString is an alias for HasReservedString.
//
//	ok := svc.HasReservedByString("reserved")
func (s *Service) HasReservedByString(name string) bool {
	return s.HasReservedString(name)
}

// HasReservedByBinary is an alias for HasReservedBinary.
//
//	ok := svc.HasReservedByBinary([]byte("reserved"))
func (s *Service) HasReservedByBinary(name []byte) bool {
	return s.HasReservedBinary(name)
}

// GetReservedByName is an alias for GetReservedName.
//
// This keeps the service-level name lookup helpers consistent: callers can
// pass either a raw catalog label or a canonical `.lthn` name.
//
//	item, ok := svc.GetReservedByName("reserved")
func (s *Service) GetReservedByName(name any) (covenant.ReservedName, bool) {
	return s.GetReservedName(name)
}

// HasReservedByName is an alias for HasReservedName.
//
//	ok := svc.HasReservedByName("reserved")
func (s *Service) HasReservedByName(name any) bool {
	return s.HasReservedName(name)
}

// HasReservedName reports whether a reserved name exists using the raw
// catalog label.
//
//	ok := svc.HasReservedName("reserved")
func (s *Service) HasReservedName(name any) bool {
	_, ok := s.GetReservedName(name)
	return ok
}

// GetReservedHash looks up a reserved .lthn name by its canonical hash.
//
// This mirrors the covenant package lookup helper for callers that already
// have the hashed form available.
//
//	hash, _ := covenant.HashString("reserved")
//	item, ok := svc.GetReservedHash(hash)
func (s *Service) GetReservedHash(hash primitives.Hash) (covenant.ReservedName, bool) {
	return covenant.GetReservedHash(hash)
}

// HasReservedHash reports whether a canonical hash is reserved.
//
//	ok := svc.HasReservedHash(hash)
func (s *Service) HasReservedHash(hash primitives.Hash) bool {
	_, ok := s.GetReservedHash(hash)
	return ok
}

// GetReservedByHash is an alias for GetReservedHash.
//
//	item, ok := svc.GetReservedByHash(hash)
func (s *Service) GetReservedByHash(hash primitives.Hash) (covenant.ReservedName, bool) {
	return s.GetReservedHash(hash)
}

// HasReservedByHash is an alias for HasReservedHash.
//
//	ok := svc.HasReservedByHash(hash)
func (s *Service) HasReservedByHash(hash primitives.Hash) bool {
	return s.HasReservedHash(hash)
}

// GetLocked looks up a locked .lthn name after canonicalisation.
//
// The lookup accepts the same inputs as Resolve and Verify, then returns the
// matching locked-name entry from the reference catalog when present.
//
//	item, ok := svc.GetLocked("nec.lthn")
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
	if label, ok := nameutil.CatalogLabel(name); ok {
		if item, ok := covenant.GetLockedName(label); ok {
			return item, true
		}
	}

	return s.GetLocked(name)
}

// GetLockedString looks up a locked name supplied as a string.
//
// The helper accepts either a raw catalog label or a canonical `.lthn` name.
//
//	item, ok := svc.GetLockedString("nec")
func (s *Service) GetLockedString(name string) (covenant.LockedName, bool) {
	return s.GetLockedName(name)
}

// GetLockedBinary looks up a locked name supplied as bytes.
//
//	item, ok := svc.GetLockedBinary([]byte("nec"))
func (s *Service) GetLockedBinary(name []byte) (covenant.LockedName, bool) {
	return s.GetLockedName(name)
}

// GetLockedByString is an alias for GetLockedString.
//
//	item, ok := svc.GetLockedByString("nec")
func (s *Service) GetLockedByString(name string) (covenant.LockedName, bool) {
	return s.GetLockedString(name)
}

// GetLockedByBinary is an alias for GetLockedBinary.
//
//	item, ok := svc.GetLockedByBinary([]byte("nec"))
func (s *Service) GetLockedByBinary(name []byte) (covenant.LockedName, bool) {
	return s.GetLockedBinary(name)
}

// HasLocked reports whether a canonicalised .lthn name is locked.
//
// This mirrors the covenant package lookup helper without exposing the hash
// plumbing to service-level callers.
//
//	ok := svc.HasLocked("nec.lthn")
func (s *Service) HasLocked(name any) bool {
	_, ok := s.GetLocked(name)
	return ok
}

// HasLockedString reports whether a locked name supplied as a string exists.
//
//	ok := svc.HasLockedString("nec")
func (s *Service) HasLockedString(name string) bool {
	_, ok := s.GetLockedName(name)
	return ok
}

// HasLockedBinary reports whether a locked name supplied as bytes exists.
//
//	ok := svc.HasLockedBinary([]byte("nec"))
func (s *Service) HasLockedBinary(name []byte) bool {
	_, ok := s.GetLockedName(name)
	return ok
}

// HasLockedByString is an alias for HasLockedString.
//
//	ok := svc.HasLockedByString("nec")
func (s *Service) HasLockedByString(name string) bool {
	return s.HasLockedString(name)
}

// HasLockedByBinary is an alias for HasLockedBinary.
//
//	ok := svc.HasLockedByBinary([]byte("nec"))
func (s *Service) HasLockedByBinary(name []byte) bool {
	return s.HasLockedBinary(name)
}

// GetLockedByName is an alias for GetLockedName.
//
// This keeps the service-level name lookup helpers consistent: callers can
// pass either a raw catalog label or a canonical `.lthn` name.
//
//	item, ok := svc.GetLockedByName("nec")
func (s *Service) GetLockedByName(name any) (covenant.LockedName, bool) {
	return s.GetLockedName(name)
}

// HasLockedByName is an alias for HasLockedName.
//
//	ok := svc.HasLockedByName("nec")
func (s *Service) HasLockedByName(name any) bool {
	return s.HasLockedName(name)
}

// HasLockedName reports whether a locked name exists using the raw catalog
// label.
//
//	ok := svc.HasLockedName("nec")
func (s *Service) HasLockedName(name any) bool {
	_, ok := s.GetLockedName(name)
	return ok
}

// GetLockedHash looks up a locked .lthn name by its canonical hash.
//
// This mirrors the covenant package lookup helper for callers that already
// have the hashed form available.
//
//	hash, _ := covenant.HashString("nec")
//	item, ok := svc.GetLockedHash(hash)
func (s *Service) GetLockedHash(hash primitives.Hash) (covenant.LockedName, bool) {
	return covenant.GetLockedHash(hash)
}

// HasLockedHash reports whether a canonical hash is locked.
//
//	ok := svc.HasLockedHash(hash)
func (s *Service) HasLockedHash(hash primitives.Hash) bool {
	_, ok := s.GetLockedHash(hash)
	return ok
}

// GetLockedByHash is an alias for GetLockedHash.
//
//	item, ok := svc.GetLockedByHash(hash)
func (s *Service) GetLockedByHash(hash primitives.Hash) (covenant.LockedName, bool) {
	return s.GetLockedHash(hash)
}

// HasLockedByHash is an alias for HasLockedHash.
//
//	ok := svc.HasLockedByHash(hash)
func (s *Service) HasLockedByHash(hash primitives.Hash) bool {
	return s.HasLockedHash(hash)
}

// ReservedCatalog exposes the reserved-name catalog used by the service.
//
// The returned catalog provides deterministic iteration helpers so callers can
// enumerate the reference set without reaching into the covenant package.
//
//	catalog := svc.ReservedCatalog()
//	items := catalog.Entries()
func (s *Service) ReservedCatalog() *covenant.ReservedCatalog {
	return covenant.DefaultReservedCatalog()
}

// GetReservedCatalog is an alias for ReservedCatalog.
//
//	catalog := svc.GetReservedCatalog()
func (s *Service) GetReservedCatalog() *covenant.ReservedCatalog {
	return s.ReservedCatalog()
}

// ReservedSize reports the number of reserved-name entries available to the service.
//
//	size := svc.ReservedSize()
func (s *Service) ReservedSize() int {
	return s.ReservedCatalog().Size()
}

// GetReservedSize is an alias for ReservedSize.
//
//	size := svc.GetReservedSize()
func (s *Service) GetReservedSize() int {
	return s.ReservedSize()
}

// ReservedNameValue reports the base value used for reserved-name lookups.
//
//	base := svc.ReservedNameValue()
func (s *Service) ReservedNameValue() uint64 {
	return s.ReservedCatalog().NameValue()
}

// GetReservedNameValue is an alias for ReservedNameValue.
//
//	base := svc.GetReservedNameValue()
func (s *Service) GetReservedNameValue() uint64 {
	return s.ReservedNameValue()
}

// ReservedRootValue reports the root-name increment used for reserved-name lookups.
//
//	root := svc.ReservedRootValue()
func (s *Service) ReservedRootValue() uint64 {
	return s.ReservedCatalog().RootValue()
}

// GetReservedRootValue is an alias for ReservedRootValue.
//
//	root := svc.GetReservedRootValue()
func (s *Service) GetReservedRootValue() uint64 {
	return s.ReservedRootValue()
}

// ReservedTopValue reports the top-100 increment used for reserved-name lookups.
//
//	top := svc.ReservedTopValue()
func (s *Service) ReservedTopValue() uint64 {
	return s.ReservedCatalog().TopValue()
}

// GetReservedTopValue is an alias for ReservedTopValue.
//
//	top := svc.GetReservedTopValue()
func (s *Service) GetReservedTopValue() uint64 {
	return s.ReservedTopValue()
}

// ReservedEntries exposes the reserved-name catalog entries used by the service.
//
//	entries := svc.ReservedEntries()
func (s *Service) ReservedEntries() []covenant.ReservedEntry {
	return s.ReservedCatalog().Entries()
}

// GetReservedEntries is an alias for ReservedEntries.
//
//	entries := svc.GetReservedEntries()
func (s *Service) GetReservedEntries() []covenant.ReservedEntry {
	return s.ReservedEntries()
}

// ReservedKeys exposes the reserved-name hashes used by the service.
//
//	keys := svc.ReservedKeys()
func (s *Service) ReservedKeys() []primitives.Hash {
	return s.ReservedCatalog().Keys()
}

// GetReservedKeys is an alias for ReservedKeys.
//
//	keys := svc.GetReservedKeys()
func (s *Service) GetReservedKeys() []primitives.Hash {
	return s.ReservedKeys()
}

// ReservedValues exposes the reserved-name entries used by the service.
//
//	values := svc.ReservedValues()
func (s *Service) ReservedValues() []covenant.ReservedName {
	return s.ReservedCatalog().Values()
}

// GetReservedValues is an alias for ReservedValues.
//
//	values := svc.GetReservedValues()
func (s *Service) GetReservedValues() []covenant.ReservedName {
	return s.ReservedValues()
}

// LockedCatalog exposes the locked-name catalog used by the service.
//
// The returned catalog provides deterministic iteration helpers so callers can
// enumerate the reference set without reaching into the covenant package.
//
//	catalog := svc.LockedCatalog()
//	items := catalog.Entries()
func (s *Service) LockedCatalog() *covenant.LockedCatalog {
	return covenant.DefaultLockedCatalog()
}

// GetLockedCatalog is an alias for LockedCatalog.
//
//	catalog := svc.GetLockedCatalog()
func (s *Service) GetLockedCatalog() *covenant.LockedCatalog {
	return s.LockedCatalog()
}

// LockedSize reports the number of locked-name entries available to the service.
//
//	size := svc.LockedSize()
func (s *Service) LockedSize() int {
	return s.LockedCatalog().Size()
}

// GetLockedSize is an alias for LockedSize.
//
//	size := svc.GetLockedSize()
func (s *Service) GetLockedSize() int {
	return s.LockedSize()
}

// LockedPrefixSize reports the byte offset where locked-name entries begin in
// the reference table.
//
//	offset := svc.LockedPrefixSize()
func (s *Service) LockedPrefixSize() int {
	return s.LockedCatalog().PrefixSize()
}

// GetLockedPrefixSize is an alias for LockedPrefixSize.
//
//	offset := svc.GetLockedPrefixSize()
func (s *Service) GetLockedPrefixSize() int {
	return s.LockedPrefixSize()
}

// ReservedPrefixSize reports the byte offset where reserved-name entries begin in
// the reference table.
//
//	offset := svc.ReservedPrefixSize()
func (s *Service) ReservedPrefixSize() int {
	return s.ReservedCatalog().PrefixSize()
}

// GetReservedPrefixSize is an alias for ReservedPrefixSize.
//
//	offset := svc.GetReservedPrefixSize()
func (s *Service) GetReservedPrefixSize() int {
	return s.ReservedPrefixSize()
}

// LockedEntries exposes the locked-name catalog entries used by the service.
//
//	entries := svc.LockedEntries()
func (s *Service) LockedEntries() []covenant.LockedEntry {
	return s.LockedCatalog().Entries()
}

// GetLockedEntries is an alias for LockedEntries.
//
//	entries := svc.GetLockedEntries()
func (s *Service) GetLockedEntries() []covenant.LockedEntry {
	return s.LockedEntries()
}

// LockedKeys exposes the locked-name hashes used by the service.
//
//	keys := svc.LockedKeys()
func (s *Service) LockedKeys() []primitives.Hash {
	return s.LockedCatalog().Keys()
}

// GetLockedKeys is an alias for LockedKeys.
//
//	keys := svc.GetLockedKeys()
func (s *Service) GetLockedKeys() []primitives.Hash {
	return s.LockedKeys()
}

// LockedValues exposes the locked-name entries used by the service.
//
//	values := svc.LockedValues()
func (s *Service) LockedValues() []covenant.LockedName {
	return s.LockedCatalog().Values()
}

// GetLockedValues is an alias for LockedValues.
//
//	values := svc.GetLockedValues()
func (s *Service) GetLockedValues() []covenant.LockedName {
	return s.LockedValues()
}
