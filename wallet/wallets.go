package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)

	err := encoder.Encode(ws)
	handlers.HandleErrors(err)

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0544)
	handlers.HandleErrors(err)
}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)
	handlers.HandleErrors(err)

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	err = decoder.Decode(&wallets)
	handlers.HandleErrors(err)

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws *Wallets) GetWallets(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAllAddresses() []string {
	var addressess []string

	for address := range ws.Wallets {
		addressess = append(addressess, address)
	}

	return addressess
}

func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}
