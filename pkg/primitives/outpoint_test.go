// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

// BEGIN v090 file-aware generated coverage

func TestOutpoint_NewOutpoint_Good(t *core.T) {
	fn := ax7ValueOf(NewOutpoint)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_NewOutpoint_Bad(t *core.T) {
	fn := ax7ValueOf(NewOutpoint)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_NewOutpoint_Ugly(t *core.T) {
	fn := ax7ValueOf(NewOutpoint)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Inject_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Inject)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Inject_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Inject)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Inject_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Inject)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Clone_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Clone)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Clone_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Clone)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Clone_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Clone)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Equals_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Equals)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Equals_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Equals)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Equals_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Equals)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Compare_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Compare)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Compare_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Compare)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_Compare_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).Compare)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_GetSize_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).GetSize)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_GetSize_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).GetSize)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_GetSize_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).GetSize)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_MarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).MarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_MarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).MarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_MarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).MarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_UnmarshalBinary_Good(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).UnmarshalBinary)
	args := ax7Args(fn.Type(), "good")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_UnmarshalBinary_Bad(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).UnmarshalBinary)
	args := ax7Args(fn.Type(), "bad")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

func TestOutpoint_Outpoint_UnmarshalBinary_Ugly(t *core.T) {
	fn := ax7ValueOf(new(Outpoint).UnmarshalBinary)
	args := ax7Args(fn.Type(), "ugly")
	outCount := -1
	core.AssertNotPanics(t, func() {
		outCount = len(ax7Invoke(fn, args))
	})
	core.AssertEqual(t, fn.Type().NumOut(), outCount)
}

// END v090 file-aware generated coverage
