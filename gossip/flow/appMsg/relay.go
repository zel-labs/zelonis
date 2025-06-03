/*
 *
 * Copyright (C) 2025 Zelonis Contributors
 *
 * This file is part of Zelonis.
 *
 * Zelonis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Zelonis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Zelonis. If not, see <https://www.gnu.org/licenses/>.
 *
 */

package appMsg

import (
	"encoding/json"
	"fmt"
	"log"
	"zelonis/external"
)

type RequestBlockRelayInfo struct {
	BlockStart uint64 `json:"block_start"`
	BlockEnd   uint64 `json:"block_end"`
}

func NewRequestBlockRelayInfo() *RequestBlockRelayInfo {
	return &RequestBlockRelayInfo{}
}

func (self *RequestBlockRelayInfo) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}

func (self *RequestBlockRelayInfo) Encode() ([]byte, error) {
	return json.Marshal(self)
}
func (self *RequestBlockRelayInfo) Process(f *flowControl) {
	currentBlock := self.BlockStart
	relayBlocks := NewResponseBlockRelayInfo()
	for currentBlock <= self.BlockEnd {
		heightStr := fmt.Sprintf("%v", currentBlock)
		block, err := f.domain.BlockManager().GetBlockById(heightStr)
		if err != nil {
			return
		}

		relayBlocks.RelayBlocks = append(relayBlocks.RelayBlocks, block)
		currentBlock++
	}
	log.Printf("Sending block from %v to %v", self.BlockStart, self.BlockEnd)
	f.encodeAndSend(relayBlocks, ResponseBlockRelay)
}

type ResponseBlockRelayInfo struct {
	RelayBlocks []*external.Block `json:"relay_blocks"`
}

func NewResponseBlockRelayInfo() *ResponseBlockRelayInfo {
	return &ResponseBlockRelayInfo{
		RelayBlocks: make([]*external.Block, 0),
	}
}
func (self *ResponseBlockRelayInfo) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}
func (self *ResponseBlockRelayInfo) Encode() ([]byte, error) {
	return json.Marshal(self)
}
