// SPDX-License-Identifier: EUPL-1.2

package lns

import (
	core "dappco.re/go"
	"dappco.re/go/lns/pkg/primitives"
)

func TestActionHelpersBehaviour_actionUint32FromAny_Good(t *core.T) {
	cases := []struct {
		in   any
		want uint32
	}{
		{uint32(7), 7},
		{uint(8), 8},
		{int(9), 9},
		{int64(10), 10},
		{float64(11), 11},
		{"12", 12},
	}
	for _, tc := range cases {
		got, ok := actionUint32FromAny(tc.in)
		core.AssertTrue(t, ok)
		core.AssertEqual(t, tc.want, got)
	}
}

func TestActionHelpersBehaviour_actionUint32FromAny_Bad(t *core.T) {
	_, ok := actionUint32FromAny(int(-1))
	core.AssertFalse(t, ok)
	_, ok = actionUint32FromAny(int64(-1))
	core.AssertFalse(t, ok)
	_, ok = actionUint32FromAny(float64(-1))
	core.AssertFalse(t, ok)
	_, ok = actionUint32FromAny("not-a-number")
	core.AssertFalse(t, ok)
	_, ok = actionUint32FromAny(nil)
	core.AssertFalse(t, ok)
}

func TestActionHelpersBehaviour_actionUint32FromAny_Ugly(t *core.T) {
	// Whitespace around a numeric string is trimmed before parsing.
	got, ok := actionUint32FromAny("  42  ")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, uint32(42), got)
}

func TestActionHelpersBehaviour_actionUint64FromAny_Good(t *core.T) {
	cases := []struct {
		in   any
		want uint64
	}{
		{uint64(7), 7},
		{uint32(8), 8},
		{uint(9), 9},
		{int(10), 10},
		{int64(11), 11},
		{float64(12), 12},
		{"13", 13},
	}
	for _, tc := range cases {
		got, ok := actionUint64FromAny(tc.in)
		core.AssertTrue(t, ok)
		core.AssertEqual(t, tc.want, got)
	}
}

func TestActionHelpersBehaviour_actionUint64FromAny_Bad(t *core.T) {
	_, ok := actionUint64FromAny(int(-1))
	core.AssertFalse(t, ok)
	_, ok = actionUint64FromAny("bad")
	core.AssertFalse(t, ok)
	_, ok = actionUint64FromAny(nil)
	core.AssertFalse(t, ok)
}

func TestActionHelpersBehaviour_actionBoolFromAny(t *core.T) {
	got, ok := actionBoolFromAny(true)
	core.AssertTrue(t, ok)
	core.AssertTrue(t, got)

	got, ok = actionBoolFromAny("true")
	core.AssertTrue(t, ok)
	core.AssertTrue(t, got)

	got, ok = actionBoolFromAny("0")
	core.AssertTrue(t, ok)
	core.AssertFalse(t, got)

	_, ok = actionBoolFromAny("maybe")
	core.AssertFalse(t, ok)
	_, ok = actionBoolFromAny(42)
	core.AssertFalse(t, ok)
}

func TestActionHelpersBehaviour_actionStringFromAny(t *core.T) {
	core.AssertEqual(t, "hello", actionStringFromAny("hello"))
	core.AssertEqual(t, "bytes", actionStringFromAny([]byte("bytes")))
	core.AssertEqual(t, "", actionStringFromAny(nil))
	// Non string/byte values are rendered then trimmed.
	core.AssertEqual(t, "42", actionStringFromAny(42))
}

func TestActionHelpersBehaviour_actionUint32MapValue(t *core.T) {
	values := map[string]any{"rolloutInterval": uint32(5)}
	core.AssertEqual(t, uint32(5), actionUint32MapValue(values, "rollout_interval", "rolloutInterval"))
	// Missing keys yield the zero value.
	core.AssertEqual(t, uint32(0), actionUint32MapValue(values, "absent"))
}

func TestActionHelpersBehaviour_actionBoolMapValue(t *core.T) {
	values := map[string]any{"noRollout": true}
	core.AssertTrue(t, actionBoolMapValue(values, "no_rollout", "noRollout"))
	core.AssertFalse(t, actionBoolMapValue(values, "absent"))
}

func TestActionHelpersBehaviour_hashFromAny_Good(t *core.T) {
	var hash primitives.Hash
	hash[0] = 0xab

	got, ok, err := hashFromAny(hash)
	core.AssertNoError(t, err)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, hash, got)

	// 32-byte slice form.
	got, ok, err = hashFromAny(append([]byte(nil), hash[:]...))
	core.AssertNoError(t, err)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, hash, got)

	// Hex string form.
	hexStr := core.HexEncode(hash[:])
	got, ok, err = hashFromAny(hexStr)
	core.AssertNoError(t, err)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, hash, got)

	// Fixed [32]byte array form.
	var arr [32]byte
	arr[0] = 0xab
	got, ok, err = hashFromAny(arr)
	core.AssertNoError(t, err)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, hash, got)
}

func TestActionHelpersBehaviour_hashFromAny_Bad(t *core.T) {
	// nil yields not-ok with no error.
	_, ok, err := hashFromAny(nil)
	core.AssertFalse(t, ok)
	core.AssertNoError(t, err)

	// Invalid hex string errors.
	_, ok, err = hashFromAny("zz")
	core.AssertFalse(t, ok)
	core.AssertError(t, err)

	// Wrong-length slice that is not valid hex errors.
	_, ok, err = hashFromAny([]byte{0x01, 0x02})
	core.AssertFalse(t, ok)
	core.AssertError(t, err)

	// Unsupported type errors.
	_, ok, err = hashFromAny(42)
	core.AssertFalse(t, ok)
	core.AssertError(t, err)
}

func TestActionHelpersBehaviour_addressHashFromAny(t *core.T) {
	got, ok := addressHashFromAny([]byte{0x01, 0x02})
	core.AssertTrue(t, ok)
	core.AssertLen(t, got, 2)

	got, ok = addressHashFromAny("0102")
	core.AssertTrue(t, ok)
	core.AssertLen(t, got, 2)

	_, ok = addressHashFromAny([]byte{})
	core.AssertFalse(t, ok)
	_, ok = addressHashFromAny("zz")
	core.AssertFalse(t, ok)
	_, ok = addressHashFromAny(42)
	core.AssertFalse(t, ok)
}

func TestActionHelpersBehaviour_addressFromAny(t *core.T) {
	addr := Address{Version: 3, Hash: []byte{0xaa}}
	got, ok := addressFromAny(addr)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, uint8(3), got.Version)

	got, ok = addressFromAny(&addr)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, uint8(3), got.Version)

	_, ok = addressFromAny((*Address)(nil))
	core.AssertFalse(t, ok)

	// Map form.
	got, ok = addressFromAny(map[string]any{"hash": "aabb", "version": uint32(2)})
	core.AssertTrue(t, ok)
	core.AssertEqual(t, uint8(2), got.Version)
	core.AssertLen(t, got.Hash, 2)

	// Map without a valid hash fails.
	_, ok = addressFromAny(map[string]any{"version": uint32(2)})
	core.AssertFalse(t, ok)

	_, ok = addressFromAny(42)
	core.AssertFalse(t, ok)
}

func TestActionHelpersBehaviour_nameRulesFromAny(t *core.T) {
	rules := NameRules{RolloutInterval: 5}
	got, ok := nameRulesFromAny(rules)
	core.AssertTrue(t, ok)
	core.AssertEqual(t, uint32(5), got.RolloutInterval)

	got, ok = nameRulesFromAny(map[string]any{
		"rolloutInterval": uint32(9),
		"noReserved":      true,
	})
	core.AssertTrue(t, ok)
	core.AssertEqual(t, uint32(9), got.RolloutInterval)
	core.AssertTrue(t, got.NoReserved)

	_, ok = nameRulesFromAny(42)
	core.AssertFalse(t, ok)
}

func TestActionHelpersBehaviour_randomActionHash(t *core.T) {
	a, err := randomActionHash()
	core.AssertNoError(t, err)
	b, err := randomActionHash()
	core.AssertNoError(t, err)
	// Two random hashes should differ.
	core.AssertNotEqual(t, a, b)
}

func TestActionHelpersBehaviour_txtLooksVPN(t *core.T) {
	core.AssertTrue(t, txtLooksVPN([]string{"this is a VPN node"}))
	core.AssertTrue(t, txtLooksVPN([]string{"wg=base64key"}))
	core.AssertTrue(t, txtLooksVPN([]string{"endpoint=1.2.3.4:51820"}))
	core.AssertFalse(t, txtLooksVPN([]string{"v=spf1 include:_spf"}))
	core.AssertFalse(t, txtLooksVPN(nil))
}

func TestActionHelpersBehaviour_dnsNameMatches(t *core.T) {
	core.AssertTrue(t, dnsNameMatches("charon.lthn", "CHARON.LTHN."))
	core.AssertFalse(t, dnsNameMatches("charon.lthn", "trade.lthn"))
}
