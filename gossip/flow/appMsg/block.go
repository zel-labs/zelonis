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
	"log"
	"zelonis/external"
)

type BlockInfo struct {
	*external.Block
}

func NewBlockInfo() *BlockInfo {
	return &BlockInfo{}
}

func (self *BlockInfo) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}

func (self *BlockInfo) Encode() ([]byte, error) {
	return json.Marshal(self)
}

type RequestBlockInfo struct {
	Hash []byte `json:"hash"`
}

func NewRequestBlockInfo() *RequestBlockInfo {
	return &RequestBlockInfo{}
}

func (self *RequestBlockInfo) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}

func (self *RequestBlockInfo) Encode() ([]byte, error) {
	return json.Marshal(self)
}

func (self *RequestBlockInfo) Process(f *flowControl) {
	block, err := f.domain.GetBlockByHash(self.Hash)
	if err != nil {
		return
	}
	resBlockInfo := NewResponseBlockInfo()
	resBlockInfo.Block = *block

	f.encodeAndSend(resBlockInfo, ResponseBlock)

	//ResponseBlockInfo
}

type ResponseBlockInfo struct {
	Block external.Block `json:"block"`
}

func NewResponseBlockInfo() *ResponseBlockInfo {
	return &ResponseBlockInfo{}
}

func (self *ResponseBlockInfo) Decode(data []byte) error {

	err := json.Unmarshal(data, self)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(data, self)
}

func (self *ResponseBlockInfo) Encode() ([]byte, error) {
	return json.Marshal(self)
}

func (self *ResponseBlockInfo) Process(f *flowControl) {
	if f.IsIDBRunning {
		return
	}
	heighestHash, err := f.domain.GetHighestBlockHash()
	if err != nil {
		return
	}
	cBlock, err := f.domain.GetBlockByHash(heighestHash)
	if err != nil {
		return
	}

	//Compare difference in recived block and current block
	if self.Block.Header.BlockHeight-cBlock.Header.BlockHeight > 1 && self.Block.Header.BlockHeight > cBlock.Header.BlockHeight {
		//Turn on IDB Sync
		log.Println("block height is ", self.Block.Header.BlockHeight, cBlock.Header.BlockHeight)

		f.IsIDBRunning = true
		status, err := self.StartIDB(cBlock, f)
		if err != nil {

			f.conn.Close()
			f.IsIDBRunning = false
			log.Println(err)
			return
		}
		if status {
			log.Println("IDB Completed successfully")

		}
	} else {
		f.domain.VerifyInsertBlockAndTransaction(&self.Block)
		f.Synced = true
	}

}
