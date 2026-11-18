// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	core "github.com/dappcore/go"
	dnspkg "github.com/dappcore/go-lns/go/pkg/dns"
)

const resourceTXTJSON = `{"records":[{"type":"TXT","txt":["hello"]}]}`

func TestActionResourceBehaviour_resourceFromAny_Good(t *core.T) {
	// *Resource passthrough.
	live := NewResource()
	got, err := resourceFromAny(live)
	core.AssertNoError(t, err)
	core.AssertSame(t, live, got)

	// Resource (value) form is copied to a pointer.
	value := *NewResource()
	got, err = resourceFromAny(value)
	core.AssertNoError(t, err)
	core.AssertNotNil(t, got)

	// JSON string form.
	got, err = resourceFromAny(resourceTXTJSON)
	core.AssertNoError(t, err)
	core.AssertLen(t, got.Records, 1)

	// map[string]any form.
	got, err = resourceFromAny(map[string]any{
		"records": []any{map[string]any{"type": "TXT", "txt": []any{"hi"}}},
	})
	core.AssertNoError(t, err)
	core.AssertLen(t, got.Records, 1)

	// Binary form via Encode round-trip.
	encoded, err := live.Encode()
	core.AssertNoError(t, err)
	got, err = resourceFromAny(encoded)
	core.AssertNoError(t, err)
	core.AssertNotNil(t, got)
}

func TestActionResourceBehaviour_resourceFromAny_Bad(t *core.T) {
	// Nil *Resource is rejected.
	_, err := resourceFromAny((*Resource)(nil))
	core.AssertError(t, err)

	// Invalid JSON string is rejected.
	_, err = resourceFromAny("{not json")
	core.AssertError(t, err)

	// Unsupported type is rejected.
	_, err = resourceFromAny(42)
	core.AssertError(t, err)
}

func TestActionResourceBehaviour_resourceFromAny_Ugly(t *core.T) {
	// ResourceJSON form is accepted directly.
	var jsonView ResourceJSON
	r := core.JSONUnmarshalString(resourceTXTJSON, &jsonView)
	core.AssertTrue(t, r.OK)
	got, err := resourceFromAny(jsonView)
	core.AssertNoError(t, err)
	core.AssertLen(t, got.Records, 1)
}

func TestActionResourceBehaviour_resourceFromJSONBytes(t *core.T) {
	got, err := resourceFromJSONBytes([]byte(resourceTXTJSON))
	core.AssertNoError(t, err)
	core.AssertLen(t, got.Records, 1)

	_, err = resourceFromJSONBytes([]byte("{broken"))
	core.AssertError(t, err)
}

func TestActionResourceBehaviour_actionResourcesValue_Good(t *core.T) {
	// map[string]*Resource passthrough.
	live := map[string]*Resource{"charon.lthn.": NewResource()}
	got, err := actionResourcesValue(live)
	core.AssertNoError(t, err)
	core.AssertLen(t, got, 1)

	// map[string]Resource is FQDN-keyed (trailing dot) and copied.
	byValue := map[string]dnspkg.Resource{"charon": *NewResource()}
	got, err = actionResourcesValue(byValue)
	core.AssertNoError(t, err)
	_, ok := got["charon."]
	core.AssertTrue(t, ok)

	// map[string]any resolves each payload via resourceFromAny.
	byAny := map[string]any{"trade": resourceTXTJSON}
	got, err = actionResourcesValue(byAny)
	core.AssertNoError(t, err)
	core.AssertLen(t, got["trade."].Records, 1)
}

func TestActionResourceBehaviour_actionResourcesValue_Bad(t *core.T) {
	// Unsupported top-level type is rejected.
	_, err := actionResourcesValue(42)
	core.AssertError(t, err)

	// An invalid nested payload propagates the error.
	_, err = actionResourcesValue(map[string]any{"bad": "{broken"})
	core.AssertError(t, err)
}
