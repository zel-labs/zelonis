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
package flow

import (
	"time"
	"zelonis/external"
	"zelonis/gossip/flow/appMsg"
)

func (f *flowv1) sendInvBlockHash(dir int) error {

	blockHash, err := f.domain.GetHighestBlockHash()

	if err != nil {

		return err
	}
	appFlow := &appMsg.Flow{
		Header:  appMsg.SendInvBlockHash,
		Payload: blockHash,
	}
	f.send(appFlow)
	f.turnOnReciver()
	return nil
}

func (f *flowv1) updateStatus(block *external.Block) {
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
