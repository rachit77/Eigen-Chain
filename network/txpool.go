package network

import (
	"sort"

	"github.com/rachit77/Eigen-Chain/core"
	"github.com/rachit77/Eigen-Chain/types"
)

//transactions are executed on first come first first serve(FIFO) basis
//unlike ethereum where gas determines the priority  of transaction execution

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

type TxMapSorter struct {
	//make a slice of transactions
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))

	i := 0
	for _, val := range txMap {
		txx[i] = val
		i++
	}

	s := &TxMapSorter{txx}

	sort.Sort(s)
	return s
}

//sort need 3 functions to be implemented i.e. len, less and Swap

func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

//swap i th and j th transactions in slice
func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

// compare i th and j th transactions in slice
//compare on basis of time the transaction is first seen locally
//The transaction first seen is executed first unlike ethereum
func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

//function to return sorted slice of transaction pointer
func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transactions
}

// Add function adds the transaction to the mempool and doesn't check if the transaction is already present
// in the mempool
//the caller should check if transaction is already present in mempool
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if p.Has(hash) {
		return nil
	}

	p.transactions[hash] = tx
	return nil
}

func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash]
	return ok
}

func (p *TxPool) Len() int {
	return len(p.transactions)
}

func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}
