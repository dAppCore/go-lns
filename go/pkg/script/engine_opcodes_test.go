// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"bytes"
	"crypto/sha1"
	"testing"

	"golang.org/x/crypto/ripemd160"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// evalStack runs ops and returns the final top-of-stack item, failing on error.
func evalTop(t *testing.T, ops ...Opcode) []byte {
	t.Helper()
	stack, err := runScript(t, VerifyNone, ops...)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	top, err := stack.Top(-1)
	if err != nil {
		t.Fatalf("empty stack")
	}
	return top
}

// evalTopNum runs ops and returns the final top as a script number.
func evalTopNum(t *testing.T, ops ...Opcode) ScriptNum {
	t.Helper()
	stack, err := runScript(t, VerifyNone, ops...)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	n, err := stack.Num(-1, false, scriptNumMaxBytes)
	if err != nil {
		t.Fatalf("top num: %v", err)
	}
	return n
}

func TestEngine_Opcodes_Numeric_Good(t *testing.T) {
	cases := []struct {
		name string
		ops  []Opcode
		want ScriptNum
	}{
		{"1ADD", []Opcode{OpFromOp(OP_5), OpFromOp(OP_1ADD)}, 6},
		{"1SUB", []Opcode{OpFromOp(OP_5), OpFromOp(OP_1SUB)}, 4},
		{"NEGATE", []Opcode{OpFromOp(OP_5), OpFromOp(OP_NEGATE)}, -5},
		{"ABS", []Opcode{OpFromOp(OP_1NEGATE), OpFromOp(OP_ABS)}, 1},
		{"NOT_true", []Opcode{OpFromOp(OP_0), OpFromOp(OP_NOT)}, 1},
		{"NOT_false", []Opcode{OpFromOp(OP_5), OpFromOp(OP_NOT)}, 0},
		{"0NOTEQUAL", []Opcode{OpFromOp(OP_5), OpFromOp(OP_0NOTEQUAL)}, 1},
		{"SUB", []Opcode{OpFromOp(OP_9), OpFromOp(OP_4), OpFromOp(OP_SUB)}, 5},
		{"MIN", []Opcode{OpFromOp(OP_9), OpFromOp(OP_4), OpFromOp(OP_MIN)}, 4},
		{"MAX", []Opcode{OpFromOp(OP_9), OpFromOp(OP_4), OpFromOp(OP_MAX)}, 9},
		{"BOOLAND", []Opcode{OpFromOp(OP_1), OpFromOp(OP_1), OpFromOp(OP_BOOLAND)}, 1},
		{"BOOLOR", []Opcode{OpFromOp(OP_0), OpFromOp(OP_1), OpFromOp(OP_BOOLOR)}, 1},
		{"NUMNOTEQUAL", []Opcode{OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_NUMNOTEQUAL)}, 1},
		{"LESSTHAN", []Opcode{OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_LESSTHAN)}, 1},
		{"GREATERTHAN", []Opcode{OpFromOp(OP_2), OpFromOp(OP_1), OpFromOp(OP_GREATERTHAN)}, 1},
		{"LESSTHANOREQUAL", []Opcode{OpFromOp(OP_2), OpFromOp(OP_2), OpFromOp(OP_LESSTHANOREQUAL)}, 1},
		{"GREATERTHANOREQUAL", []Opcode{OpFromOp(OP_2), OpFromOp(OP_2), OpFromOp(OP_GREATERTHANOREQUAL)}, 1},
	}
	for _, c := range cases {
		if got := evalTopNum(t, c.ops...); got != c.want {
			t.Fatalf("%s = %d, want %d", c.name, got, c.want)
		}
	}
}

func TestEngine_Opcodes_Stack_Good(t *testing.T) {
	// OP_DEPTH after pushing three items.
	if got := evalTopNum(t, OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_3), OpFromOp(OP_DEPTH)); got != 3 {
		t.Fatalf("DEPTH = %d, want 3", got)
	}
	// OP_DUP duplicates the top.
	if got := evalTop(t, OpFromOp(OP_7), OpFromOp(OP_DUP)); !bytes.Equal(got, []byte{0x07}) {
		t.Fatalf("DUP top = %x", got)
	}
	// OP_OVER copies the second item to the top.
	if got := evalTop(t, OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_OVER)); !bytes.Equal(got, []byte{0x01}) {
		t.Fatalf("OVER top = %x, want 01", got)
	}
	// OP_NIP removes the second item.
	st, err := runScript(t, VerifyNone, OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_NIP))
	if err != nil {
		t.Fatalf("NIP: %v", err)
	}
	if st.Len() != 1 {
		t.Fatalf("NIP stack len = %d, want 1", st.Len())
	}
	// OP_TUCK inserts a copy of the top below the second.
	st, err = runScript(t, VerifyNone, OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_TUCK))
	if err != nil {
		t.Fatalf("TUCK: %v", err)
	}
	if st.Len() != 3 {
		t.Fatalf("TUCK stack len = %d, want 3", st.Len())
	}
	// OP_IFDUP duplicates only a truthy top.
	st, _ = runScript(t, VerifyNone, OpFromOp(OP_5), OpFromOp(OP_IFDUP))
	if st.Len() != 2 {
		t.Fatalf("IFDUP(true) len = %d, want 2", st.Len())
	}
	st, _ = runScript(t, VerifyNone, OpFromOp(OP_0), OpFromOp(OP_IFDUP))
	if st.Len() != 1 {
		t.Fatalf("IFDUP(false) len = %d, want 1", st.Len())
	}
	// OP_SIZE pushes the length of the top item without consuming it.
	st, _ = runScript(t, VerifyNone, OpFromPush([]byte{0xaa, 0xbb, 0xcc}), OpFromOp(OP_SIZE))
	n, _ := st.Num(-1, false, scriptNumMaxBytes)
	if n != 3 {
		t.Fatalf("SIZE = %d, want 3", n)
	}
}

func TestEngine_Opcodes_MultiStack_Good(t *testing.T) {
	// OP_2DUP, OP_3DUP, OP_2OVER, OP_2SWAP, OP_2ROT exercised via final depth.
	st, err := runScript(t, VerifyNone,
		OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_2DUP)) // -> 1 2 1 2
	if err != nil || st.Len() != 4 {
		t.Fatalf("2DUP len = %d (err %v)", st.Len(), err)
	}
	st, err = runScript(t, VerifyNone,
		OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_3), OpFromOp(OP_3DUP)) // -> 6
	if err != nil || st.Len() != 6 {
		t.Fatalf("3DUP len = %d (err %v)", st.Len(), err)
	}
	st, err = runScript(t, VerifyNone,
		OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_3), OpFromOp(OP_4), OpFromOp(OP_2OVER))
	if err != nil || st.Len() != 6 {
		t.Fatalf("2OVER len = %d (err %v)", st.Len(), err)
	}
	// 2SWAP: 1 2 3 4 -> 3 4 1 2
	st, err = runScript(t, VerifyNone,
		OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_3), OpFromOp(OP_4), OpFromOp(OP_2SWAP))
	if err != nil {
		t.Fatalf("2SWAP: %v", err)
	}
	top, _ := st.Top(-1)
	if !bytes.Equal(top, []byte{0x02}) {
		t.Fatalf("2SWAP top = %x, want 02", top)
	}
	// 2ROT: 1 2 3 4 5 6 -> 3 4 5 6 1 2
	st, err = runScript(t, VerifyNone,
		OpFromOp(OP_1), OpFromOp(OP_2), OpFromOp(OP_3),
		OpFromOp(OP_4), OpFromOp(OP_5), OpFromOp(OP_6), OpFromOp(OP_2ROT))
	if err != nil {
		t.Fatalf("2ROT: %v", err)
	}
	top, _ = st.Top(-1)
	if !bytes.Equal(top, []byte{0x02}) {
		t.Fatalf("2ROT top = %x, want 02", top)
	}
}

func TestEngine_Opcodes_PickRoll_Good(t *testing.T) {
	// 9 8 7 OP_2 OP_PICK -> copies index 2 (the 9) to top.
	if got := evalTop(t,
		OpFromOp(OP_9), OpFromOp(OP_8), OpFromOp(OP_7), OpFromOp(OP_2), OpFromOp(OP_PICK)); !bytes.Equal(got, []byte{0x09}) {
		t.Fatalf("PICK top = %x, want 09", got)
	}
	// 9 8 7 OP_2 OP_ROLL -> moves the 9 to the top and removes it from below.
	st, err := runScript(t, VerifyNone,
		OpFromOp(OP_9), OpFromOp(OP_8), OpFromOp(OP_7), OpFromOp(OP_2), OpFromOp(OP_ROLL))
	if err != nil {
		t.Fatalf("ROLL: %v", err)
	}
	if st.Len() != 3 {
		t.Fatalf("ROLL len = %d, want 3", st.Len())
	}
	top, _ := st.Top(-1)
	if !bytes.Equal(top, []byte{0x09}) {
		t.Fatalf("ROLL top = %x, want 09", top)
	}
}

func TestEngine_Opcodes_LegacyHashes_Good(t *testing.T) {
	data := []byte("handshake")

	// OP_RIPEMD160.
	r := ripemd160.New()
	_, _ = r.Write(data)
	if got := evalTop(t, OpFromPush(data), OpFromOp(OP_RIPEMD160)); !bytes.Equal(got, r.Sum(nil)) {
		t.Fatal("RIPEMD160 mismatch")
	}
	// OP_SHA1.
	s1 := sha1.Sum(data)
	if got := evalTop(t, OpFromPush(data), OpFromOp(OP_SHA1)); !bytes.Equal(got, s1[:]) {
		t.Fatal("SHA1 mismatch")
	}
	// OP_HASH160 = RIPEMD160(SHA256(x)).
	if got := evalTop(t, OpFromPush(data), OpFromOp(OP_HASH160)); len(got) != 20 {
		t.Fatalf("HASH160 length = %d, want 20", len(got))
	}
	// OP_HASH256 = SHA256(SHA256(x)).
	if got := evalTop(t, OpFromPush(data), OpFromOp(OP_HASH256)); len(got) != 32 {
		t.Fatalf("HASH256 length = %d, want 32", len(got))
	}
}

func TestEngine_Opcodes_Type_Good(t *testing.T) {
	// OP_TYPE: empty -> 0, single byte -> 1, multi -> 2.
	if got := evalTopNum(t, OpFromOp(OP_0), OpFromOp(OP_TYPE)); got != 0 {
		t.Fatalf("TYPE(empty) = %d, want 0", got)
	}
	if got := evalTopNum(t, OpFromOp(OP_5), OpFromOp(OP_TYPE)); got != 1 {
		t.Fatalf("TYPE(single) = %d, want 1", got)
	}
	if got := evalTopNum(t, OpFromPush([]byte{0xaa, 0xbb}), OpFromOp(OP_TYPE)); got != 2 {
		t.Fatalf("TYPE(multi) = %d, want 2", got)
	}
}

func TestEngine_Opcodes_NopDiscourage_Bad(t *testing.T) {
	// Reserved NOPs are accepted without the flag, rejected with it.
	if _, err := runScript(t, VerifyNone, OpFromOp(OP_NOP4)); err != nil {
		t.Fatalf("OP_NOP4 should pass without discourage flag: %v", err)
	}
	_, err := runScript(t, VerifyDiscourageUpgradableNops, OpFromOp(OP_NOP4))
	assertScriptError(t, err, errDiscourageUpgradableNops)

	// Timelock opcodes are treated as discouraged NOPs under the flag.
	_, err = runScript(t, VerifyDiscourageUpgradableNops, OpFromOp(OP_CHECKLOCKTIMEVERIFY))
	assertScriptError(t, err, errDiscourageUpgradableNops)
}

func TestEngine_Opcodes_EqualVerify_Bad(t *testing.T) {
	// Unequal items under OP_EQUALVERIFY fail.
	_, err := runScript(t, VerifyNone,
		OpFromPush([]byte{0x01}), OpFromPush([]byte{0x02}), OpFromOp(OP_EQUALVERIFY))
	assertScriptError(t, err, errEqualVerify)
}

func TestEngine_Verify_AddressForm_Good(t *testing.T) {
	// The address-level Verify entry mirrors VerifyInput.
	priv, pub := testKey(t, 21)
	pkh := blake160Of(t, pub)
	addr := primitives.Address{Version: 0, Hash: pkh}
	tx := spendTx(addr, 1000)
	value := uint64(2000)
	redeem := NewScript().FromPubkeyhash(pkh)
	sig, _ := Sign(tx, 0, redeem, value, priv, SighashAll)
	tx.Inputs[0].Witness = [][]byte{sig, pub}

	if err := Verify(tx.Inputs[0].Witness, addr, tx, 0, value, MandatoryVerifyFlags); err != nil {
		t.Fatalf("Verify (address form): %v", err)
	}
}
