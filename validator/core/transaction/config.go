package transaction

import (
	"zelonis/stats"
	"zelonis/validator/accounts"
	"zelonis/validator/core/transaction/mempool"
	"zelonis/zeldb"
)

type Manager struct {
	db            *zeldb.ZelDB
	accountManger *accounts.Manager
	statsManager  *stats.Manager
	mempool       *mempool.TransactionsPool
}

func NewManager(db *zeldb.ZelDB, account *accounts.Manager, stats *stats.Manager) *Manager {
	tempMempool := mempool.NewMempool()

	return &Manager{
		db:            db,
		accountManger: account,
		statsManager:  stats,
		mempool:       tempMempool.NewTransactionsPool(),
	}
}
