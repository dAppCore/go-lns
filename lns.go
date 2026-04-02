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
	"dappco.re/go/lns/pkg/dns"
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
func (s *Service) Resolve(name any) (primitives.Hash, error) {
	opts := make([]dns.ServiceOption, 0, 1)
	if s != nil && s.ServiceRuntime != nil {
		opts = append(opts, dns.WithCore(s.Core()))
	}

	return dns.NewService(opts...).Resolve(name)
}
