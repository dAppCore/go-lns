// SPDX-License-Identifier: EUPL-1.2

package covenant

import (
	"encoding/binary"
	"testing"

	"dappco.re/go/lns/pkg/primitives"
)

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

	if got := CountUpdates(tx); got != 5 {
		t.Fatalf("CountUpdates() = %d, want 5", got)
	}

	if got := CountRenewals(tx); got != 3 {
		t.Fatalf("CountRenewals() = %d, want 3", got)
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
		t.Fatalf("HashString returned error: %v", err)
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
}

func TestNameSetHelpers(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	tx := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], []byte{0, 0, 0, 0}, []byte("example-name")}}},
			{Covenant: primitives.Covenant{Type: uint8(TypeBid), Items: [][]byte{hash[:], []byte{1, 0, 0, 0}, []byte("example-name"), make([]byte, 32)}}},
		},
	}

	set := map[primitives.Hash]struct{}{
		hash: {},
	}

	if !HasNames(tx, set) {
		t.Fatal("HasNames should report a matching name hash")
	}

	RemoveNames(tx, set)

	if len(set) != 0 {
		t.Fatalf("RemoveNames left %d entries, want 0", len(set))
	}

	AddNames(tx, set)

	if len(set) != 1 {
		t.Fatalf("AddNames left %d entries, want 1", len(set))
	}
}

func TestHasSaneCovenants(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	valid := primitives.Transaction{
		Inputs: []primitives.Input{
			{Prevout: primitives.Outpoint{TxHash: primitives.Hash{1}, Index: 0}},
			{Prevout: primitives.Outpoint{TxHash: primitives.Hash{2}, Index: 0}},
		},
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], []byte{0, 0, 0, 0}, []byte("example-name")}}},
			{Covenant: primitives.Covenant{Type: uint8(TypeRegister), Items: [][]byte{hash[:], []byte{1, 0, 0, 0}, []byte("record"), make([]byte, 32)}}},
		},
	}

	if !HasSaneCovenants(valid) {
		t.Fatal("HasSaneCovenants should accept structurally valid covenants")
	}

	invalid := primitives.Transaction{
		Outputs: []primitives.Output{
			{Covenant: primitives.Covenant{Type: uint8(TypeOpen), Items: [][]byte{hash[:], []byte{1, 0, 0, 0}, []byte("example-name")}}},
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
				[]byte("example-name"),
				[]byte{0},
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
					Type: 99,
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
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
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
		[]byte{0},
		committedAddress,
	}

	var finalizeCovenant primitives.Covenant
	finalizeCovenant.Type = uint8(TypeFinalize)
	finalizeCovenant.Items = [][]byte{
		hash[:],
		make([]byte, 4),
		[]byte("example-name"),
		[]byte{0},
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
			primitives.Outpoint{TxHash: prevHash, Index: 0}: primitives.Output{
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
		t.Fatalf("VerifyCovenants returned %d, want 0", got)
	}

	tx.Outputs[0].Address.Hash[0] ^= 1
	if got := VerifyCovenants(tx, view, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for an invalid finalize address, want -1", got)
	}
}

func TestVerifyCovenantsBidReveal(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
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
			[]byte{100, 0, 0, 0},
			[]byte("example-name"),
			bidBlind[:],
		},
	}

	revealCovenant := primitives.Covenant{
		Type: uint8(TypeReveal),
		Items: [][]byte{
			hash[:],
			[]byte{100, 0, 0, 0},
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
		t.Fatalf("VerifyCovenants returned %d, want 0", got)
	}

	tx.Outputs[0].Covenant.Items[2][0] ^= 1
	if got := VerifyCovenants(tx, view, 100, Network{}); got != -1 {
		t.Fatalf("VerifyCovenants returned %d for an invalid reveal nonce, want -1", got)
	}
}

func TestVerifyCovenantsRejectsInvalidStructure(t *testing.T) {
	hash, err := HashString("example-name")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
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
						[]byte{1, 0, 0, 0},
						[]byte("example-name"),
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
	reservedHash, err := HashString("reserved")
	if err != nil {
		t.Fatalf("HashString returned error: %v", err)
	}

	makeClaim := func(blockHeight uint32) primitives.Covenant {
		var blockHash primitives.Hash
		var claimHash primitives.Hash
		copy(claimHash[:], reservedHash[:])

		items := [][]byte{
			claimHash[:],
			make([]byte, 4),
			[]byte("reserved"),
			[]byte{0},
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
					Witness: [][]byte{[]byte("proof")},
				},
			},
			Outputs: []primitives.Output{
				{Covenant: primitives.Covenant{Type: uint8(TypeNone)}},
				{Covenant: makeClaim(100)},
			},
		}

		if got := VerifyCovenants(tx, testCoinView{coins: map[primitives.Outpoint]primitives.Output{}}, 100, Network{}); got != 0 {
			t.Fatalf("VerifyCovenants returned %d, want 0", got)
		}
	})

	t.Run("claim height mismatch", func(t *testing.T) {
		tx := primitives.Transaction{
			Inputs: []primitives.Input{
				{Prevout: primitives.NewOutpoint()},
				{
					Prevout: primitives.Outpoint{TxHash: primitives.Hash{2}, Index: 0},
					Witness: [][]byte{[]byte("proof")},
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
