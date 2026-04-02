// SPDX-License-Identifier: EUPL-1.2

package primitives

import "testing"

func TestOutpointIsNull(t *testing.T) {
	zero := Outpoint{Index: ^uint32(0)}
	if !zero.IsNull() {
		t.Fatal("null outpoint should be recognised")
	}

	nonZero := Outpoint{Index: 0}
	if nonZero.IsNull() {
		t.Fatal("non-zero outpoint should not be null")
	}
}

func TestInputIsCoinbase(t *testing.T) {
	coinbase := Input{Prevout: Outpoint{Index: ^uint32(0)}}
	if !coinbase.IsCoinbase() {
		t.Fatal("input with a null prevout should be coinbase")
	}

	regular := Input{Prevout: Outpoint{Index: 0}}
	if regular.IsCoinbase() {
		t.Fatal("input with a real prevout should not be coinbase")
	}
}

func TestInputIsFinal(t *testing.T) {
	final := Input{Sequence: ^uint32(0)}
	if !final.IsFinal() {
		t.Fatal("input with max sequence should be final")
	}

	nonFinal := Input{Sequence: 0}
	if nonFinal.IsFinal() {
		t.Fatal("input with a lower sequence should not be final")
	}
}
