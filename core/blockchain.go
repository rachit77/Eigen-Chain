package core

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type Blockchain struct {
	store     Storage
	lock      sync.RWMutex
	headers   []*Header
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(),
	}

	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err

}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

//validate and add the block in the chain
func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {

	// bc.lock.Lock()
	// defer bc.lock.Unlock()
	// bc.lock.RLock()
	// defer bc.lock.RUnlock()

	if height > bc.Height() {
		return nil, fmt.Errorf("given height %d is not a valid block number aka height", height)
	}

	// a read lock can't be aquired when a write lock is already aquired
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[int(height)], nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

//function to retrive height of the chain
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers) - 1)
}

// make this function local and not exposed outside
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash":   BlockHasher{}.Hash(b.Header),
	}).Info("adding new block")

	return bc.store.Put(b)
}
