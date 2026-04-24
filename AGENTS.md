# Repository Guidelines

## Overview
- Module: `dappco.re/go/lns`
- Purpose: native Go implementation of the Lethean Name System (LNS)
- Primary design docs live in `docs/CLAUDE.md`, `docs/RFC.md`, and `docs/RFC.lns.md`
- The `docs/js-*` trees are reference material from the JS/HSD lineage; treat them as porting references, not the primary source for new Go APIs

## Repository Layout
- `lns.go`: top-level service entrypoint and public API mirror layer
- `internal/nameutil`: internal helpers for canonicalizing and validating labels
- `pkg/covenant`: covenant enums, rule helpers, reserved/locked catalogs, and verification logic
- `pkg/dns`: DNS-oriented resolution, record encoding, and lookup helpers
- `pkg/primitives`: shared LNS data structures, binary/JSON forms, and state types
- `docs/`: RFCs, gap analysis, and JS reference sources

## Editing Rules
- Keep changes surgical and package-local; avoid reshaping unrelated APIs
- Preserve SPDX headers on Go files that already have them
- Match the existing style of rich doc comments for exported identifiers
- Reuse existing helpers instead of reimplementing canonicalization or validation logic
- For `.lthn` name normalization, prefer `internal/nameutil.Canonicalize` or `CatalogLabel`
- When changing public DNS or covenant surface area, check whether `lns.go` should continue mirroring the underlying types, constants, or helpers
- Maintain existing alias patterns such as `Get*` wrappers when touching code that already follows that convention

## Validation
- The parent workspace file at `/Users/snider/Code/go.work` does not include this module
- Do not rely on plain `go test ./...` from this directory; use `GOWORK=off`
- Preferred full validation command: `GOWORK=off go test ./...`
- Prefer narrower package-level runs first when making focused changes, then run the full suite if the surface area justifies it

## Design Constraints
- Keep this package self-contained as native Go code; do not introduce new runtime dependence on the JS reference implementation
- Avoid coupling LNS to unrelated blockchain or storage layers beyond the existing proof/types interfaces already used here
- Treat `docs/RFC*.md` as the contract for behavior when implementation questions arise
