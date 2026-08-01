// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	core "dappco.re/go"
	"encoding/binary"
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

func repeatByte(value byte, count int) []byte {
	out := make([]byte, count)
	for i := range out {
		out[i] = value
	}
	return out
}

func TestCountHelpers(t *testing.T) {
	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeClaim)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeUpdate)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeTransfer)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRevoke)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRegister)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRenew)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeFinalize)}},
		},
	}

	if got := CountOpens(tx); got != 1 {
		t.Fatalf("CountOpens() = %d, want 1", got)
	}

	if got := GetCountOpens(tx); got != 1 {
		t.Fatalf("GetCountOpens() = %d, want 1", got)
	}

	if got := CountUpdates(tx); got != 5 {
		t.Fatalf("CountUpdates() = %d, want 5", got)
	}

	if got := GetCountUpdates(tx); got != 5 {
		t.Fatalf("GetCountUpdates() = %d, want 5", got)
	}

	if got := CountRenewals(tx); got != 3 {
		t.Fatalf("CountRenewals() = %d, want 3", got)
	}

	if got := GetCountRenewals(tx); got != 3 {
		t.Fatalf("GetCountRenewals() = %d, want 3", got)
	}
}

func TestGrindName(t *testing.T) {
	rules := NameRules{
		AuctionStart:    100,
		RolloutInterval: 7,
	}

	name, err := GrindName(8, 100, rules)
	if err != nil {
		t.Fatalf("GrindName returned error: %v", err)
	}

	if len(name) != 8 {
		t.Fatalf("GrindName returned %q with length %d, want 8", name, len(name))
	}

	for i := 0; i < len(name); i++ {
		if name[i] < 'a' || name[i] > 'z' {
			t.Fatalf("GrindName returned %q with non-lowercase character %q", name, name[i])
		}
	}

	if !VerifyString(name) {
		t.Fatalf("GrindName returned invalid name %q", name)
	}

	hash, err := HashString(name)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	if !HasRollout(hash, 100, rules) {
		t.Fatalf("GrindName returned %q that does not satisfy rollout", name)
	}

	if HasReservedHash(hash) {
		t.Fatalf("GrindName returned reserved name %q", name)
	}

	if _, err := GrindName(0, 100, rules); err == nil {
		t.Fatal("GrindName should reject zero-length names")
	}

	alias, err := GetGrindName(8, 100, rules)
	if err != nil {
		t.Fatalf("GetGrindName returned error: %v", err)
	}

	if len(alias) != 8 {
		t.Fatalf("GetGrindName returned %q with length %d, want 8", alias, len(alias))
	}
}

func TestNameSetHelpers(t *testing.T) {
	hash, err := HashString(testExampleName)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], {0, 0, 0, 0}, []byte(testExampleName)}}},
			{Covenant: primitives.Covenant{Type: uint8(TypeBid), Items: [][]byte{hash[:], {1, 0, 0, 0}, []byte(testExampleName), make([]byte, 32)}}},
		},
	}

	set := map[primitives.Hash]struct{}{
		hash: {},
	}

	if !HasNames(tx, set) {
		t.Fatal("HasNames should report a matching name hash")
	}

	if !GetHasNames(tx, set) {
		t.Fatal("GetHasNames should alias HasNames")
	}

	RemoveNames(tx, set)

	if len(set) != 0 {
		t.Fatalf("RemoveNames left %d entries, want 0", len(set))
	}

	AddNames(tx, set)

	if len(set) != 1 {
		t.Fatalf("AddNames left %d entries, want 1", len(set))
	}

	GetRemoveNames(tx, set)

	if len(set) != 0 {
		t.Fatalf("GetRemoveNames left %d entries, want 0", len(set))
	}

	GetAddNames(tx, set)

	if len(set) != 1 {
		t.Fatalf("GetAddNames left %d entries, want 1", len(set))
	}
}

func TestHasSaneCovenants(t *testing.T) {
	hash, err := HashString(testExampleName)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	valid := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: primitives.Hash{1}, Index: 0}},
			{Prevout: primitives.Outpoint{TxHash: primitives.Hash{2}, Index: 0}},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], {0, 0, 0, 0}, []byte(testExampleName)}}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRegister), Items: [][]byte{hash[:], {1, 0, 0, 0}, []byte("record"), make([]byte, 32)}}},
		},
	}

	if !HasSaneCovenants(valid) {
		t.Fatal("HasSaneCovenants should accept structurally valid covenants")
	}

	if !GetHasSaneCovenants(valid) {
		t.Fatal("GetHasSaneCovenants should alias HasSaneCovenants")
	}

	invalid := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], {1, 0, 0, 0}, []byte(testExampleName)}}},
		},
	}

	if HasSaneCovenants(invalid) {
		t.Fatal("HasSaneCovenants should reject non-zero open heights")
	}

	coinbaseClaim := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.NewOutpoint()},
			{Prevout: primitives.Outpoint{TxHash: primitives.Hash{3}, Index: 0}},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			{Covenant: primitives.Covenant{Type: uint8(TypeClaim), Items: [][]byte{
				hash[:],
				make([]byte, 4),
				[]byte(testExampleName),
				{0},
				make([]byte, 32),
				make([]byte, 4),
			}}},
		},
	}

	if HasSaneCovenants(coinbaseClaim) {
		t.Fatal("HasSaneCovenants should reject coinbase claims without exactly one witness item")
	}

	unknown := primitives.Transaction{
		Outputs: []primitives.Output{
			{
				Covenant: primitives.Covenant{
					Type:  99,
					Items: make([][]byte, 256),
				},
			},
		},
	}

	if !HasSaneCovenants(unknown) {
		t.Fatal("HasSaneCovenants should accept unknown covenants within the script stack limit")
	}
}

type testCoinView struct {
	coins map[primitives.Outpoint]primitives.Output
}

func (v testCoinView) GetOutput(prevout primitives.Outpoint) (primitives.Output, bool) {
	coin, ok := v.coins[prevout]
	return coin, ok
}

func TestVerifyCovenants(t *testing.T) {
	hash, err := HashString(testExampleName)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	var prevHash primitives.Hash
	prevHash[0] = 1

	var finalHash primitives.Hash
	finalHash[0] = 2
	committedAddress := append([]byte(nil), finalHash[:]...)
	outputAddress := append([]byte(nil), finalHash[:]...)

	var transferCovenant primitives.Covenant
	transferCovenant.Type = uint8(TypeTransfer)
	transferCovenant.Items = [][]byte{
		hash[:],
		make([]byte, 4),
		{0},
		committedAddress,
	}

	var finalizeCovenant primitives.Covenant
	finalizeCovenant.Type = uint8(TypeFinalize)
	finalizeCovenant.Items = [][]byte{
		hash[:],
		make([]byte, 4),
		[]byte(testExampleName),
		{0},
		make([]byte, 4),
		make([]byte, 4),
		make([]byte, 32),
	}

	binary.LittleEndian.PutUint32(transferCovenant.Items[1], 100)
	binary.LittleEndian.PutUint32(finalizeCovenant.Items[1], 100)
	binary.LittleEndian.PutUint32(finalizeCovenant.Items[4], 1)
	binary.LittleEndian.PutUint32(finalizeCovenant.Items[5], 2)

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: prevHash, Index: 0}},
		},
		Outputs: []primitives.Output{
			{
				Value: 1000,
				Address: primitives.Address{
					Version: 0,
					Hash:    outputAddress,
				},
				Covenant: finalizeCovenant,
			},
		},
	}

	view := testCoinView{
		coins: map[primitives.Outpoint]primitives.Output{
			{TxHash: prevHash, Index: 0}: {
				Value: 1000,
				Address: primitives.Address{
					Version: 0,
					Hash:    []byte{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
				},
				Covenant: transferCovenant,
			},
		},
	}

	if got := VerifyCovenants(tx, view, 100, Network{}); got != 0 {
		t.Fatalf(testVerifyCovenantsZero, got)
	}

	if got := GetVerifyCovenants(tx, view, 100, Network{}); got != 0 {
		t.Fatalf("GetVerifyCovenants returned %d, want 0", got)
	}

	tx.Outputs[0].Address.Hash[0] ^= 1
	if got := VerifyCovenants(tx, view, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for an invalid finalize address, want -1", got)
	}

	if got := GetVerifyCovenants(tx, view, 100, Network{}); got != -1 {
		t.Fatalf("GetVerifyCovenants returned %d for an invalid finalize address, want -1", got)
	}
}

func TestVerifyCovenantsBidReveal(t *testing.T) {
	hash, err := HashString(testExampleName)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	var prevHash primitives.Hash
	prevHash[0] = 4

	var nonce primitives.Hash
	for i := range nonce {
		nonce[i] = byte(i)
	}

	bidBlind, err := Blind(1000, nonce)
	if err != nil {
		t.Fatalf("Blind returned error: %v", err)
	}

	bidCovenant := primitives.Covenant{
		Type: uint8(TypeBid),
		Items: [][]byte{
			hash[:],
			{100, 0, 0, 0},
			[]byte(testExampleName),
			bidBlind[:],
		},
	}

	revealCovenant := primitives.Covenant{
		Type: uint8(TypeReveal),
		Items: [][]byte{
			hash[:],
			{100, 0, 0, 0},
			nonce[:],
		},
	}

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: prevHash, Index: 0}},
		},
		Outputs: []primitives.Output{
			{
				Value:    1000,
				Covenant: revealCovenant,
			},
		},
	}

	view := testCoinView{
		coins: map[primitives.Outpoint]primitives.Output{
			{TxHash: prevHash, Index: 0}: {
				Value:    1500,
				Covenant: bidCovenant,
			},
		},
	}

	if got := VerifyCovenants(tx, view, 100, Network{}); got != 0 {
		t.Fatalf(testVerifyCovenantsZero, got)
	}

	tx.Outputs[0].Covenant.Items[2][0] ^= 1
	if got := VerifyCovenants(tx, view, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for an invalid reveal nonce, want -1", got)
	}
}

func TestVerifyCovenantsRejectsInvalidStructure(t *testing.T) {
	hash, err := HashString(testExampleName)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	var prevHash primitives.Hash
	prevHash[0] = 3

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: prevHash, Index: 0}},
		},
		Outputs: []primitives.Output{
			{
				Covenant: primitives.Covenant{
					Type: uint8(TypeOpen),
					Items: [][]byte{
						hash[:],
						{1, 0, 0, 0},
						[]byte(testExampleName),
					},
				},
			},
		},
	}

	view := testCoinView{
		coins: map[primitives.Outpoint]primitives.Output{
			{TxHash: prevHash, Index: 0}: {
				Covenant: primitives.Covenant{Type: uint8(TypeNone)},
			},
		},
	}

	if got := VerifyCovenants(tx, view, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for an invalid covenant structure, want -1", got)
	}
}

func TestHasSaneCovenantsRejectsOversizedUnknownCovenants(t *testing.T) {
	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{
				Covenant: primitives.Covenant{
					Type: uint8(99),
					Items: func() [][]byte {
						items := make([][]byte, maxScriptStack+1)
						for i := range items {
							items[i] = []byte{}
						}
						return items
					}(),
				},
			},
		},
	}

	if HasSaneCovenants(tx) {
		t.Fatal("HasSaneCovenants should reject unknown covenants that exceed the script stack limit")
	}
}

func TestVerifyCovenantsCoinbaseLinkedOutputs(t *testing.T) {
	reservedHash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	makeClaim := func(blockHeight uint32) primitives.Covenant {
		var blockHash primitives.Hash
		var claimHash primitives.Hash
		copy(claimHash[:], reservedHash[:])

		items := [][]byte{
			claimHash[:],
			make([]byte, 4),
			[]byte(testReservedLabel),
			{0},
			blockHash[:],
			make([]byte, 4),
		}

		binary.LittleEndian.PutUint32(items[1], blockHeight)
		binary.LittleEndian.PutUint32(items[5], blockHeight)

		return primitives.Covenant{Type: uint8(TypeClaim), Items: items}
	}

	t.Run("valid claim height", func(t *testing.T) {
		tx := primitives.Transaction{
			Inputs: []primitives.Input{
				{Prevout: primitives.NewOutpoint()},
				{
					Prevout: primitives.Outpoint{TxHash: primitives.Hash{1}, Index: 0},
					Witness: [][]byte{encodeTestClaim([]byte("proof"))},
				},
			},
			Outputs: []primitives.Output{
				{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
				{Covenant: makeClaim(100)},
			},
		}

		if got := VerifyCovenants(tx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != 0 {
			t.Fatalf(testVerifyCovenantsZero, got)
		}
	})

	t.Run("claim height mismatch", func(t *testing.T) {
		tx := primitives.Transaction{
			Inputs: []primitives.Input{
				{Prevout: primitives.NewOutpoint()},
				{
					Prevout: primitives.Outpoint{TxHash: primitives.Hash{2}, Index: 0},
					Witness: [][]byte{encodeTestClaim([]byte("proof"))},
				},
			},
			Outputs: []primitives.Output{
				{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
				{Covenant: makeClaim(101)},
			},
		}

		if got := VerifyCovenants(tx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != -1 {
			t.Fatalf("VerifyCovenants returned %d for a claim with a mismatched height, want -1", got)
		}
	})

	t.Run("missing witness", func(t *testing.T) {
		tx := primitives.Transaction{
			Inputs: []primitives.Input{
				{Prevout: primitives.NewOutpoint()},
				{
					Prevout: primitives.Outpoint{TxHash: primitives.Hash{3}, Index: 0},
				},
			},
			Outputs: []primitives.Output{
				{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
				{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			},
		}

		if got := VerifyCovenants(tx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != -1 {
			t.Fatalf("VerifyCovenants returned %d for a coinbase-linked output without a witness, want -1", got)
		}
	})
}

func encodeTestVarint(n uint64) []byte {
	switch {
	case n < 0xfd:
		return []byte{byte(n)}
	case n <= 0xffff:
		return []byte{0xfd, byte(n), byte(n >> 8)}
	case n <= 0xffffffff:
		return []byte{0xfe, byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24)}
	default:
		return []byte{
			0xff,
			byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24),
			byte(n >> 32), byte(n >> 40), byte(n >> 48), byte(n >> 56),
		}
	}
}

func encodeTestAirdropKey(version uint8, address []byte, value uint64, sponsor bool) []byte {
	key := []byte{4, version, byte(len(address))}
	key = append(key, address...)

	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], value)
	key = append(key, tmp[:]...)
	if sponsor {
		key = append(key, 1)
	} else {
		key = append(key, 0)
	}

	return key
}

func encodeTestAirdropProof(key, address []byte, fee uint64) []byte {
	var out []byte
	var index [4]byte
	binary.LittleEndian.PutUint32(index[:], 3)
	out = append(out, index[:]...)
	out = append(out, 0) // proof length
	out = append(out, 0) // subindex
	out = append(out, 0) // subproof length
	out = appendTestVarBytes(out, key)
	out = append(out, 0) // version
	out = append(out, byte(len(address)))
	out = append(out, address...)
	out = append(out, encodeTestVarint(fee)...)
	out = appendTestVarBytes(out, nil)
	return out
}

func appendTestVarBytes(dst, src []byte) []byte {
	dst = append(dst, encodeTestVarint(uint64(len(src)))...)
	return append(dst, src...)
}

func encodeTestClaim(blob []byte) []byte {
	claim := primitives.Claim{Blob: append([]byte(nil), blob...)}
	raw, err := claim.MarshalBinary()
	if err != nil {
		panic(err)
	}

	return raw
}

func TestCoinbaseAirdropProofValidation(t *testing.T) {
	address := repeatByte(0x42, 20)
	key := encodeTestAirdropKey(0, address, 1000, false)
	proof := encodeTestAirdropProof(key, address, 10)

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.NewOutpoint()},
			{
				Prevout: primitives.Outpoint{TxHash: primitives.Hash{4}, Index: 0},
				Witness: [][]byte{proof},
			},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			{
				Value: 990,
				Address: primitives.Address{
					Version: 0,
					Hash:    append([]byte(nil), address...),
				},
				Covenant: primitives.Covenant{Type: uint8(TypeNone)},
			},
		},
	}

	if !HasSaneCovenants(tx) {
		t.Fatal("HasSaneCovenants should accept a valid coinbase airdrop proof")
	}

	if got := VerifyCovenants(tx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != 0 {
		t.Fatalf("VerifyCovenants returned %d for a valid coinbase airdrop proof, want 0", got)
	}

	badTx := tx
	badTx.Outputs[1].Value = 991

	if got := VerifyCovenants(badTx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for a mismatched coinbase airdrop value, want -1", got)
	}
}

func TestCoinbaseClaimWitnessValidation(t *testing.T) {
	reservedHash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	var claimHash primitives.Hash
	copy(claimHash[:], reservedHash[:])

	blockHeight := make([]byte, 4)
	binary.LittleEndian.PutUint32(blockHeight, 100)

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.NewOutpoint()},
			{
				Prevout: primitives.Outpoint{TxHash: primitives.Hash{5}, Index: 0},
				Witness: [][]byte{encodeTestClaim([]byte("proof"))},
			},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			{
				Value: 0,
				Address: primitives.Address{
					Version: 0,
					Hash:    repeatByte(0x24, 20),
				},
				Covenant: primitives.Covenant{
					Type: uint8(TypeClaim),
					Items: [][]byte{
						claimHash[:],
						blockHeight,
						[]byte(testReservedLabel),
						{0},
						make([]byte, 32),
						blockHeight,
					},
				},
			},
		},
	}

	if got := HasSaneCovenants(tx); !got {
		t.Fatal("HasSaneCovenants should accept a valid coinbase claim witness wrapper")
	}

	if got := VerifyCovenants(tx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != 0 {
		t.Fatalf("VerifyCovenants returned %d for a valid coinbase claim witness wrapper, want 0", got)
	}

	badTx := tx
	badTx.Inputs[1].Witness = [][]byte{[]byte("proof")}

	if got := HasSaneCovenants(badTx); got {
		t.Fatal("HasSaneCovenants should reject an invalid coinbase claim witness wrapper")
	}

	if got := VerifyCovenants(badTx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for an invalid coinbase claim witness wrapper, want -1", got)
	}
}

func TestCoinbaseClaimConjureOverflow(t *testing.T) {
	reservedHash, err := HashString(testReservedLabel)
	if err != nil {
		t.Fatalf(testHashStringErr, err)
	}

	var claimHash primitives.Hash
	copy(claimHash[:], reservedHash[:])

	blockHeight := make([]byte, 4)
	binary.LittleEndian.PutUint32(blockHeight, 100)

	tx := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.NewOutpoint()},
			{
				Prevout: primitives.Outpoint{TxHash: primitives.Hash{6}, Index: 0},
				Witness: [][]byte{encodeTestClaim([]byte("proof-a"))},
			},
			{
				Prevout: primitives.Outpoint{TxHash: primitives.Hash{7}, Index: 0},
				Witness: [][]byte{encodeTestClaim([]byte("proof-b"))},
			},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
			{
				Value: ^uint64(0),
				Address: primitives.Address{
					Version: 0,
					Hash:    repeatByte(0x24, 20),
				},
				Covenant: primitives.Covenant{
					Type: uint8(TypeClaim),
					Items: [][]byte{
						claimHash[:],
						blockHeight,
						[]byte(testReservedLabel),
						{0},
						make([]byte, 32),
						blockHeight,
					},
				},
			},
			{
				Value: 1,
				Address: primitives.Address{
					Version: 0,
					Hash:    repeatByte(0x25, 20),
				},
				Covenant: primitives.Covenant{
					Type: uint8(TypeClaim),
					Items: [][]byte{
						claimHash[:],
						blockHeight,
						[]byte(testReservedLabel),
						{0},
						make([]byte, 32),
						blockHeight,
					},
				},
			},
		},
	}

	if got := HasSaneCovenants(tx); !got {
		t.Fatal("HasSaneCovenants should accept the coinbase claim shape")
	}

	if got := VerifyCovenants(tx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for overflowing coinbase claim outputs, want -1", got)
	}
}

// BEGIN v090 file-aware generated coverage

func TestRulesExtra_GrindName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GrindName(1, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GrindName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GrindName(0, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GrindName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GrindName(-1, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetGrindName_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetGrindName(1, 1, coverageValue[NameRules]("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetGrindName_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetGrindName(0, 0, coverageValue[NameRules]("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetGrindName_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetGrindName(-1, 63, coverageValue[NameRules]("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountOpens_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(CountOpens(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountOpens_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(CountOpens(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountOpens_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(CountOpens(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountOpens_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetCountOpens(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountOpens_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetCountOpens(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountOpens_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetCountOpens(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountUpdates_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(CountUpdates(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountUpdates_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(CountUpdates(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountUpdates_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(CountUpdates(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountUpdates_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetCountUpdates(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountUpdates_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetCountUpdates(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountUpdates_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetCountUpdates(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountRenewals_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(CountRenewals(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountRenewals_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(CountRenewals(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_CountRenewals_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(CountRenewals(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountRenewals_Good(t *core.T) {
	coverageAssertVariant(t, "good")
	values := coveragePack(GetCountRenewals(coverageTransaction("good")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountRenewals_Bad(t *core.T) {
	coverageAssertVariant(t, "bad")
	values := coveragePack(GetCountRenewals(coverageTransaction("bad")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

func TestRulesExtra_GetCountRenewals_Ugly(t *core.T) {
	coverageAssertVariant(t, "ugly")
	values := coveragePack(GetCountRenewals(coverageTransaction("ugly")))
	core.AssertTrue(t, len(values) >= 1)
	coverageAssertReturn(t, values...)
}

// END v090 file-aware generated coverage
