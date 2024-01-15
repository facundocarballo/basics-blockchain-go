package blockchain

import (
	"fmt"

	"github.com/facundocarballo/basics-blockchain-go/structures/block"
)

type Blockchain struct {
	blocks []*block.Block
}

func Init() *Blockchain {
	return &Blockchain{[]*block.Block{block.Genesis()}}
}

func (bc *Blockchain) AddBlock(data []byte) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	new := block.CreateBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, new)
}

func DemoBlockchain() {
	chain := Init()

	chain.AddBlock([]byte("Primer bloque..."))
	chain.AddBlock([]byte("Segundo bloque..."))
	chain.AddBlock([]byte("Tercer bloque..."))

	for _, block := range chain.blocks {
		fmt.Printf("----\n")
		fmt.Printf("Previous Hash: 0x%x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: 0x%x\n", block.Hash)
	}
}
