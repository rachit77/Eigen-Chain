package core

import (
	"testing"

	"github.com/rachit77/Eigen-Chain/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GenaratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)

}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GenaratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Equal(t, 1, 1)
	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GenaratePrivateKey()
	tx.From = otherPrivKey.PublicKey()
	assert.NotNil(t, tx.Verify())
}

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GenaratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privKey))
	return tx
}
