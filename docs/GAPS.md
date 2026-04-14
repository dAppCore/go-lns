# go-lns Feature Gap Report

**Date:** 2026-04-14
**Audit source:** Sonnet audit against `plans/code/core/go/lns/RFC.md` + `plans/code/core/network/RFC.lns.md`
**Current coverage:** ~8–10% — clean HSD wire-format types, zero Core action surface

## Authoritative specs for this repo

| Spec | Scope | Location |
|------|-------|----------|
| `go/lns/RFC.md` | HSD protocol layer — covenants, primitives, wallet, chain, mining, mempool | `docs/RFC.md` |
| `network/RFC.lns.md` | Redirect → canonical at `code/core/go/blockchain/RFC.lns.md` (LNS server binary) | `docs/RFC.lns.md` |

`go/lns/RFC.md` also contains an embedded second spec (`go-lns — Lethean Name Server`) that blurs the boundary with go-dns. Scope should be split:

- **go-lns** = HSD protocol layer (this repo) — covenants, primitives, DNS record encoding, shared HSD RPC client
- **go-dns** = resolver service — consumes go-lns types, serves DNS, HTTP API

## What's already here (keep)

- `pkg/covenant/` — 13 files. Covenant types, namestate, outpoint, locked/reserved lookup. ~60% of §8.2 sub-spec.
- `pkg/primitives/` — 14 files. Partial — Address, Name present. Block and TX types missing.
- `pkg/dns/` — 6 files. HSD record types (GLUE4, GLUE6, TXT, NS, DS, NSEC), resource.go, nsec.go, resolve.go. **FULL coverage of §2.**
- `internal/nameutil/` — name normalisation helpers.

## Top blockers (in priority order)

### 1. No Core action surface at all (CRITICAL)

`OnStartup()` is a no-op. `HandleIPCEvents()` is a no-op. **All 24 RFC actions absent:**

**Name actions (§3.1 — 8 actions):**
- `lns.name.register` / `lns.name.renew` / `lns.name.transfer` / `lns.name.update` / `lns.name.revoke` / `lns.name.lookup` / `lns.name.auction` / `lns.name.available`

**DNS actions (§3.2 — 7 actions):**
- `lns.dns.resolve` / `lns.dns.resolve.txt` / `lns.dns.resolve.all` / `lns.dns.reverse` / `lns.dns.serve` / `lns.dns.health` / `lns.dns.discover`

**Wallet actions (§3.3 — 7 actions):**
- `lns.wallet.*` — entire `pkg/wallet` missing

**Chain actions (§3.4 — 5 actions):**
- `lns.chain.*` — entire `pkg/chain` missing

**VPN actions (§3.5 — 4 actions):**
- `lns.vpn.*` — no wiring

### 2. Missing `pkg/` packages (§8)

| Package | Spec | Status |
|---------|------|--------|
| `pkg/wallet` | §8.4 | MISSING entirely |
| `pkg/chain` | §8.5 | MISSING entirely |
| `pkg/script` | §8.6 | MISSING entirely |
| `pkg/hd` | §8.7 | MISSING entirely (HD wallet derivation) |
| `pkg/net` | §8.8 | MISSING entirely (P2P networking) |
| `pkg/blockstore` | §8.9 | MISSING entirely |
| `pkg/mining` | §8.10 | MISSING entirely |
| `pkg/mempool` | §8.11 | MISSING entirely |

### 3. No HSD RPC client

Shared HSD RPC client missing. This blocks both go-lns and go-dns. Recommended: extract a shared `go-hsd` package consumed by both (per `go/dns/RFC.md §8`).

### 4. No live resolution pipeline

`pkg/dns/resolve.go` implements hash derivation only — not live resolution. No chain connectivity → HSD lookup → cache flow.

### 5. `pkg/primitives` gaps

- Block type missing
- TX type missing
- Current code has covenant items, namestate, outpoint only

### 6. CoreApp build tag integration (§6)

`//go:build lns` conditional compilation not wired — can't conditionally include LNS into CoreApp builds.

### 7. CoreGUI custom elements (§4)

Five web components specced, zero implemented.

### 8. CorePlay integration (§7)

Eight messaging/team actions specced for CorePlay integration, zero implemented.

### 9. Three LNS modes (blockchain/RFC.lns.md §1.1)

light / full / gateway modes — none present. This belongs in the LNS server binary (possibly a separate `core/lns` cmd in this repo or go-dns).

### 10. Advanced features (all 0%)

- DNSSEC chain-of-trust (DS → DNSKEY → RRSIG)
- SPV light-mode tree-root proof verification
- Subdomain failover (parent serves cached child)
- Trust accrual metrics (gateway mode)
- 8-PAC root zone governance via alias comment
- Wildcard record support
- Zero-downtime mode migration (light → full)
- Hash-derived cache structure

## Recommended build order

1. **Extract `go-hsd` shared RPC client** — unblocks both repos and prevents three separate implementations (currently scattered across `lthn/lns`, alias-dns-bridge, go-dns).
2. Complete `pkg/primitives` — add Block + TX types so covenants have their host structures.
3. Implement `pkg/chain` (chain sync) and `pkg/blockstore` (persistence) — required by all name and DNS actions.
4. Register Core actions: name → dns → chain → wallet → vpn (in that order of dependency depth).
5. Add `pkg/wallet` + `pkg/hd` (wallet + HD derivation) — enables name registration.
6. `pkg/mempool` + `pkg/mining` — lower priority unless go-lns runs as a full node.
7. CoreGUI components + CorePlay integration — last, depends on core action surface existing.

## Reference implementations

- `~/Code/lthn/lns/` — 672 lines of proven testnet code. Migrate or use as reference for HSD client + DNS serving path.
- HSD wire format code in this repo (`pkg/covenant`, `pkg/primitives`, `pkg/dns`) is clean — preserve.

## Known contradictions

`go/lns/RFC.md` has two specs in one file — the HSD protocol layer (lines 1-250) and an embedded LNS server binary spec (line 251+). The embedded server spec overlaps with `go/blockchain/RFC.lns.md` and `go/dns/RFC.md`. Recommendation: split the embedded LNS-server section out into its own file or attribute it to go-dns.
