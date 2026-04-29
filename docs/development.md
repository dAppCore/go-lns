# Development

Use `GOWORK=off` for Go commands in this repository. The parent workspace does
not include this module, and plain `go test ./...` can resolve the wrong module
graph from the surrounding checkout.

Recommended local loop:

```sh
GOWORK=off go test ./pkg/covenant ./pkg/dns ./pkg/primitives
GOWORK=off go test ./...
```

Before handing off a compliance branch, run the full gate:

```sh
GOWORK=off go mod tidy
GOWORK=off go vet ./...
GOWORK=off go test -count=1 ./...
gofmt -l .
bash /Users/snider/Code/core/go/tests/cli/v090-upgrade/audit.sh .
```

Code should follow the `core/go` wrapper pattern. Use `core.JSONMarshal`,
`core.JSONUnmarshal`, `core.NewBuffer`, `core.Join`, `core.Trim`, `core.E`, and
the Core assertion helpers instead of direct standard-library convenience
imports that are banned by the audit.

When adding public behavior, keep tests and examples next to the source file.
The repository uses file-aware coverage names so readers can open `foo.go`,
`foo_test.go`, and `foo_example_test.go` together and understand the contract.
