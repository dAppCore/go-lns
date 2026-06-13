// SPDX-License-Identifier: EUPL-1.2

package script

import core "dappco.re/go"

// Script error codes. They mirror the HSD ScriptError codes
// (lib/script/scripterror.js) so failures can be matched against the reference
// behaviour. The string form is the code itself.
const (
	errUnknownError                = "UNKNOWN_ERROR"
	errBadOpcode                   = "BAD_OPCODE"
	errDisabledOpcode              = "DISABLED_OPCODE"
	errPushSize                    = "PUSH_SIZE"
	errOpCount                     = "OP_COUNT"
	errStackSize                   = "STACK_SIZE"
	errScriptSize                  = "SCRIPT_SIZE"
	errMinimalData                 = "MINIMALDATA"
	errMinimalIf                   = "MINIMALIF"
	errInvalidStackOperation       = "INVALID_STACK_OPERATION"
	errInvalidAltstackOperation    = "INVALID_ALTSTACK_OPERATION"
	errUnbalancedConditional       = "UNBALANCED_CONDITIONAL"
	errVerify                      = "VERIFY"
	errEqualVerify                 = "EQUALVERIFY"
	errNumEqualVerify              = "NUMEQUALVERIFY"
	errCheckSigVerify              = "CHECKSIGVERIFY"
	errCheckMultisigVerify         = "CHECKMULTISIGVERIFY"
	errOpReturn                    = "OP_RETURN"
	errPubkeyCount                 = "PUBKEY_COUNT"
	errSigCount                    = "SIG_COUNT"
	errPubkeyEncoding              = "PUBKEY_ENCODING"
	errSigEncoding                 = "SIG_ENCODING"
	errNullFail                    = "NULLFAIL"
	errSigNullDummy                = "SIG_NULLDUMMY"
	errEvalFalse                   = "EVAL_FALSE"
	errWitnessProgramMismatch      = "WITNESS_PROGRAM_MISMATCH"
	errWitnessProgramWitnessEmpty  = "WITNESS_PROGRAM_WITNESS_EMPTY"
	errWitnessProgramWrongLength   = "WITNESS_PROGRAM_WRONG_LENGTH"
	errDiscourageUpgradableNops    = "DISCOURAGE_UPGRADABLE_NOPS"
	errDiscourageUpgradableWitness = "DISCOURAGE_UPGRADABLE_WITNESS_PROGRAM"
)

// ScriptError is an error raised by the script interpreter or the witness
// verifier. Code is one of the err* constants above; Op and IP locate the
// offending opcode when known (Op is -1, IP is -1 otherwise).
//
//	if serr, ok := err.(*script.ScriptError); ok && serr.Code == "EVAL_FALSE" { ... }
type ScriptError struct {
	// Code is the machine-readable error code (e.g. "EVAL_FALSE").
	Code string
	// Op is the opcode value where the error occurred, or OP_INVALIDOPCODE.
	Op int
	// IP is the instruction pointer (opcode index), or -1 when not applicable.
	IP int
}

// newScriptError constructs a ScriptError for the given code and location.
func newScriptError(code string, op, ip int) *ScriptError {
	return &ScriptError{Code: code, Op: op, IP: ip}
}

// withOp returns a copy of the error tagged with an opcode location. It lets the
// encoding validators (which do not know the instruction pointer) be located by
// the caller.
func (e *ScriptError) withOp(op, ip int) *ScriptError {
	return &ScriptError{Code: e.Code, Op: op, IP: ip}
}

// Error implements the error interface. The message includes the opcode symbol
// and instruction pointer when available, matching the HSD message format.
func (e *ScriptError) Error() string {
	if e.Op >= 0 && e.IP >= 0 {
		return core.Sprintf("%s (op=%s, ip=%d)", e.Code, OpcodeName(e.Op), e.IP)
	}
	return e.Code
}
