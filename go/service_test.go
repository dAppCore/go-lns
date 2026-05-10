// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	"testing"

	core "dappco.re/go"
)

// TestNewServiceFactory_DefaultOptions — happy path for the canonical
// factory shape added in service.go (Mantis #1336). Wraps lns.Register
// internally; v1 Options has no fields.
func TestNewServiceFactory_DefaultOptions(t *testing.T) {
	c := core.New(core.WithService(NewServiceFactory(Options{})))
	r := c.Service("lns")
	if !r.OK {
		t.Fatalf("lns service not registered via NewServiceFactory, got %#v", r.Value)
	}
}
