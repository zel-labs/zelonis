package transaction

import (
	"zelonis/stats"
	"zelonis/validator/accounts"
	"zelonis/zeldb"
)

type Manager struct {
	db            *zeldb.ZelDB
	accountManger *accounts.Manager
	statsManager  *stats.Manager
}

func NewManager(db *zeldb.ZelDB, account *accounts.Manager, stats *stats.Manager) *Manager {
	return &Manager{
		db:            db,
		accountManger: account,
		statsManager:  stats,
	}
}
