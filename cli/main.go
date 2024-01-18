package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/facundocarballo/basics-blockchain-go/blockchain"
	"github.com/facundocarballo/basics-blockchain-go/handlers"
	"github.com/facundocarballo/basics-blockchain-go/wallet"
)

type CommandLine struct{}

func (cli *CommandLine) help() {
	fmt.Printf("-------------------------------------------- Blockchain Facundo Help --------------------------------------------\n\n")
	fmt.Printf("getbalance -address ADDRESS               =>     Get the balance of a particular address.\n")
	fmt.Printf("createblockchain -address ADDRESS         =>     Creates a Blockchain with this a specif wallet.\n")
	fmt.Printf("send -from FROM -to TO -amount AMOUNT     =>     Add a block to the Blockchain.\n")
	fmt.Printf("print                                     =>     Prints all the blocks in the Blockchain.\n")
	fmt.Printf("createwallet                              =>     Create a new Wallet.\n")
	fmt.Printf("listaddresses                             =>     List the addresses in our wallet file.\n")
	fmt.Printf("\n-------------------------------------------- Blockchain Facundo Help --------------------------------------------\n")
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

func (cli *CommandLine) createWallet() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for idx, address := range addresses {
		fmt.Printf("[%d] => 0x%s\n", idx, address)
	}
}

func (cli *CommandLine) Run() {
	cli.validateArguments()

	getBalanceCmd := flag.NewFlagSet(GET_BALANCE_ARGUMENT, flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet(CREATE_BLOCKCHAIN_ARGUMENT, flag.ExitOnError)
	sendCmd := flag.NewFlagSet(SEND_ARGUMENT, flag.ExitOnError)
	printChainCmd := flag.NewFlagSet(PRINT_ARGUMENT, flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet(CREATE_WALLET_ARGUMENT, flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet(LIST_ADDRESSESS_ARGUMENT, flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	var err error
	switch os.Args[1] {
	case GET_BALANCE_ARGUMENT:
		err = getBalanceCmd.Parse(os.Args[2:])
	case CREATE_BLOCKCHAIN_ARGUMENT:
		err = createBlockchainCmd.Parse(os.Args[2:])
	case SEND_ARGUMENT:
		err = sendCmd.Parse(os.Args[2:])
	case PRINT_ARGUMENT:
		err = printChainCmd.Parse(os.Args[2:])
	case CREATE_WALLET_ARGUMENT:
		err = createWalletCmd.Parse(os.Args[2:])
	case LIST_ADDRESSESS_ARGUMENT:
		err = listAddressesCmd.Parse(os.Args[2:])
	default:
		cli.help()
		runtime.Goexit()
	}
	handlers.HandleErrors(err)

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

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}
}
