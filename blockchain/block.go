package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func RegisterBlock() {
	gob.Register(&Block{})
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	block := Block{[]byte{}, txs, prevHash, 0}

	pow := NewProof(&block)

	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return &block
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) IterateTransactions(
	address string,
	spentTxOut map[string][]int,
	unspentTxs *[]Transaction,
) bool {
	for _, tx := range b.Transactions {
		tx.IterateOutputs(address, spentTxOut, unspentTxs)
		if tx.IsCoinbase() == false {
			spentTxOut = tx.IterateInputs(address, spentTxOut)
		}
		return len(b.PrevHash) != 0
	}
	return true
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	handlers.HandleErrors(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	handlers.HandleErrors(err)

	return &block
}
