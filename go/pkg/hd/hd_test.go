// SPDX-License-Identifier: EUPL-1.2

package hd

import (
	"testing"

	core "dappco.re/go"
)

// BIP-32 Test Vector 1 (https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki).
// These seed, chain-code, and private-key values are implementation-independent
// and define the secp256k1 BIP-32 derivation maths that HSD (and therefore the
// Lethean name chain) inherit. They are the known-answer vectors for this
// package.
const (
	vector1SeedHex = "000102030405060708090a0b0c0d0e0f"

	// m
	vector1MasterChainHex = "873dff81c02f525623fd1fe5167eac3a55a049de3d314bb42ee227ffed37d508"
	vector1MasterPrivHex  = "e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35"
	vector1MasterPubHex   = "0339a36013301597daef41fbe593a02cc513d0b55527ec2df1050e2e8ff49c85c2"

	// m/0'
	vector1M0hChainHex = "47fdacbd0f1097043b78c63c20c34ef4ed9a111d980047ad16282c7ae6236141"
	vector1M0hPrivHex  = "edb2e14f9ee77d26dd93b4ecede8d16ed408ce149b6cd80b0715a2d911a0afea"
	vector1M0hPubHex   = "035a784662a4a20a65bf6aab9ae98a6c068a81c52e4b032c0fb5400c706cfccc56"

	// m/0'/1
	vector1M0h1ChainHex = "2a7857631386ba23dacac34180dd1983734e444fdbf774041578e9b6adb37c19"
	vector1M0h1PrivHex  = "3c6cb8d0f6a264c91ea8b5030fadaa8e538b020f0a387421a12de9319dc93368"
	vector1M0h1PubHex   = "03501e454bf00751f24b1b489aa925215d66af2234e3891c3b21a52bedb3cd711c"

	// m/0'/1/2'
	vector1M0h12hChainHex = "04466b9cc8e161e966409ca52986c584f07e9dc81f735db683c3ff6ec7b1503f"
	vector1M0h12hPrivHex  = "cbce0d719ecf7431d88e6a89fa1483e02e35092af60c042b1df2ff59fa424dca"
	vector1M0h12hPubHex   = "0357bfe1e341d01c69fe5654309956cbea516822fba8a601743a012a7896ee8dc2"
)

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	r := core.HexDecode(s)
	if !r.OK {
		t.Fatalf("HexDecode(%q) failed: %v", s, r.Error())
	}
	return r.Value.([]byte)
}

func TestHd_ser32_Good(t *core.T) {
	got := ser32(1)
	want := []byte{0, 0, 0, 1}
	core.AssertEqual(t, core.HexEncode(want), core.HexEncode(got))
}

func TestHd_ser32_Bad(t *core.T) {
	got := ser32(0)
	core.AssertEqual(t, "00000000", core.HexEncode(got))
}

func TestHd_ser32_Ugly(t *core.T) {
	got := ser32(HardenedKeyStart)
	core.AssertEqual(t, "80000000", core.HexEncode(got))
}

func TestHd_parse256_Good(t *core.T) {
	n := parse256([]byte{0x00, 0x00, 0x01})
	core.AssertTrue(t, n.Int64() == 1)
}

func TestHd_parse256_Bad(t *core.T) {
	n := parse256([]byte{})
	core.AssertTrue(t, n.Sign() == 0)
}

func TestHd_parse256_Ugly(t *core.T) {
	// All-zero 32-byte slice parses to zero.
	n := parse256(make([]byte, 32))
	core.AssertTrue(t, n.Sign() == 0)
}

func TestHd_ser256_Good(t *core.T) {
	n := parse256([]byte{0x01})
	got := ser256(n)
	core.AssertEqual(t, 32, len(got))
	core.AssertEqual(t, byte(1), got[31])
}

func TestHd_ser256_Bad(t *core.T) {
	got := ser256(parse256(nil))
	core.AssertEqual(t, 32, len(got))
	for _, b := range got {
		core.AssertEqual(t, byte(0), b)
	}
}

func TestHd_ser256_Ugly(t *core.T) {
	// Round-trip a known master private scalar through parse/ser.
	raw := mustHex(t, vector1MasterPrivHex)
	got := ser256(parse256(raw))
	core.AssertEqual(t, vector1MasterPrivHex, core.HexEncode(got))
}

func TestHd_Constants_Good(t *core.T) {
	core.AssertEqual(t, uint32(0x80000000), HardenedKeyStart)
	core.AssertEqual(t, uint32(44), Bip44Purpose)
	core.AssertEqual(t, uint32(5353), CoinTypeLethean)
}

func TestHd_Constants_Bad(t *core.T) {
	core.AssertEqual(t, uint32(0), BranchReceive)
	core.AssertEqual(t, uint32(1), BranchChange)
}

func TestHd_Constants_Ugly(t *core.T) {
	core.AssertTrue(t, SeedMinLen < SeedMaxLen)
	core.AssertEqual(t, 16, SeedMinLen)
	core.AssertEqual(t, 64, SeedMaxLen)
}
