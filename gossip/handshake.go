package gossip

import (
	"capnproto.org/go/capnp/v3"
	"encoding/json"
	"github.com/libp2p/go-libp2p/core/network"
	"log"
	"time"
	ping "zelonis/capn"
	"zelonis/external"
	flowv1 "zelonis/gossip/flow"
	"zelonis/gossip/flow/appMsg"
	"zelonis/validator/domain"
)

const (
	defaultTimeout = 5 * time.Second
)

type zelPeer struct {
	conn      network.Conn
	encoder   *capnp.Encoder
	decoder   *capnp.Decoder
	handshake bool
	domain    *domain.Domain
	validator bool
	stake     float64
	*external.NodeStatus
}

func (z *zelPeer) Close() {
	z.conn.Close()
}
func (z *zelPeer) ErrorHandler(err error) {
	defer z.Close()
	log.Println(err)

}
func (g *gossipLister) handShake(s network.Stream) {
	z := &zelPeer{
		conn:       s.Conn(),
		encoder:    capnp.NewEncoder(s),
		decoder:    capnp.NewDecoder(s),
		handshake:  false,
		domain:     g.domain,
		validator:  g.validator,
		stake:      g.stake,
		NodeStatus: g.NodeStatus,
	}
	//start handshake pattern

	if err := z.requestHandShake(); err != nil {
		z.ErrorHandler(err)
	}

	flow := flowv1.CreateFollow(z.encoder, z.decoder, z.conn, z.domain, z.validator, z.stake, g.NodeStatus)
	flow.Start(0)
	//if valid add p2phandler relay

}

func (z *zelPeer) requestHandShake() error {
	//Build handshakeMsg

	msg, err := z.getMsg()
	if err != nil {
		return err
	}

	msgHandshake := appMsg.NewHandshake()
	err = msgHandshake.Decode(msg)
	if err != nil {
		return err
	}

	handshake := appMsg.NewHandshakeWithInfo()
	msg, err = handshake.Encode()
	if err != nil {
		return err
	}

	err = z.sendMsg(msg)
	if err != nil {
		return err
	}

	if err = msgHandshake.Verify(handshake); err != nil {

		return err
	}

	return nil
}

func (z *zelPeer) encodeAndSend(msg interface{}, msgType int) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	flowMsg := &appMsg.Flow{
		Header:  msgType,
		Payload: msgByte,
	}
	msgByte, err = flowMsg.Encode()
	if err != nil {
		return err
	}
	z.sendMsg(msgByte)
	return nil
}

func (z *zelPeer) sendMsg(msgByte []byte) error {
	msg, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))

	nping, _ := ping.NewRootBlockInfo(seg)

	nping.SetMessage_(msgByte)

	if err := z.encoder.Encode(msg); err != nil {
		log.Println(err)

	}

	return nil
}

func (z *zelPeer) getMsg() ([]byte, error) {

	msg, err := z.decoder.Decode()

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
