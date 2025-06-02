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
package external

import "encoding/json"

type Transaction struct {
	Inpoint   *Inpoint    `json:"inpoint"`
	Outpoints []*Outpoint `json:"outpoints"`
	TxHash    []byte      `json:"transaction_hash"`
	Signature []byte      `json:"signature"`
	Fee       []byte      `json:"fee"`
	Status    int8        `json:"status"`
	TxType    int8        `json:"tx_type"`
}

type Inpoint struct {
	PubKey        []byte `json:"pub_key"`
	Value         []byte `json:"value"`
	PrevBlockHash []byte `json:"prev_block_hash"`
}
type Outpoint struct {
	PubKey []byte `json:"pub_key"`
	Value  []byte `json:"value"`
	TxType int8   `json:"tx_type"`
}

func (tx *Transaction) TxSerialize() []byte {
	txBytes, _ := json.Marshal(tx)
	return txBytes
}

func (tx *Transaction) DbTxToDomainTX(txByte []byte) {
	json.Unmarshal(txByte, tx)

}
