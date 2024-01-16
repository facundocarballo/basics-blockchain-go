package main

import (
	"github.com/facundocarballo/basics-blockchain-go/blockchain"
	"github.com/facundocarballo/basics-blockchain-go/cli"
)

func init() {
	blockchain.RegisterBlock()
}

func main() {
	chain := blockchain.Init()
	defer chain.Database.Close()

	commandLine := &cli.CommandLine{Blockchain: chain}
	commandLine.Run()
}
