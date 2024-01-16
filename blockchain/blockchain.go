package blockchain

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

const (
	dbPath = "./tmp/blocks"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func Init() *Blockchain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	handlers.HandleErrors(err)

	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(LAST_HASH_KEY))

		if err == badger.ErrKeyNotFound {
			return CreateDatabase(&lastHash, txn)
		}

		return GetLastHash(&lastHash, txn)
	})

	handlers.HandleErrors(err)

	return &Blockchain{LastHash: lastHash, Database: db}
}

func (bc *Blockchain) AddBlock(data []byte) {
	var lastHash []byte

	bc.Database.View(func(txn *badger.Txn) error {
		return GetLastHash(&lastHash, txn)
	})

	newBlock := CreateBlock(data, lastHash)

	err := bc.Database.Update(func(txn *badger.Txn) error {
		return SetNewBlock(newBlock, bc, txn)
	})
	handlers.HandleErrors(err)
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.LastHash, bc.Database}
}
