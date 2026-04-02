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

import core "dappco.re/go/core"

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
