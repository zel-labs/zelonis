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
	"log"
	ping "zelonis/capn"
	"zelonis/external"
)

func (self *ResponseBlockInfo) StartIDB(block *external.Block, f *flowControl) (bool, error) {

	targetBlockHeight := self.Block.Header.BlockHeight
	currentBlockHeight := block.Header.BlockHeight
	var maxBlock = uint64(199)
	for currentBlockHeight < targetBlockHeight {
		diff := targetBlockHeight - currentBlockHeight
		if diff < maxBlock {
			maxBlock = diff
		}
		rangeStart := currentBlockHeight
		rangeEnd := currentBlockHeight + maxBlock
		//GetBlocks In IDB
		requestRange := &RequestBlockRelayInfo{
			BlockStart: rangeStart,
			BlockEnd:   rangeEnd,
		}

		f.encodeAndSend(requestRange, RequestBlockRelay)

		currentBlockHeight = rangeEnd
		currentP := float64(currentBlockHeight) * 100 / float64(targetBlockHeight)

		log.Printf("Syncing Complete %v %% (Total Blocks till now %v from %v)", currentP, currentBlockHeight, targetBlockHeight)
		err := f.reciveIDBBlocks()
		if err != nil {
			return false, err
		}
	}
	log.Printf("Synced complete to block %v \n", targetBlockHeight)

	f.IsIDBRunning = false
	return true, nil
}

func (f *flowControl) reciveIDBBlocks() error {
	for {
		flow, err := f.receive()
		if err != nil {
			return err
		}
		switch flow.Header {
		case ResponseBlockRelay:
			payload := NewResponseBlockRelayInfo()
			payload.Decode(flow.Payload)
			payload.Process(f)
			return nil
		default:
			f.FilterPayload(flow)
		}
	}
}

func (self *ResponseBlockRelayInfo) Process(f *flowControl) {
	for _, block := range self.RelayBlocks {
		f.domain.VerifyInsertBlockAndTransaction(block)
	}
}

func (f *flowControl) receive() (*Flow, error) {

	flowMsg := NewFlow()
	msg, err := f.decoder.Decode()

	if err != nil {

		return nil, err
	}

	decryptedMsg, err := ping.ReadRootBlockInfo(msg)
	if err != nil {
		return nil, err
	}
	decrypted, err := decryptedMsg.Message_()
	if err != nil {

		return nil, err
	}

	err = flowMsg.Decode(decrypted)
	if err != nil {
		return nil, err
	}
	return flowMsg, nil
}
