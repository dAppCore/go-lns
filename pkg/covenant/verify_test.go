// SPDX-License-Identifier: EUPL-1.2

package covenant

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestVerify_VerifyCovenants_Good(t *core.T) {
	fn := ax7ValueOf(VerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestVerify_VerifyCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(VerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestVerify_VerifyCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(VerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestVerify_GetVerifyCovenants_Good(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenants)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestVerify_GetVerifyCovenants_Bad(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenants)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestVerify_GetVerifyCovenants_Ugly(t *core.T) {
	fn := ax7ValueOf(GetVerifyCovenants)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
