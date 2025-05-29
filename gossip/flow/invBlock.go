package flow

import (
	"errors"
	"fmt"
	"reflect"
	"time"
	"zelonis/external"
	"zelonis/gossip/appMsg"
)

func (f *flowv1) sendInvBlockHash(dir int) error {
	//Get highest block Hash
	for {

		blockHash, err := f.domain.GetHighestBlockHash()

		if err != nil {

			return err
		}
		appFlow := &appMsg.Flow{
			Header:  appMsg.SendInvBlockHash,
			Payload: blockHash,
		}
		f.send(appFlow)
		appFlow, err = f.receive()
		if err != nil {
			return err
		}
		if appFlow.Header != appMsg.SendInvBlockHash {
			return errors.New("header mismatch")
		}
		//Locate if block exists
		block, err := f.domain.GetBlockByHash(blockHash)
		if err != nil {
			return err
		}

		if reflect.DeepEqual(block.Header.BlockHash, blockHash) {

			f.updateStatus(block)

			fmt.Println("Already synced")
			continue
		}
		f.Synced = false

	}
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
