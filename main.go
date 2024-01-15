package main

import (
	"fmt"
	"strconv"

	"github.com/facundocarballo/basics-blockchain-go/structures/blockchain"
)

func main() {
	chain := blockchain.Init()

	chain.AddBlock([]byte("First block..."))
	chain.AddBlock([]byte("Second block..."))
	chain.AddBlock([]byte("Third block..."))

	for _, block := range chain.Blocks {
		fmt.Printf("1: 0x%x\n", block.PrevHash)
		fmt.Printf("2: %s\n", block.Data)
		fmt.Printf("3: 0x%x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		fmt.Printf("----\n")

	}
}
