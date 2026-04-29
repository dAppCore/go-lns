// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestHelpers_jsonraw_MarshalJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestHelpers_jsonraw_MarshalJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestHelpers_jsonraw_MarshalJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).MarshalJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestHelpers_jsonraw_UnmarshalJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestHelpers_jsonraw_UnmarshalJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestHelpers_jsonraw_UnmarshalJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(jsonraw).UnmarshalJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
