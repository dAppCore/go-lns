# LNS Documentation

The Lethean Name System Go module is organized around a small set of native Go
packages with a broad top-level mirror in `lns.go`. Start with the README for
installation and a minimal example, then use this directory for design notes and
porting contracts.

Primary references:

- `RFC.md` and `RFC.lns.md` define behavior and naming rules.
- `architecture.md` describes the Go package boundaries.
- `development.md` describes local validation and contribution workflow.
- `GAPS.md` records remaining differences against the JavaScript lineage.
- `js-*` directories preserve reference source material and fixture data.

For implementation work, treat the Go packages as authoritative unless an RFC
contradicts them. The JavaScript trees are comparison material, not a runtime
dependency.
