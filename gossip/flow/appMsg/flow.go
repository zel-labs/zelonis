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
	"capnproto.org/go/capnp/v3"
	"encoding/json"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
	ping "zelonis/capn"
	"zelonis/external"
	"zelonis/validator/domain"
)

type Flow struct {
	Header  int    `json:"header"`
	Payload []byte `json:"payload"`
}
type flowControl struct {
	conn      network.Conn
	encoder   *capnp.Encoder
	decoder   *capnp.Decoder
	domain    *domain.Domain
	validator bool
	stake     float64
	*external.NodeStatus
}

func NewFlowControl(conn network.Conn, encoder *capnp.Encoder, decorder *capnp.Decoder, domain *domain.Domain, validator bool, stake float64, nodeStatus *external.NodeStatus) *flowControl {
	return &flowControl{
		conn:       conn,
		encoder:    encoder,
		decoder:    decorder,
		domain:     domain,
		validator:  validator,
		stake:      stake,
		NodeStatus: nodeStatus,
	}
}

func NewFlow() *Flow {
	return &Flow{}
}

func (f *Flow) Encode() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Flow) Decode(b []byte) error {
	return json.Unmarshal(b, f)
}

func (f *flowControl) FilterPayload(flow *Flow) bool {

	switch flow.Header {
	case SendInvBlockHash:
		payload := NewInvBlockHash()
		payload.Decode(flow.Payload)
		return false
	case SendProposeBlock:
		payload := NewProposeBlock()
		payload.Decode(flow.Payload)
		payload.Process(f)
	case SendInviTransaction:
		payload := NewInviTransaction()
		payload.Decode(flow.Payload)
		payload.Process(f)
	}
	return true
}

func (f *flowControl) encodeAndSend(msg interface{}, header int) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	flowMsg := &Flow{
		Header:  header,
		Payload: msgByte,
	}
	msgByte, err = flowMsg.Encode()
	if err != nil {
		return err
	}
	f.sendMsg(msgByte)
	return nil
}

func (f *flowControl) sendMsg(msgByte []byte) error {
	msg, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))

	nping, _ := ping.NewRootBlockInfo(seg)

	nping.SetMessage_(msgByte)

	if err := f.encoder.Encode(msg); err != nil {
		log.Println(err)

	}

	return nil
}

func (f *flowControl) getMsg() ([]byte, error) {

	msg, err := f.decoder.Decode()

	if err != nil {

		return nil, err
	}

	decryptedMsg, err := ping.ReadRootBlockInfo(msg)
	if err != nil {

	}
	appMsg, err := decryptedMsg.Message_()
	if err != nil {

		return nil, err
	}

	return appMsg, nil
}
