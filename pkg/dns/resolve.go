// SPDX-License-Identifier: EUPL-1.2

// Package dns provides name resolution for the Lethean Name System.
// It resolves .lthn names to their registered DNS resource records by
// querying the name-chain state tree.
//
// Create a resolver and look up a name:
//
//	svc := dns.NewService()
//	// svc.Resolve("example") — implementation pending
package dns

import (
	core "dappco.re/go/core"
	"dappco.re/go/lns/pkg/covenant"
	"dappco.re/go/lns/pkg/primitives"
)

// Service handles DNS resolution for .lthn names against the LNS chain state.
// The full resolution logic (672 lines from the existing Go implementation)
// will be merged in a subsequent pass.
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
	return s.Resolve(name)
}

// ResolveName is an alias for Resolve.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveName("example.lthn")
func (s *Service) ResolveName(name any) (primitives.Hash, error) {
	return s.Resolve(name)
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
	return s.Verify(name)
}

// VerifyName is an alias for Verify.
//
//	svc := dns.NewService()
//	ok := svc.VerifyName("example.lthn")
func (s *Service) VerifyName(name any) bool {
	return s.Verify(name)
}

// ResolveByName is an alias for ResolveName.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByName("example.lthn")
func (s *Service) ResolveByName(name any) (primitives.Hash, error) {
	return s.ResolveName(name)
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
