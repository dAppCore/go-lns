// SPDX-License-Identifier: EUPL-1.2

// Package dns provides name resolution for the Lethean Name System.
// It resolves .lthn names to their registered DNS resource records by
// querying the name-chain state tree.
//
// Create a resolver and look up a name:
//
//	svc := dns.NewService()
//	hash, err := svc.Resolve("example.lthn")
package dns

import (
	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
)

// Service handles DNS resolution for .lthn names against the LNS chain state.
//
//	svc := dns.NewService()
type Service struct {
	core *core.Core
}

// NewService creates a new DNS resolution service.
// The core instance is optional during scaffolding and can be nil.
//
//	svc := dns.NewService()
//	svc := dns.NewService(dns.WithCore(c))
func NewService(opts ...ServiceOption) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ServiceOption configures a DNS Service during construction.
type ServiceOption func(*Service)

// WithCore sets the Core instance for the DNS service, enabling
// access to the service registry, logging, and configuration.
//
//	svc := dns.NewService(dns.WithCore(c))
func WithCore(c *core.Core) ServiceOption {
	return func(s *Service) {
		s.core = c
	}
}

// Resolve normalises a .lthn name and returns its canonical hash.
//
// The resolver accepts either a bare name ("example") or a fully qualified
// name ("example.lthn"). Input is normalised to lowercase before validation,
// matching DNS case-insensitivity while still enforcing the covenant name
// rules. Byte-slice inputs are accepted for parity with the covenant helper
// APIs and are interpreted as ASCII name data.
//
//	svc := dns.NewService()
//	hash, err := svc.Resolve("example.lthn")
func (s *Service) Resolve(name any) (primitives.Hash, error) {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return primitives.Hash{}, core.E("dns.Service.Resolve", "invalid name type", nil)
	}

	if !covenant.VerifyString(normalized) {
		return primitives.Hash{}, core.E("dns.Service.Resolve", "invalid name", nil)
	}

	return covenant.HashString(normalized)
}

// Hash returns the canonical hash for a .lthn name.
//
// Hash is the direct service-level alias for Resolve.
func (s *Service) Hash(name any) (primitives.Hash, error) {
	return s.Resolve(name)
}

// Verify reports whether a .lthn name is valid after canonicalisation.
//
// This mirrors Resolve's input handling but returns a boolean so callers can
// validate user input without allocating or hashing.
//
//	svc := dns.NewService()
//	ok := svc.Verify("example.lthn")
func (s *Service) Verify(name any) bool {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return false
	}

	return covenant.VerifyString(normalized)
}

// HashString returns the canonical hash for a .lthn name supplied as a string.
func (s *Service) HashString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// HashBinary returns the canonical hash for a .lthn name supplied as bytes.
func (s *Service) HashBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// ResolveByString is an alias for ResolveString.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByString("example.lthn")
func (s *Service) ResolveByString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// ResolveByBinary is an alias for ResolveBinary.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByBinary([]byte("example.lthn"))
func (s *Service) ResolveByBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// HashByString is an alias for HashString.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByString("example.lthn")
func (s *Service) HashByString(name string) (primitives.Hash, error) {
	return s.HashString(name)
}

// HashByBinary is an alias for HashBinary.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByBinary([]byte("example.lthn"))
func (s *Service) HashByBinary(name []byte) (primitives.Hash, error) {
	return s.HashBinary(name)
}

// ResolveString returns the canonical hash for a .lthn name supplied as a string.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveString("example.lthn")
func (s *Service) ResolveString(name string) (primitives.Hash, error) {
	return s.Resolve(name)
}

// ResolveBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveBinary([]byte("example.lthn"))
func (s *Service) ResolveBinary(name []byte) (primitives.Hash, error) {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return primitives.Hash{}, core.E("dns.Service.ResolveBinary", "invalid name type", nil)
	}

	if !covenant.VerifyBinary([]byte(normalized)) {
		return primitives.Hash{}, core.E("dns.Service.ResolveBinary", "invalid name", nil)
	}

	return covenant.HashBinary([]byte(normalized))
}

// ResolveName is an alias for Resolve.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveName("example.lthn")
func (s *Service) ResolveName(name any) (primitives.Hash, error) {
	return s.Resolve(name)
}

// HashName is an alias for Hash.
func (s *Service) HashName(name any) (primitives.Hash, error) {
	return s.Hash(name)
}

// VerifyString reports whether a .lthn name supplied as a string is valid.
//
//	svc := dns.NewService()
//	ok := svc.VerifyString("example.lthn")
func (s *Service) VerifyString(name string) bool {
	return s.Verify(name)
}

// VerifyBinary reports whether a .lthn name supplied as bytes is valid.
//
//	svc := dns.NewService()
//	ok := svc.VerifyBinary([]byte("example.lthn"))
func (s *Service) VerifyBinary(name []byte) bool {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return false
	}

	return covenant.VerifyBinary([]byte(normalized))
}

// VerifyName is an alias for Verify.
//
//	svc := dns.NewService()
//	ok := svc.VerifyName("example.lthn")
func (s *Service) VerifyName(name any) bool {
	return s.Verify(name)
}

// VerifyByString is an alias for VerifyString.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByString("example.lthn")
func (s *Service) VerifyByString(name string) bool {
	return s.VerifyString(name)
}

// VerifyByBinary is an alias for VerifyBinary.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByBinary([]byte("example.lthn"))
func (s *Service) VerifyByBinary(name []byte) bool {
	return s.VerifyBinary(name)
}

// ResolveByName is an alias for ResolveName.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByName("example.lthn")
func (s *Service) ResolveByName(name any) (primitives.Hash, error) {
	return s.ResolveName(name)
}

// HashByName is an alias for HashName.
func (s *Service) HashByName(name any) (primitives.Hash, error) {
	return s.HashName(name)
}

// VerifyByName is an alias for VerifyName.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByName("example.lthn")
func (s *Service) VerifyByName(name any) bool {
	return s.VerifyName(name)
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

	value = core.Lower(value)

	switch {
	case core.HasSuffix(value, ".lthn."):
		value = core.TrimSuffix(value, ".lthn.")
	case core.HasSuffix(value, ".lthn"):
		value = core.TrimSuffix(value, ".lthn")
	case core.HasSuffix(value, "."):
		value = core.TrimSuffix(value, ".")
	}

	return value, true
}
