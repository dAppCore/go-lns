// SPDX-License-Identifier: EUPL-1.2

package script

import (
	"strings"
	"testing"
)

func TestScriptError_Error_Good(t *testing.T) {
	// With opcode/ip location, the message includes the symbol and pointer.
	e := newScriptError(errEqualVerify, OP_EQUALVERIFY, 3)
	msg := e.Error()
	if !strings.Contains(msg, "EQUALVERIFY") || !strings.Contains(msg, "OP_EQUALVERIFY") || !strings.Contains(msg, "ip=3") {
		t.Fatalf("error message = %q", msg)
	}
}

func TestScriptError_Error_Bad(t *testing.T) {
	// Without a location, the message is just the code.
	e := newScriptError(errEvalFalse, OP_INVALIDOPCODE, -1)
	if e.Error() != "EVAL_FALSE" {
		t.Fatalf("error message = %q, want EVAL_FALSE", e.Error())
	}
}

func TestScriptError_WithOp_Ugly(t *testing.T) {
	base := newScriptError(errSigEncoding, OP_INVALIDOPCODE, -1)
	tagged := base.withOp(OP_CHECKSIG, 7)
	if tagged.Code != errSigEncoding {
		t.Fatalf("code lost: %q", tagged.Code)
	}
	if tagged.Op != OP_CHECKSIG || tagged.IP != 7 {
		t.Fatalf("location not applied: op=%d ip=%d", tagged.Op, tagged.IP)
	}
	// The original must be unchanged (withOp returns a copy).
	if base.Op != OP_INVALIDOPCODE || base.IP != -1 {
		t.Fatal("withOp mutated the original error")
	}
}
