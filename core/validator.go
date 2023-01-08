package core

import (
	"fmt"
)

type Validator interface {
	ValidateBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *Block) error {
	//check if that block number is already present in the blockchain
	if v.bc.HasBlock(b.Height) {
		k := new(BlockHasher)
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", b.Height, k.Hash(b.Header))
	}

	//check the block number of new block
	if b.Height != v.bc.Height()+1 {
		k := new(BlockHasher)
		return fmt.Errorf("block %s too high", k.Hash(b.Header))
	}

	//validate prevHash of current block
	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block %s is invalid", b.PrevBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}
	return nil
}
