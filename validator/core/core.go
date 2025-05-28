package core

import (
	"zelonis/stats"
	"zelonis/validator/accounts"
	"zelonis/validator/block"
	"zelonis/validator/transaction"
)

type Core struct {
	accountManager *accounts.Manager

	txManager    *transaction.Manager
	blockManager *block.Manager
	statsManager *stats.Manager
}

func New(accountManager *accounts.Manager, txManager *transaction.Manager, blockManager *block.Manager, statsManager *stats.Manager) *Core {
	return &Core{
		accountManager: accountManager,
		txManager:      txManager,
		blockManager:   blockManager,
		statsManager:   statsManager,
	}
}

//Build Mempool
