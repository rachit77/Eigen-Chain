package network

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/rachit77/Eigen-Chain/core"
	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAdTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("foooo"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, p.Len(), 1)

	_ = core.NewTransaction([]byte("foo"))
	assert.Equal(t, p.Len(), 1)

	//test transaction  flush
	p.Flush()
	assert.Equal(t, p.Len(), 0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0; i < 1000; i++ {
		tx := core.NewTransaction([]byte(strconv.Itoa(i)))

		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen, p.Len())

	//sort the transactions
	txx := p.Transactions()

	for i := 0; i < (txLen - 1); i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}

}
