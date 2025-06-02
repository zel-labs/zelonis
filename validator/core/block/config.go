/*
Copyright (C) 2025 Zelonis Contributors

This file is part of Zelonis.

Zelonis is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Zelonis is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Zelonis. If not, see <https://www.gnu.org/licenses/>.
*/

package block

import (
	"fmt"
	"strconv"
	"zelonis/external"
	"zelonis/stats"
	"zelonis/zeldb"
)

type Manager struct {
	db                *zeldb.ZelDB
	latestBlockHeight uint64
	latestBlockHash   []byte
	previousBlockHash []byte
	statsManager      *stats.Manager
}

func NewManager(db *zeldb.ZelDB, statsManager *stats.Manager) *Manager {
	return &Manager{
		db:           db,
		statsManager: statsManager,
	}
}

func (m *Manager) VerifyAndAddBlock(block *external.Block) (bool, error) {
	blockBytes := block.Serialize()

	alreadyExists, err := m.db.Has(block.Header.BlockHeightBytes())
	if err != nil {

		return alreadyExists, err
	}
	if alreadyExists {

		return alreadyExists, nil
	}

	err = m.db.Set(block.Header.BlockHeightBytes(), blockBytes)
	if err != nil {
		return alreadyExists, err
	}
	err = m.db.Set(block.Header.BlockHash, block.Header.BlockHeightBytes())
	if err != nil {
		return alreadyExists, err
	}
	m.statsManager.UpdateHighestBlockHeight(block.Header.BlockHeight)

	return alreadyExists, nil
}

func (m *Manager) GetHighestBlockHash() ([]byte, error) {

	blockHeight, err := m.statsManager.GetHighestBlockHeight()
	if err != nil {
		return nil, err
	}
	block, err := m.getBlockByBlockHeight(blockHeight)
	if err != nil {
		return nil, err
	}

	return block.Header.BlockHash, nil
}

func (m *Manager) getBlockByBlockHeight(blockHeight uint64) (*external.Block, error) {
	blockHeightStr := strconv.FormatUint(blockHeight, 10)
	key := []byte(blockHeightStr)
	blockByte, err := m.db.Get(key)
	if err != nil {
		return nil, err
	}
	block := &external.Block{}
	err = block.Deserialize(blockByte)
	if err != nil {
		return nil, err
	}

	return block, err
}

func (m *Manager) GetBlockByHash(key []byte) (*external.Block, error) {
	status, err := m.db.Has(key)
	if err != nil {
		return nil, err
	}
	if !status {
		return nil, fmt.Errorf("block hash not found exists")
	}
	blockHeightByte, err := m.db.Get(key)
	if err != nil {
		return nil, err
	}
	blockHeightStr := fmt.Sprintf("%s", blockHeightByte)
	blockHeight, err := strconv.ParseUint(blockHeightStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return m.getBlockByBlockHeight(blockHeight)
}

func (m *Manager) GetBlockById(blockHeightStr string) (*external.Block, error) {

	blockHeight, err := strconv.ParseUint(blockHeightStr, 10, 64)
	if err != nil {
		return nil, err
	}
	block, err := m.getBlockByBlockHeight(blockHeight)
	if err != nil {
		return nil, err
	}

	return block, nil
}
