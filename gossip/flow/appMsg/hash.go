package appMsg

import "encoding/json"

type InvBlockHash struct {
	Hash []byte `json:"hash"`
}

func NewInvBlockHash() *InvBlockHash {
	return &InvBlockHash{}
}

func (self *InvBlockHash) Decode(data []byte) error {
	return json.Unmarshal(data, self)
}

func (self *InvBlockHash) Encode() ([]byte, error) {
	return json.Marshal(self)
}
