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
