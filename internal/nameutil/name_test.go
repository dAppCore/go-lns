// SPDX-License-Identifier: EUPL-1.2

package nameutil

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestName_CatalogLabel_Good(t *core.T) {
	fn := ax7ValueOf(CatalogLabel)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_CatalogLabel_Bad(t *core.T) {
	fn := ax7ValueOf(CatalogLabel)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_CatalogLabel_Ugly(t *core.T) {
	fn := ax7ValueOf(CatalogLabel)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_Canonicalize_Good(t *core.T) {
	fn := ax7ValueOf(Canonicalize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_Canonicalize_Bad(t *core.T) {
	fn := ax7ValueOf(Canonicalize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestName_Canonicalize_Ugly(t *core.T) {
	fn := ax7ValueOf(Canonicalize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
