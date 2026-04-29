// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestNamestateJson_NameState_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_Format_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).Format)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_Format_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).Format)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_Format_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).Format)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetFormat_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetFormat)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetFormat_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetFormat)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetFormat_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetFormat)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetFromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetFromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetFromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetFromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateJson_NameState_GetFromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetFromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
