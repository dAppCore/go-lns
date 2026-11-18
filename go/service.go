// SPDX-License-Identifier: EUPL-1.2

// Service registration aliases for the lns package — exposes the
// canonical `NewService` shape every other go-* repo provides per
// Mantis #1336. The full Service implementation lives in lns.go and
// uses an internal serviceOptions struct + ServiceOption funcs;
// this file is the canonical-name factory surface so consumers
// calling `lns.NewService(opts)` (factory shape, returns
// `func(*core.Core) core.Result`) get the same canon every other
// canonical repo provides.
//
// The original lns.NewService(c *core.Core) → *Service signature is
// preserved for backwards compatibility — consumers using either entry
// point construct the same *Service over the same internal options.

package lns

import (
	core "github.com/dappcore/go"
)

// Options is the public canonical-shape configuration for an LNS
// service. v1 has no fields — the package's existing optional
// behaviours (catalog overrides, DNS bridge) flow through the
// ServiceOption funcs (WithReservedCatalog, WithLockedCatalog, etc.)
// applied via NewServiceWithOptions. Future fields (e.g. NodeID,
// FederationPeers) get added here as they're needed.
type Options struct{}

// NewServiceFactory returns the canonical Core service factory shape —
// `func(*core.Core) core.Result` — that wraps NewService. Use through
// core.WithService:
//
//	core.WithService(lns.NewServiceFactory(lns.Options{}))
//
// Identical to Register(c) for v1 (Options has no fields). Named
// distinctly from NewService to preserve backwards compatibility with
// the existing lns.NewService(c *core.Core) → *Service signature in
// lns.go — both shapes coexist; pick the one matching your boot
// pattern.
//
// When Options grows fields, NewServiceFactory becomes the option-
// bearing entry point and Register remains the no-config shorthand.
func NewServiceFactory(opts Options) func(*core.Core) core.Result {
	_ = opts
	return func(c *core.Core) core.Result {
		return Register(c)
	}
}
