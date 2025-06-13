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
package stats

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
	"zelonis/external"
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

func (m *Manager) GetTotalTransactions() (uint64, error) {
	byteInfo, err := m.db.Get([]byte(totalTransactionsKey))
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

func (m *Manager) UpdateTotalTransactions() error {
	txCount, _ := m.GetTotalTransactions()
	info := fmt.Sprintf("%v", txCount+1)
	infoBytes := []byte(info)
	err := m.db.Set([]byte(totalTransactionsKey), infoBytes)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetUtxoBlockHeight() (uint64, error) {
	byteInfo, err := m.db.Get([]byte(utxoBlockHeight))
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

func (m *Manager) UpdateUtxoBlockHeight(height uint64) error {

	info := fmt.Sprintf("%v", height)
	infoBytes := []byte(info)
	err := m.db.Set([]byte(utxoBlockHeight), infoBytes)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetTotalStake() ([]byte, error) {
	byteInfo, err := m.db.Get([]byte(totalStaked))
	if err != nil {
		return []byte("0"), err
	}

	return byteInfo, nil
}

func (m *Manager) UpdateTotalStake(val []byte) error {
	staked, _ := m.GetTotalStake()
	stakedT, _ := big.NewFloat(0).SetString(string(staked))
	valFloat, _ := big.NewFloat(0).SetString(string(val))
	stakedT.Add(valFloat, stakedT)

	infoBytes := []byte(stakedT.String())
	err := m.db.Set([]byte(totalStaked), infoBytes)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetTotalSupply() ([]byte, error) {
	byteInfo, err := m.db.Get([]byte(totalSupply))
	if err != nil {
		return []byte("0"), err
	}

	return byteInfo, nil
}

func (m *Manager) UpdateTotalSupply(val []byte) error {
	supply, _ := m.GetTotalSupply()
	supplyT, _ := big.NewFloat(0).SetString(string(supply))
	valFloat, _ := big.NewFloat(0).SetString(string(val))
	supplyT.Add(valFloat, supplyT)

	infoBytes := []byte(supplyT.String())
	err := m.db.Set([]byte(totalSupply), infoBytes)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetCirculating() ([]byte, error) {
	supply, _ := m.GetTotalSupply()
	staked, _ := m.GetTotalStake()
	supplyT, _ := big.NewFloat(0).SetString(string(supply))
	stakedT, _ := big.NewFloat(0).SetString(string(staked))
	supplyT.Sub(supplyT, stakedT)
	infoBytes := []byte(supplyT.String())
	return infoBytes, nil
}

func (m *Manager) GetEpoch() *external.Epoch {
	ch, _ := m.GetHighestBlockHeight()

	epoch := &external.Epoch{}
	if ch < 259200 {
		epoch.EpochNumber = 1
		epoch.EpochStart = 0
		epoch.EpochEnd = 259200
	} else {
		epochSize := uint64(math.Ceil(float64(ch-259200) / (259200 * 3)))
		epoch.EpochNumber = epochSize + 1
		epoch.EpochStart = 259200*3*epochSize - (259200 * 2)
		epoch.EpochEnd = 259200*3 + epoch.EpochStart
		epoch.EpochStart = epoch.EpochStart + 1

	}
	epoch.TimeRemaining = int64(epoch.EpochEnd-ch)*330 + time.Now().UnixMilli()
	return epoch
}
