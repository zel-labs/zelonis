package flow

import (
	"capnproto.org/go/capnp/v3"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
	ping "zelonis/capn"
	"zelonis/external"
	"zelonis/gossip/appMsg"
	"zelonis/validator/domain"
)

type flowv1 struct {
	isIncoming bool
	isOutgoing bool
	decoder    *capnp.Decoder
	encoder    *capnp.Encoder
	conn       network.Conn
	domain     *domain.Domain
	validator  bool
	stake      float64
	*external.NodeStatus
}

func CreateFollow(encoder *capnp.Encoder, decoder *capnp.Decoder, conn network.Conn, domain *domain.Domain, validator bool, stake float64, nodeStatus *external.NodeStatus) *flowv1 {
	return &flowv1{
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
}

func (f *flowv1) Start(dir int) {
	//create a ping background
	//Send Inv block

	f.sendInvBlockHash(dir)

}

func (f *flowv1) receive() (*appMsg.Flow, error) {
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
