# LNS — Lethean NameServer

> Canonical spec: `code/core/go/blockchain/RFC.lns.md`
>
> This file exists as a pointer. The LNS is a Go binary, so the authoritative spec lives with the Go blockchain specs.

See `code/core/go/blockchain/RFC.lns.md` for the full specification covering:
- Three modes: light, full, gateway
- Hash-derived cache structure with cryptographic access control
- DNSSEC via HSD DS records
- Tree-root invalidation
- Discovery via main chain + sidechain queries
- Subdomain failover
- Metrics and trust accrual
- Integration with core/agent, core/network, Blesta
