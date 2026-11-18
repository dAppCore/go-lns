// SPDX-License-Identifier: EUPL-1.2

// Package mining implements block production for the Lethean Name System: it
// assembles a candidate block from the mempool, computes the merkle and witness
// roots, and solves and verifies the HSD (Handshake) proof-of-work.
//
// # What it does
//
// Production has three stages. Template assembly drains fee-ordered transactions
// from a mempool.Pool into a Block, capped by the consensus block size and
// weight limits, with the coinbase as the first transaction. Header preparation
// fills the merkle root (over transaction ids), the witness root (over witness
// hashes), and the timestamp, leaving the nonce and extra nonce for the miner.
// Proof-of-work search iterates the nonce against the compact difficulty target
// using the HSD share-hash construction until a header hashes at or below the
// target.
//
//	tmpl, err := mining.NewTemplate(mining.TemplateOptions{
//	    Mempool:   pool,
//	    PrevBlock: tipHash,
//	    Height:    nextHeight,
//	    Bits:      tipBits,
//	    Coinbase:  coinbaseTx,
//	    Time:      uint64(time.Now().Unix()),
//	})
//	blk := tmpl.Block()
//	nonce, ok := mining.Solve(blk.Header, mining.TargetFromBits(blk.Header.Bits), 1<<20)
//	if ok {
//	    blk.Header.Nonce = nonce
//	    raw, _ := blk.Marshal()        // ready for blockstore.PutBlock
//	    hash, _ := mining.PowHash(blk.Header)
//	}
//
// # Proof-of-work scheme
//
// The PoW is the Handshake (HSD 8.99.0) scheme, NOT Bitcoin double-SHA256 and
// NOT the CryptoNote scheme used by go-blockchain. The header is split into a
// preheader (SPV-essential, miner-malleable only via nonce and time) and a
// subheader (miner-mutable). Padding is derived deterministically from the
// previous block and tree root: pad[i] = prevBlock[i%32] XOR treeRoot[i%32].
//
//	subHash    = BLAKE2b-256(extraNonce(24) ‖ reservedRoot ‖ witnessRoot ‖
//	                          merkleRoot ‖ version(4 LE) ‖ bits(4 LE))   // 128 bytes
//	maskHash   = BLAKE2b-256(prevBlock ‖ mask)
//	commitHash = BLAKE2b-256(subHash ‖ maskHash)
//	prehead    = nonce(4 LE) ‖ time(8 LE) ‖ pad(20) ‖ prevBlock ‖ treeRoot ‖
//	             commitHash                                              // 128 bytes
//	left       = BLAKE2b-512(prehead)
//	right      = SHA3-256(prehead ‖ pad(8))
//	shareHash  = BLAKE2b-256(left ‖ pad(32) ‖ right)
//	powHash    = shareHash XOR mask (first 32 bytes)
//
// A header is valid when powHash, read big-endian, is less than or equal to the
// target decoded from the compact bits field.
//
// # Scope
//
// This package implements the deterministic, self-contained core of block
// production: template assembly, the merkle/witness roots, header serialisation,
// the HSD PoW solve and verify, and the compact-target arithmetic. The following
// are deliberately NOT implemented here because they need chain or network
// context this package does not own, or are orchestration concerns the LNS RFC
// (code/core/go/lns/RFC.md §2) does not pin down; they are left to a later,
// separately-specified pass rather than stubbed:
//
//   - Coinbase construction (subsidy schedule, fee accumulation, the coinbase
//     witness commitment, claim/airdrop coinbase inputs). The caller supplies a
//     ready coinbase transaction; assembly validates only that it is structurally
//     a coinbase and lists it first.
//   - Difficulty retargeting (next-block bits from the previous window). The
//     caller supplies the bits to mine against.
//   - The Urkel tree root and reserved root. The caller supplies these via the
//     template; assembly does not recompute the name tree.
//   - Miner orchestration: thread pools, extra-nonce rolling across worker
//     ranges, getblocktemplate/submitblock RPC, and stratum. Solve searches a
//     single nonce range on one goroutine.
//
// # Local Block type
//
// pkg/primitives defines BlockHeader but no assembled Block (header plus the
// transaction vector) and no header or block serialisation. Block production
// needs both, so this package defines a minimal Block here and owns the HSD wire
// encoding of the header and body. If primitives later grows a canonical Block,
// this type should be replaced by it.
package mining

import (
	"github.com/dappcore/go-lns/go/pkg/primitives"
)

// Consensus limits, mirrored from HSD lib/protocol/consensus.js. They bound the
// transactions a template may include.
const (
	// MaxBlockSize is the maximum base (non-witness) block size in bytes.
	MaxBlockSize = 1000000

	// MaxBlockWeight is the maximum block weight, where
	// weight = base*(WitnessScaleFactor-1) + base + witness.
	MaxBlockWeight = 4000000

	// WitnessScaleFactor is the multiplier applied to base size when computing
	// weight, discounting witness bytes.
	WitnessScaleFactor = 4

	// NonceSize is the byte length of the header extra nonce (HSD NONCE_SIZE).
	NonceSize = 24

	// HashSize is the byte length of a block, transaction, or root hash.
	HashSize = 32

	// HeaderSize is the serialised byte length of a block header (HSD
	// HEADER_SIZE): preheader (44) + subheader (160) + mask (32).
	HeaderSize = 236
)

// Block is an assembled candidate block: a header plus its ordered transaction
// vector. The first transaction is the coinbase. It is the unit Solve mines and
// Marshal serialises for blockstore.PutBlock.
//
//	blk := mining.Block{Header: hdr, Txs: []primitives.Transaction{coinbase, tx1}}
type Block struct {
	// Header is the block header carrying the PoW commitment and roots.
	Header primitives.BlockHeader

	// Txs is the transaction vector. Txs[0] is the coinbase; the rest are
	// fee-paying transactions in the order they were selected.
	Txs []primitives.Transaction
}

// TxCount returns the number of transactions in the block, coinbase included.
//
//	n := blk.TxCount()
func (b Block) TxCount() int {
	return len(b.Txs)
}

// Coinbase returns the coinbase transaction (the first transaction) and whether
// the block has one.
//
//	cb, ok := blk.Coinbase()
func (b Block) Coinbase() (primitives.Transaction, bool) {
	if len(b.Txs) == 0 {
		return primitives.Transaction{}, false
	}
	return b.Txs[0], true
}
