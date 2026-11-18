// SPDX-License-Identifier: EUPL-1.2

package script

import (
	core "github.com/dappcore/go"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"golang.org/x/crypto/blake2b"

	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// Sighash type flags. They mirror HSD's common.hashType (lib/script/common.js).
const (
	// SighashAll signs all outputs.
	SighashAll = 1
	// SighashNone signs no outputs (zeroes other sequences).
	SighashNone = 2
	// SighashSingle signs the output at the same index.
	SighashSingle = 3
	// SighashSingleReverse signs the output at the reversed index.
	SighashSingleReverse = 4
	// SighashNoInput nulls this input's outpoint commitment (see note below).
	SighashNoInput = 0x40
	// SighashAnyoneCanPay signs only the current input.
	SighashAnyoneCanPay = 0x80
)

// SignatureHash computes the signing digest for one input, using the Handshake
// segwit-style commitment (docs/js-primitives via tx.js signatureHash). It is a
// BLAKE2b digest over:
//
//	version ‖ hashPrevouts ‖ hashSequence ‖ thisInput.prevout ‖
//	varbytes(prevScript) ‖ value ‖ thisInput.sequence ‖ hashOutputs ‖
//	locktime ‖ hashType
//
// prev is the script being signed: for a P2WPKH spend it is the generated P2PKH
// template (FromPubkeyhash); for a P2WSH spend it is the redeem script. value is
// the spent output's value.
//
// Note (matching the upstream): SighashNoInput zeroes only this input's
// outpoint/sequence in the per-input portion, not the aggregated hashPrevouts/
// hashSequence, so it is not malleability-safe on its own.
//
//	hash, err := script.SignatureHash(tx, 0, redeem, value, script.SighashAll)
func SignatureHash(tx primitives.Transaction, index int, prev *Script, value uint64, hashType uint32) ([]byte, error) {
	if index < 0 || index >= len(tx.Inputs) {
		return nil, core.E("script.SignatureHash", "input index out of range", nil)
	}
	if prev == nil {
		return nil, core.E("script.SignatureHash", "nil previous script", nil)
	}

	input := tx.Inputs[index]

	var zero primitives.Hash
	prevouts := zero[:]
	sequences := zero[:]
	outputs := zero[:]

	// SIGHASH_NOINPUT blanks this input's outpoint/sequence (only this input).
	if hashType&SighashNoInput != 0 {
		input = primitives.Input{Prevout: primitives.NewOutpoint()}
	}

	masked := hashType & 0x1f

	if hashType&SighashAnyoneCanPay == 0 {
		buf := make([]byte, 0, len(tx.Inputs)*36)
		for _, in := range tx.Inputs {
			out, err := in.Prevout.MarshalBinary()
			if err != nil {
				return nil, core.E("script.SignatureHash", "marshal prevout", err)
			}
			buf = append(buf, out...)
		}
		h, err := blake2bDigest(buf)
		if err != nil {
			return nil, err
		}
		prevouts = h
	}

	if hashType&SighashAnyoneCanPay == 0 &&
		masked != SighashSingle &&
		masked != SighashSingleReverse &&
		masked != SighashNone {
		buf := make([]byte, 0, len(tx.Inputs)*4)
		for _, in := range tx.Inputs {
			buf = appendU32(buf, in.Sequence)
		}
		h, err := blake2bDigest(buf)
		if err != nil {
			return nil, err
		}
		sequences = h
	}

	switch {
	case masked != SighashSingle && masked != SighashSingleReverse && masked != SighashNone:
		buf := make([]byte, 0, 256)
		for i := range tx.Outputs {
			enc, err := encodeOutput(tx.Outputs[i])
			if err != nil {
				return nil, err
			}
			buf = append(buf, enc...)
		}
		h, err := blake2bDigest(buf)
		if err != nil {
			return nil, err
		}
		outputs = h
	case masked == SighashSingle:
		if index < len(tx.Outputs) {
			enc, err := encodeOutput(tx.Outputs[index])
			if err != nil {
				return nil, err
			}
			h, err := blake2bDigest(enc)
			if err != nil {
				return nil, err
			}
			outputs = h
		}
	case masked == SighashSingleReverse:
		if index < len(tx.Outputs) {
			enc, err := encodeOutput(tx.Outputs[(len(tx.Outputs)-1)-index])
			if err != nil {
				return nil, err
			}
			h, err := blake2bDigest(enc)
			if err != nil {
				return nil, err
			}
			outputs = h
		}
	}

	prevoutBytes, err := input.Prevout.MarshalBinary()
	if err != nil {
		return nil, core.E("script.SignatureHash", "marshal input prevout", err)
	}

	preimage := make([]byte, 0, 156+prev.GetVarSize())
	preimage = appendU32(preimage, tx.Version)
	preimage = append(preimage, prevouts...)
	preimage = append(preimage, sequences...)
	preimage = append(preimage, prevoutBytes...)
	preimage = appendVarBytes(preimage, prev.Encode())
	preimage = appendU64(preimage, value)
	preimage = appendU32(preimage, input.Sequence)
	preimage = append(preimage, outputs...)
	preimage = appendU32(preimage, tx.LockTime)
	preimage = appendU32(preimage, hashType)

	return blake2bDigest(preimage)
}

// Sign produces a 65-byte script signature for an input: a 64-byte compact
// secp256k1 signature (R‖S, low-S) over SignatureHash followed by the one-byte
// sighash type. It is the witness item a P2WPKH or redeem-script spend pushes.
//
//	sig, err := script.Sign(tx, 0, redeem, value, priv, script.SighashAll)
func Sign(tx primitives.Transaction, index int, prev *Script, value uint64, priv *secp256k1.PrivateKey, hashType uint32) ([]byte, error) {
	if priv == nil {
		return nil, core.E("script.Sign", "nil private key", nil)
	}
	hash, err := SignatureHash(tx, index, prev, value, hashType)
	if err != nil {
		return nil, err
	}
	sig := ecdsa.Sign(priv, hash)
	out := encodeCompactSig(sig)
	return append(out, byte(hashType)), nil
}

// Checksig verifies a 65-byte script signature against a 33-byte compressed
// public key for one input. The trailing sighash-type byte selects the digest.
// Mirrors the interpreter's checksig + tx.checksig.
//
//	ok, err := script.Checksig(tx, 0, redeem, value, sig, pubkey)
func Checksig(tx primitives.Transaction, index int, prev *Script, value uint64, sig, pub []byte) (bool, error) {
	if len(sig) == 0 {
		return false, nil
	}
	hashType := uint32(sig[len(sig)-1])
	hash, err := SignatureHash(tx, index, prev, value, hashType)
	if err != nil {
		return false, err
	}
	return verifySig(hash, sig, pub)
}

// verifySig checks a 65-byte signature (64-byte compact body + type byte)
// against a public key over the given message hash. It does not re-derive the
// hash; the caller supplies it.
func verifySig(hash, sig, pub []byte) (bool, error) {
	if len(sig) != sigLen {
		return false, nil
	}
	pk, err := secp256k1.ParsePubKey(pub)
	if err != nil {
		return false, nil
	}
	parsed, err := decodeCompactSig(sig[:64])
	if err != nil {
		return false, nil
	}
	return parsed.Verify(hash, pk), nil
}

// encodeCompactSig serialises an ECDSA signature as the 64-byte R‖S form used by
// the chain (big-endian 32-byte R, big-endian 32-byte S).
func encodeCompactSig(sig *ecdsa.Signature) []byte {
	r := sig.R()
	s := sig.S()
	rb := r.Bytes()
	sb := s.Bytes()
	out := make([]byte, 0, 64)
	out = append(out, rb[:]...)
	out = append(out, sb[:]...)
	return out
}

// decodeCompactSig parses a 64-byte R‖S signature into an ecdsa.Signature.
func decodeCompactSig(body []byte) (*ecdsa.Signature, error) {
	if len(body) != 64 {
		return nil, core.E("script.decodeCompactSig", "signature must be 64 bytes", nil)
	}
	var r, s secp256k1.ModNScalar
	var rb, sb [32]byte
	copy(rb[:], body[:32])
	copy(sb[:], body[32:])
	if r.SetBytes(&rb) != 0 {
		return nil, core.E("script.decodeCompactSig", "signature R overflows group order", nil)
	}
	if s.SetBytes(&sb) != 0 {
		return nil, core.E("script.decodeCompactSig", "signature S overflows group order", nil)
	}
	if r.IsZero() || s.IsZero() {
		return nil, core.E("script.decodeCompactSig", "signature R or S is zero", nil)
	}
	return ecdsa.NewSignature(&r, &s), nil
}

// encodeOutput serialises a transaction output in the HSD wire form:
//
//	value(8 LE) ‖ address.version(1) ‖ address.hashLen(1) ‖ address.hash ‖ covenant
//
// Mirrors Output.write (docs/js-primitives/output.js) and reuses
// primitives.Covenant.MarshalBinary for the covenant tail.
func encodeOutput(out primitives.Output) ([]byte, error) {
	cov, err := out.Covenant.MarshalBinary()
	if err != nil {
		return nil, core.E("script.encodeOutput", "marshal covenant", err)
	}
	if len(out.Address.Hash) > 0xff {
		return nil, core.E("script.encodeOutput", "address hash too long", nil)
	}
	buf := make([]byte, 0, 8+2+len(out.Address.Hash)+len(cov))
	buf = appendU64(buf, out.Value)
	buf = append(buf, out.Address.Version, byte(len(out.Address.Hash)))
	buf = append(buf, out.Address.Hash...)
	buf = append(buf, cov...)
	return buf, nil
}

// blake2bDigest returns the 32-byte BLAKE2b-256 digest of data, using the same
// constructor the rest of the name system uses (covenant.Blind, hd).
func blake2bDigest(data []byte) ([]byte, error) {
	h, err := blake2b.New256(nil)
	if err != nil {
		return nil, core.E("script.blake2bDigest", "init blake2b", err)
	}
	if _, err := h.Write(data); err != nil {
		return nil, core.E("script.blake2bDigest", "write blake2b", err)
	}
	return h.Sum(nil), nil
}
