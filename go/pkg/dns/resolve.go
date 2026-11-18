// SPDX-License-Identifier: EUPL-1.2

// Package dns documents the DNS-only resolver path for the Lethean Name System.
// The path matters: pkg/dns/resolve.go covers name hashing and verification,
// while lns.go re-exports the broader service surface.
//
// Create a resolver and look up a name:
//
//	service := dns.NewService()
//	hash, err := service.Resolve("example.lthn")
package dns

import (
	core "github.com/dappcore/go"
	"github.com/dappcore/go-lns/go/internal/nameutil"
	"github.com/dappcore/go-lns/go/pkg/covenant"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// ServiceName is the registered name for the DNS service group.
const ServiceName = "dns"

// GetServiceName is an alias for ServiceName.
//
//	name := dns.ServiceName
func GetServiceName() string {
	return ServiceName
}

func resolveCanonicalName(method string, name any) (primitives.Hash, error) {
	switch name.(type) {
	case string, []byte:
		normalized, ok := nameutil.Canonicalize(name)
		if !ok {
			return primitives.Hash{}, core.E(method, "invalid name", nil)
		}

		return covenant.HashBinary([]byte(normalized))
	default:
		return primitives.Hash{}, core.E(method, "invalid name type", nil)
	}
}

func verifyCanonicalName(name any) bool {
	normalized, ok := nameutil.Canonicalize(name)
	if !ok {
		return false
	}

	return covenant.VerifyString(normalized)
}

// Resolve normalises a .lthn name and returns its canonical hash.
//
// The package-level helper mirrors Service.Resolve for callers that do not
// need to keep a service instance around.
//
//	hash, err := dns.Resolve("example.lthn")
func Resolve(name any) (primitives.Hash, error) {
	return resolveCanonicalName("dns.Resolve", name)
}

// GetResolve is an alias for Resolve.
//
//	hash, err := dns.Resolve("example.lthn")
func GetResolve(name any) (primitives.Hash, error) {
	return Resolve(name)
}

// Hash returns the canonical hash for a .lthn name.
//
//	hash, err := dns.Hash("example.lthn")
func Hash(name any) (primitives.Hash, error) {
	return resolveCanonicalName("dns.Hash", name)
}

// GetHash is an alias for Hash.
//
//	hash, err := dns.Hash("example.lthn")
func GetHash(name any) (primitives.Hash, error) {
	return Hash(name)
}

// Verify reports whether a .lthn name is valid after canonicalisation.
//
//	ok := dns.Verify("example.lthn")
func Verify(name any) bool {
	return verifyCanonicalName(name)
}

// GetVerify is an alias for Verify.
//
//	ok := dns.Verify("example.lthn")
func GetVerify(name any) bool {
	return Verify(name)
}

// ResolveString returns the canonical hash for a .lthn name supplied as a string.
//
//	hash, err := dns.ResolveString("example.lthn")
func ResolveString(name string) (primitives.Hash, error) {
	return resolveCanonicalName("dns.ResolveString", name)
}

// GetResolveString is an alias for ResolveString.
//
//	hash, err := dns.ResolveString("example.lthn")
func GetResolveString(name string) (primitives.Hash, error) {
	return ResolveString(name)
}

// ResolveBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	hash, err := dns.ResolveBinary([]byte("example.lthn"))
func ResolveBinary(name []byte) (primitives.Hash, error) {
	return resolveCanonicalName("dns.ResolveBinary", name)
}

// GetResolveBinary is an alias for ResolveBinary.
//
//	hash, err := dns.ResolveBinary([]byte("example.lthn"))
func GetResolveBinary(name []byte) (primitives.Hash, error) {
	return ResolveBinary(name)
}

// HashString returns the canonical hash for a .lthn name supplied as a string.
//
//	hash, err := dns.HashString("example.lthn")
func HashString(name string) (primitives.Hash, error) {
	return resolveCanonicalName("dns.HashString", name)
}

// GetHashString is an alias for HashString.
//
//	hash, err := dns.HashString("example.lthn")
func GetHashString(name string) (primitives.Hash, error) {
	return HashString(name)
}

// HashBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	hash, err := dns.HashBinary([]byte("example.lthn"))
func HashBinary(name []byte) (primitives.Hash, error) {
	return resolveCanonicalName("dns.HashBinary", name)
}

// GetHashBinary is an alias for HashBinary.
//
//	hash, err := dns.HashBinary([]byte("example.lthn"))
func GetHashBinary(name []byte) (primitives.Hash, error) {
	return HashBinary(name)
}

// ResolveName is an alias for Resolve.
//
//	hash, err := dns.ResolveName("example.lthn")
func ResolveName(name any) (primitives.Hash, error) {
	return Resolve(name)
}

// GetResolveName is an alias for ResolveName.
//
//	hash, err := dns.ResolveName("example.lthn")
func GetResolveName(name any) (primitives.Hash, error) {
	return ResolveName(name)
}

// HashName is an alias for Hash.
//
//	hash, err := dns.HashName("example.lthn")
func HashName(name any) (primitives.Hash, error) {
	return Hash(name)
}

// GetHashName is an alias for HashName.
//
//	hash, err := dns.HashName("example.lthn")
func GetHashName(name any) (primitives.Hash, error) {
	return HashName(name)
}

// VerifyName is an alias for Verify.
//
//	ok := dns.VerifyName("example.lthn")
func VerifyName(name any) bool {
	return Verify(name)
}

// GetVerifyName is an alias for VerifyName.
//
//	ok := dns.VerifyName("example.lthn")
func GetVerifyName(name any) bool {
	return VerifyName(name)
}

// ResolveByString is an alias for ResolveString.
//
//	hash, err := dns.ResolveByString("example.lthn")
func ResolveByString(name string) (primitives.Hash, error) {
	return ResolveString(name)
}

// GetResolveByString is an alias for ResolveByString.
//
//	hash, err := dns.ResolveByString("example.lthn")
func GetResolveByString(name string) (primitives.Hash, error) {
	return ResolveByString(name)
}

// ResolveByBinary is an alias for ResolveBinary.
//
//	hash, err := dns.ResolveByBinary([]byte("example.lthn"))
func ResolveByBinary(name []byte) (primitives.Hash, error) {
	return ResolveBinary(name)
}

// GetResolveByBinary is an alias for ResolveByBinary.
//
//	hash, err := dns.ResolveByBinary([]byte("example.lthn"))
func GetResolveByBinary(name []byte) (primitives.Hash, error) {
	return ResolveByBinary(name)
}

// HashByString is an alias for HashString.
//
//	hash, err := dns.HashByString("example.lthn")
func HashByString(name string) (primitives.Hash, error) {
	return HashString(name)
}

// GetHashByString is an alias for HashByString.
//
//	hash, err := dns.HashByString("example.lthn")
func GetHashByString(name string) (primitives.Hash, error) {
	return HashByString(name)
}

// HashByBinary is an alias for HashBinary.
//
//	hash, err := dns.HashByBinary([]byte("example.lthn"))
func HashByBinary(name []byte) (primitives.Hash, error) {
	return HashBinary(name)
}

// GetHashByBinary is an alias for HashByBinary.
//
//	hash, err := dns.HashByBinary([]byte("example.lthn"))
func GetHashByBinary(name []byte) (primitives.Hash, error) {
	return HashByBinary(name)
}

// VerifyString is an alias for Verify.
//
//	ok := dns.VerifyString("example.lthn")
func VerifyString(name string) bool {
	return Verify(name)
}

// GetVerifyString is an alias for VerifyString.
//
//	ok := dns.VerifyString("example.lthn")
func GetVerifyString(name string) bool {
	return VerifyString(name)
}

// VerifyBinary is an alias for Verify.
//
//	ok := dns.VerifyBinary([]byte("example.lthn"))
func VerifyBinary(name []byte) bool {
	return Verify(name)
}

// GetVerifyBinary is an alias for VerifyBinary.
//
//	ok := dns.VerifyBinary([]byte("example.lthn"))
func GetVerifyBinary(name []byte) bool {
	return VerifyBinary(name)
}

// VerifyByString is an alias for VerifyString.
//
//	ok := dns.VerifyByString("example.lthn")
func VerifyByString(name string) bool {
	return VerifyString(name)
}

// GetVerifyByString is an alias for VerifyByString.
//
//	ok := dns.VerifyByString("example.lthn")
func GetVerifyByString(name string) bool {
	return VerifyByString(name)
}

// VerifyByBinary is an alias for VerifyBinary.
//
//	ok := dns.VerifyByBinary([]byte("example.lthn"))
func VerifyByBinary(name []byte) bool {
	return VerifyBinary(name)
}

// GetVerifyByBinary is an alias for VerifyByBinary.
//
//	ok := dns.VerifyByBinary([]byte("example.lthn"))
func GetVerifyByBinary(name []byte) bool {
	return VerifyByBinary(name)
}

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
		if opt != nil {
			opt(s)
		}
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
		if s == nil || c == nil {
			return
		}

		s.core = c
	}
}

// Core returns the Core instance associated with the DNS service.
//
//	core := svc.Core()
func (s *Service) Core() *core.Core {
	if s == nil {
		return nil
	}

	return s.core
}

// GetCore is an alias for Core.
//
//	core := svc.Core()
func (s *Service) GetCore() *core.Core {
	return s.Core()
}

// ServiceName returns the registered name for the DNS service group.
//
//	name := svc.ServiceName()
func (s *Service) ServiceName() string {
	return ServiceName
}

// GetServiceName is an alias for ServiceName.
//
//	name := svc.ServiceName()
func (s *Service) GetServiceName() string {
	return s.ServiceName()
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
	return resolveCanonicalName("dns.Service.Resolve", name)
}

// GetResolve is an alias for Resolve.
//
//	svc := dns.NewService()
//	hash, err := svc.Resolve("example.lthn")
func (s *Service) GetResolve(name any) (primitives.Hash, error) {
	return s.Resolve(name)
}

// Hash returns the canonical hash for a .lthn name.
//
//	svc := dns.NewService()
//	hash, err := svc.Hash("example.lthn")
func (s *Service) Hash(name any) (primitives.Hash, error) {
	return resolveCanonicalName("dns.Service.Hash", name)
}

// GetHash is an alias for Hash.
//
//	svc := dns.NewService()
//	hash, err := svc.Hash("example.lthn")
func (s *Service) GetHash(name any) (primitives.Hash, error) {
	return s.Hash(name)
}

// Verify reports whether a .lthn name is valid after canonicalisation.
//
// This mirrors Resolve's input handling but returns a boolean so callers can
// validate user input without allocating or hashing.
//
//	svc := dns.NewService()
//	ok := svc.Verify("example.lthn")
func (s *Service) Verify(name any) bool {
	return verifyCanonicalName(name)
}

// GetVerify is an alias for Verify.
//
//	svc := dns.NewService()
//	ok := svc.Verify("example.lthn")
func (s *Service) GetVerify(name any) bool {
	return s.Verify(name)
}

// HashString returns the canonical hash for a .lthn name supplied as a string.
//
//	hash, err := svc.HashString("example.lthn")
func (s *Service) HashString(name string) (primitives.Hash, error) {
	return resolveCanonicalName("dns.Service.HashString", name)
}

// GetHashString is an alias for HashString.
//
//	hash, err := svc.HashString("example.lthn")
func (s *Service) GetHashString(name string) (primitives.Hash, error) {
	return s.HashString(name)
}

// HashBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	hash, err := svc.HashBinary([]byte("example.lthn"))
func (s *Service) HashBinary(name []byte) (primitives.Hash, error) {
	return resolveCanonicalName("dns.Service.HashBinary", name)
}

// GetHashBinary is an alias for HashBinary.
//
//	hash, err := svc.HashBinary([]byte("example.lthn"))
func (s *Service) GetHashBinary(name []byte) (primitives.Hash, error) {
	return s.HashBinary(name)
}

// ResolveByString is an alias for ResolveString.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByString("example.lthn")
func (s *Service) ResolveByString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// GetResolveByString is an alias for ResolveByString.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByString("example.lthn")
func (s *Service) GetResolveByString(name string) (primitives.Hash, error) {
	return s.ResolveByString(name)
}

// ResolveByBinary is an alias for ResolveBinary.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByBinary([]byte("example.lthn"))
func (s *Service) ResolveByBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// GetResolveByBinary is an alias for ResolveByBinary.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByBinary([]byte("example.lthn"))
func (s *Service) GetResolveByBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveByBinary(name)
}

// HashByString is an alias for HashString.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByString("example.lthn")
func (s *Service) HashByString(name string) (primitives.Hash, error) {
	return s.HashString(name)
}

// GetHashByString is an alias for HashByString.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByString("example.lthn")
func (s *Service) GetHashByString(name string) (primitives.Hash, error) {
	return s.HashByString(name)
}

// HashByBinary is an alias for HashBinary.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByBinary([]byte("example.lthn"))
func (s *Service) HashByBinary(name []byte) (primitives.Hash, error) {
	return s.HashBinary(name)
}

// GetHashByBinary is an alias for HashByBinary.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByBinary([]byte("example.lthn"))
func (s *Service) GetHashByBinary(name []byte) (primitives.Hash, error) {
	return s.HashByBinary(name)
}

// ResolveString returns the canonical hash for a .lthn name supplied as a string.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveString("example.lthn")
func (s *Service) ResolveString(name string) (primitives.Hash, error) {
	return resolveCanonicalName("dns.Service.ResolveString", name)
}

// GetResolveString is an alias for ResolveString.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveString("example.lthn")
func (s *Service) GetResolveString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// ResolveBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveBinary([]byte("example.lthn"))
func (s *Service) ResolveBinary(name []byte) (primitives.Hash, error) {
	return resolveCanonicalName("dns.Service.ResolveBinary", name)
}

// GetResolveBinary is an alias for ResolveBinary.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveBinary([]byte("example.lthn"))
func (s *Service) GetResolveBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// ResolveName is an alias for Resolve.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveName("example.lthn")
func (s *Service) ResolveName(name any) (primitives.Hash, error) {
	return s.Resolve(name)
}

// GetResolveName is an alias for ResolveName.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveName("example.lthn")
func (s *Service) GetResolveName(name any) (primitives.Hash, error) {
	return s.ResolveName(name)
}

// HashName is an alias for Hash.
//
//	hash, err := svc.HashName("example.lthn")
func (s *Service) HashName(name any) (primitives.Hash, error) {
	return s.Hash(name)
}

// GetHashName is an alias for HashName.
//
//	hash, err := svc.HashName("example.lthn")
func (s *Service) GetHashName(name any) (primitives.Hash, error) {
	return s.HashName(name)
}

// VerifyString reports whether a .lthn name supplied as a string is valid.
//
//	svc := dns.NewService()
//	ok := svc.VerifyString("example.lthn")
func (s *Service) VerifyString(name string) bool {
	return s.Verify(name)
}

// GetVerifyString is an alias for VerifyString.
//
//	svc := dns.NewService()
//	ok := svc.VerifyString("example.lthn")
func (s *Service) GetVerifyString(name string) bool {
	return s.VerifyString(name)
}

// VerifyBinary reports whether a .lthn name supplied as bytes is valid.
//
//	svc := dns.NewService()
//	ok := svc.VerifyBinary([]byte("example.lthn"))
func (s *Service) VerifyBinary(name []byte) bool {
	return s.Verify(name)
}

// GetVerifyBinary is an alias for VerifyBinary.
//
//	svc := dns.NewService()
//	ok := svc.VerifyBinary([]byte("example.lthn"))
func (s *Service) GetVerifyBinary(name []byte) bool {
	return s.VerifyBinary(name)
}

// VerifyName is an alias for Verify.
//
//	svc := dns.NewService()
//	ok := svc.VerifyName("example.lthn")
func (s *Service) VerifyName(name any) bool {
	return s.Verify(name)
}

// GetVerifyName is an alias for VerifyName.
//
//	svc := dns.NewService()
//	ok := svc.VerifyName("example.lthn")
func (s *Service) GetVerifyName(name any) bool {
	return s.VerifyName(name)
}

// ResolveByName is an alias for ResolveName.
//
//	hash, err := dns.ResolveByName("example.lthn")
func ResolveByName(name any) (primitives.Hash, error) {
	return ResolveName(name)
}

// GetResolveByName is an alias for ResolveByName.
//
//	hash, err := dns.ResolveByName("example.lthn")
func GetResolveByName(name any) (primitives.Hash, error) {
	return ResolveByName(name)
}

// HashByName is an alias for HashName.
//
//	hash, err := dns.HashByName("example.lthn")
func HashByName(name any) (primitives.Hash, error) {
	return HashName(name)
}

// GetHashByName is an alias for HashByName.
//
//	hash, err := dns.HashByName("example.lthn")
func GetHashByName(name any) (primitives.Hash, error) {
	return HashByName(name)
}

// VerifyByName is an alias for VerifyName.
//
//	ok := dns.VerifyByName("example.lthn")
func VerifyByName(name any) bool {
	return VerifyName(name)
}

// GetVerifyByName is an alias for VerifyByName.
//
//	ok := dns.VerifyByName("example.lthn")
func GetVerifyByName(name any) bool {
	return VerifyByName(name)
}

// VerifyByString is an alias for VerifyString.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByString("example.lthn")
func (s *Service) VerifyByString(name string) bool {
	return s.VerifyString(name)
}

// GetVerifyByString is an alias for VerifyByString.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByString("example.lthn")
func (s *Service) GetVerifyByString(name string) bool {
	return s.VerifyByString(name)
}

// VerifyByBinary is an alias for VerifyBinary.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByBinary([]byte("example.lthn"))
func (s *Service) VerifyByBinary(name []byte) bool {
	return s.VerifyBinary(name)
}

// GetVerifyByBinary is an alias for VerifyByBinary.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByBinary([]byte("example.lthn"))
func (s *Service) GetVerifyByBinary(name []byte) bool {
	return s.VerifyByBinary(name)
}

// ResolveByName is an alias for ResolveName.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByName("example.lthn")
func (s *Service) ResolveByName(name any) (primitives.Hash, error) {
	return s.ResolveName(name)
}

// GetResolveByName is an alias for ResolveByName.
//
//	svc := dns.NewService()
//	hash, err := svc.ResolveByName("example.lthn")
func (s *Service) GetResolveByName(name any) (primitives.Hash, error) {
	return s.ResolveByName(name)
}

// HashByName is an alias for HashName.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByName("example.lthn")
func (s *Service) HashByName(name any) (primitives.Hash, error) {
	return s.HashName(name)
}

// GetHashByName is an alias for HashByName.
//
//	svc := dns.NewService()
//	hash, err := svc.HashByName("example.lthn")
func (s *Service) GetHashByName(name any) (primitives.Hash, error) {
	return s.HashByName(name)
}

// VerifyByName is an alias for VerifyName.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByName("example.lthn")
func (s *Service) VerifyByName(name any) bool {
	return s.VerifyName(name)
}

// GetVerifyByName is an alias for VerifyByName.
//
//	svc := dns.NewService()
//	ok := svc.VerifyByName("example.lthn")
func (s *Service) GetVerifyByName(name any) bool {
	return s.VerifyByName(name)
}
