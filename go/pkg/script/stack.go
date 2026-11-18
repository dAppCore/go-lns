// SPDX-License-Identifier: EUPL-1.2

package script

import core "github.com/dappcore/go"

// Stack is the script execution stack: a LIFO of byte-string items. It backs
// both the main and alt stacks during execution and is the input form a witness
// is converted to before verification. It mirrors HSD's Stack
// (docs/js-script via lib/script/stack.js) with negative indices counting from
// the top.
//
//	st := script.NewStack()
//	st.Push([]byte{0x01})
//	top, _ := st.Top(-1)
type Stack struct {
	items [][]byte
}

// NewStack returns an empty stack.
func NewStack() *Stack {
	return &Stack{}
}

// Len returns the number of items on the stack.
func (s *Stack) Len() int {
	return len(s.items)
}

// Items returns the underlying items (top of stack is the last element). The
// returned slice aliases the stack; callers must not mutate it.
func (s *Stack) Items() [][]byte {
	return s.items
}

// abs converts a possibly-negative index (counting from the top) into an
// absolute index, returning ok=false if it is out of range.
func (s *Stack) abs(index int) (int, bool) {
	if index < 0 {
		index = len(s.items) + index
	}
	if index < 0 || index >= len(s.items) {
		return 0, false
	}
	return index, true
}

// Push appends an item to the top of the stack.
func (s *Stack) Push(item []byte) {
	s.items = append(s.items, item)
}

// PushInt pushes the minimal encoding of a script number.
func (s *Stack) PushInt(n ScriptNum) {
	s.Push(n.Encode())
}

// PushBool pushes the canonical boolean encoding (empty for false, {1} for
// true).
func (s *Stack) PushBool(b bool) {
	if b {
		s.Push([]byte{1})
	} else {
		s.Push([]byte{})
	}
}

// Pop removes and returns the top item, or an error if the stack is empty.
func (s *Stack) Pop() ([]byte, error) {
	if len(s.items) == 0 {
		return nil, core.E("script.Stack.Pop", "stack underflow", nil)
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

// Top returns the item at the given (possibly negative) index without removing
// it.
func (s *Stack) Top(index int) ([]byte, error) {
	i, ok := s.abs(index)
	if !ok {
		return nil, core.E("script.Stack.Top", "invalid stack operation", nil)
	}
	return s.items[i], nil
}

// Set replaces the item at the given (possibly negative) index.
func (s *Stack) Set(index int, item []byte) error {
	i, ok := s.abs(index)
	if !ok {
		return core.E("script.Stack.Set", "invalid stack operation", nil)
	}
	s.items[i] = item
	return nil
}

// Erase removes the item at the given (possibly negative) index.
func (s *Stack) Erase(index int) error {
	i, ok := s.abs(index)
	if !ok {
		return core.E("script.Stack.Erase", "invalid stack operation", nil)
	}
	s.items = append(s.items[:i], s.items[i+1:]...)
	return nil
}

// Insert inserts an item before the given (possibly negative) absolute index,
// shifting later items up. It matches HSD's Stack.insert (splice semantics): a
// negative index is offset by the length, and an index equal to Len appends to
// the top. OP_TUCK relies on Insert(-2, top) placing a copy below the second
// item.
func (s *Stack) Insert(index int, item []byte) error {
	if index < 0 {
		index = len(s.items) + index
	}
	if index < 0 || index > len(s.items) {
		return core.E("script.Stack.Insert", "invalid stack operation", nil)
	}
	s.items = append(s.items, nil)
	copy(s.items[index+1:], s.items[index:])
	s.items[index] = item
	return nil
}

// Bool returns the truthiness of the item at the given index.
func (s *Stack) Bool(index int) (bool, error) {
	item, err := s.Top(index)
	if err != nil {
		return false, err
	}
	return castToBool(item), nil
}

// Num decodes the item at the given index as a script number.
func (s *Stack) Num(index int, minimal bool, limit int) (ScriptNum, error) {
	item, err := s.Top(index)
	if err != nil {
		return 0, err
	}
	return DecodeScriptNum(item, minimal, limit)
}

// Int decodes the item at the given index as a script number and returns it as
// an int.
func (s *Stack) Int(index int, minimal bool, limit int) (int, error) {
	num, err := s.Num(index, minimal, limit)
	if err != nil {
		return 0, err
	}
	return num.Int(), nil
}
