package flow

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"zelonis/gossip/appMsg"
)

func (f *flowv1) sendInvBlockHash(dir int) error {
	//Get highest block Hash
	for {

		blockHash, err := f.domain.GetHighestBlockHash()
		log.Printf("%x", blockHash)
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
			fmt.Println("Already synced")
			continue
		}

	}
}
