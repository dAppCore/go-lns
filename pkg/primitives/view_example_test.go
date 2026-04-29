// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "dappco.re/go"

func ExampleNewNameView() {
	core.Println(coverageExample(NewNameView()))
	// Output: 1
}

func ExampleNameView_GetNameStateSync() {
	core.Println(coverageExample(new(NameView).GetNameStateSync(coverageValue[NameStateGetter]("good"), [32]byte{1})))
	// Output: 2
}

func ExampleNameView_GetNameState() {
	core.Println(coverageExample(new(NameView).GetNameState(coverageValue[NameStateGetter]("good"), [32]byte{1})))
	// Output: 2
}

func ExampleNameView_ToNameUndo() {
	core.Println(coverageExample(new(NameView).ToNameUndo()))
	// Output: 1
}

func ExampleNameUndo_FromView() {
	core.Println(coverageExample(new(NameUndo).FromView(nil)))
	// Output: 1
}

func ExampleNameUndo_GetSize() {
	core.Println(coverageExample(new(NameUndo).GetSize()))
	// Output: 1
}

func ExampleNameUndo_MarshalBinary() {
	core.Println(coverageExample(new(NameUndo).MarshalBinary()))
	// Output: 2
}

func ExampleNameUndo_UnmarshalBinary() {
	core.Println(coverageExample(new(NameUndo).UnmarshalBinary([]byte("lethean.lthn"))))
	// Output: 1
}
