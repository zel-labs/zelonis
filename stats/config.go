package stats

import (
	"fmt"
	"strconv"
	"zelonis/zeldb"
)

type Manager struct {
	db *zeldb.ZelDB
}

func NewManager(db *zeldb.ZelDB) *Manager {
	return &Manager{
		db: db,
	}
}

func (m *Manager) GetHighestBlockHeight() (uint64, error) {
	byteInfo, err := m.db.Get([]byte(heightestBlockHeightKey))
	if err != nil {
		return 0, err
	}
	info := fmt.Sprintf("%s", byteInfo)
	infoHeight, err := strconv.ParseUint(info, 10, 64)
	if err != nil {
		return 0, err
	}
	return infoHeight, nil
}

func (m *Manager) UpdateHighestBlockHeight(blockHeight uint64) error {
	info := fmt.Sprintf("%v", blockHeight)
	infoBytes := []byte(info)
	err := m.db.Set([]byte(heightestBlockHeightKey), infoBytes)
	if err != nil {
		return err
	}
	return nil
}
