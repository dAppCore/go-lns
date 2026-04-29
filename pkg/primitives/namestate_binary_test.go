// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestNamestateBinary_NameState_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_MarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).MarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_MarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).MarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_MarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).MarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_UnmarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(NameState).UnmarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_UnmarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(NameState).UnmarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestNamestateBinary_NameState_UnmarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(NameState).UnmarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
