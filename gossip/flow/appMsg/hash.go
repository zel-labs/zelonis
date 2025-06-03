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

package appMsg

import (
	"encoding/json"
	"errors"
	"time"
	"zelonis/external"
)

type InvBlockHash struct {
	Hash []byte `json:"hash"`
}

func NewInvBlockHash() *InvBlockHash {
	return &InvBlockHash{}
}

func (self *InvBlockHash) Decode(data []byte) error {
	self.Hash = data
	return nil
}

func (self *InvBlockHash) Encode() ([]byte, error) {
	return json.Marshal(self)
}

func (self *InvBlockHash) Process(f *flowControl) error {

	if f.IsIDBRunning {
		return nil
	}

	//Locate if block exists
	blockHash := self.Hash
	block, err := f.domain.GetBlockByHash(blockHash)
	if err != nil && !errors.Is(err, external.ErrBlockNotFound) {
		panic(err)
	}

	if block != nil {

		updateStatus(block, f)
		f.Synced = true
		//fmt.Println("Already synced")
		return nil
	}
	f.Synced = false

	requestBlockInfo := NewRequestBlockInfo()
	requestBlockInfo.Hash = blockHash
	payload, _ := requestBlockInfo.Encode()
	//Request block
	appFlow := &Flow{
		Header:  RequestBlock,
		Payload: payload,
	}
	msg, err := appFlow.Encode()
	if err != nil {
		return err
	}
	f.sendMsg(msg)
	return nil
}

func updateStatus(block *external.Block, f *flowControl) {
	if !f.Synced {
		f.SyncedTime = time.Now()
	}

	f.NodeStatus.IsConnected = true
	f.NodeStatus.LastUpdated = time.Now()
	f.NodeStatus.Synced = true
	f.NodeStatus.LastBlockHash = block.Header.BlockHash
	f.NodeStatus.LastBlockTime = time.UnixMilli(block.Header.BlockTime)
	f.NodeStatus.LastHeight = block.Header.BlockHeight

}
