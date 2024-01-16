package blockchain

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (bci *BlockchainIterator) Next() *Block {
	var block Block

	err := bci.Database.View(func(txn *badger.Txn) error {
		return GetBlock(&block, txn, bci.CurrentHash)
	})
	handlers.HandleErrors(err)
	bci.CurrentHash = block.PrevHash

	return &block
}
