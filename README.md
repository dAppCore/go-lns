# Lethean Name System for Go

This module is the native Go implementation of the Lethean Name System (LNS).
It provides name validation, name hashing, reserved and locked catalog lookup,
covenant helpers, DNS resource encoding, and a Core service wrapper that exposes
the same operations from the top-level `github.com/dappcore/go-lns/go` package.

The package is intentionally self-contained. The JavaScript and HSD source trees
under `docs/js-*` are reference material for porting and compatibility checks;
runtime code should use the Go primitives in `pkg/primitives`, `pkg/covenant`,
and `pkg/dns`.

## Quick Start

```go
package main

import (
	core "github.com/dappcore/go"
	"github.com/dappcore/go-lns/go"
)

func main() {
	app := core.New(core.WithService(lns.Register))
	_ = app
	hash, err := lns.ResolveString("example.lthn")
	if err != nil {
		panic(err)
	}
	core.Println(hash)
}
```

## Packages

`lns.go` mirrors the public API for callers that want one import path. The
subpackages hold the implementation: `pkg/covenant` owns name rules and catalog
verification, `pkg/dns` owns DNS resource projection and live bridge helpers,
`pkg/primitives` owns serializable chain data structures, and `pkg/wallet` /
`pkg/chain` provide lightweight scaffolding for action handlers.

## Validation

This repository is not included in the parent `go.work`, so run validation with
workspace mode disabled:

```sh
GOWORK=off go test ./...
```

For compliance work, also run the v0.9.0 audit script supplied by the core/go
workspace.
