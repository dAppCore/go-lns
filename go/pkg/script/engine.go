// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"crypto/sha1"
	"crypto/sha256"

	core "dappco.re/go"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"

	"dappco.re/go/lns/pkg/primitives"
)

// Verify flags. They mirror HSD's common.flags (lib/script/common.js) and
// control the strictness of script execution.
const (
	// VerifyNone disables all optional checks.
	VerifyNone = 0
	// VerifyMinimalData requires pushes and numbers to use their smallest form.
	VerifyMinimalData = 1 << 1
	// VerifyDiscourageUpgradableNops rejects the reserved NOP opcodes.
	VerifyDiscourageUpgradableNops = 1 << 2
	// VerifyDiscourageUpgradableWitnessProgram rejects unknown witness versions.
	VerifyDiscourageUpgradableWitnessProgram = 1 << 3
	// VerifyMinimalIf requires OP_IF/OP_NOTIF operands to be a clean boolean.
	VerifyMinimalIf = 1 << 4
	// VerifyNullFail requires failed CHECKSIG signatures to be empty.
	VerifyNullFail = 1 << 5

	// MandatoryVerifyFlags are the consensus (block-validation) flags.
	MandatoryVerifyFlags = VerifyMinimalData | VerifyMinimalIf | VerifyNullFail

	// StandardVerifyFlags are the policy (mempool) flags.
	StandardVerifyFlags = MandatoryVerifyFlags |
		VerifyDiscourageUpgradableNops |
		VerifyDiscourageUpgradableWitnessProgram
)

// execContext carries the transaction data the signature opcodes need. It is
// passed to Execute so the interpreter can compute signature hashes without the
// script package owning transaction-validation policy.
type execContext struct {
	tx    primitives.Transaction
	index int
	value uint64
}

// Execute runs the script against the supplied stack with the given flags and
// transaction context. It mirrors HSD's Script.execute: a stack-based evaluator
// over the opcode sequence with branch tracking, op/stack limits, and the
// Handshake crypto and signature opcodes. On any rule violation it returns a
// *ScriptError; on success the stack holds the post-execution values.
//
// ctx may be nil for scripts that never reach a signature opcode; CHECKSIG/
// CHECKMULTISIG with a nil context return a ScriptError.
func (s *Script) Execute(stack *Stack, flags int, ctx *execContext) error {
	if len(s.raw) > MaxScriptSize {
		return newScriptError(errScriptSize, OP_INVALIDOPCODE, -1)
	}

	alt := NewStack()
	var condStack []bool // executed-branch flags
	lastSep := 0
	opCount := 0
	negate := 0
	minimal := flags&VerifyMinimalData != 0

	negated := func() bool { return negate > 0 }

	for ip := 0; ip < len(s.code); ip++ {
		op := s.code[ip]

		if op.Value == opParseError {
			return newScriptError(errBadOpcode, op.Value, ip)
		}
		if op.Data != nil && len(op.Data) > MaxScriptPush {
			return newScriptError(errPushSize, op.Value, ip)
		}
		if op.Value > OP_16 {
			opCount++
			if opCount > MaxScriptOps {
				return newScriptError(errOpCount, op.Value, ip)
			}
		}
		if op.IsDisabled() {
			return newScriptError(errDisabledOpcode, op.Value, ip)
		}

		// Skip data and most opcodes inside an unexecuted branch, but still
		// process branch opcodes so nesting stays balanced.
		if negated() && !op.IsBranch() {
			if stack.Len()+alt.Len() > MaxScriptStack {
				return newScriptError(errStackSize, op.Value, ip)
			}
			continue
		}

		if op.Data != nil {
			if minimal && !op.IsMinimal() {
				return newScriptError(errMinimalData, op.Value, ip)
			}
			stack.Push(op.Data)
			if stack.Len()+alt.Len() > MaxScriptStack {
				return newScriptError(errStackSize, op.Value, ip)
			}
			continue
		}

		switch op.Value {
		case OP_0:
			stack.PushInt(0)
		case OP_1NEGATE:
			stack.PushInt(-1)
		case OP_1, OP_2, OP_3, OP_4, OP_5, OP_6, OP_7, OP_8,
			OP_9, OP_10, OP_11, OP_12, OP_13, OP_14, OP_15, OP_16:
			stack.PushInt(ScriptNum(op.Value - 0x50))

		case OP_NOP:
			// no-op

		case OP_NOP1, OP_NOP4, OP_NOP5, OP_NOP6, OP_NOP7,
			OP_NOP8, OP_NOP9, OP_NOP10:
			if flags&VerifyDiscourageUpgradableNops != 0 {
				return newScriptError(errDiscourageUpgradableNops, op.Value, ip)
			}

		case OP_CHECKLOCKTIMEVERIFY, OP_CHECKSEQUENCEVERIFY:
			// Timelock opcodes require transaction context the script package
			// does not model on its own; treat as upgradable NOP unless
			// discouraged. Full timelock enforcement belongs to chain
			// validation (deferred — see package report).
			if flags&VerifyDiscourageUpgradableNops != 0 {
				return newScriptError(errDiscourageUpgradableNops, op.Value, ip)
			}

		case OP_IF, OP_NOTIF:
			val := false
			if !negated() {
				if stack.Len() < 1 {
					return newScriptError(errUnbalancedConditional, op.Value, ip)
				}
				item, _ := stack.Top(-1)
				if flags&VerifyMinimalIf != 0 {
					if len(item) > 1 || (len(item) == 1 && item[0] != 1) {
						return newScriptError(errMinimalIf, op.Value, ip)
					}
				}
				val = castToBool(item)
				if op.Value == OP_NOTIF {
					val = !val
				}
				if _, err := stack.Pop(); err != nil {
					return newScriptError(errInvalidStackOperation, op.Value, ip)
				}
			}
			condStack = append(condStack, val)
			if !val {
				negate++
			}

		case OP_ELSE:
			if len(condStack) == 0 {
				return newScriptError(errUnbalancedConditional, op.Value, ip)
			}
			top := &condStack[len(condStack)-1]
			*top = !*top
			if !*top {
				negate++
			} else {
				negate--
			}

		case OP_ENDIF:
			if len(condStack) == 0 {
				return newScriptError(errUnbalancedConditional, op.Value, ip)
			}
			top := condStack[len(condStack)-1]
			condStack = condStack[:len(condStack)-1]
			if !top {
				negate--
			}

		case OP_VERIFY:
			if stack.Len() == 0 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			b, _ := stack.Bool(-1)
			if !b {
				return newScriptError(errVerify, op.Value, ip)
			}
			_, _ = stack.Pop()

		case OP_RETURN:
			return newScriptError(errOpReturn, op.Value, ip)

		case OP_TOALTSTACK:
			item, err := stack.Pop()
			if err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			alt.Push(item)

		case OP_FROMALTSTACK:
			item, err := alt.Pop()
			if err != nil {
				return newScriptError(errInvalidAltstackOperation, op.Value, ip)
			}
			stack.Push(item)

		case OP_2DROP:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			_, _ = stack.Pop()
			_, _ = stack.Pop()

		case OP_2DUP:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			a, _ := stack.Top(-2)
			b, _ := stack.Top(-1)
			stack.Push(a)
			stack.Push(b)

		case OP_3DUP:
			if stack.Len() < 3 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			a, _ := stack.Top(-3)
			b, _ := stack.Top(-2)
			c, _ := stack.Top(-1)
			stack.Push(a)
			stack.Push(b)
			stack.Push(c)

		case OP_2OVER:
			if stack.Len() < 4 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			a, _ := stack.Top(-4)
			b, _ := stack.Top(-3)
			stack.Push(a)
			stack.Push(b)

		case OP_2ROT:
			if stack.Len() < 6 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			a, _ := stack.Top(-6)
			b, _ := stack.Top(-5)
			if err := stack.Erase(-6); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			if err := stack.Erase(-5); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			stack.Push(a)
			stack.Push(b)

		case OP_2SWAP:
			if stack.Len() < 4 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			if err := swap(stack, -4, -2); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			if err := swap(stack, -3, -1); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}

		case OP_IFDUP:
			if stack.Len() < 1 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Top(-1)
			if castToBool(item) {
				stack.Push(item)
			}

		case OP_DEPTH:
			stack.PushInt(ScriptNum(stack.Len()))

		case OP_DROP:
			if stack.Len() < 1 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			_, _ = stack.Pop()

		case OP_DUP:
			if stack.Len() < 1 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Top(-1)
			stack.Push(item)

		case OP_NIP:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			if err := stack.Erase(-2); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}

		case OP_OVER:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Top(-2)
			stack.Push(item)

		case OP_PICK, OP_ROLL:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			num, err := stack.Num(-1, minimal, scriptNumMaxBytes)
			if err != nil {
				return wrapNumError(err, op.Value, ip)
			}
			_, _ = stack.Pop()
			n := num.Int()
			if n < 0 || n >= stack.Len() {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Top(-n - 1)
			if op.Value == OP_ROLL {
				if err := stack.Erase(-n - 1); err != nil {
					return newScriptError(errInvalidStackOperation, op.Value, ip)
				}
			}
			stack.Push(item)

		case OP_ROT:
			if stack.Len() < 3 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			if err := swap(stack, -3, -2); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			if err := swap(stack, -2, -1); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}

		case OP_SWAP:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			if err := swap(stack, -2, -1); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}

		case OP_TUCK:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Top(-1)
			if err := stack.Insert(-2, item); err != nil {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}

		case OP_SIZE:
			if stack.Len() < 1 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Top(-1)
			stack.PushInt(ScriptNum(len(item)))

		case OP_EQUAL, OP_EQUALVERIFY:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			a, _ := stack.Pop()
			b, _ := stack.Pop()
			res := equalBytes(a, b)
			stack.PushBool(res)
			if op.Value == OP_EQUALVERIFY {
				if !res {
					return newScriptError(errEqualVerify, op.Value, ip)
				}
				_, _ = stack.Pop()
			}

		case OP_1ADD, OP_1SUB, OP_NEGATE, OP_ABS, OP_NOT, OP_0NOTEQUAL:
			if stack.Len() < 1 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			num, err := stack.Num(-1, minimal, scriptNumMaxBytes)
			if err != nil {
				return wrapNumError(err, op.Value, ip)
			}
			_, _ = stack.Pop()
			var res ScriptNum
			switch op.Value {
			case OP_1ADD:
				res = num + 1
			case OP_1SUB:
				res = num - 1
			case OP_NEGATE:
				res = -num
			case OP_ABS:
				if num < 0 {
					res = -num
				} else {
					res = num
				}
			case OP_NOT:
				res = boolNum(num == 0)
			case OP_0NOTEQUAL:
				res = boolNum(num != 0)
			}
			stack.PushInt(res)

		case OP_ADD, OP_SUB, OP_BOOLAND, OP_BOOLOR,
			OP_NUMEQUAL, OP_NUMEQUALVERIFY, OP_NUMNOTEQUAL,
			OP_LESSTHAN, OP_GREATERTHAN, OP_LESSTHANOREQUAL,
			OP_GREATERTHANOREQUAL, OP_MIN, OP_MAX:
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			n1, err := stack.Num(-2, minimal, scriptNumMaxBytes)
			if err != nil {
				return wrapNumError(err, op.Value, ip)
			}
			n2, err := stack.Num(-1, minimal, scriptNumMaxBytes)
			if err != nil {
				return wrapNumError(err, op.Value, ip)
			}
			_, _ = stack.Pop()
			_, _ = stack.Pop()
			var res ScriptNum
			switch op.Value {
			case OP_ADD:
				res = n1 + n2
			case OP_SUB:
				res = n1 - n2
			case OP_BOOLAND:
				res = boolNum(n1 != 0 && n2 != 0)
			case OP_BOOLOR:
				res = boolNum(n1 != 0 || n2 != 0)
			case OP_NUMEQUAL, OP_NUMEQUALVERIFY:
				res = boolNum(n1 == n2)
			case OP_NUMNOTEQUAL:
				res = boolNum(n1 != n2)
			case OP_LESSTHAN:
				res = boolNum(n1 < n2)
			case OP_GREATERTHAN:
				res = boolNum(n1 > n2)
			case OP_LESSTHANOREQUAL:
				res = boolNum(n1 <= n2)
			case OP_GREATERTHANOREQUAL:
				res = boolNum(n1 >= n2)
			case OP_MIN:
				if n1 < n2 {
					res = n1
				} else {
					res = n2
				}
			case OP_MAX:
				if n1 > n2 {
					res = n1
				} else {
					res = n2
				}
			}
			stack.PushInt(res)
			if op.Value == OP_NUMEQUALVERIFY {
				if res == 0 {
					return newScriptError(errNumEqualVerify, op.Value, ip)
				}
				_, _ = stack.Pop()
			}

		case OP_WITHIN:
			if stack.Len() < 3 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			n1, err := stack.Num(-3, minimal, scriptNumMaxBytes)
			if err != nil {
				return wrapNumError(err, op.Value, ip)
			}
			n2, err := stack.Num(-2, minimal, scriptNumMaxBytes)
			if err != nil {
				return wrapNumError(err, op.Value, ip)
			}
			n3, err := stack.Num(-1, minimal, scriptNumMaxBytes)
			if err != nil {
				return wrapNumError(err, op.Value, ip)
			}
			_, _ = stack.Pop()
			_, _ = stack.Pop()
			_, _ = stack.Pop()
			stack.PushBool(n2 <= n1 && n1 < n3)

		case OP_RIPEMD160, OP_SHA1, OP_SHA256, OP_HASH160,
			OP_HASH256, OP_BLAKE160, OP_BLAKE256, OP_SHA3, OP_KECCAK:
			if stack.Len() < 1 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Pop()
			digest, err := hashOp(op.Value, item)
			if err != nil {
				return newScriptError(errUnknownError, op.Value, ip)
			}
			stack.Push(digest)

		case OP_CODESEPARATOR:
			lastSep = ip + 1

		case OP_CHECKSIG, OP_CHECKSIGVERIFY:
			if ctx == nil {
				return newScriptError(errUnknownError, op.Value, ip)
			}
			if stack.Len() < 2 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			sig, _ := stack.Top(-2)
			key, _ := stack.Top(-1)

			subscript, err := s.GetSubscript(lastSep)
			if err != nil {
				return newScriptError(errUnknownError, op.Value, ip)
			}
			if serr := validateKey(key); serr != nil {
				return serr.withOp(op.Value, ip)
			}
			if serr := validateSignature(sig); serr != nil {
				return serr.withOp(op.Value, ip)
			}

			res := false
			if len(sig) > 0 {
				ok, herr := ctx.checksig(subscript, sig, key)
				if herr != nil {
					return newScriptError(errUnknownError, op.Value, ip)
				}
				res = ok
			}
			if !res && flags&VerifyNullFail != 0 && len(sig) != 0 {
				return newScriptError(errNullFail, op.Value, ip)
			}
			_, _ = stack.Pop()
			_, _ = stack.Pop()
			stack.PushBool(res)
			if op.Value == OP_CHECKSIGVERIFY {
				if !res {
					return newScriptError(errCheckSigVerify, op.Value, ip)
				}
				_, _ = stack.Pop()
			}

		case OP_CHECKMULTISIG, OP_CHECKMULTISIGVERIFY:
			if ctx == nil {
				return newScriptError(errUnknownError, op.Value, ip)
			}
			if serr := s.checkMultisig(stack, op, ip, flags, ctx, lastSep, minimal, &opCount); serr != nil {
				return serr
			}

		case OP_TYPE:
			if stack.Len() < 1 {
				return newScriptError(errInvalidStackOperation, op.Value, ip)
			}
			item, _ := stack.Pop()
			stack.PushInt(ScriptNum(itemType(item)))

		default:
			return newScriptError(errBadOpcode, op.Value, ip)
		}

		if stack.Len()+alt.Len() > MaxScriptStack {
			return newScriptError(errStackSize, op.Value, ip)
		}
	}

	if len(condStack) != 0 {
		return newScriptError(errUnbalancedConditional, OP_INVALIDOPCODE, -1)
	}

	return nil
}

// checkMultisig implements OP_CHECKMULTISIG(VERIFY). It is factored out of
// Execute to keep the dispatch switch readable; the logic mirrors the HSD
// implementation including the leading null-dummy stack element.
func (s *Script) checkMultisig(stack *Stack, op Opcode, ip, flags int, ctx *execContext, lastSep int, minimal bool, opCount *int) *ScriptError {
	i := 1
	if stack.Len() < i {
		return newScriptError(errInvalidStackOperation, op.Value, ip)
	}

	n, err := stack.Int(-i, minimal, scriptNumMaxBytes)
	if err != nil {
		return wrapNumScriptError(err, op.Value, ip)
	}
	okey := n + 2
	if n < 0 || n > MaxMultisigPubkeys {
		return newScriptError(errPubkeyCount, op.Value, ip)
	}
	*opCount += n
	if *opCount > MaxScriptOps {
		return newScriptError(errOpCount, op.Value, ip)
	}

	i++
	ikey := i
	i += n
	if stack.Len() < i {
		return newScriptError(errInvalidStackOperation, op.Value, ip)
	}

	m, err := stack.Int(-i, minimal, scriptNumMaxBytes)
	if err != nil {
		return wrapNumScriptError(err, op.Value, ip)
	}
	if m < 0 || m > n {
		return newScriptError(errSigCount, op.Value, ip)
	}

	i++
	isig := i
	i += m
	if stack.Len() < i {
		return newScriptError(errInvalidStackOperation, op.Value, ip)
	}

	subscript, gerr := s.GetSubscript(lastSep)
	if gerr != nil {
		return newScriptError(errUnknownError, op.Value, ip)
	}

	res := true
	for res && m > 0 {
		sig, _ := stack.Top(-isig)
		key, _ := stack.Top(-ikey)

		if serr := validateKey(key); serr != nil {
			return serr.withOp(op.Value, ip)
		}
		if serr := validateSignature(sig); serr != nil {
			return serr.withOp(op.Value, ip)
		}

		if len(sig) > 0 {
			ok, herr := ctx.checksig(subscript, sig, key)
			if herr != nil {
				return newScriptError(errUnknownError, op.Value, ip)
			}
			if ok {
				isig++
				m--
			}
		}
		ikey++
		n--
		if m > n {
			res = false
		}
	}

	for i > 1 {
		if !res && flags&VerifyNullFail != 0 {
			top, _ := stack.Top(-1)
			if okey == 0 && len(top) != 0 {
				return newScriptError(errNullFail, op.Value, ip)
			}
		}
		if okey > 0 {
			okey--
		}
		if _, err := stack.Pop(); err != nil {
			return newScriptError(errInvalidStackOperation, op.Value, ip)
		}
		i--
	}

	if stack.Len() < 1 {
		return newScriptError(errInvalidStackOperation, op.Value, ip)
	}
	dummy, _ := stack.Top(-1)
	if len(dummy) != 0 {
		return newScriptError(errSigNullDummy, op.Value, ip)
	}
	_, _ = stack.Pop()
	stack.PushBool(res)

	if op.Value == OP_CHECKMULTISIGVERIFY {
		if !res {
			return newScriptError(errCheckMultisigVerify, op.Value, ip)
		}
		_, _ = stack.Pop()
	}
	return nil
}

// checksig computes the signature hash for subscript and verifies the 65-byte
// signature against the public key, in this context's input.
func (c *execContext) checksig(subscript *Script, sig, key []byte) (bool, error) {
	hashType := uint32(sig[len(sig)-1])
	hash, err := SignatureHash(c.tx, c.index, subscript, c.value, hashType)
	if err != nil {
		return false, err
	}
	return verifySig(hash, sig, key)
}

// swap exchanges two stack items by (possibly negative) index.
func swap(stack *Stack, i, j int) error {
	a, err := stack.Top(i)
	if err != nil {
		return err
	}
	b, err := stack.Top(j)
	if err != nil {
		return err
	}
	if err := stack.Set(i, b); err != nil {
		return err
	}
	return stack.Set(j, a)
}

// boolNum maps a boolean to the script-number 1 or 0.
func boolNum(b bool) ScriptNum {
	if b {
		return 1
	}
	return 0
}

// hashOp applies the digest opcode to data.
func hashOp(opValue int, data []byte) ([]byte, error) {
	switch opValue {
	case OP_RIPEMD160:
		h := ripemd160.New()
		_, _ = h.Write(data)
		return h.Sum(nil), nil
	case OP_SHA1:
		sum := sha1.Sum(data)
		return sum[:], nil
	case OP_SHA256:
		sum := sha256.Sum256(data)
		return sum[:], nil
	case OP_HASH160:
		first := sha256.Sum256(data)
		h := ripemd160.New()
		_, _ = h.Write(first[:])
		return h.Sum(nil), nil
	case OP_HASH256:
		first := sha256.Sum256(data)
		second := sha256.Sum256(first[:])
		return second[:], nil
	case OP_BLAKE160:
		h, err := blake2b.New(20, nil)
		if err != nil {
			return nil, err
		}
		_, _ = h.Write(data)
		return h.Sum(nil), nil
	case OP_BLAKE256:
		sum := blake2b.Sum256(data)
		return sum[:], nil
	case OP_SHA3:
		sum := sha3.Sum256(data)
		return sum[:], nil
	case OP_KECCAK:
		h := sha3.NewLegacyKeccak256()
		_, _ = h.Write(data)
		return h.Sum(nil), nil
	default:
		return nil, core.E("script.hashOp", "not a hash opcode", nil)
	}
}

// itemType returns the OP_TYPE classification of a stack item: 0 for an empty
// item, 1 for a single byte, 2 otherwise. Mirrors the HSD OP_TYPE handler.
func itemType(item []byte) int {
	switch {
	case len(item) == 0:
		return 0
	case len(item) == 1:
		return 1
	default:
		return 2
	}
}

// validateKey enforces the compressed-key encoding required for CHECKSIG keys.
func validateKey(key []byte) *ScriptError {
	if len(key) != pubkeyLen || (key[0] != 0x02 && key[0] != 0x03) {
		return newScriptError(errPubkeyEncoding, OP_INVALIDOPCODE, -1)
	}
	return nil
}

// validateSignature enforces the 65-byte signature encoding (empty allowed).
func validateSignature(sig []byte) *ScriptError {
	if len(sig) == 0 {
		return nil
	}
	if len(sig) != sigLen {
		return newScriptError(errSigEncoding, OP_INVALIDOPCODE, -1)
	}
	hashType := sig[sigLen-1]
	masked := hashType &^ byte(SighashNoInput)
	masked &^= byte(SighashAnyoneCanPay)
	if masked < SighashAll || masked > SighashSingleReverse {
		return newScriptError(errSigEncoding, OP_INVALIDOPCODE, -1)
	}
	return nil
}

// wrapNumError converts a script-number decode error into a script error tagged
// with the offending opcode (for the value-returning helpers).
func wrapNumError(err error, opValue, ip int) error {
	return wrapNumScriptError(err, opValue, ip)
}

func wrapNumScriptError(err error, opValue, ip int) *ScriptError {
	_ = err
	return newScriptError(errUnknownError, opValue, ip)
}
