# Architecture

`dappco.re/go/lns` is a native Go port of the Lethean Name System. The top-level
package is intentionally broad: `lns.go` re-exports and mirrors the public
surface of the lower packages so application code can import one package when it
does not need tighter ownership boundaries.

The implementation packages remain the source of behavior:

- `internal/nameutil` canonicalizes label input and keeps suffix handling in one
  place.
- `pkg/covenant` owns covenant types, rollout rules, name verification, reserved
  catalog loading, and locked catalog loading.
- `pkg/dns` owns LNS-to-DNS projection, binary resource encoding, JSON resource
  forms, NSEC helpers, and the live DNS bridge.
- `pkg/primitives` owns chain-shaped data structures such as claims, outpoints,
  name states, deltas, covenants, and views.
- `pkg/chain` and `pkg/wallet` are lightweight support packages used by Core
  action handlers and live bridge workflows.

The service layer in `lns.go`, `action_handlers.go`, and `action_live.go`
connects these packages to `dappco.re/go` Core actions. It should remain thin:
construct inputs, call the package-local implementation, and return typed Core
results.

Catalog data is loaded from the checked-in JS reference JSON under
`docs/js-covenants`. That data is an input fixture, not a dependency on the
JavaScript implementation.
