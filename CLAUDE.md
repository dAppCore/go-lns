# Claude Conventions

This repository follows the `core/go` v0.9 conventions. Use `github.com/dappcore/go`
wrappers for formatted output, JSON, byte buffers, path helpers, string helpers,
and test assertions. Do not add direct imports of `fmt`, `errors`, `strings`,
`path`, `os`, `log`, `encoding/json`, or `bytes` in production or tests.

Keep the top-level `lns` package as a mirror of the package-local APIs. When a
public DNS, covenant, primitive, chain, or wallet helper is added below `pkg/`,
check whether callers should also reach it through `lns.go` and preserve the
existing `Get*` alias pattern where the file already uses one.

Tests are file-aware. A public symbol in `foo.go` needs its
`TestFoo_<Symbol>_{Good,Bad,Ugly}` triplet in `foo_test.go`, and examples belong
in `foo_example_test.go`. Existing tests may use richer domain fixtures, but
new compliance coverage should stay beside the source file it exercises.

The JavaScript reference trees in `docs/js-*` are read-only porting sources.
They should inform behavior, not become runtime dependencies. Prefer the RFCs
in `docs/RFC.md` and `docs/RFC.lns.md` when resolving implementation questions.
