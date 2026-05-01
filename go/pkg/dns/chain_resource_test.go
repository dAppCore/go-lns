// SPDX-License-Identifier: EUPL-1.2

package dns

import (
	"testing"

	chainpkg "dappco.re/go/lns/pkg/chain"
)

func TestResourceFromChainResource(t *testing.T) {
	resource, err := resourceFromChainResource(&chainpkg.NameResource{
		Records: []chainpkg.Record{
			{Type: "GLUE4", NS: "charon.lthn.", Address: "192.0.2.10"},
			{Type: "TXT", TXT: []string{"hello"}},
			{Type: "NS", NS: "ns1.charon.lthn."},
		},
	})
	if err != nil {
		t.Fatalf("resourceFromChainResource returned error: %v", err)
	}

	if len(resource.Records) != 3 {
		t.Fatalf("resourceFromChainResource returned %d records, want 3", len(resource.Records))
	}
}
