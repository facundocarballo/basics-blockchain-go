package main

import (
	"os"

	"github.com/facundocarballo/basics-blockchain-go/blockchain"
	"github.com/facundocarballo/basics-blockchain-go/cli"
)

func init() {
	blockchain.RegisterBlock()
}

func main() {
	defer os.Exit(1)

	cli := cli.CommandLine{}
	cli.Run()
}
