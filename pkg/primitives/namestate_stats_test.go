// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestNamestateStats_NameState_ToStats_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).ToStats)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateStats_NameState_ToStats_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).ToStats)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateStats_NameState_ToStats_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).ToStats)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateStats_NameState_GetToStats_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetToStats)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateStats_NameState_GetToStats_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetToStats)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateStats_NameState_GetToStats_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetToStats)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
