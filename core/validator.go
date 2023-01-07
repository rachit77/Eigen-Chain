package core

import "fmt"

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
	if v.bc.HasBlock(b.Height) {
		k := new(BlockHasher)
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", b.Height, k.Hash(b))
	}

	if err := b.Verify(); err != nil {
		return err
	}
	return nil
}
