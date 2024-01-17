package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger/v3"
	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

func DBexist() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func Init(address *string) *Blockchain {
	if DBexist() {
		fmt.Printf("Blockchain already exist!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	handlers.HandleErrors(err)

	err = db.Update(func(txn *badger.Txn) error {
		return CreateDatabase(address, &lastHash, txn)
	})

	handlers.HandleErrors(err)

	return &Blockchain{LastHash: lastHash, Database: db}
}

func ContinueBlockchain() *Blockchain {
	if DBexist() == false {
		fmt.Printf("Blockchain is not exist. Create one!!")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	handlers.HandleErrors(err)

	err = db.Update(func(txn *badger.Txn) error {
		return GetLastHash(&lastHash, txn)
	})

	handlers.HandleErrors(err)

	return &Blockchain{LastHash: lastHash, Database: db}
}

func (bc *Blockchain) AddBlock(txs []*Transaction) {
	var lastHash []byte

	bc.Database.View(func(txn *badger.Txn) error {
		return GetLastHash(&lastHash, txn)
	})

	newBlock := CreateBlock(txs, lastHash)

	err := bc.Database.Update(func(txn *badger.Txn) error {
		return SetNewBlock(newBlock, bc, txn)
	})
	handlers.HandleErrors(err)
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.LastHash, bc.Database}
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	bci := bc.Iterator()

	for {
		block := bci.Next()
		if block.IterateTransactions(address, spentTXOs, &unspentTxs) {
			break
		}
	}

	return unspentTxs
}

func (bc *Blockchain) FindUnspentTX(address string) []TxOutput {
	var UTXs []TxOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXs = append(UTXs, out)
			}
		}
	}

	return UTXs
}

func IterateTxOutput(
	address string,
	accumulated *int,
	amount *int,
	tx *Transaction,
	unspentOuts map[string][]int,
) (bool, map[string][]int) {
	txID := hex.EncodeToString(tx.ID)

	for idx, out := range tx.Outputs {
		if out.CanBeUnlocked(address) && *accumulated < *amount {
			*accumulated += out.Value
			unspentOuts[txID] = append(unspentOuts[txID], idx)
			if *accumulated >= *amount {
				return true, unspentOuts
			}
		}
	}

	return false, unspentOuts
}

func (chain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0
	var res bool

	for _, tx := range unspentTxs {
		res, unspentOuts = IterateTxOutput(address, &accumulated, &amount, &tx, unspentOuts)
		if res {
			break
		}
	}

	return accumulated, unspentOuts
}
