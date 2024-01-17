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

type CommandLine struct{}

func (cli *CommandLine) help() {
	fmt.Printf("----------- Blockchain Facundo Help -----------\n")
	fmt.Printf("getbalance -address ADDRESS           =>     Get the balance of a particular address.")
	fmt.Printf("createblockchain -address ADDRESS     =>     Creates a Blockchain with this a specif wallet.")

	fmt.Printf("send -from FROM -to TO -amount AMOUNT =>     Add a block to the Blockchain.\n")
	fmt.Printf("print                                 =>     Prints all the blocks in the Blockchain.\n")
	fmt.Printf("----------- Blockchain Facundo Help -----------\n")
}

func (cli *CommandLine) validateArguments() {
	if len(os.Args) < 2 {
		cli.help()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	bc := blockchain.ContinueBlockchain()
	defer bc.Database.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()
		pow := blockchain.NewProof(block)
		fmt.Printf("Previous Hash: 0x%x\n", block.PrevHash)
		fmt.Printf("Current Hash : 0x%x\n", block.Hash)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.Init(&address)
	chain.Database.Close()
	fmt.Printf("[Create Blockchain] Finish.")
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain()
	defer chain.Database.Close()

	balance := 0
	UTXs := chain.FindUnspentTX(address)

	for _, out := range UTXs {
		balance = out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain()
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Printf("[Send Transaction] SUCCESS.")
}

func (cli *CommandLine) Run() {
	cli.validateArguments()

	getBalanceCmd := flag.NewFlagSet(GET_BALANCE_ARGUMENT, flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet(CREATE_BLOCKCHAIN_ARGUMENT, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(SEND_ARGUMENT, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PRINT_ARGUMENT, flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case GET_BALANCE_ARGUMENT:
		err := getBalanceCmd.Parse(os.Args[2:])
		handlers.HandleErrors(err)
	case CREATE_BLOCKCHAIN_ARGUMENT:
		err := createBlockchainCmd.Parse(os.Args[2:])
		handlers.HandleErrors(err)
	case SEND_ARGUMENT:
		err := sendCmd.Parse(os.Args[2:])
		handlers.HandleErrors(err)
	case PRINT_ARGUMENT:
		err := printChainCmd.Parse(os.Args[2:])
		handlers.HandleErrors(err)
	default:
		cli.help()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
