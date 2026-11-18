// SPDX-License-Identifier: EUPL-1.2

package primitives

import core "github.com/dappcore/go"

// bech32Charset is the bech32 alphabet (BIP-173): the 32 symbols data values map
// to, excluding the easily-confused 1, b, i, and o.
const bech32Charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// Encode renders the address as a bech32 segwit-style string for the given
// human-readable prefix (the network addressPrefix), matching HSD's
// Address.toString (docs/js-primitives/address.js): the witness version is the
// first data symbol and the programme hash is 8-to-5-bit regrouped, followed by
// the 6-symbol bech32 checksum. The prefix is network-specific (mainnet, testnet,
// and regtest differ), so the caller supplies it.
//
//	s, err := addr.Encode("lts")
func (a Address) Encode(hrp string) (string, error) {
	hrp = core.Lower(core.Trim(hrp))
	if hrp == "" {
		return "", core.E("primitives.Address.Encode", "human-readable prefix is required", nil)
	}
	if a.Version > 31 {
		return "", core.E("primitives.Address.Encode", "witness version out of range", nil)
	}
	if len(a.Hash) < 2 || len(a.Hash) > 40 {
		return "", core.E("primitives.Address.Encode", "programme hash length out of range", nil)
	}

	conv, err := convertBits(a.Hash, 8, 5, true)
	if err != nil {
		return "", core.E("primitives.Address.Encode", "regroup programme bits", err)
	}

	data := make([]byte, 0, len(conv)+1)
	data = append(data, a.Version)
	data = append(data, conv...)

	return bech32Encode(hrp, data)
}

// bech32Encode assembles a bech32 string from a prefix and 5-bit data values
// (the witness version followed by the regrouped programme). It appends the
// BIP-173 checksum (constant 1) and maps each value through the bech32 alphabet.
func bech32Encode(hrp string, data []byte) (string, error) {
	checksum := bech32CreateChecksum(hrp, data)
	combined := make([]byte, 0, len(data)+len(checksum))
	combined = append(combined, data...)
	combined = append(combined, checksum...)

	out := make([]byte, 0, len(hrp)+1+len(combined))
	out = append(out, hrp...)
	out = append(out, '1')
	for _, value := range combined {
		if int(value) >= len(bech32Charset) {
			return "", core.E("primitives.bech32Encode", "data value out of range", nil)
		}
		out = append(out, bech32Charset[value])
	}
	return string(out), nil
}

// bech32Polymod computes the bech32 checksum polynomial over the supplied
// values, the core of the BCH checksum (BIP-173).
func bech32Polymod(values []byte) uint32 {
	generators := [5]uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}
	chk := uint32(1)
	for _, value := range values {
		top := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ uint32(value)
		for i := range 5 {
			if (top>>uint(i))&1 == 1 {
				chk ^= generators[i]
			}
		}
	}
	return chk
}

// bech32HrpExpand expands the human-readable prefix into the value sequence the
// checksum is computed over: the high bits, a zero separator, then the low bits.
func bech32HrpExpand(hrp string) []byte {
	out := make([]byte, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		out = append(out, hrp[i]>>5)
	}
	out = append(out, 0)
	for i := 0; i < len(hrp); i++ {
		out = append(out, hrp[i]&31)
	}
	return out
}

// bech32CreateChecksum produces the 6 checksum symbols for a prefix and data,
// using checksum constant 1 (bech32, as Handshake/HSD addresses use).
func bech32CreateChecksum(hrp string, data []byte) []byte {
	values := append(bech32HrpExpand(hrp), data...)
	values = append(values, 0, 0, 0, 0, 0, 0)
	polymod := bech32Polymod(values) ^ 1
	checksum := make([]byte, 6)
	for i := range 6 {
		checksum[i] = byte((polymod >> uint(5*(5-i))) & 31)
	}
	return checksum
}

// convertBits regroups a byte slice from one bit width to another (8→5 for
// encoding), padding the final group with zero bits when pad is true. It is the
// witness-programme regrouping bech32 requires.
func convertBits(data []byte, fromBits, toBits uint, pad bool) ([]byte, error) {
	var acc uint32
	var bits uint
	maxv := uint32((1 << toBits) - 1)
	out := make([]byte, 0, len(data)*int(fromBits)/int(toBits)+1)

	for _, value := range data {
		acc = (acc << fromBits) | uint32(value)
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			out = append(out, byte((acc>>bits)&maxv))
		}
	}

	if pad {
		if bits > 0 {
			out = append(out, byte((acc<<(toBits-bits))&maxv))
		}
	} else if bits >= fromBits || ((acc<<(toBits-bits))&maxv) != 0 {
		return nil, core.E("primitives.convertBits", "invalid padding", nil)
	}

	return out, nil
}
