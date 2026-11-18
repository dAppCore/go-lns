# code/core/go/lns — Lethean Name System (Native)

> Agent context summary for `plans/code/core/go/lns/`. Read this first, then dive into RFC.md.

## What go/lns Is

Lethean Name System as native CoreGO implementation. Converting lthn/itns-sidechain (118k lines JS/HSD fork) into Go. Provides peer name resolution, trust scoring, and distributed name records without DNS or blockchain bloat. Used by core/agent peer discovery and OFM user/team names.

## Key Facts

- **Module**: `github.com/dappcore/go-lns/go`
- **Status**: Conversion from JS/HSD to Go (118k LoC → lean Go)
- **Purpose**: Peer names, user names, team names, trust scoring
- **API**: `c.Action("lns.resolve", {name, type})`
- **Depends On**: code/snider/poindexter (peer selection)

## Architecture

### Name Resolution Types

| Type | Resolves To |
|------|------------|
| Peer | Peer public key + address (Poindexter proof) |
| User | User identity + avatar + profile hash |
| Team | Team members + settings + capabilities |
| Service | Service endpoint + capability SLA |

### Trust Scoring

Each resolution includes Poindexter trust score (0.0-1.0). Cached for TTL per name. Scores degrade over time, forcing periodic re-verification.

## Critical Rules

- **This is a code/ spec** — self-contained, no project/ references
- **No blockchain coupling** — uses go-blockchain for proofs only, not primary storage
- **Trust via Poindexter** — all scores validated by peer selection logic
- **Distributed cache** — LNS records replicate across peers (P2P gossip)
- **Conversion complete** — no remaining dependencies on JS/HSD fork

## Spec Index

| File | Scope |
|------|-------|
| [RFC.md](RFC.md) | **Main spec** — name types, trust scoring, resolution protocol (~170 lines) |
