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
	"capnproto.org/go/capnp/v3"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
	ping "zelonis/capn"
	"zelonis/external"
	"zelonis/gossip/flow/appMsg"
	"zelonis/validator/domain"
)

var flowList = flows{}

type flows []*flowv1
type flowv1 struct {
	isIncoming bool
	isOutgoing bool
	decoder    *capnp.Decoder
	encoder    *capnp.Encoder
	conn       network.Conn
	domain     *domain.Domain
	validator  bool
	stake      float64
	isSyncing  bool
	*external.NodeStatus
}

func CreateFollow(encoder *capnp.Encoder, decoder *capnp.Decoder, conn network.Conn, domain *domain.Domain, validator bool, stake float64, nodeStatus *external.NodeStatus) *flowv1 {
	flow := &flowv1{
		isIncoming: false,
		isOutgoing: false,
		decoder:    decoder,
		encoder:    encoder,
		conn:       conn,
		domain:     domain,
		validator:  validator,
		stake:      stake,
		NodeStatus: nodeStatus,
	}
	flowList = append(flowList, flow)
	return flow
}

func (f *flowv1) Start(dir int) {
	//create a ping background
	//Send Inv block

	f.sendInvBlockHash(dir)

}
func (f *flowv1) BroadCastMsg(msg *appMsg.Flow) {

	for _, flow := range flowList {
		flow.send(msg)
	}
}

func (f *flowv1) turnOnReciver() error {
	flowContoller := appMsg.NewFlowControl(f.conn, f.encoder, f.decoder, f.domain, f.validator, f.stake, f.NodeStatus)

	for {

		appFlow, err := f.receive()
		if err != nil {
			return err
		}

		msg, status := flowContoller.FilterPayload(appFlow)
		if status {
			f.BroadCastMsg(msg)
		}

	}

}
func (f *flowv1) receive() (*appMsg.Flow, error) {
	defer f.isIncomingDone()
	f.isIncoming = true
	flowMsg := appMsg.NewFlow()
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

func (f *flowv1) send(flowMsg *appMsg.Flow) error {
	defer f.isOutgoingDone()
	f.isOutgoing = true
	userMsg, err := flowMsg.Encode()
	if err != nil {
		return err
	}
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return err
	}
	nping, err := ping.NewRootBlockInfo(seg)
	if err != nil {
		return err
	}
	err = nping.SetMessage_(userMsg)
	if err != nil {
		return err
	}
	if err := f.encoder.Encode(msg); err != nil {
		log.Println(err)

	}

	return nil
}

func (f *flowv1) isOutgoingDone() {
	f.isOutgoing = false
}

func (f *flowv1) isIncomingDone() {
	f.isOutgoing = false
}
