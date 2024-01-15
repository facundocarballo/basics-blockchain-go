package blockchain

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data []byte, prevHash []byte) *Block {
	block := &Block{[]byte{}, data, prevHash, 0}

	pow := NewProof(block)

	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock([]byte("Genesis"), []byte{})
}