// SPDX-License-Identifier: EUPL-1.2

package mining_test

import (
	"fmt"

	"github.com/dappcore/go-lns/go/pkg/mempool"
	"github.com/dappcore/go-lns/go/pkg/mining"
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// exampleCoinbase builds a minimal coinbase transaction (null-outpoint input,
// one output) for the examples.
func exampleCoinbase(height uint32, reward uint64) primitives.Transaction {
	return primitives.Transaction{
		Inputs: []primitives.Input{{
			Prevout:  primitives.Outpoint{Index: ^uint32(0)},
			Sequence: ^uint32(0),
			Witness:  [][]byte{{0x00}},
		}},
		Outputs:  []primitives.Output{{Value: reward, Address: primitives.Address{Hash: make([]byte, 20)}}},
		LockTime: height,
	}
}

// exampleSpend builds a distinct non-coinbase transaction seeded by id.
func exampleSpend(id byte, value uint64) primitives.Transaction {
	addr := primitives.Address{Hash: make([]byte, 20)}
	addr.Hash[0] = id
	var prev primitives.Hash
	prev[0] = id
	return primitives.Transaction{
		Inputs:  []primitives.Input{{Prevout: primitives.Outpoint{TxHash: prev}, Sequence: 0xffffffff}},
		Outputs: []primitives.Output{{Value: value, Address: addr}},
	}
}

// ExampleNewTemplate assembles a block template from a mempool and reports how
// many transactions it selected.
func ExampleNewTemplate() {
	pool := mempool.New(mempool.DefaultConfig())

	tmpl, err := mining.NewTemplate(mining.TemplateOptions{
		Mempool:   pool, // empty here, so only the coinbase is included
		Coinbase:  exampleCoinbase(100, 2000),
		PrevBlock: primitives.Hash{},
		Height:    100,
		Bits:      0x2100ffff,
		Time:      1_700_000_000,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	blk := tmpl.Block()
	fmt.Println("transactions:", blk.TxCount())
	fmt.Println("included from mempool:", tmpl.Included)
	// Output:
	// transactions: 1
	// included from mempool: 0
}

// ExampleSolve mines a header at a trivial difficulty and verifies the result.
func ExampleSolve() {
	header := primitives.BlockHeader{
		Version:    1,
		PrevBlock:  primitives.Hash{},
		MerkleRoot: primitives.Hash{},
		Time:       1_700_000_000,
		Bits:       0x2100ffff, // very easy target: solves on the first nonce
	}

	nonce, ok := mining.Solve(header, mining.TargetFromBits(header.Bits), 1<<16)
	if !ok {
		fmt.Println("no solution found")
		return
	}
	header.Nonce = nonce

	valid, err := mining.VerifyPOW(header)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("verified:", valid)
	// Output:
	// verified: true
}

// ExampleTargetFromBits decodes a compact difficulty and re-encodes it.
func ExampleTargetFromBits() {
	bits := uint32(0x1d00ffff)
	target := mining.TargetFromBits(bits)
	fmt.Printf("bit length: %d\n", target.BitLen())
	fmt.Printf("round-trips: %v\n", mining.BitsFromTarget(target) == bits)
	// Output:
	// bit length: 224
	// round-trips: true
}

// ExampleBlock_Marshal mines a tiny block and serialises it for storage.
func ExampleBlock_Marshal() {
	tmpl, err := mining.NewTemplate(mining.TemplateOptions{
		Mempool:  fixedExampleSource{txs: []primitives.Transaction{exampleSpend(0x10, 1000)}},
		Coinbase: exampleCoinbase(100, 2000),
		Bits:     0x2100ffff,
		Time:     1_700_000_000,
	})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	blk := tmpl.Block()
	nonce, _ := mining.Solve(blk.Header, mining.TargetFromBits(blk.Header.Bits), 1<<16)
	blk.Header.Nonce = nonce

	raw, err := blk.Marshal()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("serialised block larger than header:", len(raw) > mining.HeaderSize)
	// Output:
	// serialised block larger than header: true
}

// fixedExampleSource feeds a fixed transaction list to NewTemplate in the
// examples.
type fixedExampleSource struct {
	txs []primitives.Transaction
}

func (s fixedExampleSource) GetForBlock(limit int) []primitives.Transaction {
	return s.txs
}
