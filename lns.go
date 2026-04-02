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
	"context"
	"crypto/sha3"

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

var (
	_ core.Startable = (*Service)(nil)
	_ core.Stoppable = (*Service)(nil)
	_ interface {
		HandleIPCEvents(*core.Core, core.Message) core.Result
	} = (*Service)(nil)
)

// serviceOptions holds configuration for the LNS service.
type serviceOptions struct{}

// ServiceOption configures an LNS service during construction.
type ServiceOption func(*Service)

// WithCore sets the Core instance for an LNS service.
//
//	svc := lns.NewServiceWithOptions(lns.WithCore(c))
func WithCore(c *core.Core) ServiceOption {
	return func(s *Service) {
		if s == nil || c == nil {
			return
		}

		s.ServiceRuntime = core.NewServiceRuntime(c, serviceOptions{})
	}
}

// NewServiceWithOptions constructs an LNS Service from the provided options.
//
//	svc := lns.NewServiceWithOptions(lns.WithCore(c))
func NewServiceWithOptions(opts ...ServiceOption) *Service {
	svc := &Service{}

	for _, opt := range opts {
		if opt != nil {
			opt(svc)
		}
	}

	return svc
}

// NewService constructs an LNS Service bound to the provided Core instance.
//
//	c := core.New()
//	svc := lns.NewService(c)
func NewService(c *core.Core) *Service {
	return NewServiceWithOptions(WithCore(c))
}

// Register constructs an LNS Service bound to the provided Core instance.
// Pass to core.WithService during Core construction:
//
//	c := core.New(core.WithService(lns.Register))
func Register(c *core.Core) core.Result {
	if c == nil {
		return core.Result{Value: core.E("lns.Register", "core is required", nil), OK: false}
	}

	return core.Result{Value: NewService(c), OK: true}
}

// GetServiceName returns the registered name for the LNS service group.
//
//	name := lns.GetServiceName()
func GetServiceName() string {
	return ServiceName
}

// OnStartup satisfies core.Startable so LNS participates in Core lifecycle
// discovery without requiring any explicit setup.
func (s *Service) OnStartup(context.Context) core.Result {
	return core.Result{OK: true}
}

// OnShutdown satisfies core.Stoppable so LNS participates in Core lifecycle
// discovery without requiring any explicit teardown.
func (s *Service) OnShutdown(context.Context) core.Result {
	return core.Result{OK: true}
}

// HandleIPCEvents satisfies Core's IPC handler discovery.
//
//	core.New(core.WithService(lns.Register))
//	// LNS registers a no-op IPC handler so it can join the bus cleanly.
func (s *Service) HandleIPCEvents(*core.Core, core.Message) core.Result {
	return core.Result{OK: true}
}

func canonicalizeName(name any) (string, bool) {
	switch v := name.(type) {
	case string:
		return nameutil.Canonicalize(v)
	case []byte:
		return nameutil.Canonicalize(v)
	default:
		return "", false
	}
}

func resolveCanonicalName(method string, name any) (primitives.Hash, error) {
	switch v := name.(type) {
	case string, []byte:
		normalized, ok := canonicalizeName(v)
		if !ok {
			return primitives.Hash{}, core.E(method, "invalid name", nil)
		}

		return covenant.HashString(normalized)
	default:
		return primitives.Hash{}, core.E(method, "invalid name type", nil)
	}
}

func verifyCanonicalName(name any) bool {
	normalized, ok := canonicalizeName(name)
	if !ok {
		return false
	}

	return covenant.VerifyString(normalized)
}

func lookupCatalogName[T any](
	name any,
	byLabel func(string) (T, bool),
	byHash func(primitives.Hash) (T, bool),
) (T, bool) {
	label, ok := nameutil.CatalogLabel(name)
	if !ok {
		var zero T
		return zero, false
	}

	if item, ok := byLabel(label); ok {
		return item, true
	}

	normalized, ok := nameutil.Canonicalize(name)
	if !ok {
		var zero T
		return zero, false
	}

	hash := primitives.Hash(sha3.Sum256([]byte(normalized)))
	return byHash(hash)
}

// Resolve returns the canonical hash for a validated .lthn name.
//
// This is the package-level convenience wrapper around Service.Resolve.
//
//	hash, err := lns.Resolve("example.lthn")
func Resolve(name any) (primitives.Hash, error) {
	return resolveCanonicalName("lns.Resolve", name)
}

// GetResolve is an alias for Resolve.
//
//	hash, err := lns.GetResolve("example.lthn")
func GetResolve(name any) (primitives.Hash, error) {
	return Resolve(name)
}

// Hash returns the canonical hash for a validated .lthn name.
//
// This is the package-level convenience wrapper around Service.Hash.
//
//	hash, err := lns.Hash("example.lthn")
func Hash(name any) (primitives.Hash, error) {
	return Resolve(name)
}

// GetHash is an alias for Hash.
//
//	hash, err := lns.GetHash("example.lthn")
func GetHash(name any) (primitives.Hash, error) {
	return Hash(name)
}

// Verify reports whether a .lthn name is valid after canonicalisation.
//
// This is the package-level convenience wrapper around Service.Verify.
//
//	ok := lns.Verify("example.lthn")
func Verify(name any) bool {
	return verifyCanonicalName(name)
}

// GetVerify is an alias for Verify.
//
//	ok := lns.GetVerify("example.lthn")
func GetVerify(name any) bool {
	return Verify(name)
}

// HashString returns the canonical hash for a .lthn name supplied as a string.
//
//	hash, err := lns.HashString("example.lthn")
func HashString(name string) (primitives.Hash, error) {
	return ResolveString(name)
}

// GetHashString is an alias for HashString.
//
//	hash, err := lns.GetHashString("example.lthn")
func GetHashString(name string) (primitives.Hash, error) {
	return HashString(name)
}

// HashBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	hash, err := lns.HashBinary([]byte("example.lthn"))
func HashBinary(name []byte) (primitives.Hash, error) {
	return ResolveBinary(name)
}

// GetHashBinary is an alias for HashBinary.
//
//	hash, err := lns.GetHashBinary([]byte("example.lthn"))
func GetHashBinary(name []byte) (primitives.Hash, error) {
	return HashBinary(name)
}

// ResolveByString is an alias for ResolveString.
//
//	hash, err := lns.ResolveByString("example.lthn")
func ResolveByString(name string) (primitives.Hash, error) {
	return ResolveString(name)
}

// GetResolveByString is an alias for ResolveByString.
//
//	hash, err := lns.GetResolveByString("example.lthn")
func GetResolveByString(name string) (primitives.Hash, error) {
	return ResolveByString(name)
}

// ResolveByBinary is an alias for ResolveBinary.
//
//	hash, err := lns.ResolveByBinary([]byte("example.lthn"))
func ResolveByBinary(name []byte) (primitives.Hash, error) {
	return ResolveBinary(name)
}

// GetResolveByBinary is an alias for ResolveByBinary.
//
//	hash, err := lns.GetResolveByBinary([]byte("example.lthn"))
func GetResolveByBinary(name []byte) (primitives.Hash, error) {
	return ResolveByBinary(name)
}

// HashByString is an alias for HashString.
//
//	hash, err := lns.HashByString("example.lthn")
func HashByString(name string) (primitives.Hash, error) {
	return HashString(name)
}

// GetHashByString is an alias for HashByString.
//
//	hash, err := lns.GetHashByString("example.lthn")
func GetHashByString(name string) (primitives.Hash, error) {
	return HashByString(name)
}

// HashByBinary is an alias for HashBinary.
//
//	hash, err := lns.HashByBinary([]byte("example.lthn"))
func HashByBinary(name []byte) (primitives.Hash, error) {
	return HashBinary(name)
}

// GetHashByBinary is an alias for HashByBinary.
//
//	hash, err := lns.GetHashByBinary([]byte("example.lthn"))
func GetHashByBinary(name []byte) (primitives.Hash, error) {
	return HashByBinary(name)
}

// ResolveString returns the canonical hash for a .lthn name supplied as a string.
//
//	hash, err := lns.ResolveString("example.lthn")
func ResolveString(name string) (primitives.Hash, error) {
	return resolveCanonicalName("lns.ResolveString", name)
}

// GetResolveString is an alias for ResolveString.
//
//	hash, err := lns.GetResolveString("example.lthn")
func GetResolveString(name string) (primitives.Hash, error) {
	return ResolveString(name)
}

// ResolveBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	hash, err := lns.ResolveBinary([]byte("example.lthn"))
func ResolveBinary(name []byte) (primitives.Hash, error) {
	return resolveCanonicalName("lns.ResolveBinary", name)
}

// GetResolveBinary is an alias for ResolveBinary.
//
//	hash, err := lns.GetResolveBinary([]byte("example.lthn"))
func GetResolveBinary(name []byte) (primitives.Hash, error) {
	return ResolveBinary(name)
}

// ResolveName is an alias for Resolve.
//
//	hash, err := lns.ResolveName("example.lthn")
func ResolveName(name any) (primitives.Hash, error) {
	return Resolve(name)
}

// GetResolveName is an alias for ResolveName.
//
//	hash, err := lns.GetResolveName("example.lthn")
func GetResolveName(name any) (primitives.Hash, error) {
	return ResolveName(name)
}

// HashName is an alias for Hash.
//
//	hash, err := lns.HashName("example.lthn")
func HashName(name any) (primitives.Hash, error) {
	return Hash(name)
}

// GetHashName is an alias for HashName.
//
//	hash, err := lns.GetHashName("example.lthn")
func GetHashName(name any) (primitives.Hash, error) {
	return HashName(name)
}

// VerifyName is an alias for Verify.
//
//	ok := lns.VerifyName("example.lthn")
func VerifyName(name any) bool {
	return Verify(name)
}

// GetVerifyName is an alias for VerifyName.
//
//	ok := lns.GetVerifyName("example.lthn")
func GetVerifyName(name any) bool {
	return VerifyName(name)
}

// ResolveByName is an alias for ResolveName.
//
//	hash, err := lns.ResolveByName("example.lthn")
func ResolveByName(name any) (primitives.Hash, error) {
	return ResolveName(name)
}

// GetResolveByName is an alias for ResolveByName.
//
//	hash, err := lns.GetResolveByName("example.lthn")
func GetResolveByName(name any) (primitives.Hash, error) {
	return ResolveByName(name)
}

// HashByName is an alias for HashName.
//
//	hash, err := lns.HashByName("example.lthn")
func HashByName(name any) (primitives.Hash, error) {
	return HashName(name)
}

// GetHashByName is an alias for HashByName.
//
//	hash, err := lns.GetHashByName("example.lthn")
func GetHashByName(name any) (primitives.Hash, error) {
	return HashByName(name)
}

// VerifyByName is an alias for VerifyName.
//
//	ok := lns.VerifyByName("example.lthn")
func VerifyByName(name any) bool {
	return VerifyName(name)
}

// GetVerifyByName is an alias for VerifyByName.
//
//	ok := lns.GetVerifyByName("example.lthn")
func GetVerifyByName(name any) bool {
	return VerifyByName(name)
}

// VerifyString reports whether a .lthn name supplied as a string is valid.
//
//	ok := lns.VerifyString("example.lthn")
func VerifyString(name string) bool {
	return verifyCanonicalName(name)
}

// GetVerifyString is an alias for VerifyString.
//
//	ok := lns.GetVerifyString("example.lthn")
func GetVerifyString(name string) bool {
	return VerifyString(name)
}

// VerifyBinary reports whether a .lthn name supplied as bytes is valid.
//
//	ok := lns.VerifyBinary([]byte("example.lthn"))
func VerifyBinary(name []byte) bool {
	return verifyCanonicalName(name)
}

// GetVerifyBinary is an alias for VerifyBinary.
//
//	ok := lns.GetVerifyBinary([]byte("example.lthn"))
func GetVerifyBinary(name []byte) bool {
	return VerifyBinary(name)
}

// VerifyByString is an alias for VerifyString.
//
//	ok := lns.VerifyByString("example.lthn")
func VerifyByString(name string) bool {
	return VerifyString(name)
}

// GetVerifyByString is an alias for VerifyByString.
//
//	ok := lns.GetVerifyByString("example.lthn")
func GetVerifyByString(name string) bool {
	return VerifyByString(name)
}

// VerifyByBinary is an alias for VerifyBinary.
//
//	ok := lns.VerifyByBinary([]byte("example.lthn"))
func VerifyByBinary(name []byte) bool {
	return VerifyBinary(name)
}

// GetVerifyByBinary is an alias for VerifyByBinary.
//
//	ok := lns.GetVerifyByBinary([]byte("example.lthn"))
func GetVerifyByBinary(name []byte) bool {
	return VerifyByBinary(name)
}

// Blind returns the sealed-bid commitment for an amount and nonce.
//
//	bid, err := lns.Blind(1000, nonce)
func Blind(value uint64, nonce primitives.Hash) (primitives.Hash, error) {
	return covenant.Blind(value, nonce)
}

// GetBlind is an alias for Blind.
//
//	bid, err := lns.GetBlind(1000, nonce)
func GetBlind(value uint64, nonce primitives.Hash) (primitives.Hash, error) {
	return Blind(value, nonce)
}

// TypeName returns the canonical string for a covenant type.
//
//	name := lns.TypeName(covenant.TypeBid)
func TypeName(ct covenant.CovenantType) string {
	return covenant.TypeName(ct)
}

// Types mirrors the covenant type lookup table from pkg/covenant.
//
//	types := lns.Types
var Types = covenant.Types

// TypesByVal mirrors the covenant reverse lookup table from pkg/covenant.
//
//	typesByVal := lns.TypesByVal
var TypesByVal = covenant.TypesByVal

// Blacklist mirrors the covenant blacklist used by the verifier.
//
//	blacklist := lns.Blacklist
var Blacklist = covenant.Blacklist

// VerificationFlags mirrors the covenant verification flag lookup table.
//
//	flags := lns.VerificationFlags
var VerificationFlags = covenant.VerificationFlags

// VerificationFlagsByVal mirrors the covenant verification flag reverse table.
//
//	flagsByVal := lns.VerificationFlagsByVal
var VerificationFlagsByVal = covenant.VerificationFlagsByVal

// GetBlacklist is an alias for Blacklist.
//
//	blacklist := lns.GetBlacklist()
func GetBlacklist() map[string]struct{} {
	return Blacklist
}

// GetVerificationFlags is an alias for VerificationFlags.
//
//	flags := lns.GetVerificationFlags()
func GetVerificationFlags() map[string]NameFlags {
	return VerificationFlags
}

// GetVerificationFlagsByVal is an alias for VerificationFlagsByVal.
//
//	flags := lns.GetVerificationFlagsByVal()
func GetVerificationFlagsByVal() map[NameFlags]string {
	return VerificationFlagsByVal
}

// CovenantType mirrors the covenant type enum from pkg/covenant.
type CovenantType = covenant.CovenantType

// Covenant type constants mirror pkg/covenant so callers can stay in the
// top-level lns package for common rule metadata.
const (
	TypeNone     = covenant.TypeNone
	TypeClaim    = covenant.TypeClaim
	TypeOpen     = covenant.TypeOpen
	TypeBid      = covenant.TypeBid
	TypeReveal   = covenant.TypeReveal
	TypeRedeem   = covenant.TypeRedeem
	TypeRegister = covenant.TypeRegister
	TypeUpdate   = covenant.TypeUpdate
	TypeRenew    = covenant.TypeRenew
	TypeTransfer = covenant.TypeTransfer
	TypeFinalize = covenant.TypeFinalize
	TypeRevoke   = covenant.TypeRevoke
)

// NameFlags mirrors the covenant verification flag type.
type NameFlags = covenant.NameFlags

// NameRules mirrors the covenant rollout and lockup rule inputs.
type NameRules = covenant.NameRules

// NameState mirrors the name-state record type from pkg/primitives.
type NameState = primitives.NameState

// NameDelta mirrors the sparse name-state delta type from pkg/primitives.
type NameDelta = primitives.NameDelta

// NameStateRules mirrors the phase-timing rule set from pkg/primitives.
type NameStateRules = primitives.NameStateRules

// NameStateStats mirrors the derived timing summary type from pkg/primitives.
type NameStateStats = primitives.NameStateStats

// NameStateStatus mirrors the name phase enum from pkg/primitives.
type NameStateStatus = primitives.NameStateStatus

// ReservedName mirrors the reserved-name catalog entry type from pkg/covenant.
type ReservedName = covenant.ReservedName

// LockedName mirrors the locked-name catalog entry type from pkg/covenant.
type LockedName = covenant.LockedName

// ReservedEntry mirrors the reserved-name catalog entry pair from pkg/covenant.
type ReservedEntry = covenant.ReservedEntry

// LockedEntry mirrors the locked-name catalog entry pair from pkg/covenant.
type LockedEntry = covenant.LockedEntry

// Name state constants mirror pkg/primitives so callers can keep using the
// top-level lns package for common domain objects.
const (
	NameStateOpening = primitives.NameStateOpening
	NameStateLocked  = primitives.NameStateLocked
	NameStateBidding = primitives.NameStateBidding
	NameStateReveal  = primitives.NameStateReveal
	NameStateClosed  = primitives.NameStateClosed
	NameStateRevoked = primitives.NameStateRevoked
)

// CoinView mirrors the minimal UTXO lookup interface required by covenant
// verification helpers.
type CoinView = covenant.CoinView

// Network mirrors the covenant verification network inputs.
type Network = covenant.Network

// Verification and size constants mirror pkg/covenant so callers can stay in
// the top-level lns package for common rule metadata.
const (
	MaxNameSize                  = covenant.MaxNameSize
	MaxResourceSize              = covenant.MaxResourceSize
	VerifyCovenantsNone          = covenant.VerifyCovenantsNone
	VerifyCovenantsHardened      = covenant.VerifyCovenantsHardened
	VerifyCovenantsLockup        = covenant.VerifyCovenantsLockup
	MandatoryVerifyCovenantFlags = covenant.MandatoryVerifyCovenantFlags
	MaxCovenantSize              = covenant.MaxCovenantSize
	CovenantMaxSize              = covenant.CovenantMaxSize
)

// GetMaxNameSize returns the maximum supported name size.
func GetMaxNameSize() int {
	return MaxNameSize
}

// GetMaxResourceSize returns the maximum supported DNS resource size.
func GetMaxResourceSize() int {
	return MaxResourceSize
}

// GetVerifyCovenantsNone returns the disabled covenant verification mask.
func GetVerifyCovenantsNone() NameFlags {
	return VerifyCovenantsNone
}

// GetVerifyCovenantsHardened returns the hardened covenant verification mask.
func GetVerifyCovenantsHardened() NameFlags {
	return VerifyCovenantsHardened
}

// GetVerifyCovenantsLockup returns the lockup covenant verification mask.
func GetVerifyCovenantsLockup() NameFlags {
	return VerifyCovenantsLockup
}

// GetMandatoryVerifyCovenantFlags returns the default covenant verification mask.
func GetMandatoryVerifyCovenantFlags() NameFlags {
	return MandatoryVerifyCovenantFlags
}

// GetMaxCovenantSize returns the maximum covenant size.
func GetMaxCovenantSize() int {
	return MaxCovenantSize
}

// GetCovenantMaxSize returns the alias for the maximum covenant size.
func GetCovenantMaxSize() int {
	return CovenantMaxSize
}

// GetTypes is an alias for Types.
//
//	types := lns.GetTypes()
func GetTypes() map[string]covenant.CovenantType {
	return Types
}

// GetTypesByVal is an alias for TypesByVal.
//
//	typesByVal := lns.GetTypesByVal()
func GetTypesByVal() map[covenant.CovenantType]string {
	return TypesByVal
}

// GetTypeName is an alias for TypeName.
//
//	name := lns.GetTypeName(covenant.TypeBid)
func GetTypeName(ct covenant.CovenantType) string {
	return TypeName(ct)
}

// IsName reports whether a covenant type is a name operation.
//
//	ok := lns.IsName(covenant.TypeOpen)
func IsName(ct CovenantType) bool {
	return ct.IsName()
}

// GetIsName is an alias for IsName.
//
//	ok := lns.GetIsName(covenant.TypeOpen)
func GetIsName(ct CovenantType) bool {
	return IsName(ct)
}

// IsKnown reports whether a covenant type is recognized by the rules table.
//
//	ok := lns.IsKnown(covenant.TypeBid)
func IsKnown(ct CovenantType) bool {
	return ct.IsKnown()
}

// GetIsKnown is an alias for IsKnown.
//
//	ok := lns.GetIsKnown(covenant.TypeBid)
func GetIsKnown(ct CovenantType) bool {
	return IsKnown(ct)
}

// IsLinked reports whether a covenant type is a linked covenant.
//
//	ok := lns.IsLinked(covenant.TypeReveal)
func IsLinked(ct CovenantType) bool {
	return ct.IsLinked()
}

// GetIsLinked is an alias for IsLinked.
//
//	ok := lns.GetIsLinked(covenant.TypeReveal)
func GetIsLinked(ct CovenantType) bool {
	return IsLinked(ct)
}

func lookupReserved(name any) (covenant.ReservedName, bool) {
	return lookupCatalogName(
		name,
		ReservedCatalog().GetByName,
		ReservedCatalog().Get,
	)
}

func lookupLocked(name any) (covenant.LockedName, bool) {
	return lookupCatalogName(
		name,
		LockedCatalog().GetByName,
		LockedCatalog().Get,
	)
}

// MaxNameSize returns the maximum supported name size.
func (s *Service) MaxNameSize() int {
	return MaxNameSize
}

// GetMaxNameSize is an alias for MaxNameSize.
func (s *Service) GetMaxNameSize() int {
	return s.MaxNameSize()
}

// MaxResourceSize returns the maximum supported DNS resource size.
func (s *Service) MaxResourceSize() int {
	return MaxResourceSize
}

// GetMaxResourceSize is an alias for MaxResourceSize.
func (s *Service) GetMaxResourceSize() int {
	return s.MaxResourceSize()
}

// VerifyCovenantsNone returns the disabled covenant verification mask.
func (s *Service) VerifyCovenantsNone() NameFlags {
	return VerifyCovenantsNone
}

// GetVerifyCovenantsNone is an alias for VerifyCovenantsNone.
func (s *Service) GetVerifyCovenantsNone() NameFlags {
	return s.VerifyCovenantsNone()
}

// VerifyCovenantsHardened returns the hardened covenant verification mask.
func (s *Service) VerifyCovenantsHardened() NameFlags {
	return VerifyCovenantsHardened
}

// GetVerifyCovenantsHardened is an alias for VerifyCovenantsHardened.
func (s *Service) GetVerifyCovenantsHardened() NameFlags {
	return s.VerifyCovenantsHardened()
}

// VerifyCovenantsLockup returns the lockup covenant verification mask.
func (s *Service) VerifyCovenantsLockup() NameFlags {
	return VerifyCovenantsLockup
}

// GetVerifyCovenantsLockup is an alias for VerifyCovenantsLockup.
func (s *Service) GetVerifyCovenantsLockup() NameFlags {
	return s.VerifyCovenantsLockup()
}

// MandatoryVerifyCovenantFlags returns the default covenant verification mask.
func (s *Service) MandatoryVerifyCovenantFlags() NameFlags {
	return MandatoryVerifyCovenantFlags
}

// GetMandatoryVerifyCovenantFlags is an alias for MandatoryVerifyCovenantFlags.
func (s *Service) GetMandatoryVerifyCovenantFlags() NameFlags {
	return s.MandatoryVerifyCovenantFlags()
}

// MaxCovenantSize returns the maximum covenant size.
func (s *Service) MaxCovenantSize() int {
	return MaxCovenantSize
}

// GetMaxCovenantSize is an alias for MaxCovenantSize.
func (s *Service) GetMaxCovenantSize() int {
	return s.MaxCovenantSize()
}

// CovenantMaxSize returns the alias for the maximum covenant size.
func (s *Service) CovenantMaxSize() int {
	return CovenantMaxSize
}

// GetCovenantMaxSize is an alias for CovenantMaxSize.
func (s *Service) GetCovenantMaxSize() int {
	return s.CovenantMaxSize()
}

// GetRollout returns the rollout start height and rollout week for a name hash.
func GetRollout(nameHash primitives.Hash, rules NameRules) (uint32, uint32) {
	return covenant.GetRollout(nameHash, rules)
}

// HasRollout reports whether a name hash has reached its rollout height.
func HasRollout(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return covenant.HasRollout(nameHash, height, rules)
}

// IsReserved reports whether a name hash is reserved at the given height.
func IsReserved(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return covenant.IsReserved(nameHash, height, rules)
}

// IsLockedUp reports whether a name hash is locked at the given height.
func IsLockedUp(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return covenant.IsLockedUp(nameHash, height, rules)
}

// GrindName searches for an available name that satisfies the current rollout
// and reservation rules.
func GrindName(size int, height uint32, rules NameRules) (string, error) {
	return covenant.GrindName(size, height, rules)
}

// GetGrindName is an alias for GrindName.
func GetGrindName(size int, height uint32, rules NameRules) (string, error) {
	return GrindName(size, height, rules)
}

// CountOpens counts OPEN covenants in a transaction.
func CountOpens(tx primitives.Transaction) int {
	return covenant.CountOpens(tx)
}

// CountUpdates counts update-style covenants in a transaction.
func CountUpdates(tx primitives.Transaction) int {
	return covenant.CountUpdates(tx)
}

// CountRenewals counts REGISTER, RENEW, and FINALIZE covenants in a transaction.
func CountRenewals(tx primitives.Transaction) int {
	return covenant.CountRenewals(tx)
}

// HasNames reports whether any name covenant in the transaction matches the set.
func HasNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) bool {
	return covenant.HasNames(tx, set)
}

// AddNames adds the transaction's name covenants to the set.
func AddNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) {
	covenant.AddNames(tx, set)
}

// RemoveNames removes the transaction's name covenants from the set.
func RemoveNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) {
	covenant.RemoveNames(tx, set)
}

// HasSaneCovenants performs structural covenant validation on a transaction.
func HasSaneCovenants(tx primitives.Transaction) bool {
	return covenant.HasSaneCovenants(tx)
}

// VerifyCovenants performs contextual covenant validation on a transaction.
func VerifyCovenants(tx primitives.Transaction, view CoinView, height uint32, network Network) int {
	return covenant.VerifyCovenants(tx, view, height, network)
}

// GetVerifyCovenants is an alias for VerifyCovenants.
func GetVerifyCovenants(tx primitives.Transaction, view CoinView, height uint32, network Network) int {
	return VerifyCovenants(tx, view, height, network)
}

// ServiceName returns the registered name for the LNS service group.
//
//	name := svc.ServiceName()
func (s *Service) ServiceName() string {
	return ServiceName
}

// GetServiceName is an alias for ServiceName.
//
//	name := svc.GetServiceName()
func (s *Service) GetServiceName() string {
	return s.ServiceName()
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
	return resolveCanonicalName("lns.Service.Resolve", name)
}

// GetResolve is an alias for Resolve.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.GetResolve("example.lthn")
func (s *Service) GetResolve(name any) (primitives.Hash, error) {
	return s.Resolve(name)
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

// GetHash is an alias for Hash.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.GetHash("example.lthn")
func (s *Service) GetHash(name any) (primitives.Hash, error) {
	return s.Hash(name)
}

// Verify reports whether a .lthn name is valid after canonicalisation.
//
// This keeps the service-level API aligned with the DNS resolver helpers.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.Verify("example.lthn")
func (s *Service) Verify(name any) bool {
	return verifyCanonicalName(name)
}

// GetVerify is an alias for Verify.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.GetVerify("example.lthn")
func (s *Service) GetVerify(name any) bool {
	return s.Verify(name)
}

// HashString returns the canonical hash for a .lthn name supplied as a string.
//
//	ok := svc.HashString("example.lthn")
func (s *Service) HashString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// GetHashString is an alias for HashString.
//
//	ok := svc.GetHashString("example.lthn")
func (s *Service) GetHashString(name string) (primitives.Hash, error) {
	return s.HashString(name)
}

// HashBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	ok := svc.HashBinary([]byte("example.lthn"))
func (s *Service) HashBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// GetHashBinary is an alias for HashBinary.
//
//	ok := svc.GetHashBinary([]byte("example.lthn"))
func (s *Service) GetHashBinary(name []byte) (primitives.Hash, error) {
	return s.HashBinary(name)
}

// ResolveByString is an alias for ResolveString.
//
//	ok := svc.ResolveByString("example.lthn")
func (s *Service) ResolveByString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// GetResolveByString is an alias for ResolveByString.
//
//	ok := svc.GetResolveByString("example.lthn")
func (s *Service) GetResolveByString(name string) (primitives.Hash, error) {
	return s.ResolveByString(name)
}

// ResolveByBinary is an alias for ResolveBinary.
//
//	ok := svc.ResolveByBinary([]byte("example.lthn"))
func (s *Service) ResolveByBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// GetResolveByBinary is an alias for ResolveByBinary.
//
//	ok := svc.GetResolveByBinary([]byte("example.lthn"))
func (s *Service) GetResolveByBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveByBinary(name)
}

// HashByString is an alias for HashString.
//
//	ok := svc.HashByString("example.lthn")
func (s *Service) HashByString(name string) (primitives.Hash, error) {
	return s.HashString(name)
}

// GetHashByString is an alias for HashByString.
//
//	ok := svc.GetHashByString("example.lthn")
func (s *Service) GetHashByString(name string) (primitives.Hash, error) {
	return s.HashByString(name)
}

// HashByBinary is an alias for HashBinary.
//
//	ok := svc.HashByBinary([]byte("example.lthn"))
func (s *Service) HashByBinary(name []byte) (primitives.Hash, error) {
	return s.HashBinary(name)
}

// GetHashByBinary is an alias for HashByBinary.
//
//	ok := svc.GetHashByBinary([]byte("example.lthn"))
func (s *Service) GetHashByBinary(name []byte) (primitives.Hash, error) {
	return s.HashByBinary(name)
}

// ResolveString returns the canonical hash for a .lthn name supplied as a string.
//
//	hash, err := svc.ResolveString("example.lthn")
func (s *Service) ResolveString(name string) (primitives.Hash, error) {
	return resolveCanonicalName("lns.Service.ResolveString", name)
}

// GetResolveString is an alias for ResolveString.
//
//	hash, err := svc.GetResolveString("example.lthn")
func (s *Service) GetResolveString(name string) (primitives.Hash, error) {
	return s.ResolveString(name)
}

// ResolveBinary returns the canonical hash for a .lthn name supplied as bytes.
//
//	hash, err := svc.ResolveBinary([]byte("example.lthn"))
func (s *Service) ResolveBinary(name []byte) (primitives.Hash, error) {
	return resolveCanonicalName("lns.Service.ResolveBinary", name)
}

// GetResolveBinary is an alias for ResolveBinary.
//
//	hash, err := svc.GetResolveBinary([]byte("example.lthn"))
func (s *Service) GetResolveBinary(name []byte) (primitives.Hash, error) {
	return s.ResolveBinary(name)
}

// ResolveName is an alias for Resolve.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.ResolveName("example.lthn")
func (s *Service) ResolveName(name any) (primitives.Hash, error) {
	return s.Resolve(name)
}

// GetResolveName is an alias for ResolveName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.GetResolveName("example.lthn")
func (s *Service) GetResolveName(name any) (primitives.Hash, error) {
	return s.ResolveName(name)
}

// HashName is an alias for Hash.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.HashName("example.lthn")
func (s *Service) HashName(name any) (primitives.Hash, error) {
	return s.Hash(name)
}

// GetHashName is an alias for HashName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.GetHashName("example.lthn")
func (s *Service) GetHashName(name any) (primitives.Hash, error) {
	return s.HashName(name)
}

// VerifyName is an alias for Verify.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.VerifyName("example.lthn")
func (s *Service) VerifyName(name any) bool {
	return s.Verify(name)
}

// GetVerifyName is an alias for VerifyName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.GetVerifyName("example.lthn")
func (s *Service) GetVerifyName(name any) bool {
	return s.VerifyName(name)
}

// ResolveByName is an alias for ResolveName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.ResolveByName("example.lthn")
func (s *Service) ResolveByName(name any) (primitives.Hash, error) {
	return s.ResolveName(name)
}

// GetResolveByName is an alias for ResolveByName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.GetResolveByName("example.lthn")
func (s *Service) GetResolveByName(name any) (primitives.Hash, error) {
	return s.ResolveByName(name)
}

// HashByName is an alias for HashName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.HashByName("example.lthn")
func (s *Service) HashByName(name any) (primitives.Hash, error) {
	return s.HashName(name)
}

// GetHashByName is an alias for HashByName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	hash, err := svc.GetHashByName("example.lthn")
func (s *Service) GetHashByName(name any) (primitives.Hash, error) {
	return s.HashByName(name)
}

// VerifyByName is an alias for VerifyName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.VerifyByName("example.lthn")
func (s *Service) VerifyByName(name any) bool {
	return s.VerifyName(name)
}

// GetVerifyByName is an alias for VerifyByName.
//
//	c := core.New(core.WithService(lns.Register))
//	svc := c.Service("lns").(*lns.Service)
//	ok := svc.GetVerifyByName("example.lthn")
func (s *Service) GetVerifyByName(name any) bool {
	return s.VerifyByName(name)
}

// VerifyString reports whether a .lthn name supplied as a string is valid.
//
//	ok := svc.VerifyString("example.lthn")
func (s *Service) VerifyString(name string) bool {
	return verifyCanonicalName(name)
}

// GetVerifyString is an alias for VerifyString.
//
//	ok := svc.GetVerifyString("example.lthn")
func (s *Service) GetVerifyString(name string) bool {
	return s.VerifyString(name)
}

// VerifyBinary reports whether a .lthn name supplied as bytes is valid.
//
//	ok := svc.VerifyBinary([]byte("example.lthn"))
func (s *Service) VerifyBinary(name []byte) bool {
	return verifyCanonicalName(name)
}

// GetVerifyBinary is an alias for VerifyBinary.
//
//	ok := svc.GetVerifyBinary([]byte("example.lthn"))
func (s *Service) GetVerifyBinary(name []byte) bool {
	return s.VerifyBinary(name)
}

// VerifyByString is an alias for VerifyString.
//
//	ok := svc.VerifyByString("example.lthn")
func (s *Service) VerifyByString(name string) bool {
	return s.VerifyString(name)
}

// GetVerifyByString is an alias for VerifyByString.
//
//	ok := svc.GetVerifyByString("example.lthn")
func (s *Service) GetVerifyByString(name string) bool {
	return s.VerifyByString(name)
}

// VerifyByBinary is an alias for VerifyBinary.
//
//	ok := svc.VerifyByBinary([]byte("example.lthn"))
func (s *Service) VerifyByBinary(name []byte) bool {
	return s.VerifyBinary(name)
}

// GetVerifyByBinary is an alias for VerifyByBinary.
//
//	ok := svc.GetVerifyByBinary([]byte("example.lthn"))
func (s *Service) GetVerifyByBinary(name []byte) bool {
	return s.VerifyByBinary(name)
}

// Blind returns the sealed-bid commitment for an amount and nonce.
//
//	bid, err := svc.Blind(1000, nonce)
func (s *Service) Blind(value uint64, nonce primitives.Hash) (primitives.Hash, error) {
	return covenant.Blind(value, nonce)
}

// GetBlind is an alias for Blind.
//
//	bid, err := svc.GetBlind(1000, nonce)
func (s *Service) GetBlind(value uint64, nonce primitives.Hash) (primitives.Hash, error) {
	return s.Blind(value, nonce)
}

// TypeName returns the canonical string for a covenant type.
//
//	name := svc.TypeName(covenant.TypeBid)
func (s *Service) TypeName(ct covenant.CovenantType) string {
	return covenant.TypeName(ct)
}

// GetTypeName is an alias for TypeName.
//
//	name := svc.GetTypeName(covenant.TypeBid)
func (s *Service) GetTypeName(ct covenant.CovenantType) string {
	return s.TypeName(ct)
}

// IsName reports whether a covenant type is a name operation.
//
//	ok := svc.IsName(covenant.TypeOpen)
func (s *Service) IsName(ct CovenantType) bool {
	return ct.IsName()
}

// GetIsName is an alias for IsName.
//
//	ok := svc.GetIsName(covenant.TypeOpen)
func (s *Service) GetIsName(ct CovenantType) bool {
	return s.IsName(ct)
}

// IsKnown reports whether a covenant type is recognized by the rules table.
//
//	ok := svc.IsKnown(covenant.TypeBid)
func (s *Service) IsKnown(ct CovenantType) bool {
	return ct.IsKnown()
}

// GetIsKnown is an alias for IsKnown.
//
//	ok := svc.GetIsKnown(covenant.TypeBid)
func (s *Service) GetIsKnown(ct CovenantType) bool {
	return s.IsKnown(ct)
}

// IsLinked reports whether a covenant type is a linked covenant.
//
//	ok := svc.IsLinked(covenant.TypeReveal)
func (s *Service) IsLinked(ct CovenantType) bool {
	return ct.IsLinked()
}

// GetIsLinked is an alias for IsLinked.
//
//	ok := svc.GetIsLinked(covenant.TypeReveal)
func (s *Service) GetIsLinked(ct CovenantType) bool {
	return s.IsLinked(ct)
}

// GetRollout returns the rollout start height and rollout week for a name hash.
func (s *Service) GetRollout(nameHash primitives.Hash, rules NameRules) (uint32, uint32) {
	return covenant.GetRollout(nameHash, rules)
}

// HasRollout reports whether a name hash has reached its rollout height.
func (s *Service) HasRollout(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return covenant.HasRollout(nameHash, height, rules)
}

// IsReserved reports whether a name hash is reserved at the given height.
func (s *Service) IsReserved(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return covenant.IsReserved(nameHash, height, rules)
}

// IsLockedUp reports whether a name hash is locked at the given height.
func (s *Service) IsLockedUp(nameHash primitives.Hash, height uint32, rules NameRules) bool {
	return covenant.IsLockedUp(nameHash, height, rules)
}

// GrindName searches for an available name that satisfies the current rollout
// and reservation rules.
func (s *Service) GrindName(size int, height uint32, rules NameRules) (string, error) {
	return covenant.GrindName(size, height, rules)
}

// GetGrindName is an alias for GrindName.
func (s *Service) GetGrindName(size int, height uint32, rules NameRules) (string, error) {
	return s.GrindName(size, height, rules)
}

// CountOpens counts OPEN covenants in a transaction.
func (s *Service) CountOpens(tx primitives.Transaction) int {
	return covenant.CountOpens(tx)
}

// CountUpdates counts update-style covenants in a transaction.
func (s *Service) CountUpdates(tx primitives.Transaction) int {
	return covenant.CountUpdates(tx)
}

// CountRenewals counts REGISTER, RENEW, and FINALIZE covenants in a transaction.
func (s *Service) CountRenewals(tx primitives.Transaction) int {
	return covenant.CountRenewals(tx)
}

// HasNames reports whether any name covenant in the transaction matches the set.
func (s *Service) HasNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) bool {
	return covenant.HasNames(tx, set)
}

// AddNames adds the transaction's name covenants to the set.
func (s *Service) AddNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) {
	covenant.AddNames(tx, set)
}

// RemoveNames removes the transaction's name covenants from the set.
func (s *Service) RemoveNames(tx primitives.Transaction, set map[primitives.Hash]struct{}) {
	covenant.RemoveNames(tx, set)
}

// HasSaneCovenants performs structural covenant validation on a transaction.
func (s *Service) HasSaneCovenants(tx primitives.Transaction) bool {
	return covenant.HasSaneCovenants(tx)
}

// VerifyCovenants performs contextual covenant validation on a transaction.
func (s *Service) VerifyCovenants(tx primitives.Transaction, view CoinView, height uint32, network Network) int {
	return covenant.VerifyCovenants(tx, view, height, network)
}

// GetVerifyCovenants is an alias for VerifyCovenants.
func (s *Service) GetVerifyCovenants(tx primitives.Transaction, view CoinView, height uint32, network Network) int {
	return s.VerifyCovenants(tx, view, height, network)
}

// Types mirrors the covenant type lookup table from pkg/covenant.
func (s *Service) Types() map[string]covenant.CovenantType {
	return covenant.Types
}

// GetTypes is an alias for Types.
func (s *Service) GetTypes() map[string]covenant.CovenantType {
	return s.Types()
}

// TypesByVal mirrors the covenant reverse lookup table from pkg/covenant.
func (s *Service) TypesByVal() map[covenant.CovenantType]string {
	return covenant.TypesByVal
}

// GetTypesByVal is an alias for TypesByVal.
func (s *Service) GetTypesByVal() map[covenant.CovenantType]string {
	return s.TypesByVal()
}

// Blacklist mirrors the covenant verifier blacklist from pkg/covenant.
func (s *Service) Blacklist() map[string]struct{} {
	return covenant.Blacklist
}

// GetBlacklist is an alias for Blacklist.
func (s *Service) GetBlacklist() map[string]struct{} {
	return s.Blacklist()
}

// VerificationFlags mirrors the covenant verification flag lookup table.
func (s *Service) VerificationFlags() map[string]covenant.NameFlags {
	return covenant.VerificationFlags
}

// GetVerificationFlags is an alias for VerificationFlags.
func (s *Service) GetVerificationFlags() map[string]covenant.NameFlags {
	return s.VerificationFlags()
}

// VerificationFlagsByVal mirrors the covenant verification flag reverse table.
func (s *Service) VerificationFlagsByVal() map[covenant.NameFlags]string {
	return covenant.VerificationFlagsByVal
}

// GetVerificationFlagsByVal is an alias for VerificationFlagsByVal.
func (s *Service) GetVerificationFlagsByVal() map[covenant.NameFlags]string {
	return s.VerificationFlagsByVal()
}

// GetReserved looks up a reserved .lthn name after canonicalisation.
//
// The lookup accepts the same inputs as Resolve and Verify, then returns the
// matching reserved-name entry from the reference catalog when present.
//
//	item, ok := svc.GetReserved("reserved.lthn")
func (s *Service) GetReserved(name any) (covenant.ReservedName, bool) {
	return lookupReserved(name)
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
	return s.HasReservedName(name)
}

// HasReservedBinary reports whether a reserved name supplied as bytes exists.
//
//	ok := svc.HasReservedBinary([]byte("reserved"))
func (s *Service) HasReservedBinary(name []byte) bool {
	return s.HasReservedName(name)
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
	return s.HasReserved(name)
}

// GetReservedHash looks up a reserved .lthn name by its canonical hash.
//
// This mirrors the covenant package lookup helper for callers that already
// have the hashed form available.
//
//	hash, _ := covenant.HashString("reserved")
//	item, ok := svc.GetReservedHash(hash)
func (s *Service) GetReservedHash(hash primitives.Hash) (covenant.ReservedName, bool) {
	return s.ReservedCatalog().Get(hash)
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
	return lookupLocked(name)
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
	return s.HasLockedName(name)
}

// HasLockedBinary reports whether a locked name supplied as bytes exists.
//
//	ok := svc.HasLockedBinary([]byte("nec"))
func (s *Service) HasLockedBinary(name []byte) bool {
	return s.HasLockedName(name)
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
	return s.HasLocked(name)
}

// GetLockedHash looks up a locked .lthn name by its canonical hash.
//
// This mirrors the covenant package lookup helper for callers that already
// have the hashed form available.
//
//	hash, _ := covenant.HashString("nec")
//	item, ok := svc.GetLockedHash(hash)
func (s *Service) GetLockedHash(hash primitives.Hash) (covenant.LockedName, bool) {
	return s.LockedCatalog().Get(hash)
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
	return s.reservedCatalog()
}

// Core returns the Core instance this service is registered with.
//
// The accessor is nil-safe so zero-value services and nil receivers stay
// usable in lightweight tests and helper code.
//
//	core := svc.Core()
func (s *Service) Core() *core.Core {
	if s == nil || s.ServiceRuntime == nil {
		return nil
	}

	return s.ServiceRuntime.Core()
}

// reservedCatalog returns the package-backed reserved-name catalog.
func (s *Service) reservedCatalog() *covenant.ReservedCatalog {
	return Reserved
}

// lockedCatalog returns the package-backed locked-name catalog.
func (s *Service) lockedCatalog() *covenant.LockedCatalog {
	return Locked
}

// GetCore is an alias for Core.
//
//	core := svc.GetCore()
func (s *Service) GetCore() *core.Core {
	return s.Core()
}

// GetReservedCatalog is an alias for ReservedCatalog.
//
//	catalog := svc.GetReservedCatalog()
func (s *Service) GetReservedCatalog() *covenant.ReservedCatalog {
	return s.ReservedCatalog()
}

// DefaultReservedCatalog returns the reserved-name catalog used by the service.
//
//	catalog := svc.DefaultReservedCatalog()
func (s *Service) DefaultReservedCatalog() *covenant.ReservedCatalog {
	return s.ReservedCatalog()
}

// GetDefaultReservedCatalog is an alias for DefaultReservedCatalog.
//
//	catalog := svc.GetDefaultReservedCatalog()
func (s *Service) GetDefaultReservedCatalog() *covenant.ReservedCatalog {
	return s.DefaultReservedCatalog()
}

// DefaultReservedCatalog exposes the reserved-name catalog used by LNS.
//
// This returns the package-owned reserved catalog so callers stay within the
// top-level lns package.
//
//	catalog := lns.DefaultReservedCatalog()
func DefaultReservedCatalog() *covenant.ReservedCatalog {
	return Reserved
}

// GetReservedCatalog is an alias for DefaultReservedCatalog.
//
//	catalog := lns.GetReservedCatalog()
func GetReservedCatalog() *covenant.ReservedCatalog {
	return DefaultReservedCatalog()
}

// GetDefaultReservedCatalog is an alias for DefaultReservedCatalog.
//
//	catalog := lns.GetDefaultReservedCatalog()
func GetDefaultReservedCatalog() *covenant.ReservedCatalog {
	return DefaultReservedCatalog()
}

// ReservedSize reports the number of reserved-name entries available to the service.
//
//	size := svc.ReservedSize()
func (s *Service) ReservedSize() int {
	return s.reservedCatalog().Size()
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
	return s.reservedCatalog().NameValue()
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
	return s.reservedCatalog().RootValue()
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
	return s.reservedCatalog().TopValue()
}

// GetReservedTopValue is an alias for ReservedTopValue.
//
//	top := svc.GetReservedTopValue()
func (s *Service) GetReservedTopValue() uint64 {
	return s.ReservedTopValue()
}

// NameValue is an alias for ReservedNameValue.
//
//	base := svc.NameValue()
func (s *Service) NameValue() uint64 {
	return s.ReservedNameValue()
}

// GetNameValue is an alias for NameValue.
//
//	base := svc.GetNameValue()
func (s *Service) GetNameValue() uint64 {
	return s.NameValue()
}

// RootValue is an alias for ReservedRootValue.
//
//	root := svc.RootValue()
func (s *Service) RootValue() uint64 {
	return s.ReservedRootValue()
}

// GetRootValue is an alias for RootValue.
//
//	root := svc.GetRootValue()
func (s *Service) GetRootValue() uint64 {
	return s.RootValue()
}

// TopValue is an alias for ReservedTopValue.
//
//	top := svc.TopValue()
func (s *Service) TopValue() uint64 {
	return s.ReservedTopValue()
}

// GetTopValue is an alias for TopValue.
//
//	top := svc.GetTopValue()
func (s *Service) GetTopValue() uint64 {
	return s.TopValue()
}

// PrefixSize is an alias for ReservedPrefixSize.
//
//	offset := svc.PrefixSize()
func (s *Service) PrefixSize() int {
	return s.ReservedPrefixSize()
}

// GetPrefixSize is an alias for PrefixSize.
//
//	offset := svc.GetPrefixSize()
func (s *Service) GetPrefixSize() int {
	return s.PrefixSize()
}

// ReservedEntries exposes the reserved-name catalog entries used by the service.
//
//	entries := svc.ReservedEntries()
func (s *Service) ReservedEntries() []covenant.ReservedEntry {
	return s.reservedCatalog().Entries()
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
	return s.reservedCatalog().Keys()
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
	return s.reservedCatalog().Values()
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
	return s.lockedCatalog()
}

// GetLockedCatalog is an alias for LockedCatalog.
//
//	catalog := svc.GetLockedCatalog()
func (s *Service) GetLockedCatalog() *covenant.LockedCatalog {
	return s.LockedCatalog()
}

// DefaultLockedCatalog returns the locked-name catalog used by the service.
//
//	catalog := svc.DefaultLockedCatalog()
func (s *Service) DefaultLockedCatalog() *covenant.LockedCatalog {
	return s.LockedCatalog()
}

// GetDefaultLockedCatalog is an alias for DefaultLockedCatalog.
//
//	catalog := svc.GetDefaultLockedCatalog()
func (s *Service) GetDefaultLockedCatalog() *covenant.LockedCatalog {
	return s.DefaultLockedCatalog()
}

// DefaultLockedCatalog exposes the locked-name catalog used by LNS.
//
// This returns the package-owned locked catalog so callers stay within the
// top-level lns package.
//
//	catalog := lns.DefaultLockedCatalog()
func DefaultLockedCatalog() *covenant.LockedCatalog {
	return Locked
}

// GetLockedCatalog is an alias for DefaultLockedCatalog.
//
//	catalog := lns.GetLockedCatalog()
func GetLockedCatalog() *covenant.LockedCatalog {
	return DefaultLockedCatalog()
}

// GetDefaultLockedCatalog is an alias for DefaultLockedCatalog.
//
//	catalog := lns.GetDefaultLockedCatalog()
func GetDefaultLockedCatalog() *covenant.LockedCatalog {
	return DefaultLockedCatalog()
}

// Reserved exposes the reserved-name catalog instance used by the package.
//
//	catalog := lns.Reserved
var Reserved = covenant.DefaultReservedCatalog()

// Locked exposes the locked-name catalog instance used by the package.
//
//	catalog := lns.Locked
var Locked = covenant.DefaultLockedCatalog()

// ReservedCatalog exposes the reserved-name catalog used by the package.
//
//	catalog := lns.ReservedCatalog()
func ReservedCatalog() *covenant.ReservedCatalog {
	return Reserved
}

// LockedCatalog exposes the locked-name catalog used by the package.
//
//	catalog := lns.LockedCatalog()
func LockedCatalog() *covenant.LockedCatalog {
	return Locked
}

// ReservedSize reports the number of reserved-name entries available to the package.
//
//	size := lns.ReservedSize()
func ReservedSize() int {
	return ReservedCatalog().Size()
}

// GetReservedSize is an alias for ReservedSize.
//
//	size := lns.GetReservedSize()
func GetReservedSize() int {
	return ReservedSize()
}

// ReservedNameValue reports the base value used for reserved-name lookups.
//
//	base := lns.ReservedNameValue()
func ReservedNameValue() uint64 {
	return ReservedCatalog().NameValue()
}

// GetReservedNameValue is an alias for ReservedNameValue.
//
//	base := lns.GetReservedNameValue()
func GetReservedNameValue() uint64 {
	return ReservedNameValue()
}

// ReservedRootValue reports the root-name increment used for reserved-name lookups.
//
//	root := lns.ReservedRootValue()
func ReservedRootValue() uint64 {
	return ReservedCatalog().RootValue()
}

// GetReservedRootValue is an alias for ReservedRootValue.
//
//	root := lns.GetReservedRootValue()
func GetReservedRootValue() uint64 {
	return ReservedRootValue()
}

// ReservedTopValue reports the top-100 increment used for reserved-name lookups.
//
//	top := lns.ReservedTopValue()
func ReservedTopValue() uint64 {
	return ReservedCatalog().TopValue()
}

// GetReservedTopValue is an alias for ReservedTopValue.
//
//	top := lns.GetReservedTopValue()
func GetReservedTopValue() uint64 {
	return ReservedTopValue()
}

// NameValue is an alias for ReservedNameValue.
//
//	base := lns.NameValue()
func NameValue() uint64 {
	return ReservedNameValue()
}

// GetNameValue is an alias for NameValue.
//
//	base := lns.GetNameValue()
func GetNameValue() uint64 {
	return NameValue()
}

// RootValue is an alias for ReservedRootValue.
//
//	root := lns.RootValue()
func RootValue() uint64 {
	return ReservedRootValue()
}

// GetRootValue is an alias for RootValue.
//
//	root := lns.GetRootValue()
func GetRootValue() uint64 {
	return RootValue()
}

// TopValue is an alias for ReservedTopValue.
//
//	top := lns.TopValue()
func TopValue() uint64 {
	return ReservedTopValue()
}

// GetTopValue is an alias for TopValue.
//
//	top := lns.GetTopValue()
func GetTopValue() uint64 {
	return TopValue()
}

// PrefixSize is an alias for ReservedPrefixSize.
//
//	offset := lns.PrefixSize()
func PrefixSize() int {
	return ReservedPrefixSize()
}

// GetPrefixSize is an alias for PrefixSize.
//
//	offset := lns.GetPrefixSize()
func GetPrefixSize() int {
	return PrefixSize()
}

// ReservedPrefixSize reports the byte offset where reserved-name entries begin.
//
//	offset := lns.ReservedPrefixSize()
func ReservedPrefixSize() int {
	return ReservedCatalog().PrefixSize()
}

// GetReservedPrefixSize is an alias for ReservedPrefixSize.
//
//	offset := lns.GetReservedPrefixSize()
func GetReservedPrefixSize() int {
	return ReservedPrefixSize()
}

// ReservedEntries exposes the reserved-name catalog entries used by the package.
//
//	entries := lns.ReservedEntries()
func ReservedEntries() []covenant.ReservedEntry {
	return ReservedCatalog().Entries()
}

// GetReservedEntries is an alias for ReservedEntries.
//
//	entries := lns.GetReservedEntries()
func GetReservedEntries() []covenant.ReservedEntry {
	return ReservedEntries()
}

// ReservedKeys exposes the reserved-name hashes used by the package.
//
//	keys := lns.ReservedKeys()
func ReservedKeys() []primitives.Hash {
	return ReservedCatalog().Keys()
}

// GetReservedKeys is an alias for ReservedKeys.
//
//	keys := lns.GetReservedKeys()
func GetReservedKeys() []primitives.Hash {
	return ReservedKeys()
}

// ReservedValues exposes the reserved-name entries used by the package.
//
//	values := lns.ReservedValues()
func ReservedValues() []covenant.ReservedName {
	return ReservedCatalog().Values()
}

// GetReservedValues is an alias for ReservedValues.
//
//	values := lns.GetReservedValues()
func GetReservedValues() []covenant.ReservedName {
	return ReservedValues()
}

// LockedSize reports the number of locked-name entries available to the package.
//
//	size := lns.LockedSize()
func LockedSize() int {
	return LockedCatalog().Size()
}

// GetLockedSize is an alias for LockedSize.
//
//	size := lns.GetLockedSize()
func GetLockedSize() int {
	return LockedSize()
}

// LockedPrefixSize reports the byte offset where locked-name entries begin.
//
//	offset := lns.LockedPrefixSize()
func LockedPrefixSize() int {
	return LockedCatalog().PrefixSize()
}

// GetLockedPrefixSize is an alias for LockedPrefixSize.
//
//	offset := lns.GetLockedPrefixSize()
func GetLockedPrefixSize() int {
	return LockedPrefixSize()
}

// LockedEntries exposes the locked-name catalog entries used by the package.
//
//	entries := lns.LockedEntries()
func LockedEntries() []covenant.LockedEntry {
	return LockedCatalog().Entries()
}

// GetLockedEntries is an alias for LockedEntries.
//
//	entries := lns.GetLockedEntries()
func GetLockedEntries() []covenant.LockedEntry {
	return LockedEntries()
}

// LockedKeys exposes the locked-name hashes used by the package.
//
//	keys := lns.LockedKeys()
func LockedKeys() []primitives.Hash {
	return LockedCatalog().Keys()
}

// GetLockedKeys is an alias for LockedKeys.
//
//	keys := lns.GetLockedKeys()
func GetLockedKeys() []primitives.Hash {
	return LockedKeys()
}

// LockedValues exposes the locked-name entries used by the package.
//
//	values := lns.LockedValues()
func LockedValues() []covenant.LockedName {
	return LockedCatalog().Values()
}

// GetLockedValues is an alias for LockedValues.
//
//	values := lns.GetLockedValues()
func GetLockedValues() []covenant.LockedName {
	return LockedValues()
}

// GetReserved looks up a reserved .lthn name after canonicalisation.
//
//	item, ok := lns.GetReserved("reserved.lthn")
func GetReserved(name any) (covenant.ReservedName, bool) {
	return lookupReserved(name)
}

// GetReservedName looks up a reserved name after canonicalisation.
//
//	item, ok := lns.GetReservedName("reserved")
func GetReservedName(name any) (covenant.ReservedName, bool) {
	return lookupReserved(name)
}

// GetReservedString looks up a reserved name supplied as a string.
//
//	item, ok := lns.GetReservedString("reserved")
func GetReservedString(name string) (covenant.ReservedName, bool) {
	return GetReservedName(name)
}

// GetReservedBinary looks up a reserved name supplied as bytes.
//
//	item, ok := lns.GetReservedBinary([]byte("reserved"))
func GetReservedBinary(name []byte) (covenant.ReservedName, bool) {
	return GetReservedName(name)
}

// GetReservedByString is an alias for GetReservedString.
//
//	item, ok := lns.GetReservedByString("reserved")
func GetReservedByString(name string) (covenant.ReservedName, bool) {
	return GetReservedString(name)
}

// GetReservedByBinary is an alias for GetReservedBinary.
//
//	item, ok := lns.GetReservedByBinary([]byte("reserved"))
func GetReservedByBinary(name []byte) (covenant.ReservedName, bool) {
	return GetReservedBinary(name)
}

// HasReserved reports whether a canonicalised .lthn name is reserved.
//
//	ok := lns.HasReserved("reserved.lthn")
func HasReserved(name any) bool {
	_, ok := GetReserved(name)
	return ok
}

// HasReservedString reports whether a reserved name supplied as a string exists.
//
//	ok := lns.HasReservedString("reserved")
func HasReservedString(name string) bool {
	return HasReservedName(name)
}

// HasReservedBinary reports whether a reserved name supplied as bytes exists.
//
//	ok := lns.HasReservedBinary([]byte("reserved"))
func HasReservedBinary(name []byte) bool {
	return HasReservedName(name)
}

// HasReservedByString is an alias for HasReservedString.
//
//	ok := lns.HasReservedByString("reserved")
func HasReservedByString(name string) bool {
	return HasReservedString(name)
}

// HasReservedByBinary is an alias for HasReservedBinary.
//
//	ok := lns.HasReservedByBinary([]byte("reserved"))
func HasReservedByBinary(name []byte) bool {
	return HasReservedBinary(name)
}

// GetReservedByName is an alias for GetReservedName.
//
//	item, ok := lns.GetReservedByName("reserved")
func GetReservedByName(name any) (covenant.ReservedName, bool) {
	return GetReservedName(name)
}

// HasReservedByName is an alias for HasReservedName.
//
//	ok := lns.HasReservedByName("reserved")
func HasReservedByName(name any) bool {
	return HasReservedName(name)
}

// HasReservedName reports whether a reserved name exists using the raw catalog label.
//
//	ok := lns.HasReservedName("reserved")
func HasReservedName(name any) bool {
	return HasReserved(name)
}

// GetReservedHash looks up a reserved .lthn name by its canonical hash.
//
//	hash, _ := covenant.HashString("reserved")
//	item, ok := lns.GetReservedHash(hash)
func GetReservedHash(hash primitives.Hash) (covenant.ReservedName, bool) {
	return ReservedCatalog().Get(hash)
}

// HasReservedHash reports whether a canonical hash is reserved.
//
//	ok := lns.HasReservedHash(hash)
func HasReservedHash(hash primitives.Hash) bool {
	_, ok := GetReservedHash(hash)
	return ok
}

// GetReservedByHash is an alias for GetReservedHash.
//
//	item, ok := lns.GetReservedByHash(hash)
func GetReservedByHash(hash primitives.Hash) (covenant.ReservedName, bool) {
	return GetReservedHash(hash)
}

// HasReservedByHash is an alias for HasReservedHash.
//
//	ok := lns.HasReservedByHash(hash)
func HasReservedByHash(hash primitives.Hash) bool {
	return HasReservedHash(hash)
}

// GetLocked looks up a locked .lthn name after canonicalisation.
//
//	item, ok := lns.GetLocked("nec.lthn")
func GetLocked(name any) (covenant.LockedName, bool) {
	return lookupLocked(name)
}

// GetLockedName looks up a locked name after canonicalisation.
//
//	item, ok := lns.GetLockedName("nec")
func GetLockedName(name any) (covenant.LockedName, bool) {
	return lookupLocked(name)
}

// GetLockedString looks up a locked name supplied as a string.
//
//	item, ok := lns.GetLockedString("nec")
func GetLockedString(name string) (covenant.LockedName, bool) {
	return GetLockedName(name)
}

// GetLockedBinary looks up a locked name supplied as bytes.
//
//	item, ok := lns.GetLockedBinary([]byte("nec"))
func GetLockedBinary(name []byte) (covenant.LockedName, bool) {
	return GetLockedName(name)
}

// GetLockedByString is an alias for GetLockedString.
//
//	item, ok := lns.GetLockedByString("nec")
func GetLockedByString(name string) (covenant.LockedName, bool) {
	return GetLockedString(name)
}

// GetLockedByBinary is an alias for GetLockedBinary.
//
//	item, ok := lns.GetLockedByBinary([]byte("nec"))
func GetLockedByBinary(name []byte) (covenant.LockedName, bool) {
	return GetLockedBinary(name)
}

// HasLocked reports whether a canonicalised .lthn name is locked.
//
//	ok := lns.HasLocked("nec.lthn")
func HasLocked(name any) bool {
	_, ok := GetLocked(name)
	return ok
}

// HasLockedString reports whether a locked name supplied as a string exists.
//
//	ok := lns.HasLockedString("nec")
func HasLockedString(name string) bool {
	return HasLockedName(name)
}

// HasLockedBinary reports whether a locked name supplied as bytes exists.
//
//	ok := lns.HasLockedBinary([]byte("nec"))
func HasLockedBinary(name []byte) bool {
	return HasLockedName(name)
}

// HasLockedByString is an alias for HasLockedString.
//
//	ok := lns.HasLockedByString("nec")
func HasLockedByString(name string) bool {
	return HasLockedString(name)
}

// HasLockedByBinary is an alias for HasLockedBinary.
//
//	ok := lns.HasLockedByBinary([]byte("nec"))
func HasLockedByBinary(name []byte) bool {
	return HasLockedBinary(name)
}

// GetLockedByName is an alias for GetLockedName.
//
//	item, ok := lns.GetLockedByName("nec")
func GetLockedByName(name any) (covenant.LockedName, bool) {
	return GetLockedName(name)
}

// HasLockedByName is an alias for HasLockedName.
//
//	ok := lns.HasLockedByName("nec")
func HasLockedByName(name any) bool {
	return HasLockedName(name)
}

// HasLockedName reports whether a locked name exists using the raw catalog label.
//
//	ok := lns.HasLockedName("nec")
func HasLockedName(name any) bool {
	return HasLocked(name)
}

// GetLockedHash looks up a locked .lthn name by its canonical hash.
//
//	hash, _ := covenant.HashString("nec")
//	item, ok := lns.GetLockedHash(hash)
func GetLockedHash(hash primitives.Hash) (covenant.LockedName, bool) {
	return LockedCatalog().Get(hash)
}

// HasLockedHash reports whether a canonical hash is locked.
//
//	ok := lns.HasLockedHash(hash)
func HasLockedHash(hash primitives.Hash) bool {
	_, ok := GetLockedHash(hash)
	return ok
}

// GetLockedByHash is an alias for GetLockedHash.
//
//	item, ok := lns.GetLockedByHash(hash)
func GetLockedByHash(hash primitives.Hash) (covenant.LockedName, bool) {
	return GetLockedHash(hash)
}

// HasLockedByHash is an alias for HasLockedHash.
//
//	ok := lns.HasLockedByHash(hash)
func HasLockedByHash(hash primitives.Hash) bool {
	return HasLockedHash(hash)
}

// LockedSize reports the number of locked-name entries available to the service.
//
//	size := svc.LockedSize()
func (s *Service) LockedSize() int {
	return s.lockedCatalog().Size()
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
	return s.lockedCatalog().PrefixSize()
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
	return s.lockedCatalog().Entries()
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
	return s.lockedCatalog().Keys()
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
	return s.lockedCatalog().Values()
}

// GetLockedValues is an alias for LockedValues.
//
//	values := svc.GetLockedValues()
func (s *Service) GetLockedValues() []covenant.LockedName {
	return s.LockedValues()
}
