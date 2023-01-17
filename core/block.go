package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
	"crypto/sha256"

	"github.com/rachit77/Eigen-Chain/crypto"
	"github.com/rachit77/Eigen-Chain/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash //Hash of transaction data in block
	PrevBlockHash types.Hash
	Timestamp     uint64
	Height        uint32
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc:= gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

type Block struct {
	*Header
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	//cached version of header hash
	hash types.Hash
}

func(b *Block) AddTransaction(tx *Transaction) {
	b.Transactions= append(b.Transactions,tx)
}

func NewBlock(h *Header, txx []*Transaction) (*Block,error) {
	return &Block{
		Header: h,
		Transactions: txx,
	}, nil
}

func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     uint64(time.Now().UnixNano()) ,
	}

	return NewBlock(header, txx)
}


func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig,err :=privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("no signature")
	}  

	//verify signature of the block
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("invalid block signature")
	}

	//verify all the transactions in the block
	for _,tx := range(b.Transactions) {
		if err:= tx.Verify(); err !=nil {
			return err
		}
	}

	//verify data hash of the block
	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}
	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	}

	//verify prev block hash
	//already checked while adding block (refer validator.go) 

	return nil
}

func(b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func(b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func(b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}

func CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}

	for _, tx := range txx {
		if err = tx.Encode(NewGobTxEncoder(buf)); err != nil {
			return
		}
	}

	hash = sha256.Sum256(buf.Bytes())
	return
}