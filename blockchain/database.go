package blockchain

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

const (
	LAST_HASH_KEY = "./tmp/blocks"
)

// Creates
func CreateDatabase(address *string, lastHash *[]byte, txn *badger.Txn) error {
	coinbaseTx := CoinbaseTx(*address, genesisData)
	genesis := Genesis(coinbaseTx)
	err := txn.Set(genesis.Hash, genesis.Serialize())
	handlers.HandleErrors(err)
	err = txn.Set([]byte(LAST_HASH_KEY), genesis.Hash)
	*lastHash = genesis.Hash
	return err
}

// Sets
func SetNewBlock(block *Block, bc *Blockchain, txn *badger.Txn) error {
	err := txn.Set(block.Hash, block.Serialize())
	handlers.HandleErrors(err)

	err = txn.Set([]byte(LAST_HASH_KEY), block.Hash)

	bc.LastHash = block.Hash

	return err
}

// Gets
func GetLastHash(lastHash *[]byte, txn *badger.Txn) error {
	item, err := txn.Get([]byte(LAST_HASH_KEY))
	handlers.HandleErrors(err)
	item.Value(func(val []byte) error {
		*lastHash = val
		return nil
	})
	return nil
}

func GetBlock(block *Block, txn *badger.Txn, currentHash []byte) error {
	item, err := txn.Get(currentHash)
	handlers.HandleErrors(err)
	item.Value(func(val []byte) error {
		*block = *Deserialize(val)
		return nil
	})
	return nil
}
