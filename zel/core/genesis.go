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

package core

type Genesis struct {
	BlockHeight  []byte         `json:"block_height"`
	BlockTime    []byte         `json:"block_time"`
	BlockHash    []byte         `json:"hash"`
	ParentSlot   []byte         `json:"parent_slot"`
	ParentHash   []byte         `json:"parent_hash"`
	Transactions []*Transaction `json:"transactions"`
	Version      []byte         `json:"version"`
}

type Transaction struct {
	Inpoint   *Inpoint    `json:"inpoint"`
	Outpoints []*Outpoint `json:"outpoints"`
	Balance   []byte      `json:"balance"`
	TxHash    []byte      `json:"transaction_hash"`
	Signature []byte      `json:"signature"`
}

type Inpoint struct {
	PubKey        []byte `json:"pub_key"`
	Value         []byte `json:"value"`
	PrevBlockHash []byte `json:"prev_block_hash"`
}

type Outpoint struct {
	PubKey []byte `json:"pub_key"`
	Value  []byte `json:"value"`
}
