package mempool

import (
	"log"
	"os"
	"sync"
	"zelonis/external"
)

type mempool struct {
	newStaketx bool
	txMtx      sync.Mutex
	mtx        sync.RWMutex

	stagingUtxo uint64

	rejectedTx uint64
	acceptedTx uint64
}

type TransactionsPool struct {
	mempool         *mempool
	allTransactions IDToTransactionMap
}

type IDToTransactionMap map[external.DomainTransactionID]*MempoolTransaction

type MempoolTransaction struct {
	transaction *external.Transaction
}

func NewMempool() *mempool {

	return &mempool{
		newStaketx: true,
		txMtx:      sync.Mutex{},
		mtx:        sync.RWMutex{},
	}

}

func (mp *mempool) NewTransactionsPool() *TransactionsPool {
	return &TransactionsPool{
		mempool:         mp,
		allTransactions: IDToTransactionMap{},
	}
}

func (tp *TransactionsPool) AddTxToMempool(tx *external.Transaction) bool {

	tp.mempool.mtx.Lock()
	defer tp.mempool.mtx.Unlock()

	hash, err := external.NewDomainTransactionIDFromByteSlice(tx.TxHash)
	if err != nil {
		panic(err)
	}
	if (tp.allTransactions)[*hash] != nil {
		return false
	}

	mtx := newMempoolTransaction(tx)
	(tp.allTransactions)[*hash] = mtx

	return true
}

func (tp *TransactionsPool) RemoveTxFromMempool(tx *external.Transaction) bool {
	tp.mempool.mtx.Lock()
	defer tp.mempool.mtx.Unlock()
	hash, err := external.NewDomainTransactionIDFromByteSlice(tx.TxHash)
	if err != nil {
		panic(err)
	}

	delete(tp.allTransactions, *hash)
	if len(tp.allTransactions) != 0 {
		log.Println(len(tp.allTransactions))
		os.Exit(112)
	}

	return true
}

func newMempoolTransaction(transaction *external.Transaction) *MempoolTransaction {
	return &MempoolTransaction{
		transaction: transaction,
	}
}

func (tp *TransactionsPool) GetMempoolTxs() []*external.Transaction {
	txs := make([]*external.Transaction, 0)

	for _, val := range tp.allTransactions {
		txs = append(txs, val.transaction)
	}
	return txs
}
