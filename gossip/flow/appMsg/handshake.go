package appMsg

import (
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	handShakeVersion = "0.0.01"
)

type handshake struct {
	Version string `json:"version"`
	Action  int    `json:"action"`
}

func NewHandshakeWithInfo() *handshake {
	return &handshake{
		Version: handShakeVersion,
		Action:  HandshakeMsg,
	}
}
func NewHandshake() *handshake {
	return &handshake{}
}

func (h *handshake) Decode(data []byte) error {
	return json.Unmarshal(data, h)
}
func (h *handshake) Encode() ([]byte, error) {
	return json.Marshal(h)
}

func (h *handshake) Verify(msg *handshake) error {
	if h.Version != msg.Version {

		return errors.New("handshake version mismatch")
	} else if h.Action != msg.Action || h.Action != HandshakeMsg {
		return errors.New("handshake action mismatch")
	}
	return nil
}
