// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"bytes"
	"testing"
)

func TestStack_PushPop_Good(t *testing.T) {
	st := NewStack()
	st.Push([]byte{0x01})
	st.Push([]byte{0x02})
	if st.Len() != 2 {
		t.Fatalf("len = %d", st.Len())
	}
	top, err := st.Pop()
	if err != nil {
		t.Fatalf("pop: %v", err)
	}
	if !bytes.Equal(top, []byte{0x02}) {
		t.Fatalf("pop = %x, want 02", top)
	}
	if got := st.Items(); len(got) != 1 {
		t.Fatalf("Items len = %d, want 1", len(got))
	}
}

func TestStack_Pop_Bad(t *testing.T) {
	st := NewStack()
	if _, err := st.Pop(); err == nil {
		t.Fatal("expected underflow error")
	}
	if _, err := st.Top(-1); err == nil {
		t.Fatal("expected invalid-operation error on empty Top")
	}
}

func TestStack_InsertErase_Good(t *testing.T) {
	st := NewStack()
	st.Push([]byte{0x01})
	st.Push([]byte{0x03})
	// Insert 0x02 so the stack becomes 01 02 03.
	if err := st.Insert(-1, []byte{0x02}); err != nil {
		t.Fatalf("insert: %v", err)
	}
	mid, _ := st.Top(-2)
	if !bytes.Equal(mid, []byte{0x02}) {
		t.Fatalf("after insert middle = %x, want 02", mid)
	}
	// Erase the bottom item (01).
	if err := st.Erase(-3); err != nil {
		t.Fatalf("erase: %v", err)
	}
	if st.Len() != 2 {
		t.Fatalf("after erase len = %d, want 2", st.Len())
	}
}

func TestStack_InsertErase_Bad(t *testing.T) {
	st := NewStack()
	if err := st.Erase(-1); err == nil {
		t.Fatal("expected error erasing from empty stack")
	}
	if err := st.Insert(5, []byte{0x01}); err == nil {
		t.Fatal("expected error inserting out of range")
	}
	if err := st.Set(0, []byte{0x01}); err == nil {
		t.Fatal("expected error setting on empty stack")
	}
}

func TestStack_PushBoolInt_Good(t *testing.T) {
	st := NewStack()
	st.PushBool(true)
	b, _ := st.Bool(-1)
	if !b {
		t.Fatal("PushBool(true) should be truthy")
	}
	st.PushInt(-3)
	n, _ := st.Num(-1, false, scriptNumMaxBytes)
	if n != -3 {
		t.Fatalf("PushInt(-3) = %d", n)
	}
}
