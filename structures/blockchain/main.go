package blockchain

type Blockchain struct {
	Blocks []*Block
}

func Init() *Blockchain {
	return &Blockchain{[]*Block{Genesis()}}
}

func (bc *Blockchain) AddBlock(data []byte) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, new)
}
