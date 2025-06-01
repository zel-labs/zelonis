package accounts

import (
	"zelonis/stats"
	"zelonis/zeldb"
)

type Manager struct {
	db           *zeldb.ZelDB
	statsManager *stats.Manager
}

func NewManager(db *zeldb.ZelDB, statsManager *stats.Manager) *Manager {
	return &Manager{
		db:           db,
		statsManager: statsManager,
	}
}
