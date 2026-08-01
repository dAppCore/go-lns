// SPDX-License-Identifier: EUPL-1.2

package blockstore_test

import (
	"fmt"

	"dappco.re/go/store"

	"dappco.re/go/lns/pkg/blockstore"
	"dappco.re/go/lns/pkg/primitives"
)

// exampleStore builds a block store over an in-memory go-store for the examples.
func exampleStore() *blockstore.Store {
	storeResult := store.New(":memory:")
	if err := storeResult.Err(); err != nil {
		panic(err)
	}
	storeInstance := storeResult.Value.(*store.Store)
	scoped := store.NewScoped(storeInstance, "lns-blocks-example")
	bs, err := blockstore.New(scoped)
	if err != nil {
		panic(err)
	}
	return bs
}

// exampleHash builds a deterministic hash from a seed byte for the examples.
func exampleHash(seed byte) primitives.Hash {
	var h primitives.Hash
	for i := range h {
		h[i] = seed
	}
	return h
}

// ExampleStore_PutBlock stores a block and reads it back by hash and by height.
func ExampleStore_PutBlock() {
	bs := exampleStore()
	hash := exampleHash(0x01)

	if err := bs.PutBlock(hash, 100, []byte("block-at-height-100")); err != nil {
		panic(err)
	}

	byHash, _ := bs.GetBlock(hash)
	byHeight, _ := bs.GetBlockByHeight(100)

	fmt.Printf("by hash:   %s\n", byHash)
	fmt.Printf("by height: %s\n", byHeight)
	// Output:
	// by hash:   block-at-height-100
	// by height: block-at-height-100
}

// ExampleStore_Tip shows the tip advancing as blocks are stored.
func ExampleStore_Tip() {
	bs := exampleStore()

	_ = bs.PutBlock(exampleHash(0x00), 0, []byte("genesis"))
	_ = bs.PutBlock(exampleHash(0x01), 1, []byte("one"))
	_ = bs.PutBlock(exampleHash(0x02), 2, []byte("two"))

	_, height, _ := bs.Tip()
	fmt.Printf("tip height: %d\n", height)
	// Output:
	// tip height: 2
}

// ExampleStore_HasBlock shows probing for a block before deciding to refetch.
func ExampleStore_HasBlock() {
	bs := exampleStore()
	hash := exampleHash(0x07)

	before, _ := bs.HasBlock(hash)
	_ = bs.PutBlock(hash, 5, []byte("payload"))
	after, _ := bs.HasBlock(hash)

	fmt.Printf("before: %t, after: %t\n", before, after)
	// Output:
	// before: false, after: true
}

// ExampleStore_DeleteBlock shows disconnecting a block: after deletion it is gone
// from the store and its height is free again.
func ExampleStore_DeleteBlock() {
	bs := exampleStore()
	hash := exampleHash(0x09)

	_ = bs.PutBlock(hash, 9, []byte("to-be-removed"))
	_ = bs.DeleteBlock(hash)

	hasBlock, _ := bs.HasBlock(hash)
	hasHeight, _ := bs.HasHeight(9)
	fmt.Printf("block: %t, height: %t\n", hasBlock, hasHeight)
	// Output:
	// block: false, height: false
}
