package main

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{NewGenesisBlock()},
	}
}

func NewGenesisBlock() *Block {
	return NewBlock(0, "Genesis Block", "")
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(len(bc.Blocks), data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		if !bc.Blocks[i].IsValid() {
			return false
		}
	}
	return true
}

// GetBlock returns the block at index or nil if out of range.
func (bc *Blockchain) GetBlock(index int) *Block {
	if index < 0 || index >= len(bc.Blocks) {
		return nil
	}
	return bc.Blocks[index]
}
