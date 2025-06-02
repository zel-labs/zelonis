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
	"zelonis/external"
)

type JsonTx struct {
	Inpoint     *JsonInpoint    `json:"inpoint"`
	Outpoint    []*JsonOutpoint `json:"outpoint"`
	TxHash      string          `json:"txhash"`
	SigHash     string          `json:"sighash"`
	Fee         string          `json:"fee"`
	Status      int8            `json:"status"`
	TxType      int8            `json:"txtype"`
	BlockHash   string          `json:"blockhash"`
	BlockHeight string          `json:"blockheight"`
}

type JsonInpoint struct {
	Sender        string `json:"sender"`
	Amount        string `json:"amount"`
	ConcensusHash string `json:"concensushash"`
}

type JsonOutpoint struct {
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (s *RpcServer) getTxByHash(c *fiber.Ctx) error {
	hash := c.Params("hash")

	tx, err := s.domain.TxManager().GetTransactionByHash(hash)
	if err != nil {
		return err
	}
	jsonTx := s.externalTxHashToJsonTx(tx)

	height, err := s.domain.TxManager().GetTransactionBlockHeight(hash)
	if err != nil {
		return err
	}
	block, err := s.domain.BlockManager().GetBlockById(height)
	if err != nil {
		return err
	}
	jsonTx.BlockHeight = height
	jsonTx.BlockHash = fmt.Sprintf("%x", block.Header.BlockHash)

	c.JSON(jsonTx)

	return nil
}

func (s *RpcServer) externalTxHashToJsonTx(tx *external.Transaction) *JsonTx {
	outpoints := make([]*JsonOutpoint, 0)
	for _, outpoint := range tx.Outpoints {
		outpoints = append(outpoints, &JsonOutpoint{
			Receiver: fmt.Sprintf("%s", outpoint.PubKey),
			Amount:   fmt.Sprintf("%s", outpoint.Value),
		})
	}

	return &JsonTx{
		Inpoint: &JsonInpoint{
			Sender:        fmt.Sprintf("%s", tx.Inpoint.PubKey),
			Amount:        fmt.Sprintf("%s", tx.Inpoint.Value),
			ConcensusHash: fmt.Sprintf("%x", tx.Inpoint.PrevBlockHash),
		},
		Outpoint: outpoints,
		TxHash:   fmt.Sprintf("%x", tx.TxHash),

		SigHash:     fmt.Sprintf("%x", tx.Signature),
		Fee:         fmt.Sprintf("%s", tx.Fee),
		Status:      tx.Status,
		TxType:      tx.TxType,
		BlockHash:   "",
		BlockHeight: "",
	}
}
