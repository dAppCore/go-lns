// SPDX-License-Identifier: EUPL-1.2

package mining

import (
	core "dappco.re/go"

	"dappco.re/go/lns/pkg/mempool"
	"dappco.re/go/lns/pkg/primitives"
)

// Marshal serialises the whole block to its HSD wire form: the 236-byte header,
// a compact-size transaction count, then each transaction in witness
// serialisation. The result is the opaque byte slice blockstore.PutBlock
// persists.
//
//	raw, err := block.Marshal()
//	_ = bs.PutBlock(hash, height, raw)
func (b Block) Marshal() ([]byte, error) {
	head, err := MarshalHeader(b.Header)
	if err != nil {
		return nil, core.E("mining.Block.Marshal", "marshal header", err)
	}

	buf := make([]byte, 0, len(head)+9+len(b.Txs)*64)
	buf = append(buf, head...)
	buf = appendVarint(buf, uint64(len(b.Txs)))

	for i := range b.Txs {
		txBytes, err := marshalTxWitness(b.Txs[i])
		if err != nil {
			return nil, core.E("mining.Block.Marshal", "marshal transaction", err)
		}
		buf = append(buf, txBytes...)
	}

	return buf, nil
}

// Hash returns the block identifier: the proof-of-work hash of its header. Two
// blocks with the same header share a hash regardless of body, matching HSD,
// where the block hash is the header PoW hash.
//
//	hash, err := block.Hash()
func (b Block) Hash() (primitives.Hash, error) {
	hash, err := PowHash(b.Header)
	if err != nil {
		return primitives.Hash{}, core.E("mining.Block.Hash", "compute pow hash", err)
	}
	return hash, nil
}

// marshalTxWitness serialises a transaction in the HSD witness layout: the base
// (non-witness) data followed by each input's witness stack.
//
//	version(4 LE) ‖ varint(inCount) ‖ inputs(prevout36 ‖ seq4) ‖
//	varint(outCount) ‖ outputs ‖ locktime(4 LE) ‖ witnesses
//
// The base segment is byte-identical to mempool.SerializeNoWitness (which is the
// txid preimage); only the trailing witness stacks are added.
func marshalTxWitness(tx primitives.Transaction) ([]byte, error) {
	base, err := mempool.SerializeNoWitness(tx)
	if err != nil {
		return nil, core.E("mining.marshalTxWitness", "serialise base", err)
	}

	buf := make([]byte, 0, len(base)+32)
	buf = append(buf, base...)
	buf = append(buf, serializeWitnessData(tx)...)
	return buf, nil
}
