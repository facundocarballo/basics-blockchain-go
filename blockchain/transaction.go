package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/facundocarballo/basics-blockchain-go/handlers"
)

const (
	INITIAL_AMOUNT_OF_CRYPTO = 100
)

type TxInput struct {
	ID        []byte
	Out       int
	Signature string
}

type TxOutput struct {
	Value     int
	PublicKey string
}

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	handlers.HandleErrors(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) IterateOutputs(
	address string,
	spentTxOut map[string][]int,
	unspentTxs *[]Transaction,
) {
	txID := hex.EncodeToString(tx.ID)
	for idx, out := range tx.Outputs {
		if spentTxOut[txID] != nil {
			IterateSpentTxMap(spentTxOut, txID, idx)
		}
		if out.CanBeUnlocked(address) {
			*unspentTxs = append(*unspentTxs, *tx)
		}
	}
}

func (tx *Transaction) IterateInputs(
	address string,
	spentTxOut map[string][]int,
) map[string][]int {
	for _, in := range tx.Inputs {
		if in.CanUnlock(address) {
			inTxID := hex.EncodeToString(in.ID)
			spentTxOut[inTxID] = append(spentTxOut[inTxID], in.Out)
		}
	}
	return spentTxOut
}

func NewTransaction(
	from, to string,
	amount int,
	bc *Blockchain,
) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	accumulator, validOutputs := bc.FindSpendableOutputs(from, amount)

	if accumulator < amount {
		log.Panic("Error: Not enough funds.")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		handlers.HandleErrors(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if accumulator > amount {
		outputs = append(outputs, TxOutput{accumulator - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PublicKey == data
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{[]byte{}, -1, data}
	txout := TxOutput{INITIAL_AMOUNT_OF_CRYPTO, to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()

	return &tx
}

// Can be a helper (each one of this iterators)
func IterateSpentTxMap(spentTxOut map[string][]int, txID string, idx int) {
	for _, spentOut := range spentTxOut[txID] {
		if spentOut == idx {
			break
		}
	}
}
