// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

// runScript executes ops against a fresh stack with the given flags and no tx
// context, returning the resulting stack and any error. It is the harness for
// non-signature opcode tests.
func runScript(t *testing.T, flags int, ops ...Opcode) (*Stack, error) {
	t.Helper()
	s, err := FromOps(ops)
	if err != nil {
		t.Fatalf("FromOps: %v", err)
	}
	stack := NewStack()
	return stack, s.Execute(stack, flags, nil)
}

func TestEngine_Arithmetic_Good(t *testing.T) {
	// 2 3 OP_ADD -> 5
	stack, err := runScript(t, VerifyNone, OpFromOp(OP_2), OpFromOp(OP_3), OpFromOp(OP_ADD))
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if stack.Len() != 1 {
		t.Fatalf("stack len = %d", stack.Len())
	}
	n, _ := stack.Num(-1, false, scriptNumMaxBytes)
	if n != 5 {
		t.Fatalf("2+3 = %d, want 5", n)
	}
}

func TestEngine_Comparison_Good(t *testing.T) {
	// 5 5 OP_NUMEQUAL -> true
	stack, err := runScript(t, VerifyNone, OpFromOp(OP_5), OpFromOp(OP_5), OpFromOp(OP_NUMEQUAL))
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	b, _ := stack.Bool(-1)
	if !b {
		t.Fatal("5 == 5 should be true")
	}

	// 3 7 OP_WITHIN(3,7) of value 5 -> 5 3 7 within -> true
	stack, err = runScript(t, VerifyNone,
		OpFromOp(OP_5), OpFromOp(OP_3), OpFromOp(OP_7), OpFromOp(OP_WITHIN))
	if err != nil {
		t.Fatalf("within execute: %v", err)
	}
	b, _ = stack.Bool(-1)
	if !b {
		t.Fatal("5 within [3,7) should be true")
	}
}

func TestEngine_StackOps_Good(t *testing.T) {
	// 1 2 OP_SWAP -> top is 1
	stack, err := runScript(t, VerifyNone, OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_SWAP))
	if err != nil {
		t.Fatalf("swap: %v", err)
	}
	top, _ := stack.Top(-1)
	if !bytes.Equal(top, []byte{0x01}) {
		t.Fatalf("after swap top = %x, want 01", top)
	}

	// 1 2 3 OP_ROT -> 2 3 1
	stack, err = runScript(t, VerifyNone,
		OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_3), OpFromOp(OP_ROT))
	if err != nil {
		t.Fatalf("rot: %v", err)
	}
	b, _ := stack.Top(-3)
	m, _ := stack.Top(-2)
	tp, _ := stack.Top(-1)
	if !bytes.Equal(b, []byte{0x02}) || !bytes.Equal(m, []byte{0x03}) || !bytes.Equal(tp, []byte{0x01}) {
		t.Fatalf("rot result = %x %x %x, want 02 03 01", b, m, tp)
	}
}

func TestEngine_HashOps_Good(t *testing.T) {
	data := []byte("lethean")

	cases := []struct {
		op   int
		want []byte
	}{
		{OP_SHA256, sha256Sum(data)},
		{OP_BLAKE160, blake160(data)},
		{OP_BLAKE256, blake256(data)},
		{OP_SHA3, sha3Sum(data)},
		{OP_KECCAK, keccakSum(data)},
	}
	for _, c := range cases {
		stack, err := runScript(t, VerifyNone, OpFromPush(data), OpFromOp(c.op))
		if err != nil {
			t.Fatalf("%s: %v", OpcodeName(c.op), err)
		}
		got, _ := stack.Top(-1)
		if !bytes.Equal(got, c.want) {
			t.Fatalf("%s(%q) = %x\nwant %x", OpcodeName(c.op), data, got, c.want)
		}
	}
}

func TestEngine_Conditional_Good(t *testing.T) {
	// OP_1 OP_IF OP_2 OP_ELSE OP_3 OP_ENDIF -> 2
	stack, err := runScript(t, VerifyNone,
		OpFromOp(OP_1), OpFromOp(OP_IF), OpFromOp(OP_2),
		OpFromOp(OP_ELSE), OpFromOp(OP_3), OpFromOp(OP_ENDIF))
	if err != nil {
		t.Fatalf("if/else: %v", err)
	}
	n, _ := stack.Num(-1, false, scriptNumMaxBytes)
	if n != 2 {
		t.Fatalf("IF branch = %d, want 2", n)
	}

	// OP_0 OP_IF OP_2 OP_ELSE OP_3 OP_ENDIF -> 3
	stack, err = runScript(t, VerifyNone,
		OpFromOp(OP_0), OpFromOp(OP_IF), OpFromOp(OP_2),
		OpFromOp(OP_ELSE), OpFromOp(OP_3), OpFromOp(OP_ENDIF))
	if err != nil {
		t.Fatalf("else branch: %v", err)
	}
	n, _ = stack.Num(-1, false, scriptNumMaxBytes)
	if n != 3 {
		t.Fatalf("ELSE branch = %d, want 3", n)
	}
}

func TestEngine_Underflow_Bad(t *testing.T) {
	// OP_ADD with an empty stack must fail with INVALID_STACK_OPERATION.
	_, err := runScript(t, VerifyNone, OpFromOp(OP_ADD))
	assertScriptError(t, err, errInvalidStackOperation)

	// OP_DUP with empty stack.
	_, err = runScript(t, VerifyNone, OpFromOp(OP_DUP))
	assertScriptError(t, err, errInvalidStackOperation)
}

func TestEngine_Disabled_Bad(t *testing.T) {
	// OP_MUL is permanently disabled and must fail even with operands present.
	_, err := runScript(t, VerifyNone, OpFromOp(OP_2), OpFromOp(OP_3), OpFromOp(OP_MUL))
	assertScriptError(t, err, errDisabledOpcode)
}

func TestEngine_Return_Bad(t *testing.T) {
	_, err := runScript(t, VerifyNone, OpFromOp(OP_1), OpFromOp(OP_RETURN))
	assertScriptError(t, err, errOpReturn)
}

func TestEngine_Unbalanced_Bad(t *testing.T) {
	// OP_IF without OP_ENDIF leaves the conditional stack non-empty.
	_, err := runScript(t, VerifyNone, OpFromOp(OP_1), OpFromOp(OP_IF), OpFromOp(OP_2))
	assertScriptError(t, err, errUnbalancedConditional)

	// OP_ENDIF with no matching OP_IF.
	_, err = runScript(t, VerifyNone, OpFromOp(OP_ENDIF))
	assertScriptError(t, err, errUnbalancedConditional)
}

func TestEngine_OpCount_Ugly(t *testing.T) {
	// MaxScriptOps+1 executed non-push opcodes must trip OP_COUNT.
	ops := make([]Opcode, 0, MaxScriptOps+2)
	for i := 0; i < MaxScriptOps+1; i++ {
		ops = append(ops, OpFromOp(OP_NOP))
	}
	_, err := runScript(t, VerifyNone, ops...)
	assertScriptError(t, err, errOpCount)
}

func TestEngine_MinimalData_Ugly(t *testing.T) {
	// A non-minimal push under VerifyMinimalData fails.
	s := &Script{}
	s.code = []Opcode{{Value: 1, Data: []byte{0x01}}} // should be OP_1
	_ = s.compile()
	stack := NewStack()
	err := s.Execute(stack, VerifyMinimalData, nil)
	assertScriptError(t, err, errMinimalData)
}

func TestEngine_Verify_Bad(t *testing.T) {
	// OP_0 OP_VERIFY fails; OP_1 OP_VERIFY passes.
	_, err := runScript(t, VerifyNone, OpFromOp(OP_0), OpFromOp(OP_VERIFY))
	assertScriptError(t, err, errVerify)

	stack, err := runScript(t, VerifyNone, OpFromOp(OP_1), OpFromOp(OP_VERIFY))
	if err != nil {
		t.Fatalf("OP_1 OP_VERIFY: %v", err)
	}
	if stack.Len() != 0 {
		t.Fatalf("VERIFY should consume the value, stack len = %d", stack.Len())
	}
}

func TestEngine_Equal_Good(t *testing.T) {
	// push 0xdead, push 0xdead, OP_EQUAL -> true
	stack, err := runScript(t, VerifyNone,
		OpFromPush([]byte{0xde, 0xad}), OpFromPush([]byte{0xde, 0xad}), OpFromOp(OP_EQUAL))
	if err != nil {
		t.Fatalf("equal: %v", err)
	}
	b, _ := stack.Bool(-1)
	if !b {
		t.Fatal("equal bytes should be true")
	}
}

func TestEngine_AltStack_Good(t *testing.T) {
	// OP_5 OP_TOALTSTACK OP_FROMALTSTACK -> 5 back on main
	stack, err := runScript(t, VerifyNone,
		OpFromOp(OP_5), OpFromOp(OP_TOALTSTACK), OpFromOp(OP_FROMALTSTACK))
	if err != nil {
		t.Fatalf("altstack: %v", err)
	}
	n, _ := stack.Num(-1, false, scriptNumMaxBytes)
	if n != 5 {
		t.Fatalf("altstack round-trip = %d, want 5", n)
	}

	// OP_FROMALTSTACK on empty alt stack fails.
	_, err = runScript(t, VerifyNone, OpFromOp(OP_FROMALTSTACK))
	assertScriptError(t, err, errInvalidAltstackOperation)
}

// assertScriptError fails the test unless err is a *ScriptError with the given
// code.
func assertScriptError(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected ScriptError %q, got nil", code)
	}
	serr, ok := err.(*ScriptError)
	if !ok {
		t.Fatalf("expected *ScriptError, got %T: %v", err, err)
	}
	if serr.Code != code {
		t.Fatalf("ScriptError code = %q, want %q (%v)", serr.Code, code, serr)
	}
}

/*
 * Hash helpers mirroring the opcode digests for assertions.
 */

func sha256Sum(b []byte) []byte { s := sha256.Sum256(b); return s[:] }
func sha3Sum(b []byte) []byte   { s := sha3.Sum256(b); return s[:] }

func keccakSum(b []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	_, _ = h.Write(b)
	return h.Sum(nil)
}

func blake160(b []byte) []byte {
	h, _ := blake2b.New(20, nil)
	_, _ = h.Write(b)
	return h.Sum(nil)
}

func blake256(b []byte) []byte { s := blake2b.Sum256(b); return s[:] }
