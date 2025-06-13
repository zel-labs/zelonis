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

package validator

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type jsonAccountInfo struct {
	Balance             string `json:"balance"`
	Stake               string `json:"stake"`
	Reward              string `json:"reward"`
	ActivatingStake     string `json:"activating_stake"`
	DeactivatingStake   string `json:"deactivating_stake"`
	PendingActivation   string `json:"pending_activation"`
	PendingDeactivation string `json:"pending_deactivation"`
	WarmupStake         string `json:"warmup_stake"`
	CoolingDownStake    string `json:"cooling_down_stake"`
}

func (s *RpcServer) getAccountBalance(c *fiber.Ctx) error {
	account := c.Params("account")
	accountInfo, status := s.domain.GetAccountBalance([]byte(account))
	if !status {
		return fmt.Errorf("Account not found")
	}

	jsonAccount := &jsonAccountInfo{
		Balance:             s.valToJsonVal(accountInfo.Balance),
		Stake:               s.valToJsonVal(accountInfo.Stake),
		Reward:              s.valToJsonVal(accountInfo.Reward),
		ActivatingStake:     s.valToJsonVal(accountInfo.ActivatingStake),
		DeactivatingStake:   s.valToJsonVal(accountInfo.DeactivatingStake),
		PendingActivation:   s.valToJsonVal(accountInfo.PendingActivation),
		PendingDeactivation: s.valToJsonVal(accountInfo.PendingDeactivation),
		WarmupStake:         s.valToJsonVal(accountInfo.WarmupStake),
		CoolingDownStake:    s.valToJsonVal(accountInfo.CoolingDownStake),
	}
	c.JSON(jsonAccount)
	return nil
}

func (s *RpcServer) valToJsonVal(val []byte) string {
	valStr := fmt.Sprintf("%s", val)
	//bigVal, _ := big.NewFloat(0).SetString(valStr)
	valFloat, _ := strconv.ParseFloat(valStr, 64)
	return fmt.Sprintf("%.9f", valFloat)
}

func (s *RpcServer) getAccountTx(c *fiber.Ctx) error {
	account := c.Params("account")

	list, err := s.domain.AccountManager().GetTxHistoryList([]byte(account), 0, 10)
	if err != nil {
		return err
	}
	c.JSON(list)
	return nil
}

type jsonTxList struct {
	TxList []*JsonTx `json:"tx_list"`
}

func (s *RpcServer) getAccountTxList(c *fiber.Ctx) error {
	account := c.Params("account")

	list, err := s.domain.AccountManager().GetTxHistoryList([]byte(account), 0, 10)
	if err != nil {
		return err
	}
	txs := []*JsonTx{}
	for _, txHash := range list {
		tx, err := s.domain.TxManager().GetTransactionByHash(txHash)
		if err != nil {
			return err
		}
		jsonTx := s.externalTxHashToJsonTx(tx)
		height, err := s.domain.TxManager().GetTransactionBlockHeight(txHash)
		if err != nil {
			return err
		}
		block, err := s.domain.BlockManager().GetBlockById(height)
		if err != nil {
			return err
		}

		jsonTx.Timstamp = block.Header.BlockTime
		txs = append(txs, jsonTx)
	}
	JsonTxList := &jsonTxList{TxList: txs}
	c.JSON(JsonTxList)
	return nil
}

func (s *RpcServer) getAccountTxWithLimit(c *fiber.Ctx) error {
	account := c.Params("account")
	from, _ := strconv.Atoi(c.Params("from"))
	limit, _ := strconv.Atoi(c.Params("limit"))

	list, err := s.domain.AccountManager().GetTxHistoryList([]byte(account), from, limit)
	if err != nil {
		return err
	}
	c.JSON(list)
	return nil
}
