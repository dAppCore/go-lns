// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestOutpointJson_Outpoint_GetJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).GetJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpointJson_Outpoint_GetJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).GetJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpointJson_Outpoint_GetJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).GetJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpointJson_Outpoint_FromJSON_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).FromJSON)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpointJson_Outpoint_FromJSON_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).FromJSON)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpointJson_Outpoint_FromJSON_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).FromJSON)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
