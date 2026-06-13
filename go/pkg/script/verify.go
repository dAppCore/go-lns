// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"dappco.re/go/lns/pkg/primitives"
)

// VerifyInput verifies the witness of a single transaction input against the
// witness programme committed to by the output it spends. It is the entry point
// full transaction validation calls per input; the chain rejects a transaction
// if any input fails.
//
// It mirrors HSD's Script.verify (lib/script/script.js):
//
//   - Version 31 (nulldata) outputs are unspendable: OP_RETURN error.
//   - Version 0 / 32-byte: pay-to-witness-script-hash. The last witness item is
//     the redeem script; its SHA3-256 digest must equal the programme hash. The
//     remaining items form the input stack, the redeem script executes against
//     it, and the final stack must be a single truthy value.
//   - Version 0 / 20-byte: pay-to-witness-pubkey-hash. The witness must be
//     exactly [signature, pubkey]; the generated P2PKH template executes.
//   - Version 0, other lengths: WITNESS_PROGRAM_WRONG_LENGTH.
//   - Other versions: reserved upgradable programmes — accepted unless
//     VerifyDiscourageUpgradableWitnessProgram is set.
//
// tx, index, and value provide the context the signature opcodes need: the
// spending transaction, the index of this input, and the spent output's value.
//
//	err := script.VerifyInput(prevOutput, tx, 0, script.MandatoryVerifyFlags)
func VerifyInput(prevOutput primitives.Output, tx primitives.Transaction, index int, flags int) error {
	if index < 0 || index >= len(tx.Inputs) {
		return newScriptError(errUnknownError, OP_INVALIDOPCODE, -1)
	}

	addr := prevOutput.Address
	witness := tx.Inputs[index].Witness

	return verifyWitness(witness, addr, tx, index, prevOutput.Value, flags)
}

// Verify checks a witness stack against an address (witness programme) for the
// given input. It is the address-level form of VerifyInput, useful for callers
// that already hold the programme and spent value separately.
//
//	err := script.Verify(witness, addr, tx, 0, value, script.StandardVerifyFlags)
func Verify(witness [][]byte, addr primitives.Address, tx primitives.Transaction, index int, value uint64, flags int) error {
	return verifyWitness(witness, addr, tx, index, value, flags)
}

func verifyWitness(witnessItems [][]byte, addr primitives.Address, tx primitives.Transaction, index int, value uint64, flags int) error {
	if addr.Version > 31 {
		return newScriptError(errWitnessProgramWrongLength, OP_INVALIDOPCODE, -1)
	}
	if addr.Version == 31 {
		return newScriptError(errOpReturn, OP_INVALIDOPCODE, -1)
	}
	if len(witnessItems) > MaxScriptStack {
		return newScriptError(errStackSize, OP_INVALIDOPCODE, -1)
	}

	stack := NewStack()
	for _, item := range witnessItems {
		stack.Push(append([]byte(nil), item...))
	}

	var redeem *Script

	if addr.Version == 0 {
		switch len(addr.Hash) {
		case 32:
			if stack.Len() == 0 {
				return newScriptError(errWitnessProgramWitnessEmpty, OP_INVALIDOPCODE, -1)
			}
			witnessScript, _ := stack.Pop()
			if len(witnessScript) > MaxScriptSize {
				return newScriptError(errScriptSize, OP_INVALIDOPCODE, -1)
			}
			decoded, err := DecodeScript(witnessScript)
			if err != nil {
				return newScriptError(errScriptSize, OP_INVALIDOPCODE, -1)
			}
			if !equalBytes(decoded.Sha3(), addr.Hash) {
				return newScriptError(errWitnessProgramMismatch, OP_INVALIDOPCODE, -1)
			}
			redeem = decoded

		case 20:
			if stack.Len() != 2 {
				return newScriptError(errWitnessProgramMismatch, OP_INVALIDOPCODE, -1)
			}
			redeem = NewScript().FromPubkeyhash(addr.Hash)

		default:
			return newScriptError(errWitnessProgramWrongLength, OP_INVALIDOPCODE, -1)
		}
	} else {
		if flags&VerifyDiscourageUpgradableWitnessProgram != 0 {
			return newScriptError(errDiscourageUpgradableWitness, OP_INVALIDOPCODE, -1)
		}
		// Reserved upgradable witness programme: accepted.
		return nil
	}

	// Witness items still have a per-push size limit.
	for j := 0; j < stack.Len(); j++ {
		item, _ := stack.Top(j)
		if len(item) > MaxScriptPush {
			return newScriptError(errPushSize, OP_INVALIDOPCODE, -1)
		}
	}

	ctx := &execContext{tx: tx, index: index, value: value}
	if err := redeem.Execute(stack, flags, ctx); err != nil {
		return err
	}

	if stack.Len() != 1 {
		return newScriptError(errEvalFalse, OP_INVALIDOPCODE, -1)
	}
	top, _ := stack.Top(-1)
	if !castToBool(top) {
		return newScriptError(errEvalFalse, OP_INVALIDOPCODE, -1)
	}

	return nil
}
