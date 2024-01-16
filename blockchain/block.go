package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func RegisterBlock() {
	gob.Register(&Block{})
}

func CreateBlock(data []byte, prevHash []byte) *Block {
	block := Block{[]byte{}, data, prevHash, 0}

	pow := NewProof(&block)

	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return &block
}

func Genesis() *Block {
	return CreateBlock([]byte("Genesis"), []byte{})
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