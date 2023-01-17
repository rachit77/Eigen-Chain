package core

import (
	"testing"
	"time"

	"github.com/rachit77/Eigen-Chain/crypto"
	"github.com/rachit77/Eigen-Chain/types"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	privKey := crypto.GenaratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GenaratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GenaratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())
}

//our check will fail if we add transaction in a block after adding the block in the blockchain
//Because all the necesaary checks are done before block is added in blockchain
func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GenaratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	//b, err := NewBlock(header, []*Transaction{tx})
	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)

	//add data hash
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash

	//sign and return the block
	assert.Nil(t, b.Sign(privKey))

	return b
}
