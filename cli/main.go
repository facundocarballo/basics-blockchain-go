package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/facundocarballo/basics-blockchain-go/blockchain"
	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

type CommandLine struct {
	Blockchain *blockchain.Blockchain
}

func (cli *CommandLine) help() {
	fmt.Printf("----- Blockchain Facundo Help -----\n")
	fmt.Printf("add -block BLOCK_DATA     =>     Add a block to the Blockchain.\n")
	fmt.Printf("print                     =>     Prints all the blocks in the Blockchain.\n")
	fmt.Printf("----- Blockchain Facundo Help -----\n")
}

func (cli *CommandLine) validateArguments() {
	if len(os.Args) < 2 {
		cli.help()
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.Blockchain.AddBlock([]byte(data))
	fmt.Printf("Block addedd :)\n")
}

func (cli *CommandLine) printChain() {
	bci := cli.Blockchain.Iterator()
	for {
		block := bci.Next()
		pow := blockchain.NewProof(block)
		fmt.Printf("Previous Hash: 0x%x\n", block.PrevHash)
		fmt.Printf("Current Hash : 0x%x\n", block.Hash)
		fmt.Printf("Data         : %s\n", block.Data)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.validateArguments()

	addBlockCmd := flag.NewFlagSet(ADD_ARGUMENT, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PRINT_ARGUMENT, flag.ExitOnError)
	addBlockData := addBlockCmd.String(BLOCK_ARGUMENT, "", "Block data")

	switch os.Args[1] {
	case ADD_ARGUMENT:
		err := addBlockCmd.Parse(os.Args[2:])
		handlers.HandleErrors(err)
	case PRINT_ARGUMENT:
		err := printChainCmd.Parse(os.Args[2:])
		handlers.HandleErrors(err)
	default:
		cli.help()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
