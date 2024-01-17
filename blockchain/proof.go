package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

func Difficulty() uint {
	return 16
}

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.HashTransactions(),
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty())),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r0x%x", hash)
		intHash.SetBytes(hash[:])

		// If the Cmp func returns -1, it means that intHash is less than pow.Target => Block is sign.
		// Because when we apply the difficulty, brings X amounts of 0 to the left.
		// So to make sure that this new hash have at least X zeros at the left, this hash has to be less than the target.
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var initHash big.Int

	data := pow.InitData(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	initHash.SetBytes(hash[:])

	return initHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func NewProof(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// target -> 0x0000000000000000000000000000000000000000000000000000000000000001
	target.Lsh(target, uint(256-Difficulty()))
	// target -> 0x0000010000000000000000000000000000000000000000000000000000000000
	return &ProofOfWork{block, target}
}
